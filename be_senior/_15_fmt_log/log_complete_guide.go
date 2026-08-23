// ============================================================================
// FILE: log_complete_guide.go
// TITLE: راهنمای کامل پکیج log در Go - لاگ‌گیری حرفه‌ای
// HOW TO RUN: go run log_complete_guide.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - پکیج log چیست و چرا مهم است؟
// ============================================================================
//
// پکیج log امکانات کامل برای:
// 1. لاگ‌گیری در برنامه (Logging)
// 2. سطوح مختلف لاگ (Debug, Info, Warning, Error)
// 3. نوشتن در فایل و مقصدهای مختلف
// 4. فرمت‌سازی لاگ‌ها (تاریخ، زمان، فایل، خط)
// 5. سفارشی‌سازی کامل Logger
//
// قانون طلایی:
// "از log.Print برای لاگ‌های ساده، از log.Fatal برای خطاهای غیرقابل بازیابی،
//  از log.Panic برای خطاهایی که باید stack trace داشته باشند استفاده کن.
//  همیشه لاگ‌های سطح بالا را به فایل بنویس، نه فقط stdout."
// ============================================================================

package __internal_packages

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// ============================================================================
// بخش 1: لاگ‌گیری پایه - توابع استاندارد
// ============================================================================

func demonstrateBasicLogging() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📝 BASIC LOGGING - Standard Log Functions")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 1.1 توابع پایه لاگ
	// ============================================
	fmt.Println("\n--- 1.1 Basic Log Functions ---")

	// Print: مثل fmt.Print ولی با timestamp
	log.Print("This is a log message")
	log.Println("This is a log message with newline")
	log.Printf("This is a formatted log: %s=%d", "value", 42)

	// ============================================
	// 1.2 Fatal - خروج از برنامه بعد از لاگ
	// ============================================
	fmt.Println("\n--- 1.2 Fatal Functions (exits program) ---")
	fmt.Println("  (Skipping actual Fatal to not exit the demo)")
	// log.Fatal("This will exit the program")
	// log.Fatalf("Fatal error: %v", err)
	// log.Fatalln("Fatal error with newline")

	// ============================================
	// 1.3 Panic - panic بعد از لاگ
	// ============================================
	fmt.Println("\n--- 1.3 Panic Functions (panics) ---")
	fmt.Println("  (Skipping actual Panic to not crash the demo)")
	// log.Panic("This will panic")
	// log.Panicf("Panic: %v", err)
	// log.Panicln("Panic with newline")
}

// ============================================================================
// بخش 2: سفارشی‌سازی Logger
// ============================================================================

func demonstrateCustomLogger() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🔧 CUSTOM LOGGER - Flags and Configuration")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 2.1 فلگ‌های لاگ
	// ============================================
	fmt.Println("\n--- 2.1 Log Flags ---")

	flags := []struct {
		name string
		flag int
	}{
		{"date (2009/01/23)", log.Ldate},
		{"time (01:23:23)", log.Ltime},
		{"microseconds (01:23:23.123123)", log.Lmicroseconds},
		{"long file (a/b/c/d.go:23)", log.Llongfile},
		{"short file (d.go:23)", log.Lshortfile},
		{"UTC", log.LUTC},
		{"msgprefix", log.Lmsgprefix},
		{"std flags (date+time)", log.LstdFlags},
	}

	for _, f := range flags {
		logger := log.New(os.Stdout, "FLAG: ", f.flag)
		fmt.Printf("  %s:\n", f.name)
		logger.Print("sample message")
	}

	// ============================================
	// 2.2 ترکیب فلگ‌ها
	// ============================================
	fmt.Println("\n--- 2.2 Combining Flags ---")

	// ترکیب چند فلگ
	customFlags := log.Ldate | log.Ltime | log.Lshortfile
	logger := log.New(os.Stdout, "[CUSTOM] ", customFlags)
	logger.Print("This has date, time, and file info")

	// ============================================
	// 2.3 تغییر فلگ‌های default logger
	// ============================================
	fmt.Println("\n--- 2.3 Changing Default Logger Flags ---")

	oldFlags := log.Flags()
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
	log.Print("Default logger with new flags")
	log.SetFlags(oldFlags) // بازگردانی

	// ============================================
	// 2.4 تغییر Prefix
	// ============================================
	fmt.Println("\n--- 2.4 Changing Prefix ---")

	oldPrefix := log.Prefix()
	log.SetPrefix("[MAIN] ")
	log.Print("Message with new prefix")
	log.SetPrefix(oldPrefix)
}

