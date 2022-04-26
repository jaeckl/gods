package main

import (
	"github.com/jaeckl/gods/datatable"
)

func main() {
	mp := map[string][]interface{}{
		"targets": []interface{}{"a", "a", "b", "b"},
		"x1":      []interface{}{1, 0.9, 3, 4},
		"x2":      []interface{}{1.2, 0.8, 2.5, 3.4},
	}
	dt := datatable.FromMap(mp)
	cl := dt.Ge
	dt.Select(dt.Column("x1").LessThan(0.9))
}
