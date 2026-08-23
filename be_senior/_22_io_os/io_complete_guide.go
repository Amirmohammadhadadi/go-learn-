// ============================================================================
// FILE: io_complete_guide.go
// TITLE: راهنمای کامل پکیج io در Go - ورودی/خروجی، جریان‌ها، اینترفیس‌ها
// HOW TO RUN: go run io_complete_guide.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - پکیج io چیست و چرا مهم است؟
// ============================================================================
//
// پکیج io اینترفیس‌های اصلی برای عملیات ورودی/خروجی را فراهم می‌کند:
// - خواندن (Reader)
// - نوشتن (Writer)
// - بستن (Closer)
// - جستجو (Seeker)
// - و ترکیب‌های آنها
//
// فلسفه طراحی:
// "هر جا داده می‌خوانی، Reader. هر جا داده می‌نویسی، Writer.
//  این دو اینترفیس قلب I/O در Go هستند."
//
// قانون طلایی:
// "از اینترفیس‌های کوچک io.Reader و io.Writer برای پذیرش انواع مختلف
//  منابع داده استفاده کن. توابع خود را بر اساس اینترفیس بنویس، نه نوع مشخص."
// ============================================================================

package __internal_packages

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

// ============================================================================
// بخش 1: اینترفیس‌های اصلی پکیج io
// ============================================================================

func demonstrateCoreInterfaces() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🔧 CORE io INTERFACES - Reader, Writer, Closer, Seeker")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 1.1 اینترفیس Reader
	// ============================================
	fmt.Println("\n--- 1.1 io.Reader Interface ---")
	// type Reader interface {
	//     Read(p []byte) (n int, err error)
	// }

	// هر چیزی که قابل خواندن باشد، Reader است
	data := strings.NewReader("Hello, World!")
	buffer := make([]byte, 5)

	for {
		n, err := data.Read(buffer)
		if err != nil {
			if err == io.EOF {
				fmt.Println("  End of file reached")
				break
			}
			fmt.Printf("  Error: %v\n", err)
			break
		}
		fmt.Printf("  Read %d bytes: %s\n", n, buffer[:n])
	}

	// ============================================
	// 1.2 اینترفیس Writer
	// ============================================
	fmt.Println("\n--- 1.2 io.Writer Interface ---")
	// type Writer interface {
	//     Write(p []byte) (n int, err error)
	// }

	var buf bytes.Buffer
	written, err := buf.Write([]byte("Hello, "))
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	}
	written2, _ := buf.WriteString("World!")
	fmt.Printf("  Wrote %d + %d bytes, buffer: %s\n", written, written2, buf.String())

	// ============================================
	// 1.3 اینترفیس Closer
	// ============================================
	fmt.Println("\n--- 1.3 io.Closer Interface ---")
	// type Closer interface {
	//     Close() error
	// }

	// فایل نمونه (در بخش os بیشتر خواهیم دید)
	file, _ := os.CreateTemp("", "example")
	defer os.Remove(file.Name())

	file.WriteString("test data")
	file.Close()
	fmt.Printf("  File closed: %s\n", file.Name())

	// ============================================
	// 1.4 اینترفیس Seeker
	// ============================================
	fmt.Println("\n--- 1.4 io.Seeker Interface ---")
	// type Seeker interface {
	//     Seek(offset int64, whence int) (int64, error)
	// }

	// whence: 0 = از ابتدا, 1 = از当前位置, 2 = از انتها
	seekData := strings.NewReader("0123456789")

	// رفتن به اندیس 5
	seekData.Seek(5, io.SeekStart)
	pos, _ := seekData.Seek(0, io.SeekCurrent)
	fmt.Printf("  Current position: %d\n", pos)

	b := make([]byte, 3)
	seekData.Read(b)
	fmt.Printf("  Read from position 5: %s\n", b)

	// رفتن به انتها
	seekData.Seek(0, io.SeekEnd)
	endPos, _ := seekData.Seek(0, io.SeekCurrent)
	fmt.Printf("  End position: %d\n", endPos)

	// ============================================
	// 1.5 اینترفیس‌های ترکیبی
	// ============================================
	fmt.Println("\n--- 1.5 Combined Interfaces ---")
	// type ReadWriter interface { Reader; Writer }
	// type ReadCloser interface { Reader; Closer }
	// type WriteCloser interface { Writer; Closer }
	// type ReadWriteCloser interface { Reader; Writer; Closer }
	// type ReadSeeker interface { Reader; Seeker }
	// type WriteSeeker interface { Writer; Seeker }
	// type ReadWriteSeeker interface { Reader; Writer; Seeker }

	fmt.Println("  Standard combined interfaces are available in package io")
}

