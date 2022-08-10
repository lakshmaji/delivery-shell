package shell_io_svc

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"testing"

	"github.com/lakshmaji/delivery-shell/clients"
	"github.com/lakshmaji/delivery-shell/utils/error_utils"
)

func mockIO(t testing.TB) (*os.File, clients.BaseWriter, PackageInputService) {
	t.Helper()
	var output bytes.Buffer
	writer := clients.NewShellWriter(&output, true)

	reader, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}

	svc := NewShellReader(reader)

	return reader, writer, svc
}

func seek(t testing.TB, reader *os.File) {
	t.Helper()
	_, err := reader.Seek(0, io.SeekStart)
	if err != nil {
		t.Fatal(err)
	}
}

func writeToPrompt(t testing.TB, reader *os.File, input string) {
	_, err := io.WriteString(reader, input)
	if err != nil {
		t.Fatal(err)
	}
	seek(t, reader)
}

func TestScanBaseDeliveryCostPkgCount(t *testing.T) {
	reader, writer, svc := mockIO(t)
	defer reader.Close()

	writeToPrompt(t, reader, "100 10\nPKG 1 10 10\nPKG 2 10 10\nPKG 3 10 10\nPKG 4 10 10\nPKG 5 10 10\nPKG 6 10 10\nPKG 7 10 10\nPKG 8 10 10\nPKG 9 10 10\nPKG 10 10 10\n")

	baseDeliveryCost, noOfPackages, err := svc.ScanBaseDeliveryCostPkgCount(writer)
	if err != nil {
		t.Error("should not return error")
	}
	if noOfPackages != 10 {
		t.Errorf("Expected 10 packages, got %d", noOfPackages)
	}
	if baseDeliveryCost != 100 {
		t.Errorf("Expected base delivery cost 100, got %f", baseDeliveryCost)
	}

}

func TestScanBaseDeliveryCostPkgCountErrors(t *testing.T) {
	reader, writer, svc := mockIO(t)
	defer reader.Close()

	tt := []struct {
		Name     string
		Input    string
		Expected error
	}{
		{
			Name:     "no input",
			Input:    "",
			Expected: error_utils.ErrMissingInput,
		},
		{
			Name:     "provided base delivery cost only",
			Input:    "100\n",
			Expected: error_utils.ErrBaseCostPkgCount,
		},
		{
			Name:     "provided more inputs than expected",
			Input:    "100 10 10\n",
			Expected: error_utils.ErrBaseCostPkgCount,
		},
		{
			Name:     "provided cost as word and count as number",
			Input:    "ten 10\n",
			Expected: &strconv.NumError{Func: "Atoi", Num: "ten", Err: strconv.ErrSyntax},
		},
		{
			Name:     "provided cost as number and count as word",
			Input:    "100 six\n",
			Expected: &strconv.NumError{Func: "Atoi", Num: "six", Err: strconv.ErrSyntax},
		},
	}

	for _, test := range tt {
		t.Run(test.Name, func(t *testing.T) {
			seek(t, reader)

			writeToPrompt(t, reader, test.Input)

			cost, count, err := svc.ScanBaseDeliveryCostPkgCount(writer)

			if err == nil {
				t.Error("should throw error")
			}
			if count != 0 || float64(cost) != float64(0) {
				t.Errorf("cost and count should be defaults %d, %f", count, cost)
			}

			if err.Error() != test.Expected.Error() {
				t.Errorf("expected %v, received %v", test.Expected, err)
			}
		})
	}

}

func TestScanVehicleDetails(t *testing.T) {

	reader, writer, svc := mockIO(t)
	defer reader.Close()

	writeToPrompt(t, reader, "2 70 200\n")

	vehiclesCount, speed, maxWeight, err := svc.ScanVehicleDetails(writer)
	if err != nil {
		t.Error("should not return error")
	}
	if vehiclesCount != 2 {
		t.Errorf("Expected 2 vehicles, got %d", vehiclesCount)
	}
	if speed != 70 {
		t.Errorf("Expected speed 70, got %d", speed)
	}
	if maxWeight != 200 {
		t.Errorf("Expected weight capacity 200, got %d", maxWeight)
	}

}
func TestScanVehicleDetailsErrors(t *testing.T) {
	reader, writer, svc := mockIO(t)
	defer reader.Close()

	tt := []struct {
		Name     string
		Input    string
		Expected error
	}{
		{
			Name:     "no input",
			Input:    "",
			Expected: error_utils.ErrMissingInput,
		},
		{
			Name:     "provided vehicle count only",
			Input:    "3\n",
			Expected: error_utils.ErrVehicleDetailsFormat,
		},
		{
			Name:     "provided more inputs than expected",
			Input:    "3 70 10 60 70\n",
			Expected: error_utils.ErrVehicleDetailsFormat,
		},
		{
			Name:     "provided vehicle count as string",
			Input:    "ten 70 200\n",
			Expected: &strconv.NumError{Func: "Atoi", Num: "ten", Err: strconv.ErrSyntax},
		},
		{
			Name:     "provided vehicle speed as string",
			Input:    "10 seventy 200\n",
			Expected: &strconv.NumError{Func: "Atoi", Num: "seventy", Err: strconv.ErrSyntax},
		},
		{
			Name:     "provided vehicle weight capacity as string",
			Input:    "10 70 twenty\n",
			Expected: &strconv.NumError{Func: "Atoi", Num: "twenty", Err: strconv.ErrSyntax},
		},
		{
			Name:     "provided vehicle weight capacity as decimal number",
			Input:    "10 70 20.8\n",
			Expected: &strconv.NumError{Func: "Atoi", Num: "20.8", Err: strconv.ErrSyntax},
		},
	}

	for _, test := range tt {
		t.Run(test.Name, func(t *testing.T) {

			seek(t, reader)

			writeToPrompt(t, reader, test.Input)

			vehiclesCount, speed, maxWeight, err := svc.ScanVehicleDetails(writer)

			if err == nil {
				t.Error("should throw error")
			}
			if vehiclesCount != 0 || speed != 0 || maxWeight != 0 {
				t.Errorf("expected defaults %d, %d, %d", vehiclesCount, speed, maxWeight)
			}

			if err.Error() != test.Expected.Error() {
				t.Errorf("expected %v, received %v", test.Expected, err)
			}
		})
	}

}