// ============================================================================
// بخش 3: نوشتن لاگ در فایل
// ============================================================================

func demonstrateFileLogging() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📁 FILE LOGGING - Writing Logs to Files")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 3.1 لاگ در فایل ساده
	// ============================================
	fmt.Println("\n--- 3.1 Simple File Logging ---")

	// ایجاد فایل لاگ
	logFile, err := os.Create("app.log")
	if err != nil {
		log.Printf("Error creating log file: %v", err)
		return
	}
	defer logFile.Close()
	defer os.Remove("app.log")

	// ایجاد logger که در فایل می‌نویسد
	fileLogger := log.New(logFile, "[FILE] ", log.LstdFlags)

	fileLogger.Println("First log entry")
	fileLogger.Printf("User %s logged in", "ali")
	fileLogger.Println("Application started")

	// نمایش محتوای فایل
	content, _ := os.ReadFile("app.log")
	fmt.Printf("  Log file content:\n%s", string(content))

	// ============================================
	// 3.2 لاگ همزمان در فایل و کنسول
	// ============================================
	fmt.Println("\n--- 3.2 Multi-Writer Logging ---")

	logFile2, _ := os.Create("multi.log")
	defer logFile2.Close()
	defer os.Remove("multi.log")

	// نوشتن همزمان در فایل و stdout
	multiWriter := io.MultiWriter(os.Stdout, logFile2)
	multiLogger := log.New(multiWriter, "[MULTI] ", log.LstdFlags)

	multiLogger.Println("This goes to both console and file")

	// ============================================
	// 3.3 لاگ با چرخش فایل (Rotation)
	// ============================================
	fmt.Println("\n--- 3.3 Log Rotation Example ---")

	type RotatingLogger struct {
		file       *os.File
		baseName   string
		maxSize    int64
		currentSize int64
		logger     *log.Logger
	}

	func NewRotatingLogger(baseName string, maxSize int64) (*RotatingLogger, error) {
		file, err := os.OpenFile(baseName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return nil, err
		}

		info, _ := file.Stat()
		rl := &RotatingLogger{
			file:       file,
			baseName:   baseName,
			maxSize:    maxSize,
			currentSize: info.Size(),
		}
		rl.logger = log.New(file, "", log.LstdFlags)
		return rl, nil
	}

	func (rl *RotatingLogger) Write(p []byte) (n int, err error) {
		if rl.currentSize+int64(len(p)) > rl.maxSize {
			rl.rotate()
		}
		n, err = rl.file.Write(p)
		rl.currentSize += int64(n)
		return n, err
	}

	func (rl *RotatingLogger) rotate() {
		rl.file.Close()
		backupName := fmt.Sprintf("%s.%d", rl.baseName, time.Now().Unix())
		os.Rename(rl.baseName, backupName)
		rl.file, _ = os.Create(rl.baseName)
		rl.currentSize = 0
		rl.logger = log.New(rl.file, "", log.LstdFlags)
	}

	func (rl *RotatingLogger) Println(v ...interface{}) {
rl.logger.Println(v...)
}

// استفاده
rotating, _ := NewRotatingLogger("rotate.log", 1024)
defer os.Remove("rotate.log")
defer func() {
	// حذف فایل‌های backup
	os.Remove("rotate.log")
}()

rotating.Println("This will rotate when size exceeds 1KB")
fmt.Println("  (Rotation example - would rotate on size limit)")
}

// ============================================================================
// بخش 4: سطوح لاگ (Log Levels)
// ============================================================================

// LogLevel نوع سطح لاگ
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARNING
	ERROR
	FATAL
)

func (l LogLevel) String() string {
	return [...]string{"DEBUG", "INFO", "WARNING", "ERROR", "FATAL"}[l]
}

// LeveledLogger لاگر با سطوح مختلف
type LeveledLogger struct {
	*log.Logger
	level LogLevel
}

// NewLeveledLogger ایجاد لاگر با سطح
func NewLeveledLogger(out io.Writer, prefix string, flag int, level LogLevel) *LeveledLogger {
	return &LeveledLogger{
		Logger: log.New(out, prefix, flag),
		level:  level,
	}
}

