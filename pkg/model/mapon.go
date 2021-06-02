package model

import geo "github.com/kellydunn/golang-geo"

type MaponInfo struct {
	Stops []*geo.Point
	Route []*geo.Point
}
