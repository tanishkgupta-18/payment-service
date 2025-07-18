package service

import (
	"github.com/tanishkgupta-18/gofr-payment-service/store"
	"gofr.dev/pkg/gofr"
)

type PaymentService struct {
	store *store.PaymentStore
}

func NewPaymentService(s *store.PaymentStore) *PaymentService {
	return &PaymentService{store: s}
}

func (s *PaymentService) CreatePayment(ctx *gofr.Context, p *store.Payment) (int, error) {
	return s.store.CreatePayment(ctx, p)
}

func (s *PaymentService) PaymentCallback(ctx *gofr.Context, id int, status string) error {
	return s.store.UpdatePaymentStatus(ctx, id, status)
}