func (l *LeveledLogger) Debug(v ...interface{}) {
	if l.level <= DEBUG {
		l.Logger.SetPrefix("🔍 [DEBUG] ")
		l.Logger.Println(v...)
	}
}

func (l *LeveledLogger) Info(v ...interface{}) {
	if l.level <= INFO {
		l.Logger.SetPrefix("ℹ️  [INFO] ")
		l.Logger.Println(v...)
	}
}

func (l *LeveledLogger) Warning(v ...interface{}) {
	if l.level <= WARNING {
		l.Logger.SetPrefix("⚠️  [WARNING] ")
		l.Logger.Println(v...)
	}
}

func (l *LeveledLogger) Error(v ...interface{}) {
	if l.level <= ERROR {
		l.Logger.SetPrefix("❌ [ERROR] ")
		l.Logger.Println(v...)
	}
}

func (l *LeveledLogger) Fatal(v ...interface{}) {
	if l.level <= FATAL {
		l.Logger.SetPrefix("💀 [FATAL] ")
		l.Logger.Fatalln(v...)
	}
}

// فرمت‌دار
func (l *LeveledLogger) Debugf(format string, v ...interface{}) {
	if l.level <= DEBUG {
		l.Logger.SetPrefix("🔍 [DEBUG] ")
		l.Logger.Printf(format, v...)
	}
}

func (l *LeveledLogger) Infof(format string, v ...interface{}) {
	if l.level <= INFO {
		l.Logger.SetPrefix("ℹ️  [INFO] ")
		l.Logger.Printf(format, v...)
	}
}

func (l *LeveledLogger) Warningf(format string, v ...interface{}) {
	if l.level <= WARNING {
		l.Logger.SetPrefix("⚠️  [WARNING] ")
		l.Logger.Printf(format, v...)
	}
}

func (l *LeveledLogger) Errorf(format string, v ...interface{}) {
	if l.level <= ERROR {
		l.Logger.SetPrefix("❌ [ERROR] ")
		l.Logger.Printf(format, v...)
	}
}

func demonstrateLogLevels() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🎚️ LOG LEVELS - Debug, Info, Warning, Error")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 4.1 استفاده از سطوح مختلف
	// ============================================
	fmt.Println("\n--- 4.1 Different Log Levels ---")

	var buf bytes.Buffer
	logger := NewLeveledLogger(&buf, "", 0, DEBUG)

	logger.Debug("This is a debug message")
	logger.Info("This is an info message")
	logger.Warning("This is a warning message")
	logger.Error("This is an error message")

	fmt.Printf("  All levels output:\n%s", buf.String())

	// ============================================
	// 4.2 فیلتر بر اساس سطح
	// ============================================
	fmt.Println("\n--- 4.2 Filtering by Level ---")

	buf.Reset()

	// فقط WARNING و بالاتر
	warnLogger := NewLeveledLogger(&buf, "", 0, WARNING)
	warnLogger.Debug("Debug message (ignored)")
	warnLogger.Info("Info message (ignored)")
	warnLogger.Warning("Warning message (shown)")
	warnLogger.Error("Error message (shown)")

	fmt.Printf("  Only WARNING and above:\n%s", buf.String())
}

// ============================================================================
// بخش 5: لاگ ساختاریافته (Structured Logging)
// ============================================================================

// StructuredLogger لاگر با فیلدهای کلید-مقدار
type StructuredLogger struct {
	*log.Logger
	fields map[string]interface{}
}

// NewStructuredLogger ایجاد لاگر ساختاریافته
func NewStructuredLogger(out io.Writer, prefix string, flag int) *StructuredLogger {
	return &StructuredLogger{
		Logger: log.New(out, prefix, flag),
		fields: make(map[string]interface{}),
	}
}

// WithField افزودن فیلد
func (l *StructuredLogger) WithField(key string, value interface{}) *StructuredLogger {
	newLogger := &StructuredLogger{
		Logger: l.Logger,
		fields: make(map[string]interface{}),
	}
	for k, v := range l.fields {
		newLogger.fields[k] = v
	}
	newLogger.fields[key] = value
	return newLogger
}

// WithFields افزودن چند فیلد
func (l *StructuredLogger) WithFields(fields map[string]interface{}) *StructuredLogger {
	newLogger := &StructuredLogger{
		Logger: l.Logger,
		fields: make(map[string]interface{}),
	}
	for k, v := range l.fields {
		newLogger.fields[k] = v
	}
	for k, v := range fields {
		newLogger.fields[k] = v
	}
	return newLogger
}

