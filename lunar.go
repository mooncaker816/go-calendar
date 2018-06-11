package calendar

import (
	"fmt"
	"log"
	"math"
	"sort"
	"strings"

	dt "github.com/mooncaker816/learnmeeus/v3/deltat"
	"github.com/mooncaker816/learnmeeus/v3/julian"
	"github.com/mooncaker816/learnmeeus/v3/moonphase"

	sexa "github.com/soniakeys/sexagesimal"
	"github.com/soniakeys/unit"

	pp "github.com/mooncaker816/learnmeeus/v3/planetposition"
	"github.com/mooncaker816/learnmeeus/v3/solar"
	"github.com/mooncaker816/learnmeeus/v3/solstice"
)

var earth *pp.V87Planet

const (
	jianyue = 11
)

func init() {
	var err error
	earth, err = pp.LoadPlanet(pp.Earth)
	if err != nil {
		log.Fatalf("can not load planet: %v", err)
	}
}

// LunarMonth 农历月
type LunarMonth struct {
	d0   float64 // 月首儒略日数
	dn   int     // 本月的天数
	seq  int     // 月份-1
	leap bool    // 闰月
	year int     // 该月所属的年份
}

// LunarYear 农历年
type LunarYear struct {
	Year       int           // 该农历年年份
	Period                   //
	Months     []*LunarMonth // 一整个农历年月份1-12（含闰月）
	LeapN      int           // Months 中闰月的序号，-1表示没有闰月
	SpringFest float64       // 春节的儒略日数
	NatureYear               // 覆盖该农历年的2个农历自然年（冬至到冬至为一个自然年）
	debug      bool
	debugyears []int
}

// NatureYear 两个农历自然年
type NatureYear struct {
	Terms      [2][]JDPlus   // 两个自然年包含的节气
	Shuoes     [2][]JDPlus   // 两个自然年包含的朔日
	dzs        [3]JDPlus     // 划分两个自然年的三个冬至
	leap       [2]bool       // 两个自然年中是否有闰月
	months     []*LunarMonth // 两个自然年的所有月份
	springFest []float64     // 两个自然年的所有月份
}

// JDPlus 朔气JD，Avg是否为平朔气
type JDPlus struct {
	JD  float64
	Avg bool //平朔平气，已考虑了ΔT和北京时差8h
}

func (jdp JDPlus) String() string {
	jd := jdp.JD
	if !jdp.Avg {
		ΔT := dt.Interp10A(jd)
		jd = jd - ΔT.Day() + float64(8)/24
	}

	y, m, day := julian.JDToCalendar(jd)
	dz, df := math.Modf(day)
	d := int(dz)
	tm := unit.TimeFromDay(df)
	return fmt.Sprintf("%d年%d月%d日 %s", y, m, d, sexa.FmtTime(tm))
}

// Option is a function to modify the LunarYear's options
type Option func(ly *LunarYear)

// Debug turns on debug mode,
// if years are specified, then only those years will be debuged,
// otherwise each year will be debuged.
func Debug(years ...int) Option {
	return func(ly *LunarYear) {
		ly.debug = true
		ly.debugyears = years
	}
}

// GenLunarYear generates Lunar Year
func GenLunarYear(year int, opts ...Option) *LunarYear {
	ly := new(LunarYear)
	for _, opt := range opts {
		opt(ly)
	}
	ly.getPeriod(year)
	ly.Year = year
	ly.LeapN = -1
	ly.dz()            // 计算自然年起算3个冬至日
	ly.solarTerms()    // 计算从这三个冬至间共49个节气
	ly.moonShuoes()    // 对每个冬至计算从冬至之前的一个合朔日起，连续15个合朔日
	ly.genLunarMonth() // 排月
	ly.details()
	return ly
}

// 计算冬至
func (ly *LunarYear) dz() {
	y := ly.GYear - 1
	for i := 0; i < 3; i++ {
		tmp := solstice.December2(y, earth)
		if avgQiRange(tmp) {
			ly.dzs[i] = avgSQ(tmp, 7, avgQiTab)
		} else {
			ly.dzs[i] = JDPlus{tmp, false}
		}
		ly.dzs[i] = qiC(ly.dzs[i])
		y++
	}
	// ly.dzs[0] = JDPlus{solstice.December2(y-1, earth), false}
	// ly.dzs[1] = JDPlus{solstice.December2(y, earth), false}
	// ly.dzs[2] = JDPlus{solstice.December2(y+1, earth), false}
}

