package controller_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/ismtabo/mapon-viewer/pkg/controller"
	"github.com/ismtabo/mapon-viewer/pkg/errors"
	"github.com/ismtabo/mapon-viewer/pkg/service"
	mockSvcs "github.com/ismtabo/mapon-viewer/pkg/test/mocks/service"
	"github.com/kataras/go-sessions/v3"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
)

const sessionsCookie = "sessionsCookie"

func TestUserControllerLogin(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userSvc := mockSvcs.NewMockUserService(ctrl)
	ssns := sessions.New(sessions.Config{Cookie: sessionsCookie})
	securitySvc := service.NewSecurityService(ssns)
	userCtrl := controller.NewUserController(userSvc, securitySvc)
	t.Run("Test UserController correct login", func(t *testing.T) {
		rw := httptest.NewRecorder()
		body := strings.NewReader("email=email&password=password")
		req := httptest.NewRequest("POST", "/auth/login", body)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		userSvc.
			EXPECT().
			ValidateUserPassword(gomock.Any(), "email", "password").
			Return(nil)
		userCtrl.LoginUser(rw, req.WithContext(context.Background()))
		assert.EqualValues(t, rw.Result().StatusCode, http.StatusOK)
		g.Expect(rw.Result().Cookies()).To(gomega.ContainElement(gomega.ContainSubstring(sessionsCookie)))
	})
	t.Run("Test UserController login missing email", func(t *testing.T) {
		rw := httptest.NewRecorder()
		body := strings.NewReader("password=password")
		req := httptest.NewRequest("POST", "/auth/login", body)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		userCtrl.LoginUser(rw, req.WithContext(context.Background()))
		assert.EqualValues(t, http.StatusBadRequest, rw.Result().StatusCode)
		assert.EqualValues(t, "missing email form field", rw.Body.String())
		g.Expect(rw.Result().Cookies()).NotTo(gomega.ContainElement(gomega.ContainSubstring(sessionsCookie)))
	})
	t.Run("Test UserController login missing password", func(t *testing.T) {
		rw := httptest.NewRecorder()
		body := strings.NewReader("email=email")
		req := httptest.NewRequest("POST", "/auth/login", body)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		userCtrl.LoginUser(rw, req.WithContext(context.Background()))
		assert.EqualValues(t, http.StatusBadRequest, rw.Result().StatusCode)
		assert.EqualValues(t, "missing password form field", rw.Body.String())
		g.Expect(rw.Result().Cookies()).NotTo(gomega.ContainElement(gomega.ContainSubstring(sessionsCookie)))
	})
	t.Run("Test UserController login incorrect password", func(t *testing.T) {
		rw := httptest.NewRecorder()
		body := strings.NewReader("email=email&password=password")
		req := httptest.NewRequest("POST", "/auth/login", body)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		userSvc.
			EXPECT().
			ValidateUserPassword(gomock.Any(), "email", "password").
			Return(errors.NewUnauthorizedError())
		userCtrl.LoginUser(rw, req.WithContext(context.Background()))
		assert.EqualValues(t, http.StatusUnauthorized, rw.Result().StatusCode)
		g.Expect(rw.Result().Cookies()).NotTo(gomega.ContainElement(gomega.ContainSubstring(sessionsCookie)))
	})
	t.Run("Test UserController login user not found", func(t *testing.T) {
		rw := httptest.NewRecorder()
		body := strings.NewReader("email=email&password=password")
		req := httptest.NewRequest("POST", "/auth/login", body)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		userSvc.
			EXPECT().
			ValidateUserPassword(gomock.Any(), "email", "password").
			Return(errors.NewNotFoundError())
		userCtrl.LoginUser(rw, req.WithContext(context.Background()))
		assert.EqualValues(t, http.StatusNotFound, rw.Result().StatusCode)
		g.Expect(rw.Result().Cookies()).NotTo(gomega.ContainElement(gomega.ContainSubstring(sessionsCookie)))
	})
	t.Run("Test UserController login internal server error", func(t *testing.T) {
		rw := httptest.NewRecorder()
		body := strings.NewReader("email=email&password=password")
		req := httptest.NewRequest("POST", "/auth/login", body)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		userSvc.
			EXPECT().
			ValidateUserPassword(gomock.Any(), "email", "password").
			Return(errors.NewInternalServerError(nil))
		userCtrl.LoginUser(rw, req.WithContext(context.Background()))
		assert.EqualValues(t, http.StatusInternalServerError, rw.Result().StatusCode)
		g.Expect(rw.Result().Cookies()).NotTo(gomega.ContainElement(gomega.ContainSubstring(sessionsCookie)))
	})
}

