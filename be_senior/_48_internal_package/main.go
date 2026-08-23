// ============================================================================
// FILE: internal_package_guide.go
// TITLE: راهنمای کامل internal Package در Go - Encapsulation در سطح ماژول
// HOW TO RUN: go run internal_package_guide.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - internal Package چیست؟
// ============================================================================
//
// internal package یک ویژگی خاص در Go است که از Go 1.4 معرفی شده است.
// هدف: ایجاد encapsulation در سطح ماژول (module-level encapsulation)
//
// قوانین دسترسی به internal package:
// 1. کد داخل internal فقط توسط پکیج‌های داخل همان دایرکتوری والد (parent) قابل استفاده است
// 2. پکیج‌های خارج از درخت دایرکتوری که internal در آن قرار دارد، نمی‌توانند از آن import کنند
// 3. این قانون توسط کامپایلر Go强制执行 می‌شود
//
// ساختار typical:
//   /myproject
//   ├── go.mod
//   ├── internal/
//   │   ├── auth/
//   │   │   └── token.go      // فقط توسط پکیج‌های داخل myproject قابل استفاده
//   │   └── database/
//   │       └── postgres.go   // فقط توسط پکیج‌های داخل myproject قابل استفاده
//   ├── cmd/
//   │   └── server/
//   │       └── main.go        // می‌تواند internal را استفاده کند
//   └── pkg/
//       └── public/
//           └── api.go         // می‌تواند internal را استفاده کند
//
// قانون طلایی:
// "از internal برای کدهایی که نباید خارج از ماژول استفاده شوند استفاده کن.
//  هر چیزی که public API ماژول تو است را در pkg/ یا ریشه پروژه قرار بده.
//  internal بهترین روش برای encapsulation در ماژول‌های Go است."
// ============================================================================

package __internal_package

import (
	"fmt"
	"log"
	"strings"
)

// ============================================================================
// بخش 1: ساختار دایرکتوری و قوانین دسترسی
// ============================================================================

/*
ساختار مثال:

myproject/
├── go.mod                    (module: github.com/myuser/myproject)
├── internal/
│   ├── auth/
│   │   ├── token.go          // package auth
│   │   └── middleware.go     // package auth
│   ├── database/
│   │   └── postgres.go       // package database
│   ├── cache/
│   │   └── redis.go          // package cache
│   └── logger/
│       └── logger.go         // package logger
├── cmd/
│   ├── server/
│   │   └── main.go           // package main - می‌تواند internal را import کند
│   └── worker/
│       └── main.go           // package main - می‌تواند internal را import کند
├── pkg/
│   ├── api/
│   │   └── handler.go        // package api - می‌تواند internal را import کند
│   └── client/
│       └── client.go         // package client - می‌تواند internal را import کند
└── tests/
    └── integration_test.go   // می‌تواند internal را import کند (در همان ماژول)

❌ خارج از ماژول (نمی‌تواند import کند):
   otherproject/
   └── main.go                // نمی‌تواند github.com/myuser/myproject/internal را import کند
*/

// ============================================================================
// بخش 2: مثال Internal Package - Auth
// ============================================================================

// این فایل فرضی در internal/auth/token.go قرار دارد
// package auth

type TokenGenerator interface {
	Generate(userID string) (string, error)
	Validate(token string) (string, error)
}

// JWTGenerator (internal - فقط داخل ماژول قابل استفاده)
type JWTGenerator struct {
	secretKey string
}

func NewJWTGenerator(secretKey string) *JWTGenerator {
	return &JWTGenerator{secretKey: secretKey}
}

func (g *JWTGenerator) Generate(userID string) (string, error) {
	// پیاده‌سازی واقعی JWT
	return "jwt-token-" + userID, nil
}

func (g *JWTGenerator) Validate(token string) (string, error) {
	// پیاده‌سازی واقعی JWT
	if strings.HasPrefix(token, "jwt-token-") {
		return strings.TrimPrefix(token, "jwt-token-"), nil
	}
	return "", fmt.Errorf("invalid token")
}

