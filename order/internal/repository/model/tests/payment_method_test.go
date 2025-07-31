package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/andredubov/rocket-factory/order/internal/model"
	repoModel "github.com/andredubov/rocket-factory/order/internal/repository/model"
)

func TestNewPaymentMethod(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    repoModel.PaymentMethod
		expectedErr error
	}{
		{
			name:        "valid UNKNOWN",
			input:       "UNKNOWN",
			expected:    repoModel.PaymentMethodUnknown,
			expectedErr: nil,
		},
		{
			name:        "valid CARD",
			input:       "CARD",
			expected:    repoModel.PaymentMethodCard,
			expectedErr: nil,
		},
		{
			name:        "valid SBP",
			input:       "SBP",
			expected:    repoModel.PaymentMethodSBP,
			expectedErr: nil,
		},
		{
			name:        "valid CREDIT_CARD",
			input:       "CREDIT_CARD",
			expected:    repoModel.PaymentMethodCreditCard,
			expectedErr: nil,
		},
		{
			name:        "valid INVESTOR_MONEY",
			input:       "INVESTOR_MONEY",
			expected:    repoModel.PaymentMethodInvestorMoney,
			expectedErr: nil,
		},
		{
			name:        "invalid payment method",
			input:       "INVALID_METHOD",
			expected:    repoModel.PaymentMethodUnknown,
			expectedErr: model.ErrInvalidPaymentMethod,
		},
		{
			name:        "empty payment method",
			input:       "",
			expected:    repoModel.PaymentMethodUnknown,
			expectedErr: model.ErrInvalidPaymentMethod,
		},
		{
			name:        "case sensitive check",
			input:       "card", // lowercase
			expected:    repoModel.PaymentMethodUnknown,
			expectedErr: model.ErrInvalidPaymentMethod,
		},
		{
			name:        "partial match should fail",
			input:       "CREDIT", // partial match
			expected:    repoModel.PaymentMethodUnknown,
			expectedErr: model.ErrInvalidPaymentMethod,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repoModel.NewPaymentMethod(tt.input)
			assert.Equal(t, tt.expected, result)
			assert.ErrorIs(t, err, tt.expectedErr)
		})
	}
}
