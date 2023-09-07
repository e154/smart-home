package db

import (
	"bytes"
	"database/sql/driver"
	"fmt"
	"github.com/pkg/errors"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/ewkbhex"
)

type Polygon struct {
	Points []Point
}

func (p Polygon) Value() (driver.Value, error) {
	points, err := formatPoints(p.Points)
	if err != nil {
		return nil, fmt.Errorf("failed to format point for query: %v", err)
	}
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "ST_Polygon('LINESTRING(%s)'::geometry, 4326)", points)
	fmt.Println(buf.String())
	return buf.String(), nil
}

func (p *Polygon) Scan(src interface{}) (err error) {
	var geometry geom.T
	if geometry, err = ewkbhex.Decode(src.(string)); err != nil {
		return errors.Wrap(errors.New("decode value"), err.Error())
	}
	polygon, ok := geometry.(*geom.Polygon)
	if !ok {
		return errors.New("geometry is not a point")
	}
	fmt.Println(polygon.Coords())
	return nil
}

func formatPoints(polygonPoints []Point) (string, error) {
	var point []Point
	point = append(point, polygonPoints...)

	if len(point) > 0 &&
		point[0].Lon == point[len(point)-1].Lon &&
		point[0].Lat == point[len(point)-1].Lat {
		point = point[:len(point)-1]
	}

	if len(point) < 3 {
		return "", errors.New("polygon must have at least 3 unique points")
	}

	var points string
	for _, loc := range point {
		points += fmt.Sprintf("%f %f,", loc.Lon, loc.Lat)
	}

	points += fmt.Sprintf("%f %f", point[0].Lon, point[0].Lat)

	return points, nil
}
