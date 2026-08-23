// ============================================================================
// FILE: package_design_guide.go
// TITLE: راهنمای کامل Package Design در Go - Flat vs Domain-Based
// HOW TO RUN: go run package_design_guide.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - چرا Package Design مهم است؟
// ============================================================================
//
// Package design در Go تأثیر مستقیم بر:
// 1. قابلیت نگهداری (Maintainability)
// 2. قابلیت تست (Testability)
// 3. قابلیت reuse
// 4. وابستگی‌ها (Dependencies)
// 5. زمان کامپایل
//
// دو رویکرد اصلی:
//
// 1. Flat (لایه‌ای / لایه‌بندی شده)
//    ├── handlers/
//    ├── services/
//    ├── repositories/
//    ├── models/
//    └── utils/
//
// 2. Domain-Based (براساس دامنه / ویژگی)
//    ├── users/
//    │   ├── handler.go
//    │   ├── service.go
//    │   ├── repository.go
//    │   └── models.go
//    ├── orders/
//    │   ├── handler.go
//    │   ├── service.go
//    │   ├── repository.go
//    │   └── models.go
//    └── products/
//
// قانون طلایی:
// "برای پروژه‌های کوچک (کمتر از 10 پکیج) از Flat استفاده کن.
//  برای پروژه‌های بزرگ و تیمی از Domain-Based استفاده کن.
//  وابستگی‌ها را به سمت داخل نگه دار (هیچ چیز نباید به پکیج utils وابسته باشد)."
// ============================================================================

package __package_design

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// ============================================================================
// بخش 1: ساختار Flat (لایه‌بندی شده)
// ============================================================================

// ============================================
// 1.1 Models (مدل‌های مشترک)
// ============================================

// FlatModels پکیج models
type FlatUser struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type FlatProduct struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

type FlatOrder struct {
	ID        string           `json:"id"`
	UserID    string           `json:"user_id"`
	Items     []FlatOrderItem  `json:"items"`
	Total     float64          `json:"total"`
	Status    string           `json:"status"`
	CreatedAt time.Time        `json:"created_at"`
}

type FlatOrderItem struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
	Price     float64 `json:"price"`
}

// ============================================
// 1.2 Repository (دسترسی به داده)
// ============================================

// FlatRepositories پکیج repositories
type FlatUserRepository interface {
	Create(ctx context.Context, user *FlatUser) error
	GetByID(ctx context.Context, id string) (*FlatUser, error)
	GetByEmail(ctx context.Context, email string) (*FlatUser, error)
}

type FlatProductRepository interface {
	GetByID(ctx context.Context, id string) (*FlatProduct, error)
	UpdateStock(ctx context.Context, id string, quantity int) error
}

type FlatOrderRepository interface {
	Create(ctx context.Context, order *FlatOrder) error
	GetByID(ctx context.Context, id string) (*FlatOrder, error)
	UpdateStatus(ctx context.Context, id string, status string) error
}

// پیاده‌سازی In-Memory
type FlatInMemoryUserRepository struct {
	users map[string]*FlatUser
	mu    sync.RWMutex
}

func NewFlatInMemoryUserRepository() *FlatInMemoryUserRepository {
	return &FlatInMemoryUserRepository{
		users: make(map[string]*FlatUser),
	}
}

func (r *FlatInMemoryUserRepository) Create(ctx context.Context, user *FlatUser) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[user.ID] = user
	return nil
}

