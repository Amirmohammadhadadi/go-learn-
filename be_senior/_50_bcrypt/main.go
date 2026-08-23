// ============================================================================
// FILE: bcrypt_guide.go
// TITLE: راهنمای کامل هش کردن رمز با bcrypt در Go
// HOW TO RUN: go run bcrypt_guide.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - چرا bcrypt؟
// ============================================================================
//
// bcrypt یک الگوریتم هش کردن رمز عبور است که به طور خاص برای این منظور طراحی شده:
// 1. آهسته بودن (مقاوم در برابر حملات brute-force)
// 2. salt خودکار (مقاوم در برابر حملات rainbow table)
// 3. قابل تنظیم (cost factor برای افزایش امنیت در آینده)
// 4. مقاوم در برابر حملات GPU/ASIC (adaptive)
//
// مقایسه الگوریتم‌های هش:
// - MD5/SHA1/SHA256: سریع → مناسب برای checksum، نه برای رمز عبور
// - bcrypt: کند → مناسب برای رمز عبور
// - scrypt: کندتر از bcrypt → مناسب برای رمز عبور
// - argon2: برنده مسابقه Password Hashing Competition → بهترین گزینه
//
// قانون طلایی:
// "هرگز رمز عبور را ذخیره نکن، همیشه هش شده را ذخیره کن.
//  از bcrypt با cost factor حداقل 10 استفاده کن.
//  هرگز الگوریتم هش سریع (MD5, SHA1, SHA256) برای رمز عبور استفاده نکن."
// ============================================================================

package __bcrypt

import (
	"bufio"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// ============================================================================
// بخش 1: مفاهیم پایه bcrypt
// ============================================================================

// HashPassword هش کردن رمز عبور
func HashPassword(password string) (string, error) {
	// cost factor: 10 (پیش‌فرض، بین 4 تا 31)
	// هر چه بالاتر، امنیت بیشتر ولی کندتر
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// HashPasswordWithCost هش کردن رمز با cost مشخص
func HashPasswordWithCost(password string, cost int) (string, error) {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		cost = bcrypt.DefaultCost
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPasswordHash بررسی رمز عبور با هش
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// CheckPasswordHashWithDetail بررسی با جزئیات خطا
func CheckPasswordHashWithDetail(password, hash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, errors.New("invalid password")
		}
		return false, err
	}
	return true, nil
}

// ============================================================================
// بخش 2: مثال‌های پایه
// ============================================================================

func demonstrateBasicUsage() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📝 BASIC BCRYPT USAGE")
	fmt.Println(stringsRepeat("=", 80))

	password := "mySecurePassword123"

	// هش کردن رمز
	hash, err := HashPassword(password)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Original password: %s\n", password)
	fmt.Printf("Hashed password: %s\n", hash)
	fmt.Printf("Hash length: %d characters\n", len(hash))

	// بررسی رمز درست
	match := CheckPasswordHash(password, hash)
	fmt.Printf("Password match: %v\n", match)

	// بررسی رمز غلط
	wrongPassword := "wrongPassword"
	matchWrong := CheckPasswordHash(wrongPassword, hash)
	fmt.Printf("Wrong password match: %v\n", matchWrong)
}

// ============================================================================
// بخش 3: Cost Factor - امنیت vs عملکرد
// ============================================================================

func demonstrateCostFactor() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("⚙️ COST FACTOR DEMONSTRATION")
	fmt.Println(stringsRepeat("=", 80))

	password := "testPassword"

	costs := []int{4, 6, 8, 10, 12}

	fmt.Println("\nCost | Time (ms) | Hash (truncated)")
	fmt.Println("-----|-----------|------------------")

	for _, cost := range costs {
		start := time.Now()
		hash, err := HashPasswordWithCost(password, cost)
		duration := time.Since(start)

		if err != nil {
			fmt.Printf("%4d | error     | %v\n", cost, err)
			continue
		}

		hashPreview := hash[:29] + "..."
		fmt.Printf("%4d | %9d | %s\n", cost, duration.Milliseconds(), hashPreview)
	}

	fmt.Println("\n💡 Recommendation:")
	fmt.Println("   • Development/Test: cost 4-6 (faster)")
	fmt.Println("   • Production: cost 10-12 (more secure)")
	fmt.Println("   • High security: cost 14+ (slow but secure)")
}

