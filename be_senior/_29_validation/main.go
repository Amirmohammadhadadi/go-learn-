// ============================================================================
// FILE: validation_guide.go
// TITLE: راهنمای کامل Validation با go-playground/validator
// HOW TO RUN: go run validation_guide.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - Validation چیست و چرا نیاز است؟
// ============================================================================
//
// Validation فرآیند بررسی صحت و اعتبار داده‌های ورودی است.
//
// چرا Validation مهم است؟
// 1. امنیت: جلوگیری از حملات (SQL injection, XSS, etc.)
// 2. یکپارچگی داده: ذخیره داده‌های صحیح در دیتابیس
// 3. تجربه کاربری: نمایش خطاهای مفید به کاربر
// 4. کاهش خطاها: شناسایی زودهنگام مشکلات
//
// کتابخانه go-playground/validator:
// - محبوب‌ترین کتابخانه validation در Go
// - پشتیبانی از تگ‌های متنوع
// - قابلیت سفارشی‌سازی بالا
// - عملکرد بالا (با استفاده از caching)
//
// قانون طلایی:
// "همیشه داده‌های ورودی را اعتبارسنجی کن. هرگز به داده‌های کاربر اعتماد نکن.
//  از validator استفاده کن تا کد تمیز و قابل نگهداری داشته باشی."
// ============================================================================

package __validation

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

// ============================================================================
// بخش 1: نصب و راه‌اندازی اولیه
// ============================================================================

/*
نصب:
go get github.com/go-playground/validator/v10
*/

// ایجاد validator
var validate = validator.New()

// ============================================================================
// بخش 2: مدل‌های نمونه با تگ‌های Validation
// ============================================================================

// User مدل کاربر با validation tags
type User struct {
	ID        int       `json:"id" validate:"omitempty,min=1"`
	Name      string    `json:"name" validate:"required,min=3,max=50"`
	Email     string    `json:"email" validate:"required,email"`
	Age       int       `json:"age" validate:"required,min=0,max=150"`
	Password  string    `json:"password" validate:"required,min=8,containsany=!@#$%^&*"`
	Phone     string    `json:"phone" validate:"omitempty,len=11,numeric"`
	Website   string    `json:"website" validate:"omitempty,url"`
	CreatedAt time.Time `json:"created_at" validate:"omitempty"`
	Role      string    `json:"role" validate:"omitempty,oneof=admin user moderator"`
	Tags      []string  `json:"tags" validate:"omitempty,dive,min=2"`
	Address   Address   `json:"address" validate:"omitempty"`
}

// Address آدرس با validation
type Address struct {
	Street  string `json:"street" validate:"required"`
	City    string `json:"city" validate:"required"`
	Country string `json:"country" validate:"required"`
	ZipCode string `json:"zip_code" validate:"required,len=10,numeric"`
}

// Product مدل محصول
type Product struct {
	ID          int     `json:"id" validate:"omitempty"`
	Name        string  `json:"name" validate:"required,min=3,max=100"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Quantity    int     `json:"quantity" validate:"required,gte=0"`
	Category    string  `json:"category" validate:"required,oneof=electronics clothing books food"`
	Description string  `json:"description" validate:"omitempty,max=500"`
	Discount    float64 `json:"discount" validate:"omitempty,min=0,max=100"`
}

// LoginRequest درخواست لاگین
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// RegisterRequest درخواست ثبت‌نام
type RegisterRequest struct {
	Username        string `json:"username" validate:"required,min=3,max=20,alphanum"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
	TermsAccepted   bool   `json:"terms_accepted" validate:"required"`
}

// SearchRequest درخواست جستجو
type SearchRequest struct {
	Query   string                 `json:"query" validate:"required,min=2,max=100"`
	Page    int                    `json:"page" validate:"omitempty,min=1"`
	Limit   int                    `json:"limit" validate:"omitempty,min=1,max=100"`
	SortBy  string                 `json:"sort_by" validate:"omitempty,oneof=name price date"`
	Order   string                 `json:"order" validate:"omitempty,oneof=asc desc"`
	Filters map[string]interface{} `json:"filters" validate:"omitempty"`
}

// ============================================================================
// بخش 3: Validation تگ‌های پایه
// ============================================================================

