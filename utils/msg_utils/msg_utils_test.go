package msg_utils

import "testing"

func TestMessages(t *testing.T) {
	if MsgPackageStatsHeader != "Package Id, Discount, Total Delivery Cost" {
		t.Error("should not be changed")
	}

	if MsgPackageStatsEstTime != "Total Est Time" {
		t.Error("should not be changed")
	}

	if MsgBaseCostPkgCountHeader != "Enter \"base delivery cost\" and \"No of packages\":" {
		t.Error("should not be changed")
	}

	if MsgPackageDetailsHeader != "Enter package id, weight, distance and offer code:" {
		t.Error("should not be changed")
	}

	if MsgVehiclesHeader != "Enter \"vehicles count\" \"speed\" \"weight capacity\":" {
		t.Error("should not be changed")
	}

	if MsgProgramChoice != "Do you want compute est time for delivery [yes, no]" {
		t.Error("should not be changed")
	}
}
