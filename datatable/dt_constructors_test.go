package datatable_test

import (
	"fmt"
	"testing"

	"github.com/jaeckl/gods/datatable"
)

func TestFromJSON(t *testing.T) {
	dt := datatable.FromJSON("../resources/test.json")
	fmt.Println(dt)
}
