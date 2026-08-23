// ============================================================================
// FILE: clean_architecture_guide.go
// TITLE: راهنمای کامل Clean Architecture در Go - HTTP → UseCase → Repository → DB
// HOW TO RUN: go run clean_architecture_guide.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - Clean Architecture چیست؟
// ============================================================================
//
// Clean Architecture (معماری تمیز) توسط Robert C. Martin (Uncle Bob) ارائه شده است.
// هدف اصلی: جداسازی لایه‌ها و وابستگی‌ها به گونه‌ای که تغییر در یک لایه
// تأثیری روی لایه‌های دیگر نداشته باشد.
//
// لایه‌های Clean Architecture (از بیرون به داخل):
//
// 1. Framework & Drivers (HTTP, Database, External APIs)
//    - جزئیات خارجی: Gin, Chi, PostgreSQL, Redis, etc.
//    - این لایه وابسته به لایه‌های داخلی است
//
// 2. Interface Adapters (Controllers, Presenters, Repositories)
//    - تبدیل داده‌ها بین لایه‌ها
//    - پیاده‌سازی Repository interfaceها
//
// 3. Use Cases (Application Business Rules)
//    - منطق کسب و کار برنامه
//    - مستقل از دیتابیس و HTTP
//
// 4. Entities (Enterprise Business Rules)
//    - مدل‌های اصلی دامنه
//    - مستقل از همه چیز
//
// قانون طلایی:
// "وابستگی‌ها فقط به سمت داخل می‌روند.
//  لایه‌های داخلی نباید از لایه‌های خارجی چیزی بدانند.
//  اینترفیس‌ها را در سمت مصرف‌کننده تعریف کن (نه سمت پیاده‌سازی)."
// ============================================================================

package __clean_architecture

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
// بخش 1: Entities (لایه داخلی - مستقل از همه چیز)
// ============================================================================

// Entity: User (مدل اصلی دامنه)
type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Age       int       `json:"age"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// متدهای دامنه (Business Logic سطح پایین)
func (u *User) IsAdult() bool {
	return u.Age >= 18
}

func (u *User) CanAccessAdminPanel() bool {
	return u.Role == "admin" || u.Role == "super_admin"
}

func (u *User) UpdateProfile(name string, age int) error {
	if name == "" {
		return errors.New("name cannot be empty")
	}
	if age < 0 || age > 150 {
		return errors.New("invalid age")
	}
	u.Name = name
	u.Age = age
	u.UpdatedAt = time.Now()
	return nil
}

// Entity: Product
type Product struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Stock    int     `json:"stock"`
	Category string  `json:"category"`
}

func (p *Product) IsAvailable() bool {
	return p.Stock > 0
}

func (p *Product) ReduceStock(quantity int) error {
	if p.Stock < quantity {
		return errors.New("insufficient stock")
	}
	p.Stock -= quantity
	return nil
}

// Entity: Order
type Order struct {
	ID        string      `json:"id"`
	UserID    string      `json:"user_id"`
	Items     []OrderItem `json:"items"`
	Total     float64     `json:"total"`
	Status    string      `json:"status"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type OrderItem struct {
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

func (o *Order) CalculateTotal() {
	var total float64
	for _, item := range o.Items {
		total += float64(item.Quantity) * item.Price
	}
	o.Total = total
}

func (o *Order) CanBeCancelled() bool {
	return o.Status == "pending" || o.Status == "processing"
}

// ============================================================================
// بخش 2: Repository Interfaces (تعریف در لایه UseCase)
// ============================================================================

// UserRepository اینترفیس دسترسی به داده‌های کاربر
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, offset, limit int) ([]User, int64, error)
}

// ProductRepository اینترفیس دسترسی به داده‌های محصول
type ProductRepository interface {
	Create(ctx context.Context, product *Product) error
	GetByID(ctx context.Context, id string) (*Product, error)
	Update(ctx context.Context, product *Product) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, offset, limit int) ([]Product, int64, error)
	UpdateStock(ctx context.Context, id string, quantity int) error
}

