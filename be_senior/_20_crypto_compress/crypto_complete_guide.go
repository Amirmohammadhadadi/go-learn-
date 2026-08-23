// ============================================================================
// FILE: crypto_complete_guide.go
// TITLE: راهنمای کامل پکیج crypto در Go - رمزنگاری، هش، امضا
// HOW TO RUN: go run crypto_complete_guide.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - پکیج‌های crypto چیستند؟
// ============================================================================
//
// Go پکیج‌های قدرتمندی برای رمزنگاری ارائه می‌دهد:
//
// 1. crypto/md5, crypto/sha1, crypto/sha256, crypto/sha512: توابع هش (Hash)
// 2. crypto/hmac: کد احراز هویت مبتنی بر هش
// 3. crypto/aes, crypto/des: رمزنگاری متقارن (Symmetric)
// 4. crypto/rsa, crypto/ecdsa, crypto/ed25519: رمزنگاری نامتقارن (Asymmetric)
// 5. crypto/rand: تولید اعداد تصادفی امن
// 6. crypto/tls: امنیت لایه انتقال
// 7. crypto/x509: گواهی‌های دیجیتال
//
// قانون طلایی:
// "هرگز از MD5 و SHA1 برای برنامه‌های جدید استفاده نکن (ضعیف هستند).
//  برای هش از SHA256 یا SHA512 استفاده کن.
//  برای رمزنگاری متقارن از AES-GCM استفاده کن.
//  برای رمزنگاری نامتقارن از RSA با کلید حداقل 2048 بیت استفاده کن."
// ============================================================================

package __internal_packages

