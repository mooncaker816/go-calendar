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

var monthName = []string{"正月", "二月", "三月", "四月", "五月", "六月", "七月", "八月", "九月", "十月", "冬月", "腊月"}
var dayName = []string{
	"初一", "初二", "初三", "初四", "初五", "初六", "初七", "初八", "初九", "初十",
	"十一", "十二", "十三", "十四", "十五", "十六", "十七", "十八", "十九", "二十",
	"廿一", "廿二", "廿三", "廿四", "廿五", "廿六", "廿七", "廿八", "廿九", "三十",
}

func init() {
	var err error
	earth, err = pp.LoadPlanet(pp.Earth)
	if err != nil {
		log.Fatalf("can not load planet: %v", err)
	}
}

// LunarYMD contains info of Lunar Year + Month + Day
type LunarYMD struct {
	Y, M, D int
	Leap    bool
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
	Months     *[]LunarMonth // 一整个农历年月份1-12（含闰月）
	LeapN      int           // Months 中闰月的序号，-1表示没有闰月
	SpringFest float64       // 春节的儒略日数
	NatureYear               // 覆盖该农历年的2个农历自然年（冬至到冬至为一个自然年）
}

// NatureYear 两个农历自然年
type NatureYear struct {
	Terms  [2][]float64 // 两个自然年包含的节气
	Shuoes [2][]float64 // 两个自然年包含的朔日
	// shuoCnt [2]int        // 两个自然年包含的朔日的个数
	dzs    [3]float64    // 划分两个自然年的三个冬至
	leap   [2]bool       // 两个自然年中是否有闰月
	months *[]LunarMonth // 两个自然年的所有月份
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
	// ly.Stat()
	return ly
}

// 计算冬至
func (ly *LunarYear) dz(y int) {
	ly.dzs[0] = solstice.December2(y-1, earth)
	ly.dzs[1] = solstice.December2(y, earth)
	ly.dzs[2] = solstice.December2(y+1, earth)
}