// ============================================================================
// بخش 2: توابع مهم پکیج io
// ============================================================================

func demonstrateImportantFunctions() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📚 IMPORTANT io FUNCTIONS")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 2.1 io.Copy - کپی از Reader به Writer
	// ============================================
	fmt.Println("\n--- 2.1 io.Copy ---")

	src := strings.NewReader("This is the source text to copy")
	dst := &bytes.Buffer{}

	written, err := io.Copy(dst, src)
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	} else {
		fmt.Printf("  Copied %d bytes: %s\n", written, dst.String())
	}

	// ============================================
	// 2.2 io.CopyN - کپی تعداد مشخصی بایت
	// ============================================
	fmt.Println("\n--- 2.2 io.CopyN ---")

	src2 := strings.NewReader("1234567890")
	dst2 := &bytes.Buffer{}

	written, err = io.CopyN(dst2, src2, 5)
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	} else {
		fmt.Printf("  Copied %d bytes: %s\n", written, dst2.String())
	}

	// ============================================
	// 2.3 io.ReadAll - خواندن همه داده‌ها
	// ============================================
	fmt.Println("\n--- 2.3 io.ReadAll ---")

	reader := strings.NewReader("Read all of this text completely")

	data, err := io.ReadAll(reader)
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	} else {
		fmt.Printf("  Read %d bytes: %s\n", len(data), string(data))
	}

	// ============================================
	// 2.4 io.ReadAtLeast - خواندن حداقل N بایت
	// ============================================
	fmt.Println("\n--- 2.4 io.ReadAtLeast ---")

	r := strings.NewReader("Minimum read test")
	buf := make([]byte, 20)

	n, err := io.ReadAtLeast(r, buf, 10)
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	} else {
		fmt.Printf("  Read at least 10 bytes, actually read %d: %s\n", n, buf[:n])
	}

	// ============================================
	// 2.5 io.ReadFull - خواندن دقیقاً N بایت
	// ============================================
	fmt.Println("\n--- 2.5 io.ReadFull ---")

	r2 := strings.NewReader("Exactly 15 bytes!")
	buf2 := make([]byte, 15)

	n, err = io.ReadFull(r2, buf2)
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	} else {
		fmt.Printf("  Read exactly %d bytes: %s\n", n, buf2)
	}

	// ============================================
	// 2.6 io.WriteString - نوشتن رشته به Writer
	// ============================================
	fmt.Println("\n--- 2.6 io.WriteString ---")

	var buf3 bytes.Buffer
	n, err = io.WriteString(&buf3, "Hello via WriteString")
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	} else {
		fmt.Printf("  Wrote %d bytes: %s\n", n, buf3.String())
	}

	// ============================================
	// 2.7 io.LimitReader - محدود کردن خواندن
	// ============================================
	fmt.Println("\n--- 2.7 io.LimitReader ---")

	bigData := strings.NewReader("This is a very long string that we want to limit")
	limited := io.LimitReader(bigData, 20)

	data2, err := io.ReadAll(limited)
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	} else {
		fmt.Printf("  Limited read (20 bytes): %s\n", string(data2))
	}

	// ============================================
	// 2.8 io.TeeReader - خواندن و همزمان کپی کردن
	// ============================================
	fmt.Println("\n--- 2.8 io.TeeReader ---")

	source := strings.NewReader("Tee reader copies as it reads")
	var logBuffer bytes.Buffer
	tee := io.TeeReader(source, &logBuffer)

	// خواندن از tee
	result, _ := io.ReadAll(tee)
	fmt.Printf("  Read result: %s\n", string(result))
	fmt.Printf("  Logged copy: %s\n", logBuffer.String())

	// ============================================
	// 2.9 io.MultiReader - ترکیب چند Reader
	// ============================================
	fmt.Println("\n--- 2.9 io.MultiReader ---")

	r1 := strings.NewReader("Part 1, ")
	r2 := strings.NewReader("Part 2, ")
	r3 := strings.NewReader("Part 3")

	mr := io.MultiReader(r1, r2, r3)

	combined, _ := io.ReadAll(mr)
	fmt.Printf("  Combined: %s\n", string(combined))

	// ============================================
	// 2.10 io.MultiWriter - نوشتن همزمان به چند Writer
	// ============================================
	fmt.Println("\n--- 2.10 io.MultiWriter ---")

	var buf4, buf5 bytes.Buffer
	mw := io.MultiWriter(&buf4, &buf5)

	io.WriteString(mw, "Written to both buffers")

	fmt.Printf("  Buffer 1: %s\n", buf4.String())
	fmt.Printf("  Buffer 2: %s\n", buf5.String())

	// ============================================
	// 2.11 io.Pipe - کانال بین Reader و Writer
	// ============================================
	fmt.Println("\n--- 2.11 io.Pipe ---")

	pr, pw := io.Pipe()

	go func() {
		defer pw.Close()
		pw.Write([]byte("Data through pipe"))
	}()

	pipeData, _ := io.ReadAll(pr)
	fmt.Printf("  Pipe data: %s\n", string(pipeData))
}