// ============================================================================
// بخش 4: بررسی و Upgrade هش
// ============================================================================

// NeedsRehash بررسی آیا هش نیاز به به‌روزرسانی دارد
func NeedsRehash(hash string, targetCost int) bool {
	cost, err := bcrypt.Cost([]byte(hash))
	if err != nil {
		return true // هش نامعتبر، نیاز به rehash
	}
	return cost < targetCost
}

// RehashIfNeeded به‌روزرسانی هش در صورت نیاز
func RehashIfNeeded(password, currentHash string, targetCost int) (string, bool, error) {
	// بررسی صحت رمز
	if !CheckPasswordHash(password, currentHash) {
		return "", false, errors.New("invalid password")
	}

	// بررسی نیاز به rehash
	if !NeedsRehash(currentHash, targetCost) {
		return currentHash, false, nil
	}

	// هش مجدد با cost جدید
	newHash, err := HashPasswordWithCost(password, targetCost)
	if err != nil {
		return "", false, err
	}

	return newHash, true, nil
}

func demonstrateRehashing() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🔄 PASSWORD REHASHING")
	fmt.Println(stringsRepeat("=", 80))

	password := "myPassword"
	targetCost := 12

	// هش اولیه با cost پایین
	oldHash, _ := HashPasswordWithCost(password, 8)
	fmt.Printf("Old hash (cost=8): %s\n", oldHash[:40])

	// بررسی نیاز به rehash
	if NeedsRehash(oldHash, targetCost) {
		fmt.Printf("Hash needs upgrade from cost 8 to %d\n", targetCost)

		newHash, upgraded, err := RehashIfNeeded(password, oldHash, targetCost)
		if err != nil {
			log.Printf("Rehash error: %v", err)
		}

		if upgraded {
			fmt.Printf("New hash (cost=%d): %s\n", targetCost, newHash[:40])
		}
	} else {
		fmt.Println("Hash is up to date")
	}
}

// ============================================================================
// بخش 5: کار با دیتابیس (مثال عملی)
// ============================================================================

// User مدل کاربر در دیتابیس
type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // هش شده
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserRepository اینترفیس مخزن کاربران
type UserRepository interface {
	CreateUser(user *User) error
	GetUserByEmail(email string) (*User, error)
	UpdatePassword(userID int, hashedPassword string) error
}

// InMemoryUserRepository پیاده‌سازی در حافظه
type InMemoryUserRepository struct {
	users    map[int]*User
	emailMap map[string]int
	nextID   int
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users:    make(map[int]*User),
		emailMap: make(map[string]int),
		nextID:   1,
	}
}

func (r *InMemoryUserRepository) CreateUser(user *User) error {
	user.ID = r.nextID
	r.nextID++
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	r.users[user.ID] = user
	r.emailMap[user.Email] = user.ID
	return nil
}

func (r *InMemoryUserRepository) GetUserByEmail(email string) (*User, error) {
	id, exists := r.emailMap[email]
	if !exists {
		return nil, errors.New("user not found")
	}
	return r.users[id], nil
}

func (r *InMemoryUserRepository) UpdatePassword(userID int, hashedPassword string) error {
	user, exists := r.users[userID]
	if !exists {
		return errors.New("user not found")
	}
	user.Password = hashedPassword
	user.UpdatedAt = time.Now()
	return nil
}

// AuthService سرویس احراز هویت
type AuthService struct {
	userRepo UserRepository
	cost     int
}

// NewAuthService ایجاد سرویس احراز هویت
func NewAuthService(userRepo UserRepository, cost int) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		cost:     cost,
	}
}

