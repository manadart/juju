// Copyright 2021 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package application

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/juju/errors"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/juju/juju/caas"
	"github.com/juju/juju/internal/provider/kubernetes/resources"
	"github.com/juju/juju/internal/provider/kubernetes/scale"
	"github.com/juju/juju/internal/provider/kubernetes/utils"
	"github.com/juju/juju/internal/storage"
)

// Scale reconciles the application's replica count and stops when ctx is
// cancelled.
func (a *app) Scale(ctx context.Context, scaleTo int) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	switch a.deploymentType {
	case caas.DeploymentStateful:
		return scale.PatchReplicasToScale(
			ctx,
			a.name,
			int32(scaleTo),
			scale.StatefulSetScalePatcher(a.client.AppsV1().StatefulSets(a.namespace)),
		)
	case caas.DeploymentStateless:
		return scale.PatchReplicasToScale(
			ctx,
			a.name,
			int32(scaleTo),
			scale.DeploymentScalePatcher(a.client.AppsV1().Deployments(a.namespace)),
		)
	default:
		return errors.NotSupportedf(
			"application %q deployment type %q cannot be scaled",
			a.name, a.deploymentType)
	}
}

// ScaleRange reconciles the application to the ordinal range
// [startOrdinal, startOrdinal+replicaCount) and stops when ctx is cancelled.
// Sparse ordinal ranges are not supported. Both values must be non-negative.
func (a *app) ScaleRange(ctx context.Context, replicaCount, startOrdinal int) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if replicaCount < 0 || startOrdinal < 0 {
		return errors.NotValidf("replica count %d or start ordinal %d", replicaCount, startOrdinal)
	}
	switch a.deploymentType {
	case caas.DeploymentStateful:
		patch, err := json.Marshal(map[string]any{
			"spec": map[string]any{
				"replicas": replicaCount,
				"ordinals": map[string]any{
					"start": startOrdinal,
				},
			},
		})
		if err != nil {
			return errors.Trace(err)
		}
		_, err = a.client.AppsV1().StatefulSets(a.namespace).Patch(
			ctx, a.name, types.MergePatchType, patch, meta.PatchOptions{})
		if k8serrors.IsNotFound(err) {
			return errors.NotFoundf("statefulset %q", a.name)
		} else if err != nil {
			return errors.Annotatef(err, "scaling statefulset %q", a.name)
		}
		return nil
	case caas.DeploymentStateless:
		if startOrdinal != 0 {
			return errors.NotSupportedf("scaling deployment %q from ordinal %d", a.name, startOrdinal)
		}
		return scale.PatchReplicasToScale(
			ctx,
			a.name,
			int32(replicaCount),
			scale.DeploymentScalePatcher(a.client.AppsV1().Deployments(a.namespace)),
		)
	default:
		return errors.NotSupportedf(
			"application %q deployment type %q cannot be scaled",
			a.name, a.deploymentType)
	}
}

// currentScale returns the current scale in use for the applications. i.e how
// many units is Kubernetes currently running for application x.
func (a *app) currentScale(ctx context.Context) (int, error) {
	switch a.deploymentType {
	case caas.DeploymentStateful:
		ss, err := a.client.AppsV1().StatefulSets(a.namespace).Get(ctx, a.name, meta.GetOptions{})
		if k8serrors.IsNotFound(err) {
			err = errors.WithType(err, errors.NotFound)
		}
		if err != nil {
			return 0, fmt.Errorf("fetching scale for application %q statefuleset: %w",
				a.name, err)
		}

		return int(*ss.Spec.Replicas), nil

	default:
		return 0, fmt.Errorf("application %q deployment type %q is not supported for fetching scale",
			a.name, a.deploymentType)
	}
}

// UnitsToRemove returns the names of units that need to be removed to reach the desired scale.
func (a *app) UnitsToRemove(ctx context.Context, desiredScale int) ([]string, error) {
	var unitsToRemove []string
	currentScale, err := a.currentScale(ctx)
	if err != nil {
		return unitsToRemove, err
	}

	numUnitsToRemove := desiredScale - currentScale
	if numUnitsToRemove >= 0 {
		return unitsToRemove, nil
	}

	for ; numUnitsToRemove != 0; numUnitsToRemove++ {
		unitsToRemove = append(unitsToRemove, fmt.Sprintf("%s/%d", a.name, currentScale+numUnitsToRemove))
	}

	return unitsToRemove, nil
}

// EnsurePVCs creates missing claims for the desired units before scaling.
// Existing claims are checked but never patched or replaced.
func (a *app) EnsurePVCs(
	ctx context.Context,
	filesystems []storage.KubernetesFilesystemParams,
	filesystemUnitAttachments map[string][]storage.KubernetesFilesystemUnitAttachmentParams,
	storageUniqueID string,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if len(filesystemUnitAttachments) == 0 {
		return nil
	}

	pvcTemplateNames, err := a.storageNameToPVCTemplateNames(filesystems)
	if err != nil {
		return errors.Annotatef(err, "mapping pvc template names for app %q", a.name)
	}
	// The StatefulSet retains its original claim templates. Use their names
	// even when no existing PVC has been reported to Juju yet.
	sts, err := a.client.AppsV1().StatefulSets(a.namespace).Get(ctx, a.name, meta.GetOptions{})
	if err != nil {
		return errors.Annotate(err, "getting claim templates before scaling")
	}
	for _, fs := range filesystems {
		for _, template := range sts.Spec.VolumeClaimTemplates {
			if _, err := utils.MatchStorageMetaLabelVersion(template.ObjectMeta, fs.StorageName); err == nil {
				pvcTemplateNames[a.appStorageName(fs.StorageName)] = template.Name
				break
			}
		}
	}

	pvcNameGetter := a.pvcNameGetter(pvcTemplateNames, storageUniqueID)

	storageClasses, err := resources.ListStorageClass(ctx, a.client.StorageV1().StorageClasses(), meta.ListOptions{})
	if err != nil {
		return errors.Trace(err)
	}
	storageClassMap := make(map[string]resources.StorageClass)
	for _, v := range storageClasses {
		storageClassMap[v.Name] = v
	}

	for _, fs := range filesystems {
		if len(filesystemUnitAttachments[fs.StorageName]) == 0 {
			continue
		}
		name := a.appStorageName(fs.StorageName)
		_, pvc, sc, err := a.filesystemToVolumeInfo(name, fs, storageClassMap, pvcNameGetter)
		if err != nil {
			return errors.Trace(err)
		}
		if pvc != nil {
			if sc != nil {
				if err := resources.NewStorageClass(a.client.StorageV1().StorageClasses(), sc.Name, sc).Apply(ctx); err != nil {
					return err
				}
				storageClassMap[sc.Name] = resources.StorageClass{StorageClass: *sc}
			}
			err := a.handleVolumeAttachment(ctx, filesystemUnitAttachments, *pvc, fs.StorageName)
			if err != nil {
				return errors.Trace(err)
			}
		}
	}
	return nil
}
