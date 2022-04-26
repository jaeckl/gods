package datatable

import (
	"encoding/json"
	"os"
)

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

func FromJSON(filepath string) *DataTable {
	var data []map[string]interface{}
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil {
		//panic(err)
	}
	if len(data) != 0 {
		dt := &DataTable{
			columns:       make(map[int]*Slice),
			columnNameIds: make(map[string]int),
			shape:         [2]int{1, len(data[0])},
		}
		i := 0
		for k, v := range data[0] {
			dt.columnNameIds[k] = i
			dt.columns[i] = &Slice{
				size:   1,
				fields: make(map[int]*Cell),
			}
			dt.columns[i].fields[0] = &Cell{
				value: v,
			}
			i++
		}
		j := 1
		for _, v := range data[1:] {
			for k, val := range v {
				if c, ok := dt.columns[dt.columnNameIds[k]]; ok {
					c.fields[j] = &Cell{
						value: val,
					}
					c.size = j + 1
				} else {
					panic("invalid column name")
				}
			}
			j++
		}
		return dt
	}
	panic("No data in file")
}
