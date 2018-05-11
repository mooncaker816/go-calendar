package calendar

import (
	"math"
)

//计算月亮
func XL1_calc(zn uint, t float64, n int) float64 {
	var ob = XL1[zn]
	var v float64
	t2 := t * t
	t3 := t2 * t
	t4 := t3 * t
	t5 := t4 * t
	tx := t - 10
	if zn == 0 {
		v += (3.81034409 + 8399.684730072*t - 3.319e-05*t2 + 3.11e-08*t3 - 2.033e-10*t4) * rad //月球平黄经(弧度)
		v += 5028.792262*t + 1.1124406*t2 + 0.00007699*t3 - 0.000023479*t4 - 0.0000000178*t5   //岁差(角秒)
		if tx > 0 {                                                                            //对公元3000年至公元5000年的拟合,最大误差小于10角秒
			v += -0.866 + 1.43*tx + 0.054*tx*tx
		}
	}
	t2 /= 1e4
	t3 /= 1e8
	t4 /= 1e8
	n *= 6
	if n < 0 {
		n = len(ob[0])
	}
	tn := 1.0
	// for i := 0; i < len(ob); i++ {
	// 	F := ob[i]
	for i, F := range ob {
		N := int(math.Floor(float64(n)*float64(len(F))/float64(len(ob[0])) + 0.5))
		if i > 0 {
			N += 6
		}
		if N >= len(F) {
			N = len(F)
		}
		var c float64
		for j := 0; j < N; j += 6 {
			c += F[j] * math.Cos(F[j+1]+t*F[j+2]+t2*F[j+3]+t3*F[j+4]+t4*F[j+5])
		}
		v += c * tn
		tn *= t
	}
	if zn != 2 {
		v /= rad
	}
	return v
}

//返回月球坐标,n1,n2,n3为各坐标所取的项数
func m_coord(t float64, n1, n2, n3 int) JW {
	var jw JW
	jw.j = XL1_calc(0, t, n1)
	jw.w = XL1_calc(1, t, n2)
	jw.r = XL1_calc(2, t, n3)
	return jw
}
