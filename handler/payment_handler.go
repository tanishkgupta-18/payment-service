package handler

import (
	"strconv"

	"gofr.dev/pkg/gofr"

	"github.com/tanishkgupta-18/gofr-payment-service/service"
	"github.com/tanishkgupta-18/gofr-payment-service/store"
)

type PaymentHandler struct {
	service *service.PaymentService
}

func NewPaymentHandler(s *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{service: s}
}

func (h *PaymentHandler) CreatePayment(ctx *gofr.Context) (interface{}, error) {
	var p store.Payment

	err := ctx.Bind(&p)
	if err != nil {
		return nil, err
	}

	id, err := h.service.CreatePayment(ctx, &p)
	if err != nil {
		return nil, err
	}

	return map[string]int{"paymentID": id}, nil
}

func (h *PaymentHandler) PaymentCallback(ctx *gofr.Context) (interface{}, error) {
	idStr := ctx.Param("id")
	status := ctx.Param("status")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, err
	}

	err = h.service.PaymentCallback(ctx, id, status)
	if err != nil {
		return nil, err
	}

	return map[string]string{"message": "Payment status updated"}, nil
}
