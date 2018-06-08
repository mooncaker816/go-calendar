package calendar

import (
	"fmt"
	"log"
	"math"

	"github.com/mooncaker816/learnmeeus/v3/julian"
	"github.com/mooncaker816/learnmeeus/v3/moonphase"

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
	Months     []*LunarMonth // 一整个农历年月份1-12（含闰月）
	LeapN      int           // Months 中闰月的序号，-1表示没有闰月
	SpringFest float64       // 春节的儒略日数
	NatureYear               // 覆盖该农历年的2个农历自然年（冬至到冬至为一个自然年）
}

// NatureYear 两个农历自然年
type NatureYear struct {
	Terms  [2][]SQ       // 两个自然年包含的节气
	Shuoes [2][]SQ       // 两个自然年包含的朔日
	dzs    [3]SQ         // 划分两个自然年的三个冬至
	leap   [2]bool       // 两个自然年中是否有闰月
	months []*LunarMonth // 两个自然年的所有月份
}

// SQ 朔气JD，Avg是否为平朔气
type SQ struct {
	JD  float64
	Avg bool
}

// GenLunarYear generates Lunar Year
func GenLunarYear(year int) *LunarYear {
	ly := new(LunarYear)
	ly.Year = year
	ly.LeapN = -1
	ly.dz(year)        // 计算year-1年,year年和year+1年的冬至
	ly.solarTerms()    // 计算从year-1年冬至起到year+1年冬至的49个节气
	ly.moonShuoes()    // 计算从上一年冬至之前的一个合朔日起，连续15个合朔日
	ly.genLunarMonth() // 建月
	return ly
}

// 计算冬至
func (ly *LunarYear) dz(y int) {
	ly.dzs[0] = SQ{solstice.December2(y-1, earth), false}
	ly.dzs[1] = SQ{solstice.December2(y, earth), false}
	ly.dzs[2] = SQ{solstice.December2(y+1, earth), false}
}

// 计算节气
func (ly *LunarYear) solarTerms() {
	for i := 0; i < 2; i++ {
		j0 := ly.dzs[i].JD
		for j := 0; j < 25; j++ {
			t := calcTermI(j0, j)
			ly.Terms[i] = append(ly.Terms[i], SQ{t, false})
		}
	}
}

// 计算单个节气
func calcTermI(jde float64, i int) float64 {
	q := unit.Angle(1.5*math.Pi + float64(i)*math.Pi/12).Mod1() // 节气对应的黄经
	jde += 15.2184 * float64(i)
	for {
		λ, _, _ := solar.ApparentVSOP87(earth, jde)
		c := 58 * (q - λ).Sin()
		jde += c
		if math.Abs(c) < .000005 {
			break
		}
	}
	return jde
}

// 计算合朔日
func (ly *LunarYear) moonShuoes() {
	for i := 0; i < 2; i++ {
		dz0, dz1 := ly.dzs[i], ly.dzs[i+1]
		jde0 := dz0

		nm0 := newmoonI(jde0.JD, 0, 0)
		// fmt.Println("first cal newmoon:", jd2year(jde0), DT2SolarTime(dz0), DT2SolarTime(nm0))

		if !sLEq(nm0, jde0) { // nm0>jde0 获得离冬至最近的前一个朔日，当冬至和朔日重合时，默认朔在前
			// jde0 -= 29.5306
			// jde0 -= 29.5306 / 2
			jde0.JD -= 29.5306 - nm0.JD + jde0.JD
		}
		prevnm := 0.0
		for j := 0; j < 15; j++ { //计算第一个朔日（十一月初一）起的连续15个朔日
			shuo := newmoonI(jde0.JD, prevnm, j)
			shuo = shuoC(shuo, shuoCorrect) //按古历修正少数朔日
			prevnm = shuo.JD
			ly.Shuoes[i] = append(ly.Shuoes[i], shuo)
		}
		// 判断第14个朔日是否在第二个冬至之前（重合），是则说明两冬至之间多余12个农历月，需要安排闰月
		// if sInDZs(ly.Shuoes[i][13], dz0, dz1) {
		if sLEq(ly.Shuoes[i][13], dz1) {
			ly.leap[i] = true
		}
	}
	return
}

func jd2year(jd float64) float64 {
	year, m, d := julian.JDToCalendar(jd)
	z, f := math.Modf(d)
	yeardays := 365.
	if julian.LeapYearGregorian(year) {
		yeardays++
	}
	return float64(year) + (float64(julian.DayOfYearGregorian(year, m, int(z)))+f)/yeardays
}

func avgRange(jde float64) bool {
	return jde < 1947168-14 && jde >= 1457698-14
}