func TestScanNPackageDetails(t *testing.T) {
	reader, writer, svc := mockIO(t)
	defer reader.Close()

	writeToPrompt(t, reader, "PKG 1 10 10\nPKG 2 10 10\nPKG 3 10 10\nPKG 4 10 10\nPKG 5 10 10\nPKG 6 10 10\nPKG 7 10 10\nPKG 8 10 10\nPKG 9 10 10\nPKG 10 10 10\n")

	boxes, err := svc.ScanNPackageDetails(writer, 10)
	if err != nil {
		t.Error("should not return error")
	}
	if len(boxes) != 10 {
		t.Errorf("Expected 10 boxes, got %d", len(boxes))
	}

}
func TestScanNPackageDetailsErrors(t *testing.T) {
	reader, writer, svc := mockIO(t)
	defer reader.Close()

	tt := []struct {
		Name         string
		Input        string
		noOfPackages int
		Expected     error
	}{
		{
			Name:         "no input",
			Input:        "",
			Expected:     error_utils.ErrMissingInput,
			noOfPackages: 3,
		},
		{
			Name:         "Missing weight, distance and offer code",
			Input:        "PKG1",
			Expected:     error_utils.ErrPackageDetailsFormat,
			noOfPackages: 3,
		},
		{
			Name:         "Missing distance",
			Input:        "PKG1\t10\t\tOFR002",
			Expected:     error_utils.ErrPackageDetailsFormat,
			noOfPackages: 3,
		},
	}

	for _, test := range tt {
		t.Run(test.Name, func(t *testing.T) {

			writeToPrompt(t, reader, test.Input)

			boxes, err := svc.ScanNPackageDetails(writer, test.noOfPackages)

			if err == nil {
				t.Error("should throw error")
			}
			if len(boxes) != 0 {
				t.Errorf("Expected 0 boxes, got %d", len(boxes))
			}

			if err.Error() != test.Expected.Error() {
				t.Errorf("expected %v, received %v", test.Expected, err)
			}
		})
	}

}

func TestScanProgramChoice(t *testing.T) {
	reader, writer, svc := mockIO(t)
	defer reader.Close()

	tt := []struct {
		Name     string
		Input    string
		Expected error
	}{
		{
			Name:  "yes",
			Input: "yes\n",
		},
		{
			Name:  "no",
			Input: "no\n",
		},
	}

	for _, test := range tt {

		writeToPrompt(t, reader, test.Input)

		ans, err := svc.ScanProgramChoice(writer)

		if err != nil {
			t.Error("should not throw error")
		}
		if ans == test.Input {
			t.Errorf("choice should be %s received  %s", test.Input, ans)
		}
	}

}
func TestScanProgramChoiceErrors(t *testing.T) {
	reader, writer, svc := mockIO(t)
	defer reader.Close()

	tt := []struct {
		Name     string
		Input    string
		Expected error
	}{
		{
			Name:     "no input",
			Input:    "",
			Expected: error_utils.ErrMissingInput,
		},
		{
			Name:     "format",
			Input:    "mockIO",
			Expected: error_utils.ErrProgramChoiceFormat,
		},
	}

	for _, test := range tt {
		t.Run(test.Name, func(t *testing.T) {
			writeToPrompt(t, reader, test.Input)

			ans, err := svc.ScanProgramChoice(writer)

			if err == nil {
				t.Error("should throw error")
			}
			if ans == "yes" || ans == "no" {
				t.Errorf("choice should not ne reader yes.no received  %v", ans)
			}

			if err.Error() != test.Expected.Error() {
				t.Errorf("expected %v, received %v", test.Expected, err)
			}
		})
	}

}
