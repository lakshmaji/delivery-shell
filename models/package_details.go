package models

type OfferCode string
type PackageID string
type Weight = float64
type Distance = float64

type PackageDetails struct {
	Id          PackageID
	Weight      Weight
	Distance    Distance
	Code        OfferCode // offer code which is applied on this package
	DeliveredIn float64
}

type BaseDeliveryCost float64

func (p *PackageDetails) IsSamePackage(box PackageDetails) bool {
	return p.Id == box.Id
}

func (p *PackageDetails) IsValid() bool {
	if p.Weight <= 0 || p.Distance <= 0 || p.Id == "" {
		return false
	}
	return true
}

type PackageDeliveryTime map[PackageID]float64
