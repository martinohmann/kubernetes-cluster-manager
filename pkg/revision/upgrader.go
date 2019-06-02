package revision

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/gammazero/workerpool"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/kr/text"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/hook"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kubernetes"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/manifest"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/resource"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	// MaxWorkers is the maximum number of go-routines to use for concurrent
	// actions.
	MaxWorkers = 10
)

// Client applies and deletes manifests from a cluster.
type Client interface {
	// ApplyManifest applies raw manifest bytes.
	ApplyManifest(context.Context, []byte) error

	// DeleteManifest deletes raw manifest bytes.
	DeleteManifest(context.Context, []byte) error

	// DeleteResource deletes a resource by its kind, name and namespace.
	DeleteResource(context.Context, resource.Head) error

	// Wait waits for a resource condition to be met.
	Wait(context.Context, kubernetes.WaitOptions) error
}

// Upgrader handles revision upgrades.
type Upgrader interface {
	// Upgrader takes a context and a revision and performs an upgrade.
	// Depending on the type of revision it will carry out a complete creation
	// or deletion of the revision's manifest resources or just do partial
	// updates of resources that have been changed. It also executes hooks
	// before and after processing the revision.
	Upgrade(context.Context, *Revision) error
}

// UpgraderOptions configure an Upgrader.
type UpgraderOptions struct {
	DryRun           bool
	IncludeUnchanged bool
	NoHooks          bool
	NoSave           bool
	ManifestsDir     string
	FullDiff         bool
}

// upgrader is an implementations of Upgrader.
type upgrader struct {
	client           Client
	printer          *resource.Printer
	dryRun           bool
	includeUnchanged bool
	noHooks          bool
	noSave           bool
	manifestsDir     string
	fullDiff         bool
}

// NewUpgrader creates a new Upgrader with client and options.
func NewUpgrader(client Client, o *UpgraderOptions) Upgrader {
	if o == nil {
		o = &UpgraderOptions{}
	}

	u := &upgrader{
		client:           client,
		printer:          resource.NewPrinter(writerFunc(log.Info)),
		dryRun:           o.DryRun,
		includeUnchanged: o.IncludeUnchanged,
		noHooks:          o.NoHooks,
		noSave:           o.NoSave,
		manifestsDir:     o.ManifestsDir,
		fullDiff:         o.FullDiff,
	}

	return u
}

// Upgrade implements Upgrader.
func (u *upgrader) Upgrade(ctx context.Context, rev *Revision) error {
	var err error

	if !rev.IsValid() {
		return errors.New("cannot perform upgrade on invalid revision")
	}

	manifest := rev.Manifest()
	filename := filepath.Join(u.manifestsDir, manifest.Filename())

	changeSet := rev.ChangeSet()

	if u.fullDiff {
		if diff := rev.diff(); diff != "" {
			log.Infof("changes to component %s:\n\n%s\n", color.YellowString(manifest.Name), text.Indent(diff, "  "))
		} else {
			log.Infof("no changes to component %s", color.YellowString(manifest.Name))
		}
	}

	if rev.IsRemoval() {
		err = u.processManifestDeletion(ctx, manifest)
		if err == nil && !u.dryRun {
			err := os.Remove(filename)
			if err != nil && !os.IsNotExist(err) {
				return errors.WithStack(err)
			}
		}

		return err
	}

	if rev.IsInitial() {
		err = u.processManifestCreation(ctx, manifest)
	} else if changeSet.HasResourceChanges() || u.includeUnchanged {
		err = u.processManifestUpdate(ctx, manifest, changeSet)
	}

	if err == nil && !u.dryRun && !u.noSave {
		return ioutil.WriteFile(filename, manifest.Content(), 0660)
	}

	return err
}

// processManifestDeletion delete all manifest resources from the cluster. It
// will run the pre-delete and post-delete hooks and also remove
// PersistentVolumeClaims of StatefulSets that enabled the delete-pvcs deletion
// policy.
func (u *upgrader) processManifestDeletion(ctx context.Context, manifest *manifest.Manifest) error {
	return u.wrapHooks(ctx, manifest.Hooks, hook.TypeDelete, func() error {
		log.Warnf("removing component %s", color.YellowString(manifest.Name))

		u.printer.PrintSlice(manifest.Resources)

		err := u.deleteResources(ctx, manifest.Resources)
		if err != nil {
			return err
		}

		claims := manifest.Resources.PersistentVolumeClaimsForDeletion()

		return u.deletePersistentVolumeClaims(ctx, claims)
	})
}

// processManifestCreation applies all manifest resources to the cluster. It
// will run the pre-create and post-create hooks.
func (u *upgrader) processManifestCreation(ctx context.Context, manifest *manifest.Manifest) error {
	return u.wrapHooks(ctx, manifest.Hooks, hook.TypeCreate, func() error {
		log.Warnf("creating component %s", color.YellowString(manifest.Name))

		u.printer.PrintSlice(manifest.Resources)

		return u.applyResources(ctx, manifest.Resources)
	})
}

