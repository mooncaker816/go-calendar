package calendar

import (
	"FunOfSinoGraph/src/ichang"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/mooncaker816/learnmeeus/v3/julian"
)

// 干支推算起始点
const (
	J1984lc  = 2445735 // 1984-2-4 立春 甲子年（鼠年）
	J1998dx  = 2451155 // 1998-12-7 大雪 甲子月
	J2000jzr = 2451551 // 2000-1-7 甲子日
)

// Year contains 1 Gregorian year's calendar info including Lunar info
type Year struct {
	Num    int
	Months []Month
	Leap   bool
}

// Month contains 1 Gregorian month's calendar info including Lunar info
type Month struct {
	Num   int     //公历月份
	D0    float64 //月首儒略日数
	Dn    int     //本月的天数
	Week0 int     //月首的星期
	WeekN int     //本月的总周数
	Days  []Day   //该月的日
}

// Day contains 1 Day's calendar info including Lunar info
type Day struct {
	// 公历信息
	Jd    float64 // 儒略日数,北京时12:00
	DN    int     // 所在公历月内日数
	MN    int     // 所在公历月
	YN    int     // 所在公历年
	Week  int     // 星期
	Weeki int     // 在本月中的周序号
	XZ    int     // 星座序号
	// 农历信息
	LDN    int  //所在农历月内日数
	LMN    int  //农历月数
	LMDn   int  //农历月天数
	LMleap bool //闰月标志
	LYN    int  //农历年数，以春节为界
	LYSX   ichang.Shengxiao
	GZInfo
}

// GZInfo 干支信息
type GZInfo struct {
	LYGZ0 GZ //年干支，以春节为界，用于纪年
	LMGZ0 GZ //月干支，以朔日为界，用于纪月
	LYGZ1 GZ //年干支，以立春为界，用于四柱
	LMGZ1 GZ //月干支，以节为界，用于四柱
	LDGZ  GZ //日干支
	LTGZ  GZ //时干支
}

// GZ 干支组合
type GZ struct {
	G ichang.Tiangan
	Z ichang.Dizhi
}

// GenDay generates the details for a specific JD
func genDay(jd float64, ly *LunarYear) Day {
	var day Day
	jdN := jd2jdN(jd)
	// 近似处理，精确到1毫秒，主要处理因截断导致的如59.99999秒在时辰交替点的判断出现的误差
	tm := julian.JDToTime(jd).Round(time.Millisecond)
	_ = tm
	// 公历信息
	var d float64
	day.Jd = jdN
	day.YN, day.MN, d = julian.JDToCalendar(jdN)
	day.DN = int(d)
	mDay0Jd := julian.CalendarGregorianToJD(day.YN, day.MN, 1)
	mDay0W := julian.DayOfWeek(mDay0Jd)
	day.Week = julian.DayOfWeek(jdN)
	day.Weeki = int(math.Floor(float64(mDay0W+day.DN-1) / 7))

	ly = checkLY(ly, day.YN, jdN)

	// 农历信息
	prev := ly.months[0]
	for _, m := range ly.months {
		if jdN < m.d0 {
			break
		}
		prev = m
	}
	day.LDN = int(jdN-prev.d0) + 1
	day.LMN = prev.seq + 1
	day.LMDn = prev.dn
	day.LMleap = prev.leap
	day.LYN = prev.year

	lc := jd2jdN(beijingTime(ly.Terms[0][3])) // 立春
	sf := ly.SpringFest                       // 春节
	// 年干支，春节为界
	dCnt := sf - J1984lc // 计算日所在农历自然年的春节与1984年平均春节(立春附近)相差天数估计
	yCnt := math.Floor(dCnt/365.2422 + 0.5)
	if jdN < sf {
		yCnt--
	}
	g, z := mod(int(yCnt), 10), mod(int(yCnt)+10, 12)
	day.LYGZ0 = GZ{ichang.Tiangan(g), ichang.Dizhi(z)}
	day.LYSX = ichang.Shengxiao(mod(int(yCnt), 12))
	// 年干支，立春为界
	dCnt = lc - J1984lc // 计算日所在农历自然年的立春距离1984年2月4日立春的天数
	yCnt = math.Floor(dCnt/365.2422 + 0.5)
	if jdN < lc {
		yCnt--
	}
	g, z = mod(int(yCnt), 10), mod(int(yCnt)+10, 12)
	day.LYGZ1 = GZ{ichang.Tiangan(g), ichang.Dizhi(z)}

	dz := jd2jdN(beijingTime(ly.Terms[0][0]))
	xz := jd2jdN(beijingTime(ly.Terms[0][12]))
	yCnt = math.Floor((xz - J1998dx) / 365.2422) // 用夏至点算离1998年12月7(大雪)的完整年数，确保不会有误差
	ymCnt := yCnt * 12                           // 从1998年12月7(大雪)到计算日前一个大雪的累计月数

	// 月干支，朔为界
	mCnt := int(ymCnt) + mod(day.LMN+1, 12)
	g, z = mod(mCnt, 10), mod(mCnt+10, 12)
	day.LMGZ0 = GZ{ichang.Tiangan(g), ichang.Dizhi(z)}
	// 月干支，节为界
	mk := int(math.Floor((day.Jd - dz) / 30.43685))
	// fmt.Println(day.Jd, dz, "->", mk)
	if mk < 12 && day.Jd >= jd2jdN(beijingTime(ly.Terms[0][2*mk+1])) {
		mk++ //相对计算日前一个大雪的月数计算,mk的取值范围0-12
	}
	mCnt = int(ymCnt) + mk
	g, z = mod(mCnt, 10), mod(mCnt+10, 12)
	day.LMGZ1 = GZ{ichang.Tiangan(g), ichang.Dizhi(z)}
	// 日干支,2000年1月7日起算
	dCnt = day.Jd - J2000jzr
	g, z = mod(int(dCnt), 10), mod(int(dCnt)+10, 12)
	day.LDGZ = GZ{ichang.Tiangan(g), ichang.Dizhi(z)}
	// 时干支，日上起时
	// 甲己还加甲，乙庚丙做初。
	// 丙辛从戊起，丁壬庚子居。
	// 戊癸何处去？壬子是真途。
	scI := time2sci(tm)
	g, z = mod(mod(g, 5)*2+scI, 10), mod(scI+10, 12)
	day.LTGZ = GZ{ichang.Tiangan(g), ichang.Dizhi(z)}
	xzI := int(math.Floor((jdN - dz - 15) / 30.43685))
	if xzI < 11 && jdN >= jd2jdN(beijingTime(ly.Terms[0][2*xzI+2])) {
		xzI++
	}
	day.XZ = (xzI + 12) % 12
	return day
}

