// ============================================================================
// FILE: postgresql_guide.go
// TITLE: راهنمای کامل PostgreSQL در Go - Migration, Connection, Pooling
// HOW TO RUN: go run postgresql_guide.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - PostgreSQL در Go
// ============================================================================
//
// PostgreSQL یک دیتابیس رابطه‌ای قدرتمند و متن‌باز است.
// در Go برای کار با PostgreSQL از درایورهای مختلفی استفاده می‌شود:
//
// 1. database/sql + pq (lib/pq)
//    - قدیمی‌ترین و پایدارترین
//    - پشتیبانی از ویژگی‌های پایه PostgreSQL
//    - دیگر maintenance فعالی ندارد
//
// 2. database/sql + pgx (jackc/pgx)
//    - مدرن‌تر و سریع‌تر
//    - پشتیبانی از ویژگی‌های پیشرفته PostgreSQL
//    -推荐的 (توصیه شده)
//
// 3. sqlx (jmoiron/sqlx)
//    - افزونه بر روی database/sql
//    - امکانات راحت‌تر برای mapping struct به جدول
//    - می‌تواند با pq یا pgx کار کند
//
// 4. GORM (ORM)
//    - ORM کامل برای Go
//    - مناسب برای پروژه‌های ساده
//    - برای پروژه‌های پیچیده، sqlx یا pgx توصیه می‌شود
//
// قانون طلایی:
// "برای پروژه‌های جدید از pgx استفاده کن (سریع‌ترین و کامل‌ترین).
//  برای migration از golang-migrate استفاده کن.
//  همیشه connection pooling را فعال کن و max connections را محدود کن."
// ============================================================================

package __postgresql

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // درایور pgx برای database/sql
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

// ============================================================================
// بخش 1: نصب وابستگی‌ها
// ============================================================================

/*
نصب:

# golang-migrate (CLI)
$ go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# pgx (توصیه شده)
$ go get github.com/jackc/pgx/v5
$ go get github.com/jackc/pgx/v5/stdlib

# sqlx (اختیاری - برای راحتی)
$ go get github.com/jmoiron/sqlx

# pq (قدیمی - در صورت نیاز)
$ go get github.com/lib/pq

# godotenv (برای متغیرهای محیطی)
$ go get github.com/joho/godotenv
*/

// ============================================================================
// بخش 2: مدل‌های داده (Models)
// ============================================================================

