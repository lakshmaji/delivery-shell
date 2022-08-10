package common_utils

import (
	"strconv"
	"testing"
)

func TestConvertStrToFloat64(t *testing.T) {
	tt := []struct {
		Name     string
		Input    interface{}
		Expected interface{}
	}{
		{Name: "float64", Input: "23.9", Expected: float64(23.9)},
		{Name: "int", Input: "23", Expected: float64(23)},
	}

	for _, test := range tt {

		t.Run(t.Name(), func(t *testing.T) {
			result, err := ConvertStrToFloat64(test.Input.(string))
			if err != nil {
				t.Error("should not return error")
			}
			if result != test.Expected {
				t.Errorf("expected %d received %v", test.Expected, result)
			}
		})
	}
}

func TestConvertStrToFloat64Errors(t *testing.T) {

	tt := []struct {
		Name     string
		Input    interface{}
		Expected *strconv.NumError
	}{
		{Name: "alphanumeric", Input: "23a.9", Expected: &strconv.NumError{Func: "ParseFloat", Num: "23a.9", Err: strconv.ErrSyntax}},
		{Name: "alpha", Input: "two", Expected: &strconv.NumError{Func: "ParseFloat", Num: "two", Err: strconv.ErrSyntax}},
		{Name: "special characters", Input: "$%", Expected: &strconv.NumError{Func: "ParseFloat", Num: "$%", Err: strconv.ErrSyntax}},
		{Name: "string with special characters", Input: "abc@", Expected: &strconv.NumError{Func: "ParseFloat", Num: "abc@", Err: strconv.ErrSyntax}},
		{Name: "number with spaces", Input: " 55    ", Expected: &strconv.NumError{Func: "ParseFloat", Num: " 55    ", Err: strconv.ErrSyntax}},
		{Name: "boolean string", Input: "true", Expected: &strconv.NumError{Func: "ParseFloat", Num: "true", Err: strconv.ErrSyntax}},
	}

	for _, test := range tt {
		t.Run(t.Name(), func(t *testing.T) {
			result, err := ConvertStrToFloat64(test.Input.(string))
			if err.Error() != test.Expected.Error() {
				t.Errorf("should return error, received %v", err)
			}
			if result != 0 {
				t.Errorf("expected default value received %v", result)
			}
		})
	}
}

func TestConvertStrToInt(t *testing.T) {

	tt := []struct {
		Name     string
		Input    interface{}
		Expected interface{}
	}{
		{Name: "int", Input: "23", Expected: 23},
	}

	for _, test := range tt {
		t.Run(t.Name(), func(t *testing.T) {
			result, err := ConvertStrToInt(test.Input.(string))
			if err != nil {
				t.Error("should not return error")
			}
			if result != test.Expected {
				t.Errorf("expected %d received %v", test.Expected, result)
			}
		})
	}
}
func TestConvertStrToIntErrors(t *testing.T) {

	tt := []struct {
		Name     string
		Input    interface{}
		Expected *strconv.NumError
	}{
		{Name: "float", Input: "23.00", Expected: &strconv.NumError{Func: "Atoi", Num: "23.00", Err: strconv.ErrSyntax}},
		{Name: "alphanumeric", Input: "23a.9", Expected: &strconv.NumError{Func: "Atoi", Num: "23a.9", Err: strconv.ErrSyntax}},
		{Name: "alpha", Input: "two", Expected: &strconv.NumError{Func: "Atoi", Num: "two", Err: strconv.ErrSyntax}},
		{Name: "special characters", Input: "$%", Expected: &strconv.NumError{Func: "Atoi", Num: "$%", Err: strconv.ErrSyntax}},
		{Name: "string with special characters", Input: "abc@", Expected: &strconv.NumError{Func: "Atoi", Num: "abc@", Err: strconv.ErrSyntax}},
		{Name: "number with spaces", Input: " 55    ", Expected: &strconv.NumError{Func: "Atoi", Num: " 55    ", Err: strconv.ErrSyntax}},
		{Name: "boolean string", Input: "true", Expected: &strconv.NumError{Func: "Atoi", Num: "true", Err: strconv.ErrSyntax}},
	}

	for _, test := range tt {
		t.Run(t.Name(), func(t *testing.T) {
			result, err := ConvertStrToInt(test.Input.(string))
			if err.Error() != test.Expected.Error() {
				t.Errorf("should return error, received %v", err)
			}
			if result != 0 {
				t.Errorf("expected default value received %v", result)
			}
		})
	}

}

func TestToFixed(t *testing.T) {
	expected := 23.987
	output := ToFixed(23.98734864736, 3)
	if output != expected {
		t.Errorf("expected %f received %f", expected, output)
	}
}

func TestMaxVal(t *testing.T) {
	output := MaxVal(23, 70)
	expected := 70
	if output != expected {
		t.Errorf("expected %d received %d", expected, output)
	}
}