func (gz GZ) String() string {
	return gz.G.String() + gz.Z.String()
}

// 公元2000年1月1日
// 星期六 摩羯座
// JD 2451545
// 农历[狗年] 四月（大）三十
// 甲子年 甲子月 甲子日
// 四柱：甲子 甲子 甲子 甲子
func (d Day) String() string {
	var b strings.Builder
	y := d.YN
	yh := "公元"
	if y <= 0 {
		yh += "前"
		y = -y + 1
	}
	b.WriteString(fmt.Sprintf("%s%d年%d月%d日\n", yh, y, d.MN, d.DN))
	b.WriteString("星期" + weekName[d.Week] + " " + xzName[d.XZ] + "座" + "\n")
	b.WriteString(fmt.Sprintf("JD %d\n", int(d.Jd)))
	b.WriteString("农历【" + d.LYSX.String() + "】")
	leap := ""
	if d.LMleap {
		leap = "闰"
	}
	size := "小"
	if d.LMDn > 29 {
		size = "大"
	}
	b.WriteString(leap + monthName[d.LMN-1] + "（" + size + "）" + dayName[d.LDN-1] + "\n")
	b.WriteString(d.LYGZ0.String() + "年 " + d.LMGZ0.String() + "月 " + d.LDGZ.String() + "日\n")
	b.WriteString("四柱：" + d.LYGZ1.String() + " " + d.LMGZ1.String() + " " + d.LDGZ.String() + " " + d.LTGZ.String())
	return b.String()
}

func (m Month) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("📅%13s%d月\n", " ", m.Num))
	b.WriteString("   日  一  二  三  四  五  六\n")

	k := 1
	idx := 0
	cnt := 7 - m.Week0
