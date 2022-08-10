package handlers

import (
	"io"

	"github.com/lakshmaji/delivery-shell/clients"
	"github.com/lakshmaji/delivery-shell/models"
	"github.com/lakshmaji/delivery-shell/services/shell_io_svc"
)

type mockInputServiceData struct {
	reader                          io.Reader
	baseDeliveryCost                float64
	noOfPackages                    int
	boxes                           []*models.PackageDetails
	ErrScanBaseDeliveryCostPkgCount error
	ErrScanNPackageDetails          error
	ErrScanVehicleDetails           error
	noOfVehicles                    int
	speed                           int
	maxWeight                       int
	ErrScanProgramChoice            error
	choice                          string
}

type mockDeliveryPrgmInputs struct {
	reader                          io.Reader
	baseDeliveryCost                float64
	noOfPackages                    int
	boxes                           []*models.PackageDetails
	ErrScanBaseDeliveryCostPkgCount error
	ErrScanNPackageDetails          error
	ErrScanVehicleDetails           error
	noOfVehicles                    int
	speed                           int
	maxWeight                       int
	ErrScanProgramChoice            error
	choice                          string
}

func newInputSvc(data mockInputServiceData) shell_io_svc.PackageInputService {
	return &mockDeliveryPrgmInputs{
		reader:                          data.reader,
		baseDeliveryCost:                data.baseDeliveryCost,
		noOfPackages:                    data.noOfPackages,
		boxes:                           data.boxes,
		ErrScanBaseDeliveryCostPkgCount: data.ErrScanBaseDeliveryCostPkgCount,
		ErrScanNPackageDetails:          data.ErrScanNPackageDetails,
		ErrScanVehicleDetails:           data.ErrScanVehicleDetails,
		noOfVehicles:                    data.noOfVehicles,
		speed:                           data.speed,
		maxWeight:                       data.maxWeight,
		ErrScanProgramChoice:            data.ErrScanProgramChoice,
		choice:                          data.choice,
	}
}

func (d *mockDeliveryPrgmInputs) ScanBaseDeliveryCostPkgCount(writer clients.BaseWriter) (models.BaseDeliveryCost, int, error) {

	if d.ErrScanBaseDeliveryCostPkgCount != nil {
		return 0, 0, d.ErrScanBaseDeliveryCostPkgCount
	}

	return models.BaseDeliveryCost(d.baseDeliveryCost), d.noOfPackages, nil
}

func (d *mockDeliveryPrgmInputs) ScanNPackageDetails(writer clients.BaseWriter, noOfPackages int) ([]*models.PackageDetails, error) {
	if d.ErrScanNPackageDetails != nil {
		return nil, d.ErrScanNPackageDetails
	}
	return d.boxes, nil
}

func (d *mockDeliveryPrgmInputs) ScanVehicleDetails(writer clients.BaseWriter) (int, int, int, error) {
	if d.ErrScanVehicleDetails != nil {
		return 0, 0, 0, d.ErrScanVehicleDetails
	}
	return d.noOfVehicles, d.speed, d.maxWeight, nil
}

func (d *mockDeliveryPrgmInputs) ScanProgramChoice(writer clients.BaseWriter) (string, error) {
	if d.ErrScanProgramChoice != nil {
		return "", d.ErrScanProgramChoice
	}
	return d.choice, nil
}
