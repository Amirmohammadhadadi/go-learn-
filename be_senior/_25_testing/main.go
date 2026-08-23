// ============================================================================
// FILE: testing_complete_guide.go
// TITLE: راهنمای کامل تست در Go - Table-Driven Tests, Benchmark, Testify, Mocking
// HOW TO RUN: go test -v
// HOW TO BENCHMARK: go test -bench=. -benchmem
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - تست در Go
// ============================================================================
//
// Go دارای ابزارهای تست داخلی قدرتمند است:
//
// 1. testing package: پکیج استاندارد تست
// 2. Table-driven tests: تست با جداول داده
// 3. Benchmarking: اندازه‌گیری عملکرد
// 4. Testify: کتابخانه شخص ثالث برای断言‌های خواناتر
// 5. Mocking: شبیه‌سازی وابستگی‌ها با interface
//
// قانون طلایی:
// "همیشه تست بنویس. از table-driven tests برای پوشش موارد مختلف استفاده کن.
//  از interfaceها برای قابلیت mock کردن استفاده کن.
//  همیشه با -race تست کن تا race conditions را پیدا کنی."
// ============================================================================

package __test

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// ============================================================================
// بخش 1: کدهای نمونه برای تست (Sample Code to Test)
// ============================================================================

// Calculator - یک ماشین حساب ساده برای تست
type Calculator struct{}

func (c *Calculator) Add(a, b int) int {
	return a + b
}

func (c *Calculator) Subtract(a, b int) int {
	return a - b
}

func (c *Calculator) Multiply(a, b int) int {
	return a * b
}

func (c *Calculator) Divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}

// StringUtils - توابع کاربردی روی رشته
type StringUtils struct{}