func demonstrateBasicValidation() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📝 BASIC VALIDATION TAGS")
	fmt.Println(strings.Repeat("=", 80))

	// مثال validation ساده
	user := User{
		Name:  "A",       // خیلی کوتاه
		Email: "invalid", // نامعتبر
		Age:   200,       // بیشتر از حد
	}

	err := validate.Struct(user)
	if err != nil {
		fmt.Println("Validation errors:")
		for _, err := range err.(validator.ValidationErrors) {
			fmt.Printf("  Field: %s, Tag: %s, Value: %v\n",
				err.Field(), err.Tag(), err.Value())
		}
	}
}

// ============================================================================
// بخش 4: Validation تگ‌های پیشرفته
// ============================================================================

func demonstrateAdvancedValidation() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🔧 ADVANCED VALIDATION TAGS")
	fmt.Println(strings.Repeat("=", 80))

	// ثبت تگ‌های سفارشی
	validate.RegisterValidation("custompassword", func(fl validator.FieldLevel) bool {
		password := fl.Field().String()
		// حداقل یک حرف بزرگ، یک حرف کوچک، یک عدد
		hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
		hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
		hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
		return hasUpper && hasLower && hasNumber
	})

	// مثال با تگ سفارشی
	type CustomUser struct {
		Password string `validate:"required,custompassword"`
	}

	customUser := CustomUser{Password: "weak"}
	err := validate.Struct(customUser)
	if err != nil {
		fmt.Println("Custom validation error:", err)
	}

	// مثال با struct level validation
	type DateRange struct {
		StartDate time.Time `validate:"required"`
		EndDate   time.Time `validate:"required,gtfield=StartDate"`
	}

	dateRange := DateRange{
		StartDate: time.Now(),
		EndDate:   time.Now().Add(-24 * time.Hour), // قبل از StartDate
	}

	err = validate.Struct(dateRange)
	if err != nil {
		fmt.Println("Date range error:", err)
	}
}

// ============================================================================
// بخش 5: Custom Validation Functions
// ============================================================================

// 5.1 تابع validation سفارشی برای شماره ملی ایران
func iranianNationalCodeValidation(fl validator.FieldLevel) bool {
	code := fl.Field().String()
	if len(code) != 10 {
		return false
	}
	// الگوریتم ساده validation کد ملی
	// (در عمل باید الگوریتم کامل پیاده‌سازی شود)
	return true
}

// 5.2 تابع validation برای شماره تلفن ایران
func iranianPhoneValidation(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	match, _ := regexp.MatchString(`^09[0-9]{9}$`, phone)
	return match
}

// 5.3 تابع validation برای کد پستی ایران
func iranianPostalCodeValidation(fl validator.FieldLevel) bool {
	code := fl.Field().String()
	match, _ := regexp.MatchString(`^\d{5}-\d{5}$|^\d{10}$`, code)
	return match
}

// 5.4 تابع validation برای سن
func ageValidation(fl validator.FieldLevel) bool {
	age := fl.Field().Int()
	return age >= 0 && age <= 150
}

// 5.5 تابع validation برای URL سفارشی
func customURLValidation(fl validator.FieldLevel) bool {
	url := fl.Field().String()
	match, _ := regexp.MatchString(`^https?://[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`, url)
	return match
}

func setupCustomValidations() {
	// ثبت validation‌های سفارشی
	validate.RegisterValidation("irannationalcode", iranianNationalCodeValidation)
	validate.RegisterValidation("iranphone", iranianPhoneValidation)
	validate.RegisterValidation("iranpostalcode", iranianPostalCodeValidation)
	validate.RegisterValidation("validage", ageValidation)
	validate.RegisterValidation("customurl", customURLValidation)

	fmt.Println("Custom validations registered successfully!")
}

// ============================================================================
// بخش 6: Struct Level Validation (اعتبارسنجی بین فیلدها)
// ============================================================================

// Order سفارش با validation بین فیلدها
type Order struct {
	Items         []OrderItem `json:"items" validate:"required,min=1,dive"`
	Total         float64     `json:"total" validate:"required,gt=0"`
	Discount      float64     `json:"discount" validate:"min=0,max=100"`
	FinalTotal    float64     `json:"final_total"`
	PaymentMethod string      `json:"payment_method" validate:"required,oneof=credit_card cash bank_transfer"`
}

