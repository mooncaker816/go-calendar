package calendar

type vector3 struct {
	x float64
	y float64
	z float64
}

type JW struct {
	j float64
	w float64
	r float64
}

//中升降时刻
type zsjsk struct {
	z float64 //
	x float64 //
	s float64 //
	j float64 //
	c float64 //
	h float64 //
}

type zdxxsjcjw struct {
	H  float64
	H1 float64
	H2 float64
	H3 float64
	H4 float64
}

type tyszj struct {
	sm string
	s  float64 //升起
	j  float64 //降落
	c  float64 //民用晨
	h  float64 //民用昏
	c2 float64 //航海晨
	h2 float64 //航海昏
	c3 float64 //天文晨
	h3 float64 //天文昏
	z  float64 //中天
	x  float64 //下中天
}

type cs_szj struct {
	s  float64 //升
	z  float64 //中
	j  float64 //降
	c  float64 //晨
	h  float64 //昏
	ch float64 //光照时间,timeStr()内部+0.5,所以这里补上-0.5
	sj float64 //昼长
	Ms float64 //月亮升
	Mz float64 //月亮中
	Mj float64 //月亮降
}
