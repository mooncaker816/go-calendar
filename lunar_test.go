package calendar

import (
	"sort"
	"sync"
	"testing"
)

// func ExampleGenLunarYear() {
// 	ly := GenLunarYear(2100)
// 	for _, term := range ly.Terms[0] {
// 		fmt.Println(DT2SolarTime(term))
// 	}
// 	for _, term := range ly.Terms[1] {
// 		fmt.Println(DT2SolarTime(term))
// 	}
// 	fmt.Println()
// 	for _, shuo := range ly.Shuoes[0] {
// 		fmt.Println(DT2SolarTime(shuo))
// 	}
// 	for _, shuo := range ly.Shuoes[1] {
// 		fmt.Println(DT2SolarTime(shuo))
// 	}
// 	// Output:
// 	// x
// }

// func ExampleGregorianToLunarDate() {
// 	fmt.Println(DayCalendar(2018, 6, 2, true, nil))
// 	fmt.Println(DayCalendar(2017, 7, 25, true, nil))
// 	fmt.Println(DayCalendar(2017, 12, 21, true, nil))
// 	fmt.Println(DayCalendar(2017, 12, 22, true, nil))
// 	fmt.Println(DayCalendar(2017, 12, 23, true, nil))
// 	fmt.Println(DayCalendar(2017, 12, 31, true, nil))
// 	fmt.Println(DayCalendar(2018, 1, 1, true, nil))
// 	fmt.Println(DayCalendar(2018, 1, 31, true, nil))
// 	fmt.Println(DayCalendar(2018, 2, 15, true, nil))
// 	fmt.Println(DayCalendar(2018, 2, 16, true, nil))
// 	// Output:
// 	// 2018年四月十九
// 	// 2017年闰六月初三
// 	// 2017年冬月初四
// 	// 2017年冬月初五
// 	// 2017年冬月初六
// 	// 2017年冬月十四
// 	// 2017年冬月十五
// 	// 2017年腊月十五
// 	// 2017年腊月三十
// 	// 2018年正月初一
// }

