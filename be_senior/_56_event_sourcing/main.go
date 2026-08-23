// ============================================================================
// FILE: event_sourcing_guide.go
// TITLE: راهنمای کامل Event Sourcing در Go - رویدادمحور و بازسازی وضعیت
// HOW TO RUN: go run event_sourcing_guide.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - Event Sourcing چیست؟
// ============================================================================
//
// Event Sourcing یک الگوی معماری است که به جای ذخیره وضعیت فعلی،
// تمام تغییرات (رویدادها) را ذخیره می‌کند و وضعیت فعلی از بازپخش رویدادها ساخته می‌شود.
//
// مفاهیم کلیدی:
// - Event (رویداد): چیزی که در سیستم رخ داده است (غیرقابل تغییر)
// - Aggregate: تجمعی از رویدادها که وضعیت فعلی یک Entity را نشان می‌دهد
// - Event Store: جایی که رویدادها ذخیره می‌شوند (append-only)
// - Snapshot: تصویر لحظه‌ای از وضعیت برای بهبود بازپخش
// - Projection: نمایش داده شده از رویدادها برای query
//
// مزایا:
// 1. تاریخچه کامل: همه تغییرات ثبت می‌شود
// 2. بازپخش (Replay): می‌توان وضعیت را در هر نقطه‌ای بازسازی کرد
// 3. وارونگی (Rollback): با replay معکوس می‌توان تغییرات را برگرداند
// 4. تحلیل: می‌توان رویدادها را برای تحلیل استفاده کرد
// 5. Event-Driven: به راحتی می‌توان با CQRS و Event-Driven Architecture ترکیب کرد
//
// معایب:
// 1. پیچیدگی بیشتر
// 2. نیاز به فضای ذخیره‌سازی بیشتر
// 3. eventual consistency
//
// قانون طلایی:
// "رویدادها را به صورت append-only ذخیره کن، هرگز رویدادها را تغییر یا حذف نکن.
//  از Snapshots برای بهبود بازپخش استفاده کن.
//  رویدادها را version بندی کن تا امکان تغییر schema داشته باشی."
// ============================================================================

package __event_sourcing

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
)

// ============================================================================
// بخش 1: Core Types (انواع پایه)
// ============================================================================

// Event رویداد پایه
type Event interface {
	GetID() string
	GetAggregateID() string
	GetType() string
	GetVersion() int
	GetTimestamp() time.Time
	GetData() interface{}
}

// BaseEvent پیاده‌سازی پایه Event
type BaseEvent struct {
	ID          string      `json:"id"`
	AggregateID string      `json:"aggregate_id"`
	Type        string      `json:"type"`
	Version     int         `json:"version"`
	Timestamp   time.Time   `json:"timestamp"`
	Data        interface{} `json:"data"`
}

func (e BaseEvent) GetID() string           { return e.ID }
func (e BaseEvent) GetAggregateID() string  { return e.AggregateID }
func (e BaseEvent) GetType() string         { return e.Type }
func (e BaseEvent) GetVersion() int         { return e.Version }
func (e BaseEvent) GetTimestamp() time.Time { return e.Timestamp }
func (e BaseEvent) GetData() interface{}    { return e.Data }

// EventStore ذخیره‌سازی رویدادها
type EventStore interface {
	Save(ctx context.Context, events []Event, expectedVersion int) error
	Load(ctx context.Context, aggregateID string) ([]Event, error)
	LoadFromVersion(ctx context.Context, aggregateID string, version int) ([]Event, error)
}

// Snapshot تصویر لحظه‌ای از Aggregate
type Snapshot struct {
	AggregateID   string      `json:"aggregate_id"`
	AggregateType string      `json:"aggregate_type"`
	State         interface{} `json:"state"`
	Version       int         `json:"version"`
	Timestamp     time.Time   `json:"timestamp"`
}

// ============================================================================
// بخش 2: Domain Events (رویدادهای دامنه)
// ============================================================================

// رویدادهای سفارش (Order Events)

// OrderCreated رویداد ایجاد سفارش
type OrderCreatedData struct {
	UserID    string  `json:"user_id"`
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	Total     float64 `json:"total"`
}

type OrderCreated struct {
	BaseEvent
	Data OrderCreatedData `json:"data"`
}