func TestUserControllerLogout(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userSvc := mockSvcs.NewMockUserService(ctrl)
	ssns := sessions.New(sessions.Config{Cookie: sessionsCookie})
	securitySvc := service.NewSecurityService(ssns)
	userCtrl := controller.NewUserController(userSvc, securitySvc)
	t.Run("Test UserController correct logout", func(t *testing.T) {
		rw := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/auth/logout", nil)
		userCtrl.LogoutUser(rw, req.WithContext(context.Background()))
		assert.EqualValues(t, rw.Result().StatusCode, http.StatusOK)
		g.Expect(rw.Result().Cookies()).NotTo(gomega.ContainElement(gomega.ContainSubstring(sessionsCookie)))
	})
	t.Run("Test UserController correct logout after login", func(t *testing.T) {
		rw := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/auth/logout", nil)
		ssns.Start(rw, req)
		cookie := getCookie(rw, sessionsCookie)
		req.AddCookie(cookie)
		rw = httptest.NewRecorder()
		userCtrl.LogoutUser(rw, req.WithContext(context.Background()))
		assert.EqualValues(t, rw.Result().StatusCode, http.StatusOK)
		cookie = getCookie(rw, sessionsCookie)
		assert.NotNil(t, cookie)
		assert.True(t, cookie.Expires.Before(time.Now()))
	})
}

func TestUserControllerRegister(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userSvc := mockSvcs.NewMockUserService(ctrl)
	ssns := sessions.New(sessions.Config{Cookie: sessionsCookie})
	securitySvc := service.NewSecurityService(ssns)
	userCtrl := controller.NewUserController(userSvc, securitySvc)
	t.Run("Test UserController correct register", func(t *testing.T) {
		rw := httptest.NewRecorder()
		body := strings.NewReader("email=email&password=password&confirm=password")
		req := httptest.NewRequest("POST", "/auth/register", body)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		userSvc.
			EXPECT().
			CreateUser(gomock.Any(), "email", "password").
			Return(nil)
		userCtrl.RegisterUser(rw, req.WithContext(context.Background()))
		assert.EqualValues(t, rw.Result().StatusCode, http.StatusOK)
		g.Expect(rw.Result().Cookies()).To(gomega.ContainElement(gomega.ContainSubstring(sessionsCookie)))
	})
	t.Run("Test UserController register missing email", func(t *testing.T) {
		rw := httptest.NewRecorder()
		body := strings.NewReader("password=password")
		req := httptest.NewRequest("POST", "/auth/register", body)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		userCtrl.RegisterUser(rw, req.WithContext(context.Background()))
		assert.EqualValues(t, http.StatusBadRequest, rw.Result().StatusCode)
		assert.EqualValues(t, "missing email form field", rw.Body.String())
		g.Expect(rw.Result().Cookies()).NotTo(gomega.ContainElement(gomega.ContainSubstring(sessionsCookie)))
	})
	t.Run("Test UserController register missing password", func(t *testing.T) {
		rw := httptest.NewRecorder()
		body := strings.NewReader("email=email")
		req := httptest.NewRequest("POST", "/auth/register", body)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		userCtrl.RegisterUser(rw, req.WithContext(context.Background()))
		assert.EqualValues(t, http.StatusBadRequest, rw.Result().StatusCode)
		assert.EqualValues(t, "missing password form field", rw.Body.String())
		g.Expect(rw.Result().Cookies()).NotTo(gomega.ContainElement(gomega.ContainSubstring(sessionsCookie)))
	})
	t.Run("Test UserController register internal server error", func(t *testing.T) {
		rw := httptest.NewRecorder()
		body := strings.NewReader("email=email&password=password")
		req := httptest.NewRequest("POST", "/auth/register", body)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		userSvc.
			EXPECT().
			CreateUser(gomock.Any(), "email", "password").
			Return(errors.NewInternalServerError(nil))
		userCtrl.RegisterUser(rw, req.WithContext(context.Background()))
		assert.EqualValues(t, http.StatusInternalServerError, rw.Result().StatusCode)
		g.Expect(rw.Result().Cookies()).NotTo(gomega.ContainElement(gomega.ContainSubstring(sessionsCookie)))
	})
	t.Run("Test UserController register conflict error", func(t *testing.T) {
		expectedErrorMsg := "user 'email' already exists"
		rw := httptest.NewRecorder()
		body := strings.NewReader("email=email&password=password")
		req := httptest.NewRequest("POST", "/auth/register", body)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		userSvc.
			EXPECT().
			CreateUser(gomock.Any(), "email", "password").
			Return(errors.NewConflictError(expectedErrorMsg))
		userCtrl.RegisterUser(rw, req.WithContext(context.Background()))
		assert.EqualValues(t, http.StatusConflict, rw.Result().StatusCode)
		assert.EqualValues(t, expectedErrorMsg, rw.Body.String())
		g.Expect(rw.Result().Cookies()).NotTo(gomega.ContainElement(gomega.ContainSubstring(sessionsCookie)))
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
