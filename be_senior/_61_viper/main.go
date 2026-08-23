// ============================================================================
// FILE: viper_guide.go
// TITLE: راهنمای کامل Viper در Go - مدیریت پیکربندی قدرتمند
// HOW TO RUN: go run viper_guide.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - Viper چیست و چرا نیاز است؟
// ============================================================================
//
// Viper یک کتابخانه کامل برای مدیریت پیکربندی در Go است که توسط spf13 ساخته شده.
// ویژگی‌های کلیدی:
// 1. پشتیبانی از فرمت‌های مختلف: JSON, TOML, YAML, HCL, envfile, Java properties
// 2. خواندن از Environment Variables
// 3. خواندن از فایل‌های پیکربندی
// 4. خواندن از Command Line Flags
// 5. مقادیر پیش‌فرض
// 6. Watch کردن تغییرات فایل (hot reload)
// 7. ادغام چندین منبع
// 8. پشتیبانی از alias (بدون case-sensitive)
//
// قانون طلایی:
// "از Viper برای جداسازی تنظیمات از کد استفاده کن.
//  ترتیب اولویت: flag > env > config file > default
//  از Viper در init() یا در آغاز برنامه راه‌اندازی کن."
// ============================================================================

package __viper

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// ============================================================================
// بخش 1: نصب و راه‌اندازی
// ============================================================================

/*
نصب:
$ go get github.com/spf13/viper
*/

// ============================================================================
// بخش 2: مدل‌های پیکربندی (Configuration Structs)
// ============================================================================

// Config ساختار اصلی پیکربندی
type Config struct {
	App       AppConfig       `mapstructure:"app"`
	Server    ServerConfig    `mapstructure:"server"`
	Database  DatabaseConfig  `mapstructure:"database"`
	Redis     RedisConfig     `mapstructure:"redis"`
	JWT       JWTConfig       `mapstructure:"jwt"`
	Logging   LoggingConfig   `mapstructure:"logging"`
	RateLimit RateLimitConfig `mapstructure:"rate_limit"`
	Features  FeaturesConfig  `mapstructure:"features"`
	External  ExternalConfig  `mapstructure:"external"`
}

// AppConfig تنظیمات برنامه
type AppConfig struct {
	Name        string `mapstructure:"name"`
	Version     string `mapstructure:"version"`
	Environment string `mapstructure:"environment"`
	Debug       bool   `mapstructure:"debug"`
}

// ServerConfig تنظیمات سرور
type ServerConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

// DatabaseConfig تنظیمات دیتابیس
type DatabaseConfig struct {
	Driver          string        `mapstructure:"driver"`
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	Username        string        `mapstructure:"username"`
	Password        string        `mapstructure:"password"`
	Database        string        `mapstructure:"database"`
	SSLMode         string        `mapstructure:"ssl_mode"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

// RedisConfig تنظیمات Redis
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

// JWTConfig تنظیمات JWT
type JWTConfig struct {
	Secret     string        `mapstructure:"secret"`
	Expiration time.Duration `mapstructure:"expiration"`
	Issuer     string        `mapstructure:"issuer"`
}

// LoggingConfig تنظیمات لاگ
type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"` // json, text
	Output string `mapstructure:"output"` // stdout, file
	File   string `mapstructure:"file"`
}

// RateLimitConfig تنظیمات محدودیت نرخ
type RateLimitConfig struct {
	Enabled   bool          `mapstructure:"enabled"`
	Requests  int           `mapstructure:"requests"`
	PerSecond time.Duration `mapstructure:"per_second"`
}

// FeaturesConfig تنظیمات ویژگی‌ها (feature flags)
type FeaturesConfig struct {
	EnableMetrics   bool `mapstructure:"enable_metrics"`
	EnableProfiling bool `mapstructure:"enable_profiling"`
	EnableSwagger   bool `mapstructure:"enable_swagger"`
	EnableTracing   bool `mapstructure:"enable_tracing"`
}

