// ============================================================================
// FILE: api_versioning_guide.go
// TITLE: راهنمای کامل نسخه‌بندی API - URL, Header, Query Parameter
// HOW TO RUN: go run api_versioning_guide.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - چرا نسخه‌بندی API نیاز است؟
// ============================================================================
//
// نسخه‌بندی API دلایل مهمی دارد:
// 1. تغییرات backward-incompatible (تغییر فرمت پاسخ، حذف فیلد، تغییر نوع داده)
// 2. اضافه کردن قابلیت‌های جدید بدون شکستن کلاینت‌های قدیمی
// 3. نگهداری همزمان چند نسخه برای کلاینت‌های مختلف
// 4. تست و استقرار تدریجی تغییرات
// 5. منسوخ کردن تدریجی (deprecation) APIهای قدیمی
//
// روش‌های اصلی نسخه‌بندی:
// 1. URL Path (مثل /v1/users, /v2/users) - رایج‌ترین روش
// 2. Query Parameter (مثل /users?version=1) - ساده اما کمتر استفاده می‌شود
// 3. Custom Header (مثل Accept: application/vnd.myapi.v1+json) - RESTful
// 4. Content Negotiation (Accept header با vendor type)
//
// قانون طلایی:
// "از URL Path برای نسخه‌بندی استفاده کن (ساده و واضح).
//  تغییرات backward-compatible را بدون افزایش نسخه انجام بده.
//  نسخه‌های قدیمی را با هشدار deprecation نگه دار و بعد از مدتی حذف کن."
// ============================================================================

package __api_versioning

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// ============================================================================
// بخش 1: مدل‌های داده با نسخه‌های مختلف
// ============================================================================

// UserV1 نسخه اول کاربر (فیلدهای ساده)
type UserV1 struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UserV2 نسخه دوم کاربر (اضافه شدن فیلدهای جدید)
type UserV2 struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone,omitempty"`
	Avatar    string    `json:"avatar,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserV3 نسخه سوم (تغییر ساختار)
type UserV3 struct {
	ID        int       `json:"id"`
	FullName  string    `json:"full_name"` // تغییر از Name به FullName
	Email     string    `json:"email"`
	Contact   Contact   `json:"contact"` // فیلد ترکیبی جدید
	Metadata  Metadata  `json:"metadata"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Contact struct {
	Phone string `json:"phone"`
	Email string `json:"email"`
}

type Metadata struct {
	Avatar   string   `json:"avatar"`
	Roles    []string `json:"roles"`
	LastSeen string   `json:"last_seen,omitempty"`
}

// ProductV1 نسخه اول محصول
type ProductV1 struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// ProductV2 نسخه دوم (تغییر نوع قیمت)
type ProductV2 struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Price    string `json:"price"`    // تبدیل به string برای دقت بیشتر
	Currency string `json:"currency"` // افزودن واحد پول
	InStock  bool   `json:"in_stock"`
}

// ============================================================================
// بخش 2: Service Layer
// ============================================================================

type UserService struct {
	usersV1 map[int]UserV1
	usersV2 map[int]UserV2
	usersV3 map[int]UserV3
}