// ============================================================================
// بخش 3: مثال Internal Package - Database
// ============================================================================

// فرضی در internal/database/postgres.go
// package database

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

type Connection struct {
	config Config
	pool   interface{} // در واقع *sql.DB
}

func NewConnection(config Config) (*Connection, error) {
	// ایجاد اتصال واقعی
	return &Connection{config: config}, nil
}

func (c *Connection) Query(query string, args ...interface{}) ([]map[string]interface{}, error) {
	// اجرای کوئری واقعی
	return nil, nil
}

func (c *Connection) Exec(query string, args ...interface{}) error {
	// اجرای کوئری واقعی
	return nil
}

func (c *Connection) Close() error {
	// بستن اتصال
	return nil
}

// ============================================================================
// بخش 4: مثال Internal Package - Logger
// ============================================================================

// فرضی در internal/logger/logger.go
// package logger

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

type Logger struct {
	level Level
}

func NewLogger(level Level) *Logger {
	return &Logger{level: level}
}

func (l *Logger) Debug(msg string, keysAndValues ...interface{}) {
	if l.level <= LevelDebug {
		log.Printf("[DEBUG] %s %v", msg, keysAndValues)
	}
}

func (l *Logger) Info(msg string, keysAndValues ...interface{}) {
	if l.level <= LevelInfo {
		log.Printf("[INFO] %s %v", msg, keysAndValues)
	}
}

func (l *Logger) Warn(msg string, keysAndValues ...interface{}) {
	if l.level <= LevelWarn {
		log.Printf("[WARN] %s %v", msg, keysAndValues)
	}
}

func (l *Logger) Error(msg string, keysAndValues ...interface{}) {
	if l.level <= LevelError {
		log.Printf("[ERROR] %s %v", msg, keysAndValues)
	}
}

// ============================================================================
// بخش 5: استفاده از Internal Packages در cmd/server
// ============================================================================

// فرضی در cmd/server/main.go
// package main

/*
import (
    "net/http"
    "github.com/myuser/myproject/internal/auth"
    "github.com/myuser/myproject/internal/database"
    "github.com/myuser/myproject/internal/logger"
)

func main() {
    // استفاده از internal packages (مجاز است)
    dbConfig := database.Config{
        Host: "localhost",
        Port: 5432,
        User: "postgres",
        DBName: "mydb",
    }
    db, _ := database.NewConnection(dbConfig)
    defer db.Close()

    logger := logger.NewLogger(logger.LevelInfo)
    logger.Info("Application started")

    jwtGen := auth.NewJWTGenerator("my-secret-key")

    // استفاده در HTTP handler
    http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
        token, _ := jwtGen.Generate("user-123")
        w.Write([]byte(token))
    })

    http.ListenAndServe(":8080", nil)
}
*/

// ============================================================================
// بخش 6: استفاده از Internal Packages در pkg/api
// ============================================================================

// فرضی در pkg/api/handler.go
// package api

/*
import (
    "net/http"
    "github.com/myuser/myproject/internal/auth"
    "github.com/myuser/myproject/internal/logger"
)

type Handler struct {
    auth   *auth.JWTGenerator
    logger *logger.Logger
}

func NewHandler(auth *auth.JWTGenerator, logger *logger.Logger) *Handler {
    return &Handler{
        auth:   auth,
        logger: logger,
    }
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    token := r.Header.Get("Authorization")
    userID, err := h.auth.Validate(token)
    if err != nil {
        h.logger.Error("Invalid token", "error", err)
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    h.logger.Info("Request authenticated", "user_id", userID)
    // ادامه پردازش
}
*/

// ============================================================================
// بخش 7: سلسله‌مراتب Internal (Nested Internal)
// ============================================================================

/*
ساختار با nested internal:

myproject/
├── go.mod
├── internal/
│   ├── auth/
│   │   ├── token.go
│   │   └── internal/           # nested internal
│   │       └── crypto.go       # فقط توسط internal/auth قابل استفاده
│   └── database/
│       └── internal/           # nested internal
│           └── pool.go         # فقط توسط internal/database قابل استفاده
└── cmd/
    └── server/
        └── main.go              # می‌تواند internal/auth را استفاده کند
                                 # ولی نمی‌تواند internal/auth/internal را استفاده کند
*/

