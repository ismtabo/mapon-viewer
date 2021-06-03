package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"github.com/ismtabo/mapon-viewer/pkg/cfg"
	"github.com/ismtabo/mapon-viewer/pkg/ctxt"
	"github.com/ismtabo/mapon-viewer/pkg/errors"
	"github.com/ismtabo/mapon-viewer/pkg/model"
	"github.com/ismtabo/mapon-viewer/pkg/repository/dao"
	geo "github.com/kellydunn/golang-geo"
)

type MaponRepository interface {
	GetInfo(ctx context.Context, from, till time.Time) ([]*model.MaponRoute, error)
}

type maponRepository struct {
	*cfg.MaponConfig
}

func NewMaponRespository(config *cfg.MaponConfig) MaponRepository {
	return &maponRepository{MaponConfig: config}
}

func (r maponRepository) GetInfo(ctx context.Context, from, till time.Time) ([]*model.MaponRoute, error) {
	log := ctxt.GetLogger(ctx)
	url := fmt.Sprintf("%s/%s", r.URL, r.Endpoints.Route)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.NewInternalServerError(err)
	}
	q := req.URL.Query()
	q.Add("key", string(r.Key))
	q.Add("from", from.UTC().Format(time.RFC3339))
	q.Add("till", till.UTC().Format(time.RFC3339))
	req.URL.RawQuery = q.Encode()
	b, _ := httputil.DumpRequest(req, true)
	log.Debug().Msg(string(b))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.NewInternalServerError(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := httputil.DumpResponse(resp, true)
		err := fmt.Errorf("failed retrieving info from mapon: %s", string(b))
		return nil, errors.NewInternalServerError(err)
	}
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.NewInternalServerError(err)
	}
	bodyString := string(bytes)
	if strings.Contains(bodyString, "error") {
		err := parseMaponError(bodyString)
		return nil, errors.NewInternalServerError(err)
	}
	routeDAO := &dao.Route{}
	if err := json.Unmarshal(bytes, routeDAO); err != nil {
		return nil, errors.NewInternalServerError(err)
	}
	maponRoutes := make([]*model.MaponRoute, 0)
	if routeDAO == nil {
		return maponRoutes, nil
	}
	for _, unit := range routeDAO.Data.Units {
		routes := unit.Routes
		stopsPoints := make([]*geo.Point, 0)
		routesTrack := make([]*model.Track, 0)
		for _, route := range routes {
			switch route.Type {
			case dao.Stop:
				stopsPoints = append(stopsPoints, geo.NewPoint(*route.Start.Lat, *route.Start.Lng))
			case dao.TypeRoute:
				routeTrack := &model.Track{}
				routeTrack.Start = geo.NewPoint(*route.Start.Lat, *route.Start.Lng)
				if route.End != nil {
					routeTrack.End = geo.NewPoint(*route.End.Lat, *route.End.Lng)
				}
				routesTrack = append(routesTrack, routeTrack)
			}
		}
		maponRoutes = append(maponRoutes, &model.MaponRoute{Stops: stopsPoints, Routes: routesTrack})
	}
	return maponRoutes, nil
}

func parseMaponError(body string) error {
	maponErr := &dao.MaponError{}
	if err := json.Unmarshal([]byte(body), maponErr); err != nil {
		return err
	}
	return fmt.Errorf("error: code: %d msg: '%s'", maponErr.Err.Code, maponErr.Err.Msg)
}
