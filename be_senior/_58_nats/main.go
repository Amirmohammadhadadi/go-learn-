// ============================================================================
// FILE: nats_guide.go
// TITLE: راهنمای کامل NATS در Go - سیستم پیام‌رسانی سبک و پرسرعت
// HOW TO RUN: go run nats_guide.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - NATS چیست؟
// ============================================================================
//
// NATS (Natural language-based messaging system) یک سیستم پیام‌رسانی open-source،
// سبک و پرسرعت است که توسط Cloud Native Computing Foundation (CNCF) نگهداری می‌شود.
//
// ویژگی‌های کلیدی NATS:
// 1. سرعت بسیار بالا (قابل مقایسه با TCP)
// 2. سبک بودن (باینری حدود 20MB)
// 3. سادگی (پروتکل متنی ساده)
// 4. At-most-once و At-least-once (با JetStream)
// 5.پشتیبانی از Request-Reply, Pub/Sub, Queue Groups
//
// اجزای NATS:
// - Client: برنامهای که به NATS متصل می‌شود
// - Server: سرور NATS (یک یا چند node)
// - Subject: موضوع برای routing پیام‌ها (مثل اینترنتی topic)
// - Publisher: sender (انتشار دهنده پیام)
// - Subscriber: receiver (دریافت کننده پیام)
// - Queue Group: گروهی از subscriberها که پیام‌ها را به صورت round-robin تقسیم می‌کنند
//
// مقایسه با سایر message brokers:
// ┌─────────────────────────────────────────────────────────────────────────────┐
// │ ویژگی              │ NATS         │ RabbitMQ    │ Kafka                    │
// ├────────────────────┼──────────────┼─────────────┼──────────────────────────┤
// │ Latency            │ microsecond  │ millisecond │ millisecond              │
// │ Throughput         │ بسیار بالا   │ متوسط      │ بسیار بالا               │
// │ Persistence        │ با JetStream │ ✅          │ ✅                       │
// │ Message Ordering   │ ✅           │ ✅          │ ✅ (در partition)        │
// │ Protocol           │ TCP          │ AMQP        │ Kafka Protocol           │
// │ Use Case           │ Microservices│ RPC/Tasks   │ Event Streaming          │
// │ Complexity         │ پایین        │ متوسط      │ بالا                     │
// └─────────────────────────────────────────────────────────────────────────────┘
//
// قانون طلایی:
// "برای ارتباطات low-latency و high-throughput بین میکروسرویس‌ها از NATS استفاده کن.
//  برای persistence از JetStream استفاده کن.
//  از Queue Groups برای load balancing استفاده کن."
// ============================================================================

package __nats

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	_ "github.com/nats-io/nats.go"
)

// ============================================================================
// بخش 1: مدل‌های داده
// ============================================================================