// ============================================================================
// بخش 3: پیاده‌سازی اینترفیس‌های io (سفارشی)
// ============================================================================

// 3.1 Reader سفارشی (تولید اعداد تصادفی)
type RandomReader struct {
	Count int
}

func (r *RandomReader) Read(p []byte) (n int, err error) {
	if r.Count <= 0 {
		return 0, io.EOF
	}

	// تولید بایت‌های تصادفی ساده
	for i := range p {
		if r.Count <= 0 {
			break
		}
		p[i] = byte(time.Now().UnixNano() % 256)
		n++
		r.Count--
	}

	return n, nil
}

// 3.2 Writer سفارشی (نوشتن با شمارش)
type CountingWriter struct {
	Writer io.Writer
	Count  int64
}

func (cw *CountingWriter) Write(p []byte) (n int, err error) {
	n, err = cw.Writer.Write(p)
	cw.Count += int64(n)
	return n, err
}

// 3.3 Reader از نوع محدودکننده نرخ (Rate Limiter)
type RateLimitedReader struct {
	Reader   io.Reader
	Rate     int // bytes per second
	lastRead time.Time
}

func (r *RateLimitedReader) Read(p []byte) (n int, err error) {
	if r.lastRead.IsZero() {
		r.lastRead = time.Now()
	}

	// محدودیت نرخ
	elapsed := time.Since(r.lastRead)
	expectedRate := float64(r.Rate) * elapsed.Seconds()

	if expectedRate < 1 && len(p) > 0 {
		time.Sleep(time.Duration(float64(time.Second) / float64(r.Rate)))
	}

	r.lastRead = time.Now()
	return r.Reader.Read(p)
}

// 3.4 Writer با هش (برای محاسبه MD5 همزمان)
type HashWriter struct {
	Writer io.Writer
	Hash   md5.Hash
}

func NewHashWriter(w io.Writer) *HashWriter {
	return &HashWriter{
		Writer: w,
		Hash:   md5.New(),
	}
}

func (hw *HashWriter) Write(p []byte) (n int, err error) {
	// همزمان در Writer اصلی و در هش می‌نویسد
	hw.Hash.Write(p)
	return hw.Writer.Write(p)
}

func (hw *HashWriter) Sum() []byte {
	return hw.Hash.Sum(nil)
}

func demonstrateCustomReadersWriters() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🔨 CUSTOM io.Reader AND io.Writer IMPLEMENTATIONS")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 3.1 Random Reader
	// ============================================
	fmt.Println("\n--- 3.1 Custom RandomReader ---")

	randReader := &RandomReader{Count: 10}
	randData, _ := io.ReadAll(randReader)
	fmt.Printf("  Random bytes: %x\n", randData[:min(10, len(randData))])

	// ============================================
	// 3.2 Counting Writer
	// ============================================
	fmt.Println("\n--- 3.2 Custom CountingWriter ---")

	var buf bytes.Buffer
	countingWriter := &CountingWriter{Writer: &buf}

	io.WriteString(countingWriter, "Hello")
	io.WriteString(countingWriter, " ")
	io.WriteString(countingWriter, "World")

	fmt.Printf("  Total bytes written: %d, content: %s\n",
		countingWriter.Count, buf.String())

	// ============================================
	// 3.3 Hash Writer (MD5)
	// ============================================
	fmt.Println("\n--- 3.3 Custom HashWriter ---")

	var output bytes.Buffer
	hashWriter := NewHashWriter(&output)

	io.WriteString(hashWriter, "Data to hash")

	fmt.Printf("  Written data: %s\n", output.String())
	fmt.Printf("  MD5 hash: %x\n", hashWriter.Sum())
}

