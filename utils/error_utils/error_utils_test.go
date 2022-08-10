package error_utils

import (
	"fmt"
	"testing"

	"github.com/lakshmaji/delivery-shell/models"
)

func TestErrors(t *testing.T) {
	if ErrMissingInput.Error() != "Missing input" {
		t.Error("Value changed")
	}

	if ErrBaseCostPkgCount.Error() != "Format Error:  \"base delivery cost\" and \"No of packages\" separated by space delimiter" {
		t.Error("Value changed")
	}

	if ErrPackageDetailsFormat.Error() != "Format Error: \"box_id\" \"box_weight_in_kg\" \"distance_in_km\" \"offer_code\"" {
		t.Error("Value changed")
	}

	if ErrVehicleDetailsFormat.Error() != "Format Error: \"vehicles count\" \"speed\" \"weight capacity\"" {
		t.Error("Value changed")
	}

	if ErrPackageDetailsInValid.Error() != "Package weight wont be considered for delivery" {
		t.Error("Value changed")
	}

	box := &models.PackageDetails{
		Id:     "PKG 1",
		Weight: 26,
	}
	maxWeight := 20
	if ErrVehicleMaxWeightCapacity(box, maxWeight).Error() != fmt.Sprintf("Box %s weight %f exceed vehicle max weight capacity of %d", box.Id, box.Weight, maxWeight) {
		t.Error("Value changed")
	}

	if ErrCalculateDiscount.Error() != "Error while applying discount" {
		t.Error("Value changed")
	}

}
