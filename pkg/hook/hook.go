package hook

import (
	"strings"
	"time"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/resource"
	"github.com/pkg/errors"
)

const (
	// Annotation contains the hook type. This is also what the parser looks to
	// decide whether a resource is a hook or not. The WaitForAnnotation sets
	// an optional condition to wait for with the timeout specified in
	// WaitTimeoutAnnotation. The PolicyAnnotation can contain policies that
	// should be enforced, e.g. deletion of the hook resource after the wait
	// condition was met.
	Annotation            = "kcm/hook"
	WaitForAnnotation     = "kcm/hook-wait-for"
	WaitTimeoutAnnotation = "kcm/hook-wait-timeout"
	PolicyAnnotation      = "kcm/hook-policy"

	// Types of hooks.
	TypePreCreate   = "pre-create"
	TypePreDelete   = "pre-delete"
	TypePreUpgrade  = "pre-upgrade"
	TypePostCreate  = "post-create"
	TypePostDelete  = "post-delete"
	TypePostUpgrade = "post-upgrade"

	// Policies for hooks.
	PolicyDeleteAfterCompletion = "delete-after-completion"
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

	// Policies contains all valid hook policies.
	Policies = []string{
		PolicyDeleteAfterCompletion,
	}
)

// TypePair is a pair of associated hooks that are applied before and after a
// revision upgrade.
type TypePair struct {
	Pre, Post string
}

// Hook is a resource that is applied during revision upgrade.
type Hook struct {
	Resource              *resource.Resource
	Type                  string
	WaitFor               string
	WaitTimeout           time.Duration
	DeleteAfterCompletion bool
}

// New creates a new hook with given resource and annotations. Will return an
// error if the annotations are invalid or the resource does not match any of
// the allowed hook kinds.
func New(r *resource.Resource, annotations map[string]string) (*Hook, error) {
	var err error

	if r.Kind != resource.KindJob {
		return nil, errors.Errorf(`unsupported hook kind %q, currently only %q is supported.`, r.Kind, resource.KindJob)
	}

	typ := annotations[Annotation]
	if !isValidType(typ) {
		return nil, errors.Errorf(`invalid hook type %q, allowed values: %s`, typ, strings.Join(Types, ", "))
	}

	h := &Hook{
		Resource: r,
		Type:     typ,
		WaitFor:  annotations[WaitForAnnotation],
	}

	wt, ok := annotations[WaitTimeoutAnnotation]
	if ok {
		h.WaitTimeout, err = time.ParseDuration(wt)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse annotation %s: %s", WaitTimeoutAnnotation, wt)
		}
	}

	ps, ok := annotations[PolicyAnnotation]
	if ok {
		policies := strings.Split(ps, ",")
		for _, p := range policies {
			if !isValidPolicy(p) {
				return nil, errors.Errorf(`invalid hook policy %q, allowed values: %s`, p, strings.Join(Policies, ", "))
			}

			switch p {
			case PolicyDeleteAfterCompletion:
				if h.WaitFor == "" {
					return nil, errors.Errorf(`policy %q requires to also specify the %s annotation with a valid condition`, p, WaitForAnnotation)
				}

				h.DeleteAfterCompletion = true
			}
		}
	}

	return h, nil
}

// String implements fmt.Stringer
func (h *Hook) String() string {
	var sb strings.Builder

	sb.WriteString(h.Type)
	sb.WriteRune('/')
	sb.WriteString(h.Resource.String())

	if h.WaitFor != "" {
		sb.WriteString(" (wait-for=")
		sb.WriteString(h.WaitFor)

		if h.WaitTimeout > 0 {
			sb.WriteString(",wait-timeout=")
			sb.WriteString(h.WaitTimeout.String())
		}

		sb.WriteRune(')')
	}

	return sb.String()
}

func isValidType(typ string) bool {
	for _, t := range Types {
		if t == typ {
			return true
		}
	}

	return false
}

func isValidPolicy(policy string) bool {
	for _, p := range Policies {
		if p == policy {
			return true
		}
	}

	return false
}
