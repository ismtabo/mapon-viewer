package ctxt

import (
	"context"
)

type contextKey string

var appCtxKey = contextKey("appCtx")

// ApplicationContext contains the context of a request for logging and functional purposes.
type ApplicationContext struct {
	TransactionID string
	Correlator    string
	Operation     string
	Alarm         string
}

// InitApplicationContext returns a new context with a new ApplicationContext.
func InitApplicationContext(ctx context.Context) context.Context {
	appCtx := &ApplicationContext{}
	return context.WithValue(ctx, appCtxKey, appCtx)
}

// GetApplicationContext retrieves the ApplicationContext from standard context.
func GetApplicationContext(ctx context.Context) *ApplicationContext {
	if appCtx, ok := ctx.Value(appCtxKey).(*ApplicationContext); ok {
		return appCtx
	}
	return nil
}