// OrderRepository اینترفیس دسترسی به داده‌های سفارش
type OrderRepository interface {
	Create(ctx context.Context, order *Order) error
	GetByID(ctx context.Context, id string) (*Order, error)
	GetByUserID(ctx context.Context, userID string, offset, limit int) ([]Order, int64, error)
	Update(ctx context.Context, order *Order) error
	UpdateStatus(ctx context.Context, id string, status string) error
}

// ============================================================================
// بخش 3: Use Cases (لایه منطق کسب و کار)
// ============================================================================

// 3.1 User Use Cases
type UserUseCases struct {
	userRepo UserRepository
	// می‌توان وابستگی‌های دیگر اضافه کرد
}

func NewUserUseCases(userRepo UserRepository) *UserUseCases {
	return &UserUseCases{
		userRepo: userRepo,
	}
}

// RegisterUser ثبت‌نام کاربر جدید
func (uc *UserUseCases) RegisterUser(ctx context.Context, name, email string, age int) (*User, error) {
	// اعتبارسنجی
	if name == "" {
		return nil, errors.New("name is required")
	}
	if email == "" {
		return nil, errors.New("email is required")
	}
	if age < 0 || age > 150 {
		return nil, errors.New("invalid age")
	}

	// بررسی تکراری نبودن ایمیل
	existing, _ := uc.userRepo.GetByEmail(ctx, email)
	if existing != nil {
		return nil, errors.New("email already taken")
	}

	// ایجاد کاربر
	user := &User{
		ID:        generateID(),
		Name:      name,
		Email:     email,
		Age:       age,
		Role:      "user",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// GetUserProfile دریافت پروفایل کاربر
func (uc *UserUseCases) GetUserProfile(ctx context.Context, userID string) (*User, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return user, nil
}

// UpdateUserProfile به‌روزرسانی پروفایل کاربر
func (uc *UserUseCases) UpdateUserProfile(ctx context.Context, userID, name string, age int) (*User, error) {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if err := user.UpdateProfile(name, age); err != nil {
		return nil, err
	}

	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

// DeleteUser حذف کاربر
func (uc *UserUseCases) DeleteUser(ctx context.Context, userID string) error {
	// بررسی وجود کاربر
	_, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return errors.New("user not found")
	}

	return uc.userRepo.Delete(ctx, userID)
}

// ListUsers لیست کاربران با pagination
func (uc *UserUseCases) ListUsers(ctx context.Context, page, pageSize int) ([]User, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	return uc.userRepo.List(ctx, offset, pageSize)
}

// 3.2 Order Use Cases
type OrderUseCases struct {
	orderRepo   OrderRepository
	productRepo ProductRepository
	userRepo    UserRepository
}

func NewOrderUseCases(orderRepo OrderRepository, productRepo ProductRepository, userRepo UserRepository) *OrderUseCases {
	return &OrderUseCases{
		orderRepo:   orderRepo,
		productRepo: productRepo,
		userRepo:    userRepo,
	}
}

// CreateOrder ایجاد سفارش جدید
func (uc *OrderUseCases) CreateOrder(ctx context.Context, userID string, items []OrderItem) (*Order, error) {
	// بررسی وجود کاربر
	_, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if len(items) == 0 {
		return nil, errors.New("order must have at least one item")
	}

	// بررسی موجودی محصولات و محاسبه قیمت
	var orderItems []OrderItem
	var total float64

	for i, item := range items {
		product, err := uc.productRepo.GetByID(ctx, item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("product %s not found", item.ProductID)
		}

		if !product.IsAvailable() {
			return nil, fmt.Errorf("product %s is out of stock", product.Name)
		}

		if product.Stock < item.Quantity {
			return nil, fmt.Errorf("insufficient stock for product %s", product.Name)
		}

		orderItems = append(orderItems, OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Price,
		})
		total += float64(item.Quantity) * product.Price

		// کاهش موجودی (در یک تراکنش واقعی باید انجام شود)
		if err := uc.productRepo.UpdateStock(ctx, item.ProductID, -item.Quantity); err != nil {
			return nil, fmt.Errorf("failed to update stock: %w", err)
		}
	}

	order := &Order{
		ID:        generateID(),
		UserID:    userID,
		Items:     orderItems,
		Total:     total,
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := uc.orderRepo.Create(ctx, order); err != nil {
		// در صورت خطا، موجودی را برگردان
		for _, item := range items {
			uc.productRepo.UpdateStock(ctx, item.ProductID, item.Quantity)
		}
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	return order, nil
}

// GetOrder دریافت سفارش
func (uc *OrderUseCases) GetOrder(ctx context.Context, orderID string) (*Order, error) {
	return uc.orderRepo.GetByID(ctx, orderID)
}

// CancelOrder لغو سفارش
func (uc *OrderUseCases) CancelOrder(ctx context.Context, orderID string) error {
	order, err := uc.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return errors.New("order not found")
	}

	if !order.CanBeCancelled() {
		return errors.New("order cannot be cancelled")
	}

	// برگرداندن موجودی محصولات
	for _, item := range order.Items {
		if err := uc.productRepo.UpdateStock(ctx, item.ProductID, item.Quantity); err != nil {
			log.Printf("failed to restore stock for product %s: %v", item.ProductID, err)
		}
	}

	return uc.orderRepo.UpdateStatus(ctx, orderID, "cancelled")
}

// ============================================================================
// بخش 4: Repository Implementations (لایه Adapter - پیاده‌سازی Concrete)
// ============================================================================

// 4.1 In-Memory Repository (برای تست و توسعه)
type InMemoryUserRepository struct {
	users map[string]*User
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: make(map[string]*User),
	}
}

func (r *InMemoryUserRepository) Create(ctx context.Context, user *User) error {
	if _, exists := r.users[user.ID]; exists {
		return errors.New("user already exists")
	}
	r.users[user.ID] = user
	return nil
}

func (r *InMemoryUserRepository) GetByID(ctx context.Context, id string) (*User, error) {
	user, exists := r.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (r *InMemoryUserRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (r *InMemoryUserRepository) Update(ctx context.Context, user *User) error {
	if _, exists := r.users[user.ID]; !exists {
		return errors.New("user not found")
	}
	r.users[user.ID] = user
	return nil
}

func (r *InMemoryUserRepository) Delete(ctx context.Context, id string) error {
	if _, exists := r.users[id]; !exists {
		return errors.New("user not found")
	}
	delete(r.users, id)
	return nil
}

func (r *InMemoryUserRepository) List(ctx context.Context, offset, limit int) ([]User, int64, error) {
	users := make([]User, 0, len(r.users))
	for _, u := range r.users {
		users = append(users, *u)
	}

	total := int64(len(users))

	start := offset
	if start > len(users) {
		start = len(users)
	}
	end := start + limit
	if end > len(users) {
		end = len(users)
	}

	return users[start:end], total, nil
}

// 4.2 In-Memory Product Repository
type InMemoryProductRepository struct {
	products map[string]*Product
	mu       sync.RWMutex
}

func NewInMemoryProductRepository() *InMemoryProductRepository {
	return &InMemoryProductRepository{
		products: make(map[string]*Product),
	}
}

func (r *InMemoryProductRepository) Create(ctx context.Context, product *Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.products[product.ID] = product
	return nil
}

func (r *InMemoryProductRepository) GetByID(ctx context.Context, id string) (*Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	product, exists := r.products[id]
	if !exists {
		return nil, errors.New("product not found")
	}
	return product, nil
}

func (r *InMemoryProductRepository) Update(ctx context.Context, product *Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.products[product.ID]; !exists {
		return errors.New("product not found")
	}
	r.products[product.ID] = product
	return nil
}

func (r *InMemoryProductRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.products, id)
	return nil
}

func (r *InMemoryProductRepository) List(ctx context.Context, offset, limit int) ([]Product, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	products := make([]Product, 0, len(r.products))
	for _, p := range r.products {
		products = append(products, *p)
	}

	total := int64(len(products))
	start := offset
	if start > len(products) {
		start = len(products)
	}
	end := start + limit
	if end > len(products) {
		end = len(products)
	}

	return products[start:end], total, nil
}

func (r *InMemoryProductRepository) UpdateStock(ctx context.Context, id string, delta int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	product, exists := r.products[id]
	if !exists {
		return errors.New("product not found")
	}
	product.Stock += delta
	return nil
}

// 4.3 In-Memory Order Repository
type InMemoryOrderRepository struct {
	orders map[string]*Order
	mu     sync.RWMutex
}

func NewInMemoryOrderRepository() *InMemoryOrderRepository {
	return &InMemoryOrderRepository{
		orders: make(map[string]*Order),
	}
}

func (r *InMemoryOrderRepository) Create(ctx context.Context, order *Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.orders[order.ID] = order
	return nil
}

func (r *InMemoryOrderRepository) GetByID(ctx context.Context, id string) (*Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	order, exists := r.orders[id]
	if !exists {
		return nil, errors.New("order not found")
	}
	return order, nil
}

func (r *InMemoryOrderRepository) GetByUserID(ctx context.Context, userID string, offset, limit int) ([]Order, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var userOrders []Order
	for _, order := range r.orders {
		if order.UserID == userID {
			userOrders = append(userOrders, *order)
		}
	}

	total := int64(len(userOrders))
	start := offset
	if start > len(userOrders) {
		start = len(userOrders)
	}
	end := start + limit
	if end > len(userOrders) {
		end = len(userOrders)
	}

	return userOrders[start:end], total, nil
}

func (r *InMemoryOrderRepository) Update(ctx context.Context, order *Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.orders[order.ID]; !exists {
		return errors.New("order not found")
	}
	r.orders[order.ID] = order
	return nil
}

func (r *InMemoryOrderRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	order, exists := r.orders[id]
	if !exists {
		return errors.New("order not found")
	}
	order.Status = status
	order.UpdatedAt = time.Now()
	return nil
}

// ============================================================================
// بخش 5: Controllers/Handlers (لایه Interface Adapters - HTTP)
// ============================================================================

// 5.1 DTOs (Data Transfer Objects) - جدا از Entities
type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

type UpdateUserRequest struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type CreateOrderRequest struct {
	Items []OrderItemRequest `json:"items"`
}

type OrderItemRequest struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

// 5.2 User HTTP Handler
type UserHandler struct {
	userUseCases *UserUseCases
}

func NewUserHandler(userUseCases *UserUseCases) *UserHandler {
	return &UserHandler{
		userUseCases: userUseCases,
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := h.userUseCases.RegisterUser(r.Context(), req.Name, req.Email, req.Age)
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	sendJSON(w, http.StatusCreated, Response{
		Success: true,
		Data:    user,
		Message: "User created successfully",
	})
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("id")
	if userID == "" {
		sendError(w, http.StatusBadRequest, "user id is required")
		return
	}

	user, err := h.userUseCases.GetUserProfile(r.Context(), userID)
	if err != nil {
		sendError(w, http.StatusNotFound, err.Error())
		return
	}

	sendJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    user,
	})
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("id")
	if userID == "" {
		sendError(w, http.StatusBadRequest, "user id is required")
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := h.userUseCases.UpdateUserProfile(r.Context(), userID, req.Name, req.Age)
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	sendJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    user,
		Message: "User updated successfully",
	})
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("id")
	if userID == "" {
		sendError(w, http.StatusBadRequest, "user id is required")
		return
	}

	if err := h.userUseCases.DeleteUser(r.Context(), userID); err != nil {
		sendError(w, http.StatusNotFound, err.Error())
		return
	}

	sendJSON(w, http.StatusOK, Response{
		Success: true,
		Message: "User deleted successfully",
	})
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	page := getIntParam(r, "page", 1)
	pageSize := getIntParam(r, "page_size", 10)

	users, total, err := h.userUseCases.ListUsers(r.Context(), page, pageSize)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendJSON(w, http.StatusOK, Response{
		Success: true,
		Data: map[string]interface{}{
			"users":       users,
			"total":       total,
			"page":        page,
			"page_size":   pageSize,
			"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// 5.3 Order HTTP Handler
type OrderHandler struct {
	orderUseCases *OrderUseCases
}

func NewOrderHandler(orderUseCases *OrderUseCases) *OrderHandler {
	return &OrderHandler{
		orderUseCases: orderUseCases,
	}
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		sendError(w, http.StatusBadRequest, "user_id is required")
		return
	}

	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	items := make([]OrderItem, len(req.Items))
	for i, item := range req.Items {
		items[i] = OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		}
	}

	order, err := h.orderUseCases.CreateOrder(r.Context(), userID, items)
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	sendJSON(w, http.StatusCreated, Response{
		Success: true,
		Data:    order,
		Message: "Order created successfully",
	})
}

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	orderID := r.URL.Query().Get("id")
	if orderID == "" {
		sendError(w, http.StatusBadRequest, "order id is required")
		return
	}

	order, err := h.orderUseCases.GetOrder(r.Context(), orderID)
	if err != nil {
		sendError(w, http.StatusNotFound, err.Error())
		return
	}

	sendJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    order,
	})
}

