package calendar

import (
	"testing"
)

var tt = []struct {
	gd GregDate
	jd JulianDate
}{
	{gd: GregDate{year: 2000, month: 1, day: 1.5},
		jd: 2451545.0,
	},
	{gd: GregDate{year: 1987, month: 1, day: 27.0},
		jd: 2446822.5,
	},
	{gd: GregDate{year: 837, month: 4, day: 10.3},
		jd: 2026871.8,
	},
	{gd: GregDate{year: -1001, month: 8, day: 17.9},
		jd: 1355671.4,
	},
	{gd: GregDate{year: -4712, month: 1, day: 1.5},
		jd: 0.0,
	},
}

func TestToJulianDate(t *testing.T) {
	for _, tc := range tt {
		if v := tc.gd.ToJulianDate(); v != tc.jd {
			t.Errorf("%v ToJulianDate() got %v, want %v\n", tc.gd, v, tc.jd)
		}
	}
}

func TestToGergDate(t *testing.T) {
	for _, tc := range tt {
		if v := tc.jd.ToGregDate(); v.String() != tc.gd.String() {
			t.Errorf("%v ToGregDate() got %v, want %v\n", tc.jd, v, tc.gd)
		}
	}
}
