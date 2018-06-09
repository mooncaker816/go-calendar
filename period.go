package calendar

import (
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
	case p.GYear >= -221 && p.GYear < -104:
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
	// 公元前170年为闰年，为第一式闰章的终止年
	v := mod((y - (-169)), 19) //闰章索引
	if y+1 > -169 {
		i := sort.SearchInts(r2, v)
		if i < len(r2) && r2[i] == v {
			return true
		}
	} else {
		i := sort.SearchInts(r1, v)
		if i < len(r1) && r1[i] == v {
			return true
		}
	}
	return false
}

var r1 = []int{3, 6, 9, 11, 14, 17, 19} //闰章第一式
var r2 = []int{3, 6, 8, 11, 14, 17, 19} //闰章第二式