func (s *StringUtils) Reverse(str string) string {
	runes := []rune(str)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func (s *StringUtils) IsPalindrome(str string) bool {
	cleaned := strings.ToLower(strings.ReplaceAll(str, " ", ""))
	return cleaned == s.Reverse(cleaned)
}

func (s *StringUtils) CountVowels(str string) int {
	vowels := "aeiouAEIOU"
	count := 0
	for _, ch := range str {
		if strings.ContainsRune(vowels, ch) {
			count++
		}
	}
	return count
}

// UserService - سرویس کاربر با وابستگی به دیتابیس (برای mocking)
type User struct {
	ID    int
	Name  string
	Email string
}

type UserRepository interface {
	GetByID(id int) (*User, error)
	Save(user *User) error
	Delete(id int) error
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUser(id int) (*User, error) {
	if id <= 0 {
		return nil, errors.New("invalid user id")
	}
	return s.repo.GetByID(id)
}

func (s *UserService) CreateUser(name, email string) (*User, error) {
	if name == "" || email == "" {
		return nil, errors.New("name and email are required")
	}
	user := &User{
		ID:    time.Now().UnixNano(),
		Name:  name,
		Email: email,
	}
	err := s.repo.Save(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// ============================================================================
// بخش 2: تست‌های پایه (Basic Tests)
// ============================================================================

// TestAdd - تست ساده
func TestAdd(t *testing.T) {
	calc := &Calculator{}
	result := calc.Add(2, 3)
	expected := 5

	if result != expected {
		t.Errorf("Add(2,3) = %d; want %d", result, expected)
	}
}

// TestSubtract - تست با sub-tests
func TestSubtract(t *testing.T) {
	calc := &Calculator{}

	t.Run("positive numbers", func(t *testing.T) {
		result := calc.Subtract(10, 3)
		if result != 7 {
			t.Errorf("got %d, want 7", result)
		}
	})

	t.Run("negative result", func(t *testing.T) {
		result := calc.Subtract(3, 10)
		if result != -7 {
			t.Errorf("got %d, want -7", result)
		}
	})
}

// ============================================================================
// بخش 3: Table-Driven Tests (تست با جدول داده)
// ============================================================================

// TestAddTableDriven - تست Add با روش table-driven
func TestAddTableDriven(t *testing.T) {
	calc := &Calculator{}

	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"positive numbers", 2, 3, 5},
		{"negative numbers", -2, -3, -5},
		{"positive and negative", 5, -3, 2},
		{"zero values", 0, 0, 0},
		{"large numbers", 1000000, 2000000, 3000000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calc.Add(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Add(%d, %d) = %d; want %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

// TestDivideTableDriven - تست Divide با table-driven و مدیریت خطا
func TestDivideTableDriven(t *testing.T) {
	calc := &Calculator{}

	tests := []struct {
		name      string
		a, b      int
		expected  int
		expectErr bool
		errMsg    string
	}{
		{"normal division", 10, 2, 5, false, ""},
		{"division by zero", 10, 0, 0, true, "division by zero"},
		{"negative division", -10, 2, -5, false, ""},
		{"zero divided", 0, 5, 0, false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := calc.Divide(tt.a, tt.b)

			if tt.expectErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				if err.Error() != tt.errMsg {
					t.Errorf("expected error %q, got %q", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if result != tt.expected {
					t.Errorf("Divide(%d, %d) = %d; want %d", tt.a, tt.b, result, tt.expected)
				}
			}
		})
	}
}

// TestStringUtilsTableDriven - تست توابع رشته با table-driven
func TestStringUtilsTableDriven(t *testing.T) {
	utils := &StringUtils{}

	t.Run("Reverse", func(t *testing.T) {
		tests := []struct {
			name     string
			input    string
			expected string
		}{
			{"normal string", "hello", "olleh"},
			{"palindrome", "racecar", "racecar"},
			{"unicode", "Hello 世界", "界世 olleH"},
			{"empty", "", ""},
			{"single char", "a", "a"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := utils.Reverse(tt.input)
				if result != tt.expected {
					t.Errorf("Reverse(%q) = %q; want %q", tt.input, result, tt.expected)
				}
			})
		}
	})

	t.Run("IsPalindrome", func(t *testing.T) {
		tests := []struct {
			name     string
			input    string
			expected bool
		}{
			{"palindrome word", "racecar", true},
			{"palindrome phrase", "A man a plan a canal panama", true},
			{"not palindrome", "hello", false},
			{"empty", "", true},
			{"single char", "a", true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := utils.IsPalindrome(tt.input)
				if result != tt.expected {
					t.Errorf("IsPalindrome(%q) = %v; want %v", tt.input, result, tt.expected)
				}
			})
		}
	})

	t.Run("CountVowels", func(t *testing.T) {
		tests := []struct {
			name     string
			input    string
			expected int
		}{
			{"only vowels", "aeiou", 5},
			{"mixed", "hello world", 3},
			{"uppercase", "HELLO", 2},
			{"no vowels", "xyz", 0},
			{"empty", "", 0},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := utils.CountVowels(tt.input)
				if result != tt.expected {
					t.Errorf("CountVowels(%q) = %d; want %d", tt.input, result, tt.expected)
				}
			})
		}
	})
}

// ============================================================================
// بخش 4: Benchmarking (سنجش عملکرد)
// ============================================================================

// BenchmarkAdd - بنچمارک تابع Add
func BenchmarkAdd(b *testing.B) {
	calc := &Calculator{}

	// ResetTimer برای نادیده گرفتن کارهای setup
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		calc.Add(i, i+1)
	}
}

// BenchmarkAddParallel - بنچمارک موازی
func BenchmarkAddParallel(b *testing.B) {
	calc := &Calculator{}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			calc.Add(10, 20)
		}
	})
}

// BenchmarkReverse - بنچمارک با رشته‌های مختلف
func BenchmarkReverse(b *testing.B) {
	utils := &StringUtils{}

	benchmarks := []struct {
		name  string
		input string
	}{
		{"short", "hello"},
		{"medium", "hello world from go"},
		{"long", strings.Repeat("a", 1000)},
		{"unicode", "Hello 世界 123"},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				utils.Reverse(bm.input)
			}
		})
	}
}

// BenchmarkIsPalindrome - بنچمارک با داده‌های مختلف
func BenchmarkIsPalindrome(b *testing.B) {
	utils := &StringUtils{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		utils.IsPalindrome("A man a plan a canal panama")
	}
}

// BenchmarkCountVowels - بنچمارک با حافظه
func BenchmarkCountVowels(b *testing.B) {
	utils := &StringUtils{}
	input := "The quick brown fox jumps over the lazy dog"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		utils.CountVowels(input)
	}
}