func NewUserService() *UserService {
	return &UserService{
		usersV1: map[int]UserV1{
			1: {ID: 1, Name: "Ali", Email: "ali@example.com"},
			2: {ID: 2, Name: "Sara", Email: "sara@example.com"},
		},
		usersV2: map[int]UserV2{
			1: {ID: 1, Name: "Ali", Email: "ali@example.com", Phone: "+989123456789", CreatedAt: time.Now(), UpdatedAt: time.Now()},
			2: {ID: 2, Name: "Sara", Email: "sara@example.com", Phone: "+989987654321", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		},
		usersV3: map[int]UserV3{
			1: {ID: 1, FullName: "Ali Rezaei", Email: "ali@example.com", Contact: Contact{Phone: "+989123456789", Email: "ali@example.com"}, Metadata: Metadata{Avatar: "/avatars/ali.jpg", Roles: []string{"admin", "user"}}, CreatedAt: time.Now(), UpdatedAt: time.Now()},
			2: {ID: 2, FullName: "Sara Mohammadi", Email: "sara@example.com", Contact: Contact{Phone: "+989987654321", Email: "sara@example.com"}, Metadata: Metadata{Avatar: "/avatars/sara.jpg", Roles: []string{"user"}}, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		},
	}
}

func (s *UserService) GetUserV1(id int) (*UserV1, error) {
	if user, ok := s.usersV1[id]; ok {
		return &user, nil
	}
	return nil, fmt.Errorf("user not found")
}

func (s *UserService) GetUserV2(id int) (*UserV2, error) {
	if user, ok := s.usersV2[id]; ok {
		return &user, nil
	}
	return nil, fmt.Errorf("user not found")
}

func (s *UserService) GetUserV3(id int) (*UserV3, error) {
	if user, ok := s.usersV3[id]; ok {
		return &user, nil
	}
	return nil, fmt.Errorf("user not found")
}

func (s *UserService) CreateUserV3(user UserV3) (*UserV3, error) {
	user.ID = len(s.usersV3) + 1
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	s.usersV3[user.ID] = user
	return &user, nil
}

// ============================================================================
// بخش 3: Method 1 - URL Path Versioning (رایج‌ترین)
// ============================================================================

// URL Path: /api/v1/users, /api/v2/users, /api/v3/users

func handleV1Users(w http.ResponseWriter, r *http.Request) {
	userService := NewUserService()

	// GET /api/v1/users/1
	if strings.HasPrefix(r.URL.Path, "/api/v1/users/") {
		idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/users/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "INVALID_ID", "Invalid user ID")
			return
		}

		user, err := userService.GetUserV1(id)
		if err != nil {
			writeJSONError(w, http.StatusNotFound, "NOT_FOUND", "User not found")
			return
		}

		writeJSONResponse(w, http.StatusOK, user)
		return
	}

	// GET /api/v1/users
	if r.URL.Path == "/api/v1/users" {
		var users []UserV1
		for _, u := range userService.usersV1 {
			users = append(users, u)
		}
		writeJSONResponse(w, http.StatusOK, users)
		return
	}

	writeJSONError(w, http.StatusNotFound, "NOT_FOUND", "Endpoint not found")
}

func handleV2Users(w http.ResponseWriter, r *http.Request) {
	userService := NewUserService()

	// GET /api/v2/users/1
	if strings.HasPrefix(r.URL.Path, "/api/v2/users/") {
		idStr := strings.TrimPrefix(r.URL.Path, "/api/v2/users/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "INVALID_ID", "Invalid user ID")
			return
		}

		user, err := userService.GetUserV2(id)
		if err != nil {
			writeJSONError(w, http.StatusNotFound, "NOT_FOUND", "User not found")
			return
		}

		writeJSONResponse(w, http.StatusOK, user)
		return
	}

	// GET /api/v2/users
	if r.URL.Path == "/api/v2/users" {
		var users []UserV2
		for _, u := range userService.usersV2 {
			users = append(users, u)
		}
		writeJSONResponse(w, http.StatusOK, users)
		return
	}

	writeJSONError(w, http.StatusNotFound, "NOT_FOUND", "Endpoint not found")
}

func handleV3Users(w http.ResponseWriter, r *http.Request) {
	userService := NewUserService()

	// GET /api/v3/users/1
	if strings.HasPrefix(r.URL.Path, "/api/v3/users/") && r.Method == http.MethodGet {
		idStr := strings.TrimPrefix(r.URL.Path, "/api/v3/users/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "INVALID_ID", "Invalid user ID")
			return
		}

		user, err := userService.GetUserV3(id)
		if err != nil {
			writeJSONError(w, http.StatusNotFound, "NOT_FOUND", "User not found")
			return
		}

		writeJSONResponse(w, http.StatusOK, user)
		return
	}

	// GET /api/v3/users
	if r.URL.Path == "/api/v3/users" && r.Method == http.MethodGet {
		var users []UserV3
		for _, u := range userService.usersV3 {
			users = append(users, u)
		}
		writeJSONResponse(w, http.StatusOK, users)
		return
	}

	// POST /api/v3/users
	if r.URL.Path == "/api/v3/users" && r.Method == http.MethodPost {
		var user UserV3
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			writeJSONError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid JSON")
			return
		}

		created, err := userService.CreateUserV3(user)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "CREATE_ERROR", err.Error())
			return
		}

		writeJSONResponse(w, http.StatusCreated, created)
		return
	}

	writeJSONError(w, http.StatusNotFound, "NOT_FOUND", "Endpoint not found")
}

