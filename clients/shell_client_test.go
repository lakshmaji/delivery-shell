package clients

import (
	"bytes"
	"testing"
)

func TestWrite(t *testing.T) {

	tt := []struct {
		description string
		input       interface{}
		expected    interface{}
	}{
		{
			description: "TestWrite",
			input:       "Hello World",
			expected:    "Hello World\n",
		},
		{
			description: "TestWrite",
			input:       10,
			expected:    "10\n",
		},
		{
			description: "TestWrite",
			input:       true,
			expected:    "true\n",
		},
	}

	for _, tc := range tt {
		t.Run(tc.description, func(t *testing.T) {
			var output bytes.Buffer
			NewShellWriter(&output, true).Write(tc.input)
			if output.String() != tc.expected {
				t.Errorf("Expected %v, got %v", tc.expected, output)
			}
		})
	}
}

func TestWriteWithMultiLineString(t *testing.T) {

	tt := []struct {
		description string
		input       interface{}
		expected    interface{}
	}{
		{
			description: "TestWriteWithMultiLineString",
			input:       "Hello World\nHello World",
			expected:    "Hello World\nHello World\n",
		},
	}

	for _, tc := range tt {
		t.Run(tc.description, func(t *testing.T) {
			var output bytes.Buffer
			NewShellWriter(&output, true).Write(tc.input)
			if output.String() != tc.expected {
				t.Errorf("Expected %v, got %v", tc.expected, output)
			}
		})
	}
}

func TestWriteWithPackageStruct(t *testing.T) {
	str := ("Package Id, Discount, Total Delivery Cost\n")
	str += ("pkg1, 0.00, 175.00\n")
	str += ("pkg2, 0.00, 275.00\n")
	str += ("pkg3, 35.00, 665.00")

	expected := "Package Id, Discount, Total Delivery Cost\npkg1, 0.00, 175.00\npkg2, 0.00, 275.00\npkg3, 35.00, 665.00\n"
	var output bytes.Buffer
	NewShellWriter(&output, true).Write(str)

	if output.String() != expected {
		t.Errorf("Expected %v, got %v", str, output)
	}
}

func TestWriteError(t *testing.T) {

	tt := []struct {
		description string
		input       interface{}
		expected    interface{}
	}{
		{
			description: "TestWrite",
			input:       "Hello World",
			expected:    "Hello World",
		},
		{
			description: "TestWrite",
			input:       10,
			expected:    "10",
		},
		{
			description: "TestWrite",
			input:       true,
			expected:    "true",
		},
	}

	for _, tc := range tt {
		t.Run(tc.description, func(t *testing.T) {
			var output bytes.Buffer
			defer func() {
				r := recover()
				if output.String() != tc.expected {
					t.Errorf("Expected %v, got %v", tc.expected, output)
				}
				if r == nil {
					t.Errorf("Should panic")
				}
			}()
			NewShellWriter(&output, true).WriteError(tc.input)
		})
	}
}
