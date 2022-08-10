package offers_svc

import (
	"errors"
	"testing"

	"github.com/lakshmaji/delivery-shell/models"
)

func TestApplicableDiscount(t *testing.T) {
	mockIoReadFile := func(filename string) ([]models.Offer, error) {
		var offersSlice []models.Offer = []models.Offer{
			{
				Code:     "A",
				Discount: 0.20,
				Conditions: []models.Condition{
					{
						Fact:     "weight",
						Operator: "lessThan",
						Value:    models.Weight(20),
					},
					{
						Fact:     "distance",
						Operator: "lessThan",
						Value:    models.Distance(10),
					},
				},
			},
		}

		return offersSlice, nil

	}
	svc := NewOffersService(mockIoReadFile)

	type args struct {
		deliveryCost float64
		code         models.OfferCode
		weight       models.Weight
		distance     models.Distance
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "TestApplicableDiscount",
			args: args{
				deliveryCost: float64(100),
				code:         models.OfferCode("A"),
				weight:       models.Weight(10),
				distance:     models.Distance(5),
			},
			want: float64(20),
		},
		{
			name: "TestApplicableDiscount",
			args: args{
				deliveryCost: float64(100),
				code:         models.OfferCode("INVALID"),
				weight:       models.Weight(10),
				distance:     models.Distance(5),
			},
			want: float64(0),
		},
		{
			name: "TestApplicableDiscount",
			args: args{
				deliveryCost: float64(100),
				weight:       models.Weight(10),
				distance:     models.Distance(5),
			},
			want: float64(0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := svc.ApplicableDiscount(tt.args.deliveryCost, tt.args.code, tt.args.weight, tt.args.distance)
			if err != nil {
				t.Error("should not throw error")
			}
			if got != tt.want {
				t.Errorf("ApplicableDiscount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApplicableDiscountFail(t *testing.T) {

	mockIoReadFile := func(filename string) ([]models.Offer, error) {
		return nil, errors.New("unable to read contents")

	}

	svc := NewOffersService(mockIoReadFile)

	type args struct {
		deliveryCost float64
		code         models.OfferCode
		weight       models.Weight
		distance     models.Distance
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "TestApplicableDiscount",
			args: args{
				deliveryCost: float64(100),
				code:         models.OfferCode("A"),
				weight:       models.Weight(10),
				distance:     models.Distance(5),
			},
			want: float64(0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			totalCost, err := svc.ApplicableDiscount(tt.args.deliveryCost, tt.args.code, tt.args.weight, tt.args.distance)
			if err == nil {
				t.Error("Should throw error")
			}
			if totalCost != tt.want {
				t.Errorf("ApplicableDiscount() = %f, want %f", totalCost, tt.want)
			}
		})
	}
}