// ExternalConfig تنظیمات سرویس‌های خارجی
type ExternalConfig struct {
	APIKey    string            `mapstructure:"api_key"`
	Endpoints map[string]string `mapstructure:"endpoints"`
	Timeouts  map[string]int    `mapstructure:"timeouts"`
}

// ============================================================================
// بخش 3: نمونه فایل‌های پیکربندی
// ============================================================================

func configFileExamples() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📄 CONFIGURATION FILE EXAMPLES")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
# ============================================================================
# config.yaml (مثال فایل پیکربندی YAML)
# ============================================================================

app:
  name: "myapp"
  version: "1.0.0"
  environment: "development"  # development, staging, production
  debug: true

server:
  host: "0.0.0.0"
  port: 8080
  read_timeout: "30s"
  write_timeout: "30s"
  idle_timeout: "120s"

database:
  driver: "postgres"
  host: "localhost"
  port: 5432
  username: "postgres"
  password: "postgres"
  database: "myapp"
  ssl_mode: "disable"
  max_open_conns: 25
  max_idle_conns: 10
  conn_max_lifetime: "5m"

redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0
  pool_size: 10

jwt:
  secret: "your-secret-key-change-in-production"
  expiration: "24h"
  issuer: "myapp"

logging:
  level: "info"      # debug, info, warn, error
  format: "json"     # json, text
  output: "stdout"   # stdout, file
  file: "/var/log/myapp/app.log"

rate_limit:
  enabled: true
  requests: 100
  per_second: "1m"

features:
  enable_metrics: true
  enable_profiling: false
  enable_swagger: true
  enable_tracing: false

external:
  api_key: "your-api-key"
  endpoints:
    users: "https://api.example.com/users"
    products: "https://api.example.com/products"
  timeouts:
    users: 30
    products: 60
`)

	fmt.Println(`
# ============================================================================
# config.json (مثال فایل پیکربندی JSON)
# ============================================================================

{
  "app": {
    "name": "myapp",
    "version": "1.0.0",
    "environment": "development",
    "debug": true
  },
  "server": {
    "host": "0.0.0.0",
    "port": 8080,
    "read_timeout": "30s",
    "write_timeout": "30s"
  },
  "database": {
    "driver": "postgres",
    "host": "localhost",
    "port": 5432,
    "username": "postgres",
    "password": "postgres",
    "database": "myapp"
  }
}
`)

	fmt.Println(`
# ============================================================================
# .env (فایل متغیرهای محیطی)
# ============================================================================

APP_NAME=myapp
APP_VERSION=1.0.0
APP_ENVIRONMENT=development
APP_DEBUG=true

SERVER_HOST=0.0.0.0
SERVER_PORT=8080
SERVER_READ_TIMEOUT=30s
SERVER_WRITE_TIMEOUT=30s

DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USERNAME=postgres
DB_PASSWORD=postgres
DB_NAME=myapp

REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

JWT_SECRET=your-secret-key
JWT_EXPIRATION=24h

LOG_LEVEL=info
LOG_FORMAT=json
`)

	fmt.Println(`
# ============================================================================
# config.toml (مثال فایل پیکربندی TOML)
# ============================================================================

[app]
  name = "myapp"
  version = "1.0.0"
  environment = "development"
  debug = true

[server]
  host = "0.0.0.0"
  port = 8080
  read_timeout = "30s"
  write_timeout = "30s"

[database]
  driver = "postgres"
  host = "localhost"
  port = 5432
  username = "postgres"
  password = "postgres"
  database = "myapp"
