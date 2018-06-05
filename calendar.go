package calendar

import (
	"FunOfSinoGraph/src/ichang"
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

type Calendar struct {
}

type Year struct {
	Year   int
	Months []Month
	Leap   bool
}

type Month struct {
	D0    float64 //月首儒略日数
	Dn    int     //本月的天数
	Week0 int     //月首的星期
	WeekN int     //本月的总周数
	Num   int     //公历月份
	Days  []Day   //该月的日
}

type Day struct {
	// 公历信息
	Jd    float64 // 儒略日数,北京时12:00
	DN    int     // 所在公历月内日数
	MN    int     // 所在公历月
	YN    int     // 所在公历年
	Week  int     // 星期
	Weeki int     // 在本月中的周序号
	// 农历信息
	LDN int //所在农历月内日数
	// ob.cur_dz 距冬至的天数
	// ob.cur_xz 距夏至的天数
	// ob.cur_lq 距立秋的天数
	// ob.cur_mz 距芒种的天数
	// ob.cur_xs 距小暑的天数
	LMN    int  //农历月数
	LMDn   int  //农历月天数
	LMleap bool //闰月标志
	LYN    int  //农历年数，以春节为界
	LYSX   ichang.Shengxiao
	GZInfo
	// ob.Ltime2 纪时
	// ob.XiZ 星座
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

func genDay(jd float64, ly *LunarYear) Day {
	var day Day
	var d float64
	jdN := jd2jdN(jd)
	// 近似处理，精确到1毫秒，主要处理因截断导致的如59.99999秒在时辰交替点的判断出现的误差
	tm := julian.JDToTime(jd).Round(time.Millisecond)

	// 公历信息
	day.Jd = jdN
	day.YN, day.MN, d = julian.JDToCalendar(jdN)
	day.DN = int(d)
	mDay0Jd := julian.CalendarGregorianToJD(day.YN, day.MN, 1)
	mDay0W := julian.DayOfWeek(mDay0Jd)
	day.Week = julian.DayOfWeek(jdN)
	day.Weeki = int(math.Floor(float64(mDay0W+day.DN-1) / 7))
	// 农历信息
	prev := (*(ly.months))[0]
	for _, m := range *ly.months {
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
	return day
}

func (gz GZ) String() string {
	return gz.G.String() + gz.Z.String()
}

// 农历[狗年] 四月（大）三十
// 戊戌年 戊戌月 戊戌日
// 四柱：乙未 戊寅 丙辰 丁酉
func (d Day) String() string {
	var b strings.Builder
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

func time2sci(t time.Time) int {
	return ((t.Hour() + 1) / 2) % 12
}
