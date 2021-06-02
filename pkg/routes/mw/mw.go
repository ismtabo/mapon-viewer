package mw

import (
	"context"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/ismtabo/mapon-viewer/pkg/ctxt"
	"github.com/rs/zerolog"
)

// InitAppCtxHandler is a middleware to initialize the ApplicationContext in the standard context.
func InitAppCtxHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := ctxt.InitApplicationContext(r.Context())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// InitLogCtxHandler is a middleware to initialize the LogContext in the standard context.
func InitLogCtxHandler(logger *zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := ctxt.InitLogContext(r.Context(), logger)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// CorrelatorHTTPHeader is the HTTP header name to store the correlator.
const CorrelatorHTTPHeader = "Unica-Correlator"

// CorrelatorHandler registers the correlator and transactionID in the application context.
// It generates a transactionID (as a UUID).
// It also obtains the correlator from a HTTP header. If the header is not available, the correlator
// is set to the transactionID.
// It also adds a HTTP header with the correlator in the response.
func CorrelatorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		trans := ""
		if guid, err := uuid.NewUUID(); err == nil {
			trans = guid.String()
		}
		corr := r.Header.Get(CorrelatorHTTPHeader)
		if corr == "" {
			corr = trans
		}
		appCtx := ctxt.GetApplicationContext(r.Context())
		if appCtx == nil {
			getLogger(r.Context()).Fatal().Caller().Msg("Application context is not initialized. Fix HTTP middlewares")
		} else {
			appCtx.TransactionID = trans
			appCtx.Correlator = corr
		}
		w.Header().Set(CorrelatorHTTPHeader, corr)
		next.ServeHTTP(w, r)
	})
}

// LogContextHandler is a middleware to update the logger with a context using the operation
// and the current ApplicationContext in the Context.
func LogContextHandler(op string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			appCtx := ctxt.GetApplicationContext(r.Context())
			if appCtx == nil {
				getLogger(r.Context()).Fatal().Caller().Msg("Application context is not initialized. Fix HTTP middlewares")
				return
			}
			appCtx.Operation = op
			logger := ctxt.GetLogger(r.Context())
			if logger == nil {
				getLogger(r.Context()).Fatal().Caller().Msg("Log context is not initialized. Fix HTTP middlewares")
				return
			}
			subLoggerCtx := logger.With().
				Str("corr", appCtx.Correlator).
				Str("trans", appCtx.TransactionID).
				Str("op", op)
			subLogger := subLoggerCtx.Logger()
			ctxt.SetLogger(r.Context(), &subLogger)
			next.ServeHTTP(w, r)
		})
	}
}

func getLogger(ctx context.Context) *zerolog.Logger {
	if log := ctxt.GetLogger(ctx); log != nil {
		return log
	}
	log := zerolog.New(os.Stdout).With().Timestamp().Logger()
	return &log
}
