package delivery_svc

import (
	"reflect"
	"testing"

	"github.com/lakshmaji/delivery-shell/models"
)

func TestCalculateDeliveryCost(t *testing.T) {
	type args struct {
		weight           models.Weight
		distance         models.Distance
		baseDeliveryCost models.BaseDeliveryCost
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "TestCalculateDeliveryCost",
			args: args{
				weight:           models.Weight(10),
				distance:         models.Distance(5),
				baseDeliveryCost: models.BaseDeliveryCost(100),
			},
			want: float64(100) + (float64(10) * 10) + (float64(5) * 5),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewDeliveryService(NewOffersSvcMock())
			if got := svc.CalculateDeliveryCost(tt.args.weight, tt.args.distance, tt.args.baseDeliveryCost); got != tt.want {
				t.Errorf("CalculateDeliveryCost() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculateDiscount(t *testing.T) {
	type args struct {
		weight       models.Weight
		distance     models.Distance
		code         models.OfferCode
		deliveryCost float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "TestCalculateDiscount",
			args: args{
				weight:       models.Weight(10),
				distance:     models.Distance(100),
				code:         models.OfferCode("OFR003"),
				deliveryCost: 100,
			},
			want: 0.05,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewDeliveryService(NewOffersSvcMock())
			got, _ := svc.CalculateDiscount(tt.args.weight, tt.args.distance, tt.args.code, tt.args.deliveryCost)
			if got != tt.want {
				t.Errorf("CalculateDiscount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculateEstTime(t *testing.T) {
	type args struct {
		// assumption input will be always sanitized (weight)
		items        []*models.PackageDetails
		maxWeight    int
		noOfVehicles int
		maxSpeed     int
	}

	tests := []struct {
		name string
		args args
		want models.PackageDeliveryTime
	}{
		{
			name: "Sample 1",
			args: args{
				items: []*models.PackageDetails{
					{
						Id:       "PKG1",
						Weight:   3,
						Distance: 20,
						Code:     "OFR001",
					},
					{
						Id:       "PKG2",
						Weight:   3,
						Distance: 10,
						Code:     "OFR002",
					},
				},
				noOfVehicles: 1,
				maxSpeed:     70,
				maxWeight:    5,
			},
			want: models.PackageDeliveryTime{
				"PKG1": 0.56,
				"PKG2": 0.14,
			},
		},
		{
			name: "Sample 2",
			args: args{
				items: []*models.PackageDetails{
					{
						Id:       "PKG1",
						Weight:   3,
						Distance: 10,
						Code:     "OFR001",
					},
					{
						Id:       "PKG2",
						Weight:   3,
						Distance: 10,
						Code:     "OFR002",
					},
				},
				noOfVehicles: 1,
				maxSpeed:     70,
				maxWeight:    5,
			},
			want: models.PackageDeliveryTime{
				"PKG1": 0.14,
				"PKG2": 0.42,
			},
		},
		{
			name: "Sample 3 (where package wont be shipped)",
			args: args{
				items: []*models.PackageDetails{
					{
						Id:       "PKG1",
						Weight:   3,
						Distance: 10,
						Code:     "OFR001",
					},
					{
						Id:       "PKG2",
						Weight:   0,
						Distance: 5,
						Code:     "OFR002",
					},
				},
				noOfVehicles: 1,
				maxSpeed:     70,
				maxWeight:    5,
			},
			want: models.PackageDeliveryTime{
				"PKG1": 0.14,
				// "PKG2": wont be shipped (because its weight 0, which means not a package 😂),
			},
		},
		{
			name: "Sample 4 (where package wont be shipped)",
			args: args{
				items: []*models.PackageDetails{
					{
						Id:       "PKG1",
						Weight:   3,
						Distance: 10,
						Code:     "OFR001",
					},
					{
						Id:       "PKG2",
						Weight:   28,
						Distance: 5,
						Code:     "OFR002",
					},
				},
				noOfVehicles: 1,
				maxSpeed:     70,
				maxWeight:    5,
			},
			want: models.PackageDeliveryTime{
				"PKG1": 0.14,
				// "PKG2": wont be shipped (because its weight 28, beyond vehicle cap),
			},
		},
		{
			name: "Sample 5",
			args: args{
				items: []*models.PackageDetails{
					{
						Id:       "PKG1",
						Weight:   50,
						Distance: 30,
						Code:     "OFR001",
					},
					{
						Id:       "PKG2",
						Weight:   75,
						Distance: 125,
						Code:     "OFR002",
					},
					{
						Id:       "PKG3",
						Weight:   175,
						Distance: 100,
						Code:     "OFR003",
					},
					{
						Id:       "PKG4",
						Weight:   110,
						Distance: 60,
						Code:     "OFR002",
					},
					{
						Id:       "PKG5",
						Weight:   155,
						Distance: 95,
						Code:     "NA",
					},
				},
				noOfVehicles: 2,
				maxSpeed:     70,
				maxWeight:    200,
			},
			want: models.PackageDeliveryTime{
				"PKG1": 3.98,
				"PKG2": 1.78,
				"PKG3": 1.42, "PKG4": 0.85, "PKG5": 4.19,
			},
		},
		{
			name: "Sample 6 (when distance is zero)",
			args: args{
				items: []*models.PackageDetails{
					{
						Id:       "PKG1",
						Weight:   3,
						Distance: 0,
						Code:     "OFR001",
					},
					{
						Id:       "PKG2",
						Weight:   2,
						Distance: 0,
						Code:     "OFR002",
					},
				},
				noOfVehicles: 1,
				maxSpeed:     70,
				maxWeight:    5,
			},
			want: models.PackageDeliveryTime{
				"PKG1": 0,
				"PKG2": 0,
			},
		},
		{
			name: "Sample 7 (max no of items reader shipping)",
			args: args{
				items: []*models.PackageDetails{

					{
						Id:       "PKG2",
						Weight:   4,
						Distance: 10,
						Code:     "OFR002",
					},
					{
						Id:       "PKG3",
						Weight:   1,
						Distance: 10,
						Code:     "OFR002",
					},
					{
						Id:       "PKG1",
						Weight:   2,
						Distance: 10,
						Code:     "OFR001",
					},
					{
						Id:       "PKG4",
						Weight:   2,
						Distance: 5,
						Code:     "OFR002",
					},
					{
						Id:       "PKG5",
						Weight:   3,
						Distance: 10,
						Code:     "OFR002",
					},
				},
				noOfVehicles: 1,
				maxSpeed:     5,
				maxWeight:    6,
			},
			want: models.PackageDeliveryTime{
				"PKG1": 6,
				"PKG2": 6,
				"PKG3": 2,
				"PKG4": 1,
				"PKG5": 2,
			},
		},
		{
			name: "Sample 8 (when shipments weights are same and pick less distant one first)",
			args: args{
				items: []*models.PackageDetails{
					{
						Id:       "PKG1",
						Weight:   3,
						Distance: 10,
						Code:     "OFR001",
					},
					{
						Id:       "PKG2",
						Weight:   3,
						Distance: 10,
						Code:     "OFR002",
					},
					{
						Id:       "PKG3",
						Weight:   2,
						Distance: 10,
						Code:     "OFR002",
					},
					{
						Id:       "PKG4",
						Weight:   4,
						Distance: 5,
						Code:     "OFR002",
					},
				},
				noOfVehicles: 1,
				maxSpeed:     5,
				maxWeight:    6,
			},
			want: models.PackageDeliveryTime{
				"PKG1": 2,
				"PKG2": 2,
				"PKG3": 6,
				"PKG4": 5,
			},
		},
		{
			name: "Sample 9 (when shipments weights are same and pick less distant one first)",
			args: args{
				items: []*models.PackageDetails{
					{
						Id:       "PKG1",
						Weight:   3,
						Distance: 10,
						Code:     "OFR001",
					},
					{
						Id:       "PKG2",
						Weight:   3,
						Distance: 5,
						Code:     "OFR002",
					},
					{
						Id:       "PKG3",
						Weight:   3,
						Distance: 10,
						Code:     "OFR002",
					},
					{
						Id:       "PKG4",
						Weight:   3,
						Distance: 5,
						Code:     "OFR002",
					},
				},
				noOfVehicles: 1,
				maxSpeed:     5,
				maxWeight:    6,
			},
			want: models.PackageDeliveryTime{
				"PKG1": 4,
				"PKG2": 1,
				"PKG3": 4,
				"PKG4": 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewDeliveryService(NewOffersSvcMock())
			got := svc.EstDeliveryTime(tt.args.items, tt.args.maxWeight, tt.args.noOfVehicles, tt.args.maxSpeed)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EstDeliveryTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
