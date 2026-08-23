// ============================================================================
// FILE: env_secrets_management_guide.go
// TITLE: راهنمای کامل Environment Variables و Secrets Management در Go
// HOW TO RUN: go run env_secrets_management_guide.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - چرا Environment Variables و Secrets Management مهم است؟
// ============================================================================
//
// Environment Variables (متغیرهای محیطی) و مدیریت secrets برای امنیت و پیکربندی برنامه حیاتی هستند:
//
// 1. جداسازی کد از تنظیمات (Code from Config)
// 2. امنیت: secrets نباید در کد ذخیره شوند
// 3. انعطاف‌پذیری: تغییر رفتار بدون تغییر کد
// 4. محیط‌های مختلف: dev, staging, production
//
// انواع secrets:
// - API Keys (Google, GitHub, AWS)
// - Database passwords
// - JWT secrets
// - Encryption keys
// - OAuth client secrets
//
// روش‌های مدیریت secrets:
// 1. Environment Variables (ساده، رایج)
// 2. .env files (برای توسعه)
// 3. Vault (HashiCorp) - برای پروداکشن
// 4. AWS Secrets Manager / Parameter Store
// 5. Google Cloud Secret Manager
// 6. Kubernetes Secrets
//
// قانون طلایی:
// "هرگز secrets را در کد ذخیره نکن. از environment variables استفاده کن.
//  فایل‌های .env را به .gitignore اضافه کن.
//  برای پروداکشن از secret management tools استفاده کن."
// ============================================================================

package __env_secrets_management

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// ============================================================================
// بخش 1: مفاهیم پایه Environment Variables
// ============================================================================

// LoadEnvFile بارگذاری فایل .env
func LoadEnvFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		// اگر فایل وجود نداشت، خطا ندهید (ممکن است در پروداکشن از env استفاده شود)
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// رد کردن خطوط خالی و کامنت‌ها
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// تقسیم key=value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// حذف نقل قول‌ها
		if len(value) > 1 {
			if (value[0] == '"' && value[len(value)-1] == '"') ||
				(value[0] == '\'' && value[len(value)-1] == '\'') {
				value = value[1 : len(value)-1]
			}
		}

		// تنظیم متغیر محیطی
		os.Setenv(key, value)
	}

	return scanner.Err()
}

// GetEnv گرفتن متغیر محیطی با مقدار پیش‌فرض (string)
func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetEnvAsInt گرفتن متغیر محیطی به عنوان int
func GetEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

