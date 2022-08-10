package offer_utils

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/lakshmaji/delivery-shell/models"
	"github.com/lakshmaji/delivery-shell/utils/error_utils"
)

func isValidFact(condition models.Condition, value float64) bool {
	var isValid bool
	switch condition.Operator {
	case models.LessThan:
		isValid = value < condition.Value
	case models.GreaterThanOrEqual:
		isValid = value >= condition.Value
	case models.LessThanOrEqual:
		isValid = value <= condition.Value
	}
	return isValid
}

// If any one of the condition fails to satisfy, then the offer is not applicable
func IsOfferApplicable(conditions []models.Condition, facts models.Facts, fact models.Fact) bool {
	// When both conditions and facts empty then offer is applicable
	if len(conditions) == 0 && len(facts) == 0 {
		return true
	}

	if len(conditions) == 0 {
		return false
	}
	if fact == (models.Fact{}) {
		return false
	}
	var isApplicable bool = true
	for _, condition := range conditions {
		var isValid bool
		switch condition.Fact {
		case "distance":
			isValid = isValidFact(condition, fact.Distance)
		case "weight":
			isValid = isValidFact(condition, fact.Weight)
		}
		isApplicable = isValid && isApplicable
	}
	return isApplicable
}

func LoadOffers(filename string) ([]models.Offer, error) {
	if len(strings.TrimSpace(filename)) == 0 {
		return nil, error_utils.ErrMissingInput
	}
	var OffersSlice []models.Offer
	content, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(content, &OffersSlice)
	if err != nil {
		return nil, err
	}
	return OffersSlice, nil
}