// ============================================================================
// بخش 4: کار با فایل‌ها (ادغام با os)
// ============================================================================

func demonstrateFileIO() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📁 FILE I/O - Working with Files (os + io)")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 4.1 نوشتن در فایل
	// ============================================
	fmt.Println("\n--- 4.1 Writing to File ---")

	// ایجاد فایل
	file, err := os.CreateTemp("", "example_*.txt")
	if err != nil {
		log.Printf("  Error creating file: %v", err)
		return
	}
	defer os.Remove(file.Name())
	defer file.Close()

	// نوشتن با io.WriteString
	written, _ := io.WriteString(file, "Hello from io.WriteString\n")
	fmt.Printf("  Wrote %d bytes\n", written)

	// نوشتن با fmt.Fprintf
	fmt.Fprintf(file, "Writing with Fprintf: %d\n", 42)

	// نوشتن با Write
	file.Write([]byte("Direct write\n"))

	// ============================================
	// 4.2 خواندن از فایل
	// ============================================
	fmt.Println("\n--- 4.2 Reading from File ---")

	// بازگشت به ابتدا
	file.Seek(0, io.SeekStart)

	// خواندن کل فایل
	data, _ := io.ReadAll(file)
	fmt.Printf("  File content:\n%s", string(data))

	// خواندن خط به خط
	file.Seek(0, io.SeekStart)
	buf := make([]byte, 32)
	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("  Error: %v\n", err)
			break
		}
		fmt.Printf("  Read chunk: %s", buf[:n])
	}
	fmt.Println()

	// ============================================
	// 4.3 کپی فایل
	// ============================================
	fmt.Println("\n--- 4.3 Copy File ---")

	srcFile, _ := os.CreateTemp("", "src_*.txt")
	defer os.Remove(srcFile.Name())
	srcFile.WriteString("Source content to copy")
	srcFile.Close()

	dstFile, _ := os.CreateTemp("", "dst_*.txt")
	defer os.Remove(dstFile.Name())

	// باز کردن مجدد
	srcFile, _ = os.Open(srcFile.Name())
	defer srcFile.Close()

	written, _ = io.Copy(dstFile, srcFile)
	fmt.Printf("  Copied %d bytes\n", written)
}

// ============================================================================
// بخش 5: بافر و بهینه‌سازی
// ============================================================================

func demonstrateBuffering() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("⚡ BUFFERING AND OPTIMIZATION")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 5.1 bytes.Buffer - بافر در حافظه
	// ============================================
	fmt.Println("\n--- 5.1 bytes.Buffer ---")

	var buffer bytes.Buffer

	// نوشتن
	buffer.WriteString("First line\n")
	buffer.Write([]byte("Second line\n"))
	buffer.WriteString("Third line\n")

	// خواندن
	line, _ := buffer.ReadString('\n')
	fmt.Printf("  Read line: %s", line)

	// باقیمانده
	fmt.Printf("  Remaining: %s", buffer.String())

	// کارایی: preallocation
	var efficientBuf bytes.Buffer
	efficientBuf.Grow(1024) // pre-allocate 1KB
	efficientBuf.WriteString("Efficient write")

	// reset
	efficientBuf.Reset()
	fmt.Printf("  After reset: len=%d, cap=%d\n", efficientBuf.Len(), efficientBuf.Cap())

	// ============================================
	// 5.2 strings.Builder - برای ساخت کارآمد رشته
	// ============================================
	fmt.Println("\n--- 5.2 strings.Builder ---")

	var builder strings.Builder

	for i := 0; i < 5; i++ {
		builder.WriteString(fmt.Sprintf("Item %d, ", i))
	}

	result := builder.String()
	fmt.Printf("  Built string: %s\n", result)
	fmt.Printf("  Builder stats: len=%d, cap=%d\n", builder.Len(), builder.Cap())

	// reset
	builder.Reset()
	builder.WriteString("After reset")
	fmt.Printf("  After reset: %s\n", builder.String())
}

