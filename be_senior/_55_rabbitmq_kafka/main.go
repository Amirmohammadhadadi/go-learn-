// ============================================================================
// FILE: rabbitmq_kafka_guide.go
// TITLE: راهنمای کامل RabbitMQ و Apache Kafka در Go - Message Queue & Event Streaming
// HOW TO RUN: go run rabbitmq_kafka_guide.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - Message Queue و Event Streaming چیست؟
// ============================================================================
//
// Message Queue و Event Streaming برای ارتباط ناهمگام (asynchronous) بین سرویس‌ها استفاده می‌شوند.
//
// RabbitMQ (Message Queue):
// - مدل: Producer → Queue → Consumer
// - پیام‌ها پس از مصرف از صف حذف می‌شوند
// - مناسب برای: Task distribution, RPC, Decoupling services
// - ویژگی‌ها: Routing, Topics, Dead Letter Queue, Publisher Confirms
//
// Apache Kafka (Event Streaming):
// - مدل: Producer → Topic (Partitioned Log) → Consumer
// - پیام‌ها پس از مصرف باقی می‌مانند (با retention)
// - مناسب برای: Event sourcing, Real-time analytics, Log aggregation
// - ویژگی‌ها: Partitioning, Replication, Exactly-once semantics, Kafka Streams
//
// مقایسه سریع:
// ┌─────────────────────────────────────────────────────────────────────────────┐
// │ ویژگی              │ RabbitMQ                    │ Kafka                    │
// ├────────────────────┼─────────────────────────────┼──────────────────────────┤
// │ مدل                │ Queue                       │ Topic (log)              │
// │ مصرف پیام          │ Destructive (حذف می‌شود)     │ Non-destructive (باقی می‌ماند)│
// │ ترتیب پیام         │ تضمین شده در یک queue        │ تضمین شده در یک partition │
// │ Throughput         │ متوسط (~50k msg/sec)        │ بسیار بالا (~1M msg/sec)  │
// │ Persistence        │ پیام‌ها تا مصرف در دیسک     │ پیام‌ها با retention      │
// │ Use cases          │ Task distribution, RPC      │ Event streaming, Analytics│
// │ Language           │ Erlang                      │ Scala/Java               │
// └─────────────────────────────────────────────────────────────────────────────┘
//
// قانون طلایی:
// "برای task distribution و RPC از RabbitMQ استفاده کن.
//  برای event streaming و real-time analytics از Kafka استفاده کن.
//  همیشه از dead letter queue برای پیام‌های ناموفق استفاده کن."
// ============================================================================

package __rabbitmq_kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	// Importهای RabbitMQ
	amqp "github.com/rabbitmq/amqp091-go"

	// Importهای Kafka
	"github.com/segmentio/kafka-go"
)

// ============================================================================
// بخش 1: Models و Types مشترک
// ============================================================================

// Event رویداد عمومی
type Event struct {
	ID        string      `json:"id"`
	Type      string      `json:"type"`
	Timestamp time.Time   `json:"timestamp"`
	Source    string      `json:"source"`
	Data      interface{} `json:"data"`
}

// OrderEvent رویداد سفارش
type OrderEvent struct {
	OrderID     string    `json:"order_id"`
	UserID      string    `json:"user_id"`
	ProductID   string    `json:"product_id"`
	Quantity    int       `json:"quantity"`
	TotalAmount float64   `json:"total_amount"`
	Status      string    `json:"status"`
	Timestamp   time.Time `json:"timestamp"`
}

// PaymentEvent رویداد پرداخت
type PaymentEvent struct {
	PaymentID     string    `json:"payment_id"`
	OrderID       string    `json:"order_id"`
	Amount        float64   `json:"amount"`
	PaymentMethod string    `json:"payment_method"`
	Status        string    `json:"status"`
	Timestamp     time.Time `json:"timestamp"`
}