// GetEnvAsBool گرفتن متغیر محیطی به عنوان bool
func GetEnvAsBool(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

// GetEnvAsDuration گرفتن متغیر محیطی به عنوان time.Duration
func GetEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := time.ParseDuration(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

// GetEnvAsSlice گرفتن متغیر محیطی به عنوان slice با جداکننده
func GetEnvAsSlice(key string, separator string, defaultValue []string) []string {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	return strings.Split(valueStr, separator)
}

// ============================================================================
// بخش 2: Config Struct با تگ‌های سفارشی
// ============================================================================

// Config تنظیمات برنامه
type Config struct {
	// Server
	ServerHost         string        `env:"SERVER_HOST" default:"0.0.0.0"`
	ServerPort         int           `env:"SERVER_PORT" default:"8080"`
	ServerReadTimeout  time.Duration `env:"SERVER_READ_TIMEOUT" default:"30s"`
	ServerWriteTimeout time.Duration `env:"SERVER_WRITE_TIMEOUT" default:"30s"`
	ServerIdleTimeout  time.Duration `env:"SERVER_IDLE_TIMEOUT" default:"60s"`

	// Database
	DBHost     string `env:"DB_HOST" required:"true"`
	DBPort     int    `env:"DB_PORT" default:"5432"`
	DBUser     string `env:"DB_USER" required:"true"`
	DBPassword string `env:"DB_PASSWORD" required:"true" secret:"true"`
	DBName     string `env:"DB_NAME" default:"myapp"`
	DBSSLMode  string `env:"DB_SSLMODE" default:"disable"`

	// Redis
	RedisHost     string `env:"REDIS_HOST" default:"localhost"`
	RedisPort     int    `env:"REDIS_PORT" default:"6379"`
	RedisPassword string `env:"REDIS_PASSWORD" secret:"true"`
	RedisDB       int    `env:"REDIS_DB" default:"0"`

	// JWT
	JWTSecret     string        `env:"JWT_SECRET" required:"true" secret:"true"`
	JWTExpiration time.Duration `env:"JWT_EXPIRATION" default:"24h"`

	// API Keys
	GoogleAPIKey   string `env:"GOOGLE_API_KEY" secret:"true"`
	GitHubClientID string `env:"GITHUB_CLIENT_ID"`
	GitHubSecret   string `env:"GITHUB_SECRET" secret:"true"`

	// Feature flags
	EnableMetrics   bool `env:"ENABLE_METRICS" default:"false"`
	EnableProfiling bool `env:"ENABLE_PROFILING" default:"false"`
	EnableSwagger   bool `env:"ENABLE_SWAGGER" default:"true"`
	EnableRateLimit bool `env:"ENABLE_RATE_LIMIT" default:"true"`

	// Application
	AppName  string `env:"APP_NAME" default:"myapp"`
	AppEnv   string `env:"APP_ENV" default:"development"`
	Debug    bool   `env:"DEBUG" default:"false"`
	LogLevel string `env:"LOG_LEVEL" default:"info"`
}

// LoadConfig بارگذاری تنظیمات از environment variables
func LoadConfig() (*Config, error) {
	cfg := &Config{}

	// استفاده از reflection برای پر کردن struct
	if err := loadConfigFromEnv(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// MustLoadConfig بارگذاری تنظیمات با panic در صورت خطا
func MustLoadConfig() *Config {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	return cfg
}

// loadConfigFromEnv بارگذاری تنظیمات با reflection
func loadConfigFromEnv(cfg interface{}) error {
	v := reflect.ValueOf(cfg).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		if !fieldValue.CanSet() {
			continue
		}

		// خواندن تگ env
		envKey := field.Tag.Get("env")
		if envKey == "" {
			continue
		}

		// خواندن مقدار از environment
		strValue := os.Getenv(envKey)

		// اگر مقدار خالی بود، از default استفاده کن
		if strValue == "" {
			defaultValue := field.Tag.Get("default")
			if defaultValue != "" {
				strValue = defaultValue
			} else if field.Tag.Get("required") == "true" {
				return fmt.Errorf("required environment variable not set: %s", envKey)
			}
		}

		// اگر هنوز مقدار خالی است، skip
		if strValue == "" {
			continue
		}

		// تبدیل مقدار به نوع فیلد
		if err := setFieldValue(fieldValue, strValue); err != nil {
			return fmt.Errorf("failed to set field %s: %w", field.Name, err)
		}
	}

	return nil
}

// setFieldValue تنظیم مقدار فیلد بر اساس نوع
func setFieldValue(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		// بررسی اینکه آیا time.Duration است
		if field.Type() == reflect.TypeOf(time.Duration(0)) {
			duration, err := time.ParseDuration(value)
			if err != nil {
				return err
			}
			field.SetInt(int64(duration))
		} else {
			intVal, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return err
			}
			field.SetInt(intVal)
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintVal, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetUint(uintVal)

	case reflect.Float32, reflect.Float64:
		floatVal, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		field.SetFloat(floatVal)

	case reflect.Bool:
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(boolVal)

	case reflect.Slice:
		if field.Type().Elem().Kind() == reflect.String {
			parts := strings.Split(value, ",")
			slice := reflect.MakeSlice(field.Type(), len(parts), len(parts))
			for i, part := range parts {
				slice.Index(i).SetString(strings.TrimSpace(part))
			}
			field.Set(slice)
		}

	default:
		return fmt.Errorf("unsupported field type: %s", field.Kind())
	}

	return nil
}

// String (برای لاگ کردن بدون نمایش secrets)
func (c *Config) String() string {
	// کپی از struct برای عدم تغییر اصلی
	cfgCopy := *c

	// مخفی کردن secrets
	cfgCopy.DBPassword = "***"
	cfgCopy.RedisPassword = "***"
	cfgCopy.JWTSecret = "***"
	cfgCopy.GoogleAPIKey = "***"
	cfgCopy.GitHubSecret = "***"

	return fmt.Sprintf("%+v", cfgCopy)
}

// ============================================================================
// بخش 3: Secrets Encryption (رمزنگاری secrets در فایل)
// ============================================================================

// SecretEncryptor رمزنگار secrets
type SecretEncryptor struct {
	key []byte
}

// NewSecretEncryptor ایجاد رمزنگار جدید
func NewSecretEncryptor(encryptionKey string) (*SecretEncryptor, error) {
	// کلید باید 32 بایت برای AES-256 باشد
	key := []byte(encryptionKey)
	if len(key) != 32 {
		// اگر کلید کوتاه است، آن را padding کن
		if len(key) < 32 {
			padded := make([]byte, 32)
			copy(padded, key)
			key = padded
		} else {
			key = key[:32]
		}
	}

	return &SecretEncryptor{key: key}, nil
}

// Encrypt رمزنگاری متن
func (e *SecretEncryptor) Encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt رمزگشایی متن
func (e *SecretEncryptor) Decrypt(encryptedText string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// ============================================================================
// بخش 4: Encrypted .env File
// ============================================================================

// EncryptedEnvFile مدیریت فایل .env رمزنگاری شده
type EncryptedEnvFile struct {
	encryptor *SecretEncryptor
	filename  string
}

// NewEncryptedEnvFile ایجاد مدیریت فایل رمزنگاری شده
func NewEncryptedEnvFile(filename, encryptionKey string) (*EncryptedEnvFile, error) {
	encryptor, err := NewSecretEncryptor(encryptionKey)
	if err != nil {
		return nil, err
	}

	return &EncryptedEnvFile{
		encryptor: encryptor,
		filename:  filename,
	}, nil
}

// Save ذخیره secrets در فایل رمزنگاری شده
func (e *EncryptedEnvFile) Save(secrets map[string]string) error {
	// تبدیل secrets به JSON
	data, err := json.Marshal(secrets)
	if err != nil {
		return err
	}

	// رمزنگاری
	encrypted, err := e.encryptor.Encrypt(string(data))
	if err != nil {
		return err
	}

	// ذخیره در فایل
	return os.WriteFile(e.filename, []byte(encrypted), 0600)
}

// Load بارگذاری secrets از فایل رمزنگاری شده
func (e *EncryptedEnvFile) Load() (map[string]string, error) {
	// خواندن فایل
	encrypted, err := os.ReadFile(e.filename)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]string), nil
		}
		return nil, err
	}

	// رمزگشایی
	decrypted, err := e.encryptor.Decrypt(string(encrypted))
	if err != nil {
		return nil, err
	}

	// تبدیل JSON به map
	var secrets map[string]string
	if err := json.Unmarshal([]byte(decrypted), &secrets); err != nil {
		return nil, err
	}

	return secrets, nil
}