// User مدل کاربر
type User struct {
	ID        int       `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Email     string    `db:"email" json:"email"`
	Age       int       `db:"age" json:"age"`
	IsActive  bool      `db:"is_active" json:"is_active"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// Product مدل محصول
type Product struct {
	ID         int     `db:"id" json:"id"`
	Name       string  `db:"name" json:"name"`
	Price      float64 `db:"price" json:"price"`
	Quantity   int     `db:"quantity" json:"quantity"`
	CategoryID int     `db:"category_id" json:"category_id"`
}

// Category مدل دسته‌بندی
type Category struct {
	ID   int    `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

// Order مدل سفارش
type Order struct {
	ID        int       `db:"id" json:"id"`
	UserID    int       `db:"user_id" json:"user_id"`
	Total     float64   `db:"total" json:"total"`
	Status    string    `db:"status" json:"status"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// OrderItem مدل آیتم سفارش
type OrderItem struct {
	ID        int     `db:"id" json:"id"`
	OrderID   int     `db:"order_id" json:"order_id"`
	ProductID int     `db:"product_id" json:"product_id"`
	Quantity  int     `db:"quantity" json:"quantity"`
	Price     float64 `db:"price" json:"price"`
}

// ============================================================================
// بخش 3: Migration با golang-migrate
// ============================================================================

/*
ساختار پوشه migration:

migrations/
├── 000001_create_users_table.up.sql
├── 000001_create_users_table.down.sql
├── 000002_create_products_table.up.sql
├── 000002_create_products_table.down.sql
├── 000003_create_orders_table.up.sql
└── 000003_create_orders_table.down.sql
*/

// فایل‌های SQL مثال:

// ============================================================================
// 000001_create_users_table.up.sql
// ============================================================================
const createUsersTableUp = `
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    age INT CHECK (age >= 0 AND age <= 150),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_name ON users(name);
`

// 000001_create_users_table.down.sql
const createUsersTableDown = `
DROP TABLE IF EXISTS users;
`

// ============================================================================
// 000002_create_products_table.up.sql
// ============================================================================
const createProductsTableUp = `
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    price DECIMAL(10,2) NOT NULL CHECK (price >= 0),
    quantity INT NOT NULL DEFAULT 0 CHECK (quantity >= 0),
    category_id INT REFERENCES categories(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_products_name ON products(name);
CREATE INDEX idx_products_category ON products(category_id);
`

// 000002_create_products_table.down.sql
const createProductsTableDown = `
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS categories;
`

// ============================================================================
// 000003_create_orders_table.up.sql
// ============================================================================
const createOrdersTableUp = `
CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    total DECIMAL(10,2) NOT NULL DEFAULT 0,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS order_items (
    id SERIAL PRIMARY KEY,
    order_id INT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id INT NOT NULL REFERENCES products(id),
    quantity INT NOT NULL CHECK (quantity > 0),
    price DECIMAL(10,2) NOT NULL CHECK (price >= 0)
);

CREATE INDEX idx_orders_user ON orders(user_id);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_order_items_order ON order_items(order_id);
`

// 000003_create_orders_table.down.sql
const createOrdersTableDown = `
DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS orders;
`

// ============================================================================
// بخش 4: Connection Pooling با database/sql + pgx
// ============================================================================

// Config تنظیمات دیتابیس
type DBConfig struct {
	Host         string
	Port         string
	User         string
	Password     string
	DBName       string
	SSLMode      string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  time.Duration
	MaxIdleTime  time.Duration
}

// LoadConfig بارگذاری تنظیمات از محیط
func LoadConfig() *DBConfig {
	// بارگذاری فایل .env
	godotenv.Load()

	return &DBConfig{
		Host:         getEnv("DB_HOST", "localhost"),
		Port:         getEnv("DB_PORT", "5432"),
		User:         getEnv("DB_USER", "postgres"),
		Password:     getEnv("DB_PASSWORD", "postgres"),
		DBName:       getEnv("DB_NAME", "testdb"),
		SSLMode:      getEnv("DB_SSLMODE", "disable"),
		MaxOpenConns: getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
		MaxIdleConns: getEnvAsInt("DB_MAX_IDLE_CONNS", 10),
		MaxLifetime:  time.Duration(getEnvAsInt("DB_MAX_LIFETIME", 60)) * time.Minute,
		MaxIdleTime:  time.Duration(getEnvAsInt("DB_MAX_IDLE_TIME", 10)) * time.Minute,
	}
}

// getEnv گرفتن متغیر محیطی با مقدار پیش‌فرض
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt گرفتن متغیر محیطی به صورت int
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var intVal int
		fmt.Sscanf(value, "%d", &intVal)
		return intVal
	}
	return defaultValue
}

// ConnectionString ایجاد رشته اتصال
func (c *DBConfig) ConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

// ConnectionStringPGX رشته اتصال مخصوص pgx
func (c *DBConfig) ConnectionStringPGX() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.DBName, c.SSLMode,
	)
}

// NewDatabase ایجاد اتصال دیتابیس با database/sql و pgx
func NewDatabase(config *DBConfig) (*sql.DB, error) {
	// باز کردن اتصال
	db, err := sql.Open("pgx", config.ConnectionStringPGX())
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// تنظیمات Connection Pooling
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.MaxLifetime)
	db.SetConnMaxIdleTime(config.MaxIdleTime)

	// تست اتصال
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	log.Println("Database connected successfully")
	return db, nil
}

// ============================================================================
// بخش 5: اتصال با SQLx (راحت‌تر از database/sql)
// ============================================================================

