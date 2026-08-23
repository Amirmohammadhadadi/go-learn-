// ============================================================================
// FILE: response_guide.go
// TITLE: راهنمای کامل Response استاندارد و Error Handling در لایه HTTP
// HOW TO RUN: go run response_guide.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - چرا Response استاندارد نیاز داریم؟
// ============================================================================
//
// یک Response استاندارد مزایای زیر را دارد:
// 1. قابلیت پیش‌بینی: کلاینت می‌داند چه فرمتی دریافت می‌کند
// 2. یکپارچگی: تمام APIها از یک فرمت پیروی می‌کنند
// 3. خطایابی آسان: خطاها ساختار یکسانی دارند
// 4. مستندسازی ساده: یک فرمت برای همه endpoints
// 5. کاهش کد تکراری: یک تابع برای ارسال پاسخ
//
// قانون طلایی:
// "همیشه از یک ساختار استاندارد برای پاسخ‌ها استفاده کن.
//  خطاها را با کد HTTP مناسب و پیام مفید برگردان.
//  هرگز جزئیات داخلی سرور را در خطاها لو نده."
// ============================================================================

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
	"time"
)

// ============================================================================
// بخش 1: ساختار Response استاندارد
// ============================================================================

// Response ساختار استاندارد پاسخ API
type Response struct {
	Success   bool        `json:"success"`              // موفقیت آمیز بودن درخواست
	Status    int         `json:"status"`               // کد وضعیت HTTP
	Message   string      `json:"message,omitempty"`    // پیام عمومی
	Data      interface{} `json:"data,omitempty"`       // داده‌های پاسخ
	Error     *APIError   `json:"error,omitempty"`      // اطلاعات خطا (در صورت وجود)
	Timestamp time.Time   `json:"timestamp"`            // زمان پاسخ
	Path      string      `json:"path,omitempty"`       // مسیر درخواست
	RequestID string      `json:"request_id,omitempty"` // شناسه درخواست (برای tracing)
}

// APIError ساختار خطا
type APIError struct {
	Code     string            `json:"code"`               // کد خطا (مثل "VALIDATION_ERROR")
	Message  string            `json:"message"`            // پیام خطا (قابل نمایش به کاربر)
	Details  map[string]string `json:"details,omitempty"`  // جزئیات بیشتر (مثل validation errors)
	Internal string            `json:"internal,omitempty"` // خطای داخلی (فقط برای debugging)
}

// ============================================================================
// بخش 2: Response Helper Functions
// ============================================================================

// ResponseWriter ساختار کمکی برای نوشتن پاسخ
type ResponseWriter struct {
	http.ResponseWriter
	requestID string
}