// ============================================================================
// بخش 4: Method 2 - Query Parameter Versioning
// ============================================================================

// Query Parameter: /api/users?version=1, /api/users?version=2

func handleUsersWithQueryVersion(w http.ResponseWriter, r *http.Request) {
	version := r.URL.Query().Get("version")
	if version == "" {
		version = "1" // default version
	}

	userService := NewUserService()

	// GET /api/users/1?version=1
	if strings.HasPrefix(r.URL.Path, "/api/users/") {
		idStr := strings.TrimPrefix(r.URL.Path, "/api/users/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "INVALID_ID", "Invalid user ID")
			return
		}

		switch version {
		case "1":
			user, err := userService.GetUserV1(id)
			if err != nil {
				writeJSONError(w, http.StatusNotFound, "NOT_FOUND", "User not found")
				return
			}
			writeJSONResponse(w, http.StatusOK, user)
		case "2":
			user, err := userService.GetUserV2(id)
			if err != nil {
				writeJSONError(w, http.StatusNotFound, "NOT_FOUND", "User not found")
				return
			}
			writeJSONResponse(w, http.StatusOK, user)
		case "3":
			user, err := userService.GetUserV3(id)
			if err != nil {
				writeJSONError(w, http.StatusNotFound, "NOT_FOUND", "User not found")
				return
			}
			writeJSONResponse(w, http.StatusOK, user)
		default:
			writeJSONError(w, http.StatusBadRequest, "UNSUPPORTED_VERSION", "Version not supported")
		}
		return
	}

	writeJSONError(w, http.StatusNotFound, "NOT_FOUND", "Endpoint not found")
}

// ============================================================================
// بخش 5: Method 3 - Custom Header Versioning
// ============================================================================

// Custom Header: Accept: application/vnd.myapi.v1+json

const (
	HeaderAccept = "Accept"
)

func parseVersionFromHeader(r *http.Request) (string, error) {
	accept := r.Header.Get(HeaderAccept)

	// مثال: application/vnd.myapi.v1+json
	if strings.Contains(accept, "vnd.myapi.v") {
		// استخراج نسخه
		parts := strings.Split(accept, "vnd.myapi.v")
		if len(parts) > 1 {
			versionPart := strings.Split(parts[1], "+")[0]
			return versionPart, nil
		}
	}

	return "1", nil // default version
}

func handleUsersWithHeaderVersion(w http.ResponseWriter, r *http.Request) {
	version, err := parseVersionFromHeader(r)
	if err != nil {
		version = "1"
	}

	userService := NewUserService()

	// GET /api/users/1
	if strings.HasPrefix(r.URL.Path, "/api/users/") {
		idStr := strings.TrimPrefix(r.URL.Path, "/api/users/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "INVALID_ID", "Invalid user ID")
			return
		}

		switch version {
		case "1":
			user, err := userService.GetUserV1(id)
			if err != nil {
				writeJSONError(w, http.StatusNotFound, "NOT_FOUND", "User not found")
				return
			}
			w.Header().Set("Content-Type", "application/vnd.myapi.v1+json")
			writeJSONResponse(w, http.StatusOK, user)
		case "2":
			user, err := userService.GetUserV2(id)
			if err != nil {
				writeJSONError(w, http.StatusNotFound, "NOT_FOUND", "User not found")
				return
			}
			w.Header().Set("Content-Type", "application/vnd.myapi.v2+json")
			writeJSONResponse(w, http.StatusOK, user)
		case "3":
			user, err := userService.GetUserV3(id)
			if err != nil {
				writeJSONError(w, http.StatusNotFound, "NOT_FOUND", "User not found")
				return
			}
			w.Header().Set("Content-Type", "application/vnd.myapi.v3+json")
			writeJSONResponse(w, http.StatusOK, user)
		default:
			writeJSONError(w, http.StatusBadRequest, "UNSUPPORTED_VERSION", "Version not supported")
		}
		return
	}

	writeJSONError(w, http.StatusNotFound, "NOT_FOUND", "Endpoint not found")
}

