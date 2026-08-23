// ============================================================================
// FILE: routing_guide.go
// TITLE: راهنمای کامل Routing در Go - Chi, Gorilla/Mux, Gin
// HOW TO RUN: go run routing_guide.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - مقایسه Routing Libraries
// ============================================================================
//
// سه کتابخانه محبوب Routing در Go:
//
// 1. Chi (توصیه شده برای پروژه‌های جدید)
//    - سبک و سریع (100% compatible با net/http)
//    - Middlewareهای built-in عالی
//    - Route grouping و sub-router
//    - URL پارامتر و regex support
//    - محبوبیت بالا و active maintenance
//
// 2. Gorilla/Mux (کلاسیک و پایدار)
//    - قدیمی‌ترین و پایدارترین
//    - Full-featured router
//    - Reverse routing
//    - هنوز هم کار می‌کند ولی maintenance کمتری دارد
//
// 3. Gin (Full-featured framework)
//    - بسیار سریع (fasthttp based)
//    - Built-in validation و binding
//    - Middleware ecosystem بزرگ
//    - Performance-critical applications
//
// قانون طلایی:
// "برای APIهای جدید از Chi استفاده کن (ساده و استاندارد).
//  برای پروژه‌های legacy از Gorilla/Mux استفاده کن.
//  برای performance-critical یا full-featured از Gin استفاده کن."
// ============================================================================

package __routing

// برای اجرا: هر بخش را جداگانه در فایل جداگانه تست کنید
// go run chi_example.go
// go run gorilla_example.go
// go run gin_example.go

import (
	"fmt"
	"strings"
)

// ============================================================================
// بخش 1: Chi Router (توصیه شده)
// ============================================================================

/*
نصب:
go get github.com/go-chi/chi/v5
go get github.com/go-chi/chi/v5/middleware
go get github.com/go-chi/cors
*/

/*
// مثال کامل با Chi
package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var users = map[int]User{
	1: {ID: 1, Name: "Ali", Age: 30},
	2: {ID: 2, Name: "Sara", Age: 25},
}

// 1.1 Basic Chi Router
func basicChiRouter() {
	r := chi.NewRouter()

	// Middlewareها
	r.Use(middleware.Logger)    // لاگ کردن درخواست‌ها
	r.Use(middleware.Recoverer) // بازیابی از panic
	r.Use(middleware.Timeout(60 * time.Second))

	// Routes
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello Chi!"))
	})

	r.Get("/users", getUsers)
	r.Get("/users/{id}", getUserByID)
	r.Post("/users", createUser)
	r.Put("/users/{id}", updateUser)
	r.Delete("/users/{id}", deleteUser)

	http.ListenAndServe(":8080", r)
}

// 1.2 URL Parameters
func getUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	user, exists := users[id]
	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

// 1.3 Query Parameters
func getUsers(w http.ResponseWriter, r *http.Request) {
	// ?page=1&limit=10
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	response := map[string]interface{}{
		"users": users,
		"page":  page,
		"limit": limit,
	}
	json.NewEncoder(w).Encode(response)
}

// 1.4 Request Body (JSON)
func createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	user.ID = len(users) + 1
	users[user.ID] = user

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// 1.5 Route Grouping و Sub-router
func routeGrouping() {
	r := chi.NewRouter()

	// گروه API v1
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/users", getUsers)
		r.Get("/users/{id}", getUserByID)

		// زیرگروه admin
		r.Route("/admin", func(r chi.Router) {
			// Admin-only middleware
			r.Use(adminAuthMiddleware)
			r.Get("/stats", getStats)
			r.Post("/users", createUser)
		})
	})

	// گروه API v2
	r.Route("/api/v2", func(r chi.Router) {
		r.Get("/users", getUsersV2)
	})
}

// 1.6 Custom Middleware
func adminAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")
		if apiKey != "secret-key" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func getStats(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]int{"users": len(users)})
}

func getUsersV2(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]interface{}{
		"version": "v2",
		"data":    users,
	})
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)

	var user User
	json.NewDecoder(r.Body).Decode(&user)
	user.ID = id
	users[id] = user

	json.NewEncoder(w).Encode(user)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)

	delete(users, id)
	w.WriteHeader(http.StatusNoContent)
}

// 1.7 Chi با Regex و Wildcard
func chiWithPatterns() {
	r := chi.NewRouter()

	// Regex در مسیر
	r.Get("/users/{id:[0-9]+}", getUserByID)

	// Wildcard
	r.Get("/files/*", func(w http.ResponseWriter, r *http.Request) {
		path := chi.URLParam(r, "*")
		w.Write([]byte("File path: " + path))
	})

	// Optional parameter
	r.Get("/articles/{slug}-{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "slug")
		id := chi.URLParam(r, "id")
		w.Write([]byte(fmt.Sprintf("Article: %s (ID: %s)", slug, id)))
	})
}

func main() {
	basicChiRouter()
}
*/