import (
	"bytes"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("🎯 COMPLETE crypto PACKAGES GUIDE IN GO")
	fmt.Println("Hashing, Encryption, Signatures, Certificates")
	fmt.Println(strings.Repeat("=", 80))

	// ============================================================================
	// بخش 1: توابع هش (Hash Functions)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🔐 SECTION 1: HASH FUNCTIONS (MD5, SHA1, SHA256, SHA512)")
	fmt.Println(strings.Repeat("=", 80))

	data := []byte("Hello, World! This is a test message for hashing.")

	// 1.1 MD5 (ضعیف - فقط برای سازگاری با سیستم‌های قدیمی)
	fmt.Println("\n--- 1.1 MD5 (Weak - Legacy only) ---")
	md5Hash := md5.Sum(data)
	fmt.Printf("  MD5: %x\n", md5Hash)
	fmt.Printf("  MD5 (base64): %s\n", base64.StdEncoding.EncodeToString(md5Hash[:]))

	// 1.2 SHA1 (ضعیف - فقط برای سازگاری)
	fmt.Println("\n--- 1.2 SHA1 (Weak - Legacy only) ---")
	sha1Hash := sha1.Sum(data)
	fmt.Printf("  SHA1: %x\n", sha1Hash)

	// 1.3 SHA256 (توصیه شده)
	fmt.Println("\n--- 1.3 SHA256 (Recommended) ---")
	sha256Hash := sha256.Sum256(data)
	fmt.Printf("  SHA256: %x\n", sha256Hash)
	fmt.Printf("  SHA256 length: %d bytes\n", len(sha256Hash))

	// 1.4 SHA512
	fmt.Println("\n--- 1.4 SHA512 ---")
	sha512Hash := sha512.Sum512(data)
	fmt.Printf("  SHA512: %x...\n", sha512Hash[:32])
	fmt.Printf("  SHA512 length: %d bytes\n", len(sha512Hash))

	// 1.5 هش با استفاده از New (برای جریان داده)
	fmt.Println("\n--- 1.5 Streaming Hash ---")
	hasher := sha256.New()
	hasher.Write([]byte("First part "))
	hasher.Write([]byte("Second part "))
	hasher.Write([]byte("Third part"))
	streamingHash := hasher.Sum(nil)
	fmt.Printf("  Streaming SHA256: %x\n", streamingHash)

	// 1.6 هش فایل
	fmt.Println("\n--- 1.6 File Hashing ---")

	// ایجاد فایل تست
	tmpFile, _ := os.CreateTemp("", "hash_test_*.txt")
	tmpFile.WriteString("Content to hash")
	tmpFileName := tmpFile.Name()
	tmpFile.Close()
	defer os.Remove(tmpFileName)

	// باز کردن و هش کردن فایل
	file, _ := os.Open(tmpFileName)
	defer file.Close()

	fileHasher := sha256.New()
	io.Copy(fileHasher, file)
	fileHash := fileHasher.Sum(nil)
	fmt.Printf("  File SHA256: %x\n", fileHash)

	// ============================================================================
	// بخش 2: HMAC (Hash-based Message Authentication Code)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🔑 SECTION 2: HMAC (Message Authentication)")
	fmt.Println(strings.Repeat("=", 80))

	// 2.1 ایجاد HMAC
	fmt.Println("\n--- 2.1 Creating HMAC ---")

	secretKey := []byte("my-secret-key-12345")
	message := []byte("Important message to authenticate")

	hmacHasher := hmac.New(sha256.New, secretKey)
	hmacHasher.Write(message)
	signature := hmacHasher.Sum(nil)

	fmt.Printf("  Message: %s\n", message)
	fmt.Printf("  HMAC-SHA256: %x\n", signature)

	// 2.2 اعتبارسنجی HMAC
	fmt.Println("\n--- 2.2 Verifying HMAC ---")

	// گیرنده با داشتن کلید می‌تواند اعتبارسنجی کند
	verifier := hmac.New(sha256.New, secretKey)
	verifier.Write(message)
	expectedMAC := verifier.Sum(nil)

	isValid := hmac.Equal(signature, expectedMAC)
	fmt.Printf("  Signature valid: %v\n", isValid)

	// با پیام تغییر یافته
	tamperedMessage := []byte("Tampered message")
	verifier.Reset()
	verifier.Write(tamperedMessage)
	tamperedMAC := verifier.Sum(nil)

	isValidTampered := hmac.Equal(signature, tamperedMAC)
	fmt.Printf("  Tampered message valid: %v\n", isValidTampered)

	// ============================================================================
	// بخش 3: رمزنگاری متقارن (AES)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🔒 SECTION 3: SYMMETRIC ENCRYPTION (AES)")
	fmt.Println(strings.Repeat("=", 80))

	// 3.1 AES-GCM (توصیه شده - حالت احراز هویت)
	fmt.Println("\n--- 3.1 AES-GCM (Authenticated Encryption) ---")

	plaintext := []byte("This is a secret message that needs encryption.")

	// تولید کلید 32 بایت (AES-256)
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		fmt.Printf("  Error generating key: %v\n", err)
		return
	}

	// ایجاد cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Printf("  Error creating cipher: %v\n", err)
		return
	}

	// ایجاد GCM (Galois/Counter Mode)
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Printf("  Error creating GCM: %v\n", err)
		return
	}

	// تولید nonce (عدد یکبار مصرف)
	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		fmt.Printf("  Error generating nonce: %v\n", err)
		return
	}

	// رمزنگاری (Seal)
	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)

	fmt.Printf("  Original: %s\n", plaintext)
	fmt.Printf("  Original size: %d bytes\n", len(plaintext))
	fmt.Printf("  Ciphertext size: %d bytes\n", len(ciphertext))
	fmt.Printf("  Nonce: %x\n", nonce)

	// رمزگشایی (Open)
	decrypted, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		fmt.Printf("  Error decrypting: %v\n", err)
		return
	}

	fmt.Printf("  Decrypted: %s\n", decrypted)
	fmt.Printf("  Decryption successful: %v\n", string(decrypted) == string(plaintext))

	// 3.2 AES-CBC (حالت قدیمی‌تر)
	fmt.Println("\n--- 3.2 AES-CBC (Legacy Mode) ---")

	// تولید کلید و IV
	keyCBC := make([]byte, 32) // AES-256
	iv := make([]byte, aes.BlockSize)
	rand.Read(keyCBC)
	rand.Read(iv)

	blockCBC, _ := aes.NewCipher(keyCBC)

	// Padding برای داده (PKCS#7)
	plaintextCBC := []byte("Data for CBC mode")
	padding := aes.BlockSize - len(plaintextCBC)%aes.BlockSize
	paddedText := append(plaintextCBC, bytes.Repeat([]byte{byte(padding)}, padding)...)

	// رمزنگاری CBC
	ciphertextCBC := make([]byte, len(paddedText))
	mode := cipher.NewCBCEncrypter(blockCBC, iv)
	mode.CryptBlocks(ciphertextCBC, paddedText)

	fmt.Printf("  CBC encrypted size: %d bytes\n", len(ciphertextCBC))

	// رمزگشایی CBC
	decrypter := cipher.NewCBCDecrypter(blockCBC, iv)
	decryptedCBC := make([]byte, len(ciphertextCBC))
	decrypter.CryptBlocks(decryptedCBC, ciphertextCBC)

	// حذف padding
	unpadded := decryptedCBC[:len(decryptedCBC)-int(decryptedCBC[len(decryptedCBC)-1])]
	fmt.Printf("  CBC decrypted: %s\n", unpadded)

	// ============================================================================
	// بخش 4: رمزنگاری نامتقارن (RSA)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🔓 SECTION 4: ASYMMETRIC ENCRYPTION (RSA)")
	fmt.Println(strings.Repeat("=", 80))

	// 4.1 تولید کلید RSA
	fmt.Println("\n--- 4.1 Generating RSA Keys ---")

	// تولید کلید خصوصی 2048 بیت
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Printf("  Error generating RSA key: %v\n", err)
		return
	}

	// استخراج کلید عمومی
	publicKey := &privateKey.PublicKey

	fmt.Printf("  Private key size: %d bits\n", privateKey.N.BitLen())
	fmt.Printf("  Public key exponent: %d\n", publicKey.E)

	// 4.2 رمزنگاری با کلید عمومی
	fmt.Println("\n--- 4.2 Encryption with Public Key ---")

	rsaPlaintext := []byte("Secret message for RSA encryption")

	// رمزنگاری با کلید عمومی
	rsaCiphertext, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, rsaPlaintext)
	if err != nil {
		fmt.Printf("  Encryption error: %v\n", err)
		return
	}

	fmt.Printf("  Original: %s\n", rsaPlaintext)
	fmt.Printf("  Original size: %d bytes\n", len(rsaPlaintext))
	fmt.Printf("  Encrypted size: %d bytes\n", len(rsaCiphertext))

	// 4.3 رمزگشایی با کلید خصوصی
	fmt.Println("\n--- 4.3 Decryption with Private Key ---")

	rsaDecrypted, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, rsaCiphertext)
	if err != nil {
		fmt.Printf("  Decryption error: %v\n", err)
		return
	}

	fmt.Printf("  Decrypted: %s\n", rsaDecrypted)
	fmt.Printf("  Decryption successful: %v\n", string(rsaDecrypted) == string(rsaPlaintext))

	// 4.4 امضای دیجیتال با RSA
	fmt.Println("\n--- 4.4 Digital Signatures with RSA ---")

	signData := []byte("Document to sign")

	// هش کردن داده
	hash := sha256.Sum256(signData)

	// امضا کردن
	signatureRSA, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash[:])
	if err != nil {
		fmt.Printf("  Sign error: %v\n", err)
		return
	}

	fmt.Printf("  Signature: %x...\n", signatureRSA[:32])

	// اعتبارسنجی امضا
	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash[:], signatureRSA)
	if err == nil {
		fmt.Println("  Signature valid: true")
	} else {
		fmt.Printf("  Signature valid: false (%v)\n", err)
	}

	// 4.5 سریالایز کردن کلیدها (PEM)
	fmt.Println("\n--- 4.5 Serializing Keys (PEM format) ---")

	// تبدیل کلید خصوصی به PEM
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privatePEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	// تبدیل کلید عمومی به PEM
	publicKeyBytes, _ := x509.MarshalPKIXPublicKey(publicKey)
	publicPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	fmt.Printf("  Private key PEM (first 100 chars):\n%s...\n", string(privatePEM[:100]))
	fmt.Printf("  Public key PEM (first 100 chars):\n%s...\n", string(publicPEM[:100]))

	// ============================================================================
	// بخش 5: اعداد تصادفی امن (crypto/rand)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🎲 SECTION 5: CRYPTOGRAPHIC RANDOM NUMBERS")
	fmt.Println(strings.Repeat("=", 80))

	// 5.1 تولید بایت‌های تصادفی
	fmt.Println("\n--- 5.1 Random Bytes ---")

	randomBytes := make([]byte, 32)
	_, err = rand.Read(randomBytes)
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
		return
	}

	fmt.Printf("  Random bytes: %x\n", randomBytes)

	// 5.2 تولید عدد تصادفی در بازه
	fmt.Println("\n--- 5.2 Random Integer in Range ---")

	// تابع تولید عدد تصادفی در بازه [0, max)
	randomInt := func(max int64) (int64, error) {
		if max <= 0 {
			return 0, fmt.Errorf("max must be positive")
		}
		var n int64
		bytes := make([]byte, 8)
		_, err := rand.Read(bytes)
		if err != nil {
			return 0, err
		}
		for _, b := range bytes {
			n = (n << 8) | int64(b)
		}
		return n % max, nil
	}

	num, _ := randomInt(100)
	fmt.Printf("  Random number 0-99: %d\n", num)

	num2, _ := randomInt(1000000)
	fmt.Printf("  Random number 0-999999: %d\n", num2)

	// 5.3 تولید توکن امن
	fmt.Println("\n--- 5.3 Secure Token Generation ---")

	token := make([]byte, 32)
	rand.Read(token)
	tokenString := hex.EncodeToString(token)
	fmt.Printf("  Secure token: %s\n", tokenString)

	// ============================================================================
	// بخش 6: کاربردهای عملی
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("💼 SECTION 6: PRACTICAL APPLICATIONS")
	fmt.Println(strings.Repeat("=", 80))

	// 6.1 هش رمز عبور با Salt
	fmt.Println("\n--- 6.1 Password Hashing with Salt ---")

	hashPassword := func(password string) (string, string) {
		// تولید salt (16 بایت)
		salt := make([]byte, 16)
		rand.Read(salt)

		// ترکیب password + salt
		saltedPassword := append([]byte(password), salt...)

		// هش SHA256
		hash := sha256.Sum256(saltedPassword)

		return hex.EncodeToString(hash[:]), hex.EncodeToString(salt)
	}

	verifyPassword := func(password, storedHash, storedSalt string) bool {
		salt, _ := hex.DecodeString(storedSalt)
		saltedPassword := append([]byte(password), salt...)
		hash := sha256.Sum256(saltedPassword)
		return hex.EncodeToString(hash[:]) == storedHash
	}

	password := "mySecurePassword123"
	hash, salt := hashPassword(password)
	fmt.Printf("  Password: %s\n", password)
	fmt.Printf("  Salt: %s\n", salt)
	fmt.Printf("  Hash: %s\n", hash)

	valid := verifyPassword(password, hash, salt)
	fmt.Printf("  Verification: %v\n", valid)

	invalid := verifyPassword("wrongPassword", hash, salt)
	fmt.Printf("  Wrong password verification: %v\n", invalid)

	// 6.2 رمزنگاری کوکی‌ها
	fmt.Println("\n--- 6.2 Cookie Encryption ---")

	// کلید ثابت برای مثال (در عمل از کلید امن استفاده کنید)
	cookieKey := []byte("1234567890123456") // 16 بایت برای AES-128

	encryptCookie := func(value string) (string, error) {
		block, _ := aes.NewCipher(cookieKey)
		gcm, _ := cipher.NewGCM(block)

		nonce := make([]byte, gcm.NonceSize())
		rand.Read(nonce)

		ciphertext := gcm.Seal(nil, nonce, []byte(value), nil)
		return base64.StdEncoding.EncodeToString(append(nonce, ciphertext...)), nil
	}

	decryptCookie := func(encrypted string) (string, error) {
		data, _ := base64.StdEncoding.DecodeString(encrypted)

		block, _ := aes.NewCipher(cookieKey)
		gcm, _ := cipher.NewGCM(block)

		nonceSize := gcm.NonceSize()
		nonce, ciphertext := data[:nonceSize], data[nonceSize:]

		plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
		if err != nil {
			return "", err
		}
		return string(plaintext), nil
	}

	cookieValue := "session_id=abc123; user_id=42"
	encryptedCookie, _ := encryptCookie(cookieValue)
	decryptedCookie, _ := decryptCookie(encryptedCookie)

	fmt.Printf("  Original cookie: %s\n", cookieValue)
	fmt.Printf("  Encrypted: %s\n", encryptedCookie[:50]+"...")
	fmt.Printf("  Decrypted: %s\n", decryptedCookie)

	// 6.3 تایید یکپارچگی فایل
	fmt.Println("\n--- 6.3 File Integrity Check ---")

	integrityFile, _ := os.CreateTemp("", "integrity_*.txt")
	integrityFile.WriteString("Important data that needs integrity check")
	integrityFile.Close()
	defer os.Remove(integrityFile.Name())

	// محاسبه HMAC فایل
	integrityKey := []byte("integrity-check-key")

	fileForHash, _ := os.Open(integrityFile.Name())
	hmacHasher2 := hmac.New(sha256.New, integrityKey)
	io.Copy(hmacHasher2, fileForHash)
	fileHMAC := hmacHasher2.Sum(nil)
	fileForHash.Close()

	fmt.Printf("  File HMAC: %x\n", fileHMAC[:16])

	// اعتبارسنجی (در عمل HMAC را جداگانه ذخیره می‌کنیم)
	fmt.Println("  File integrity can be verified by recomputing HMAC")

	// ============================================================================
	// بخش 7: اشتباهات رایج
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("❌ SECTION 7: COMMON MISTAKES")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println("\n❌ Mistake 1: Using MD5 or SHA1 for security")
	fmt.Println("   These are cryptographically broken!")
	fmt.Println("   ✅ Use SHA256 or SHA512")

	fmt.Println("\n❌ Mistake 2: Using ECB mode")
	fmt.Println("   ECB is insecure for most applications")
	fmt.Println("   ✅ Use GCM (authenticated) or CBC with random IV")

	fmt.Println("\n❌ Mistake 3: Reusing nonce/IV with same key")
	fmt.Println("   Nonce must be unique for each encryption")
	fmt.Println("   ✅ Generate new random nonce for each operation")

	fmt.Println("\n❌ Mistake 4: Not using authenticated encryption")
	fmt.Println("   AES-CBC alone doesn't detect tampering")
	fmt.Println("   ✅ Use AES-GCM which includes authentication")

	fmt.Println("\n❌ Mistake 5: Hardcoding cryptographic keys")
	fmt.Println("   Keys in source code are not secure")
	fmt.Println("   ✅ Use environment variables or key management services")

	fmt.Println("\n❌ Mistake 6: Using rand package for security")
	fmt.Println("   math/rand is predictable!")
	fmt.Println("   ✅ Use crypto/rand for security-critical random numbers")

	fmt.Println("\n❌ Mistake 7: Not using constant-time comparison")
	fmt.Println("   bytes.Equal may leak timing information")
	fmt.Println("   ✅ Use hmac.Equal for comparing MACs")

	fmt.Println("\n❌ Mistake 8: Small RSA keys")
	fmt.Println("   RSA 1024-bit is considered broken")
	fmt.Println("   ✅ Use RSA 2048-bit or larger")

	// ============================================================================
	// بخش 8: جمع‌بندی و جدول مرجع
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📚 SECTION 8: QUICK REFERENCE")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ HASH FUNCTIONS                                                │")
	fmt.Println("├─────────────────────────────────────────────────────────────────┤")
	fmt.Println("│ md5.Sum(data)        - MD5 hash (32 hex chars) - WEAK!        │")
	fmt.Println("│ sha1.Sum(data)       - SHA1 hash (40 hex chars) - WEAK!       │")
	fmt.Println("│ sha256.Sum256(data)  - SHA256 hash (64 hex chars) - RECOMMENDED│")
	fmt.Println("│ sha512.Sum512(data)  - SHA512 hash (128 hex chars)            │")
	fmt.Println("└─────────────────────────────────────────────────────────────────┘")

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ HMAC (Message Authentication)                                 │")
	fmt.Println("├─────────────────────────────────────────────────────────────────┤")
	fmt.Println("│ hmac.New(hashFunc, key) - Create HMAC                        │")
	fmt.Println("│ hmac.Equal(mac1, mac2)   - Constant-time comparison          │")
	fmt.Println("└─────────────────────────────────────────────────────────────────┘")

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ SYMMETRIC ENCRYPTION (AES)                                   │")
	fmt.Println("├─────────────────────────────────────────────────────────────────┤")
	fmt.Println("│ aes.NewCipher(key)    - Create cipher block                  │")
	fmt.Println("│ cipher.NewGCM(block)  - GCM mode (authenticated)             │")
	fmt.Println("│ gcm.Seal()            - Encrypt + authenticate               │")
	fmt.Println("│ gcm.Open()            - Decrypt + verify                     │")
	fmt.Println("└─────────────────────────────────────────────────────────────────┘")

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ ASYMMETRIC ENCRYPTION (RSA)                                  │")
	fmt.Println("├─────────────────────────────────────────────────────────────────┤")
	fmt.Println("│ rsa.GenerateKey(rand, bits) - Generate key pair              │")
	fmt.Println("│ rsa.EncryptPKCS1v15()     - Encrypt with public key          │")
	fmt.Println("│ rsa.DecryptPKCS1v15()     - Decrypt with private key         │")
	fmt.Println("│ rsa.SignPKCS1v15()        - Sign with private key            │")
	fmt.Println("│ rsa.VerifyPKCS1v15()      - Verify with public key           │")
	fmt.Println("└─────────────────────────────────────────────────────────────────┘")

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ RANDOM NUMBERS (crypto/rand)                                 │")
	fmt.Println("├─────────────────────────────────────────────────────────────────┤")
	fmt.Println("│ rand.Read(bytes)       - Fill bytes with random data         │")
	fmt.Println("│ rand.Int(rand.Reader, max) - Random integer 0..max-1         │")
	fmt.Println("│ rand.Prime(rand.Reader, bits) - Random prime number          │")
	fmt.Println("└─────────────────────────────────────────────────────────────────┘")

	fmt.Println("\n💡 GOLDEN RULES:")
	fmt.Println("  1. NEVER use MD5 or SHA1 for security-critical applications")
	fmt.Println("  2. ALWAYS use authenticated encryption (AES-GCM, not AES-CBC alone)")
	fmt.Println("  3. Generate random nonce/IV for each encryption operation")
	fmt.Println("  4. Use crypto/rand for security, not math/rand")
	fmt.Println("  5. Use hmac.Equal for comparing MACs (constant-time)")
	fmt.Println("  6. Store keys in environment variables or KMS, never in code")
	fmt.Println("  7. Use RSA with at least 2048-bit keys")
	fmt.Println("  8. Always add salt when hashing passwords")
	fmt.Println("  9. Use TLS 1.3 for network communication")
	fmt.Println("  10. Keep crypto libraries updated")

	fmt.Println("\n🎯 RECOMMENDATIONS BY USE CASE:")
	fmt.Println("  • Password storage → bcrypt or scrypt (not in standard library)")
	fmt.Println("  • File integrity → SHA256 + HMAC")
	fmt.Println("  • Data at rest → AES-256-GCM")
	fmt.Println("  • Data in transit → TLS")
	fmt.Println("  • Digital signatures → RSA with SHA256")
	fmt.Println("  • API authentication → HMAC-SHA256")
	fmt.Println("  • Session tokens → crypto/rand + hex encoding")
}

// تابع کمکی برای تکرار رشته
func stringsRepeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}

// اضافه کردن import crypto برای RSA
// (در فایل واقعی باید import "crypto" هم باشد)
type cryptoSigner interface {
	crypto.Signer
}

// bytes.Repeat برای padding (در فایل واقعی از bytes.Repeat استفاده کنید)
func bytesRepeat(b byte, count int) []byte {
	result := make([]byte, count)
	for i := 0; i < count; i++ {
		result[i] = b
	}
	return result
}