// NewSQLXDatabase ایجاد اتصال با sqlx
func NewSQLXDatabase(config *DBConfig) (*sqlx.DB, error) {
	// باز کردن اتصال
	db, err := sqlx.Connect("pgx", config.ConnectionStringPGX())
	if err != nil {
		return nil, fmt.Errorf("error connecting database: %w", err)
	}

	// تنظیمات Connection Pooling
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.MaxLifetime)
	db.SetConnMaxIdleTime(config.MaxIdleTime)

	log.Println("SQLx database connected successfully")
	return db, nil
}

// ============================================================================
// بخش 6: CRUD Operations با database/sql
// ============================================================================

// UserRepository سرویس کاربران
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository ایجاد repository جدید
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create ایجاد کاربر جدید
func (r *UserRepository) Create(user *User) error {
	query := `
		INSERT INTO users (name, email, age, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	err := r.db.QueryRow(
		query,
		user.Name,
		user.Email,
		user.Age,
		user.IsActive,
		time.Now(),
		time.Now(),
	).Scan(&user.ID)

	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}
	return nil
}

// GetByID دریافت کاربر با ID
func (r *UserRepository) GetByID(id int) (*User, error) {
	query := `SELECT id, name, email, age, is_active, created_at, updated_at FROM users WHERE id = $1`
	var user User
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Age,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("error getting user: %w", err)
	}
	return &user, nil
}

// GetAll دریافت همه کاربران
func (r *UserRepository) GetAll() ([]User, error) {
	query := `SELECT id, name, email, age, is_active, created_at, updated_at FROM users ORDER BY id`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying users: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Age,
			&user.IsActive,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning user: %w", err)
		}
		users = append(users, user)
	}
	return users, nil
}

// Update به‌روزرسانی کاربر
func (r *UserRepository) Update(user *User) error {
	query := `
		UPDATE users 
		SET name = $1, email = $2, age = $3, is_active = $4, updated_at = $5
		WHERE id = $6
	`
	result, err := r.db.Exec(
		query,
		user.Name,
		user.Email,
		user.Age,
		user.IsActive,
		time.Now(),
		user.ID,
	)
	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

// Delete حذف کاربر
func (r *UserRepository) Delete(id int) error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

// ============================================================================
// بخش 7: CRUD Operations با SQLx (راحت‌تر)
// ============================================================================

// SQLxUserRepository سرویس کاربران با sqlx
type SQLxUserRepository struct {
	db *sqlx.DB
}

// NewSQLxUserRepository ایجاد repository جدید
func NewSQLxUserRepository(db *sqlx.DB) *SQLxUserRepository {
	return &SQLxUserRepository{db: db}
}

// Create ایجاد کاربر
func (r *SQLxUserRepository) Create(user *User) error {
	query := `
		INSERT INTO users (name, email, age, is_active, created_at, updated_at)
		VALUES (:name, :email, :age, :is_active, :created_at, :updated_at)
		RETURNING id
	`
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	rows, err := r.db.NamedQuery(query, user)
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		rows.Scan(&user.ID)
	}
	return nil
}

// GetByID دریافت کاربر با ID
func (r *SQLxUserRepository) GetByID(id int) (*User, error) {
	query := `SELECT * FROM users WHERE id = $1`
	var user User
	err := r.db.Get(&user, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("error getting user: %w", err)
	}
	return &user, nil
}

// GetAll دریافت همه کاربران
func (r *SQLxUserRepository) GetAll() ([]User, error) {
	query := `SELECT * FROM users ORDER BY id`
	var users []User
	err := r.db.Select(&users, query)
	if err != nil {
		return nil, fmt.Errorf("error getting users: %w", err)
	}
	return users, nil
}

// Update به‌روزرسانی کاربر
func (r *SQLxUserRepository) Update(user *User) error {
	user.UpdatedAt = time.Now()
	query := `UPDATE users SET name = :name, email = :email, age = :age, is_active = :is_active, updated_at = :updated_at WHERE id = :id`
	result, err := r.db.NamedExec(query, user)
	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

// Delete حذف کاربر
func (r *SQLxUserRepository) Delete(id int) error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

// ============================================================================
// بخش 8: Transactions (تراکنش‌ها)
// ============================================================================

// OrderRepository سرویس سفارشات با تراکنش
type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

// CreateOrderWithItems ایجاد سفارش با آیتم‌ها (با تراکنش)
func (r *OrderRepository) CreateOrderWithItems(userID int, items []OrderItem) (*Order, error) {
	// شروع تراکنش
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback() // اگر commit نشود، rollback می‌شود

	// محاسبه total
	var total float64
	for _, item := range items {
		// گرفتن قیمت محصول
		var price float64
		err := tx.QueryRow(`SELECT price FROM products WHERE id = $1`, item.ProductID).Scan(&price)
		if err != nil {
			return nil, fmt.Errorf("error getting product price: %w", err)
		}
		item.Price = price
		total += price * float64(item.Quantity)
	}

	// ایجاد سفارش
	var order Order
	order.UserID = userID
	order.Total = total
	order.Status = "pending"
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

	err = tx.QueryRow(`
		INSERT INTO orders (user_id, total, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, order.UserID, order.Total, order.Status, order.CreatedAt, order.UpdatedAt).Scan(&order.ID)
	if err != nil {
		return nil, fmt.Errorf("error creating order: %w", err)
	}

	// ایجاد آیتم‌های سفارش
	for _, item := range items {
		item.OrderID = order.ID
		_, err = tx.Exec(`
			INSERT INTO order_items (order_id, product_id, quantity, price)
			VALUES ($1, $2, $3, $4)
		`, item.OrderID, item.ProductID, item.Quantity, item.Price)
		if err != nil {
			return nil, fmt.Errorf("error creating order item: %w", err)
		}

		// کاهش موجودی محصول
		_, err = tx.Exec(`
			UPDATE products SET quantity = quantity - $1 WHERE id = $2 AND quantity >= $1
		`, item.Quantity, item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("error updating product quantity: %w", err)
		}
	}

	// commit تراکنش
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	return &order, nil
}

