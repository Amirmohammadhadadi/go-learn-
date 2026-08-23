// ============================================================================
// FILE: structured_logging_guide.go
// TITLE: راهنمای کامل Logging ساختاریافته در Go - Zap و Logrus
// HOW TO RUN: go run structured_logging_guide.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - Logging ساختاریافته چیست و چرا مهم است؟
// ============================================================================
//
// Logging ساختاریافته روشی برای ثبت لاگ‌ها به صورت key-value pairs یا JSON است.
//
// مزایا نسبت به لاگ سنتی (خطوط متنی):
// 1. قابلیت جستجو پیشرفته (مثلاً جستجو بر اساس `level=error`)
// 2. Parse کردن آسان توسط ابزارهایی مانند ELK, Loki, Splunk
// 3. امکان اضافه کردن metadata به هر لاگ (request_id, user_id, etc.)
// 4. نمایش بهتر در ابزارهای visualization
// 5. استانداردسازی بین سرویس‌های مختلف
//
// مقایسه Zap و Logrus:
// ┌─────────────────────────────────────────────────────────────────────────────┐
// │ ویژگی              │ Zap                      │ Logrus                    │
// ├────────────────────┼──────────────────────────┼───────────────────────────┤
// │ Performance        │ بسیار بالا (zero alloc)  │ متوسط                     │
// │ Ease of use        │ متوسط (API پیچیده‌تر)    │ بسیار آسان                │
// │ Structured fields  │ ✅ (سریع)                │ ✅ (با fields)            │
// │ Hooks              │ محدود                     │ ✅ (بسیار غنی)            │
// │ Output formats     │ JSON, Console            │ JSON, Text                │
// │ Level parsing      │ ✅                       │ ✅                        │
// │ Sampling           │ ✅                       │ ❌                       │
// │ Stack traces       │ ✅                       │ ✅                        │
// │ Use case           │ High-performance apps    │ General purpose          │
// └─────────────────────────────────────────────────────────────────────────────┘
//
// قانون طلایی:
// "برای برنامه‌های با performance بالا از Zap استفاده کن.
//  برای سادگی و قابلیت توسعه از Logrus استفاده کن.
//  همیشه از context-aware logging استفاده کن (request_id, user_id).
//  لاگ‌های حاوی secrets را هرگز ننویس."
// ============================================================================

package __structured_logging

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ============================================================================
// بخش 1: Zap Logger - تنظیمات و راه‌اندازی اولیه
// ============================================================================

// NewZapLogger ایجاد Zap Logger با تنظیمات بهینه
func NewZapLogger(level string, isProduction bool) (*zap.Logger, error) {
	var config zap.Config

	if isProduction {
		// تنظیمات پروداکشن (JSON, 高性能)
		config = zap.NewProductionConfig()
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		// تنظیمات توسعه (human-readable, colorful)
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// تنظیم سطح لاگ
	var zapLevel zapcore.Level
	switch level {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "info":
		zapLevel = zapcore.InfoLevel
	case "warn":
		zapLevel = zapcore.WarnLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	default:
		zapLevel = zapcore.InfoLevel
	}
	config.Level = zap.NewAtomicLevelAt(zapLevel)

	// ایجاد logger
	logger, err := config.Build(
		zap.AddCaller(),                       // افزودن caller info
		zap.AddStacktrace(zapcore.ErrorLevel), // افزودن stacktrace برای error
	)
	if err != nil {
		return nil, err
	}

	return logger, nil
}

// NewZapLoggerWithOptions ایجاد Zap Logger با تنظیمات سفارشی
func NewZapLoggerWithOptions() *zap.Logger {
	// تنظیم encoder
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// تنظیم core
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(zapcore.Lock(zapcore.AddSync(consoleWriter{}))),
		zapcore.DebugLevel,
	)

	// ایجاد logger با sampler (کاهش لاگ‌های تکراری)
	logger := zap.New(core,
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	return logger
}

// consoleWriter برای نوشتن در console
type consoleWriter struct{}

func (w consoleWriter) Write(p []byte) (n int, err error) {
	return fmt.Fprint(zapcore.Lock(os.Stdout), string(p))
}

func (w consoleWriter) Sync() error {
	return nil
}

// ============================================================================
// بخش 2: Zap Logger - مثال‌های عملی
// ============================================================================

func demonstrateZapExample() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("⚡ ZAP LOGGER EXAMPLES")
	fmt.Println(strings.Repeat("=", 80))

	// ایجاد logger
	logger, err := NewZapLogger("debug", false)
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	// لاگ پایه
	logger.Debug("Debug message",
		zap.String("module", "auth"),
		zap.Int("attempt", 1),
	)

	logger.Info("Server started",
		zap.String("host", "localhost"),
		zap.Int("port", 8080),
		zap.Duration("timeout", 30*time.Second),
	)

	logger.Warn("High memory usage",
		zap.Float64("memory_usage_mb", 512.5),
		zap.Float64("threshold_mb", 256.0),
	)

	logger.Error("Database connection failed",
		zap.String("error", "connection timeout"),
		zap.String("db_host", "postgres.example.com"),
		zap.Int("retry_count", 3),
	)

	// لاگ با struct
	type User struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	user := User{ID: "123", Name: "Ali", Email: "ali@example.com"}
	logger.Info("User created",
		zap.Any("user", user),
		zap.String("request_id", "req-123"),
	)

	// لاگ با grouped fields
	logger.Info("HTTP request",
		zap.Namespace("request"),
		zap.String("method", "POST"),
		zap.String("path", "/api/users"),
		zap.Int("status", 201),
		zap.Namespace("client"),
		zap.String("ip", "192.168.1.1"),
		zap.String("user_agent", "curl/7.68.0"),
	)

	// لاگ با sugar (راحت‌تر اما کمی کندتر)
	sugar := logger.Sugar()
	sugar.Infow("Sugar logging example",
		"key1", "value1",
		"key2", 42,
		"key3", true,
	)
}

