// ============================================================================
// FILE: oauth2_guide.go
// TITLE: راهنمای کامل OAuth2 / OpenID Connect در Go - ورود با گوگل، گیت‌هاب
// HOW TO RUN: go run oauth2_guide.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - OAuth2 و OpenID Connect چیست؟
// ============================================================================
//
// OAuth2 یک پروتکل استاندارد برای授权 (authorization) است که به برنامه‌ها اجازه می‌دهد
// بدون دسترسی به رمز عبور کاربر، به منابع او دسترسی داشته باشند.
//
// OpenID Connect (OIDC) لایه‌ای روی OAuth2 است که احراز هویت (authentication) را اضافه می‌کند.
//
// مفاهیم کلیدی:
// - Client: برنامه شما
// - Resource Owner: کاربر
// - Authorization Server: سرویس دهنده (گوگل، گیت‌هاب)
// - Resource Server: API که داده‌ها را ارائه می‌دهد
// - Access Token: توکن دسترسی به منابع
// - Refresh Token: توکن برای دریافت Access Token جدید
// - ID Token: توکن حاوی اطلاعات کاربر (در OIDC)
//
// Flow اصلی (Authorization Code):
// 1. کاربر روی "ورود با گوگل" کلیک می‌کند
// 2. به صفحه گوگل هدایت می‌شود
// 3. پس از تأیید، به سایت شما برمی‌گردد با کد (code)
// 4. سرور شما کد را با Access Token تعویض می‌کند
// 5. با Access Token اطلاعات کاربر را دریافت می‌کند
//
// قانون طلایی:
// "هرگز Client Secret را در کلاینت ذخیره نکن (فقط در سرور).
//  همیشه state parameter را برای جلوگیری از CSRF استفاده کن.
//  از HTTPS در همه جای استفاده کن."
// ============================================================================

package __oauth2

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

// ============================================================================
// بخش 1: Configuration و State Management
// ============================================================================

// Config تنظیمات OAuth2
type OAuth2Config struct {
	GoogleConfig *oauth2.Config
	GitHubConfig *oauth2.Config
	SessionStore *sessions.CookieStore
	RedirectURL  string
}

// UserInfo اطلاعات کاربر از provider
type UserInfo struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
	Provider  string `json:"provider"`
}

// AppConfig تنظیمات برنامه
type AppConfig struct {
	GoogleClientID     string
	GoogleClientSecret string
	GitHubClientID     string
	GitHubClientSecret string
	SessionSecret      string
	RedirectURL        string
}