// ============================================================================
// بخش 6: فشرده‌سازی (Compression) با io
// ============================================================================

func demonstrateCompression() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🗜️ COMPRESSION - gzip with io")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 6.1 فشرده‌سازی
	// ============================================
	fmt.Println("\n--- 6.1 Compression (gzip) ---")

	var compressed bytes.Buffer

	// Writer فشرده‌ساز
	gzWriter := gzip.NewWriter(&compressed)

	// نوشتن داده
	original := "This is the data to compress. It should become smaller."
	gzWriter.Write([]byte(original))
	gzWriter.Close()

	fmt.Printf("  Original size: %d bytes\n", len(original))
	fmt.Printf("  Compressed size: %d bytes\n", compressed.Len())

	// ============================================
	// 6.2 خارج کردن از فشرده‌سازی
	// ============================================
	fmt.Println("\n--- 6.2 Decompression ---")

	gzReader, err := gzip.NewReader(&compressed)
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
		return
	}
	defer gzReader.Close()

	decompressed, _ := io.ReadAll(gzReader)
	fmt.Printf("  Decompressed: %s\n", string(decompressed))
}

// ============================================================================
// بخش 7: رمزنگاری و هش با io
// ============================================================================

func demonstrateCryptoWithIO() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🔐 CRYPTO & HASHING with io")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 7.1 محاسبه هش همزمان با نوشتن
	// ============================================
	fmt.Println("\n--- 7.1 Hashing While Writing ---")

	data := []byte("Data to hash with multiple algorithms")

	// هش MD5
	md5Hash := md5.New()
	// هش SHA256
	shaHash := sha256.New()

	// MultiWriter برای نوشتن همزمان در هر دو هش
	multi := io.MultiWriter(md5Hash, shaHash)

	// نوشتن داده
	multi.Write(data)

	fmt.Printf("  MD5: %x\n", md5Hash.Sum(nil))
	fmt.Printf("  SHA256: %x\n", shaHash.Sum(nil))

	// ============================================
	// 7.2 Base64 encoding/decoding
	// ============================================
	fmt.Println("\n--- 7.2 Base64 Encoding ---")

	original2 := []byte("Encode this to base64")

	// Encoding
	encoded := base64.StdEncoding.EncodeToString(original2)
	fmt.Printf("  Encoded: %s\n", encoded)

	// Decoding
	decoded, _ := base64.StdEncoding.DecodeString(encoded)
	fmt.Printf("  Decoded: %s\n", string(decoded))
}

// ============================================================================
// بخش 8: تبدیل بین انواع داده با binary
// ============================================================================

func demonstrateBinaryIO() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🔢 BINARY I/O - Reading/Writing Binary Data")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 8.1 نوشتن اعداد به صورت باینری
	// ============================================
	fmt.Println("\n--- 8.1 Writing Binary Data ---")

	var buf bytes.Buffer

	// نوشتن int32
	err := binary.Write(&buf, binary.LittleEndian, int32(12345))
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	}

	// نوشتن float64
	err = binary.Write(&buf, binary.LittleEndian, float64(3.14159))
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	}

	// نوشتن string
	buf.WriteString("Hello")

	fmt.Printf("  Binary size: %d bytes\n", buf.Len())
	fmt.Printf("  Raw bytes: %x\n", buf.Bytes())

	// ============================================
	// 8.2 خواندن اعداد از باینری
	// ============================================
	fmt.Println("\n--- 8.2 Reading Binary Data ---")

	var intVal int32
	var floatVal float64
	stringVal := make([]byte, 5)

	// خواندن int32
	err = binary.Read(&buf, binary.LittleEndian, &intVal)
	if err == nil {
		fmt.Printf("  Read int32: %d\n", intVal)
	}

	// خواندن float64
	err = binary.Read(&buf, binary.LittleEndian, &floatVal)
	if err == nil {
		fmt.Printf("  Read float64: %.5f\n", floatVal)
	}

	// خواندن string
	buf.Read(stringVal)
	fmt.Printf("  Read string: %s\n", string(stringVal))
}

// ============================================================================
// بخش 9: الگوهای پیشرفته
// ============================================================================

