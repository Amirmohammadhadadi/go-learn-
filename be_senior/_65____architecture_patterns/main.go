// ============================================================================
// FILE: architecture_patterns_guide.go
// TITLE: راهنمای معماری پیشرفته - Monolith vs Microservices, API Gateway, Service Discovery, CQRS, Distributed Tracing
// HOW TO RUN: go run architecture_patterns_guide.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - معماری‌های مدرن
// ============================================================================
//
// انتخاب معماری مناسب یکی از مهم‌ترین تصمیمات در طراحی سیستم است.
//
// 1. Monolithic vs Microservices
//    - Monolith: همه چیز در یک برنامه
//    - Microservices: سرویس‌های کوچک و مستقل
//
// 2. API Gateway
//    - نقطه ورود واحد برای همه کلاینت‌ها
//    - مسئولیت: routing, auth, rate limiting, logging
//
// 3. Service Discovery
//    - یافتن آدرس سرویس‌ها به صورت پویا
//    - ثبت و کشف سرویس‌ها
//
// 4. CQRS (Command Query Responsibility Segregation)
//    - جداسازی عملیات write (Command) از read (Query)
//    - بهینه‌سازی جداگانه برای هر کدام
//
// 5. Distributed Tracing
//    - ردیابی درخواست در چندین سرویس
//    - شناسایی bottlenecks و خطاها
//
// قانون طلایی:
// "از Monolith شروع کن، زمانی که نیاز شد به Microservices مهاجرت کن.
//  API Gateway را همیشه در مقابل Microservices قرار بده.
//  برای درخواست‌های پیچیده از CQRS استفاده کن.
//  همیشه distributed tracing داشته باش."
// ============================================================================

package _____architecture_patterns

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

// ============================================================================
// بخش 1: Monolith vs Microservices - مقایسه و راهنمای انتخاب
// ============================================================================

func comparisonMonolithMicroservices() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🏛️ MONOLITH vs MICROSERVICES")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ CRITERION              │ MONOLITH              │ MICROSERVICES             │
├────────────────────────┼───────────────────────┼───────────────────────────┤
│ Complexity             │ Low                   │ High                      │
│ Deployment             │ Simple (single unit)  │ Complex (many services)   │
│ Scalability            │ Vertical (scale up)   │ Horizontal (scale out)    │
│ Development Speed      │ Fast (early stages)   │ Slow (initial setup)      │
│ Team Structure         │ Small team (2-5)      │ Multiple teams (per service)│
│ Technology Stack       │ Single language       │ Polyglot (多 زبان)        │
│ Fault Isolation        │ Poor (full outage)    │ Good (service isolation)  │
│ Testing                │ Easier                │ Harder (integration)      │
│ Time to Market         │ Faster                │ Slower                    │
│ Operational Overhead   │ Low                   │ High                      │
│ Data Management        │ Single DB             │ Database per service      │
│ Communication          │ In-memory calls       │ Network calls (HTTP/gRPC) │
└────────────────────────┴───────────────────────┴───────────────────────────┘

📌 WHEN TO USE MONOLITH:

   ✅ Startup / MVP
   ✅ Small team (< 5 developers)
   ✅ Simple domain
   ✅ Tight deadlines
   ✅ No need for independent scaling
   ✅ When you're not sure about boundaries

📌 WHEN TO USE MICROSERVICES:

   ✅ Large team (> 10 developers)
   ✅ Complex domain with clear boundaries
   ✅ Need independent scaling
   ✅ Different teams own different features
   ✅ Need polyglot technologies
   ✅ High availability requirements
   ✅ Can afford operational complexity

💡 RECOMMENDATION:

   Start with a well-structured monolith!
   You can always extract microservices later.
   Don't start with microservices unless you have the team and expertise.