// RegisterRequest درخواست ثبت‌نام
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Register ثبت‌نام کاربر جدید
func (s *AuthService) Register(req RegisterRequest) (*User, error) {
	// اعتبارسنجی رمز عبور
	if len(req.Password) < 8 {
		return nil, errors.New("password must be at least 8 characters")
	}

	// هش کردن رمز
	hashedPassword, err := HashPasswordWithCost(req.Password, s.cost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	if err := s.userRepo.CreateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

// LoginRequest درخواست ورود
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login ورود کاربر
func (s *AuthService) Login(req LoginRequest) (*User, error) {
	// دریافت کاربر از دیتابیس
	user, err := s.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// بررسی رمز عبور
	if !CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	// بررسی نیاز به rehash (اختیاری)
	if NeedsRehash(user.Password, s.cost) {
		newHash, err := HashPasswordWithCost(req.Password, s.cost)
		if err == nil {
			s.userRepo.UpdatePassword(user.ID, newHash)
		}
	}

	return user, nil
}

// ============================================================================
// بخش 6: HTTP Handlers (مثال API)
// ============================================================================

type AuthHandler struct {
	authService *AuthService
}

func NewAuthHandler(authService *AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user, err := h.authService.Register(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// حذف رمز از پاسخ
	user.Password = ""

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user, err := h.authService.Login(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// حذف رمز از پاسخ
	user.Password = ""

	// در عمل، اینجا JWT صادر می‌شود
	response := map[string]interface{}{
		"user":  user,
		"token": "jwt-token-here", // در واقع JWT
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ============================================================================
// بخش 7: Constant Time Comparison (امنیت در برابر timing attacks)
// ============================================================================

// ConstantTimeCompare مقایسه امن در برابر timing attacks
func ConstantTimeCompare(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

// SecurePasswordVerify بررسی امن رمز عبور
func SecurePasswordVerify(password, hash string) bool {
	// bcrypt.CompareHashAndPassword درون خود constant time است
	// اما برای مثال، اینجا نحوه استفاده از subtle را نشان می‌دهیم
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func demonstrateConstantTime() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🔒 CONSTANT TIME COMPARISON")
	fmt.Println(stringsRepeat("=", 80))

	password := "secret123"
	hash, _ := HashPassword(password)

	// مقایسه معمولی (آسیب‌پذیر به timing attack)
	normalCompare := password == "secret123"

	// مقایسه constant time (امن)
	secureCompare := ConstantTimeCompare(password, "secret123")

	// بررسی رمز با bcrypt (constant time internally)
	bcryptVerify := SecurePasswordVerify(password, hash)

	fmt.Printf("Normal compare: %v\n", normalCompare)
	fmt.Printf("Constant time compare: %v\n", secureCompare)
	fmt.Printf("bcrypt verify: %v\n", bcryptVerify)
	fmt.Println("\n💡 bcrypt's CompareHashAndPassword is already constant-time")
}

// ============================================================================
// بخش 8: Bulk Password Hashing (برای import داده)
// ============================================================================

// BulkHashResult نتیجه هش دسته‌ای
type BulkHashResult struct {
	Original string
	Hash     string
	Error    error
}

// HashPasswordsBulk هش کردن دسته‌ای رمزهای عبور
func HashPasswordsBulk(passwords []string, cost int) []BulkHashResult {
	results := make([]BulkHashResult, len(passwords))

	for i, pwd := range passwords {
		hash, err := HashPasswordWithCost(pwd, cost)
		results[i] = BulkHashResult{
			Original: pwd,
			Hash:     hash,
			Error:    err,
		}
	}

	return results
}

// HashPasswordsParallel هش کردن موازی (برای تعداد زیاد)
func HashPasswordsParallel(passwords []string, cost int, workers int) []BulkHashResult {
	results := make([]BulkHashResult, len(passwords))
	jobs := make(chan int, len(passwords))
	resultsChan := make(chan BulkHashResult, len(passwords))

	// شروع workers
	for w := 0; w < workers; w++ {
		go func() {
			for idx := range jobs {
				hash, err := HashPasswordWithCost(passwords[idx], cost)
				resultsChan <- BulkHashResult{
					Original: passwords[idx],
					Hash:     hash,
					Error:    err,
				}
			}
		}()
	}

	// ارسال jobs
	for i := range passwords {
		jobs <- i
	}
	close(jobs)

	// جمع‌آوری نتایج
	for i := 0; i < len(passwords); i++ {
		result := <-resultsChan
		// پیدا کردن index اصلی (در عمل بهتر است index هم ارسال شود)
		for j, p := range passwords {
			if p == result.Original {
				results[j] = result
				break
			}
		}
	}

	return results
}

func demonstrateBulkHashing() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📦 BULK PASSWORD HASHING")
	fmt.Println(stringsRepeat("=", 80))

	passwords := []string{
		"password123",
		"mySecurePass",
		"admin123",
		"user456",
	}

	// هش کردن دسته‌ای
	results := HashPasswordsBulk(passwords, 10)

	fmt.Println("Bulk hash results:")
	for _, r := range results {
		if r.Error != nil {
			fmt.Printf("  %s: ERROR - %v\n", r.Original, r.Error)
		} else {
			fmt.Printf("  %s: %s...\n", r.Original, r.Hash[:30])
		}
	}
}

// ============================================================================
// بخش 9: Validation و Security Checks
// ============================================================================

// PasswordValidator اعتبارسنجی رمز عبور
type PasswordValidator struct {
	MinLength        int
	RequireUppercase bool
	RequireLowercase bool
	RequireNumber    bool
	RequireSpecial   bool
}

// DefaultPasswordValidator validator پیش‌فرض
func DefaultPasswordValidator() *PasswordValidator {
	return &PasswordValidator{
		MinLength:        8,
		RequireUppercase: true,
		RequireLowercase: true,
		RequireNumber:    true,
		RequireSpecial:   false,
	}
}

// Validate اعتبارسنجی رمز عبور
func (v *PasswordValidator) Validate(password string) error {
	if len(password) < v.MinLength {
		return fmt.Errorf("password must be at least %d characters", v.MinLength)
	}

	hasUpper := false
	hasLower := false
	hasNumber := false
	hasSpecial := false

	for _, c := range password {
		switch {
		case c >= 'A' && c <= 'Z':
			hasUpper = true
		case c >= 'a' && c <= 'z':
			hasLower = true
		case c >= '0' && c <= '9':
			hasNumber = true
		case c >= '!' && c <= '/' || c >= ':' && c <= '@' || c >= '[' && c <= '`' || c >= '{' && c <= '~':
			hasSpecial = true
		}
	}

	if v.RequireUppercase && !hasUpper {
		return errors.New("password must contain at least one uppercase letter")
	}
	if v.RequireLowercase && !hasLower {
		return errors.New("password must contain at least one lowercase letter")
	}
	if v.RequireNumber && !hasNumber {
		return errors.New("password must contain at least one number")
	}
	if v.RequireSpecial && !hasSpecial {
		return errors.New("password must contain at least one special character")
	}

	return nil
}

// CommonPasswordCheck بررسی رمزهای رایج
func CommonPasswordCheck(password string) bool {
	// لیست رمزهای رایج (در عمل از فایل یا دیتابیس بخوانید)
	commonPasswords := map[string]bool{
		"password":    true,
		"123456":      true,
		"qwerty":      true,
		"admin":       true,
		"welcome":     true,
		"password123": true,
	}

	return commonPasswords[password]
}

func demonstrateValidation() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("✅ PASSWORD VALIDATION")
	fmt.Println(stringsRepeat("=", 80))

	validator := DefaultPasswordValidator()

	testPasswords := []string{
		"weak",
		"nouppercase123",
		"NOLOWER123",
		"NoNumbers",
		"ValidPass123",
	}

	for _, pwd := range testPasswords {
		err := validator.Validate(pwd)
		if err != nil {
			fmt.Printf("  %s: ❌ %v\n", pwd, err)
		} else {
			fmt.Printf("  %s: ✅ valid\n", pwd)
		}
	}

	// بررسی رمزهای رایج
	fmt.Println("\nCommon password check:")
	common := "password123"
	if CommonPasswordCheck(common) {
		fmt.Printf("  %s: ❌ is a common password, please choose a stronger one\n", common)
	}
}

// ============================================================================
// بخش 10: Best Practices
// ============================================================================

func bestPractices() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("💡 BCRYPT BEST PRACTICES")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ 1. CHOOSE THE RIGHT COST                                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│    • Cost 10: Good default for most applications                          │
│    • Cost 12-14: High security (banks, healthcare)                        │
│    • Cost 4-8: Testing/development only                                   │
│    • Benchmark on your hardware: aim for 100-200ms per hash               │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 2. PASSWORD VALIDATION                                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│    • Minimum length: 8-12 characters                                      │
│    • Require mixed case, numbers, specials                                │
│    • Check against common password lists                                  │
│    • Never validate on client side only                                   │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 3. STORAGE                                                                │
├─────────────────────────────────────────────────────────────────────────────┤
│    • Store only the hash, never the plain password                        │
│    • Use a separate column for hash (length: 60-72 chars)                 │
│    • Consider storing the cost factor used                                │
│    • Implement password upgrade mechanism                                 │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 4. SECURITY                                                               │
├─────────────────────────────────────────────────────────────────────────────┤
│    • Use constant-time comparison (bcrypt does this)                      │
│    • Implement rate limiting on login endpoints                           │
│    • Use account lockout after failed attempts                            │
│    • Require re-authentication for sensitive actions                      │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 5. MIGRATION                                                              │
├─────────────────────────────────────────────────────────────────────────────┤
│    • Rehash passwords on successful login                                 │
│    • Support multiple hash formats during migration                       │
│    • Gradually increase cost factor                                       │
│    • Never downgrade hash security                                        │
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
│ MISTAKE 1: Using weak hash algorithms                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│    ❌ hash := sha256.Sum256([]byte(password))                             │
│    ✅ hash, _ := bcrypt.GenerateFromPassword([]byte(password), cost)      │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ MISTAKE 2: Not using salt                                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│    ❌ hash := hash(password) // no salt                                   │
│    ✅ bcrypt includes salt automatically                                  │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ MISTAKE 3: Too low cost factor                                            │
├─────────────────────────────────────────────────────────────────────────────┤
│    ❌ hash, _ := bcrypt.GenerateFromPassword(pwd, 4) // too fast          │
│    ✅ Use cost >= 10 for production                                       │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ MISTAKE 4: Storing passwords in logs                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│    ❌ log.Printf("User login: %s, password: %s", email, password)         │
│    ✅ Never log passwords, not even hashes                                │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ MISTAKE 5: Not handling bcrypt errors                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│    ❌ hash, _ := bcrypt.GenerateFromPassword(pwd, cost)                   │
│    ✅ Always check and handle errors                                      │
└─────────────────────────────────────────────────────────────────────────────┘
`)
}

// ============================================================================
// بخش 12: Command Line Tool
// ============================================================================

func runCLI() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🔧 BCRYPT COMMAND LINE TOOL")
	fmt.Println(stringsRepeat("=", 80))

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter password: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	if password == "" {
		fmt.Println("No password entered")
		return
	}

	// هش کردن
	hash, err := HashPassword(password)
	if err != nil {
		fmt.Printf("Error hashing: %v\n", err)
		return
	}

	fmt.Printf("\nHashed password: %s\n", hash)
	fmt.Printf("Hash length: %d characters\n", len(hash))

	// verify
	fmt.Print("\nVerify password (enter again): ")
	verify, _ := reader.ReadString('\n')
	verify = strings.TrimSpace(verify)

	if CheckPasswordHash(verify, hash) {
		fmt.Println("✅ Password matches!")
	} else {
		fmt.Println("❌ Password does NOT match!")
	}
}

// ============================================================================
// بخش 13: Main
// ============================================================================

func main() {
	fmt.Println(stringsRepeat("=", 80))
	fmt.Println("🎯 COMPLETE BCRYPT GUIDE")
	fmt.Println("Password Hashing in Go")
	fmt.Println(stringsRepeat("=", 80))

	demonstrateBasicUsage()
	demonstrateCostFactor()
	demonstrateRehashing()
	demonstrateConstantTime()
	demonstrateBulkHashing()
	demonstrateValidation()
	bestPractices()
	commonMistakes()

	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🚀 BCRYPT COMMAND LINE EXAMPLE")
	fmt.Println(stringsRepeat("=", 80))
	runCLI()

	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🎯 BCRYPT - COMPLETE")
	fmt.Println("Secure password hashing for Go applications")
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