`)
}

// ============================================================================
// بخش 4: راه‌اندازی Viper
// ============================================================================

// LoadConfig بارگذاری تنظیمات با Viper
func LoadConfig() (*Config, error) {
	// 1. تنظیم نام فایل و مسیر
	viper.SetConfigName("config")      // نام فایل (بدون پسوند)
	viper.SetConfigType("yaml")        // نوع فایل (yaml, json, toml, env)
	viper.AddConfigPath(".")           // مسیر فعلی
	viper.AddConfigPath("./configs")   // دایرکتوری configs
	viper.AddConfigPath("/etc/myapp/") // مسیر سیستم

	// 2. مقادیر پیش‌فرض
	setDefaults()

	// 3. خواندن فایل پیکربندی
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Config file not found, using defaults and environment variables")
		} else {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	// 4. خواندن از Environment Variables
	viper.SetEnvPrefix("MYAPP")                            // پیشوند متغیرها (MYAPP_...)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // تبدیل نقطه به زیرخط
	viper.AutomaticEnv()                                   // فعال کردن خودکار env vars

	// 5. Watch کردن تغییرات فایل (اختیاری)
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("Config file changed: %s", e.Name)
		// در اینجا می‌توان برنامه را ری‌لود کرد
	})

	// 6. تبدیل به struct
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// 7. اعتبارسنجی
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}

// setDefaults تنظیم مقادیر پیش‌فرض
func setDefaults() {
	// App defaults
	viper.SetDefault("app.name", "myapp")
	viper.SetDefault("app.version", "1.0.0")
	viper.SetDefault("app.environment", "development")
	viper.SetDefault("app.debug", true)

	// Server defaults
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.read_timeout", "30s")
	viper.SetDefault("server.write_timeout", "30s")
	viper.SetDefault("server.idle_timeout", "120s")

	// Database defaults
	viper.SetDefault("database.driver", "postgres")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.max_open_conns", 25)
	viper.SetDefault("database.max_idle_conns", 10)
	viper.SetDefault("database.conn_max_lifetime", "5m")

	// Redis defaults
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("redis.pool_size", 10)

	// JWT defaults
	viper.SetDefault("jwt.expiration", "24h")
	viper.SetDefault("jwt.issuer", "myapp")

	// Logging defaults
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "json")
	viper.SetDefault("logging.output", "stdout")

	// Rate limit defaults
	viper.SetDefault("rate_limit.enabled", true)
	viper.SetDefault("rate_limit.requests", 100)
	viper.SetDefault("rate_limit.per_second", "1m")

	// Features defaults
	viper.SetDefault("features.enable_metrics", true)
	viper.SetDefault("features.enable_profiling", false)
	viper.SetDefault("features.enable_swagger", true)
}

// validateConfig اعتبارسنجی تنظیمات
func validateConfig(config *Config) error {
	// اعتبارسنجی محیط
	validEnvs := map[string]bool{
		"development": true,
		"staging":     true,
		"production":  true,
	}
	if !validEnvs[config.App.Environment] {
		return fmt.Errorf("invalid environment: %s (must be development, staging, or production)", config.App.Environment)
	}

	// اعتبارسنجی پورت
	if config.Server.Port < 1 || config.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", config.Server.Port)
	}

	// اعتبارسنجی دیتابیس
	if config.Database.Driver == "" {
		return fmt.Errorf("database driver is required")
	}
	if config.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if config.Database.Username == "" {
		return fmt.Errorf("database username is required")
	}
	if config.Database.Database == "" {
		return fmt.Errorf("database name is required")
	}

	// اعتبارسنجی JWT
	if config.JWT.Secret == "" {
		// در پروداکشن باید error بدهد، در توسعه می‌تواند مقدار پیش‌فرض داشته باشد
		if config.App.Environment == "production" {
			return fmt.Errorf("JWT secret is required in production")
		}
	}

	return nil
}

// ============================================================================
// بخش 5: دسترسی به مقادیر (با روش‌های مختلف)
// ============================================================================

func demonstrateValueAccess() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🔍 ACCESSING CONFIGURATION VALUES")
	fmt.Println(strings.Repeat("=", 80))

	// 1. دسترسی مستقیم
	appName := viper.GetString("app.name")
	appVersion := viper.GetString("app.version")
	serverPort := viper.GetInt("server.port")
	debug := viper.GetBool("app.debug")
	timeout := viper.GetDuration("server.read_timeout")

	fmt.Println("\n--- Direct Access ---")
	fmt.Printf("  app.name: %s\n", appName)
	fmt.Printf("  app.version: %s\n", appVersion)
	fmt.Printf("  server.port: %d\n", serverPort)
	fmt.Printf("  app.debug: %v\n", debug)
	fmt.Printf("  server.read_timeout: %v\n", timeout)

	// 2. با مقدار پیش‌فرض
	nonExistent := viper.GetString("non.existent.key")
	defaultValue := viper.GetString("non.existent.key")
	fmt.Printf("\n  Non-existent key: %q (zero value)\n", nonExistent)
	fmt.Printf("  With default: %q\n", defaultValue)

	// 3. بررسی وجود کلید
	if viper.IsSet("database.password") {
		fmt.Printf("\n  database.password is set\n")
	}

	// 4. دسترسی به nested values
	dbHost := viper.GetString("database.host")
	dbPort := viper.GetInt("database.port")
	fmt.Printf("  Database: %s:%d\n", dbHost, dbPort)

	// 5. دسترسی به map
	externalEndpoints := viper.GetStringMapString("external.endpoints")
	fmt.Printf("\n  External endpoints: %v\n", externalEndpoints)

	// 6. دسترسی به slice
	// (در صورت وجود در فایل پیکربندی)
	// servers := viper.GetStringSlice("servers")
}

// ============================================================================
// بخش 6: Environment Variables Mapping
// ============================================================================

func demonstrateEnvMapping() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🌍 ENVIRONMENT VARIABLES MAPPING")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ CONFIG KEY                │ ENV VARIABLE              │ EXAMPLE VALUE        │
├───────────────────────────┼───────────────────────────┼──────────────────────┤
│ app.name                  │ MYAPP_APP_NAME           │ myapp                │
│ app.environment           │ MYAPP_APP_ENVIRONMENT    │ production           │
│ server.port               │ MYAPP_SERVER_PORT        │ 8080                 │
│ database.host             │ MYAPP_DATABASE_HOST      │ postgres.example.com │
│ database.password         │ MYAPP_DATABASE_PASSWORD  │ secret123            │
│ jwt.secret                │ MYAPP_JWT_SECRET         │ super-secret-key     │
│ redis.password            │ MYAPP_REDIS_PASSWORD     │ redis-pass           │
│ logging.level             │ MYAPP_LOGGING_LEVEL      │ debug                │
└─────────────────────────────────────────────────────────────────────────────┘

📝 USAGE:

   # Development
   export MYAPP_APP_ENVIRONMENT=development
   export MYAPP_SERVER_PORT=8080
   export MYAPP_DATABASE_HOST=localhost
   go run main.go

   # Production
   export MYAPP_APP_ENVIRONMENT=production
   export MYAPP_SERVER_PORT=80
   export MYAPP_DATABASE_HOST=prod-db.example.com
   export MYAPP_DATABASE_PASSWORD=${DB_PASSWORD}
   export MYAPP_JWT_SECRET=${JWT_SECRET}
   go run main.go
`)
}

