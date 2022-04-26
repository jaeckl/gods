package datatable

import "math"

type Slice struct {
	size   int
	fields map[int]*Cell
}

func (sl *Slice) Size() int {
	return sl.size
}

func (sl *Slice) Append(value interface{}) *Slice {
	sl2 := *sl
	sl2.fields[sl2.size] = &Cell{value: value}
	sl2.size++
	return &sl2
}

func (sl *Slice) Concate(sl2 *Slice) {
	for k, v := range sl2.fields {
		sl.fields[sl.size+k-1] = v
	}
	sl.size += sl2.size
}

func (sl *Slice) Get(index int) *Cell {
	if index < 0 || index >= sl.size {
		panic("index out of bounds")
	}
	return sl.fields[index]
}

func (sl *Slice) Region(interval [2]int) *Slice {
	if interval[0] < 0 || interval[1] > sl.size || interval[0] > interval[1] {
		panic("invalid interval")
	}
	slice := &Slice{
		size:   interval[1] - interval[0],
		fields: make(map[int]*Cell),
	}
	for i := interval[0]; i < interval[1]; i++ {
		slice.fields[i-interval[0]] = sl.fields[i]
	}
	return slice
}

func (sl *Slice) Max() float64 {
	max := math.Inf(-1)
	for _, v := range sl.fields {
		if v.value.(float64) > max {
			max = v.value.(float64)
		}
	}
	return max
}

func (sl *Slice) Min() float64 {
	min := math.Inf(1)
	for _, v := range sl.fields {
		if v.value.(float64) < min {
			min = v.value.(float64)
		}
	}
	return min
}

func (sl *Slice) Mean() float64 {
	sum := 0.0
	for _, v := range sl.fields {
		sum += v.value.(float64)
	}
	return sum / float64(sl.size)
}

func (sl *Slice) Var() float64 {
	var sum float64
	for _, v := range sl.fields {
		sum += v.value.(float64)
	}
	mean := sum / float64(sl.size)
	var sum2 float64
	for _, v := range sl.fields {
		sum2 += math.Pow(v.value.(float64)-mean, 2)
	}
	return sum2 / float64(sl.size)
}

func (sl *Slice) Std() float64 {
	return math.Sqrt(sl.Var())
}

// some comments
func (sl *Slice) Mode() float64 {
	m := make(map[float64]int)
	for _, v := range sl.fields {
		m[v.value.(float64)]++
	}
	max := 0
	var mode float64
	for k, v := range m {
		if v > max {
			max = v
			mode = k
		}
	}
	return mode
}

func (sl *Slice) Quartile(q float64) float64 {
	if q < 0 || q > 1 {
		panic("invalid quartile")
	}
	if q == 0 {
		return sl.Min()
	}
	if q == 1 {
		return sl.Max()
	}
	return sl.Min() + (sl.Max()-sl.Min())*q
}
