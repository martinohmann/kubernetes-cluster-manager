package cluster

import (
	"os"
	"path/filepath"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/file"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kcm"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kubernetes"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

const (
	dirMode os.FileMode = 0775
)

var (
	emptyCredentials = kcm.Credentials{}
)

// Manager is a kcm.Manager.
type Manager struct {
	credentialSource kcm.CredentialSource
	provisioner      kcm.Provisioner
	renderer         kcm.Renderer
	logger           *log.Logger
}

// NewManager creates a new cluster manager.
func NewManager(
	credentialSource kcm.CredentialSource,
	provisioner kcm.Provisioner,
	renderer kcm.Renderer,
	logger *log.Logger,
) *Manager {
	return &Manager{
		credentialSource: credentialSource,
		provisioner:      provisioner,
		renderer:         renderer,
		logger:           logger,
	}
}

// Provision implements Provision from the kcm.ClusterManager interface.
func (m *Manager) Provision(o *kcm.Options) error {
	var err error

	if o.DryRun {
		err = m.provisioner.Reconcile()
	} else {
		err = m.provisioner.Provision()
	}

	if err != nil || o.SkipManifests {
		return err
	}

	return m.ApplyManifests(o)
}

// ApplyManifests implements ApplyManifests from the kcm.ClusterManager interface.
func (m *Manager) ApplyManifests(o *kcm.Options) error {
	values, err := m.readValues(o.Values)
	if err != nil {
		return err
	}

	deletions, err := m.readDeletions(o.Deletions)
	if err != nil {
		return err
	}

	creds, err := m.readCredentials()
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
		m.finalizeChanges(o, o.Deletions, deletions.FilterPending())
	}()

	err = processResourceDeletions(o, m.logger, kubectl, deletions.PreApply)
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
			if err := kubectl.ApplyManifest(manifest); err != nil {
				return err
			}

			if err := changeSet.Apply(); err != nil {
				return err
			}
		}
	}

	return processResourceDeletions(o, m.logger, kubectl, deletions.PostApply)
}

// Destroy implements Destroy from the kcm.ClusterManager interface.
func (m *Manager) Destroy(o *kcm.Options) error {
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

// DeleteManifests implements DeleteManifests from the kcm.ClusterManager interface.
func (m *Manager) DeleteManifests(o *kcm.Options) error {
	values, err := m.readValues(o.Values)
	if err != nil {
		return err
	}

	deletions, err := m.readDeletions(o.Deletions)
	if err != nil {
		return err
	}

	creds, err := m.readCredentials()
	if err != nil {
		return err
	}

	manifests, err := m.renderer.RenderManifests(values)
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
			if err := kubectl.DeleteManifest(manifest); err != nil {
				return err
			}

			err = os.Remove(filename)
			if err != nil && !os.IsNotExist(err) {
				return errors.WithStack(err)
			}
		}
	}

	processResourceDeletions(o, m.logger, kubectl, deletions.PreDestroy)

	return m.finalizeChanges(o, o.Deletions, deletions.FilterPending())
}

func (m *Manager) logChanges(changeSet *file.ChangeSet) {
	filename := changeSet.Filename()
	if changeSet.HasChanges() {
		m.logger.Infof("Changes to %s:\n%s", filename, changeSet.Diff())
	} else {
		m.logger.Infof("No changes to %s", filename)
	}
}

func (m *Manager) finalizeChanges(o *kcm.Options, filename string, v interface{}) error {
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

	additional, err := m.provisioner.Fetch()
	if err == nil {
		v.Merge(additional)
	}

	return
}

func (m *Manager) readDeletions(filename string) (d *kcm.Deletions, err error) {
	err = file.ReadYAML(filename, &d)

	return
}

func (m *Manager) readCredentials() (*kcm.Credentials, error) {
	creds, err := m.credentialSource.GetCredentials()
	if err != nil {
		return nil, err
	}

	if *creds == emptyCredentials {
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