func (h *OrderHandler) CancelOrder(w http.ResponseWriter, r *http.Request) {
	orderID := r.URL.Query().Get("id")
	if orderID == "" {
		sendError(w, http.StatusBadRequest, "order id is required")
		return
	}

	if err := h.orderUseCases.CancelOrder(r.Context(), orderID); err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	sendJSON(w, http.StatusOK, Response{
		Success: true,
		Message: "Order cancelled successfully",
	})
}

// ============================================================================
// بخش 6: Dependency Injection / Composition Root
// ============================================================================

// App تمام وابستگی‌های برنامه را مدیریت می‌کند
type App struct {
	userHandler  *UserHandler
	orderHandler *OrderHandler
	httpServer   *http.Server
}

// NewApp ایجاد برنامه جدید (Composition Root)
func NewApp() *App {
	// لایه Repository (In-Memory برای مثال)
	userRepo := NewInMemoryUserRepository()
	productRepo := NewInMemoryProductRepository()
	orderRepo := NewInMemoryOrderRepository()

	// اضافه کردن داده نمونه
	seedData(productRepo)

	// لایه UseCase
	userUseCases := NewUserUseCases(userRepo)
	orderUseCases := NewOrderUseCases(orderRepo, productRepo, userRepo)

	// لایه Handler
	userHandler := NewUserHandler(userUseCases)
	orderHandler := NewOrderHandler(orderUseCases)

	// Router
	mux := http.NewServeMux()
	mux.HandleFunc("POST /users", userHandler.CreateUser)
	mux.HandleFunc("GET /users", userHandler.GetUser)
	mux.HandleFunc("PUT /users", userHandler.UpdateUser)
	mux.HandleFunc("DELETE /users", userHandler.DeleteUser)
	mux.HandleFunc("GET /users/list", userHandler.ListUsers)

	mux.HandleFunc("POST /orders", orderHandler.CreateOrder)
	mux.HandleFunc("GET /orders", orderHandler.GetOrder)
	mux.HandleFunc("POST /orders/cancel", orderHandler.CancelOrder)

	// Health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		sendJSON(w, http.StatusOK, Response{Success: true, Message: "OK"})
	})

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	return &App{
		userHandler:  userHandler,
		orderHandler: orderHandler,
		httpServer:   server,
	}
}

