package renderer

import (
	"bytes"
	"io"
	"io/ioutil"
	"path/filepath"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/kcm"
	"github.com/pkg/errors"
)

// GoTemplate uses the text/template package to render manifests
type GoTemplate struct {
	TemplatesDir string
}

// NewGoTemplate creates a new go-template renderer.
func NewGoTemplate(o *Options) Renderer {
	return &GoTemplate{
		TemplatesDir: o.TemplatesDir,
	}
}

// RenderManifests implements Renderer.
func (r *GoTemplate) RenderManifests(v kcm.Values) ([]*Manifest, error) {
	return renderManifests(r.TemplatesDir, v, renderDirectory)
}

// renderDirectory renders all templates in a directory. It satisfies the
// signature of renderManifestFunc.
func renderDirectory(dir string, v kcm.Values) (*Manifest, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	manifestName := filepath.Base(dir)

	var buf bytes.Buffer

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		ext := filepath.Ext(f.Name())
		if ext != ".yaml" && ext != ".yml" {
			continue
		}

		source := filepath.Join(manifestName, f.Name())
		templatePath := filepath.Join(dir, f.Name())

		writeSourceHeader(&buf, source)

		err := renderTemplate(templatePath, &buf, v)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		buf.WriteString("\n")
	}

	m := &Manifest{
		Name:    manifestName,
		Content: buf.Bytes(),
	}

	return m, nil
}

// renderTemplate renders a single template file into w.
func renderTemplate(path string, w io.Writer, data interface{}) error {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	t, err := template.New(filepath.Base(path)).
		Funcs(sprig.TxtFuncMap()).
		Parse(string(buf))

	if err != nil {
		return err
	}

	return t.Execute(w, data)
}
