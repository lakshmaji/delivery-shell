package models

import "testing"

func TestIsValid(t *testing.T) {

	tt := []struct {
		name     string
		expected bool
		input    PackageDetails
	}{
		{
			name:     "valid product attributes",
			expected: true,
			input: PackageDetails{
				Id:       "PKG1",
				Weight:   20,
				Distance: 20,
			},
		},
		{
			name:     "invalid weight",
			expected: false,
			input: PackageDetails{
				Id:       "PKG1",
				Weight:   0,
				Distance: 20,
			},
		},
		{
			name:     "invalid distance",
			expected: false,
			input: PackageDetails{
				Id:       "PKG1",
				Weight:   50,
				Distance: 0,
			},
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			output := test.input.IsValid()
			if output != test.expected {
				t.Errorf("should be %t received %t", test.expected, output)
			}
		})
	}
}

func TestIsSamePackage(t *testing.T) {

	tt := []struct {
		name     string
		expected bool
		input    PackageDetails
		other    PackageDetails
	}{
		{
			name:     "Same package",
			expected: true,
			input: PackageDetails{
				Id:       "PKG1",
				Weight:   20,
				Distance: 20,
			},
			other: PackageDetails{
				Id: "PKG1",
			},
		},
		{
			name:     "not same package",
			expected: false,
			input: PackageDetails{
				Id:       "PKG1",
				Weight:   0,
				Distance: 20,
			},
			other: PackageDetails{
				Id: "PKG5",
			},
		},
		{
			name:     "not valid package",
			expected: false,
			input: PackageDetails{
				Id:       "PKG1",
				Weight:   50,
				Distance: 0,
			},
			other: PackageDetails{},
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			output := test.input.IsSamePackage(test.other)
			if output != test.expected {
				t.Errorf("should be %t received %t", test.expected, output)
			}
		})
	}
}
