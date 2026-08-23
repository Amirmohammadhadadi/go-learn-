// ============================================================================
// FILE: jwt_guide.go
// TITLE: راهنمای کامل JWT در Go - ذخیره در HTTP-Only Cookie یا Header
// HOW TO RUN: go run jwt_guide.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - JWT چیست و چرا نیاز است؟
// ============================================================================
//
// JWT (JSON Web Token) یک استاندارد باز برای انتقال امن اطلاعات بین parties است.
//
// ساختار JWT: Header.Payload.Signature
// - Header: نوع token و الگوریتم (مثل HS256, RS256)
// - Payload: Claims (داده‌های کاربر، مثل user_id, role, exp)
// - Signature: امضای دیجیتال برای تأیید authenticity
//
// دو روش اصلی ذخیره JWT:
//
// 1. HTTP-Only Cookie (توصیه شده برای SPA و وب‌اپلیکیشن‌ها)
//    - امنیت بالاتر (دسترسی ناپذیر توسط JavaScript)
//    - محافظت در برابر XSS
//    - خودکار در هر درخواست ارسال می‌شود
//    - نیاز به CSRF protection دارد
//
// 2. Authorization Header (توصیه شده برای APIها)
//    - ساده‌تر برای کلاینت‌های موبایل و API
//    - بدون مشکل CSRF
//    - نیاز به مدیریت manual ارسال token
//
// قانون طلایی:
// "برای برنامه‌های وب از HTTP-Only Cookie استفاده کن (امن‌تر در برابر XSS).
//  برای APIهای عمومی از Authorization Header استفاده کن (ساده‌تر و stateless).
//  همیشه از HTTPS استفاده کن و tokenها را کوتاه‌مدت نگه دار."
// ============================================================================

package __jwt

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// ============================================================================
// بخش 1: Constants و Types
// ============================================================================

const (
	// Cookie names
	AccessTokenCookieName  = "access_token"
	RefreshTokenCookieName = "refresh_token"

	// Header names
	AuthorizationHeader = "Authorization"
	BearerPrefix        = "Bearer "

	// Token durations
	AccessTokenDuration  = 15 * time.Minute
	RefreshTokenDuration = 7 * 24 * time.Hour

	// Context keys
	UserContextKey = "user"
)

// CustomClaims ساختار claims سفارشی
type CustomClaims struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	DeviceID string `json:"device_id,omitempty"`
	jwt.RegisteredClaims
}

// User مدل کاربر
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

// TokenPair جفت توکن‌ها
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

// TokenResponse پاسخ توکن
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}

// ============================================================================
// بخش 2: JWT Service
// ============================================================================

// JWTService سرویس مدیریت JWT
type JWTService struct {
	secretKey       []byte
	accessDuration  time.Duration
	refreshDuration time.Duration
	issuer          string
}

// NewJWTService ایجاد سرویس JWT جدید
func NewJWTService(secretKey string, accessDuration, refreshDuration time.Duration, issuer string) *JWTService {
	return &JWTService{
		secretKey:       []byte(secretKey),
		accessDuration:  accessDuration,
		refreshDuration: refreshDuration,
		issuer:          issuer,
	}
}

// GenerateTokenPair تولید جفت توکن (access + refresh)
func (s *JWTService) GenerateTokenPair(userID, email, role, deviceID string) (*TokenPair, error) {
	// Access Token
	accessToken, err := s.generateToken(userID, email, role, deviceID, s.accessDuration)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Refresh Token
	refreshToken, err := s.generateToken(userID, email, role, deviceID, s.refreshDuration)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.accessDuration.Seconds()),
		TokenType:    "Bearer",
	}, nil
}

