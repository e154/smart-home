// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2023, Filippov Alex
//
// This library is free software: you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 3 of the License, or (at your option) any later version.
//
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Library General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this library.  If not, see
// <https://www.gnu.org/licenses/>.

package location

import (
	m "github.com/e154/smart-home/models"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetDistanceBetweenPoints(t *testing.T) {

	point1 := m.Point{
		Lon: 80.113927,
		Lat: 6.131738,
	}
	point2 := m.Point{
		Lon: 80.108525,
		Lat: 6.125416,
	}

	distance := GetDistanceBetweenPoints(point1, point2)
	require.Equal(t, 0.93, math.Ceil(distance*100)/100)

}

func TestCalculateDistanceToPolygon(t *testing.T) {

	point1 := m.Point{
		Lon: 80.108620,
		Lat: 6.125419,
	}
	polygon := []m.Point{
		{
			Lon: 80.112948,
			Lat: 6.131909,
		},
		{
			Lon: 80.113697,
			Lat: 6.132449,
		},
		{
			Lon: 80.114448,
			Lat: 6.131780,
		},
		{
			Lon: 80.113730,
			Lat: 6.131185,
		},
	}

	distance := GetDistanceToPolygon(point1, polygon)
	require.Equal(t, 0.86, math.Ceil(distance*100)/100)

}

func TestPointInsidePolygon(t *testing.T) {

	point1 := m.Point{
		Lon: 80.103999,
		Lat: 6.121958,
	}
	polygon := []m.Point{
		{
			Lon: 80.112948,
			Lat: 6.131909,
		},
		{
			Lon: 80.113697,
			Lat: 6.132449,
		},
		{
			Lon: 80.114448,
			Lat: 6.131780,
		},
		{
			Lon: 80.113730,
			Lat: 6.131185,
		},
		{
			Lon: 80.112948,
			Lat: 6.131909,
		},
	}

	contains := PointInsidePolygon(point1, polygon)
	require.False(t, contains)

	point1 = m.Point{
		Lon: 77.782887,
		Lat: 14.242355,
	}

	contains = PointInsidePolygon(point1, polygon)
	require.False(t, contains)

	point1 = m.Point{
		Lon: 80.113927,
		Lat: 6.131738,
	}

	contains = PointInsidePolygon(point1, polygon)
	require.True(t, contains)

	point1 = m.Point{
		Lon: 80.113564,
		Lat: 6.131762,
	}

	contains = PointInsidePolygon(point1, polygon)
	require.True(t, contains)

}
