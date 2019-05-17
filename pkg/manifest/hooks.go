package manifest

import (
	"strings"

	"github.com/pkg/errors"
)

const (
	HooksAnnotation      = "kcm/hooks"
	HookPolicyAnnotation = "kcm/hook-policy"

	HookTypePreApply   HookType = "pre-apply"
	HookTypePreDelete  HookType = "pre-delete"
	HookTypePostApply  HookType = "post-apply"
	HookTypePostDelete HookType = "post-delete"
)

type HookPolicy string

type HookType string

type HookSlice []*Hook

type HookSliceMap map[HookType]HookSlice

type Hook struct {
	*Resource

	types  []HookType
	policy HookPolicy
}

func newHook(r *Resource, head resourceHead) (*Hook, error) {
	if head.Kind != "Job" {
		return nil, errors.Errorf(`Unsupported hook kind %q. Currently only "Job" is supported.`, head.Kind)
	}

	h := &Hook{
		Resource: r,
		types:    make([]HookType, 0),
	}

	p, ok := head.Metadata.Annotations[HookPolicyAnnotation]
	if ok {
		h.policy = HookPolicy(p)
	}

	hooks := head.Metadata.Annotations[HooksAnnotation]

	parts := strings.Split(hooks, ",")
	for _, p := range parts {
		hookType := HookType(strings.TrimSpace(p))
		h.types = append(h.types, hookType)
	}

	return h, nil
}

func (s HookSlice) Bytes() []byte {
	var buf resourceBuffer

	for _, h := range s {
		buf.Write(h.Content)
	}

	return buf.Bytes()
}

func (m HookSliceMap) Get(typ HookType) HookSlice {
	return m[typ]
}

func (m HookSliceMap) Has(typ HookType) bool {
	hooks, ok := m[typ]

	return ok && len(hooks) > 0
}
