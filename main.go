package main

import (
	"os"

	"github.com/lakshmaji/delivery-shell/clients"
	"github.com/lakshmaji/delivery-shell/handlers"
	"github.com/lakshmaji/delivery-shell/services/delivery_svc"
	"github.com/lakshmaji/delivery-shell/services/offers_svc"
	"github.com/lakshmaji/delivery-shell/services/shell_io_svc"
	"github.com/lakshmaji/delivery-shell/utils/offer_utils"
)

func main() {
	// production or development
	appEnv := os.Getenv("APP_ENVIRONMENT")
	if appEnv == "" {
		appEnv = "production"
	}
	// IO (std)
	reader := shell_io_svc.NewShellReader(os.Stdin)
	writer := clients.NewShellWriter(os.Stdout, appEnv == "development")

	// Deps (go-way)
	offers_svc_with_data := offers_svc.NewOffersService(offer_utils.LoadOffers)
	delivery_svc := delivery_svc.NewDeliveryService(offers_svc_with_data)

	handlers.PackageHandler(writer, delivery_svc, reader)
}
