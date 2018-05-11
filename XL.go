package calendar

import "math"

//日月星历基本函数类
//=====================
//星历函数(日月球面坐标计算)
//地球经度计算,返回Date分点黄经,传入世纪数、取项数
func E_Lon(t, n float64) float64 { return XL0_calc(0, 0, t, n) }

//月球经度计算,返回Date分点黄经,传入世纪数,n是项数比例
func M_Lon(t, n float64) float64 { return XL1_calc(0, t, int(n)) }

//=========================
//地球速度,t是世纪数,误差小于万分3
func E_v(t float64) float64 {
	var f = 628.307585 * t
	return 628.332 + 21*math.Sin(1.527+f) + 0.44*math.Sin(1.48+f*2) +
		0.129*math.Sin(5.82+f)*t + 0.00055*math.Sin(4.21+f)*t*t
}

//月球速度计算,传入世经数
func M_v(t float64) float64 {
	var v = 8399.71 - 914*math.Sin(0.7848+8328.691425*t+0.0001523*t*t) //误差小于5%
	v -= 179*math.Sin(2.543+15542.7543*t) +                            //误差小于0.3%
		160*math.Sin(0.1874+7214.0629*t) +
		62*math.Sin(3.14+16657.3828*t) +
		34*math.Sin(4.827+16866.9323*t) +
		22*math.Sin(4.9+23871.4457*t) +
		12*math.Sin(2.59+14914.4523*t) +
		7*math.Sin(0.23+6585.7609*t) +
		5*math.Sin(0.9+25195.624*t) +
		5*math.Sin(2.32-7700.3895*t) +
		5*math.Sin(3.88+8956.9934*t) +
		5*math.Sin(0.49+7771.3771*t)
	return v
}

//=========================
//月日视黄经的差值
func MS_aLon(t, Mn, Sn float64) float64 {
	return M_Lon(t, Mn) + gxc_moonLon(t) - (E_Lon(t, Sn) + gxc_sunLon(t) + math.Pi)
}

//太阳视黄经
func S_aLon(t, n float64) float64 {
	return E_Lon(t, n) + nutationLon2(t) + gxc_sunLon(t) + math.Pi //注意，这里的章动计算很耗时
}

//=========================
//已知地球真黄经求时间
func E_Lon_t(W float64) float64 {
	var v = 628.3319653318
	t := (W - 1.75347) / v
	v = E_v(t) //v的精度0.03%，详见原文
	t += (W - E_Lon(t, 10)) / v
	v = E_v(t) //再算一次v有助于提高精度,不算也可以
	t += (W - E_Lon(t, -1)) / v
	return t
}

//已知真月球黄经求时间
func M_Lon_t(W float64) float64 {
	var v = 8399.70911033384
	t := (W - 3.81034) / v
	t += (W - M_Lon(t, 3)) / v
	v = M_v(t) //v的精度0.5%，详见原文
	t += (W - M_Lon(t, 20)) / v
	t += (W - M_Lon(t, -1)) / v
	return t
}

//已知月日视黄经差求时间
func MS_aLon_t(W float64) float64 {
	var v = 7771.37714500204
	t := (W + 1.08472) / v
	t += (W - MS_aLon(t, 3, 3)) / v
	v = M_v(t) - E_v(t) //v的精度0.5%，详见原文
	t += (W - MS_aLon(t, 20, 10)) / v
	t += (W - MS_aLon(t, -1, 60)) / v
	return t
}

//已知太阳视黄经反求时间
func S_aLon_t(W float64) float64 {
	var v = 628.3319653318
	t := (W - 1.75347 - math.Pi) / v
	v = E_v(t) //v的精度0.03%，详见原文
	t += (W - S_aLon(t, 10)) / v
	v = E_v(t) //再算一次v有助于提高精度,不算也可以
	t += (W - S_aLon(t, -1)) / v
	return t
}

/****
MS_aLon_t1:func(W){ //已知月日视黄经差求时间,高速低精度,误差不超过40秒
  var t,v = 7771.37714500204;
  t  = ( W + 1.08472               )/v;
  t += ( W - this.MS_aLon(t, 3, 3) )/v;  v=this.M_v(t)-this.E_v(t);  //v的精度0.5%，详见原文
  t += ( W - this.MS_aLon(t,50,20) )/v;
  return t;
},
S_aLon_t1:func(W){ //已知太阳视黄经反求时间,高速低精度,最大误差不超过50秒,平均误差15秒
  var t,v= 628.3319653318;
  t  = ( W - 1.75347-math.Pi   )/v; v = 628.332 + 21*math.Sin( 1.527+628.307585*t );
  t += ( W - this.S_aLon(t,3) )/v;
  t += ( W - this.S_aLon(t,40))/v;
  return t;
},
****/
//已知月日视黄经差求时间,高速低精度,误差不超过600秒(只验算了几千年)
func MS_aLon_t2(W float64) float64 {
	var v = 7771.37714500204
	t := (W + 1.08472) / v
	var t2 = t * t
	t -= (-0.00003309*t2 + 0.10976*math.Cos(0.784758+8328.6914246*t+0.000152292*t2) + 0.02224*math.Cos(0.18740+7214.0628654*t-0.00021848*t2) - 0.03342*math.Cos(4.669257+628.307585*t)) / v
	L := M_Lon(t, 20) - (4.8950632 + 628.3319653318*t + 0.000005297*t*t + 0.0334166*math.Cos(4.669257+628.307585*t) + 0.0002061*math.Cos(2.67823+628.307585*t)*t + 0.000349*math.Cos(4.6261+1256.61517*t) - 20.5/rad)
	v = 7771.38 - 914*math.Sin(0.7848+8328.691425*t+0.0001523*t*t) - 179*math.Sin(2.543+15542.7543*t) - 160*math.Sin(0.1874+7214.0629*t)
	t += (W - L) / v
	return t
}