// فرضی در internal/auth/internal/crypto.go
// package internal (این package درون internal/auth قرار دارد)

/*
// این کد فقط توسط package auth قابل استفاده است
type cryptoHelper struct{}

func (c *cryptoHelper) encrypt(data []byte) ([]byte, error) {
    // پیاده‌سازی
    return data, nil
}

func (c *cryptoHelper) decrypt(data []byte) ([]byte, error) {
    // پیاده‌سازی
    return data, nil
}
*/

// ============================================================================
// بخش 8: Internal vs Pkg - تفاوت‌ها و کاربردها
// ============================================================================

func compareInternalAndPkg() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📊 INTERNAL vs PKG")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ CRITERION              │ internal/          │ pkg/                         │
├────────────────────────┼────────────────────┼──────────────────────────────┤
│ Visibility             │ Module-private     │ Public (other modules can use)│
│ Enforced by            │ Go compiler        │ Convention (not enforced)    │
│ Use case               │ Internal implementation│ Reusable libraries        │
│ Stability guarantee    │ None (can change)  │ Semantic versioning expected │
│ Documentation          │ Internal only      │ Public documentation         │
│ Example                │ internal/auth/     │ pkg/api/client/              │
└────────────────────────┴────────────────────┴──────────────────────────────┘

📁 RECOMMENDED STRUCTURE:

   mymodule/
   ├── go.mod
   ├── internal/           # Private to this module
   │   ├── auth/          # Authentication helpers
   │   ├── database/      # Database layer
   │   ├── cache/         # Caching layer
   │   ├── queue/         # Message queue
   │   ├── logger/        # Logging (if not public)
   │   └── config/        # Configuration loading
   ├── pkg/               # Public, reusable packages
   │   ├── api/           # Public API client
   │   │   └── client.go
   │   ├── errors/        # Public error types
   │   │   └── errors.go
   │   └── utils/         # Public utilities (use sparingly)
   │       └── string.go
   └── cmd/               # Executables (can use internal)
       ├── server/
       │   └── main.go
       └── cli/
           └── main.go
`)
}

// ============================================================================
// بخش 9: Best Practices
// ============================================================================

func bestPractices() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("💡 INTERNAL PACKAGE BEST PRACTICES")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ 1. WHAT GOES IN internal/                                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│    ✅ Database connection code                                            │
│    ✅ Repository implementations                                          │
│    ✅ Authentication/authorization logic                                  │
│    ✅ Configuration loading                                               │
│    ✅ Internal service implementations                                    │
│    ✅ Shared utilities that shouldn't be public                           │
│    ✅ Code that may change frequently                                     │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 2. WHAT GOES IN pkg/ (or root)                                           │
├─────────────────────────────────────────────────────────────────────────────┤
│    ✅ Public API clients                                                  │
│    ✅ SDK for external use                                                │
│    ✅ Shared types that other modules need                                │
│    ✅ Stable, well-documented code                                        │
│    ✅ Code that follows semantic versioning                               │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 3. WHAT GOES IN cmd/                                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│    ✅ main package for executables                                        │
│    ✅ Can import internal/ and pkg/                                       │
│    ✅ Keep main.go thin (just wiring)                                     │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 4. DON'TS                                                                │
├─────────────────────────────────────────────────────────────────────────────┤
│    ❌ Don't put everything in internal/                                   │
│    ❌ Don't use internal for code that should be reusable                 │
│    ❌ Don't import internal from outside the module                       │
│    ❌ Don't create deep nested internal (internal/internal/...)          │
│    ❌ Don't put test code in internal/                                    │
└─────────────────────────────────────────────────────────────────────────────┘
`)
}

// ============================================================================
// بخش 10: Common Mistakes
// ============================================================================

