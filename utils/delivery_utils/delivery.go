package delivery_utils

// Pre-conditions (for the application)
// The discount will be always applied on delivery cost, so assuming the delivery cost greater than discount.
// But this function will handle other scenarios as well.
// This could be handled by delivery service itself.
func TotalDeliveryCost(deliveryCost float64, discount float64) float64 {
	if deliveryCost == 0 {
		return deliveryCost
	}
	if deliveryCost < discount {
		return 0
	}
	return deliveryCost - discount
}