// NotificationEvent رویداد نوتیفیکیشن
type NotificationEvent struct {
	UserID    string    `json:"user_id"`
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Email     string    `json:"email,omitempty"`
	Phone     string    `json:"phone,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

// ============================================================================
// بخش 2: RabbitMQ - اتصال و تنظیمات پایه
// ============================================================================

// RabbitMQConfig تنظیمات RabbitMQ
type RabbitMQConfig struct {
	URL        string
	Exchange   string
	Queue      string
	RoutingKey string
}

// DefaultRabbitMQConfig تنظیمات پیش‌فرض
func DefaultRabbitMQConfig() *RabbitMQConfig {
	return &RabbitMQConfig{
		URL:        "amqp://guest:guest@localhost:5672/",
		Exchange:   "events.exchange",
		Queue:      "events.queue",
		RoutingKey: "events.key",
	}
}

// RabbitMQClient کلاینت RabbitMQ
type RabbitMQClient struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	config  *RabbitMQConfig
	mu      sync.RWMutex
}

// NewRabbitMQClient ایجاد کلاینت RabbitMQ جدید
func NewRabbitMQClient(config *RabbitMQConfig) (*RabbitMQClient, error) {
	conn, err := amqp.Dial(config.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	client := &RabbitMQClient{
		conn:    conn,
		channel: channel,
		config:  config,
	}

	// تنظیم exchange
	if err := client.setupExchange(); err != nil {
		return nil, err
	}

	// تنظیم queue
	if err := client.setupQueue(); err != nil {
		return nil, err
	}

	return client, nil
}

// setupExchange تنظیم exchange
func (c *RabbitMQClient) setupExchange() error {
	return c.channel.ExchangeDeclare(
		c.config.Exchange, // name
		"topic",           // type
		true,              // durable
		false,             // auto-deleted
		false,             // internal
		false,             // no-wait
		nil,               // arguments
	)
}

// setupQueue تنظیم queue
func (c *RabbitMQClient) setupQueue() error {
	_, err := c.channel.QueueDeclare(
		c.config.Queue, // name
		true,           // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		return err
	}

	// bind queue به exchange
	return c.channel.QueueBind(
		c.config.Queue,      // queue name
		c.config.RoutingKey, // routing key
		c.config.Exchange,   // exchange
		false,
		nil,
	)
}

// Publish ارسال پیام به RabbitMQ
func (c *RabbitMQClient) Publish(ctx context.Context, routingKey string, message interface{}) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	return c.channel.PublishWithContext(
		ctx,
		c.config.Exchange,
		routingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			Timestamp:    time.Now(),
			DeliveryMode: amqp.Persistent, // ذخیره در دیسک
		},
	)
}

// Consume مصرف پیام از RabbitMQ
func (c *RabbitMQClient) Consume(queue string, handler func([]byte) error) error {
	messages, err := c.channel.Consume(
		queue,
		"",    // consumer
		false, // auto-ack (false برای manual ack)
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range messages {
			if err := handler(msg.Body); err != nil {
				log.Printf("Error handling message: %v", err)
				// Negative acknowledgement (requeue)
				msg.Nack(false, true)
				continue
			}
			// Acknowledgement
			msg.Ack(false)
		}
	}()

	return nil
}

// Close بستن اتصال
func (c *RabbitMQClient) Close() error {
	if err := c.channel.Close(); err != nil {
		return err
	}
	return c.conn.Close()
}

// ============================================================================
// بخش 3: RabbitMQ - Producer و Consumer مثال
// ============================================================================

// OrderProducer تولید کننده سفارشات
type OrderProducer struct {
	client *RabbitMQClient
}

func NewOrderProducer(client *RabbitMQClient) *OrderProducer {
	return &OrderProducer{client: client}
}

func (p *OrderProducer) SendOrder(ctx context.Context, order OrderEvent) error {
	order.Timestamp = time.Now()
	return p.client.Publish(ctx, "order.created", order)
}

// PaymentConsumer مصرف کننده پرداخت
type PaymentConsumer struct {
	client *RabbitMQClient
}

func NewPaymentConsumer(client *RabbitMQClient) *PaymentConsumer {
	return &PaymentConsumer{client: client}
}

func (c *PaymentConsumer) Start() error {
	return c.client.Consume("payment.queue", func(data []byte) error {
		var payment PaymentEvent
		if err := json.Unmarshal(data, &payment); err != nil {
			return err
		}

		log.Printf("Processing payment: %+v", payment)
		// پردازش پرداخت
		// ...

		return nil
	})
}

// ============================================================================
// بخش 4: RabbitMQ - Worker Pool (برای پردازش موازی)
// ============================================================================

// WorkerPool Worker pool با RabbitMQ
type RabbitMQWorkerPool struct {
	client  *RabbitMQClient
	workers int
	queue   string
	handler func([]byte) error
	wg      sync.WaitGroup
	stopCh  chan struct{}
}

func NewRabbitMQWorkerPool(client *RabbitMQClient, workers int, queue string, handler func([]byte) error) *RabbitMQWorkerPool {
	return &RabbitMQWorkerPool{
		client:  client,
		workers: workers,
		queue:   queue,
		handler: handler,
		stopCh:  make(chan struct{}),
	}
}

func (p *RabbitMQWorkerPool) Start() error {
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go p.worker(i)
	}
	return nil
}

func (p *RabbitMQWorkerPool) worker(id int) {
	defer p.wg.Done()

	messages, err := p.client.channel.Consume(
		p.queue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("Worker %d error: %v", id, err)
		return
	}

	for {
		select {
		case msg := <-messages:
			if err := p.handler(msg.Body); err != nil {
				msg.Nack(false, true)
			} else {
				msg.Ack(false)
			}
		case <-p.stopCh:
			return
		}
	}
}

func (p *RabbitMQWorkerPool) Stop() {
	close(p.stopCh)
	p.wg.Wait()
}

// ============================================================================
// بخش 5: Apache Kafka - اتصال و تنظیمات پایه
// ============================================================================

// KafkaConfig تنظیمات Kafka
type KafkaConfig struct {
	Brokers   []string
	Topic     string
	GroupID   string
	Partition int
}

// DefaultKafkaConfig تنظیمات پیش‌فرض
func DefaultKafkaConfig() *KafkaConfig {
	return &KafkaConfig{
		Brokers:   []string{"localhost:9092"},
		Topic:     "events.topic",
		GroupID:   "events.group",
		Partition: 0,
	}
}

// KafkaProducer تولید کننده Kafka
type KafkaProducer struct {
	writer *kafka.Writer
	config *KafkaConfig
}

// NewKafkaProducer ایجاد تولید کننده Kafka
func NewKafkaProducer(config *KafkaConfig) *KafkaProducer {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(config.Brokers...),
		Topic:        config.Topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
		Async:        false,
		Compression:  kafka.Snappy,
		BatchSize:    100,
		BatchTimeout: 100 * time.Millisecond,
	}

	return &KafkaProducer{
		writer: writer,
		config: config,
	}
}

// Publish ارسال پیام به Kafka
func (p *KafkaProducer) Publish(ctx context.Context, key string, value interface{}) error {
	body, err := json.Marshal(value)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Key:   []byte(key),
		Value: body,
		Time:  time.Now(),
		Headers: []kafka.Header{
			{Key: "content-type", Value: []byte("application/json")},
		},
	}

	return p.writer.WriteMessages(ctx, msg)
}

// PublishBatch ارسال دسته‌ای پیام‌ها
func (p *KafkaProducer) PublishBatch(ctx context.Context, messages []kafka.Message) error {
	return p.writer.WriteMessages(ctx, messages...)
}

// Close بستن تولید کننده
func (p *KafkaProducer) Close() error {
	return p.writer.Close()
}

// KafkaConsumer مصرف کننده Kafka
type KafkaConsumer struct {
	reader  *kafka.Reader
	config  *KafkaConfig
	handler func(kafka.Message) error
	stopCh  chan struct{}
	wg      sync.WaitGroup
}

// NewKafkaConsumer ایجاد مصرف کننده Kafka
func NewKafkaConsumer(config *KafkaConfig, handler func(kafka.Message) error) *KafkaConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     config.Brokers,
		Topic:       config.Topic,
		GroupID:     config.GroupID,
		Partition:   config.Partition,
		MinBytes:    10e3, // 10KB
		MaxBytes:    10e6, // 10MB
		MaxWait:     1 * time.Second,
		StartOffset: kafka.LastOffset,
	})

	return &KafkaConsumer{
		reader:  reader,
		config:  config,
		handler: handler,
		stopCh:  make(chan struct{}),
	}
}

// Start شروع مصرف
func (c *KafkaConsumer) Start() {
	c.wg.Add(1)
	go c.consume()
}

func (c *KafkaConsumer) consume() {
	defer c.wg.Done()

	for {
		select {
		case <-c.stopCh:
			return
		default:
			msg, err := c.reader.ReadMessage(context.Background())
			if err != nil {
				log.Printf("Error reading message: %v", err)
				continue
			}

			if err := c.handler(msg); err != nil {
				log.Printf("Error handling message: %v", err)
			}
		}
	}
}

// Stop توقف مصرف
func (c *KafkaConsumer) Stop() error {
	close(c.stopCh)
	c.wg.Wait()
	return c.reader.Close()
}

// ============================================================================
// بخش 6: Kafka - Partitioned Consumer (مصرف از partition خاص)
// ============================================================================

// PartitionConsumer مصرف کننده از partition خاص
type PartitionConsumer struct {
	reader    *kafka.Reader
	partition int
	handler   func(kafka.Message) error
}

func NewPartitionConsumer(config *KafkaConfig, partition int, handler func(kafka.Message) error) *PartitionConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     config.Brokers,
		Topic:       config.Topic,
		Partition:   partition,
		MinBytes:    10e3,
		MaxBytes:    10e6,
		StartOffset: kafka.LastOffset,
	})

	return &PartitionConsumer{
		reader:    reader,
		partition: partition,
		handler:   handler,
	}
}

func (c *PartitionConsumer) Consume(ctx context.Context) error {
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			return err
		}

		if err := c.handler(msg); err != nil {
			log.Printf("Error handling message from partition %d: %v", c.partition, err)
		}
	}
}

func (c *PartitionConsumer) Close() error {
	return c.reader.Close()
}

// ============================================================================
// بخش 7: Kafka - Consumer Group (برای مصرف موازی)
// ============================================================================

// ConsumerGroupHandler مدیریت consumer group
type ConsumerGroupHandler struct {
	handler func(kafka.Message) error
}

func (h *ConsumerGroupHandler) Setup(session kafka.ConsumerGroupSession) error {
	log.Printf("Consumer group setup: %s", session.MemberID())
	return nil
}

func (h *ConsumerGroupHandler) Cleanup(session kafka.ConsumerGroupSession) error {
	log.Printf("Consumer group cleanup: %s", session.MemberID())
	return nil
}

func (h *ConsumerGroupHandler) ConsumeClaim(session kafka.ConsumerGroupSession, claim kafka.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		if err := h.handler(msg); err != nil {
			log.Printf("Error handling message: %v", err)
		}
		session.MarkMessage(msg, "")
	}
	return nil
}

// ConsumerGroup مصرف کننده گروهی
type ConsumerGroup struct {
	group   kafka.ConsumerGroup
	handler *ConsumerGroupHandler
	config  *KafkaConfig
}

func NewConsumerGroup(config *KafkaConfig, handler func(kafka.Message) error) (*ConsumerGroup, error) {
	group, err := kafka.NewConsumerGroup(kafka.ConsumerGroupConfig{
		ID:      config.GroupID,
		Brokers: config.Brokers,
		Topics:  []string{config.Topic},
	})
	if err != nil {
		return nil, err
	}

	return &ConsumerGroup{
		group:   group,
		handler: &ConsumerGroupHandler{handler: handler},
		config:  config,
	}, nil
}

func (cg *ConsumerGroup) Start(ctx context.Context) error {
	return cg.group.Consume(ctx, []string{cg.config.Topic}, cg.handler)
}

func (cg *ConsumerGroup) Close() error {
	return cg.group.Close()
}

// ============================================================================
// بخش 8: Event Processing Pipeline (مثال یکپارچه)
// ============================================================================

// EventPipeline خط لوله پردازش رویداد
type EventPipeline struct {
	rabbitProducer *RabbitMQClient
	kafkaProducer  *KafkaProducer
	eventHandlers  map[string]func(Event) error
}

func NewEventPipeline(rabbitConfig *RabbitMQConfig, kafkaConfig *KafkaConfig) (*EventPipeline, error) {
	rabbitClient, err := NewRabbitMQClient(rabbitConfig)
	if err != nil {
		return nil, err
	}

	kafkaProducer := NewKafkaProducer(kafkaConfig)

	// همچنین می‌توان مصرف کننده Kafka را تنظیم کرد
	// kafkaConsumer := NewKafkaConsumer(kafkaConfig, pipeline.handleEvent)

	return &EventPipeline{
		rabbitProducer: rabbitClient,
		kafkaProducer:  kafkaProducer,
		eventHandlers:  make(map[string]func(Event) error),
	}, nil
}

// RegisterHandler ثبت handler برای نوع خاص رویداد
func (p *EventPipeline) RegisterHandler(eventType string, handler func(Event) error) {
	p.eventHandlers[eventType] = handler
}

// ProcessEvent پردازش رویداد
func (p *EventPipeline) ProcessEvent(event Event) error {
	handler, exists := p.eventHandlers[event.Type]
	if !exists {
		return fmt.Errorf("no handler for event type: %s", event.Type)
	}
	return handler(event)
}

// EmitEvent ارسال رویداد (با RabbitMQ و Kafka)
func (p *EventPipeline) EmitEvent(ctx context.Context, event Event) error {
	// ارسال به RabbitMQ برای پردازش همزمان
	if err := p.rabbitProducer.Publish(ctx, event.Type, event); err != nil {
		log.Printf("Failed to publish to RabbitMQ: %v", err)
	}

	// ارسال به Kafka برای ذخیره و تحلیل
	if err := p.kafkaProducer.Publish(ctx, event.ID, event); err != nil {
		log.Printf("Failed to publish to Kafka: %v", err)
	}

	return nil
}

// ============================================================================
// بخش 9: Dead Letter Queue (RabbitMQ)
// ============================================================================

// DeadLetterConfig تنظیمات Dead Letter Queue
type DeadLetterConfig struct {
	Exchange   string
	Queue      string
	RoutingKey string
	TTL        int // milliseconds
	MaxRetries int
}

// SetupDeadLetterQueue تنظیم Dead Letter Queue
func SetupDeadLetterQueue(channel *amqp.Channel, config *DeadLetterConfig) error {
	// اصلی exchange
	if err := channel.ExchangeDeclare(
		config.Exchange,
		"topic",
		true, false, false, false, nil,
	); err != nil {
		return err
	}

	// Dead Letter exchange
	dlxExchange := config.Exchange + ".dlx"
	if err := channel.ExchangeDeclare(
		dlxExchange,
		"fanout",
		true, false, false, false, nil,
	); err != nil {
		return err
	}

	// اصلی queue با dead letter تنظیمات
	args := amqp.Table{
		"x-dead-letter-exchange": dlxExchange,
		"x-message-ttl":          int32(config.TTL),
		"x-max-retries":          int16(config.MaxRetries),
	}

	_, err := channel.QueueDeclare(
		config.Queue,
		true, false, false, false, args,
	)
	return err
}

// ============================================================================
// بخش 10: Data Streaming with Kafka (مثال real-time)
// ============================================================================

// StreamProcessor پردازنده جریانی
type StreamProcessor struct {
	consumer *KafkaConsumer
	window   *SlidingWindow
	metrics  map[string]int64
	mu       sync.RWMutex
}

type SlidingWindow struct {
	duration time.Duration
	events   []Event
	mu       sync.RWMutex
}

func NewSlidingWindow(duration time.Duration) *SlidingWindow {
	return &SlidingWindow{
		duration: duration,
		events:   make([]Event, 0),
	}
}

func (sw *SlidingWindow) Add(event Event) {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	sw.events = append(sw.events, event)

	// حذف رویدادهای قدیمی
	cutoff := time.Now().Add(-sw.duration)
	var valid []Event
	for _, e := range sw.events {
		if e.Timestamp.After(cutoff) {
			valid = append(valid, e)
		}
	}
	sw.events = valid
}

func (sw *SlidingWindow) Count() int {
	sw.mu.RLock()
	defer sw.mu.RUnlock()
	return len(sw.events)
}

func NewStreamProcessor(config *KafkaConfig) *StreamProcessor {
	return &StreamProcessor{
		window:  NewSlidingWindow(5 * time.Minute),
		metrics: make(map[string]int64),
	}
}

func (p *StreamProcessor) ProcessMessage(msg kafka.Message) error {
	var event Event
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		return err
	}

	// اضافه کردن به پنجره لغزنده
	p.window.Add(event)

	// به‌روزرسانی متریک‌ها
	p.mu.Lock()
	p.metrics[event.Type]++
	p.mu.Unlock()

	// تحلیل real-time
	p.analyzeEvent(event)

	return nil
}

func (p *StreamProcessor) analyzeEvent(event Event) {
	// تحلیل جریانی
	switch event.Type {
	case "order.created":
		log.Printf("New order detected: %s", event.ID)
	case "payment.failed":
		log.Printf("Payment failure alert: %s", event.ID)
		// ارسال هشدار
	}
}

func (p *StreamProcessor) GetMetrics() map[string]int64 {
	p.mu.RLock()
	defer p.mu.RUnlock()

	result := make(map[string]int64)
	for k, v := range p.metrics {
		result[k] = v
	}
	return result
}

// ============================================================================
// بخش 11: Best Practices
// ============================================================================

func bestPractices() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("💡 MESSAGE QUEUE BEST PRACTICES")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ 1. RABBITMQ BEST PRACTICES                                                │
├─────────────────────────────────────────────────────────────────────────────┤
│    • Use persistent messages for important data                            │
│    • Implement publisher confirms                                          │
│    • Use dead letter queues for failed messages                            │
│    • Set appropriate TTL for messages                                      │
│    • Monitor queue length and consumer lag                                 │
│    • Use connection and channel pooling                                    │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 2. KAFKA BEST PRACTICES                                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│    • Choose appropriate partition count                                    │
│    • Use meaningful keys for partitioning                                  │
│    • Set proper replication factor (>=2 for production)                    │
│    • Configure retention policy based on use case                          │
│    • Use consumer groups for parallel processing                           │
│    • Handle rebalancing gracefully                                         │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 3. MESSAGE DESIGN                                                         │
├─────────────────────────────────────────────────────────────────────────────┤
│    • Use idempotent message processing                                     │
│    • Include message ID and timestamp                                      │
│    • Version your messages                                                 │
│    • Keep messages small (Kafka: <1MB, RabbitMQ: < 4MB)                    │
│    • Use schemas (Avro, Protobuf) for evolution                            │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 4. ERROR HANDLING                                                         │
├─────────────────────────────────────────────────────────────────────────────┤
│    • Implement retry with exponential backoff                              │
│    • Use dead letter queues for unprocessable messages                     │
│    • Log errors with correlation ID                                        │
│    • Monitor dead letter queues                                            │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 5. MONITORING                                                             │
├─────────────────────────────────────────────────────────────────────────────┤
│    • Monitor queue/topic depth                                             │
│    • Track consumer lag                                                    │
│    • Alert on message 堆积                                                │
│    • Monitor throughput and latency                                        │
└─────────────────────────────────────────────────────────────────────────────┘
`)
}

