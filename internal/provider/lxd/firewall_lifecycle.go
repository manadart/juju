// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package lxd

import (
	"context"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/canonical/lxd/shared/api"
	"github.com/juju/collections/set"
	"github.com/juju/errors"
)

func cleanupFailedInstanceACLs(ctx context.Context, srv Server, instanceName string) {
	if instanceName == "" {
		return
	}
	// Provisioning cancellation must not prevent rollback, but cleanup
	// should not hold the provisioner indefinitely either.
	cleanupCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 30*time.Second)
	defer cancel()
	if err := removeInstanceNetworkACLs(cleanupCtx, srv, []string{instanceName}); err != nil {
		logger.Errorf(cleanupCtx, "cleaning up firewall for instance %q: %v", instanceName, err)
	}
}

func removeInstanceNetworkACLs(ctx context.Context, srv Server, instanceNames []string) error {
	if len(instanceNames) == 0 {
		return nil
	}
	owners := set.NewStrings(instanceNames...)
	return removeNetworkACLs(ctx, srv, func(acl api.NetworkACL) bool {
		return acl.Config[jujuACLKindKey] == "machine" &&
			owners.Contains(acl.Config[jujuACLInstanceKey])
	})
}

func removeNetworkACLs(ctx context.Context, srv Server, matches func(api.NetworkACL) bool) error {
	client, acls, err := ownedNetworkACLs(ctx, srv)
	if err != nil {
		return err
	}
	for _, acl := range acls {
		if !matches(acl) {
			continue
		}
		if err := ctx.Err(); err != nil {
			return err
		}
		// LXD refuses deletion while an ACL is attached or referenced by
		// another ACL. Do not detach or alter somebody else's resources.
		op, err := client.DeleteNetworkACL(acl.Name)
		if err == nil {
			err = op.WaitContext(ctx)
		}
		if err != nil && !api.StatusErrorCheck(err, http.StatusNotFound) {
			return errors.Annotatef(err, "deleting firewall ACL %q", acl.Name)
		}
	}
	return nil
}

func (env *environ) adoptNetworkACLs(ctx context.Context, controllerUUID string) error {
	client, acls, err := ownedNetworkACLs(ctx, env.server())
	if err != nil {
		return err
	}
	for _, acl := range acls {
		if acl.Config[jujuACLModelKey] != env.uuid {
			continue
		}
		err := updateNetworkACL(ctx, client, acl.Name, acl.Config, func(put *api.NetworkACLPut) error {
			put.Config[jujuACLControllerKey] = controllerUUID
			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// ownedNetworkACLs uses both ownership metadata and deterministic names.
// Listing ACLs also finds leftovers after a failed provision or deletion,
// when there is no longer an instance from which to discover the ACL.
func ownedNetworkACLs(ctx context.Context, srv Server) (networkACLServer, []api.NetworkACL, error) {
	if err := ctx.Err(); err != nil {
		return nil, nil, err
	}
	if !srv.HasExtension("network_acl") {
		return nil, nil, nil
	}
	client, ok := srv.(networkACLServer)
	if !ok {
		return nil, nil, errors.NotSupportedf("LXD server ACL operations")
	}
	all, err := client.GetNetworkACLs()
	if err != nil {
		return nil, nil, errors.Annotate(err, "listing firewall ACLs")
	}
	var owned []api.NetworkACL
	for _, acl := range all {
		modelUUID := acl.Config[jujuACLModelKey]
		if modelUUID == "" {
			continue
		}
		var expectedName string
		switch acl.Config[jujuACLKindKey] {
		case "model":
			expectedName = modelACLName(modelUUID)
		case "global":
			expectedName = globalACLName(modelUUID)
		case "machine":
			instanceName := acl.Config[jujuACLInstanceKey]
			if instanceName == "" {
				continue
			}
			expectedName = machineACLName(modelUUID, instanceName)
		default:
			continue
		}
		if acl.Name == expectedName {
			owned = append(owned, acl)
		}
	}
	slices.SortFunc(owned, func(a, b api.NetworkACL) int {
		return strings.Compare(a.Name, b.Name)
	})
	return client, owned, nil
}
