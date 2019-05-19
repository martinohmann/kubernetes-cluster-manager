package hook

import (
	"strings"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/resource"
	"github.com/pkg/errors"
)

const (
	HooksAnnotation  = "kcm/hooks"
	PolicyAnnotation = "kcm/hook-policy"

	TypePreApply   Type = "pre-apply"
	TypePreDelete  Type = "pre-delete"
	TypePostApply  Type = "post-apply"
	TypePostDelete Type = "post-delete"
)

type Policy string

type Type string

type Hook struct {
	*resource.Resource

	Types  []Type
	policy Policy
}

func New(r *resource.Resource, annotations map[string]string) (*Hook, error) {
	if r.Kind != "Job" {
		return nil, errors.Errorf(`Unsupported hook kind %q. Currently only "Job" is supported.`, r.Kind)
	}

	h := &Hook{
		Resource: r,
		Types:    make([]Type, 0),
	}

	p, ok := annotations[PolicyAnnotation]
	if ok {
		h.policy = Policy(p)
	}

	hooks := annotations[HooksAnnotation]

	parts := strings.Split(hooks, ",")
	for _, p := range parts {
		hookType := Type(strings.TrimSpace(p))
		h.Types = append(h.Types, hookType)
	}

	return h, nil
}