// NewResponseWriter ایجاد ResponseWriter جدید
func NewResponseWriter(w http.ResponseWriter, r *http.Request) *ResponseWriter {
	requestID := r.Header.Get("X-Request-ID")
	if requestID == "" {
		requestID = fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return &ResponseWriter{
		ResponseWriter: w,
		requestID:      requestID,
	}
}

// JSON ارسال پاسخ JSON
func (rw *ResponseWriter) JSON(status int, data interface{}) {
	rw.Header().Set("Content-Type", "application/json")
	rw.Header().Set("X-Request-ID", rw.requestID)
	rw.WriteHeader(status)

	// اگر data از نوع Response است، timestamp و path را اضافه کن
	if resp, ok := data.(Response); ok {
		resp.Timestamp = time.Now()
		resp.Path = "" // در صورت نیاز مسیر را اضافه کنید
		resp.RequestID = rw.requestID
		json.NewEncoder(rw.ResponseWriter).Encode(resp)
		return
	}

	json.NewEncoder(rw.ResponseWriter).Encode(data)
}

// Success ارسال پاسخ موفق
func (rw *ResponseWriter) Success(status int, message string, data interface{}) {
	rw.JSON(status, Response{
		Success: true,
		Status:  status,
		Message: message,
		Data:    data,
	})
}

// Error ارسال پاسخ خطا
func (rw *ResponseWriter) Error(status int, code, message string, details map[string]string) {
	rw.JSON(status, Response{
		Success: false,
		Status:  status,
		Error: &APIError{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}

// ErrorWithInternal ارسال پاسخ خطا با خطای داخلی (برای debugging)
func (rw *ResponseWriter) ErrorWithInternal(status int, code, message string, internalErr error, details map[string]string) {
	err := &APIError{
		Code:     code,
		Message:  message,
		Details:  details,
		Internal: internalErr.Error(),
	}

	rw.JSON(status, Response{
		Success: false,
		Status:  status,
		Error:   err,
	})
}

// ValidationError ارسال خطای اعتبارسنجی
func (rw *ResponseWriter) ValidationError(errors map[string]string) {
	rw.Error(http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", errors)
}

// NotFound ارسال خطای 404
func (rw *ResponseWriter) NotFound(resource string) {
	rw.Error(http.StatusNotFound, "NOT_FOUND", fmt.Sprintf("%s not found", resource), nil)
}

// Unauthorized ارسال خطای 401
func (rw *ResponseWriter) Unauthorized(message string) {
	if message == "" {
		message = "Authentication required"
	}
	rw.Error(http.StatusUnauthorized, "UNAUTHORIZED", message, nil)
}

// Forbidden ارسال خطای 403
func (rw *ResponseWriter) Forbidden(message string) {
	if message == "" {
		message = "Access denied"
	}
	rw.Error(http.StatusForbidden, "FORBIDDEN", message, nil)
}

// InternalError ارسال خطای 500
func (rw *ResponseWriter) InternalError(err error) {
	rw.ErrorWithInternal(http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error", err, nil)
}

// ============================================================================
// بخش 3: Error Types (خطاهای سفارشی)
// ============================================================================

// AppError خطای سفارشی برنامه
type AppError struct {
	Code    string
	Message string
	Status  int
	Details map[string]string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s - %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// خطاهای از پیش تعریف شده
var (
	ErrNotFound        = &AppError{Code: "NOT_FOUND", Message: "Resource not found", Status: http.StatusNotFound}
	ErrUnauthorized    = &AppError{Code: "UNAUTHORIZED", Message: "Authentication required", Status: http.StatusUnauthorized}
	ErrForbidden       = &AppError{Code: "FORBIDDEN", Message: "Access denied", Status: http.StatusForbidden}
	ErrBadRequest      = &AppError{Code: "BAD_REQUEST", Message: "Invalid request", Status: http.StatusBadRequest}
	ErrValidation      = &AppError{Code: "VALIDATION_ERROR", Message: "Validation failed", Status: http.StatusBadRequest}
	ErrConflict        = &AppError{Code: "CONFLICT", Message: "Resource already exists", Status: http.StatusConflict}
	ErrInternal        = &AppError{Code: "INTERNAL_ERROR", Message: "Internal server error", Status: http.StatusInternalServerError}
	ErrTooManyRequests = &AppError{Code: "TOO_MANY_REQUESTS", Message: "Rate limit exceeded", Status: http.StatusTooManyRequests}
)

// ============================================================================
// بخش 4: Error Handler Middleware
// ============================================================================

// ErrorHandlerMiddleware میدلور مدیریت خطا
func ErrorHandlerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// لاگ کردن panic
				log.Printf("PANIC: %v\n%s", err, debug.Stack())

				rw := NewResponseWriter(w, r)
				rw.Error(http.StatusInternalServerError, "INTERNAL_ERROR", "An unexpected error occurred", nil)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// ============================================================================
// بخش 5: Example Service و Handlers
// ============================================================================

// UserService سرویس نمونه
type UserService struct {
	users map[int]User
}

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func NewUserService() *UserService {
	return &UserService{
		users: map[int]User{
			1: {ID: 1, Name: "Ali", Email: "ali@example.com"},
			2: {ID: 2, Name: "Sara", Email: "sara@example.com"},
		},
	}
}

// GetUser دریافت کاربر
func (s *UserService) GetUser(id int) (*User, error) {
	if id <= 0 {
		return nil, ErrBadRequest
	}

	user, exists := s.users[id]
	if !exists {
		return nil, ErrNotFound
	}

	return &user, nil
}

// CreateUser ایجاد کاربر
func (s *UserService) CreateUser(user User) (*User, error) {
	if user.Name == "" {
		return nil, &AppError{
			Code:    "VALIDATION_ERROR",
			Message: "Name is required",
			Status:  http.StatusBadRequest,
			Details: map[string]string{"name": "Name cannot be empty"},
		}
	}

	if user.Email == "" {
		return nil, &AppError{
			Code:    "VALIDATION_ERROR",
			Message: "Email is required",
			Status:  http.StatusBadRequest,
			Details: map[string]string{"email": "Email cannot be empty"},
		}
	}

	// بررسی ایمیل تکراری
	for _, existing := range s.users {
		if existing.Email == user.Email {
			return nil, ErrConflict
		}
	}

	user.ID = len(s.users) + 1
	s.users[user.ID] = user

	return &user, nil
}

// ============================================================================
// بخش 6: HTTP Handlers با Response استاندارد
// ============================================================================

var userService = NewUserService()

// GetUserHandler هندلر دریافت کاربر
func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	rw := NewResponseWriter(w, r)

	// استخراج ID از URL (در اینجا ساده فرض شده)
	idStr := strings.TrimPrefix(r.URL.Path, "/users/")
	var id int
	fmt.Sscanf(idStr, "%d", &id)

	if id == 0 {
		rw.Error(http.StatusBadRequest, "INVALID_ID", "Invalid user ID", nil)
		return
	}

	user, err := userService.GetUser(id)
	if err != nil {
		switch e := err.(type) {
		case *AppError:
			rw.Error(e.Status, e.Code, e.Message, e.Details)
		default:
			rw.Error(http.StatusInternalServerError, "UNKNOWN_ERROR", err.Error(), nil)
		}
		return
	}

	rw.Success(http.StatusOK, "User retrieved successfully", user)
}

// CreateUserHandler هندلر ایجاد کاربر
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	rw := NewResponseWriter(w, r)

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		rw.Error(http.StatusBadRequest, "INVALID_JSON", "Invalid JSON format", nil)
		return
	}

	created, err := userService.CreateUser(user)
	if err != nil {
		switch e := err.(type) {
		case *AppError:
			rw.Error(e.Status, e.Code, e.Message, e.Details)
		default:
			rw.Error(http.StatusInternalServerError, "UNKNOWN_ERROR", err.Error(), nil)
		}
		return
	}

	rw.Success(http.StatusCreated, "User created successfully", created)
}

// ============================================================================
// بخش 7: Chi Router با Response استاندارد
// ============================================================================

/*
// Chi Router با Response استاندارد
func chiRouter() {
	r := chi.NewRouter()

	// Middlewareها
	r.Use(ErrorHandlerMiddleware)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Routes
	r.Get("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		rw := NewResponseWriter(w, r)

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			rw.Error(http.StatusBadRequest, "INVALID_ID", "Invalid user ID", nil)
			return
		}

		user, err := userService.GetUser(id)
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				rw.NotFound("User")
				return
			}
			rw.InternalError(err)
			return
		}

		rw.Success(http.StatusOK, "User found", user)
	})

	r.Post("/users", func(w http.ResponseWriter, r *http.Request) {
		rw := NewResponseWriter(w, r)

		var user User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			rw.Error(http.StatusBadRequest, "INVALID_JSON", "Invalid JSON", nil)
			return
		}

		created, err := userService.CreateUser(user)
		if err != nil {
			if appErr, ok := err.(*AppError); ok {
				rw.Error(appErr.Status, appErr.Code, appErr.Message, appErr.Details)
				return
			}
			rw.InternalError(err)
			return
		}

		rw.Success(http.StatusCreated, "User created", created)
	})

	http.ListenAndServe(":8080", r)
}
*/

// ============================================================================
// بخش 8: Gin Framework با Response استاندارد
// ============================================================================

/*
// Gin Router با Response استاندارد
func ginRouter() {
	r := gin.Default()

	// Middleware برای تبدیل Response استاندارد
	r.Use(func(c *gin.Context) {
		c.Set("response_writer", NewResponseWriter(c.Writer, c.Request))
		c.Next()
	})

	r.GET("/users/:id", func(c *gin.Context) {
		rw := c.MustGet("response_writer").(*ResponseWriter)

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			rw.Error(http.StatusBadRequest, "INVALID_ID", "Invalid user ID", nil)
			return
		}

		user, err := userService.GetUser(id)
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				rw.NotFound("User")
				return
			}
			rw.InternalError(err)
			return
		}

		rw.Success(http.StatusOK, "User found", user)
	})

	r.POST("/users", func(c *gin.Context) {
		rw := c.MustGet("response_writer").(*ResponseWriter)

		var user User
		if err := c.ShouldBindJSON(&user); err != nil {
			rw.Error(http.StatusBadRequest, "INVALID_JSON", "Invalid JSON", nil)
			return
		}

		created, err := userService.CreateUser(user)
		if err != nil {
			if appErr, ok := err.(*AppError); ok {
				rw.Error(appErr.Status, appErr.Code, appErr.Message, appErr.Details)
				return
			}
			rw.InternalError(err)
			return
		}

		rw.Success(http.StatusCreated, "User created", created)
	})

	r.Run(":8080")
}
*/

// ============================================================================
// بخش 9: کدهای خطا (Error Codes)
// ============================================================================

func errorCodes() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📋 STANDARD ERROR CODES")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────┐
│ CODE                    │ HTTP STATUS │ DESCRIPTION              │
├─────────────────────────┼─────────────┼──────────────────────────┤
│ BAD_REQUEST             │ 400         │ Invalid request format   │
│ VALIDATION_ERROR        │ 400         │ Validation failed        │
│ INVALID_PARAM           │ 400         │ Invalid parameter        │
│ MISSING_PARAM           │ 400         │ Required parameter missing│
│ INVALID_JSON            │ 400         │ Invalid JSON format      │
│ UNAUTHORIZED            │ 401         │ Authentication required  │
│ FORBIDDEN               │ 403         │ Access denied            │
│ NOT_FOUND               │ 404         │ Resource not found       │
│ METHOD_NOT_ALLOWED      │ 405         │ HTTP method not allowed  │
│ CONFLICT                │ 409         │ Resource conflict        │
│ TOO_MANY_REQUESTS       │ 429         │ Rate limit exceeded      │
│ INTERNAL_ERROR          │ 500         │ Internal server error    │
│ SERVICE_UNAVAILABLE     │ 503         │ Service unavailable      │
└─────────────────────────────────────────────────────────────────┘
`)
}

// ============================================================================
// بخش 10: Best Practices
// ============================================================================

func bestPractices() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("💡 RESPONSE BEST PRACTICES")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────┐
│ 1. Use Consistent Response Format                              │
│    Always use the same structure for all responses             │
│                                                                 │
│ 2. Use Appropriate HTTP Status Codes                           │
│    200 - OK                                                     │
│    201 - Created                                                │
│    400 - Bad Request                                            │
│    401 - Unauthorized                                           │
│    403 - Forbidden                                              │
│    404 - Not Found                                              │
│    409 - Conflict                                               │
│    422 - Unprocessable Entity                                   │
│    429 - Too Many Requests                                      │
│    500 - Internal Server Error                                  │
│                                                                 │
│ 3. Don't Expose Internal Details                               │
│    ❌ "error": "database connection failed"                     │
│    ✅ "error": "internal server error"                          │
│                                                                 │
│ 4. Include Request ID for Tracing                              │
│    Helps debugging across services                              │
│                                                                 │
│ 5. Validate Response Before Sending                            │
│    Check for nil data, empty arrays, etc.                      │
│                                                                 │
│ 6. Set Proper Headers                                          │
│    Content-Type: application/json                              │
│    X-Request-ID: <id>                                          │
│    Cache-Control: no-cache (for dynamic responses)             │
│                                                                 │
│ 7. Document Your Response Format                               │
│    Use OpenAPI/Swagger for API documentation                   │
│                                                                 │
│ 8. Return Meaningful Error Messages                            │
│    ✅ "Email format is invalid"                                 │
│    ❌ "Validation failed"                                       │
│                                                                 │
│ 9. Include Validation Details                                  │
│    {                                                           │
│      "error": {                                                │
│        "code": "VALIDATION_ERROR",                             │
│        "details": {                                            │
│          "email": "Email is required",                         │
│          "age": "Age must be at least 18"                      │
│        }                                                       │
│      }                                                         │
│    }                                                           │
│                                                                 │
│ 10. Log Errors But Don't Return to Client                     │
│     Internal errors should be logged, not exposed              │
└─────────────────────────────────────────────────────────────────┘
`)
}

// ============================================================================
// بخش 11: Testing Response
// ============================================================================

func testResponse() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🧪 TESTING RESPONSE HANDLERS")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
// Example test for response handler
func TestGetUserHandler(t *testing.T) {
    req := httptest.NewRequest("GET", "/users/1", nil)
    rr := httptest.NewRecorder()
    
    GetUserHandler(rr, req)
    
    assert.Equal(t, http.StatusOK, rr.Code)
    
    var resp Response
    json.NewDecoder(rr.Body).Decode(&resp)
    
    assert.True(t, resp.Success)
    assert.NotNil(t, resp.Data)
}

// Example test for error response
func TestGetUserHandlerNotFound(t *testing.T) {
    req := httptest.NewRequest("GET", "/users/999", nil)
    rr := httptest.NewRecorder()
    
    GetUserHandler(rr, req)
    
    assert.Equal(t, http.StatusNotFound, rr.Code)
    
    var resp Response
    json.NewDecoder(rr.Body).Decode(&resp)
    
    assert.False(t, resp.Success)
    assert.Equal(t, "NOT_FOUND", resp.Error.Code)
}
`)
}

// ============================================================================
// بخش 12: سرور کامل
// ============================================================================

func runServer() {
	mux := http.NewServeMux()

	// اعمال ErrorHandlerMiddleware روی همه handlers
	handler := ErrorHandlerMiddleware(mux)

	// Routes
	mux.HandleFunc("/users/", GetUserHandler)
	mux.HandleFunc("/users", CreateUserHandler)

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		rw := NewResponseWriter(w, r)
		rw.Success(http.StatusOK, "Service is healthy", nil)
	})

	// Root
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		rw := NewResponseWriter(w, r)
		rw.Success(http.StatusOK, "Welcome to API", map[string]string{
			"version": "1.0.0",
			"docs":    "/docs",
		})
	})

	log.Println("Server starting on :8080")
	log.Println("Endpoints:")
	log.Println("  GET    /users/{id}  - Get user by ID")
	log.Println("  POST   /users       - Create new user")
	log.Println("  GET    /health      - Health check")
	log.Println("  GET    /            - API info")

	log.Fatal(http.ListenAndServe(":8080", handler))
}

// ============================================================================
// بخش 13: main
// ============================================================================

func main() {
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("🎯 STANDARD RESPONSE & ERROR HANDLING GUIDE")
	fmt.Println(strings.Repeat("=", 80))

	errorCodes()
	bestPractices()
	testResponse()

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🚀 Starting server on :8080")
	fmt.Println(strings.Repeat("=", 80))

	runServer()
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
// Response موفق
{
  "success": true,
  "status": 200,
  "message": "User retrieved successfully",
  "data": { "id": 1, "name": "Ali" },
  "timestamp": "2024-01-15T10:30:00Z",
  "request_id": "1234567890"
}

// Response خطا
{
  "success": false,
  "status": 400,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "details": {
      "email": "Email is required",
      "age": "Age must be at least 18"
    }
  },
  "timestamp": "2024-01-15T10:30:00Z",
  "request_id": "1234567890"
}
*/

/*
کدهای HTTP مهم
کد	معنی	استفاده
200	OK	درخواست موفق
201	Created	منبع جدید ایجاد شد
204	No Content	موفق بدون محتوا
400	Bad Request	درخواست نامعتبر
401	Unauthorized	نیاز به احراز هویت
403	Forbidden	دسترسی غیرمجاز
404	Not Found	منبع یافت نشد
409	Conflict	تداخل با منبع موجود
422	Unprocessable Entity	داده معتبر نیست
429	Too Many Requests	محدودیت نرخ
500	Internal Server Error	خطای سرور
آیا نیاز به توضیح بیشتری در مورد Response استاندارد دارید؟


*/
