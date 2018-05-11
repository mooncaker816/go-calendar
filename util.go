package calendar

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

//将弧度转为字串,ext为小数保留位数
//tim=0输出格式示例: -23°59' 48.23"
//tim=1输出格式示例:  18h 29m 44.52s
func rad2strE(d, tim, ext float64) string {
	s := " "
	w1, w2, w3 := "°", "'", "\""
	if d < 0 {
		d = -d
		s = "-"
	}
	if tim == 1 {
		d *= 12 / math.Pi
		w1, w2, w3 = "h ", "m", "s"
	} else {
		d *= 180 / math.Pi
	}
	a := math.Floor(d)
	d = (d - a) * 60
	b := math.Floor(d)
	d = (d - b) * 60
	c := math.Floor(d)

	Q := math.Pow(10, ext)

	d = math.Floor((d-c)*Q + 0.5)
	if d >= Q {
		d -= Q
		c++
	}
	if c >= 60 {
		c -= 60
		b++
	}
	if b >= 60 {
		b -= 60
		a++
	}
	s += fmt.Sprintf("%03d%s%02d%s%02d", int(a), w1, int(b), w2, int(c))
	if ext > 0 {
		s += fmt.Sprintf(".%0*d%s", ext, d, w3)
	}
	return s
}

//将弧度转为字串,保留2位
func rad2str(d, tim float64) string { return rad2strE(d, tim, 2) }

//将弧度转为字串,精确到分
//输出格式示例: -23°59'
func rad2str2(d float64) string {
	s := "+"
	w1, w2 := "°", "'"
	if d < 0 {
		d = -d
		s = "-"
	}
	d *= 180 / math.Pi
	a := math.Floor(d)
	b := math.Floor((d-a)*60 + 0.5)
	if b >= 60 {
		b -= 60
		a++
	}
	s += fmt.Sprintf("%03d%s%02d%s", int(a), w1, int(b), w2)
	return s
}

//秒转为分秒,fx为小数点位数,fs为1转为"分秒"格式否则转为"角度分秒"格式
func m2fm(v float64, fx, flag int) string {
	s := ""
	w1, w2 := "'", "\""
	if v < 0 {
		v = -v
		s = "-"
	}
	min := math.Floor(v / 60)
	sec := v - min*60
	switch flag {
	case 1:
		w1, w2 = "分", "秒"
	case 2:
		w1, w2 = "m", "s"
	default:
	}
	return fmt.Sprintf("%s%d%s%.*f%s", s, int(min), w1, fx, sec, w2)
}

//串转弧度, f=1表示输入的s为时分秒
func str2rad(str string, f float64) (float64, error) {
	fh := 1
	if f == 1 {
		f = 15
	}
	if strings.Index(str, "-") != -1 {
		fh = -1
	}
	pat := `h|m|s|(-)|(°)|\'|\"`
	reg := regexp.MustCompile(pat)
	str = strings.TrimSpace(reg.ReplaceAllString(str, " "))
	sl := strings.Fields(str)
	h, err := strconv.Atoi(sl[0])
	if err != nil {
		return 0, fmt.Errorf("can not parse string: %s %v", sl[0], err)
	}
	m, err := strconv.Atoi(sl[1])
	if err != nil {
		return 0, fmt.Errorf("can not parse string: %s %v", sl[1], err)
	}
	s, err := strconv.Atoi(sl[2])
	if err != nil {
		return 0, fmt.Errorf("can not parse string: %s %v", sl[2], err)
	}
	return float64(fh*(h*3600+m*60+s*1)) / rad * f, nil
}

//对超过0-2Pi的角度转为0-2Pi
func rad2mrad(v float64) float64 {
	if v <= 2*math.Pi && v >= 0 {
		return v
	}
	step := -2 * math.Pi
	if v < 0 {
		step = -step
	}
	for v > 2*math.Pi || v < 0 {
		v += step
	}
	return v
}

