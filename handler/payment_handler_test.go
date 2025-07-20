package handler

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
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

func TestCreatePayment(t *testing.T) {
	type gofrResponse struct {
		result any
		err    error
	}

	tests := []struct {
		name        string
		body        string
		mockExpect  func(sqlmock.Sqlmock)
		expectedRes gofrResponse
	}{
		{
			name:        "Invalid JSON",
			body:        `{"amount":100.5,"status":"initiated"`, // malformed
			mockExpect:  func(_ sqlmock.Sqlmock) {},
			expectedRes: gofrResponse{nil, &json.SyntaxError{}},
		},
		{
			name: "DB Insertion Error",
			body: `{"amount":100.5,"status":"initiated"}`,
			mockExpect: func(m sqlmock.Sqlmock) {
				m.ExpectExec("INSERT INTO payments (amount, status) VALUES (?, ?)").
					WithArgs(100.5, "initiated").
					WillReturnError(sql.ErrConnDone)
			},
			expectedRes: gofrResponse{nil, sql.ErrConnDone},
		},
		{
			name: "Successful Payment",
			body: `{"amount":100.5,"status":"initiated"}`,
			mockExpect: func(m sqlmock.Sqlmock) {
				m.ExpectExec("INSERT INTO payments (amount, status) VALUES (?, ?)").
					WithArgs(100.5, "initiated").
					WillReturnResult(sqlmock.NewResult(12, 1))
			},
			expectedRes: gofrResponse{map[string]int{"paymentID": 12}, nil},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, ctx, mock := setupHandler(t)
			tt.mockExpect(mock)

			req := httptest.NewRequest(http.MethodPost, "/payments", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			ctx.Request = gofrHttp.NewRequest(req)

			got, err := h.CreatePayment(ctx)
			if tt.name == "Invalid JSON" {
				assert.IsType(t, tt.expectedRes.err, err)
			} else {
				assert.Equal(t, tt.expectedRes.result, got)
				assert.Equal(t, tt.expectedRes.err, err)
			}
		})
	}
}

func TestPaymentCallback(t *testing.T) {
	tests := []struct {
		name       string
		query      string
		mockExpect func(sqlmock.Sqlmock)
		wantErr    bool
	}{
		{"Invalid ID", "id=abc&status=ok", func(_ sqlmock.Sqlmock) {}, true},
		{
			"DB Update Fails", "id=1&status=completed",
			func(m sqlmock.Sqlmock) {
				m.ExpectExec("UPDATE payments SET status = ? WHERE id = ?").
					WithArgs("completed", 1).
					WillReturnError(errors.New("fail"))
			}, true,
		},
		{
			"Success", "id=1&status=completed",
			func(m sqlmock.Sqlmock) {
				m.ExpectExec("UPDATE payments SET status = ? WHERE id = ?").
					WithArgs("completed", 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
			}, false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, ctx, mock := setupHandler(t)
			tt.mockExpect(mock)

			req := httptest.NewRequest(http.MethodPost, "/payments/callback?"+tt.query, http.NoBody)
			ctx.Request = gofrHttp.NewRequest(req)

			_, err := h.PaymentCallback(ctx)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
