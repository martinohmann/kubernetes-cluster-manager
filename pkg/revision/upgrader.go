package revision

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/kr/text"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/hook"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/resource"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Client applies and deletes manifests from a cluster.
type Client interface {
	// ApplyManifest applies raw manifest bytes.
	ApplyManifest(context.Context, []byte) error

	// DeleteManifest deletes raw manifest bytes.
	DeleteManifest(context.Context, []byte) error
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

	c := rev.ChangeSet()

	if diff := rev.Diff(); diff != "" {
		log.Infof("Changes to component %s:\n%s", color.YellowString(manifest.Name), diff)
	}

	if rev.IsRemoval() {
		err = u.wrapHooks(ctx, manifest.Hooks, hook.TypeDelete, func() error {
			log.Warnf("Removing component %s", color.YellowString(manifest.Name))

			return u.deleteResources(ctx, manifest.Resources)
		})

		if err == nil && !u.dryRun {
			err := os.Remove(filename)
			if err != nil && !os.IsNotExist(err) {
				return errors.WithStack(err)
			}
		}

		return err
	}

	if rev.IsInitial() {
		err = u.wrapHooks(ctx, manifest.Hooks, hook.TypeCreate, func() error {
			log.Warnf("Creating component %s", color.YellowString(manifest.Name))

			return u.applyResources(ctx, manifest.Resources)
		})
	} else if c.HasResourceChanges() || u.includeUnchanged {
		err = u.wrapHooks(ctx, manifest.Hooks, hook.TypeUpgrade, func() error {
			log.Infof("Updating component %s", color.YellowString(manifest.Name))

			err := u.deleteResources(ctx, c.RemovedResources)
			if err != nil {
				return err
			}

			resources := append(c.AddedResources, c.UpdatedResources...)

			if u.includeUnchanged {
				resources = append(resources, c.UnchangedResources...)
			}

			return u.applyResources(ctx, resources)
		})
	}

	if err == nil && !u.dryRun && !u.noSave {
		return ioutil.WriteFile(filename, manifest.Content(), 0660)
	}

	return err
}

func (u *upgrader) wrapHooks(ctx context.Context, hooks hook.SliceMap, hookTypes hook.TypePair, f func() error) error {
	err := u.execHooks(ctx, hookTypes.Pre, hooks)
	if err != nil {
		return err
	}

	if err := f(); err != nil {
		return err
	}

	return u.execHooks(ctx, hookTypes.Post, hooks)
}

func (u *upgrader) deleteResources(ctx context.Context, r resource.Slice) error {
	if len(r) == 0 {
		return nil
	}

	if u.dryRun {
		log.Infof("Would delete %d resources:\n%s", len(r), text.Indent(r.String(), "  "))
		return nil
	}

	log.Infof("Deleting %d resources", len(r))

	return u.client.DeleteManifest(ctx, r.Sort(resource.DeleteOrder).Bytes())
}

func (u *upgrader) applyResources(ctx context.Context, r resource.Slice) error {
	if len(r) == 0 {
		return nil
	}

	if u.dryRun {
		log.Infof("Would apply %d resources:\n%s", len(r), text.Indent(r.String(), "  "))
		return nil
	}

	log.Infof("Applying %d resources", len(r))

	return u.client.ApplyManifest(ctx, r.Sort(resource.ApplyOrder).Bytes())
}

func (u *upgrader) execHooks(ctx context.Context, typ string, hookMap hook.SliceMap) error {
	hooks, ok := hookMap[typ]
	if u.noHooks || !ok {
		return nil
	}

	r := hooks.Resources()

	if len(r) == 0 {
		return nil
	}

	if u.dryRun {
		log.Infof("Would execute %d %s hooks:\n%s", len(r), typ, text.Indent(r.String(), "  "))
		return nil
	}

	log.Infof("Executing %d %s hooks", len(r), typ)

	err := u.client.DeleteManifest(ctx, r.Sort(resource.DeleteOrder).Bytes())
	if err != nil {
		return err
	}

	return u.client.ApplyManifest(ctx, r.Sort(resource.ApplyOrder).Bytes())
}
