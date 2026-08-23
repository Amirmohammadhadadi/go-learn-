// ============================================================================
// FILE: solid_go_guide.go
// TITLE: راهنمای کامل SOLID Principles در Go - با تأکید بر Interfaceهای کوچک
// HOW TO RUN: go run solid_go_guide.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - SOLID چیست؟
// ============================================================================
//
// SOLID پنج اصل اساسی طراحی شیء‌گرا هستند که توسط Robert C. Martin معرفی شده‌اند:
//
// S - Single Responsibility Principle (SRP) - اصل مسئولیت واحد
// O - Open/Closed Principle (OCP) - اصل باز/بسته
// L - Liskov Substitution Principle (LSP) - اصل جایگزینی لیسکوف
// I - Interface Segregation Principle (ISP) - اصل جداسازی اینترفیس
// D - Dependency Inversion Principle (DIP) - اصل وارونگی وابستگی
//
// در Go، این اصول کمی متفاوت پیاده‌سازی می‌شوند:
// - به جای کلاس، از struct استفاده می‌کنیم
// - به جای ارث‌بری، از composition استفاده می‌کنیم
// - اینترفیس‌ها به صورت implicit پیاده‌سازی می‌شوند
// - تأکید بر اینترفیس‌های کوچک (1-3 متد)
//
// قانون طلایی:
// "اینترفیس‌ها را کوچک نگه دار (هر چه کوچک‌تر، بهتر).
//  از composition به جای inheritance استفاده کن.
//  وابستگی‌ها را از طریق اینترفیس تزریق کن."
// ============================================================================

package __solid

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

// ============================================================================
// بخش 1: Single Responsibility Principle (SRP)
// ============================================================================
//
// هر struct باید فقط یک دلیل برای تغییر داشته باشد.
// یعنی هر struct فقط یک مسئولیت را بر عهده بگیرد.
//
// ❌ بد: یک struct چند کار مختلف انجام می‌دهد
// ✅ خوب: هر struct فقط یک کار انجام می‌دهد

// ============================================
// ❌ مثال بد: UserService چند مسئولیت دارد
// ============================================
type BadUserService struct {
	db     *sql.DB
	cache  *redis.Client
	logger *log.Logger
}

// این struct چند کار مختلف انجام می‌دهد:
// 1. منطق کسب و کار
// 2. دسترسی به دیتابیس
// 3. کش کردن
// 4. لاگینگ
// 5. ارسال ایمیل
func (s *BadUserService) RegisterUser(name, email string) error {
	// اعتبارسنجی
	if name == "" || email == "" {
		s.logger.Println("Validation failed")
		return errors.New("invalid input")
	}

	// ذخیره در دیتابیس
	_, err := s.db.Exec("INSERT INTO users (name, email) VALUES ($1, $2)", name, email)
	if err != nil {
		s.logger.Printf("DB error: %v", err)
		return err
	}

	// ذخیره در کش
	err = s.cache.Set(context.Background(), "user:"+email, name, time.Hour).Err()
	if err != nil {
		s.logger.Printf("Cache error: %v", err)
	}

	// ارسال ایمیل خوش‌آمدگویی
	if err := s.sendWelcomeEmail(email); err != nil {
		s.logger.Printf("Email error: %v", err)
	}

	return nil
}

func (s *BadUserService) sendWelcomeEmail(email string) error {
	// منطق ارسال ایمیل
	return nil
}

// ============================================
// ✅ خوب: تفکیک مسئولیت‌ها
// ============================================

// 1. Entity - فقط داده و منطق دامنه
type User struct {
	ID        string
	Name      string
	Email     string
	CreatedAt time.Time
}

func (u *User) Validate() error {
	if u.Name == "" {
		return errors.New("name is required")
	}
	if u.Email == "" {
		return errors.New("email is required")
	}
	return nil
}

