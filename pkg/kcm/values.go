package kcm

import (
	"github.com/imdario/mergo"
)

// Values contains the output values of an infrastructure manager.
type Values map[string]interface{}

// Merge deep merges other on top of v and overrides values already present in
// v.
func (v Values) Merge(other Values) error {
	return mergo.Merge(&v, other, mergo.WithOverride)
}
