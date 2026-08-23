// ============================================================================
// FILE: go_design_patterns_guide.go
// TITLE: راهنمای کامل الگوهای طراحی در Go - Repository, Service, DI, Options, Middleware, Worker
// HOW TO RUN: go run go_design_patterns_guide.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - چرا الگوهای طراحی در Go مهم هستند؟
// ============================================================================
//
// الگوهای طراحی راه‌حل‌های اثبات‌شده برای مشکلات رایج در طراحی نرم‌افزار هستند.
// در Go، به دلیل سادگی زبان، الگوها اغلب ساده‌تر از زبان‌های دیگر پیاده‌سازی می‌شوند.
//
// 1. Repository Pattern: جداسازی منطق دسترسی به داده از منطق کسب و کار
// 2. Service Layer: encapsulate کردن منطق کسب و کار
// 3. Dependency Injection: وابستگی‌ها را از بیرون تزریق کنید (نه داخل ساخت)
// 4. Options Pattern: پیکربندی انعطاف‌پذیر توابع و structها
// 5. Middleware Pattern: پردازش درخواست‌ها در یک زنجیره
// 6. Worker Pattern: پردازش موازی کارها با تعداد محدود worker
//
// قانون طلایی:
// "از الگوها برای کاهش وابستگی و افزایش قابلیت تست استفاده کن.
//  در Go، از embedding به جای ارث‌بری استفاده کن.
//  اینترفیس‌ها را در سمت مصرف‌کننده تعریف کن، نه سمت تولیدکننده."
// ============================================================================

package __go_design_patterns

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// ============================================================================
// بخش 1: Repository Pattern
// ============================================================================

