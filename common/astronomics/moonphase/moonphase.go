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

package moonphase

import (
	"math"
	"time"
)

// Moon ...
type Moon struct {
	phase     float64
	illum     float64
	age       float64
	dist      float64
	angdia    float64
	sundist   float64
	sunangdia float64
	pdata     float64
	quarters  [8]float64
	timespace float64
}

var synmonth float64 = 29.53058868 // Synodic month (new Moon to new Moon)

// New ...
func New(t time.Time) (moonP *Moon) {
	moonP = new(Moon)

	// Astronomical constants
	var epoch float64 = 2444238.5 // 1989 January 0.0

	//Constants defining the Sun's apparent orbit
	var elonge float64 = 278.833540  // Ecliptic longitude of the Sun at epoch 1980.0
	var elongp float64 = 282.596403  // Ecliptic longitude of the Sun at perigee
	var eccent float64 = 0.016718    // Eccentricity of Earth's orbit
	var sunsmax float64 = 1.495985e8 // Sun's angular size, degrees, at semi-major axis distance
	var sunangsiz float64 = 0.533128

	// Elements of the Moon's orbit, epoch 1980.0
	var mmlong float64 = 64.975464   // Moon's mean longitude at the epoch
	var mmlongp float64 = 349.383063 // Mean longitude of the perigee at the epoch
	var mecc float64 = 0.054900      // Eccentricity of the Moon's orbit
	var mangsiz float64 = 0.5181     // Moon's angular size at distance a from Earth
	var msmax float64 = 384401       // Semi-major axis of Moon's orbit in km

	moonP.timespace = float64(t.Unix())
	moonP.pdata = utcToJulian(float64(t.Unix()))
	// Calculation of the Sun's position
	var day = moonP.pdata - epoch // Date within epoch

	var n float64 = fixangle((360 / 365.2422) * day) // Mean anomaly of the Sun
	var m float64 = fixangle(n + elonge - elongp)    // Convert from perigee co-orginates to epoch 1980.0
	var ec = kepler(m, eccent)                       // Solve equation of Kepler
	ec = math.Sqrt((1+eccent)/(1-eccent)) * math.Tan(ec/2)
	ec = 2 * rad2deg(math.Atan(ec))               // True anomaly
	var lambdasun float64 = fixangle(ec + elongp) // Sun's geocentric ecliptic longitude

	var f float64 = ((1 + eccent*cos(deg2rad(ec))) / (1 - eccent*eccent)) // Orbital distance factor
	var sunDist float64 = sunsmax / f                                     // Distance to Sun in km
	var sunAng float64 = f * sunangsiz                                    // Sun's angular size in degrees

	// Calsulation of the Moon's position
	var ml float64 = fixangle(13.1763966*day + mmlong)          // Moon's mean longitude
	var mm float64 = fixangle(ml - 0.1114041*day - mmlongp)     // Moon's mean anomaly
	var ev float64 = 1.2739 * sin(deg2rad(2*(ml-lambdasun)-mm)) // Evection
	var ae float64 = 0.1858 * sin(deg2rad(m))                   // Annual equation
	var a3 float64 = 0.37 * sin(deg2rad(m))                     // Correction term
	var mmP float64 = mm + ev - ae - a3                         // Corrected anomaly
	var mec float64 = 6.2886 * sin(deg2rad(mmP))                // Correction for the equation of the centre
	var a4 float64 = 0.214 * sin(deg2rad(2*mmP))                // Another correction term
	var lP float64 = ml + ev + mec - ae + a4                    // Corrected longitude
	var v float64 = 0.6583 * sin(deg2rad(2*(lP-lambdasun)))     // Variation
	var lPP float64 = lP + v                                    // True longitude

	// Calculation of the phase of the Moon
	var moonAge float64 = lPP - lambdasun                   // Age of the Moon in degrees
	var moonPhase float64 = (1 - cos(deg2rad(moonAge))) / 2 // Phase of the Moon

	// Distance of moon from the centre of the Earth
	var moonDist float64 = (msmax * (1 - mecc*mecc)) / (1 + mecc*cos(deg2rad(mmP+mec)))

	var moonDFrac float64 = moonDist / msmax
	var moonAng float64 = mangsiz / moonDFrac // Moon's angular diameter

	// store result
	moonP.phase = fixangle(moonAge) / 360 // Phase (0 to 1)
	moonP.illum = moonPhase               // Illuminated fraction (0 to 1)
	moonP.age = synmonth * moonP.phase    // Age of moon (days)
	moonP.dist = moonDist                 // Distance (kilometres)
	moonP.angdia = moonAng                // Angular diameter (degreees)
	moonP.sundist = sunDist               // Distance to Sun (kilometres)
	moonP.sunangdia = sunAng              // Sun's angular diameter (degrees)
	moonP.phaseHunt()
	return moonP
}

