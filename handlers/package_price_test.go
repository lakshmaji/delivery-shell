package handlers

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/lakshmaji/delivery-shell/clients"
	"github.com/lakshmaji/delivery-shell/models"
	"github.com/lakshmaji/delivery-shell/services/delivery_svc"
	"github.com/lakshmaji/delivery-shell/services/offers_svc"
	"github.com/lakshmaji/delivery-shell/utils/error_utils"
)

var offersSlice []models.Offer = []models.Offer{
	{
		Code:     "OFR001",
		Discount: 0.01,
		Conditions: []models.Condition{
			{
				Fact:     "distance",
				Operator: "lessThan",
				Value:    200,
			},
			{
				Fact:     "weight",
				Operator: "greaterThanOrEqual",
				Value:    70,
			},
			{
				Fact:     "weight",
				Operator: "lessThanOrEqual",
				Value:    200,
			},
		},
	},
	{
		Code:     "OFR002",
		Discount: 0.07,
		Conditions: []models.Condition{
			{
				Fact:     "distance",
				Operator: "greaterThanOrEqual",
				Value:    50,
			},
			{
				Fact:     "distance",
				Operator: "lessThanOrEqual",
				Value:    150,
			},
			{
				Fact:     "weight",
				Operator: "greaterThanOrEqual",
				Value:    100,
			},
			{
				Fact:     "weight",
				Operator: "lessThanOrEqual",
				Value:    250,
			},
		},
	},
	{
		Code:     "OFR003",
		Discount: 0.05,
		Conditions: []models.Condition{
			{
				Fact:     "distance",
				Operator: "greaterThanOrEqual",
				Value:    50,
			},
			{
				Fact:     "distance",
				Operator: "lessThanOrEqual",
				Value:    250,
			},
			{
				Fact:     "weight",
				Operator: "greaterThanOrEqual",
				Value:    10,
			},
			{
				Fact:     "weight",
				Operator: "lessThanOrEqual",
				Value:    150,
			},
		},
	},
}

func TestPkgDiscount(t *testing.T) {
	reader, output, mockWriter, mockPkgDeliveryComputeService := mockIO(t)
	defer reader.Close()

	_, err := io.WriteString(reader, "yes\n"+"100 3\n"+"PKG1 5 5 OFR001\n"+"PKG2 15 5 OFR002\n"+"PKG3 10 100 OFR003")
	if err != nil {
		t.Fatal(err)
	}

	_, err = reader.Seek(0, io.SeekCurrent)
	if err != nil {
		t.Fatal(err)
	}
	packages := []*models.PackageDetails{
		{
			Id:       "PKG1",
			Weight:   5,
			Distance: 5,
			Code:     "OFR001",
		},
		{
			Id:       "PKG2",
			Weight:   15,
			Distance: 5,
			Code:     "OFR002",
		},
		{
			Id:       "PKG3",
			Weight:   10,
			Distance: 100,
			Code:     "OFR003",
		},
	}

	mockInput := mockInputServiceData{
		reader:           &bytes.Buffer{},
		baseDeliveryCost: 100,
		noOfPackages:     3,
		boxes:            packages,
	}
	inputSvc := newInputSvc(mockInput)

	PackageHandler(mockWriter, mockPkgDeliveryComputeService, inputSvc)

	expected := "Package Id, Discount, Total Delivery Cost\n"
	expected += "PKG1, 0.00, 175.00\n"
	expected += "PKG2, 0.00, 275.00\n"
	expected += "PKG3, 35.00, 665.00\n\n"
	e := output.String()
	if e != expected {
		t.Errorf("Expected %v, got %v", expected, output.String())
	}

}