// ============================================================================
// بخش 3: Zap Logger با Context (برای Request Tracing)
// ============================================================================

type contextKey string

const RequestIDKey contextKey = "request_id"
const UserIDKey contextKey = "user_id"

// ZapContextLogger لاگر آگاه از context
type ZapContextLogger struct {
	logger *zap.Logger
}

func NewZapContextLogger(logger *zap.Logger) *ZapContextLogger {
	return &ZapContextLogger{logger: logger}
}

// getFieldsFromContext استخراج فیلدها از context
func (l *ZapContextLogger) getFieldsFromContext(ctx context.Context) []zap.Field {
	var fields []zap.Field

	if requestID := ctx.Value(RequestIDKey); requestID != nil {
		fields = append(fields, zap.String("request_id", requestID.(string)))
	}
	if userID := ctx.Value(UserIDKey); userID != nil {
		fields = append(fields, zap.String("user_id", userID.(string)))
	}

	return fields
}

func (l *ZapContextLogger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(l.getFieldsFromContext(ctx), fields...)
	l.logger.Debug(msg, allFields...)
}

func (l *ZapContextLogger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(l.getFieldsFromContext(ctx), fields...)
	l.logger.Info(msg, allFields...)
}

func (l *ZapContextLogger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(l.getFieldsFromContext(ctx), fields...)
	l.logger.Warn(msg, allFields...)
}

func (l *ZapContextLogger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(l.getFieldsFromContext(ctx), fields...)
	l.logger.Error(msg, allFields...)
}

func demonstrateZapContextLogger() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🔗 ZAP CONTEXT-AWARE LOGGING")
	fmt.Println(strings.Repeat("=", 80))

	logger, _ := NewZapLogger("info", false)
	ctxLogger := NewZapContextLogger(logger)

	ctx := context.Background()
	ctx = context.WithValue(ctx, RequestIDKey, "req-abc-123")
	ctx = context.WithValue(ctx, UserIDKey, "user-456")

	ctxLogger.Info(ctx, "Processing request",
		zap.String("action", "login"),
		zap.String("ip", "192.168.1.1"),
	)

	ctxLogger.Error(ctx, "Authentication failed",
		zap.String("reason", "invalid password"),
	)
}

// ============================================================================
// بخش 4: Logrus Logger - تنظیمات و راه‌اندازی اولیه
// ============================================================================

// NewLogrusLogger ایجاد Logrus Logger
func NewLogrusLogger(level string, format string) *logrus.Logger {
	logger := logrus.New()

	// تنظیم سطح لاگ
	switch level {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}

	// تنظیم فرمت
	if format == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339Nano,
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message",
			},
		})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: time.RFC3339,
		})
	}

	// افزودن caller info
	logger.SetReportCaller(true)

	return logger
}

