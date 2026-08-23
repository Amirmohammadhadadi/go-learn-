// ============================================================================
// FILE: advanced_testing_guide.go
// TITLE: راهنمای تست پیشرفته در Go - Integration Tests با Testcontainers و Load Testing با k6/Vegeta
// HOW TO RUN: go test -v -tags=integration ./...
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - تست پیشرفته چیست و چرا نیاز است؟
// ============================================================================
//
// تست پیشرفته شامل دو حوزه اصلی است:
//
// 1. Integration Tests (تست یکپارچگی)
//    - تست تعامل بین سرویس‌ها و وابستگی‌ها
//    - تست با دیتابیس واقعی (PostgreSQL, Redis, etc.)
//    - تست API endpoints واقعی
//    - Testcontainers: کتابخانه‌ای برای اجرای کانتینرهای داکر در تست‌ها
//
// 2. Load Testing (تست بار)
//    - شبیه‌سازی ترافیک سنگین
//    - اندازه‌گیری performance تحت فشار
//    - شناسایی bottlenecks
//    - k6: ابزار مدرن تست بار (خروجی‌های غنی)
//    - Vegeta: ابزار ساده و سریع خط فرمان
//
// قانون طلایی:
// "همیشه integration tests را با محیط واقعی (دیتابیس واقعی) اجرا کن.
//  از testcontainers برای ایزوله کردن تست‌ها استفاده کن.
//  قبل از هر release، load testing انجام بده.
//  تست‌های integration را در CI/CD اجرا کن."
// ============================================================================

//go:build integration
// +build integration

package __advanced_testing

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/suite"
)

// ============================================================================
// بخش 1: Testcontainers - راه‌اندازی
// ============================================================================

const (
	// TestContainers images
	postgresImage = "postgres:16-alpine"
	redisImage    = "redis:7-alpine"
	mysqlImage    = "mysql:8.0"
	mongoImage    = "mongo:7-alpine"

	// Test databases
	testDBName     = "testdb"
	testDBUser     = "testuser"
	testDBPassword = "testpass"
)

// TestContainersManager مدیریت کانتینرهای تست
type TestContainersManager struct {
	pool      *dockertest.Pool
	resources []*dockertest.Resource
}

// NewTestContainersManager ایجاد مدیر کانتینر
func NewTestContainersManager() (*TestContainersManager, error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, fmt.Errorf("could not connect to docker: %w", err)
	}

	pool.MaxWait = 2 * time.Minute

	return &TestContainersManager{
		pool:      pool,
		resources: make([]*dockertest.Resource, 0),
	}, nil
}

// StartPostgresContainer راه‌اندازی کانتینر PostgreSQL
func (m *TestContainersManager) StartPostgresContainer(ctx context.Context) (string, int, error) {
	resource, err := m.pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "16-alpine",
		Env: []string{
			fmt.Sprintf("POSTGRES_DB=%s", testDBName),
			fmt.Sprintf("POSTGRES_USER=%s", testDBUser),
			fmt.Sprintf("POSTGRES_PASSWORD=%s", testDBPassword),
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		return "", 0, err
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	m.resources = append(m.resources, resource)

	// منتظر آماده شدن دیتابیس
	if err := m.pool.Retry(func() error {
		db, err := sqlx.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
			testDBUser, testDBPassword, hostAndPort, testDBName))
		if err != nil {
			return err
		}
		defer db.Close()
		return db.Ping()
	}); err != nil {
		return "", 0, err
	}

	return hostAndPort, resource.GetPort("5432/tcp"), nil
}

// StartRedisContainer راه‌اندازی کانتینر Redis
func (m *TestContainersManager) StartRedisContainer(ctx context.Context) (string, int, error) {
	resource, err := m.pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "redis",
		Tag:        "7-alpine",
		Cmd:        []string{"redis-server", "--appendonly", "yes"},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
	})
	if err != nil {
		return "", 0, err
	}

	hostAndPort := resource.GetHostPort("6379/tcp")
	m.resources = append(m.resources, resource)

	return hostAndPort, resource.GetPort("6379/tcp"), nil
}

