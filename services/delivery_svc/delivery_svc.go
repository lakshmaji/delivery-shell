package delivery_svc

import "github.com/lakshmaji/delivery-shell/models"

type DeliveryService interface {

	//  Calculates delivery cost for a package
	//
	//  Formulae:
	//   delivery cost = (package total weight * 10) + (distance to destination *5)
	//
	//  We could any other factors impacting delivery service charge like weather, surge etc. (without any offer service related code)
	//
	//  @param weight Weight of package
	//  @param distance Distance to destination
	//  @param baseDeliveryCost	Base Delivery Cost
	//
	//  @return delivery cost
	CalculateDeliveryCost(weight models.Weight, distance models.Distance, baseDeliveryCost models.BaseDeliveryCost) float64
	//  Calculates discount for a package based on applicable Offer code.
	//  Validate Offer code using offers service. (So that the other consuming services don't require to do this)
	//
	//
	//  ENHANCEMENT: Do other validations like whether pkg is valid, or pkg it self is eligible to accept discounts etc...!
	//  any additional discount logic apart from OfferCode
	//
	//
	//  @param weight Weight of package
	//  @param distance Distance to destination
	//  @param code Offer code
	//  @param deliveryCost Delivery cost
	//
	//  @return discount
	CalculateDiscount(weight models.Weight, distance models.Distance, code models.OfferCode, deliveryCost float64) (float64, error)

	EstDeliveryTime(items []*models.PackageDetails, maxWeight int, noOfVehicles int, maxSpeed int) models.PackageDeliveryTime
}