//对超过-Pi到Pi的角度转为-Pi到Pi
func rad2rrad(v float64) float64 {
	if v <= math.Pi && v >= -math.Pi {
		return v
	}
	step := -2 * math.Pi
	if v < -math.Pi {
		step = -step
	}
	for v > math.Pi || v < -math.Pi {
		v += step
	}
	return v
}

//临界余数(a与最近的整倍数b相差的距离)
func mod2(a, b float64) float64 {
	c := a / b
	c -= math.Floor(c)
	if c > 0.5 {
		c -= 1
	}
	return c * b
}

//球面转直角坐标
func llr2xyz(jw JW) vector3 {
	var v3 vector3
	v3.x = jw.r * math.Cos(jw.w) * math.Cos(jw.j)
	v3.y = jw.r * math.Cos(jw.w) * math.Sin(jw.j)
	v3.z = jw.r * math.Sin(jw.w)
	return v3
}

//直角坐标转球
func xyz2llr(v3 vector3) JW {
	var jw JW
	jw.r = math.Sqrt(v3.x*v3.x + v3.y*v3.y + v3.z*v3.z)
	jw.w = math.Asin(v3.z / jw.r)
	jw.j = rad2mrad(math.Atan2(v3.y, v3.x))
	return jw
}

//球面坐标旋转
//黄道赤道坐标变换,赤到黄E取负
func llrConv(jw JW, E float64) JW {
	var njw JW
	njw.j = math.Atan2(math.Sin(jw.j)*math.Cos(E)-math.Tan(jw.w)*math.Sin(E), math.Cos(jw.j))
	njw.w = math.Asin(math.Cos(E)*math.Sin(jw.w) + math.Sin(E)*math.Cos(jw.w)*math.Sin(jw.j))
	njw.r = jw.r
	njw.j = rad2mrad(njw.j)
	return njw
}

//赤道坐标转为地平坐标
//转到相对于地平赤道分点的赤道坐标
func CD2DP(jw JW, L, fa, gst float64) JW {
	var njw JW
	njw.j = jw.j + math.Pi/2 - gst - L
	njw.w = jw.w
	njw.r = jw.r
	njw = llrConv(njw, math.Pi/2-fa)
	njw.j = rad2mrad(math.Pi/2 - njw.j)
	return njw
}

//求角度差
func j1_j2(J1, W1, J2, W2 float64) float64 {
	dJ := rad2rrad(J1 - J2)
	dW := W1 - W2
	if math.Abs(dJ) < 1/1000 && math.Abs(dW) < 1/1000 {
		dJ *= math.Cos((W1 + W2) / 2)
		return math.Sqrt(dJ*dJ + dW*dW)
	}
	return math.Acos(math.Sin(W1)*math.Sin(W2) + math.Cos(W1)*math.Cos(W2)*math.Cos(dJ))
}

//日心球面转地心球面,Z星体球面坐标,A地球球面坐标
//本函数是通用的球面坐标中心平移函数,行星计算中将反复使用
func h2g(z, a JW) JW {
	aa := llr2xyz(a) //地球
	zz := llr2xyz(z) //星体
	zz.x -= aa.x
	zz.y -= aa.y
	zz.z -= aa.z
	return xyz2llr(zz)
}

//视差角(不是视差)
func shiChaJ(gst, L, fa, J, W float64) float64 {
	H := gst + L - J //天体的时角
	return rad2mrad(math.Atan2(math.Sin(H), math.Tan(fa)*math.Cos(W)-math.Sin(W)*math.Cos(H)))
}

//=============================一些天文基本问题=====================================
//==================================================================================
//返回朔日的编号,jd应在朔日附近，允许误差数天
func suoN(jd JulianDate) float64 { return math.Floor((float64(jd) + 8) / 29.5306) }

//太阳光行差,t是世纪数
func gxc_sunLon(t float64) float64 {
	v := -0.043126 + 628.301955*t - 0.000002732*t*t //平近点角
	e := 0.016708634 - 0.000042037*t - 0.0000001267*t*t
	return (-20.49552 * (1 + e*math.Cos(v))) / rad //黄经光行差
}

