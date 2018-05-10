package calendar

import (
	"fmt"
	"math"
)

// GregDate type
type GregDate struct {
	year       float64
	month      float64
	day        float64
	hour       float64
	min        float64
	sec        float64
	decimalday bool
}

//2000-01-01 12:00:00
func (gd GregDate) String() string {
	t := (gd.day - math.Floor(gd.day))
	if t > 0 { // if day has float value , use that value as hh:mm:ss
		s := t * 86400
		h := math.Floor(s / 3600)
		s -= h * 3600
		m := math.Floor(s / 60)
		s -= m * 60
		return fmt.Sprintf("%5d-%02d-%02d %02d:%02d:%02d", int(gd.year), int(gd.month), int(gd.day), int(h), int(m), int(s))
	}
	return fmt.Sprintf("%5d-%02d-%02d %02d:%02d:%02d", int(gd.year), int(gd.month), int(gd.day), int(gd.hour), int(gd.min), int(gd.sec))
}

// ToJulianDate converts Gregorian date to Julian date
func (gd GregDate) ToJulianDate() JulianDate {
	return g2j(gd.year, gd.month, gd.day+((gd.sec/60+gd.min)/60+gd.hour)/24)
}

func g2j(y, m, d float64) JulianDate {
	greg := false
	n := 0.0
	// check if it's Gregorian Calendar or not
	if y*372+m*31+math.Floor(d) >= 588829 {
		greg = true
	}
	if m <= 2 {
		m += 12
		y--
	}
	if greg { // handle centry leap year
		n = math.Floor(y / 100)
		n = 2 - n + math.Floor(n/4)
	}
	return JulianDate(math.Floor(365.25*(y+4716)) + math.Floor(30.6001*(m+1)) + d + n - 1524.5)
}

// SetYear update year of the gregorian date
func (gd *GregDate) SetYear(y int) {
	gd.year = float64(y)
}

// SetMonth update month of the gregorian date
func (gd *GregDate) SetMonth(m int) {
	gd.month = float64(m)
}

// SetDay update day of the gregorian date, day can be a decimal.
// Once the day is set to decimal, then hh:mm:ss can not be set separately.
func (gd *GregDate) SetDay(d float64) {
	fd := math.Floor(d)
	if d != fd {
		gd.decimalday = true
	}
	gd.day = d
	// keep the decimal in GregDate struct, only expand to hh:mm:ss when printing
	// gd.decimalday = true
	// gd.day = fd
	// s := (d - fd) * 86400
	// gd.hour = math.Floor(s / 3600)
	// s -= gd.hour * 3600
	// gd.min = math.Floor(s / 60)
	// s -= gd.min * 60
	// gd.sec = math.Floor(s)
}

// SetHour update hour of the gregorian date
// If the day is a decimal, then this hour setting will be ignored.
func (gd *GregDate) SetHour(h int) {
	if gd.decimalday {
		return
	}
	gd.hour = float64(h)
}

// SetMin update minute of the gregorian date
// If the day is a decimal, then this minute setting will be ignored.
func (gd *GregDate) SetMin(m int) {
	if gd.decimalday {
		return
	}
	gd.min = float64(m)
}

// SetSec update second of the gregorian date
// If the day is a decimal, then this minute setting will be ignored.
func (gd *GregDate) SetSec(s int) {
	if gd.decimalday {
		return
	}
	gd.sec = float64(s)
}
