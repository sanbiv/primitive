package primitive

import (
	"fmt"
	"math"

	"github.com/fogleman/gg"
)

type Rectangle struct {
	X1, Y1 int
	X2, Y2 int
}

func NewRandomRectangle(worker *Worker) *Rectangle {
	rnd := worker.Rnd
	m := int(worker.Mutation)
	x1 := rnd.Intn(worker.W)
	y1 := rnd.Intn(worker.H)
	if worker.X != 0 {
		x1 = worker.X
	}
	if worker.Y != 0 {
		y1 = worker.Y
	}
	x2 := clampInt(x1+rnd.Intn(m)+1, 0, worker.W-1)
	y2 := clampInt(y1+rnd.Intn(m)+1, 0, worker.H-1)
	return &Rectangle{x1, y1, x2, y2}
}

func (r *Rectangle) bounds() (x1, y1, x2, y2 int) {
	x1, y1 = r.X1, r.Y1
	x2, y2 = r.X2, r.Y2
	if x1 > x2 {
		x1, x2 = x2, x1
	}
	if y1 > y2 {
		y1, y2 = y2, y1
	}
	return
}

func (r *Rectangle) Draw(dc *gg.Context, scale float64) {
	x1, y1, x2, y2 := r.bounds()
	dc.DrawRectangle(float64(x1), float64(y1), float64(x2-x1+1), float64(y2-y1+1))
	dc.Fill()
}

func (r *Rectangle) SVG(attrs string) string {
	x1, y1, x2, y2 := r.bounds()
	w := x2 - x1 + 1
	h := y2 - y1 + 1
	return fmt.Sprintf(
		"<rect %s x=\"%d\" y=\"%d\" width=\"%d\" height=\"%d\" />",
		attrs, x1, y1, w, h)
}

func (r *Rectangle) Command() string {
	return fmt.Sprintf("rectangle %d %d %d %d", r.X1, r.Y1, r.X2, r.Y2)
}

func (r *Rectangle) Copy() Shape {
	a := *r
	return &a
}

func (r *Rectangle) Scale(s float64) Shape {
	a := r.Copy().(*Rectangle)
	a.X1 = scaleInt(a.X1, s)
	a.Y1 = scaleInt(a.Y1, s)
	a.X2 = scaleInt(a.X2, s)
	a.Y2 = scaleInt(a.Y2, s)
	return a
}

func (r *Rectangle) Mutate(worker *Worker) {
	w := worker.W
	h := worker.H
	m := worker.Mutation
	rnd := worker.Rnd
	switch rnd.Intn(2) {
	case 0:
		r.X1 = clampInt(r.X1+int(rnd.NormFloat64()*m), 0, w-1)
		r.Y1 = clampInt(r.Y1+int(rnd.NormFloat64()*m), 0, h-1)
	case 1:
		r.X2 = clampInt(r.X2+int(rnd.NormFloat64()*m), 0, w-1)
		r.Y2 = clampInt(r.Y2+int(rnd.NormFloat64()*m), 0, h-1)
	}
}

func (r *Rectangle) Rasterize(worker *Worker) []Scanline {
	x1, y1, x2, y2 := r.bounds()
	lines := worker.Lines[:0]
	for y := y1; y <= y2; y++ {
		lines = append(lines, Scanline{y, x1, x2, 0xffff})
	}
	return lines
}

type RotatedRectangle struct {
	X, Y   int
	Sx, Sy int
	Angle  int
}

func NewRandomRotatedRectangle(worker *Worker) *RotatedRectangle {
	rnd := worker.Rnd
	m := int(worker.Mutation)
	x := rnd.Intn(worker.W)
	y := rnd.Intn(worker.H)
	if worker.X != 0 {
		x = worker.X
	}
	if worker.Y != 0 {
		y = worker.Y
	}
	sx := rnd.Intn(m) + 1
	sy := rnd.Intn(m) + 1
	a := rnd.Intn(360)
	r := &RotatedRectangle{x, y, sx, sy, a}
	r.Mutate(worker)
	return r
}

func (r *RotatedRectangle) Draw(dc *gg.Context, scale float64) {
	sx, sy := float64(r.Sx), float64(r.Sy)
	dc.Push()
	dc.Translate(float64(r.X), float64(r.Y))
	dc.Rotate(radians(float64(r.Angle)))
	dc.DrawRectangle(-sx/2, -sy/2, sx, sy)
	dc.Pop()
	dc.Fill()
}

