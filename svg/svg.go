package main

import (
	"fmt"
	"os"

	"github.com/ajstarks/svgo"
)

var dayName = []string{
	"初一", "初二", "初三", "初四", "初五", "初六", "初七", "初八", "初九", "初十",
	"十一", "十二", "十三", "十四", "十五", "十六", "十七", "十八", "十九", "二十",
	"廿一", "廿二", "廿三", "廿四", "廿五", "廿六", "廿七", "廿八", "廿九", "三十",
}

func main() {
	canvas := svg.New(os.Stdout)
	canvas.Start(840, 600)
	// week0 := 3
	// dayN := 31
	// lday0 := 23
	style := "font-size:30pt;fill:Red;text-anchor:middle"
	canvas.Text(420, 45, "六月", style)

	for i, v := range []string{"日", "一", "二", "三", "四", "五", "六", "七"} {
		canvas.Text(120*(i+1)-60, 105, v, style)
	}
	style = "font-size:20pt;fill:black;text-anchor:middle"
	for i := 1; i <= 7; i++ {
		canvas.Text(120*i-60, 165, fmt.Sprintf("%d", i), style)
		canvas.Text(120*i-60, 225, dayName[i-1], style)
	}
	canvas.Grid(0, 0, 840, 600, 60, "stroke:black")
	canvas.End()
}
