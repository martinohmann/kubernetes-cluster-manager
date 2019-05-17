package upgrader

import "github.com/martinohmann/kubernetes-cluster-manager/pkg/manifest"

type Applier interface {
	ApplyManifest([]byte) error
	DeleteManifest([]byte) error
}

type Upgrader interface {
	Upgrade(*manifest.Revision) error
}

type Options struct {
	DryRun           bool
	IncludeUnchanged bool
}

type upgrader struct {
	applier          Applier
	dryRun           bool
	includeUnchanged bool
}

func New(applier Applier, o *Options) Upgrader {
	return &upgrader{
		applier:          applier,
		dryRun:           o.DryRun,
		includeUnchanged: o.IncludeUnchanged,
	}
}

func (u *upgrader) Upgrade(rev *manifest.Revision) error {
	c := rev.ChangeSet()
	err := u.deleteResources(c.RemovedResources, c.Hooks)
	if err != nil {
		return err
	}

	updates := append(c.ChangedResources, c.AddedResources...)

	if u.includeUnchanged {
		updates = append(updates, c.UnchangedResources...)
	}

	return u.applyResources(updates, c.Hooks)
}

func (u *upgrader) deleteResources(r manifest.ResourceSlice, h manifest.HookSliceMap) error {
	if len(r) == 0 {
		return nil
	}

	err := u.execHooks(manifest.HookTypePreDelete, h)
	if err != nil {
		return err
	}

	err = u.applier.DeleteManifest(r.Sort(manifest.DeleteOrder).Bytes())
	if err != nil {
		return err
	}

	return u.execHooks(manifest.HookTypePostDelete, h)
}

func (u *upgrader) applyResources(r manifest.ResourceSlice, h manifest.HookSliceMap) error {
	if len(r) == 0 {
		return nil
	}

	err := u.execHooks(manifest.HookTypePreDelete, h)
	if err != nil {
		return err
	}

	err = u.applier.ApplyManifest(r.Bytes())
	if err != nil {
		return err
	}

	return u.execHooks(manifest.HookTypePostDelete, h)
}

func (u *upgrader) execHooks(typ manifest.HookType, hooks manifest.HookSliceMap) error {
	if !hooks.Has(typ) {
		return nil
	}

	typeHooks := hooks.Get(typ)

	return u.applier.ApplyManifest(typeHooks.Bytes())
}
