// ============================================================================
// FILE: query_optimization_index_guide.go
// TITLE: راهنمای کامل بهینه‌سازی کوئری و ایندکس در PostgreSQL
// HOW TO RUN: go run query_optimization_index_guide.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - چرا بهینه‌سازی کوئری و ایندکس اهمیت دارد؟
// ============================================================================
//
// بهینه‌سازی کوئری و ایندکس برای عملکرد دیتابیس حیاتی است:
//
// مشکلات رایج:
// 1. Full Table Scan: اسکن کامل جدول بدون استفاده از ایندکس
// 2. N+1 Query: کوئری‌های زیاد در حلقه‌ها
// 3. Missing Index: عدم وجود ایندکس مناسب برای فیلترها
// 4. Wrong Index Type: استفاده از ایندکس نامناسب برای نوع داده
// 5. Unused Index: ایندکس‌هایی که استفاده نمی‌شوند
//
// قانون طلایی:
// "قبل از بهینه‌سازی، همیشه EXPLAIN ANALYZE را اجرا کن.
//  برای فیلدهایی که در WHERE, JOIN, ORDER BY استفاده می‌شوند ایندکس بساز.
//  ایندکس‌های اضافی را حذف کن (هر ایندکس هزینه نگهداری دارد)."
// ============================================================================

package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// ============================================================================
// بخش 1: ساختار دیتابیس نمونه
// ============================================================================

// Schema برای مثال‌ها
const createTablesSQL = `
-- جدول کاربران
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL,
    full_name VARCHAR(200),
    age INT,
    city VARCHAR(100),
    country VARCHAR(100),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- جدول محصولات
CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL,
    category_id INT,
    stock INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- جدول سفارشات
CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    total DECIMAL(10,2) NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- جدول آیتم‌های سفارش
CREATE TABLE IF NOT EXISTS order_items (
    id SERIAL PRIMARY KEY,
    order_id INT NOT NULL,
    product_id INT NOT NULL,
    quantity INT NOT NULL,
    price DECIMAL(10,2) NOT NULL
);

-- جدول لاگ
CREATE TABLE IF NOT EXISTS query_logs (
    id SERIAL PRIMARY KEY,
    query_text TEXT,
    execution_time_ms INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
`

// ============================================================================
// بخش 2: انواع ایندکس (Index Types)
// ============================================================================

func indexTypes() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📊 INDEX TYPES IN POSTGRESQL")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ INDEX TYPE          │ USE CASE                           │ SYNTAX            │
├─────────────────────┼────────────────────────────────────┼───────────────────┤
│ B-Tree (default)    │ Equality, range, sorting           │ CREATE INDEX ... │
│ Hash                │ Equality only (faster than B-Tree) │ USING HASH        │
│ GIN                 │ Full-text search, arrays, JSON     │ USING GIN         │
│ GiST                │ Geometric data, full-text          │ USING GiST        │
│ BRIN                │ Very large tables with correlation │ USING BRIN        │
│ Partial             │ Only subset of rows                │ WHERE condition  │
│ Expression          │ Functions or expressions           │ ON (LOWER(email)) │
│ Covering           │ Include extra columns               │ INCLUDE (col)    │
└─────────────────────┴────────────────────────────────────┴───────────────────┘

📝 EXAMPLES:

-- B-Tree (پیش‌فرض)
CREATE INDEX idx_users_email ON users(email);

-- Hash (فقط معادله)
CREATE INDEX idx_users_hash ON users USING HASH (username);

-- GIN (جستجوی متن)
CREATE INDEX idx_products_search ON products USING GIN (to_tsvector('english', name || ' ' || description));

-- Partial Index (فقط کاربران فعال)
CREATE INDEX idx_active_users ON users(email) WHERE is_active = true;

-- Expression Index (جستجوی case-insensitive)
CREATE INDEX idx_users_lower_email ON users(LOWER(email));

-- Covering Index (شامل فیلدهای اضافی)
CREATE INDEX idx_users_covering ON users(email) INCLUDE (full_name, city);