Loop:
	for i := 0; i < m.WeekN; i++ {
		if i > 0 {
			cnt = 7
		}
		b.WriteString("☀️  ")
		if i == 0 {
			for j := 0; j < m.Week0; j++ {
				b.WriteString(fmt.Sprintf("%4s", " "))
			}
		}
		for j := 0; j < cnt; j++ {
			width := 2
			if k < 10 && j == cnt-1 {
				width = 1
			}
			b.WriteString(fmt.Sprintf("%-*d", width, k)) //左对齐
			k++
			if k > m.Dn {
				// b.WriteString("\n")
				break
			}
			if j == cnt-1 {
				continue
			}
			b.WriteString(fmt.Sprintf("%2s", " "))
		}
		b.WriteString("\n")
		b.WriteString("🌛  ")
		if i == 0 {
			for j := 0; j < m.Week0; j++ {
				b.WriteString(fmt.Sprintf("%4s", " "))
			}
		}
		for j := 0; j < cnt; j++ {
			d := m.Days[idx]
			switch {
			case d.LDN == 1:
				b.WriteString(yueNames[d.LMN-1])
				if d.LMleap {
					b.WriteString("®")
				}
			case d.LDN > 1 && d.LDN < 10 && (j == cnt-1 || idx == m.Dn-1):
				b.WriteString(fmt.Sprintf("%-d", d.LDN)) //左对齐
			default:
				b.WriteString(fmt.Sprintf("%-2d", d.LDN)) //左对齐
			}
			idx++
			if idx > m.Dn-1 {
				// b.WriteString("\n")
				break Loop
			}
			if j == cnt-1 {
				continue
			}
			if d.LDN == 1 && d.LMleap {
				b.WriteString(fmt.Sprintf("%s", " "))
			} else {
				b.WriteString("  ")
			}
		}
		b.WriteString("\n")
	}

	return b.String()
}

func (y Year) String() string {
	var b strings.Builder
	leap := ""
	if y.Leap {
		leap = "（闰）"
	}
	b.WriteString(fmt.Sprintf("%13s%d年%s\n", " ", y.Num, leap))
	for i := 0; i < 12; i++ {
		b.WriteString(y.Months[i].String())
		b.WriteString("\n")
	}
	return b.String()
}

func time2sci(t time.Time) int {
	return ((t.Hour() + 1) / 2) % 12
}

// DayCalendar 日历, 单独调用时ly可置nil，ly只是为了方便需要多次调用（如建月历）的时候无需多次建立农历
func DayCalendar(y, m, d int, AD bool, ly *LunarYear) (Day, error) {
	var day Day
	if y <= 0 {
		return day, errors.New("year should be positive num")
	}
	if !AD {
		y = -y + 1
	}
	if m <= 0 || m > 12 {
		return day, errors.New("month number is not valid")
	}
	cnt := monthDayCnt[m-1]
	if m == 2 && julian.LeapYearGregorian(y) {
		cnt++
	}
	if d > cnt {
		return day, errors.New("invalid day number for this month")
	}
	jd00 := jd2jd00(julian.CalendarGregorianToJD(y, m, float64(d)))
	jd := jd00 + float64(time.Now().Hour())/24

	ly = checkLY(ly, y, jd00-0.5)

	day = genDay(jd, ly)
	return day, nil
}

// MonthCalendar 月历,单独调用时ly可置nil，ly只是为了方便需要多次调用（如建年历）的时候无需多次建立农历
func MonthCalendar(y, m int, AD bool, ly *LunarYear) (Month, error) {
	var month Month
	if y <= 0 {
		return month, errors.New("year should be positive num")
	}
	if !AD {
		y = -y + 1
	}
	if m <= 0 || m > 12 {
		return month, errors.New("month number is not valid")
	}
	jdN0 := julian.CalendarGregorianToJD(y, m, 1.5)
	month.Num = m   //公历月份
	month.D0 = jdN0 //月首儒略日数
	cnt := monthDayCnt[m-1]
	if m == 2 && julian.LeapYearGregorian(y) {
		cnt++
	}
	month.Dn = cnt                          //本月的天数
	month.Week0 = julian.DayOfWeek(jdN0)    //月首的星期
	month.WeekN = (month.Week0+cnt-1)/7 + 1 //本月的总周数
	h := time.Now().Hour()
	jd := jd2jd00(jdN0) + float64(h)/24
	days := make([]Day, cnt)
	ly = checkLY(ly, y, jdN0)
	for i := 0; i < cnt; i++ {
		days[i] = genDay(jd+float64(i), ly)
	}
	month.Days = days
	return month, nil
}

// YearCalendar 年历
func YearCalendar(y int, AD bool) (Year, error) {
	var year Year
	yN := y
	if y <= 0 {
		return year, errors.New("year should be positive num")
	}
	if !AD {
		yN = -y + 1
	}
	year.Num = yN
	year.Leap = julian.LeapYearGregorian(y)
	year.Months = make([]Month, 12)
	ly := GenLunarYear(yN)
	for i := 0; i < 12; i++ {
		m, err := MonthCalendar(y, i+1, AD, ly)
		if err != nil {
			return year, err
		}
		year.Months[i] = m
	}
	return year, nil
}
