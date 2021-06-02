package ctxt_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/ismtabo/mapon-viewer/pkg/ctxt"
	"github.com/stretchr/testify/assert"
)

func TestApplicationContext(t *testing.T) {
	t.Run("Test GetApplicationContext. Empty context", func(t *testing.T) {
		appCtx := ctxt.GetApplicationContext(context.Background())
		assert.Nil(t, appCtx)
	})
	t.Run("Test init application context then get it. Empty context", func(t *testing.T) {
		ctx := ctxt.InitApplicationContext(context.Background())
		appCtx := ctxt.GetApplicationContext(ctx)
		assert.NotNil(t, appCtx)
	})
	t.Run("Test init application context then get it. Overrides already initialize context", func(t *testing.T) {
		ctx := ctxt.InitApplicationContext(context.Background())
		expectedAppCtx := ctxt.GetApplicationContext(ctx)
		expectedAppCtx.Correlator = uuid.NewString()
		ctx = ctxt.InitApplicationContext(context.Background())
		actualAppCtx := ctxt.GetApplicationContext(ctx)
		assert.NotEqual(t, expectedAppCtx, actualAppCtx)
	})
}
