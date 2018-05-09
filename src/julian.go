package calendar

import (
	"fmt"
	"math"
)

// JulianDate type
type JulianDate float64

func (jd JulianDate) String() string {
	return jd.ToGregDate().String()
}

// ToGregDate convert JulianDate to GregDate
func (jd JulianDate) ToGregDate() GregDate {
	var gd GregDate
	d := math.Floor(float64(jd + 0.5))
	f := float64(jd) + 0.5 - d
	if d >= 2299161 {
		c := math.Floor((d - 1867216.25) / 36524.25)
		d += 1 + c - math.Floor(c/4)
	}
	d += 1524
	gd.year = math.Floor((d - 122.1) / 365.25)
	d -= math.Floor(365.25 * gd.year)
	gd.month = math.Floor(d / 30.601)
	d -= math.Floor(30.601 * gd.month)
	gd.day = d
	if gd.month > 13 {
		gd.month -= 13
		gd.year -= 4715
	} else {
		gd.month--
		gd.year -= 4716
	}
	// hour
	f *= 24
	gd.hour = math.Floor(f)
	// minute
	f -= gd.hour
	f *= 60
	gd.min = math.Floor(f)
	// second
	f -= gd.min
	f *= 60
	gd.sec = f
	return gd
}

// TimeStr extract time info from Julian date, such as 12:00:00
func (jd JulianDate) TimeStr() string {
	t := float64(jd + 0.5)
	t = (t - math.Floor(t))
	s := math.Floor(t*86400 + 0.5)
	h := math.Floor(s / 3600)
	s -= h * 3600
	m := math.Floor(s / 60)
	s -= m * 60
	return fmt.Sprintf("%02d:%02d:%02d", int(h), int(m), int(s))
}

// GetWeek returns the week day of that Julian date
func (jd JulianDate) GetWeek() int {
	return int(math.Floor(float64(jd)+1.5+7000000)) % 7
}

// get julian date of year y, month m, the nth weekday w
func nnweek(y, m, n, w float64) JulianDate {
	jd := g2j(y, m, 1.5)                     //julian date of the first day of month m year y
	w0 := int(jd+1+7000000) % 7              // the week day of the first day of that month
	r := float64(jd) - float64(w0) + 7*n + w //jd-w0+7*n是和n个星期0,起算下本月第一行的星期日(可能落在上一月)。加w后为第n个星期w
	if w >= float64(w0) {                    //第1个星期w可能落在上个月,造成多算1周,所以考虑减1周
		r -= 7
	}
	if n == 5 {
		m++
		if m > 12 { //下个月
			m = 1
			y++
		}
		if r >= float64(g2j(y, m, 1.5)) { //r跑到下个月则减1周
			r -= 7
		}
	}
	return JulianDate(r)
}