`)
}

// ============================================================================
// بخش 2: API Gateway با Traefik
// ============================================================================

// APIGateway سرویس Gateway ساده
type APIGateway struct {
	routes     map[string]RouteConfig
	upstreams  map[string]string
	mu         sync.RWMutex
}

type RouteConfig struct {
	Path        string
	Upstream    string
	Methods     []string
	RequireAuth bool
	RateLimit   int
}

// NewAPIGateway ایجاد API Gateway جدید
func NewAPIGateway() *APIGateway {
	return &APIGateway{
		routes:    make(map[string]RouteConfig),
		upstreams: make(map[string]string),
	}
}

// AddRoute افزودن route جدید
func (g *APIGateway) AddRoute(path, upstream string, methods []string, requireAuth bool, rateLimit int) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.routes[path] = RouteConfig{
		Path:        path,
		Upstream:    upstream,
		Methods:     methods,
		RequireAuth: requireAuth,
		RateLimit:   rateLimit,
	}
}

// HandleRequest پردازش درخواست
func (g *APIGateway) HandleRequest(w http.ResponseWriter, r *http.Request) {
	// یافتن route
	g.mu.RLock()
	route, exists := g.routes[r.URL.Path]
	g.mu.RUnlock()

	if !exists {
		http.Error(w, "Route not found", http.StatusNotFound)
		return
	}

	// بررسی متد
	methodAllowed := false
	for _, m := range route.Methods {
		if m == r.Method {
			methodAllowed = true
			break
		}
	}
	if !methodAllowed {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// بررسی احراز هویت (مثال ساده)
	if route.RequireAuth {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}

	// پروکسی به upstream
	g.proxyRequest(w, r, route.Upstream)
}

func (g *APIGateway) proxyRequest(w http.ResponseWriter, r *http.Request, upstream string) {
	// ساخت URL مقصد
	url := fmt.Sprintf("http://%s%s", upstream, r.URL.Path)

	// ایجاد درخواست جدید
	req, err := http.NewRequest(r.Method, url, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	// کپی هدرها
	req.Header = r.Header

	// ارسال درخواست
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// کپی پاسخ
	for k, v := range resp.Header {
		w.Header()[k] = v
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

// RunGateway راه‌اندازی Gateway
func (g *APIGateway) RunGateway(port string) error {
	handler := http.HandlerFunc(g.HandleRequest)
	log.Printf("API Gateway running on port %s", port)
	return http.ListenAndServe(":"+port, handler)
}

// Traefik configuration example
func traefikConfigExample() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🔄 TRAEFIK CONFIGURATION (Docker)")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
# ============================================================================
# traefik.yml (Static config)
# ============================================================================

global:
  sendAnonymousUsage: false

api:
  dashboard: true
  debug: true

entryPoints:
  web:
    address: ":80"
  websecure:
    address: ":443"

providers:
  docker:
    endpoint: "unix:///var/run/docker.sock"
    exposedByDefault: false
  file:
    filename: /etc/traefik/dynamic.yml

certificatesResolvers:
  letsencrypt:
    acme:
      email: admin@example.com
      storage: /etc/traefik/acme.json
      httpChallenge:
        entryPoint: web

# ============================================================================
# docker-compose.yml with Traefik
# ============================================================================

version: '3.8'

services:
  traefik:
    image: traefik:v3.0
    command:
      - "--api.insecure=true"
      - "--providers.docker=true"
      - "--providers.docker.exposedbydefault=false"
      - "--entrypoints.web.address=:80"
      - "--entrypoints.websecure.address=:443"
    ports:
      - "80:80"
      - "443:443"
      - "8080:8080"  # Dashboard
    volumes:
      - "/var/run/docker.sock:/var/run/docker.sock:ro"
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.dashboard.rule=Host(\\`traefik.example.com\\`)"
      - "traefik.http.routers.dashboard.service=api@internal"

  user-service:
    build: ./services/user-service
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.user.rule=Host(\\`api.example.com\\`) && PathPrefix(\\`/users\\`)"
      - "traefik.http.services.user.loadbalancer.server.port=8080"

  order-service:
    build: ./services/order-service
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.order.rule=Host(\\`api.example.com\\`) && PathPrefix(\\`/orders\\`)"
      - "traefik.http.services.order.loadbalancer.server.port=8080"

  product-service:
    build: ./services/product-service
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.product.rule=Host(\\`api.example.com\\`) && PathPrefix(\\`/products\\`)"
      - "traefik.http.services.product.loadbalancer.server.port=8080"
`)
}

// ============================================================================
// بخش 3: Service Discovery (ساده با Consul)
// ============================================================================

