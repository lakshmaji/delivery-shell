package shell_io_svc

import (
	"github.com/lakshmaji/delivery-shell/clients"
	"github.com/lakshmaji/delivery-shell/models"
)

// Handles responsibility of capturing inputs from given scanner
type PackageInputService interface {
	ScanBaseDeliveryCostPkgCount(clients.BaseWriter) (models.BaseDeliveryCost, int, error)
	ScanNPackageDetails(clients.BaseWriter, int) ([]*models.PackageDetails, error)
	ScanVehicleDetails(clients.BaseWriter) (int, int, int, error)
	ScanProgramChoice(clients.BaseWriter) (string, error)
}
