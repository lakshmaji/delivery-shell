package offer_utils

import (
	"errors"
	"os"
	"testing"

	"github.com/lakshmaji/delivery-shell/models"
	"github.com/lakshmaji/delivery-shell/utils/error_utils"
)

func TestIsOfferApplicable(t *testing.T) {

	tt := []struct {
		desc       string
		conditions []models.Condition
		facts      models.Facts
		fact       models.Fact
		expected   bool
	}{
		{
			desc: "distance and weight satisfied offer code conditions",
			conditions: []models.Condition{
				{
					Fact:     "distance",
					Operator: models.GreaterThanOrEqual,
					Value:    10,
				},
				{
					Fact:     "weight",
					Operator: models.GreaterThanOrEqual,
					Value:    10,
				},
			},
			facts:    []string{"distance", "weight"},
			fact:     models.Fact{Distance: 10.3, Weight: 10.5},
			expected: true,
		},
		{
			desc: "when distance fails to satisfy",
			conditions: []models.Condition{
				{
					Fact:     "distance",
					Operator: models.GreaterThanOrEqual,
					Value:    10,
				},
				{
					Fact:     "weight",
					Operator: models.GreaterThanOrEqual,
					Value:    10,
				},
			},
			facts: []string{"distance", "weight"},
			fact: models.Fact{
				Distance: 9,
				Weight:   10,
			},
			expected: false,
		},
		{
			desc: "when weight failed to satisfy",
			conditions: []models.Condition{
				{
					Fact:     "distance",
					Operator: models.GreaterThanOrEqual,
					Value:    10,
				},
				{
					Fact:     "weight",
					Operator: models.GreaterThanOrEqual,
					Value:    10,
				},
			},
			facts: []string{"distance", "weight"},
			fact: models.Fact{
				Distance: 10,
				Weight:   9,
			},
			expected: false,
		},
		{
			desc: "when un even condition is provided",
			conditions: []models.Condition{
				{
					Fact:     "distance",
					Operator: models.GreaterThanOrEqual,
					Value:    10,
				},
				{
					Fact:     "distance",
					Operator: models.LessThan,
					Value:    10,
				},
			},
			facts: []string{"distance"},
			fact: models.Fact{
				Distance: 10,
				Weight:   22,
			},
			expected: false,
		},
		{
			desc: "with different operators",
			conditions: []models.Condition{
				{
					Fact:     "distance",
					Operator: models.LessThanOrEqual,
					Value:    10,
				},
			},
			facts: []string{"distance"},
			fact: models.Fact{
				Distance: 10,
			},
			expected: true,
		},
		{
			desc: "when values are not provided",
			conditions: []models.Condition{
				{
					Fact:     "distance",
					Operator: models.GreaterThanOrEqual,
					Value:    10,
				},
				{
					Fact:     "weight",
					Operator: models.GreaterThanOrEqual,
					Value:    10,
				},
			},
			facts:    []string{"distance", "weight"},
			fact:     models.Fact{},
			expected: false,
		},
		{
			desc:       "when conditions are zero",
			conditions: []models.Condition{},
			facts:      []string{"distance", "weight"},
			fact: models.Fact{
				Distance: 10,
				Weight:   10,
			},
			expected: false,
		},
		{
			desc:       "when facts and conditions are empty",
			conditions: []models.Condition{},
			facts:      []string{},
			fact: models.Fact{
				Distance: 10,
				Weight:   10,
			},
			expected: true,
		},
	}
	for _, test := range tt {
		t.Run(test.desc, func(t *testing.T) {
			result := IsOfferApplicable(test.conditions, test.facts, test.fact)
			if test.expected != result {
				t.Errorf("expected %t received %t", result, test.expected)
			}
		})
	}
}

func TestLoadOffers(t *testing.T) {
	// TODO: not mocking `ioutil.ReadFile`
	// This is a proper solution for now, to not to impose deps on loadOffers() function
	tt := []struct {
		desc        string
		offersFile  string
		expectedErr error
		expected    []models.Offer
	}{
		{
			desc:        "when offers file is empty",
			offersFile:  "",
			expectedErr: error_utils.ErrMissingInput,
		},
		{
			desc:       "when offers file is not empty",
			offersFile: "./testdata/offers.json",
			expected:   []models.Offer{},
		},
		{
			desc:        "when offers file is not available",
			offersFile:  "./testdata/nooffers.json",
			expectedErr: &os.PathError{Op: "open", Path: "./testdata/nooffers.json", Err: errors.New("no such file or directory")},
		},
	}
	for _, test := range tt {
		t.Run(test.desc, func(t *testing.T) {
			result, err := LoadOffers(test.offersFile)
			if test.expectedErr != nil {
				if err.Error() != test.expectedErr.Error() {
					t.Errorf("expected %v, received %v", test.expectedErr, err)
				}
			}
			if test.expected != nil {
				if len(result) != len(test.expected) {
					t.Errorf("expected %d received %d", len(result), len(test.expected))
				}
			}
		})
	}
}
