// ============================================================================
// FILE: mongodb_complete_guide.go
// TITLE: راهنمای کامل MongoDB در Go - Official MongoDB Go Driver v2
// HOW TO RUN: go run mongodb_complete_guide.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - MongoDB Go Driver v2
// ============================================================================
//
// MongoDB یک دیتابیس NoSQL مستندگرا (Document-oriented) است.
// از نسخه 1.17.8 به بعد، درایور قدیمی (v1) رسماً منسوخ (deprecated) اعلام شده است.
// کلیه پروژه‌های جدید باید از نسخه 2 به بعد استفاده کنند.
//
// نصب:
// $ go get go.mongodb.org/mongo-driver/v2/mongo
// $ go get go.mongodb.org/mongo-driver/v2/bson
// $ go get go.mongodb.org/mongo-driver/v2/mongo/options
//
// تغییرات مهم در v2 نسبت به v1:
// 1. مسیر import تغییر کرده: go.mongodb.org/mongo-driver/v2/...
// 2. بهبود performance با استفاده از sync.Pool در BSON encoding/decoding [citation:2]
// 3. پشتیبانی از omitempty سراسری بدون نیاز به تگ در هر فیلد [citation:2]
// 4. پشتیبانی از errors.Is و errors.As برای مدیریت خطاها [citation:2]
// 5. MarshalBSONValue دیگر پشتیبانی نمی‌شود - به جای آن از MarshalBSON استفاده کنید [citation:4]
//
// قانون طلایی:
// "برای پروژه‌های جدید همیشه از نسخه v2 استفاده کن.
//  از bson.D برای documentهای ordered (مثل دستورات MongoDB) و از bson.M برای داده‌های معمولی استفاده کن.
//  همیشه context را به تمام عملیات پاس بده و timeout تنظیم کن."
// ============================================================================

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

// ============================================================================
// بخش 1: مدل‌های داده (Models)
// ============================================================================

