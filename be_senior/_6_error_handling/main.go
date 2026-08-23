// ============================================================================
// FILE: error_handling_guide.go
// PURPOSE: راهنمای کامل مدیریت خطا در Go - قابل اجرا و کامنت‌گذاری شده
// HOW TO RUN: go run error_handling_guide.go
// ============================================================================

package _6_error_handling

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

// ============================================================================
// بخش 1: تعریف خطاهای سراسری (Global Errors)
// ============================================================================

// ErrNotFound - زمانی که رکوردی پیدا نمی‌شود
var ErrNotFound = errors.New("record not found")

// ErrInvalidInput - زمانی که ورودی نامعتبر است
var ErrInvalidInput = errors.New("invalid input provided")

// ErrPermissionDenied - زمانی که دسترسی وجود ندارد
var ErrPermissionDenied = errors.New("permission denied")

// ============================================================================
// بخش 2: خطای سفارشی با ساختار (Custom Error Struct)
// ============================================================================

// ValidationError خطای اعتبارسنجی با جزئیات کامل
type ValidationError struct {
	Field   string      // نام فیلد خطادار
	Value   interface{} // مقدار ارسال‌شده
	Message string      // توضیح خطا
	Code    int         // کد خطا (مثلاً 400)
}

// Error پیاده‌سازی اینترفیس error
func (e ValidationError) Error() string {
	return fmt.Sprintf("[%d] validation failed on '%s': %s (value: %v)",
		e.Code, e.Field, e.Message, e.Value)
}

// ============================================================================
// بخش 3: خطای حرفه‌ای با Operation Tracking
// ============================================================================

// AppError خطای استاندارد اپلیکیشن با قابلیت ردیابی عملیات
type AppError struct {
	Op   string // نام عملیات (مثلاً "db.insert", "auth.login")
	Err  error  // خطای اصلی (wrapped error)
	Code int    // کد HTTP یا کد داخلی
}

// Error پیاده‌سازی اینترفیس error
func (e AppError) Error() string {
	if e.Op == "" {
		return e.Err.Error()
	}
	return fmt.Sprintf("%s: %v", e.Op, e.Err)
}

// Unwrap برای compatibility با errors.Is و errors.As
func (e AppError) Unwrap() error {
	return e.Err
}

// ============================================================================
// بخش 4: توابع با مدیریت خطای صحیح (بدون panic معمولی)
// ============================================================================

// divide تقسیم دو عدد با بررسی خطای تقسیم بر صفر
func divide(a, b float64) (float64, error) {
	if b == 0 {
		// برگرداندن خطای توصیفی با fmt.Errorf
		return 0, fmt.Errorf("division by zero: %.2f / %.2f", a, b)
	}
	return a / b, nil
}

// findUser پیدا کردن کاربر با خطای استاندارد ErrNotFound
func findUser(id int) (string, error) {
	if id <= 0 {
		return "", ErrInvalidInput
	}
	if id != 42 {
		// پیچیدن خطا با اطلاعات بیشتر
		return "", fmt.Errorf("findUser: user with id %d: %w", id, ErrNotFound)
	}
	return "Ali Rezaei", nil
}

// validateAge اعتبارسنجی سن با خطای سفارشی ValidationError
func validateAge(age int) error {
	if age < 0 {
		return ValidationError{
			Field:   "age",
			Value:   age,
			Message: "age cannot be negative",
			Code:    400,
		}
	}
	if age < 18 {
		return ValidationError{
			Field:   "age",
			Value:   age,
			Message: "user must be at least 18 years old",
			Code:    400,
		}
	}
	return nil
}

// processOrder یک تابع لایه بالا که خطاهای مختلف را مدیریت می‌کند
func processOrder(orderID int) error {
	// استفاده از الگوی AppError
	if orderID <= 0 {
		return AppError{
			Op:   "processOrder.validate",
			Err:  ErrInvalidInput,
			Code: 400,
		}
	}

	// شبیه‌سازی خطای دیتابیس
	if orderID == 999 {
		return AppError{
			Op:   "processOrder.findOrder",
			Err:  ErrNotFound,
			Code: 404,
		}
	}

	return nil
}

// ============================================================================
// بخش 5: panic و recover (فقط در موارد خاص)
// ============================================================================

// Config خطای panic فقط در زمان راه‌اندازی برنامه مجاز است
type Config struct {
	Port string
}

// LoadConfig بارگذاری کانفیگ - اگر فایل نباشد panic می‌کند
func LoadConfig(path string) Config {
	// در یک برنامه واقعی، فایل را می‌خواند
	// اینجا شبیه‌سازی می‌کنیم
	if path == "" {
		// ✅ تنها جای درست برای panic: مقداردهی اولیه اجباری
		panic("config path cannot be empty - application cannot start")
	}
	return Config{Port: "8080"}
}

// SafeHandler میان‌افزار recovery برای جلوگیری از کرش سرور
// این فقط در لایه بالا (مثل main یا router) استفاده می‌شود
func SafeHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// لاگ کردن خطا
				log.Printf("PANIC RECOVERED: %v", err)
				// برگرداندن خطای 500 به کاربر
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
		}()
		fn(w, r)
	}
}

// riskyHandler یک هندلر که ممکن است panic کند
func riskyHandler(w http.ResponseWriter, r *http.Request) {
	// ❌ این کار غلط است - نباید در منطق عادی panic کرد
	// اما برای تست recovery اینجا می‌گذاریم
	panic("something unexpected happened!")
}

// ============================================================================
// بخش 6: بررسی خطا با errors.Is و errors.As
// ============================================================================