// 计算节气
func (ly *LunarYear) solarTerms() {
	for i := 0; i < 2; i++ {
		j0 := ly.dzs[i]
		for j := 0; j < 25; j++ {
			t := calcTermI(j0, j)
			ly.Terms[i] = append(ly.Terms[i], t)
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
		jd0 := dz0
		nm0 := newmoonI(jd0, 0, 0)
		nm0 = shuoC(nm0, shuoCorrect)

		// ly.shuoCnt[i] = 0
		if !sLEq(nm0, jd0) { // nm0>jd0
			// if math.Floor(nm0+0.5) > math.Floor(beijingTime(jd0)+0.5) {
			jd0 -= 29.5306
			// ly.shuoCnt[i]--
		}
		prevnm := 0.0
		for j := 0; j < 15; j++ {
			shuo := newmoonI(jd0, prevnm, j)
			shuo = shuoC(shuo, shuoCorrect)
			prevnm = shuo
			// fmt.Println(DT2SolarTime(shuo))
			// shuop := math.Floor(shuo + 0.5)
			// if shuo >= jd0 && shuo < jd1 {
			// if sLEq(shuo, jd0) >= 0 && sq(shuo, jd1) == -1 { // shuo>=jd0 && shuo<jd1
			// if sInDZs(shuo, dz0, dz1) {
			// 	// if shuop >= math.Floor(beijingTime(jd0)+0.5) && shuop < math.Floor(beijingTime(jd1)+0.5) {
			// 	ly.shuoCnt[i]++
			// }
			ly.Shuoes[i] = append(ly.Shuoes[i], shuo)
			// fmt.Println(ly.Shuoes[i][0])
		}
		// if ly.shuoCnt[i] >= 13 { //冬至之间是否有13个朔日
		if sInDZs(ly.Shuoes[i][13], dz0, dz1) {
			ly.leap[i] = true
		}
	}
	return
}

func jd2year(jd float64) float64 {
	year, m, d := julian.JDToCalendar(jd)
	yeardays := 365.
	if julian.LeapYearGregorian(year) {
		yeardays++
	}
	return float64(year) + float64(julian.DayOfYearGregorian(year, m, int(d)))/yeardays
}

func newmoonI(jd, prevnm float64, i int) float64 {
	nmjd := jd + 29.5306*float64(i)
	y := jd2year(nmjd)
	s := moonphase.New(y)
	if s == prevnm { //可能前后两次计算出的最近新月是相同的，这时就要再加一天进行计算
		s = moonphase.New(jd2year(nmjd + 1))
	}
	// fmt.Println(jd, i, "->", s)
	return s
}

// 建月，两个自然年之间的所有月份，包含一整个农历年月份
func (ly *LunarYear) genLunarMonth() {
	var ms []LunarMonth
	monthnum := []int{12, 12}
	var leapI [2]int
	leapI[0], leapI[1] = -1, -1
	// fmt.Println(ly.leap[0], ly.leap[1])
	for i := 0; i < 2; i++ {
		if ly.leap[i] {
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

	// fmt.Println(ly.leap[0], ly.leap[1])
	// fmt.Println(leapI[0], leapI[1])
	for i := 0; i < 2; i++ {
		for j := 0; j < monthnum[i]; j++ {
			var lm LunarMonth
			lm.d0 = math.Floor(beijingTime(ly.Shuoes[i][j]) + 0.5) //儒略日数
			lm.dn = int(math.Floor(beijingTime(ly.Shuoes[i][j+1])+0.5) - lm.d0)
			lm.year = ly.Year + i
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
			if lm.seq >= 10 {
				lm.year--
			}
			// year, _, _ := julian.JDToCalendar(ly.SpringFest)
			// if lm.d0 < ly.SpringFest && ly.Year == year {
			// 	lm.year = ly.Year - 1
			// }
			ms = append(ms, lm)
		}
	}

	ly.months = &ms
	length := 12
	if ly.LeapN > -1 {
		length++
	}
	Ms := make([]LunarMonth, length)
	copy(Ms, ms[offset:offset+length])
	// fmt.Println("haha", offset, length, len(ms), len(Ms))
	ly.Months = &Ms
	// ms = append(ms, m)
	return
}

// func checkLeap(shuoes, terms [2][]float64) [2]int {
// 	var leapI [2]int
// 	for i := 0; i < 2; i++ {
// 		leapI[i] = -1
// 		leapI[i] = getLeapI(shuoes[i], terms[i])
// 		fmt.Println(leapI[i])
// 	}
// 	return leapI
// }

// 中气与合朔日发生在同一天，是用“发生时刻的先后顺序确定某月是否包中气”还是用“日期来确定包含关系”。
// 从原理上说，这两种方法都是可行的，不过，传统上为了降低历算的精度要求，采用后者来判断一个月中是否包含中气，紫金历也是如此。
func getLeapI(shuoes, terms []float64) int {
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

func getSpringFest(shuoes []float64, leapI int) (float64, int) {
	springFest := shuoes[2]
	offset := 2
	if leapI != -1 && leapI <= 2 { //闰11或闰12月
		springFest = shuoes[3]
	}
	return math.Floor(beijingTime(springFest) + 0.5), offset
}

func i2monthseq(i int) int {
	l := (i + jianyue) % 12
	if l == 0 {
		l = 12
	}
	return l - 1
}

// GregorianToLunarDate 北京时间转农历
func GregorianToLunarDate(y, m, d int, ly *LunarYear) LunarYMD {
	var lymd LunarYMD
	day := float64(d) + 0.5
	jd := julian.CalendarGregorianToJD(y, m, day) // 儒略日数
	if ly == nil || jd < math.Floor(ly.dzs[0]+0.5) {
		ly = GenLunarYear(y)
	}
	if jd >= math.Floor(ly.dzs[2]+0.5) {
		ly = GenLunarYear(y + 1)
	}
	// fmt.Println(ly)
	// fmt.Println(len(*ly.Months))
	prev := (*(ly.months))[0]
	for _, m := range *ly.months {
		if jd < m.d0 {
			break
		}
		prev = m
	}
	lymd.D = int(jd - prev.d0)
	lymd.M = prev.seq
	lymd.Y = prev.year
	lymd.Leap = prev.leap
	// year, _, _ := julian.JDToCalendar(ly.SpringFest)
	// if jd < ly.SpringFest && y == year {
	// 	lymd.Y = y - 1
	// }
	return lymd
}

func (lymd LunarYMD) String() string {
	leap := ""
	if lymd.Leap {
		leap = "闰"
	}
	return fmt.Sprintf("%d年%s%s%s", lymd.Y, leap, monthName[lymd.M], dayName[lymd.D])
}

// Stat 列出基本信息
func (ly LunarYear) Stat() {
	fmt.Println("年：", ly.Year)
	fmt.Println("闰月：", ly.LeapN+1)
	fmt.Println("春节：")
	fmt.Println(julian.JDToCalendar(ly.SpringFest))
	fmt.Println("xxxxxxxxxxxxxxxxxx")
	for _, m := range *ly.Months {
		fmt.Println("月首：")
		fmt.Println(julian.JDToCalendar(m.d0))
		fmt.Println("月长：", m.dn)
		fmt.Println("闰：", m.leap)
		fmt.Println("月：", m.seq+1)
		fmt.Println("年：", m.year)
		fmt.Println("==============")
	}
	fmt.Println("两个自然年是否有闰：", ly.leap)
	// fmt.Println("两个自然年中冬至之间包含的朔日个数：", ly.shuoCnt[0], ly.shuoCnt[1])
	for i, dz := range ly.dzs {
		fmt.Println("冬至：", i, math.Floor(beijingTime(dz)+0.5))
		// fmt.Println(julian.JDToCalendar(dz))
	}
	for _, shuo := range ly.Shuoes {
		for i, v := range shuo {
			fmt.Println("朔：", i, math.Floor(beijingTime(v)+0.5))
			// fmt.Println(julian.JDToCalendar(v))
		}
	}
	for _, term := range ly.Terms {
		for i := 0; i < len(term); i = i + 2 {
			fmt.Println("气：", i/2, math.Floor(beijingTime(term[i])+0.5))
			// fmt.Println(julian.JDToCalendar(v))
		}
	}
	for _, m := range *ly.months {
		fmt.Println("月首：")
		fmt.Println(julian.JDToCalendar(m.d0))
		fmt.Println("月长：", m.dn)
		fmt.Println("闰：", m.leap)
		fmt.Println("月：", m.seq+1)
		fmt.Println("年：", m.year)
		fmt.Println("==============")
	}
}

func beijingTime(jd float64) float64 {
	return jd - deltat(jd) + float64(8)/24
}

// 判断气相对于朔的关系，用于判断具体闰哪一个月
// 若朔=气，则认为该气属于该朔
// 若朔<气，则认为该气可能属于该朔，还要检查该气是否小于下一个朔
// 若朔>气，则认为该气不属于该朔
func sLEq(s, q float64) bool {
	s = math.Floor(beijingTime(s) + 0.5)
	// s = shuoC(s, shuoCorrect)
	q = math.Floor(beijingTime(q) + 0.5)
	switch {
	case s <= q:
		return true
	default:
		return false
	}
}

// 判断朔是否在两冬至之间
func sInDZs(s, dz0, dz1 float64) bool {
	s = math.Floor(beijingTime(s) + 0.5)
	// s = shuoC(s, shuoCorrect)
	dz0 = math.Floor(beijingTime(dz0) + 0.5)
	dz1 = math.Floor(beijingTime(dz1) + 0.5)
	// 注意，此处如果朔和冬至重合，则将该朔记入该冬至之前的那个自然年中，否则会出现偏差（如1984）
	if s > dz0 && s <= dz1 {
		// fmt.Println(s, dz0, dz1)
		return true
	}
	return false
}

func shuoC(shuo float64, a []struct{ jd, delta float64 }) float64 {
	key := math.Floor(beijingTime(shuo) + 0.5)
	lo := 0
	hi := len(a) - 1
	for lo <= hi {
		mid := lo + (hi-lo)/2
		if key < a[mid].jd {
			hi = mid - 1
		} else if key > a[mid].jd {
			lo = mid + 1
		} else {
			// fmt.Println(shuo, shuo+a[mid].delta)
			return shuo + a[mid].delta
		}
	}
	return shuo
}
