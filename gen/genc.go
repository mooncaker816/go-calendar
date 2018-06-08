package main

import (
	"fmt"
	"math"
	"strings"

	dt "github.com/mooncaker816/learnmeeus/v3/deltat"

	"github.com/mooncaker816/learnmeeus/v3/moonphase"
	"github.com/soniakeys/unit"

	"github.com/mooncaker816/learnmeeus/v3/julian"
)

func init() {
	suoS = "EqoFscDcrFpmEsF2DfFideFelFpFfFfFiaipqti1ksttikptikqckstekqttgkqttgkqteksttikptikq2fjstgjqttjkqttgkqt"
	suoS += "ekstfkptikq2tijstgjiFkirFsAeACoFsiDaDiADc1AFbBfgdfikijFifegF1FhaikgFag1E2btaieeibggiffdeigFfqDfaiBkF"
	suoS += "1kEaikhkigeidhhdiegcFfakF1ggkidbiaedksaFffckekidhhdhdikcikiakicjF1deedFhFccgicdekgiFbiaikcfi1kbFibef"
	suoS += "gEgFdcFkFeFkdcfkF1kfkcickEiFkDacFiEfbiaejcFfffkhkdgkaiei1ehigikhdFikfckF1dhhdikcfgjikhfjicjicgiehdik"
	suoS += "cikggcifgiejF1jkieFhegikggcikFegiegkfjebhigikggcikdgkaFkijcfkcikfkcifikiggkaeeigefkcdfcfkhkdgkegieid"
	suoS += "hijcFfakhfgeidieidiegikhfkfckfcjbdehdikggikgkfkicjicjF1dbidikFiggcifgiejkiegkigcdiegfggcikdbgfgefjF1"
	suoS += "kfegikggcikdgFkeeijcfkcikfkekcikdgkabhkFikaffcfkhkdgkegbiaekfkiakicjhfgqdq2fkiakgkfkhfkfcjiekgFebicg"
	suoS += "gbedF1jikejbbbiakgbgkacgiejkijjgigfiakggfggcibFifjefjF1kfekdgjcibFeFkijcfkfhkfkeaieigekgbhkfikidfcje"
	suoS += "aibgekgdkiffiffkiakF1jhbakgdki1dj1ikfkicjicjieeFkgdkicggkighdF1jfgkgfgbdkicggfggkidFkiekgijkeigfiski"
	suoS += "ggfaidheigF1jekijcikickiggkidhhdbgcfkFikikhkigeidieFikggikhkffaffijhidhhakgdkhkijF1kiakF1kfheakgdkif"
	suoS += "iggkigicjiejkieedikgdfcggkigieeiejfgkgkigbgikicggkiaideeijkefjeijikhkiggkiaidheigcikaikffikijgkiahi1"
	suoS += "hhdikgjfifaakekighie1hiaikggikhkffakicjhiahaikggikhkijF1kfejfeFhidikggiffiggkigicjiekgieeigikggiffig"
	suoS += "gkidheigkgfjkeigiegikifiggkidhedeijcfkFikikhkiggkidhh1ehigcikaffkhkiggkidhh1hhigikekfiFkFikcidhh1hit"
	suoS += "cikggikhkfkicjicghiediaikggikhkijbjfejfeFhaikggifikiggkigiejkikgkgieeigikggiffiggkigieeigekijcijikgg"
	suoS += "ifikiggkideedeijkefkfckikhkiggkidhh1ehijcikaffkhkiggkidhh1hhigikhkikFikfckcidhh1hiaikgjikhfjicjicgie"
	suoS += "hdikcikggifikigiejfejkieFhegikggifikiggfghigkfjeijkhigikggifikiggkigieeijcijcikfksikifikiggkidehdeij"
	suoS += "cfdckikhkiggkhghh1ehijikifffffkhsFngErD1pAfBoDd1BlEtFqA2AqoEpDqElAEsEeB2BmADlDkqBtC1FnEpDqnEmFsFsAFn"
	suoS += "llBbFmDsDiCtDmAB2BmtCgpEplCpAEiBiEoFqFtEqsDcCnFtADnFlEgdkEgmEtEsCtDmADqFtAFrAtEcCqAE1BoFqC1F1DrFtBmF"
	suoS += "tAC2ACnFaoCgADcADcCcFfoFtDlAFgmFqBq2bpEoAEmkqnEeCtAE1bAEqgDfFfCrgEcBrACfAAABqAAB1AAClEnFeCtCgAADqDoB"
	suoS += "mtAAACbFiAAADsEtBqAB2FsDqpFqEmFsCeDtFlCeDtoEpClEqAAFrAFoCgFmFsFqEnAEcCqFeCtFtEnAEeFtAAEkFnErAABbFkAD"
	suoS += "nAAeCtFeAfBoAEpFtAABtFqAApDcCGJ"

	//1645-09-23开始7567个节气修正表
	qiS = "FrcFs22AFsckF2tsDtFqEtF1posFdFgiFseFtmelpsEfhkF2anmelpFlF1ikrotcnEqEq2FfqmcDsrFor22FgFrcgDscFs22FgEe"
	qiS += "FtE2sfFs22sCoEsaF2tsD1FpeE2eFsssEciFsFnmelpFcFhkF2tcnEqEpFgkrotcnEqrEtFermcDsrE222FgBmcmr22DaEfnaF22"
	qiS += "2sD1FpeForeF2tssEfiFpEoeFssD1iFstEqFppDgFstcnEqEpFg11FscnEqrAoAF2ClAEsDmDtCtBaDlAFbAEpAAAAAD2FgBiBqo"
	qiS += "BbnBaBoAAAAAAAEgDqAdBqAFrBaBoACdAAf1AACgAAAeBbCamDgEifAE2AABa1C1BgFdiAAACoCeE1ADiEifDaAEqAAFe1AcFbcA"
	qiS += "AAAAF1iFaAAACpACmFmAAAAAAAACrDaAAADG0"
}