// generateToken تولید یک token
func (s *JWTService) generateToken(userID, email, role, deviceID string, duration time.Duration) (string, error) {
	claims := &CustomClaims{
		UserID:   userID,
		Email:    email,
		Role:     role,
		DeviceID: deviceID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    s.issuer,
			Subject:   userID,
			ID:        generateTokenID(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

// VerifyToken اعتبارسنجی token
func (s *JWTService) VerifyToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// بررسی الگوریتم
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// RefreshToken تولید توکن جدید از refresh token
func (s *JWTService) RefreshToken(refreshTokenString string) (*TokenPair, error) {
	claims, err := s.VerifyToken(refreshTokenString)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// تولید جفت توکن جدید
	return s.GenerateTokenPair(claims.UserID, claims.Email, claims.Role, claims.DeviceID)
}

// ============================================================================
// بخش 3: Token Storage - HTTP-Only Cookie (روش اول)
// ============================================================================

// CookieTokenManager مدیریت توکن‌ها با Cookie
type CookieTokenManager struct {
	jwtService *JWTService
	secure     bool // true در production (HTTPS)
	httpOnly   bool // true (دسترسی ناپذیر توسط JavaScript)
	sameSite   http.SameSite
	domain     string
	path       string
}

// NewCookieTokenManager ایجاد manager کوکی
func NewCookieTokenManager(jwtService *JWTService, secure bool, domain string) *CookieTokenManager {
	return &CookieTokenManager{
		jwtService: jwtService,
		secure:     secure,
		httpOnly:   true,
		sameSite:   http.SameSiteStrictMode,
		domain:     domain,
		path:       "/",
	}
}

// SetTokens ذخیره توکن‌ها در کوکی
func (m *CookieTokenManager) SetTokens(w http.ResponseWriter, accessToken, refreshToken string) {
	// Access Token Cookie (کوتاه‌مدت)
	m.setCookie(w, AccessTokenCookieName, accessToken, AccessTokenDuration)

	// Refresh Token Cookie (بلندمدت)
	m.setCookie(w, RefreshTokenCookieName, refreshToken, RefreshTokenDuration)
}

// setCookie تنظیم یک کوکی
func (m *CookieTokenManager) setCookie(w http.ResponseWriter, name, value string, maxAge time.Duration) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     m.path,
		Domain:   m.domain,
		MaxAge:   int(maxAge.Seconds()),
		HttpOnly: m.httpOnly,
		Secure:   m.secure,
		SameSite: m.sameSite,
	}
	http.SetCookie(w, cookie)
}

// GetTokens دریافت توکن‌ها از کوکی
func (m *CookieTokenManager) GetTokens(r *http.Request) (string, string, error) {
	accessCookie, err := r.Cookie(AccessTokenCookieName)
	if err != nil {
		return "", "", errors.New("access token not found")
	}

	refreshCookie, err := r.Cookie(RefreshTokenCookieName)
	if err != nil {
		return "", "", errors.New("refresh token not found")
	}

	return accessCookie.Value, refreshCookie.Value, nil
}

// ClearTokens حذف توکن‌ها از کوکی
func (m *CookieTokenManager) ClearTokens(w http.ResponseWriter) {
	m.setCookie(w, AccessTokenCookieName, "", -1)
	m.setCookie(w, RefreshTokenCookieName, "", -1)
}

// GetUserFromCookie دریافت کاربر از کوکی (middleware)
func (m *CookieTokenManager) GetUserFromCookie(r *http.Request) (*CustomClaims, error) {
	accessCookie, err := r.Cookie(AccessTokenCookieName)
	if err != nil {
		return nil, errors.New("no access token")
	}

	return m.jwtService.VerifyToken(accessCookie.Value)
}

// ============================================================================
// بخش 4: Token Storage - Authorization Header (روش دوم)
// ============================================================================

// HeaderTokenManager مدیریت توکن‌ها با Header
type HeaderTokenManager struct {
	jwtService *JWTService
}

// NewHeaderTokenManager ایجاد manager هدر
func NewHeaderTokenManager(jwtService *JWTService) *HeaderTokenManager {
	return &HeaderTokenManager{
		jwtService: jwtService,
	}
}

// ExtractToken استخراج token از هدر
func (m *HeaderTokenManager) ExtractToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get(AuthorizationHeader)
	if authHeader == "" {
		return "", errors.New("authorization header is missing")
	}

	// بررسی فرمت Bearer
	if !strings.HasPrefix(authHeader, BearerPrefix) {
		return "", errors.New("invalid authorization format")
	}

	token := strings.TrimPrefix(authHeader, BearerPrefix)
	if token == "" {
		return "", errors.New("token is empty")
	}

	return token, nil
}

// GetUserFromHeader دریافت کاربر از هدر (middleware)
func (m *HeaderTokenManager) GetUserFromHeader(r *http.Request) (*CustomClaims, error) {
	token, err := m.ExtractToken(r)
	if err != nil {
		return nil, err
	}

	return m.jwtService.VerifyToken(token)
}

