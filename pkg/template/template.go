package template

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"path"
)

// TemplateManager implements a html template management.
type TemplateManager interface {
	StaticFiles(file string) ([]byte, error)
	RenderFile(file string, data interface{}) ([]byte, error)
}

type templateManager struct{}

// NewTemplateManager creates a new instance of a TemplateManager.
func NewTemplateManager() TemplateManager {
	return &templateManager{}
}

// StaticFiles returns index html.
func (m templateManager) StaticFiles(file string) ([]byte, error) {
	return ioutil.ReadFile(path.Join("client", "assets", file))
}

// RenderBlockedAccessPage returns index html.
func (m templateManager) RenderFile(file string, data interface{}) ([]byte, error) {
	template, err := template.ParseFiles(path.Join("client", "pages", file))
	if err != nil {
		return nil, err
	}
	var buffer bytes.Buffer
	err = template.Execute(&buffer, data)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