// processManifestUpdate will update resources that have changed and delete
// resources that disappeared from the manifest and also remove
// PersistentVolumeClaims of StatefulSets that were removed and had the
// delete-pvcs deletion policy enabled. It will run the pre-upgrade and
// post-upgrade hooks.
func (u *upgrader) processManifestUpdate(ctx context.Context, manifest *manifest.Manifest, changeSet *ChangeSet) error {
	return u.wrapHooks(ctx, manifest.Hooks, hook.TypeUpgrade, func() error {
		log.Infof("updating component %s", color.YellowString(manifest.Name))

		u.printer.PrintSlice(changeSet.RemovedResources)

		err := u.deleteResources(ctx, changeSet.RemovedResources)
		if err != nil {
			return err
		}

		claims := changeSet.RemovedResources.PersistentVolumeClaimsForDeletion()

		err = u.deletePersistentVolumeClaims(ctx, claims)
		if err != nil {
			return err
		}

		resources := append(changeSet.AddedResources, changeSet.UpdatedResources...)

		if u.includeUnchanged {
			resources = append(resources, changeSet.UnchangedResources...)
		}

		u.printer.PrintSlice(resources)

		return u.applyResources(ctx, resources)
	})
}

// deletePersistentVolumeClaims removes the PersistentVolumeClaims in the
// claims slice from the cluster. This will be a no-op when dry-run mode is
// enabled.
func (u *upgrader) deletePersistentVolumeClaims(ctx context.Context, claims resource.Slice) error {
	if len(claims) == 0 {
		return nil
	}

	log.Info("removing persistent volume claims")

	u.printer.PrintSlice(claims)

	if u.dryRun {
		log.Debug("skipping pvc deletions due to dry run")
		return nil
	}

	for _, claim := range claims {
		err := u.client.DeleteResource(ctx, resource.Head{
			Kind: claim.Kind,
			Metadata: resource.Metadata{
				Name:      claim.Name,
				Namespace: claim.Namespace,
			},
		})

		if err != nil {
			return err
		}
	}

	return nil
}

// wrapHooks wraps given func f with hooks of hookType.
func (u *upgrader) wrapHooks(ctx context.Context, hooks hook.SliceMap, hookTypes hook.TypePair, f func() error) error {
	err := u.execHooks(ctx, hooks[hookTypes.Pre])
	if err != nil {
		return err
	}

	if err := f(); err != nil {
		return err
	}

	return u.execHooks(ctx, hooks[hookTypes.Post])
}

// deleteResources deletes all resources in r from the cluster. This will be a
// no-op when dry-run mode is enabled.
func (u *upgrader) deleteResources(ctx context.Context, r resource.Slice) error {
	if len(r) == 0 {
		return nil
	}

	if u.dryRun {
		log.Debug("skipping resource deletions due to dry run")
		return nil
	}

	return u.client.DeleteManifest(ctx, r.Sort(resource.DeleteOrder).Bytes())
}

// applyResources applies all resources in r to the cluster. This will be a
// no-op when dry-run mode is enabled.
func (u *upgrader) applyResources(ctx context.Context, r resource.Slice) error {
	if len(r) == 0 {
		return nil
	}

	if u.dryRun {
		log.Debug("skipping resource updates due to dry run")
		return nil
	}

	return u.client.ApplyManifest(ctx, r.Sort(resource.ApplyOrder).Bytes())
}

// execHooks executes given hooks. It will delete the hooks from the cluster
// prior to applying them to ensure that Job resources are recreated properly.
// This is a no-op if dry-run mode is enabled.
func (u *upgrader) execHooks(ctx context.Context, hooks hook.Slice) error {
	if u.noHooks || hooks == nil {
		return nil
	}

	r := hooks.Resources()

	if len(r) == 0 {
		return nil
	}

	log.Infof("executing %d hooks", len(hooks))

	u.printer.PrintSlice(r)

	if u.dryRun {
		log.Debug("skipping hooks due to dry run")
		return nil
	}

	err := u.deleteResources(ctx, r)
	if err != nil {
		return err
	}

	err = u.applyResources(ctx, r)
	if err != nil {
		return err
	}

	return u.waitForHooks(ctx, hooks)
}

// waitForHooks waits for the hooks WaitFor condition to be met. Will wait for
// a maximum of MaxWorker jobs in parallel.
func (u *upgrader) waitForHooks(ctx context.Context, hooks hook.Slice) error {
	pool := workerpool.New(MaxWorkers)
	errs := &multierror.Error{}

	for _, hook := range hooks {
		if hook.WaitFor == "" {
			continue
		}

		h := hook

		pool.Submit(func() {
			log.Infof("waiting for hook %s", h.String())

			err := u.client.Wait(ctx, kubernetes.WaitOptions{
				Kind:      h.Resource.Kind,
				Name:      h.Resource.Name,
				Namespace: h.Resource.Namespace,
				For:       h.WaitFor,
				Timeout:   h.WaitTimeout,
			})

			if err != nil {
				errs = multierror.Append(errs, errors.Wrapf(err, "waiting for hook %s failed", h))
			} else if h.DeleteAfterCompletion {
				err := u.client.DeleteManifest(ctx, h.Resource.Content)
				if err != nil {
					errs = multierror.Append(errs, errors.Wrapf(err, "failed to delete hook %s", h))
				}
			}
		})
	}

	pool.StopWait()

	return errs.ErrorOrNil()
}

// writerFunc is a printf-style func which satisfies the io.Writer interface.
type writerFunc func(args ...interface{})

// Write implements io.Writer.
func (w writerFunc) Write(p []byte) (n int, err error) {
	w(string(p))

	return len(p), nil
}