// User مدل کاربر (با تگ‌های bson)
type User struct {
	ID        string    `bson:"_id,omitempty" json:"id"`           // omitempty برای ID جدید
	Name      string    `bson:"name" json:"name" validate:"required"`
	Email     string    `bson:"email" json:"email" validate:"required,email"`
	Age       int       `bson:"age" json:"age"`
	IsActive  bool      `bson:"is_active" json:"is_active"`
	Role      string    `bson:"role" json:"role"`
	Tags      []string  `bson:"tags,omitempty" json:"tags,omitempty"`
	Metadata  bson.M    `bson:"metadata,omitempty" json:"metadata,omitempty"` // داده داینامیک
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

// Product مدل محصول
type Product struct {
	ID          string    `bson:"_id,omitempty" json:"id"`
	Name        string    `bson:"name" json:"name"`
	Price       float64   `bson:"price" json:"price"`
	Quantity    int       `bson:"quantity" json:"quantity"`
	Category    string    `bson:"category" json:"category"`
	Description string    `bson:"description,omitempty" json:"description,omitempty"`
	Tags        []string  `bson:"tags,omitempty" json:"tags,omitempty"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
}

// Order مدل سفارش
type Order struct {
	ID         string      `bson:"_id,omitempty" json:"id"`
	UserID     string      `bson:"user_id" json:"user_id"`
	Items      []OrderItem `bson:"items" json:"items"`
	Total      float64     `bson:"total" json:"total"`
	Status     string      `bson:"status" json:"status"`
	CreatedAt  time.Time   `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time   `bson:"updated_at" json:"updated_at"`
}

// OrderItem آیتم سفارش
type OrderItem struct {
	ProductID string  `bson:"product_id" json:"product_id"`
	Name      string  `bson:"name" json:"name"`
	Quantity  int     `bson:"quantity" json:"quantity"`
	Price     float64 `bson:"price" json:"price"`
}

// ============================================================================
// بخش 2: اتصال به MongoDB (Connection)
// ============================================================================

// Config تنظیمات اتصال
type MongoConfig struct {
	URI      string
	Database string
	Timeout  time.Duration
}

// DefaultConfig تنظیمات پیش‌فرض
func DefaultConfig() *MongoConfig {
	return &MongoConfig{
		URI:      "mongodb://localhost:27017",
		Database: "testdb",
		Timeout:  10 * time.Second,
	}
}

// NewMongoClient ایجاد کلاینت MongoDB (نسخه v2)
func NewMongoClient(config *MongoConfig) (*mongo.Client, error) {
	// تنظیمات کلاینت
	clientOptions := options.Client().
		ApplyURI(config.URI).
		SetTimeout(config.Timeout).
		SetRetryWrites(true).
		SetRetryReads(true).
		SetMaxPoolSize(100).
		SetMinPoolSize(10)

	// اتصال به MongoDB
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// تست اتصال با Ping
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Println("Connected to MongoDB successfully")
	return client, nil
}

// ============================================================================
// بخش 3: عملیات CRUD پایه (با استفاده از v2)
// ============================================================================

// UserRepository سرویس کاربران
type UserRepository struct {
	collection *mongo.Collection
}

// NewUserRepository ایجاد repository جدید
func NewUserRepository(client *mongo.Client, dbName string) *UserRepository {
	collection := client.Database(dbName).Collection("users")
	return &UserRepository{collection: collection}
}

// Create ایجاد کاربر جدید
func (r *UserRepository) Create(ctx context.Context, user *User) error {
	user.ID = bson.NewObjectID().Hex()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	user.ID = result.InsertedID.(string)
	return nil
}

// GetByID دریافت کاربر با ID
func (r *UserRepository) GetByID(ctx context.Context, id string) (*User, error) {
	var user User
	filter := bson.M{"_id": id}

	err := r.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// GetAll دریافت همه کاربران (با pagination)
func (r *UserRepository) GetAll(ctx context.Context, page, pageSize int) ([]User, int64, error) {
	var users []User

	// محاسبه offset
	skip := int64((page - 1) * pageSize)
	limit := int64(pageSize)

	// تنظیمات find
	opts := options.Find().
		SetSkip(skip).
		SetLimit(limit).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find users: %w", err)
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &users); err != nil {
		return nil, 0, fmt.Errorf("failed to decode users: %w", err)
	}

	// شمارش کل
	total, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	return users, total, nil
}

// Update به‌روزرسانی کاربر
func (r *UserRepository) Update(ctx context.Context, user *User) error {
	user.UpdatedAt = time.Now()

	filter := bson.M{"_id": user.ID}
	update := bson.M{
		"$set": bson.M{
			"name":       user.Name,
			"email":      user.Email,
			"age":        user.Age,
			"is_active":  user.IsActive,
			"role":       user.Role,
			"tags":       user.Tags,
			"metadata":   user.Metadata,
			"updated_at": user.UpdatedAt,
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// Delete حذف کاربر
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	filter := bson.M{"_id": id}
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// ============================================================================
// بخش 4: Queryهای پیشرفته (فیلترها، پروجکشن، سورت، پیجینیشن)
// ============================================================================

// FindByFilters جستجوی پیشرفته با فیلترهای مختلف
func (r *UserRepository) FindByFilters(ctx context.Context, filters map[string]interface{}, sortField string, sortOrder int, page, pageSize int) ([]User, int64, error) {
	var users []User

	// ساخت فیلتر
	filter := bson.M{}
	for k, v := range filters {
		filter[k] = v
	}

	// تنظیمات sort
	sort := bson.D{{Key: sortField, Value: sortOrder}}

	// تنظیمات pagination
	skip := int64((page - 1) * pageSize)
	limit := int64(pageSize)

	opts := options.Find().
		SetSort(sort).
		SetSkip(skip).
		SetLimit(limit)

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find users: %w", err)
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &users); err != nil {
		return nil, 0, fmt.Errorf("failed to decode users: %w", err)
	}

	// شمارش کل
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	return users, total, nil
}

// FindByAgeRange جستجوی کاربران در بازه سنی
func (r *UserRepository) FindByAgeRange(ctx context.Context, minAge, maxAge int) ([]User, error) {
	var users []User

	filter := bson.M{
		"age": bson.M{
			"$gte": minAge,
			"$lte": maxAge,
		},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find users: %w", err)
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &users); err != nil {
		return nil, fmt.Errorf("failed to decode users: %w", err)
	}

	return users, nil
}

// SearchUsers جستجوی متن در نام و ایمیل
func (r *UserRepository) SearchUsers(ctx context.Context, keyword string) ([]User, error) {
	var users []User

	filter := bson.M{
		"$or": []bson.M{
			{"name": bson.M{"$regex": keyword, "$options": "i"}},
			{"email": bson.M{"$regex": keyword, "$options": "i"}},
		},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &users); err != nil {
		return nil, fmt.Errorf("failed to decode users: %w", err)
	}

	return users, nil
}

// ============================================================================
// بخش 5: Aggregation Pipeline (برای گزارشات و تحلیل داده)
// ============================================================================

// UserStats آمار کاربران
type UserStats struct {
	TotalUsers    int64   `bson:"total_users" json:"total_users"`
	ActiveUsers   int64   `bson:"active_users" json:"active_users"`
	InactiveUsers int64   `bson:"inactive_users" json:"inactive_users"`
	AverageAge    float64 `bson:"average_age" json:"average_age"`
	MinAge        int     `bson:"min_age" json:"min_age"`
	MaxAge        int     `bson:"max_age" json:"max_age"`
}

// GetUserStats دریافت آمار کاربران با aggregation pipeline
func (r *UserRepository) GetUserStats(ctx context.Context) (*UserStats, error) {
	var stats UserStats

	// Pipeline برای aggregation
	pipeline := mongo.Pipeline{
		// مرحله 1: گروه‌بندی برای محاسبات
		{{Key: "$group", Value: bson.M{
			"_id":          nil,
			"total_users":  bson.M{"$sum": 1},
			"active_users": bson.M{"$sum": bson.M{"$cond": []interface{}{"$is_active", 1, 0}}},
			"inactive_users": bson.M{"$sum": bson.M{"$cond": []interface{}{bson.M{"$eq": []interface{}{"$is_active", false}}, 1, 0}}},
			"average_age":  bson.M{"$avg": "$age"},
			"min_age":      bson.M{"$min": "$age"},
			"max_age":      bson.M{"$max": "$age"},
		}},
		}}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to aggregate: %w", err)
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		if err := cursor.Decode(&stats); err != nil {
			return nil, fmt.Errorf("failed to decode stats: %w", err)
		}
	}

	return &stats, nil
}

// AgeGroupCount تعداد کاربران در گروه‌های سنی مختلف
type AgeGroupCount struct {
	Group string `bson:"_id" json:"group"`
	Count int64  `bson:"count" json:"count"`
}

// GetAgeGroupCount دریافت توزیع سنی کاربران
func (r *UserRepository) GetAgeGroupCount(ctx context.Context) ([]AgeGroupCount, error) {
	var results []AgeGroupCount

	pipeline := mongo.Pipeline{
		// مرحله 1: اضافه کردن فیلد age_group
		{{Key: "$addFields", Value: bson.M{
			"age_group": bson.M{
				"$switch": bson.M{
					"branches": []bson.M{
						{"case": bson.M{"$lt": []interface{}{"$age", 18}}, "then": "Under 18"},
						{"case": bson.M{"$and": []interface{}{
							bson.M{"$gte": []interface{}{"$age", 18}},
							bson.M{"$lt": []interface{}{"$age", 30}},
						}}, "then": "18-29"},
						{"case": bson.M{"$and": []interface{}{
							bson.M{"$gte": []interface{}{"$age", 30}},
							bson.M{"$lt": []interface{}{"$age", 50}},
						}}, "then": "30-49"},
					},
					"default": "50+",
				},
			},
		}},
			// مرحله 2: گروه‌بندی
			{{Key: "$group", Value: bson.M{
				"_id":   "$age_group",
				"count": bson.M{"$sum": 1},
			}}},
			// مرحله 3: مرتب‌سازی
			{{Key: "$sort", Value: bson.M{"_id": 1}}},
		}

		cursor, err := r.collection.Aggregate(ctx, pipeline)
		if err != nil {
		return nil, fmt.Errorf("failed to aggregate: %w", err)
	}
		defer cursor.Close(ctx)

		if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to decode results: %w", err)
	}

		return results, nil
	}

	// ============================================================================
	// بخش 6: Bulk Operations (عملیات دسته‌ای)
	// ============================================================================

	// BulkInsert درج دسته‌ای کاربران
	func (r *UserRepository) BulkInsert(ctx context.Context, users []User) error {
		if len(users) == 0 {
		return nil
	}

		// تبدیل به interface{} برای InsertMany
		documents := make([]interface{}, len(users))
		for i, user := range users {
		user.ID = bson.NewObjectID().Hex()
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()
		documents[i] = user
	}

		_, err := r.collection.InsertMany(ctx, documents)
		if err != nil {
		return fmt.Errorf("failed to bulk insert: %w", err)
	}

		return nil
	}

	// BulkUpdate به‌روزرسانی دسته‌ای
	func (r *UserRepository) BulkUpdate(ctx context.Context, updates map[string]bson.M) error {
		// ایجاد مدل‌های bulk write
		var models []mongo.WriteModel

		for id, update := range updates {
		filter := bson.M{"_id": id}
		updateModel := mongo.NewUpdateOneModel().
		SetFilter(filter).
		SetUpdate(bson.M{"$set": update})
		models = append(models, updateModel)
	}

		if len(models) == 0 {
		return nil
	}

		_, err := r.collection.BulkWrite(ctx, models)
		if err != nil {
		return fmt.Errorf("failed to bulk update: %w", err)
	}

		return nil
	}

	// ============================================================================
	// بخش 7: Transactions (تراکنش‌ها - نیاز به replica set)
	// ============================================================================

	// CreateOrderWithTransaction ایجاد سفارش با تراکنش
	func CreateOrderWithTransaction(ctx context.Context, client *mongo.Client, dbName string, order *Order) error {
		// شروع تراکنش
		session, err := client.StartSession()
		if err != nil {
		return fmt.Errorf("failed to start session: %w", err)
	}
		defer session.EndSession(ctx)

		// اجرای عملیات در تراکنش
		callback := func(sessCtx context.Context) (interface{}, error) {
		// دریافت collection‌ها
		ordersCollection := client.Database(dbName).Collection("orders")
		productsCollection := client.Database(dbName).Collection("products")

		// ایجاد سفارش
		order.ID = bson.NewObjectID().Hex()
		order.CreatedAt = time.Now()
		order.UpdatedAt = time.Now()

		if _, err := ordersCollection.InsertOne(sessCtx, order); err != nil {
		return nil, err
	}

		// کاهش موجودی محصولات
		for _, item := range order.Items {
		filter := bson.M{"_id": item.ProductID}
		update := bson.M{"$inc": bson.M{"quantity": -item.Quantity}}

		result, err := productsCollection.UpdateOne(sessCtx, filter, update)
		if err != nil {
		return nil, err
	}
		if result.MatchedCount == 0 {
		return nil, fmt.Errorf("product not found: %s", item.ProductID)
	}
		if result.ModifiedCount == 0 {
		return nil, fmt.Errorf("insufficient stock for product: %s", item.ProductID)
	}
	}

		return nil, nil
	}

		_, err = session.WithTransaction(ctx, callback)
		if err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}

		return nil
	}

	// ============================================================================
	// بخش 8: Index Management (مدیریت ایندکس‌ها)
	// ============================================================================

	// CreateIndexes ایجاد ایندکس‌های مورد نیاز
	func CreateIndexes(ctx context.Context, client *mongo.Client, dbName string) error {
		usersCollection := client.Database(dbName).Collection("users")
		productsCollection := client.Database(dbName).Collection("products")

		// ایندکس برای email (unique)
		emailIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("idx_email_unique"),
	}

		// ایندکس کامپوزیت برای name و age
		nameAgeIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "name", Value: 1}, {Key: "age", Value: -1}},
		Options: options.Index().SetName("idx_name_age"),
	}

		// ایندکس برای جستجوی متنی
		textIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "name", Value: "text"}, {Key: "email", Value: "text"}},
		Options: options.Index().SetName("idx_text_search"),
	}

		// ایجاد ایندکس‌ها
		_, err := usersCollection.Indexes().CreateMany(ctx, []mongo.IndexModel{emailIndex, nameAgeIndex, textIndex})
		if err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

		// ایندکس برای products
		priceIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "price", Value: 1}},
		Options: options.Index().SetName("idx_price"),
	}

		_, err = productsCollection.Indexes().CreateOne(ctx, priceIndex)
		if err != nil {
		return fmt.Errorf("failed to create product index: %w", err)
	}

		log.Println("Indexes created successfully")
		return nil
	}

	// ============================================================================
	// بخش 9: BSON Types (انواع داده خاص MongoDB)
	// ============================================================================

	// BSONExamples مثال‌های استفاده از انواع BSON
	func BSONExamples() {
		// ObjectID (شناسه یکتا)
		objID := bson.NewObjectID()
		fmt.Printf("ObjectID: %s\n", objID.Hex())
		fmt.Printf("ObjectID Timestamp: %v\n", objID.Timestamp())

		// DateTime
		now := time.Now()
		datetime := bson.NewDateTimeFromTime(now)
		fmt.Printf("BSON DateTime: %d\n", datetime)
		fmt.Printf("Back to Time: %v\n", datetime.Time())

		// Decimal128 (برای اعداد دقیق مالی)
		price, _ := bson.ParseDecimal128("99.99")
		fmt.Printf("Decimal128: %s\n", price.String())

		// Binary
		binaryData := bson.Binary{
			Subtype: 0x00,
			Data:    []byte("hello world"),
		}
		fmt.Printf("Binary data length: %d\n", len(binaryData.Data))
	}

	// ============================================================================
	// بخش 10: Best Practices
	// ============================================================================

	func bestPractices() {
		fmt.Println("\n" + stringsRepeat("=", 80))
		fmt.Println("💡 MONGODB BEST PRACTICES (v2)")
		fmt.Println(stringsRepeat("=", 80))

		fmt.Println(`
┌─────────────────────────────────────────────────────────────────┐
│ CONNECTION MANAGEMENT                                         │
├─────────────────────────────────────────────────────────────────┤
│ • Use connection pooling (default is enabled)                 │
│ • Set MaxPoolSize based on your application load              │
│ • Always defer Disconnect()                                   │
│ • Use context with timeout for all operations                 │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ DOCUMENT DESIGN                                               │
├─────────────────────────────────────────────────────────────────┤
│ • Use meaningful _id values when possible                     │
│ • Keep documents small (avoid unbounded arrays)               │
│ • Use embedded documents for related data                     │
│ • Use references for many-to-many relationships               │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ QUERY OPTIMIZATION                                            │
├─────────────────────────────────────────────────────────────────┤
│ • Create indexes for frequently queried fields                │
│ • Use projection to return only needed fields                 │
│ • Use explain() to analyze query performance                  │
│ • Avoid using skip for large offsets (use cursor pagination)  │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ BSON TYPES                                                    │
├─────────────────────────────────────────────────────────────────┤
│ • Use bson.D for ordered documents (commands)                 │
│ • Use bson.M for unordered documents (queries)                │
│ • Use bson.A for arrays                                       │
│ • Use primitive.ObjectID for _id fields                       │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ ERROR HANDLING                                                │
├─────────────────────────────────────────────────────────────────┤
│ • Check for mongo.ErrNoDocuments when using FindOne          │
│ • Use errors.Is and errors.As for error handling [citation:2]│
│ • Always check cursor.Err() after iteration                  │
│ • Close cursors with defer                                   │
└─────────────────────────────────────────────────────────────────┘
`)
	}

	// ============================================================================
	// بخش 11: Main
	// ============================================================================

	func main() {
		fmt.Println(stringsRepeat("=", 80))
		fmt.Println("🎯 COMPLETE MONGODB GUIDE")
		fmt.Println("Official MongoDB Go Driver v2 - BSON | CRUD | Aggregation | Transactions")
		fmt.Println(stringsRepeat("=", 80))

		bestPractices()

		fmt.Println("\n" + stringsRepeat("=", 80))
		fmt.Println("📦 IMPORTANT: VERSION 1 IS DEPRECATED!")
		fmt.Println(stringsRepeat("=", 80))

		fmt.Println(`
╔═══════════════════════════════════════════════════════════════════════════════╗
║  ⚠️  MongoDB Go Driver v1 is DEPRECATED as of version 1.17.8 [citation:1][citation:6]  ║
║                                                                               ║
║  All new projects MUST use v2:                                               ║
║  $ go get go.mongodb.org/mongo-driver/v2/mongo                               ║
║                                                                               ║
║  Migration guide: https://www.mongodb.com/docs/drivers/go/upgrade-v2/        ║
╚═══════════════════════════════════════════════════════════════════════════════╝
`)

		fmt.Println("\n" + stringsRepeat("=", 80))
		fmt.Println("📝 EXAMPLE USAGE")
		fmt.Println(stringsRepeat("=", 80))

		fmt.Println(`
// مثال کامل استفاده از MongoDB Driver v2
func main() {
    ctx := context.Background()
    
    // اتصال به MongoDB
    config := DefaultConfig()
    client, err := NewMongoClient(config)
    if err != nil {
        log.Fatal(err)
    }
    defer client.Disconnect(ctx)
    
    // ایجاد repository
    userRepo := NewUserRepository(client, config.Database)
    
    // ایجاد کاربر جدید
    user := &User{
        Name:     "Ali Rezaei",
        Email:    "ali@example.com",
        Age:      30,
        IsActive: true,
        Role:     "user",
        Tags:     []string{"golang", "mongodb"},
    }
    
    if err := userRepo.Create(ctx, user); err != nil {
        log.Fatal(err)
    }
    fmt.Printf("User created: %s\n", user.ID)
    
    // دریافت کاربر
    fetched, err := userRepo.GetByID(ctx, user.ID)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Fetched user: %+v\n", fetched)
    
    // جستجو
    users, err := userRepo.SearchUsers(ctx, "ali")
    if err != nil {
        log.Fatal(err)
    }
    
    // آمار
    stats, err := userRepo.GetUserStats(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Stats: %+v\n", stats)
}
`)

		fmt.Println("\n" + stringsRepeat("=", 80))
		fmt.Println("🎯 MONGODB GUIDE - COMPLETE")
		fmt.Println("Ready to build scalable NoSQL applications with MongoDB!")
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