// ============================================================================
// بخش 12: Comparison Table
// ============================================================================

func comparisonTable() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📊 RABBITMQ vs KAFKA COMPARISON")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ FEATURE              │ RABBITMQ              │ KAFKA                        │
├──────────────────────┼───────────────────────┼──────────────────────────────┤
│ Messaging Model      │ Point-to-point,       │ Publish-Subscribe,           │
│                      │ Publish-Subscribe     │ Event Streaming              │
├──────────────────────┼───────────────────────┼──────────────────────────────┤
│ Message Ordering     │ Order preserved       │ Order preserved             │
│                      │ in single queue       │ in single partition          │
├──────────────────────┼───────────────────────┼──────────────────────────────┤
│ Message Retention    │ Until consumed        │ Based on time/size          │
├──────────────────────┼───────────────────────┼──────────────────────────────┤
│ Throughput           │ ~50k msg/sec          │ ~1M msg/sec                  │
├──────────────────────┼───────────────────────┼──────────────────────────────┤
│ Latency              │ Low (microseconds)    │ Higher (milliseconds)        │
├──────────────────────┼───────────────────────┼──────────────────────────────┤
│ Message Routing      │ Flexible (exchanges)  │ Limited (topics)             │
├──────────────────────┼───────────────────────┼──────────────────────────────┤
│ Exactly-once         │ Limited               │ Yes (idempotent producer)    │
├──────────────────────┼───────────────────────┼──────────────────────────────┤
│ Best For             │ Task distribution,    │ Event sourcing,              │
│                      │ RPC, Decoupling       │ Analytics, Logs              │
└──────────────────────┴───────────────────────┴──────────────────────────────┘
`)
}

// ============================================================================
// بخش 13: Complete Example
// ============================================================================

func runCompleteExample() {
	// ========== RabbitMQ Example ==========
	log.Println("=== RabbitMQ Example ===")

	rabbitConfig := DefaultRabbitMQConfig()
	rabbitClient, err := NewRabbitMQClient(rabbitConfig)
	if err != nil {
		log.Printf("RabbitMQ connection failed: %v (make sure RabbitMQ is running)", err)
	} else {
		defer rabbitClient.Close()

		// ارسال پیام
		ctx := context.Background()
		order := OrderEvent{
			OrderID:     "ORD-001",
			UserID:      "USR-001",
			ProductID:   "PRD-001",
			Quantity:    2,
			TotalAmount: 99.98,
			Status:      "pending",
		}

		if err := rabbitClient.Publish(ctx, "order.created", order); err != nil {
			log.Printf("Failed to publish: %v", err)
		} else {
			log.Println("Message published to RabbitMQ")
		}
	}

	// ========== Kafka Example ==========
	log.Println("\n=== Kafka Example ===")

	kafkaConfig := DefaultKafkaConfig()
	kafkaProducer := NewKafkaProducer(kafkaConfig)
	defer kafkaProducer.Close()

	event := Event{
		ID:        "EVT-001",
		Type:      "user.login",
		Timestamp: time.Now(),
		Source:    "auth-service",
		Data: map[string]interface{}{
			"user_id": "123",
			"ip":      "192.168.1.1",
		},
	}

	if err := kafkaProducer.Publish(context.Background(), event.ID, event); err != nil {
		log.Printf("Kafka publish failed: %v (make sure Kafka is running)", err)
	} else {
		log.Println("Message published to Kafka")
	}

	// ========== Pipeline Example ==========
	log.Println("\n=== Pipeline Example ===")

	// تنظیم pipeline
	eventPipeline, err := NewEventPipeline(rabbitConfig, kafkaConfig)
	if err != nil {
		log.Printf("Pipeline creation failed: %v", err)
	} else {
		// ثبت handler
		eventPipeline.RegisterHandler("order.created", func(e Event) error {
			log.Printf("Processing order event: %+v", e)
			return nil
		})

		// ارسال رویداد
		testEvent := Event{
			ID:        "EVT-002",
			Type:      "order.created",
			Timestamp: time.Now(),
			Source:    "api-gateway",
			Data:      order,
		}

		eventPipeline.EmitEvent(context.Background(), testEvent)
	}
}

// ============================================================================
// بخش 14: Main
// ============================================================================

func main() {
	fmt.Println(stringsRepeat("=", 80))
	fmt.Println("🎯 RABBITMQ & APACHE KAFKA GUIDE")
	fmt.Println("Message Queue & Event Streaming in Go")
	fmt.Println(stringsRepeat("=", 80))

	bestPractices()
	comparisonTable()

	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🚀 Running Examples")
	fmt.Println(stringsRepeat("=", 80))

	// اجرای مثال‌ها در گوروتین جدا (برای عدم blocking)
	go runCompleteExample()

	// منتظر سیگنال برای خروج
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down...")
}

// تابع کمکی
func stringsRepeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}
