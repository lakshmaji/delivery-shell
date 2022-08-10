package models

import (
	"fmt"

	"github.com/lakshmaji/delivery-shell/utils/msg_utils"
)

type PackageStats struct {
	Id                PackageID
	Discount          float64
	TotalDeliveryCost float64
	EstDeliveryTime   float64
}

type PackageStatsList []PackageStats

// Convert PackageStats to string, so that it can be written to stdout
// This way we can allow loose coupling among clients implementations (stdout, http etc)
func (pList PackageStatsList) FmtOutput(computesDeliveryTime bool) string {
	var finalStr string
	finalStr = fmt.Sprintf(msg_utils.MsgPackageStatsHeader)
	if computesDeliveryTime {
		finalStr += fmt.Sprintf(", %s", msg_utils.MsgPackageStatsEstTime)
	}
	finalStr += "\n"
	for _, pkg := range pList {
		finalStr += fmt.Sprintf("%s, %.2f, %.2f", pkg.Id, pkg.Discount, pkg.TotalDeliveryCost)
		if computesDeliveryTime {
			finalStr += fmt.Sprintf(", %.2f", pkg.EstDeliveryTime)
		}
		finalStr += "\n"
	}
	return finalStr
}