// StartMySQLContainer راه‌اندازی کانتینر MySQL
func (m *TestContainersManager) StartMySQLContainer(ctx context.Context) (string, int, error) {
	resource, err := m.pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "mysql",
		Tag:        "8.0",
		Env: []string{
			"MYSQL_ROOT_PASSWORD=rootpass",
			"MYSQL_DATABASE=testdb",
			"MYSQL_USER=testuser",
			"MYSQL_PASSWORD=testpass",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
	})
	if err != nil {
		return "", 0, err
	}

	hostAndPort := resource.GetHostPort("3306/tcp")
	m.resources = append(m.resources, resource)

	return hostAndPort, resource.GetPort("3306/tcp"), nil
}

// Cleanup پاک کردن همه کانتینرها
func (m *TestContainersManager) Cleanup() {
	for _, resource := range m.resources {
		if err := m.pool.Purge(resource); err != nil {
			log.Printf("Failed to purge resource: %v", err)
		}
	}
}

// ============================================================================
// بخش 2: Test Suite (با Testcontainers)
// ============================================================================

// IntegrationTestSuite تست یکپارچگی با دیتابیس واقعی
type IntegrationTestSuite struct {
	suite.Suite
	containerManager *TestContainersManager
	db               *sqlx.DB
	postgresHost     string
	postgresPort     int
	redisHost        string
	redisPort        int
	server           *httptest.Server
}

// SetupSuite قبل از همه تست‌ها (یک بار)
func (s *IntegrationTestSuite) SetupSuite() {
	var err error

	// راه‌اندازی Testcontainers manager
	s.containerManager, err = NewTestContainersManager()
	s.Require().NoError(err)

	// راه‌اندازی PostgreSQL
	ctx := context.Background()
	host, port, err := s.containerManager.StartPostgresContainer(ctx)
	s.Require().NoError(err)
	s.postgresHost = host
	s.postgresPort = port

	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		testDBUser, testDBPassword, host, testDBName)
	s.db, err = sqlx.Connect("postgres", dsn)
	s.Require().NoError(err)

	// ایجاد جداول
	s.createTables()

	// راه‌اندازی HTTP server برای تست
	s.server = httptest.NewServer(s.setupRouter())
}

// TearDownSuite بعد از همه تست‌ها (یک بار)
func (s *IntegrationTestSuite) TearDownSuite() {
	if s.db != nil {
		s.db.Close()
	}
	if s.containerManager != nil {
		s.containerManager.Cleanup()
	}
	if s.server != nil {
		s.server.Close()
	}
}

// SetupTest قبل از هر تست
func (s *IntegrationTestSuite) SetupTest() {
	// پاک کردن داده‌ها قبل از هر تست
	s.clearTables()
}

// createTables ایجاد جداول دیتابیس
func (s *IntegrationTestSuite) createTables() {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			email VARCHAR(255) NOT NULL UNIQUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS products (
			id SERIAL PRIMARY KEY,
			name VARCHAR(200) NOT NULL,
			price DECIMAL(10,2) NOT NULL,
			stock INT NOT NULL DEFAULT 0
		)`,
		`CREATE TABLE IF NOT EXISTS orders (
			id SERIAL PRIMARY KEY,
			user_id INT NOT NULL REFERENCES users(id),
			total DECIMAL(10,2) NOT NULL,
			status VARCHAR(50) DEFAULT 'pending',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	for _, query := range queries {
		_, err := s.db.Exec(query)
		s.Require().NoError(err)
	}
}

// clearTables پاک کردن داده‌ها
func (s *IntegrationTestSuite) clearTables() {
	_, err := s.db.Exec("TRUNCATE TABLE orders, products, users RESTART IDENTITY CASCADE")
	s.Require().NoError(err)
}

// setupRouter تنظیم router برای تست
func (s *IntegrationTestSuite) setupRouter() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/users", s.createUserHandler).Methods("POST")
	r.HandleFunc("/users/{id}", s.getUserHandler).Methods("GET")
	r.HandleFunc("/users", s.listUsersHandler).Methods("GET")
	return r
}

