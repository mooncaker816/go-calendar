package main

import (
	"fmt"
	"log"
	"math"
	"strings"

	"github.com/RoaringBitmap/roaring"

	"github.com/mooncaker816/go-calendar"

	pp "github.com/mooncaker816/learnmeeus/v3/planetposition"

	"github.com/mooncaker816/learnmeeus/v3/solar"

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
	q := extract(qiS)
	// fmt.Println(len(s), len(q))
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
	qis := getQi([]byte(q))
	counts := 0
	counts1 := 0
	counts2 := 0
	countq := 0
	countq1 := 0
	countq2 := 0
	ra1 := roaring.New()
	for _, v := range shuos {

		jdN := math.Floor(v.real + 0.5)
		if v.old != jdN {
			counts++
			delta := v.old - jdN
			// y, m, d := julian.JDToCalendar(v.old)
			// fmt.Println(i, y, m, d)
			// y, m, d = julian.JDToCalendar(math.Floor(v.real + 0.5))
			// fmt.Println(i, y, m, d, julian.JDToTime(v.real))
			// 1947877
			if delta == -1 {
				counts1++
				ra1.Add(uint32(jdN))
				fmt.Printf("%7d,%s,%d\n", int(jdN), calendar.JDPlus{jdN, true}, int(delta))
			}
		}
	}
	fmt.Println("++++++++++++++++++++++++++++++++++++")
	ra2 := roaring.New()
	for _, v := range shuos {
		jdN := math.Floor(v.real + 0.5)
		if v.old != jdN {
			delta := v.old - jdN
			// y, m, d := julian.JDToCalendar(v.old)
			// fmt.Println(i, y, m, d)
			// y, m, d = julian.JDToCalendar(math.Floor(v.real + 0.5))
			// fmt.Println(i, y, m, d, julian.JDToTime(v.real))
			// 1949825
			if delta == 1 {
				counts2++
				ra2.Add(uint32(jdN))
				fmt.Printf("%7d,%s,%d\n", int(jdN), calendar.JDPlus{jdN, true}, int(delta))
			}
		}
	}
	fmt.Println("++++++++++++++++++++++++++++++++++++")
	ra3 := roaring.New()
	for _, v := range qis {
		jdN := math.Floor(v.real + 0.5)
		if v.old != jdN {
			delta := v.old - jdN
			countq++
			// y, m, d := julian.JDToCalendar(v.old)
			// fmt.Println(i, y, m, d)
			// y, m, d = julian.JDToCalendar(math.Floor(v.real + 0.5))
			// fmt.Println(i, y, m, d, julian.JDToTime(v.real))
			// 2322344
			if delta == -1 {
				countq1++
				ra3.Add(uint32(jdN))
				fmt.Printf("%7d,%s,%d\n", int(jdN), calendar.JDPlus{jdN, true}, int(delta))
			}
		}
	}
	fmt.Println("++++++++++++++++++++++++++++++++++++")
	ra4 := roaring.New()
	for _, v := range qis {
		jdN := math.Floor(v.real + 0.5)
		if v.old != jdN {
			delta := v.old - jdN
			// y, m, d := julian.JDToCalendar(v.old)
			// fmt.Println(i, y, m, d)
			// y, m, d = julian.JDToCalendar(math.Floor(v.real + 0.5))
			// fmt.Println(i, y, m, d, julian.JDToTime(v.real))
			// 2322468
			if delta == 1 {
				countq2++
				ra4.Add(uint32(jdN))
				fmt.Printf("%7d,%s,%d\n", int(jdN), calendar.JDPlus{jdN, true}, int(delta))
			}
		}
	}
	str1, _ := ra1.ToBase64()
	fmt.Println(str1)
	str2, _ := ra2.ToBase64()
	fmt.Println(str2)
	str3, _ := ra3.ToBase64()
	fmt.Println(str3)
	str4, _ := ra4.ToBase64()
	fmt.Println(str4)

	fmt.Println(len(shuos), counts, counts1, counts2)
	fmt.Println(len(qis), countq, countq1, countq2)
	genLeaps()
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

func avgSQ(jde float64, delta float64, avgTab []float64) float64 {
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

func lowQi(jde float64, correct []byte) float64 {
	W := math.Floor((jde+7-2451259)/365.2422*24) * math.Pi / 12
	n := correct[int(math.Floor((jde-2322147.76+7)/365.2422*24))] //找定朔修正值
	v := 628.3319653318
	t := (W - 4.895062166) / v                                                                            //第一次估算,误差2天以内
	t -= (53*t*t + 334116*math.Cos(4.67+628.307585*t) + 2061*math.Cos(2.678+628.3076*t)*t) / v / 10000000 //第二次估算,误差2小时以内

	L := 48950621.66 + 6283319653.318*t + 53*t*t + //平黄经
		+334166*math.Cos(4.669257+628.307585*t) + //地球椭圆轨道级数展开
		+3489*math.Cos(4.6261+1256.61517*t) + //地球椭圆轨道级数展开
		+2060.6*math.Cos(2.67823+628.307585*t)*t + //一次泊松项
		-994 - 834*math.Sin(2.1824-33.75705*t) //光行差与章动修正

	t -= (L/10000000-W)/628.332 + (32*(t+1.8)*(t+1.8)-20)/86400/36525
	qilow := math.Floor(t*36525 + float64(8)/24 + 0.5)
	if n == '1' {
		qilow++
	}
	if n == '2' {
		qilow--
	}
	return qilow + 2451545
}

type shuo struct {
	old, real float64
}

type qi struct {
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

func getQi(correct []byte) []qi {
	var q []qi
	// for i := 1457698.; i < 1947168-14; i += 29.5306 {
	// 	y := jd2year(i)
	// 	s = append(s, shuo{avgShuo(i), math.Floor(beijingTime(moonphase.MeanNew(y)) + 0.5)})
	// }
	hj := unit.Angle(math.Pi)
	for i := 2322147.76 + 7; i < 2436935; i += 15.2184 {
		lowqi := lowQi(i, correct)
		hj += unit.Angle(math.Pi / 12) // 节气对应的黄经
		hj = hj.Mod1()
		// k := hj
		earth, err := pp.LoadPlanet(pp.Earth)
		if err != nil {
			log.Fatalf("can not load planet: %v", err)
		}
		highqi := 0.
		for {
			λ, _, _ := solar.ApparentVSOP87(earth, i)
			c := 58 * (hj - λ).Sin()
			i += c
			if math.Abs(c) < .000005 {
				break
			}
		}
		highqi = beijingTime(i)
		// for math.Abs(highqi-lowqi) > 15 {
		// 	k +=
		// }
		q = append(q, qi{lowqi, highqi})
	}
	return q
}

func jd2year(jd float64) float64 {
	year, m, d := julian.JDToCalendar(jd)
	z, f := math.Modf(d)
	yeardays := 365.
	leap := calendar.LeapYear(year)
	if leap {
		yeardays++
	}
	return float64(year) + (float64(julian.DayOfYear(year, m, int(z), leap))+f)/yeardays
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

func genLeaps() {
	prev := 0.
	var leaps []int
	for i := -723.; i <= -104; i++ {
		jd := 0.
		k := 0
		if i >= -220 {
			jd = 1640641 + math.Floor(0.866+(i+220)*12.369000)*29.5306
		} else if i >= -480 {
			jd = 1546083 + math.Floor(0.500+(i+479)*12.368422)*29.5306
		} else {
			jd = 1457698 + math.Floor(0.342+(i+721)*12.368422)*29.5306
		}
		sf := avgSQ(jd, 14, avgShuoTab)
		y, _, _ := julian.JDToCalendar(sf)
		if math.Floor((sf-prev+15)/29.5306) == 13 {
			k++
			leaps = append(leaps, y)
			fmt.Printf("{%d, %v},", y, true)
		} else {
			fmt.Printf("{%d, %v},", y, false)
		}
		if (y+720)%19 == 0 {
			fmt.Println()
		}
		prev = sf
	}
	// sort.Sort(sort.Reverse(sort.IntSlice(leaps)))
	// fmt.Println(leaps)
	// 	base := []int{-359, -361, -364, -367, -369, -372, -375}
	// 	k := 0
	// loop:
	// 	for i := 0; ; i++ {
	// 		for _, v := range base {
	// 			leap := v - 19*i
	// 			fmt.Printf("%d ", leap)
	// 			if leap != leaps[k] {
	// 				fmt.Printf("*%d ", leap-leaps[k])
	// 			}
	// 			k++
	// 			if v-19*i < -722 || k >= len(leaps) {
	// 				break loop
	// 			}
	// 		}
	// 		fmt.Println()
	// 	}

}
