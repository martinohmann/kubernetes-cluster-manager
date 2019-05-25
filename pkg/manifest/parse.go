package manifest

import (
	"bytes"
	"io"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/hook"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/resource"
	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

func Parse(buf []byte) (resource.Slice, hook.SliceMap, error) {
	resources := make(resource.Slice, 0)
	hooks := make(hook.SliceMap)

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
			log.Debugf("error while parsing resources, skipping resource: %s", err.Error())
			continue
		}

		buf, err := yaml.Marshal(v)
		if err != nil {
			return nil, nil, err
		}

		var head resource.Head

		err = yaml.Unmarshal(buf, &head)
		if err != nil {
			return nil, nil, err
		}

		if head.Kind == "" || head.Metadata.Name == "" {
			continue
		}

		r := resource.New(buf, head)

		if _, ok := head.Metadata.Annotations[hook.Annotation]; ok {
			h, err := hook.New(r, head.Metadata.Annotations)
			if err != nil {
				return nil, nil, err
			}

			if _, ok := hooks[h.Type]; ok {
				hooks[h.Type] = append(hooks[h.Type], h)
			} else {
				hooks[h.Type] = hook.Slice{h}
			}

			continue
		}

		resources = append(resources, r)
	}

	resources.Sort(resource.ApplyOrder)
	hooks.Sort()

	return resources, hooks, nil
}