// ============================================================================
// بخش 3: Handlers برای تست
// ============================================================================

type User struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (s *IntegrationTestSuite) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user User
	query := `INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id, name, email, created_at`
	err := s.db.QueryRowx(query, req.Name, req.Email).StructScan(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (s *IntegrationTestSuite) getUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var user User
	err := s.db.Get(&user, "SELECT id, name, email, created_at FROM users WHERE id = $1", id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (s *IntegrationTestSuite) listUsersHandler(w http.ResponseWriter, r *http.Request) {
	var users []User
	err := s.db.Select(&users, "SELECT id, name, email, created_at FROM users ORDER BY id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
}

// ============================================================================
// بخش 4: Integration Tests (تست واقعی)
// ============================================================================

// TestCreateUser تست ایجاد کاربر
func (s *IntegrationTestSuite) TestCreateUser() {
	// ایجاد درخواست
	userReq := CreateUserRequest{
		Name:  "Ali Rezaei",
		Email: "ali@example.com",
	}

	body, _ := json.Marshal(userReq)
	resp, err := http.Post(s.server.URL+"/users", "application/json", bytes.NewBuffer(body))
	s.Require().NoError(err)
	defer resp.Body.Close()

	// بررسی status code
	s.Equal(http.StatusCreated, resp.StatusCode)

	// بررسی پاسخ
	var user User
	err = json.NewDecoder(resp.Body).Decode(&user)
	s.Require().NoError(err)

	s.Equal(userReq.Name, user.Name)
	s.Equal(userReq.Email, user.Email)
	s.NotZero(user.ID)

	// بررسی ذخیره در دیتابیس
	var dbUser User
	err = s.db.Get(&dbUser, "SELECT * FROM users WHERE id = $1", user.ID)
	s.Require().NoError(err)
	s.Equal(userReq.Name, dbUser.Name)
}

// TestGetUser تست دریافت کاربر
func (s *IntegrationTestSuite) TestGetUser() {
	// ایجاد کاربر تست در دیتابیس
	testUser := User{Name: "Test User", Email: "test@example.com"}
	err := s.db.QueryRowx(
		"INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id",
		testUser.Name, testUser.Email,
	).Scan(&testUser.ID)
	s.Require().NoError(err)

	// درخواست GET
	resp, err := http.Get(fmt.Sprintf("%s/users/%d", s.server.URL, testUser.ID))
	s.Require().NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusOK, resp.StatusCode)

	var user User
	err = json.NewDecoder(resp.Body).Decode(&user)
	s.Require().NoError(err)

	s.Equal(testUser.Name, user.Name)
	s.Equal(testUser.Email, user.Email)
}

// TestGetUserNotFound تست کاربر ناموجود
func (s *IntegrationTestSuite) TestGetUserNotFound() {
	resp, err := http.Get(s.server.URL + "/users/99999")
	s.Require().NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusNotFound, resp.StatusCode)
}

// TestListUsers تست لیست کاربران
func (s *IntegrationTestSuite) TestListUsers() {
	// ایجاد چند کاربر تست
	for i := 0; i < 3; i++ {
		_, err := s.db.Exec("INSERT INTO users (name, email) VALUES ($1, $2)",
			fmt.Sprintf("User%d", i), fmt.Sprintf("user%d@example.com", i))
		s.Require().NoError(err)
	}

	// درخواست GET
	resp, err := http.Get(s.server.URL + "/users")
	s.Require().NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusOK, resp.StatusCode)

	var users []User
	err = json.NewDecoder(resp.Body).Decode(&users)
	s.Require().NoError(err)

	s.GreaterOrEqual(len(users), 3)
}

// ============================================================================
// بخش 5: اجرای Integration Tests
// ============================================================================

// TestIntegrationSuite اجرای تست‌های یکپارچگی
func TestIntegrationSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}
	suite.Run(t, new(IntegrationTestSuite))
}