// OrderEvent رویداد سفارش
type OrderEvent struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	ProductID string    `json:"product_id"`
	Quantity  int       `json:"quantity"`
	Amount    float64   `json:"amount"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

// PaymentEvent رویداد پرداخت
type PaymentEvent struct {
	ID      string  `json:"id"`
	OrderID string  `json:"order_id"`
	Amount  float64 `json:"amount"`
	Method  string  `json:"method"`
	Status  string  `json:"status"`
}

// NotificationRequest درخواست نوتیفیکیشن
type NotificationRequest struct {
	UserID  string `json:"user_id"`
	Type    string `json:"type"`
	Title   string `json:"title"`
	Message string `json:"message"`
}

// Response پاسخ کلی
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ============================================================================
// بخش 2: اتصال و تنظیمات پایه NATS
// ============================================================================

// NATSConfig تنظیمات NATS
type NATSConfig struct {
	URL           string
	Username      string
	Password      string
	Token         string
	MaxReconnects int
	ReconnectWait time.Duration
	Timeout       time.Duration
}

// DefaultNATSConfig تنظیمات پیش‌فرض
func DefaultNATSConfig() *NATSConfig {
	return &NATSConfig{
		URL:           nats.DefaultURL,
		MaxReconnects: -1, // نامحدود
		ReconnectWait: 2 * time.Second,
		Timeout:       10 * time.Second,
	}
}

// NATSClient کلاینت NATS
type NATSClient struct {
	conn   *nats.Conn
	js     nats.JetStreamContext
	config *NATSConfig
	mu     sync.RWMutex
	subs   []*nats.Subscription
}

// NewNATSClient ایجاد کلاینت NATS جدید
func NewNATSClient(config *NATSConfig) (*NATSClient, error) {
	opts := []nats.Option{
		nats.MaxReconnects(config.MaxReconnects),
		nats.ReconnectWait(config.ReconnectWait),
		nats.Timeout(config.Timeout),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			log.Printf("NATS disconnected: %v", err)
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			log.Printf("NATS reconnected to %s", nc.ConnectedUrl())
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			log.Printf("NATS connection closed")
		}),
	}

	// احراز هویت
	if config.Username != "" && config.Password != "" {
		opts = append(opts, nats.UserInfo(config.Username, config.Password))
	}
	if config.Token != "" {
		opts = append(opts, nats.Token(config.Token))
	}

	// اتصال
	conn, err := nats.Connect(config.URL, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	// JetStream (برای persistence)
	js, err := conn.JetStream()
	if err != nil {
		log.Printf("Warning: JetStream not available: %v", err)
	}

	return &NATSClient{
		conn:   conn,
		js:     js,
		config: config,
		subs:   make([]*nats.Subscription, 0),
	}, nil
}

// Close بستن اتصال
func (c *NATSClient) Close() {
	for _, sub := range c.subs {
		sub.Unsubscribe()
	}
	c.conn.Close()
}

// IsConnected بررسی اتصال
func (c *NATSClient) IsConnected() bool {
	return c.conn != nil && c.conn.IsConnected()
}

// ============================================================================
// بخش 3: Publish-Subscribe (Pub/Sub) با NATS
// ============================================================================

// Publisher منتشر کننده پیام
type Publisher struct {
	client *NATSClient
}

// NewPublisher ایجاد publisher جدید
func NewPublisher(client *NATSClient) *Publisher {
	return &Publisher{client: client}
}

// Publish انتشار پیام
func (p *Publisher) Publish(subject string, data interface{}) error {
	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}
	return p.client.conn.Publish(subject, body)
}

// PublishAsync انتشار پیام به صورت غیرهمزمان
func (p *Publisher) PublishAsync(subject string, data interface{}) error {
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return p.client.conn.PublishAsync(subject, body)
}

// PublishWithReply انتشار پیام و دریافت پاسخ
func (p *Publisher) PublishWithReply(subject string, data interface{}, timeout time.Duration) (*nats.Msg, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return p.client.conn.Request(subject, body, timeout)
}

// Subscriber مشترک پیام
type Subscriber struct {
	client *NATSClient
}

// NewSubscriber ایجاد subscriber جدید
func NewSubscriber(client *NATSClient) *Subscriber {
	return &Subscriber{client: client}
}

// Subscribe اشتراک در subject
func (s *Subscriber) Subscribe(subject string, handler func(*nats.Msg)) (*nats.Subscription, error) {
	sub, err := s.client.conn.Subscribe(subject, func(msg *nats.Msg) {
		handler(msg)
	})
	if err != nil {
		return nil, err
	}
	s.client.subs = append(s.client.subs, sub)
	return sub, nil
}

// SubscribeWithQueue اشتراک با queue group (load balancing)
func (s *Subscriber) SubscribeWithQueue(subject, queue string, handler func(*nats.Msg)) (*nats.Subscription, error) {
	sub, err := s.client.conn.QueueSubscribe(subject, queue, func(msg *nats.Msg) {
		handler(msg)
	})
	if err != nil {
		return nil, err
	}
	s.client.subs = append(s.client.subs, sub)
	return sub, nil
}

// ============================================================================
// بخش 4: Request-Reply (RPC) با NATS
// ============================================================================

// RPCServer سرویس RPC
type RPCServer struct {
	client   *NATSClient
	handlers map[string]func(*nats.Msg) interface{}
}

// NewRPCServer ایجاد RPC server جدید
func NewRPCServer(client *NATSClient) *RPCServer {
	return &RPCServer{
		client:   client,
		handlers: make(map[string]func(*nats.Msg) interface{}),
	}
}

// Register ثبت handler
func (r *RPCServer) Register(method string, handler func(*nats.Msg) interface{}) {
	r.handlers[method] = handler

	subject := fmt.Sprintf("rpc.%s", method)
	r.client.conn.Subscribe(subject, func(msg *nats.Msg) {
		result := handler(msg)
		response, _ := json.Marshal(result)
		msg.Respond(response)
	})
}

// Start شروع RPC server
func (r *RPCServer) Start() error {
	log.Println("RPC server started")
	return nil
}

// RPCClient کلاینت RPC
type RPCClient struct {
	client *NATSClient
}

// NewRPCClient ایجاد RPC client جدید
func NewRPCClient(client *NATSClient) *RPCClient {
	return &RPCClient{client: client}
}

// Call فراخوانی متد RPC
func (c *RPCClient) Call(method string, request interface{}, timeout time.Duration) ([]byte, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	subject := fmt.Sprintf("rpc.%s", method)
	msg, err := c.client.conn.Request(subject, body, timeout)
	if err != nil {
		return nil, err
	}

	return msg.Data, nil
}

// ============================================================================
// بخش 5: JetStream (Persistent Messaging)
// ============================================================================

// JetStreamManager مدیریت JetStream
type JetStreamManager struct {
	client *NATSClient
}

// NewJetStreamManager ایجاد JetStream manager
func NewJetStreamManager(client *NATSClient) *JetStreamManager {
	return &JetStreamManager{client: client}
}

// CreateStream ایجاد stream
func (jsm *JetStreamManager) CreateStream(streamName, subject string, replicas int) error {
	if jsm.client.js == nil {
		return fmt.Errorf("JetStream not available")
	}

	cfg := &nats.StreamConfig{
		Name:      streamName,
		Subjects:  []string{subject},
		Retention: nats.LimitsPolicy,
		MaxAge:    7 * 24 * time.Hour,
		MaxBytes:  1024 * 1024 * 1024, // 1GB
		Storage:   nats.FileStorage,
		Replicas:  replicas,
	}

	_, err := jsm.client.js.AddStream(cfg)
	if err != nil {
		return fmt.Errorf("failed to create stream: %w", err)
	}

	log.Printf("Stream %s created successfully", streamName)
	return nil
}

// PublishPersistent انتشار پیام پایدار
func (jsm *JetStreamManager) PublishPersistent(streamName, subject string, data interface{}) error {
	if jsm.client.js == nil {
		return fmt.Errorf("JetStream not available")
	}

	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = jsm.client.js.Publish(subject, body)
	if err != nil {
		return fmt.Errorf("failed to publish persistent message: %w", err)
	}

	return nil
}

// SubscribePersistent اشتراک پایدار
func (jsm *JetStreamManager) SubscribePersistent(streamName, subject, durableName string, handler func(*nats.Msg)) error {
	if jsm.client.js == nil {
		return fmt.Errorf("JetStream not available")
	}

	sub, err := jsm.client.js.Subscribe(subject, handler, nats.Durable(durableName), nats.ManualAck())
	if err != nil {
		return fmt.Errorf("failed to subscribe: %w", err)
	}

	go func() {
		defer sub.Unsubscribe()
		<-make(chan struct{})
	}()

	return nil
}

// ============================================================================
// بخش 6: NATS with Context (برای timeout و cancellation)
// ============================================================================

// RequestWithContext ارسال درخواست با context
func (c *NATSClient) RequestWithContext(ctx context.Context, subject string, data []byte) (*nats.Msg, error) {
	msg := nats.NewMsg(subject)
	msg.Data = data

	// ایجاد کانال برای دریافت پاسخ
	respChan := make(chan *nats.Msg, 1)
	var err error

	sub, err := c.conn.Subscribe(msg.Reply, func(m *nats.Msg) {
		respChan <- m
	})
	if err != nil {
		return nil, err
	}
	defer sub.Unsubscribe()

	// ارسال درخواست
	if err := c.conn.PublishMsg(msg); err != nil {
		return nil, err
	}

	// منتظر پاسخ یا timeout
	select {
	case resp := <-respChan:
		return resp, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// ============================================================================
// بخش 7: مثال کاربردی: سیستم سفارشات
// ============================================================================

// OrderService سرویس سفارشات با NATS
type OrderService struct {
	publisher  *Publisher
	subscriber *Subscriber
	rpcServer  *RPCServer
	rpcClient  *RPCClient
	jetStream  *JetStreamManager
}

// NewOrderService ایجاد سرویس سفارشات
func NewOrderService(client *NATSClient) *OrderService {
	return &OrderService{
		publisher:  NewPublisher(client),
		subscriber: NewSubscriber(client),
		rpcServer:  NewRPCServer(client),
		rpcClient:  NewRPCClient(client),
		jetStream:  NewJetStreamManager(client),
	}
}

// StartOrderService شروع سرویس سفارشات
func (s *OrderService) StartOrderService() error {
	// 1. ایجاد stream برای JetStream
	if err := s.jetStream.CreateStream("ORDERS", "order.*", 1); err != nil {
		log.Printf("Warning: JetStream stream creation failed: %v", err)
	}

	// 2. RPC handler برای ایجاد سفارش
	s.rpcServer.Register("createOrder", func(msg *nats.Msg) interface{} {
		var order OrderEvent
		if err := json.Unmarshal(msg.Data, &order); err != nil {
			return Response{Success: false, Message: err.Error()}
		}

		order.ID = generateID()
		order.Timestamp = time.Now()
		order.Status = "created"

		// انتشار رویداد سفارش جدید
		s.publisher.Publish("order.created", order)

		return Response{Success: true, Data: order}
	})

	// 3. Subscriber برای پردازش سفارش‌ها (با queue group برای load balancing)
	_, err := s.subscriber.SubscribeWithQueue("order.created", "order-processors", func(msg *nats.Msg) {
		var order OrderEvent
		if err := json.Unmarshal(msg.Data, &order); err != nil {
			log.Printf("Error unmarshaling order: %v", err)
			return
		}

		log.Printf("Processing order: %s for user %s", order.ID, order.UserID)

		// شبیه‌سازی پردازش
		time.Sleep(100 * time.Millisecond)

		// انتشار رویداد پردازش شده
		order.Status = "processed"
		s.publisher.Publish("order.processed", order)
	})

	if err != nil {
		return err
	}

	// 4. Subscriber برای سفارشات پردازش شده (برای metrics)
	_, err = s.subscriber.Subscribe("order.processed", func(msg *nats.Msg) {
		var order OrderEvent
		json.Unmarshal(msg.Data, &order)
		log.Printf("Order %s processed successfully", order.ID)
	})

	return err
}

// CreateOrderViaRPC ایجاد سفارش از طریق RPC
func (s *OrderService) CreateOrderViaRPC(order OrderEvent) (*Response, error) {
	resp, err := s.rpcClient.Call("createOrder", order, 5*time.Second)
	if err != nil {
		return nil, err
	}

	var response Response
	if err := json.Unmarshal(resp, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// ============================================================================
// بخش 8: Queue Groups (فرآیندهای موازی)
// ============================================================================

// QueueWorker worker در queue group
type QueueWorker struct {
	ID      int
	Client  *NATSClient
	Subject string
	Queue   string
	Handler func(*nats.Msg)
}

// StartWorker شروع worker
func (w *QueueWorker) StartWorker(wg *sync.WaitGroup) {
	defer wg.Done()

	_, err := w.Client.conn.QueueSubscribe(w.Subject, w.Queue, func(msg *nats.Msg) {
		w.Handler(msg)
	})
	if err != nil {
		log.Printf("Worker %d error: %v", w.ID, err)
	}

	select {} // منتظر ماندن
}

// StartQueueGroup شروع گروه workers
func StartQueueGroup(client *NATSClient, subject, queue string, handler func(*nats.Msg), workerCount int) {
	var wg sync.WaitGroup

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		worker := &QueueWorker{
			ID:      i + 1,
			Client:  client,
			Subject: subject,
			Queue:   queue,
			Handler: handler,
		}
		go worker.StartWorker(&wg)
	}

	log.Printf("Started %d workers for queue '%s' on subject '%s'", workerCount, queue, subject)
	wg.Wait()
}

// ============================================================================
// بخش 9: Heartbeat و Health Checks
// ============================================================================

// HeartbeatSender ارسال کننده heartbeat
type HeartbeatSender struct {
	client    *NATSClient
	serviceID string
	interval  time.Duration
	stopCh    chan struct{}
}

// NewHeartbeatSender ایجاد heartbeat sender جدید
func NewHeartbeatSender(client *NATSClient, serviceID string, interval time.Duration) *HeartbeatSender {
	return &HeartbeatSender{
		client:    client,
		serviceID: serviceID,
		interval:  interval,
		stopCh:    make(chan struct{}),
	}
}

// Start شروع ارسال heartbeat
func (h *HeartbeatSender) Start() {
	ticker := time.NewTicker(h.interval)
	go func() {
		for {
			select {
			case <-h.stopCh:
				ticker.Stop()
				return
			case <-ticker.C:
				heartbeat := map[string]interface{}{
					"service_id": h.serviceID,
					"timestamp":  time.Now(),
					"status":     "alive",
				}
				h.client.conn.Publish("health.heartbeat", mustMarshal(heartbeat))
			}
		}
	}()
}

// Stop توقف heartbeat
func (h *HeartbeatSender) Stop() {
	close(h.stopCh)
}

// HeartbeatReceiver دریافت کننده heartbeat
type HeartbeatReceiver struct {
	client   *NATSClient
	services map[string]time.Time
	mu       sync.RWMutex
}

// NewHeartbeatReceiver ایجاد heartbeat receiver جدید
func NewHeartbeatReceiver(client *NATSClient) *HeartbeatReceiver {
	receiver := &HeartbeatReceiver{
		client:   client,
		services: make(map[string]time.Time),
	}

	client.conn.Subscribe("health.heartbeat", func(msg *nats.Msg) {
		var heartbeat map[string]interface{}
		if err := json.Unmarshal(msg.Data, &heartbeat); err != nil {
			return
		}

		serviceID, ok := heartbeat["service_id"].(string)
		if !ok {
			return
		}

		receiver.mu.Lock()
		receiver.services[serviceID] = time.Now()
		receiver.mu.Unlock()
	})

	return receiver
}

// GetActiveServices دریافت سرویس‌های فعال
func (r *HeartbeatReceiver) GetActiveServices(timeout time.Duration) []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var active []string
	now := time.Now()
	for id, lastSeen := range r.services {
		if now.Sub(lastSeen) <= timeout {
			active = append(active, id)
		}
	}
	return active
}

// ============================================================================
// بخش 10: Monitoring and Metrics
// ============================================================================

// MetricsCollector جمع‌آوری کننده متریک‌ها
type MetricsCollector struct {
	client         *NATSClient
	publishCount   int64
	subscribeCount int64
	errorCount     int64
	mu             sync.RWMutex
}

// NewMetricsCollector ایجاد metrics collector
func NewMetricsCollector(client *NATSClient) *MetricsCollector {
	return &MetricsCollector{
		client: client,
	}
}

// RecordPublish ثبت publish
func (m *MetricsCollector) RecordPublish() {
	atomic.AddInt64(&m.publishCount, 1)
}

// RecordSubscribe ثبت subscribe
func (m *MetricsCollector) RecordSubscribe() {
	atomic.AddInt64(&m.subscribeCount, 1)
}

// RecordError ثبت خطا
func (m *MetricsCollector) RecordError() {
	atomic.AddInt64(&m.errorCount, 1)
}

// GetStats دریافت آمار
func (m *MetricsCollector) GetStats() map[string]int64 {
	return map[string]int64{
		"publish_count":   atomic.LoadInt64(&m.publishCount),
		"subscribe_count": atomic.LoadInt64(&m.subscribeCount),
		"error_count":     atomic.LoadInt64(&m.errorCount),
	}
}

// ============================================================================
// بخش 11: Full Example Application
// ============================================================================

func runNATSExample() {
	log.Println("=== NATS Messaging System Example ===\n")

	// 1. اتصال به NATS
	config := DefaultNATSConfig()
	client, err := NewNATSClient(config)
	if err != nil {
		log.Printf("Failed to connect to NATS: %v (make sure NATS server is running)", err)
		log.Println("To start NATS server: docker run -d --name nats -p 4222:4222 nats:latest")
		return
	}
	defer client.Close()

	log.Println("Connected to NATS server")

	// 2. Setup order service
	orderService := NewOrderService(client)
	if err := orderService.StartOrderService(); err != nil {
		log.Printf("Error starting order service: %v", err)
	}

	// 3. Create order via RPC
	order := OrderEvent{
		UserID:    "user-123",
		ProductID: "prod-001",
		Quantity:  2,
		Amount:    99.98,
	}

	response, err := orderService.CreateOrderViaRPC(order)
	if err != nil {
		log.Printf("RPC error: %v", err)
	} else {
		log.Printf("RPC Response: %+v", response)
	}

	// 4. Direct publish
	publisher := NewPublisher(client)
	notification := NotificationRequest{
		UserID:  "user-123",
		Type:    "order_confirmation",
		Title:   "Order Confirmed",
		Message: "Your order has been confirmed",
	}

	if err := publisher.Publish("notifications.email", notification); err != nil {
		log.Printf("Publish error: %v", err)
	}

	// 5. Subscribe to various subjects
	subscriber := NewSubscriber(client)

	// گوش دادن به نوتیفیکیشن‌ها
	subscriber.Subscribe("notifications.*", func(msg *nats.Msg) {
		log.Printf("Received notification on subject %s: %s", msg.Subject, string(msg.Data))
	})

	// گوش دادن به رویدادهای سفارش
	subscriber.Subscribe("order.*", func(msg *nats.Msg) {
		var order OrderEvent
		if err := json.Unmarshal(msg.Data, &order); err == nil {
			log.Printf("Order event: %s - %s", msg.Subject, order.Status)
		}
	})

	// 6. Queue group example
	log.Println("\nStarting queue group workers...")
	StartQueueGroup(client, "tasks.process", "task-workers", func(msg *nats.Msg) {
		log.Printf("Processing task: %s", string(msg.Data))
		time.Sleep(100 * time.Millisecond)
	}, 3)

	// 7. Publish tasks to queue
	for i := 0; i < 10; i++ {
		task := fmt.Sprintf("Task %d", i+1)
		publisher.Publish("tasks.process", task)
	}

	// 8. Keep running for a while to see messages
	log.Println("\nListening for messages for 5 seconds...")
	time.Sleep(5 * time.Second)
}

// ============================================================================
// بخش 12: Best Practices
// ============================================================================

func bestPractices() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("💡 NATS BEST PRACTICES")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ 1. CONNECTION MANAGEMENT                                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│    • Reconnect handlers for handling disconnections                       │
│    • Configure reasonable timeouts                                        │
│    • Use connection pooling for high traffic                              │
│    • Monitor connection status                                            │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 2. SUBJECT NAMING                                                         │
├─────────────────────────────────────────────────────────────────────────────┤
│    • Use dot-separated hierarchies (e.g., "orders.created")               │
│    • Use wildcards: "*" for single token, ">" for multiple                │
│    • Keep subjects short and meaningful                                   │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 3. MESSAGE DESIGN                                                         │
├─────────────────────────────────────────────────────────────────────────────┤
│    • Keep messages small (preferably <1MB)                                │
│    • Use schemas for message evolution                                    │
│    • Include message ID and timestamp                                     │
│    • Use request-reply for RPC patterns                                   │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 4. QUEUE GROUPS                                                           │
├─────────────────────────────────────────────────────────────────────────────┤
│    • Use queue groups for load balancing                                  │
│    • Ensure idempotent processing                                         │
│    • Monitor queue depth and worker health                                │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 5. JETSTREAM                                                              │
├─────────────────────────────────────────────────────────────────────────────┤
│    • Use for persistent storage and replay                                │
│    • Set appropriate retention policies                                   │
│    • Monitor stream depth                                                 │
│    • Configure replicas for high availability                             │
└─────────────────────────────────────────────────────────────────────────────┘
`)
}