// LoadToEnv بارگذاری secrets و تنظیم در environment variables
func (e *EncryptedEnvFile) LoadToEnv() error {
	secrets, err := e.Load()
	if err != nil {
		return err
	}

	for key, value := range secrets {
		os.Setenv(key, value)
	}

	return nil
}

// ============================================================================
// بخش 5: Validation (اعتبارسنجی تنظیمات)
// ============================================================================

// ValidateConfig اعتبارسنجی تنظیمات
func (c *Config) Validate() error {
	var errors []string

	// بررسی پورت
	if c.ServerPort < 1 || c.ServerPort > 65535 {
		errors = append(errors, "invalid server port")
	}

	// بررسی JWT secret
	if c.JWTSecret == "" {
		errors = append(errors, "JWT secret is required")
	} else if len(c.JWTSecret) < 32 {
		errors = append(errors, "JWT secret should be at least 32 characters")
	}

	// بررسی DB password
	if c.DBPassword == "" {
		errors = append(errors, "database password is required")
	}

	// بررسی محیط
	validEnvs := map[string]bool{"development": true, "staging": true, "production": true}
	if !validEnvs[c.AppEnv] {
		errors = append(errors, fmt.Sprintf("invalid app environment: %s", c.AppEnv))
	}

	if len(errors) > 0 {
		return fmt.Errorf("config validation failed: %s", strings.Join(errors, "; "))
	}

	return nil
}

