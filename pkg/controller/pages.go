package controller

import (
	"mime"
	"net/http"
	"path"
	"strings"

	"github.com/ismtabo/mapon-viewer/pkg/service"
	"github.com/ismtabo/mapon-viewer/pkg/template"
	"github.com/rs/zerolog/log"
)

// PagesController impelements HTTP controller for web pages requests.
type PagesController interface {
	StaticFiles(w http.ResponseWriter, r *http.Request)
	IndexPage(w http.ResponseWriter, r *http.Request)
	LoginPage(w http.ResponseWriter, r *http.Request)
}

type pagesController struct {
	ssnsSvc  service.SessionsService
	tmplMngr template.Manager
}

// NewPagesController creates a new instance of PagesController.
func NewPagesController(ssnsSvc service.SessionsService, tmplMngr template.Manager) PagesController {
	return &pagesController{ssnsSvc: ssnsSvc, tmplMngr: tmplMngr}
}

// StaticFiles returns statics files of web application.
func (c *pagesController) StaticFiles(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Path
	filePath = strings.Replace(filePath, "/static/", "", 1)
	file, err := c.tmplMngr.StaticFiles(filePath)
	if err != nil {
		RenderError(r.Context(), w, err)
		return
	}
	log.Info().Msgf("Serving static file %s", filePath)
	mimeType := mime.TypeByExtension(path.Ext(filePath))
	w.Header().Add("Content-Type", mimeType)
	w.Write(file)
}

// IndexPage returns HTML index page.
func (c *pagesController) IndexPage(w http.ResponseWriter, r *http.Request) {
	if !c.ssnsSvc.IsAuthenticated(w, r) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	body, err := c.tmplMngr.RenderFile("index.html", nil)
	if err != nil {
		RenderError(r.Context(), w, err)
		return
	}
	w.Write(body)
}

// LoginPage returns HTML login page.
func (c *pagesController) LoginPage(w http.ResponseWriter, r *http.Request) {
	if c.ssnsSvc.IsAuthenticated(w, r) {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	body, err := c.tmplMngr.RenderFile("login.html", nil)
	if err != nil {
		RenderError(r.Context(), w, err)
		return
	}
	w.Write(body)
}