// ============================================================================
// بخش 9: Prepared Statements و Bulk Operations
// ============================================================================

// BulkInsertUsers درج دسته‌ای کاربران
func (r *UserRepository) BulkInsertUsers(users []User) error {
	// آماده‌سازی statement
	stmt, err := r.db.Prepare(`
		INSERT INTO users (name, email, age, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	// اجرای bulk insert
	for _, user := range users {
		_, err := stmt.Exec(user.Name, user.Email, user.Age, user.IsActive, time.Now(), time.Now())
		if err != nil {
			return fmt.Errorf("error inserting user: %w", err)
		}
	}
	return nil
}

// BulkInsertUsersSQLx درج دسته‌ای با sqlx
func (r *SQLxUserRepository) BulkInsertUsersSQLx(users []User) error {
	query := `
		INSERT INTO users (name, email, age, is_active, created_at, updated_at)
		VALUES (:name, :email, :age, :is_active, :created_at, :updated_at)
	`
	_, err := r.db.NamedExec(query, users)
	if err != nil {
		return fmt.Errorf("error bulk inserting users: %w", err)
	}
	return nil
}

// ============================================================================
// بخش 10: JSON Operations با PostgreSQL
// ============================================================================

// ProductWithMetadata محصول با متادیتا (JSON)
type ProductWithMetadata struct {
	Product
	Metadata map[string]interface{} `db:"metadata" json:"metadata"`
}

// CreateProductWithMetadata ایجاد محصول با متادیتا JSON
func CreateProductWithMetadata(db *sql.DB, product ProductWithMetadata) error {
	metadataJSON, err := json.Marshal(product.Metadata)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO products (name, price, quantity, category_id, metadata)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err = db.Exec(query, product.Name, product.Price, product.Quantity, product.CategoryID, metadataJSON)
	return err
}

// ============================================================================
// بخش 11: Context و Timeout
// ============================================================================

// GetUserWithContext دریافت کاربر با context (برای timeout و cancellation)
func (r *UserRepository) GetUserWithContext(ctx context.Context, id int) (*User, error) {
	query := `SELECT id, name, email, age, is_active, created_at, updated_at FROM users WHERE id = $1`
	var user User

	// اجرا با context
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Age,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("error getting user: %w", err)
	}
	return &user, nil
}

// ============================================================================
// بخش 12: Health Check و Monitoring
// ============================================================================

// HealthCheck بررسی سلامت دیتابیس
func HealthCheck(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	return db.PingContext(ctx)
}

// GetDBStats دریافت آمار connection pool
func GetDBStats(db *sql.DB) map[string]interface{} {
	stats := db.Stats()
	return map[string]interface{}{
		"max_open_connections": stats.MaxOpenConnections,
		"open_connections":     stats.OpenConnections,
		"in_use":               stats.InUse,
		"idle":                 stats.Idle,
		"wait_count":           stats.WaitCount,
		"wait_duration":        stats.WaitDuration.String(),
		"max_idle_closed":      stats.MaxIdleClosed,
		"max_lifetime_closed":  stats.MaxLifetimeClosed,
	}
}

// ============================================================================
// بخش 13: Migrate CLI Commands
// ============================================================================

func migrateCommands() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📋 MIGRATE CLI COMMANDS")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println(`
# ایجاد migration جدید
$ migrate create -ext sql -dir migrations -seq create_users_table

# اجرای migration‌ها (up)
$ migrate -path migrations -database "postgres://user:pass@localhost:5432/dbname?sslmode=disable" up

# اجرای تعداد مشخصی migration
$ migrate -path migrations -database "postgres://user:pass@localhost:5432/dbname?sslmode=disable" up 2

# rollback (down)
$ migrate -path migrations -database "postgres://user:pass@localhost:5432/dbname?sslmode=disable" down

# rollback همه migration‌ها
$ migrate -path migrations -database "postgres://user:pass@localhost:5432/dbname?sslmode=disable" down -all

# نمایش نسخه فعلی
$ migrate -path migrations -database "postgres://user:pass@localhost:5432/dbname?sslmode=disable" version

#强制执行 migration در صورت خطا
$ migrate -path migrations -database "postgres://user:pass@localhost:5432/dbname?sslmode=disable" force 20240101000000
`)
}

// ============================================================================
// بخش 14: Best Practices
// ============================================================================

func bestPractices() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("💡 POSTGRESQL BEST PRACTICES")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────┐
│ CONNECTION POOLING                                            │
├─────────────────────────────────────────────────────────────────┤
│ • SetMaxOpenConns: محدود بر اساس توان دیتابیس (معمولاً 10-50) │
│ • SetMaxIdleConns: معمولاً برابر یا کمتر از MaxOpenConns      │
│ • SetConnMaxLifetime: حداکثر 1 ساعت                           │
│ • SetConnMaxIdleTime: حداکثر 10-15 دقیقه                      │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ QUERY OPTIMIZATION                                            │
├─────────────────────────────────────────────────────────────────┤
│ • از prepared statements برای queryهای تکراری استفاده کن      │
│ • برای bulk insert از COPY استفاده کن                         │
│ • از LIMIT و OFFSET برای pagination استفاده کن               │
│ • ایندکس مناسب برای فیلدهای جستجو ایجاد کن                   │
│ • از EXPLAIN ANALYZE برای بررسی performance queryها استفاده کن│
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ TRANSACTIONS                                                  │
├─────────────────────────────────────────────────────────────────┤
│ • تراکنش‌ها را کوتاه نگه دار                                  │
│ • همیشه از defer tx.Rollback() استفاده کن                    │
│ • از isolation levels مناسب استفاده کن                        │
│ • از deadlock آگاه باش                                         │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ MIGRATION                                                     │
├─────────────────────────────────────────────────────────────────┤
│ • همیشه up و down migration بنویس                             │
│ • migrationها را version control کن                           │
│ • از sequential numbering استفاده کن                         │
│ • هر migration یک change خاص داشته باشد                       │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ ERROR HANDLING                                                │
├─────────────────────────────────────────────────────────────────┤
│ • همیشه sql.ErrNoRows را بررسی کن                            │
│ • خطاهای connection را handle کن                              │
│ • از context cancellation آگاه باش                            │
│ • لاگ خطاهای دیتابیس برای debugging                           │
└─────────────────────────────────────────────────────────────────┘
`)
}

