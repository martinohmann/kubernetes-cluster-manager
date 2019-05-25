package cluster

import (
	"context"
	"os"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/credentials"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/file"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kcm"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kubernetes"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/manifest"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/provisioner"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/revision"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/template"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

const (
	dirMode os.FileMode = 0775
)

// Options are used to configure the cluster manager.
type Options struct {
	DryRun        bool   `json:"dryRun,omitempty" yaml:"dryRun,omitempty"`
	Values        string `json:"values,omitempty" yaml:"values,omitempty"`
	Deletions     string `json:"deletions,omitempty" yaml:"deletions,omitempty"`
	ManifestsDir  string `json:"manifestsDir,omitempty" yaml:"manifestsDir,omitempty"`
	TemplatesDir  string `json:"templatesDir,omitempty" yaml:"templatesDir,omitempty"`
	SkipManifests bool   `json:"skipManifests,omitempty" yaml:"skipManifests,omitempty"`
	AllManifests  bool   `json:"allManifests,omitempty" yaml:"allManifests,omitempty"`
	NoSave        bool   `json:"noSave,omitempty" yaml:"noSave,omitempty"`
}

// Manager is a Kubernetes cluster manager that will orchestrate changes to the
// cluster infrastructure and the cluster itself.
type Manager struct {
	credentialSource credentials.Source
	provisioner      provisioner.Provisioner
	renderer         template.Renderer
}

// NewManager creates a new cluster manager.
func NewManager(
	credentialSource credentials.Source,
	provisioner provisioner.Provisioner,
	renderer template.Renderer,
) *Manager {
	return &Manager{
		credentialSource: credentialSource,
		provisioner:      provisioner,
		renderer:         renderer,
	}
}

// Provision performs all steps necessary to create and setup a cluster and
// the required infrastructure. If a cluster already exists, it should
// update it if there are pending changes to be rolled out. Depending on
// the options it may or may not perform a dry run of the pending changes.
func (m *Manager) Provision(ctx context.Context, o *Options) error {
	var err error

	if !o.DryRun {
		err = m.provisioner.Provision(ctx)
	} else if r, ok := m.provisioner.(provisioner.Reconciler); ok {
		err = r.Reconcile(ctx)
	}

	if err != nil || o.SkipManifests {
		return err
	}

	return m.ApplyManifests(ctx, o)
}

// ApplyManifests renders and applies all manifests to the cluster. It also
// takes care of pending resource deletions that should be performed before
// and after applying.
func (m *Manager) ApplyManifests(ctx context.Context, o *Options) error {
	values, err := m.readValues(ctx, o.Values)
	if err != nil {
		return err
	}

	deletions, err := m.readDeletions(o.Deletions)
	if err != nil {
		return err
	}

	err = m.finalizeChanges(o, o.Values, values)
	if err != nil {
		return err
	}

	nextManifests, err := manifest.RenderDir(m.renderer, o.TemplatesDir, values)
	if err != nil {
		return err
	}

	currentManifests, err := manifest.ReadDir(o.ManifestsDir)
	if err != nil {
		return err
	}

	revisions := revision.NewSlice(currentManifests, nextManifests)

	creds, err := m.readCredentials(ctx, o)
	if err != nil {
		return err
	}

	kubectl := kubernetes.NewKubectl(creds)

	if !o.DryRun {
		if err := os.MkdirAll(o.ManifestsDir, dirMode); err != nil {
			return errors.WithStack(err)
		}

		log.Info("Waiting for cluster to become available...")

		if err := kubectl.WaitForCluster(ctx); err != nil {
			return err
		}
	}

	defer func() {
		m.finalizeChanges(o, o.Deletions, deletions)
	}()

	deletions.PreApply, err = processResourceDeletions(ctx, o, kubectl, deletions.PreApply)
	if err != nil {
		return err
	}

	upgrader := revision.NewUpgrader(kubectl, &revision.UpgraderOptions{
		DryRun:           o.DryRun,
		ManifestsDir:     o.ManifestsDir,
		NoSave:           o.NoSave,
		IncludeUnchanged: o.AllManifests,
		NoHooks:          false,
	})

	for _, revision := range revisions {
		if err = upgrader.Upgrade(ctx, revision); err != nil {
			return err
		}
	}

	deletions.PostApply, err = processResourceDeletions(ctx, o, kubectl, deletions.PostApply)

	return err
}

// Destroy deletes all applied manifests from a cluster and tears down the
// cluster infrastructure. Depending on the options it may or may not
// perform a dry run of the destruction process.
func (m *Manager) Destroy(ctx context.Context, o *Options) error {
	if !o.SkipManifests {
		if err := m.DeleteManifests(ctx, o); err != nil {
			return err
		}
	}

	if o.DryRun {
		log.Warn("Would destroy cluster infrastructure")
		return nil
	}

	return m.provisioner.Destroy(ctx)
}

// DeleteManifests renders and deletes all manifests from the cluster. It
// also takes care of other resource deletions that should be performed
// after the manifests have been deleted from the cluster.
func (m *Manager) DeleteManifests(ctx context.Context, o *Options) error {
	values, err := m.readValues(ctx, o.Values)
	if err != nil {
		return err
	}

	deletions, err := m.readDeletions(o.Deletions)
	if err != nil {
		return err
	}

	manifests, err := manifest.RenderDir(template.NewRenderer(), o.TemplatesDir, values)
	if err != nil {
		return err
	}

	revisions := revision.NewSlice(manifests, nil)

	creds, err := m.readCredentials(ctx, o)
	if err != nil {
		return err
	}

	kubectl := kubernetes.NewKubectl(creds)

	if !o.DryRun {
		if _, err := kubectl.ClusterInfo(ctx); err != nil {
			return err
		}
	}

	upgrader := revision.NewUpgrader(kubectl, &revision.UpgraderOptions{
		DryRun:           o.DryRun,
		ManifestsDir:     o.ManifestsDir,
		NoSave:           o.NoSave,
		IncludeUnchanged: o.AllManifests,
		NoHooks:          false,
	})

	for _, revision := range revisions.Reverse() {
		if err = upgrader.Upgrade(ctx, revision); err != nil {
			return err
		}
	}

	deletions.PreDestroy, _ = processResourceDeletions(ctx, o, kubectl, deletions.PreDestroy)

	return m.finalizeChanges(o, o.Deletions, deletions)
}

func (m *Manager) finalizeChanges(o *Options, filename string, v interface{}) error {
	buf, err := yaml.Marshal(v)
	if err != nil {
		return err
	}

	changeSet, err := file.NewChangeSet(filename, buf)
	if err != nil {
		return err
	}

	if changeSet.HasChanges() {
		log.Infof("Changes to %s:\n%s", filename, changeSet.Diff())
	} else {
		log.Infof("No changes to %s", filename)
	}

	if o.DryRun || o.NoSave {
		return nil
	}

	return changeSet.Apply()
}

func (m *Manager) readValues(ctx context.Context, filename string) (v kcm.Values, err error) {
	if err = file.ReadYAML(filename, &v); err != nil {
		return
	}

	if o, ok := m.provisioner.(provisioner.Outputter); ok {
		values, err := o.Output(ctx)
		if err == nil && len(values) > 0 {
			log.Info("Merging values from provisioner")
			v.Merge(values)
		}
	}

	return
}

func (m *Manager) readDeletions(filename string) (d *Deletions, err error) {
	err = file.ReadYAML(filename, &d)

	return
}

func (m *Manager) readCredentials(ctx context.Context, o *Options) (*credentials.Credentials, error) {
	creds, err := m.credentialSource.GetCredentials(ctx)
	if err != nil {
		return nil, err
	}

	if !o.DryRun && creds.Empty() {
		return nil, errors.New("Empty kubernetes credentials found! " +
			"Provide `kubeconfig` (and optionally `context`) or " +
			"`server` and `token` via the provisioner or set the corresponding --cluster-* flags")
	}

	c := *creds
	if c.Token != "" {
		c.Token = "<sensitive>"
	}

	log.Debugf("Using kubernetes credentials: %#v", c)

	return creds, nil
}