// ============================================================================
// بخش 5: Logrus - مثال‌های عملی
// ============================================================================

func demonstrateLogrusExample() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📝 LOGRUS LOGGER EXAMPLES")
	fmt.Println(strings.Repeat("=", 80))

	logger := NewLogrusLogger("debug", "text")

	// لاگ پایه
	logger.Debug("Debug message")
	logger.Info("Info message")
	logger.Warn("Warning message")
	logger.Error("Error message")

	// لاگ با fields
	logger.WithFields(logrus.Fields{
		"module":  "auth",
		"user_id": "123",
		"attempt": 1,
	}).Info("Login attempt")

	// لاگ با struct
	type Order struct {
		ID     string  `json:"id"`
		Amount float64 `json:"amount"`
	}
	order := Order{ID: "ORD-001", Amount: 99.99}
	logger.WithField("order", order).Info("Order created")

	// لاگ با multiple fields
	logEntry := logger.WithFields(logrus.Fields{
		"request_id":  "req-456",
		"method":      "GET",
		"path":        "/api/users",
		"status":      200,
		"duration_ms": 45,
	})
	logEntry.Info("HTTP request completed")

	// لاگ با error
	err := fmt.Errorf("connection refused")
	logger.WithError(err).Error("Database connection failed")

	// لاگ با context values
	ctx := context.Background()
	ctxLogger := logger.WithContext(ctx)
	ctxLogger.WithFields(logrus.Fields{
		"request_id": "ctx-req-789",
	}).Info("Request processed")
}

// ============================================================================
// بخش 6: Logrus Hooks (برای افزودن قابلیت‌های اضافی)
// ============================================================================

// CustomHook هوک سفارشی برای Logrus
type CustomHook struct{}

func (hook *CustomHook) Fire(entry *logrus.Entry) error {
	// افزودن فیلدهای اضافی به همه لاگ‌ها
	entry.Data["app_name"] = "myapp"
	entry.Data["app_version"] = "1.0.0"
	return nil
}

func (hook *CustomHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// demonstrateLogrusHooks نمایش هوک‌های Logrus
func demonstrateLogrusHooks() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🪝 LOGRUS HOOKS")
	fmt.Println(strings.Repeat("=", 80))

	logger := NewLogrusLogger("info", "json")
	logger.AddHook(&CustomHook{})

	logger.Info("This log includes custom fields from hook")
}

// ============================================================================
// بخش 7: Logrus با File Rotation
// ============================================================================

func demonstrateLogrusFileRotation() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📁 LOGRUS WITH FILE ROTATION")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
// برای چرخش فایل لاگ، از کتابخانه‌های زیر استفاده کنید:
// go get github.com/natefinch/lumberjack

import (
    "github.com/natefinch/lumberjack"
    "github.com/sirupsen/logrus"
)

func setupFileRotation() {
    logger := logrus.New()
    
    // تنظیم lumberjack برای چرخش فایل
    logger.SetOutput(&lumberjack.Logger{
        Filename:   "/var/log/myapp/app.log",
        MaxSize:    100, // مگابایت
        MaxBackups: 3,
        MaxAge:     28,  // روز
        Compress:   true,
    })
    
    logger.Info("Log with rotation")
}
`)
}

// ============================================================================
// بخش 8: Zap vs Logrus Performance Comparison
// ============================================================================

func performanceComparison() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("⚡ ZAP vs LOGRUS PERFORMANCE")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ OPERATION              │ ZAP (ns/op)        │ LOGRUS (ns/op)     │ SPEEDUP  │
├────────────────────────┼────────────────────┼────────────────────┼──────────┤
│ Print with fields      │ ~200-300           │ ~800-1000          │ 3-4x     │
│ Print without fields   │ ~100-150           │ ~300-400           │ 2-3x     │
│ With context fields    │ ~250-350           │ ~1000-1200         │ 3-4x     │
│ Sampling enabled       │ ~50-100            │ N/A                │ N/A      │
└────────────────────────┴────────────────────┴────────────────────┴──────────┘

📊 BENCHMARK CODE:

func BenchmarkZap(b *testing.B) {
    logger, _ := zap.NewProduction()
    defer logger.Sync()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        logger.Info("test", zap.String("key", "value"), zap.Int("count", i))
    }
}

func BenchmarkLogrus(b *testing.B) {
    logger := logrus.New()
    logger.SetLevel(logrus.InfoLevel)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        logger.WithFields(logrus.Fields{
            "key": "value",
            "count": i,
        }).Info("test")
    }
}

💡 CONCLUSION:
   • For high-performance systems: Use Zap
   • For simplicity and features: Use Logrus
   • For libraries: Use interface (like logr)
`)
}