func (r *RotatedRectangle) SVG(attrs string) string {
	return fmt.Sprintf(
		"<g transform=\"translate(%d %d) rotate(%d) scale(%d %d)\"><rect %s x=\"-0.5\" y=\"-0.5\" width=\"1\" height=\"1\" /></g>",
		r.X, r.Y, r.Angle, r.Sx, r.Sy, attrs)
}

func (r *RotatedRectangle) Command() string {
	return fmt.Sprintf("rotatedrectangle %d %d %d %d %d",
		r.X, r.Y, r.Sx, r.Sy, r.Angle)
}

func (r *RotatedRectangle) Copy() Shape {
	a := *r
	return &a
}

func (r *RotatedRectangle) Scale(s float64) Shape {
	a := r.Copy().(*RotatedRectangle)
	a.X = scaleInt(a.X, s)
	a.Y = scaleInt(a.Y, s)
	a.Sx = scaleInt(a.Sx, s)
	a.Sy = scaleInt(a.Sy, s)
	return a
}

func (r *RotatedRectangle) Mutate(worker *Worker) {
	w := worker.W
	h := worker.H
	m := worker.Mutation
	rnd := worker.Rnd
	switch rnd.Intn(3) {
	case 0:
		r.X = clampInt(r.X+int(rnd.NormFloat64()*m), 0, w-1)
		r.Y = clampInt(r.Y+int(rnd.NormFloat64()*m), 0, h-1)
	case 1:
		r.Sx = clampInt(r.Sx+int(rnd.NormFloat64()*m), 1, w-1)
		r.Sy = clampInt(r.Sy+int(rnd.NormFloat64()*m), 1, h-1)
	case 2:
		r.Angle = r.Angle + int(rnd.NormFloat64()*m)
	}
	// for !r.Valid() {
	// 	r.Sx = clampInt(r.Sx+int(rnd.NormFloat64()*m), 0, w-1)
	// 	r.Sy = clampInt(r.Sy+int(rnd.NormFloat64()*m), 0, h-1)
	// }
}

func (r *RotatedRectangle) Valid() bool {
	a, b := r.Sx, r.Sy
	if a < b {
		a, b = b, a
	}
	aspect := float64(a) / float64(b)
	return aspect <= 5
}

func (r *RotatedRectangle) Rasterize(worker *Worker) []Scanline {
	w := worker.W
	h := worker.H
	sx, sy := float64(r.Sx), float64(r.Sy)
	angle := radians(float64(r.Angle))
	rx1, ry1 := rotate(-sx/2, -sy/2, angle)
	rx2, ry2 := rotate(sx/2, -sy/2, angle)
	rx3, ry3 := rotate(sx/2, sy/2, angle)
	rx4, ry4 := rotate(-sx/2, sy/2, angle)
	x1, y1 := int(rx1)+r.X, int(ry1)+r.Y
	x2, y2 := int(rx2)+r.X, int(ry2)+r.Y
	x3, y3 := int(rx3)+r.X, int(ry3)+r.Y
	x4, y4 := int(rx4)+r.X, int(ry4)+r.Y
	miny := minInt(y1, minInt(y2, minInt(y3, y4)))
	maxy := maxInt(y1, maxInt(y2, maxInt(y3, y4)))
	n := maxy - miny + 1
	min := make([]int, n)
	max := make([]int, n)
	for i := range min {
		min[i] = w
	}
	xs := []int{x1, x2, x3, x4, x1}
	ys := []int{y1, y2, y3, y4, y1}
	// TODO: this could be better probably
	for i := 0; i < 4; i++ {
		x, y := float64(xs[i]), float64(ys[i])
		dx, dy := float64(xs[i+1]-xs[i]), float64(ys[i+1]-ys[i])
		count := int(math.Sqrt(dx*dx+dy*dy)) * 2
		for j := 0; j < count; j++ {
			t := float64(j) / float64(count-1)
			xi := int(x + dx*t)
			yi := int(y+dy*t) - miny
			min[yi] = minInt(min[yi], xi)
			max[yi] = maxInt(max[yi], xi)
		}
	}
	lines := worker.Lines[:0]
	for i := 0; i < n; i++ {
		y := miny + i
		if y < 0 || y >= h {
			continue
		}
		a := maxInt(min[i], 0)
		b := minInt(max[i], w-1)
		if b >= a {
			lines = append(lines, Scanline{y, a, b, 0xffff})
		}
	}
	return lines
}
