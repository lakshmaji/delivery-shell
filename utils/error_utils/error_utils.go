package error_utils

import (
	"errors"
	"fmt"

	"github.com/lakshmaji/delivery-shell/models"
)

var (
	ErrMissingInput          = errors.New("Missing input")
	ErrBaseCostPkgCount      = errors.New("Format Error:  \"base delivery cost\" and \"No of packages\" separated by space delimiter")
	ErrPackageDetailsFormat  = errors.New("Format Error: \"box_id\" \"box_weight_in_kg\" \"distance_in_km\" \"offer_code\"")
	ErrVehicleDetailsFormat  = errors.New("Format Error: \"vehicles count\" \"speed\" \"weight capacity\"")
	ErrProgramChoiceFormat   = errors.New("Format Error: enter one of them yes, no")
	ErrPackageDetailsInValid = errors.New("Package weight wont be considered for delivery")
	ErrCalculateDiscount     = errors.New("Error while applying discount")
)

func ErrVehicleMaxWeightCapacity(box *models.PackageDetails, maxWeight int) error {
	//nolint:gosimple
	return errors.New(fmt.Sprintf("Box %s weight %f exceed vehicle max weight capacity of %d", box.Id, box.Weight, maxWeight))
}