func (l *StructuredLogger) formatFields() string {
	if len(l.fields) == 0 {
		return ""
	}
	var parts []string
	for k, v := range l.fields {
		parts = append(parts, fmt.Sprintf("%s=%v", k, v))
	}
	return " {" + strings.Join(parts, ", ") + "}"
}

func (l *StructuredLogger) Info(v ...interface{}) {
	msg := fmt.Sprint(v...)
	l.Logger.Print(msg + l.formatFields())
}

func (l *StructuredLogger) Infof(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	l.Logger.Print(msg + l.formatFields())
}

func (l *StructuredLogger) Error(v ...interface{}) {
	msg := fmt.Sprint(v...)
	l.Logger.Print(msg + l.formatFields())
}

func (l *StructuredLogger) Errorf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	l.Logger.Print(msg + l.formatFields())
}

func demonstrateStructuredLogging() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🏗️ STRUCTURED LOGGING - Key-Value Pairs")
	fmt.Println(stringsRepeat("=", 80))

	var buf bytes.Buffer
	logger := NewStructuredLogger(&buf, "", log.Ltime)

	// لاگ با فیلد
	requestLogger := logger.WithField("request_id", "req-123")
	requestLogger.Info("Request received")

	userLogger := requestLogger.WithFields(map[string]interface{}{
		"user_id": 1001,
		"action":  "login",
		"ip":      "192.168.1.1",
	})
	userLogger.Info("User action")

	errorLogger := userLogger.WithField("error", "invalid credentials")
	errorLogger.Error("Authentication failed")

	fmt.Printf("  Structured logs:\n%s", buf.String())
}

// ============================================================================
// بخش 6: لاگ با Context (Request Tracing)
// ============================================================================

type ContextLogger struct {
	*log.Logger
	requestID string
	userID    string
}

func NewContextLogger(out io.Writer, prefix string, flag int) *ContextLogger {
	return &ContextLogger{
		Logger: log.New(out, prefix, flag),
	}
}

func (l *ContextLogger) WithRequestID(id string) *ContextLogger {
	return &ContextLogger{
		Logger:    l.Logger,
		requestID: id,
		userID:    l.userID,
	}
}

func (l *ContextLogger) WithUserID(id string) *ContextLogger {
	return &ContextLogger{
		Logger:    l.Logger,
		requestID: l.requestID,
		userID:    id,
	}
}

func (l *ContextLogger) formatContext() string {
	var parts []string
	if l.requestID != "" {
		parts = append(parts, fmt.Sprintf("req=%s", l.requestID))
	}
	if l.userID != "" {
		parts = append(parts, fmt.Sprintf("user=%s", l.userID))
	}
	if len(parts) > 0 {
		return " [" + strings.Join(parts, " ") + "]"
	}
	return ""
}

func (l *ContextLogger) Println(v ...interface{}) {
	msg := fmt.Sprint(v...)
	l.Logger.Print(msg + l.formatContext())
}

func (l *ContextLogger) Printf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	l.Logger.Print(msg + l.formatContext())
}

func demonstrateContextLogging() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🔗 CONTEXT LOGGING - Request Tracing")
	fmt.Println(stringsRepeat("=", 80))

	var buf bytes.Buffer
	logger := NewContextLogger(&buf, "", log.Ltime)

	// شبیه‌سازی یک درخواست HTTP
	requestLogger := logger.WithRequestID("req-2024-001")
	requestLogger.Println("Request started")

	userLogger := requestLogger.WithUserID("user-ali")
	userLogger.Println("Authenticating user")

	// درخواست دیگر
	requestLogger2 := logger.WithRequestID("req-2024-002")
	requestLogger2.Println("Another request")

	fmt.Printf("  Context logs:\n%s", buf.String())
}

// ============================================================================
// بخش 7: لاگ با رنگ (Colorized Logging)
// ============================================================================

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
)

type ColorLogger struct {
	*log.Logger
}

func NewColorLogger(out io.Writer, prefix string, flag int) *ColorLogger {
	return &ColorLogger{
		Logger: log.New(out, prefix, flag),
	}
}

func (l *ColorLogger) Debug(v ...interface{}) {
	l.SetPrefix(colorCyan + "[DEBUG] " + colorReset)
	l.Logger.Println(v...)
}

