package revision

import (
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/hook"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/resource"
)

type Applier interface {
	ApplyManifest([]byte) error
	DeleteManifest([]byte) error
}

type Upgrader interface {
	Upgrade(*Revision) error
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

func (u *upgrader) Upgrade(rev *Revision) error {
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

func (u *upgrader) deleteResources(r resource.Slice, h hook.SliceMap) error {
	if len(r) == 0 {
		return nil
	}

	err := u.execHooks(hook.TypePreDelete, h)
	if err != nil {
		return err
	}

	err = u.applier.DeleteManifest(r.Sort(resource.DeleteOrder).Bytes())
	if err != nil {
		return err
	}

	return u.execHooks(hook.TypePostDelete, h)
}

func (u *upgrader) applyResources(r resource.Slice, h hook.SliceMap) error {
	if len(r) == 0 {
		return nil
	}

	err := u.execHooks(hook.TypePreApply, h)
	if err != nil {
		return err
	}

	err = u.applier.ApplyManifest(r.Bytes())
	if err != nil {
		return err
	}

	return u.execHooks(hook.TypePostApply, h)
}

func (u *upgrader) execHooks(typ hook.Type, hooks hook.SliceMap) error {
	if !hooks.Has(typ) {
		return nil
	}

	typeHooks := hooks.Get(typ)

	return u.applier.ApplyManifest(typeHooks.Resources().Bytes())
}
