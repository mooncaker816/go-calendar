package calendar_test

import (
	"testing"

	. "github.com/mooncaker816/go-calendar"
)

var tt []struct {
	year  int
	month int
	day   float64
	hour  int
	min   int
	sec   int
	jd    JulianDate
}

func init() {
	tt = []struct {
		year  int
		month int
		day   float64
		hour  int
		min   int
		sec   int
		jd    JulianDate
	}{
		{year: 2000, month: 1, day: 1.5, jd: 2451545.0},
		{year: 1987, month: 1, day: 27.0, jd: 2446822.5},
		{year: 837, month: 4, day: 10.3, jd: 2026871.8},
		{year: -1001, month: 8, day: 17.9, jd: 1355671.4},
		{year: -4712, month: 1, day: 1.5, jd: 0.0},
		{year: -4712, month: 1, day: 1.5, hour: 10, jd: 0.0}, // hour 10 will be ignored
	}
}

type pair struct {
	gd GregDate
	jd JulianDate
}

func creteCases() []pair {
	var cases []pair
	for _, tc := range tt {
		gd := new(GregDate)
		gd.SetYear(tc.year)
		gd.SetMonth(tc.month)
		gd.SetDay(tc.day)
		gd.SetHour(tc.hour)
		gd.SetMin(tc.min)
		gd.SetSec(tc.sec)
		jd := JulianDate(tc.jd)
		cases = append(cases, pair{*gd, jd})
	}
	return cases
}
func TestToJulianDate(t *testing.T) {
	tcs := creteCases()
	for _, tc := range tcs {
		if v := tc.gd.ToJulianDate(); v != tc.jd {
			t.Errorf("%v ToJulianDate() got %v, want %v\n", tc.gd, v, tc.jd)
		}
	}
}

func TestToGergDate(t *testing.T) {
	tcs := creteCases()
	for _, tc := range tcs {
		if v := tc.jd.ToGregDate(); v.String() != tc.gd.String() {
			t.Errorf("%v ToGregDate() got %v, want %v\n", tc.jd, v, tc.gd)
			// fmt.Println(tc.gd.year, v.year)
			// fmt.Println(tc.gd.month, v.month)
			// fmt.Println(tc.gd.day, v.day)
			// fmt.Println(tc.gd.hour, v.hour)
			// fmt.Println(tc.gd.min, v.min)
			// fmt.Println(tc.gd.sec, v.sec)
			// fmt.Println(tc.gd.decimalday, v.decimalday)
		}
	}
}
