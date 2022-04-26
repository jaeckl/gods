package datatable

import "math"

type Alignment = int

const (
	Rows = iota
	Cols = iota
)

type Cell struct {
	value interface{}
}

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

type DataTable struct {
	shape         [2]int
	columns       map[int]*Slice
	columnNameIds map[string]int
	selection     Selection
}

func FromMap(m map[string][]interface{}) *DataTable {
	dt := &DataTable{
		columns:       make(map[int]*Slice),
		columnNameIds: make(map[string]int),
		shape:         [2]int{len(m), 0},
	}
	i := 0
	for k, v := range m {
		dt.shape[1] = len(v)
		dt.columns[i] = &Slice{
			size:   len(v),
			fields: make(map[int]*Cell),
		}
		for j, val := range v {
			dt.columns[i].fields[j] = &Cell{
				value: val,
			}
		}
		dt.columnNameIds[k] = i
		i++
	}
	return dt
}

func _validate_colums_unique(a, b map[int]*Slice) {
	for k := range a {
		if _, ok := b[k]; ok {
			panic("column already exists")
		}
	}
}

func _validate_equal_columns(a, b map[int]*Slice) {
	if len(a) != len(b) {
		panic("columns count mismatch")
	}
	for k := range a {
		if _, ok := b[k]; !ok {
			panic("columns name mismatch")
		}
	}
}

// DataTable Shape Manipulation
func (dt *DataTable) Concate(dt2 *DataTable, align Alignment) {
	if align == Cols {
		_validate_colums_unique(dt.columns, dt2.columns)
		for k, v := range dt2.columns {
			dt.columns[k] = v
		}
	}
	if align == Rows {
		_validate_equal_columns(dt.columns, dt2.columns)
		for k, v := range dt.columns {
			v.Concate(dt2.columns[k])
		}

	} else {
		panic("invalid alignment")
	}
}

type Selection struct {
	data         *DataTable
	columnNames  []string
	rowIntervals [][2]int
	selection    map[string]*Slice
}

func (s *Selection) Or(selection Selection) {
	for k, v := range selection.selection {
		for i, v2 := range v.fields {
			if v2.value.(bool) {
				s.selection[k].fields[i] = &Cell{value: true}
			}
		}
	}
}

func (s *Selection) And(selection Selection) {
	for k, v := range selection.selection {
		for i, v2 := range v.fields {
			if v2.value.(bool) && s.selection[k].fields[i].value.(bool) {
				s.selection[k].fields[i] = &Cell{value: true}
			} else {
				s.selection[k].fields[i] = &Cell{value: false}
			}
		}
	}
}

func (s *Selection) LessThan(value float64) {
	s.selection = make(map[string]*Slice)
	for name, i := range s.data.columnNameIds {
		s.selection[name] = &Slice{
			size:   s.data.shape[0],
			fields: make(map[int]*Cell),
		}
		for j := 0; j < s.data.shape[0]; j++ {
			if s.data.columns[i].fields[j].value.(float64) < value {
				s.selection[name].fields[j] = &Cell{value: true}
			} else {
				s.selection[name].fields[j] = &Cell{value: false}
			}
		}
	}
}

func (s *Selection) LessThanEquals(value float64) {
	s.selection = make(map[string]*Slice)
	for name, i := range s.data.columnNameIds {
		s.selection[name] = &Slice{
			size:   s.data.shape[0],
			fields: make(map[int]*Cell),
		}
		for j := 0; j < s.data.shape[0]; j++ {
			if s.data.columns[i].fields[j].value.(float64) <= value {
				s.selection[name].fields[j] = &Cell{value: true}
			} else {
				s.selection[name].fields[j] = &Cell{value: false}
			}
		}
	}
}

