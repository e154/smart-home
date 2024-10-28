// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2023, Filippov Alex
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

package suncalc

// Translated in GO from the NPM library:
//   https://github.com/mourner/suncalc

import (
	"math"
	"time"
)

// date/DayTime constants and conversions
const millyToNano = 1000000
const dayMs = 1000 * 60 * 60 * 24

// J1970 ...
const J1970 = 2440588

// J2000 ...
const J2000 = 2451545

func timeToUnixMillis(date time.Time) int64   { return int64(float64(date.UnixNano()) / millyToNano) }
func unixMillisToTime(date float64) time.Time { return time.Unix(0, int64(date*millyToNano)) }
func toJulian(date time.Time) float64         { return float64(timeToUnixMillis(date))/dayMs - 0.5 + J1970 }
func fromJulian(j float64) time.Time          { return unixMillisToTime((j + 0.5 - J1970) * dayMs) }
func toDays(date time.Time) float64           { return toJulian(date) - J2000 }

// general calculations for position
const rad = math.Pi / 180
const e = rad * 23.4397 // obliquity of the Earth

func rightAscension(l float64, b float64) float64 {
	return math.Atan2(math.Sin(l)*math.Cos(e)-math.Tan(b)*math.Sin(e), math.Cos(l))
}

func declination(l float64, b float64) float64 {
	return math.Asin(math.Sin(b)*math.Cos(e) + math.Cos(b)*math.Sin(e)*math.Sin(l))
}

func azimuth(H float64, phi float64, dec float64) float64 {
	return math.Atan2(math.Sin(H), math.Cos(H)*math.Sin(phi)-math.Tan(dec)*math.Cos(phi))
}

func altitude(H float64, phi float64, dec float64) float64 {
	return math.Asin(math.Sin(phi)*math.Sin(dec) + math.Cos(phi)*math.Cos(dec)*math.Cos(H))
}

func siderealTime(d float64, lw float64) float64 { return rad*(280.16+360.9856235*d) - lw }

func astroRefraction(h float64) float64 {
	if h < 0.0 {
		h = 0 // if h = -0.08901179 a div/0 would occur.
	} // the following formula works for positive altitudes only.

	// formula 16.4 of "Astronomical Algorithms" 2nd edition by Jean Meeus (Willmann-Bell, Richmond) 1998.
	// 1.02 / tan(h + 10.26 / (h + 5.10)) h in degrees, result in arc minutes -> converted to rad:
	return 0.0002967 / math.Tan(h+0.00312536/(h+0.08901179))
}

// general sun calculations

func solarMeanAnomalyI(d float64) float64 { return solarMeanAnomalyF(d) }
func solarMeanAnomalyF(d float64) float64 { return rad * (357.5291 + 0.98560028*d) }

func eclipticLongitude(M float64) float64 {

	var C = rad * (1.9148*math.Sin(M) + 0.02*math.Sin(2*M) + 0.0003*math.Sin(3*M)) // equation of center
	var P = rad * 102.9372                                                         // perihelion of the Earth

	return M + C + P + math.Pi
}

// DayTimeName ...
type DayTimeName string

const (
	// Sunrise ...
	Sunrise DayTimeName = "sunrise"
	// Sunset ...
	Sunset DayTimeName = "sunset"

	// SunriseEnd ...
	SunriseEnd DayTimeName = "sunriseEnd"
	// SunsetStart ...
	SunsetStart DayTimeName = "sunsetStart"

	// Dawn ...
	Dawn DayTimeName = "dawn"
	// Dusk ...
	Dusk DayTimeName = "dusk"

	// NauticalDawn ...
	NauticalDawn DayTimeName = "nauticalDawn"
	// NauticalDusk ...
	NauticalDusk DayTimeName = "nauticalDusk"

	// NightEnd ...
	NightEnd DayTimeName = "nightEnd"
	// Night ...
	Night DayTimeName = "night"

	// GoldenHourEnd ...
	GoldenHourEnd DayTimeName = "goldenHourEnd"
	// GoldenHour ...
	GoldenHour DayTimeName = "goldenHour"

	// SolarNoon ...
	SolarNoon DayTimeName = "solarNoon"
	// Nadir ...
	Nadir DayTimeName = "nadir"
)

// DayTime ...
type DayTime struct {
	MorningName DayTimeName
	Time        time.Time
}

type dayTimeConf struct {
	angle       float64
	morningName DayTimeName
	eveningName DayTimeName
}

type coord struct {
	dec float64
	ra  float64
}

func sunCoords(d float64) coord {
	var M = solarMeanAnomalyI(d)
	var L = eclipticLongitude(M)

	return coord{
		declination(L, 0),
		rightAscension(L, 0),
	}
}