var leaps = []struct {
	year int
	m    string
}{
	{1645, "五月"}, {1648, "四月"}, {1651, "正月"}, {1653, "六月"}, {1656, "五月"},
	{1659, "三月"}, {1661, "八月"}, {1664, "六月"}, {1667, "四月"}, {1670, "二月"},
	{1672, "七月"}, {1675, "五月"}, {1678, "三月"}, {1680, "八月"}, {1683, "六月"},
	{1686, "四月"}, {1689, "三月"}, {1691, "七月"}, {1694, "五月"}, {1697, "三月"},
	{1699, "七月"}, {1702, "六月"}, {1705, "四月"}, {1708, "三月"}, {1710, "七月"},
	{1713, "五月"}, {1716, "三月"}, {1718, "八月"}, {1721, "六月"}, {1724, "四月"},
	{1727, "二月"}, {1729, "七月"}, {1732, "五月"}, {1735, "四月"}, {1737, "九月"},
	{1740, "六月"}, {1743, "四月"}, {1746, "三月"}, {1748, "七月"}, {1751, "五月"},
	{1754, "四月"}, {1756, "九月"}, {1759, "六月"}, {1762, "五月"}, {1765, "二月"},
	{1767, "七月"}, {1770, "五月"}, {1773, "三月"}, {1775, "十月"}, {1778, "六月"},
	{1781, "五月"}, {1784, "三月"}, {1786, "七月"}, {1789, "五月"}, {1792, "四月"},
	{1795, "二月"}, {1797, "六月"}, {1800, "四月"}, {1803, "二月"}, {1805, "六月"}, //{1805, "七月"}原文有误，应该闰六月
	{1808, "五月"}, {1811, "三月"}, {1814, "二月"}, {1816, "六月"}, {1819, "四月"},
	{1822, "三月"}, {1824, "七月"}, {1827, "五月"}, {1830, "四月"}, {1832, "九月"},
	{1835, "六月"}, {1838, "四月"}, {1841, "三月"}, {1843, "七月"}, {1846, "五月"},
	{1849, "四月"}, {1851, "八月"}, {1854, "七月"}, {1857, "五月"}, {1860, "三月"},
	{1862, "八月"}, {1865, "五月"}, {1868, "四月"}, {1870, "十月"}, {1873, "六月"},
	{1876, "五月"}, {1879, "三月"}, {1881, "七月"}, {1884, "五月"}, {1887, "四月"},
	{1890, "二月"}, {1892, "六月"}, {1895, "五月"}, {1898, "三月"}, {1900, "八月"},
	{1903, "五月"}, {1906, "四月"}, {1909, "二月"}, {1911, "六月"}, {1914, "五月"},
	{1917, "二月"}, {1919, "七月"}, {1922, "五月"}, {1925, "四月"}, {1928, "二月"},
	{1930, "六月"}, {1933, "五月"}, {1936, "三月"}, {1938, "七月"}, {1941, "六月"},
	{1944, "四月"}, {1947, "二月"}, {1949, "七月"}, {1952, "五月"}, {1955, "三月"},
	{1957, "八月"}, {1960, "六月"}, {1963, "四月"}, {1966, "三月"}, {1968, "七月"},
	{1971, "五月"}, {1974, "四月"}, {1976, "八月"}, {1979, "六月"}, {1982, "四月"},
	{1984, "十月"}, {1987, "六月"}, {1990, "五月"}, {1993, "三月"}, {1995, "八月"},
	{1998, "五月"}, {2001, "四月"}, {2004, "二月"}, {2006, "七月"}, {2009, "五月"},
	{2012, "四月"}, {2014, "九月"}, {2017, "六月"}, {2020, "四月"}, {2023, "二月"},
	{2025, "六月"}, {2028, "五月"}, {2031, "三月"}, {2033, "冬月"}, {2036, "六月"},
	{2039, "五月"}, {2042, "二月"}, {2044, "七月"}, {2047, "五月"}, {2050, "三月"},
	{2052, "八月"}, {2055, "六月"}, {2058, "四月"}, {2061, "三月"}, {2063, "七月"},
	{2066, "五月"}, {2069, "四月"}, {2071, "八月"}, {2074, "六月"}, {2077, "四月"},
	{2080, "三月"}, {2082, "七月"}, {2085, "五月"}, {2088, "四月"}, {2090, "八月"},
	{2093, "六月"}, {2096, "四月"}, {2099, "二月"}, {2101, "七月"}, {2104, "五月"},
	{2107, "四月"}, {2109, "九月"}, {2112, "六月"}, {2115, "四月"}, {2118, "三月"},
	{2120, "七月"}, {2123, "五月"}, {2126, "四月"}, {2128, "冬月"}, {2131, "六月"},
	{2134, "五月"}, {2137, "二月"}, {2139, "七月"}, {2142, "五月"}, {2145, "四月"},
	{2147, "冬月"}, {2150, "六月"}, {2153, "五月"}, {2156, "三月"}, {2158, "七月"},
	{2161, "六月"}, {2164, "四月"}, {2166, "十月"}, {2169, "六月"}, {2172, "五月"},
	{2175, "三月"}, {2177, "七月"}, {2180, "六月"}, {2183, "四月"}, {2186, "二月"},
	{2188, "六月"}, {2191, "五月"}, {2194, "三月"}, {2196, "七月"}, {2199, "六月"},
	{2202, "四月"}, {2204, "九月"}, {2207, "六月"}, {2210, "四月"}, {2213, "三月"},
	{2215, "七月"}, {2218, "五月"}, {2221, "四月"}, {2223, "九月"}, {2226, "七月"},
	{2229, "五月"}, {2232, "三月"}, {2234, "八月"}, {2237, "五月"}, {2240, "四月"},
	{2242, "冬月"}, {2245, "六月"}, {2248, "五月"}, {2251, "三月"}, {2253, "七月"},
	{2256, "六月"}, {2259, "五月"}, {2262, "正月"}, {2264, "七月"}, {2267, "五月"},
	{2270, "三月"}, {2272, "八月"}, {2275, "六月"}, {2278, "四月"}, {2281, "二月"},
	{2283, "六月"}, {2286, "五月"}, {2289, "三月"}, {2291, "七月"}, {2294, "六月"},
	{2297, "四月"}, {2300, "二月"}, {2302, "六月"}, {2305, "五月"}, {2308, "三月"},
	{2310, "七月"}, {2313, "六月"}, {2316, "四月"}, {2318, "十月"}, {2321, "七月"},
	{2324, "五月"}, {2327, "三月"}, {2329, "八月"}, {2332, "六月"}, {2335, "四月"},
	{2338, "三月"}, {2340, "七月"}, {2343, "五月"}, {2346, "四月"}, {2348, "八月"},
	{2351, "六月"}, {2354, "五月"}, {2357, "正月"}, {2359, "七月"}, {2362, "五月"},
	{2365, "四月"}, {2367, "八月"}, {2370, "六月"}, {2373, "五月"}, {2376, "二月"},
	{2378, "七月"}, {2381, "五月"}, {2384, "四月"}, {2386, "十月"}, {2389, "六月"},
	{2392, "四月"}, {2395, "二月"}, {2397, "六月"}, {2400, "五月"}, {2403, "三月"},
	{2405, "八月"}, {2408, "六月"}, {2411, "五月"}, {2414, "二月"}, {2416, "七月"},
	{2419, "五月"}, {2422, "三月"}, {2424, "八月"}, {2427, "六月"}, {2430, "四月"},
	{2433, "三月"}, {2435, "七月"}, {2438, "五月"}, {2441, "四月"}, {2443, "八月"},
	{2446, "七月"}, {2449, "五月"}, {2452, "三月"}, {2454, "八月"}, {2457, "五月"},
	{2460, "四月"}, {2462, "八月"}, {2465, "六月"}, {2468, "五月"}, {2471, "三月"},
	{2473, "七月"}, {2476, "五月"}, {2479, "四月"}, {2481, "十月"}, {2484, "六月"},
	{2487, "五月"}, {2490, "三月"}, {2492, "七月"}, {2495, "五月"}, {2498, "四月"},
	{2500, "十月"}, {2503, "六月"}, {2506, "五月"}, {2509, "二月"}, {2511, "七月"},
	{2514, "五月"}, {2517, "四月"}, {2520, "正月"}, {2522, "六月"}, {2525, "五月"},
	{2528, "三月"}, {2530, "七月"}, {2533, "六月"}, {2536, "四月"}, {2539, "正月"},
	{2541, "七月"}, {2544, "五月"}, {2547, "三月"}, {2549, "七月"}, {2552, "六月"},
	{2555, "四月"}, {2557, "八月"}, {2560, "七月"}, {2563, "五月"}, {2566, "四月"},
	{2568, "七月"}, {2571, "六月"}, {2574, "四月"}, {2576, "九月"}, {2579, "六月"},
	{2582, "四月"}, {2585, "三月"}, {2587, "七月"}, {2590, "五月"}, {2593, "四月"},
	{2595, "十月"}, {2598, "七月"}, {2601, "五月"}, {2604, "三月"}, {2606, "八月"},
	{2609, "六月"}, {2612, "四月"}, {2614, "冬月"}, {2617, "六月"}, {2620, "五月"},
	{2623, "三月"}, {2625, "八月"}, {2628, "六月"}, {2631, "五月"}, {2634, "正月"},
	{2636, "七月"}, {2639, "五月"}, {2642, "三月"}, {2644, "八月"}, {2647, "六月"},
	{2650, "四月"}, {2653, "二月"}, {2655, "七月"}, {2658, "五月"}, {2661, "三月"},
	{2663, "七月"}, {2666, "六月"}, {2669, "四月"}, {2672, "三月"}, {2674, "七月"},
	{2677, "五月"}, {2680, "三月"}, {2682, "七月"}, {2685, "六月"}, {2688, "四月"},
	{2691, "三月"}, {2693, "七月"}, {2696, "五月"}, {2699, "三月"}, {2701, "八月"},
	{2704, "六月"}, {2707, "四月"}, {2710, "三月"}, {2712, "七月"}, {2715, "五月"},
	{2718, "四月"}, {2720, "九月"}, {2723, "六月"}, {2726, "五月"}, {2728, "冬月"},
	{2731, "七月"}, {2734, "五月"}, {2737, "四月"}, {2739, "九月"}, {2742, "六月"},
	{2745, "五月"}, {2748, "二月"}, {2750, "七月"}, {2753, "六月"}, {2756, "四月"},
	{2758, "八月"}, {2761, "六月"}, {2764, "五月"}, {2767, "三月"}, {2769, "七月"},
	{2772, "六月"}, {2775, "三月"}, {2777, "八月"}, {2780, "六月"}, {2783, "五月"},
	{2786, "三月"}, {2788, "七月"}, {2791, "六月"}, {2794, "三月"}, {2796, "八月"},
}