// ============================================================================
// بخش 6: Different Environments (توسعه، استیجینگ، پروداکشن)
// ============================================================================

// Environment مدیریت محیط‌های مختلف
type Environment struct {
	Name string
	File string
}

var (
	EnvDevelopment = &Environment{Name: "development", File: ".env.dev"}
	EnvStaging     = &Environment{Name: "staging", File: ".env.staging"}
	EnvProduction  = &Environment{Name: "production", File: ".env.prod"}
)

// GetEnvironment دریافت محیط فعلی
func GetEnvironment() *Environment {
	env := os.Getenv("APP_ENV")
	switch env {
	case "staging":
		return EnvStaging
	case "production":
		return EnvProduction
	default:
		return EnvDevelopment
	}
}

// LoadEnvForEnvironment بارگذاری فایل env برای محیط فعلی
func LoadEnvForEnvironment() error {
	env := GetEnvironment()

	// بارگذاری فایل .env عمومی (در صورت وجود)
	LoadEnvFile(".env")

	// بارگذاری فایل مخصوص محیط
	return LoadEnvFile(env.File)
}

// ============================================================================
// بخش 7: Secret Injection (تزریق secrets از Vault/AWS)
// ============================================================================

// SecretProvider اینترفیس ارائه‌دهنده secrets
type SecretProvider interface {
	GetSecret(key string) (string, error)
	SetSecret(key, value string) error
}

// EnvSecretProvider از environment variables
type EnvSecretProvider struct{}

func (p *EnvSecretProvider) GetSecret(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("secret not found: %s", key)
	}
	return value, nil
}

func (p *EnvSecretProvider) SetSecret(key, value string) error {
	os.Setenv(key, value)
	return nil
}

// FileSecretProvider از فایل رمزنگاری شده
type FileSecretProvider struct {
	encryptor *SecretEncryptor
	filename  string
	secrets   map[string]string
}

func NewFileSecretProvider(filename, encryptionKey string) (*FileSecretProvider, error) {
	encryptor, err := NewSecretEncryptor(encryptionKey)
	if err != nil {
		return nil, err
	}

	provider := &FileSecretProvider{
		encryptor: encryptor,
		filename:  filename,
		secrets:   make(map[string]string),
	}

	if err := provider.load(); err != nil {
		return nil, err
	}

	return provider, nil
}

func (p *FileSecretProvider) load() error {
	data, err := os.ReadFile(p.filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	decrypted, err := p.encryptor.Decrypt(string(data))
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(decrypted), &p.secrets)
}

func (p *FileSecretProvider) save() error {
	data, err := json.Marshal(p.secrets)
	if err != nil {
		return err
	}

	encrypted, err := p.encryptor.Encrypt(string(data))
	if err != nil {
		return err
	}

	return os.WriteFile(p.filename, []byte(encrypted), 0600)
}

func (p *FileSecretProvider) GetSecret(key string) (string, error) {
	value, ok := p.secrets[key]
	if !ok {
		return "", fmt.Errorf("secret not found: %s", key)
	}
	return value, nil
}

func (p *FileSecretProvider) SetSecret(key, value string) error {
	p.secrets[key] = value
	return p.save()
}