// مدل‌های دامنه (Domain Models)
type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Age       int       `json:"age"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Product struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Stock    int     `json:"stock"`
	Category string  `json:"category"`
}

// 1.1 تعریف Repository Interface (در سمت مصرف‌کننده)
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, offset, limit int) ([]User, int64, error)
}

// 1.2 پیاده‌سازی In-Memory Repository (برای تست و توسعه)
type InMemoryUserRepository struct {
	users map[string]*User
	mu    sync.RWMutex
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: make(map[string]*User),
	}
}

func (r *InMemoryUserRepository) Create(ctx context.Context, user *User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.ID]; exists {
		return errors.New("user already exists")
	}
	r.users[user.ID] = user
	return nil
}

func (r *InMemoryUserRepository) GetByID(ctx context.Context, id string) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (r *InMemoryUserRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (r *InMemoryUserRepository) Update(ctx context.Context, user *User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.ID]; !exists {
		return errors.New("user not found")
	}
	user.UpdatedAt = time.Now()
	r.users[user.ID] = user
	return nil
}

func (r *InMemoryUserRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[id]; !exists {
		return errors.New("user not found")
	}
	delete(r.users, id)
	return nil
}

func (r *InMemoryUserRepository) List(ctx context.Context, offset, limit int) ([]User, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]User, 0, len(r.users))
	for _, u := range r.users {
		users = append(users, *u)
	}

	total := int64(len(users))

	// اعمال offset و limit
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

// 1.3 پیاده‌سازی PostgreSQL Repository
type PostgresUserRepository struct {
	db *sql.DB // در مثال واقعی از sql.DB استفاده می‌شود
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

// (بدنبال پیاده‌سازی واقعی - فقط برای نمایش)
func (r *PostgresUserRepository) Create(ctx context.Context, user *User) error {
	// در واقع: INSERT INTO users ...
	return nil
}

// ============================================================================
// بخش 2: Service Layer
// ============================================================================

// 2.1 تعریف Service Interface
type UserService interface {
	Register(ctx context.Context, name, email string, age int) (*User, error)
	GetProfile(ctx context.Context, userID string) (*User, error)
	UpdateProfile(ctx context.Context, userID string, name string, age int) (*User, error)
	DeleteAccount(ctx context.Context, userID string) error
	IsEmailTaken(ctx context.Context, email string) (bool, error)
}

// 2.2 پیاده‌سازی Service با وابستگی به Repository
type userService struct {
	userRepo UserRepository
	// می‌توان وابستگی‌های دیگر اضافه کرد: logger, cache, etc.
}

// NewUserService سازنده با Dependency Injection
func NewUserService(userRepo UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) Register(ctx context.Context, name, email string, age int) (*User, error) {
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
	taken, err := s.IsEmailTaken(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email: %w", err)
	}
	if taken {
		return nil, errors.New("email already taken")
	}

	// ایجاد کاربر
	user := &User{
		ID:        generateID(),
		Name:      name,
		Email:     email,
		Age:       age,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// می‌توان رویدادهایی مثل UserRegistered منتشر کرد
	// s.eventBus.Publish(UserRegisteredEvent{UserID: user.ID})

	return user, nil
}

func (s *userService) GetProfile(ctx context.Context, userID string) (*User, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// نمی‌خواهیم اطلاعات حساس برگردانیم
	return user, nil
}

func (s *userService) UpdateProfile(ctx context.Context, userID string, name string, age int) (*User, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if name != "" {
		user.Name = name
	}
	if age > 0 {
		user.Age = age
	}
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

func (s *userService) DeleteAccount(ctx context.Context, userID string) error {
	if userID == "" {
		return errors.New("user ID is required")
	}

	// قبل از حذف، می‌توان بررسی‌های امنیتی انجام داد
	// مثلاً آیا کاربر اجازه حذف دارد؟

	if err := s.userRepo.Delete(ctx, userID); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	// ارسال رویداد UserDeleted
	return nil
}

func (s *userService) IsEmailTaken(ctx context.Context, email string) (bool, error) {
	_, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if err.Error() == "user not found" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// ============================================================================
// بخش 3: Dependency Injection (با Struct Embedding و Interface)
// ============================================================================

// 3.1 وابستگی‌های برنامه
type Dependencies struct {
	UserRepository UserRepository
	UserService    UserService
	Logger         Logger
	Cache          Cache
}

// 3.2 Logger Interface
type Logger interface {
	Info(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
	Debug(msg string, keysAndValues ...interface{})
}

// 3.3 پیاده‌سازی Logger
type SimpleLogger struct {
	prefix string
}

func NewSimpleLogger(prefix string) *SimpleLogger {
	return &SimpleLogger{prefix: prefix}
}

func (l *SimpleLogger) Info(msg string, keysAndValues ...interface{}) {
	log.Printf("[INFO] %s: %s %v", l.prefix, msg, keysAndValues)
}

func (l *SimpleLogger) Error(msg string, keysAndValues ...interface{}) {
	log.Printf("[ERROR] %s: %s %v", l.prefix, msg, keysAndValues)
}

func (l *SimpleLogger) Debug(msg string, keysAndValues ...interface{}) {
	log.Printf("[DEBUG] %s: %s %v", l.prefix, msg, keysAndValues)
}

// 3.4 Cache Interface
type Cache interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{}, ttl time.Duration)
	Delete(key string)
}

// 3.5 پیاده‌سازی In-Memory Cache
type InMemoryCache struct {
	data map[string]cacheItem
	mu   sync.RWMutex
}

type cacheItem struct {
	value      interface{}
	expiration time.Time
}

func NewInMemoryCache() *InMemoryCache {
	cache := &InMemoryCache{
		data: make(map[string]cacheItem),
	}
	// شروع cleanup routine
	go cache.cleanup()
	return cache
}

func (c *InMemoryCache) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for key, item := range c.data {
			if now.After(item.expiration) {
				delete(c.data, key)
			}
		}
		c.mu.Unlock()
	}
}

func (c *InMemoryCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.data[key]
	if !exists {
		return nil, false
	}
	if time.Now().After(item.expiration) {
		delete(c.data, key)
		return nil, false
	}
	return item.value, true
}

func (c *InMemoryCache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = cacheItem{
		value:      value,
		expiration: time.Now().Add(ttl),
	}
}

func (c *InMemoryCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
}

// 3.6 Struct Embedding برای اشتراک‌گذاری وابستگی‌ها
type BaseService struct {
	Logger Logger
	Cache  Cache
}

type EnhancedUserService struct {
	BaseService // Embedding
	userRepo    UserRepository
}

func NewEnhancedUserService(userRepo UserRepository, logger Logger, cache Cache) *EnhancedUserService {
	return &EnhancedUserService{
		BaseService: BaseService{Logger: logger, Cache: cache},
		userRepo:    userRepo,
	}
}

func (s *EnhancedUserService) GetUserWithCache(ctx context.Context, id string) (*User, error) {
	// ابتدا از کش می‌خوانیم
	cacheKey := "user:" + id
	if cached, found := s.Cache.Get(cacheKey); found {
		s.Logger.Debug("Cache hit", "key", cacheKey)
		if user, ok := cached.(User); ok {
			return &user, nil
		}
	}

	// از دیتابیس می‌خوانیم
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		s.Logger.Error("Failed to get user", "error", err, "id", id)
		return nil, err
	}

	// ذخیره در کش
	s.Cache.Set(cacheKey, *user, 5*time.Minute)
	s.Logger.Debug("Cache set", "key", cacheKey)

	return user, nil
}

// ============================================================================
// بخش 4: Options Pattern
// ============================================================================

// 4.1 تنظیمات سرور با Options Pattern
type ServerConfig struct {
	Host           string
	Port           int
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	MaxHeaderBytes int
	EnableTLS      bool
	CertFile       string
	KeyFile        string
	Middleware     []func(http.Handler) http.Handler
}

// Option تابع تنظیم‌کننده
type Option func(*ServerConfig)

// Options پیش‌فرض
func DefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		Host:           "0.0.0.0",
		Port:           8080,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
		EnableTLS:      false,
	}
}

// WithHost تنظیم host
func WithHost(host string) Option {
	return func(c *ServerConfig) {
		c.Host = host
	}
}

// WithPort تنظیم port
func WithPort(port int) Option {
	return func(c *ServerConfig) {
		c.Port = port
	}
}

// WithTimeout تنظیم timeouts
func WithTimeout(read, write time.Duration) Option {
	return func(c *ServerConfig) {
		c.ReadTimeout = read
		c.WriteTimeout = write
	}
}

// WithTLS تنظیم TLS
func WithTLS(certFile, keyFile string) Option {
	return func(c *ServerConfig) {
		c.EnableTLS = true
		c.CertFile = certFile
		c.KeyFile = keyFile
	}
}

// WithMiddleware افزودن middleware
func WithMiddleware(middleware ...func(http.Handler) http.Handler) Option {
	return func(c *ServerConfig) {
		c.Middleware = append(c.Middleware, middleware...)
	}
}

// NewServer ایجاد سرور با استفاده از options
func NewServer(handler http.Handler, opts ...Option) *http.Server {
	config := DefaultServerConfig()
	for _, opt := range opts {
		opt(config)
	}

	// اعمال middlewareها
	finalHandler := handler
	for i := len(config.Middleware) - 1; i >= 0; i-- {
		finalHandler = config.Middleware[i](finalHandler)
	}

	return &http.Server{
		Addr:           fmt.Sprintf("%s:%d", config.Host, config.Port),
		Handler:        finalHandler,
		ReadTimeout:    config.ReadTimeout,
		WriteTimeout:   config.WriteTimeout,
		MaxHeaderBytes: config.MaxHeaderBytes,
	}
}

// 4.2 Options Pattern برای توابع (Functional Options)
type Connection struct {
	url      string
	timeout  time.Duration
	retries  int
	useSSL   bool
	username string
	password string
}

type ConnectionOption func(*Connection)

func WithTimeout(timeout time.Duration) ConnectionOption {
	return func(c *Connection) {
		c.timeout = timeout
	}
}

func WithRetries(retries int) ConnectionOption {
	return func(c *Connection) {
		c.retries = retries
	}
}

func WithSSL(useSSL bool) ConnectionOption {
	return func(c *Connection) {
		c.useSSL = useSSL
	}
}

func WithAuth(username, password string) ConnectionOption {
	return func(c *Connection) {
		c.username = username
		c.password = password
	}
}

func NewConnection(url string, opts ...ConnectionOption) *Connection {
	conn := &Connection{
		url:     url,
		timeout: 30 * time.Second,
		retries: 3,
		useSSL:  false,
	}

	for _, opt := range opts {
		opt(conn)
	}

	return conn
}

// ============================================================================
// بخش 5: Middleware Pattern
// ============================================================================

// 5.1 Middleware تعریف
type Middleware func(http.Handler) http.Handler

// 5.2 Chain کردن Middlewareها
func Chain(handler http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

// 5.3 Logging Middleware
func LoggingMiddleware(logger Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := &responseWriter{ResponseWriter: w, status: http.StatusOK}

			next.ServeHTTP(ww, r)

			logger.Info("HTTP Request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", ww.status,
				"duration", time.Since(start),
				"ip", r.RemoteAddr,
			)
		})
	}
}

// responseWriter برای گرفتن status code
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// 5.4 Recovery Middleware
func RecoveryMiddleware(logger Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error("Panic recovered", "error", err, "path", r.URL.Path)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// 5.5 Authentication Middleware
func AuthMiddleware(requiredRole string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if token == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// اعتبارسنجی token (ساده برای مثال)
			// در واقع باید JWT یا token واقعی بررسی شود
			if token != "Bearer valid-token" {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// 5.6 Rate Limit Middleware
func RateLimitMiddleware(limit int, window time.Duration) Middleware {
	var mu sync.Mutex
	requests := make(map[string][]time.Time)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr

			mu.Lock()
			defer mu.Unlock()

			// پاک کردن درخواست‌های قدیمی
			now := time.Now()
			cutoff := now.Add(-window)
			var validRequests []time.Time
			for _, t := range requests[ip] {
				if t.After(cutoff) {
					validRequests = append(validRequests, t)
				}
			}
			requests[ip] = validRequests

			if len(requests[ip]) >= limit {
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}

			requests[ip] = append(requests[ip], now)
			next.ServeHTTP(w, r)
		})
	}
}

// 5.7 CORS Middleware
func CORSMiddleware(allowedOrigins []string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			for _, allowed := range allowedOrigins {
				if allowed == "*" || allowed == origin {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					break
				}
			}

			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Max-Age", "86400")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// 5.8 Request ID Middleware
type contextKey string

const RequestIDKey contextKey = "request_id"

func RequestIDMiddleware() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = generateID()
			}
			w.Header().Set("X-Request-ID", requestID)
			ctx := context.WithValue(r.Context(), RequestIDKey, requestID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// ============================================================================
// بخش 6: Worker Pattern
// ============================================================================

// 6.1 Job و Worker
type Job struct {
	ID      string
	Payload interface{}
}

type Result struct {
	JobID string
	Data  interface{}
	Error error
}

// 6.2 Worker Pool
type WorkerPool struct {
	jobQueue    chan Job
	resultQueue chan Result
	workerCount int
	wg          sync.WaitGroup
	stopCh      chan struct{}
	stopped     bool
	mu          sync.Mutex
}

// NewWorkerPool ایجاد Worker Pool جدید
func NewWorkerPool(workerCount int, queueSize int) *WorkerPool {
	return &WorkerPool{
		jobQueue:    make(chan Job, queueSize),
		resultQueue: make(chan Result, queueSize),
		workerCount: workerCount,
		stopCh:      make(chan struct{}),
	}
}

// Start راه‌اندازی Worker Pool
func (wp *WorkerPool) Start() {
	for i := 0; i < wp.workerCount; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

// worker تابع اجرای هر worker
func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()

	for {
		select {
		case job, ok := <-wp.jobQueue:
			if !ok {
				return
			}
			result := wp.processJob(job)
			wp.resultQueue <- result

		case <-wp.stopCh:
			return
		}
	}
}

// processJob پردازش یک Job
func (wp *WorkerPool) processJob(job Job) Result {
	// شبیه‌سازی پردازش سنگین
	time.Sleep(100 * time.Millisecond)

	// پردازش واقعی
	var data interface{}
	var err error

	switch v := job.Payload.(type) {
	case string:
		data = fmt.Sprintf("Processed: %s", v)
	case int:
		data = v * v
	default:
		err = fmt.Errorf("unsupported payload type")
	}

	return Result{
		JobID: job.ID,
		Data:  data,
		Error: err,
	}
}

// Submit ارسال Job به Pool
func (wp *WorkerPool) Submit(job Job) error {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if wp.stopped {
		return errors.New("worker pool is stopped")
	}

	select {
	case wp.jobQueue <- job:
		return nil
	default:
		return errors.New("job queue is full")
	}
}

// SubmitWithTimeout ارسال Job با تایم‌اوت
func (wp *WorkerPool) SubmitWithTimeout(job Job, timeout time.Duration) error {
	select {
	case wp.jobQueue <- job:
		return nil
	case <-time.After(timeout):
		return errors.New("submit timeout")
	}
}

// Results دریافت کانال نتایج
func (wp *WorkerPool) Results() <-chan Result {
	return wp.resultQueue
}

// Stop توقف Worker Pool
func (wp *WorkerPool) Stop() {
	wp.mu.Lock()
	if wp.stopped {
		wp.mu.Unlock()
		return
	}
	wp.stopped = true
	wp.mu.Unlock()

	close(wp.stopCh)
	close(wp.jobQueue)
	wp.wg.Wait()
	close(wp.resultQueue)
}

// 6.3 Worker Pool با Context (قابل کنسل شدن)
type WorkerPoolWithContext struct {
	jobQueue chan Job
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
}

func NewWorkerPoolWithContext(workerCount int, queueSize int) *WorkerPoolWithContext {
	ctx, cancel := context.WithCancel(context.Background())
	wp := &WorkerPoolWithContext{
		jobQueue: make(chan Job, queueSize),
		ctx:      ctx,
		cancel:   cancel,
	}

	for i := 0; i < workerCount; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}

	return wp
}

func (wp *WorkerPoolWithContext) worker(id int) {
	defer wp.wg.Done()

	for {
		select {
		case job, ok := <-wp.jobQueue:
			if !ok {
				return
			}
			wp.processJob(job)
		case <-wp.ctx.Done():
			return
		}
	}
}

func (wp *WorkerPoolWithContext) processJob(job Job) {
	// پردازش با context
	select {
	case <-wp.ctx.Done():
		return
	default:
		time.Sleep(50 * time.Millisecond) // شبیه‌سازی کار
	}
}

func (wp *WorkerPoolWithContext) Submit(job Job) error {
	select {
	case wp.jobQueue <- job:
		return nil
	case <-wp.ctx.Done():
		return wp.ctx.Err()
	}
}

func (wp *WorkerPoolWithContext) Stop() {
	wp.cancel()
	close(wp.jobQueue)
	wp.wg.Wait()
}

// ============================================================================
// بخش 7: مثال کامل استفاده از همه الگوها
// ============================================================================

// 7.1 HTTP Handler با استفاده از Service
type UserHandler struct {
	userService UserService
	logger      Logger
}

func NewUserHandler(userService UserService, logger Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
	}
}

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Age   int    `json:"age"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode request", "error", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user, err := h.userService.Register(r.Context(), req.Name, req.Email, req.Age)
	if err != nil {
		h.logger.Error("Failed to register user", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("id")
	if userID == "" {
		http.Error(w, "user id required", http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetProfile(r.Context(), userID)
	if err != nil {
		if err.Error() == "user not found" {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// 7.2 راه‌اندازی کامل برنامه
func setupApplication() {
	// ایجاد وابستگی‌ها
	userRepo := NewInMemoryUserRepository()
	logger := NewSimpleLogger("APP")
	cache := NewInMemoryCache()

	// Service با Dependency Injection
	userService := NewUserService(userRepo)

	// Handler
	userHandler := NewUserHandler(userService, logger)

	// Middlewareها
	middlewares := []Middleware{
		RequestIDMiddleware(),
		RecoveryMiddleware(logger),
		LoggingMiddleware(logger),
		CORSMiddleware([]string{"*"}),
		RateLimitMiddleware(100, time.Minute),
	}

	// Router
	mux := http.NewServeMux()
	mux.HandleFunc("POST /users", userHandler.RegisterUser)
	mux.HandleFunc("GET /users", userHandler.GetUser)

	// اعمال middlewareها
	handler := Chain(mux, middlewares...)

	// ایجاد سرور با Options Pattern
	server := NewServer(handler,
		WithHost("localhost"),
		WithPort(8080),
		WithTimeout(30*time.Second, 30*time.Second),
	)

	// Worker Pool برای پردازش پس‌زمینه
	workerPool := NewWorkerPool(5, 100)
	workerPool.Start()

	// شروع سرور
	go func() {
		log.Println("Server starting on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// منتظر سیگنال برای shutdown
	// (در مثال واقعی از signal استفاده می‌شود)
	time.Sleep(5 * time.Second)
	log.Println("Shutting down...")

	workerPool.Stop()
	server.Shutdown(context.Background())
}

// ============================================================================
// بخش 8: Helper Functions
// ============================================================================

func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// ============================================================================
// بخش 9: Main
// ============================================================================

func main() {
	fmt.Println(stringsRepeat("=", 80))
	fmt.Println("🎯 GO DESIGN PATTERNS GUIDE")
	fmt.Println("Repository | Service | DI | Options | Middleware | Worker")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ PATTERN SUMMARY                                                           │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│ 1. Repository Pattern                                                      │
│    • جدا کردن منطق دسترسی به داده از منطق کسب و کار                         │
│    • امکان تعویض دیتابیس بدون تغییر سرویس                                   │
│    • تست راحت‌تر با mock repository                                         │
│                                                                             │
│ 2. Service Layer                                                           │
│    • encapsulate کردن منطق کسب و کار                                        │
│    • استفاده از repositoryها برای دسترسی به داده                            │
│    • مدیریت تراکنش‌ها و اعتبارسنجی                                          │
│                                                                             │
│ 3. Dependency Injection                                                    │
│    • وابستگی‌ها را از بیرون تزریق کنید (نه داخل)                           │
│    • استفاده از interfaceها برای انعطاف‌پذیری                               │
│    • استفاده از struct embedding برای اشتراک‌گذاری وابستگی‌ها               │
│                                                                             │
│ 4. Options Pattern                                                         │
│    • پیکربندی انعطاف‌پذیر توابع و structها                                  │
│    • جلوگیری از explosion پارامترها                                         │
│    • مقادیر پیش‌فرض منطقی                                                   │
│                                                                             │
│ 5. Middleware Pattern                                                      │
│    • پردازش درخواست‌ها در یک زنجیره                                         │
│    • قابلیت ترکیب و reuse                                                  │
│    • جداسازی cross-cutting concerns                                        │
│                                                                             │
│ 6. Worker Pattern                                                          │
│    • پردازش موازی کارها با تعداد محدود worker                              │
│    • کنترل همزمانی و جلوگیری از overwhelm                                  │
│    • صف‌بندی و مدیریت نتایج                                                │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

💡 GOLDEN RULES:

   1. Repository: Interface را در سمت مصرف‌کننده تعریف کن
   2. Service: سرویس‌ها stateless باشند (بدون state داخلی)
   3. DI: از constructor injection استفاده کن (نه field injection)
   4. Options: مقادیر پیش‌فرض منطقی داشته باش
   5. Middleware: ترتیب middlewareها مهم است
   6. Worker: همیشه راهی برای stop کردن worker pool داشته باش
`)

	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🎯 GO DESIGN PATTERNS - COMPLETE")
	fmt.Println("Ready to build maintainable Go applications!")
	fmt.Println(stringsRepeat("=", 80))
}

// تابع کمکی برای تکرار رشته
func stringsRepeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}

// برای کامپایل شدن (در واقع باید sql را import کنید)
type sqlDB struct{}
type sqlTx struct{}
type sqlDBInterface interface{}

// NOTE: در فایل واقعی باید import "database/sql" داشته باشید
var _ = sqlDB{}
