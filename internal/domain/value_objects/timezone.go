package value_objects

import (
	"time"
)

type TimeZone string

func (t TimeZone) LoadLocation() *time.Location {
	location, _ := time.LoadLocation(string(t))
	return location
}