// ============================================================================
// بخش 15: Connection String Examples
// ============================================================================

func connectionStrings() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🔗 CONNECTION STRING EXAMPLES")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────┐
│ FORMATS                                                        │
├─────────────────────────────────────────────────────────────────┤
│                                                                │
│ 1. Standard (pgx/pq)                                         │
│    postgres://username:password@localhost:5432/dbname?sslmode=disable│
│                                                                │
│ 2. Key-Value (pq)                                            │
│    host=localhost port=5432 user=postgres password=postgres dbname=testdb sslmode=disable│
│                                                                │
│ 3. With SSL                                                    │
│    postgres://user:pass@localhost:5432/dbname?sslmode=require  │
│                                                                │
│ 4. With Connection Pool Parameters                             │
│    postgres://user:pass@localhost:5432/dbname?pool_max_conns=10│
│                                                                │
│ 5. Unix Socket                                                 │
│    postgres:///postgres?host=/var/run/postgresql              │
│                                                                │
└─────────────────────────────────────────────────────────────────┘

🎯 COMPARISON: pq vs pgx

   FEATURE           | pq         | pgx
   ------------------|------------|------------------
   Performance       | Good       | Excellent (2-3x faster)
   PostgreSQL 10+    | Partial    | Full
   Logical Replication| No        | Yes
   COPY protocol     | Basic      | Advanced
   SCRAM auth        | No         | Yes
   Connection pool   | Limited    | Built-in
   Maintenance       | Low        | Active
   
   RECOMMENDATION: Use pgx for new projects ✅