// ============================================================================
// بخش 2: Gorilla/Mux Router
// ============================================================================

/*
نصب:
go get github.com/gorilla/mux
go get github.com/gorilla/handlers
*/

/*
// مثال کامل با Gorilla/Mux
package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

var products = map[int]Product{
	1: {ID: 1, Name: "Laptop", Price: 999.99},
	2: {ID: 2, Name: "Mouse", Price: 19.99},
}

// 2.1 Basic Mux Router
func basicMuxRouter() {
	r := mux.NewRouter()

	// Middlewareها
	r.Use(loggingMiddleware)
	r.Use(handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
	))

	// Routes با روش‌های مختلف
	r.HandleFunc("/", homeHandler).Methods("GET")
	r.HandleFunc("/products", getProducts).Methods("GET")
	r.HandleFunc("/products/{id:[0-9]+}", getProductByID).Methods("GET")
	r.HandleFunc("/products", createProduct).Methods("POST")
	r.HandleFunc("/products/{id:[0-9]+}", updateProduct).Methods("PUT")
	r.HandleFunc("/products/{id:[0-9]+}", deleteProduct).Methods("DELETE")

	// Route با Query Parameters
	r.HandleFunc("/search", searchProducts).Queries("q", "{q}", "page", "{page:[0-9]+}")

	// Subrouter
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/products", getProducts).Methods("GET")
	api.HandleFunc("/products/{id}", getProductByID).Methods("GET")

	// Static files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	// Custom 404
	r.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	http.ListenAndServe(":8080", r)
}

// 2.2 URL Parameters (Mux)
func getProductByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	product, exists := products[id]
	if !exists {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(product)
}

// 2.3 Path Prefix و Subrouter
func subrouterExample() {
	r := mux.NewRouter()

	// Admin subrouter
	admin := r.PathPrefix("/admin").Subrouter()
	admin.Use(adminAuthMiddleware)
	admin.HandleFunc("/users", getUsers)
	admin.HandleFunc("/stats", getStats)

	// Public subrouter
	public := r.PathPrefix("/public").Subrouter()
	public.HandleFunc("/info", getPublicInfo)
}

// 2.4 Reverse Routing (نام‌گذاری Routes)
func namedRoutes() {
	r := mux.NewRouter()

	// نام‌گذاری route
	r.HandleFunc("/users/{id}", getUser).Name("getUser")
	r.HandleFunc("/users", createUser).Name("createUser")

	// استفاده از نام برای ساخت URL
	url, _ := r.Get("getUser").URL("id", "123")
	// url.String() = "/users/123"
}

// 2.5 Matcher Functions
func matcherExample() {
	r := mux.NewRouter()

	// بر اساس Header
	r.HandleFunc("/api", apiHandler).Headers("X-API-Version", "2.0")

	// بر اساس Scheme (http/https)
	r.HandleFunc("/secure", secureHandler).Schemes("https")

	// بر اساس Host
	r.HandleFunc("/", hostHandler).Host("api.example.com")

	// ترکیب شرط‌ها
	r.HandleFunc("/special", specialHandler).
		Methods("POST").
		Headers("Content-Type", "application/json").
		Host("api.example.com")
}

// 2.6 Middleware در Mux
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}

func adminAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token != "Bearer secret" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Home Page"))
}

func getProducts(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(products)
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	var product Product
	json.NewDecoder(r.Body).Decode(&product)
	product.ID = len(products) + 1
	products[product.ID] = product
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

func updateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var product Product
	json.NewDecoder(r.Body).Decode(&product)
	product.ID = id
	products[id] = product

	json.NewEncoder(w).Encode(product)
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	delete(products, id)
	w.WriteHeader(http.StatusNoContent)
}

func searchProducts(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	page := r.URL.Query().Get("page")
	w.Write([]byte(fmt.Sprintf("Search: %s, Page: %s", query, page)))
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Admin: Users list"))
}

func getStats(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Admin: Stats"))
}

func getPublicInfo(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Public Info"))
}

func getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Write([]byte("User ID: " + vars["id"]))
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("API v2.0"))
}

func secureHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Secure"))
}

func hostHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("API endpoint"))
}

func specialHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Special endpoint"))
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Custom 404 - Page not found", http.StatusNotFound)
}
*/

