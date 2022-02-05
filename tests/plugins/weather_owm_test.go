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

package plugins

import (
	"fmt"
	"github.com/e154/smart-home/system/event_bus/events"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	weatherPlugin "github.com/e154/smart-home/plugins/weather"
	"github.com/e154/smart-home/plugins/weather_owm"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/system/scripts"
	. "github.com/smartystreets/goconvey/convey"
)

func TestWeatherOwm(t *testing.T) {

	const serverData = `{
  "name":"home","lat":54.9022,"lon":83.0335,"elevation":150,"weatherdata":{"lat":54.9022,"lon":83.0335,"timezone":"Asia/Novosibirsk","timezone_offset":25200,"current":{"dt":1635266148,"sunrise":1635211117,"sunset":1635246289,"temp":3.05,"feels_like":-3.22,"pressure":1011,"humidity":65,"dew_point":-2.55,"uvi":0,"clouds":100,"visibility":10000,"wind_speed":10.8,"wind_deg":272,"wind_gust":15.99,"weather":[{"id":804,"main":"Clouds","description":"пасмурно","icon":"04n"}]},"hourly":[{"dt":1635264000,"temp":3.56,"feels_like":-2.25,"pressure":1011,"humidity":67,"dew_point":-1.75,"uvi":0,"clouds":100,"visibility":10000,"wind_speed":9.75,"wind_deg":268,"wind_gust":14.98,"weather":[{"id":500,"main":"Rain","description":"небольшой дождь","icon":"10n"}],"pop":0.28,"rain":{"1h":0.12}},{"dt":1635267600,"temp":3.05,"feels_like":-3.22,"pressure":1011,"humidity":65,"dew_point":-2.55,"uvi":0,"clouds":100,"visibility":10000,"wind_speed":10.8,"wind_deg":272,"wind_gust":15.99,"weather":[{"id":804,"main":"Clouds","description":"пасмурно","icon":"04n"}],"pop":0.07},{"dt":1635271200,"temp":2.88,"feels_like":-3.12,"pressure":1011,"humidity":66,"dew_point":-2.51,"uvi":0,"clouds":98,"visibility":10000,"wind_speed":9.68,"wind_deg":274,"wind_gust":15.13,"weather":[{"id":804,"main":"Clouds","description":"пасмурно","icon":"04n"}],"pop":0.02},{"dt":1635274800,"temp":2.46,"feels_like":-3.66,"pressure":1013,"humidity":66,"dew_point":-2.86,"uvi":0,"clouds":88,"visibility":10000,"wind_speed":9.62,"wind_deg":274,"wind_gust":14.53,"weather":[{"id":804,"main":"Clouds","description":"пасмурно","icon":"04n"}],"pop":0},{"dt":1635278400,"temp":1.78,"feels_like":-4.5,"pressure":1015,"humidity":65,"dew_point":-3.62,"uvi":0,"clouds":77,"visibility":10000,"wind_speed":9.42,"wind_deg":273,"wind_gust":14.22,"weather":[{"id":803,"main":"Clouds","description":"облачно с прояснениями","icon":"04n"}],"pop":0},{"dt":1635282000,"temp":0.89,"feels_like":-5.46,"pressure":1017,"humidity":66,"dew_point":-4.2,"uvi":0,"clouds":66,"visibility":10000,"wind_speed":8.81,"wind_deg":275,"wind_gust":13.59,"weather":[{"id":803,"main":"Clouds","description":"облачно с прояснениями","icon":"04n"}],"pop":0},{"dt":1635285600,"temp":-0.22,"feels_like":-6.65,"pressure":1020,"humidity":67,"dew_point":-5.79,"uvi":0,"clouds":56,"visibility":10000,"wind_speed":8.09,"wind_deg":276,"wind_gust":12.63,"weather":[{"id":803,"main":"Clouds","description":"облачно с прояснениями","icon":"04n"}],"pop":0},{"dt":1635289200,"temp":-0.73,"feels_like":-7.01,"pressure":1022,"humidity":69,"dew_point":-5.99,"uvi":0,"clouds":48,"visibility":10000,"wind_speed":7.35,"wind_deg":272,"wind_gust":11.82,"weather":[{"id":802,"main":"Clouds","description":"переменная облачность","icon":"03n"}],"pop":0},{"dt":1635292800,"temp":-1.08,"feels_like":-7.28,"pressure":1023,"humidity":71,"dew_point":-6.03,"uvi":0,"clouds":41,"visibility":10000,"wind_speed":6.94,"wind_deg":268,"wind_gust":11.17,"weather":[{"id":802,"main":"Clouds","description":"переменная облачность","icon":"03n"}],"pop":0},{"dt":1635296400,"temp":-1.3,"feels_like":-7.36,"pressure":1024,"humidity":71,"dew_point":-6.1,"uvi":0,"clouds":8,"visibility":10000,"wind_speed":6.51,"wind_deg":266,"wind_gust":10.68,"weather":[{"id":800,"main":"Clear","description":"ясно","icon":"01n"}],"pop":0},{"dt":1635300000,"temp":-1.33,"feels_like":-7.11,"pressure":1025,"humidity":71,"dew_point":-6.23,"uvi":0,"clouds":11,"visibility":10000,"wind_speed":5.94,"wind_deg":265,"wind_gust":9.98,"weather":[{"id":801,"main":"Clouds","description":"небольшая облачность","icon":"02d"}],"pop":0},{"dt":1635303600,"temp":-0.76,"feels_like":-6.36,"pressure":1026,"humidity":66,"dew_point":-6.61,"uvi":0.21,"clouds":11,"visibility":10000,"wind_speed":5.89,"wind_deg":264,"wind_gust":9.05,"weather":[{"id":801,"main":"Clouds","description":"небольшая облачность","icon":"02d"}],"pop":0},{"dt":1635307200,"temp":0.03,"feels_like":-5.24,"pressure":1026,"humidity":59,"dew_point":-7.34,"uvi":0.49,"clouds":12,"visibility":10000,"wind_speed":5.66,"wind_deg":264,"wind_gust":7.72,"weather":[{"id":801,"main":"Clouds","description":"небольшая облачность","icon":"02d"}],"pop":0},{"dt":1635310800,"temp":0.87,"feels_like":-4.09,"pressure":1027,"humidity":52,"dew_point":-8.09,"uvi":0.82,"clouds":13,"visibility":10000,"wind_speed":5.49,"wind_deg":260,"wind_gust":7.12,"weather":[{"id":801,"main":"Clouds","description":"небольшая облачность","icon":"02d"}],"pop":0},{"dt":1635314400,"temp":1.67,"feels_like":-3.09,"pressure":1026,"humidity":46,"dew_point":-9.08,"uvi":1.02,"clouds":13,"visibility":10000,"wind_speed":5.51,"wind_deg":257,"wind_gust":7.11,"weather":[{"id":801,"main":"Clouds","description":"небольшая облачность","icon":"02d"}],"pop":0},{"dt":1635318000,"temp":2.21,"feels_like":-2.42,"pressure":1026,"humidity":43,"dew_point":-9.53,"uvi":1.07,"clouds":71,"visibility":10000,"wind_speed":5.52,"wind_deg":254,"wind_gust":6.88,"weather":[{"id":803,"main":"Clouds","description":"облачно с прояснениями","icon":"04d"}],"pop":0},{"dt":1635321600,"temp":2.59,"feels_like":-1.9,"pressure":1026,"humidity":41,"dew_point":-9.56,"uvi":0.84,"clouds":47,"visibility":10000,"wind_speed":5.43,"wind_deg":250,"wind_gust":6.76,"weather":[{"id":802,"main":"Clouds","description":"переменная облачность","icon":"03d"}],"pop":0},{"dt":1635325200,"temp":2.73,"feels_like":-1.67,"pressure":1026,"humidity":41,"dew_point":-9.35,"uvi":0.49,"clouds":39,"visibility":10000,"wind_speed":5.33,"wind_deg":248,"wind_gust":6.75,"weather":[{"id":802,"main":"Clouds","description":"переменная облачность","icon":"03d"}],"pop":0},{"dt":1635328800,"temp":2.52,"feels_like":-1.65,"pressure":1026,"humidity":44,"dew_point":-8.62,"uvi":0.2,"clouds":31,"visibility":10000,"wind_speed":4.79,"wind_deg":248,"wind_gust":6.73,"weather":[{"id":802,"main":"Clouds","description":"переменная облачность","icon":"03d"}],"pop":0},{"dt":1635332400,"temp":1.4,"feels_like":-2.52,"pressure":1026,"humidity":50,"dew_point":-8.03,"uvi":0,"clouds":28,"visibility":10000,"wind_speed":3.93,"wind_deg":249,"wind_gust":6.25,"weather":[{"id":802,"main":"Clouds","description":"переменная облачность","icon":"03d"}],"pop":0},{"dt":1635336000,"temp":0.84,"feels_like":-3.17,"pressure":1026,"humidity":53,"dew_point":-7.9,"uvi":0,"clouds":25,"visibility":10000,"wind_speed":3.88,"wind_deg":246,"wind_gust":5.66,"weather":[{"id":802,"main":"Clouds","description":"переменная облачность","icon":"03n"}],"pop":0},{"dt":1635339600,"temp":0.43,"feels_like":-3.76,"pressure":1027,"humidity":53,"dew_point":-8.27,"uvi":0,"clouds":12,"visibility":10000,"wind_speed":4.01,"wind_deg":249,"wind_gust":7.18,"weather":[{"id":801,"main":"Clouds","description":"небольшая облачность","icon":"02n"}],"pop":0},{"dt":1635343200,"temp":0.11,"feels_like":-4.16,"pressure":1027,"humidity":53,"dew_point":-8.66,"uvi":0,"clouds":10,"visibility":10000,"wind_speed":4.01,"wind_deg":248,"wind_gust":7.59,"weather":[{"id":800,"main":"Clear","description":"ясно","icon":"01n"}],"pop":0},{"dt":1635346800,"temp":-0.08,"feels_like":-4.51,"pressure":1027,"humidity":53,"dew_point":-8.68,"uvi":0,"clouds":9,"visibility":10000,"wind_speed":4.19,"wind_deg":249,"wind_gust":8.14,"weather":[{"id":800,"main":"Clear","description":"ясно","icon":"01n"}],"pop":0},{"dt":1635350400,"temp":-0.28,"feels_like":-4.77,"pressure":1027,"humidity":55,"dew_point":-8.53,"uvi":0,"clouds":9,"visibility":10000,"wind_speed":4.21,"wind_deg":251,"wind_gust":8.14,"weather":[{"id":800,"main":"Clear","description":"ясно","icon":"01n"}],"pop":0},{"dt":1635354000,"temp":-0.47,"feels_like":-4.93,"pressure":1027,"humidity":57,"dew_point":-8.25,"uvi":0,"clouds":9,"visibility":10000,"wind_speed":4.1,"wind_deg":251,"wind_gust":8.02,"weather":[{"id":800,"main":"Clear","description":"ясно","icon":"01n"}],"pop":0},{"dt":1635357600,"temp":-0.65,"feels_like":-5.04,"pressure":1027,"humidity":58,"dew_point":-8.07,"uvi":0,"clouds":9,"visibility":10000,"wind_speed":3.94,"wind_deg":249,"wind_gust":7.8,"weather":[{"id":800,"main":"Clear","description":"ясно","icon":"01n"}],"pop":0},{"dt":1635361200,"temp":-0.78,"feels_like":-5.07,"pressure":1027,"humidity":60,"dew_point":-7.96,"uvi":0,"clouds":14,"visibility":10000,"wind_speed":3.77,"wind_deg":246,"wind_gust":7.16,"weather":[{"id":801,"main":"Clouds","description":"небольшая облачность","icon":"02n"}],"pop":0},{"dt":1635364800,"temp":-0.86,"feels_like":-4.98,"pressure":1028,"humidity":60,"dew_point":-7.89,"uvi":0,"clouds":15,"visibility":10000,"wind_speed":3.52,"wind_deg":242,"wind_gust":5.34,"weather":[{"id":801,"main":"Clouds","description":"небольшая облачность","icon":"02n"}],"pop":0},{"dt":1635368400,"temp":-0.96,"feels_like":-4.89,"pressure":1028,"humidity":61,"dew_point":-7.88,"uvi":0,"clouds":15,"visibility":10000,"wind_speed":3.27,"wind_deg":237,"wind_gust":4.41,"weather":[{"id":801,"main":"Clouds","description":"небольшая облачность","icon":"02n"}],"pop":0},{"dt":1635372000,"temp":-1.05,"feels_like":-4.94,"pressure":1028,"humidity":61,"dew_point":-7.96,"uvi":0,"clouds":17,"visibility":10000,"wind_speed":3.2,"wind_deg":228,"wind_gust":3.93,"weather":[{"id":801,"main":"Clouds","description":"небольшая облачность","icon":"02n"}],"pop":0},{"dt":1635375600,"temp":-1.12,"feels_like":-4.92,"pressure":1028,"humidity":61,"dew_point":-7.91,"uvi":0,"clouds":27,"visibility":10000,"wind_speed":3.08,"wind_deg":220,"wind_gust":3.82,"weather":[{"id":802,"main":"Clouds","description":"переменная облачность","icon":"03n"}],"pop":0},{"dt":1635379200,"temp":-1.13,"feels_like":-4.83,"pressure":1028,"humidity":61,"dew_point":-7.95,"uvi":0,"clouds":39,"visibility":10000,"wind_speed":2.97,"wind_deg":212,"wind_gust":3.74,"weather":[{"id":802,"main":"Clouds","description":"переменная облачность","icon":"03n"}],"pop":0},{"dt":1635382800,"temp":-1.14,"feels_like":-5,"pressure":1028,"humidity":62,"dew_point":-7.85,"uvi":0,"clouds":98,"visibility":10000,"wind_speed":3.14,"wind_deg":205,"wind_gust":4.26,"weather":[{"id":804,"main":"Clouds","description":"пасмурно","icon":"04n"}],"pop":0},{"dt":1635386400,"temp":-0.87,"feels_like":-4.85,"pressure":1028,"humidity":62,"dew_point":-7.63,"uvi":0,"clouds":98,"visibility":10000,"wind_speed":3.35,"wind_deg":201,"wind_gust":5.47,"weather":[{"id":804,"main":"Clouds","description":"пасмурно","icon":"04d"}],"pop":0},{"dt":1635390000,"temp":0.23,"feels_like":-4,"pressure":1028,"humidity":58,"dew_point":-7.39,"uvi":0.18,"clouds":99,"visibility":10000,"wind_speed":4,"wind_deg":200,"wind_gust":7.65,"weather":[{"id":804,"main":"Clouds","description":"пасмурно","icon":"04d"}],"pop":0},{"dt":1635393600,"temp":1.48,"feels_like":-3.2,"pressure":1028,"humidity":52,"dew_point":-7.57,"uvi":0.43,"clouds":99,"visibility":10000,"wind_speed":5.25,"wind_deg":203,"wind_gust":7.69,"weather":[{"id":804,"main":"Clouds","description":"пасмурно","icon":"04d"}],"pop":0},{"dt":1635397200,"temp":2.47,"feels_like":-2.27,"pressure":1028,"humidity":47,"dew_point":-7.84,"uvi":0.73,"clouds":93,"visibility":10000,"wind_speed":5.89,"wind_deg":207,"wind_gust":7.8,"weather":[{"id":804,"main":"Clouds","description":"пасмурно","icon":"04d"}],"pop":0},{"dt":1635400800,"temp":3.21,"feels_like":-1.46,"pressure":1027,"humidity":45,"dew_point":-7.98,"uvi":0.9,"clouds":81,"visibility":10000,"wind_speed":6.2,"wind_deg":207,"wind_gust":7.89,"weather":[{"id":803,"main":"Clouds","description":"облачно с прояснениями","icon":"04d"}],"pop":0},{"dt":1635404400,"temp":3.67,"feels_like":-0.94,"pressure":1027,"humidity":43,"dew_point":-8,"uvi":1.01,"clouds":18,"visibility":10000,"wind_speed":6.36,"wind_deg":207,"wind_gust":7.84,"weather":[{"id":801,"main":"Clouds","description":"небольшая облачность","icon":"02d"}],"pop":0},{"dt":1635408000,"temp":3.95,"feels_like":-0.54,"pressure":1026,"humidity":42,"dew_point":-7.97,"uvi":0.8,"clouds":18,"visibility":10000,"wind_speed":6.24,"wind_deg":207,"wind_gust":7.8,"weather":[{"id":801,"main":"Clouds","description":"небольшая облачность","icon":"02d"}],"pop":0},{"dt":1635411600,"temp":4.08,"feels_like":-0.28,"pressure":1026,"humidity":43,"dew_point":-7.7,"uvi":0.47,"clouds":18,"visibility":10000,"wind_speed":6.03,"wind_deg":204,"wind_gust":8.26,"weather":[{"id":801,"main":"Clouds","description":"небольшая облачность","icon":"02d"}],"pop":0},{"dt":1635415200,"temp":3.6,"feels_like":-0.84,"pressure":1025,"humidity":44,"dew_point":-7.89,"uvi":0.2,"clouds":19,"visibility":10000,"wind_speed":5.9,"wind_deg":200,"wind_gust":8.94,"weather":[{"id":801,"main":"Clouds","description":"небольшая облачность","icon":"02d"}],"pop":0},{"dt":1635418800,"temp":2.57,"feels_like":-2.02,"pressure":1025,"humidity":46,"dew_point":-8.07,"uvi":0,"clouds":20,"visibility":10000,"wind_speed":5.64,"wind_deg":194,"wind_gust":9.63,"weather":[{"id":801,"main":"Clouds","description":"небольшая облачность","icon":"02d"}],"pop":0},{"dt":1635422400,"temp":2.13,"feels_like":-2.73,"pressure":1025,"humidity":48,"dew_point":-8.07,"uvi":0,"clouds":20,"visibility":10000,"wind_speed":5.97,"wind_deg":194,"wind_gust":10.65,"weather":[{"id":801,"main":"Clouds","description":"небольшая облачность","icon":"02n"}],"pop":0},{"dt":1635426000,"temp":2.01,"feels_like":-3.11,"pressure":1025,"humidity":48,"dew_point":-8.15,"uvi":0,"clouds":50,"visibility":10000,"wind_speed":6.48,"wind_deg":193,"wind_gust":11.58,"weather":[{"id":802,"main":"Clouds","description":"переменная облачность","icon":"03n"}],"pop":0},{"dt":1635429600,"temp":1.93,"feels_like":-3.3,"pressure":1025,"humidity":47,"dew_point":-8.33,"uvi":0,"clouds":60,"visibility":10000,"wind_speed":6.68,"wind_deg":194,"wind_gust":12.15,"weather":[{"id":803,"main":"Clouds","description":"облачно с прояснениями","icon":"04n"}],"pop":0},{"dt":1635433200,"temp":1.97,"feels_like":-3.3,"pressure":1024,"humidity":46,"dew_point":-8.55,"uvi":0,"clouds":73,"visibility":10000,"wind_speed":6.8,"wind_deg":194,"wind_gust":12.51,"weather":[{"id":803,"main":"Clouds","description":"облачно с прояснениями","icon":"04n"}],"pop":0}],"daily":[{"dt":1635228000,"sunrise":1635211117,"sunset":1635246289,"moonrise":1635256020,"moonset":1635233220,"moon_phase":0.67,"temp":{"day":10.1,"min":3.56,"max":11.38,"night":3.56,"eve":10.98,"morn":6.58},"feels_like":{"day":8.39,"night":-2.25,"eve":9.33,"morn":3.77},"pressure":1007,"humidity":47,"dew_point":-1.05,"wind_speed":10.09,"wind_deg":274,"wind_gust":17.44,"weather":[{"id":500,"main":"Rain","description":"небольшой дождь","icon":"10d"}],"clouds":100,"pop":0.31,"rain":0.25,"uvi":0.74},{"dt":1635314400,"sunrise":1635297639,"sunset":1635332556,"moonrise":1635345780,"moonset":1635322620,"moon_phase":0.7,"temp":{"day":1.67,"min":-1.33,"max":3.05,"night":-0.28,"eve":0.84,"morn":-1.08},"feels_like":{"day":-3.09,"night":-4.77,"eve":-3.17,"morn":-7.28},"pressure":1026,"humidity":46,"dew_point":-9.08,"wind_speed":10.8,"wind_deg":272,"wind_gust":15.99,"weather":[{"id":801,"main":"Clouds","description":"небольшая облачность","icon":"02d"}],"clouds":13,"pop":0.07,"uvi":1.07},{"dt":1635400800,"sunrise":1635384161,"sunset":1635418825,"moonrise":1635436320,"moonset":1635411180,"moon_phase":0.73,"temp":{"day":3.21,"min":-1.14,"max":4.08,"night":2.1,"eve":2.13,"morn":-1.13},"feels_like":{"day":-1.46,"night":-3.18,"eve":-2.73,"morn":-4.83},"pressure":1027,"humidity":45,"dew_point":-7.98,"wind_speed":6.92,"wind_deg":196,"wind_gust":12.77,"weather":[{"id":803,"main":"Clouds","description":"облачно с прояснениями","icon":"04d"}],"clouds":81,"pop":0,"uvi":1.01},{"dt":1635487200,"sunrise":1635470683,"sunset":1635505096,"moonrise":0,"moonset":1635499140,"moon_phase":0.75,"temp":{"day":5.94,"min":2.12,"max":6.75,"night":5.57,"eve":6.08,"morn":3},"feels_like":{"day":1.05,"night":0.53,"eve":1.33,"morn":-2.56},"pressure":1020,"humidity":48,"dew_point":-4.54,"wind_speed":9.31,"wind_deg":219,"wind_gust":16.36,"weather":[{"id":803,"main":"Clouds","description":"облачно с прояснениями","icon":"04d"}],"clouds":61,"pop":0,"uvi":1.02},{"dt":1635573600,"sunrise":1635557205,"sunset":1635591368,"moonrise":1635527280,"moonset":1635586740,"moon_phase":0.79,"temp":{"day":-0.88,"min":-1.12,"max":5.67,"night":-0.92,"eve":-1.02,"morn":3.78},"feels_like":{"day":-7.56,"night":-6.71,"eve":-7.05,"morn":-2.03},"pressure":1027,"humidity":50,"dew_point":-10.22,"wind_speed":10,"wind_deg":240,"wind_gust":15.02,"weather":[{"id":600,"main":"Snow","description":"небольшой снег","icon":"13d"}],"clouds":63,"pop":0.86,"snow":0.58,"uvi":0.81},{"dt":1635660000,"sunrise":1635643728,"sunset":1635677641,"moonrise":1635618600,"moonset":1635674040,"moon_phase":0.83,"temp":{"day":0.76,"min":-1.93,"max":1.4,"night":-1.93,"eve":-0.54,"morn":-0.7},"feels_like":{"day":-4.33,"night":-7.76,"eve":-6.24,"morn":-6.38},"pressure":1016,"humidity":95,"dew_point":-0.2,"wind_speed":6.23,"wind_deg":262,"wind_gust":10.39,"weather":[{"id":600,"main":"Snow","description":"небольшой снег","icon":"13d"}],"clouds":100,"pop":1,"snow":2.76,"uvi":1},{"dt":1635746400,"sunrise":1635730250,"sunset":1635763916,"moonrise":1635710040,"moonset":1635761220,"moon_phase":0.86,"temp":{"day":-0.07,"min":-2.4,"max":1.72,"night":1.72,"eve":0.22,"morn":-2.4},"feels_like":{"day":-4.53,"night":-3.2,"eve":-4.78,"morn":-6.83},"pressure":1023,"humidity":77,"dew_point":-3.97,"wind_speed":5.85,"wind_deg":217,"wind_gust":10.1,"weather":[{"id":601,"main":"Snow","description":"снег","icon":"13d"}],"clouds":98,"pop":1,"snow":3.62,"uvi":1},{"dt":1635832800,"sunrise":1635816773,"sunset":1635850192,"moonrise":1635801540,"moonset":1635848280,"moon_phase":0.9,"temp":{"day":3.87,"min":2.18,"max":4.97,"night":4.97,"eve":4.97,"morn":3.26},"feels_like":{"day":-1.1,"night":-0.26,"eve":-0.31,"morn":-1.44},"pressure":1010,"humidity":97,"dew_point":3.11,"wind_speed":9.47,"wind_deg":243,"wind_gust":15.37,"weather":[{"id":500,"main":"Rain","description":"небольшой дождь","icon":"10d"}],"clouds":100,"pop":1,"rain":4.07,"uvi":1}],"alerts":[{"sender_name":"","event":"Wind","start":1635260400,"end":1635307200,"description":"","tags":["Wind"]},{"sender_name":"","event":"Ветер","start":1635260400,"end":1635307200,"description":"Местами порывы до 22 м/с","tags":["Wind"]}]},
  "loaded_at": "LOADED_AT"
}`

	Convey("weather_owm", t, func(ctx C) {
		_ = container.Invoke(func(adaptors *adaptors.Adaptors,
			migrations *migrations.Migrations,
			scriptService scripts.ScriptService,
			entityManager entity_manager.EntityManager,
			eventBus event_bus.EventBus,
			pluginManager common.PluginManager) {

			eventBus.Purge()
			scriptService.Purge()

			err := migrations.Purge()
			ctx.So(err, ShouldBeNil)

			// register plugins
			err = AddPlugin(adaptors, "weather")
			ctx.So(err, ShouldBeNil)
			settings := weather_owm.NewSettings()
			settings[weather_owm.AttrAppid].Value = "qweqweqwe"
			settings[weather_owm.AttrUnits].Value = "metric"
			settings[weather_owm.AttrLang].Value = "ru"
			err = AddPlugin(adaptors, "weather_owm", settings)
			ctx.So(err, ShouldBeNil)

			// add entity
			// ------------------------------------------------
			weatherEnt := GetNewWeather("home")
			weatherEnt.Settings[weatherPlugin.AttrPlugin].Value = "weather_owm"
			err = adaptors.Entity.Add(weatherEnt)
			ctx.So(err, ShouldBeNil)

			weatherEnt, err = adaptors.Entity.GetById(weatherEnt.Id)
			So(err, ShouldBeNil)

			// add weather vars
			// ------------------------------------------------

			err = adaptors.Variable.CreateOrUpdate(m.Variable{
				Name: "weather_owm.home",
				//Value: serverData,
				Value: strings.Replace(serverData, "LOADED_AT", time.Now().Format(time.RFC3339), -1),
			})
			ctx.So(err, ShouldBeNil)

			// ------------------------------------------------
			pluginManager.Start()
			entityManager.SetPluginManager(pluginManager)

			defer func() {
				pluginManager.Shutdown()
			}()

			time.Sleep(time.Second * 1)

			t.Run("add weather", func(t *testing.T) {
				Convey("weather_owm", t, func(ctx C) {

					// subscribe
					// ------------------------------------------------
					ch := make(chan events.EventRequestState, 3)
					fn := func(topic string, msg interface{}) {

						switch v := msg.(type) {
						case events.EventRequestState:
							ch <- v
						case events.EventAddedActor:
						default:
							fmt.Printf("unknown type %s\n", reflect.TypeOf(v).String())

						}
					}
					err = eventBus.Subscribe(event_bus.TopicEntities, fn)
					So(err, ShouldBeNil)

					//settings := weatherPlugin.NewSettings()
					//settings[weatherPlugin.AttrLat].Value = 54.9022
					//settings[weatherPlugin.AttrLon].Value = 83.0335
					//settings[weatherPlugin.AttrPlugin].Value = "weather_owm"
					//eventBus.Publish(event_bus.TopicEntities, event_bus.EventAddedActor{
					//	Type:       weatherPlugin.EntityWeather,
					//	EntityId:   "weather.home",
					//	Attributes: weatherPlugin.BaseForecast(),
					//	Settings:   settings,
					//})

					err = entityManager.Add(weatherEnt)
					So(err, ShouldBeNil)

					ticker := time.NewTimer(time.Second * 3)
					defer ticker.Stop()

					var msg events.EventRequestState
					var ok bool
					select {
					case msg = <-ch:
						ok = true
						break
					case <-ticker.C:
						break
					}

					ctx.So(ok, ShouldBeTrue)

					So(msg.From, ShouldEqual, "weather_owm.home")
					So(msg.To, ShouldEqual, "weather.home")
					So(msg.Attributes[weatherPlugin.AttrWeatherAttribution].String(), ShouldEqual, "Weather forecast from openweathermap api")

					err = eventBus.Unsubscribe(event_bus.TopicEntities, fn)
					So(err, ShouldBeNil)

					time.Sleep(time.Second * 1)
				})
			})

			t.Run("update weather", func(t *testing.T) {
				Convey("weather_owm", t, func(ctx C) {

					// subscribe
					// ------------------------------------------------
					ch := make(chan events.EventRequestState, 3)
					fn := func(topic string, msg interface{}) {

						switch v := msg.(type) {
						case events.EventRequestState:
							ch <- v
						case events.EventAddedActor:

						}
					}
					err = eventBus.Subscribe(event_bus.TopicEntities, fn)
					So(err, ShouldBeNil)

					//settings := weatherPlugin.NewSettings()
					//settings[weatherPlugin.AttrLat].Value = 54.9022
					//settings[weatherPlugin.AttrLon].Value = 83.0335
					//settings[weatherPlugin.AttrPlugin].Value = "weather_owm"
					//eventBus.Publish(weatherPlugin.TopicPluginWeather, weatherPlugin.EventStateChanged{
					//	Type:       weatherPlugin.EntityWeather,
					//	EntityId:   "weather.home",
					//	State:      weatherPlugin.StatePositionUpdate,
					//	//Attributes: weatherPlugin.BaseForecast(),
					//	Settings:   settings,
					//})

					weatherEnt.Settings[weatherPlugin.AttrLat].Value = 54.9033
					err = adaptors.Entity.Update(weatherEnt)
					ctx.So(err, ShouldBeNil)

					err = entityManager.Update(weatherEnt)
					So(err, ShouldBeNil)

					ticker := time.NewTimer(time.Second * 2)
					defer ticker.Stop()

					var msg events.EventRequestState
					var ok bool
					select {
					case msg = <-ch:
						ok = true
						break
					case <-ticker.C:
						break
					}

					ctx.So(ok, ShouldBeTrue)

					So(msg.From, ShouldEqual, "weather_owm.home")
					So(msg.To, ShouldEqual, "weather.home")
					So(msg.Attributes[weatherPlugin.AttrWeatherAttribution].String(), ShouldEqual, "Weather forecast from openweathermap api")

					err = eventBus.Unsubscribe(event_bus.TopicEntities, fn)
					So(err, ShouldBeNil)
				})
			})

			t.Run("remove weather", func(t *testing.T) {
				Convey("weather_owm", t, func(ctx C) {

					// subscribe
					// ------------------------------------------------
					ch := make(chan events.EventRemoveActor)
					fn := func(topic string, msg interface{}) {

						switch v := msg.(type) {
						case events.EventRequestState:
						case events.EventAddedActor:
						case events.EventRemoveActor:
							if v.PluginName == "weather_owm" {
								ch <- v
							}
						}
					}
					err = eventBus.Subscribe(event_bus.TopicEntities, fn)
					So(err, ShouldBeNil)

					//eventBus.Publish(event_bus.TopicEntities, event_bus.EventRemoveActor{
					//	Type:     weatherPlugin.EntityWeather,
					//	EntityId: "weather.home",
					//})

					entityManager.Remove(weatherEnt.Id)

					ticker := time.NewTimer(time.Second * 2)
					defer ticker.Stop()

					var msg events.EventRemoveActor
					var ok bool
					select {
					case msg = <-ch:
						ok = true
						break
					case <-ticker.C:
						break
					}

					ctx.So(ok, ShouldBeTrue)

					So(msg.EntityId, ShouldEqual, "weather_owm.home")

					err = eventBus.Unsubscribe(event_bus.TopicEntities, fn)
					So(err, ShouldBeNil)
				})
			})

			t.Run("weather_owm", func(t *testing.T) {
				Convey("weather_owm", t, func(ctx C) {

					w := weather_owm.NewWeatherOwm(eventBus, adaptors, weather_owm.NewSettings())
					w.AddWeather("weather.home", weatherEnt.Settings)

					loc, _ := time.LoadLocation("Asia/Novosibirsk")
					now := time.Date(2021, 5, 29, 0, 0, 0, 0, loc)
					f, err := w.GetForecast(weather_owm.Zone{
						Name: "home",
						Lat:  weatherEnt.Settings[weatherPlugin.AttrLat].Float64(),
						Lon:  weatherEnt.Settings[weatherPlugin.AttrLon].Float64(),
					}, now)
					So(err, ShouldEqual, nil)

					//fmt.Println("------")
					//debug.Println(f)

					attr := weatherPlugin.BaseForecast()
					ch, err := attr.Deserialize(f)
					So(err, ShouldEqual, nil)
					So(ch, ShouldEqual, true)

					So(attr[weatherPlugin.AttrWeatherAttribution].String(), ShouldEqual, "Weather forecast from openweathermap api")
					//So(attr[weatherPlugin.AttrWeatherDatetime].Time().String(), ShouldEqual, "2021-10-26 23:35:48 +0700 +07")
					So(attr[weatherPlugin.AttrWeatherMain].String(), ShouldEqual, "cloudy")
					So(attr[weatherPlugin.AttrWeatherDescription].String(), ShouldEqual, "cloudy")
					So(attr[weatherPlugin.AttrWeatherIcon].String(), ShouldEqual, "data/static/weather/yr/04.svg")
					So(attr[weatherPlugin.AttrWeatherHumidity].Float64(), ShouldEqual, 65)
					So(attr[weatherPlugin.AttrWeatherTemperature].Float64(), ShouldEqual, 3.05)
					So(attr[weatherPlugin.AttrWeatherMaxTemperature].Float64(), ShouldEqual, 11.38)
					So(attr[weatherPlugin.AttrWeatherMinTemperature].Float64(), ShouldEqual, 3.56)
					So(attr[weatherPlugin.AttrWeatherPressure].Float64(), ShouldEqual, 1011)
					So(attr[weatherPlugin.AttrWeatherVisibility].Int64(), ShouldEqual, 10000)
					So(attr[weatherPlugin.AttrWeatherWindBearing].Float64(), ShouldEqual, 272)
					So(attr[weatherPlugin.AttrWeatherWindSpeed].Float64(), ShouldEqual, 10.8)

					// day1
					day1 := attr[weatherPlugin.AttrForecastDay1].Map()
					So(day1[weatherPlugin.AttrWeatherMain].String(), ShouldEqual, "fair")
					So(day1[weatherPlugin.AttrWeatherDescription].String(), ShouldEqual, "fair")
					So(day1[weatherPlugin.AttrWeatherIcon].String(), ShouldEqual, "data/static/weather/yr/02d.svg")
					//So(day1[weatherPlugin.AttrWeatherDatetime].Time().String(), ShouldEqual, "2021-10-27 13:00:00 +0700 +07")
					So(day1[weatherPlugin.AttrWeatherHumidity].Float64(), ShouldEqual, 46)
					So(day1[weatherPlugin.AttrWeatherMaxTemperature].Float64(), ShouldEqual, 3.05)
					So(day1[weatherPlugin.AttrWeatherMinTemperature].Float64(), ShouldEqual, -1.33)
					So(day1[weatherPlugin.AttrWeatherPressure].Float64(), ShouldEqual, 1026)
					So(day1[weatherPlugin.AttrWeatherWindBearing].Float64(), ShouldEqual, 272)
					So(day1[weatherPlugin.AttrWeatherWindSpeed].Float64(), ShouldEqual, 10.8)

					// day2
					day2 := attr[weatherPlugin.AttrForecastDay2].Map()
					So(day2[weatherPlugin.AttrWeatherMain].String(), ShouldEqual, "partlyCloudy")
					So(day2[weatherPlugin.AttrWeatherDescription].String(), ShouldEqual, "partly cloudy")
					So(day2[weatherPlugin.AttrWeatherIcon].String(), ShouldEqual, "data/static/weather/yr/03d.svg")
					//So(day2[weatherPlugin.AttrWeatherDatetime].Time().String(), ShouldEqual, "2021-10-28 13:00:00 +0700 +07")
					So(day2[weatherPlugin.AttrWeatherHumidity].Float64(), ShouldEqual, 45)
					So(day2[weatherPlugin.AttrWeatherMaxTemperature].Float64(), ShouldEqual, 4.08)
					So(day2[weatherPlugin.AttrWeatherMinTemperature].Float64(), ShouldEqual, -1.14)
					So(day2[weatherPlugin.AttrWeatherPressure].Float64(), ShouldEqual, 1027)
					So(day2[weatherPlugin.AttrWeatherWindBearing].Float64(), ShouldEqual, 196)
					So(day2[weatherPlugin.AttrWeatherWindSpeed].Float64(), ShouldEqual, 6.92)

					// day3
					day3 := attr[weatherPlugin.AttrForecastDay3].Map()
					So(day3[weatherPlugin.AttrWeatherMain].String(), ShouldEqual, "partlyCloudy")
					So(day3[weatherPlugin.AttrWeatherDescription].String(), ShouldEqual, "partly cloudy")
					So(day3[weatherPlugin.AttrWeatherIcon].String(), ShouldEqual, "data/static/weather/yr/03d.svg")
					//So(day3[weatherPlugin.AttrWeatherDatetime].Time().String(), ShouldEqual, "2021-10-29 13:00:00 +0700 +07")
					So(day3[weatherPlugin.AttrWeatherHumidity].Float64(), ShouldEqual, 48)
					So(day3[weatherPlugin.AttrWeatherMaxTemperature].Float64(), ShouldEqual, 6.75)
					So(day3[weatherPlugin.AttrWeatherMinTemperature].Float64(), ShouldEqual, 2.12)
					So(day3[weatherPlugin.AttrWeatherPressure].Float64(), ShouldEqual, 1020)
					So(day3[weatherPlugin.AttrWeatherWindBearing].Float64(), ShouldEqual, 219)
					So(day3[weatherPlugin.AttrWeatherWindSpeed].Float64(), ShouldEqual, 9.31)

					// day4
					day4 := attr[weatherPlugin.AttrForecastDay4].Map()
					So(day4[weatherPlugin.AttrWeatherMain].String(), ShouldEqual, "lightSnow")
					So(day4[weatherPlugin.AttrWeatherDescription].String(), ShouldEqual, "light snow")
					So(day4[weatherPlugin.AttrWeatherIcon].String(), ShouldEqual, "data/static/weather/yr/49.svg")
					//So(day4[weatherPlugin.AttrWeatherDatetime].Time().String(), ShouldEqual, "2021-10-30 13:00:00 +0700 +07")
					So(day4[weatherPlugin.AttrWeatherHumidity].Float64(), ShouldEqual, 50)
					So(day4[weatherPlugin.AttrWeatherMaxTemperature].Float64(), ShouldEqual, 5.67)
					So(day4[weatherPlugin.AttrWeatherMinTemperature].Float64(), ShouldEqual, -1.12)
					So(day4[weatherPlugin.AttrWeatherPressure].Float64(), ShouldEqual, 1027)
					So(day4[weatherPlugin.AttrWeatherWindBearing].Float64(), ShouldEqual, 240)
					So(day4[weatherPlugin.AttrWeatherWindSpeed].Float64(), ShouldEqual, 10)

					// day5
					day5 := attr[weatherPlugin.AttrForecastDay5].Map()
					So(day4[weatherPlugin.AttrWeatherMain].String(), ShouldEqual, "lightSnow")
					So(day5[weatherPlugin.AttrWeatherDescription].String(), ShouldEqual, "light snow")
					So(day5[weatherPlugin.AttrWeatherIcon].String(), ShouldEqual, "data/static/weather/yr/49.svg")
					//So(day5[weatherPlugin.AttrWeatherDatetime].Time().String(), ShouldEqual, "2021-10-31 13:00:00 +0700 +07")
					So(day5[weatherPlugin.AttrWeatherHumidity].Float64(), ShouldEqual, 95)
					So(day5[weatherPlugin.AttrWeatherMaxTemperature].Float64(), ShouldEqual, 1.4)
					So(day5[weatherPlugin.AttrWeatherMinTemperature].Float64(), ShouldEqual, -1.93)
					So(day5[weatherPlugin.AttrWeatherPressure].Float64(), ShouldEqual, 1016)
					So(day5[weatherPlugin.AttrWeatherWindBearing].Float64(), ShouldEqual, 262)
					So(day5[weatherPlugin.AttrWeatherWindSpeed].Float64(), ShouldEqual, 6.23)
				})
			})

		})
	})
}
