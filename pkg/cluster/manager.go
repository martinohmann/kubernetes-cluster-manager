package cluster

import (
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/file"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kcm"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kubernetes"
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
func (p *Manager) Provision(o *kcm.Options) error {
	var err error

	if o.DryRun {
		err = p.provisioner.Reconcile()
	} else {
		err = p.provisioner.Provision()
	}

	if err != nil || o.SkipManifests {
		return err
	}

	return p.ApplyManifests(o)
}

// ApplyManifests implements ApplyManifests from the kcm.ClusterManager interface.
func (p *Manager) ApplyManifests(o *kcm.Options) error {
	if err := p.prepareManifests(o); err != nil {
		return err
	}

	creds, err := p.credentialSource.GetCredentials()
	if err != nil {
		return err
	}

	valueBytes, err := yaml.Marshal(p.values)
	if err != nil {
		return err
	}

	err = p.finalizeChanges(o, o.Values, valueBytes)
	if err != nil {
		return err
	}

	manifest, err := p.renderer.RenderManifest(p.values)
	if err != nil {
		return err
	}

	kubectl := kubernetes.NewKubectl(creds)

	if !o.DryRun {
		p.logger.Info("Waiting for cluster to become available...")

		if err := kubectl.WaitForCluster(); err != nil {
			return err
		}
	}

	err = p.finalizeChanges(o, o.Manifest, manifest)
	if err != nil {
		return err
	}

	defer p.finalizeDeletions(o, p.deletions)

	if err := processResourceDeletions(o, p.logger, kubectl, p.deletions.PreApply); err != nil {
		return err
	}

	if o.DryRun {
		p.logger.Warn("Would apply manifest")
		p.logger.Debug(string(manifest))
	} else if err := kubectl.ApplyManifest(manifest); err != nil {
		return err
	}

	return processResourceDeletions(o, p.logger, kubectl, p.deletions.PostApply)
}

// Destroy implements Destroy from the kcm.ClusterManager interface.
func (p *Manager) Destroy(o *kcm.Options) error {
	if !o.SkipManifests {
		if err := p.DeleteManifests(o); err != nil {
			return err
		}
	}

	if o.DryRun {
		p.logger.Warn("Would destroy infrastructure")
		return nil
	}

	return p.provisioner.Destroy()

}

// DeleteManifests implements DeleteManifests from the kcm.ClusterManager interface.
func (p *Manager) DeleteManifests(o *kcm.Options) error {
	if err := p.prepareManifests(o); err != nil {
		return err
	}

	creds, err := p.credentialSource.GetCredentials()
	if err != nil {
		return err
	}

	manifest, err := p.renderer.RenderManifest(p.values)
	if err != nil {
		return err
	}

	kubectl := kubernetes.NewKubectl(creds)

	if o.DryRun {
		p.logger.Warn("Would delete manifest")
		p.logger.Debug(string(manifest))
	} else if err := kubectl.DeleteManifest(manifest); err != nil {
		return err
	}

	defer p.finalizeDeletions(o, p.deletions)

	return processResourceDeletions(o, p.logger, kubectl, p.deletions.PreDestroy)
}

func (p *Manager) finalizeChanges(o *kcm.Options, filename string, content []byte) error {
	changes := file.NewChanges(filename, content)

	diff, err := changes.Diff()
	if err != nil {
		return err
	}

	if len(diff) > 0 {
		p.logger.Infof("Changes to %s:\n%s", filename, diff)
	} else {
		p.logger.Infof("No changes to %s", filename)
	}

	if o.DryRun {
		return nil
	}

	return changes.Apply()
}

func (p *Manager) finalizeDeletions(o *kcm.Options, deletions *kcm.Deletions) error {
	buf, err := yaml.Marshal(deletions.FilterPending())
	if err != nil {
		return err
	}

	return p.finalizeChanges(o, o.Deletions, buf)
}

func (p *Manager) prepareManifests(o *kcm.Options) error {
	if err := file.LoadYAML(o.Values, &p.values); err != nil {
		return err
	}

	if err := file.LoadYAML(o.Deletions, &p.deletions); err != nil {
		return err
	}

	v, err := p.provisioner.Fetch()
	if err != nil {
		return err
	}

	return p.values.Merge(v)
}
