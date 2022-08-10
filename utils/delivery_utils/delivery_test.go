package delivery_utils

import (
	"testing"
)

type testCase struct {
	deliveryCost      float64
	discount          float64
	totalDeliveryCost float64
	description       string
}

var calculateDeliveryCostTests = []testCase{
	{deliveryCost: 10, discount: 0.01, totalDeliveryCost: 9.99, description: "when both inputs are decimals"},
	{deliveryCost: 10, discount: 3, totalDeliveryCost: 7, description: "when both inputs are integers"},
	{deliveryCost: 0, discount: 0, totalDeliveryCost: 0, description: "when both inputs are 0"},
	{deliveryCost: 1, discount: 3, totalDeliveryCost: 0, description: "when delivery cost is less than discount"},
}

func TestCalculateTotalDeliveryCost(t *testing.T) {
	for _, test := range calculateDeliveryCostTests {
		t.Run(test.description, func(t *testing.T) {
			totalDeliveryCost := TotalDeliveryCost(test.deliveryCost, test.discount)
			if totalDeliveryCost != test.totalDeliveryCost {
				t.Errorf("\t\tExpected: %v, got: %v", test.totalDeliveryCost, totalDeliveryCost)
			}
		})
	}
}