// ============================================================================
// بخش 6: Load Testing با Vegeta (Code Example)
// ============================================================================

// VegetaLoadTest مثال تست بار با Vegeta
func VegetaLoadTest() error {
	// نصب: go get -u github.com/tsenart/vegeta/v12

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📊 LOAD TESTING WITH VEGETA")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
# Vegeta commands:

# Single endpoint test
$ echo "GET http://localhost:8080/api/users" | vegeta attack -rate=100 -duration=10s | vegeta report

# Multiple endpoints from file
$ cat targets.txt
GET http://localhost:8080/api/users
GET http://localhost:8080/api/products/1
POST http://localhost:8080/api/users

$ vegeta attack -targets=targets.txt -rate=50 -duration=30s -output=results.bin
$ vegeta report -type=json results.bin
$ vegeta plot -title="Load Test" results.bin > plot.html

# Advanced options
$ echo "GET http://localhost:8080/api/users" | vegeta attack \
    -rate=100 \
    -duration=60s \
    -workers=10 \
    -timeout=5s \
    -lazy \
    -header="Authorization: Bearer token123" \
    | vegeta report -type=histogram
`)

	return nil
}

// ============================================================================
// بخش 7: Load Testing با k6 (JavaScript)
// ============================================================================

func k6LoadTestExample() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📊 LOAD TESTING WITH K6")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
# Install k6:
# macOS: brew install k6
# Linux: sudo apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys C5AD17C747E3415A3642D57D77C6C491D6AC1D69
#        echo "deb https://dl.k6.io/deb stable main" | sudo tee /etc/apt/sources.list.d/k6.list
#        sudo apt-get update && sudo apt-get install k6

# Run test:
$ k6 run load_test.js

# Example k6 script (load_test.js):
# ============================================================================
# import http from 'k6/http';
# import { check, sleep } from 'k6';
#
# export const options = {
#   stages: [
#     { duration: '30s', target: 20 },  // ramp up
#     { duration: '1m', target: 20 },   // steady
#     { duration: '30s', target: 0 },   // ramp down
#   ],
#   thresholds: {
#     http_req_duration: ['p(95)<500'],  // 95% requests < 500ms
#     http_req_failed: ['rate<0.01'],    // error rate < 1%
#   },
# };
#
# export default function () {
#   const res = http.get('http://localhost:8080/api/users');
#   check(res, {
#     'status is 200': (r) => r.status === 200,
#     'response time < 200ms': (r) => r.timings.duration < 200,
#   });
#   sleep(1);
# }
`)
}

// ============================================================================
// بخش 8: Load Testing Implementation in Go
// ============================================================================

// LoadTestRunner اجرای تست بار در Go
type LoadTestRunner struct {
	url         string
	concurrency int
	duration    time.Duration
	method      string
	headers     map[string]string
	body        []byte
}

// NewLoadTestRunner ایجاد runner جدید
func NewLoadTestRunner(url string, concurrency int, duration time.Duration) *LoadTestRunner {
	return &LoadTestRunner{
		url:         url,
		concurrency: concurrency,
		duration:    duration,
		method:      "GET",
		headers:     make(map[string]string),
	}
}

// SetMethod تنظیم متد HTTP
func (l *LoadTestRunner) SetMethod(method string) {
	l.method = method
}

// SetHeader تنظیم header
func (l *LoadTestRunner) SetHeader(key, value string) {
	l.headers[key] = value
}

// SetBody تنظیم body
func (l *LoadTestRunner) SetBody(body []byte) {
	l.body = body
}

// Run اجرای تست بار
func (l *LoadTestRunner) Run() *LoadTestResult {
	ctx, cancel := context.WithTimeout(context.Background(), l.duration)
	defer cancel()

	var wg sync.WaitGroup
	results := make(chan RequestResult, 1000)

	// شروع workerها
	for i := 0; i < l.concurrency; i++ {
		wg.Add(1)
		go l.worker(ctx, &wg, results)
	}

	// منتظر اتمام
	go func() {
		wg.Wait()
		close(results)
	}()

	// جمع‌آوری نتایج
	return l.collectResults(results)
}