-- Composite Index (چند فیلد)
CREATE INDEX idx_users_city_age ON users(city, age);
`)
}

// ============================================================================
// بخش 3: ایجاد ایندکس‌های بهینه
// ============================================================================

// IndexManager مدیریت ایندکس‌ها
type IndexManager struct {
	db *sql.DB
}

// NewIndexManager ایجاد مدیر ایندکس
func NewIndexManager(db *sql.DB) *IndexManager {
	return &IndexManager{db: db}
}

// CreateOptimalIndexes ایجاد ایندکس‌های بهینه
func (im *IndexManager) CreateOptimalIndexes() error {
	indexes := []string{
		// ایندکس برای جستجوی ایمیل (بیشترین استفاده)
		`CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_email ON users(email)`,

		// ایندکس ترکیبی برای city و age
		`CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_city_age ON users(city, age)`,

		// ایندکس partial برای کاربران فعال
		`CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_active_users ON users(email) WHERE is_active = true`,

		// ایندکس expression برای جستجوی case-insensitive
		`CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_lower_username ON users(LOWER(username))`,

		// ایندکس covering برای کوئری‌های رایج
		`CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_covering ON users(email) INCLUDE (full_name, city, age)`,

		// ایندکس برای foreign key
		`CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_orders_user_id ON orders(user_id)`,
		`CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_order_items_order_id ON order_items(order_id)`,
		`CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_order_items_product_id ON order_items(product_id)`,

		// ایندکس برای مرتب‌سازی
		`CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_orders_created_at ON orders(created_at DESC)`,

		// ایندکس برای جستجوی قیمت
		`CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_products_price ON products(price)`,

		// ایندکس GIN برای جستجوی متن
		`CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_products_search ON products USING GIN (to_tsvector('english', name || ' ' || COALESCE(description, '')))`,
	}

	for _, idx := range indexes {
		if _, err := im.db.Exec(idx); err != nil {
			log.Printf("Warning: failed to create index: %v", err)
		}
	}

	log.Println("Indexes created successfully")
	return nil
}

// ============================================================================
// بخش 4: تحلیل کوئری با EXPLAIN
// ============================================================================

// QueryAnalyzer تحلیل‌گر کوئری
type QueryAnalyzer struct {
	db *sql.DB
}

// ExplainResult نتیجه EXPLAIN
type ExplainResult struct {
	Plan  string
	Cost  string
	Rows  int64
	Width int64
}

// AnalyzeQuery تحلیل کوئری با EXPLAIN ANALYZE
func (qa *QueryAnalyzer) AnalyzeQuery(query string, args ...interface{}) error {
	explainQuery := "EXPLAIN (ANALYZE, BUFFERS, FORMAT JSON) " + query

	var explainJSON string
	err := qa.db.QueryRow(explainQuery, args...).Scan(&explainJSON)
	if err != nil {
		return fmt.Errorf("failed to explain query: %w", err)
	}

	fmt.Printf("EXPLAIN ANALYZE result:\n%s\n", explainJSON[:min(500, len(explainJSON))])
	return nil
}

// FindSlowQueries پیدا کردن کوئری‌های کند
func (qa *QueryAnalyzer) FindSlowQueries(thresholdMs int) ([]string, error) {
	query := `
		SELECT query_text 
		FROM query_logs 
		WHERE execution_time_ms > $1 
		ORDER BY execution_time_ms DESC 
		LIMIT 10
	`

	rows, err := qa.db.Query(query, thresholdMs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var slowQueries []string
	for rows.Next() {
		var queryText string
		if err := rows.Scan(&queryText); err != nil {
			continue
		}
		slowQueries = append(slowQueries, queryText)
	}

	return slowQueries, nil
}

// ============================================================================
// بخش 5: کوئری‌های بهینه (Optimized Queries)
// ============================================================================

// OptimizedQueries کلاس کوئری‌های بهینه
type OptimizedQueries struct {
	db *sql.DB
}

// GetUsersWithPagination دریافت کاربران با pagination (بهینه)
func (oq *OptimizedQueries) GetUsersWithPagination(limit, offset int) ([]User, error) {
	// استفاده از ایندکس برای ORDER BY و LIMIT
	query := `
		SELECT id, username, email, full_name, age, city, country, is_active, created_at
		FROM users
		ORDER BY id
		LIMIT $1 OFFSET $2
	`

	rows, err := oq.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.FullName, &u.Age, &u.City, &u.Country, &u.IsActive, &u.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

// GetUsersByCityAndAge دریافت کاربران بر اساس شهر و سن (با ایندکس ترکیبی)
func (oq *OptimizedQueries) GetUsersByCityAndAge(city string, minAge, maxAge int) ([]User, error) {
	// از ایندکس idx_users_city_age استفاده می‌کند
	query := `
		SELECT id, username, email, full_name, age, city, country, is_active
		FROM users
		WHERE city = $1 AND age BETWEEN $2 AND $3
		ORDER BY age
	`

	rows, err := oq.db.Query(query, city, minAge, maxAge)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.FullName, &u.Age, &u.City, &u.Country, &u.IsActive)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

// SearchUsers جستجوی کاربران (با استفاده از Expression Index)
func (oq *OptimizedQueries) SearchUsers(keyword string) ([]User, error) {
	// از ایندکس LOWER استفاده می‌کند
	query := `
		SELECT id, username, email, full_name
		FROM users
		WHERE LOWER(username) LIKE $1 OR LOWER(email) LIKE $1
		LIMIT 50
	`

	searchPattern := "%" + keyword + "%"
	rows, err := oq.db.Query(query, searchPattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.FullName)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

// GetUserOrdersWithDetails دریافت سفارشات کاربر با جزئیات (بهینه با JOIN)
func (oq *OptimizedQueries) GetUserOrdersWithDetails(userID int) ([]OrderDetail, error) {
	// استفاده از ایندکس‌های مناسب
	query := `
		SELECT 
			o.id as order_id,
			o.total,
			o.status,
			o.created_at,
			oi.product_id,
			p.name as product_name,
			oi.quantity,
			oi.price
		FROM orders o
		INNER JOIN order_items oi ON o.id = oi.order_id
		INNER JOIN products p ON oi.product_id = p.id
		WHERE o.user_id = $1
		ORDER BY o.created_at DESC
	`

	rows, err := oq.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var details []OrderDetail
	for rows.Next() {
		var d OrderDetail
		err := rows.Scan(&d.OrderID, &d.Total, &d.Status, &d.CreatedAt,
			&d.ProductID, &d.ProductName, &d.Quantity, &d.Price)
		if err != nil {
			return nil, err
		}
		details = append(details, d)
	}

	return details, nil
}

// GetTopProducts دریافت محصولات پرفروش (با aggregation)
func (oq *OptimizedQueries) GetTopProducts(limit int) ([]TopProduct, error) {
	query := `
		SELECT 
			p.id,
			p.name,
			SUM(oi.quantity) as total_sold,
			SUM(oi.quantity * oi.price) as total_revenue
		FROM products p
		INNER JOIN order_items oi ON p.id = oi.product_id
		INNER JOIN orders o ON oi.order_id = o.id
		WHERE o.status = 'completed'
		GROUP BY p.id, p.name
		ORDER BY total_sold DESC
		LIMIT $1
	`

	rows, err := oq.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []TopProduct
	for rows.Next() {
		var p TopProduct
		err := rows.Scan(&p.ID, &p.Name, &p.TotalSold, &p.TotalRevenue)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

// ============================================================================
// بخش 6: کوئری‌های غیربهینه (برای مقایسه)
// ============================================================================

// UnoptimizedQueries کوئری‌های ضعیف (برای نشان دادن مشکلات)
type UnoptimizedQueries struct {
	db *sql.DB
}

// GetUsersSlow دریافت کاربران با FULL TABLE SCAN (بدون ایندکس)
func (uq *UnoptimizedQueries) GetUsersSlow() ([]User, error) {
	// ❌ مشکل: استفاده از LIKE با wildcard در ابتدا
	query := `SELECT * FROM users WHERE email LIKE '%gmail.com'`

	rows, err := uq.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		rows.Scan(&u.ID, &u.Username, &u.Email, &u.FullName, &u.Age, &u.City, &u.Country, &u.IsActive, &u.CreatedAt)
		users = append(users, u)
	}

	return users, nil
}

// GetOrdersWithNPlusOneProblem مثال از مشکل N+1
func (uq *UnoptimizedQueries) GetOrdersWithNPlusOneProblem() ([]OrderDetail, error) {
	// ❌ مشکل: N+1 query
	// ابتدا همه سفارشات را می‌گیرد
	ordersQuery := `SELECT id, user_id, total, status FROM orders`
	ordersRows, err := uq.db.Query(ordersQuery)
	if err != nil {
		return nil, err
	}
	defer ordersRows.Close()

	var details []OrderDetail
	for ordersRows.Next() {
		var order OrderDetail
		ordersRows.Scan(&order.OrderID, &order.UserID, &order.Total, &order.Status)

		// سپس برای هر سفارش یک کوئری جداگانه می‌زند (N+1)
		itemsQuery := `SELECT product_id, quantity, price FROM order_items WHERE order_id = $1`
		itemsRows, err := uq.db.Query(itemsQuery, order.OrderID)
		if err != nil {
			continue
		}
		itemsRows.Close()
		details = append(details, order)
	}

	return details, nil
}

// ============================================================================
// بخش 7: Maintenance ایندکس
// ============================================================================

// IndexMaintenance نگهداری ایندکس‌ها
type IndexMaintenance struct {
	db *sql.DB
}

// RebuildIndexes بازسازی ایندکس‌ها
func (im *IndexMaintenance) RebuildIndexes(tableName string) error {
	query := fmt.Sprintf("REINDEX INDEX CONCURRENTLY idx_%s", tableName)
	_, err := im.db.Exec(query)
	return err
}

// AnalyzeTables به‌روزرسانی آمار دیتابیس
func (im *IndexMaintenance) AnalyzeTables() error {
	_, err := im.db.Exec("ANALYZE")
	return err
}

// FindUnusedIndexes پیدا کردن ایندکس‌های استفاده نشده
func (im *IndexMaintenance) FindUnusedIndexes() ([]string, error) {
	query := `
		SELECT schemaname, indexname 
		FROM pg_stat_user_indexes 
		WHERE idx_scan = 0 
		AND indexname NOT LIKE 'pg_%'
		ORDER BY idx_scan
	`

	rows, err := im.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var indexes []string
	for rows.Next() {
		var schema, name string
		if err := rows.Scan(&schema, &name); err != nil {
			continue
		}
		indexes = append(indexes, fmt.Sprintf("%s.%s", schema, name))
	}

	return indexes, nil
}

// GetIndexSizeAndUsage دریافت حجم و کاربرد ایندکس
func (im *IndexMaintenance) GetIndexSizeAndUsage() ([]IndexInfo, error) {
	query := `
		SELECT 
			indexname,
			pg_size_pretty(pg_indexes_size(indexname::regclass)) as size,
			idx_scan as scans,
			idx_tup_read as tuples_read,
			idx_tup_fetch as tuples_fetched
		FROM pg_stat_user_indexes
		ORDER BY idx_scan DESC
	`

	rows, err := im.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var indexes []IndexInfo
	for rows.Next() {
		var idx IndexInfo
		err := rows.Scan(&idx.Name, &idx.Size, &idx.Scans, &idx.TuplesRead, &idx.TuplesFetched)
		if err != nil {
			continue
		}
		indexes = append(indexes, idx)
	}

	return indexes, nil
}

// ============================================================================
// بخش 8: Query Monitoring (مانیتورینگ)
// ============================================================================

// QueryMonitor مانیتور کوئری‌ها
type QueryMonitor struct {
	db *sql.DB
}

// LogSlowQueries لاگ کوئری‌های کند
func (qm *QueryMonitor) LogSlowQueries(thresholdMs int) error {
	// فعال کردن logging برای کوئری‌های کند
	_, err := qm.db.Exec(fmt.Sprintf("SET log_min_duration_statement = %d", thresholdMs))
	return err
}

// GetRunningQueries دریافت کوئری‌های در حال اجرا
func (qm *QueryMonitor) GetRunningQueries() ([]RunningQuery, error) {
	query := `
		SELECT 
			pid,
			usename,
			application_name,
			state,
			query,
			now() - query_start as duration
		FROM pg_stat_activity
		WHERE state = 'active' 
		AND query NOT LIKE '%pg_stat_activity%'
		ORDER BY duration DESC
	`

	rows, err := qm.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var queries []RunningQuery
	for rows.Next() {
		var q RunningQuery
		err := rows.Scan(&q.PID, &q.Username, &q.AppName, &q.State, &q.Query, &q.Duration)
		if err != nil {
			continue
		}
		queries = append(queries, q)
	}

	return queries, nil
}

// ============================================================================
// بخش 9: Best Practices
// ============================================================================

func bestPractices() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("💡 QUERY OPTIMIZATION BEST PRACTICES")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────┐
│ INDEX DESIGN                                                  │
├─────────────────────────────────────────────────────────────────┤
│ • Create indexes for WHERE, JOIN, ORDER BY, GROUP BY         │
│ • Order columns by selectivity (most selective first)         │
│ • Use partial indexes for subsets of data                     │
│ • Use covering indexes to avoid reading the table             │
│ • Drop unused indexes (they slow down writes)                 │
│ • Use CONCURRENTLY to avoid locking in production             │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ QUERY WRITING                                                 │
├─────────────────────────────────────────────────────────────────┤
│ • Avoid SELECT * - only select needed columns                 │
│ • Use LIMIT for large result sets                             │
│ • Avoid functions on indexed columns in WHERE                 │
│   ❌ WHERE LOWER(email) = 'user@example.com'                  │
│   ✅ WHERE email = 'user@example.com'                         │
│ • Use EXISTS instead of IN for subqueries                     │
│ • Use UNION ALL instead of UNION when possible                │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ MAINTENANCE                                                   │
├─────────────────────────────────────────────────────────────────┤
│ • Run ANALYZE regularly to update statistics                  │
│ • Reindex periodically (especially after large deletes)       │
│ • Monitor index usage and remove unused ones                  │
│ • Set fillfactor for frequently updated tables                │
│ • Use partitioning for very large tables                      │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ MONITORING                                                    │
├─────────────────────────────────────────────────────────────────┤
│ • Enable slow query logging                                   │
│ • Monitor query execution time                                │
│ • Check for sequential scans                                  │
│ • Track index hit ratio                                       │
│ • Use pg_stat_statements extension                            │
└─────────────────────────────────────────────────────────────────┘
`)
}

