package main

import (
	"github.com/tanishkgupta-18/gofr-payment-service/handler"
	"github.com/tanishkgupta-18/gofr-payment-service/migrations"
	"gofr.dev/pkg/gofr"
)

func main() {
	app := gofr.New()

	app.Migrate(migrations.All())

	h := handler.New()
	app.GET("/health", h.Health)

	app.Run()
}
