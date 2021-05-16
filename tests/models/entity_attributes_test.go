// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2021, Filippov Alex
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

package models

import (
	"encoding/json"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestEntityAttributes(t *testing.T) {

	const data = `
{
  "s": "string",
  "i": 123,
  "f": 456.123,
  "b": true,
  "m": {
    "s2": "string",
    "i2": 123,
    "f2": 456.123,
    "b2": true,
    "m2": {
      "s3": "string",
      "i3": 123,
      "f3": 456.123,
      "b3": true
    }
  }
}`

	const sourceAttrs = `
{
  "ozone": {
    "name": "ozone",
    "type": "float"
  },
  "datetime": {
    "name": "datetime",
    "type": "time"
  },
  "humidity": {
    "name": "humidity",
    "type": "float"
  },
  "attribution": {
    "name": "attribution",
    "type": "string"
  },
  "forecast_day1": {
    "name": "forecast_day1",
    "type": "map",
    "value": {
      "ozone": {
        "name": "ozone",
        "type": "float"
      },
      "datetime": {
        "name": "datetime",
        "type": "time"
      },
      "humidity": {
        "name": "humidity",
        "type": "float"
      },
      "pressure": {
        "name": "pressure",
        "type": "float"
      },
      "visibility": {
        "name": "visibility",
        "type": "float"
      },
      "wind_speed": {
        "name": "wind_speed",
        "type": "float"
      },
      "temperature": {
        "name": "temperature",
        "type": "float"
      },
      "wind_bearing": {
        "name": "wind_bearing",
        "type": "float"
      },
      "max_temperature": {
        "name": "max_temperature",
        "type": "float"
      },
      "min_temperature": {
        "name": "min_temperature",
        "type": "float"
      }
    }
  }
}`
	const sourceAttrsValue = `{
  "ozone": null,
  "datetime": "2021-04-21T23:20:07.829185+07:00",
  "humidity": 78.4,
  "attribution": "Weather forecast from met.no, delivered by the Norwegian Meteorological Institute.",
  "forecast_day1": {
    "ozone": null,
    "datetime": "2021-04-22T00:00:00+07:00",
    "humidity": 78.4,
    "pressure": 1037.8,
    "visibility": null,
    "wind_speed": 9.72,
    "temperature": null,
    "wind_bearing": 27.3,
    "max_temperature": 5.8,
    "min_temperature": -4.2
  },
  "max_temperature": 5.4,
  "min_temperature": -4.2
}`

	t.Run("deserialize", func(t *testing.T) {
		Convey("deserialize", t, func(ctx C) {

			obj := map[string]interface{}{}
			err := json.Unmarshal([]byte(data), &obj)
			So(err, ShouldBeNil)

			var attrs = NetEntityAttr()
			changed, err := attrs.Deserialize(obj)
			So(err, ShouldBeNil)
			So(changed, ShouldEqual, true)

			So(attrs["s"].Value, ShouldEqual, "string")
			So(attrs["s"].String(), ShouldEqual, "string")
			So(attrs["i"].Value, ShouldEqual, 123)
			So(attrs["i"].Int64(), ShouldEqual, 123)
			So(attrs["f"].Value, ShouldEqual, 456.123)
			So(attrs["f"].Float64(), ShouldEqual, 456.123)
			So(attrs["b"].Value, ShouldEqual, true)
			So(attrs["b"].Bool(), ShouldEqual, true)
			m1 := attrs["m"].Map()
			So(m1, ShouldNotBeNil)
			So(m1["s2"].Value, ShouldEqual, "string")
			So(m1["s2"].String(), ShouldEqual, "string")
			So(m1["i2"].Value, ShouldEqual, 123)
			So(m1["i2"].Int64(), ShouldEqual, 123)
			So(m1["f2"].Value, ShouldEqual, 456.123)
			So(m1["f2"].Float64(), ShouldEqual, 456.123)
			So(m1["b2"].Value, ShouldEqual, true)
			So(m1["b2"].Bool(), ShouldEqual, true)
			m2 := m1["m2"].Map()
			So(m2["s3"].Value, ShouldEqual, "string")
			So(m2["s3"].String(), ShouldEqual, "string")
			So(m2["i3"].Value, ShouldEqual, 123)
			So(m2["i3"].Int64(), ShouldEqual, 123)
			So(m2["f3"].Value, ShouldEqual, 456.123)
			So(m2["f3"].Float64(), ShouldEqual, 456.123)
			So(m2["b3"].Value, ShouldEqual, true)
			So(m2["b3"].Bool(), ShouldEqual, true)
		})
	})

	t.Run("deserialize from string", func(t *testing.T) {
		Convey("deserialize from string", t, func(ctx C) {

			attrVal := make(m.EntityAttributeValue)
			err := json.Unmarshal([]byte(sourceAttrsValue), &attrVal)
			So(err, ShouldBeNil)

			var attrs = make(m.EntityAttributes)
			err = json.Unmarshal([]byte(sourceAttrs), &attrs)
			So(err, ShouldBeNil)

			changed, err := attrs.Deserialize(attrVal)
			So(err, ShouldBeNil)
			So(changed, ShouldEqual, true)

			So(attrs["attribution"].Value, ShouldEqual, "Weather forecast from met.no, delivered by the Norwegian Meteorological Institute.")
			So(attrs["ozone"].Value, ShouldBeEmpty)
			So(attrs["datetime"].Value, ShouldEqual, "2021-04-21T23:20:07.829185+07:00")
			So(attrs["humidity"].Value, ShouldEqual, 78.4)
			m := attrs["forecast_day1"].Map()
			So(m["datetime"].Value, ShouldEqual, "2021-04-22T00:00:00+07:00")
			So(m["humidity"].Value, ShouldEqual, 78.4)
			So(m["max_temperature"].Value, ShouldEqual, 5.8)
			So(m["min_temperature"].Value, ShouldEqual, -4.2)
			So(m["pressure"].Value, ShouldEqual, 1037.8)
			So(m["wind_bearing"].Value, ShouldEqual, 27.3)
			So(m["wind_speed"].Value, ShouldEqual, 9.72)
		})
	})

	t.Run("serialize", func(t *testing.T) {
		Convey("serialize", t, func(ctx C) {

			obj := map[string]interface{}{}
			err := json.Unmarshal([]byte(data), &obj)
			So(err, ShouldBeNil)

			var attrs = NetEntityAttr()
			changed, err := attrs.Deserialize(obj)
			So(err, ShouldBeNil)
			So(changed, ShouldEqual, true)

			s := attrs.Serialize()
			So(s["s"], ShouldEqual, "string")
			So(s["i"], ShouldEqual, 123)
			So(s["f"], ShouldEqual, 456.123)
			So(s["b"], ShouldEqual, true)
			m1, ok := s["m"].(m.EntityAttributeValue)
			So(ok, ShouldEqual, true)
			So(m1["s2"], ShouldEqual, "string")
			So(m1["i2"], ShouldEqual, 123)
			So(m1["f2"], ShouldEqual, 456.123)
			So(m1["b2"], ShouldEqual, true)
			So(m1["m2"], ShouldNotBeNil)
			m2, ok := m1["m2"].(m.EntityAttributeValue)
			So(ok, ShouldEqual, true)
			So(m2["s3"], ShouldEqual, "string")
			So(m2["i3"], ShouldEqual, 123)
			So(m2["f3"], ShouldEqual, 456.123)
			So(m2["b3"], ShouldEqual, true)
		})
	})

	t.Run("signature", func(t *testing.T) {
		Convey("signature", t, func(ctx C) {

			obj := map[string]interface{}{}
			err := json.Unmarshal([]byte(data), &obj)
			So(err, ShouldBeNil)

			var attrs = NetEntityAttr()
			changed, err := attrs.Deserialize(obj)
			So(err, ShouldBeNil)
			So(changed, ShouldEqual, true)

			s := attrs.Signature()
			So(s["s"].Name, ShouldEqual, "s")
			So(s["s"].Value, ShouldBeNil)
			So(s["s"].Type, ShouldEqual, common.EntityAttributeString)
			So(s["i"].Name, ShouldEqual, "i")
			So(s["i"].Value, ShouldBeNil)
			So(s["i"].Type, ShouldEqual, common.EntityAttributeInt)
			So(s["f"].Name, ShouldEqual, "f")
			So(s["f"].Value, ShouldBeNil)
			So(s["f"].Type, ShouldEqual, common.EntityAttributeFloat)
			So(s["b"].Name, ShouldEqual, "b")
			So(s["b"].Value, ShouldBeNil)
			So(s["b"].Type, ShouldEqual, common.EntityAttributeBool)
			So(s["m"].Name, ShouldEqual, "m")
			So(s["m"].Value, ShouldNotBeNil)
			So(s["m"].Type, ShouldEqual, common.EntityAttributeMap)

			m1 := s["m"].Map()
			So(m1["s2"].Name, ShouldEqual, "s2")
			So(m1["s2"].Value, ShouldBeNil)
			So(m1["s2"].Type, ShouldEqual, common.EntityAttributeString)
			So(m1["i2"].Name, ShouldEqual, "i2")
			So(m1["i2"].Value, ShouldBeNil)
			So(m1["i2"].Type, ShouldEqual, common.EntityAttributeInt)
			So(m1["f2"].Name, ShouldEqual, "f2")
			So(m1["f2"].Value, ShouldBeNil)
			So(m1["f2"].Type, ShouldEqual, common.EntityAttributeFloat)
			So(m1["b2"].Name, ShouldEqual, "b2")
			So(m1["b2"].Value, ShouldBeNil)
			So(m1["b2"].Type, ShouldEqual, common.EntityAttributeBool)
			So(m1["m2"].Name, ShouldEqual, "m2")
			So(m1["m2"].Value, ShouldNotBeNil)
			So(m1["m2"].Type, ShouldEqual, common.EntityAttributeMap)

			m2 := m1["m2"].Map()
			So(m2["s3"].Name, ShouldEqual, "s3")
			So(m2["s3"].Value, ShouldBeNil)
			So(m2["s3"].Type, ShouldEqual, common.EntityAttributeString)
			So(m2["i3"].Name, ShouldEqual, "i3")
			So(m2["i3"].Value, ShouldBeNil)
			So(m2["i3"].Type, ShouldEqual, common.EntityAttributeInt)
			So(m2["f3"].Name, ShouldEqual, "f3")
			So(m2["f3"].Value, ShouldBeNil)
			So(m2["f3"].Type, ShouldEqual, common.EntityAttributeFloat)
			So(m2["b3"].Name, ShouldEqual, "b3")
			So(m2["b3"].Value, ShouldBeNil)
			So(m2["b3"].Type, ShouldEqual, common.EntityAttributeBool)
		})
	})

	t.Run("copy", func(t *testing.T) {
		Convey("copy", t, func(ctx C) {

			obj := map[string]interface{}{}
			err := json.Unmarshal([]byte(data), &obj)
			So(err, ShouldBeNil)

			var attrs = NetEntityAttr()
			changed, err := attrs.Deserialize(obj)
			So(err, ShouldBeNil)
			So(changed, ShouldEqual, true)

			cpy := attrs.Copy()

			So(cpy["s"].Value, ShouldEqual, "string")
			So(cpy["s"].String(), ShouldEqual, "string")
			So(cpy["i"].Value, ShouldEqual, 123)
			So(cpy["i"].Int64(), ShouldEqual, 123)
			So(cpy["f"].Value, ShouldEqual, 456.123)
			So(cpy["f"].Float64(), ShouldEqual, 456.123)
			So(cpy["b"].Value, ShouldEqual, true)
			So(cpy["b"].Bool(), ShouldEqual, true)
			m1 := attrs["m"].Map()
			So(m1, ShouldNotBeNil)
			So(m1["s2"].Value, ShouldEqual, "string")
			So(m1["s2"].String(), ShouldEqual, "string")
			So(m1["i2"].Value, ShouldEqual, 123)
			So(m1["i2"].Int64(), ShouldEqual, 123)
			So(m1["f2"].Value, ShouldEqual, 456.123)
			So(m1["f2"].Float64(), ShouldEqual, 456.123)
			So(m1["b2"].Value, ShouldEqual, true)
			So(m1["b2"].Bool(), ShouldEqual, true)
			m2 := m1["m2"].Map()
			So(m2["s3"].Value, ShouldEqual, "string")
			So(m2["s3"].String(), ShouldEqual, "string")
			So(m2["i3"].Value, ShouldEqual, 123)
			So(m2["i3"].Int64(), ShouldEqual, 123)
			So(m2["f3"].Value, ShouldEqual, 456.123)
			So(m2["f3"].Float64(), ShouldEqual, 456.123)
			So(m2["b3"].Value, ShouldEqual, true)
			So(m2["b3"].Bool(), ShouldEqual, true)
		})
	})

	t.Run("serialize+copy", func(t *testing.T) {
		Convey("serialize + copy", t, func(ctx C) {

			obj := map[string]interface{}{}
			err := json.Unmarshal([]byte(data), &obj)
			So(err, ShouldBeNil)

			var attrs = NetEntityAttr()
			changed, err := attrs.Deserialize(obj)
			So(err, ShouldBeNil)
			So(changed, ShouldEqual, true)

			attrs.Serialize()

			cpy := attrs.Copy()

			So(cpy["s"].Value, ShouldEqual, "string")
			So(cpy["s"].String(), ShouldEqual, "string")
			So(cpy["i"].Value, ShouldEqual, 123)
			So(cpy["i"].Int64(), ShouldEqual, 123)
			So(cpy["f"].Value, ShouldEqual, 456.123)
			So(cpy["f"].Float64(), ShouldEqual, 456.123)
			So(cpy["b"].Value, ShouldEqual, true)
			So(cpy["b"].Bool(), ShouldEqual, true)
			m1 := attrs["m"].Map()
			So(m1, ShouldNotBeNil)
			So(m1["s2"].Value, ShouldEqual, "string")
			So(m1["s2"].String(), ShouldEqual, "string")
			So(m1["i2"].Value, ShouldEqual, 123)
			So(m1["i2"].Int64(), ShouldEqual, 123)
			So(m1["f2"].Value, ShouldEqual, 456.123)
			So(m1["f2"].Float64(), ShouldEqual, 456.123)
			So(m1["b2"].Value, ShouldEqual, true)
			So(m1["b2"].Bool(), ShouldEqual, true)
			m2 := m1["m2"].Map()
			So(m2["s3"].Value, ShouldEqual, "string")
			So(m2["s3"].String(), ShouldEqual, "string")
			So(m2["i3"].Value, ShouldEqual, 123)
			So(m2["i3"].Int64(), ShouldEqual, 123)
			So(m2["f3"].Value, ShouldEqual, 456.123)
			So(m2["f3"].Float64(), ShouldEqual, 456.123)
			So(m2["b3"].Value, ShouldEqual, true)
			So(m2["b3"].Bool(), ShouldEqual, true)
		})
	})

	t.Run("signature+copy", func(t *testing.T) {
		Convey("signature + copy", t, func(ctx C) {

			obj := map[string]interface{}{}
			err := json.Unmarshal([]byte(data), &obj)
			So(err, ShouldBeNil)

			var attrs = NetEntityAttr()
			changed, err := attrs.Deserialize(obj)
			So(err, ShouldBeNil)
			So(changed, ShouldEqual, true)

			attrs.Signature()

			cpy := attrs.Copy()

			So(cpy["s"].Value, ShouldEqual, "string")
			So(cpy["s"].String(), ShouldEqual, "string")
			So(cpy["i"].Value, ShouldEqual, 123)
			So(cpy["i"].Int64(), ShouldEqual, 123)
			So(cpy["f"].Value, ShouldEqual, 456.123)
			So(cpy["f"].Float64(), ShouldEqual, 456.123)
			So(cpy["b"].Value, ShouldEqual, true)
			So(cpy["b"].Bool(), ShouldEqual, true)
			m1 := attrs["m"].Map()
			So(m1, ShouldNotBeNil)
			So(m1["s2"].Value, ShouldEqual, "string")
			So(m1["s2"].String(), ShouldEqual, "string")
			So(m1["i2"].Value, ShouldEqual, 123)
			So(m1["i2"].Int64(), ShouldEqual, 123)
			So(m1["f2"].Value, ShouldEqual, 456.123)
			So(m1["f2"].Float64(), ShouldEqual, 456.123)
			So(m1["b2"].Value, ShouldEqual, true)
			So(m1["b2"].Bool(), ShouldEqual, true)
			m2 := m1["m2"].Map()
			So(m2["s3"].Value, ShouldEqual, "string")
			So(m2["s3"].String(), ShouldEqual, "string")
			So(m2["i3"].Value, ShouldEqual, 123)
			So(m2["i3"].Int64(), ShouldEqual, 123)
			So(m2["f3"].Value, ShouldEqual, 456.123)
			So(m2["f3"].Float64(), ShouldEqual, 456.123)
			So(m2["b3"].Value, ShouldEqual, true)
			So(m2["b3"].Bool(), ShouldEqual, true)
		})
	})
}
