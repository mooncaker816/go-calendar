package calendar

import (
	"fmt"
	"math"

	dt "github.com/mooncaker816/learnmeeus/v3/deltat"

	"github.com/mooncaker816/learnmeeus/v3/julian"
	sexa "github.com/soniakeys/sexagesimal"
	"github.com/soniakeys/unit"
)

// SolarTime is a local UTC contains year,month,date and time
type SolarTime struct {
	Y, M, D int
	T       unit.Time
}

// DT2SolarTime converts DT to local time
func DT2SolarTime(jde float64) SolarTime {
	// log.Println(julian.JDToTime(jde))
	var st SolarTime
	ΔT := dt.Interp10A(jde)
	jd := jde - ΔT.Day() // UT
	// jd0h := math.Floor(jd+0.5) - 0.5 //当天0点 jd
	var day float64
	st.Y, st.M, day = julian.JDToCalendar(jd + float64(8)/24)
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
	default:
		//  year >= 1600 && year <= y+50:
		return dt.Interp10A(jde).Day()
		// default:
		// 	return dt.PolyAfter2000(jd2year(jde)).Day()
	}
}