// ServiceInstance اطلاعات یک سرویس
type ServiceInstance struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Address  string            `json:"address"`
	Port     int               `json:"port"`
	Tags     []string          `json:"tags"`
	Metadata map[string]string `json:"metadata"`
	Healthy  bool              `json:"healthy"`
	LastSeen time.Time         `json:"last_seen"`
}

// ServiceRegistry ثبت و کشف سرویس
type ServiceRegistry struct {
	services map[string][]ServiceInstance
	mu       sync.RWMutex
}

// NewServiceRegistry ایجاد registry جدید
func NewServiceRegistry() *ServiceRegistry {
	return &ServiceRegistry{
		services: make(map[string][]ServiceInstance),
	}
}

// Register ثبت سرویس
func (r *ServiceRegistry) Register(instance ServiceInstance) {
	r.mu.Lock()
	defer r.mu.Unlock()

	instance.Healthy = true
	instance.LastSeen = time.Now()

	// حذف نمونه قبلی
	instances := r.services[instance.Name]
	for i, inst := range instances {
		if inst.ID == instance.ID {
			instances = append(instances[:i], instances[i+1:]...)
			break
		}
	}

	r.services[instance.Name] = append(instances, instance)
	log.Printf("Service registered: %s (%s) at %s:%d",
		instance.Name, instance.ID, instance.Address, instance.Port)
}

// Deregister حذف سرویس
func (r *ServiceRegistry) Deregister(name, id string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	instances := r.services[name]
	for i, inst := range instances {
		if inst.ID == id {
			r.services[name] = append(instances[:i], instances[i+1:]...)
			log.Printf("Service deregistered: %s (%s)", name, id)
			break
		}
	}
}

// Discover کشف سرویس
func (r *ServiceRegistry) Discover(name string) ([]ServiceInstance, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	instances, exists := r.services[name]
	if !exists || len(instances) == 0 {
		return nil, fmt.Errorf("no instances found for service: %s", name)
	}

	// فیلتر نمونه‌های healthy
	var healthy []ServiceInstance
	for _, inst := range instances {
		if inst.Healthy && time.Since(inst.LastSeen) < 30*time.Second {
			healthy = append(healthy, inst)
		}
	}

	if len(healthy) == 0 {
		return nil, fmt.Errorf("no healthy instances for service: %s", name)
	}

	return healthy, nil
}

// GetRandomInstance دریافت یک نمونه تصادفی (برای load balancing)
func (r *ServiceRegistry) GetRandomInstance(name string) (*ServiceInstance, error) {
	instances, err := r.Discover(name)
	if err != nil {
		return nil, err
	}
	idx := rand.Intn(len(instances))
	return &instances[idx], nil
}

// HealthCheck بررسی سلامت سرویس‌ها (heartbeat)
func (r *ServiceRegistry) HealthCheck(name, id string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	instances := r.services[name]
	for i, inst := range instances {
		if inst.ID == id {
			instances[i].LastSeen = time.Now()
			instances[i].Healthy = true
			break
		}
	}
}

// StartHeartbeat شروع heartbeat برای سرویس
func (r *ServiceRegistry) StartHeartbeat(name, id string, interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			r.HealthCheck(name, id)
		}
	}()
}

// ============================================================================
// بخش 4: CQRS (Command Query Responsibility Segregation)
// ============================================================================

// Command و Query

// Command دستور (تغییر وضعیت)
type Command interface {
	GetType() string
}

// Query پرس و جو (خواندن وضعیت)
type Query interface {
	GetType() string
}

// Command Bus
type CommandBus struct {
	handlers map[string]func(Command) error
	mu       sync.RWMutex
}

func NewCommandBus() *CommandBus {
	return &CommandBus{
		handlers: make(map[string]func(Command) error),
	}
}

func (b *CommandBus) Register(commandType string, handler func(Command) error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[commandType] = handler
}

func (b *CommandBus) Dispatch(cmd Command) error {
	b.mu.RLock()
	handler, exists := b.handlers[cmd.GetType()]
	b.mu.RUnlock()

	if !exists {
		return fmt.Errorf("no handler registered for command: %s", cmd.GetType())
	}

	return handler(cmd)
}