// ============================================================================
// بخش 7: Config Files per Environment
// ============================================================================

// LoadConfigForEnvironment بارگذاری تنظیمات مخصوص محیط
func LoadConfigForEnvironment(env string) error {
	// بارگذاری فایل پایه
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./configs")

	if err := viper.MergeInConfig(); err != nil {
		return err
	}

	// بارگذاری فایل مخصوص محیط (مثل config.production.yaml)
	viper.SetConfigName(fmt.Sprintf("config.%s", env))
	if err := viper.MergeInConfig(); err != nil {
		// فایل محیط وجود ندارد - مشکلی نیست
		log.Printf("Environment config file not found: config.%s.yaml", env)
	}

	return nil
}

// ============================================================================
// بخش 8: Command Line Flags Integration
// ============================================================================

func demonstrateFlags() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🚩 COMMAND LINE FLAGS")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
// اتصال Viper با command line flags
// نیاز به پکیج flag (یا cobra)

import (
	"flag"
	"github.com/spf13/viper"
)

func init() {
	// تعریف flags
	configPath := flag.String("config", "", "Path to config file")
	port := flag.Int("port", 8080, "Server port")
	debug := flag.Bool("debug", false, "Enable debug mode")
	
	flag.Parse()
	
	// اتصال به Viper
	if *configPath != "" {
		viper.SetConfigFile(*configPath)
	}
	
	// Bind flags به Viper
	viper.BindPFlag("server.port", flag.Lookup("port"))
	viper.BindPFlag("app.debug", flag.Lookup("debug"))
}