type OrderItem struct {
	ProductID int     `json:"product_id" validate:"required,gt=0"`
	Quantity  int     `json:"quantity" validate:"required,gt=0"`
	Price     float64 `json:"price" validate:"required,gt=0"`
}

// Struct level validation برای Order
func (o *Order) ValidateStruct() error {
	// محاسبه total و بررسی
	var calculatedTotal float64
	for _, item := range o.Items {
		calculatedTotal += float64(item.Quantity) * item.Price
	}

	if calculatedTotal != o.Total {
		return fmt.Errorf("total mismatch: calculated %.2f, provided %.2f", calculatedTotal, o.Total)
	}

	// محاسبه final total
	o.FinalTotal = o.Total * (1 - o.Discount/100)

	return nil
}

// ============================================================================
// بخش 7: Error Handling و فرمت خطاها
// ============================================================================

// ValidationError ساختار خطای سفارشی
type ValidationError struct {
	Field   string      `json:"field"`
	Message string      `json:"message"`
	Value   interface{} `json:"value,omitempty"`
}

// FormatValidationErrors تبدیل خطاهای validator به فرمت خوانا
func FormatValidationErrors(err error) []ValidationError {
	var errors []ValidationError

	if err == nil {
		return errors
	}

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, ve := range validationErrors {
			errorMsg := getErrorMessage(ve)
			errors = append(errors, ValidationError{
				Field:   ve.Field(),
				Message: errorMsg,
				Value:   ve.Value(),
			})
		}
	}

	return errors
}

// getErrorMessage تولید پیام خطای خوانا
func getErrorMessage(ve validator.FieldError) string {
	switch ve.Tag() {
	case "required":
		return "This field is required"
	case "min":
		return fmt.Sprintf("Minimum length is %s", ve.Param())
	case "max":
		return fmt.Sprintf("Maximum length is %s", ve.Param())
	case "email":
		return "Invalid email format"
	case "url":
		return "Invalid URL format"
	case "oneof":
		return fmt.Sprintf("Must be one of: %s", ve.Param())
	case "gt":
		return fmt.Sprintf("Must be greater than %s", ve.Param())
	case "gte":
		return fmt.Sprintf("Must be greater than or equal to %s", ve.Param())
	case "lt":
		return fmt.Sprintf("Must be less than %s", ve.Param())
	case "lte":
		return fmt.Sprintf("Must be less than or equal to %s", ve.Param())
	case "eqfield":
		return "Fields do not match"
	case "containsany":
		return fmt.Sprintf("Must contain at least one of: %s", ve.Param())
	case "alphanum":
		return "Must contain only alphanumeric characters"
	case "numeric":
		return "Must contain only numbers"
	case "len":
		return fmt.Sprintf("Must be exactly %s characters", ve.Param())
	case "iranphone":
		return "Invalid Iranian phone number (format: 09123456789)"
	case "iranpostalcode":
		return "Invalid postal code (format: 1234567890 or 12345-67890)"
	case "validage":
		return "Age must be between 0 and 150"
	default:
		return fmt.Sprintf("Validation failed on %s tag", ve.Tag())
	}
}

// ============================================================================
// بخش 8: Validation در HTTP Handlers (net/http)
// ============================================================================

// validateAndDecode تابع کمکی برای validation
func validateAndDecode(r *http.Request, target interface{}) []ValidationError {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(target); err != nil {
		return []ValidationError{
			{Field: "body", Message: "Invalid JSON: " + err.Error()},
		}
	}

	if err := validate.Struct(target); err != nil {
		return FormatValidationErrors(err)
	}

	return nil
}

// createUserHandler هندلر ایجاد کاربر با validation
func createUserHandler(w http.ResponseWriter, r *http.Request) {
	var user User

	// Validation
	if errors := validateAndDecode(r, &user); len(errors) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "error",
			"errors": errors,
		})
		return
	}

	// داده معتبر است
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data":   user,
	})
}

// registerHandler هندلر ثبت‌نام
func registerHandler(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest

	if errors := validateAndDecode(r, &req); len(errors) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "error",
			"errors": errors,
		})
		return
	}

	// ثبت‌نام موفق
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "User registered successfully",
	})
}

