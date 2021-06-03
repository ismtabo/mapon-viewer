package controller_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/ismtabo/mapon-viewer/pkg/controller"
	"github.com/ismtabo/mapon-viewer/pkg/controller/dto"
	"github.com/ismtabo/mapon-viewer/pkg/errors"
	"github.com/ismtabo/mapon-viewer/pkg/model"
	"github.com/ismtabo/mapon-viewer/pkg/test/mocks/repository"
	"github.com/jinzhu/copier"
	geo "github.com/kellydunn/golang-geo"
	"github.com/stretchr/testify/assert"
)

func TestMaponController(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := repository.NewMockMaponRepository(ctrl)
	maponCtrl := controller.NewMaponController(repo)
	t.Run("Test MaponController.", func(t *testing.T) {
		rw := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)
		q := r.URL.Query()
		q.Add("from", time.Now().Format(time.RFC3339))
		q.Add("till", time.Now().Format(time.RFC3339))
		r.URL.RawQuery = q.Encode()
		expectedRoutes := []*model.MaponRoute{
			{
				Stops: []*geo.Point{geo.NewPoint(0, 0)},
				Routes: []*model.Track{
					{
						Start: geo.NewPoint(0, 0),
						End:   geo.NewPoint(0, 0),
					},
				},
			},
		}
		expectedRoutesDTO := []*dto.MaponRoute{}
		for _, route := range expectedRoutes {
			routeDTO := &dto.MaponRoute{}
			if err := copier.Copy(routeDTO, route); err != nil {
				t.Error(err)
			}
			expectedRoutesDTO = append(expectedRoutesDTO, routeDTO)
		}
		expectedJson, err := json.Marshal(expectedRoutesDTO)
		if err != nil {
			t.Error(err)
		}
		repo.EXPECT().
			GetInfo(gomock.Any(), gomock.AssignableToTypeOf(time.Time{}), gomock.AssignableToTypeOf(time.Time{})).
			Return(expectedRoutes, nil)
		maponCtrl.GetMaponInfo(rw, r)
		assert.Equal(t, http.StatusOK, rw.Result().StatusCode)
		assert.Equal(t, string(expectedJson), rw.Body.String())
	})
	t.Run("Test MaponController. Empty repo data", func(t *testing.T) {
		rw := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)
		q := r.URL.Query()
		q.Add("from", time.Now().Format(time.RFC3339))
		q.Add("till", time.Now().Format(time.RFC3339))
		r.URL.RawQuery = q.Encode()
		repo.EXPECT().
			GetInfo(gomock.Any(), gomock.AssignableToTypeOf(time.Time{}), gomock.AssignableToTypeOf(time.Time{})).
			Return([]*model.MaponRoute{}, nil)
		maponCtrl.GetMaponInfo(rw, r)
		assert.Equal(t, http.StatusOK, rw.Result().StatusCode)
		assert.Equal(t, "[]", rw.Body.String())
	})
	t.Run("Test MaponController. Empty from query param", func(t *testing.T) {
		rw := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)
		q := r.URL.Query()
		q.Add("till", time.Now().Format(time.RFC3339))
		r.URL.RawQuery = q.Encode()
		maponCtrl.GetMaponInfo(rw, r)
		assert.Equal(t, http.StatusBadRequest, rw.Result().StatusCode)
	})
	t.Run("Test MaponController. Empty till query param", func(t *testing.T) {
		rw := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)
		q := r.URL.Query()
		q.Add("from", time.Now().Format(time.RFC3339))
		r.URL.RawQuery = q.Encode()
		maponCtrl.GetMaponInfo(rw, r)
		assert.Equal(t, http.StatusBadRequest, rw.Result().StatusCode)
	})
	t.Run("Test MaponController. Invalid from query param", func(t *testing.T) {
		rw := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)
		q := r.URL.Query()
		q.Add("from", time.Now().Format(time.ANSIC))
		q.Add("till", time.Now().Format(time.RFC3339))
		r.URL.RawQuery = q.Encode()
		maponCtrl.GetMaponInfo(rw, r)
		assert.Equal(t, http.StatusBadRequest, rw.Result().StatusCode)
	})
	t.Run("Test MaponController. Invalid till query param", func(t *testing.T) {
		rw := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)
		q := r.URL.Query()
		q.Add("from", time.Now().Format(time.RFC3339))
		q.Add("till", time.Now().Format(time.ANSIC))
		r.URL.RawQuery = q.Encode()
		maponCtrl.GetMaponInfo(rw, r)
		assert.Equal(t, http.StatusBadRequest, rw.Result().StatusCode)
	})
	t.Run("Test MaponController. Error returned from repo", func(t *testing.T) {
		rw := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)
		q := r.URL.Query()
		q.Add("from", time.Now().Format(time.RFC3339))
		q.Add("till", time.Now().Format(time.RFC3339))
		r.URL.RawQuery = q.Encode()
		expectedError := errors.NewInternalServerError(nil)
		repo.EXPECT().
			GetInfo(gomock.Any(), gomock.AssignableToTypeOf(time.Time{}), gomock.AssignableToTypeOf(time.Time{})).
			Return(nil, expectedError)
		maponCtrl.GetMaponInfo(rw, r)
		assert.Equal(t, http.StatusInternalServerError, rw.Result().StatusCode)
	})
}