// 2. Repository - فقط دسترسی به داده
type UserRepository interface {
	Save(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
}

// 3. Cache - فقط کش کردن
type UserCache interface {
	Get(ctx context.Context, key string) (*User, error)
	Set(ctx context.Context, key string, user *User, ttl time.Duration) error
}

// 4. EmailService - فقط ارسال ایمیل
type EmailService interface {
	SendWelcomeEmail(ctx context.Context, email string) error
}

// 5. Logger - فقط لاگینگ
type Logger interface {
	Info(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
}

// 6. UserService - فقط منطق ثبت‌نام (هماهنگ‌کننده)
type UserService struct {
	repo   UserRepository
	cache  UserCache
	email  EmailService
	logger Logger
}

func NewUserService(repo UserRepository, cache UserCache, email EmailService, logger Logger) *UserService {
	return &UserService{
		repo:   repo,
		cache:  cache,
		email:  email,
		logger: logger,
	}
}

func (s *UserService) RegisterUser(ctx context.Context, name, email string) error {
	user := &User{
		ID:        generateID(),
		Name:      name,
		Email:     email,
		CreatedAt: time.Now(),
	}

	// اعتبارسنجی
	if err := user.Validate(); err != nil {
		s.logger.Error("Validation failed", "error", err)
		return err
	}

	// ذخیره در دیتابیس
	if err := s.repo.Save(ctx, user); err != nil {
		s.logger.Error("Failed to save user", "error", err)
		return err
	}

	// ذخیره در کش (اختیاری)
	go func() {
		if err := s.cache.Set(context.Background(), "user:"+email, user, time.Hour); err != nil {
			s.logger.Error("Failed to cache user", "error", err)
		}
	}()

	// ارسال ایمیل (غیرهمزمان)
	go func() {
		if err := s.email.SendWelcomeEmail(context.Background(), email); err != nil {
			s.logger.Error("Failed to send welcome email", "error", err)
		}
	}()

	s.logger.Info("User registered successfully", "email", email)
	return nil
}

// ============================================================================
// بخش 2: Open/Closed Principle (OCP)
// ============================================================================
//
// کلاس‌ها باید برای گسترش باز، ولی برای تغییر بسته باشند.
// یعنی باید بتوانیم قابلیت‌های جدید اضافه کنیم بدون اینکه کد موجود را تغییر دهیم.
//
// ❌ بد: با اضافه کردن نوع جدید، کد موجود تغییر می‌کند
// ✅ خوب: از اینترفیس و composition برای گسترش استفاده می‌کنیم

// ============================================
// ❌ مثال بد: هر بار با نوع جدید، کد تغییر می‌کند
// ============================================
type BadPaymentProcessor struct{}

func (p *BadPaymentProcessor) ProcessPayment(method string, amount float64) error {
	switch method {
	case "credit_card":
		// منطق کارت اعتباری
		return nil
	case "paypal":
		// منطق PayPal
		return nil
	case "crypto":
		// منطق ارز دیجیتال
		return nil
	default:
		return errors.New("unsupported payment method")
	}
}

// ============================================
// ✅ خوب: باز برای گسترش، بسته برای تغییر
// ============================================

// PaymentMethod اینترفیس برای همه روش‌های پرداخت
type PaymentMethod interface {
	Pay(amount float64) error
	GetName() string
}

// CreditCardPayment
type CreditCardPayment struct {
	CardNumber string
	CVV        string
	ExpiryDate string
}

func (p *CreditCardPayment) Pay(amount float64) error {
	fmt.Printf("Paying %.2f via Credit Card (%s)\n", amount, p.CardNumber)
	// منطق واقعی پرداخت
	return nil
}

func (p *CreditCardPayment) GetName() string {
	return "Credit Card"
}

// PayPalPayment
type PayPalPayment struct {
	Email string
}

func (p *PayPalPayment) Pay(amount float64) error {
	fmt.Printf("Paying %.2f via PayPal (%s)\n", amount, p.Email)
	return nil
}

func (p *PayPalPayment) GetName() string {
	return "PayPal"
}

// CryptoPayment (نوع جدید - بدون تغییر در کد موجود)
type CryptoPayment struct {
	WalletAddress string
	Currency      string
}

func (p *CryptoPayment) Pay(amount float64) error {
	fmt.Printf("Paying %.2f %s via Crypto (%s)\n", amount, p.Currency, p.WalletAddress)
	return nil
}

func (p *CryptoPayment) GetName() string {
	return "Cryptocurrency"
}

// PaymentProcessor (بسته برای تغییر)
type PaymentProcessor struct {
	paymentMethods map[string]PaymentMethod
}

func NewPaymentProcessor() *PaymentProcessor {
	return &PaymentProcessor{
		paymentMethods: make(map[string]PaymentMethod),
	}
}

func (p *PaymentProcessor) RegisterMethod(method PaymentMethod) {
	p.paymentMethods[method.GetName()] = method
}

func (p *PaymentProcessor) ProcessPayment(methodName string, amount float64) error {
	method, exists := p.paymentMethods[methodName]
	if !exists {
		return fmt.Errorf("payment method %s not supported", methodName)
	}
	return method.Pay(amount)
}

// ============================================================================
// بخش 3: Liskov Substitution Principle (LSP)
// ============================================================================
//
// اشیاء یک برنامه باید قابل جایگزینی با نمونه‌های زیرنوع خود باشند
// بدون اینکه correctness برنامه تحت تأثیر قرار گیرد.
//
// ❌ بد: زیرنوع رفتار نوع پایه را تغییر می‌دهد
// ✅ خوب: زیرنوع می‌تواند جایگزین نوع پایه شود

// ============================================
// ❌ مثال بد: Square نمی‌تواند جایگزین Rectangle شود
// ============================================
type Rectangle struct {
	Width  int
	Height int
}

func (r *Rectangle) SetWidth(width int) {
	r.Width = width
}

func (r *Rectangle) SetHeight(height int) {
	r.Height = height
}

func (r *Rectangle) Area() int {
	return r.Width * r.Height
}

type Square struct {
	Rectangle // embedding - ❌ نقض LSP
}

func (s *Square) SetWidth(width int) {
	s.Width = width
	s.Height = width
}

func (s *Square) SetHeight(height int) {
	s.Width = height
	s.Height = height
}

// تابعی که با Rectangle کار می‌کند
func processRectangle(r *Rectangle) {
	r.SetWidth(5)
	r.SetHeight(4)
	// انتظار داریم Area = 20
	if r.Area() != 20 {
		fmt.Printf("LSP violation: expected 20, got %d\n", r.Area())
	}
}

// ============================================
// ✅ خوب: طراحی بهتر با اینترفیس
// ============================================

// Shape اینترفیس (کوچک)
type Shape interface {
	Area() float64
}

// Rectangle (پیاده‌سازی Shape)
type GoodRectangle struct {
	Width  float64
	Height float64
}

func (r *GoodRectangle) Area() float64 {
	return r.Width * r.Height
}

// Square (پیاده‌سازی Shape)
type GoodSquare struct {
	Side float64
}

func (s *GoodSquare) Area() float64 {
	return s.Side * s.Side
}

// تابعی که با Shape کار می‌کند (با هر Shapeای کار می‌کند)
func printArea(s Shape) {
	fmt.Printf("Area: %.2f\n", s.Area())
}

// ============================================================================
// بخش 4: Interface Segregation Principle (ISP)
// ============================================================================
//
// اینترفیس‌ها را کوچک نگه دار. هیچ کلاسی نباید مجبور باشد
// متدهایی را پیاده‌سازی کند که استفاده نمی‌کند.
//
// در Go، این اصل بسیار مهم است. اینترفیس‌های کوچک (1-3 متد) بهترین practice هستند.
//
// ❌ بد: اینترفیس بزرگ با متدهای زیاد
// ✅ خوب: چند اینترفیس کوچک

// ============================================
// ❌ مثال بد: اینترفیس بزرگ
// ============================================
type BigInterface interface {
	Save() error
	Delete() error
	Find() error
	Update() error
	Validate() error
	SendEmail() error
	GenerateReport() error
	ExportCSV() error
}

// هر struct باید همه متدها را پیاده‌سازی کند (حتی اگر نیازی نداشته باشد)

// ============================================
// ✅ خوب: اینترفیس‌های کوچک
// ============================================

// اینترفیس‌های کوچک و متمرکز
type Saver interface {
	Save() error
}

type Finder interface {
	Find(id string) (interface{}, error)
}

type Deleter interface {
	Delete(id string) error
}

type Validator interface {
	Validate() error
}

type EmailSender interface {
	SendEmail(to, subject, body string) error
}

// ترکیب اینترفیس‌ها (در صورت نیاز)
type Repository interface {
	Saver
	Finder
	Deleter
}

// هر struct فقط اینترفیس‌هایی که نیاز دارد را پیاده‌سازی می‌کند
type UserRepository struct{}

func (r *UserRepository) Save() error {
	// پیاده‌سازی
	return nil
}

func (r *UserRepository) Find(id string) (interface{}, error) {
	// پیاده‌سازی
	return nil, nil
}

func (r *UserRepository) Delete(id string) error {
	// پیاده‌سازی
	return nil
}

// Service فقط به اینترفیس‌های مورد نیازش وابسته است
type UserServiceISP struct {
	saver  Saver
	finder Finder
	mailer EmailSender
}

func NewUserServiceISP(saver Saver, finder Finder, mailer EmailSender) *UserServiceISP {
	return &UserServiceISP{
		saver:  saver,
		finder: finder,
		mailer: mailer,
	}
}

// ============================================================================
// بخش 5: Dependency Inversion Principle (DIP)
// ============================================================================
//
// 1. ماژول‌های سطح بالا نباید به ماژول‌های سطح پایین وابسته باشند.
//    هر دو باید به abstractions وابسته باشند.
// 2. Abstractions نباید به details وابسته باشند.
//    Details باید به abstractions وابسته باشند.
//
// در Go، این یعنی وابستگی‌ها را از طریق اینترفیس تزریق کنیم.
//
// ❌ بد: وابستگی مستقیم به پیاده‌سازی concrete
// ✅ خوب: وابستگی به اینترفیس

// ============================================
// ❌ مثال بد: وابستگی مستقیم
// ============================================

type BadNotificationService struct {
	// وابستگی مستقیم به پیاده‌سازی concrete
	emailClient *EmailClient
	smsClient   *SMSClient
}

func NewBadNotificationService() *BadNotificationService {
	return &BadNotificationService{
		emailClient: &EmailClient{},
		smsClient:   &SMSClient{},
	}
}

func (s *BadNotificationService) Notify(userID string, message string) error {
	// وابستگی به concrete implementations
	if err := s.emailClient.Send(userID, message); err != nil {
		return err
	}
	return s.smsClient.Send(userID, message)
}

// ============================================
// ✅ خوب: وابستگی به اینترفیس
// ============================================

// اینترفیس سطح بالا
type Notifier interface {
	Send(recipient string, message string) error
}

// ماژول سطح بالا (Service) فقط به اینترفیس وابسته است
type NotificationService struct {
	notifiers []Notifier // وابستگی به abstraction
}

func NewNotificationService(notifiers ...Notifier) *NotificationService {
	return &NotificationService{
		notifiers: notifiers,
	}
}

func (s *NotificationService) Notify(recipient string, message string) error {
	for _, notifier := range s.notifiers {
		if err := notifier.Send(recipient, message); err != nil {
			return fmt.Errorf("failed to send via %T: %w", notifier, err)
		}
	}
	return nil
}

// ماژول‌های سطح پایین (Details) اینترفیس را پیاده‌سازی می‌کنند
type EmailNotifier struct {
	SMTPHost string
	SMTPPort int
}

func (n *EmailNotifier) Send(recipient string, message string) error {
	fmt.Printf("Sending email to %s: %s\n", recipient, message)
	return nil
}

type SMSNotifier struct {
	APIKey    string
	APISecret string
}

func (n *SMSNotifier) Send(recipient string, message string) error {
	fmt.Printf("Sending SMS to %s: %s\n", recipient, message)
	return nil
}

type PushNotifier struct {
	AppID     string
	AppSecret string
}

func (n *PushNotifier) Send(recipient string, message string) error {
	fmt.Printf("Sending push notification to %s: %s\n", recipient, message)
	return nil
}

// ============================================================================
// بخش 6: تأکید بر Interfaceهای کوچک در Go
// ============================================================================

// در Go، بهترین practice این است که اینترفیس‌ها را بسیار کوچک نگه داریم.
// بسیاری از اینترفیس‌های استاندارد Go فقط 1-2 متد دارند:

// مثال از اینترفیس‌های استاندارد Go (کوچک):
// type Reader interface { Read(p []byte) (n int, err error) }
// type Writer interface { Write(p []byte) (n int, err error) }
// type Closer interface { Close() error }
// type Stringer interface { String() string }
// type Error interface { Error() string }

// 6.1 مزایای اینترفیس‌های کوچک:

// 1. قابلیت reuse بالا
// 2. تست آسان‌تر
// 3. پیاده‌سازی ساده‌تر
// 4. ترکیب‌پذیری (composability) بالا

// 6.2 مثال: اینترفیس‌های کوچک برای یک API

// اینترفیس‌های کوچک و متمرکز
type Creator interface {
	Create(ctx context.Context, data interface{}) (interface{}, error)
}

type Reader interface {
	Read(ctx context.Context, id string) (interface{}, error)
}

type Updater interface {
	Update(ctx context.Context, id string, data interface{}) error
}

type Deleter interface {
	Delete(ctx context.Context, id string) error
}

type Lister interface {
	List(ctx context.Context, filter interface{}) ([]interface{}, error)
}

// ترکیب اینترفیس‌ها برای نیازهای خاص
type CRUDRepository interface {
	Creator
	Reader
	Updater
	Deleter
}

type FullRepository interface {
	CRUDRepository
	Lister
}

// 6.3 مثال: اینترفیس‌های کوچک برای کش

type CacheGetter interface {
	Get(key string) (interface{}, bool)
}

type CacheSetter interface {
	Set(key string, value interface{}, ttl time.Duration)
}

type CacheDeleter interface {
	Delete(key string)
}

// ترکیب برای نیازهای مختلف
type Cache interface {
	CacheGetter
	CacheSetter
	CacheDeleter
}

// 6.4 مثال: اینترفیس‌های کوچک برای لاگینگ

type LoggerInfo interface {
	Info(msg string, keysAndValues ...interface{})
}

type LoggerError interface {
	Error(msg string, keysAndValues ...interface{})
}

type LoggerDebug interface {
	Debug(msg string, keysAndValues ...interface{})
}

// ترکیب
type Logger interface {
	LoggerInfo
	LoggerError
	LoggerDebug
}

// ============================================================================
// بخش 7: Anti-Patterns در Go
// ============================================================================

func antiPatterns() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("⚠️ COMMON SOLID ANTI-PATTERNS IN GO")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ 1. God Object (نقض SRP)                                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│    type GodService struct { ... } // همه چیز در یک struct                  │
│    ✅ هر struct را به چند struct کوچک با مسئولیت واحد تقسیم کن            │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 2. Switch on Type (نقض OCP)                                              │
├─────────────────────────────────────────────────────────────────────────────┤
│    switch v := data.(type) { ... } // هر نوع جدید نیاز به تغییر دارد       │
│    ✅ از اینترفیس و polymorphism استفاده کن                                │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 3. Fat Interface (نقض ISP)                                               │
├─────────────────────────────────────────────────────────────────────────────┤
│    type BigInterface interface { A(); B(); C(); D(); E() }                │
│    ✅ اینترفیس را به چند اینترفیس کوچک تقسیم کن                           │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 4. Leaky Abstraction (نقض DIP)                                           │
├─────────────────────────────────────────────────────────────────────────────┤
│    type Service struct { db *postgres.DB } // وابستگی مستقیم به دیتابیس   │
│    ✅ به اینترفیس وابسته باش: type DB interface { ... }                   │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 5. Interface on Producer Side (نقض ISP)                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│    // بسته تولیدکننده اینترفیس را تعریف می‌کند                            │
│    type UserRepository interface { ... }                                  │
│    ✅ اینترفیس را در سمت مصرف‌کننده (consumer) تعریف کن                    │
└─────────────────────────────────────────────────────────────────────────────┘
`)
}

// ============================================================================
// بخش 8: Practical Example - Complete SOLID Implementation
// ============================================================================

// مثال کامل از یک سیستم سفارش‌گیری که اصول SOLID را رعایت می‌کند

// 8.1 Domain Entities (SRP - فقط داده و منطق دامنه)
type OrderItem struct {
	ProductID string
	Quantity  int
	Price     float64
}

func (i *OrderItem) Total() float64 {
	return float64(i.Quantity) * i.Price
}

type OrderStatus string

const (
	StatusPending   OrderStatus = "pending"
	StatusConfirmed OrderStatus = "confirmed"
	StatusShipped   OrderStatus = "shipped"
	StatusDelivered OrderStatus = "delivered"
	StatusCancelled OrderStatus = "cancelled"
)

type Order struct {
	ID        string
	UserID    string
	Items     []OrderItem
	Status    OrderStatus
	Total     float64
	CreatedAt time.Time
}

func (o *Order) CalculateTotal() {
	var total float64
	for _, item := range o.Items {
		total += item.Total()
	}
	o.Total = total
}

func (o *Order) CanBeCancelled() bool {
	return o.Status == StatusPending || o.Status == StatusConfirmed
}

// 8.2 اینترفیس‌های کوچک (ISP)

// Repository اینترفیس‌ها
type OrderSaver interface {
	SaveOrder(ctx context.Context, order *Order) error
}

type OrderFinder interface {
	FindOrderByID(ctx context.Context, id string) (*Order, error)
	FindOrdersByUser(ctx context.Context, userID string) ([]Order, error)
}

type OrderUpdater interface {
	UpdateOrderStatus(ctx context.Context, id string, status OrderStatus) error
}

// ترکیب برای نیازهای خاص
type OrderRepository interface {
	OrderSaver
	OrderFinder
	OrderUpdater
}

// 8.3 اینترفیس‌های کوچک برای Notification
type OrderNotifier interface {
	NotifyOrderCreated(ctx context.Context, order *Order) error
	NotifyOrderStatusChanged(ctx context.Context, order *Order) error
}

// 8.4 پیاده‌سازی Concrete (Details)
type InMemoryOrderRepository struct {
	orders map[string]*Order
}

func NewInMemoryOrderRepository() *InMemoryOrderRepository {
	return &InMemoryOrderRepository{
		orders: make(map[string]*Order),
	}
}

func (r *InMemoryOrderRepository) SaveOrder(ctx context.Context, order *Order) error {
	r.orders[order.ID] = order
	return nil
}

func (r *InMemoryOrderRepository) FindOrderByID(ctx context.Context, id string) (*Order, error) {
	order, exists := r.orders[id]
	if !exists {
		return nil, errors.New("order not found")
	}
	return order, nil
}

func (r *InMemoryOrderRepository) FindOrdersByUser(ctx context.Context, userID string) ([]Order, error) {
	var userOrders []Order
	for _, order := range r.orders {
		if order.UserID == userID {
			userOrders = append(userOrders, *order)
		}
	}
	return userOrders, nil
}

func (r *InMemoryOrderRepository) UpdateOrderStatus(ctx context.Context, id string, status OrderStatus) error {
	order, exists := r.orders[id]
	if !exists {
		return errors.New("order not found")
	}
	order.Status = status
	return nil
}

// 8.5 Use Case (OCP - باز برای گسترش با اینترفیس‌های کوچک)
type OrderUseCase struct {
	repo     OrderRepository
	notifier OrderNotifier
	logger   Logger
}

func NewOrderUseCase(repo OrderRepository, notifier OrderNotifier, logger Logger) *OrderUseCase {
	return &OrderUseCase{
		repo:     repo,
		notifier: notifier,
		logger:   logger,
	}
}

func (uc *OrderUseCase) CreateOrder(ctx context.Context, userID string, items []OrderItem) (*Order, error) {
	// اعتبارسنجی
	if userID == "" {
		return nil, errors.New("user ID is required")
	}
	if len(items) == 0 {
		return nil, errors.New("order must have at least one item")
	}

	// ایجاد سفارش
	order := &Order{
		ID:        generateID(),
		UserID:    userID,
		Items:     items,
		Status:    StatusPending,
		CreatedAt: time.Now(),
	}
	order.CalculateTotal()

	// ذخیره
	if err := uc.repo.SaveOrder(ctx, order); err != nil {
		uc.logger.Error("Failed to save order", "error", err)
		return nil, err
	}

	// ارسال نوتیفیکیشن (اختیاری - وابستگی به اینترفیس)
	go func() {
		if err := uc.notifier.NotifyOrderCreated(context.Background(), order); err != nil {
			uc.logger.Error("Failed to send notification", "error", err)
		}
	}()

	uc.logger.Info("Order created", "order_id", order.ID, "user_id", userID)
	return order, nil
}

func (uc *OrderUseCase) CancelOrder(ctx context.Context, orderID string) error {
	order, err := uc.repo.FindOrderByID(ctx, orderID)
	if err != nil {
		return err
	}

	if !order.CanBeCancelled() {
		return errors.New("order cannot be cancelled")
	}

	if err := uc.repo.UpdateOrderStatus(ctx, orderID, StatusCancelled); err != nil {
		return err
	}

	// ارسال نوتیفیکیشن
	go func() {
		order.Status = StatusCancelled
		uc.notifier.NotifyOrderStatusChanged(context.Background(), order)
	}()

	uc.logger.Info("Order cancelled", "order_id", orderID)
	return nil
}

// 8.6 Dependency Injection (DIP)
func setupOrderSystem() {
	// لایه Repository (جزئیات سطح پایین)
	repo := NewInMemoryOrderRepository()

	// لایه Notifier (جزئیات سطح پایین)
	notifier := &ConsoleNotifier{}

	// لایه Logger (جزئیات سطح پایین)
	logger := &ConsoleLogger{}

	// لایه UseCase (سطح بالا) - وابستگی‌ها از طریق اینترفیس تزریق می‌شوند
	orderUseCase := NewOrderUseCase(repo, notifier, logger)

	// استفاده
	ctx := context.Background()
	items := []OrderItem{
		{ProductID: "prod-1", Quantity: 2, Price: 25.00},
		{ProductID: "prod-2", Quantity: 1, Price: 50.00},
	}

	order, _ := orderUseCase.CreateOrder(ctx, "user-123", items)
	fmt.Printf("Order created: %s, Total: %.2f\n", order.ID, order.Total)

	orderUseCase.CancelOrder(ctx, order.ID)
	fmt.Println("Order cancelled")
}

type ConsoleNotifier struct{}

func (n *ConsoleNotifier) NotifyOrderCreated(ctx context.Context, order *Order) error {
	fmt.Printf("[NOTIFICATION] Order %s created\n", order.ID)
	return nil
}

func (n *ConsoleNotifier) NotifyOrderStatusChanged(ctx context.Context, order *Order) error {
	fmt.Printf("[NOTIFICATION] Order %s status changed to %s\n", order.ID, order.Status)
	return nil
}

type ConsoleLogger struct{}

func (l *ConsoleLogger) Info(msg string, keysAndValues ...interface{}) {
	fmt.Printf("[INFO] %s %v\n", msg, keysAndValues)
}

func (l *ConsoleLogger) Error(msg string, keysAndValues ...interface{}) {
	fmt.Printf("[ERROR] %s %v\n", msg, keysAndValues)
}

func (l *ConsoleLogger) Debug(msg string, keysAndValues ...interface{}) {
	fmt.Printf("[DEBUG] %s %v\n", msg, keysAndValues)
}

// ============================================================================
// بخش 9: Best Practices Summary
// ============================================================================

func bestPractices() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("💡 SOLID BEST PRACTICES IN GO")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ 1. INTERFACE SIZE                                                         │
├─────────────────────────────────────────────────────────────────────────────┤
│    • اینترفیس‌ها را کوچک نگه دار (1-3 متد)                                 │
│    • اینترفیس‌ها را در سمت مصرف‌کننده تعریف کن                             │
│    • از composition اینترفیس‌ها برای ساخت اینترفیس بزرگ‌تر استفاده کن       │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 2. DEPENDENCY INJECTION                                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│    • از constructor injection استفاده کن (نه field injection)             │
│    • وابستگی‌ها را از طریق اینترفیس تزریق کن                               │
│    • از functional options برای dependencies اختیاری استفاده کن           │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 3. COMPOSITION OVER INHERITANCE                                           │
├─────────────────────────────────────────────────────────────────────────────┤
│    • از struct embedding استفاده کن (نه inheritance)                       │
│    • از interface composition برای اشتراک رفتار استفاده کن                 │
│    • هرگز structها را برای reuse کد embedding نکن                          │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 4. TESTING                                                                │
├─────────────────────────────────────────────────────────────────────────────┤
│    • اینترفیس‌های کوچک = تست آسان‌تر                                       │
│    • از mockgen برای تولید mock استفاده کن                                 │
│    • هر UseCase را با mock repository تست کن                               │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 5. PACKAGE STRUCTURE                                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│    • اینترفیس‌ها را در پکیجی که از آن‌ها استفاده می‌کند تعریف کن          │
│    • پیاده‌سازی concrete را در پکیج جداگانه قرار بده                       │
│    • از internal package برای کدهای خصوصی استفاده کن                      │
└─────────────────────────────────────────────────────────────────────────────┘
`)
}

// ============================================================================
// بخش 10: Main
// ============================================================================

func main() {
	fmt.Println(stringsRepeat("=", 80))
	fmt.Println("🎯 SOLID PRINCIPLES IN GO")
	fmt.Println("With Emphasis on Small Interfaces")
	fmt.Println(stringsRepeat("=", 80))

	antiPatterns()
	bestPractices()

	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📝 EXAMPLE: Running the SOLID Order System")
	fmt.Println(stringsRepeat("=", 80))

	setupOrderSystem()

	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🎯 SOLID - COMPLETE")
	fmt.Println("Write maintainable and testable Go code!")
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

// برای رفع خطاهای undefined
type sqlDB struct{}
type redisClient struct{}

var _ = sqlDB{}
var _ = redisClient{}
