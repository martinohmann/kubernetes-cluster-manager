package cluster

import (
	"os"
	"path/filepath"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/credentials"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/file"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kcm"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kubernetes"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/provisioner"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/renderer"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

const (
	dirMode os.FileMode = 0775
)

// Options are used to configure the cluster manager.
type Options struct {
	DryRun        bool   `json:"dryRun" yaml:"dryRun"`
	Values        string `json:"values" yaml:"values"`
	Deletions     string `json:"deletions" yaml:"deletions"`
	ManifestsDir  string `json:"manifestsDir" yaml:"manifestsDir"`
	SkipManifests bool   `json:"skipManifests" yaml:"skipManifests"`
	AllManifests  bool   `json:"allManifests" yaml:"allManifests"`
}

// Manager is a Kubernetes cluster manager that will orchestrate changes to the
// cluster infrastructure and the cluster itself.
type Manager struct {
	credentialSource credentials.Source
	provisioner      provisioner.Provisioner
	renderer         renderer.Renderer
	logger           *log.Logger
}

// NewManager creates a new cluster manager.
func NewManager(
	credentialSource credentials.Source,
	provisioner provisioner.Provisioner,
	renderer renderer.Renderer,
	logger *log.Logger,
) *Manager {
	return &Manager{
		credentialSource: credentialSource,
		provisioner:      provisioner,
		renderer:         renderer,
		logger:           logger,
	}
}

// Provision performs all steps necessary to create and setup a cluster and
// the required infrastructure. If a cluster already exists, it should
// update it if there are pending changes to be rolled out. Depending on
// the options it may or may not perform a dry run of the pending changes.
func (m *Manager) Provision(o *Options) error {
	var err error

	if !o.DryRun {
		err = m.provisioner.Provision()
	} else if r, ok := m.provisioner.(provisioner.Reconciler); ok {
		err = r.Reconcile()
	}

	if err != nil || o.SkipManifests {
		return err
	}

	return m.ApplyManifests(o)
}

// ApplyManifests renders and applies all manifests to the cluster. It also
// takes care of pending resource deletions that should be performed before
// and after applying.
func (m *Manager) ApplyManifests(o *Options) error {
	values, err := m.readValues(o.Values)
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

	manifests, err := m.renderer.RenderManifests(values)
	if err != nil {
		return err
	}

	creds, err := m.readCredentials(o)
	if err != nil {
		return err
	}

	kubectl := kubernetes.NewKubectl(creds)

	if !o.DryRun {
		if err := os.MkdirAll(o.ManifestsDir, dirMode); err != nil {
			return errors.WithStack(err)
		}

		m.logger.Info("Waiting for cluster to become available...")

		if err := kubectl.WaitForCluster(); err != nil {
			return err
		}
	}

	defer func() {
		m.finalizeChanges(o, o.Deletions, deletions)
	}()

	deletions.PreApply, err = processResourceDeletions(o, m.logger, kubectl, deletions.PreApply)
	if err != nil {
		return err
	}

	for _, manifest := range manifests {
		filename := filepath.Join(o.ManifestsDir, manifest.Filename)
		changeSet, err := file.NewChangeSet(filename, manifest.Content)
		if err != nil {
			return err
		}

		m.logChanges(changeSet)

		if !o.AllManifests && !changeSet.HasChanges() {
			continue
		}

		if o.DryRun {
			m.logger.Warnf("Would apply manifest %s", filename)
			m.logger.Debug(string(manifest.Content))
		} else {
			m.logger.Infof("Applying manifest %s", filename)
			if err := kubectl.ApplyManifest(manifest.Content); err != nil {
				return err
			}

			if err := changeSet.Apply(); err != nil {
				return err
			}
		}
	}

	deletions.PostApply, err = processResourceDeletions(o, m.logger, kubectl, deletions.PostApply)

	return err
}

// Destroy deletes all applied manifests from a cluster and tears down the
// cluster infrastructure. Depending on the options it may or may not
// perform a dry run of the destruction process.
func (m *Manager) Destroy(o *Options) error {
	if !o.SkipManifests {
		if err := m.DeleteManifests(o); err != nil {
			return err
		}
	}

	if o.DryRun {
		m.logger.Warn("Would destroy cluster infrastructure")
		return nil
	}

	return m.provisioner.Destroy()
}

// DeleteManifests renders and deletes all manifests from the cluster. It
// also takes care of other resource deletions that should be performed
// after the manifests have been deleted from the cluster.
func (m *Manager) DeleteManifests(o *Options) error {
	values, err := m.readValues(o.Values)
	if err != nil {
		return err
	}

	deletions, err := m.readDeletions(o.Deletions)
	if err != nil {
		return err
	}

	manifests, err := m.renderer.RenderManifests(values)
	if err != nil {
		return err
	}

	creds, err := m.readCredentials(o)
	if err != nil {
		return err
	}

	kubectl := kubernetes.NewKubectl(creds)

	for _, manifest := range manifests {
		filename := filepath.Join(o.ManifestsDir, manifest.Filename)

		if o.DryRun {
			m.logger.Warnf("Would delete manifest %s", filename)
			m.logger.Debug(string(manifest.Content))
		} else {
			m.logger.Infof("Deleting manifest %s", filename)
			if err := kubectl.DeleteManifest(manifest.Content); err != nil {
				return err
			}

			err = os.Remove(filename)
			if err != nil && !os.IsNotExist(err) {
				return errors.WithStack(err)
			}
		}
	}

	deletions.PreDestroy, _ = processResourceDeletions(o, m.logger, kubectl, deletions.PreDestroy)

	return m.finalizeChanges(o, o.Deletions, deletions)
}

func (m *Manager) logChanges(changeSet *file.ChangeSet) {
	filename := changeSet.Filename()
	if changeSet.HasChanges() {
		m.logger.Infof("Changes to %s:\n%s", filename, changeSet.Diff())
	} else {
		m.logger.Infof("No changes to %s", filename)
	}
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

	m.logChanges(changeSet)

	if o.DryRun {
		return nil
	}

	return changeSet.Apply()
}

func (m *Manager) readValues(filename string) (v kcm.Values, err error) {
	if err = file.ReadYAML(filename, &v); err != nil {
		return
	}

	if o, ok := m.provisioner.(provisioner.Outputter); ok {
		values, err := o.Output()
		if err == nil && len(values) > 0 {
			m.logger.Info("Merging values from provisioner")
			v.Merge(values)
		}
	}

	return
}

func (m *Manager) readDeletions(filename string) (d *Deletions, err error) {
	err = file.ReadYAML(filename, &d)

	return
}

func (m *Manager) readCredentials(o *Options) (*credentials.Credentials, error) {
	creds, err := m.credentialSource.GetCredentials()
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

	m.logger.Debugf("Using kubernetes credentials: %#v", c)

	return creds, nil
}
