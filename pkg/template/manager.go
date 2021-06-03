package template

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"os"
	"path"

	"github.com/ismtabo/mapon-viewer/pkg/errors"
)

// Manager implements a html template management.
type Manager interface {
	StaticFiles(file string) ([]byte, error)
	RenderFile(file string, data interface{}) ([]byte, error)
}

type templateManager struct{}

// NewTemplateManager creates a new instance of a Manager.
func NewTemplateManager() Manager {
	return &templateManager{}
}

// StaticFiles returns index html.
func (m templateManager) StaticFiles(file string) ([]byte, error) {
	bytes, err := ioutil.ReadFile(path.Join("client", "assets", file))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.NewNotFoundError()
		}
		return nil, errors.NewInternalServerError(err)
	}
	return bytes, nil
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
