package models

const (
	LessThan           = "lessThan"
	GreaterThanOrEqual = "greaterThanOrEqual"
	LessThanOrEqual    = "lessThanOrEqual"
)

type Condition struct {
	Fact     string  `json:"fact"`     // distance weight
	Operator string  `json:"operator"` // lessThan greaterThanOrEqual lessThanOrEqual
	Value    float64 `json:"value"`
}

type Offer struct {
	Code       OfferCode
	Conditions []Condition
	Discount   float64
}

type Facts []string

func (o Offer) FactsToValidate() Facts {
	return []string{"distance", "weight"}
}

type Fact struct {
	Distance Distance
	Weight   Weight
}
