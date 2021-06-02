package controller

import (
	"mime"
	"net/http"
	"path"
	"strings"

	"github.com/ismtabo/mapon-viewer/pkg/template"
	"github.com/kataras/go-sessions/v3"
	"github.com/rs/zerolog/log"
)

type PagesController interface {
	StaticFiles(w http.ResponseWriter, r *http.Request)
	IndexPage(w http.ResponseWriter, r *http.Request)
	LoginPage(w http.ResponseWriter, r *http.Request)
}

type pagesController struct {
	sessions *sessions.Sessions
	template template.TemplateManager
}

func NewPagesController(sessions *sessions.Sessions, template template.TemplateManager) PagesController {
	return &pagesController{sessions: sessions, template: template}
}

func (c *pagesController) StaticFiles(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Path
	filePath = strings.Replace(filePath, "/static/", "", 1)
	file, err := c.template.StaticFiles(filePath)
	if err != nil {
		RenderError(r.Context(), w, err)
		return
	}
	log.Info().Msgf("Serving static file %s", filePath)
	mimeType := mime.TypeByExtension(path.Ext(filePath))
	w.Header().Add("Content-Type", mimeType)
	w.Write(file)
}

func (c *pagesController) IndexPage(w http.ResponseWriter, r *http.Request) {
	if auth, _ := c.sessions.Start(w, r).GetBoolean(authAttribute); !auth {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	body, err := c.template.RenderFile("index.html", nil)
	if err != nil {
		RenderError(r.Context(), w, err)
		return
	}
	w.Write(body)
}

func (c *pagesController) LoginPage(w http.ResponseWriter, r *http.Request) {
	if auth, _ := c.sessions.Start(w, r).GetBoolean(authAttribute); auth {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	body, err := c.template.RenderFile("login.html", nil)
	if err != nil {
		RenderError(r.Context(), w, err)
		return
	}
	w.Write(body)
}