func (l *LoadTestRunner) worker(ctx context.Context, wg *sync.WaitGroup, results chan<- RequestResult) {
	defer wg.Done()

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			start := time.Now()

			req, _ := http.NewRequest(l.method, l.url, bytes.NewBuffer(l.body))
			for k, v := range l.headers {
				req.Header.Set(k, v)
			}

			resp, err := client.Do(req)
			duration := time.Since(start)

			result := RequestResult{
				Duration: duration,
				Success:  err == nil && resp != nil && resp.StatusCode < 400,
			}

			if resp != nil {
				result.StatusCode = resp.StatusCode
				resp.Body.Close()
			}

			if err != nil {
				result.Error = err.Error()
			}

			results <- result
		}
	}
}

type RequestResult struct {
	Duration   time.Duration
	Success    bool
	StatusCode int
	Error      string
}

type LoadTestResult struct {
	TotalRequests   int
	SuccessCount    int
	FailureCount    int
	MinDuration     time.Duration
	MaxDuration     time.Duration
	AvgDuration     time.Duration
	P50Duration     time.Duration
	P95Duration     time.Duration
	P99Duration     time.Duration
	RequestsPerSec  float64
	Errors          map[string]int
	StatusCodeCount map[int]int
}

func (l *LoadTestRunner) collectResults(results <-chan RequestResult) *LoadTestResult {
	res := &LoadTestResult{
		MinDuration:     time.Hour,
		Errors:          make(map[string]int),
		StatusCodeCount: make(map[int]int),
	}

	var durations []time.Duration
	var totalDuration time.Duration

	for result := range results {
		res.TotalRequests++

		if result.Success {
			res.SuccessCount++
		} else {
			res.FailureCount++
			if result.Error != "" {
				res.Errors[result.Error]++
			}
		}

		if result.StatusCode > 0 {
			res.StatusCodeCount[result.StatusCode]++
		}

		durations = append(durations, result.Duration)
		totalDuration += result.Duration

		if result.Duration < res.MinDuration {
			res.MinDuration = result.Duration
		}
		if result.Duration > res.MaxDuration {
			res.MaxDuration = result.Duration
		}
	}

	if res.TotalRequests > 0 {
		res.AvgDuration = totalDuration / time.Duration(res.TotalRequests)

		// محاسبه percentiles
		sortDurations(durations)
		res.P50Duration = percentile(durations, 50)
		res.P95Duration = percentile(durations, 95)
		res.P99Duration = percentile(durations, 99)
	}

	return res
}

func percentile(durations []time.Duration, p float64) time.Duration {
	if len(durations) == 0 {
		return 0
	}
	idx := int(float64(len(durations)) * p / 100)
	if idx >= len(durations) {
		idx = len(durations) - 1
	}
	return durations[idx]
}

func sortDurations(durations []time.Duration) {
	for i := 0; i < len(durations)-1; i++ {
		for j := i + 1; j < len(durations); j++ {
			if durations[i] > durations[j] {
				durations[i], durations[j] = durations[j], durations[i]
			}
		}
	}
}

