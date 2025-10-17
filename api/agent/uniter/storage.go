// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package uniter

import (
	"context"

	"github.com/juju/errors"
	"github.com/juju/names/v6"

	"github.com/juju/juju/api/base"
	apiwatcher "github.com/juju/juju/api/watcher"
	apiservererrors "github.com/juju/juju/apiserver/errors"
	"github.com/juju/juju/core/watcher"
	"github.com/juju/juju/rpc/params"
)

type StorageAccessor struct {
	facade base.FacadeCaller
}

// NewStorageAccessor creates a StorageAccessor on the specified facade,
// and uses this name when calling through the caller.
func NewStorageAccessor(facade base.FacadeCaller) *StorageAccessor {
	return &StorageAccessor{facade}
}

// UnitStorageAttachments returns the IDs of a unit's storage attachments.
func (sa *StorageAccessor) UnitStorageAttachments(
	ctx context.Context, unitTag names.UnitTag,
) ([]params.StorageAttachmentId, error) {
	args := params.Entities{
		Entities: []params.Entity{{Tag: unitTag.String()}},
	}
	var results params.StorageAttachmentIdsResults
	err := sa.facade.FacadeCall(ctx, "UnitStorageAttachments", args, &results)
	if err != nil {
		return nil, errors.Trace(err)
	}
	if len(results.Results) != 1 {
		return nil, errors.Errorf("expected 1 result, got %d", len(results.Results))
	}
	result := results.Results[0]
	if result.Error != nil {
		return nil, result.Error
	}
	return result.Result.Ids, nil
}

// DestroyUnitStorageAttachments ensures that the specified unit's storage
// attachments will be removed at some point in the future.
// THIS SHOULD CALL REMOVESTORAGE (DYING + JOB).
func (sa *StorageAccessor) DestroyUnitStorageAttachments(ctx context.Context, unitTag names.UnitTag) error {
	args := params.Entities{
		Entities: []params.Entity{{Tag: unitTag.String()}},
	}
	var results params.ErrorResults
	err := sa.facade.FacadeCall(ctx, "DestroyUnitStorageAttachments", args, &results)
	if err != nil {
		return errors.Trace(err)
	}
	if len(results.Results) != 1 {
		return errors.Errorf("expected 1 result, got %d", len(results.Results))
	}
	result := results.Results[0]
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// WatchUnitStorageAttachments starts a watcher for changes to storage
// attachments related to the unit. The watcher will return the
// IDs of the corresponding storage instances.
func (sa *StorageAccessor) WatchUnitStorageAttachments(
	ctx context.Context, unitTag names.UnitTag,
) (watcher.StringsWatcher, error) {
	var results params.StringsWatchResults
	args := params.Entities{
		Entities: []params.Entity{{Tag: unitTag.String()}},
	}
	err := sa.facade.FacadeCall(ctx, "WatchUnitStorageAttachments", args, &results)
	if err != nil {
		return nil, err
	}
	if len(results.Results) != 1 {
		return nil, errors.Errorf("expected 1 result, got %d", len(results.Results))
	}
	result := results.Results[0]
	if result.Error != nil {
		return nil, result.Error
	}
	w := apiwatcher.NewStringsWatcher(sa.facade.RawAPICaller(), result)
	return w, nil
}

// StorageAttachment returns the storage attachment with the specified
// unit and storage tags.
func (sa *StorageAccessor) StorageAttachment(
	ctx context.Context, storageTag names.StorageTag, unitTag names.UnitTag,
) (params.StorageAttachment, error) {
	args := params.StorageAttachmentIds{
		Ids: []params.StorageAttachmentId{{
			StorageTag: storageTag.String(),
			UnitTag:    unitTag.String(),
		}},
	}
	var results params.StorageAttachmentResults
	err := sa.facade.FacadeCall(ctx, "StorageAttachments", args, &results)
	if err != nil {
		return params.StorageAttachment{}, errors.Trace(err)
	}
	if len(results.Results) != 1 {
		return params.StorageAttachment{}, errors.Errorf("expected 1 result, got %d", len(results.Results))
	}
	result := results.Results[0]
	if result.Error != nil {
		return params.StorageAttachment{}, apiservererrors.RestoreError(result.Error)
	}
	return result.Result, nil
}

// StorageAttachmentLife returns the lifecycle state of the storage attachments
// with the specified IDs.
func (sa *StorageAccessor) StorageAttachmentLife(
	ctx context.Context, ids []params.StorageAttachmentId,
) ([]params.LifeResult, error) {
	args := params.StorageAttachmentIds{Ids: ids}
	var results params.LifeResults
	err := sa.facade.FacadeCall(ctx, "StorageAttachmentLife", args, &results)
	if err != nil {
		return nil, errors.Trace(err)
	}
	if len(results.Results) != len(ids) {
		return nil, errors.Errorf("expected %d results, got %d", len(ids), len(results.Results))
	}
	return results.Results, nil
}

// WatchStorageAttachments starts a watcher for changes to the info
// of the storage attachment with the specified unit and storage tags.
func (sa *StorageAccessor) WatchStorageAttachment(
	ctx context.Context, storageTag names.StorageTag, unitTag names.UnitTag,
) (watcher.NotifyWatcher, error) {
	var results params.NotifyWatchResults
	args := params.StorageAttachmentIds{
		Ids: []params.StorageAttachmentId{{
			StorageTag: storageTag.String(),
			UnitTag:    unitTag.String(),
		}},
	}
	err := sa.facade.FacadeCall(ctx, "WatchStorageAttachments", args, &results)
	if err != nil {
		return nil, err
	}
	if len(results.Results) != 1 {
		return nil, errors.Errorf("expected 1 result, got %d", len(results.Results))
	}
	result := results.Results[0]
	if result.Error != nil {
		return nil, result.Error
	}
	w := apiwatcher.NewNotifyWatcher(sa.facade.RawAPICaller(), result)
	return w, nil
}

// RemoveStorageAttachment removes the storage attachment with the
// specified unit and storage tags from state. This method is only
// expected to succeed if the storage attachment is Dead.
// THIS SHOULD MARK DEAD.
func (sa *StorageAccessor) RemoveStorageAttachment(
	ctx context.Context, storageTag names.StorageTag, unitTag names.UnitTag,
) error {
	var results params.ErrorResults
	args := params.StorageAttachmentIds{
		Ids: []params.StorageAttachmentId{{
			StorageTag: storageTag.String(),
			UnitTag:    unitTag.String(),
		}},
	}
	err := sa.facade.FacadeCall(ctx, "RemoveStorageAttachments", args, &results)
	if err != nil {
		return err
	}
	if len(results.Results) != 1 {
		return errors.Errorf("expected 1 result, got %d", len(results.Results))
	}
	result := results.Results[0]
	if result.Error != nil {
		return result.Error
	}
	return nil
}
