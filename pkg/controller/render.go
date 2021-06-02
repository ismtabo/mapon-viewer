package controller

import (
	"context"
	"net/http"

	"github.com/ismtabo/mapon-viewer/pkg/ctxt"
	"github.com/ismtabo/mapon-viewer/pkg/errors"
	"github.com/rs/zerolog"
)

// RenderError log the error and it writes the corresponding HTTP response.
func RenderError(ctx context.Context, rw http.ResponseWriter, err error) {
	appErr := assertApplicationError(err)
	logError(ctx, appErr)
	writeError(ctx, rw, appErr)
}

func assertApplicationError(err error) *errors.Error {
	appErr, ok := err.(*errors.Error)
	if !ok {
		return errors.NewInternalServerError(err)
	}
	return appErr
}

func logError(ctx context.Context, err *errors.Error) {
	log := ctxt.GetLogger(ctx)
	if log == nil {
		return
	}
	var event *zerolog.Event
	if err.StatusCode < http.StatusInternalServerError {
		event = log.Warn()
	} else {
		event = log.Error().Err(err.Err)
	}
	event.Str("code", string(err.Code)).Msg(err.Error())
}

func writeError(ctx context.Context, rw http.ResponseWriter, err *errors.Error) {
	rw.WriteHeader(err.StatusCode)
	if err.StatusCode < http.StatusInternalServerError {
		rw.Write([]byte(err.Message))
	}
}