var shuoes = []struct {
	y, m, d, delta int
}{
	{1501, 6, 16, -1}, {1508, 1, 2, 1},
	{1513, 10, 29, -1}, {1516, 10, 25, 1},
	{1521, 10, 1, -1}, {1526, 7, 10, -1},
	{1527, 6, 29, -1}, {1534, 6, 12, -1}, // need check
	{1535, 8, 29, -1}, {1535, 10, 26, 1},
	{1544, 5, 22, -1}, {1546, 1, 2, 1},
	{1546, 7, 28, -1}, {1571, 8, 21, -1},
	{1572, 8, 9, -1}, {1581, 10, 27, 1},
	{1582, 7, 20, -1}, {1588, 4, 26, -1},
	{1589, 1, 16, 1}, {1591, 9, 18, -1},
	{1599, 1, 26, 1}, {1600, 2, 15, -1},
	{1612, 3, 2, 1}, {1616, 5, 16, -1},
	{1622, 7, 9, -1}, {1627, 9, 10, -1},
	{1628, 1, 6, 1}, {1630, 4, 12, 1},
	{1634, 8, 24, -1}, {1643, 2, 18, 1},
	{1649, 5, 12, -1}, {1650, 11, 23, 1}, // need check
	{1652, 2, 9, 1} /* need check */, {1653, 9, 22, -1}, // need check
	{1654, 1, 18, 1} /* need check */, {1662, 2, 19, -1},
	{1673, 11, 8, 1}, {1685, 2, 4, -1},
	{1687, 3, 14, -1}, {1694, 6, 23, -1},
	{1704, 10, 28, 1}, {1708, 2, 22, -1},
	{1720, 7, 6, -1}, {1759, 3, 29, -1},
	{1763, 9, 8, -1}, {1778, 3, 29, -1},
	{1779, 7, 14, -1}, {1787, 12, 10, -1},
	{1789, 7, 23, -1}, {1796, 6, 6, -1},
	{1804, 8, 6, -1}, {1821, 6, 29, 1}, // need check
	{1831, 4, 13, -1}, {1842, 1, 12, -1},
	{1863, 1, 20, -1}, {1880, 11, 2, 1},
	{1896, 2, 14, -1}, {1914, 11, 18, -1},
	{1916, 2, 4, -1}, {1920, 11, 11, -1},
}
var oldnew = []string{
	"J", "00",
	"I", "000",
	"H", "0000",
	"G", "00000",
	"t", "02",
	"s", "002",
	"r", "0002",
	"q", "00002",
	"p", "000002",
	"o", "0000002",
	"n", "00000002",
	"m", "000000002",
	"l", "0000000002",
	"k", "01",
	"j", "0101",
	"i", "001",
	"h", "001001",
	"g", "0001",
	"f", "00001",
	"e", "000001",
	"d", "0000001",
	"c", "00000001",
	"b", "000000001",
	"a", "0000000001",
	"A", "000000000000000000000000000000000000000000000000000000000000",
	"B", "00000000000000000000000000000000000000000000000000",
	"C", "0000000000000000000000000000000000000000",
	"D", "000000000000000000000000000000",
	"E", "00000000000000000000",
	"F", "0000000000",
}
var suoS, qiS string