📝 USAGE:

   # Using flags
   $ go run main.go --config=config.yaml --port=9090 --debug
   
   # Flags override config files and env vars
   # Priority: flag > env > config file > default
`)
}

// ============================================================================
// بخش 9: Hot Reload (Watch Config Changes)
// ============================================================================

// WatchConfig شروع watching تغییرات فایل
func WatchConfig(callback func()) {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("Config file changed: %s", e.Name)
		callback()
	})
}

// ReloadConfigCallback تابع callback برای reload
func ReloadConfigCallback() {
	// بازخوانی تنظیمات
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Printf("Failed to reload config: %v", err)
		return
	}

	log.Printf("Config reloaded successfully")
	// در اینجا می‌توان تنظیمات جدید را اعمال کرد
}

// ============================================================================
// بخش 10: Remote Config (Consul, Etcd, etc.)
// ============================================================================

func demonstrateRemoteConfig() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("☁️ REMOTE CONFIGURATION")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
// خواندن تنظیمات از remote sources مانند Consul یا Etcd

import (
	"github.com/spf13/viper/remote"
	_ "github.com/hashicorp/consul/api"
)

func init() {
	// خواندن از Consul
	viper.AddRemoteProvider("consul", "localhost:8500", "myapp/config")
	viper.SetConfigType("yaml")
	
	if err := viper.ReadRemoteConfig(); err != nil {
		log.Fatal(err)
	}
	
	// Watching changes
	go func() {
		for {
			time.Sleep(30 * time.Second)
			if err := viper.WatchRemoteConfig(); err != nil {
				log.Printf("Error watching remote config: %v", err)
			}
		}
	}()
}
`)
}

// ============================================================================
// بخش 11: Complete Example
// ============================================================================

func runViperExample() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🚀 RUNNING VIPER EXAMPLE")
	fmt.Println(strings.Repeat("=", 80))

	// بارگذاری تنظیمات
	config, err := LoadConfig()
	if err != nil {
		log.Printf("Error loading config: %v", err)
		return
	}

	// نمایش تنظیمات
	fmt.Println("\n--- Loaded Configuration ---")
	fmt.Printf("App Name: %s\n", config.App.Name)
	fmt.Printf("App Version: %s\n", config.App.Version)
	fmt.Printf("Environment: %s\n", config.App.Environment)
	fmt.Printf("Debug: %v\n", config.App.Debug)
	fmt.Printf("Server: %s:%d\n", config.Server.Host, config.Server.Port)
	fmt.Printf("Database: %s@%s:%d/%s\n",
		config.Database.Username,
		config.Database.Host,
		config.Database.Port,
		config.Database.Database)
	fmt.Printf("Redis: %s:%d\n", config.Redis.Host, config.Redis.Port)
	fmt.Printf("JWT Expiration: %v\n", config.JWT.Expiration)
	fmt.Printf("Log Level: %s\n", config.Logging.Level)
	fmt.Printf("Rate Limit: %v req per %v\n", config.RateLimit.Requests, config.RateLimit.PerSecond)

	// دسترسی به مقادیر
	demonstrateValueAccess()
}

// ============================================================================
// بخش 12: Best Practices
// ============================================================================

