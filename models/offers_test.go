package models

import (
	"reflect"
	"testing"
)

func TestFactsToValidate(t *testing.T) {
	offer := Offer{}

	result := offer.FactsToValidate()
	expected := Facts{"distance", "weight"}

	if len(result) != 2 {
		t.Error("There should be only two facts")
	}
	if !reflect.DeepEqual(result, expected) {
		t.Error("facts are invalidated or modified, verify if this is intentional")
	}
}
