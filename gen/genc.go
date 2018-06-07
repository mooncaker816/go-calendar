package main

import (
	"fmt"
	"math"
	"strings"

	calendar "github.com/mooncaker816/go-calendar"
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
	{1527, 6, 29, -1}, {1534, 6, 12, -1},
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
	{1649, 5, 12, -1}, {1650, 11, 23, 1},
	{1652, 2, 9, 1}, {1653, 9, 22, -1},
	{1654, 1, 18, 1}, {1662, 2, 19, -1},
	{1673, 11, 8, 1}, {1685, 2, 4, -1},
	{1687, 3, 14, -1}, {1694, 6, 23, -1},
	{1704, 10, 28, 1}, {1708, 2, 22, -1},
	{1720, 7, 6, -1}, {1759, 3, 29, -1},
	{1763, 9, 8, -1}, {1778, 3, 29, -1},
	{1779, 7, 14, -1}, {1787, 12, 10, -1},
	{1789, 7, 23, -1}, {1796, 6, 6, -1},
	{1804, 8, 6, -1}, {1821, 6, 29, 1},
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
	for _, shuo := range shuoes {
		jdN := julian.CalendarGregorianToJD(shuo.y, shuo.m, float64(shuo.d)+0.5)
		fmt.Printf("%7.f %d\n", jdN, shuo.delta)
		offset := int(math.Floor((jdN - 1947168 + 14) / 29.5306))
		// fmt.Println("offset:", offset)
		fmt.Println(string(s[offset]))
	}
	fmt.Println(calendar.GenLunarYear(1685))
}
