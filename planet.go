package calendar

import "math"

//xt星体,zn坐标号,t儒略世纪数,n计算项数
func XL0_calc(xt, zn int, t, n float64) float64 {
	t /= 10 //转为儒略千年数
	var v, N float64
	var tn float64 = 1
	var F = XL0[xt]
	pn := zn*6 + 1
	var N0 = F[pn+1] - F[pn] //N0序列总数
	for i := 0; i < 6; i++ {
		n1 := F[pn+i]
		n2 := F[pn+1+i]
		n0 := n2 - n1
		if n0 == 0 {
			continue
		}
		if n < 0 { //确定项数
			N = n2
		} else {
			N = math.Floor(3*n*n0/N0+0.5) + n1
			if i != 0 {
				N += 3
			}
			if N > n2 {
				N = n2
			}
		}
		var c float64
		for j := int(n1); j < int(N); j += 3 {
			c += F[j] * math.Cos(F[j+1]+t*F[j+2])
		}
		v += c * tn
		tn *= t
	}
	v /= F[0]
	if xt == 0 { //地球
		t2 := t * t
		t3 := t2 * t //千年数的各次方
		if zn == 0 {
			v += (-0.0728 - 2.7702*t - 1.1019*t2 - 0.0996*t3) / rad
		}
		if zn == 1 {
			v += (+0.0000 + 0.0004*t + 0.0004*t2 - 0.0026*t3) / rad
		}
		if zn == 2 {
			v += (-0.0020 + 0.0044*t + 0.0213*t2 - 0.0250*t3) / 1000
		}
	} else { //其它行星
		var dv = XL0_xzb[(xt-1)*3+zn]
		if zn == 0 {
			v += -3 * t / rad
		} else if zn == 2 {
			v += dv / 100000
		} else {
			v += dv / rad
		}
	}
	return v
}

//返回冥王星J2000直角坐标
func pluto_coord(t float64) vector3 {
	var c0 = math.Pi / 180 / 100000
	var x = -1 + 2*(t*36525+1825394.5)/2185000
	var T = t / 100000000
	var v3 vector3
	for i := 0; i < 9; i++ {
		ob := XL0Pluto[i]
		var v float64
		for j := 0; j < len(ob); j += 3 {
			v += ob[j] * math.Sin(ob[j+1]*T+ob[j+2]*c0)
		}
		if i%3 == 1 {
			v *= x
		}
		if i%3 == 2 {
			v *= x * x
		}
		switch {
		case i < 3:
			v3.x += v / 100000000
		case i < 6 && i >= 3:
			v3.y += v / 100000000
		default:
			v3.z += v / 100000000
		}
	}
	v3.x += 9.922274 + 0.154154*x
	v3.y += 10.016090 + 0.064073*x
	v3.z += -3.947474 - 0.042746*x
	return v3
}

func p_coord(xt int, t, n1, n2, n3 float64) vector3 {
	var v vector3
	if xt < 8 {
		v.x = XL0_calc(xt, 0, t, n1)
		v.y = XL0_calc(xt, 1, t, n2)
		v.z = XL0_calc(xt, 2, t, n3)
	}
	if xt == 8 { //冥王星
		z := pluto_coord(t)
		r := xyz2llr(z)
		r = HDllr_J2D(t, r, "P03")
		v = vector3{r.j, r.w, r.r}
	}
	if xt == 9 { //太阳
		v.x = 0
		v.y = 0
		v.z = 0
	}
	return v
}

//返回地球坐标,t为世纪数
func e_coord(t, n1, n2, n3 float64) vector3 {
	var v3 vector3
	v3.x = XL0_calc(0, 0, t, n1)
	v3.y = XL0_calc(0, 1, t, n2)
	v3.z = XL0_calc(0, 2, t, n3)
	return v3
}