// ============================================================================
// بخش 9: Custom Logger Interface (برای قابلیت تعویض)
// ============================================================================

// Logger اینترفیس عمومی (برای استفاده در کل برنامه)
type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	WithField(key string, value interface{}) Logger
	WithFields(fields map[string]interface{}) Logger
}

// Field ساختار فیلد
type Field struct {
	Key   string
	Value interface{}
}

// ZapAdapter آداپتور Zap برای اینترفیس Logger
type ZapAdapter struct {
	logger *zap.Logger
	sugar  *zap.SugaredLogger
}

func NewZapAdapter(logger *zap.Logger) *ZapAdapter {
	return &ZapAdapter{
		logger: logger,
		sugar:  logger.Sugar(),
	}
}

func (z *ZapAdapter) fieldsToZap(fields []Field) []zap.Field {
	zapFields := make([]zap.Field, len(fields))
	for i, f := range fields {
		zapFields[i] = zap.Any(f.Key, f.Value)
	}
	return zapFields
}

func (z *ZapAdapter) Debug(msg string, fields ...Field) {
	if len(fields) > 0 {
		z.logger.Debug(msg, z.fieldsToZap(fields)...)
	} else {
		z.logger.Debug(msg)
	}
}

func (z *ZapAdapter) Info(msg string, fields ...Field) {
	if len(fields) > 0 {
		z.logger.Info(msg, z.fieldsToZap(fields)...)
	} else {
		z.logger.Info(msg)
	}
}

func (z *ZapAdapter) Warn(msg string, fields ...Field) {
	if len(fields) > 0 {
		z.logger.Warn(msg, z.fieldsToZap(fields)...)
	} else {
		z.logger.Warn(msg)
	}
}

func (z *ZapAdapter) Error(msg string, fields ...Field) {
	if len(fields) > 0 {
		z.logger.Error(msg, z.fieldsToZap(fields)...)
	} else {
		z.logger.Error(msg)
	}
}

func (z *ZapAdapter) WithField(key string, value interface{}) Logger {
	return &ZapAdapter{
		logger: z.logger.With(zap.Any(key, value)),
		sugar:  z.sugar.With(key, value),
	}
}

func (z *ZapAdapter) WithFields(fields map[string]interface{}) Logger {
	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	return &ZapAdapter{
		logger: z.logger.With(zapFields...),
		sugar:  z.sugar.With(zapFields...),
	}
}

// LogrusAdapter آداپتور Logrus برای اینترفیس Logger
type LogrusAdapter struct {
	logger *logrus.Entry
}

func NewLogrusAdapter(logger *logrus.Logger) *LogrusAdapter {
	return &LogrusAdapter{
		logger: logrus.NewEntry(logger),
	}
}

func (l *LogrusAdapter) fieldsToLogrus(fields []Field) logrus.Fields {
	logrusFields := make(logrus.Fields, len(fields))
	for _, f := range fields {
		logrusFields[f.Key] = f.Value
	}
	return logrusFields
}

func (l *LogrusAdapter) Debug(msg string, fields ...Field) {
	if len(fields) > 0 {
		l.logger.WithFields(l.fieldsToLogrus(fields)).Debug(msg)
	} else {
		l.logger.Debug(msg)
	}
}

func (l *LogrusAdapter) Info(msg string, fields ...Field) {
	if len(fields) > 0 {
		l.logger.WithFields(l.fieldsToLogrus(fields)).Info(msg)
	} else {
		l.logger.Info(msg)
	}
}

func (l *LogrusAdapter) Warn(msg string, fields ...Field) {
	if len(fields) > 0 {
		l.logger.WithFields(l.fieldsToLogrus(fields)).Warn(msg)
	} else {
		l.logger.Warn(msg)
	}
}

