package model

import geo "github.com/kellydunn/golang-geo"

// Track represents a route track with start and end points.
type Track struct {
	Start *geo.Point
	End   *geo.Point
}

// MaponRoute represents a mapon route information with stops points and router tracks.
type MaponRoute struct {
	Stops  []*geo.Point
	Routes []*Track
}