// ============================================================================
// بخش 6: Router با پشتیبانی از چند نسخه
// ============================================================================

func versionedRouter() http.Handler {
	mux := http.NewServeMux()

	// Method 1: URL Path Versioning
	mux.HandleFunc("/api/v1/users", handleV1Users)
	mux.HandleFunc("/api/v1/users/", handleV1Users)
	mux.HandleFunc("/api/v2/users", handleV2Users)
	mux.HandleFunc("/api/v2/users/", handleV2Users)
	mux.HandleFunc("/api/v3/users", handleV3Users)
	mux.HandleFunc("/api/v3/users/", handleV3Users)

	// Method 2: Query Parameter Versioning
	mux.HandleFunc("/api/users", handleUsersWithQueryVersion)
	mux.HandleFunc("/api/users/", handleUsersWithQueryVersion)

	// Method 3: Custom Header Versioning
	mux.HandleFunc("/api/header/users", handleUsersWithHeaderVersion)
	mux.HandleFunc("/api/header/users/", handleUsersWithHeaderVersion)

	// Documentation endpoint
	mux.HandleFunc("/docs", handleDocumentation)

	return mux
}

// ============================================================================
// بخش 7: Deprecation Warning
// ============================================================================

// DeprecationWarning افزودن هدر هشدار برای نسخه‌های قدیمی
func DeprecationWarning(version string, deprecationDate string, sunsetDate string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// افزودن هدر هشدار
			w.Header().Set("Deprecation", "true")
			w.Header().Set("Deprecation-Date", deprecationDate)
			w.Header().Set("Sunset", sunsetDate)
			w.Header().Set("Link", "</docs/migration>; rel=\"deprecation\"; type=\"text/html\"")
			w.Header().Set("Warning", `299 - "API version `+version+` is deprecated. Please upgrade to v3."`)

			next.ServeHTTP(w, r)
		})
	}
}

// ============================================================================
// بخش 8: Chi Router Versioning (مثال)
// ============================================================================

/*
// Chi Router با نسخه‌بندی
func chiVersionedRouter() {
	r := chi.NewRouter()

	// Middleware برای لاگ کردن نسخه
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// استخراج نسخه از URL
			var version string
			if strings.HasPrefix(r.URL.Path, "/api/v1") {
				version = "v1"
			} else if strings.HasPrefix(r.URL.Path, "/api/v2") {
				version = "v2"
			}
			ctx := context.WithValue(r.Context(), "api_version", version)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	// نسخه v1
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/users", handleV1Users)
		r.Get("/users/{id}", handleV1Users)
	})

	// نسخه v2
	r.Route("/api/v2", func(r chi.Router) {
		// افزودن هشدار deprecation
		r.Use(DeprecationWarning("v2", "2024-01-01", "2024-06-01"))
		r.Get("/users", handleV2Users)
		r.Get("/users/{id}", handleV2Users)
	})

	// نسخه v3 (جدید)
	r.Route("/api/v3", func(r chi.Router) {
		r.Get("/users", handleV3Users)
		r.Get("/users/{id}", handleV3Users)
		r.Post("/users", handleV3Users)
	})

	http.ListenAndServe(":8080", r)
}
*/

// ============================================================================
// بخش 9: Gin Versioning (مثال)
// ============================================================================