func extract(s string) string {
	return strings.NewReplacer(oldnew...).Replace(s)
}
func main() {
	s := extract(suoS)
	fmt.Println(len(s))
	// fmt.Println(s)
	// for _, shuo := range shuoes {
	// 	jdN := julian.CalendarGregorianToJD(shuo.y, shuo.m, float64(shuo.d)+0.5)
	// 	fmt.Printf("%7.f %d\n", jdN, shuo.delta)
	// 	offset := int(math.Floor((jdN - 1947168 + 14) / 29.5306))
	// 	// fmt.Println("offset:", offset)
	// 	fmt.Println(string(s[offset]))
	// }
	// fmt.Println(calendar.GenLunarYear(1685))
	shuos := getShuo([]byte(s))
	count := 0
	for i, v := range shuos {
		if v.old != math.Floor(v.real+0.5) {
			count++
			y, m, d := julian.JDToCalendar(v.old)
			fmt.Println(i, y, m, d)
			y, m, d = julian.JDToCalendar(math.Floor(v.real + 0.5))
			fmt.Println(i, y, m, d, julian.JDToTime(v.real))

		}
	}
	fmt.Println(len(shuos), count)
}

//朔直线拟合参数
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
	if jde >= 1947168.00-14 || jde < 1457698.231017-14 {
		return 0
	}
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

func lowShuo(jde float64, correct []byte) float64 {
	if jde >= 2436935 || jde < 1947168.00-14 {
		return 0
	}
	W := math.Floor((jde+14-2451551)/29.5306) * math.Pi * 2
	n := correct[int(math.Floor((jde-1947168.00+14)/29.5306))] //找定朔修正值
	v := 7771.37714500204
	t := (W + 1.08472) / v
	t -= (-0.0000331*t*t+
		0.10976*math.Cos(0.785+8328.6914*t)+
		0.02224*math.Cos(0.187+7214.0629*t)-
		0.03342*math.Cos(4.669+628.3076*t))/v +
		(32*(t+1.8)*(t+1.8)-20)/86400/36525
	shuolow := math.Floor(t*36525 + float64(8)/24 + 0.5)

	if n == '1' {
		shuolow++
	}
	if n == '2' {
		shuolow--
	}
	return shuolow + 2451545
}

type shuo struct {
	old, real float64
}

func getShuo(correct []byte) []shuo {
	var s []shuo
	// for i := 1457698.; i < 1947168-14; i += 29.5306 {
	// 	y := jd2year(i)
	// 	s = append(s, shuo{avgShuo(i), math.Floor(beijingTime(moonphase.MeanNew(y)) + 0.5)})
	// }
	for i := 1947182.; i < 2436935; i += 29.5306 {
		y := jd2year(i)
		k := i
		lowshuo := lowShuo(i, correct)
		highshuo := beijingTime(moonphase.New(y))
		for math.Abs(highshuo-lowshuo) > 15 {
			k--
			y = jd2year(k)
			highshuo = beijingTime(moonphase.New(y))
		}
		s = append(s, shuo{lowshuo, highshuo})
	}
	return s
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

func beijingTime(jde float64) float64 {
	return jde - deltat(jde) + float64(8)/24
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