func sin(a float64) float64 {
	return math.Sin(a)
}

func cos(a float64) float64 {
	return math.Cos(a)
}

func rad2deg(r float64) float64 {
	return (r * 180) / math.Pi
}

func deg2rad(d float64) float64 {
	return (d * math.Pi) / 180
}

func fixangle(a float64) float64 {
	return (a - 360*math.Floor(a/360))
}

func kepler(m float64, ecc float64) float64 {
	epsilon := 0.000001
	m = deg2rad(m)
	e := m
	var delta float64
	delta = e - ecc*math.Sin(e) - m
	e -= delta / (1 - ecc*math.Cos(e))
	for math.Abs(delta) > epsilon {
		delta = e - ecc*math.Sin(e) - m
		e -= delta / (1 - ecc*math.Cos(e))
	}
	return e
}

func (m *Moon) phaseHunt() {
	var sdate = utcToJulian(m.timespace)
	var adate = sdate - 45
	var ats = m.timespace - 86400*45
	t := time.Unix(int64(ats), 0)
	var yy = float64(t.Year())
	var mm = float64(t.Month())

	//var k1 = math.Floor(float64(yy+((mm-1)*(1/12))-1900) * 12.3685)
	var k1 = math.Floor(float64(yy+((mm-1)*(1.0/12.0))-1900) * 12.3685)
	var nt1 = meanPhase(adate, k1)
	adate = nt1
	var nt2, k2 float64

	for {
		adate += synmonth
		k2 = k1 + 1
		nt2 = meanPhase(adate, k2)
		if math.Abs(nt2-sdate) < 0.75 {
			nt2 = truePhase(k2, 0.0)
		}
		if nt1 <= sdate && nt2 > sdate {
			break
		}
		nt1 = nt2
		k1 = k2
	}

	var data [8]float64

	data[0] = truePhase(k1, 0.0)
	data[1] = truePhase(k1, 0.25)
	data[2] = truePhase(k1, 0.5)
	data[3] = truePhase(k1, 0.75)
	data[4] = truePhase(k2, 0.0)
	data[5] = truePhase(k2, 0.25)
	data[6] = truePhase(k2, 0.5)
	data[7] = truePhase(k2, 0.75)

	for i := 0; i < 8; i++ {
		m.quarters[i] = (data[i] - 2440587.5) * 86400 // convert to UNIX time
	}
}

func utcToJulian(t float64) float64 {
	return t/86400 + 2440587.5
}

//func julianToUtc(t float64) float64 {
//	return t*86400 + 2440587.5
//}

/**
  Calculates time of the mean new Moon for a given
  base date. This argument K to this function is the
  precomputed synodic month index, given by:
      K = (year - 1900) * 12.3685
  where year is expressed as a year aand fractional year
*/
func meanPhase(sdate float64, k float64) float64 {
	// Time in Julian centuries from 1900 January 0.5
	var t float64 = (sdate - 2415020.0) / 36525
	var t2 float64 = t * t
	var t3 float64 = t2 * t

	nt := 2415020.75933 + synmonth*k +
		0.0001178*t2 -
		0.000000155*t3 +
		0.00033*sin(deg2rad(166.56+132.87*t-0.009173*t2))

	return nt
}

