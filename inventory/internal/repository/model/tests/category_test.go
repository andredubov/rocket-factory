package tests

import (
	"testing"

	repoModel "github.com/andredubov/rocket-factory/inventory/internal/repository/model"
)

func TestNewPartCategory(t *testing.T) {
	tests := []struct {
		name     string
		input    int32
		expected repoModel.PartCategory
		wantErr  bool
	}{
		{
			name:     "Unknown category",
			input:    0,
			expected: repoModel.PartCategoryUnknown,
			wantErr:  false,
		},
		{
			name:     "Engine category",
			input:    1,
			expected: repoModel.PartCategoryEngine,
			wantErr:  false,
		},
		{
			name:     "Fuel category",
			input:    2,
			expected: repoModel.PartCategoryFuel,
			wantErr:  false,
		},
		{
			name:     "Porthole category",
			input:    3,
			expected: repoModel.PartCategoryPorthole,
			wantErr:  false,
		},
		{
			name:     "Wing category",
			input:    4,
			expected: repoModel.PartCategoryWing,
			wantErr:  false,
		},
		{
			name:     "Invalid category negative number",
			input:    -1,
			expected: repoModel.PartCategoryUnknown,
			wantErr:  true,
		},
		{
			name:     "Invalid category large number",
			input:    100,
			expected: repoModel.PartCategoryUnknown,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repoModel.NewPartCategory(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPartCategory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.expected {
				t.Errorf("NewPartCategory() = %v, want %v", got, tt.expected)
			}
		})
	}
}
