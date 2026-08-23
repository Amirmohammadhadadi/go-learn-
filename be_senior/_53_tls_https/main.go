// ============================================================================
// FILE: tls_https_guide.go
// TITLE: راهنمای کامل TLS/HTTPS در Go - توسعه و پروداکشن
// HOW TO RUN: go run tls_https_guide.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - TLS/HTTPS چیست و چرا مهم است؟
// ============================================================================
//
// TLS (Transport Layer Security) پروتکلی برای امنیت ارتباطات شبکه است.
// HTTPS = HTTP + TLS
//
// مزایای TLS/HTTPS:
// 1. رمزنگاری (Encryption): محافظت از داده در برابر eavesdropping
// 2. احراز هویت (Authentication): تأیید هویت سرور (و گاهی کلاینت)
// 3. یکپارچگی (Integrity): جلوگیری از تغییر داده در حین انتقال
//
// انواع گواهی (Certificates):
// - Self-Signed: برای توسعه و تست
// - Let's Encrypt: رایگان، خودکار، مناسب برای پروداکشن
// - Commercial: خریداری شده از CAs معتبر
//
// قانون طلایی:
// "همیشه در پروداکشن از HTTPS استفاده کن.
//  برای توسعه از self-signed certificate استفاده کن (یا mkcert).
//  از TLS 1.2 یا بالاتر استفاده کن، TLS 1.0 و 1.1 منسوخ شده‌اند."
// ============================================================================

package __tls_https

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// ============================================================================
// بخش 1: تولید Self-Signed Certificate (برای توسعه)
// ============================================================================

// GenerateSelfSignedCert تولید گواهی self-signed برای توسعه
func GenerateSelfSignedCert(host string, validDays int) (certPEM, keyPEM []byte, err error) {
	// تنظیمات گواهی
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(time.Duration(validDays) * 24 * time.Hour)

	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, nil, err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Development"},
			CommonName:   host,
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{host, "localhost", "127.0.0.1"},
	}

	// ایجاد گواهی
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return nil, nil, err
	}

	// PEM encoding
	certPEM = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	keyPEM = pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})

	return certPEM, keyPEM, nil
}

// SaveCertFiles ذخیره گواهی در فایل
func SaveCertFiles(certPEM, keyPEM []byte, certPath, keyPath string) error {
	if err := os.WriteFile(certPath, certPEM, 0644); err != nil {
		return err
	}
	if err := os.WriteFile(keyPath, keyPEM, 0600); err != nil {
		return err
	}
	return nil
}

// SetupDevCertificates تنظیم گواهی برای توسعه
func SetupDevCertificates(host string) (certFile, keyFile string, err error) {
	certDir := filepath.Join(os.TempDir(), "dev-certs")
	if err := os.MkdirAll(certDir, 0700); err != nil {
		return "", "", err
	}

	certPath := filepath.Join(certDir, "cert.pem")
	keyPath := filepath.Join(certDir, "key.pem")

	// اگر گواهی وجود دارد و معتبر است، آن را برگردان
	if _, err := os.Stat(certPath); err == nil {
		if _, err := os.Stat(keyPath); err == nil {
			return certPath, keyPath, nil
		}
	}

	// تولید گواهی جدید
	certPEM, keyPEM, err := GenerateSelfSignedCert(host, 365)
	if err != nil {
		return "", "", err
	}

	if err := SaveCertFiles(certPEM, keyPEM, certPath, keyPath); err != nil {
		return "", "", err
	}

	log.Printf("Development certificates generated at %s", certDir)
	return certPath, keyPath, nil
}

// ============================================================================
// بخش 2: TLS Configuration بهینه برای پروداکشن
// ============================================================================

// ProductionTLSConfig تنظیمات TLS برای پروداکشن
func ProductionTLSConfig(certFile, keyFile string) (*tls.Config, error) {
	// بارگذاری گواهی
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load certificate: %w", err)
	}

	// تنظیمات امنیتی
	return &tls.Config{
		Certificates: []tls.Certificate{cert},

		// فقط TLS 1.2 و 1.3
		MinVersion: tls.VersionTLS12,
		MaxVersion: tls.VersionTLS13,

		// Cipher suites امن (برای TLS 1.2)
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},

		// ترجیح cipher suites سرور
		PreferServerCipherSuites: true,

		// تنظیمات session tickets
		SessionTicketsDisabled: false,
		SessionTicketKey:       [32]byte{}, // در عمل از کلید تصادفی استفاده کن

		// تنظیمات curves
		CurvePreferences: []tls.CurveID{
			tls.CurveP256,
			tls.X25519,
			tls.CurveP384,
		},

		// تنظیمات Next Protocol Negotiation (برای HTTP/2)
		NextProtos: []string{"h2", "http/1.1"},

		// احراز هویت کلاینت (اختیاری)
		ClientAuth: tls.NoClientCert,

		// تنظیمات زمان
		// (می‌توان با GetConfigForClient دینامیک کرد)
	}, nil
}