func NewOrderCreated(aggregateID string, data OrderCreatedData) *OrderCreated {
	return &OrderCreated{
		BaseEvent: BaseEvent{
			ID:          uuid.New().String(),
			AggregateID: aggregateID,
			Type:        "OrderCreated",
			Version:     1,
			Timestamp:   time.Now(),
		},
		Data: data,
	}
}

// OrderStatusChanged رویداد تغییر وضعیت سفارش
type OrderStatusChangedData struct {
	OldStatus string `json:"old_status"`
	NewStatus string `json:"new_status"`
	Reason    string `json:"reason,omitempty"`
}

type OrderStatusChanged struct {
	BaseEvent
	Data OrderStatusChangedData `json:"data"`
}

func NewOrderStatusChanged(aggregateID string, version int, data OrderStatusChangedData) *OrderStatusChanged {
	return &OrderStatusChanged{
		BaseEvent: BaseEvent{
			ID:          uuid.New().String(),
			AggregateID: aggregateID,
			Type:        "OrderStatusChanged",
			Version:     version,
			Timestamp:   time.Now(),
		},
		Data: data,
	}
}

// OrderItemAdded رویداد افزودن آیتم به سفارش
type OrderItemAddedData struct {
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type OrderItemAdded struct {
	BaseEvent
	Data OrderItemAddedData `json:"data"`
}

func NewOrderItemAdded(aggregateID string, version int, data OrderItemAddedData) *OrderItemAdded {
	return &OrderItemAdded{
		BaseEvent: BaseEvent{
			ID:          uuid.New().String(),
			AggregateID: aggregateID,
			Type:        "OrderItemAdded",
			Version:     version,
			Timestamp:   time.Now(),
		},
		Data: data,
	}
}

// OrderItemRemoved رویداد حذف آیتم از سفارش
type OrderItemRemovedData struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
	Reason    string `json:"reason,omitempty"`
}

type OrderItemRemoved struct {
	BaseEvent
	Data OrderItemRemovedData `json:"data"`
}

// PaymentReceived رویداد دریافت پرداخت
type PaymentReceivedData struct {
	PaymentID     string  `json:"payment_id"`
	Amount        float64 `json:"amount"`
	PaymentMethod string  `json:"payment_method"`
}

type PaymentReceived struct {
	BaseEvent
	Data PaymentReceivedData `json:"data"`
}

func NewPaymentReceived(aggregateID string, version int, data PaymentReceivedData) *PaymentReceived {
	return &PaymentReceived{
		BaseEvent: BaseEvent{
			ID:          uuid.New().String(),
			AggregateID: aggregateID,
			Type:        "PaymentReceived",
			Version:     version,
			Timestamp:   time.Now(),
		},
		Data: data,
	}
}

// ============================================================================
// بخش 3: Aggregate (مدیریت وضعیت)
// ============================================================================

// Order Aggregate برای سفارش
type Order struct {
	ID         string
	UserID     string
	Items      map[string]*OrderItem
	Total      float64
	Status     string
	Version    int
	Events     []Event
	eventMutex sync.RWMutex
}

type OrderItem struct {
	ProductID string
	Quantity  int
	Price     float64
	Subtotal  float64
}

// NewOrder ایجاد سفارش جدید
func NewOrder(id, userID string) *Order {
	return &Order{
		ID:      id,
		UserID:  userID,
		Items:   make(map[string]*OrderItem),
		Status:  "pending",
		Version: 0,
		Events:  make([]Event, 0),
	}
}

// ApplyEvent اعمال رویداد روی Aggregate
func (o *Order) ApplyEvent(event Event) {
	switch e := event.(type) {
	case *OrderCreated:
		o.UserID = e.Data.UserID
		o.Status = "pending"
		o.Total = 0

	case *OrderItemAdded:
		item := &OrderItem{
			ProductID: e.Data.ProductID,
			Quantity:  e.Data.Quantity,
			Price:     e.Data.Price,
			Subtotal:  float64(e.Data.Quantity) * e.Data.Price,
		}
		o.Items[e.Data.ProductID] = item
		o.recalculateTotal()

	case *OrderItemRemoved:
		delete(o.Items, e.Data.ProductID)
		o.recalculateTotal()

	case *OrderStatusChanged:
		o.Status = e.Data.NewStatus

	case *PaymentReceived:
		// پردازش پرداخت
	}

	o.Version = event.GetVersion()
}

