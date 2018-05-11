package calendar

import (
	"math"
)

// SZJ 日月的升中天降,不考虑气温和气压的影响
type SZJ struct {
	L   float64 //站点地理经度,向东测量为正
	fa  float64 //站点地理纬度
	dt  float64 //TD-UT
	E   float64 //黄赤交角
	rts []cs_szj
}

//时角及赤经纬
type sjcjw struct {
	H  float64 //天体时角
	H0 float64 //升起对应的时角
}

//h地平纬度,w赤纬,返回时角
func (szj *SZJ) getH(h, w float64) float64 {
	var c = (math.Sin(h) - math.Sin(szj.fa)*math.Sin(w)) / math.Cos(szj.fa) / math.Cos(w)
	if math.Abs(c) > 1 {
		return math.Pi
	}
	return math.Acos(c)
}

// Mcoord 章动同时影响恒星时和天体坐标,所以不计算章动。返回时角及赤经纬
// H0 == 0,不反回:升起对应的时角H0
// H0 == 1,  反回:升起对应的时角H0
func (szj *SZJ) Mcoord(jd float64, H0 int, r *sjcjw) {
	z := m_coord((jd+szj.dt)/36525, 40, 30, 8)      //低精度月亮赤经纬
	nz := llrConv(z, szj.E)                         //转为赤道坐标
	r.H = rad2rrad(pGST(jd, szj.dt) + szj.L - nz.j) //得到此刻天体时角
	if H0 == 1 {
		r.H0 = szj.getH(0.7275*cs_rEar/z.r-34*60/rad, z.w) //升起对应的时角
	}
}

// Mt 月亮到中升降时刻计算,传入jd含义与St()函数相同
func (szj *SZJ) Mt(jd float64) zsjsk {
	var r zsjsk
	szj.dt = dt_T(jd)
	szj.E = hcjj(jd / 36525)
	jd -= mod2(0.1726222+0.966136808032357*jd-0.0366*szj.dt+szj.L/pi2, 1) //查找最靠近当日中午的月上中天,mod2的第1参数为本地时角近似值
	rr := new(sjcjw)
	sv := pi2 * 0.966
	r.z, r.x, r.s, r.j, r.c, r.h = jd, jd, jd, jd, jd, jd
	szj.Mcoord(jd, 1, rr) //月亮坐标
	r.s += (-rr.H0 - rr.H) / sv
	r.j += (rr.H0 - rr.H) / sv
	r.z += (0 - rr.H) / sv
	r.x += (math.Pi - rr.H) / sv
	szj.Mcoord(r.s, 1, rr)
	r.s += rad2rrad(-rr.H0-rr.H) / sv
	szj.Mcoord(r.j, 1, rr)
	r.j += rad2rrad(+rr.H0-rr.H) / sv
	szj.Mcoord(r.z, 0, rr)
	r.z += rad2rrad(0-rr.H) / sv
	szj.Mcoord(r.x, 0, rr)
	r.x += rad2rrad(math.Pi-rr.H) / sv
	return r
}

//章动同时影响恒星时和天体坐标,所以不计算章动。返回时角及赤经纬
func (szj *SZJ) Scoord(jd, xm float64) zdxxsjcjw {
	var r zdxxsjcjw
	var z = JW{E_Lon((jd+szj.dt)/36525, 5) + math.Pi - 20.5/rad, 0, 1} //太阳坐标(修正了光行差)
	z = llrConv(z, szj.E)                                              //转为赤道坐标
	r.H = rad2rrad(pGST(jd, szj.dt) + szj.L - z.j)                     //得到此刻天体时角

	if xm == 10 || xm == 1 {
		r.H1 = szj.getH(-50*60/rad, z.w)
	} //地平以下50分
	if xm == 10 || xm == 2 {
		r.H2 = szj.getH(-6*3600/rad, z.w)
	} //地平以下6度
	if xm == 10 || xm == 3 {
		r.H3 = szj.getH(-12*3600/rad, z.w)
	} //地平以下12度
	if xm == 10 || xm == 4 {
		r.H4 = szj.getH(-18*3600/rad, z.w)
	} //地平以下18度
	return r
}