func (l *ColorLogger) Info(v ...interface{}) {
	l.SetPrefix(colorGreen + "[INFO] " + colorReset)
	l.Logger.Println(v...)
}

func (l *ColorLogger) Warning(v ...interface{}) {
	l.SetPrefix(colorYellow + "[WARNING] " + colorReset)
	l.Logger.Println(v...)
}

func (l *ColorLogger) Error(v ...interface{}) {
	l.SetPrefix(colorRed + "[ERROR] " + colorReset)
	l.Logger.Println(v...)
}

func demonstrateColorLogging() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🎨 COLORIZED LOGGING - Terminal Colors")
	fmt.Println(stringsRepeat("=", 80))

	var buf bytes.Buffer

	// برای نمایش، به جای ترمینال از buffer استفاده می‌کنیم
	// در عمل، رنگ‌ها در ترمینال دیده می‌شوند
	fmt.Println("  (Colors visible in terminal only)")

	logger := NewColorLogger(os.Stdout, "", 0)
	logger.Debug("Debug message (cyan)")
	logger.Info("Info message (green)")
	logger.Warning("Warning message (yellow)")
	logger.Error("Error message (red)")
}

// ============================================================================
// بخش 8: لاگ با اطلاعات تماس (Caller Info)
// ============================================================================

func demonstrateCallerInfo() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📞 CALLER INFO - File, Line, Function")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 8.1 استفاده از فلگ Lshortfile
	// ============================================
	fmt.Println("\n--- 8.1 Using Lshortfile Flag ---")

	callerLogger := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	callerLogger.Println("This shows file and line number")

	// ============================================
	// 8.2 گرفتن دستی caller info
	// ============================================
	fmt.Println("\n--- 8.2 Manual Caller Info ---")

	logWithCaller := func(level, msg string) {
		_, file, line, ok := runtime.Caller(1)
		if !ok {
			file = "unknown"
			line = 0
		}
		file = filepath.Base(file)
		fmt.Printf("[%s] %s:%d: %s\n", level, file, line, msg)
	}

	logWithCaller("INFO", "This shows caller info")

	// ============================================
	// 8.3 Logger با caller info خودکار
	// ============================================
	fmt.Println("\n--- 8.3 Logger with Auto Caller Info ---")

	type CallerLogger struct {
		*log.Logger
	}

	func NewCallerLogger(out io.Writer, prefix string, flag int) *CallerLogger {
		return &CallerLogger{
		Logger: log.New(out, prefix, flag|log.Lshortfile),
	}
	}

	func (l *CallerLogger) Info(v ...interface{}) {
l.SetPrefix("[INFO] ")
l.Logger.Println(v...)
}

var buf bytes.Buffer
cl := NewCallerLogger(&buf, "", 0)
cl.Info("This message shows the caller")

fmt.Printf("  Output: %s", buf.String())
}

// ============================================================================
// بخش 9: کاربردهای عملی
// ============================================================================

func demonstratePracticalUses() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("💡 PRACTICAL LOGGING USE CASES")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 9.1 لاگ در برنامه وب
	// ============================================
	fmt.Println("\n--- 9.1 Web Application Logging ---")

	type HTTPLogEntry struct {
		Method   string
		Path     string
		Status   int
		Duration time.Duration
		ClientIP string
	}

	logHTTPRequest := func(entry HTTPLogEntry) {
		log.Printf("[HTTP] %s %s → %d (%v) from %s",
			entry.Method, entry.Path, entry.Status, entry.Duration, entry.ClientIP)
	}

	logHTTPRequest(HTTPLogEntry{
		Method:   "GET",
		Path:     "/api/users",
		Status:   200,
		Duration: 45 * time.Millisecond,
		ClientIP: "192.168.1.100",
	})

	// ============================================
	// 9.2 لاگ خطا با stack trace
	// ============================================
	fmt.Println("\n--- 9.2 Error Logging with Stack Trace ---")

	logError := func(err error) {
		buf := make([]byte, 4096)
		n := runtime.Stack(buf, false)
		log.Printf("ERROR: %v\nStack trace:\n%s", err, buf[:n])
	}

	logError(fmt.Errorf("database connection failed"))

	// ============================================
	// 9.3 لاگ با زمان‌سنجی
	// ============================================
	fmt.Println("\n--- 9.3 Performance Logging ---")

	logDuration := func(name string, fn func()) {
		start := time.Now()
		fn()
		duration := time.Since(start)
		log.Printf("[PERF] %s took %v", name, duration)
	}

	logDuration("sleep 50ms", func() {
		time.Sleep(50 * time.Millisecond)
	})

	// ============================================
	// 9.4 لاگ با JSON (برای ELK و سایر سیستم‌ها)
	// ============================================
	fmt.Println("\n--- 9.4 JSON Logging (for ELK/Loki) ---")

	type JSONLog struct {
		Timestamp string                 `json:"timestamp"`
		Level     string                 `json:"level"`
		Message   string                 `json:"message"`
		Fields    map[string]interface{} `json:"fields,omitempty"`
	}

	logJSON := func(level, message string, fields map[string]interface{}) {
		entry := JSONLog{
			Timestamp: time.Now().Format(time.RFC3339Nano),
			Level:     level,
			Message:   message,
			Fields:    fields,
		}
		// در عمل از json.Marshal استفاده می‌شود
		fmt.Printf("  JSON log: %+v\n", entry)
	}

	logJSON("info", "User logged in", map[string]interface{}{
		"user_id": 123,
		"ip":      "10.0.0.1",
	})
}