func newmoonI(jde, prevnm float64, i int) SQ {
	nmjd := jde + 29.5306*float64(i)
	if avgRange(nmjd) {
		return SQ{avgShuo(nmjd), true}
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
	return SQ{s, false}
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

func avgShuo(jde float64) float64 {
	s := 0.
	for i := 0; i < len(avgShuoTab); i += 2 {
		if jde+14 > avgShuoTab[i] {
			s = avgShuoTab[i] + avgShuoTab[i+1]*math.Floor((jde+14-avgShuoTab[i])/avgShuoTab[i+1])
			break
		}
	}
	s = math.Floor(s + 0.5)
	if s == 1683460 { //如果使用太初历计算-103年1月24日的朔日,结果得到的是23日,这里修正为24日(实历)。修正后仍不影响-103的无中置闰。如果使用秦汉历，得到的是24日，本行D不会被执行。
		s++
	}
	return s
}

// 建月，两个自然年之间的所有月份，包含一整个农历年月份
func (ly *LunarYear) genLunarMonth() {
	monthnum := [2]int{12, 12}
	leapI := [2]int{-1, -1}
	for i := 0; i < 2; i++ {
		if ly.leap[i] { // 如果该自然年有闰，则获取该闰月朔日的索引号
			leapI[i] = getLeapI(ly.Shuoes[i], ly.Terms[i])
			monthnum[i]++
		}
	}
	if leapI[0] >= 3 && leapI[0] <= 12 { // 闰农历年的1-10月
		ly.LeapN = i2monthseq(leapI[0] - 1)
	}
	if leapI[1] <= 2 && leapI[1] > 0 { //闰农历年的11，12月
		ly.LeapN = i2monthseq(leapI[1] - 1)
	}

	// 定农历年首（春节）
	var offset int
	ly.SpringFest, offset = getSpringFest(ly.Shuoes[0], leapI[0])

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
					lm.seq = i2monthseq(j)
					lm.leap = false
				case j == leapI[i]:
					lm.seq = i2monthseq(j - 1)
					lm.leap = true
				default:
					lm.seq = i2monthseq(j - 1)
					lm.leap = false
				}
			} else {
				lm.seq = i2monthseq(j)
				lm.leap = false
			}
			// 农历11，12月记为上一年
			if lm.seq >= 10 {
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
func getLeapI(shuoes, terms []SQ) int {
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
func getSpringFest(shuoes []SQ, leapI int) (float64, int) {
	springFest := shuoes[2]
	offset := 2
	if leapI != -1 && leapI <= 2 { //闰11或闰12月
		springFest = shuoes[3]
	}
	if springFest.Avg {
		return springFest.JD, offset
	}
	return jd2jdN(beijingTime(springFest)), offset
}

// 将朔表中月份序号映射到正常月数-1
func i2monthseq(i int) int {
	l := (i + jianyue) % 12
	if l == 0 {
		l = 12
	}
	return l - 1
}

// 将力学时转为北京时间
func beijingTime(sq SQ) float64 {
	if sq.Avg {
		return sq.JD
	}
	return sq.JD - deltat(sq.JD) + float64(8)/24
}

// 判断气相对于朔的关系
// 若朔=气，默认为朔在前
func sLEq(s, q SQ) bool {
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
func sInDZs(s, dz0, dz1 SQ) bool {
	st := jd2jdN(beijingTime(s))
	dz0t := jd2jdN(beijingTime(dz0))
	dz1t := jd2jdN(beijingTime(dz1))

	if st > dz0t && st <= dz1t {
		return true
	}
	return false
}

// 由于古历算法的局限性，少数朔日实际有误，此处仍按古历进行修正
func shuoC(shuo SQ, a []struct{ jdN, delta float64 }) SQ {
	if shuo.Avg {
		return shuo
	}
	key := jd2jdN(beijingTime(shuo))
	lo := 0
	hi := len(a) - 1
	for lo <= hi {
		mid := lo + (hi-lo)/2
		if key < a[mid].jdN {
			hi = mid - 1
		} else if key > a[mid].jdN {
			lo = mid + 1
		} else {
			shuo.JD += a[mid].delta
			return shuo
		}
	}
	return shuo
}

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
func (ly LunarYear) debug() {
	fmt.Println("年：", ly.Year)
	fmt.Println("闰月：", ly.LeapN+1)
	fmt.Println("春节：")
	fmt.Println(julian.JDToCalendar(ly.SpringFest))
	fmt.Println("xxxxxxxxxxxxxxxxxx")
	for _, m := range ly.Months {
		fmt.Println("月首：")
		fmt.Println(julian.JDToCalendar(m.d0))
		fmt.Println("月长：", m.dn)
		fmt.Println("闰：", m.leap)
		fmt.Println("月：", m.seq+1)
		fmt.Println("年：", m.year)
		fmt.Println("==============")
	}
	fmt.Println("两个自然年是否有闰：", ly.leap)
	fmt.Printf("ΔT≈%fs,寿星ΔT≈%fs\n", deltat(ly.dzs[0].JD)*86400, deltat2(ly.dzs[0].JD)*86400)
	for i, dz := range ly.dzs {
		fmt.Printf("冬至:%d %6.f %s\n", i, jd2jdN(beijingTime(dz)), DT2SolarTime(dz))
	}
	for _, shuo := range ly.Shuoes {
		for i, v := range shuo {
			fmt.Printf("朔:%d %6.f %s\n", i, jd2jdN(beijingTime(v)), DT2SolarTime(v))
		}
	}
	for _, term := range ly.Terms {
		for i := 0; i < len(term); i = i + 2 {
			fmt.Printf("气:%d %6.f %s\n", i/2, jd2jdN(beijingTime(term[i])), DT2SolarTime(term[i]))
		}
	}
	for _, m := range ly.months {
		fmt.Println("月首：")
		fmt.Println(julian.JDToCalendar(m.d0))
		fmt.Println("月长：", m.dn)
		fmt.Println("闰：", m.leap)
		fmt.Println("月：", m.seq+1)
		fmt.Println("年：", m.year)
		fmt.Println("==============")
	}
}