/*
// Gin Router با نسخه‌بندی
func ginVersionedRouter() {
	r := gin.Default()

	// Group by version
	v1 := r.Group("/api/v1")
	{
		v1.GET("/users", func(c *gin.Context) {
			// logic
		})
		v1.GET("/users/:id", func(c *gin.Context) {
			// logic
		})
	}

	v2 := r.Group("/api/v2")
	{
		// Deprecation warning
		v2.Use(func(c *gin.Context) {
			c.Header("Deprecation", "true")
			c.Header("Sunset", "2024-06-01")
			c.Next()
		})
		v2.GET("/users", func(c *gin.Context) {
			// logic
		})
	}

	v3 := r.Group("/api/v3")
	{
		v3.GET("/users", func(c *gin.Context) {
			// logic
		})
		v3.POST("/users", func(c *gin.Context) {
			// logic
		})
	}

	r.Run(":8080")
}
*/

// ============================================================================
// بخش 10: Documentation Endpoint
// ============================================================================

func handleDocumentation(w http.ResponseWriter, r *http.Request) {
	docs := map[string]interface{}{
		"versions": []map[string]interface{}{
			{
				"version":  "v1",
				"status":   "stable",
				"base_url": "/api/v1",
				"endpoints": []string{
					"GET /users",
					"GET /users/{id}",
				},
				"deprecated": false,
			},
			{
				"version":  "v2",
				"status":   "deprecated",
				"base_url": "/api/v2",
				"endpoints": []string{
					"GET /users",
					"GET /users/{id}",
				},
				"deprecated":       true,
				"deprecation_date": "2024-01-01",
				"sunset_date":      "2024-06-01",
				"migration_url":    "/docs/migration-v2-to-v3",
			},
			{
				"version":  "v3",
				"status":   "latest",
				"base_url": "/api/v3",
				"endpoints": []string{
					"GET /users",
					"GET /users/{id}",
					"POST /users",
				},
				"deprecated": false,
			},
		},
		"migration_guide": map[string]string{
			"v1_to_v2": "/docs/migration/v1-to-v2",
			"v2_to_v3": "/docs/migration/v2-to-v3",
		},
		"changelog": "/docs/changelog",
	}

	writeJSONResponse(w, http.StatusOK, docs)
}

// ============================================================================
// بخش 11: Helper Functions
// ============================================================================

func writeJSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeJSONError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": map[string]string{
			"code":    code,
			"message": message,
		},
	})
}

// ============================================================================
// بخش 12: Best Practices
// ============================================================================

