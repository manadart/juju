// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package lxd

import (
	"fmt"
	"net/netip"
	"strconv"
	"strings"

	"github.com/canonical/lxd/shared/api"
	"github.com/juju/collections/set"
	"github.com/juju/errors"

	"github.com/juju/juju/core/network"
	"github.com/juju/juju/core/network/firewall"
)

const (
	jujuIngressDescription  = "Juju ingress rule"
	jujuBaselineDescription = "Juju model network access"
)

// changeACLIngress applies a delta to the Juju rules, preserving baseline
// rules and any rules added by the operator. Source CIDRs are removed
// individually so closing one source cannot remove access for another.
func changeACLIngress(acl *api.NetworkACLPut, rules firewall.IngressRules, open bool) error {
	current, err := aclIngressRules(acl.Ingress)
	if err != nil {
		return err
	}
	byPort, err := ingressCIDRs(current)
	if err != nil {
		return err
	}
	delta, err := ingressCIDRs(rules)
	if err != nil {
		return err
	}
	for ports, cidrs := range delta {
		if open {
			if byPort[ports] == nil {
				byPort[ports] = set.NewStrings()
			}
			byPort[ports] = byPort[ports].Union(cidrs)
		} else {
			byPort[ports] = byPort[ports].Difference(cidrs)
		}
	}
	encoded, err := ingressACLRules(ingressRulesFromCIDRs(byPort))
	if err != nil {
		return err
	}
	var preserved []api.NetworkACLRule
	for _, rule := range acl.Ingress {
		if rule.Description != jujuIngressDescription {
			preserved = append(preserved, rule)
		}
	}
	acl.Ingress = append(preserved, encoded...)
	return nil
}

// ingressACLRules emits one rule per CIDR. In particular, ICMP must be
// separated by address family because LXD uses icmp4 and icmp6 protocols.
func ingressACLRules(rules firewall.IngressRules) ([]api.NetworkACLRule, error) {
	byPort, err := ingressCIDRs(rules)
	if err != nil {
		return nil, err
	}
	var result []api.NetworkACLRule
	for _, rule := range ingressRulesFromCIDRs(byPort) {
		for _, cidr := range rule.SourceCIDRs.SortedValues() {
			entry := api.NetworkACLRule{
				Action: "allow", State: "enabled",
				Description: jujuIngressDescription,
				Protocol:    rule.PortRange.Protocol, Source: cidr,
			}
			if entry.Protocol == "icmp" {
				entry.Protocol = "icmp4"
				if strings.Contains(cidr, ":") {
					entry.Protocol = "icmp6"
				}
			} else {
				entry.DestinationPort = strconv.Itoa(rule.PortRange.FromPort)
				if rule.PortRange.FromPort != rule.PortRange.ToPort {
					entry.DestinationPort += "-" + strconv.Itoa(rule.PortRange.ToPort)
				}
			}
			result = append(result, entry)
		}
	}
	return result, nil
}

// aclIngressRules reads only dynamic Juju rules. Baseline rules must not be
// returned to the firewaller, which would otherwise reconcile them away.
func aclIngressRules(entries []api.NetworkACLRule) (firewall.IngressRules, error) {
	var result firewall.IngressRules
	for _, entry := range entries {
		if entry.Description != jujuIngressDescription {
			continue
		}
		if entry.State == "disabled" {
			continue
		}
		if entry.Action != "allow" || entry.Destination != "" ||
			entry.SourcePort != "" || entry.ICMPType != "" || entry.ICMPCode != "" {
			return nil, errors.NotSupportedf("modified Juju ACL rule %+v", entry)
		}
		protocol := entry.Protocol
		var ports []network.PortRange
		switch protocol {
		case "icmp4", "icmp6":
			if entry.DestinationPort != "" {
				return nil, errors.NotValidf("ICMP ACL destination port")
			}
			ports = []network.PortRange{{Protocol: "icmp", FromPort: -1, ToPort: -1}}
		case "tcp", "udp":
			portList := entry.DestinationPort
			if portList == "" {
				portList = "0-65535"
			}
			for _, value := range strings.Split(portList, ",") {
				portRange, err := network.ParsePortRange(strings.TrimSpace(value) + "/" + protocol)
				if err != nil {
					return nil, errors.Annotate(err, "parsing ACL destination ports")
				}
				ports = append(ports, portRange)
			}
		default:
			return nil, errors.NotSupportedf("ACL protocol %q", protocol)
		}
		var sources []string
		if entry.Source == "" {
			if protocol != "icmp6" {
				sources = append(sources, firewall.AllNetworksIPV4CIDR)
			}
			if protocol != "icmp4" {
				sources = append(sources, firewall.AllNetworksIPV6CIDR)
			}
		} else {
			for _, source := range strings.Split(entry.Source, ",") {
				prefix, err := netip.ParsePrefix(strings.TrimSpace(source))
				if err != nil {
					return nil, errors.Annotate(err, "parsing ACL source CIDR")
				}
				if (protocol == "icmp4" && !prefix.Addr().Is4()) ||
					(protocol == "icmp6" && !prefix.Addr().Is6()) {
					return nil, errors.NotValidf("ACL source address family for %q", protocol)
				}
				sources = append(sources, prefix.Masked().String())
			}
		}
		for _, portRange := range ports {
			result = append(result, firewall.NewIngressRule(portRange, sources...))
		}
	}
	byPort, err := ingressCIDRs(result)
	if err != nil {
		return nil, err
	}
	return ingressRulesFromCIDRs(byPort), nil
}

func ingressCIDRs(rules firewall.IngressRules) (map[network.PortRange]set.Strings, error) {
	if err := rules.Validate(); err != nil {
		return nil, errors.Trace(err)
	}
	result := make(map[network.PortRange]set.Strings)
	for _, rule := range rules {
		ports := rule.PortRange
		ports.Protocol = strings.ToLower(ports.Protocol)
		if result[ports] == nil {
			result[ports] = set.NewStrings()
		}
		if len(rule.SourceCIDRs) == 0 {
			result[ports].Add(firewall.AllNetworksIPV4CIDR)
			result[ports].Add(firewall.AllNetworksIPV6CIDR)
		}
		for cidr := range rule.SourceCIDRs {
			prefix, err := netip.ParsePrefix(cidr)
			if err != nil {
				return nil, fmt.Errorf("parsing ingress source CIDR: %w", err)
			}
			result[ports].Add(prefix.Masked().String())
		}
	}
	return result, nil
}

func ingressRulesFromCIDRs(byPort map[network.PortRange]set.Strings) firewall.IngressRules {
	var result firewall.IngressRules
	for ports, cidrs := range byPort {
		if len(cidrs) > 0 {
			result = append(result, firewall.NewIngressRule(ports, cidrs.SortedValues()...))
		}
	}
	result.Sort()
	return result
}