//太阳到中升降时刻计算,传入jd是当地中午12点时间对应的2000年首起算的格林尼治时间UT
func (szj *SZJ) St(jd float64) tyszj {
	var r tyszj
	szj.dt = dt_T(jd)
	szj.E = hcjj(jd / 36525)
	jd -= mod2(jd+szj.L/pi2, 1) //查找最靠近当日中午的日上中天,mod2的第1参数为本地时角近似值

	//  var r = new Array(),
	sv := pi2
	r.z, r.x, r.s, r.j, r.c, r.h, r.c2, r.h2, r.c3, r.h3 = jd, jd, jd, jd, jd, jd, jd, jd, jd, jd
	r.sm = " "
	t := szj.Scoord(jd, 10)   //太阳坐标
	r.s += (-t.H1 - t.H) / sv //升起
	r.j += (t.H1 - t.H) / sv  //降落

	r.c += (-t.H2 - t.H) / sv  //民用晨
	r.h += (t.H2 - t.H) / sv   //民用昏
	r.c2 += (-t.H3 - t.H) / sv //航海晨
	r.h2 += (t.H3 - t.H) / sv  //航海昏
	r.c3 += (-t.H4 - t.H) / sv //天文晨
	r.h3 += (t.H4 - t.H) / sv  //天文昏

	r.z += (0 - t.H) / sv       //中天
	r.x += (math.Pi - t.H) / sv //下中天
	t = szj.Scoord(r.s, 1)
	r.s += rad2rrad(-t.H1-t.H) / sv
	if t.H1 == math.Pi {
		r.sm += "无升起."
	}
	t = szj.Scoord(r.j, 1)
	r.j += rad2rrad(+t.H1-t.H) / sv
	if t.H1 == math.Pi {
		r.sm += "无降落."
	}
	t = szj.Scoord(r.c, 2)
	r.c += rad2rrad(-t.H2-t.H) / sv
	if t.H2 == math.Pi {
		r.sm += "无民用晨."
	}
	t = szj.Scoord(r.h, 2)
	r.h += rad2rrad(+t.H2-t.H) / sv
	if t.H2 == math.Pi {
		r.sm += "无民用昏."
	}
	t = szj.Scoord(r.c2, 3)
	r.c2 += rad2rrad(-t.H3-t.H) / sv
	if t.H3 == math.Pi {
		r.sm += "无航海晨."
	}
	t = szj.Scoord(r.h2, 3)
	r.h2 += rad2rrad(+t.H3-t.H) / sv
	if t.H3 == math.Pi {
		r.sm += "无航海昏."
	}
	t = szj.Scoord(r.c3, 4)
	r.c3 += rad2rrad(-t.H4-t.H) / sv
	if t.H4 == math.Pi {
		r.sm += "无天文晨."
	}
	t = szj.Scoord(r.h3, 4)
	r.h3 += rad2rrad(+t.H4-t.H) / sv
	if t.H4 == math.Pi {
		r.sm += "无天文昏."
	}

	t = szj.Scoord(r.z, 0)
	r.z += (0 - t.H) / sv
	t = szj.Scoord(r.x, 0)
	r.x += rad2rrad(math.Pi-t.H) / sv
	return r
}

// rts:new Array(),//多天的升中降
//多天升中降计算,jd是当地起始略日(中午时刻),sq是时区
func (szj *SZJ) calcRTS(jd float64, n int, Jdl, Wdl, sq float64) {
	//  var c,r;
	if szj.rts == nil {
		szj.rts = make([]cs_szj, 31)
	}
	//设置站点参数
	szj.L = Jdl
	szj.fa = Wdl
	sq /= 24
	for i := -1; i <= n; i++ {
		if i >= 0 && i < n { //太阳
			r1 := szj.St(jd + float64(i) + sq)
			szj.rts[i].s = r1.s - sq          //升
			szj.rts[i].z = r1.z - sq          //中
			szj.rts[i].j = r1.j - sq          //降
			szj.rts[i].c = r1.c - sq          //晨
			szj.rts[i].h = r1.h - sq          //昏
			szj.rts[i].ch = r1.h - r1.c - 0.5 //光照时间,timeStr()内部+0.5,所以这里补上-0.5
			szj.rts[i].sj = r1.j - r1.s - 0.5 //昼长
		}
		r2 := szj.Mt(jd + float64(i) + sq) //月亮

		c := int(math.Floor(r2.s-sq+0.5) - jd)
		if c >= 0 && c < n {
			szj.rts[c].Ms = r2.s - sq
		}
		c = int(math.Floor(r2.z-sq+0.5) - jd)
		if c >= 0 && c < n {
			szj.rts[c].Mz = r2.z - sq
		}
		c = int(math.Floor(r2.j-sq+0.5) - jd)
		if c >= 0 && c < n {
			szj.rts[c].Mj = r2.j - sq
		}
	}
	//  szj.rts.dn = n;
}
