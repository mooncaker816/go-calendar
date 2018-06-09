package calendar

import (
	"fmt"
	"math"

	"github.com/mooncaker816/learnmeeus/v3/base"

	dt "github.com/mooncaker816/learnmeeus/v3/deltat"

	"github.com/mooncaker816/learnmeeus/v3/julian"
	sexa "github.com/soniakeys/sexagesimal"
	"github.com/soniakeys/unit"
)

var dtat = []float64{ // TD - UT1 计算表
	-4000, 108371.7, -13036.80, 392.000, 0.0000,
	-500, 17201.0, -627.82, 16.170, -0.3413,
	-150, 12200.6, -346.41, 5.403, -0.1593,
	150, 9113.8, -328.13, -1.647, 0.0377,
	500, 5707.5, -391.41, 0.915, 0.3145,
	900, 2203.4, -283.45, 13.034, -0.1778,
	1300, 490.1, -57.35, 2.085, -0.0072,
	1600, 120.0, -9.81, -1.532, 0.1403,
	1700, 10.2, -0.91, 0.510, -0.0370,
	1800, 13.4, -0.72, 0.202, -0.0193,
	1830, 7.8, -1.81, 0.416, -0.0247,
	1860, 8.3, -0.13, -0.406, 0.0292,
	1880, -5.4, 0.32, -0.183, 0.0173,
	1900, -2.3, 2.06, 0.169, -0.0135,
	1920, 21.2, 1.69, -0.304, 0.0167,
	1940, 24.2, 1.22, -0.064, 0.0031,
	1960, 33.2, 0.51, 0.231, -0.0109,
	1980, 51.0, 1.29, -0.026, 0.0032,
	2000, 63.87, 0.1, 0, 0,
	2005, 64.7, 0.4, 0, 0, //一次项记为x,则 10x=0.4秒/年*(2015-2005),解得x=0.4
	2015, 69,
}

// SolarTime is a local UTC contains year,month,date and time
type SolarTime struct {
	Y, M, D int
	T       unit.Time
}

// DT2SolarTime converts DT to local time
func DT2SolarTime(sq JDPlus) SolarTime {
	// log.Println(julian.JDToTime(jde))
	var st SolarTime
	jd := sq.JD
	if !sq.Avg {
		ΔT := dt.Interp10A(sq.JD)
		jd = sq.JD - ΔT.Day() + float64(8)/24
	}
	// jd0h := math.Floor(jd+0.5) - 0.5 //当天0点 jd
	var day float64
	st.Y, st.M, day = julian.JDToCalendar(jd)
	dz, f := math.Modf(day)
	st.D = int(dz)
	st.T = unit.TimeFromDay(f)
	return st
}

func (st SolarTime) String() string {
	return fmt.Sprintf("%d年%d月%d日 %s", st.Y, st.M, st.D, sexa.FmtTime(st.T))
}

func deltat(jde float64) float64 {
	// y, _, _ := julian.JDToCalendar(julian.TimeToJD(time.Now()))
	year, _, _ := julian.JDToCalendar(jde)
	switch {
	case year < 948:
		return dt.PolyBefore948(jd2year(jde)).Day()
	case year >= 948 && year < 1600:
		return dt.Poly948to1600(jd2year(jde)).Day()
	case year >= 1600 && year <= 2100:
		return dt.Interp10A(jde).Day()
	default:
		// return dt.PolyAfter2000(jd2year(jde)).Day()
		return unit.Time(float64(31*(year-1820)*(year-1820))/10000 - float64(20)).Day()
	}
}

func mod(x, y int) int {
	m := x % y
	if m < 0 {
		m += y
	}
	return m
}

// 将 jd 转换成同一天12点的儒略日
func jd2jdN(jd float64) float64 {
	return math.Floor(jd + 0.5)
}

// 将 jd 转换成同一天0点的儒略日
func jd2jd00(jd float64) float64 {
	return jd2jdN(jd) - 0.5
}

// 寿星ΔT
func deltat2(jde float64) float64 {
	year, _, _ := julian.JDToCalendar(jde)
	y := float64(year)
	y0 := dtat[len(dtat)-2]
	t0 := dtat[len(dtat)-1]
	if y >= y0 {
		jsd := 31.
		if y > y0+100 {
			return unit.Time(dtAfter100(y, jsd)).Day()
		}
		v := dtAfter100(y, jsd)        //二次曲线外推
		dv := dtAfter100(y0, jsd) - t0 //ye年的二次外推与te的差
		return unit.Time(v - dv*(y0+100-y)/100).Day()
	}
	i := 0
	for ; i < len(dtat); i += 5 {
		if y < dtat[i+5] {
			break
		}
	}
	t1 := (y - dtat[i]) / (dtat[i+5] - dtat[i]) * 10
	return unit.Time(base.Horner(t1, dtat[i+1:i+5]...)).Day()
}

func dtAfter100(y, jsd float64) float64 {
	dy := (y - 1820) / 100
	return -20 + jsd*dy*dy
}