// SunPosition ...
type SunPosition struct {
	Azimuth  float64
	Altitude float64
}

// calculates sun position for a given date and latitude/longitude
func GetSunPosition(date time.Time, lat float64, lng float64) SunPosition {

	var lw = rad * -lng
	var phi = rad * lat
	var d = toDays(date)
	var c = sunCoords(d)
	var H = siderealTime(d, lw) - c.ra

	return SunPosition{
		azimuth(H, phi, c.dec),
		altitude(H, phi, c.dec),
	}
}

// IsNight ...
func IsNight(sunrisePos SunPosition) bool {
	solarElevation := sunrisePos.Altitude * 180 / math.Pi
	return solarElevation < -0.833
}

// sun times configuration (angle, morning name, evening name)
var times = []dayTimeConf{
	{-0.833, Sunrise, Sunset},
	{-0.3, SunriseEnd, SunsetStart},
	{-6, Dawn, Dusk},
	{-12, NauticalDawn, NauticalDusk},
	{-18, NightEnd, Night},
	{6, GoldenHourEnd, GoldenHour},
}

// calculations for sun times
const J0 = 0.0009

func julianCycle(d float64, lw float64) float64 { return math.Round(d - J0 - lw/(2*math.Pi)) }

func approxTransit(Ht float64, lw float64, n float64) float64 {
	return J0 + (Ht+lw)/(2*math.Pi) + n
}
func solarTransitJ(ds float64, M float64, L float64) float64 {
	return J2000 + ds + 0.0053*math.Sin(M) - 0.0069*math.Sin(2*L)
}
func hourAngle(h float64, phi float64, d float64) float64 {
	return math.Acos((math.Sin(h) - math.Sin(phi)*math.Sin(d)) / (math.Cos(phi) * math.Cos(d)))
}

// returns set DayTime for the given sun altitude
func getSetJ(h float64, lw float64, phi float64, dec float64, n float64, M float64, L float64) float64 {
	var w = hourAngle(h, phi, dec)
	var a = approxTransit(w, lw, n)
	return solarTransitJ(a, M, L)
}

// calculates sun times for a given date and latitude/longitude
func GetTimes(date time.Time, lat float64, lng float64) map[DayTimeName]DayTime {
	lw := rad * -lng
	phi := rad * lat

	d := toDays(date)
	n := julianCycle(d, lw)
	ds := approxTransit(0, lw, n)

	M := solarMeanAnomalyF(ds)
	L := eclipticLongitude(M)
	dec := declination(L, 0)

	Jnoon := solarTransitJ(ds, M, L)

	//i, len, DayTime, Jset, Jrise;
	var oneTime dayTimeConf
	result := make(map[DayTimeName]DayTime)

	result[SolarNoon] = DayTime{SolarNoon, fromJulian(Jnoon)}
	result[Nadir] = DayTime{Nadir, fromJulian(Jnoon - 0.5)}

	for i := 0; i < len(times); i++ {
		oneTime = times[i]

		Jset := getSetJ(oneTime.angle*rad, lw, phi, dec, n, M, L)
		Jrise := Jnoon - (Jset - Jnoon)

		result[oneTime.morningName] = DayTime{oneTime.morningName, fromJulian(Jrise)}
		result[oneTime.eveningName] = DayTime{oneTime.eveningName, fromJulian(Jset)}
	}

	return result
}

type moonCoordinates struct {
	rightAscension float64
	declination    float64
	distance       float64
}

// moon calculations, based on http://aa.quae.nl/en/reken/hemelpositie.html formulas
func moonCoords(d float64) moonCoordinates { // geocentric ecliptic coordinates of the moon
	L := rad * (218.316 + 13.176396*d) // ecliptic longitude
	M := rad * (134.963 + 13.064993*d) // mean anomaly
	F := rad * (93.272 + 13.229350*d)  // mean distance

	l := L + rad*6.289*math.Sin(M)   // longitude
	b := rad * 5.128 * math.Sin(F)   // latitude
	dt := 385001 - 20905*math.Cos(M) // distance to the moon in km

	return moonCoordinates{
		rightAscension(l, b),
		declination(l, b),
		dt,
	}
}

// MoonPosition ...
type MoonPosition struct {
	Azimuth          float64
	Altitude         float64
	Distance         float64
	ParallacticAngle float64
}