func commonMistakes() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("⚠️ COMMON MISTAKES WITH internal")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ MISTAKE 1: Putting reusable libraries in internal                         │
├─────────────────────────────────────────────────────────────────────────────┤
│    ❌ internal/stringutil/ (can't be used by other projects)              │
│    ✅ pkg/stringutil/ (can be imported by others)                         │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ MISTAKE 2: Putting everything in internal/                                │
├─────────────────────────────────────────────────────────────────────────────┤
│    ❌ internal/                                                           │
│    │   ├── models/                                                        │
│    │   ├── handlers/                                                      │
│    │   ├── services/                                                      │
│    │   └── repositories/                                                  │
│    ✅ Use internal for truly private code only                            │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ MISTAKE 3: Deep nesting                                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│    ❌ internal/database/internal/pool/internal/connection/                │
│    ✅ internal/database/ (keep it flat)                                   │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ MISTAKE 4: Importing internal from outside the module                     │
├─────────────────────────────────────────────────────────────────────────────┤
│    ❌ otherproject/                                                       │
│        └── main.go: import "myproject/internal/auth"  // COMPILER ERROR! │
│    ✅ Use pkg/ for code that should be shared                             │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ MISTAKE 5: Not using internal for code that should be private            │
├─────────────────────────────────────────────────────────────────────────────┤
│    ❌ pkg/auth/ (external projects can import)                            │
│    ✅ internal/auth/ (truly private)                                      │
└─────────────────────────────────────────────────────────────────────────────┘
`)
}

// ============================================================================
// بخش 11: Real-World Example Structure
// ============================================================================

func realWorldExample() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🏗️ REAL-WORLD EXAMPLE")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ EXAMPLE: e-commerce platform                                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   github.com/mycompany/ecommerce/                                         │
│   ├── go.mod                                                              │
│   ├── internal/                     # Private to ecommerce module         │
│   │   ├── auth/                                                           │
│   │   │   ├── jwt.go                # JWT generation/validation           │
│   │   │   ├── middleware.go         # Auth middleware                     │
│   │   │   └── internal/                                                   │
│   │   │       └── crypto.go         # Internal crypto (private)           │
│   │   ├── database/                                                        │
│   │   │   ├── postgres.go           # PostgreSQL connection               │
│   │   │   ├── migrations/           # SQL migrations                      │
│   │   │   └── repository/           # Base repository                     │
│   │   ├── cache/                                                          │
│   │   │   ├── redis.go              # Redis connection                    │
│   │   │   └── cache.go              # Cache interfaces                    │
│   │   ├── queue/                                                          │
│   │   │   └── rabbitmq.go           # Message queue                       │
│   │   ├── logger/                                                         │
│   │   │   └── zap.go                # Logging setup                       │
│   │   ├── payment/                                                        │
│   │   │   ├── gateway.go            # Payment gateway integration         │
│   │   │   └── webhook.go            # Payment webhook handler             │
│   │   └── config/                                                         │
│   │       └── config.go             # Configuration loading               │
│   ├── pkg/                          # Public (other projects can use)    │
│   │   ├── api/                                                           │
│   │   │   └── client/               # Public API client                   │
│   │   │       └── client.go                                               │
│   │   ├── models/                   # Shared data models                  │
│   │   │   ├── user.go                                                     │
│   │   │   ├── product.go                                                  │
│   │   │   └── order.go                                                    │
│   │   └── errors/                  # Public error types                   │
│   │       └── errors.go                                                   │
│   ├── cmd/                                                                 │
│   │   ├── api-server/              # HTTP API server                      │
│   │   │   ├── main.go                                                     │
│   │   │   └── wire.go              # Dependency injection                 │
│   │   ├── worker/                   # Background worker                   │
│   │   │   └── main.go                                                     │
│   │   └── migrate/                  # Migration tool                      │
│   │       └── main.go                                                     │
│   ├── api/                          # API definitions                     │
│   │   ├── openapi/                                                        │
│   │   │   └── swagger.yaml                                                │
│   │   └── grpc/                                                           │
│   │       └── proto/                                                      │
│   ├── web/                          # Web assets (if any)                 │
│   │   ├── static/                                                         │
│   │   └── templates/                                                      │
│   ├── scripts/                      # Build/deploy scripts                │
│   │   ├── build.sh                                                        │
│   │   └── deploy.sh                                                       │
│   ├── test/                         # Integration tests                   │
│   │   └── integration_test.go                                             │
│   ├── go.mod                                                              │
│   └── go.sum                                                              │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
`)
}

