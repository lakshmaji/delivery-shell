package offers_svc

import (
	"github.com/lakshmaji/delivery-shell/models"
)

type OffersService interface {
	// Validate whether discount is applicable or not
	// returns computed discount (applicable)
	ApplicableDiscount(deliveryCost float64, code models.OfferCode, wt models.Weight, dt models.Distance) (float64, error)
}