// Query Bus
type QueryBus struct {
	handlers map[string]func(Query) (interface{}, error)
	mu       sync.RWMutex
}

func NewQueryBus() *QueryBus {
	return &QueryBus{
		handlers: make(map[string]func(Query) (interface{}, error)),
	}
}

func (b *QueryBus) Register(queryType string, handler func(Query) (interface{}, error)) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[queryType] = handler
}

func (b *QueryBus) Dispatch(query Query) (interface{}, error) {
	b.mu.RLock()
	handler, exists := b.handlers[query.GetType()]
	b.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("no handler registered for query: %s", query.GetType())
	}

	return handler(query)
}

// Example Commands and Queries
type CreateOrderCommand struct {
	OrderID   string
	UserID    string
	ProductID string
	Quantity  int
	Amount    float64
}

func (c CreateOrderCommand) GetType() string { return "CreateOrder" }

type GetOrderQuery struct {
	OrderID string
}

func (g GetOrderQuery) GetType() string { return "GetOrder" }

// CQRS Handlers
type OrderCommandHandler struct {
	eventStore EventStore
}

func (h *OrderCommandHandler) HandleCreateOrder(cmd Command) error {
	order := cmd.(CreateOrderCommand)

	// اعتبارسنجی
	if order.Quantity <= 0 {
		return fmt.Errorf("invalid quantity")
	}

	// ذخیره در write database
	// (در مثال واقعی در دیتابیس write ذخیره می‌شود)
	log.Printf("CQRS Command: Creating order %s for user %s", order.OrderID, order.UserID)

	// انتشار رویداد برای به‌روزرسانی read model
	// eventBus.Publish(OrderCreatedEvent{...})

	return nil
}

type OrderQueryHandler struct {
	readModel ReadModel
}

func (h *OrderQueryHandler) HandleGetOrder(query Query) (interface{}, error) {
	q := query.(GetOrderQuery)

	// خواندن از read database
	// (در مثال واقعی از دیتابیس read می‌خواند)
	log.Printf("CQRS Query: Getting order %s", q.OrderID)

	return map[string]interface{}{
		"order_id": q.OrderID,
		"status":   "pending",
	}, nil
}

// ReadModel model خواندن (دنرمالایز شده)
type ReadModel struct {
	mu      sync.RWMutex
	orders  map[string]OrderReadModel
}

type OrderReadModel struct {
	OrderID   string  `json:"order_id"`
	UserID    string  `json:"user_id"`
	Status    string  `json:"status"`
	Total     float64 `json:"total"`
	CreatedAt string  `json:"created_at"`
}

func NewReadModel() *ReadModel {
	return &ReadModel{
		orders: make(map[string]OrderReadModel),
	}
}

// UpdateFromEvent به‌روزرسانی read model از رویدادها
func (rm *ReadModel) UpdateFromEvent(event interface{}) {
	// اعمال رویداد به read model
	// (در کانال رویدادها مصرف می‌شود)
}

// ============================================================================
// بخش 5: Distributed Tracing با OpenTelemetry + Jaeger
// ============================================================================

// initTracer راه‌اندازی tracer برای distributed tracing
func initTracer(serviceName string) (*sdktrace.TracerProvider, error) {
	// ایجاد Jaeger exporter
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://localhost:14268/api/traces")))
	if err != nil {
		return nil, err
	}

	// ایجاد TracerProvider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		)),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return tp, nil
}

// DistributedTracer ردیاب توزیع شده
type DistributedTracer struct {
	tracer trace.Tracer
}

func NewDistributedTracer(serviceName string) *DistributedTracer {
	return &DistributedTracer{
		tracer: otel.Tracer(serviceName),
	}
}

// StartSpan شروع یک span جدید
func (dt *DistributedTracer) StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return dt.tracer.Start(ctx, name, opts...)
}

// TraceHTTPRequest ردیابی یک درخواست HTTP
func (dt *DistributedTracer) TraceHTTPRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// استخراج context از هدرها
		ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))

		// ایجاد span
		ctx, span := dt.tracer.Start(ctx, fmt.Sprintf("%s %s", r.Method, r.URL.Path),
			trace.WithAttributes(
				attribute.String("http.method", r.Method),
				attribute.String("http.url", r.URL.Path),
				attribute.String("http.user_agent", r.UserAgent()),
			))
		defer span.End()

		// ذخیره context در request
		r = r.WithContext(ctx)

		// فراخوانی handler بعدی
		next.ServeHTTP(w, r)

		// ثبت نتیجه
		// (در wrapper واقعی status code را می‌گیریم)
	})
}

