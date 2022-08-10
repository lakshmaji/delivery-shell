package common_utils

import (
	"math"
	"strconv"
)

func ConvertStrToInt(value string) (int, error) {
	val, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func ConvertStrToFloat64(value string) (float64, error) {
	val, err := (strconv.ParseFloat(value, 64))
	if err != nil {
		return 0, err
	}
	return float64(val), nil
}

func CanComputeDeliveryTime(value string) bool {
	return value == "yes"

}
func IsValidChoice(value string) bool {
	return value == "yes" || value == "no"
}

func ToFixed(value float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(int(value*output)) / output
}

func MaxVal(a, b int) int {
	if a > b {
		return a
	}
	return b
}
