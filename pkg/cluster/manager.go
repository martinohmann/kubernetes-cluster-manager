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

// Manager is a kcm.Manager.
type Manager struct {
	credentialSource kcm.CredentialSource
	provisioner      kcm.Provisioner
	renderer         kcm.Renderer
	values           kcm.Values
	deletions        *kcm.Deletions
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
		deletions:        &kcm.Deletions{},
		values:           kcm.Values{},
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
	if err := m.prepare(o); err != nil {
		return err
	}

	creds, err := m.credentialSource.GetCredentials()
	if err != nil {
		return err
	}

	valueBytes, err := yaml.Marshal(m.values)
	if err != nil {
		return err
	}

	err = m.finalizeChanges(o, o.Values, valueBytes)
	if err != nil {
		return err
	}

	manifests, err := m.renderer.RenderManifests(m.values)
	if err != nil {
		return err
	}

	kubectl := kubernetes.NewKubectl(creds)

	if !o.DryRun {
		if err := os.MkdirAll(o.ManifestsDir, 0775); err != nil {
			return errors.WithStack(err)
		}

		m.logger.Info("Waiting for cluster to become available...")

		if err := kubectl.WaitForCluster(); err != nil {
			return err
		}
	}

	defer m.finalizeDeletions(o, m.deletions)

	err = processResourceDeletions(o, m.logger, kubectl, m.deletions.PreApply)
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

		if o.OnlyChanges && !changeSet.HasChanges() {
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

	return processResourceDeletions(o, m.logger, kubectl, m.deletions.PostApply)
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
	if err := m.prepare(o); err != nil {
		return err
	}

	creds, err := m.credentialSource.GetCredentials()
	if err != nil {
		return err
	}

	manifests, err := m.renderer.RenderManifests(m.values)
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

	defer m.finalizeDeletions(o, m.deletions)

	return processResourceDeletions(o, m.logger, kubectl, m.deletions.PreDestroy)
}

func (m *Manager) logChanges(changeSet *file.ChangeSet) {
	filename := changeSet.Filename
	if changeSet.HasChanges() {
		m.logger.Infof("Changes to %s:\n%s", filename, changeSet.Diff())
	} else {
		m.logger.Infof("No changes to %s", filename)
	}
}

func (m *Manager) finalizeChanges(o *kcm.Options, filename string, content []byte) error {
	cs, err := file.NewChangeSet(filename, content)
	if err != nil {
		return err
	}

	m.logChanges(cs)

	if o.DryRun {
		return nil
	}

	return cs.Apply()
}

func (m *Manager) finalizeDeletions(o *kcm.Options, deletions *kcm.Deletions) error {
	buf, err := yaml.Marshal(deletions.FilterPending())
	if err != nil {
		return err
	}

	return m.finalizeChanges(o, o.Deletions, buf)
}

func (m *Manager) prepare(o *kcm.Options) error {
	if err := file.ReadYAML(o.Values, &m.values); err != nil {
		return err
	}

	if err := file.ReadYAML(o.Deletions, &m.deletions); err != nil {
		return err
	}

	v, err := m.provisioner.Fetch()
	if err != nil {
		return err
	}

	return m.values.Merge(v)
}