// 计算节气
func (ly *LunarYear) solarTerms() {
	for i := 0; i < 2; i++ {
		j0 := ly.dzs[i].JD
		for j := 0; j < 25; j++ {
			t := calcTermI(j0, j)
			t = qiC(t)
			ly.Terms[i] = append(ly.Terms[i], t)
		}
	}
}

// 计算单个节气
func calcTermI(jde float64, i int) JDPlus {
	q := unit.Angle(1.5*math.Pi + float64(i)*math.Pi/12).Mod1() // 节气对应的黄经
	jde += 15.2184 * float64(i)
	if avgQiRange(jde) {
		return avgSQ(jde, 7, avgQiTab)
	}
	for {
		λ, _, _ := solar.ApparentVSOP87(earth, jde)
		c := 58 * (q - λ).Sin()
		jde += c
		if math.Abs(c) < .000005 {
			break
		}
	}
	return JDPlus{jde, false}
}

// 计算合朔日
func (ly *LunarYear) moonShuoes() {
	for i := 0; i < 2; i++ {
		dz0, dz1 := ly.dzs[i], ly.dzs[i+1]
		jde0 := dz0

		nm0 := newmoonI(jde0.JD, 0, 0)

		if !sLEq(nm0, jde0) { // nm0>jde0 获得离冬至最近的前一个朔日，当冬至和朔日重合时，默认朔在前
			// jde0 -= 29.5306
			// jde0 -= 29.5306 / 2
			jde0.JD -= 29.5306 - nm0.JD + jde0.JD
		}
		prevnm := 0.0
		for j := 0; j < 15; j++ { //计算第一个朔日（十一月初一）起的连续15个朔日
			shuo := newmoonI(jde0.JD, prevnm, j)
			prevnm = shuo.JD
			shuo = shuoC(shuo) //按古历修正少数朔日
			ly.Shuoes[i] = append(ly.Shuoes[i], shuo)
		}
		// 判断第14个朔日是否在第二个冬至之前（重合），是则说明两冬至之间多余12个农历月，需要安排闰月
		// if sInDZs(ly.Shuoes[i][13], dz0, dz1) {
		if ly.ZhiRun == RNoZQ && sLEq(ly.Shuoes[i][13], dz1) {
			ly.leap[i] = true
		}
	}
	return
}

func jd2year(jd float64) float64 {
	year, m, d := julian.JDToCalendar(jd)
	z, f := math.Modf(d)
	yeardays := 365.
	leap := LeapYear(year)
	if leap {
		yeardays++
	}
	return float64(year) + (float64(julian.DayOfYear(year, m, int(z), leap))+f)/yeardays
}

func avgShuoRange(jde float64) bool {
	return jde < 1947168-14 && jde >= 1457698-14
}

func avgQiRange(jde float64) bool {
	return jde < 2322147-7 && jde >= 1640650-7
}

func newmoonI(jde, prevnm float64, i int) JDPlus {
	nmjd := jde + 29.5306*float64(i)
	if avgShuoRange(nmjd) {
		return avgSQ(nmjd, 14, avgShuoTab)
	}
	y := jd2year(nmjd)
	s := moonphase.New(y)
	// 测试过程中碰到不同的日期计算出的最近新月相同，这时只要再加一天进行计算
	// 也有可能碰到计算出的最近新月并不是最近，这时要只减一天进行计算
	for prevnm > 0 && (s-prevnm > 40 || s == prevnm) {
		if s == prevnm {
			// fmt.Println("same as prevnm", nmjd, prevnm, s)
			nmjd++
		}
		if s-prevnm > 40 {
			// fmt.Println("miss a newmoon", nmjd, prevnm, s)
			if s > nmjd {
				nmjd--
			}
		}
		s = moonphase.New(jd2year(nmjd))
		// fmt.Println("after correction:", nmjd, prevnm, s)
	}
	return JDPlus{s, false}
}