// SetTokenResponse تنظیم پاسخ با توکن (برای API)
func (m *HeaderTokenManager) SetTokenResponse(w http.ResponseWriter, accessToken, refreshToken string) {
	response := TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(AccessTokenDuration.Seconds()),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// ============================================================================
// بخش 5: Authentication Middleware (یکپارچه)
// ============================================================================

// AuthMiddleware میدلور احراز هویت (پشتیبانی از هر دو روش)
type AuthMiddleware struct {
	cookieManager *CookieTokenManager
	headerManager *HeaderTokenManager
}

// NewAuthMiddleware ایجاد میدلور احراز هویت
func NewAuthMiddleware(jwtService *JWTService, secure bool, domain string) *AuthMiddleware {
	return &AuthMiddleware{
		cookieManager: NewCookieTokenManager(jwtService, secure, domain),
		headerManager: NewHeaderTokenManager(jwtService),
	}
}

// Authenticate میدلور اصلی - تلاش هر دو روش
func (m *AuthMiddleware) Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var claims *CustomClaims
		var err error

		// روش اول: بررسی Cookie
		claims, err = m.cookieManager.GetUserFromCookie(r)
		if err == nil {
			// ذخیره user در context
			ctx := context.WithValue(r.Context(), UserContextKey, claims)
			next(w, r.WithContext(ctx))
			return
		}

		// روش دوم: بررسی Header
		claims, err = m.headerManager.GetUserFromHeader(r)
		if err == nil {
			ctx := context.WithValue(r.Context(), UserContextKey, claims)
			next(w, r.WithContext(ctx))
			return
		}

		// هر دو روش شکست خوردند
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}

// RequireRole میدلور بررسی نقش کاربر
func RequireRole(allowedRoles ...string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(UserContextKey).(*CustomClaims)
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			for _, role := range allowedRoles {
				if claims.Role == role {
					next(w, r)
					return
				}
			}

			http.Error(w, "Forbidden", http.StatusForbidden)
		}
	}
}

// ============================================================================
// بخش 6: Auth Handlers (Login, Logout, Refresh)
// ============================================================================

// AuthHandler مدیریت احراز هویت
type AuthHandler struct {
	jwtService     *JWTService
	cookieManager  *CookieTokenManager
	headerManager  *HeaderTokenManager
	userRepository UserRepository // فرضی
}

// LoginRequest درخواست لاگین
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	DeviceID string `json:"device_id,omitempty"`
}

// LoginHandler (با Cookie - برای وب)
func (h *AuthHandler) LoginWithCookie(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// اعتبارسنجی کاربر (در واقع از دیتابیس)
	user, err := h.userRepository.FindByEmail(r.Context(), req.Email)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// بررسی رمز عبور
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// تولید توکن
	tokenPair, err := h.jwtService.GenerateTokenPair(user.ID, user.Email, user.Role, req.DeviceID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// ذخیره در کوکی
	h.cookieManager.SetTokens(w, tokenPair.AccessToken, tokenPair.RefreshToken)

	// پاسخ موفقیت
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Login successful",
		"user": map[string]interface{}{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
			"role":  user.Role,
		},
	})
}

// LoginWithHeader (برای API)
func (h *AuthHandler) LoginWithHeader(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// اعتبارسنجی کاربر
	user, err := h.userRepository.FindByEmail(r.Context(), req.Email)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// تولید توکن
	tokenPair, err := h.jwtService.GenerateTokenPair(user.ID, user.Email, user.Role, req.DeviceID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// پاسخ با توکن در body
	h.headerManager.SetTokenResponse(w, tokenPair.AccessToken, tokenPair.RefreshToken)
}

// RefreshHandler تمدید توکن
func (h *AuthHandler) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	var refreshToken string

	// تلاش از کوکی
	if cookie, err := r.Cookie(RefreshTokenCookieName); err == nil {
		refreshToken = cookie.Value
	} else {
		// تلاش از هدر
		authHeader := r.Header.Get(AuthorizationHeader)
		if strings.HasPrefix(authHeader, BearerPrefix) {
			refreshToken = strings.TrimPrefix(authHeader, BearerPrefix)
		}
	}

	if refreshToken == "" {
		http.Error(w, "Refresh token required", http.StatusBadRequest)
		return
	}

	// تمدید توکن
	newTokenPair, err := h.jwtService.RefreshToken(refreshToken)
	if err != nil {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	// ذخیره در کوکی (اگر از کوکی استفاده می‌شود)
	if _, err := r.Cookie(AccessTokenCookieName); err == nil {
		h.cookieManager.SetTokens(w, newTokenPair.AccessToken, newTokenPair.RefreshToken)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "refreshed"})
		return
	}

	// پاسخ با توکن در body (برای API)
	h.headerManager.SetTokenResponse(w, newTokenPair.AccessToken, newTokenPair.RefreshToken)
}

