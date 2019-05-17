package manifest

import (
	"bytes"
	"io"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kubernetes"
	"gopkg.in/yaml.v2"
)

type Resource struct {
	Kind        string
	Name        string
	Namespace   string
	Annotations map[string]string
	Content     []byte
}

func newResource(content []byte, head resourceHead) *Resource {
	return &Resource{
		Kind:        head.Kind,
		Name:        head.Metadata.Name,
		Namespace:   head.Metadata.Namespace,
		Annotations: head.Metadata.Annotations,
		Content:     content,
	}
}

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

func (r *Resource) Selector() kubernetes.ResourceSelector {
	return kubernetes.ResourceSelector{
		Name:      r.Name,
		Namespace: r.Namespace,
		Kind:      r.Kind,
	}
}

type ResourceSlice []*Resource

func (s ResourceSlice) Bytes() []byte {
	var buf resourceBuffer

	for _, r := range s {
		buf.Write(r.Content)
	}

	return buf.Bytes()
}

func (s ResourceSlice) Selectors() []kubernetes.ResourceSelector {
	rs := make([]kubernetes.ResourceSelector, 0)

	for _, res := range s {
		rs = append(rs, res.Selector())
	}

	return rs
}

func (s ResourceSlice) Sort(order ResourceOrder) ResourceSlice {
	return sortResources(s, order)
}

type resourceHead struct {
	Kind     string `yaml:"kind"`
	Metadata struct {
		Name        string            `yaml:"name"`
		Namespace   string            `yaml:"namespace"`
		Annotations map[string]string `yaml:"annotations"`
	} `yaml:"metadata"`
}

func parseResources(buf []byte) (ResourceSlice, HookSliceMap, error) {
	resources := make(ResourceSlice, 0)
	hooks := make(HookSliceMap)

	r := bytes.NewBuffer(buf)
	d := yaml.NewDecoder(r)

	for {
		// since the yaml parser does not support unmarshaling into raw bytes
		// (e.g. like json.RawMessage) we have to unmarshal into a map first,
		// then marshal to get the content of the single manifest and then
		// unmarshal again to get the manifest metadata.
		var v map[string]interface{}

		err := d.Decode(&v)
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, nil, err
		}

		buf, err := yaml.Marshal(v)
		if err != nil {
			return nil, nil, err
		}

		var head resourceHead

		err = yaml.Unmarshal(buf, &head)
		if err != nil {
			return nil, nil, err
		}

		if head.Kind == "" || head.Metadata.Name == "" {
			continue
		}

		resource := newResource(buf, head)

		if _, ok := head.Metadata.Annotations[HooksAnnotation]; ok {
			h, err := newHook(resource, head)
			if err != nil {
				return nil, nil, err
			}

			for _, t := range h.types {
				if hooks[t] == nil {
					hooks[t] = make(HookSlice, 0)
				}

				hooks[t] = append(hooks[t], h)
			}
			continue
		}

		resources = append(resources, resource)
	}

	return resources, hooks, nil
}
