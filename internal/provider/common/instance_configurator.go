// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package common

import (
	"context"
	"os"
	"strings"

	"github.com/juju/errors"
	"github.com/juju/utils/v4/ssh"

	"github.com/juju/juju/core/network/firewall"
	"github.com/juju/juju/internal/network/iptables"
)

// InstanceConfigurator describes methods for manipulating firewall
// rules directly on a single instance.

type InstanceConfigurator interface {

	// DropAllPorts denies access to all ports.
	DropAllPorts(exceptPorts []int, addr string) error

	// ChangeIngressRules opens and/or closes ports.
	ChangeIngressRules(ipAddress string, insert bool, rules firewall.IngressRules) error

	// FindIngressRules returns all firewall rules.
	FindIngressRules() (firewall.IngressRules, error)
}

type sshInstanceConfigurator struct {
	client  ssh.Client
	host    string
	options *ssh.Options
}

// NewSshInstanceConfigurator creates new sshInstanceConfigurator.
func NewSshInstanceConfigurator(host string) InstanceConfigurator {
	options := ssh.Options{}
	options.SetIdentities("/var/lib/juju/system-identity")

	// Disable host key checking. We're not sending any sensitive data
	// across, and we don't have access to the host's keys from here.
	//
	// TODO(axw) 2017-12-07 #1732665
	// Stop using SSH, instead manage iptables on the machine
	// itself. This will also provide firewalling for MAAS and
	// LXD machines.
	options.SetStrictHostKeyChecking(ssh.StrictHostChecksNo)
	options.SetKnownHostsFile(os.DevNull)

	return &sshInstanceConfigurator{
		client:  ssh.DefaultClient,
		host:    "ubuntu@" + host,
		options: &options,
	}
}

func (c *sshInstanceConfigurator) runCommand(cmd string) (string, error) {
	command := c.client.Command(c.host, []string{"/bin/bash"}, c.options)
	command.Stdin = strings.NewReader(cmd)
	output, err := command.CombinedOutput()
	if err != nil {
		return "", errors.Trace(err)
	}
	return string(output), nil
}

// DropAllPorts implements InstanceConfigurator interface.
func (c *sshInstanceConfigurator) DropAllPorts(exceptPorts []int, addr string) error {
	cmds := []string{
		iptables.DropCommand{DestinationAddress: addr}.Render(),
	}
	for _, port := range exceptPorts {
		cmds = append(cmds, iptables.AcceptInternalCommand{
			Protocol:           "tcp",
			DestinationAddress: addr,
			DestinationPort:    port,
		}.Render())
	}

	output, err := c.runCommand(strings.Join(cmds, "\n"))
	if err != nil {
		return errors.Errorf("failed to drop all ports: %s", output)
	}
	logger.Tracef(context.TODO(), "drop all ports output: %s", output)
	return nil
}

// ChangeIngressRules implements InstanceConfigurator interface.
func (c *sshInstanceConfigurator) ChangeIngressRules(ipAddress string, insert bool, rules firewall.IngressRules) error {
	var cmds []string
	for _, rule := range rules {
		cmds = append(cmds, iptables.IngressRuleCommand{
			Rule:               rule,
			DestinationAddress: ipAddress,
			Delete:             !insert,
		}.Render())
	}

	output, err := c.runCommand(strings.Join(cmds, "\n"))
	if err != nil {
		return errors.Annotatef(err, "configuring ports for address %q: %s", ipAddress, output)
	}
	logger.Tracef(context.TODO(), "change ports output: %s", output)
	return nil
}

// FindIngressRules implements InstanceConfigurator interface.
func (c *sshInstanceConfigurator) FindIngressRules() (firewall.IngressRules, error) {
	output, err := c.runCommand("sudo iptables -L INPUT -n")
	if err != nil {
		return nil, errors.Errorf("failed to list open ports: %s", output)
	}
	logger.Tracef(context.TODO(), "find open ports output: %s", output)
	return iptables.ParseIngressRules(strings.NewReader(output))
}