//黄纬光行差
func gxc_sunLat(t float64) float64 { return 0 }

//月球经度光行差,误差0.07"
func gxc_moonLon(t float64) float64 { return -3.4E-6 }

//月球纬度光行差,误差0.006"
func gxc_moonLat(t float64) float64 { return 0.063 * math.Sin(0.057+8433.4662*t+0.000064*t*t) / rad }

//传入T是2000年首起算的日数(UT),dt是deltatT(日),精度要求不高时dt可取值为0
//返回格林尼治平恒星时(不含赤经章动及非多项式部分),即格林尼治子午圈的平春风点起算的赤经
func pGST(T, dt float64) float64 {
	t := (T + dt) / 36525
	t2 := t * t
	t3 := t2 * t
	t4 := t3 * t
	return pi2*(0.7790572732640+1.00273781191135448*T) + //T是UT,下一行的t是力学时(世纪数)
		(0.014506+4612.15739966*t+1.39667721*t2-0.00009344*t3+0.00001882*t4)/rad
}

//传入力学时J2000起算日数，返回平恒星时
func pGST2(jd float64) float64 {
	var dt = dt_T(jd)
	return pGST(jd-dt, dt)
}

//太阳升降计算。jd儒略日(须接近L当地平午UT)，L地理经度，fa地理纬度，sj=-1升,sj=1降
func sunShengJ(jd, L, fa, sj float64) float64 {
	jd = math.Floor(jd+0.5) - L/pi2
	for i := 0; i < 2; i++ {
		T := jd / 36525
		E := (84381.4060 - 46.836769*T) / rad        //黄赤交角
		t := T + (32*(T+1.8)*(T+1.8)-20)/86400/36525 //儒略世纪年数,力学时
		J := (48950621.66 + 6283319653.318*t + 53*t*t - 994 +
			334166*math.Cos(4.669257+628.307585*t) +
			3489*math.Cos(4.6261+1256.61517*t) +
			2060.6*math.Cos(2.67823+628.307585*t)*t) / 10000000
		//太阳黄经以及它的正余弦值
		sinJ := math.Sin(J)
		cosJ := math.Cos(J)
		//恒星时(子午圈位置)
		gst := (0.7790572732640+1.00273781191135448*jd)*pi2 + (0.014506+4612.15739966*T+1.39667721*T*T)/rad
		A := math.Atan2(sinJ*math.Cos(E), cosJ) //太阳赤经
		D := math.Asin(math.Sin(E) * sinJ)      //太阳赤纬
		cosH0 := (math.Sin(-50*60/rad) - math.Sin(fa)*math.Sin(D)) / (math.Cos(fa) * math.Cos(D))
		//太阳在地平线上的cos(时角)计算
		if math.Abs(cosH0) >= 1 {
			return 0
		}
		//(升降时角-太阳时角)/太阳速度
		jd += rad2rrad(sj*math.Acos(cosH0)-(gst+L-A)) / 6.28
	}
	return jd //反回格林尼治UT
}

//时差计算(高精度),t力学时儒略世纪数
func pty_zty(t float64) float64 {
	t2 := t * t
	t3 := t2 * t
	t4 := t3 * t
	t5 := t4 * t
	L := (1753470142+628331965331.8*t+5296.74*t2+0.432*t3-0.1124*t4-0.00009*t5)/1000000000 + math.Pi - 20.5/rad

	//   var f,z=new Array();
	dL := -17.2 * math.Sin(2.1824-33.75705*t) / rad //黄经章
	dE := 9.2 * math.Cos(2.1824-33.75705*t) / rad   //交角章
	E := hcjj(t) + dE                               //真黄赤交角

	var z JW
	//地球坐标
	z.j = XL0_calc(0, 0, t, 50) + math.Pi + gxc_sunLon(t) + dL
	z.w = -(2796*math.Cos(3.1987+8433.46616*t) + 1016*math.Cos(5.4225+550.75532*t) + 804*math.Cos(3.88+522.3694*t)) / 1000000000

	z = llrConv(z, E) //z太阳地心赤道坐标
	z.j -= dL * math.Cos(E)

	L = rad2rrad(L - z.j)
	return L / pi2 //单位是周(天)
}

