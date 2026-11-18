// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package lxd

import (
	"context"
	"crypto/sha256"
	"fmt"
	"maps"
	"net/http"
	"reflect"

	lxdclient "github.com/canonical/lxd/client"
	"github.com/canonical/lxd/shared/api"
	"github.com/juju/errors"

	"github.com/juju/juju/environs/tags"
	"github.com/juju/juju/internal/container/lxd"
)

const (
	jujuACLModelKey      = lxd.UserNamespacePrefix + tags.JujuModel
	jujuACLControllerKey = lxd.UserNamespacePrefix + tags.JujuController
	jujuACLInstanceKey   = lxd.UserNamespacePrefix + "juju-instance"
	jujuACLKindKey       = lxd.UserNamespacePrefix + "juju-firewall-kind"
	aclUpdateAttempts    = 5
)

// networkACLServer is the optional LXD API required only for OVN firewalling.
// Keeping it separate allows servers without ACL support to provide ordinary
// bridge instances. The production server supplies these methods through its
// embedded LXD client, in the same project as its instance operations.
type networkACLServer interface {
	GetNetworkACLs() ([]api.NetworkACL, error)
	GetNetworkACL(string) (*api.NetworkACL, string, error)
	CreateNetworkACL(api.NetworkACLsPost) (lxdclient.Operation, error)
	UpdateNetworkACL(string, api.NetworkACLPut, string) (lxdclient.Operation, error)
	DeleteNetworkACL(string) (lxdclient.Operation, error)
	UpdateInstance(string, api.InstancePut, string) (lxdclient.Operation, error)
}

var _ networkACLServer = (*providerServer)(nil)

func aclServer(srv Server) (networkACLServer, error) {
	for _, extension := range []string{"network_acl", "network_ovn_acl", "network_ovn_acl_defaults"} {
		if !srv.HasExtension(extension) {
			return nil, errors.NotSupportedf("OVN firewalling requires the %q LXD API extension", extension)
		}
	}
	client, ok := srv.(networkACLServer)
	if !ok {
		return nil, errors.NotSupportedf("LXD server ACL operations")
	}
	return client, nil
}

func modelACLName(modelUUID string) string {
	return "juju-" + modelUUID + "-model"
}

func globalACLName(modelUUID string) string {
	return "juju-" + modelUUID + "-global"
}

func machineACLName(modelUUID, instanceName string) string {
	// LXD ACL names are limited to 63 characters and cannot contain dots or
	// slashes. Instance names may contain both, and need not be short.
	digest := sha256.Sum256([]byte(instanceName))
	return fmt.Sprintf("juju-%s-m-%x", modelUUID, digest[:8])
}

// ensureNetworkACL creates an owned ACL if absent. Creation uses no subject
// selectors: LXD validates those against existing ACL names, including self
// references. Baseline selectors are installed afterwards, before attachment.
func ensureNetworkACL(ctx context.Context, client networkACLServer, name string, metadata map[string]string, ingress []api.NetworkACLRule) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	acl, _, err := client.GetNetworkACL(name)
	if err == nil {
		return checkACLOwner(acl, metadata)
	}
	if !api.StatusErrorCheck(err, http.StatusNotFound) {
		return errors.Annotatef(err, "reading ACL %q", name)
	}
	op, err := client.CreateNetworkACL(api.NetworkACLsPost{
		Name: name,
		NetworkACLPut: api.NetworkACLPut{
			Description: "Juju OVN firewall",
			Config:      maps.Clone(metadata), Ingress: ingress,
		},
	})
	if err == nil {
		err = op.WaitContext(ctx)
	}
	if err != nil {
		// Another provisioner can create the shared model ACL concurrently.
		// Verify its ownership before accepting a duplicate or lost reply.
		if ctx.Err() != nil {
			return ctx.Err()
		}
		acl, _, readErr := client.GetNetworkACL(name)
		if readErr != nil {
			return errors.Annotatef(err, "creating ACL %q", name)
		}
		return checkACLOwner(acl, metadata)
	}
	return nil
}

// updateNetworkACL retries the entire read/modify/write on an ETag conflict.
// Provisioners and the firewaller may update the same ACL from separate
// processes; a process-local mutex cannot prevent lost updates.
func updateNetworkACL(ctx context.Context, client networkACLServer, name string, metadata map[string]string, change func(*api.NetworkACLPut) error) error {
	for range aclUpdateAttempts {
		if err := ctx.Err(); err != nil {
			return err
		}
		acl, etag, err := client.GetNetworkACL(name)
		if err != nil {
			return errors.Annotatef(err, "reading ACL %q", name)
		}
		if err := checkACLOwner(acl, metadata); err != nil {
			return err
		}
		before := acl.Writable()
		after := acl.Writable()
		after.Config = maps.Clone(before.Config)
		if err := change(&after); err != nil {
			return errors.Annotatef(err, "updating ACL %q", name)
		}
		if reflect.DeepEqual(before, after) {
			return nil
		}
		op, err := client.UpdateNetworkACL(name, after, etag)
		if err == nil {
			err = op.WaitContext(ctx)
		}
		if api.StatusErrorCheck(err, http.StatusPreconditionFailed) {
			continue
		}
		return errors.Annotatef(err, "updating ACL %q", name)
	}
	return errors.Errorf("ACL %q kept changing during update", name)
}

func checkACLOwner(acl *api.NetworkACL, metadata map[string]string) error {
	for _, key := range []string{jujuACLModelKey, jujuACLInstanceKey, jujuACLKindKey} {
		if acl.Config[key] != metadata[key] {
			return errors.Errorf("ACL %q belongs to another owner (%s)", acl.Name, key)
		}
	}
	return nil
}
