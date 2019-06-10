package cluster

import (
	"context"
	"io/ioutil"
	"os"

	"github.com/imdario/mergo"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/credentials"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/diff"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/file"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kubernetes"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/log"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/manifest"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/provisioner"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/revision"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/template"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

const (
	dirMode os.FileMode = 0775
)

// Options are used to configure the cluster manager.
type Options struct {
	DryRun        bool   `json:"dryRun,omitempty" yaml:"dryRun,omitempty"`
	Values        string `json:"values,omitempty" yaml:"values,omitempty"`
	ManifestsDir  string `json:"manifestsDir,omitempty" yaml:"manifestsDir,omitempty"`
	TemplatesDir  string `json:"templatesDir,omitempty" yaml:"templatesDir,omitempty"`
	SkipManifests bool   `json:"skipManifests,omitempty" yaml:"skipManifests,omitempty"`
	AllManifests  bool   `json:"allManifests,omitempty" yaml:"allManifests,omitempty"`
	NoSave        bool   `json:"noSave,omitempty" yaml:"noSave,omitempty"`
	NoHooks       bool   `json:"noHooks,omitempty" yaml:"noHooks,omitempty"`
	FullDiff      bool   `json:"fullDiff,omitempty" yaml:"fullDiff,omitempty"`
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

// ApplyManifests applies all manifests to the cluster.
func (m *Manager) ApplyManifests(ctx context.Context, o *Options) error {
	values, err := m.readValues(ctx, o.Values)
	if err != nil {
		return err
	}

	err = m.updateValuesFile(o.Values, values, o)
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

		logrus.Info("waiting for cluster to become available...")

		if err := kubectl.WaitForCluster(ctx); err != nil {
			return err
		}
	}

	upgrader := revision.NewUpgrader(kubectl, &revision.UpgraderOptions{
		DryRun:           o.DryRun,
		ManifestsDir:     o.ManifestsDir,
		NoSave:           o.NoSave,
		IncludeUnchanged: o.AllManifests,
		NoHooks:          o.NoHooks,
		FullDiff:         o.FullDiff,
	})

	for _, revision := range revisions {
		if err = upgrader.Upgrade(ctx, revision); err != nil {
			return err
		}
	}

	return nil
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
		logrus.Warn("would destroy cluster infrastructure")
		return nil
	}

	return m.provisioner.Destroy(ctx)
}

// DeleteManifests deletes all manifests from the cluster in reverse apply
// order.
func (m *Manager) DeleteManifests(ctx context.Context, o *Options) error {
	var manifests []*manifest.Manifest
	var err error

	if o.AllManifests {
		// To be able to attempt the deletion of manifests that are already
		// removed from the manifests dir we render them again.
		var values map[string]interface{}

		values, err = m.readValues(ctx, o.Values)
		if err != nil {
			return err
		}

		manifests, err = manifest.RenderDir(template.NewRenderer(), o.TemplatesDir, values)
	} else {
		manifests, err = manifest.ReadDir(o.ManifestsDir)
	}

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
		NoHooks:          o.NoHooks,
		FullDiff:         o.FullDiff,
	})

	for _, revision := range revisions.Reverse() {
		if err = upgrader.Upgrade(ctx, revision); err != nil {
			return err
		}
	}

	return nil
}

func (m *Manager) updateValuesFile(filename string, v map[string]interface{}, o *Options) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	buf, err := yaml.Marshal(v)
	if err != nil {
		return err
	}

	diffOptions := diff.Options{
		Filename: filename,
		A:        content,
		B:        buf,
	}

	diff.NewPrinter(log.LineWriter(logrus.Info)).Print(diffOptions)

	if o.DryRun || o.NoSave {
		return nil
	}

	return ioutil.WriteFile(filename, buf, 0660)
}

func (m *Manager) readValues(ctx context.Context, filename string) (v map[string]interface{}, err error) {
	if err = file.ReadYAML(filename, &v); err != nil {
		return
	}

	if o, ok := m.provisioner.(provisioner.Outputter); ok {
		var values map[string]interface{}

		values, err = o.Output(ctx)
		if err == nil && len(values) > 0 {
			logrus.Info("merging values from provisioner")
			err = mergo.Merge(&v, values, mergo.WithOverride)
		}
	}

	return
}

func (m *Manager) readCredentials(ctx context.Context, o *Options) (*credentials.Credentials, error) {
	creds, err := m.credentialSource.GetCredentials(ctx)
	if err != nil {
		return nil, err
	}

	if !o.DryRun && creds.Empty() {
		return nil, errors.New("empty kubernetes credentials found, " +
			"provide `kubeconfig` (and optionally `context`) or " +
			"`server` and `token` via the provisioner or set the corresponding --cluster-* flags")
	}

	c := *creds
	if c.Token != "" {
		c.Token = "<sensitive>"
	}

	logrus.WithFields(logrus.Fields{
		"kubeconfig": c.Kubeconfig,
		"context":    c.Context,
		"server":     c.Server,
		"token":      c.Token,
	}).Debugf("using kubernetes credentials")

	return creds, nil
}