// ============================================================================
// بخش 3: Gin Framework
// ============================================================================

/*
نصب:
go get github.com/gin-gonic/gin
*/

/*
// مثال کامل با Gin
package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Task struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

var tasks = map[int]Task{
	1: {ID: 1, Title: "Learn Go", Completed: false, CreatedAt: time.Now()},
	2: {ID: 2, Title: "Build API", Completed: false, CreatedAt: time.Now()},
}

// 3.1 Basic Gin Router
func basicGinRouter() {
	r := gin.Default()

	// Middlewareها
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Routes
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello Gin!")
	})

	r.GET("/tasks", getTasks)
	r.GET("/tasks/:id", getTaskByID)
	r.POST("/tasks", createTask)
	r.PUT("/tasks/:id", updateTask)
	r.DELETE("/tasks/:id", deleteTask)

	// Grouping
	v1 := r.Group("/api/v1")
	{
		v1.GET("/tasks", getTasks)
		v1.POST("/tasks", createTask)
	}

	// Static files
	r.Static("/static", "./static")
	r.StaticFile("/favicon.ico", "./static/favicon.ico")

	r.Run(":8080")
}

// 3.2 URL Parameters (Gin)
func getTaskByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	task, exists := tasks[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, task)
}

// 3.3 Query Parameters
func getTasks(c *gin.Context) {
	// ?page=1&limit=10&search=golang
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")

	c.JSON(http.StatusOK, gin.H{
		"tasks":  tasks,
		"page":   page,
		"limit":  limit,
		"search": search,
	})
}

// 3.4 Request Body Binding
type CreateTaskRequest struct {
	Title     string    `json:"title" binding:"required"`
	Completed bool      `json:"completed"`
	DueDate   time.Time `json:"due_date" binding:"required"`
}

func createTask(c *gin.Context) {
	var req CreateTaskRequest

	// Validation خودکار
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task := Task{
		ID:        len(tasks) + 1,
		Title:     req.Title,
		Completed: req.Completed,
		CreatedAt: time.Now(),
	}
	tasks[task.ID] = task

	c.JSON(http.StatusCreated, task)
}

// 3.5 Form Data Binding
type LoginForm struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func loginHandler(c *gin.Context) {
	var form LoginForm
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

// 3.6 File Upload
func uploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save file
	c.SaveUploadedFile(file, "./uploads/"+file.Filename)

	c.JSON(http.StatusOK, gin.H{
		"filename": file.Filename,
		"size":     file.Size,
	})
}

// 3.7 Custom Middleware in Gin
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token != "Bearer secret-token" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		c.Next()
	}
}

func adminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetHeader("X-User-Role")
		if role != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Admin only"})
			return
		}
		c.Next()
	}
}

// 3.8 Route Grouping with Middleware
func groupWithMiddleware() {
	r := gin.Default()

	// Public routes
	public := r.Group("/api")
	{
		public.GET("/health", healthCheck)
		public.POST("/login", loginHandler)
	}

	// Protected routes
	protected := r.Group("/api")
	protected.Use(authMiddleware())
	{
		protected.GET("/profile", getProfile)
		protected.PUT("/profile", updateProfile)
	}

	// Admin routes
	admin := r.Group("/api/admin")
	admin.Use(authMiddleware(), adminMiddleware())
	{
		admin.GET("/users", listUsers)
		admin.DELETE("/users/:id", deleteUser)
	}
}

// 3.9 Gin with HTML Templates
func ginWithTemplates() {
	r := gin.Default()

	// Load templates
	r.LoadHTMLGlob("templates/*")

	r.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title":   "My App",
			"message": "Hello from Gin!",
		})
	})
}

// 3.10 Gin with Custom HTTP Config
func customHTTPConfig() {
	r := gin.Default()

	s := &http.Server{
		Addr:           ":8080",
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	s.ListenAndServe()
}

// Handler functions
func updateTask(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var task Task
	c.BindJSON(&task)
	task.ID = id
	tasks[id] = task
	c.JSON(http.StatusOK, task)
}

func deleteTask(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	delete(tasks, id)
	c.JSON(http.StatusNoContent, nil)
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func getProfile(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"user": "profile"})
}

func updateProfile(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "profile updated"})
}

func listUsers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"users": []string{"admin", "user1"}})
}

func deleteUser(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "user " + id + " deleted"})
}
*/

