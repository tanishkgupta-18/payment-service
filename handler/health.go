package handler

import "gofr.dev/pkg/gofr"

type Health struct{}

func New() *Health {
	return &Health{}
}

func (*Health) Health(_ *gofr.Context) (any, error) {
	return map[string]string{"status": "UP"}, nil
}
