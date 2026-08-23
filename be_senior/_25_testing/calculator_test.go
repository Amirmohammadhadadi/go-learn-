// ============================================================================
// FILE: calculator_test.go
// ============================================================================

package __test_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalculator(t *testing.T) {
	calc := &Calculator{}

	t.Run("Add", func(t *testing.T) {
		tests := []struct {
			name string
			a, b int
			want int
		}{
			{"positive", 2, 3, 5},
			{"negative", -2, -3, -5},
			{"zero", 0, 5, 5},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got := calc.Add(tt.a, tt.b)
				assert.Equal(t, tt.want, got)
			})
		}
	})

	t.Run("Divide", func(t *testing.T) {
		tests := []struct {
			name      string
			a, b      int
			want      int
			expectErr bool
		}{
			{"normal", 10, 2, 5, false},
			{"by zero", 10, 0, 0, true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := calc.Divide(tt.a, tt.b)
				if tt.expectErr {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
					assert.Equal(t, tt.want, got)
				}
			})
		}
	})
}

// ============================================================================
// FILE: benchmark_test.go
// ============================================================================

func BenchmarkCalculatorAdd(b *testing.B) {
	calc := &Calculator{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calc.Add(i, i+1)
	}
}

// ============================================================================
// FILE: user_service_test.go
// ============================================================================

func TestUserService(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	t.Run("GetUser success", func(t *testing.T) {
		expected := &User{ID: 1, Name: "Ali"}
		mockRepo.On("GetByID", 1).Return(expected, nil)

		user, err := service.GetUser(1)

		assert.NoError(t, err)
		assert.Equal(t, expected, user)
		mockRepo.AssertExpectations(t)
	})
}

/*
# اجرای همه تست‌ها
go test -v

# اجرای تست خاص
go test -v -run TestCalculator

# اجرای benchmarkها
go test -bench=. -benchmem

# اجرای با race detector
go test -race -v

# نمایش پوشش تست
go test -cover -coverprofile=coverage.out
go tool cover -html=coverage.out

# اجرای با short mode (برای skip کردن تست‌های طولانی)
go test -short -v
*/
