package mw_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/google/uuid"
	"github.com/ismtabo/mapon-viewer/pkg/ctxt"
	"github.com/ismtabo/mapon-viewer/pkg/routes/mw"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestInitAppCtxHandler(t *testing.T) {
	rw := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	handler := mw.InitAppCtxHandler(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		appCtx := ctxt.GetApplicationContext(r.Context())
		assert.NotNil(t, appCtx)
	}))
	handler.ServeHTTP(rw, req)
}

func TestInitLogCtxHandler(t *testing.T) {
	buffer := bytes.NewBuffer([]byte{})
	logger := zerolog.New(iotest.NewWriteLogger("", buffer))
	expectedLogger := &logger
	rw := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	handler := mw.InitLogCtxHandler(expectedLogger)(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		actualLogger := ctxt.GetLogger(r.Context())
		assert.Equal(t, expectedLogger, actualLogger)
	}))
	handler.ServeHTTP(rw, req)
}

func TestCorrelatorHandler(t *testing.T) {
	t.Run("Test CorrelatorHandler. Empty correlator header", func(t *testing.T) {
		rw := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", strings.NewReader(""))
		ctx := ctxt.InitApplicationContext(context.Background())
		correlatorCh := make(chan string, 1)
		handler := mw.CorrelatorHandler(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			appCtx := ctxt.GetApplicationContext(r.Context())
			assert.NotEmpty(t, appCtx.TransactionID)
			assert.Equal(t, appCtx.TransactionID, appCtx.Correlator)
			correlatorCh <- appCtx.Correlator
		}))
		handler.ServeHTTP(rw, req.WithContext(ctx))
		expectedCorrelator := <-correlatorCh
		corrHeader := rw.Result().Header[mw.CorrelatorHTTPHeader]
		assert.Len(t, corrHeader, 1)
		assert.Contains(t, corrHeader, expectedCorrelator)
	})
	t.Run("Test CorrelatorHandler. Present correlator header", func(t *testing.T) {
		rw := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", strings.NewReader(""))
		expectedCorrelator := uuid.NewString()
		req.Header.Add(mw.CorrelatorHTTPHeader, expectedCorrelator)
		ctx := ctxt.InitApplicationContext(context.Background())
		correlatorCh := make(chan string, 1)
		handler := mw.CorrelatorHandler(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			appCtx := ctxt.GetApplicationContext(r.Context())
			assert.NotEmpty(t, appCtx.TransactionID)
			assert.NotEqual(t, appCtx.TransactionID, appCtx.Correlator)
			assert.Equal(t, appCtx.Correlator, expectedCorrelator)
			correlatorCh <- appCtx.Correlator
		}))
		handler.ServeHTTP(rw, req.WithContext(ctx))
		corrHeader := rw.Result().Header[mw.CorrelatorHTTPHeader]
		assert.Len(t, corrHeader, 1)
		assert.Contains(t, corrHeader, expectedCorrelator)
	})
}
