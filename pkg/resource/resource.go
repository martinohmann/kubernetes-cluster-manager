package resource

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

const (
	// DeletionPolicyAnnotation can be set on a resources to control the
	// behaviour of the resource deletion operation.
	DeletionPolicyAnnotation = "kcm/deletion-policy"

	// DeletePersistentVolumeClaimsPolicy will cause the deletion of all PVCs
	// that were created for a StatefulSet.
	DeletePersistentVolumeClaimsPolicy = "delete-pvcs"
)

const (
	// Kinds of Kubernetes resources that are treated in a special way by kcm.
	Job                   = "Job"
	PersistentVolumeClaim = "PersistentVolumeClaim"
	StatefulSet           = "StatefulSet"
)

// Resource is a kubernetes resource.
type Resource struct {
	Kind      string
	Name      string
	Namespace string
	Content   []byte

	DeletePersistentVolumeClaims bool

	contentHint []byte
	hint        Hint
}

// Head defines the yaml structure of a resource head. This is used
// for parsing metadata from raw yaml documents.
type Head struct {
	Kind     string   `yaml:"kind"`
	Metadata Metadata `yaml:"metadata"`
}

// String implements fmt.Stringer
func (h Head) String() string {
	if h.Metadata.Namespace == "" {
		return fmt.Sprintf("%s/%s", strings.ToLower(h.Kind), h.Metadata.Name)
	}

	return fmt.Sprintf("%s/%s/%s", h.Metadata.Namespace, strings.ToLower(h.Kind), h.Metadata.Name)
}

// Metadata is the resource metadata we are interested in.
type Metadata struct {
	Name        string            `yaml:"name"`
	Namespace   string            `yaml:"namespace"`
	Annotations map[string]string `yaml:"annotations"`
}

// New creates a new resource value with content and head.
func New(content []byte, head Head) (*Resource, error) {
	r := &Resource{
		Kind:      head.Kind,
		Name:      head.Metadata.Name,
		Namespace: head.Metadata.Namespace,
		Content:   content,
	}

	policy, ok := head.Metadata.Annotations[DeletionPolicyAnnotation]
	if ok {
		if policy != DeletePersistentVolumeClaimsPolicy {
			return nil, errors.Errorf("unsupported deletion policy %q", policy)
		}

		if r.Kind != StatefulSet {
			return nil, errors.Errorf("deletion policy %q can only be applied to StatefulSets, got %s", policy, r.Kind)
		}

		r.DeletePersistentVolumeClaims = true
	}

	return r, nil
}

// String implements fmt.Stringer
func (r *Resource) String() string {
	if r.Namespace == "" {
		return fmt.Sprintf("%s/%s", strings.ToLower(r.Kind), r.Name)
	}

	return fmt.Sprintf("%s/%s/%s", r.Namespace, strings.ToLower(r.Kind), r.Name)
}

// matches returns true if other matches r. Two resources match if their name,
// kind and namespace are the same.
func (r *Resource) matches(other *Resource) bool {
	if r == other {
		return true
	}

	if r == nil || other == nil {
		return false
	}

	if r.Kind != other.Kind || r.Namespace != other.Namespace {
		return false
	}

	return r.Name == other.Name
}