// func pty_zty2(t){ //时差计算(低精度),误差约在1秒以内,t力学时儒略世纪数
//   var L = ( 1753470142 + 628331965331.8*t + 5296.74*t*t )/1000000000 + math.PI;
//   var z=new Array();
//   var E= (84381.4088 -46.836051*t)/rad;
//   z[0]=XL0_calc(0,0,t,5)+math.PI, z[1]=0; //地球坐标
//   z = llrConv( z, E ); //z太阳地心赤道坐标
//   L = rad2rrad(L-z[0]);
//   return L/pi2; //单位是周(天)
// }

//=================================deltat T计算=====================================
//==================================================================================
// TD - UT1 计算表
var dt_at = []float64{
	-4000, 108371.7, -13036.80, 392.000, 0.0000,
	-500, 17201.0, -627.82, 16.170, -0.3413,
	-150, 12200.6, -346.41, 5.403, -0.1593,
	150, 9113.8, -328.13, -1.647, 0.0377,
	500, 5707.5, -391.41, 0.915, 0.3145,
	900, 2203.4, -283.45, 13.034, -0.1778,
	1300, 490.1, -57.35, 2.085, -0.0072,
	1600, 120.0, -9.81, -1.532, 0.1403,
	1700, 10.2, -0.91, 0.510, -0.0370,
	1800, 13.4, -0.72, 0.202, -0.0193,
	1830, 7.8, -1.81, 0.416, -0.0247,
	1860, 8.3, -0.13, -0.406, 0.0292,
	1880, -5.4, 0.32, -0.183, 0.0173,
	1900, -2.3, 2.06, 0.169, -0.0135,
	1920, 21.2, 1.69, -0.304, 0.0167,
	1940, 24.2, 1.22, -0.064, 0.0031,
	1960, 33.2, 0.51, 0.231, -0.0109,
	1980, 51.0, 1.29, -0.026, 0.0032,
	2000, 63.87, 0.1, 0, 0,
	2005, 64.7, 0.4, 0, 0, //一次项记为x,则 10x=0.4秒/年*(2015-2005),解得x=0.4
	2015, 69,
}

//二次曲线外推
func dt_ext(y, jsd float64) float64 {
	dy := (y - 1820) / 100
	return -20 + jsd*dy*dy
}

func dt_calc(y float64) float64 { //计算世界时与原子时之差,传入年
	y0 := dt_at[len(dt_at)-2] //表中最后一年
	t0 := dt_at[len(dt_at)-1] //表中最后一年的deltatT
	if y >= y0 {
		var jsd float64 = 31 //sjd是y1年之后的加速度估计。瑞士星历表jsd=31,NASA网站jsd=32,skmap的jsd=29
		if y > y0+100 {
			return dt_ext(y, jsd)
		}
		v := dt_ext(y, jsd)        //二次曲线外推
		dv := dt_ext(y0, jsd) - t0 //ye年的二次外推与te的差
		return v - dv*(y0+100-y)/100
	}
	d := dt_at
	i := 0
	for i = 0; i < len(d); i += 5 {
		if y < d[i+5] {
			break
		}
	}
	t1 := (y - d[i]) / (d[i+5] - d[i]) * 10
	t2 := t1 * t1
	t3 := t2 * t1
	return d[i+1] + d[i+2]*t1 + d[i+3]*t2 + d[i+4]*t3
}

func dt_T(t float64) float64 {
	return dt_calc(t/365.2425+2000) / 86400.0
} //传入儒略日(J2000起算),计算TD-UT(单位:日)