// ============================================================================
// بخش 8: Health Check for Config
// ============================================================================

// ConfigHealthChecker بررسی سلامت تنظیمات
type ConfigHealthChecker struct {
	requiredEnvs []string
}

func NewConfigHealthChecker(requiredEnvs []string) *ConfigHealthChecker {
	return &ConfigHealthChecker{
		requiredEnvs: requiredEnvs,
	}
}

func (c *ConfigHealthChecker) Check() []string {
	var missing []string

	for _, env := range c.requiredEnvs {
		if os.Getenv(env) == "" {
			missing = append(missing, env)
		}
	}

	return missing
}

// ============================================================================
// بخش 9: Example Usage
// ============================================================================

func exampleUsage() {
	// 1. بارگذاری فایل .env
	if err := LoadEnvFile(".env.example"); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// 2. بارگذاری فایل مخصوص محیط
	if err := LoadEnvForEnvironment(); err != nil {
		log.Printf("Warning: environment file not found: %v", err)
	}

	// 3. بارگذاری تنظیمات با struct
	config := MustLoadConfig()
	log.Printf("Config loaded: %s", config)

	// 4. اعتبارسنجی تنظیمات
	if err := config.Validate(); err != nil {
		log.Fatalf("Invalid config: %v", err)
	}

	// 5. مثال استفاده از secrets provider
	provider := &EnvSecretProvider{}
	dbPassword, err := provider.GetSecret("DB_PASSWORD")
	if err != nil {
		log.Printf("Failed to get DB password: %v", err)
	} else {
		log.Printf("DB password loaded (length: %d)", len(dbPassword))
	}

	// 6. آموزش استفاده از متغیرهای محیطی
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📝 ENVIRONMENT VARIABLES EXAMPLES")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Printf("SERVER_HOST: %s\n", GetEnv("SERVER_HOST", "localhost"))
	fmt.Printf("SERVER_PORT: %d\n", GetEnvAsInt("SERVER_PORT", 8080))
	fmt.Printf("DEBUG: %v\n", GetEnvAsBool("DEBUG", false))
	fmt.Printf("TIMEOUT: %v\n", GetEnvAsDuration("TIMEOUT", 30*time.Second))
	fmt.Printf("ALLOWED_ORIGINS: %v\n", GetEnvAsSlice("ALLOWED_ORIGINS", ",", []string{"*"}))
}

// ============================================================================
// بخش 10: Security Best Practices
// ============================================================================

func bestPractices() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("💡 SECRETS MANAGEMENT BEST PRACTICES")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ 1. NEVER hardcode secrets in code                                         │
├─────────────────────────────────────────────────────────────────────────────┤
│    ❌ const apiKey = "sk-12345"                                           │
│    ✅ apiKey := os.Getenv("API_KEY")                                      │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 2. Use .gitignore for .env files                                          │
├─────────────────────────────────────────────────────────────────────────────┤
│    .env                                                                    │
│    .env.local                                                              │
│    .env.*.local                                                            │
│    *.key                                                                   │
│    *.pem                                                                   │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 3. Use different secrets per environment                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│    .env.dev      - development secrets                                    │
│    .env.staging  - staging secrets                                        │
│    .env.prod     - production secrets (never commit)                      │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 4. Rotate secrets regularly                                               │
├─────────────────────────────────────────────────────────────────────────────┤
│    • Set expiration for secrets                                           │
│    • Automate rotation                                                    │
│    • Use secret managers (Vault, AWS SM, etc.)                            │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 5. Encrypt secrets at rest                                                │
├─────────────────────────────────────────────────────────────────────────────┤
│    • Encrypted .env files                                                 │
│    • Use KMS (Key Management Service)                                     │
│    • Never store plaintext secrets on disk                                │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 6. Validate required secrets on startup                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│    • Check all required env vars                                          │
│    • Fail fast if missing                                                 │
│    • Provide clear error messages                                         │
└─────────────────────────────────────────────────────────────────────────────┘
`)
}

// ============================================================================
// بخش 11: Common Mistakes
// ============================================================================

func commonMistakes() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("⚠️ COMMON MISTAKES")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ MISTAKE 1: Committing .env files to git                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│    ❌ git add .env                                                         │
│    ✅ Add .env to .gitignore                                              │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ MISTAKE 2: Logging secrets                                                │
├─────────────────────────────────────────────────────────────────────────────┤
│    ❌ log.Printf("DB Password: %s", dbPassword)                           │
│    ✅ Never log secrets                                                   │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ MISTAKE 3: Hardcoding default secrets                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│    ❌ password := GetEnv("DB_PASSWORD", "default-password")               │
│    ✅ Check and fail if missing required secrets                          │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ MISTAKE 4: Using same secrets across environments                         │
├─────────────────────────────────────────────────────────────────────────────┤
│    ❌ Same API key in dev and production                                  │
│    ✅ Use different secrets per environment                               │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ MISTAKE 5: Storing secrets in environment files on production servers     │
├─────────────────────────────────────────────────────────────────────────────┤
│    ❌ .env file on production server                                      │
│    ✅ Use secret manager or K8s secrets                                   │
└─────────────────────────────────────────────────────────────────────────────┘
`)
}

