package delivery_svc

import (
	"github.com/lakshmaji/delivery-shell/models"
	"github.com/lakshmaji/delivery-shell/services/offers_svc"
)

var OffersSliceMock []models.Offer

type offerServiceMock struct {
}

func init() {
	OffersSliceMock = []models.Offer{
		{
			Code:     models.OfferCode("A"),
			Discount: 20,
			Conditions: []models.Condition{
				{
					Fact:     "weight",
					Operator: "lessThan",
					Value:    models.Weight(20),
				},
			},
		},
	}

}

func NewOffersSvcMock() offers_svc.OffersService {
	return &offerServiceMock{}
}

func (*offerServiceMock) ApplicableDiscount(deliveryCost float64, code models.OfferCode, weight models.Weight, distance models.Distance) (float64, error) {
	return 0.05, nil
}
