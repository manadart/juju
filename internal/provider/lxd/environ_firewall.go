// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package lxd

import (
	"context"
	"net/http"

	"github.com/canonical/lxd/shared/api"
	"github.com/juju/errors"

	"github.com/juju/juju/core/network/firewall"
	"github.com/juju/juju/environs"
	"github.com/juju/juju/environs/config"
	"github.com/juju/juju/environs/instances"
	"github.com/juju/juju/environs/models"
)

var (
	_ environs.Firewaller             = (*environ)(nil)
	_ models.ModelFirewaller          = (*environ)(nil)
	_ environs.FirewallFeatureQuerier = (*environ)(nil)
	_ instances.InstanceFirewaller    = (*environInstance)(nil)
)

// SupportsRulesWithIPV6CIDRs implements environs.FirewallFeatureQuerier.
func (env *environ) SupportsRulesWithIPV6CIDRs(ctx context.Context) (bool, error) {
	return true, ctx.Err()
}

// OpenPorts implements environs.Firewaller for global firewall mode.
func (env *environ) OpenPorts(ctx context.Context, rules firewall.IngressRules) error {
	return env.changeSharedFirewall(ctx, true, rules, true)
}

// ClosePorts implements environs.Firewaller for global firewall mode.
func (env *environ) ClosePorts(ctx context.Context, rules firewall.IngressRules) error {
	return env.changeSharedFirewall(ctx, true, rules, false)
}

// IngressRules implements environs.Firewaller for global firewall mode.
func (env *environ) IngressRules(ctx context.Context) (firewall.IngressRules, error) {
	return env.sharedFirewallRules(ctx, true)
}

// OpenModelPorts implements models.ModelFirewaller.
func (env *environ) OpenModelPorts(ctx context.Context, rules firewall.IngressRules) error {
	return env.changeSharedFirewall(ctx, false, rules, true)
}

// CloseModelPorts implements models.ModelFirewaller.
func (env *environ) CloseModelPorts(ctx context.Context, rules firewall.IngressRules) error {
	return env.changeSharedFirewall(ctx, false, rules, false)
}

// ModelIngressRules implements models.ModelFirewaller.
func (env *environ) ModelIngressRules(ctx context.Context) (firewall.IngressRules, error) {
	return env.sharedFirewallRules(ctx, false)
}

func (env *environ) changeSharedFirewall(ctx context.Context, global bool, rules firewall.IngressRules, open bool) error {
	if err := rules.Validate(); err != nil {
		return errors.Trace(err)
	}
	client, name, owner, err := env.sharedFirewall(ctx, global)
	if err != nil || client == nil {
		return err
	}
	err = updateNetworkACL(ctx, client, name, owner, func(acl *api.NetworkACLPut) error {
		return changeACLIngress(acl, rules, open)
	})
	return env.HandleCredentialError(ctx, err)
}

func (env *environ) sharedFirewallRules(ctx context.Context, global bool) (firewall.IngressRules, error) {
	client, name, owner, err := env.sharedFirewall(ctx, global)
	if err != nil || client == nil {
		return nil, err
	}
	rules, err := readACLIngress(ctx, client, name, owner)
	return rules, env.HandleCredentialError(ctx, err)
}