func bestPractices() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("💡 VIPER BEST PRACTICES")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ 1. CONFIGURATION STRUCTURE                                                │
├─────────────────────────────────────────────────────────────────────────────┤
│    ✔ Use nested structs for better organization                           │
│    ✔ Use mapstructure tags for struct fields                              │
│    ✔ Always set defaults for all fields                                   │
│    ✔ Validate config after loading                                        │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 2. FILE ORGANIZATION                                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│    ✔ Store config files in ./configs directory                            │
│    ✔ Use different files per environment                                  │
│    ✔ Use YAML for human-readable configs                                  │
│    ✔ Use JSON for programmatic generation                                 │
│    ✔ Use .env for simple configurations                                   │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 3. ENVIRONMENT VARIABLES                                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│    ✔ Use consistent naming convention                                     │
│    ✔ Set prefix to avoid conflicts                                        │
│    ✔ Never store secrets in config files                                  │
│    ✔ Use env vars for sensitive data                                      │
│    ✔ Document all required env vars                                       │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 4. PERFORMANCE                                                            │
├─────────────────────────────────────────────────────────────────────────────┤
│    ✔ Load config once at startup                                          │
│    ✔ Cache config in global variable                                      │
│    ✔ Use Get* methods for single values                                   │
│    ✔ Avoid unmarshaling multiple times                                    │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 5. SECURITY                                                               │
├─────────────────────────────────────────────────────────────────────────────┤
│    ✔ Never log config values with secrets                                 │
│    ✔ Use different secrets per environment                                │
│    ✔ Rotate secrets regularly                                             │
│    ✔ Validate all inputs from config                                      │
│    ✔ Use read-only config for production                                  │
└─────────────────────────────────────────────────────────────────────────────┘
`)
}

// ============================================================================
// بخش 13: Common Mistakes
// ============================================================================

func commonMistakes() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("⚠️ COMMON MISTAKES")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ MISTAKE 1: Not checking if config file exists                             │
├─────────────────────────────────────────────────────────────────────────────┤
│    ❌ viper.ReadInConfig() // panics if file not found                     │
│    ✅ Check error and use defaults                                         │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ MISTAKE 2: Ignoring environment variable prefix                           │
├─────────────────────────────────────────────────────────────────────────────┤
│    ❌ viper.AutomaticEnv() // can conflict with system env vars            │
│    ✅ viper.SetEnvPrefix("MYAPP")                                          │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ MISTAKE 3: Not using mapstructure tags                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│    ❌ type Config struct { Port int ` + "`" + `json:"port"` + "`" + ` }                         │
│    ✅ type Config struct { Port int ` + "`" + `mapstructure:"port"` + "`" + ` }                 │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ MISTAKE 4: Hardcoding config file path                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│    ❌ viper.AddConfigPath("/app/config.yaml")                             │
│    ✅ viper.AddConfigPath(".", "./configs", "/etc/myapp")                 │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ MISTAKE 5: Not handling config reload                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│    ❌ WatchConfig without callback                                         │
│    ✅ Implement reload logic                                               │
└─────────────────────────────────────────────────────────────────────────────┘
`)
}

// ============================================================================
// بخش 14: Main
// ============================================================================

func main() {
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("🎯 COMPLETE VIPER GUIDE")
	fmt.Println("Configuration Management for Go Applications")
	fmt.Println(strings.Repeat("=", 80))

	configFileExamples()
	bestPractices()
	commonMistakes()
	demonstrateEnvMapping()
	demonstrateFlags()
	demonstrateRemoteConfig()

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🚀 Running Viper Example")
	fmt.Println(strings.Repeat("=", 80))

	runViperExample()

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🎯 VIPER - COMPLETE")
	fmt.Println("Flexible configuration for your Go apps!")
	fmt.Println(strings.Repeat("=", 80))
}

// تابع کمکی برای تبدیل به JSON (برای نمایش)
func toJSON(v interface{}) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}