// TraceAPICall ردیابی یک فراخوانی API بین سرویس‌ها
func (dt *DistributedTracer) TraceAPICall(ctx context.Context, service, method, url string) func() {
	ctx, span := dt.tracer.Start(ctx, fmt.Sprintf("Call %s.%s", service, method),
		trace.WithAttributes(
			attribute.String("service", service),
			attribute.String("method", method),
			attribute.String("url", url),
		))

	return func() {
		span.End()
	}
}

// RecordError ثبت خطا در span
func (dt *DistributedTracer) RecordError(span trace.Span, err error) {
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
}

// ============================================================================
// بخش 6: Example - Microservices با Tracing
// ============================================================================

// UserService سرویس کاربران
type UserService struct {
	tracer *DistributedTracer
	registry *ServiceRegistry
	instanceID string
}

func NewUserService(tracer *DistributedTracer, registry *ServiceRegistry) *UserService {
	id := uuid.New().String()
	return &UserService{
		tracer:     tracer,
		registry:   registry,
		instanceID: id,
	}
}

func (s *UserService) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	_, span := s.tracer.tracer.Start(ctx, "UserService.GetUser")
	defer span.End()

	userID := chi.URLParam(r, "id")
	span.SetAttributes(attribute.String("user_id", userID))

	// شبیه‌سازی کار
	time.Sleep(50 * time.Millisecond)

	response := map[string]interface{}{
		"id":    userID,
		"name":  "Ali Rezaei",
		"email": "ali@example.com",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// OrderService سرویس سفارشات (که user service را صدا می‌زند)
type OrderService struct {
	tracer    *DistributedTracer
	registry  *ServiceRegistry
	httpClient *http.Client
}

func NewOrderService(tracer *DistributedTracer, registry *ServiceRegistry) *OrderService {
	return &OrderService{
		tracer:     tracer,
		registry:   registry,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *OrderService) GetOrderHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := s.tracer.tracer.Start(ctx, "OrderService.GetOrder")
	defer span.End()

	orderID := chi.URLParam(r, "id")
	span.SetAttributes(attribute.String("order_id", orderID))

	// شبیه‌سازی دریافت اطلاعات کاربر از user service
	userID := "user-123"

	// دریافت آدرس user service از service discovery
	instances, err := s.registry.Discover("user-service")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// فراخوانی user service (با propagation context)
	url := fmt.Sprintf("http://%s:%d/users/%s", instances[0].Address, instances[0].Port, userID)

	// ایجاد درخواست با context (برای propagation)
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)

	// Inject context به هدرها
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	resp, err := s.httpClient.Do(req)
	if err != nil {
		s.tracer.RecordError(span, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// شبیه‌سازی ایجاد سفارش
	time.Sleep(100 * time.Millisecond)

	response := map[string]interface{}{
		"order_id": orderID,
		"status":   "created",
		"total":    99.99,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ============================================================================
// بخش 7: Main - راه‌اندازی سیستم کامل
// ============================================================================

func setupCompleteSystem() {
	// 1. راه‌اندازی tracer
	tp, err := initTracer("api-gateway")
	if err != nil {
		log.Fatal(err)
	}
	defer tp.Shutdown(context.Background())

	// 2. ایجادサービス registry
	registry := NewServiceRegistry()

	// 3. ایجاد tracer
	tracer := NewDistributedTracer("api-gateway")

	// 4. راه‌اندازی microservices در گوروتین‌های جدا
	// User Service
	userService := NewUserService(NewDistributedTracer("user-service"), registry)
	userRouter := chi.NewRouter()
	userRouter.Get("/users/{id}", userService.GetUserHandler)
	go func() {
		// ثبت در service registry
		registry.Register(ServiceInstance{
			ID:      uuid.New().String(),
			Name:    "user-service",
			Address: "localhost",
			Port:    8081,
			Healthy: true,
		})
		log.Println("User service running on :8081")
		http.ListenAndServe(":8081", userRouter)
	}()

	// Order Service
	orderService := NewOrderService(NewDistributedTracer("order-service"), registry)
	orderRouter := chi.NewRouter()
	orderRouter.Get("/orders/{id}", orderService.GetOrderHandler)
	go func() {
		registry.Register(ServiceInstance{
			ID:      uuid.New().String(),
			Name:    "order-service",
			Address: "localhost",
			Port:    8082,
			Healthy: true,
		})
		log.Println("Order service running on :8082")
		http.ListenAndServe(":8082", orderRouter)
	}()

	// 5. API Gateway
	gateway := NewAPIGateway()
	gateway.AddRoute("/users", "localhost:8081", []string{"GET"}, true, 100)
	gateway.AddRoute("/orders", "localhost:8082", []string{"GET"}, true, 100)

	// 6. راه‌اندازی Gateway با tracing middleware
	handler := tracer.TraceHTTPRequest(http.HandlerFunc(gateway.HandleRequest))

	log.Println("API Gateway running on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

// ============================================================================
// بخش 8: Best Practices
// ============================================================================

func bestPractices() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("💡 ARCHITECTURE BEST PRACTICES")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ 1. MONOLITH to MICROSERVICES                                              │
├─────────────────────────────────────────────────────────────────────────────┤
│    • Start with modular monolith (clear boundaries)                       │
│    • Extract one service at a time                                        │
│    • Use API Gateway for routing                                           │
│    • Share nothing (database per service)                                 │
│    • Automate everything (CI/CD, deployment)                              │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 2. API GATEWAY                                                           │
├─────────────────────────────────────────────────────────────────────────────┤
│    • Single entry point for all clients                                   │
│    • Handle cross-cutting concerns (auth, logging, rate limiting)         │
│    • Use reverse proxy (Traefik, Nginx, Envoy)                            │
│    • Implement circuit breakers                                           │
│    • Cache responses when possible                                        │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 3. SERVICE DISCOVERY                                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│    • Use Consul, etcd, or Kubernetes for production                       │
│    • Implement health checks                                              │
│    • Use client-side or server-side load balancing                        │
│    • Cache discovery results                                              │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 4. CQRS                                                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│    • Use only when needed (complex read models)                          │
│    • Separate databases for read/write                                   │
│    • Use event sourcing for consistency                                  │
│    • Handle eventual consistency                                         │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 5. DISTRIBUTED TRACING                                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│    • Sample only a percentage of requests                                │
│    • Add business-relevant attributes                                    │
│    • Use correlation IDs                                                 │
│    • Monitor trace data                                                  │
│    • Set up alerts on error patterns                                     │
└─────────────────────────────────────────────────────────────────────────────┘
`)
}

// ============================================================================
// بخش 9: Main
// ============================================================================

func main() {
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("🎯 ADVANCED ARCHITECTURE PATTERNS")
	fmt.Println("Monolith vs Microservices | API Gateway | Service Discovery | CQRS | Distributed Tracing")
	fmt.Println(strings.Repeat("=", 80))

	comparisonMonolithMicroservices()
	traefikConfigExample()
	bestPractices()

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🚀 Starting Demo (requires Jaeger and services)")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
📋 Prerequisites:

   # Start Jaeger for tracing
   $ docker run -d --name jaeger \
     -e COLLECTOR_ZIPKIN_HOST_PORT=:9411 \
     -p 5775:5775/udp \
     -p 6831:6831/udp \
     -p 6832:6832/udp \
     -p 5778:5778 \
     -p 16686:16686 \
     -p 14268:14268 \
     -p 14250:14250 \
     -p 9411:9411 \
     jaegertracing/all-in-one:latest

   # Then run the application
   $ go run architecture_patterns_guide.go

   # Test endpoints
   $ curl http://localhost:8080/users/123
   $ curl http://localhost:8080/orders/456

   # View traces
   $ open http://localhost:16686
`)

	// در صورت نیاز به اجرای واقعی، این خط را uncomment کنید
	// setupCompleteSystem()
}

// تابع کمکی برای تکرار رشته
func stringsRepeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}

// برای کامپایل
var _ = io.Copy
type ioWriter struct{}
var _ = io.Writer