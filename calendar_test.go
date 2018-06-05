package calendar

import (
	"fmt"

	"github.com/mooncaker816/learnmeeus/v3/julian"
)

func ExampleGenDay() {
	ly := GenLunarYear(1987)
	jd0 := julian.CalendarGregorianToJD(1987, 6, 5.5)
	jd1 := julian.CalendarGregorianToJD(1987, 6, 6.5)
	fmt.Println(genDay(jd0, ly))
	fmt.Println(genDay(jd1, ly))
	// Output:
	//
}
