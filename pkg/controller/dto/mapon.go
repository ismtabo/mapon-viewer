package dto

import geo "github.com/kellydunn/golang-geo"

type MaponInfo struct {
	Stops []*geo.Point `json:"stops"`
	Route []*geo.Point `json:"route"`
}