// TestLunarLeap checks all the leap year from 1645 to 2796
func Test1000yearsLunarLeap(t *testing.T) {
	count := 0
	addcount := 0
	misscount := 0
	wrongcount := 0
	goodcount := 0
	addChan := make(chan int)
	wrongChan := make(chan int)
	goodChan := make(chan int)
	totalChan := make(chan int)
	query := make(chan int)
	workerexit := make(chan struct{})
	counterexit := make(chan struct{})
	seen := make(map[int]int)
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			for {
				select {
				case <-workerexit:
					wg.Done()
					return
				case y := <-query:
					ly := GenLunarYear(y)
					if ly.LeapN > -1 {
						totalChan <- 1
						leapMonth := ly.Months[ly.LeapN]
						f := func(i int) bool {
							return leaps[i].year >= leapMonth.year
						}
						idx := sort.Search(len(leaps), f)
						if idx < len(leaps) && leaps[idx].year == leapMonth.year {
							if leaps[idx].m == monthName[leapMonth.seq] {
								goodChan <- leapMonth.year
							} else {
								t.Errorf("wrong leap month : %d %s, expect: %d %s", leapMonth.year, monthName[leapMonth.seq], leaps[idx].year, leaps[idx].m)
								wrongChan <- leapMonth.year
							}
						} else {
							t.Errorf("additional leap year: %d %s", leapMonth.year, monthName[leapMonth.seq])
							addChan <- leapMonth.year
						}
					}
				}
			}
		}()
	}
	var wg1 sync.WaitGroup
	wg1.Add(1)
	go func() {
		for {
			select {
			case y := <-addChan:
				addcount++
				seen[y]++
			case y := <-wrongChan:
				wrongcount++
				seen[y]++
			case y := <-goodChan:
				goodcount++
				seen[y]++
			case <-totalChan:
				count++
			case <-counterexit:
				wg1.Done()
				return
			}
		}
	}()

	for i := 1645; i <= 2796; i++ {
		query <- i
	}
	close(workerexit)
	wg.Wait()
	close(counterexit)
	wg1.Wait()
	for _, v := range leaps {
		if _, ok := seen[v.year]; !ok {
			misscount++
			t.Errorf("miss leap year: %d %s", v.year, v.m)
		}
	}

	if addcount > 0 || misscount > 0 || wrongcount > 0 {
		t.Logf("======================================\n")
		t.Logf("total cal leaps: %d\n", count)
		t.Logf("total base leaps: %d\n", len(leaps))
		t.Logf("miss leaps: %d\n", misscount)
		t.Logf("additional leaps: %d\n", addcount)
		t.Logf("wrong leaps: %d\n", wrongcount)
		t.Logf("good leaps: %d\n", goodcount)
	}
}

// // TestGenLunarYear list all the informations while generating the LunarYear struct
// func TestGenLunarYear(t *testing.T) {
// 	ly := GenLunarYear(2700)
// 	ly.debug()
// }
