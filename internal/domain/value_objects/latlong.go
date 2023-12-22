package value_objects

import (
	"database/sql"
	"database/sql/driver"
	"strconv"
	"strings"
)

type LatLong []string

var _ sql.Scanner = &LatLong{}
var _ driver.Valuer = LatLong{}

func (l LatLong) Latitude() float64 {
	lat, _ := strconv.ParseFloat(strings.Trim(l[0], " "), 32)
	return lat
}

func (l LatLong) Longitude() float64 {
	lng, _ := strconv.ParseFloat(strings.Trim(l[1], " "), 32)
	return lng
}

func (l LatLong) Value() (driver.Value, error) {
	return strings.Join(l, ","), nil
}

func (l *LatLong) Scan(src interface{}) error {
	if src == nil {
		return nil
	}

	latLong := strings.Split(src.(string), ",")
	*l = latLong

	return nil
}