// 9.1 TeeReader برای لاگینگ
type LoggingReader struct {
	Reader io.Reader
	Logger io.Writer
}

func (lr *LoggingReader) Read(p []byte) (n int, err error) {
	n, err = lr.Reader.Read(p)
	if n > 0 {
		lr.Logger.Write(p[:n])
	}
	return n, err
}

// 9.2 Progress Reader (نشان دادن پیشرفت)
type ProgressReader struct {
	Reader   io.Reader
	Total    int64
	Current  int64
	OnUpdate func(current, total int64)
}

func (pr *ProgressReader) Read(p []byte) (n int, err error) {
	n, err = pr.Reader.Read(p)
	pr.Current += int64(n)
	if pr.OnUpdate != nil {
		pr.OnUpdate(pr.Current, pr.Total)
	}
	return n, err
}

// 9.3 LimitWriter (محدود کردن نوشتن)
type LimitWriter struct {
	Writer io.Writer
	Limit  int64
	Wrote  int64
}

func (lw *LimitWriter) Write(p []byte) (n int, err error) {
	if lw.Wrote >= lw.Limit {
		return 0, fmt.Errorf("write limit exceeded")
	}

	remaining := lw.Limit - lw.Wrote
	if int64(len(p)) > remaining {
		p = p[:remaining]
	}

	n, err = lw.Writer.Write(p)
	lw.Wrote += int64(n)
	return n, err
}

func demonstrateAdvancedPatterns() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🎨 ADVANCED IO PATTERNS")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 9.1 Logging Reader
	// ============================================
	fmt.Println("\n--- 9.1 LoggingReader ---")

	var logBuf bytes.Buffer
	source := strings.NewReader("Data to be logged")
	loggingReader := &LoggingReader{
		Reader: source,
		Logger: &logBuf,
	}

	data, _ := io.ReadAll(loggingReader)
	fmt.Printf("  Read data: %s\n", string(data))
	fmt.Printf("  Logged: %s\n", logBuf.String())

	// ============================================
	// 9.2 Progress Reader
	// ============================================
	fmt.Println("\n--- 9.2 ProgressReader ---")

	source2 := strings.NewReader(strings.Repeat("x", 1000))

	progressReader := &ProgressReader{
		Reader: source2,
		Total:  1000,
		OnUpdate: func(current, total int64) {
			percent := float64(current) / float64(total) * 100
			fmt.Printf("\r  Progress: %.1f%%", percent)
		},
	}

	io.Copy(io.Discard, progressReader)
	fmt.Println("\n  Complete!")

	// ============================================
	// 9.3 LimitWriter
	// ============================================
	fmt.Println("\n--- 9.3 LimitWriter ---")

	var limitBuf bytes.Buffer
	limitWriter := &LimitWriter{
		Writer: &limitBuf,
		Limit:  10,
	}

	written, err := limitWriter.Write([]byte("This is a long text"))
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	}
	fmt.Printf("  Wrote %d bytes (limit was 10)\n", written)
	fmt.Printf("  Buffer: %s\n", limitBuf.String())
}

// ============================================================================
// بخش 10: اشتباهات رایج
// ============================================================================

func demonstrateCommonMistakes() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("❌ COMMON MISTAKES WITH io")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println("\n❌ Mistake 1: Not handling EOF properly")
	fmt.Println("   for { n, err := reader.Read(buf); err != nil { break } }")
	fmt.Println("   ✅ Check for io.EOF explicitly")

	fmt.Println("\n❌ Mistake 2: Assuming Read reads all requested bytes")
	fmt.Println("   n, _ := reader.Read(buf)  // may read less")
	fmt.Println("   ✅ Use io.ReadFull for exact reads")

	fmt.Println("\n❌ Mistake 3: Not closing resources")
	fmt.Println("   file, _ := os.Open(\"file.txt\")")
	fmt.Println("   // missing defer file.Close()")
	fmt.Println("   ✅ defer file.Close()")

	fmt.Println("\n❌ Mistake 4: Ignoring bytes written count")
	fmt.Println("   writer.Write(data)  // ignoring return")
	fmt.Println("   ✅ Always check n and err")

	fmt.Println("\n❌ Mistake 5: Using bytes.Buffer without preallocation")
	fmt.Println("   var buf bytes.Buffer")
	fmt.Println("   for i := 0; i < 1000000; i++ { buf.WriteString(\"x\") }")
	fmt.Println("   ✅ buf.Grow(1000000) for preallocation")

	fmt.Println("\n❌ Mistake 6: Reading entire large files")
	fmt.Println("   data, _ := io.ReadAll(file)  // memory blow")
	fmt.Println("   ✅ Use io.Copy with chunked reading")
}

