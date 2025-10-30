package transaction

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rikw22/challenge-money/internal/domain/account"
)

type mockRepository struct {
	createFunc                             func(ctx context.Context, transaction *Transaction) error
	getTransactionsWithNegativeBalanceFunc func(ctx context.Context, accountId int) ([]Transaction, error)
	updateTransactionBalanceFunc           func(ctx context.Context, uuid pgtype.UUID, balance int) error
}

func (m *mockRepository) Create(ctx context.Context, transaction *Transaction) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, transaction)
	}
	return errors.New("not implemented")
}

func (m *mockRepository) GetTransactionsWithNegativeBalance(ctx context.Context, accountId int) ([]Transaction, error) {
	if m.getTransactionsWithNegativeBalanceFunc != nil {
		return m.getTransactionsWithNegativeBalanceFunc(ctx, accountId)
	}
	return nil, errors.New("not implemented")
}

func (m *mockRepository) UpdateTransactionBalance(ctx context.Context, uuid pgtype.UUID, balance int) error {
	if m.updateTransactionBalanceFunc != nil {
		return m.updateTransactionBalanceFunc(ctx, uuid, balance)
	}
	return errors.New("not implemented")
}

type mockAccountRepository struct {
	existFunc func(ctx context.Context, id int) (bool, error)
}

func (m *mockAccountRepository) GetByID(ctx context.Context, id string) (account.Account, error) {
	return account.Account{}, errors.New("not implemented")
}

func (m *mockAccountRepository) Create(ctx context.Context, acc *account.Account) error {
	return errors.New("not implemented")
}

func (m *mockAccountRepository) Exist(ctx context.Context, id int) (bool, error) {
	if m.existFunc != nil {
		return m.existFunc(ctx, id)
	}
	return false, errors.New("not implemented")
}

type mockOperationTypeRepository struct {
	existFunc func(ctx context.Context, id int) (bool, error)
}

func (m *mockOperationTypeRepository) Exist(ctx context.Context, id int) (bool, error) {
	if m.existFunc != nil {
		return m.existFunc(ctx, id)
	}
	return false, errors.New("not implemented")
}