var avgShuoTab = []float64{
	1457698.231017, 29.53067166, // -721-12-17 h=0.00032 古历·春秋
	1546082.512234, 29.53085106, // -479-12-11 h=0.00053 古历·战国
	1640640.735300, 29.53060000, // -221-10-31 h=0.01010 古历·秦汉
	1642472.151543, 29.53085439, // -216-11-04 h=0.00040 古历·秦汉

	1683430.509300, 29.53086148, // -104-12-25 h=0.00313 汉书·律历志(太初历)平气平朔
	1752148.041079, 29.53085097, //   85-02-13 h=0.00049 后汉书·律历志(四分历)
	1807665.420323, 29.53059851, //  237-02-12 h=0.00033 晋书·律历志(景初历)
	1883618.114100, 29.53060000, //  445-01-24 h=0.00030 宋书·律历志(何承天元嘉历)
	1907360.704700, 29.53060000, //  510-01-26 h=0.00030 宋书·律历志(祖冲之大明历)
	1936596.224900, 29.53060000, //  590-02-10 h=0.01010 随书·律历志(开皇历)
	1939135.675300, 29.53060000, //  597-01-24 h=0.00890 随书·律历志(大业历)
	1947168.00, //  619-01-21
}
var avgQiTab = []float64{ //气直线拟合参数
	1640650.479938, 15.21842500, // -221-11-09 h=0.01709 古历·秦汉
	1642476.703182, 15.21874996, // -216-11-09 h=0.01557 古历·秦汉

	1683430.515601, 15.218750011, // -104-12-25 h=0.01560 汉书·律历志(太初历)平气平朔 回归年=365.25000
	1752157.640664, 15.218749978, //   85-02-23 h=0.01559 后汉书·律历志(四分历) 回归年=365.25000
	1807675.003759, 15.218620279, //  237-02-22 h=0.00010 晋书·律历志(景初历) 回归年=365.24689
	1883627.765182, 15.218612292, //  445-02-03 h=0.00026 宋书·律历志(何承天元嘉历) 回归年=365.24670
	1907369.128100, 15.218449176, //  510-02-03 h=0.00027 宋书·律历志(祖冲之大明历) 回归年=365.24278
	1936603.140413, 15.218425000, //  590-02-17 h=0.00149 随书·律历志(开皇历) 回归年=365.24220
	1939145.524180, 15.218466998, //  597-02-03 h=0.00121 随书·律历志(大业历) 回归年=365.24321
	1947180.798300, 15.218524844, //  619-02-03 h=0.00052 新唐书·历志(戊寅元历)平气定朔 回归年=365.24460
	1964362.041824, 15.218533526, //  666-02-17 h=0.00059 新唐书·历志(麟德历) 回归年=365.24480
	1987372.340971, 15.218513908, //  729-02-16 h=0.00096 新唐书·历志(大衍历,至德历) 回归年=365.24433
	1999653.819126, 15.218530782, //  762-10-03 h=0.00093 新唐书·历志(五纪历) 回归年=365.24474
	2007445.469786, 15.218535181, //  784-02-01 h=0.00059 新唐书·历志(正元历,观象历) 回归年=365.24484
	2021324.917146, 15.218526248, //  822-02-01 h=0.00022 新唐书·历志(宣明历) 回归年=365.24463
	2047257.232342, 15.218519654, //  893-01-31 h=0.00015 新唐书·历志(崇玄历) 回归年=365.24447
	2070282.898213, 15.218425000, //  956-02-16 h=0.00149 旧五代·历志(钦天历) 回归年=365.24220
	2073204.872850, 15.218515221, //  964-02-16 h=0.00166 宋史·律历志(应天历) 回归年=365.24437
	2080144.500926, 15.218530782, //  983-02-16 h=0.00093 宋史·律历志(乾元历) 回归年=365.24474
	2086703.688963, 15.218523776, // 1001-01-31 h=0.00067 宋史·律历志(仪天历,崇天历) 回归年=365.24457
	2110033.182763, 15.218425000, // 1064-12-15 h=0.00669 宋史·律历志(明天历) 回归年=365.24220
	2111190.300888, 15.218425000, // 1068-02-15 h=0.00149 宋史·律历志(崇天历) 回归年=365.24220
	2113731.271005, 15.218515671, // 1075-01-30 h=0.00038 李锐补修(奉元历) 回归年=365.24438
	2120670.840263, 15.218425000, // 1094-01-30 h=0.00149 宋史·律历志 回归年=365.24220
	2123973.309063, 15.218425000, // 1103-02-14 h=0.00669 李锐补修(占天历) 回归年=365.24220
	2125068.997336, 15.218477932, // 1106-02-14 h=0.00056 宋史·律历志(纪元历) 回归年=365.24347
	2136026.312633, 15.218472436, // 1136-02-14 h=0.00088 宋史·律历志(统元历,乾道历,淳熙历) 回归年=365.24334
	2156099.495538, 15.218425000, // 1191-01-29 h=0.00149 宋史·律历志(会元历) 回归年=365.24220
	2159021.324663, 15.218425000, // 1199-01-29 h=0.00149 宋史·律历志(统天历) 回归年=365.24220
	2162308.575254, 15.218461742, // 1208-01-30 h=0.00146 宋史·律历志(开禧历) 回归年=365.24308
	2178485.706538, 15.218425000, // 1252-05-15 h=0.04606 淳祐历 回归年=365.24220
	2178759.662849, 15.218445786, // 1253-02-13 h=0.00231 会天历 回归年=365.24270
	2185334.020800, 15.218425000, // 1271-02-13 h=0.00520 宋史·律历志(成天历) 回归年=365.24220
	2187525.481425, 15.218425000, // 1277-02-12 h=0.00520 本天历 回归年=365.24220
	2188621.191481, 15.218437494, // 1280-02-13 h=0.00015 元史·历志(郭守敬授时历) 回归年=365.24250
	2322147.76, // 1645-09-21
}

