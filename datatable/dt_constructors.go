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
	var data map[string][]interface{}
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil {
		panic(err)
	}
	return FromMap(data)
}
