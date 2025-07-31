package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/andredubov/rocket-factory/order/internal/model"
	repoModel "github.com/andredubov/rocket-factory/order/internal/repository/model"
)

func TestNewOrderStatus(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    repoModel.OrderStatus
		expectedErr error
	}{
		{
			name:        "valid PENDING_PAYMENT",
			input:       "PENDING_PAYMENT",
			expected:    repoModel.OrderStatusPending,
			expectedErr: nil,
		},
		{
			name:        "valid PAID",
			input:       "PAID",
			expected:    repoModel.OrderStatusPaid,
			expectedErr: nil,
		},
		{
			name:        "valid CANCELLED",
			input:       "CANCELLED",
			expected:    repoModel.OrderStatusCancelled,
			expectedErr: nil,
		},
		{
			name:        "invalid status",
			input:       "INVALID_STATUS",
			expected:    repoModel.OrderStatusPending,
			expectedErr: model.ErrInvalidOrderStatus,
		},
		{
			name:        "empty status",
			input:       "",
			expected:    repoModel.OrderStatusPending,
			expectedErr: model.ErrInvalidOrderStatus,
		},
		{
			name:        "case sensitive check",
			input:       "paid", // lowercase
			expected:    repoModel.OrderStatusPending,
			expectedErr: model.ErrInvalidOrderStatus,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repoModel.NewOrderStatus(tt.input)
			assert.Equal(t, tt.expected, result)
			assert.ErrorIs(t, err, tt.expectedErr)
		})
	}
}
