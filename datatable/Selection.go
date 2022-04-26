package datatable

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