// ============================================================================
// بخش 5: تست با Testify (Assert/Require)
// ============================================================================

// TestCalculatorWithTestify - تست با testify/assert
func TestCalculatorWithTestify(t *testing.T) {
	calc := &Calculator{}

	// assert: خطا را ثبت می‌کند ولی ادامه می‌دهد
	assert.Equal(t, 5, calc.Add(2, 3), "Add(2,3) should equal 5")
	assert.Equal(t, 7, calc.Subtract(10, 3), "Subtract(10,3) should equal 7")
	assert.Equal(t, 6, calc.Multiply(2, 3), "Multiply(2,3) should equal 6")

	// تست با require: اگر خطا داشت، تست را متوقف می‌کند
	require.NotNil(t, calc, "Calculator should not be nil")

	// تست divide با خطا
	result, err := calc.Divide(10, 2)
	assert.NoError(t, err, "Divide(10,2) should not return error")
	assert.Equal(t, 5, result)

	_, err = calc.Divide(10, 0)
	assert.Error(t, err, "Divide(10,0) should return error")
	assert.EqualError(t, err, "division by zero", "error message should match")
}

// TestStringUtilsWithTestify - تست testify با struct
func TestStringUtilsWithTestify(t *testing.T) {
	utils := &StringUtils{}

	// تست Reverse
	assert.Equal(t, "olleh", utils.Reverse("hello"))
	assert.Equal(t, "界世 olleH", utils.Reverse("Hello 世界"))
	assert.Empty(t, utils.Reverse(""))

	// تست IsPalindrome
	assert.True(t, utils.IsPalindrome("racecar"))
	assert.True(t, utils.IsPalindrome("A man a plan a canal panama"))
	assert.False(t, utils.IsPalindrome("hello"))

	// تست CountVowels
	assert.Equal(t, 3, utils.CountVowels("hello world"))
	assert.Equal(t, 0, utils.CountVowels("xyz"))
	assert.Equal(t, 5, utils.CountVowels("aeiou"))
}

// TestTableDrivenWithTestify - table-driven با testify
func TestTableDrivenWithTestify(t *testing.T) {
	calc := &Calculator{}

	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"add positive", 2, 3, 5},
		{"add negative", -2, -3, -5},
		{"add zero", 0, 5, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calc.Add(tt.a, tt.b)
			assert.Equal(t, tt.expected, result,
				"Add(%d, %d) should equal %d", tt.a, tt.b, tt.expected)
		})
	}
}

// ============================================================================
// بخش 6: Test Suite با Testify
// ============================================================================

// CalculatorTestSuite - تست‌های گروهی با suite
type CalculatorTestSuite struct {
	suite.Suite
	calc *Calculator
}

// SetupSuite - قبل از همه تست‌ها یک بار اجرا می‌شود
func (s *CalculatorTestSuite) SetupSuite() {
	s.calc = &Calculator{}
}

// TearDownSuite - بعد از همه تست‌ها یک بار اجرا می‌شود
func (s *CalculatorTestSuite) TearDownSuite() {
	// cleanup code
}

// SetupTest - قبل از هر تست اجرا می‌شود
func (s *CalculatorTestSuite) SetupTest() {
	// setup before each test
}

// TearDownTest - بعد از هر تست اجرا می‌شود
func (s *CalculatorTestSuite) TearDownTest() {
	// cleanup after each test
}

func (s *CalculatorTestSuite) TestAdd() {
	s.Equal(5, s.calc.Add(2, 3))
	s.Equal(0, s.calc.Add(-2, 2))
}

func (s *CalculatorTestSuite) TestSubtract() {
	s.Equal(7, s.calc.Subtract(10, 3))
	s.Equal(-7, s.calc.Subtract(3, 10))
}

func (s *CalculatorTestSuite) TestDivide() {
	result, err := s.calc.Divide(10, 2)
	s.NoError(err)
	s.Equal(5, result)

	_, err = s.calc.Divide(10, 0)
	s.Error(err)
	s.Equal("division by zero", err.Error())
}

// اجرای تست‌های suite
// در main نیست، توسط go test اجرا می‌شود:
// func TestCalculatorSuite(t *testing.T) {
//     suite.Run(t, new(CalculatorTestSuite))
// }

// ============================================================================
// بخش 7: Mocking با Testify/Mock
// ============================================================================

