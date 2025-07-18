package main

import (
	"github.com/tanishkgupta-18/gofr-payment-service/handler"
	"github.com/tanishkgupta-18/gofr-payment-service/migrations"
	"github.com/tanishkgupta-18/gofr-payment-service/service"
	"github.com/tanishkgupta-18/gofr-payment-service/store"
	"gofr.dev/pkg/gofr"
)

func main() {
	app := gofr.New()

	// Register migrations
	app.Migrate(migrations.All())

	// Set up dependencies
	paymentStore := store.NewPaymentStore()
	paymentService := service.NewPaymentService(paymentStore)
	paymentHandler := handler.NewPaymentHandler(paymentService)

	// Routes
	app.POST("/payments", paymentHandler.CreatePayment)
	app.POST("/payments/callback", paymentHandler.PaymentCallback)

	app.GET("/health", func(ctx *gofr.Context) (interface{}, error) {
		return map[string]string{"status": "ok"}, nil
	})

	app.Run()
}