// LogoutHandler خروج از سیستم
func (h *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// حذف کوکی‌ها (اگر وجود داشته باشند)
	h.cookieManager.ClearTokens(w)

	// پاسخ
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Logged out successfully",
	})
}

// ============================================================================
// بخش 7: Example User Repository (فرضی)
// ============================================================================

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id string) (*User, error)
	Create(ctx context.Context, user *User) error
}

type InMemoryUserRepository struct {
	users map[string]*User
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	// ایجاد کاربر نمونه
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	return &InMemoryUserRepository{
		users: map[string]*User{
			"1": {
				ID:        "1",
				Email:     "user@example.com",
				Password:  string(hashedPassword),
				Name:      "Test User",
				Role:      "user",
				CreatedAt: time.Now(),
			},
			"2": {
				ID:        "2",
				Email:     "admin@example.com",
				Password:  string(hashedPassword),
				Name:      "Admin User",
				Role:      "admin",
				CreatedAt: time.Now(),
			},
		},
	}
}

func (r *InMemoryUserRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (r *InMemoryUserRepository) FindByID(ctx context.Context, id string) (*User, error) {
	user, exists := r.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (r *InMemoryUserRepository) Create(ctx context.Context, user *User) error {
	r.users[user.ID] = user
	return nil
}

// ============================================================================
// بخش 8: Protected Handlers (مثال)
// ============================================================================

func profileHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(UserContextKey).(*CustomClaims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": claims.UserID,
		"email":   claims.Email,
		"role":    claims.Role,
		"exp":     claims.ExpiresAt,
	})
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Welcome to admin panel",
	})
}

// ============================================================================
// بخش 9: Setup Server (مثال کامل)
// ============================================================================

func setupServer() {
	// JWT Secret (در production از environment variable استفاده کن)
	jwtSecret := "your-256-bit-secret-key-here-please-change-me"

	// ایجاد JWT Service
	jwtService := NewJWTService(
		jwtSecret,
		AccessTokenDuration,
		RefreshTokenDuration,
		"myapp",
	)

	// ایجاد User Repository
	userRepo := NewInMemoryUserRepository()

	// ایجاد Auth Handler
	authHandler := &AuthHandler{
		jwtService:     jwtService,
		cookieManager:  NewCookieTokenManager(jwtService, false, ""),
		headerManager:  NewHeaderTokenManager(jwtService),
		userRepository: userRepo,
	}

	// ایجاد Auth Middleware
	authMiddleware := NewAuthMiddleware(jwtService, false, "")

	// Router
	mux := http.NewServeMux()

	// Public routes (بدون احراز هویت)
	mux.HandleFunc("POST /api/login", authHandler.LoginWithCookie)
	mux.HandleFunc("POST /api/login/header", authHandler.LoginWithHeader)
	mux.HandleFunc("POST /api/refresh", authHandler.RefreshHandler)
	mux.HandleFunc("POST /api/logout", authHandler.LogoutHandler)

	// Protected routes (نیاز به احراز هویت)
	mux.HandleFunc("GET /api/profile", authMiddleware.Authenticate(profileHandler))
	mux.HandleFunc("GET /api/admin", authMiddleware.Authenticate(RequireRole("admin")(adminHandler)))

	// Start server
	log.Println("Server starting on :8080")
	log.Println("Endpoints:")
	log.Println("  POST   /api/login          - Login with cookie (web)")
	log.Println("  POST   /api/login/header   - Login with header (API)")
	log.Println("  POST   /api/refresh        - Refresh token")
	log.Println("  POST   /api/logout         - Logout")
	log.Println("  GET    /api/profile        - Get profile (requires auth)")
	log.Println("  GET    /api/admin          - Admin only (requires admin role)")

	log.Fatal(http.ListenAndServe(":8080", mux))
}

// ============================================================================
// بخش 10: Helper Functions
// ============================================================================

func generateTokenID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// ============================================================================
// بخش 11: Best Practices
// ============================================================================