// DevelopmentTLSConfig تنظیمات TLS برای توسعه (کمتر restrictive)
func DevelopmentTLSConfig(certFile, keyFile string) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		},
	}, nil
}

// ============================================================================
// بخش 3: HTTP Server with TLS
// ============================================================================

// SecureServer سرور امن با TLS
type SecureServer struct {
	server   *http.Server
	config   *tls.Config
	certFile string
	keyFile  string
}

// NewSecureServer ایجاد سرور امن جدید
func NewSecureServer(addr string, handler http.Handler, certFile, keyFile string, isProduction bool) (*SecureServer, error) {
	var tlsConfig *tls.Config
	var err error

	if isProduction {
		tlsConfig, err = ProductionTLSConfig(certFile, keyFile)
	} else {
		tlsConfig, err = DevelopmentTLSConfig(certFile, keyFile)
	}

	if err != nil {
		return nil, err
	}

	server := &http.Server{
		Addr:         addr,
		Handler:      handler,
		TLSConfig:    tlsConfig,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &SecureServer{
		server:   server,
		config:   tlsConfig,
		certFile: certFile,
		keyFile:  keyFile,
	}, nil
}

// Start شروع سرور
func (s *SecureServer) Start() error {
	log.Printf("Starting HTTPS server on %s", s.server.Addr)
	return s.server.ListenAndServeTLS(s.certFile, s.keyFile)
}

// StartWithAutoCert شروع سرور با گواهی خودکار (اگر قبلاً بارگذاری شده)
func (s *SecureServer) StartWithAutoCert() error {
	log.Printf("Starting HTTPS server on %s", s.server.Addr)
	return s.server.ListenAndServeTLS("", "")
}

// Shutdown خاموش کردن سرور
func (s *SecureServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

// ============================================================================
// بخش 4: HTTP to HTTPS Redirect (HSTS)
// ============================================================================

// RedirectToHTTPSMiddleware میدلور هدایت HTTP به HTTPS
func RedirectToHTTPSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// اگر درخواست HTTP است و در محیط پروداکشن
		if r.Header.Get("X-Forwarded-Proto") == "http" || r.URL.Scheme == "http" {
			// ساخت URL HTTPS
			httpsURL := "https://" + r.Host + r.URL.Path
			if r.URL.RawQuery != "" {
				httpsURL += "?" + r.URL.RawQuery
			}

			// هدایت با status 301 (Moved Permanently)
			http.Redirect(w, r, httpsURL, http.StatusMovedPermanently)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// HSTSMiddleware افزودن HSTS header
func HSTSMiddleware(maxAge int, includeSubdomains, preload bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := fmt.Sprintf("max-age=%d", maxAge)
			if includeSubdomains {
				header += "; includeSubDomains"
			}
			if preload {
				header += "; preload"
			}
			w.Header().Set("Strict-Transport-Security", header)
			next.ServeHTTP(w, r)
		})
	}
}

// ============================================================================
// بخش 5: HTTP Server (برای redirect)
// ============================================================================

// StartHTTPServer سرور HTTP برای redirect به HTTPS
func StartHTTPServer(httpAddr, httpsAddr string) {
	redirectHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ساخت URL HTTPS
		target := "https://" + httpsAddr + r.URL.Path
		if r.URL.RawQuery != "" {
			target += "?" + r.URL.RawQuery
		}
		http.Redirect(w, r, target, http.StatusMovedPermanently)
	})

	server := &http.Server{
		Addr:         httpAddr,
		Handler:      redirectHandler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	log.Printf("HTTP redirect server starting on %s (redirecting to https://%s)", httpAddr, httpsAddr)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP server error: %v", err)
		}
	}()
}

// ============================================================================
// بخش 6: Let's Encrypt Integration (با golang.org/x/crypto/acme/autocert)
// ============================================================================

/*
برای استفاده از Let's Encrypt، پکیج زیر را نصب کنید:
go get golang.org/x/crypto/acme/autocert

// AutoTLSServer سرور با گواهی خودکار Let's Encrypt
func AutoTLSServer(domains []string, cacheDir string) *http.Server {
    manager := &autocert.Manager{
        Prompt:     autocert.AcceptTOS,
        HostPolicy: autocert.HostWhitelist(domains...),
        Cache:      autocert.DirCache(cacheDir),
        Email:      "admin@example.com", // ایمیل برای اطلاع‌رسانی
    }

    server := &http.Server{
        Addr:      ":443",
        Handler:   yourHandler,
        TLSConfig: manager.TLSConfig(),
    }

    return server
}

// همچنین یک سرور HTTP برای challenge verification نیاز است:
func startHTTPChallengeServer(manager *autocert.Manager) {
    httpServer := &http.Server{
        Addr:    ":80",
        Handler: manager.HTTPHandler(nil),
    }
    go httpServer.ListenAndServe()
}
*/