func (r *FlatInMemoryUserRepository) GetByID(ctx context.Context, id string) (*FlatUser, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	user, exists := r.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (r *FlatInMemoryUserRepository) GetByEmail(ctx context.Context, email string) (*FlatUser, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

// ============================================
// 1.3 Service (منطق کسب و کار)
// ============================================

// FlatServices پکیج services
type FlatUserService struct {
	userRepo    FlatUserRepository
	orderRepo   FlatOrderRepository
	productRepo FlatProductRepository
}

func NewFlatUserService(
	userRepo FlatUserRepository,
	orderRepo FlatOrderRepository,
	productRepo FlatProductRepository,
) *FlatUserService {
	return &FlatUserService{
		userRepo:    userRepo,
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}

func (s *FlatUserService) RegisterUser(ctx context.Context, name, email string) (*FlatUser, error) {
	// بررسی تکراری نبودن ایمیل
	existing, _ := s.userRepo.GetByEmail(ctx, email)
	if existing != nil {
		return nil, errors.New("email already taken")
	}

	user := &FlatUser{
		ID:        generateID(),
		Name:      name,
		Email:     email,
		CreatedAt: time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *FlatUserService) GetUserOrders(ctx context.Context, userID string) ([]FlatOrder, error) {
	// اینجا باید از orderRepo استفاده کند
	return nil, nil
}

// ============================================
// 1.4 Handler (HTTP)
// ============================================

// FlatHandlers پکیج handlers
type FlatUserHandler struct {
	userService *FlatUserService
}

func NewFlatUserHandler(userService *FlatUserService) *FlatUserHandler {
	return &FlatUserHandler{
		userService: userService,
	}
}

func (h *FlatUserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user, err := h.userService.RegisterUser(r.Context(), req.Name, req.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// ============================================
// 1.5 Main (Composition Root)
// ============================================
func setupFlatArchitecture() {
	// ایجاد وابستگی‌ها
	userRepo := NewFlatInMemoryUserRepository()
	// ... سایر repositoryها

	userService := NewFlatUserService(userRepo, nil, nil)
	userHandler := NewFlatUserHandler(userService)

	// Router
	mux := http.NewServeMux()
	mux.HandleFunc("POST /users", userHandler.CreateUser)

	// شروع سرور
	log.Println("Flat architecture server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

// ============================================================================
// بخش 2: ساختار Domain-Based (براساس دامنه)
// ============================================================================

// ============================================
// 2.1 Domain: Users (پکیج مستقل)
// ============================================

// DomainUsers پکیج users
package users

import (
"context"
"encoding/json"
"errors"
"net/http"
"sync"
"time"
)

// Entity
type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// Repository Interface (تعریف در همین پکیج)
type Repository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
}

// Service Interface (منطق کسب و کار)
type Service interface {
	Register(ctx context.Context, name, email string) (*User, error)
	GetProfile(ctx context.Context, id string) (*User, error)
}

// Service Implementation
type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Register(ctx context.Context, name, email string) (*User, error) {
	// اعتبارسنجی
	if name == "" || email == "" {
		return nil, errors.New("name and email are required")
	}

	// بررسی تکراری نبودن
	existing, _ := s.repo.GetByEmail(ctx, email)
	if existing != nil {
		return nil, errors.New("email already taken")
	}

	user := &User{
		ID:        generateID(),
		Name:      name,
		Email:     email,
		CreatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) GetProfile(ctx context.Context, id string) (*User, error) {
	return s.repo.GetByID(ctx, id)
}

// In-Memory Repository Implementation
type inMemoryRepo struct {
	users map[string]*User
	mu    sync.RWMutex
}

func NewInMemoryRepository() Repository {
	return &inMemoryRepo{
		users: make(map[string]*User),
	}
}

func (r *inMemoryRepo) Create(ctx context.Context, user *User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[user.ID] = user
	return nil
}

func (r *inMemoryRepo) GetByID(ctx context.Context, id string) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	user, exists := r.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (r *inMemoryRepo) GetByEmail(ctx context.Context, email string) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

// Handler (HTTP)
type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user, err := h.service.Register(r.Context(), req.Name, req.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	user, err := h.service.GetProfile(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// ============================================
// 2.2 Domain: Orders (پکیج مستقل دیگر)
// ============================================

// DomainOrders پکیج orders
package orders

import (
"context"
"errors"
"sync"
"time"
)

// Entity
type Order struct {
	ID        string       `json:"id"`
	UserID    string       `json:"user_id"`
	Items     []OrderItem  `json:"items"`
	Total     float64      `json:"total"`
	Status    string       `json:"status"`
	CreatedAt time.Time    `json:"created_at"`
}

type OrderItem struct {
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

// Repository Interface
type Repository interface {
	Create(ctx context.Context, order *Order) error
	GetByID(ctx context.Context, id string) (*Order, error)
	GetByUserID(ctx context.Context, userID string) ([]Order, error)
	UpdateStatus(ctx context.Context, id string, status string) error
}

// Service Interface
type Service interface {
	CreateOrder(ctx context.Context, userID string, items []OrderItem) (*Order, error)
	CancelOrder(ctx context.Context, orderID string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateOrder(ctx context.Context, userID string, items []OrderItem) (*Order, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}
	if len(items) == 0 {
		return nil, errors.New("order must have items")
	}

	var total float64
	for i := range items {
		total += float64(items[i].Quantity) * items[i].Price
	}

	order := &Order{
		ID:        generateID(),
		UserID:    userID,
		Items:     items,
		Total:     total,
		Status:    "pending",
		CreatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, order); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *service) CancelOrder(ctx context.Context, orderID string) error {
	order, err := s.repo.GetByID(ctx, orderID)
	if err != nil {
		return err
	}

	if order.Status != "pending" && order.Status != "processing" {
		return errors.New("order cannot be cancelled")
	}

	return s.repo.UpdateStatus(ctx, orderID, "cancelled")
}

// In-Memory Repository
type inMemoryRepo struct {
	orders map[string]*Order
	mu     sync.RWMutex
}

func NewInMemoryRepository() Repository {
	return &inMemoryRepo{
		orders: make(map[string]*Order),
	}
}

func (r *inMemoryRepo) Create(ctx context.Context, order *Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.orders[order.ID] = order
	return nil
}

func (r *inMemoryRepo) GetByID(ctx context.Context, id string) (*Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	order, exists := r.orders[id]
	if !exists {
		return nil, errors.New("order not found")
	}
	return order, nil
}

func (r *inMemoryRepo) GetByUserID(ctx context.Context, userID string) ([]Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var userOrders []Order
	for _, order := range r.orders {
		if order.UserID == userID {
			userOrders = append(userOrders, *order)
		}
	}
	return userOrders, nil
}

func (r *inMemoryRepo) UpdateStatus(ctx context.Context, id string, status string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	order, exists := r.orders[id]
	if !exists {
		return errors.New("order not found")
	}
	order.Status = status
	return nil
}

// ============================================
// 2.3 Domain: Products
// ============================================

// DomainProducts پکیج products
package products

import (
"context"
"errors"
"sync"
)

type Product struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

type Repository interface {
	GetByID(ctx context.Context, id string) (*Product, error)
	UpdateStock(ctx context.Context, id string, delta int) error
}

type Service interface {
	GetProduct(ctx context.Context, id string) (*Product, error)
	ReserveStock(ctx context.Context, productID string, quantity int) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetProduct(ctx context.Context, id string) (*Product, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) ReserveStock(ctx context.Context, productID string, quantity int) error {
	product, err := s.repo.GetByID(ctx, productID)
	if err != nil {
		return err
	}
	if product.Quantity < quantity {
		return errors.New("insufficient stock")
	}
	return s.repo.UpdateStock(ctx, productID, -quantity)
}

type inMemoryRepo struct {
	products map[string]*Product
	mu       sync.RWMutex
}

func NewInMemoryRepository() Repository {
	return &inMemoryRepo{
		products: make(map[string]*Product),
	}
}

func (r *inMemoryRepo) GetByID(ctx context.Context, id string) (*Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	product, exists := r.products[id]
	if !exists {
		return nil, errors.New("product not found")
	}
	return product, nil
}

func (r *inMemoryRepo) UpdateStock(ctx context.Context, id string, delta int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	product, exists := r.products[id]
	if !exists {
		return errors.New("product not found")
	}
	product.Quantity += delta
	return nil
}

// ============================================
// 2.4 Main (Composition Root) با Domain-Based
// ============================================

func setupDomainBasedArchitecture() {
	// ایجاد repositoryها (هر domain مستقل)
	userRepo := users.NewInMemoryRepository()
	orderRepo := orders.NewInMemoryRepository()
	productRepo := products.NewInMemoryRepository()

	// ایجاد serviceها
	userService := users.NewService(userRepo)
	orderService := orders.NewService(orderRepo)
	productService := products.NewService(productRepo)

	// ایجاد handlerها
	userHandler := users.NewHandler(userService)
	// ... handlerهای دیگر

	// استفاده در router
	mux := http.NewServeMux()
	mux.HandleFunc("POST /users", userHandler.CreateUser)
	mux.HandleFunc("GET /users", userHandler.GetUser)

	log.Println("Domain-based architecture server starting on :8081")
	log.Fatal(http.ListenAndServe(":8081", mux))
}

// ============================================================================
// بخش 3: مقایسه Flat vs Domain-Based
// ============================================================================

func compareArchitectures() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📊 FLAT vs DOMAIN-BASED COMPARISON")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ CRITERION              │ FLAT              │ DOMAIN-BASED                  │
├────────────────────────┼───────────────────┼───────────────────────────────┤
│ Project Size           │ Small (<10 pkgs)  │ Medium to Large               │
│ Team Size              │ 1-3 developers    │ 3+ developers                 │
│ Module Coupling        │ Higher            │ Lower                          │
│ Test Isolation         │ Harder            │ Easier                         │
│ Parallel Development   │ Harder            │ Easier                         │
│ Circular Dependencies  │ Possible          │ Prevented by design           │
│ Code Reuse             │ Less              │ More (across domains)         │
│ Onboarding New Devs    │ Easier            │ Harder (more packages)        │
│ Compilation Time       │ Faster            │ Slower (more packages)        │
│ Package Visibility     │ Internal only     │ Cross-domain visibility       │
└────────────────────────┴───────────────────┴───────────────────────────────┘

📁 FLAT STRUCTURE (LAYERED):

   .
   ├── cmd/
   │   └── server/
   │       └── main.go
   ├── internal/
   │   ├── models/           # Shared models
   │   │   ├── user.go
   │   │   ├── product.go
   │   │   └── order.go
   │   ├── repositories/     # Data access
   │   │   ├── user_repo.go
   │   │   ├── product_repo.go
   │   │   └── order_repo.go
   │   ├── services/         # Business logic
   │   │   ├── user_service.go
   │   │   ├── product_service.go
   │   │   └── order_service.go
   │   └── handlers/         # HTTP handlers
   │       ├── user_handler.go
   │       ├── product_handler.go
   │       └── order_handler.go
   └── pkg/
       └── utils/            # Shared utilities
           ├── id.go
           └── errors.go

📁 DOMAIN-BASED STRUCTURE:

   .
   ├── cmd/
   │   └── server/
   │       └── main.go
   ├── internal/
   │   ├── users/            # Complete user domain
   │   │   ├── handler.go
   │   │   ├── service.go
   │   │   ├── repository.go
   │   │   ├── models.go
   │   │   └── errors.go
   │   ├── orders/           # Complete order domain
   │   │   ├── handler.go
   │   │   ├── service.go
   │   │   ├── repository.go
   │   │   ├── models.go
   │   │   └── errors.go
   │   ├── products/         # Complete product domain
   │   │   ├── handler.go
   │   │   ├── service.go
   │   │   ├── repository.go
   │   │   ├── models.go
   │   │   └── errors.go
   │   └── shared/           # Shared across domains
   │       ├── errors/
   │       ├── middleware/
   │       └── utils/
   └── pkg/
       └── (public packages)
`)
}

// ============================================================================
// بخش 4: Guidelines برای انتخاب ساختار
// ============================================================================

func selectionGuidelines() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📋 SELECTION GUIDELINES")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ WHEN TO USE FLAT ARCHITECTURE                                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ✅ پروژه‌های کوچک (کمتر از 10 پکیج)                                      │
│  ✅ تیم‌های کوچک (1-3 نفر)                                                 │
│  ✅ پروژه‌های یک‌بار مصرف یا MVP                                           │
│  ✅ زمانی که سادگی مهم‌تر از مقیاس‌پذیری است                               │
│  ✅ زمانی که دامنه‌ها وابستگی زیادی به هم دارند                             │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ WHEN TO USE DOMAIN-BASED ARCHITECTURE                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ✅ پروژه‌های بزرگ (بیش از 10 پکیج)                                       │
│  ✅ تیم‌های بزرگ (3+ نفر)                                                  │
│  ✅ پروژه‌های بلندمدت با نیاز به نگهداری بالا                             │
│  ✅ زمانی که دامنه‌ها مستقل هستند                                          │
│  ✅ زمانی که قابلیت reuse بین پروژه‌ها مهم است                             │
│  ✅ معماری میکروسرویس                                                      │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ HYBRID APPROACH (Mixed)                                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  بسیاری از پروژه‌ها از ترکیب هر دو استفاده می‌کنند:                       │
│                                                                             │
│  internal/                                                                 │
│  ├── domain/                    # Domain-based برای core domains           │
│  │   ├── users/                                                           │
│  │   ├── orders/                                                          │
│  │   └── products/                                                        │
│  ├── shared/                     # Shared across domains (flat)           │
│  │   ├── database/                                                        │
│  │   ├── middleware/                                                      │
│  │   └── utils/                                                           │
│  └── infrastructure/             # Technical concerns (flat)              │
│      ├── cache/                                                           │
│      ├── queue/                                                           │
│      └── logger/                                                          │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
`)
}

// ============================================================================
// بخش 5: Best Practices برای Package Design
// ============================================================================

func bestPractices() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("💡 PACKAGE DESIGN BEST PRACTICES")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ 1. NAMING CONVENTIONS                                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│    • Package names should be lowercase, single-word                       │
│    • Avoid generic names: "common", "util", "base"                        │
│    • Use plural for packages that contain multiple files: "users"         │
│    • Use meaningful names: "httputil" not "helper"                        │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 2. EXPORTED VS UNEXPORTED                                                │
├─────────────────────────────────────────────────────────────────────────────┤
│    • Keep exported API small                                             │
│    • Unexport what you can                                               │
│    • Use internal package for private code                               │
│    • Document all exported identifiers                                   │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 3. DEPENDENCY MANAGEMENT                                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│    • Avoid circular dependencies                                          │
│    • Use interfaces to break dependencies                                 │
│    • Keep dependency graph shallow                                        │
│    • "internal" packages cannot be imported from outside                  │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 4. PACKAGE SIZE                                                          │
├─────────────────────────────────────────────────────────────────────────────┤
│    • One responsibility per package                                       │
│    • Files should be 200-400 lines on average                            │
│    • Split when package grows too large                                  │
│    • Don't split prematurely                                             │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 5. CYCLE DEPENDENCIES                                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│    • ❌ users → orders → users (circular)                                 │
│    • ✅ users → orders (one-way)                                          │
│    • ✅ Use interfaces in separate package to break cycles                │
│    • ✅ Move shared types to a new package                                │
└─────────────────────────────────────────────────────────────────────────────┘
`)
}

// ============================================================================
// بخش 6: Anti-Patterns
// ============================================================================

func antiPatterns() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("⚠️ COMMON PACKAGE DESIGN ANTI-PATTERNS")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ 1. GOD PACKAGE                                                            │
├─────────────────────────────────────────────────────────────────────────────┤
│    ❌ package models (everything in one package)                          │
│    ✅ Separate packages: users, products, orders                          │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 2. UTILITY PACKAGES                                                       │
├─────────────────────────────────────────────────────────────────────────────┤
│    ❌ package utils (dump of unrelated functions)                         │
│    ✅ package stringutil, package timeutil, package fileutil              │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 3. CIRCULAR DEPENDENCIES                                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│    ❌ users → orders and orders → users                                   │
│    ✅ Use interfaces or move shared types to separate package             │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 4. DEEP NESTING                                                          │
├─────────────────────────────────────────────────────────────────────────────┤
│    ❌ internal/services/domain/users/repository/postgres/                 │
│    ✅ internal/users/postgres_repository.go                               │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 5. PREMATURE ABSTRACTION                                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│    ❌ Creating interfaces before needed                                   │
│    ✅ Start with concrete, abstract when needed                           │
└─────────────────────────────────────────────────────────────────────────────┘
`)
}

// ============================================================================
// بخش 7: Example Project Structure Templates
// ============================================================================

func projectTemplates() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📁 PROJECT STRUCTURE TEMPLATES")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ TEMPLATE 1: Simple Web API (Flat)                                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   myproject/                                                               │
│   ├── cmd/                                                                │
│   │   └── server/                                                         │
│   │       └── main.go                                                     │
│   ├── internal/                                                           │
│   │   ├── config/                                                         │
│   │   │   └── config.go                                                   │
│   │   ├── models/                                                         │
│   │   │   ├── user.go                                                     │
│   │   │   └── product.go                                                  │
│   │   ├── repository/                                                     │
│   │   │   ├── user.go                                                     │
│   │   │   └── product.go                                                  │
│   │   ├── service/                                                        │
│   │   │   ├── user.go                                                     │
│   │   │   └── product.go                                                  │
│   │   ├── handler/                                                        │
│   │   │   ├── user.go                                                     │
│   │   │   └── product.go                                                  │
│   │   └── middleware/                                                     │
│   │       ├── auth.go                                                     │
│   │       └── logging.go                                                  │
│   ├── pkg/                                                                │
│   │   ├── errors/                                                         │
│   │   │   └── errors.go                                                   │
│   │   └── utils/                                                          │
│   │       └── id.go                                                       │
│   ├── go.mod                                                              │
│   └── go.sum                                                              │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ TEMPLATE 2: Enterprise API (Domain-Based)                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   myproject/                                                               │
│   ├── cmd/                                                                │
│   │   ├── api/                                                            │
│   │   │   └── main.go                                                     │
│   │   └── worker/                                                         │
│   │       └── main.go                                                     │
│   ├── internal/                                                           │
│   │   ├── users/               # Complete user domain                     │
│   │   │   ├── handler.go                                                  │
│   │   │   ├── service.go                                                  │
│   │   │   ├── repository.go                                               │
│   │   │   ├── models.go                                                   │
│   │   │   └── errors.go                                                   │
│   │   ├── orders/              # Complete order domain                    │
│   │   │   ├── handler.go                                                  │
│   │   │   ├── service.go                                                  │
│   │   │   ├── repository.go                                               │
│   │   │   ├── models.go                                                   │
│   │   │   └── events.go                                                   │
│   │   ├── products/            # Complete product domain                  │
│   │   │   ├── handler.go                                                  │
│   │   │   ├── service.go                                                  │
│   │   │   ├── repository.go                                               │
│   │   │   ├── models.go                                                   │
│   │   │   └── search.go                                                   │
│   │   ├── payments/            # Complete payment domain                  │
│   │   │   ├── handler.go                                                  │
│   │   │   ├── service.go                                                  │
│   │   │   ├── repository.go                                               │
│   │   │   ├── models.go                                                   │
│   │   │   └── gateway.go                                                  │
│   │   ├── shared/              # Shared across domains                    │
│   │   │   ├── database/                                                   │
│   │   │   │   └── postgres.go                                             │
│   │   │   ├── cache/                                                      │
│   │   │   │   └── redis.go                                                │
│   │   │   ├── queue/                                                      │
│   │   │   │   └── rabbitmq.go                                             │
│   │   │   ├── logger/                                                     │
│   │   │   │   └── zap.go                                                  │
│   │   │   ├── middleware/                                                 │
│   │   │   │   ├── auth.go                                                 │
│   │   │   │   ├── logging.go                                              │
│   │   │   │   └── metrics.go                                              │
│   │   │   └── errors/                                                     │
│   │   │       └── errors.go                                               │
│   │   └── infrastructure/      # Technical infrastructure                 │
│   │       ├── grpc/                                                       │
│   │       │   └── server.go                                               │
│   │       └── http/                                                       │
│   │           └── server.go                                               │
│   ├── pkg/                                                                │
│   │   ├── api/                     # Public API clients                   │
│   │   │   └── client/                                                     │
│   │   │       └── client.go                                               │
│   │   └── sdk/                     # SDK for external use                 │
│   │       └── sdk.go                                                      │
│   ├── api/                                                                │
│   │   ├── openapi/                                                        │
│   │   │   └── swagger.yaml                                                │
│   │   └── grpc/                                                           │
│   │       └── proto/                                                      │
│   ├── migrations/                                                         │
│   │   ├── 001_create_users_table.up.sql                                   │
│   │   └── 001_create_users_table.down.sql                                 │
│   ├── scripts/                                                            │
│   │   ├── build.sh                                                        │
│   │   └── deploy.sh                                                       │
│   ├── go.mod                                                              │
│   └── go.sum                                                              │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
`)
}

// ============================================================================
// بخش 8: Summary
// ============================================================================

func summary() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📚 PACKAGE DESIGN SUMMARY")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ DECISION MAKING FLOW                                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   Start                                                                    │
│     │                                                                      │
│     ▼                                                                      │
│   Is project small (<10 packages)? ──YES──► Use FLAT Architecture         │
│     │                                                                      │
│     NO                                                                     │
│     ▼                                                                      │
│   Is team small (1-3 people)? ──────YES──► Use FLAT Architecture          │
│     │                                                                      │
│     NO                                                                     │
│     ▼                                                                      │
│   Use DOMAIN-BASED Architecture                                            │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

💡 GOLDEN RULES:

   1. Start simple, refactor as needed
   2. Keep packages focused (single responsibility)
   3. Avoid circular dependencies
   4. Use internal packages for private code
   5. Document exported symbols
   6. Keep dependency graph shallow
   7. Prefer composition over embedding for reuse
   8. Don't over-engineer - you can always restructure
   9. Use interfaces to break dependencies
   10. Follow Go proverb: "A little copying is better than a little dependency"
`)
}

// ============================================================================
// بخش 9: Main
// ============================================================================

func main() {
	fmt.Println(stringsRepeat("=", 80))
	fmt.Println("🎯 PACKAGE DESIGN IN GO")
	fmt.Println("Flat vs Domain-Based Architecture")
	fmt.Println(stringsRepeat("=", 80))

	compareArchitectures()
	selectionGuidelines()
	bestPractices()
	antiPatterns()
	projectTemplates()
	summary()

	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🎯 PACKAGE DESIGN - COMPLETE")
	fmt.Println("Choose the right structure for your Go project!")
	fmt.Println(stringsRepeat("=", 80))
}

// Helper functions
func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func stringsRepeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}