func bestPractices() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("💡 JWT BEST PRACTICES")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ 1. TOKEN LIFETIME                                                         │
├─────────────────────────────────────────────────────────────────────────────┤
│    • Access Token: 15 minutes (short-lived)                               │
│    • Refresh Token: 7-30 days (long-lived)                                │
│    • Rotate refresh tokens on each use                                    │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 2. STORAGE SECURITY                                                       │
├─────────────────────────────────────────────────────────────────────────────┤
│    • Web Apps: Use HTTP-Only Cookies (secure, httpOnly, sameSite)         │
│    • Mobile/API: Use Authorization Header                                 │
│    • Always use HTTPS (secure flag for cookies)                           │
│    • Never store tokens in localStorage (XSS vulnerable)                  │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 3. TOKEN CONTENTS                                                         │
├─────────────────────────────────────────────────────────────────────────────┤
│    • Minimal claims (user_id, role, exp)                                  │
│    • Don't store sensitive data (password, SSN, etc.)                     │
│    • Include jti (token ID) for revocation                                │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 4. SIGNING ALGORITHMS                                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│    • HS256: Good for single service                                       │
│    • RS256: Better for microservices                                      │
│    • ES256: Best for high security                                        │
│    • Never use none algorithm                                             │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 5. SECURITY CONSIDERATIONS                                                │
├─────────────────────────────────────────────────────────────────────────────┤
│    • Implement token revocation (blacklist)                               │
│    • Use refresh token rotation                                           │
│    • Implement logout (clear cookies/blacklist)                           │
│    • CSRF protection for cookie-based auth                                │
│    • Rate limiting on auth endpoints                                      │
└─────────────────────────────────────────────────────────────────────────────┘
`)
}

// ============================================================================
// بخش 12: Comparison Table
// ============================================================================

func comparisonTable() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📊 COOKIE vs HEADER COMPARISON")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ FEATURE              │ HTTP-Only Cookie    │ Authorization Header          │
├──────────────────────┼─────────────────────┼───────────────────────────────┤
│ XSS Protection       │ ✅ Yes              │ ❌ No (JS can access)         │
│ CSRF Protection      │ ❌ Required         │ ✅ Not needed                 │
│ Automatic sending    │ ✅ Yes              │ ❌ Manual                     │
│ Mobile API           │ ⚠️ Difficult       │ ✅ Easy                       │
│ Stateless            │ ✅ Yes              │ ✅ Yes                        │
│ Token size limit     │ 4KB                 │ None (practically)            │
│ Browser support      │ ✅ All              │ ✅ All                        │
│ Logout               │ Easy (clear cookie)│ Hard (token revocation)       │
│ Cross-domain         │ ⚠️ Limited         │ ✅ Easy                       │
│ RECOMMENDED FOR      │ Web Apps (SPA)      │ Mobile/Public APIs            │
└──────────────────────┴─────────────────────┴───────────────────────────────┘
`)
}

// ============================================================================
// بخش 13: Main
// ============================================================================

func main() {
	fmt.Println(stringsRepeat("=", 80))
	fmt.Println("🎯 COMPLETE JWT GUIDE")
	fmt.Println("HTTP-Only Cookie vs Authorization Header")
	fmt.Println(stringsRepeat("=", 80))

	bestPractices()
	comparisonTable()

	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🚀 Starting JWT example server on :8080")
	fmt.Println(stringsRepeat("=", 80))

	setupServer()
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
# Login with Cookie
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}' \
  -c cookies.txt

# Access protected endpoint with Cookie
curl http://localhost:8080/api/profile -b cookies.txt

# Login with Header
curl -X POST http://localhost:8080/api/login/header \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'

# Access with Header
curl http://localhost:8080/api/profile \
  -H "Authorization: Bearer <access_token>"
*/

/*
خلاصه JWT Storage Methods
ویژگی	HTTP-Only Cookie	Authorization Header
امنیت در برابر XSS	✅ عالی	❌ ضعیف
CSRF Protection	❌ نیاز دارد	✅ نیاز ندارد
ارسال خودکار	✅ بله	❌ دستی
موبایل/API	⚠️ سخت	✅ آسان
Logout	آسان (پاک کردن کوکی)	سخت (revocation)
حجم توکن	محدود (4KB)	بدون محدودیت
توصیه برای	وب اپلیکیشن‌ها	APIهای عمومی
*/