// ============================================================================
// بخش 7: TLS Client (برای اتصال امن)
// ============================================================================

// SecureClient کلاینت HTTPS امن
type SecureClient struct {
	client *http.Client
}

// NewSecureClient ایجاد کلاینت امن
func NewSecureClient(skipVerify bool) *SecureClient {
	// تنظیمات TLS کلاینت
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
		// در پروداکشن، این را true نگذارید
		InsecureSkipVerify: skipVerify,
	}

	// برای اتصال به سرورهای خاص، می‌توان Root CA اضافه کرد
	// caCertPool := x509.NewCertPool()
	// caCertPool.AppendCertsFromPEM(caCert)
	// tlsConfig.RootCAs = caCertPool

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
		Timeout: 30 * time.Second,
	}

	return &SecureClient{client: client}
}

// Get درخواست GET امن
func (c *SecureClient) Get(url string) (*http.Response, error) {
	return c.client.Get(url)
}

// Post درخواست POST امن
func (c *SecureClient) Post(url, contentType string, body io.Reader) (*http.Response, error) {
	return c.client.Post(url, contentType, body)
}

// Do انجام درخواست
func (c *SecureClient) Do(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}

// ============================================================================
// بخش 8: Certificate Inspection (بررسی گواهی)
// ============================================================================

// InspectCertificate بررسی جزئیات گواهی
func InspectCertificate(certFile string) error {
	certPEM, err := os.ReadFile(certFile)
	if err != nil {
		return err
	}

	block, _ := pem.Decode(certPEM)
	if block == nil {
		return fmt.Errorf("failed to parse certificate PEM")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return err
	}

	fmt.Println("Certificate Details:")
	fmt.Printf("  Subject: %s\n", cert.Subject.String())
	fmt.Printf("  Issuer: %s\n", cert.Issuer.String())
	fmt.Printf("  DNS Names: %v\n", cert.DNSNames)
	fmt.Printf("  Not Before: %s\n", cert.NotBefore.Format(time.RFC3339))
	fmt.Printf("  Not After: %s\n", cert.NotAfter.Format(time.RFC3339))
	fmt.Printf("  Serial Number: %s\n", cert.SerialNumber.String())
	fmt.Printf("  Signature Algorithm: %s\n", cert.SignatureAlgorithm.String())

	// بررسی انقضا
	if time.Now().After(cert.NotAfter) {
		fmt.Println("  ⚠️ Certificate has EXPIRED!")
	} else {
		daysLeft := int(time.Until(cert.NotAfter).Hours() / 24)
		fmt.Printf("  ✅ Valid for %d more days\n", daysLeft)
	}

	return nil
}

// ============================================================================
// بخش 9: Example Handlers
// ============================================================================

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `
<!DOCTYPE html>
<html>
<head>
    <title>HTTPS Demo</title>
</head>
<body>
    <h1>✅ Connected via HTTPS!</h1>
    <p>Protocol: %s</p>
    <p>TLS Version: %s</p>
    <p>Server: Go HTTPS Server</p>
</body>
</html>
`, r.Proto, GetTLSVersion(r.TLS))
}

func GetTLSVersion(tlsState *tls.ConnectionState) string {
	if tlsState == nil {
		return "None (HTTP)"
	}
	switch tlsState.Version {
	case tls.VersionTLS10:
		return "TLS 1.0"
	case tls.VersionTLS11:
		return "TLS 1.1"
	case tls.VersionTLS12:
		return "TLS 1.2"
	case tls.VersionTLS13:
		return "TLS 1.3"
	default:
		return "Unknown"
	}
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status": "secure", "protocol": "%s", "tls": "%s"}`,
		r.Proto, GetTLSVersion(r.TLS))
}

// ============================================================================
// بخش 10: Best Practices
// ============================================================================

