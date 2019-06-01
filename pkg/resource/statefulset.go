package resource

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

// partialStatefulSet only contains the neccessary fields for parsing PVCs from
// a StatefulSet resource.
type partialStatefulSet struct {
	Kind     string   `yaml:"kind"`
	Metadata Metadata `yaml:"metadata"`
	Spec     struct {
		Replicas             int                           `yaml:"replicas"`
		VolumeClaimTemplates []*partialVolumeClaimTemplate `yaml:"volumeClaimTemplates"`
	} `yaml:"spec"`
}

// partialVolumeClaimTemplate only contains the metadata of a volume claim template.
type partialVolumeClaimTemplate struct {
	Metadata Metadata `yaml:"metadata"`
}

// buildPersistentVolumeClaims build a PersistentVolumeClaim Head for each
// volume claim template it finds in statefuleSet, taking the number of
// replicas into account.
func buildPersistentVolumeClaims(statefulSet *partialStatefulSet) []Head {
	replicas := statefulSet.Spec.Replicas
	if replicas < 1 {
		replicas = 1
	}

	volumeClaimTemplates := statefulSet.Spec.VolumeClaimTemplates

	claims := make([]Head, 0, len(volumeClaimTemplates)*replicas)

	for _, vct := range volumeClaimTemplates {
		for i := 0; i < replicas; i++ {
			claims = append(claims, Head{
				Kind: KindPersistentVolumeClaim,
				Metadata: Metadata{
					Name:      buildPersistentVolumeClaimName(statefulSet, vct, i),
					Namespace: statefulSet.Metadata.Namespace,
				},
			})
		}
	}

	return claims
}

// buildPersistentVolumeClaimName assembles the PVC name generated by a volume
// claim template. The name has the format {stateful-set-name}-{volume-claim-name}-{ordinal}.
func buildPersistentVolumeClaimName(statefulSet *partialStatefulSet, vct *partialVolumeClaimTemplate, ordinal int) string {
	return fmt.Sprintf("%s-%s-%d", statefulSet.Metadata.Name, vct.Metadata.Name, ordinal)
}

// PersistentVolumeClaimsForDeletion extracts all PersistentVolumeClaims for
// StatefulSets that have been annotated with:
//
//   kcm/deletion-policy: delete-pvcs
//
// It returns a slice of Head containing name and namespace of every PVC that
// should be deleted after the StatefulSet is deleted. This is a workaround
// until a similar feature (https://github.com/kubernetes/kubernetes/issues/55045)
// is implemented in Kubernetes itself.
func (s Slice) PersistentVolumeClaimsForDeletion() []Head {
	claims := make([]Head, 0)

	for _, r := range s {
		if r.Kind != KindStatefulSet || !r.DeletePersistentVolumeClaims {
			continue
		}

		var statefulSet partialStatefulSet

		err := yaml.Unmarshal(r.Content, &statefulSet)
		if err != nil {
			log.Errorf("error while parsing stateful set %q: %s", r.Name, err.Error())
			continue
		}

		claims = append(claims, buildPersistentVolumeClaims(&statefulSet)...)
	}

	return claims
}
