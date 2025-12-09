package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/container"
	gofrHttp "gofr.dev/pkg/gofr/http"

	"github.com/tanishkgupta-18/gofr-payment-service/service"
	storePkg "github.com/tanishkgupta-18/gofr-payment-service/store"
)

var errDBFailure = errors.New("mock db failure")

// setupHandler initializes the handler with a mock context and SQL mock.
func setupHandler(t *testing.T) (*PaymentHandler, *gofr.Context, sqlmock.Sqlmock) {
	t.Helper()

	mockContainer, mock := container.NewMockContainer(t)

	ctx := &gofr.Context{
		Context:   context.Background(),
		Container: mockContainer,
	}

	s := storePkg.NewPaymentStore()
	svc := service.NewPaymentService(s)
	handler := NewPaymentHandler(svc)

	return handler, ctx, mock.SQL
}

func TestCreatePayment_InvalidJSON(t *testing.T) {
	h, ctx, _ := setupHandler(t)

	body := `{"amount":100.5,"status":"initiated"`
	req := httptest.NewRequest(http.MethodPost, "/payments", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	ctx.Request = gofrHttp.NewRequest(req)

	_, err := h.CreatePayment(ctx)
	require.Error(t, err)

	var syntaxErr *json.SyntaxError

	assert.ErrorAs(t, err, &syntaxErr)
}

func TestCreatePayment_DBError(t *testing.T) {
	h, ctx, mock := setupHandler(t)

	mock.ExpectExec("INSERT INTO payments").
		WithArgs(100.5, "initiated").
		WillReturnError(errDBFailure)

	body := `{"amount":100.5,"status":"initiated"}`
	req := httptest.NewRequest(http.MethodPost, "/payments", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	ctx.Request = gofrHttp.NewRequest(req)

	_, err := h.CreatePayment(ctx)
	require.Error(t, err)
	assert.Equal(t, errDBFailure, err)
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

	result, err := h.CreatePayment(ctx)
	require.NoError(t, err)
	assert.Equal(t, map[string]int{"paymentID": 12}, result)
}

func TestPaymentCallback_InvalidID(t *testing.T) {
	h, ctx, _ := setupHandler(t)

	req := httptest.NewRequest(http.MethodPost, "/payments/callback?id=abc&status=ok", http.NoBody)
	ctx.Request = gofrHttp.NewRequest(req)

	_, err := h.PaymentCallback(ctx)
	require.Error(t, err)
}

func TestPaymentCallback_DBError(t *testing.T) {
	h, ctx, mock := setupHandler(t)

	mock.ExpectExec("UPDATE payments SET status = ? WHERE id = ?").
		WithArgs("completed", 1).
		WillReturnError(errDBFailure)

	req := httptest.NewRequest(http.MethodPost, "/payments/callback?id=1&status=completed", http.NoBody)
	ctx.Request = gofrHttp.NewRequest(req)

	_, err := h.PaymentCallback(ctx)
	require.Error(t, err)
	assert.Equal(t, errDBFailure, err)
}

func TestPaymentCallback_Success(t *testing.T) {
	h, ctx, mock := setupHandler(t)

	mock.ExpectExec("UPDATE payments SET status = ? WHERE id = ?").
		WithArgs("completed", 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	req := httptest.NewRequest(http.MethodPost, "/payments/callback?id=1&status=completed", http.NoBody)
	ctx.Request = gofrHttp.NewRequest(req)

	_, err := h.PaymentCallback(ctx)
	require.NoError(t, err)
}