func (s *Selection) GreaterThan(value float64) {
	s.selection = make(map[string]*Slice)
	for name, i := range s.data.columnNameIds {
		s.selection[name] = &Slice{
			size:   s.data.shape[0],
			fields: make(map[int]*Cell),
		}
		for j := 0; j < s.data.shape[0]; j++ {
			if s.data.columns[i].fields[j].value.(float64) > value {
				s.selection[name].fields[j] = &Cell{value: true}
			} else {
				s.selection[name].fields[j] = &Cell{value: false}
			}
		}
	}
}

func (s *Selection) GreaterThanEquals(value float64) {
	s.selection = make(map[string]*Slice)
	for name, i := range s.data.columnNameIds {
		s.selection[name] = &Slice{
			size:   s.data.shape[0],
			fields: make(map[int]*Cell),
		}
		for j := 0; j < s.data.shape[0]; j++ {
			if s.data.columns[i].fields[j].value.(float64) >= value {
				s.selection[name].fields[j] = &Cell{value: true}
			} else {
				s.selection[name].fields[j] = &Cell{value: false}
			}
		}
	}
}

func (s *Selection) EqualsTo(value interface{}) {
	s.selection = make(map[string]*Slice)
	for name, i := range s.data.columnNameIds {
		s.selection[name] = &Slice{
			size:   s.data.shape[0],
			fields: make(map[int]*Cell),
		}
		for j := 0; j < s.data.shape[0]; j++ {
			if s.data.columns[i].fields[j].value.(float64) == value {
				s.selection[name].fields[j] = &Cell{value: true}
			} else {
				s.selection[name].fields[j] = &Cell{value: false}
			}
		}
	}
}

func (dt *DataTable) Select(selector ...Selection) *DataTable {
	data := &DataTable{
		shape:         dt.shape,
		columns:       dt.columns,
		columnNameIds: dt.columnNameIds,
		selection:     dt.selection,
	}
	for _, s := range selector {
		data.selection.And(s)
	}
	return data
}

func (dt *DataTable) Set(value interface{}) {
	for str, v := range dt.columnNameIds {
		for i := 0; i < dt.columns[v].size; i++ {
			for _, s := range dt.selection {
				if s.selection[str].fields[i].value.(bool) {
					dt.columns[v].fields[i].value = value
					continue
				}
			}
		}
	}
}

func (dt *DataTable) Column(name ...string) *Selection {
	if len(name) == 0 {
		panic("invalid column name")
	}
	for _, v := range name {
		if _, ok := dt.columnNameIds[v]; !ok {
			panic("invalid column name")
		}
	}
	return &Selection{
		data:        dt,
		columnNames: name,
	}
}

func (dt *DataTable) Row(interval ...[2]int) *Selection {
	if len(interval) == 0 {
		panic("invalid interval")
	}

	for i := range interval {
		if interval[i][0] < 0 || interval[i][0] >= dt.shape[0] {
			panic("invalid interval")
		}
		if interval[i][1] < 0 || interval[i][1] >= dt.shape[0] {
			panic("invalid interval")
		}
	}
	return &Selection{
		data:         dt,
		rowIntervals: interval,
	}
}

type Region struct {
	data      *DataTable
	selection *Selection
}

type TableView interface {
	RowEntry() int
	ColEntry() int
	RowExit() int
	ColExit() int
}

type View struct {
	data *DataTable
}

func (dt *DataTable) RowView(row int) *Slice {
	if row < 0 || row >= dt.shape[0] {
		panic("row out of bounds")
	}
	slice := &Slice{
		size:   dt.shape[1],
		fields: make(map[int]*Cell),
	}
	i := 0
	for _, sl := range dt.columns {
		slice.fields[i] = sl.fields[row]
		i++
	}
	return slice
}

func (dt *DataTable) ColView(col int) *Slice {
	if col < 0 || col >= dt.shape[1] {
		panic("col out of bounds")
	}
	slice := &Slice{
		size:   dt.shape[0],
		fields: make(map[int]*Cell),
	}
	i := 0
	for _, sl := range dt.columns[col].fields {
		slice.fields[i] = sl
		i++
	}
	return slice
}