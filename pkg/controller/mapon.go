package controller

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ismtabo/mapon-viewer/pkg/controller/dto"
	"github.com/ismtabo/mapon-viewer/pkg/errors"
	"github.com/ismtabo/mapon-viewer/pkg/repository"
	"github.com/jinzhu/copier"
)

type MaponController interface {
	GetMaponInfo(rw http.ResponseWriter, r *http.Request)
}

type maponController struct {
	repo repository.MaponRepository
}

func NewMaponController(repo repository.MaponRepository) MaponController {
	return &maponController{repo: repo}
}

func (c maponController) GetMaponInfo(rw http.ResponseWriter, r *http.Request) {
	f := r.URL.Query().Get("from")
	if f == "" {
		RenderError(r.Context(), rw, errors.NewBadRequestError("missing 'from' query param"))
		return
	}
	from, err := time.Parse(time.RFC3339, f)
	if err != nil {
		RenderError(r.Context(), rw, errors.NewBadRequestError("invalid 'from' query param"))
		return
	}
	t := r.URL.Query().Get("till")
	if t == "" {
		RenderError(r.Context(), rw, errors.NewBadRequestError("missing 'till' query param"))
		return
	}
	till, err := time.Parse(time.RFC3339, t)
	if err != nil {
		RenderError(r.Context(), rw, errors.NewBadRequestError("invalid 'till' query param"))
		return
	}
	info, err := c.repo.GetInfo(r.Context(), from, till)
	if err != nil {
		RenderError(r.Context(), rw, err)
		return
	}
	infoDTO := &dto.MaponInfo{}
	if err := copier.Copy(infoDTO, info); err != nil {
		RenderError(r.Context(), rw, err)
		return
	}
	body, err := json.Marshal(infoDTO)
	if err != nil {
		RenderError(r.Context(), rw, err)
		return
	}
	rw.Header().Add("Content-Type", "application/json")
	rw.Write(body)
}
