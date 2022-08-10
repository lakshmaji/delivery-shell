package delivery_svc

import (
	"math"
	"sort"

	"github.com/lakshmaji/delivery-shell/models"
	"github.com/lakshmaji/delivery-shell/services/offers_svc"
	"github.com/lakshmaji/delivery-shell/utils/common_utils"
)

type defaultService struct {
	offer_svc offers_svc.OffersService
}

func NewDeliveryService(offer_svc offers_svc.OffersService) DeliveryService {
	return &defaultService{
		offer_svc: offer_svc,
	}
}

func (p *defaultService) CalculateDeliveryCost(weight models.Weight, distance models.Distance, baseDeliveryCost models.BaseDeliveryCost) float64 {
	return float64(baseDeliveryCost) + (weight * 10) + (distance * 5)
}

func (p *defaultService) CalculateDiscount(weight models.Weight, distance models.Distance, code models.OfferCode, deliveryCost float64) (float64, error) {
	return p.offer_svc.ApplicableDiscount(deliveryCost, code, weight, distance)
}

func (p *defaultService) EstDeliveryTime(items []*models.PackageDetails, maxWeight int, noOfVehicles int, maxSpeed int) models.PackageDeliveryTime {
	vehicles := initVehicles(noOfVehicles)
	var itemsDeliveryTime models.PackageDeliveryTime = make(models.PackageDeliveryTime)
	var minVehicle *models.Vehicle

	for len(items) > 0 {
		buffer := pickItemByMaxNetWeight(items, maxWeight)
		shipmentItems := getShipmentItems(items, buffer, maxWeight)

		if len(shipmentItems) == 0 {
			break
		}

		maxDeliveryTime := float64(math.MinInt64)
		minWaitTime := float64(math.MaxFloat64)

		for _, vehicle := range vehicles {
			if vehicle.WaitTime < minWaitTime {
				minWaitTime = vehicle.WaitTime
				minVehicle = vehicle
			}
		}

		var maxRoundTripTime float64
		for _, item := range shipmentItems {
			deliveredIn := float64(item.Distance) / float64(maxSpeed)
			if maxDeliveryTime < deliveredIn {
				maxDeliveryTime = common_utils.ToFixed(deliveredIn, 2)
			}

			DeliveredIn := common_utils.ToFixed(deliveredIn+minVehicle.WaitTime, 2)
			maxRoundTripTime = minVehicle.WaitTime + common_utils.ToFixed(maxDeliveryTime*2, 2)

			itemsDeliveryTime[item.Id] = DeliveredIn
		}
		minVehicle.WaitTime = maxRoundTripTime

		items = removeItems(items, shipmentItems)

	}

	return itemsDeliveryTime

}

func pickItemByMaxNetWeight(items []*models.PackageDetails, maxWeight int) [][]weightBuffer {
	buffer := make([][]weightBuffer, len(items)+1)
	for i := 0; i < len(buffer); i++ {
		buffer[i] = make([]weightBuffer, maxWeight+1)
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].Weight == items[j].Weight {
			return items[i].Distance < items[j].Distance
		}
		return items[i].Weight < items[j].Weight
	})

	for i := 1; i <= len(items); i++ {
		for j := 1; j <= maxWeight; j++ {
			if int(items[i-1].Weight) > j {
				buffer[i][j] = buffer[i-1][j]
			} else {
				prevItem := buffer[i-1][j-int(items[i-1].Weight)]
				filledWeight := int(prevItem.computedWeight) + int(items[i-1].Weight)

				buffer[i][j] = weightBuffer{
					computedWeight: common_utils.MaxVal(filledWeight, buffer[i-1][j].computedWeight),
				}
			}
		}
	}

	return buffer
}

func removeItems(items []*models.PackageDetails, shippedItems []*models.PackageDetails) []*models.PackageDetails {
	for _, shippedItem := range shippedItems {
		for i, item := range items {
			if item.IsSamePackage(*shippedItem) {
				items = append(items[:i], items[i+1:]...)
				break
			}
		}
	}
	return items
}

func getShipmentItems(items []*models.PackageDetails, buffer [][]weightBuffer, maxWeight int) []*models.PackageDetails {
	var bag []*models.PackageDetails
	i := len(items)
	j := maxWeight

	for i > 0 && j > 0 {
		if buffer[i][j].computedWeight == buffer[i-1][j].computedWeight {
			i--
		} else {
			bag = append(bag, items[i-1])
			j -= int(items[i-1].Weight)
			i--
		}
	}

	return bag
}

type weightBuffer struct {
	computedWeight int
}

func initVehicles(noOfVehicles int) []*models.Vehicle {
	vehicles := make([]*models.Vehicle, noOfVehicles)
	for i := 0; i < noOfVehicles; i++ {
		vehicles[i] = &models.Vehicle{WaitTime: 0}
	}
	return vehicles
}
