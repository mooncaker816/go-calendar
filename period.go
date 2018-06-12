package calendar

import (
	"fmt"
	"sort"
)

type YueJian uint8

type ZhiRun uint8

const (
	ZZ         YueJian = iota //子正
	CZ                        //丑正
	YZ                        //寅正
	ZZYY                      //子正寅一
	R7in19st1  ZhiRun  = iota //19年7闰，正月为首，闰在岁末
	R7in19st10                //19年7闰，十月为首，闰在岁末
	RNoZQ                     //无中气置闰
)

type Period struct {
	GYear int //冬至对应的格里历年份
	YueJian
	ZhiRun
}

func (ly *LunarYear) getPeriod(year int) {
	var p Period
	p.GYear = year
	p.YueJian = YZ
	p.ZhiRun = RNoZQ
	switch {
	case p.GYear < -221:
		// p.GYear++
		p.YueJian = ZZ
		p.ZhiRun = R7in19st1
	case p.GYear >= -221 && p.GYear <= -104:
		p.ZhiRun = R7in19st10
	case p.GYear >= 8 && p.GYear < 23:
		p.YueJian = CZ
	case p.GYear >= 237 && p.GYear < 239:
		p.YueJian = CZ
	case p.GYear >= 689 && p.GYear < 700:
		// p.GYear++
		p.YueJian = ZZYY
	case p.GYear == 761:
		// p.GYear++
		p.YueJian = ZZ
	}
	ly.Period = p
}

func chkR7in19(y int) bool {
	v := mod((y - (-169)), 19) //闰章索引
	if v == 0 {
		v = 19
	}
	switch {
	//闰章第二式
	case y <= -226 || y > -169 && y <= -104:
		i := sort.SearchInts(r2, v)
		if i < len(r2) && r2[i] == v {
			return true
		}
	//闰章第一式
	case y > -226 && y <= -169:
		i := sort.SearchInts(r1, v)
		if i < len(r1) && r1[i] == v {
			return true
		}
	}
	return false
}

var r1 = []int{3, 6, 9, 11, 14, 17, 19} //闰章第一式
var r2 = []int{3, 6, 8, 11, 14, 17, 19} //闰章第二式

func (p Period) String() string {
	var yjStr, leapStr string
	switch p.YueJian {
	case ZZ:
		yjStr = "子正"
	case CZ:
		yjStr = "丑正"
	case YZ:
		yjStr = "寅正"
	case ZZYY:
		yjStr = "子正寅一"
	}
	switch p.ZhiRun {
	case R7in19st1:
		leapStr = "19年7闰，正月为首，末尾置闰"
	case R7in19st10:
		leapStr = "19年7闰，十月为首，末尾置闰"
	case RNoZQ:
		leapStr = "无中气置闰"
	}
	return fmt.Sprintf("%d %s%s", p.GYear, yjStr, leapStr)
}