// ============================================================================
// بخش 10: Types
// ============================================================================

type User struct {
	ID        int
	Username  string
	Email     string
	FullName  string
	Age       int
	City      string
	Country   string
	IsActive  bool
	CreatedAt time.Time
}

type OrderDetail struct {
	OrderID     int
	UserID      int
	Total       float64
	Status      string
	CreatedAt   time.Time
	ProductID   int
	ProductName string
	Quantity    int
	Price       float64
}

type TopProduct struct {
	ID           int
	Name         string
	TotalSold    int
	TotalRevenue float64
}

type IndexInfo struct {
	Name          string
	Size          string
	Scans         int64
	TuplesRead    int64
	TuplesFetched int64
}

type RunningQuery struct {
	PID      int
	Username string
	AppName  string
	State    string
	Query    string
	Duration interface{}
}

// ============================================================================
// بخش 11: Main
// ============================================================================

func main() {
	fmt.Println(stringsRepeat("=", 80))
	fmt.Println("🎯 QUERY OPTIMIZATION & INDEX GUIDE")
	fmt.Println("PostgreSQL Performance Tuning")
	fmt.Println(stringsRepeat("=", 80))

	indexTypes()
	bestPractices()

	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📊 QUICK REFERENCE - INDEX DECISIONS")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ QUERY PATTERN                    │ BEST INDEX TYPE          │
├──────────────────────────────────┼──────────────────────────┤
│ WHERE column = value             │ B-Tree                   │
│ WHERE column > value             │ B-Tree                   │
│ ORDER BY column                  │ B-Tree                   │
│ LIKE 'value%'                    │ B-Tree                   │
│ LIKE '%value%'                   │ GIN (trigram)           │
│ Full-text search                 │ GIN (tsvector)          │
│ JSON field access                │ GIN                      │
│ Array operations                 │ GIN                      │
│ Geometric data                   │ GiST                     │
│ Very large, correlated data      │ BRIN                     │
│ Partial data                     │ Partial Index            │
│ Function result                  │ Expression Index         │
└──────────────────────────────────┴──────────────────────────┘

🔍 COMMANDS:

   # Analyze query
   EXPLAIN (ANALYZE, BUFFERS) SELECT * FROM users WHERE email = 'test';

   # Find missing indexes
   SELECT * FROM pg_stat_user_indexes WHERE idx_scan = 0;

   # Show index usage
   SELECT schemaname, indexname, idx_scan FROM pg_stat_user_indexes;

   # Rebuild index
   REINDEX INDEX CONCURRENTLY idx_name;

   # Update statistics
   ANALYZE users;

   # Show table size
   SELECT pg_size_pretty(pg_total_relation_size('users'));
`)

	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🎯 QUERY OPTIMIZATION - COMPLETE")
	fmt.Println("Ready to optimize your database performance!")
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