// GetMoonPosition ...
func GetMoonPosition(date time.Time, lat float64, lng float64) MoonPosition {
	lw := rad * -lng
	phi := rad * lat
	d := toDays(date)

	c := moonCoords(d)
	H := siderealTime(d, lw) - c.rightAscension
	h := altitude(H, phi, c.declination)
	// formula 14.1 of "Astronomical Algorithms" 2nd edition by Jean Meeus (Willmann-Bell, Richmond) 1998.
	pa := math.Atan2(math.Sin(H), math.Tan(phi)*math.Cos(c.declination)-math.Sin(c.declination)*math.Cos(H))
	h = h + astroRefraction(h) // altitude correction for refraction

	return MoonPosition{
		azimuth(H, phi, c.declination),
		h,
		c.distance,
		pa,
	}
}

// MoonIllumination ...
type MoonIllumination struct {
	fraction float64
	phase    float64
	angle    float64
}

// calculations for illumination parameters of the moon,
// based on http://idlastro.gsfc.nasa.gov/ftp/pro/astro/mphase.pro formulas and
// Chapter 48 of "Astronomical Algorithms" 2nd edition by Jean Meeus (Willmann-Bell, Richmond) 1998.
func GetMoonIllumination(date time.Time) MoonIllumination {

	d := toDays(date)
	s := sunCoords(d)
	m := moonCoords(d)

	sdist := 149598000. // distance from Earth to Sun in km

	phi := math.Acos(math.Sin(s.dec)*math.Sin(m.declination) + math.Cos(s.dec)*math.Cos(m.declination)*math.Cos(s.ra-m.rightAscension))
	inc := math.Atan2(sdist*math.Sin(phi), m.distance-sdist*math.Cos(phi))
	angle := math.Atan2(math.Cos(s.dec)*math.Sin(s.ra-m.rightAscension), math.Sin(s.dec)*math.Cos(m.declination)-math.Cos(s.dec)*math.Sin(m.declination)*math.Cos(s.ra-m.rightAscension))
	phaseAngle := 1.
	if angle < 0 {
		phaseAngle = -1.
	}

	return MoonIllumination{
		(1 + math.Cos(inc)) / 2,
		0.5 + 0.5*inc*phaseAngle/math.Pi,
		angle,
	}
}

func hoursLater(date time.Time, h int64) time.Time {
	return date.Add(time.Duration(h * dayMs / 24 * millyToNano))
}

// MoonTimes ...
type MoonTimes struct {
	Rise       time.Time
	Set        time.Time
	AlwaysUp   bool
	AlwaysDown bool
}

// calculations for moon rise/set times are based on http://www.stargazing.net/kepler/moonrise.html article
func GetMoonTimes(date time.Time, lat float64, lng float64, inUTC bool) MoonTimes {
	var t time.Time
	if inUTC {
		t = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	} else {
		t = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	}

	hc := 0.133 * rad
	h0 := GetMoonPosition(t, lat, lng).Altitude - hc
	//h1, h2, rise, set, a, b, xe, ye, d, roots, x1, x2, dx;
	var ye float64
	var x1 float64
	var x2 float64
	var rise float64
	var set float64

	// go in 2-hour chunks, each DayTime seeing if a 3-point quadratic curve crosses zero (which means rise or set)
	i := int64(0)
	for i <= 24 {

		h1 := GetMoonPosition(hoursLater(t, i), lat, lng).Altitude - hc
		h2 := GetMoonPosition(hoursLater(t, i+1), lat, lng).Altitude - hc
		a := (h0+h2)/2 - h1
		b := (h2 - h0) / 2
		xe := -b / (2 * a)
		ye = (a*xe+b)*xe + h1
		d := b*b - 4*a*h1
		roots := 0
		if d >= 0 {
			dx := math.Sqrt(d) / (math.Abs(a) * 2)
			x1 = xe - dx
			x2 = xe + dx
			if math.Abs(x1) <= 1 {
				roots++
			}
			if math.Abs(x2) <= 1 {
				roots++
			}
			if x1 < -1 {
				x1 = x2
			}
		}

		if roots == 1 {
			if h0 < 0 {
				rise = float64(i) + x1
			} else {
				set = float64(i) + x1
			}

		} else {
			if roots == 2 {
				if ye < 0 {
					rise = float64(i) + x2
					set = float64(i) + x1
				} else {
					rise = float64(i) + x1
					set = float64(i) + x2
				}
			}
		}
		if rise != 0 && set != 0 {
			break
		}

		h0 = h2
		i += 2
	}

	var result = MoonTimes{}

	if rise != 0 {
		result.Rise = hoursLater(t, int64(rise))
	}
	if set != 0 {
		result.Set = hoursLater(t, int64(set))
	}
	if rise == 0 && set == 0 {
		if ye > 0 {
			result.AlwaysUp = true
		} else {
			result.AlwaysDown = true
		}
	}

	return result
}