func mockIO(t testing.TB) (*os.File, *bytes.Buffer, clients.BaseWriter, delivery_svc.DeliveryService) {
	t.Helper()
	var output bytes.Buffer

	mockOffersSvc := offers_svc.NewOffersService(
		func(filename string) ([]models.Offer, error) {

			return offersSlice, nil

		},
	)
	reader, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}

	mockPkgDeliveryComputeService := delivery_svc.NewDeliveryService(mockOffersSvc)
	mockWriter := clients.NewShellWriter(&output, true)
	return reader, &output, mockWriter, mockPkgDeliveryComputeService
}
func TestPkgDiscountFail(t *testing.T) {
	reader, output, mockWriter, mockPkgDeliveryComputeService := mockIO(t)
	defer reader.Close()
	tt := []struct {
		description                     string
		baseDeliveryCost                float64
		noOfBoxes                       int
		packages                        []*models.PackageDetails
		expected                        error
		ErrScanBaseDeliveryCostPkgCount error
		ErrScanNPackageDetails          error
		ErrScanProgramChoice            error
		choice                          string
	}{
		{
			choice:                          "no",
			description:                     "ScanBaseDeliveryCostPkgCount() fails",
			noOfBoxes:                       3,
			ErrScanBaseDeliveryCostPkgCount: error_utils.ErrMissingInput,
			expected:                        error_utils.ErrMissingInput,
		},
		{
			choice:                 "no",
			description:            "ScanNPackageDetails() fails",
			baseDeliveryCost:       100,
			noOfBoxes:              3,
			ErrScanNPackageDetails: error_utils.ErrPackageDetailsFormat,
			expected:               error_utils.ErrPackageDetailsFormat,
		},
	}

	for _, test := range tt {
		t.Run(test.description, func(t *testing.T) {
			output.Truncate(0)
			mockInput := mockInputServiceData{
				reader:           &bytes.Buffer{},
				baseDeliveryCost: test.baseDeliveryCost,
				noOfPackages:     test.noOfBoxes,
				boxes:            test.packages,

				ErrScanBaseDeliveryCostPkgCount: test.ErrScanBaseDeliveryCostPkgCount,
				ErrScanNPackageDetails:          test.ErrScanNPackageDetails,
			}
			inputSvc := newInputSvc(mockInput)

			defer func() {
				r := recover()
				expectedOutput := test.expected.Error()
				if output.String() != expectedOutput {
					t.Errorf("Expected %v, received %v", expectedOutput, output.String())
				}
				if r == nil {
					t.Errorf("Should panic")
				}
			}()

			PackageHandler(mockWriter, mockPkgDeliveryComputeService, inputSvc)

		})
	}

}

func TestDeliveryTimeEstimate(t *testing.T) {

	reader, output, mockWriter, mockPkgDeliveryComputeService := mockIO(t)
	defer reader.Close()

	tt := []struct {
		description                     string
		baseDeliveryCost                float64
		noOfBoxes                       int
		packages                        []*models.PackageDetails
		expected                        string
		ErrScanBaseDeliveryCostPkgCount error
		ErrScanNPackageDetails          error
		ErrScanVehicleDetails           error
		noOfVehicles                    int
		speed                           int
		maxWeight                       int
		ErrScanProgramChoice            error
		choice                          string
	}{
		{
			choice:           "yes",
			description:      "Sample 1",
			baseDeliveryCost: 100,
			noOfBoxes:        7,
			packages: []*models.PackageDetails{
				{
					Id:       "PKG1",
					Weight:   3,
					Distance: 30,
					Code:     "OFR001",
				},
				{
					Id:       "PKG2",
					Weight:   2,
					Distance: 125,
					Code:     "OFR002",
				},
				{
					Id:       "PKG3",
					Weight:   3,
					Distance: 100,
					Code:     "OFR008",
				},
				{
					Id:       "PKG4",
					Weight:   4,
					Distance: 60,
					Code:     "OFR002",
				},
				{
					Id:       "PKG5",
					Weight:   1,
					Distance: 95,
					Code:     "NA",
				},
				{
					Id:       "PKG6",
					Weight:   5,
					Distance: 95,
					Code:     "NA",
				},
				{
					Id:       "PKG7",
					Weight:   6,
					Distance: 95,
					Code:     "NA",
				},
			},
			noOfVehicles: 2,
			speed:        70,
			maxWeight:    6,
			expected:     "Package Id, Discount, Total Delivery Cost, Total Est Time\nPKG1, 0.00, 280.00, 0.42\nPKG2, 0.00, 745.00, 1.78\nPKG3, 0.00, 630.00, 6.68\nPKG4, 0.00, 440.00, 4.41\nPKG5, 0.00, 585.00, 1.35\nPKG6, 0.00, 625.00, 4.05\nPKG7, 0.00, 635.00, 1.35\n\n",
		},
		{
			choice:           "yes",
			description:      "Sample 2",
			baseDeliveryCost: 100,
			noOfBoxes:        5,
			packages: []*models.PackageDetails{
				{
					Id:       "PKG1",
					Weight:   50,
					Distance: 30,
					Code:     "OFR001",
				},
				{
					Id:       "PKG2",
					Weight:   75,
					Distance: 125,
					Code:     "OFR002",
				},
				{
					Id:       "PKG3",
					Weight:   175,
					Distance: 100,
					Code:     "OFR003",
				},
				{
					Id:       "PKG4",
					Weight:   110,
					Distance: 60,
					Code:     "OFR002",
				},
				{
					Id:       "PKG5",
					Weight:   155,
					Distance: 95,
					Code:     "NA",
				},
			},
			noOfVehicles: 2,
			speed:        70,
			maxWeight:    200,
			expected:     "Package Id, Discount, Total Delivery Cost, Total Est Time\nPKG1, 0.00, 750.00, 3.98\nPKG2, 0.00, 1475.00, 1.78\nPKG3, 0.00, 2350.00, 1.42\nPKG4, 105.00, 1395.00, 0.85\nPKG5, 0.00, 2125.00, 4.19\n\n",
		},
	}

	for _, test := range tt {
		t.Run(test.description, func(t *testing.T) {
			output.Truncate(0)
			mockInput := mockInputServiceData{
				reader:               &bytes.Buffer{},
				baseDeliveryCost:     test.baseDeliveryCost,
				noOfPackages:         test.noOfBoxes,
				boxes:                test.packages,
				noOfVehicles:         test.noOfVehicles,
				speed:                test.speed,
				maxWeight:            test.maxWeight,
				choice:               test.choice,
				ErrScanProgramChoice: test.ErrScanProgramChoice,
			}
			inputSvc := newInputSvc(mockInput)

			PackageHandler(mockWriter, mockPkgDeliveryComputeService, inputSvc)

			e := output.String()
			if e != test.expected {
				t.Errorf("Expected %v, received %v", test.expected, output.String())
			}
		})
	}

}

