package hook

import (
	"strings"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/resource"
	"github.com/pkg/errors"
)

const (
	// Annotation contains the hook type. This is also what the parser looks to
	// decide whether a resource is a hook or not.
	Annotation = "kcm/hook"

	// PolicyAnnotation specifies the hook policy. This is currently not used.
	PolicyAnnotation = "kcm/hook-policy"

	// Types of hooks.
	TypePreCreate   = "pre-create"
	TypePreDelete   = "pre-delete"
	TypePreUpgrade  = "pre-upgrade"
	TypePostCreate  = "post-create"
	TypePostDelete  = "post-delete"
	TypePostUpgrade = "post-upgrade"
)

var (
	// Types contains all valid hook types.
	Types = []string{
		TypePreCreate,
		TypePreDelete,
		TypePreUpgrade,
		TypePostCreate,
		TypePostDelete,
		TypePostUpgrade,
	}

	// Pairs of associated hook types.
	TypeCreate  = TypePair{TypePreCreate, TypePostCreate}
	TypeDelete  = TypePair{TypePreDelete, TypePostDelete}
	TypeUpgrade = TypePair{TypePreUpgrade, TypePostUpgrade}
)

// TypePair is a pair of associated hooks that are applied before and after a
// revision upgrade.
type TypePair struct {
	Pre, Post string
}

// Hook is a resource that is applied during revision upgrade.
type Hook struct {
	Resource *resource.Resource
	Type     string
	Policy   string
}

// New creates a new hook with given resource and annotations. Will return an
// error if the annotations are invalid or the resource does not match any of
// the allowed hook kinds.
func New(r *resource.Resource, annotations map[string]string) (*Hook, error) {
	if r.Kind != "Job" {
		return nil, errors.Errorf(`Unsupported hook kind %q. Currently only "Job" is supported.`, r.Kind)
	}

	typ := annotations[Annotation]
	if !isValidType(typ) {
		return nil, errors.Errorf(`Invalid hook type %q. Allowed values: %s`, typ, strings.Join(Types, ", "))
	}

	h := &Hook{
		Resource: r,
		Type:     typ,
		Policy:   annotations[PolicyAnnotation],
	}

	return h, nil
}

func isValidType(typ string) bool {
	for _, t := range Types {
		if t == typ {
			return true
		}
	}

	return false
}