// ============================================================================
// بخش 10: لاگ در کتابخانه‌های شخص ثالث (مرور)
// ============================================================================

func demonstrateThirdPartyLogging() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📦 THIRD-PARTY LOGGING - Popular Libraries")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println("\nPopular logging libraries in Go ecosystem:")
	fmt.Println("  1. logrus - Structured logging with levels")
	fmt.Println("     go get github.com/sirupsen/logrus")
	fmt.Println()
	fmt.Println("  2. zap - High-performance structured logging")
	fmt.Println("     go get go.uber.org/zap")
	fmt.Println()
	fmt.Println("  3. zerolog - Zero-allocation JSON logging")
	fmt.Println("     go get github.com/rs/zerolog")
	fmt.Println()
	fmt.Println("  4. slog - Standard structured logging (Go 1.21+)")
	fmt.Println("     import log/slog")
	fmt.Println()
	fmt.Println("💡 Recommendation:")
	fmt.Println("  • Simple apps → standard log package")
	fmt.Println("  • Structured logs → slog (Go 1.21+) or logrus")
	fmt.Println("  • High performance → zap or zerolog")
}

// ============================================================================
// بخش 11: اشتباهات رایج
// ============================================================================

func demonstrateCommonMistakes() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("❌ COMMON MISTAKES WITH log PACKAGE")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println("\n❌ Mistake 1: Using log.Fatal without cleanup")
	fmt.Println("   log.Fatal(\"error\")  // exits immediately")
	fmt.Println("   ✅ Defer cleanup before Fatal")

	fmt.Println("\n❌ Mistake 2: Not using log levels")
	fmt.Println("   log.Print(\"debug info\")  // always shown")
	fmt.Println("   ✅ Implement levels or use leveled logger")

	fmt.Println("\n❌ Mistake 3: Logging sensitive information")
	fmt.Println("   log.Printf(\"Password: %s\", password)  // security risk")
	fmt.Println("   ✅ Redact sensitive data before logging")

	fmt.Println("\n❌ Mistake 4: Ignoring log output")
	fmt.Println("   log.SetOutput(io.Discard)  // throws away logs")
	fmt.Println("   ✅ Always have a log destination")

	fmt.Println("\n❌ Mistake 5: Not rotating log files")
	fmt.Println("   // logs grow indefinitely")
	fmt.Println("   ✅ Use log rotation (lumberjack, etc.)")

	fmt.Println("\n❌ Mistake 6: Logging in tight loops")
	fmt.Println("   for i := 0; i < 1000000; i++ { log.Print(i) }")
	fmt.Println("   ✅ Use sampling or higher log level")

	fmt.Println("\n❌ Mistake 7: Not including context in logs")
	fmt.Println("   log.Print(\"user action\")  // no context")
	fmt.Println("   ✅ Include request ID, user ID, etc.")
}

// ============================================================================
// بخش 12: جمع‌بندی و جدول مرجع
// ============================================================================

