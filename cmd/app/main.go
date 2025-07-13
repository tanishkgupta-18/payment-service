package main

import "gofr.dev/pkg/gofr"

func main() {
	app := gofr.New()

	app.GET("/health", func(_ *gofr.Context) (interface{}, error) {
		return map[string]string{"status": "UP"}, nil
	})

	app.Run()
}
