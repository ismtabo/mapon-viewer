package dao

import "time"

// Generated by https://quicktype.io

type Route struct {
	Data RouteData `json:"data"`
}

type RouteData struct {
	Units []*RouteUnit `json:"units"`
}

type RouteUnit struct {
	UnitID int64           `json:"unit_id"`
	Routes []*RouteElement `json:"routes"`
}

type RouteElement struct {
	Type     Type   `json:"type"`
	RouteID  int64  `json:"route_id"`
	Start    Point  `json:"start"`
	End      *Point `json:"end,omitempty"`
	AvgSpeed *int64 `json:"avg_speed,omitempty"`
	MaxSpeed *int64 `json:"max_speed,omitempty"`
	Distance *int64 `json:"distance,omitempty"`
}

type Point struct {
	Time    time.Time `json:"time"`
	Address *string   `json:"address,omitempty"`
	Lat     *float64  `json:"lat,omitempty"`
	Lng     *float64  `json:"lng,omitempty"`
}

type Type string

const (
	Stop      Type = "stop"
	TypeRoute Type = "route"
)