// ============================================================================
// بخش 4: مقایسه و راهنمای انتخاب
// ============================================================================

func compareRouters() {
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("📊 COMPARISON TABLE")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ FEATURE              │ Chi          │ Gorilla/Mux  │ Gin                      │
├──────────────────────┼──────────────┼──────────────┼──────────────────────────┤
│ Speed                │ Very Fast    │ Fast         │ Fastest                  │
│ net/http compatible  │ ✅ 100%      │ ✅ 100%      │ ⚠️ Wrapper               │
│ Middleware           │ ✅ Built-in  │ ✅ External  │ ✅ Built-in              │
│ Route grouping       │ ✅ Excellent │ ✅ Good      │ ✅ Excellent             │
│ URL parameters       │ ✅           │ ✅           │ ✅                       │
│ Regex in routes      │ ✅           │ ✅           │ ❌ Limited               │
│ Reverse routing      │ ❌           │ ✅           │ ❌                       │
│ Validation           │ ❌           │ ❌           │ ✅ Built-in              │
│ Binding              │ Manual       │ Manual       │ ✅ Automatic             │
│ Error handling       │ Manual       │ Manual       │ ✅ Built-in              │
│ Documentation        │ ✅ Good      │ ✅ Good      │ ✅ Excellent             │
│ Maintenance          │ ✅ Active    │ ⚠️ Slow      │ ✅ Active                │
│ Learning curve       │ Easy         │ Easy         │ Medium                   │
│ Ecosystem            │ Medium       │ Large        │ Very Large               │
│ Performance          │ 10-20 ns/op  │ 50-100 ns/op │ 1-5 ns/op                │
└──────────────────────┴──────────────┴──────────────┴──────────────────────────┘

📌 RECOMMENDATIONS:

   🏆 Chi (توصیه اصلی):
   • پروژه‌های جدید Go
   • Microservices
   • APIهایی که با net/http استاندارد کار می‌کنند
   • تیم‌هایی که سادگی را ترجیح می‌دهند

   📦 Gorilla/Mux:
   • پروژه‌های legacy که از Gorilla استفاده می‌کنند
   • نیاز به reverse routing دارید
   • پروژه‌های پایدار با تغییرات کم

   ⚡ Gin:
   • Performance-critical applications
   • نیاز به validation و binding خودکار دارید
   • تیم به framework full-featured عادت دارد (مثل Express.js)
   • پروژه‌های بزرگ با نیاز به ecosystem گسترده

🎯 SELECTION FLOWCHART:

   Start
     │
     ▼
   نیاز به max performance؟ ────YES───► Use Gin
     │
     NO
     ▼
   پروژه legacy با Gorilla؟ ────YES───► Use Gorilla/Mux
     │
     NO
     ▼
   Want standard net/http? ────YES───► Use Chi
     │
     NO
     ▼
   Need built-in validation? ────YES───► Use Gin
     │
     NO
     ▼
   Default: Use Chi ✅
`)
}

// ============================================================================
// بخش 5: Common Patterns
// ============================================================================

func commonPatterns() {
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("🔧 COMMON ROUTING PATTERNS")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────┐
│ 1. CRUD API Pattern                                            │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  GET    /resource          → List all resources                │
│  GET    /resource/{id}     → Get single resource               │
│  POST   /resource          → Create new resource               │
│  PUT    /resource/{id}     → Update entire resource            │
│  PATCH  /resource/{id}     → Partial update                    │
│  DELETE /resource/{id}     → Delete resource                   │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ 2. Nested Resources                                            │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  GET    /users/{userId}/posts              → User's posts      │
│  GET    /users/{userId}/posts/{postId}     → Specific post     │
│  POST   /users/{userId}/posts              → Create post       │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ 3. API Versioning                                              │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  /api/v1/users         → Version 1 API                         │
│  /api/v2/users         → Version 2 API                         │
│  /api/latest/users     → Latest version                        │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ 4. Middleware Chain                                            │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  Logging → Recovery → RateLimit → Auth → Handler               │
│                                                                 │
│  Order matters!                                                │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
`)
}