func truePhase(k float64, phase float64) float64 {
	k += phase                  // Add phase to new moon time
	var t float64 = k / 1236.85 // Time in Julian centures from 1900 January 0.5
	var t2 float64 = t * t
	var t3 float64 = t2 * t
	var pt float64
	pt = 2415020.75933 + synmonth*k +
		0.0001178*t2 -
		0.000000155*t3 +
		0.00033*sin(deg2rad(166.56+132.87*t-0.009173*t2))

	var m, mprime, f float64
	m = 359.2242 + 29.10535608*k - 0.0000333*t2 - 0.00000347*t3       // Sun's mean anomaly
	mprime = 306.0253 + 385.81691806*k + 0.0107306*t2 + 0.00001236*t3 // Moon's mean anomaly
	f = 21.2964 + 390.67050646*k - 0.0016528*t2 - 0.00000239*t3       // Moon's argument of latitude

	if phase < 0.01 || math.Abs(phase-0.5) < 0.01 {
		// Corrections for New and Full Moon
		pt += (0.1734-0.000393*t)*sin(deg2rad(m)) +
			0.0021*sin(deg2rad(2*m)) -
			0.4068*sin(deg2rad(mprime)) +
			0.0161*sin(deg2rad(2*mprime)) -
			0.0004*sin(deg2rad(3*mprime)) +
			0.0104*sin(deg2rad(2*f)) -
			0.0051*sin(deg2rad(m+mprime)) -
			0.0074*sin(deg2rad(m-mprime)) +
			0.0004*sin(deg2rad(2*f+m)) -
			0.0004*sin(deg2rad(2*f-m)) -
			0.0006*sin(deg2rad(2*f+mprime)) +
			0.0010*sin(deg2rad(2*f-mprime)) +
			0.0005*sin(deg2rad(m+2*mprime))
	} else if math.Abs(phase-0.25) < 0.01 || math.Abs(phase-0.75) < 0.01 {
		pt += (0.1721-0.0004*t)*sin(deg2rad(m)) +
			0.0021*sin(deg2rad(2*m)) -
			0.6280*sin(deg2rad(mprime)) +
			0.0089*sin(deg2rad(2*mprime)) -
			0.0004*sin(deg2rad(3*mprime)) +
			0.0079*sin(deg2rad(2*f)) -
			0.0119*sin(deg2rad(m+mprime)) -
			0.0047*sin(deg2rad(m-mprime)) +
			0.0003*sin(deg2rad(2*f+m)) -
			0.0004*sin(deg2rad(2*f-m)) -
			0.0006*sin(deg2rad(2*f+mprime)) +
			0.0021*sin(deg2rad(2*f-mprime)) +
			0.0003*sin(deg2rad(m+2*mprime)) +
			0.0004*sin(deg2rad(m-2*mprime)) -
			0.0003*sin(deg2rad(2*m+mprime))
		if phase < 0.5 { // First quarter correction
			pt += 0.0028 - 0.0004*cos(deg2rad(m)) + 0.0003*cos(deg2rad(mprime))
		} else { // Last quarter correction
			pt += -0.0028 + 0.0004*cos(deg2rad(m)) - 0.0003*cos(deg2rad(mprime))
		}
	}

	return pt
}

// Phase ...
func (m *Moon) Phase() float64 {
	return m.phase
}

// Illumination ...
func (m *Moon) Illumination() float64 {
	return m.illum
}

// Age ...
func (m *Moon) Age() float64 {
	return m.age
}

// Distance ...
func (m *Moon) Distance() float64 {
	return m.dist
}

// Diameter ...
func (m *Moon) Diameter() float64 {
	return m.angdia
}

// SunDistance ...
func (m *Moon) SunDistance() float64 {
	return m.sundist
}

// SunDiameter ...
func (m *Moon) SunDiameter() float64 {
	return m.sunangdia
}

// NewMoon ...
func (m *Moon) NewMoon() float64 {
	return m.quarters[0]
}

// FirstQuarter ...
func (m *Moon) FirstQuarter() float64 {
	return m.quarters[1]
}

// FullMoon ...
func (m *Moon) FullMoon() float64 {
	return m.quarters[2]
}

// LastQuarter ...
func (m *Moon) LastQuarter() float64 {
	return m.quarters[3]
}

// NextNewMoon ...
func (m *Moon) NextNewMoon() float64 {
	return m.quarters[4]
}

// NextFirstQuarter ...
func (m *Moon) NextFirstQuarter() float64 {
	return m.quarters[1]
}

// NextFullMoon ...
func (m *Moon) NextFullMoon() float64 {
	return m.quarters[6]
}

// NextLastQuarter ...
func (m *Moon) NextLastQuarter() float64 {
	return m.quarters[7]
}

// PhaseName ...
func (m *Moon) PhaseName() string {
	names := map[int]string{
		0: "New Moon",
		1: "Waxing Crescent",
		2: "First Quarter",
		3: "Waxing Gibbous",
		4: "Full Moon",
		5: "Waning Gibbous",
		6: "Third Quarter",
		7: "Waning Crescent",
		8: "New Moon",
	}

	i := int(math.Floor((m.phase + 0.0625) * 8))
	return names[i]
}
