package revision

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/gammazero/workerpool"
	multierror "github.com/hashicorp/go-multierror"
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
}

// upgrader is an implementations of Upgrader.
type upgrader struct {
	client           Client
	dryRun           bool
	includeUnchanged bool
	noHooks          bool
	noSave           bool
	manifestsDir     string
}

// NewUpgrader creates a new Upgrader with client and options.
func NewUpgrader(client Client, o *UpgraderOptions) Upgrader {
	if o == nil {
		o = &UpgraderOptions{}
	}

	u := &upgrader{
		client:           client,
		dryRun:           o.DryRun,
		includeUnchanged: o.IncludeUnchanged,
		noHooks:          o.NoHooks,
		noSave:           o.NoSave,
		manifestsDir:     o.ManifestsDir,
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

	if diff := rev.Diff(); diff != "" {
		log.Infof("changes to component %s:\n%s", color.YellowString(manifest.Name), diff)
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

func (u *upgrader) processManifestDeletion(ctx context.Context, manifest *manifest.Manifest) error {
	return u.wrapHooks(ctx, manifest.Hooks, hook.TypeDelete, func() error {
		log.Warnf("removing component %s", color.YellowString(manifest.Name))

		err := u.deleteResources(ctx, manifest.Resources)
		if err != nil {
			return err
		}

		claims := manifest.Resources.PersistentVolumeClaimsForDeletion()

		return u.deletePersistentVolumeClaims(ctx, claims)
	})
}

func (u *upgrader) processManifestCreation(ctx context.Context, manifest *manifest.Manifest) error {
	return u.wrapHooks(ctx, manifest.Hooks, hook.TypeCreate, func() error {
		log.Warnf("creating component %s", color.YellowString(manifest.Name))

		return u.applyResources(ctx, manifest.Resources)
	})
}

func (u *upgrader) processManifestUpdate(ctx context.Context, manifest *manifest.Manifest, changeSet *ChangeSet) error {
	return u.wrapHooks(ctx, manifest.Hooks, hook.TypeUpgrade, func() error {
		log.Infof("updating component %s", color.YellowString(manifest.Name))

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

		return u.applyResources(ctx, resources)
	})
}

func (u *upgrader) deletePersistentVolumeClaims(ctx context.Context, claims []resource.Head) error {
	for _, claim := range claims {
		if u.dryRun {
			log.Warnf("would delete %s", claim.String())
			continue
		}

		err := u.client.DeleteResource(ctx, claim)
		if err != nil {
			return err
		}
	}

	return nil
}

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

func (u *upgrader) deleteResources(ctx context.Context, r resource.Slice) error {
	if len(r) == 0 {
		return nil
	}

	if u.dryRun {
		log.Infof("would delete %d resources:\n%s", len(r), indent(r.String()))
		return nil
	}

	log.Infof("deleting %d resources", len(r))

	return u.client.DeleteManifest(ctx, r.Sort(resource.DeleteOrder).Bytes())
}

func (u *upgrader) applyResources(ctx context.Context, r resource.Slice) error {
	if len(r) == 0 {
		return nil
	}

	if u.dryRun {
		log.Infof("would apply %d resources:\n%s", len(r), indent(r.String()))
		return nil
	}

	log.Infof("applying %d resources", len(r))

	return u.client.ApplyManifest(ctx, r.Sort(resource.ApplyOrder).Bytes())
}

func (u *upgrader) execHooks(ctx context.Context, hooks hook.Slice) error {
	if u.noHooks || hooks == nil {
		return nil
	}

	r := hooks.Resources()

	if len(r) == 0 {
		return nil
	}

	if u.dryRun {
		log.Infof("would execute %d hooks:\n%s", len(hooks), indent(hooks.String()))
		return nil
	}

	log.Infof("executing %d hooks", len(hooks))

	err := u.client.DeleteManifest(ctx, r.Sort(resource.DeleteOrder).Bytes())
	if err != nil {
		return err
	}

	err = u.client.ApplyManifest(ctx, r.Sort(resource.ApplyOrder).Bytes())
	if err != nil {
		return err
	}

	return u.waitForHooks(ctx, hooks)
}

func (u *upgrader) waitForHooks(ctx context.Context, hooks hook.Slice) error {
	pool := workerpool.New(MaxWorkers)
	errs := &multierror.Error{ErrorFormat: errorFormatFunc}

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