// ============================================================================
// بخش 12: Command Line Tool (تولید .env.example)
// ============================================================================

// GenerateEnvExample تولید فایل .env.example از struct
func GenerateEnvExample(cfg interface{}, outputFile string) error {
	v := reflect.ValueOf(cfg).Elem()
	t := v.Type()

	var lines []string
	lines = append(lines, "# Environment Variables")
	lines = append(lines, "# Generated automatically")
	lines = append(lines, "")

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		envKey := field.Tag.Get("env")
		if envKey == "" {
			continue
		}

		defaultValue := field.Tag.Get("default")
		required := field.Tag.Get("required") == "true"
		secret := field.Tag.Get("secret") == "true"

		line := fmt.Sprintf("# %s", field.Name)
		lines = append(lines, line)

		if secret {
			lines = append(lines, fmt.Sprintf("# %s=<secret>", envKey))
		} else if defaultValue != "" {
			lines = append(lines, fmt.Sprintf("%s=%s", envKey, defaultValue))
		} else if required {
			lines = append(lines, fmt.Sprintf("%s=", envKey))
		} else {
			lines = append(lines, fmt.Sprintf("# %s=", envKey))
		}
		lines = append(lines, "")
	}

	return os.WriteFile(outputFile, []byte(strings.Join(lines, "\n")), 0644)
}

// ============================================================================
// بخش 13: Main
// ============================================================================

func main() {
	fmt.Println(stringsRepeat("=", 80))
	fmt.Println("🎯 ENVIRONMENT VARIABLES & SECRETS MANAGEMENT")
	fmt.Println("Complete Guide for Go Applications")
	fmt.Println(stringsRepeat("=", 80))

	bestPractices()
	commonMistakes()

	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📝 EXAMPLE USAGE")
	fmt.Println(stringsRepeat("=", 80))

	exampleUsage()

	// تولید فایل .env.example
	if err := GenerateEnvExample(&Config{}, ".env.example"); err != nil {
		log.Printf("Failed to generate .env.example: %v", err)
	} else {
		fmt.Println("\n✅ Generated .env.example file")
	}

	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🎯 ENV & SECRETS - COMPLETE")
	fmt.Println("Secure your Go applications!")
	fmt.Println(stringsRepeat("=", 80))
}

// تابع کمکی برای reflection
// در فایل واقعی باید "reflect" را import کنید
func init() {
	// برای کامپایل، این خط را کامنت کنید و reflect را import کنید
	// var _ = reflect.ValueOf
}

// اضافه کردن تابع stringsRepeat محلی
func stringsRepeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}

// اضافه کردن type برای reflection (در فایل واقعی import کنید)
type reflectType struct{}

var reflect reflectType