// Run اجرای برنامه
func (a *App) Run() error {
	log.Println("Server starting on :8080")
	log.Println("Endpoints:")
	log.Println("  POST   /users           - Create user")
	log.Println("  GET    /users?id=xxx    - Get user")
	log.Println("  PUT    /users?id=xxx    - Update user")
	log.Println("  DELETE /users?id=xxx    - Delete user")
	log.Println("  GET    /users/list      - List users")
	log.Println("  POST   /orders?user_id=xxx - Create order")
	log.Println("  GET    /orders?id=xxx   - Get order")
	log.Println("  POST   /orders/cancel?id=xxx - Cancel order")
	log.Println("  GET    /health          - Health check")

	return a.httpServer.ListenAndServe()
}

// Shutdown خاموش کردن برنامه
func (a *App) Shutdown(ctx context.Context) error {
	return a.httpServer.Shutdown(ctx)
}

// ============================================================================
// بخش 7: Helper Functions
// ============================================================================

func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func sendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func sendError(w http.ResponseWriter, status int, message string) {
	sendJSON(w, status, Response{
		Success: false,
		Error:   message,
	})
}

func getIntParam(r *http.Request, key string, defaultValue int) int {
	val := r.URL.Query().Get(key)
	if val == "" {
		return defaultValue
	}
	var result int
	fmt.Sscanf(val, "%d", &result)
	if result <= 0 {
		return defaultValue
	}
	return result
}