// handleUserRequest بررسی انواع مختلف خطا
func handleUserRequest(userID int) error {
	_, err := findUser(userID)
	if err != nil {
		// استفاده از errors.Is برای خطاهای از پیش تعریف شده
		if errors.Is(err, ErrNotFound) {
			return fmt.Errorf("user not found: please check user ID %d", userID)
		}
		if errors.Is(err, ErrInvalidInput) {
			return fmt.Errorf("invalid user ID: %d", userID)
		}
		return err
	}
	return nil
}

// processWithValidation استفاده از errors.As برای خطاهای سفارشی
func processWithValidation(age int) error {
	err := validateAge(age)
	if err != nil {
		var valErr ValidationError
		// استخراج خطای سفارشی
		if errors.As(err, &valErr) {
			// حالا می‌توانیم به فیلدهای خاص دسترسی داشته باشیم
			log.Printf("Validation failed on field: %s", valErr.Field)
			return fmt.Errorf("input error: %s", valErr.Message)
		}
		return err
	}
	return nil
}

// handleAppError استفاده از AppError با errors.As
func handleAppError(orderID int) error {
	err := processOrder(orderID)
	if err != nil {
		var appErr AppError
		if errors.As(err, &appErr) {
			// دسترسی به کد خطا
			switch appErr.Code {
			case 400:
				return fmt.Errorf("bad request: %s", appErr.Error())
			case 404:
				return fmt.Errorf("not found: %s", appErr.Error())
			default:
				return fmt.Errorf("internal error: %s", appErr.Error())
			}
		}
		return err
	}
	return nil
}

// ============================================================================
// بخش 7: الگوی پیچیدن خطا (Wrapping) و حفظ زنجیره
// ============================================================================

// dbLayer شبیه‌سازی لایه دیتابیس
func dbLayer(id int) error {
	if id == 0 {
		return errors.New("connection timeout")
	}
	return nil
}

// serviceLayer لایه سرویس - خطا را با اطلاعات بیشتر می‌پیچد
func serviceLayer(id int) error {
	err := dbLayer(id)
	if err != nil {
		// استفاده از %w برای حفظ زنجیره خطا
		return fmt.Errorf("serviceLayer: failed to fetch data for id %d: %w", id, err)
	}
	return nil
}

// handlerLayer لایه هندلر - خطا را برای کاربر نهایی آماده می‌کند
func handlerLayer(id int) error {
	err := serviceLayer(id)
	if err != nil {
		// بررسی زنجیره خطا با errors.Is
		if errors.Is(err, errors.New("connection timeout")) {
			return fmt.Errorf("database is unavailable, please try later")
		}
		return fmt.Errorf("internal server error")
	}
	return nil
}

// ============================================================================
// بخش 8: جمع‌بندی و اجرای مثال‌ها
// ============================================================================

func main() {
	fmt.Println("========== ERROR HANDLING GUIDE IN GO ==========\n")

	// مثال 1: تقسیم
	fmt.Println("--- Example 1: Division ---")
	result, err := divide(10, 0)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Result: %.2f\n", result)
	}

	// مثال 2: پیدا کردن کاربر
	fmt.Println("\n--- Example 2: Find User ---")
	user, err := findUser(99)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		if errors.Is(err, ErrNotFound) {
			fmt.Println(">> Action: Show 404 page to user")
		}
	} else {
		fmt.Printf("User found: %s\n", user)
	}

	// مثال 3: اعتبارسنجی سن
	fmt.Println("\n--- Example 3: Age Validation ---")
	err = validateAge(15)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		var valErr ValidationError
		if errors.As(err, &valErr) {
			fmt.Printf(">> Field: %s, Message: %s\n", valErr.Field, valErr.Message)
		}
	}

	// مثال 4: خطای اپلیکیشن
	fmt.Println("\n--- Example 4: App Error ---")
	err = handleAppError(999)
	if err != nil {
		fmt.Printf("App Error: %v\n", err)
	}

	// مثال 5: پیچیدن خطا
	fmt.Println("\n--- Example 5: Error Wrapping ---")
	err = handlerLayer(0)
	if err != nil {
		fmt.Printf("Final error to user: %v\n", err)
	}

	// مثال 6: panic/recovery در سرور (اجرا نمی‌شود، فقط کد نشان داده شده)
	fmt.Println("\n--- Example 6: Panic/Recover Pattern ---")
	fmt.Println("See SafeHandler middleware in code above")
	fmt.Println("✅ Correct: panic only in init/config")
	fmt.Println("❌ Wrong: panic in business logic")

	// قانون طلایی نهایی
	fmt.Println("\n========== GOLDEN RULES ==========")
	fmt.Println("1. Always return error, never panic (except init)")
	fmt.Println("2. Always check if err != nil")
	fmt.Println("3. Use errors.Is for sentinel errors (ErrNotFound)")
	fmt.Println("4. Use errors.As for custom error types")
	fmt.Println("5. Wrap errors with %w to preserve chain")
	fmt.Println("6. Use recover only in top-level goroutines (main, HTTP handlers)")
	fmt.Println("7. Never ignore errors with _")
	fmt.Println("8. Log error only once at the caller level")
}

// ============================================================================
// بخش 9: توابع کمکی که در کد واقعی استفاده می‌شوند
// ============================================================================

// MustParse یک تابع کمکی که panic می‌کند اگر parsing شکست بخورد
// فقط برای استفاده در زمان initialization (مثل flags, env vars)
func MustParse(value string) int {
	var result int
	// شبیه‌سازی parsing
	if value == "" {
		panic("empty value cannot be parsed")
	}
	return result
}

// IgnoreError خطا را نادیده می‌گیرد - فقط برای مواردی که واقعاً خطا不重要 است
// مثال: Close یک فایل فقط برای logging
func IgnoreError(err error) {
	if err != nil {
		log.Printf("ignored error: %v", err)
	}
}