func avgSQ(jde float64, delta float64, avgTab []float64) JDPlus {
	s := 0.
	i := 0
	for ; i < len(avgTab); i += 2 {
		if jde+delta < avgTab[i+2] {
			break
		}
	}
	s = avgTab[i] + avgTab[i+1]*math.Floor((jde+delta-avgTab[i])/avgTab[i+1])
	s = math.Floor(s + 0.5)
	if s == 1683460 { //如果使用太初历计算-103年1月24日的朔日,结果得到的是23日,这里修正为24日(实历)。修正后仍不影响-103的无中置闰。如果使用秦汉历，得到的是24日，本行D不会被执行。
		s++
	}
	return JDPlus{s, true}
}

// 建月，两个自然年之间的所有月份，包含一整个农历年月份
func (ly *LunarYear) genLunarMonth() {
	monthnum := [2]int{12, 12}
	leapI := [2]int{-1, -1}
	switch ly.ZhiRun {
	case R7in19st10, R7in19st1:
		if chkR7in19(ly.GYear) {
			if ly.YueJian == ZZ {
				leapI[0] = 12 //闰十三月
				ly.LeapN = 11
			}
			if ly.YueJian == YZ {
				leapI[0] = 11 //后九月
				ly.LeapN = 8
			}
			monthnum[0]++
			ly.leap[0] = true
		}
	default:
		for i := 0; i < 2; i++ {
			if ly.leap[i] { // 如果该自然年有闰，则获取该闰月朔日的索引号
				leapI[i] = getLeapI(ly.Shuoes[i], ly.Terms[i])
				monthnum[i]++
			}
		}
		switch ly.YueJian {
		case YZ:
			if leapI[0] >= 3 && leapI[0] <= 12 { // 闰农历年的1-10月
				ly.LeapN = i2monthseq(leapI[0]-1, YZ)
			}
			if leapI[1] <= 2 && leapI[1] > 0 { //闰农历年的11，12月
				ly.LeapN = i2monthseq(leapI[1]-1, YZ)
			}
		case CZ:
			if leapI[0] >= 2 && leapI[0] <= 12 {
				ly.LeapN = i2monthseq(leapI[0]-1, CZ)
			}
			if leapI[1] <= 1 && leapI[1] > 0 { //闰农历年的12月
				ly.LeapN = i2monthseq(leapI[1]-1, CZ)
			}
		case ZZ, ZZYY:
			if leapI[0] >= 1 && leapI[0] <= 12 {
				ly.LeapN = i2monthseq(leapI[0]-1, ZZ)
			}
		}
	}

	// 定农历年首（春节）
	var offset int
	for i := 0; i < 2; i++ {
		sf, tmp := getSpringFest(ly.Shuoes[i], leapI[i], ly.YueJian)
		ly.springFest = append(ly.springFest, sf)
		if i == 0 {
			ly.SpringFest, offset = sf, tmp
		}
	}

	yjNum := int(ly.YueJian)
	if ly.YueJian == ZZYY {
		yjNum = 0
	}
	// var ms []LunarMonth
	for i := 0; i < 2; i++ {
		for j := 0; j < monthnum[i]; j++ {
			var lm LunarMonth
			lm.d0 = jd2jdN(beijingTime(ly.Shuoes[i][j]))                // 月首儒略日数
			lm.dn = int(jd2jdN(beijingTime(ly.Shuoes[i][j+1])) - lm.d0) // 月长（天）
			lm.year = ly.Year + i                                       // 月所属年份
			if ly.leap[i] {
				switch {
				case j < leapI[i]:
					lm.seq = i2monthseq(j, ly.YueJian)
					lm.leap = false
				case j == leapI[i]:
					lm.seq = i2monthseq(j-1, ly.YueJian)
					lm.leap = true
				default:
					lm.seq = i2monthseq(j-1, ly.YueJian)
					lm.leap = false
				}
			} else {
				lm.seq = i2monthseq(j, ly.YueJian)
				lm.leap = false
			}
			// 春节之前
			// if lm.d0 < ly.SpringFest && ly.YueJian != ZZ && ly.YueJian != ZZYY {

			if lm.seq >= 12-yjNum {
				lm.year--
			}
			// ms = append(ms, lm)
			ly.months = append(ly.months, &lm)
		}
	}

	length := 12
	if ly.LeapN > -1 {
		length++
	}

	ly.Months = ly.months[offset : offset+length]

	return
}