// ============================================================================
// بخش 13: Comparison Table
// ============================================================================

func comparisonTable() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📊 NATS vs RABBITMQ vs KAFKA")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ FEATURE              │ NATS         │ RabbitMQ    │ Kafka                    │
├──────────────────────┼──────────────┼─────────────┼──────────────────────────┤
│ Use Case             │ Microservices│ RPC/Tasks   │ Event Streaming          │
│ Latency              │ μs           │ ms          │ ms                       │
│ Throughput           │ >1M msg/sec  │ ~50K msg/sec│ >1M msg/sec              │
│ Persistence          │ Optional     │ Yes         │ Yes                      │
│ Message Ordering     │ Per subject  │ Per queue   │ Per partition            │
│ Complexity           │ Low          │ Medium      │ High                     │
│ Learning Curve       │ Easy         │ Moderate    │ Difficult                │
│ Language Support     │ Many         │ Many        │ Many                     │
│ Best For             │ Real-time    │ Task queue  │ Analytics, Audit         │
└─────────────────────────────────────────────────────────────────────────────┘
`)
}

// ============================================================================
// بخش 14: Helper Functions
// ============================================================================

func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func mustMarshal(v interface{}) []byte {
	data, _ := json.Marshal(v)
	return data
}

// ============================================================================
// بخش 15: Main
// ============================================================================

func main() {
	fmt.Println(stringsRepeat("=", 80))
	fmt.Println("🎯 NATS MESSAGING SYSTEM IN GO")
	fmt.Println("Lightweight, High-Performance Message Broker")
	fmt.Println(stringsRepeat("=", 80))

	bestPractices()
	comparisonTable()

	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🚀 Running NATS Examples")
	fmt.Println(stringsRepeat("=", 80))
	fmt.Println()
	fmt.Println("⚠️  Make sure NATS server is running:")
	fmt.Println("   docker run -d --name nats -p 4222:4222 nats:latest")
	fmt.Println("   or")
	fmt.Println("   nats-server")
	fmt.Println()

	runNATSExample()

	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🎯 NATS GUIDE - COMPLETE")
	fmt.Println("Build fast, scalable microservices with NATS!")
	fmt.Println(stringsRepeat("=", 80))
}

// تابع کمکی برای atomic operations
var _ = &sync.Mutex{}
var atomicInt64 = int64(0)

func atomicAddInt64(addr *int64, delta int64) int64 {
	// در عمل از atomic.AddInt64 استفاده کنید
	*addr += delta
	return *addr
}

func atomicLoadInt64(addr *int64) int64 {
	return *addr
}

// تابع کمکی برای تکرار رشته
func stringsRepeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}