func main() {
	fmt.Println(stringsRepeat("=", 80))
	fmt.Println("🎯 COMPLETE log PACKAGE GUIDE IN GO")
	fmt.Println(stringsRepeat("=", 80))

	// بخش 1: Basic Logging
	demonstrateBasicLogging()

	// بخش 2: Custom Logger
	demonstrateCustomLogger()

	// بخش 3: File Logging
	demonstrateFileLogging()

	// بخش 4: Log Levels
	demonstrateLogLevels()

	// بخش 5: Structured Logging
	demonstrateStructuredLogging()

	// بخش 6: Context Logging
	demonstrateContextLogging()

	// بخش 7: Colorized Logging
	demonstrateColorLogging()

	// بخش 8: Caller Info
	demonstrateCallerInfo()

	// بخش 9: Practical Uses
	demonstratePracticalUses()

	// بخش 10: Third-party
	demonstrateThirdPartyLogging()

	// بخش 11: Common Mistakes
	demonstrateCommonMistakes()

	// بخش 12: Quick Reference
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📚 log PACKAGE QUICK REFERENCE")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ BASIC FUNCTIONS                                               │")
	fmt.Println("├─────────────────────────────────────────────────────────────────┤")
	fmt.Println("│ log.Print(v...)     - Log message                              │")
	fmt.Println("│ log.Println(v...)   - Log message with newline                 │")
	fmt.Println("│ log.Printf(f, v...) - Formatted log message                    │")
	fmt.Println("│ log.Fatal(v...)     - Log + os.Exit(1)                         │")
	fmt.Println("│ log.Fatalf(f, v...) - Formatted log + exit                     │")
	fmt.Println("│ log.Panic(v...)     - Log + panic()                            │")
	fmt.Println("│ log.Panicf(f, v...) - Formatted log + panic                    │")
	fmt.Println("└─────────────────────────────────────────────────────────────────┘")

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ LOGGER CONFIGURATION                                          │")
	fmt.Println("├─────────────────────────────────────────────────────────────────┤")
	fmt.Println("│ log.New(out, prefix, flag) - Create custom logger             │")
	fmt.Println("│ log.SetOutput(w)           - Set output destination           │")
	fmt.Println("│ log.SetPrefix(p)           - Set prefix                       │")
	fmt.Println("│ log.SetFlags(flag)         - Set flags                        │")
	fmt.Println("│ log.Flags()                - Get current flags                │")
	fmt.Println("│ log.Prefix()               - Get current prefix               │")
	fmt.Println("└─────────────────────────────────────────────────────────────────┘")

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ FLAGS                                                          │")
	fmt.Println("├─────────────────────────────────────────────────────────────────┤")
	fmt.Println("│ log.Ldate         - Date (2009/01/23)                          │")
	fmt.Println("│ log.Ltime         - Time (01:23:23)                            │")
	fmt.Println("│ log.Lmicroseconds - Time with microseconds (01:23:23.123123)   │")
	fmt.Println("│ log.Llongfile     - Full file name and line                    │")
	fmt.Println("│ log.Lshortfile    - Short file name and line                   │")
	fmt.Println("│ log.LUTC          - Use UTC for timestamps                     │")
	fmt.Println("│ log.Lmsgprefix    - Move prefix before message                 │")
	fmt.Println("│ log.LstdFlags     - Ldate | Ltime (default)                    │")
	fmt.Println("└─────────────────────────────────────────────────────────────────┘")

	fmt.Println("\n💡 GOLDEN RULES:")
	fmt.Println("  1. Use log.Fatal only for unrecoverable errors")
	fmt.Println("  2. Use log.Panic only when you need stack trace")
	fmt.Println("  3. Always close log files with defer")
	fmt.Println("  4. Use log levels for different environments")
	fmt.Println("  5. Never log sensitive information (passwords, tokens)")
	fmt.Println("  6. Use structured logging for production")
	fmt.Println("  7. Include request ID in logs for tracing")
	fmt.Println("  8. Rotate log files to avoid disk fill")
	fmt.Println("  9. Log errors with context (file, line, function)")
	fmt.Println("  10. Use third-party libs for complex logging needs")

	fmt.Println("\n🎯 BEST PRACTICES:")
	fmt.Println("  • Development: Human-readable logs with colors")
	fmt.Println("  • Production: Structured logs (JSON) for machine parsing")
	fmt.Println("  • Debug: Verbose logs with caller info")
	fmt.Println("  • Performance: Async logging, avoid allocations")
	fmt.Println("  • Security: Redact sensitive data, audit logs")
}

// ============================================================================
// بخش 13: توابع کمکی
// ============================================================================

func stringsRepeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}