package data

import "testing"

func TestValidation(t *testing.T) {
	p := &Product{Name: "jaskeerat", Price: 12, SKU: "abs-af-yh"}
	err := p.Validate()
	if err != nil {
		t.Fatal(err)
	}
}