// createProductHandler هندلر ایجاد محصول
func createProductHandler(w http.ResponseWriter, r *http.Request) {
	var product Product

	if errors := validateAndDecode(r, &product); len(errors) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "error",
			"errors": errors,
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data":   product,
	})
}

// ============================================================================
// بخش 9: Validation در Chi Router
// ============================================================================

/*
// Chi Router با validation
func chiValidationRouter() {
	r := chi.NewRouter()

	// Middleware برای validation
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// می‌توان validation عمومی اضافه کرد
			next.ServeHTTP(w, r)
		})
	})

	r.Post("/users", func(w http.ResponseWriter, r *http.Request) {
		var user User
		if errors := validateAndDecode(r, &user); len(errors) > 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"errors": errors,
			})
			return
		}
		// process user
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
	})

	r.Post("/register", func(w http.ResponseWriter, r *http.Request) {
		var req RegisterRequest
		if errors := validateAndDecode(r, &req); len(errors) > 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"errors": errors,
			})
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"status": "registered"})
	})

	http.ListenAndServe(":8080", r)
}
*/

// ============================================================================
// بخش 10: Validation در Gin Framework
// ============================================================================

/*
// Gin Router با validation یکپارچه
func ginValidationRouter() {
	r := gin.Default()

	// 1. Basic validation
	r.POST("/users", func(c *gin.Context) {
		var user User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Validation با validator
		if err := validate.Struct(user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"errors": FormatValidationErrors(err),
			})
			return
		}

		c.JSON(http.StatusCreated, user)
	})

	// 2. Register with validation
	r.POST("/register", func(c *gin.Context) {
		var req RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := validate.Struct(req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"errors": FormatValidationErrors(err),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"status": "registered"})
	})

	// 3. Query parameter validation
	r.GET("/search", func(c *gin.Context) {
		var req SearchRequest
		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := validate.Struct(req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"errors": FormatValidationErrors(err),
			})
			return
		}

		c.JSON(http.StatusOK, req)
	})

	// 4. Custom validator in Gin
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("custom", func(fl validator.FieldLevel) bool {
			return fl.Field().String() == "custom"
		})
	}

	r.Run(":8080")
}
*/

// ============================================================================
// بخش 11: Validation با Query Parameters
// ============================================================================

func validateQueryParams(r *http.Request, target interface{}) []ValidationError {
	query := r.URL.Query()

	// تبدیل query params به map برای validation
	// (در عمل باید از یک کتابخانه مثل go-querystring استفاده کرد)

	if err := validate.Struct(target); err != nil {
		return FormatValidationErrors(err)
	}

	return nil
}

type SearchParams struct {
	Q     string `validate:"required,min=2,max=100"`
	Page  int    `validate:"omitempty,min=1"`
	Limit int    `validate:"omitempty,min=1,max=100"`
	Sort  string `validate:"omitempty,oneof=asc desc"`
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	params := SearchParams{
		Q:     r.URL.Query().Get("q"),
		Page:  atoi(r.URL.Query().Get("page")),
		Limit: atoi(r.URL.Query().Get("limit")),
		Sort:  r.URL.Query().Get("sort"),
	}

	if errors := validateQueryParams(r, &params); len(errors) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "error",
			"errors": errors,
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(params)
}

// تابع کمکی برای تبدیل string به int
func atoi(s string) int {
	var i int
	fmt.Sscanf(s, "%d", &i)
	return i
}

// ============================================================================
// بخش 12: Validation Middleware
// ============================================================================

// ValidationMiddleware میدلور عمومی برای validation
type ValidationMiddleware struct {
	validate *validator.Validate
}

func NewValidationMiddleware() *ValidationMiddleware {
	return &ValidationMiddleware{
		validate: validator.New(),
	}
}

func (vm *ValidationMiddleware) ValidateBody(target interface{}) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if err := json.NewDecoder(r.Body).Decode(target); err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON"})
				return
			}

			if err := vm.validate.Struct(target); err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"errors": FormatValidationErrors(err),
				})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// ============================================================================
// بخش 13: Best Practices و Tips
// ============================================================================