// 中气与合朔日发生在同一天，是用“发生时刻的先后顺序确定某月是否包中气”还是用“日期来确定包含关系”。
// 从原理上说，这两种方法都是可行的，不过，传统上为了降低历算的精度要求，采用后者来判断一个月中是否包含中气，紫金历也是如此。
// 朔气同天，朔在前
func getLeapI(shuoes, terms []JDPlus) int {
	j := 0
	t := terms[0]
	for i := 0; i < 13; i++ {
		s0 := shuoes[i]
		s1 := shuoes[i+1]
		for !sLEq(s0, t) { // s0>t
			j += 2
			if j >= len(terms) {
				return i
			}
			t = terms[j]
		}
		if sLEq(s0, t) && !sLEq(s1, t) { // s0 <= t && s1 > t
			continue
		}
		return i // [1,12]
	}
	return -1
}

// 定农历年首
func getSpringFest(shuoes []JDPlus, leapI int, yj YueJian) (float64, int) {
	var springFest JDPlus
	var offset int
	switch yj {
	case ZZ, ZZYY:
		springFest = shuoes[0]
	case CZ:
		springFest = shuoes[1]
		offset = 1
		if leapI != -1 && leapI <= 1 { //闰子月（12月）
			springFest = shuoes[2]
			offset++
		}
	case YZ:
		springFest = shuoes[2]
		offset = 2
		if leapI != -1 && leapI <= 2 { //闰11或闰12月
			springFest = shuoes[3]
			offset++
		}
	}
	if springFest.Avg {
		return springFest.JD, offset
	}
	return jd2jdN(beijingTime(springFest)), offset
}

// 将朔表中月份序号映射到正常月数-1
func i2monthseq(i int, yj YueJian) int {
	j := 11
	switch yj {
	case ZZ:
		j = 1
	case CZ:
		j = 12
	}
	l := (i + j) % 12
	if l == 0 {
		l = 12
	}
	return l - 1
}

// 将力学时转为北京时间
func beijingTime(sq JDPlus) float64 {
	if sq.Avg {
		return sq.JD
	}
	return sq.JD - deltat(sq.JD) + float64(8)/24
}

// 判断气相对于朔的关系
// 若朔=气，默认为朔在前
func sLEq(s, q JDPlus) bool {
	st := jd2jdN(beijingTime(s))
	qt := jd2jdN(beijingTime(q))

	switch {
	case st <= qt:
		return true
	default:
		return false
	}
}

// 判断朔是否在两冬至之间
func sInDZs(s, dz0, dz1 JDPlus) bool {
	st := jd2jdN(beijingTime(s))
	dz0t := jd2jdN(beijingTime(dz0))
	dz1t := jd2jdN(beijingTime(dz1))

	if st > dz0t && st <= dz1t {
		return true
	}
	return false
}

