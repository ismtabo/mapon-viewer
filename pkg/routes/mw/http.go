package mw

import (
	"net/http"
	"time"

	"github.com/ismtabo/mapon-viewer/pkg/ctxt"
	"github.com/ismtabo/mapon-viewer/pkg/service"
)

// LoggableResponseWriter is a ResponseWriter wrapper to log the response status code.
type LoggableResponseWriter struct {
	Status int
	http.ResponseWriter
}

// WriteHeader overwrites ResponseWriter's WriteHeader to store the response status code.
func (w *LoggableResponseWriter) WriteHeader(statusCode int) {
	w.Status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// LogHTTPHandler is a middleware to log the HTTP request and response.
func LogHTTPHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		log := ctxt.GetLogger(r.Context())
		if log == nil {
			getLogger(r.Context()).Fatal().Caller().Msg("Log context is not initialized. Fix HTTP middlewares")
		}
		log.Info().
			Str("path", r.RequestURI).
			Str("method", r.Method).
			Str("remoteAddr", r.RemoteAddr).
			Msg("request")
		lw := &LoggableResponseWriter{Status: http.StatusOK, ResponseWriter: w}
		next.ServeHTTP(lw, r)
		latency := int(time.Since(now).Nanoseconds() / 1000000)
		log.Info().
			Int("status", lw.Status).
			Int("latency", latency).
			Msg("response")
	})
}

// MethodsHandler is a middleware to filter incoming HTTP request by its method.
func MethodsHandler(methods ...string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			allowMethod := false
			for _, method := range methods {
				if method == r.Method {
					allowMethod = true
					break
				}
			}
			if !allowMethod {
				rw.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			next.ServeHTTP(rw, r)
		})
	}
}

// SecurityHandler is a middleware to filter incoming HTTP request returning 401 if not authenticated.
func SecurityHandler(securitySvc service.SecurityService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			if !securitySvc.IsAuthenticated(rw, r) {
				http.Error(rw, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(rw, r)
		})
	}
}