func TestDeliveryTimeEstimateFail(t *testing.T) {
	reader, output, mockWriter, mockPkgDeliveryComputeService := mockIO(t)
	defer reader.Close()

	tt := []struct {
		description                     string
		baseDeliveryCost                float64
		noOfBoxes                       int
		packages                        []*models.PackageDetails
		expected                        error
		ErrScanBaseDeliveryCostPkgCount error
		ErrScanNPackageDetails          error
		ErrScanVehicleDetails           error
		noOfVehicles                    int
		speed                           int
		maxWeight                       int
		ErrScanProgramChoice            error
		choice                          string
	}{
		{
			choice:               "test",
			description:          "Sample 1",
			ErrScanProgramChoice: error_utils.ErrProgramChoiceFormat,
			expected:             error_utils.ErrProgramChoiceFormat,
		},
		{
			choice:                "yes",
			description:           "ScanNPackageDetails() fails",
			baseDeliveryCost:      100,
			noOfBoxes:             3,
			ErrScanVehicleDetails: error_utils.ErrVehicleDetailsFormat,
			expected:              error_utils.ErrVehicleDetailsFormat,
		},
		{
			choice:           "yes",
			description:      "Package weight exceeds vehicle max weight capacity",
			baseDeliveryCost: 100,
			noOfBoxes:        3,
			packages: []*models.PackageDetails{
				{
					Id:       "PKG1",
					Weight:   5,
					Distance: 5,
					Code:     "OFR001",
				},
				{
					Id:       "PKG2",
					Weight:   15,
					Distance: 5,
					Code:     "OFR002",
				},
				{
					Id:       "PKG3",
					Weight:   10,
					Distance: 100,
					Code:     "OFR003",
				},
			},
			noOfVehicles: 2,
			speed:        70,
			maxWeight:    5,
			expected:     errors.New("Box PKG2 weight 15.000000 exceed vehicle max weight capacity of 5"),
		},
	}

	for _, test := range tt {
		t.Run(test.description, func(t *testing.T) {
			output.Truncate(0)
			mockInput := mockInputServiceData{
				reader:                          &bytes.Buffer{},
				baseDeliveryCost:                test.baseDeliveryCost,
				noOfPackages:                    test.noOfBoxes,
				boxes:                           test.packages,
				ErrScanBaseDeliveryCostPkgCount: test.ErrScanBaseDeliveryCostPkgCount,
				ErrScanNPackageDetails:          test.ErrScanNPackageDetails,
				ErrScanVehicleDetails:           test.ErrScanVehicleDetails,
				noOfVehicles:                    test.noOfVehicles,
				speed:                           test.speed,
				maxWeight:                       test.maxWeight,
				choice:                          test.choice,
				ErrScanProgramChoice:            test.ErrScanProgramChoice,
			}
			inputSvc := newInputSvc(mockInput)

			defer func() {
				r := recover()
				expectedOutput := test.expected.Error()
				if output.String() != expectedOutput {
					t.Errorf("Expected %v, received %v", expectedOutput, output.String())
				}
				if r == nil {
					t.Errorf("Should panic")
				}
			}()

			PackageHandler(mockWriter, mockPkgDeliveryComputeService, inputSvc)
		})
	}

}