// ============================================================================
// بخش 11: جمع‌بندی و جدول مرجع
// ============================================================================

func main() {
	fmt.Println(stringsRepeat("=", 80))
	fmt.Println("🎯 COMPLETE io PACKAGE GUIDE IN GO")
	fmt.Println(stringsRepeat("=", 80))

	// بخش 1: اینترفیس‌های اصلی
	demonstrateCoreInterfaces()

	// بخش 2: توابع مهم
	demonstrateImportantFunctions()

	// بخش 3: پیاده‌سازی سفارشی
	demonstrateCustomReadersWriters()

	// بخش 4: کار با فایل
	demonstrateFileIO()

	// بخش 5: بافر
	demonstrateBuffering()

	// بخش 6: فشرده‌سازی
	demonstrateCompression()

	// بخش 7: رمزنگاری
	demonstrateCryptoWithIO()

	// بخش 8: باینری
	demonstrateBinaryIO()

	// بخش 9: الگوهای پیشرفته
	demonstrateAdvancedPatterns()

	// بخش 10: اشتباهات رایج
	demonstrateCommonMistakes()

	// بخش 11: جمع‌بندی
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📚 io PACKAGE QUICK REFERENCE")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ INTERFACE      │ METHODS                      │ USE CASE        │")
	fmt.Println("├────────────────┼──────────────────────────────┼─────────────────┤")
	fmt.Println("│ io.Reader      │ Read(p []byte) (n int, err)  │ Read data       │")
	fmt.Println("│ io.Writer      │ Write(p []byte) (n int, err) │ Write data      │")
	fmt.Println("│ io.Closer      │ Close() error                │ Close resource  │")
	fmt.Println("│ io.Seeker      │ Seek(offset, whence) (int64) │ Jump position   │")
	fmt.Println("│ io.ReaderFrom  │ ReadFrom(r Reader) (n int64) │ Read from       │")
	fmt.Println("│ io.WriterTo    │ WriteTo(w Writer) (n int64)  │ Write to        │")
	fmt.Println("└────────────────┴──────────────────────────────┴─────────────────┘")

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ FUNCTION               │ PURPOSE                                │")
	fmt.Println("├────────────────────────┼────────────────────────────────────────┤")
	fmt.Println("│ io.Copy(dst, src)      │ Copy from Reader to Writer            │")
	fmt.Println("│ io.ReadAll(r)          │ Read everything into []byte           │")
	fmt.Println("│ io.ReadFull(r, buf)    │ Read exactly len(buf) bytes           │")
	fmt.Println("│ io.LimitReader(r, n)   │ Limit reading to n bytes              │")
	fmt.Println("│ io.TeeReader(r, w)     │ Read and copy simultaneously         │")
	fmt.Println("│ io.MultiReader(...)    │ Combine multiple readers              │")
	fmt.Println("│ io.MultiWriter(...)    │ Write to multiple writers             │")
	fmt.Println("│ io.Pipe()              │ Create synchronous in-memory pipe     │")
	fmt.Println("└────────────────────────┴────────────────────────────────────────┘")

	fmt.Println("\n💡 GOLDEN RULES:")
	fmt.Println("  1. Always check errors, especially io.EOF")
	fmt.Println("  2. Use io.ReadFull when you need exact bytes")
	fmt.Println("  3. Prefer io.Copy over manual reading/writing loops")
	fmt.Println("  4. Close resources with defer")
	fmt.Println("  5. Use bytes.Buffer and strings.Builder for efficient building")
	fmt.Println("  6. Implement small interfaces (Reader, Writer) for reusability")
	fmt.Println("  7. Use io.Reader/io.Writer in APIs, not concrete types")
	fmt.Println("  8. Preallocate buffers when size is known")
	fmt.Println("  9. Use io.Discard to throw away data")
	fmt.Println("  10. Combine interfaces with embedding for flexibility")
}

// ============================================================================
// بخش 12: توابع کمکی
// ============================================================================

func stringsRepeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
