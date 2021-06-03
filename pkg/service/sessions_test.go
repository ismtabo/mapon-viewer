package service_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ismtabo/mapon-viewer/pkg/service"
	"github.com/kataras/go-sessions/v3"
	"github.com/stretchr/testify/assert"
)

func TestSecurityService(t *testing.T) {
	ssnCookie := t.Name()
	ssns := sessions.New(sessions.Config{Cookie: ssnCookie})
	secSvc := service.NewSessionsService(ssns)
	t.Run("Test SecurityService Login", func(t *testing.T) {
		rw := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/", nil)
		secSvc.Login(rw, r)
		cookie := getCookie(rw, ssnCookie)
		assert.NotNil(t, cookie)
	})
	t.Run("Test SecurityService Logout", func(t *testing.T) {
		rw := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/", nil)
		secSvc.Logout(rw, r)
		cookie := getCookie(rw, ssnCookie)
		assert.Nil(t, cookie)
	})
	t.Run("Test SecurityService IsAuthenticated", func(t *testing.T) {
		t.Run("Empty session", func(t *testing.T) {
			rw := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/", nil)
			assert.False(t, secSvc.IsAuthenticated(rw, r))
		})
		t.Run("Login then authenticated", func(t *testing.T) {
			rw := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/", nil)
			secSvc.Login(rw, r)
			copyCookies(rw, r)
			assert.True(t, secSvc.IsAuthenticated(rw, r))
		})
		t.Run("Login then Logout then Not authenticated", func(t *testing.T) {
			rw := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/", nil)
			secSvc.Login(rw, r)
			copyCookies(rw, r)
			secSvc.Logout(rw, r)
			copyCookies(rw, r)
			assert.False(t, secSvc.IsAuthenticated(rw, r))
		})
	})
}

func getCookie(rw *httptest.ResponseRecorder, name string) *http.Cookie {
	for _, cookie := range rw.Result().Cookies() {
		if cookie.Name == name {
			return cookie
		}
	}
	return nil
}

func copyCookies(rw *httptest.ResponseRecorder, r *http.Request) {
	cookies := rw.Result().Cookies()
	for _, cookie := range cookies {
		r.AddCookie(cookie)
	}
}