// ============================================================================
// بخش 12: Testing with Internal Packages
// ============================================================================

func testingConsiderations() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🧪 TESTING WITH INTERNAL PACKAGES")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ WHITE-BOX TESTS                                                           │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   // tests can be inside the same package                                 │
│   // internal/auth/auth_test.go                                           │
│   package auth                                                             │
│                                                                             │
│   func TestJWTGenerator(t *testing.T) {                                    │
│       // can access unexported fields                                      │
│       gen := &JWTGenerator{secretKey: "test"}                             │
│       // test...                                                          │
│   }                                                                        │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ BLACK-BOX TESTS                                                           │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   // tests in external package                                            │
│   // internal/auth/auth_external_test.go                                  │
│   package auth_test                                                        │
│                                                                             │
│   import "myproject/internal/auth"                                        │
│                                                                             │
│   func TestAuth(t *testing.T) {                                            │
│       // can only access exported symbols                                  │
│       gen := auth.NewJWTGenerator("test")                                 │
│       // test...                                                          │
│   }                                                                        │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ INTEGRATION TESTS                                                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   // test/integration_test.go                                             │
│   package test                                                             │
│                                                                             │
│   import (                                                                │
│       "myproject/internal/database"                                       │
│       "myproject/internal/auth"                                           │
│   )                                                                        │
│                                                                             │
│   func TestIntegration(t *testing.T) {                                     │
│       // can use internal packages in tests                               │
│       db := database.NewConnection(...)                                   │
│       auth := auth.NewJWTGenerator(...)                                   │
│       // test integration                                                 │
│   }                                                                        │
│                                                                             │
│   Note: Integration tests should be in the same module                    │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
`)
}

// ============================================================================
// بخش 13: Summary
// ============================================================================

func summary() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📚 INTERNAL PACKAGE SUMMARY")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ KEY POINTS                                                                │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   • internal packages are enforced by Go compiler                         │
│   • Only accessible within the same module                                 │
│   • Perfect for implementation details                                     │
│   • Helps create clean public APIs                                         │
│   • No versioning guarantees for internal code                            │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ DECISION TREE                                                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   Is this code needed outside the module?                                 │
│     │                                                                      │
│     ├── YES ──► Should it be stable?                                      │
│     │           ├── YES ──► pkg/ (public, versioned)                      │
│     │           └── NO  ──► root or cmd/ (internal to module)            │
│     │                                                                      │
│     └── NO  ──► internal/ (truly private)                                 │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

💡 GOLDEN RULES:

   1. Use internal for code that should NOT be imported by other modules
   2. Use pkg/ for code that CAN be imported by other modules
   3. Keep internal packages focused and not too deep
   4. Document why something is internal (if not obvious)
   5. Internal code can change without notice - no semver guarantees
   6. Tests can (and should) import internal packages
   7. Avoid circular dependencies within internal
   8. Don't put everything in internal - only truly private code
   9. Use nested internal sparingly (only when necessary)
   10. Remember: internal is enforced by compiler, not just convention
`)
}

// ============================================================================
// بخش 14: Main
// ============================================================================

func main() {
	fmt.Println(stringsRepeat("=", 80))
	fmt.Println("🎯 INTERNAL PACKAGE IN GO")
	fmt.Println("Module-Level Encapsulation")
	fmt.Println(stringsRepeat("=", 80))

	compareInternalAndPkg()
	bestPractices()
	commonMistakes()
	realWorldExample()
	testingConsiderations()
	summary()

	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🎯 INTERNAL PACKAGE - COMPLETE")
	fmt.Println("Keep your module's implementation details private!")
	fmt.Println(stringsRepeat("=", 80))
}

// Helper function
func stringsRepeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}
