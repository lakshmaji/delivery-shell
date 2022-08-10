package shell_io_svc

import (
	"bufio"
	"io"
	"strings"

	"github.com/lakshmaji/delivery-shell/clients"
	"github.com/lakshmaji/delivery-shell/models"
	"github.com/lakshmaji/delivery-shell/utils/common_utils"
	"github.com/lakshmaji/delivery-shell/utils/error_utils"
	"github.com/lakshmaji/delivery-shell/utils/msg_utils"
)

type packageInputSvc struct {
	reader io.Reader
}

// Handles responsibility of capturing inputs from **stdin**
func NewShellReader(reader io.Reader) PackageInputService {
	return &packageInputSvc{reader}
}

// Reads base delivery cost and no of packages
func (d *packageInputSvc) ScanBaseDeliveryCostPkgCount(writer clients.BaseWriter) (models.BaseDeliveryCost, int, error) {
	scanner := bufio.NewScanner(d.reader)
	writer.Write(msg_utils.MsgBaseCostPkgCountHeader)
	scanner.Scan()
	text := scanner.Text()
	if len(text) == 0 {
		return 0, 0, error_utils.ErrMissingInput
	}

	input := strings.Fields(text)
	if len(input) != 2 {
		return 0, 0, error_utils.ErrBaseCostPkgCount
	}

	cost, err := common_utils.ConvertStrToInt(input[0])
	if err != nil {
		return 0, 0, err
	}

	noOfPackages, err := common_utils.ConvertStrToInt(input[1])
	if err != nil {
		return 0, 0, err
	}

	return models.BaseDeliveryCost(cost), noOfPackages, nil
}

// Reads package details from user input
func (d *packageInputSvc) ScanNPackageDetails(writer clients.BaseWriter, noOfPackages int) ([]*models.PackageDetails, error) {
	scanner := bufio.NewScanner(d.reader)
	var packages []*models.PackageDetails
	for i := 0; i < noOfPackages; i++ {
		writer.Write(msg_utils.MsgPackageDetailsHeader)
		scanner.Scan()
		text := scanner.Text()
		if len(text) == 0 {
			return nil, error_utils.ErrMissingInput
		}

		input := strings.Fields(text)
		if len(input) != 4 {
			return nil, error_utils.ErrPackageDetailsFormat
		}
		weight, err := common_utils.ConvertStrToFloat64(input[1])
		if err != nil {
			return nil, err
		}
		distance, err := common_utils.ConvertStrToFloat64(input[2])
		if err != nil {
			return nil, err
		}
		box := models.PackageDetails{
			Id:       models.PackageID(input[0]),
			Weight:   weight,
			Distance: distance,
			Code:     models.OfferCode(input[3]),
		}

		packages = append(packages, &box)
	}

	return packages, nil
}

// no_of_vehicles <space> max_speed_of_all_vehicles_in_km_per_hour <space> max_capacity_of_all_vehicles_in_kg
// Reads base delivery cost and no of packages
func (d *packageInputSvc) ScanVehicleDetails(writer clients.BaseWriter) (int, int, int, error) {
	scanner := bufio.NewScanner(d.reader)
	writer.Write(msg_utils.MsgVehiclesHeader)
	scanner.Scan()
	text := scanner.Text()
	if len(text) == 0 {
		return 0, 0, 0, error_utils.ErrMissingInput
	}

	input := strings.Fields(text)
	if len(input) != 3 {
		return 0, 0, 0, error_utils.ErrVehicleDetailsFormat
	}

	noOfVehicles, err := common_utils.ConvertStrToInt(input[0])
	if err != nil {
		return 0, 0, 0, err
	}
	speed, err := common_utils.ConvertStrToInt(input[1])
	if err != nil {
		return 0, 0, 0, err
	}
	maxWeight, err := common_utils.ConvertStrToInt(input[2])
	if err != nil {
		return 0, 0, 0, err
	}

	return noOfVehicles, speed, maxWeight, nil
}

// Which version of program to run
// no - Discount only
// yes - Discount and Est time of delivery
func (d *packageInputSvc) ScanProgramChoice(writer clients.BaseWriter) (string, error) {
	writer.Write(msg_utils.MsgProgramChoice)

	scanner := bufio.NewScanner(d.reader)
	scanner.Scan()
	timeComputeDecisionInput := scanner.Text()

	if len(timeComputeDecisionInput) == 0 {
		return "", error_utils.ErrMissingInput
	}

	if !common_utils.IsValidChoice(timeComputeDecisionInput) {
		return "", error_utils.ErrProgramChoiceFormat
	}

	return timeComputeDecisionInput, nil
}
