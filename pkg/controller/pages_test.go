package controller_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ismtabo/mapon-viewer/pkg/controller"
	"github.com/ismtabo/mapon-viewer/pkg/errors"
	"github.com/ismtabo/mapon-viewer/pkg/test/mocks/service"
	"github.com/ismtabo/mapon-viewer/pkg/test/mocks/template"
	"github.com/stretchr/testify/assert"
)

func TestPagesControllerIndex(t *testing.T) {
	ctrl := gomock.NewController(t)
	tmplMngr := template.NewMockManager(ctrl)
	ssnsSvc := service.NewMockSessionsService(ctrl)
	pagesCtrl := controller.NewPagesController(ssnsSvc, tmplMngr)
	t.Run("TestPagesControllerIndex", func(t *testing.T) {
		expectedBytes := []byte("index.html")
		rw := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)
		ssnsSvc.EXPECT().IsAuthenticated(rw, r).Return(true)
		tmplMngr.EXPECT().RenderFile("index.html", nil).Return(expectedBytes, nil)
		pagesCtrl.IndexPage(rw, r)
		assert.Equal(t, http.StatusOK, rw.Result().StatusCode)
		assert.Equal(t, expectedBytes, rw.Body.Bytes())
	})
	t.Run("TestPagesControllerIndex. No authenticated", func(t *testing.T) {
		rw := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)
		ssnsSvc.EXPECT().IsAuthenticated(rw, r).Return(false)
		pagesCtrl.IndexPage(rw, r)
		assert.Equal(t, http.StatusFound, rw.Result().StatusCode)
		location := rw.Header().Get("Location")
		assert.Equal(t, "/login", location)
	})
	t.Run("TestPagesControllerIndex. Error occurs", func(t *testing.T) {
		expectedErr := errors.NewInternalServerError(nil)
		rw := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)
		ssnsSvc.EXPECT().IsAuthenticated(rw, r).Return(true)
		tmplMngr.EXPECT().RenderFile("index.html", nil).Return(nil, expectedErr)
		pagesCtrl.IndexPage(rw, r)
		assert.Equal(t, http.StatusInternalServerError, rw.Result().StatusCode)
	})
}

func TestPagesControllerLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	tmplMngr := template.NewMockManager(ctrl)
	ssnsSvc := service.NewMockSessionsService(ctrl)
	pagesCtrl := controller.NewPagesController(ssnsSvc, tmplMngr)
	t.Run("TestPagesControllerLogin", func(t *testing.T) {
		expectedBytes := []byte("login.html")
		rw := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)
		ssnsSvc.EXPECT().IsAuthenticated(rw, r).Return(false)
		tmplMngr.EXPECT().RenderFile("login.html", nil).Return(expectedBytes, nil)
		pagesCtrl.LoginPage(rw, r)
		assert.Equal(t, http.StatusOK, rw.Result().StatusCode)
		assert.Equal(t, expectedBytes, rw.Body.Bytes())
	})
	t.Run("TestPagesControllerLogin. Already authenticated", func(t *testing.T) {
		rw := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)
		ssnsSvc.EXPECT().IsAuthenticated(rw, r).Return(true)
		pagesCtrl.LoginPage(rw, r)
		assert.Equal(t, http.StatusFound, rw.Result().StatusCode)
		location := rw.Header().Get("Location")
		assert.Equal(t, "/", location)
	})
	t.Run("TestPagesControllerLogin. Error occurs", func(t *testing.T) {
		expectedErr := errors.NewInternalServerError(nil)
		rw := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)
		ssnsSvc.EXPECT().IsAuthenticated(rw, r).Return(false)
		tmplMngr.EXPECT().RenderFile("login.html", nil).Return(nil, expectedErr)
		pagesCtrl.LoginPage(rw, r)
		assert.Equal(t, http.StatusInternalServerError, rw.Result().StatusCode)
	})
}
