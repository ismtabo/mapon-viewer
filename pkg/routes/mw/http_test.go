package mw_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/ismtabo/mapon-viewer/pkg/routes/mw"
	"github.com/ismtabo/mapon-viewer/pkg/service"
	"github.com/kataras/go-sessions/v3"
	"github.com/stretchr/testify/assert"
)

func TestMethodsHandler(t *testing.T) {
	handler := mw.MethodsHandler("GET")(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {}))
	t.Run("Test MethodsHandler", func(t *testing.T) {
		rw := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)
		handler.ServeHTTP(rw, r)
		assert.EqualValues(t, http.StatusOK, rw.Result().StatusCode)
	})
	t.Run("Test MethodsHandler. Method Not Allowed", func(t *testing.T) {
		rw := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/", nil)
		handler.ServeHTTP(rw, r)
		assert.EqualValues(t, http.StatusMethodNotAllowed, rw.Result().StatusCode)
	})
}

func TestSecurityHandler(t *testing.T) {
	ssns := sessions.New(sessions.Config{
		Cookie: uuid.NewString(),
	})
	secSvc := service.NewSecurityService(ssns)
	handler := mw.SecurityHandler(secSvc)(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {}))
	t.Run("Test SecurityHandler", func(t *testing.T) {
		rw := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)
		secSvc.Login(rw, r)
		copyCookies(rw, r)
		handler.ServeHTTP(rw, r)
		assert.EqualValues(t, http.StatusOK, rw.Result().StatusCode)
	})
	t.Run("Test SecurityHandler. Unauthorized", func(t *testing.T) {
		rw := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)
		handler.ServeHTTP(rw, r)
		assert.EqualValues(t, http.StatusUnauthorized, rw.Result().StatusCode)
	})
}

func copyCookies(rw *httptest.ResponseRecorder, r *http.Request) {
	cookies := rw.Result().Cookies()
	for _, cookie := range cookies {
		r.AddCookie(cookie)
	}
}