func TestHandler_Create(t *testing.T) {
	tests := []struct {
		name                   string
		body                   interface{}
		setupMock              func(*mockRepository)
		setupAccountMock       func(*mockAccountRepository)
		setupOperationTypeMock func(*mockOperationTypeRepository)
		expectedStatus         int
	}{
		{
			name: "valid request",
			body: CreateTransactionRequest{
				AccountId:       1,
				OperationTypeId: 4,
				Amount:          123.45,
			},
			setupMock: func(m *mockRepository) {
				m.createFunc = func(ctx context.Context, t *Transaction) error {
					// Simulate database setting ID and timestamp
					t.ID = pgtype.UUID{Bytes: [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}, Valid: true}
					t.EventDate = time.Now()
					return nil
				}
				m.getTransactionsWithNegativeBalanceFunc = func(ctx context.Context, accountId int) ([]Transaction, error) {
					return []Transaction{}, nil
				}
			},
			setupAccountMock: func(m *mockAccountRepository) {
				m.existFunc = func(ctx context.Context, id int) (bool, error) {
					return true, nil
				}
			},
			setupOperationTypeMock: func(m *mockOperationTypeRepository) {
				m.existFunc = func(ctx context.Context, id int) (bool, error) {
					return true, nil
				}
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "valid request with large account id",
			body: CreateTransactionRequest{
				AccountId:       999999,
				OperationTypeId: 1,
				Amount:          50.00,
			},
			setupMock: func(m *mockRepository) {
				m.createFunc = func(ctx context.Context, t *Transaction) error {
					t.ID = pgtype.UUID{Bytes: [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}, Valid: true}
					t.EventDate = time.Now()
					return nil
				}
			},
			setupAccountMock: func(m *mockAccountRepository) {
				m.existFunc = func(ctx context.Context, id int) (bool, error) {
					return true, nil
				}
			},
			setupOperationTypeMock: func(m *mockOperationTypeRepository) {
				m.existFunc = func(ctx context.Context, id int) (bool, error) {
					return true, nil
				}
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:                   "empty body",
			body:                   map[string]interface{}{},
			setupMock:              nil,
			setupAccountMock:       nil,
			setupOperationTypeMock: nil,
			expectedStatus:         http.StatusBadRequest,
		},
		{
			name:                   "invalid json",
			body:                   "invalid json",
			setupMock:              nil,
			setupAccountMock:       nil,
			setupOperationTypeMock: nil,
			expectedStatus:         http.StatusBadRequest,
		},
		{
			name: "missing account_id",
			body: map[string]interface{}{
				"operation_type_id": 4,
				"amount":            123.45,
			},
			setupMock:              nil,
			setupAccountMock:       nil,
			setupOperationTypeMock: nil,
			expectedStatus:         http.StatusBadRequest,
		},
		{
			name: "account_id zero",
			body: map[string]interface{}{
				"account_id":        0,
				"operation_type_id": 4,
				"amount":            123.45,
			},
			setupMock:              nil,
			setupAccountMock:       nil,
			setupOperationTypeMock: nil,
			expectedStatus:         http.StatusBadRequest,
		},
		{
			name: "account_id negative",
			body: map[string]interface{}{
				"account_id":        -1,
				"operation_type_id": 4,
				"amount":            123.45,
			},
			setupMock:              nil,
			setupAccountMock:       nil,
			setupOperationTypeMock: nil,
			expectedStatus:         http.StatusBadRequest,
		},
		{
			name: "invalid account_id type",
			body: map[string]interface{}{
				"account_id":        "invalid",
				"operation_type_id": 4,
				"amount":            123.45,
			},
			setupMock:              nil,
			setupAccountMock:       nil,
			setupOperationTypeMock: nil,
			expectedStatus:         http.StatusBadRequest,
		},
		{
			name: "account does not exist",
			body: CreateTransactionRequest{
				AccountId:       999,
				OperationTypeId: 4,
				Amount:          123.45,
			},
			setupMock: nil,
			setupAccountMock: func(m *mockAccountRepository) {
				m.existFunc = func(ctx context.Context, id int) (bool, error) {
					return false, nil
				}
			},
			setupOperationTypeMock: nil,
			expectedStatus:         http.StatusBadRequest,
		},
		{
			name: "account validation error",
			body: CreateTransactionRequest{
				AccountId:       1,
				OperationTypeId: 4,
				Amount:          123.45,
			},
			setupMock: nil,
			setupAccountMock: func(m *mockAccountRepository) {
				m.existFunc = func(ctx context.Context, id int) (bool, error) {
					return false, errors.New("database connection error")
				}
			},
			setupOperationTypeMock: nil,
			expectedStatus:         http.StatusInternalServerError,
		},
		{
			name: "operation type does not exist",
			body: CreateTransactionRequest{
				AccountId:       1,
				OperationTypeId: 5,
				Amount:          123.45,
			},
			setupMock: nil,
			setupAccountMock: func(m *mockAccountRepository) {
				m.existFunc = func(ctx context.Context, id int) (bool, error) {
					return true, nil
				}
			},
			setupOperationTypeMock: func(m *mockOperationTypeRepository) {
				m.existFunc = func(ctx context.Context, id int) (bool, error) {
					return false, nil
				}
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "operation type validation error",
			body: CreateTransactionRequest{
				AccountId:       1,
				OperationTypeId: 4,
				Amount:          123.45,
			},
			setupMock: nil,
			setupAccountMock: func(m *mockAccountRepository) {
				m.existFunc = func(ctx context.Context, id int) (bool, error) {
					return true, nil
				}
			},
			setupOperationTypeMock: func(m *mockOperationTypeRepository) {
				m.existFunc = func(ctx context.Context, id int) (bool, error) {
					return false, errors.New("database connection error")
				}
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "missing operation_type_id",
			body: map[string]interface{}{
				"account_id": 1,
				"amount":     123.45,
			},
			setupMock:              nil,
			setupAccountMock:       nil,
			setupOperationTypeMock: nil,
			expectedStatus:         http.StatusBadRequest,
		},
		{
			name: "invalid operation_type_id",
			body: map[string]interface{}{
				"account_id":        1,
				"operation_type_id": 5,
				"amount":            123.45,
			},
			setupMock:              nil,
			setupAccountMock:       nil,
			setupOperationTypeMock: nil,
			expectedStatus:         http.StatusBadRequest,
		},
		{
			name: "missing amount",
			body: map[string]interface{}{
				"account_id":        1,
				"operation_type_id": 4,
			},
			setupMock:              nil,
			setupAccountMock:       nil,
			setupOperationTypeMock: nil,
			expectedStatus:         http.StatusBadRequest,
		},
		{
			name: "repository error",
			body: CreateTransactionRequest{
				AccountId:       1,
				OperationTypeId: 4,
				Amount:          123.45,
			},
			setupMock: func(m *mockRepository) {
				m.createFunc = func(ctx context.Context, t *Transaction) error {
					return errors.New("database error")
				}
				m.getTransactionsWithNegativeBalanceFunc = func(ctx context.Context, accountId int) ([]Transaction, error) {
					return []Transaction{}, nil
				}
			},
			setupAccountMock: func(m *mockAccountRepository) {
				m.existFunc = func(ctx context.Context, id int) (bool, error) {
					return true, nil
				}
			},
			setupOperationTypeMock: func(m *mockOperationTypeRepository) {
				m.existFunc = func(ctx context.Context, id int) (bool, error) {
					return true, nil
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

			mockAccountRepo := &mockAccountRepository{}
			if tt.setupAccountMock != nil {
				tt.setupAccountMock(mockAccountRepo)
			}

			mockOperationTypeRepo := &mockOperationTypeRepository{}
			if tt.setupOperationTypeMock != nil {
				tt.setupOperationTypeMock(mockOperationTypeRepo)
			}

			handler := NewHandler(validator.New(), mockRepo, mockAccountRepo, mockOperationTypeRepo)

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

			req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.Create(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d. Response body: %s", tt.expectedStatus, w.Code, w.Body.String())
			}

			if tt.expectedStatus == http.StatusCreated {
				var response CreateTransactionResponse
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}

				if response.ID.String() == "00000000-0000-0000-0000-000000000000" {
					t.Error("expected non-zero ID in response")
				}

				if response.AccountId == 0 {
					t.Error("expected non-zero account_id in response")
				}

				if response.OperationTypeId == 0 {
					t.Error("expected non-zero operation_type_id in response")
				}

				if response.Amount == 0 {
					t.Error("expected non-zero amount in response")
				}
			}
		})
	}
}

func TestHandler_Create_AmountSign(t *testing.T) {
	tests := []struct {
		name            string
		operationTypeId int
		inputAmount     float64
		expectedAmount  int
	}{
		{
			name:            "operation type 1 converts positive to negative",
			operationTypeId: 1,
			inputAmount:     50.00,
			expectedAmount:  -5000,
		},
		{
			name:            "operation type 2 converts positive to negative",
			operationTypeId: 2,
			inputAmount:     100.50,
			expectedAmount:  -10050,
		},
		{
			name:            "operation type 3 converts positive to negative",
			operationTypeId: 3,
			inputAmount:     75.25,
			expectedAmount:  -7525,
		},
		{
			name:            "operation type 4 keeps positive amount",
			operationTypeId: 4,
			inputAmount:     200.00,
			expectedAmount:  20000,
		},
		{
			name:            "operation type 1 with decimal amount",
			operationTypeId: 1,
			inputAmount:     50.99,
			expectedAmount:  -5099,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var capturedAmount int

			mockRepo := &mockRepository{
				createFunc: func(ctx context.Context, transaction *Transaction) error {
					capturedAmount = transaction.Amount
					transaction.ID = pgtype.UUID{Bytes: [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}, Valid: true}
					transaction.EventDate = time.Now()
					return nil
				},
				getTransactionsWithNegativeBalanceFunc: func(ctx context.Context, accountId int) ([]Transaction, error) {
					return []Transaction{}, nil
				},
			}

			mockAccountRepo := &mockAccountRepository{
				existFunc: func(ctx context.Context, id int) (bool, error) {
					return true, nil
				},
			}

			mockOperationTypeRepo := &mockOperationTypeRepository{
				existFunc: func(ctx context.Context, id int) (bool, error) {
					return true, nil
				},
			}

			handler := NewHandler(validator.New(), mockRepo, mockAccountRepo, mockOperationTypeRepo)

			body := CreateTransactionRequest{
				AccountId:       1,
				OperationTypeId: tt.operationTypeId,
				Amount:          tt.inputAmount,
			}

			bodyBytes, err := json.Marshal(body)
			if err != nil {
				t.Fatalf("failed to marshal request body: %v", err)
			}

			req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.Create(w, req)

			if w.Code != http.StatusCreated {
				t.Errorf("expected status %d, got %d. Response body: %s", http.StatusCreated, w.Code, w.Body.String())
			}

			if capturedAmount != tt.expectedAmount {
				t.Errorf("expected amount %d, got %d", tt.expectedAmount, capturedAmount)
			}

			// Verify response always returns positive amount
			var response CreateTransactionResponse
			if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if response.Amount < 0 {
				t.Errorf("response amount should always be positive, got %f", response.Amount)
			}

			expectedResponseAmount := math.Abs(float64(tt.expectedAmount) / 100)
			if response.Amount != expectedResponseAmount {
				t.Errorf("expected response amount %f, got %f", expectedResponseAmount, response.Amount)
			}
		})
	}
}

func TestHandler_Create_PaymentAllocation(t *testing.T) {
	tests := []struct {
		name                     string
		paymentAmount            float64
		existingNegativeBalances []Transaction
		expectedBalanceUpdates   map[string]int
	}{
		{
			name:          "payment fully covers single debt",
			paymentAmount: 100.00,
			existingNegativeBalances: []Transaction{
				{
					ID:      pgtype.UUID{Bytes: [16]byte{1}, Valid: true},
					Balance: -10000,
				},
			},
			expectedBalanceUpdates: map[string]int{
				"01000000-0000-0000-0000-000000000000": 0,
			},
		},
		{
			name:          "payment partially covers single debt",
			paymentAmount: 50.00,
			existingNegativeBalances: []Transaction{
				{
					ID:      pgtype.UUID{Bytes: [16]byte{1}, Valid: true},
					Balance: -10000,
				},
			},
			expectedBalanceUpdates: map[string]int{
				"01000000-0000-0000-0000-000000000000": -5000,
			},
		},
		{
			name:          "payment covers multiple debts fully",
			paymentAmount: 300.00,
			existingNegativeBalances: []Transaction{
				{
					ID:      pgtype.UUID{Bytes: [16]byte{1}, Valid: true},
					Balance: -10000,
				},
				{
					ID:      pgtype.UUID{Bytes: [16]byte{2}, Valid: true},
					Balance: -15000,
				},
				{
					ID:      pgtype.UUID{Bytes: [16]byte{3}, Valid: true},
					Balance: -5000,
				},
			},
			expectedBalanceUpdates: map[string]int{
				"01000000-0000-0000-0000-000000000000": 0,
				"02000000-0000-0000-0000-000000000000": 0,
				"03000000-0000-0000-0000-000000000000": 0,
			},
		},
		{
			name:          "payment covers some debts but not all",
			paymentAmount: 180.00,
			existingNegativeBalances: []Transaction{
				{
					ID:      pgtype.UUID{Bytes: [16]byte{1}, Valid: true},
					Balance: -10000, // R$ 100,00
				},
				{
					ID:      pgtype.UUID{Bytes: [16]byte{2}, Valid: true},
					Balance: -15000, // R$ 150,00
				},
				{
					ID:      pgtype.UUID{Bytes: [16]byte{3}, Valid: true},
					Balance: -5000, // R$ 50,00
				},
			},
			expectedBalanceUpdates: map[string]int{
				"01000000-0000-0000-0000-000000000000": 0,
				"02000000-0000-0000-0000-000000000000": -7000,
			},
		},
		{
			name:          "payment covers first debt and part of second",
			paymentAmount: 120.00,
			existingNegativeBalances: []Transaction{
				{
					ID:      pgtype.UUID{Bytes: [16]byte{1}, Valid: true},
					Balance: -10000,
				},
				{
					ID:      pgtype.UUID{Bytes: [16]byte{2}, Valid: true},
					Balance: -15000,
				},
			},
			expectedBalanceUpdates: map[string]int{
				"01000000-0000-0000-0000-000000000000": 0,
				"02000000-0000-0000-0000-000000000000": -13000,
			},
		},
		{
			name:                     "payment when no debts exist",
			paymentAmount:            100.00,
			existingNegativeBalances: []Transaction{},
			expectedBalanceUpdates:   map[string]int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			balanceUpdates := make(map[string]int)

			mockRepo := &mockRepository{
				createFunc: func(ctx context.Context, transaction *Transaction) error {
					transaction.ID = pgtype.UUID{Bytes: [16]byte{99}, Valid: true}
					transaction.EventDate = time.Now()
					return nil
				},
				getTransactionsWithNegativeBalanceFunc: func(ctx context.Context, accountId int) ([]Transaction, error) {
					return tt.existingNegativeBalances, nil
				},
				updateTransactionBalanceFunc: func(ctx context.Context, uuid pgtype.UUID, balance int) error {
					balanceUpdates[uuid.String()] = balance
					return nil
				},
			}

			mockAccountRepo := &mockAccountRepository{
				existFunc: func(ctx context.Context, id int) (bool, error) {
					return true, nil
				},
			}

			mockOperationTypeRepo := &mockOperationTypeRepository{
				existFunc: func(ctx context.Context, id int) (bool, error) {
					return true, nil
				},
			}

			handler := NewHandler(validator.New(), mockRepo, mockAccountRepo, mockOperationTypeRepo)

			body := CreateTransactionRequest{
				AccountId:       1,
				OperationTypeId: 4,
				Amount:          tt.paymentAmount,
			}

			bodyBytes, err := json.Marshal(body)
			if err != nil {
				t.Fatalf("failed to marshal request body: %v", err)
			}

			req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.Create(w, req)

			if w.Code != http.StatusCreated {
				t.Errorf("expected status %d, got %d. Response body: %s", http.StatusCreated, w.Code, w.Body.String())
			}

			if len(balanceUpdates) != len(tt.expectedBalanceUpdates) {
				t.Errorf("expected %d balance updates, got %d", len(tt.expectedBalanceUpdates), len(balanceUpdates))
			}

			for uuid, expectedBalance := range tt.expectedBalanceUpdates {
				actualBalance, ok := balanceUpdates[uuid]
				if !ok {
					t.Errorf("expected balance update for transaction %s, but none found", uuid)
					continue
				}
				if actualBalance != expectedBalance {
					t.Errorf("transaction %s: expected balance %d, got %d", uuid, expectedBalance, actualBalance)
				}
			}
		})
	}
}