// recalculateTotal محاسبه مجدد مجموع
func (o *Order) recalculateTotal() {
	total := 0.0
	for _, item := range o.Items {
		total += item.Subtotal
	}
	o.Total = total
}

// AddItem افزودن آیتم به سفارش
func (o *Order) AddItem(productID string, quantity int, price float64) error {
	if o.Status != "pending" {
		return errors.New("cannot add items to non-pending order")
	}

	if quantity <= 0 {
		return errors.New("quantity must be positive")
	}

	if price <= 0 {
		return errors.New("price must be positive")
	}

	event := NewOrderItemAdded(o.ID, o.Version+1, OrderItemAddedData{
		ProductID: productID,
		Quantity:  quantity,
		Price:     price,
	})

	o.addEvent(event)
	return nil
}

// RemoveItem حذف آیتم از سفارش
func (o *Order) RemoveItem(productID string, quantity int, reason string) error {
	if o.Status != "pending" {
		return errors.New("cannot remove items from non-pending order")
	}

	item, exists := o.Items[productID]
	if !exists {
		return errors.New("item not found")
	}

	if quantity > item.Quantity {
		quantity = item.Quantity
	}

	event := &OrderItemRemoved{
		BaseEvent: BaseEvent{
			ID:          uuid.New().String(),
			AggregateID: o.ID,
			Type:        "OrderItemRemoved",
			Version:     o.Version + 1,
			Timestamp:   time.Now(),
		},
		Data: OrderItemRemovedData{
			ProductID: productID,
			Quantity:  quantity,
			Reason:    reason,
		},
	}

	o.addEvent(event)
	return nil
}

// ChangeStatus تغییر وضعیت سفارش
func (o *Order) ChangeStatus(newStatus, reason string) error {
	validStatuses := map[string]bool{
		"pending": true, "processing": true, "shipped": true,
		"delivered": true, "cancelled": true, "refunded": true,
	}

	if !validStatuses[newStatus] {
		return errors.New("invalid status")
	}

	event := NewOrderStatusChanged(o.ID, o.Version+1, OrderStatusChangedData{
		OldStatus: o.Status,
		NewStatus: newStatus,
		Reason:    reason,
	})

	o.addEvent(event)
	return nil
}

// ReceivePayment ثبت پرداخت
func (o *Order) ReceivePayment(paymentID, paymentMethod string, amount float64) error {
	if o.Status != "pending" && o.Status != "processing" {
		return errors.New("cannot receive payment for this order status")
	}

	if amount <= 0 {
		return errors.New("amount must be positive")
	}

	if amount > o.Total {
		return errors.New("amount exceeds order total")
	}

	event := NewPaymentReceived(o.ID, o.Version+1, PaymentReceivedData{
		PaymentID:     paymentID,
		Amount:        amount,
		PaymentMethod: paymentMethod,
	})

	o.addEvent(event)
	return nil
}

// addEvent افزودن رویداد به Aggregate
func (o *Order) addEvent(event Event) {
	o.eventMutex.Lock()
	defer o.eventMutex.Unlock()
	o.Events = append(o.Events, event)
	o.ApplyEvent(event)
}

// GetUncommittedEvents دریافت رویدادهای ذخیره نشده
func (o *Order) GetUncommittedEvents() []Event {
	o.eventMutex.RLock()
	defer o.eventMutex.RUnlock()
	return o.Events
}

// MarkEventsAsCommitted علامت ذخیره شدن رویدادها
func (o *Order) MarkEventsAsCommitted() {
	o.eventMutex.Lock()
	defer o.eventMutex.Unlock()
	o.Events = make([]Event, 0)
}

// LoadFromHistory بازسازی وضعیت از تاریخچه رویدادها
func (o *Order) LoadFromHistory(events []Event) {
	for _, event := range events {
		o.ApplyEvent(event)
	}
	o.Events = make([]Event, 0)
}

// ============================================================================
// بخش 4: In-Memory Event Store (برای مثال)
// ============================================================================