// ============================================================================
// بخش 6: Performance Comparison
// ============================================================================

func performanceComparison() {
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("⚡ PERFORMANCE COMPARISON")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
Benchmark results (typical):

BenchmarkChi_Simple-8        5000000    250 ns/op    0 B/op    0 allocs/op
BenchmarkChi_Param-8         3000000    450 ns/op    0 B/op    0 allocs/op

BenchmarkGorilla_Simple-8    2000000    800 ns/op    0 B/op    0 allocs/op
BenchmarkGorilla_Param-8     1000000   1200 ns/op    0 B/op    0 allocs/op

BenchmarkGin_Simple-8       10000000    120 ns/op    0 B/op    0 allocs/op
BenchmarkGin_Param-8         8000000    200 ns/op    0 B/op    0 allocs/op

Note: Results vary based on hardware and use case.
Always benchmark with your specific workload.
`)
}

// ============================================================================
// بخش 7: جمع‌بندی
// ============================================================================

func main() {
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("🎯 COMPLETE ROUTING GUIDE")
	fmt.Println("Chi | Gorilla/Mux | Gin")
	fmt.Println(strings.Repeat("=", 80))

	compareRouters()
	commonPatterns()
	performanceComparison()

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📚 QUICK REFERENCE")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
💡 INSTALLATION:

   # Chi
   go get github.com/go-chi/chi/v5
   go get github.com/go-chi/chi/v5/middleware

   # Gorilla/Mux
   go get github.com/gorilla/mux

   # Gin
   go get github.com/gin-gonic/gin

💡 BASIC SERVER:

   // Chi
   r := chi.NewRouter()
   r.Get("/", handler)
   http.ListenAndServe(":8080", r)

   // Gorilla
   r := mux.NewRouter()
   r.HandleFunc("/", handler)
   http.ListenAndServe(":8080", r)

   // Gin
   r := gin.Default()
   r.GET("/", handler)
   r.Run(":8080")

💡 BEST PRACTICES:

   1. Always group related routes
   2. Use middleware for cross-cutting concerns
   3. Keep handlers thin (move logic to services)
   4. Use proper HTTP status codes
   5. Validate input early
   6. Log requests (but not sensitive data)
   7. Set timeouts for all handlers
   8. Use context for cancellation
   9. Handle panics with recovery middleware
   10. Document your API (OpenAPI/Swagger)

🎯 MY RECOMMENDATION:

   Start with Chi for new projects ✅
   It's simple, fast, and compatible with standard library.
   Switch to Gin if you need built-in validation/performance.
   Use Gorilla only for existing projects.
`)

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🎯 ROUTING GUIDE - COMPLETE")
	fmt.Println("Choose the right router for your project!")
	fmt.Println(strings.Repeat("=", 80))
}

// تابع کمکی برای تکرار رشته
func stringsRepeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}
