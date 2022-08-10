package offers_svc

import (
	"github.com/lakshmaji/delivery-shell/models"
	"github.com/lakshmaji/delivery-shell/utils/offer_utils"
)

type offerService struct {
	fn func(filename string) ([]models.Offer, error)
}

// Offers service which works with a local json file
func NewOffersService(fn func(string) ([]models.Offer, error)) OffersService {
	return &offerService{
		fn: fn,
	}
}

// Retrieves the offer object for a given offer-code
func (o *offerService) retrieveOfferBy(code models.OfferCode) (models.Offer, error) {
	OffersSlice, err := o.fn("offers.json")
	if err != nil {
		return models.Offer{}, err
	}
	// convert OffersSlice to map with offer.Code as key
	offersMap := make(map[models.OfferCode]models.Offer)
	for _, offer := range OffersSlice {
		offersMap[offer.Code] = offer
	}

	return offersMap[code], nil
}

func (o *offerService) ApplicableDiscount(deliveryCost float64, code models.OfferCode, wt models.Weight, dt models.Distance) (float64, error) {
	// faking our local database call
	offer, err := o.retrieveOfferBy(code)
	if err != nil {
		return 0, err
	}
	fact := models.Fact{
		Weight:   wt,
		Distance: dt,
	}
	var canApplyDiscount bool = offer_utils.IsOfferApplicable(offer.Conditions, offer.FactsToValidate(), fact)

	if canApplyDiscount {
		return deliveryCost * offer.Discount, nil
	}
	return 0, nil
}