// 由于古历算法的局限性，少数朔日实际有误，此处仍按古历进行修正
// func shuoC(shuo JDPlus, a []struct{ jdN, delta float64 }) JDPlus {
// 	if shuo.Avg {
// 		return shuo
// 	}
// 	key := jd2jdN(beijingTime(shuo))
// 	lo := 0
// 	hi := len(a) - 1
// 	for lo <= hi {
// 		mid := lo + (hi-lo)/2
// 		if key < a[mid].jdN {
// 			hi = mid - 1
// 		} else if key > a[mid].jdN {
// 			lo = mid + 1
// 		} else {
// 			shuo.JD += a[mid].delta
// 			return shuo
// 		}
// 	}
// 	return shuo
// }

func checkLY(ly *LunarYear, year int, jdN float64) *LunarYear {
	if ly == nil || jdN < jd2jdN(ly.dzs[0].JD) {
		ly = GenLunarYear(year)
	}
	if jdN >= jd2jdN(ly.dzs[2].JD) {
		ly = GenLunarYear(year + 1)
	}
	return ly
}

// debug testing only
type tag struct {
	key  string
	desc string
	JDPlus
}

func (ly LunarYear) details() {
	if !ly.debug {
		return
	}
	if len(ly.debugyears) == 0 { // debug mode for each year
		goto debug
	}
	for _, y := range ly.debugyears { // debug mode only for provided years
		if ly.Year == y {
			goto debug
		}
	}
	return
debug:
	fmt.Println("=====================DEBUG BEGIN========================")
	fmt.Println("=====================Year Information===================")
	fmt.Println("农历年：", ly.Year)
	fmt.Println("闰月：", ly.LeapN+1)
	fmt.Println("春节：", JDPlus{ly.SpringFest - 0.5, true})
	fmt.Println("=====================Month Information==================")
	for _, m := range ly.Months {
		fmt.Println("月首：", JDPlus{m.d0 - 0.5, true})
		fmt.Println("月长：", m.dn)
		fmt.Println("闰：", m.leap)
		fmt.Println("月：", m.seq+1)
		fmt.Println("年：", m.year)
		fmt.Println("---------------------------------------------------")
	}
	fmt.Println("=====================Other Information==================")
	fmt.Printf("ΔT≈%fs,寿星ΔT≈%fs\n", deltat(ly.dzs[0].JD)*86400, deltat2(ly.dzs[0].JD)*86400)
	fmt.Println("两个自然年是否有闰：", ly.leap)
	for i, dz := range ly.dzs {
		fmt.Printf("冬至:%d %6.f %s\n", i, jd2jdN(beijingTime(dz)), dz)
	}

	var combines [2][]tag
	for i, shuo := range ly.Shuoes {
		for j, s := range shuo {
			key := fmt.Sprintf("%7.f%d", jd2jdN(beijingTime(s)), 1)
			desc := fmt.Sprintf("朔:%d", j)
			combines[i] = append(combines[i], tag{key, desc, s})
		}
	}
	for i, term := range ly.Terms {
		for j, t := range term {
			if j%2 != 0 {
				continue
			}
			key := fmt.Sprintf("%7.f%d", jd2jdN(beijingTime(t)), 2)
			desc := fmt.Sprintf("气:%d", j)
			combines[i] = append(combines[i], tag{key, desc, t})
		}
	}
	for i := 0; i < 2; i++ {
		sort.Slice(combines[i], func(m, n int) bool { return combines[i][m].key < combines[i][n].key })
		for _, v := range combines[i] {
			if strings.HasPrefix(v.desc, "朔") {
				fmt.Println("----------------------------------------------")
			}
			fmt.Printf("%s %s %s\n", v.desc, v.key[0:len(v.key)-1], v.JDPlus)
		}
		fmt.Println()
	}

	for _, m := range ly.months {
		fmt.Println("月首：", JDPlus{m.d0 - 0.5, true})
		fmt.Println("月长：", m.dn)
		fmt.Println("闰：", m.leap)
		fmt.Println("月：", m.seq+1)
		fmt.Println("年：", m.year)
		fmt.Println("---------------------------------------------------")
	}
	fmt.Println("=====================DEBUG END========================")
}
