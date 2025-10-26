package account

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
)

type mockRepository struct {
	createFunc func(ctx context.Context, documentNumber string) (Account, error)
}

func (m *mockRepository) Create(ctx context.Context, documentNumber string) (Account, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, documentNumber)
	}
	return Account{}, errors.New("not implemented")
}

func TestHandler_Create(t *testing.T) {
	tests := []struct {
		name           string
		body           interface{}
		setupMock      func(*mockRepository)
		expectedStatus int
	}{
		{
			name: "valid request",
			body: CreateRequest{
				DocumentNumber: "12345678900",
			},
			setupMock: func(m *mockRepository) {
				m.createFunc = func(ctx context.Context, documentNumber string) (Account, error) {
					return Account{
						ID:             1,
						DocumentNumber: documentNumber,
						CreatedAt:      time.Now(),
					}, nil
				}
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "empty body",
			body:           map[string]interface{}{},
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid json",
			body:           "invalid json",
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "missing document_number",
			body: map[string]interface{}{
				"other_field": "value",
			},
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid document_number type",
			body: map[string]interface{}{
				"document_number": 123,
			},
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "repository error",
			body: CreateRequest{
				DocumentNumber: "12345678900",
			},
			setupMock: func(m *mockRepository) {
				m.createFunc = func(ctx context.Context, documentNumber string) (Account, error) {
					return Account{}, errors.New("database error")
				}
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockRepository{}
			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}
			handler := NewHandler(validator.New(), mockRepo)

			var bodyBytes []byte
			var err error

			if strBody, ok := tt.body.(string); ok {
				bodyBytes = []byte(strBody)
			} else {
				bodyBytes, err = json.Marshal(tt.body)
				if err != nil {
					t.Fatalf("failed to marshal request body: %v", err)
				}
			}

			req := httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.Create(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d. Response body: %s", tt.expectedStatus, w.Code, w.Body.String())
			}

			if tt.expectedStatus == http.StatusCreated {
				var response GetResponse
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}

				if response.AccountId == 0 {
					t.Error("expected non-zero account_id in response")
				}

				if response.DocumentNumber == "" {
					t.Error("expected non-empty document_number in response")
				}
			}
		})
	}
}