func seedData(productRepo ProductRepository) {
	products := []*Product{
		{ID: generateID(), Name: "Laptop", Price: 999.99, Stock: 10, Category: "Electronics"},
		{ID: generateID(), Name: "Mouse", Price: 19.99, Stock: 50, Category: "Electronics"},
		{ID: generateID(), Name: "Keyboard", Price: 49.99, Stock: 30, Category: "Electronics"},
	}

	for _, p := range products {
		productRepo.Create(context.Background(), p)
	}
}

// ============================================================================
// بخش 8: Best Practices
// ============================================================================

func bestPractices() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("💡 CLEAN ARCHITECTURE BEST PRACTICES")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ DEPENDENCY RULE                                                           │
├─────────────────────────────────────────────────────────────────────────────┤
│ • وابستگی‌ها فقط به سمت داخل می‌روند                                        │
│ • لایه‌های داخلی نباید از لایه‌های خارجی چیزی بدانند                        │
│ • اینترفیس‌ها را در سمت مصرف‌کننده تعریف کن                                 │
│ • از Dependency Injection استفاده کن                                        │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ PACKAGE STRUCTURE                                                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   ├── cmd/                                                                 │
│   │   └── server/                                                          │
│   │       └── main.go                                                      │
│   ├── internal/                                                            │
│   │   ├── domain/        # Entities                                        │
│   │   │   ├── user.go                                                      │
│   │   │   ├── product.go                                                   │
│   │   │   └── order.go                                                     │
│   │   ├── usecase/       # Use Cases                                       │
│   │   │   ├── user.go                                                      │
│   │   │   └── order.go                                                     │
│   │   ├── repository/    # Repository Interfaces                           │
│   │   │   └── interfaces.go                                                │
│   │   ├── adapter/       # Repository Implementations                      │
│   │   │   ├── postgres/                                                    │
│   │   │   └── inmemory/                                                    │
│   │   └── handler/       # HTTP Handlers                                   │
│   │       ├── user.go                                                      │
│   │       └── order.go                                                     │
│   └── pkg/               # Shared utilities                                │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ TESTING STRATEGY                                                          │
├─────────────────────────────────────────────────────────────────────────────┤
│ • UseCaseها را با mock Repository تست کن                                   │
│ • Repositoryها را با دیتابیس تست واقعی تست کن                              │
│ • Handlerها را با httptest تست کن                                          │
│ • Entityها را با unit test تست کن                                          │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ COMMON MISTAKES                                                           │
├─────────────────────────────────────────────────────────────────────────────┤
│ ❌ استفاده از structهای HTTP در UseCase                                     │
│ ❌ استفاده از ORM در Entity                                                 │
│ ❌ وابستگی UseCase به Repository پیاده‌سازی خاص                              │
│ ❌ قرار دادن منطق دامنه در Handler                                          │
│ ❌ برگرداندن خطاهای دیتابیس به کاربر                                        │
└─────────────────────────────────────────────────────────────────────────────┘
`)
}

// ============================================================================
// بخش 9: Main
// ============================================================================

func main() {
	fmt.Println(stringsRepeat("=", 80))
	fmt.Println("🎯 CLEAN ARCHITECTURE GUIDE")
	fmt.Println("HTTP → UseCase → Repository → DB")
	fmt.Println(stringsRepeat("=", 80))

	bestPractices()

	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🚀 Starting Clean Architecture Application")
	fmt.Println(stringsRepeat("=", 80))

	app := NewApp()

	// راه‌اندازی سرور
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

// تابع کمکی برای تکرار رشته
func stringsRepeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}

// برای رفع خطای undefined: sync
var _ = sync.RWMutex{}
