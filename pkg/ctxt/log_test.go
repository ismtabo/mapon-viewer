package ctxt_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/ismtabo/mapon-viewer/pkg/ctxt"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestLogContext(t *testing.T) {
	buffer := bytes.NewBufferString("")
	logger := zerolog.New(buffer)
	t.Run("Test GetLogContext. Empty context", func(t *testing.T) {
		appCtx := ctxt.GetLogger(context.Background())
		assert.Nil(t, appCtx)
	})
	t.Run("Test init log context then get it. Empty context", func(t *testing.T) {
		expected := &logger
		ctx := ctxt.InitLogContext(context.Background(), expected)
		actual := ctxt.GetLogger(ctx)
		assert.Equal(t, expected, actual)
	})
	t.Run("Test init log context then get it. Overrides already initialize context", func(t *testing.T) {
		expected := &logger
		ctx := ctxt.InitLogContext(context.Background(), expected)
		ctx = ctxt.InitLogContext(ctx, nil)
		actual := ctxt.GetLogger(ctx)
		assert.NotEqual(t, expected, actual)
	})
	t.Run("Test set logger", func(t *testing.T) {
		expected := &logger
		ctx := ctxt.InitLogContext(context.Background(), expected)
		ctxt.SetLogger(ctx, nil)
		actual := ctxt.GetLogger(ctx)
		assert.NotEqual(t, expected, actual)
	})
	t.Run("Test set logger. Empty context", func(t *testing.T) {
		expected := &logger
		ctx := context.Background()
		ctxt.SetLogger(ctx, expected)
		actual := ctxt.GetLogger(ctx)
		assert.Nil(t, actual)
	})
	t.Run("Test get logger", func(t *testing.T) {
		expected := &logger
		ctx := ctxt.InitLogContext(context.Background(), expected)
		actual := ctxt.GetLogger(ctx)
		assert.Equal(t, expected, actual)
	})
	t.Run("Test get logger. Empty context", func(t *testing.T) {
		actual := ctxt.GetLogger(context.Background())
		assert.Nil(t, actual)
	})
}