//已知太阳视黄经反求时间,高速低精度,最大误差不超过600秒
func S_aLon_t2(W float64) float64 {
	var v = 628.3319653318
	t := (W - 1.75347 - math.Pi) / v
	t -= (0.000005297*t*t + 0.0334166*math.Cos(4.669257+628.307585*t) + 0.0002061*math.Cos(2.67823+628.307585*t)*t) / v
	t += (W - E_Lon(t, 8) - math.Pi + (20.5+17.2*math.Sin(2.1824-33.75705*t))/rad) / v
	return t
}

//月亮被照亮部分的比例
func moonIll(t float64) float64 {
	t2 := t * t
	t3 := t2 * t
	t4 := t3 * t
	dm := math.Pi / 180
	D := (297.8502042 + 445267.1115168*t - 0.0016300*t2 + t3/545868 - t4/113065000) * dm //日月平距角
	M := (357.5291092 + 35999.0502909*t - 0.0001536*t2 + t3/24490000) * dm               //太阳平近点
	m := (134.9634114 + 477198.8676313*t + 0.0089970*t2 + t3/69699 - t4/14712000) * dm   //月亮平近点
	a := math.Pi - D + (-6.289*math.Sin(m)+2.100*math.Sin(M)-1.274*math.Sin(D*2-m)-0.658*math.Sin(D*2)-0.214*math.Sin(m*2)-0.110*math.Sin(D))*dm
	return (1 + math.Cos(a)) / 2
}

//转入地平纬度及地月质心距离,返回站心视半径(角秒)
func moonRad(r, h float64) float64 {
	return cs_sMoon / r * (1 + math.Sin(h)*cs_rEar/r)
}

//求月亮近点时间和距离,t为儒略世纪数力学时
func moonMinR(t, min float64) vector3 {
	var a = 27.55454988 / 36525
	var b float64
	if min != 0 {
		b = -10.3302 / 36525
	} else {
		b = 3.4471 / 36525
	}
	t = b + a*math.Floor((t-b)/a+0.5) //平近(远)点时间
	//初算二次
	var dt float64 = 1 / 36525
	r1 := XL1_calc(2, t-dt, 10)
	r2 := XL1_calc(2, t, 10)
	r3 := XL1_calc(2, t+dt, 10)
	t += (r1 - r3) / (r1 + r3 - 2*r2) * dt / 2
	dt = 0.5 / 36525
	r1 = XL1_calc(2, t-dt, 20)
	r2 = XL1_calc(2, t, 20)
	r3 = XL1_calc(2, t+dt, 20)
	t += (r1 - r3) / (r1 + r3 - 2*r2) * dt / 2
	//精算
	dt = 1200 / 86400 / 36525
	r1 = XL1_calc(2, t-dt, -1)
	r2 = XL1_calc(2, t, -1)
	r3 = XL1_calc(2, t+dt, -1)
	t += (r1 - r3) / (r1 + r3 - 2*r2) * dt / 2
	r2 += (r1 - r3) / (r1 + r3 - 2*r2) * (r3 - r1) / 8
	return vector3{t, r2, 0}
}

//月亮升交点
func moonNode(t, asc float64) vector3 {
	var a = 27.21222082 / 36525
	var b float64
	if asc != 0 {
		b = 21 / 36525
	} else {
		b = 35 / 36525
	}
	t = b + a*math.Floor((t-b)/a+0.5) //平升(降)交点时间
	dt := 0.5 / 36525
	w := XL1_calc(1, t, 10)
	w2 := XL1_calc(1, t+dt, 10)
	v := (w2 - w) / dt
	t -= w / v
	dt = 0.05 / 36525
	w = XL1_calc(1, t, 40)
	w2 = XL1_calc(1, t+dt, 40)
	v = (w2 - w) / dt
	t -= w / v
	w = XL1_calc(1, t, -1)
	t -= w / v
	return vector3{t, XL1_calc(0, t, -1), 0}
}

//地球近远点
func earthMinR(t, min float64) vector3 {
	var a = 365.25963586 / 36525
	var b float64
	if min != 0 {
		b = 1.7 / 36525
	} else {
		b = 184.5 / 36525
	}
	t = b + a*math.Floor((t-b)/a+0.5) //平近(远)点时间
	var r1, r2, r3, dt float64
	//初算二次
	dt = 3 / 36525
	r1 = XL0_calc(0, 2, t-dt, 10)
	r2 = XL0_calc(0, 2, t, 10)
	r3 = XL0_calc(0, 2, t+dt, 10)
	t += (r1 - r3) / (r1 + r3 - 2*r2) * dt / 2 //误差几个小时
	dt = 0.2 / 36525
	r1 = XL0_calc(0, 2, t-dt, 80)
	r2 = XL0_calc(0, 2, t, 80)
	r3 = XL0_calc(0, 2, t+dt, 80)
	t += (r1 - r3) / (r1 + r3 - 2*r2) * dt / 2 //误差几分钟
	//精算
	dt = 0.01 / 36525
	r1 = XL0_calc(0, 2, t-dt, -1)
	r2 = XL0_calc(0, 2, t, -1)
	r3 = XL0_calc(0, 2, t+dt, -1)
	t += (r1 - r3) / (r1 + r3 - 2*r2) * dt / 2 //误差小于秒
	r2 += (r1 - r3) / (r1 + r3 - 2*r2) * (r3 - r1) / 8
	return vector3{t, r2, 0}
}