// MockUserRepository - پیاده‌سازی mock از UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetByID(id int) (*User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockUserRepository) Save(user *User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

// TestUserServiceWithMock - تست سرویس با mock
func TestUserServiceWithMock(t *testing.T) {
	// ایجاد mock
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	t.Run("GetUser success", func(t *testing.T) {
		expectedUser := &User{ID: 1, Name: "Ali", Email: "ali@test.com"}

		// تنظیم mock: وقتی GetByID(1) صدا زده شد، user مورد نظر را برگردان
		mockRepo.On("GetByID", 1).Return(expectedUser, nil)

		// اجرای تست
		user, err := service.GetUser(1)

		// بررسی‌ها
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)

		// اطمینان از اینکه متد mock صدا زده شده
		mockRepo.AssertCalled(t, "GetByID", 1)
		mockRepo.AssertExpectations(t)
	})

	t.Run("GetUser not found", func(t *testing.T) {
		// تنظیم mock: خطای not found برگردان
		mockRepo.On("GetByID", 999).Return(nil, errors.New("user not found"))

		user, err := service.GetUser(999)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, "user not found", err.Error())

		mockRepo.AssertCalled(t, "GetByID", 999)
	})

	t.Run("GetUser invalid id", func(t *testing.T) {
		user, err := service.GetUser(-1)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, "invalid user id", err.Error())

		// اطمینان از اینکه mock صدا زده نشده
		mockRepo.AssertNotCalled(t, "GetByID", -1)
	})

	t.Run("CreateUser success", func(t *testing.T) {
		// تنظیم mock برای Save (هر Userای قبول کن)
		mockRepo.On("Save", mock.AnythingOfType("*main.User")).Return(nil)

		user, err := service.CreateUser("Sara", "sara@test.com")

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "Sara", user.Name)
		assert.Equal(t, "sara@test.com", user.Email)

		mockRepo.AssertCalled(t, "Save", mock.AnythingOfType("*main.User"))
	})

	t.Run("CreateUser validation error", func(t *testing.T) {
		user, err := service.CreateUser("", "")

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, "name and email are required", err.Error())

		// اطمینان از اینکه Save صدا زده نشده
		mockRepo.AssertNotCalled(t, "Save", mock.Anything)
	})
}

// ============================================================================
// بخش 8: Mocking با Interface-based Design
// ============================================================================

// مثال از طراحی مبتنی بر interface برای قابلیت تست

// Logger interface - برای تست می‌توانیم mock کنیم
type Logger interface {
	Info(msg string)
	Error(msg string)
}

// واقعی
type RealLogger struct{}

func (l *RealLogger) Info(msg string)  { fmt.Println("[INFO]", msg) }
func (l *RealLogger) Error(msg string) { fmt.Println("[ERROR]", msg) }

// Mock برای تست
type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Info(msg string) {
	m.Called(msg)
}

func (m *MockLogger) Error(msg string) {
	m.Called(msg)
}

// سرویس با وابستگی به Logger
type PaymentService struct {
	logger Logger
}

func NewPaymentService(logger Logger) *PaymentService {
	return &PaymentService{logger: logger}
}

func (p *PaymentService) ProcessPayment(amount float64) error {
	if amount <= 0 {
		p.logger.Error("Invalid payment amount")
		return errors.New("amount must be positive")
	}
	p.logger.Info(fmt.Sprintf("Processing payment of $%.2f", amount))
	return nil
}

// تست PaymentService با mock logger
func TestPaymentServiceWithMockLogger(t *testing.T) {
	mockLogger := new(MockLogger)
	service := NewPaymentService(mockLogger)

	t.Run("successful payment", func(t *testing.T) {
		mockLogger.On("Info", "Processing payment of $100.00").Return()

		err := service.ProcessPayment(100)

		assert.NoError(t, err)
		mockLogger.AssertCalled(t, "Info", "Processing payment of $100.00")
	})

	t.Run("invalid payment", func(t *testing.T) {
		mockLogger.On("Error", "Invalid payment amount").Return()

		err := service.ProcessPayment(-10)

		assert.Error(t, err)
		mockLogger.AssertCalled(t, "Error", "Invalid payment amount")
	})
}

