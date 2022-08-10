package handlers

import (
	"github.com/lakshmaji/delivery-shell/clients"
	"github.com/lakshmaji/delivery-shell/models"
	"github.com/lakshmaji/delivery-shell/services/delivery_svc"
	"github.com/lakshmaji/delivery-shell/services/shell_io_svc"
	"github.com/lakshmaji/delivery-shell/utils/common_utils"
	"github.com/lakshmaji/delivery-shell/utils/delivery_utils"
	"github.com/lakshmaji/delivery-shell/utils/error_utils"
)

func PackageHandler(writer clients.BaseWriter, boxService delivery_svc.DeliveryService, packageInputSvc shell_io_svc.PackageInputService) {
	computesDeliveryTime, baseDeliveryCost, packages, noOfVehicles, maxSpeed, maxWeightCapacity := readInputs(writer, packageInputSvc)

	for _, box := range packages {
		if !box.IsValid() {
			writer.WriteError(error_utils.ErrPackageDetailsInValid)
		}
		if computesDeliveryTime {
			if box.Weight > float64(maxWeightCapacity) {
				writer.WriteError(error_utils.ErrVehicleMaxWeightCapacity(box, maxWeightCapacity))
			}
		}
	}

	packageStats, err := handlePackageStats(boxService, packages, baseDeliveryCost, noOfVehicles, maxSpeed, maxWeightCapacity, computesDeliveryTime)
	if err != nil {
		writer.WriteError(err)
	}
	writer.Write(packageStats.FmtOutput(computesDeliveryTime))
}

func readInputs(writer clients.BaseWriter, packageInputSvc shell_io_svc.PackageInputService) (computesDeliveryTime bool, baseDeliveryCost models.BaseDeliveryCost, packages []*models.PackageDetails, noOfVehicles, maxSpeed, maxWeightCapacity int) {
	var err error
	var noOfPackages int
	var timeComputeDecisionInput string
	timeComputeDecisionInput, err = packageInputSvc.ScanProgramChoice(writer)
	computesDeliveryTime = common_utils.CanComputeDeliveryTime(timeComputeDecisionInput)
	if err != nil {
		writer.WriteError(err)
	}
	baseDeliveryCost, noOfPackages, err = packageInputSvc.ScanBaseDeliveryCostPkgCount(writer)
	if err != nil {
		writer.WriteError(err)
	}
	packages, err = packageInputSvc.ScanNPackageDetails(writer, noOfPackages)
	if err != nil {
		writer.WriteError(err)
	}
	if computesDeliveryTime {
		noOfVehicles, maxSpeed, maxWeightCapacity, err = packageInputSvc.ScanVehicleDetails(writer)
		if err != nil {
			writer.WriteError(err)
		}
	}
	return
}

// Computes discounts, est delivery time
func handlePackageStats(boxService delivery_svc.DeliveryService, boxes []*models.PackageDetails, baseDeliveryCost models.BaseDeliveryCost, noOfVehicles int, maxSpeed int, maxWeightCapacity int, computesDeliveryTime bool) (models.PackageStatsList, error) {
	var packageStats []models.PackageStats

	// clone pointer variable boxes without modifying the original
	boxesClone := make(models.Shipment, len(boxes))
	copy(boxesClone, boxes)

	var itemsDeliveryTime models.PackageDeliveryTime
	if computesDeliveryTime {
		// calculate est time
		itemsDeliveryTime = boxService.EstDeliveryTime(boxesClone, int(maxWeightCapacity), noOfVehicles, maxSpeed)
	}

	for _, pkg := range boxes {
		weight := pkg.Weight
		distance := pkg.Distance
		code := pkg.Code
		// TODO: these 3 methods can be refactored to a single method
		// get delivery cost
		deliveryCost := boxService.CalculateDeliveryCost(weight, distance, baseDeliveryCost)
		// Apply offer code if applicable
		discount, err := boxService.CalculateDiscount(weight, distance, code, deliveryCost)
		if err != nil {
			return nil, error_utils.ErrCalculateDiscount
		}
		totalDeliveryCost := delivery_utils.TotalDeliveryCost(deliveryCost, discount)
		packageStat := models.PackageStats{Id: pkg.Id, Discount: discount, TotalDeliveryCost: totalDeliveryCost}
		if computesDeliveryTime {
			packageStat.EstDeliveryTime = itemsDeliveryTime[pkg.Id]
		}
		packageStats = append(packageStats, packageStat)
	}
	return packageStats, nil
}