// PrintResult چاپ نتایج
func (r *LoadTestResult) PrintResult() {
	fmt.Println("\n=== Load Test Results ===")
	fmt.Printf("Total Requests:  %d\n", r.TotalRequests)
	fmt.Printf("Successful:      %d (%.2f%%)\n", r.SuccessCount, float64(r.SuccessCount)/float64(r.TotalRequests)*100)
	fmt.Printf("Failed:          %d (%.2f%%)\n", r.FailureCount, float64(r.FailureCount)/float64(r.TotalRequests)*100)
	fmt.Printf("Requests/sec:    %.2f\n", r.RequestsPerSec)
	fmt.Println("\n--- Latency ---")
	fmt.Printf("Min:            %v\n", r.MinDuration)
	fmt.Printf("Max:            %v\n", r.MaxDuration)
	fmt.Printf("Avg:            %v\n", r.AvgDuration)
	fmt.Printf("P50:            %v\n", r.P50Duration)
	fmt.Printf("P95:            %v\n", r.P95Duration)
	fmt.Printf("P99:            %v\n", r.P99Duration)

	if len(r.StatusCodeCount) > 0 {
		fmt.Println("\n--- Status Codes ---")
		for code, count := range r.StatusCodeCount {
			fmt.Printf("%d: %d (%.2f%%)\n", code, count, float64(count)/float64(r.TotalRequests)*100)
		}
	}

	if len(r.Errors) > 0 {
		fmt.Println("\n--- Errors ---")
		for err, count := range r.Errors {
			fmt.Printf("%s: %d\n", err, count)
		}
	}
}

// ============================================================================
// بخش 9: Example Usage
// ============================================================================

func runLoadTestExample() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🚀 Running Load Test Example")
	fmt.Println(strings.Repeat("=", 80))

	// ایجاد سرور تست
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// شبیه‌سازی کار
		time.Sleep(10 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}))
	defer server.Close()

	// اجرای تست بار
	runner := NewLoadTestRunner(server.URL, 50, 5*time.Second)
	result := runner.Run()
	result.PrintResult()
}

// ============================================================================
// بخش 10: Best Practices
// ============================================================================

func bestPractices() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("💡 ADVANCED TESTING BEST PRACTICES")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ INTEGRATION TESTS                                                         │
├─────────────────────────────────────────────────────────────────────────────┤
│  • Always use testcontainers for real dependencies                        │
│  • Clean up between tests                                                  │
│  • Use build tags to separate integration tests                           │
│  • Run integration tests in CI/CD                                         │
│  • Keep tests independent and idempotent                                   │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ LOAD TESTING                                                              │
├─────────────────────────────────────────────────────────────────────────────┤
│  • Test in production-like environment                                    │
│  • Use realistic traffic patterns                                          │
│  • Monitor system metrics during tests                                    │
│  • Define SLOs (Service Level Objectives)                                 │
│  • Gradually increase load                                                │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ COMMAND REFERENCE                                                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  # Run integration tests                                                   │
│  $ go test -tags=integration -v ./...                                      │
│                                                                             │
│  # Run with short mode (skip long tests)                                   │
│  $ go test -short -v ./...                                                 │
│                                                                             │
│  # Run load test with vegeta                                               │
│  $ echo "GET http://localhost:8080/api/users" | vegeta attack -rate=100 -duration=30s | vegeta report│
│                                                                             │
│  # Run load test with k6                                                   │
│  $ k6 run load_test.js                                                     │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
`)
}

// ============================================================================
// بخش 11: Main
// ============================================================================

func main() {
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("🎯 ADVANCED TESTING IN GO")
	fmt.Println("Integration Tests with Testcontainers | Load Testing with k6/Vegeta")
	fmt.Println(strings.Repeat("=", 80))

	bestPractices()
	VegetaLoadTest()
	k6LoadTestExample()

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📝 QUICK START")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
# 1. Install test dependencies
$ go get github.com/ory/dockertest/v3
$ go get github.com/stretchr/testify
$ go get github.com/jmoiron/sqlx

# 2. Install load testing tools
$ go install github.com/tsenart/vegeta/v12@latest
$ brew install k6  # or download from https://k6.io

# 3. Run integration tests
$ go test -tags=integration -v ./...

# 4. Run load test
$ echo "GET http://localhost:8080/api/users" | vegeta attack -rate=100 -duration=30s | vegeta report
`)

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🎯 ADVANCED TESTING - COMPLETE")
	fmt.Println("Ready to test your Go applications with confidence!")
	fmt.Println(strings.Repeat("=", 80))
}

// تابع کمکی
func stringsRepeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}

// اضافه کردن importهای لازم
var _ = sqlx.Connect
var _ = dockertest.NewPool
