package db

import (
	"bytes"
	"fmt"

	"database/sql/driver"

	"github.com/pkg/errors"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/ewkbhex"
)

type Point struct {
	Lon float64 // X
	Lat float64 // Y
}

func (p Point) Value() (driver.Value, error) {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "SRID=4326;POINT(%f %f)", p.Lon, p.Lat)
	return buf.String(), nil
}

func (p *Point) Scan(src interface{}) (err error) {
	var geometry geom.T
	if geometry, err = ewkbhex.Decode(src.(string)); err != nil {
		return errors.Wrap(errors.New("decode value"), err.Error())
	}
	point, ok := geometry.(*geom.Point)
	if !ok {
		return errors.New("geometry is not a point")
	}
	p.Lon = point.X()
	p.Lat = point.Y()
	return nil
}