func (env *environ) sharedFirewall(ctx context.Context, global bool) (networkACLServer, string, map[string]string, error) {
	if err := ctx.Err(); err != nil {
		return nil, "", nil, err
	}
	mode := env.Config().FirewallMode()
	if mode == config.FwNone {
		return nil, "", nil, nil
	}
	if global && mode != config.FwGlobal {
		return nil, "", nil, errors.Errorf("global firewall operation in %q mode", mode)
	}

	// Reconcile attachments as well as rules. This adopts existing OVN
	// instances on upgrade, including machines with no exposed ports, and
	// ensures global mode reaches every machine in the model.
	all, err := env.allInstances()
	if err != nil {
		return nil, "", nil, env.HandleCredentialError(ctx, err)
	}
	var client networkACLServer
	for _, inst := range all {
		attached, err := env.ensureInstanceFirewall(ctx, string(inst.Id()))
		if err != nil {
			return nil, "", nil, env.HandleCredentialError(ctx, err)
		}
		if attached != nil {
			client = attached
		}
	}
	// Bridge-only models deliberately retain the provider's unmanaged
	// firewall behaviour. No ACL is created or attached for their NICs.
	if client == nil {
		return nil, "", nil, nil
	}
	name, kind := modelACLName(env.uuid), "model"
	if global {
		name, kind = globalACLName(env.uuid), "global"
	}
	return client, name, aclMetadata(env.uuid, "", "", kind), nil
}

// OpenPorts implements instances.InstanceFirewaller.
func (i *environInstance) OpenPorts(ctx context.Context, machineID string, rules firewall.IngressRules) error {
	return i.changeFirewall(ctx, machineID, rules, true)
}

// ClosePorts implements instances.InstanceFirewaller.
func (i *environInstance) ClosePorts(ctx context.Context, machineID string, rules firewall.IngressRules) error {
	return i.changeFirewall(ctx, machineID, rules, false)
}

// IngressRules implements instances.InstanceFirewaller.
func (i *environInstance) IngressRules(ctx context.Context, machineID string) (firewall.IngressRules, error) {
	client, name, owner, err := i.instanceFirewall(ctx, machineID)
	if err != nil || client == nil {
		return nil, err
	}
	rules, err := readACLIngress(ctx, client, name, owner)
	return rules, i.env.HandleCredentialError(ctx, err)
}

func (i *environInstance) changeFirewall(ctx context.Context, machineID string, rules firewall.IngressRules, open bool) error {
	if err := rules.Validate(); err != nil {
		return errors.Trace(err)
	}
	client, name, owner, err := i.instanceFirewall(ctx, machineID)
	if err != nil || client == nil {
		return err
	}
	err = updateNetworkACL(ctx, client, name, owner, func(acl *api.NetworkACLPut) error {
		return changeACLIngress(acl, rules, open)
	})
	return i.env.HandleCredentialError(ctx, err)
}

func (i *environInstance) instanceFirewall(ctx context.Context, machineID string) (networkACLServer, string, map[string]string, error) {
	if err := ctx.Err(); err != nil {
		return nil, "", nil, err
	}
	mode := i.env.Config().FirewallMode()
	if mode == config.FwNone {
		return nil, "", nil, nil
	}
	if mode != config.FwInstance {
		return nil, "", nil, errors.Errorf("instance firewall operation in %q mode", mode)
	}
	expectedName, err := i.env.namespace.Hostname(machineID)
	if err != nil {
		return nil, "", nil, errors.Trace(err)
	}
	if expectedName != i.container.Name {
		return nil, "", nil, errors.Errorf("instance %q does not belong to machine %q", i.container.Name, machineID)
	}
	client, err := i.env.ensureInstanceFirewall(ctx, i.container.Name)
	if err != nil {
		return nil, "", nil, i.env.HandleCredentialError(ctx, err)
	}
	return client, machineACLName(i.env.uuid, i.container.Name),
		aclMetadata(i.env.uuid, "", i.container.Name, "machine"), nil
}

func readACLIngress(ctx context.Context, client networkACLServer, name string, owner map[string]string) (firewall.IngressRules, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	acl, _, err := client.GetNetworkACL(name)
	if api.StatusErrorCheck(err, http.StatusNotFound) {
		return nil, errors.NotFoundf("firewall ACL %q", name)
	}
	if err != nil {
		return nil, errors.Annotatef(err, "reading firewall ACL %q", name)
	}
	if err := checkACLOwner(acl, owner); err != nil {
		return nil, err
	}
	return aclIngressRules(acl.Ingress)
}
