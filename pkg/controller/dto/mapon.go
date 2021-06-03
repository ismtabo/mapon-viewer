package dto

import geo "github.com/kellydunn/golang-geo"

type Track struct {
	Start *geo.Point `json:"start"`
	End   *geo.Point `json:"end"`
}

type MaponRoute struct {
	Stops  []*geo.Point `json:"stops"`
	Routes []*Track     `json:"routes"`
}