// ============================================================================
// بخش 9: Helper Functions در تست
// ============================================================================

// Helper function برای ایجاد تست کاربر
func createTestUser(name, email string) *User {
	return &User{
		ID:    123,
		Name:  name,
		Email: email,
	}
}

// TestWithHelper - استفاده از helper function
func TestWithHelper(t *testing.T) {
	user := createTestUser("Ali", "ali@test.com")

	assert.Equal(t, 123, user.ID)
	assert.Equal(t, "Ali", user.Name)
	assert.Equal(t, "ali@test.com", user.Email)
}

// Must関数 برای تست‌ها (panic اگر خطا داشت)
func mustDivide(a, b int) int {
	result, err := (&Calculator{}).Divide(a, b)
	if err != nil {
		panic(err)
	}
	return result
}

func TestMustDivide(t *testing.T) {
	assert.Equal(t, 5, mustDivide(10, 2))

	// این خط panic می‌کند - در تست باید recover کنیم
	defer func() {
		if r := recover(); r != nil {
			assert.Equal(t, "division by zero", r)
		}
	}()
	mustDivide(10, 0)
}

// ============================================================================
// بخش 10: تست با Skip و Parallel
// ============================================================================

func TestSkipExample(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping long test in short mode")
	}

	// تست طولانی
	time.Sleep(1 * time.Second)
	assert.True(t, true)
}

func TestParallelExample(t *testing.T) {
	t.Parallel() // این تست موازی اجرا می‌شود

	time.Sleep(100 * time.Millisecond)
	assert.True(t, true)
}

func TestParallelSubTests(t *testing.T) {
	tests := []struct {
		name string
		a, b int
	}{
		{"case1", 1, 2},
		{"case2", 3, 4},
		{"case3", 5, 6},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel() // هر sub-test موازی اجرا می‌شود
			result := (&Calculator{}).Add(tt.a, tt.b)
			assert.Equal(t, tt.a+tt.b, result)
		})
	}
}

// ============================================================================
// بخش 11: تست فایل‌های مثال (Example Functions)
// ============================================================================

// ExampleAdd - مثال برای documentation
func ExampleAdd() {
	calc := &Calculator{}
	result := calc.Add(2, 3)
	fmt.Println(result)
	// Output: 5
}

// ExampleReverse - مثال با چند خط خروجی
func ExampleReverse() {
	utils := &StringUtils{}
	result := utils.Reverse("hello")
	fmt.Println(result)
	// Output: olleh
}

// ============================================================================
// بخش 12: جمع‌بندی و نکات
// ============================================================================

/*
نکات مهم برای تست در Go:

1. Table-Driven Tests:
   - بهترین روش برای تست موارد مختلف
   - از struct slice برای تعریف تست‌ها استفاده کن
   - از t.Run برای sub-tests استفاده کن

2. Benchmarking:
   - نام تابع باید با Benchmark شروع شود
   - از b.N برای تعداد تکرار استفاده کن
   - از b.ResetTimer() برای نادیده گرفتن setup
   - از b.RunParallel برای تست موازی

3. Testify:
   - assert: خطا ثبت می‌شود ولی ادامه می‌دهد
   - require: در صورت خطا، تست متوقف می‌شود
   - suite: برای گروه‌بندی تست‌ها
   - mock: برای شبیه‌سازی وابستگی‌ها

4. Mocking:
   - همیشه از interfaceها استفاده کن
   - Mock فقط وابستگی‌های خارجی (DB, API, etc.)
   - از testify/mock برای ایجاد mock استفاده کن

5. Commandes مفید:
   go test                    - اجرای همه تست‌ها
   go test -v                 - با خروجی verbose
   go test -run TestName      - فقط تست خاص
   go test -bench=.           - اجرای benchmarkها
   go test -benchmem          - نمایش آمار حافظه
   go test -cover             - نمایش پوشش تست
   go test -race              - بررسی race conditions
   go test -short             -跳过 تست‌های طولانی
*/

// اجرای تست‌ها (این توابع در فایل *_test.go قرار می‌گیرند)
// در اینجا فقط برای نمایش هستند

// برای اجرای تست‌های suite (در فایل واقعی)
// func TestCalculatorSuite(t *testing.T) {
//     suite.Run(t, new(CalculatorTestSuite))
// }