// LoadConfig بارگذاری تنظیمات از محیط
func LoadConfig() *AppConfig {
	return &AppConfig{
		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		GitHubClientID:     getEnv("GITHUB_CLIENT_ID", ""),
		GitHubClientSecret: getEnv("GITHUB_CLIENT_SECRET", ""),
		SessionSecret:      getEnv("SESSION_SECRET", "your-secret-key-change-me"),
		RedirectURL:        getEnv("REDIRECT_URL", "http://localhost:8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// NewOAuth2Config ایجاد تنظیمات OAuth2
func NewOAuth2Config(cfg *AppConfig) *OAuth2Config {
	return &OAuth2Config{
		GoogleConfig: &oauth2.Config{
			ClientID:     cfg.GoogleClientID,
			ClientSecret: cfg.GoogleClientSecret,
			RedirectURL:  cfg.RedirectURL + "/auth/google/callback",
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
		GitHubConfig: &oauth2.Config{
			ClientID:     cfg.GitHubClientID,
			ClientSecret: cfg.GitHubClientSecret,
			RedirectURL:  cfg.RedirectURL + "/auth/github/callback",
			Scopes: []string{
				"user:email",
				"read:user",
			},
			Endpoint: github.Endpoint,
		},
		SessionStore: sessions.NewCookieStore([]byte(cfg.SessionSecret)),
		RedirectURL:  cfg.RedirectURL,
	}
}

// ============================================================================
// بخش 2: State Management (برای جلوگیری از CSRF)
// ============================================================================

// generateStateToken تولید token تصادفی برای state
func generateStateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// storeState ذخیره state در session
func (c *OAuth2Config) storeState(r *http.Request, w http.ResponseWriter, state, provider string) error {
	session, _ := c.SessionStore.Get(r, "oauth-state")
	session.Values["state"] = state
	session.Values["provider"] = provider
	return session.Save(r, w)
}

// verifyState بررسی و حذف state
func (c *OAuth2Config) verifyState(r *http.Request, state string) (string, bool) {
	session, _ := c.SessionStore.Get(r, "oauth-state")
	storedState, ok := session.Values["state"].(string)
	if !ok || storedState != state {
		return "", false
	}
	provider, _ := session.Values["provider"].(string)

	// پاک کردن state بعد از استفاده
	delete(session.Values, "state")
	delete(session.Values, "provider")
	session.Save(r, nil)

	return provider, true
}

// ============================================================================
// بخش 3: User Repository (ذخیره کاربران)
// ============================================================================

// User مدل کاربر
type User struct {
	ID          string    `json:"id"`
	Email       string    `json:"email"`
	Name        string    `json:"name"`
	AvatarURL   string    `json:"avatar_url"`
	Provider    string    `json:"provider"`
	ProviderID  string    `json:"provider_id"`
	CreatedAt   time.Time `json:"created_at"`
	LastLoginAt time.Time `json:"last_login_at"`
}

// UserRepository مخزن کاربران
type UserRepository interface {
	FindOrCreate(userInfo *UserInfo) (*User, error)
	FindByID(id string) (*User, error)
}

// InMemoryUserRepository پیاده‌سازی در حافظه
type InMemoryUserRepository struct {
	users    map[string]*User
	emailMap map[string]string
	mu       chan struct{}
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users:    make(map[string]*User),
		emailMap: make(map[string]string),
		mu:       make(chan struct{}, 1),
	}
}

func (r *InMemoryUserRepository) FindOrCreate(userInfo *UserInfo) (*User, error) {
	r.mu <- struct{}{}
	defer func() { <-r.mu }()

	// بررسی وجود کاربر با email
	if id, exists := r.emailMap[userInfo.Email]; exists {
		user := r.users[id]
		user.LastLoginAt = time.Now()
		return user, nil
	}

	// ایجاد کاربر جدید
	user := &User{
		ID:          generateID(),
		Email:       userInfo.Email,
		Name:        userInfo.Name,
		AvatarURL:   userInfo.AvatarURL,
		Provider:    userInfo.Provider,
		ProviderID:  userInfo.ID,
		CreatedAt:   time.Now(),
		LastLoginAt: time.Now(),
	}

	r.users[user.ID] = user
	r.emailMap[user.Email] = user.ID

	return user, nil
}

func (r *InMemoryUserRepository) FindByID(id string) (*User, error) {
	r.mu <- struct{}{}
	defer func() { <-r.mu }()

	user, exists := r.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

// ============================================================================
// بخش 4: Provider Clients (Google, GitHub)
// ============================================================================

// GoogleUserInfo دریافت اطلاعات کاربر از Google
func getGoogleUserInfo(ctx context.Context, client *http.Client) (*UserInfo, error) {
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &UserInfo{
		ID:        data.ID,
		Email:     data.Email,
		Name:      data.Name,
		AvatarURL: data.Picture,
		Provider:  "google",
	}, nil
}

// GitHubUserInfo دریافت اطلاعات کاربر از GitHub
func getGitHubUserInfo(ctx context.Context, client *http.Client) (*UserInfo, error) {
	// دریافت اطلاعات کاربر
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userData struct {
		ID     int    `json:"id"`
		Login  string `json:"login"`
		Name   string `json:"name"`
		Email  string `json:"email"`
		Avatar string `json:"avatar_url"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&userData); err != nil {
		return nil, err
	}

	// اگر ایمیل عمومی نبود، از endpoint جداگانه بگیریم
	email := userData.Email
	if email == "" {
		emailResp, err := client.Get("https://api.github.com/user/emails")
		if err == nil {
			defer emailResp.Body.Close()
			var emails []struct {
				Email   string `json:"email"`
				Primary bool   `json:"primary"`
			}
			if json.NewDecoder(emailResp.Body).Decode(&emails) == nil {
				for _, e := range emails {
					if e.Primary {
						email = e.Email
						break
					}
				}
			}
		}
	}

	name := userData.Name
	if name == "" {
		name = userData.Login
	}

	return &UserInfo{
		ID:        fmt.Sprintf("%d", userData.ID),
		Email:     email,
		Name:      name,
		AvatarURL: userData.Avatar,
		Provider:  "github",
	}, nil
}

// ============================================================================
// بخش 5: Session Management
// ============================================================================

// SessionManager مدیریت session کاربر
type SessionManager struct {
	store *sessions.CookieStore
}

func NewSessionManager(secret string) *SessionManager {
	return &SessionManager{
		store: sessions.NewCookieStore([]byte(secret)),
	}
}

func (sm *SessionManager) SetUser(w http.ResponseWriter, r *http.Request, user *User) error {
	session, err := sm.store.Get(r, "user-session")
	if err != nil {
		return err
	}

	session.Values["user_id"] = user.ID
	session.Values["user_email"] = user.Email
	session.Values["user_name"] = user.Name
	session.Values["authenticated"] = true

	return session.Save(r, w)
}

func (sm *SessionManager) GetUser(r *http.Request) (*User, error) {
	session, err := sm.store.Get(r, "user-session")
	if err != nil {
		return nil, err
	}

	authenticated, ok := session.Values["authenticated"].(bool)
	if !ok || !authenticated {
		return nil, errors.New("not authenticated")
	}

	userID, ok := session.Values["user_id"].(string)
	if !ok {
		return nil, errors.New("user ID not found")
	}

	userEmail, _ := session.Values["user_email"].(string)
	userName, _ := session.Values["user_name"].(string)

	return &User{
		ID:    userID,
		Email: userEmail,
		Name:  userName,
	}, nil
}

func (sm *SessionManager) ClearUser(w http.ResponseWriter, r *http.Request) error {
	session, err := sm.store.Get(r, "user-session")
	if err != nil {
		return err
	}

	session.Values = make(map[interface{}]interface{})
	return session.Save(r, w)
}

// ============================================================================
// بخش 6: HTTP Handlers
// ============================================================================

// AuthHandler مدیریت احراز هویت
type AuthHandler struct {
	oauth2Config *OAuth2Config
	userRepo     UserRepository
	sessionMgr   *SessionManager
}

func NewAuthHandler(oauth2Config *OAuth2Config, userRepo UserRepository, sessionMgr *SessionManager) *AuthHandler {
	return &AuthHandler{
		oauth2Config: oauth2Config,
		userRepo:     userRepo,
		sessionMgr:   sessionMgr,
	}
}

// LoginPage صفحه ورود
func (h *AuthHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	tmpl := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Login - OAuth2 Demo</title>
		<style>
			body {
				font-family: Arial, sans-serif;
				display: flex;
				justify-content: center;
				align-items: center;
				height: 100vh;
				margin: 0;
				background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
			}
			.login-container {
				background: white;
				padding: 40px;
				border-radius: 10px;
				box-shadow: 0 4px 15px rgba(0,0,0,0.2);
				text-align: center;
				min-width: 300px;
			}
			h1 {
				margin-bottom: 30px;
				color: #333;
			}
			.btn {
				display: block;
				width: 100%;
				padding: 12px;
				margin: 10px 0;
				border: none;
				border-radius: 5px;
				font-size: 16px;
				cursor: pointer;
				text-decoration: none;
				text-align: center;
				color: white;
			}
			.btn-google {
				background: #DB4437;
			}
			.btn-google:hover {
				background: #c23321;
			}
			.btn-github {
				background: #333;
			}
			.btn-github:hover {
				background: #222;
			}
			.info {
				margin-top: 20px;
				font-size: 12px;
				color: #666;
			}
		</style>
	</head>
	<body>
		<div class="login-container">
			<h1>Login with</h1>
			<a href="/auth/google" class="btn btn-google">Continue with Google</a>
			<a href="/auth/github" class="btn btn-github">Continue with GitHub</a>
			<div class="info">
				Demo OAuth2 Application<br>
				Your data is not stored permanently
			</div>
		</div>
	</body>
	</html>
	`
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, tmpl)
}

// AuthRedirect هدایت به provider
func (h *AuthHandler) AuthRedirect(provider string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var config *oauth2.Config
		switch provider {
		case "google":
			config = h.oauth2Config.GoogleConfig
		case "github":
			config = h.oauth2Config.GitHubConfig
		default:
			http.Error(w, "Unknown provider", http.StatusBadRequest)
			return
		}

		state := generateStateToken()
		h.oauth2Config.storeState(r, w, state, provider)

		url := config.AuthCodeURL(state, oauth2.AccessTypeOffline)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

// AuthCallback برگشت از provider
func (h *AuthHandler) AuthCallback(provider string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// بررسی state
		state := r.URL.Query().Get("state")
		returnedProvider, ok := h.oauth2Config.verifyState(r, state)
		if !ok || returnedProvider != provider {
			http.Error(w, "Invalid state parameter", http.StatusBadRequest)
			return
		}

		// دریافت کد
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "Code not found", http.StatusBadRequest)
			return
		}

		// انتخاب config مناسب
		var config *oauth2.Config
		var getUserInfo func(context.Context, *http.Client) (*UserInfo, error)

		switch provider {
		case "google":
			config = h.oauth2Config.GoogleConfig
			getUserInfo = getGoogleUserInfo
		case "github":
			config = h.oauth2Config.GitHubConfig
			getUserInfo = getGitHubUserInfo
		default:
			http.Error(w, "Unknown provider", http.StatusBadRequest)
			return
		}

		// تعویض کد با توکن
		ctx := context.Background()
		token, err := config.Exchange(ctx, code)
		if err != nil {
			http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// ایجاد HTTP client با توکن
		client := config.Client(ctx, token)

		// دریافت اطلاعات کاربر
		userInfo, err := getUserInfo(ctx, client)
		if err != nil {
			http.Error(w, "Failed to get user info: "+err.Error(), http.StatusInternalServerError)
			return
		}
		userInfo.Provider = provider

		// ذخیره یا پیدا کردن کاربر در دیتابیس
		user, err := h.userRepo.FindOrCreate(userInfo)
		if err != nil {
			http.Error(w, "Failed to save user: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// ایجاد session
		if err := h.sessionMgr.SetUser(w, r, user); err != nil {
			http.Error(w, "Failed to create session: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// هدایت به صفحه پروفایل
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
	}
}

// ProfilePage صفحه پروفایل کاربر
func (h *AuthHandler) ProfilePage(w http.ResponseWriter, r *http.Request) {
	user, err := h.sessionMgr.GetUser(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	tmpl := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Profile - OAuth2 Demo</title>
		<style>
			body {
				font-family: Arial, sans-serif;
				background: #f0f2f5;
				margin: 0;
				padding: 20px;
			}
			.container {
				max-width: 600px;
				margin: 0 auto;
				background: white;
				border-radius: 10px;
				padding: 30px;
				box-shadow: 0 2px 10px rgba(0,0,0,0.1);
			}
			h1 {
				color: #333;
				border-bottom: 2px solid #667eea;
				padding-bottom: 10px;
			}
			.info {
				margin: 20px 0;
			}
			.label {
				font-weight: bold;
				color: #666;
				width: 100px;
				display: inline-block;
			}
			.value {
				color: #333;
			}
			.btn {
				background: #667eea;
				color: white;
				padding: 10px 20px;
				border: none;
				border-radius: 5px;
				cursor: pointer;
				text-decoration: none;
				display: inline-block;
				margin-top: 20px;
			}
			.btn:hover {
				background: #5a67d8;
			}
			.btn-logout {
				background: #e53e3e;
			}
			.btn-logout:hover {
				background: #c53030;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<h1>Welcome, {{.Name}}!</h1>
			<div class="info">
				<div><span class="label">ID:</span> <span class="value">{{.ID}}</span></div>
				<div><span class="label">Email:</span> <span class="value">{{.Email}}</span></div>
				<div><span class="label">Name:</span> <span class="value">{{.Name}}</span></div>
				<div><span class="label">Provider:</span> <span class="value">{{.Provider}}</span></div>
			</div>
			<a href="/logout" class="btn btn-logout">Logout</a>
		</div>
	</body>
	</html>
	`

	t := template.Must(template.New("profile").Parse(tmpl))
	t.Execute(w, user)
}

// LogoutHandler خروج از سیستم
func (h *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	h.sessionMgr.ClearUser(w, r)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// ============================================================================
// بخش 7: Middleware احراز هویت
// ============================================================================

// AuthMiddleware میدلور بررسی احراز هویت
func (h *AuthHandler) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := h.sessionMgr.GetUser(r)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}

// ============================================================================
// بخش 8: Setup Server
// ============================================================================

func setupServer() {
	// بارگذاری تنظیمات
	cfg := LoadConfig()

	// بررسی وجود تنظیمات
	if cfg.GoogleClientID == "" && cfg.GitHubClientID == "" {
		log.Println("⚠️ Warning: No OAuth2 credentials configured!")
		log.Println("Please set environment variables:")
		log.Println("  GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET")
		log.Println("  GITHUB_CLIENT_ID, GITHUB_CLIENT_SECRET")
		log.Println("Or get them from:")
		log.Println("  Google: https://console.cloud.google.com/apis/credentials")
		log.Println("  GitHub: https://github.com/settings/developers")
	}

	// ایجاد وابستگی‌ها
	oauth2Config := NewOAuth2Config(cfg)
	userRepo := NewInMemoryUserRepository()
	sessionMgr := NewSessionManager(cfg.SessionSecret)
	authHandler := NewAuthHandler(oauth2Config, userRepo, sessionMgr)

	// Routes
	mux := http.NewServeMux()

	// صفحه اصلی و ورود
	mux.HandleFunc("/", authHandler.LoginPage)
	mux.HandleFunc("/login", authHandler.LoginPage)

	// OAuth2 routes
	mux.HandleFunc("/auth/google", authHandler.AuthRedirect("google"))
	mux.HandleFunc("/auth/google/callback", authHandler.AuthCallback("google"))
	mux.HandleFunc("/auth/github", authHandler.AuthRedirect("github"))
	mux.HandleFunc("/auth/github/callback", authHandler.AuthCallback("github"))

	// Protected routes
	mux.HandleFunc("/profile", authHandler.AuthMiddleware(authHandler.ProfilePage))
	mux.HandleFunc("/logout", authHandler.LogoutHandler)

	// Start server
	port := "8080"
	log.Printf("Server starting on http://localhost:%s", port)
	log.Println("Endpoints:")
	log.Println("  GET  /                    - Login page")
	log.Println("  GET  /auth/google         - Login with Google")
	log.Println("  GET  /auth/github         - Login with GitHub")
	log.Println("  GET  /profile             - User profile (protected)")
	log.Println("  GET  /logout              - Logout")

	log.Fatal(http.ListenAndServe(":"+port, mux))
}

// ============================================================================
// بخش 9: Helper Functions
// ============================================================================

func generateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// ============================================================================
// بخش 10: Best Practices
// ============================================================================

func bestPractices() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("💡 OAUTH2 / OIDC BEST PRACTICES")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ 1. SECURITY                                                               │
├─────────────────────────────────────────────────────────────────────────────┤
│    • Always use HTTPS in production                                       │
│    • Use state parameter to prevent CSRF                                  │
│    • Store client secrets securely (env variables)                        │
│    • Use PKCE for public clients                                          │
│    • Validate ID tokens (if using OIDC)                                   │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 2. TOKEN MANAGEMENT                                                       │
├─────────────────────────────────────────────────────────────────────────────┤
│    • Store refresh tokens securely                                        │
│    • Implement token rotation                                             │
│    • Set appropriate token expiration times                               │
│    • Use short-lived access tokens                                        │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 3. USER DATA                                                             │
├─────────────────────────────────────────────────────────────────────────────┤
│    • Never trust user data from provider without validation              │
│    • Create local user record on first login                             │
│    • Handle email conflicts across providers                              │
│    • Allow users to connect multiple providers                           │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 4. PROVIDER SETUP                                                        │
├─────────────────────────────────────────────────────────────────────────────┤
│    Google:                                                                │
│    • Create project at console.cloud.google.com                          │
│    • Enable Google+ API or People API                                    │
│    • Configure redirect URIs                                             │
│                                                                           │
│    GitHub:                                                               │
│    • Register OAuth App at github.com/settings/developers                │
│    • Set Authorization callback URL                                      │
│    • Request user:email scope                                            │
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
│ MISTAKE 1: Not using state parameter                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│    ❌ Redirect without state                                              │
│    ✅ Always include and verify state                                     │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ MISTAKE 2: Exposing client secret in client-side code                     │
├─────────────────────────────────────────────────────────────────────────────┤
│    ❌ Using implicit flow with secret in JS                               │
│    ✅ Use authorization code flow on server                               │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ MISTAKE 3: Not validating email domain                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│    ❌ Accepting any email from provider                                   │
│    ✅ Restrict to specific domains if needed                              │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ MISTAKE 4: Storing provider tokens insecurely                             │
├─────────────────────────────────────────────────────────────────────────────┤
│    ❌ Tokens in localStorage                                              │
│    ✅ Tokens on server with encryption                                    │
└─────────────────────────────────────────────────────────────────────────────┘
`)
}

// ============================================================================
// بخش 12: Main
// ============================================================================

func main() {
	fmt.Println(stringsRepeat("=", 80))
	fmt.Println("🎯 COMPLETE OAUTH2 / OPENID CONNECT GUIDE")
	fmt.Println("Login with Google, GitHub in Go")
	fmt.Println(stringsRepeat("=", 80))

	bestPractices()
	commonMistakes()

	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🚀 Starting OAuth2 Server")
	fmt.Println(stringsRepeat("=", 80))
	fmt.Println()
	fmt.Println("To use Google Login:")
	fmt.Println("  1. Go to https://console.cloud.google.com/apis/credentials")
	fmt.Println("  2. Create OAuth 2.0 Client ID")
	fmt.Println("  3. Add redirect URI: http://localhost:8080/auth/google/callback")
	fmt.Println()
	fmt.Println("To use GitHub Login:")
	fmt.Println("  1. Go to https://github.com/settings/developers")
	fmt.Println("  2. Create OAuth App")
	fmt.Println("  3. Set redirect URI: http://localhost:8080/auth/github/callback")
	fmt.Println()
	fmt.Println("Set environment variables:")
	fmt.Println("  export GOOGLE_CLIENT_ID=your_client_id")
	fmt.Println("  export GOOGLE_CLIENT_SECRET=your_client_secret")
	fmt.Println("  export GITHUB_CLIENT_ID=your_client_id")
	fmt.Println("  export GITHUB_CLIENT_SECRET=your_client_secret")
	fmt.Println()

	setupServer()
}

// تابع کمکی
func stringsRepeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}
