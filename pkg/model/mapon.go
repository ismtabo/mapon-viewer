package model

import geo "github.com/kellydunn/golang-geo"

type Track struct {
	Start *geo.Point
	End   *geo.Point
}

type MaponRoute struct {
	Stops  []*geo.Point
	Routes []*Track
}