func versioningBestPractices() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("💡 API VERSIONING BEST PRACTICES")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────┐
│ 1. Choose a Strategy and Stick to It                          │
│    - URL Path is most common and easiest to understand        │
│    - Use same strategy across all endpoints                    │
│                                                                 │
│ 2. Keep Backward Compatibility                                 │
│    - Adding new fields is safe (backward-compatible)           │
│    - Removing fields requires new version                      │
│    - Changing field types requires new version                 │
│                                                                 │
│ 3. Version at the API Level, Not Resource Level               │
│    ✅ /api/v1/users, /api/v1/products                          │
│    ❌ /api/users/v1, /api/products/v1                          │
│                                                                 │
│ 4. Communicate Deprecation Clearly                             │
│    - Add Deprecation header                                    │
│    - Add Sunset header                                         │
│    - Provide migration guide                                   │
│    - Give足够的时间 (minimum 6 months)                           │
│                                                                 │
│ 5. Maintain Old Versions for a Period                          │
│    - Support至少 2 versions at a time                           │
│    - Document sunset date                                      │
│    - Monitor usage before removing                             │
│                                                                 │
│ 6. Use Semantic Versioning                                     │
│    - Major version for breaking changes                        │
│    - Minor version for new features (backward-compatible)      │
│    - Patch version for bug fixes                               │
│                                                                 │
│ 7. Document Everything                                         │
│    - Changelog                                                 │
│    - Migration guides                                          │
│    - Deprecation schedule                                      │
│    - Example requests/responses                                │
│                                                                 │
│ 8. Test All Versions                                           │
│    - Integration tests for each version                        │
│    - Ensure old versions still work                            │
│                                                                 │
│ 9. Monitor Version Usage                                       │
│    - Track which versions are used                             │
│    - Plan deprecation based on usage                           │
│                                                                 │
│ 10. Use API Gateway for Version Routing                        │
│     - Can route different versions to different services       │
│     - Easier to deprecate old versions                         │
└─────────────────────────────────────────────────────────────────┘
`)
}

// ============================================================================
// بخش 13: Comparison Table
// ============================================================================

func versioningComparison() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📊 VERSIONING METHODS COMPARISON")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ METHOD              │ PROS                    │ CONS                         │
├─────────────────────┼─────────────────────────┼──────────────────────────────┤
│ URL Path            │ • Simple and clear      │ • Clutters URLs              │
│ (/v1/resource)      │ • Easy to cache         │ • Not RESTful (some argue)   │
│                     │ • Easy to route         │                              │
│                     │ • Bookmarkable          │                              │
├─────────────────────┼─────────────────────────┼──────────────────────────────┤
│ Query Parameter     │ • Simple to implement   │ • Not obvious                │
│ (?version=1)        │ • Doesn't change URL    │ • Can be forgotten           │
│                     │ • Easy to default       │ • Harder to cache            │
├─────────────────────┼─────────────────────────┼──────────────────────────────┤
│ Custom Header       │ • Clean URLs            │ • Not visible in browser     │
│ (Accept: vnd.api.v1)│ • RESTful               │ • Harder to test             │
│                     │ • Content negotiation   │ • Requires special client    │
├─────────────────────┼─────────────────────────┼──────────────────────────────┤
│ RECOMMENDATION: Use URL Path for simplicity and wide compatibility.         │
└─────────────────────────────────────────────────────────────────────────────┘
`)
}

// ============================================================================
// بخش 14: main
// ============================================================================

func main() {
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("🎯 API VERSIONING GUIDE")
	fmt.Println("URL Path | Query Parameter | Custom Header")
	fmt.Println(strings.Repeat("=", 80))

	versioningBestPractices()
	versioningComparison()

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🚀 Starting versioned API server on :8080")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println("\n📌 Available Endpoints:")
	fmt.Println("  Method 1 - URL Path:")
	fmt.Println("    GET  /api/v1/users         - List users (v1)")
	fmt.Println("    GET  /api/v1/users/1       - Get user (v1)")
	fmt.Println("    GET  /api/v2/users         - List users (v2)")
	fmt.Println("    GET  /api/v3/users         - List users (v3)")
	fmt.Println("    POST /api/v3/users         - Create user (v3)")
	fmt.Println()
	fmt.Println("  Method 2 - Query Parameter:")
	fmt.Println("    GET  /api/users?version=1  - Get users (v1)")
	fmt.Println("    GET  /api/users/1?version=2 - Get user (v2)")
	fmt.Println()
	fmt.Println("  Method 3 - Custom Header:")
	fmt.Println("    GET  /api/header/users/1   - With Accept: application/vnd.myapi.v1+json")
	fmt.Println()
	fmt.Println("  Documentation:")
	fmt.Println("    GET  /docs                  - API documentation")

	log.Fatal(http.ListenAndServe(":8080", versionedRouter()))
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
خلاصه روش‌های نسخه‌بندی
روش	مثال	مزایا	معایب
URL Path	/api/v1/users	ساده، واضح، کش شدن آسان	URL طولانی‌تر
Query Parameter	/api/users?version=1	ساده، تغییر نمی‌دهد URL	نامشخص، فراموش شدن
Custom Header	Accept: vnd.api.v1+json	تمیز، RESTful	نامرئی در مرورگر
Deprecation Headers
Header	مثال	توضیح
Deprecation	true	نشان می‌دهد API منسوخ شده
Deprecation-Date	2024-01-01	تاریخ شروع deprecation
Sunset	2024-06-01	تاریخ حذف نهایی
Link	</docs/migration>; rel="deprecation"	لینک راهنمای مهاجرت
Warning	299 - "API version v2 is deprecated"	پیام هشدار

*/