func (l *LogrusAdapter) Error(msg string, fields ...Field) {
	if len(fields) > 0 {
		l.logger.WithFields(l.fieldsToLogrus(fields)).Error(msg)
	} else {
		l.logger.Error(msg)
	}
}

func (l *LogrusAdapter) WithField(key string, value interface{}) Logger {
	return &LogrusAdapter{
		logger: l.logger.WithField(key, value),
	}
}

func (l *LogrusAdapter) WithFields(fields map[string]interface{}) Logger {
	return &LogrusAdapter{
		logger: l.logger.WithFields(fields),
	}
}

// ============================================================================
// بخش 10: Best Practices
// ============================================================================

func bestPractices() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("💡 LOGGING BEST PRACTICES")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ 1. LOG LEVELS                                                             │
├─────────────────────────────────────────────────────────────────────────────┤
│    • DEBUG: Detailed info for debugging (only in development)             │
│    • INFO: General application events (startup, shutdown, major ops)      │
│    • WARN: Unexpected but recoverable situations                          │
│    • ERROR: Errors that need investigation                                │
│    • FATAL: Critical errors, application will exit                        │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 2. WHAT TO LOG                                                            │
├─────────────────────────────────────────────────────────────────────────────┤
│    ✅ Request ID (trace ID)                                               │
│    ✅ User ID (if authenticated)                                          │
│    ✅ Operation duration                                                  │
│    ✅ HTTP status codes                                                   │
│    ✅ Error stacks                                                        │
│    ✅ System metrics (memory, CPU)                                        │
│                                                                           │
│    ❌ Passwords, tokens, secrets                                          │
│    ❌ Credit card numbers                                                 │
│    ❌ Personal identifiable information (PII)                             │
│    ❌ Health data                                                         │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 3. PERFORMANCE                                                            │
├─────────────────────────────────────────────────────────────────────────────┤
│    • Use structured logging (JSON) for parsing                            │
│    • Sample high-frequency logs                                           │
│    • Use async logging for high throughput                                │
│    • Avoid logging in hot paths                                           │
│    • Use log levels effectively                                           │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 4. PRODUCTION SETUP                                                       │
├─────────────────────────────────────────────────────────────────────────────┤
│    • Use JSON format for log aggregation                                  │
│    • Set appropriate log level (INFO or WARN)                            │
│    • Implement log rotation                                               │
│    • Forward logs to central system (ELK, Loki)                          │
│    • Set up alerts on error patterns                                      │
└─────────────────────────────────────────────────────────────────────────────┘
`)
}

// ============================================================================
// بخش 11: Main
// ============================================================================

func main() {
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("🎯 STRUCTURED LOGGING IN GO")
	fmt.Println("Zap | Logrus")
	fmt.Println(strings.Repeat("=", 80))

	bestPractices()
	performanceComparison()

	demonstrateZapExample()
	demonstrateZapContextLogger()
	demonstrateLogrusExample()
	demonstrateLogrusHooks()
	demonstrateLogrusFileRotation()

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📊 QUICK REFERENCE")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ INSTALLATION                                                              │
├─────────────────────────────────────────────────────────────────────────────┤
│   # Zap                                                                   │
│   go get -u go.uber.org/zap                                               │
│                                                                           │
│   # Logrus                                                                │
│   go get -u github.com/sirupsen/logrus                                   │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ QUICK START - ZAP                                                         │
├─────────────────────────────────────────────────────────────────────────────┤
│   logger, _ := zap.NewProduction()                                        │
│   defer logger.Sync()                                                     │
│                                                                           │
│   logger.Info("Server started",                                           │
│       zap.String("host", "localhost"),                                    │
│       zap.Int("port", 8080))                                             │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ QUICK START - LOGRUS                                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│   logger := logrus.New()                                                  │
│   logger.SetFormatter(&logrus.JSONFormatter{})                            │
│                                                                           │
│   logger.WithFields(logrus.Fields{                                        │
│       "host": "localhost",                                                │
│       "port": 8080,                                                       │
│   }).Info("Server started")                                              │
└─────────────────────────────────────────────────────────────────────────────┘
`)

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🎯 STRUCTURED LOGGING - COMPLETE")
	fmt.Println("Choose the right logger for your Go application!")
	fmt.Println(strings.Repeat("=", 80))
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
var _ = os.Stdout
