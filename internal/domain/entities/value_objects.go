package entities

import (
	"strconv"
	"strings"
	"time"
)

type TimeZone string
type LatLong []string

func (l LatLong) Latitude() float64 {
	lat, _ := strconv.ParseFloat(strings.Trim(l[0], " "), 32)
	return lat
}

func (l LatLong) Longitude() float64 {
	lng, _ := strconv.ParseFloat(strings.Trim(l[1], " "), 32)
	return lng
}

func (t TimeZone) LoadLocation() *time.Location {
	location, _ := time.LoadLocation(string(t))
	return location
}
