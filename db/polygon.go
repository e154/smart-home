package db

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/ewkbhex"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Polygon struct {
	Points []Point
}

func (p Polygon) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	points, _ := formatPoints(p.Points)
	return clause.Expr{
		SQL:  "ST_Polygon(?::geometry, 4326)",
		Vars: []interface{}{fmt.Sprintf("LINESTRING(%s)", points)},
	}
}

func (p *Polygon) Scan(src any) (err error) {
	var geometry geom.T
	if geometry, err = ewkbhex.Decode(src.(string)); err != nil {
		return errors.Wrap(errors.New("decode value"), err.Error())
	}
	polygon, ok := geometry.(*geom.Polygon)
	if !ok {
		return errors.New("geometry is not a point")
	}
	fmt.Println(polygon.Coords())
	for _, point := range polygon.Coords()[0] {
		p.Points = append(p.Points, Point{
			Lon: point.X(),
			Lat: point.Y(),
		})
	}
	return nil
}

func (Polygon) GormDataType() string {
	return "geometry"
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
