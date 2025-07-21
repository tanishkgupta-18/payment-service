package handler

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/container"
	gofrHttp "gofr.dev/pkg/gofr/http"

	"github.com/tanishkgupta-18/gofr-payment-service/service"
	"github.com/tanishkgupta-18/gofr-payment-service/store"
)

// setupHandler initializes the handler with a mock context and SQL mock
func setupHandler(t *testing.T) (*PaymentHandler, *gofr.Context, sqlmock.Sqlmock) {
	mockContainer, mock := container.NewMockContainer(t)

	ctx := &gofr.Context{
		Context:   context.Background(),
		Container: mockContainer,
	}

	store := store.NewPaymentStore()
	svc := service.NewPaymentService(store)
	handler := NewPaymentHandler(svc)

	return handler, ctx, mock.SQL
}

// -------------------- CreatePayment Tests --------------------

func TestCreatePayment_InvalidJSON(t *testing.T) {
	h, ctx, _ := setupHandler(t)

	body := `{"amount":100.5,"status":"initiated"` 
	req := httptest.NewRequest(http.MethodPost, "/payments", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	ctx.Request = gofrHttp.NewRequest(req)

	_, err := h.CreatePayment(ctx)
	assert.Error(t, err)
}

func TestCreatePayment_DBError(t *testing.T) {
	h, ctx, mock := setupHandler(t)

	mock.ExpectExec("INSERT INTO payments").
		WithArgs(100.5, "initiated").
		WillReturnError(sql.ErrConnDone)

	body := `{"amount":100.5,"status":"initiated"}`
	req := httptest.NewRequest(http.MethodPost, "/payments", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	ctx.Request = gofrHttp.NewRequest(req)

	_, err := h.CreatePayment(ctx)
	assert.Equal(t, sql.ErrConnDone, err)
}

func TestCreatePayment_Success(t *testing.T) {
	h, ctx, mock := setupHandler(t)

	mock.ExpectExec("INSERT INTO payments").
		WithArgs(100.5, "initiated").
		WillReturnResult(sqlmock.NewResult(12, 1))

	body := `{"amount":100.5,"status":"initiated"}`
	req := httptest.NewRequest(http.MethodPost, "/payments", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	ctx.Request = gofrHttp.NewRequest(req)

	resp, err := h.CreatePayment(ctx)
	assert.NoError(t, err)
	assert.Equal(t, map[string]int{"paymentID": 12}, resp)
}

// -------------------- PaymentCallback Tests --------------------

func TestPaymentCallback_InvalidID(t *testing.T) {
	h, ctx, _ := setupHandler(t)

	req := httptest.NewRequest(http.MethodPost, "/payments/callback?id=abc&status=ok", nil)
	ctx.Request = gofrHttp.NewRequest(req)

	_, err := h.PaymentCallback(ctx)
	assert.Error(t, err)
}

func TestPaymentCallback_DBError(t *testing.T) {
	h, ctx, mock := setupHandler(t)

	mock.ExpectExec("UPDATE payments SET status = ? WHERE id = ?").
		WithArgs("completed", 1).
		WillReturnError(errors.New("fail"))

	req := httptest.NewRequest(http.MethodPost, "/payments/callback?id=1&status=completed", nil)
	ctx.Request = gofrHttp.NewRequest(req)

	_, err := h.PaymentCallback(ctx)
	assert.Error(t, err)
}

func TestPaymentCallback_Success(t *testing.T) {
	h, ctx, mock := setupHandler(t)

	mock.ExpectExec("UPDATE payments SET status = ? WHERE id = ?").
		WithArgs("completed", 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	req := httptest.NewRequest(http.MethodPost, "/payments/callback?id=1&status=completed", nil)
	ctx.Request = gofrHttp.NewRequest(req)

	_, err := h.PaymentCallback(ctx)
	assert.NoError(t, err)
}
