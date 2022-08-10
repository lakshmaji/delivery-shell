package models

import "testing"

func TestMapPackageStatsOutput(t *testing.T) {

	tt := []struct {
		description    string
		boxes          PackageStatsList
		computeEstTime bool
		expected       string
	}{
		{
			description:    "TestMapPackageStatsOutput",
			boxes:          PackageStatsList{PackageStats{Id: "PKG 1", Discount: 10, TotalDeliveryCost: 100}, PackageStats{Id: "PKG 10", Discount: 13, TotalDeliveryCost: 70}},
			computeEstTime: false,
			expected:       "Package Id, Discount, Total Delivery Cost\nPKG 1, 10.00, 100.00\nPKG 10, 13.00, 70.00\n",
		},
		{
			description:    "TestMapPackageStatsOutput",
			boxes:          PackageStatsList{PackageStats{Id: "PKG 1", Discount: 10, TotalDeliveryCost: 100, EstDeliveryTime: 0.43}, PackageStats{Id: "PKG 10", Discount: 13, TotalDeliveryCost: 70, EstDeliveryTime: 1.78}},
			computeEstTime: true,
			expected:       "Package Id, Discount, Total Delivery Cost, Total Est Time\nPKG 1, 10.00, 100.00, 0.43\nPKG 10, 13.00, 70.00, 1.78\n",
		},
	}

	for _, tc := range tt {
		t.Run(tc.description, func(t *testing.T) {
			output := tc.boxes.FmtOutput(tc.computeEstTime)
			if output != tc.expected {
				t.Errorf("Expected %v, got %v", tc.expected, output)
			}
		})
	}
}