func validationBestPractices() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("💡 VALIDATION BEST PRACTICES")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────┐
│ 1. Use Omitempty for Optional Fields                          │
│    ` + "`validate:\"omitempty,min=3\"`" + `                         │
│                                                                 │
│ 2. Validate on Both Client and Server                         │
│    - Client: User experience                                   │
│    - Server: Security                                          │
│                                                                 │
│ 3. Return User-Friendly Error Messages                        │
│    ❌ "Field validation for 'Email' failed on the 'email' tag" │
│    ✅ "Email format is invalid"                                │
│                                                                 │
│ 4. Use Struct Level Validation for Cross-Field Rules          │
│    - Confirm password matching                                 │
│    - Date range validation                                     │
│                                                                 │
│ 5. Create Custom Validators for Business Logic                │
│    - Iranian national code                                     │
│    - Product availability                                      │
│                                                                 │
│ 6. Validate Early (Before Database Operations)                │
│                                                                 │
│ 7. Use Different Validation Rules for Different Operations    │
│    - Create vs Update                                          │
│    - Admin vs User                                             │
│                                                                 │
│ 8. Cache Validation Results When Possible                     │
│                                                                 │
│ 9. Log Validation Failures for Security Monitoring            │
│                                                                 │
│ 10. Use Context for Dynamic Validation Rules                  │
└─────────────────────────────────────────────────────────────────┘

🎯 COMMON VALIDATION TAGS:

   required     - Field must not be zero value
   omitempty    - Skip validation if field is zero
   min          - Minimum value (numbers) or length (strings)
   max          - Maximum value (numbers) or length (strings)
   len          - Exact length
   eq           - Equal to value
   ne           - Not equal to value
   gt           - Greater than
   gte          - Greater than or equal
   lt           - Less than
   lte          - Less than or equal
   oneof        - One of the values
   email        - Valid email format
   url          - Valid URL format
   alphanum     - Alphanumeric characters only
   numeric      - Numeric characters only
   contains     - Contains substring
   containsany  - Contains any of characters
   excluded     - Excluded value
   eqfield      - Equal to another field
   nefield      - Not equal to another field
   gtfield      - Greater than another field
   dive         - Validate nested struct/slice
`)
}

// ============================================================================
// بخش 14: سرور کامل
// ============================================================================

func runValidationServer() {
	// ثبت validation‌های سفارشی
	setupCustomValidations()

	mux := http.NewServeMux()

	// Endpoints
	mux.HandleFunc("/users", createUserHandler)
	mux.HandleFunc("/register", registerHandler)
	mux.HandleFunc("/products", createProductHandler)
	mux.HandleFunc("/search", searchHandler)

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Println("Validation server starting on :8080")
	log.Println("Test endpoints:")
	log.Println("  POST /users     - Create user with validation")
	log.Println("  POST /register  - Register user")
	log.Println("  POST /products  - Create product")
	log.Println("  GET  /search?q=... - Search with query validation")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

// ============================================================================
// بخش 15: main
// ============================================================================

func main() {
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("🎯 GO-PLAYGROUND/VALIDATOR GUIDE")
	fmt.Println("Complete validation for Go applications")
	fmt.Println(strings.Repeat("=", 80))

	// نمایش validation tags پایه
	demonstrateBasicValidation()

	// نمایش validation پیشرفته
	demonstrateAdvancedValidation()

	// Best practices
	validationBestPractices()

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🚀 Starting validation server on :8080")
	fmt.Println(strings.Repeat("=", 80))

	runValidationServer()
}

// تابع کمکی برای تکرار رشته
func stringsRepeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}

/*
خلاصه Validation Tags مهم
دسته	تگ	توضیح
الزام	required	فیلد الزامی است
omitempty	اگر مقدار صفر باشد، validation نادیده گرفته می‌شود
طول/مقدار	min	حداقل مقدار/طول
max	حداکثر مقدار/طول
len	طول دقیق
مقایسه	gt	بزرگتر از
gte	بزرگتر یا مساوی
lt	کوچکتر از
lte	کوچکتر یا مساوی
eq	مساوی با
ne	نامساوی با
دسته‌بندی	oneof	یکی از مقادیر مجاز
فرمت	email	ایمیل معتبر
url	آدرس اینترنتی معتبر
alphanum	فقط حروف و اعداد
numeric	فقط اعداد
مقایسه فیلدها	eqfield	مساوی با فیلد دیگر
nefield	نامساوی با فیلد دیگر
gtfield	بزرگتر از فیلد دیگر
تودرتو	dive	اعمال validation روی اعضای slice/struct

*/
