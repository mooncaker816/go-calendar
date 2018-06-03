package calendar

import (
	"fmt"
	"log"
	"math"

	"github.com/mooncaker816/learnmeeus/v3/base"
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
	Terms   [2][]float64  // 两个自然年包含的节气
	mTerms  [2][]float64  // 两个自然年包含的节气的J1900起算儒略日数
	Shuoes  [2][]float64  // 两个自然年包含的朔日
	mShuoes [2][]float64  // 两个自然年包含的朔日的J1900起算儒略日数
	shuoCnt [2]int        // 两个自然年包含的朔日的个数
	dzs     [3]float64    // 划分两个自然年的三个冬至
	leap    [2]bool       // 两个自然年中是否有闰月
	months  *[]LunarMonth // 两个自然年的所有月份
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
			ly.mTerms[i] = append(ly.mTerms[i], math.Floor(t+0.5)-base.J1900)
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
	return beijingTime(jde)
}

// 计算合朔日
func (ly *LunarYear) moonShuoes() {
	for i := 0; i < 2; i++ {
		jd0, jd1 := ly.dzs[i], ly.dzs[i+1]
		nm0 := newmoonI(jd0, 0)
		ly.shuoCnt[i] = 0
		// fmt.Println(nm0, jd0)
		if math.Floor(nm0+0.5) > math.Floor(jd0+0.5) {
			jd0 -= 29.5306
			ly.shuoCnt[i]--
		}
		for j := 0; j < 15; j++ {
			shuo := newmoonI(jd0, j)
			// fmt.Println(shuo)
			shuop := math.Floor(shuo + 0.5)
			// if shuo >= jd0 && shuo < jd1 {
			if shuop >= math.Floor(jd0+0.5) && shuop < math.Floor(jd1+0.5) {
				ly.shuoCnt[i]++
			}
			ly.Shuoes[i] = append(ly.Shuoes[i], shuo)
			ly.mShuoes[i] = append(ly.mShuoes[i], shuop-base.J1900)
		}
		if ly.shuoCnt[i] >= 13 { //冬至之间是否有13个朔日
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

func newmoonI(jd float64, i int) float64 {
	nmjd := jd + 29.5306*float64(i)
	y := jd2year(nmjd)
	return beijingTime(moonphase.New(y))
}

// 建月，两个自然年之间的所有月份，包含一整个农历年月份
func (ly *LunarYear) genLunarMonth() {
	var ms []LunarMonth
	monthnum := []int{12, 12}
	var leapI [2]int
	leapI[0], leapI[1] = -1, -1
	if ly.leap[0] || ly.leap[1] {
		leapI = checkLeap(ly.mShuoes, ly.mTerms) //检查两个农历自然年的闰月情况
		if ly.leap[0] {
			monthnum[0]++
		}
		if ly.leap[1] {
			monthnum[1]++
		}
		if leapI[0] >= 3 && leapI[0] <= 12 { // 闰农历年的1-10月
			ly.LeapN = i2monthseq(leapI[0] - 1)
		}
		if leapI[1] <= 2 && leapI[1] > 0 { //闰农历年的11，12月
			ly.LeapN = i2monthseq(leapI[1] - 1)
		}
	}
	// 定农历年首（春节）
	var offset int
	ly.SpringFest, offset = getSpringFest(ly.mShuoes[0], leapI[0])

	// fmt.Println(ly.leap[0], ly.leap[1])
	// fmt.Println(leapI[0], leapI[1])
	for i := 0; i < 2; i++ {
		for j := 0; j < monthnum[i]; j++ {
			var lm LunarMonth
			lm.d0 = ly.mShuoes[i][j] + base.J1900 //儒略日数
			lm.dn = int(ly.mShuoes[i][j+1] + base.J1900 - lm.d0)
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
func checkLeap(shuoes, terms [2][]float64) [2]int {
	var leapI [2]int
	for i := 0; i < 2; i++ {
		leapI[i] = -1
		leapI[i] = getLeapI(shuoes[i], terms[i])
	}
	return leapI
}

// 中气与合朔日发生在同一天，是用“发生时刻的先后顺序确定某月是否包中气”还是用“日期来确定包含关系”。
// 从原理上说，这两种方法都是可行的，不过，传统上为了降低历算的精度要求，采用后者来判断一个月中是否包含中气，紫金历也是如此。
func getLeapI(shuoes, terms []float64) int {
	j := 0
	t := terms[0]
	for i := 0; i < 13; i++ {
		s0 := shuoes[i]
		s1 := shuoes[i+1]
		for t < s0 {
			j += 2
			if j >= len(terms) {
				return i
			}
			t = terms[j]
		}
		if s0 <= t && t < s1 {
			continue
		}
		return i // [1,12]
	}
	return -1
}

func getSpringFest(shuoes []float64, leapI int) (float64, int) {
	springFest := shuoes[2] //春节儒略日数
	offset := 2
	if leapI != -1 && leapI <= 2 { //闰11或闰12月
		springFest = shuoes[3] //春节儒略日数
	}
	return springFest + base.J1900, offset
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
	fmt.Println("两个自然年中冬至之间包含的朔日个数：", ly.shuoCnt[0], ly.shuoCnt[1])
	for i, dz := range ly.dzs {
		fmt.Println("冬至：", i, math.Floor(dz)-base.J1900)
		// fmt.Println(julian.JDToCalendar(dz))
	}
	for _, shuo := range ly.mShuoes {
		for i, v := range shuo {
			fmt.Println("朔：", i, v)
			// fmt.Println(julian.JDToCalendar(v))
		}
	}
	for _, term := range ly.mTerms {
		for i := 0; i < len(term); i = i + 2 {
			fmt.Println("气：", i/2, term[i])
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
	return jd + deltat(jd) + float64(8)/24
}