// InMemoryEventStore ذخیره‌سازی رویداد در حافظه
type InMemoryEventStore struct {
	events    map[string][]Event // aggregateID -> events
	snapshots map[string]Snapshot
	mu        sync.RWMutex
}

func NewInMemoryEventStore() *InMemoryEventStore {
	return &InMemoryEventStore{
		events:    make(map[string][]Event),
		snapshots: make(map[string]Snapshot),
	}
}

func (s *InMemoryEventStore) Save(ctx context.Context, events []Event, expectedVersion int) error {
	if len(events) == 0 {
		return nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	aggregateID := events[0].GetAggregateID()

	// بررسی نسخه مورد انتظار
	existingEvents := s.events[aggregateID]
	if expectedVersion != len(existingEvents) {
		return errors.New("concurrency conflict: unexpected version")
	}

	// ذخیره رویدادها
	s.events[aggregateID] = append(existingEvents, events...)
	return nil
}

func (s *InMemoryEventStore) Load(ctx context.Context, aggregateID string) ([]Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	events, ok := s.events[aggregateID]
	if !ok {
		return []Event{}, nil
	}

	return events, nil
}

func (s *InMemoryEventStore) LoadFromVersion(ctx context.Context, aggregateID string, version int) ([]Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	events, ok := s.events[aggregateID]
	if !ok {
		return []Event{}, nil
	}

	if version < 0 || version > len(events) {
		return nil, errors.New("invalid version")
	}

	return events[version:], nil
}

// SaveSnapshot ذخیره snapshot
func (s *InMemoryEventStore) SaveSnapshot(snapshot Snapshot) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.snapshots[snapshot.AggregateID] = snapshot
	return nil
}

// LoadSnapshot بارگذاری snapshot
func (s *InMemoryEventStore) LoadSnapshot(aggregateID string) (*Snapshot, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	snapshot, ok := s.snapshots[aggregateID]
	if !ok {
		return nil, errors.New("snapshot not found")
	}
	return &snapshot, nil
}

// ============================================================================
// بخش 5: Repository (با پشتیبانی از Snapshot)
// ============================================================================

// OrderRepository مخزن سفارشات
type OrderRepository struct {
	eventStore   EventStore
	snapshotFreq int // هر چند رویداد یک snapshot بگیرد
}

func NewOrderRepository(eventStore EventStore, snapshotFreq int) *OrderRepository {
	return &OrderRepository{
		eventStore:   eventStore,
		snapshotFreq: snapshotFreq,
	}
}

// Save ذخیره سفارش
func (r *OrderRepository) Save(ctx context.Context, order *Order) error {
	events := order.GetUncommittedEvents()
	if len(events) == 0 {
		return nil
	}

	// ذخیره رویدادها
	if err := r.eventStore.Save(ctx, events, order.Version-len(events)); err != nil {
		return err
	}

	// ذخیره snapshot در صورت نیاز
	if order.Version > 0 && order.Version%r.snapshotFreq == 0 {
		snapshot := Snapshot{
			AggregateID:   order.ID,
			AggregateType: "Order",
			State:         order,
			Version:       order.Version,
			Timestamp:     time.Now(),
		}

		if store, ok := r.eventStore.(*InMemoryEventStore); ok {
			store.SaveSnapshot(snapshot)
		}
	}

	order.MarkEventsAsCommitted()
	return nil
}

// Load بارگذاری سفارش
func (r *OrderRepository) Load(ctx context.Context, id string) (*Order, error) {
	var events []Event
	var startVersion int = 0

	// تلاش برای بارگذاری snapshot
	if store, ok := r.eventStore.(*InMemoryEventStore); ok {
		if snapshot, err := store.LoadSnapshot(id); err == nil {
			if order, ok := snapshot.State.(*Order); ok {
				order.Events = make([]Event, 0)
				startVersion = snapshot.Version
				events, _ = r.eventStore.LoadFromVersion(ctx, id, startVersion)
				order.LoadFromHistory(events)
				return order, nil
			}
		}
	}

	// بارگذاری از اول
	events, err := r.eventStore.Load(ctx, id)
	if err != nil {
		return nil, err
	}

	if len(events) == 0 {
		return nil, errors.New("order not found")
	}

	// بازسازی وضعیت
	order := NewOrder(id, "")
	order.LoadFromHistory(events)

	return order, nil
}