func bestPractices() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("💡 TLS/HTTPS BEST PRACTICES")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ 1. TLS VERSIONS                                                           │
├─────────────────────────────────────────────────────────────────────────────┤
│    ✅ TLS 1.2 and TLS 1.3 only                                            │
│    ❌ TLS 1.0 and TLS 1.1 (deprecated)                                    │
│    ❌ SSLv2, SSLv3 (completely broken)                                    │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 2. CERTIFICATE MANAGEMENT                                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│    • Development: use self-signed or mkcert                               │
│    • Production: use Let's Encrypt (free) or commercial CA                │
│    • Renew certificates before expiration (auto-renewal)                  │
│    • Use strong key sizes: RSA 2048+ or ECDSA P-256                       │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 3. HSTS (HTTP Strict Transport Security)                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│    • Add HSTS header: Strict-Transport-Security                           │
│    • Start with: max-age=31536000; includeSubDomains                      │
│    • Consider preload after testing                                       │
│    • Be careful: HSTS is hard to undo!                                    │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 4. CIPHER SUITES                                                          │
├─────────────────────────────────────────────────────────────────────────────┤
│    • Prefer AEAD ciphers (GCM, ChaCha20-Poly1305)                         │
│    • Disable weak ciphers (RC4, 3DES, CBC mode)                           │
│    • Let server choose cipher suite (PreferServerCipherSuites=true)      │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 5. DEVELOPMENT vs PRODUCTION                                              │
├─────────────────────────────────────────────────────────────────────────────┤
│    Development:                                                           │
│    • Self-signed certificates OK                                          │
│    • InsecureSkipVerify for clients                                       │
│    • Less restrictive cipher suites                                       │
│                                                                           │
│    Production:                                                            │
│    • Valid CA-signed certificates                                         │
│    • Strict TLS configuration                                             │
│    • Monitor certificate expiration                                       │
│    • Enable HSTS                                                          │
└─────────────────────────────────────────────────────────────────────────────┘
`)
}

// ============================================================================
// بخش 11: Common Mistakes
// ============================================================================

func commonMistakes() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("⚠️ COMMON TLS/HTTPS MISTAKES")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ MISTAKE 1: InsecureSkipVerify in production                               │
├─────────────────────────────────────────────────────────────────────────────┤
│    ❌ client := &http.Client{Transport: &http.Transport{                 │
│           TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}       │
│    ✅ Always verify certificates in production                            │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ MISTAKE 2: Not redirecting HTTP to HTTPS                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│    ❌ Allowing HTTP traffic                                                │
│    ✅ Always redirect to HTTPS (301)                                       │
│    ✅ Enable HSTS                                                          │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ MISTAKE 3: Allowing weak TLS versions                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│    ❌ MinVersion: tls.VersionTLS10                                        │
│    ✅ MinVersion: tls.VersionTLS12                                        │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ MISTAKE 4: Not monitoring certificate expiry                              │
├─────────────────────────────────────────────────────────────────────────────┤
│    ❌ Let certificates expire                                              │
│    ✅ Set up monitoring, auto-renewal                                     │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ MISTAKE 5: Using self-signed in production                                │
├─────────────────────────────────────────────────────────────────────────────┤
│    ❌ Self-signed certificates in production                              │
│    ✅ Use Let's Encrypt or commercial CA                                  │
└─────────────────────────────────────────────────────────────────────────────┘
`)
}

// ============================================================================
// بخش 12: Complete Example Setup
// ============================================================================

func runExample() {
	// 1. تولید گواهی برای توسعه
	certFile, keyFile, err := SetupDevCertificates("localhost")
	if err != nil {
		log.Fatalf("Failed to setup certificates: %v", err)
	}

	// 2. ایجاد handler
	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/api", apiHandler)

	// 3. اعمال middlewareها
	handler := http.Handler(mux)
	handler = HSTSMiddleware(31536000, true, false)(handler)
	handler = RedirectToHTTPSMiddleware(handler) // فقط در صورتی که سرور HTTP جدا دارید

	// 4. ایجاد سرور امن
	isProduction := false // در توسعه false، در پروداکشن true
	server, err := NewSecureServer(":8443", handler, certFile, keyFile, isProduction)
	if err != nil {
		log.Fatalf("Failed to create secure server: %v", err)
	}

	// 5. (اختیاری) سرور HTTP برای redirect
	// StartHTTPServer(":8080", "localhost:8443")

	// 6. شروع سرور
	log.Println("HTTPS server starting...")
	log.Println("Open https://localhost:8443 in your browser")
	log.Println("Note: You may need to accept the self-signed certificate")

	if err := server.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

// ============================================================================
// بخش 13: Main
// ============================================================================

func main() {
	fmt.Println(stringsRepeat("=", 80))
	fmt.Println("🎯 TLS/HTTPS GUIDE")
	fmt.Println("Secure Go Applications - Development & Production")
	fmt.Println(stringsRepeat("=", 80))

	bestPractices()
	commonMistakes()

	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🚀 Starting HTTPS Example Server")
	fmt.Println(stringsRepeat("=", 80))

	runExample()
}

// تابع کمکی
func stringsRepeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}