`)
}

// ============================================================================
// بخش 16: Main
// ============================================================================

func main() {
	fmt.Println(stringsRepeat("=", 80))
	fmt.Println("🎯 POSTGRESQL GUIDE")
	fmt.Println("Migration | database/sql | sqlx | Connection Pooling | pgx")
	fmt.Println(stringsRepeat("=", 80))

	migrateCommands()
	bestPractices()
	connectionStrings()

	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📝 EXAMPLE USAGE")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println(`
// مثال اتصال با pgx
func main() {
    config := LoadConfig()
    db, err := NewDatabase(config)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // CRUD operations
    userRepo := NewUserRepository(db)
    
    // Create user
    user := &User{
        Name:     "Ali Rezaei",
        Email:    "ali@example.com",
        Age:      30,
        IsActive: true,
    }
    if err := userRepo.Create(user); err != nil {
        log.Fatal(err)
    }
    
    // Get user
    user, err := userRepo.GetByID(1)
    if err != nil {
        log.Fatal(err)
    }
    
    // Update user
    user.Name = "Ali Mohammadi"
    if err := userRepo.Update(user); err != nil {
        log.Fatal(err)
    }
    
    // Delete user
    if err := userRepo.Delete(1); err != nil {
        log.Fatal(err)
    }
}
`)

	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🎯 POSTGRESQL GUIDE - COMPLETE")
	fmt.Println("Ready to build robust database applications!")
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

/*
خلاصه درایورها و کتابخانه‌ها
کتابخانه	کاربرد	مزایا	معایب
database/sql	استاندارد Go	پایدار، استاندارد	کد تکراری زیاد
pgx	درایور مدرن	سریع، کامل، active	وابستگی خارجی
sqlx	افزونه بر database/sql	راحت‌تر، mapping خودکار	لایه اضافی
pq	درایور قدیمی	پایدار، ساده	دیگر active نیست
golang-migrate	مدیریت migration	کامل، CLI و کتابخانه	نیاز به یادگیری

*/
/*
تنظیمات Connection Pool
متد	توضیح	مقدار پیشنهادی
SetMaxOpenConns	حداکثر اتصالات همزمان	10-50
SetMaxIdleConns	حداکثر اتصالات idle	5-10
SetConnMaxLifetime	حداکثر عمر هر ا
*/
