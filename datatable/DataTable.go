package datatable

type Alignment = int

const (
	Rows = iota
	Cols = iota
)

type Cell struct {
	value interface{}
}

type DataTable struct {
	shape         [2]int
	columns       map[int]*Slice
	columnNameIds map[string]int
	selection     Selection
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

/*
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
}*/

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

func (dt *DataTable) IterRows() []*Slice {
	slice := make([]*Slice, dt.shape[0])
	for i := 0; i < dt.shape[0]; i++ {
		slice[i] = dt.RowView(i)
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