// ============================================================================
// بخش 6: Event Handlers و Projections (نمایش داده‌ها)
// ============================================================================

// OrderProjection نمایش سفارشات برای query
type OrderProjection struct {
	mu     sync.RWMutex
	orders map[string]*OrderReadModel
}

type OrderReadModel struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Items     []ItemDTO `json:"items"`
	Total     float64   `json:"total"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ItemDTO struct {
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	Subtotal  float64 `json:"subtotal"`
}

func NewOrderProjection() *OrderProjection {
	return &OrderProjection{
		orders: make(map[string]*OrderReadModel),
	}
}

// HandleEvent پردازش رویداد برای به‌روزرسانی projection
func (p *OrderProjection) HandleEvent(event Event) {
	p.mu.Lock()
	defer p.mu.Unlock()

	switch e := event.(type) {
	case *OrderCreated:
		p.orders[e.AggregateID] = &OrderReadModel{
			ID:        e.AggregateID,
			UserID:    e.Data.UserID,
			Items:     []ItemDTO{},
			Total:     0,
			Status:    "pending",
			CreatedAt: e.Timestamp,
			UpdatedAt: e.Timestamp,
		}

	case *OrderItemAdded:
		if order, ok := p.orders[e.AggregateID]; ok {
			order.Items = append(order.Items, ItemDTO{
				ProductID: e.Data.ProductID,
				Quantity:  e.Data.Quantity,
				Price:     e.Data.Price,
				Subtotal:  float64(e.Data.Quantity) * e.Data.Price,
			})
			order.recalculateTotal()
			order.UpdatedAt = e.Timestamp
		}

	case *OrderStatusChanged:
		if order, ok := p.orders[e.AggregateID]; ok {
			order.Status = e.Data.NewStatus
			order.UpdatedAt = e.Timestamp
		}
	}
}

func (p *OrderProjection) GetOrder(id string) (*OrderReadModel, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	order, ok := p.orders[id]
	return order, ok
}

func (p *OrderProjection) GetUserOrders(userID string) []*OrderReadModel {
	p.mu.RLock()
	defer p.mu.RUnlock()

	var result []*OrderReadModel
	for _, order := range p.orders {
		if order.UserID == userID {
			result = append(result, order)
		}
	}
	return result
}

func (o *OrderReadModel) recalculateTotal() {
	total := 0.0
	for _, item := range o.Items {
		total += item.Subtotal
	}
	o.Total = total
}

// ============================================================================
// بخش 7: Rebuilder (بازسازی Projection)
// ============================================================================

// ProjectionRebuilder بازسازی کننده projection
type ProjectionRebuilder struct {
	eventStore EventStore
	handlers   map[string][]func(Event)
}

func NewProjectionRebuilder(eventStore EventStore) *ProjectionRebuilder {
	return &ProjectionRebuilder{
		eventStore: eventStore,
		handlers:   make(map[string][]func(Event)),
	}
}

// RegisterHandler ثبت handler برای نوع رویداد
func (r *ProjectionRebuilder) RegisterHandler(eventType string, handler func(Event)) {
	r.handlers[eventType] = append(r.handlers[eventType], handler)
}

// Rebuild بازسازی از ابتدا
func (r *ProjectionRebuilder) Rebuild(ctx context.Context) error {
	// دریافت همه aggregateIDها (در عمل باید از دیتابیس کوئری زد)
	// برای مثال، همه رویدادها را بگیریم

	// پیمایش همه رویدادها و اعمال handlerها
	// (در عمل بهتر است از offset/cursor استفاده کرد)

	return nil
}

// ============================================================================
// بخش 8: Event Publisher (برای انتشار رویدادها)
// ============================================================================

// EventPublisher انتشار‌دهنده رویداد
type EventPublisher struct {
	subscribers map[string][]func(Event)
	mu          sync.RWMutex
}

func NewEventPublisher() *EventPublisher {
	return &EventPublisher{
		subscribers: make(map[string][]func(Event)),
	}
}

// Subscribe اشتراک در رویدادها
func (p *EventPublisher) Subscribe(eventType string, handler func(Event)) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.subscribers[eventType] = append(p.subscribers[eventType], handler)
}

// Publish انتشار رویداد
func (p *EventPublisher) Publish(event Event) {
	p.mu.RLock()
	handlers := p.subscribers[event.GetType()]
	p.mu.RUnlock()

	for _, handler := range handlers {
		go handler(event) // غیرهمزمان برای عدم blocking
	}
}

// ============================================================================
// بخش 9: Command Handler (مدیریت دستورات)
// ============================================================================

// CommandHandler مدیریت دستورات
type CommandHandler struct {
	repository     *OrderRepository
	eventPublisher *EventPublisher
}

func NewCommandHandler(repo *OrderRepository, publisher *EventPublisher) *CommandHandler {
	return &CommandHandler{
		repository:     repo,
		eventPublisher: publisher,
	}
}

// CreateOrderCommand دستور ایجاد سفارش
type CreateOrderCommand struct {
	OrderID   string
	UserID    string
	ProductID string
	Quantity  int
	Price     float64
}

func (h *CommandHandler) HandleCreateOrder(ctx context.Context, cmd CreateOrderCommand) error {
	// ایجاد Aggregate جدید
	order := NewOrder(cmd.OrderID, cmd.UserID)

	// افزودن آیتم
	if err := order.AddItem(cmd.ProductID, cmd.Quantity, cmd.Price); err != nil {
		return err
	}

	// ذخیره
	if err := h.repository.Save(ctx, order); err != nil {
		return err
	}

	// انتشار رویدادها
	for _, event := range order.GetUncommittedEvents() {
		h.eventPublisher.Publish(event)
	}

	return nil
}

// AddItemCommand دستور افزودن آیتم
type AddItemCommand struct {
	OrderID   string
	ProductID string
	Quantity  int
	Price     float64
}

func (h *CommandHandler) HandleAddItem(ctx context.Context, cmd AddItemCommand) error {
	order, err := h.repository.Load(ctx, cmd.OrderID)
	if err != nil {
		return err
	}

	if err := order.AddItem(cmd.ProductID, cmd.Quantity, cmd.Price); err != nil {
		return err
	}

	if err := h.repository.Save(ctx, order); err != nil {
		return err
	}

	for _, event := range order.GetUncommittedEvents() {
		h.eventPublisher.Publish(event)
	}

	return nil
}

// ChangeOrderStatusCommand دستور تغییر وضعیت
type ChangeOrderStatusCommand struct {
	OrderID string
	Status  string
	Reason  string
}

func (h *CommandHandler) HandleChangeStatus(ctx context.Context, cmd ChangeOrderStatusCommand) error {
	order, err := h.repository.Load(ctx, cmd.OrderID)
	if err != nil {
		return err
	}

	if err := order.ChangeStatus(cmd.Status, cmd.Reason); err != nil {
		return err
	}

	if err := h.repository.Save(ctx, order); err != nil {
		return err
	}

	for _, event := range order.GetUncommittedEvents() {
		h.eventPublisher.Publish(event)
	}

	return nil
}

// ============================================================================
// بخش 10: Complete Example
// ============================================================================

func runEventSourcingExample() {
	log.Println("=== Event Sourcing Example ===")

	// ایجاد components
	eventStore := NewInMemoryEventStore()
	repository := NewOrderRepository(eventStore, 5)
	publisher := NewEventPublisher()
	projection := NewOrderProjection()

	// اشتراک projection در رویدادها
	publisher.Subscribe("OrderCreated", projection.HandleEvent)
	publisher.Subscribe("OrderItemAdded", projection.HandleEvent)
	publisher.Subscribe("OrderStatusChanged", projection.HandleEvent)

	// ایجاد command handler
	commandHandler := NewCommandHandler(repository, publisher)

	ctx := context.Background()
	orderID := uuid.New().String()

	// 1. ایجاد سفارش
	log.Println("\n1. Creating order...")
	err := commandHandler.HandleCreateOrder(ctx, CreateOrderCommand{
		OrderID:   orderID,
		UserID:    "user-123",
		ProductID: "prod-001",
		Quantity:  2,
		Price:     49.99,
	})
	if err != nil {
		log.Printf("Error creating order: %v", err)
		return
	}

	// 2. افزودن آیتم دیگر
	log.Println("\n2. Adding another item...")
	err = commandHandler.HandleAddItem(ctx, AddItemCommand{
		OrderID:   orderID,
		ProductID: "prod-002",
		Quantity:  1,
		Price:     29.99,
	})
	if err != nil {
		log.Printf("Error adding item: %v", err)
	}

	// 3. تغییر وضعیت
	log.Println("\n3. Changing order status...")
	err = commandHandler.HandleChangeStatus(ctx, ChangeOrderStatusCommand{
		OrderID: orderID,
		Status:  "processing",
		Reason:  "Payment received",
	})
	if err != nil {
		log.Printf("Error changing status: %v", err)
	}

	// 4. نمایش وضعیت نهایی
	log.Println("\n4. Final order state:")
	order, err := repository.Load(ctx, orderID)
	if err != nil {
		log.Printf("Error loading order: %v", err)
	} else {
		log.Printf("  Order ID: %s", order.ID)
		log.Printf("  User ID: %s", order.UserID)
		log.Printf("  Status: %s", order.Status)
		log.Printf("  Total: %.2f", order.Total)
		log.Printf("  Version: %d", order.Version)
		log.Printf("  Items:")
		for _, item := range order.Items {
			log.Printf("    - Product: %s, Quantity: %d, Price: %.2f",
				item.ProductID, item.Quantity, item.Price)
		}
	}

	// 5. نمایش projection
	log.Println("\n5. Order projection:")
	if proj, ok := projection.GetOrder(orderID); ok {
		log.Printf("  Order: %+v", proj)
	}

	// 6. نمایش تاریخچه رویدادها
	log.Println("\n6. Event History:")
	events, _ := eventStore.Load(ctx, orderID)
	for i, event := range events {
		log.Printf("  %d. %s (version %d)", i+1, event.GetType(), event.GetVersion())
	}
}

// ============================================================================
// بخش 11: Best Practices
// ============================================================================

func bestPractices() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("💡 EVENT SOURCING BEST PRACTICES")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ 1. EVENT DESIGN                                                           │
├─────────────────────────────────────────────────────────────────────────────┤
│    • Use past tense for event names (OrderCreated, PaymentReceived)       │
│    • Include all relevant data in events                                   │
│    • Version your events for schema evolution                              │
│    • Keep events small and focused                                        │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 2. EVENT STORE                                                            │
├─────────────────────────────────────────────────────────────────────────────┤
│    • Append-only storage                                                   │
│    • Never modify or delete events                                        │
│    • Use optimistic concurrency control                                    │
│    • Consider using dedicated event store (EventStoreDB, PostgreSQL)      │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 3. SNAPSHOTS                                                              │
├─────────────────────────────────────────────────────────────────────────────┤
│    • Take snapshots periodically (e.g., every 100 events)                 │
│    • Store snapshots to speed up recovery                                 │
│    • Version snapshots with events                                        │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 4. PROJECTIONS                                                           │
├─────────────────────────────────────────────────────────────────────────────┤
│    • Use separate read models for queries                                 │
│    • Rebuild projections when schema changes                             │
│    • Handle events idempotently                                           │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 5. CONCURRENCY                                                            │
├─────────────────────────────────────────────────────────────────────────────┤
│    • Use optimistic concurrency with version checks                       │
│    • Handle concurrency conflicts gracefully                              │
│    • Consider using outbox pattern for event publishing                  │
└─────────────────────────────────────────────────────────────────────────────┘
`)
}

// ============================================================================
// بخش 12: Main
// ============================================================================

func main() {
	fmt.Println(stringsRepeat("=", 80))
	fmt.Println("🎯 EVENT SOURCING IN GO")
	fmt.Println("Event-Driven Architecture & State Reconstruction")
	fmt.Println(stringsRepeat("=", 80))

	bestPractices()

	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🚀 Running Event Sourcing Example")
	fmt.Println(stringsRepeat("=", 80))

	runEventSourcingExample()

	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🎯 EVENT SOURCING - COMPLETE")
	fmt.Println("Build event-driven systems with Go!")
	fmt.Println(stringsRepeat("=", 80))
}

// تابع کمکی
func stringsRepeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}
