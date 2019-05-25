package hook

import (
	"strings"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/resource"
	"github.com/pkg/errors"
)

const (
	Annotation       = "kcm/hook"
	PolicyAnnotation = "kcm/hook-policy"

	TypePreApply    = "pre-apply"
	TypePreCreate   = "pre-create"
	TypePreDelete   = "pre-delete"
	TypePreUpgrade  = "pre-upgrade"
	TypePostApply   = "post-apply"
	TypePostCreate  = "post-create"
	TypePostDelete  = "post-delete"
	TypePostUpgrade = "post-upgrade"
)

var (
	Types = []string{
		TypePreApply,
		TypePreCreate,
		TypePreDelete,
		TypePreUpgrade,
		TypePostApply,
		TypePostCreate,
		TypePostDelete,
		TypePostUpgrade,
	}

	TypeApply   = TypePair{TypePreApply, TypePostApply}
	TypeCreate  = TypePair{TypePreCreate, TypePostCreate}
	TypeDelete  = TypePair{TypePreDelete, TypePostDelete}
	TypeUpgrade = TypePair{TypePreUpgrade, TypePostUpgrade}
)

type TypePair struct {
	Pre, Post string
}

type Hook struct {
	Resource *resource.Resource
	Type     string
	Policy   string
}

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
