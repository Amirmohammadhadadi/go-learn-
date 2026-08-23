// ============================================================================
// FILE: compress_complete_guide.go
// TITLE: راهنمای کامل پکیج compress در Go - فشرده‌سازی و از حالت فشرده خارج کردن
// HOW TO RUN: go run compress_complete_guide.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - پکیج‌های compress چیستند؟
// ============================================================================
//
// Go پکیج‌های متعددی برای فشرده‌سازی (compression) ارائه می‌دهد:
//
// 1. compress/gzip: فشرده‌سازی به فرمت GNU zip (معروف‌ترین فرمت)
// 2. compress/zlib: فشرده‌سازی به فرمت zlib (مورد استفاده در TLS/SSL)
// 3. compress/bzip2: فشرده‌سازی به فرمت bzip2 (فقط خواندن، بدون نوشتن)
// 4. compress/flate: پیاده‌سازی پایه الگوریتم DEFLATE
// 5. compress/lzw: فشرده‌سازی LZW (مورد استفاده در GIF و PDF)
//
// قانون طلایی:
// "برای فشرده‌سازی عمومی از gzip استفاده کن. برای APIها و پروتکل‌ها از zlib.
//  همیشه defer writer.Close() را فراموش نکن تا داده‌ها به درستی فلاش شوند."
// ============================================================================

package __internal_packages

import (
	"bytes"
	"compress/bzip2"
	"compress/flate"
	"compress/gzip"
	"compress/lzw"
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

func main() {
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("🎯 COMPLETE compress PACKAGES GUIDE IN GO")
	fmt.Println("Gzip, Zlib, Bzip2, Flate, LZW - Compression & Decompression")
	fmt.Println(strings.Repeat("=", 80))

	// ============================================================================
	// بخش 1: Gzip (GNU zip) - محبوب‌ترین فرمت فشرده‌سازی
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📦 SECTION 1: GZIP COMPRESSION")
	fmt.Println(strings.Repeat("=", 80))

	// 1.1 فشرده‌سازی با Gzip
	fmt.Println("\n--- 1.1 Gzip Compression ---")

	originalData := []byte("This is the original data that we want to compress. " +
		"It contains repetitive patterns and can be compressed efficiently. " +
		"Go is a great language for compression and decompression operations.")

	var gzipBuf bytes.Buffer

	// ایجاد Gzip writer
	gzipWriter := gzip.NewWriter(&gzipBuf)

	// تنظیم هدرهای Gzip (اختیاری)
	gzipWriter.Name = "example.txt"
	gzipWriter.Comment = "Example compression"
	gzipWriter.ModTime = time.Now()

	// نوشتن داده
	_, err := gzipWriter.Write(originalData)
	if err != nil {
		fmt.Printf("  Error writing to gzip: %v\n", err)
		return
	}

	// همیشه Close کنید تا داده‌ها فلاش شوند
	gzipWriter.Close()

	compressedSize := gzipBuf.Len()
	fmt.Printf("  Original size: %d bytes\n", len(originalData))
	fmt.Printf("  Compressed size: %d bytes (%.1f%%)\n",
		compressedSize, float64(compressedSize)/float64(len(originalData))*100)

	// 1.2 از حالت فشرده خارج کردن Gzip
	fmt.Println("\n--- 1.2 Gzip Decompression ---")

	gzipReader, err := gzip.NewReader(&gzipBuf)
	if err != nil {
		fmt.Printf("  Error creating gzip reader: %v\n", err)
		return
	}
	defer gzipReader.Close()

	// خواندن هدرها
	fmt.Printf("  Original filename: %s\n", gzipReader.Name)
	fmt.Printf("  Comment: %s\n", gzipReader.Comment)
	fmt.Printf("  Modified time: %v\n", gzipReader.ModTime)

	// خواندن داده
	decompressed, err := io.ReadAll(gzipReader)
	if err != nil {
		fmt.Printf("  Error reading decompressed data: %v\n", err)
		return
	}

	fmt.Printf("  Decompressed size: %d bytes\n", len(decompressed))
	fmt.Printf("  Data matches: %v\n", bytes.Equal(originalData, decompressed))

	// 1.3 Gzip با سطوح مختلف فشرده‌سازی
	fmt.Println("\n--- 1.3 Gzip Compression Levels ---")

	testData := strings.Repeat("ABCDEFGHIJKLMNOPQRSTUVWXYZ", 100)

	levels := []struct {
		level int
		name  string
	}{
		{gzip.NoCompression, "No Compression"},
		{gzip.BestSpeed, "Best Speed"},
		{gzip.BestCompression, "Best Compression"},
		{gzip.DefaultCompression, "Default"},
	}

	for _, l := range levels {
		var buf bytes.Buffer
		writer, _ := gzip.NewWriterLevel(&buf, l.level)
		writer.Write([]byte(testData))
		writer.Close()

		fmt.Printf("  %s: %d bytes (%.1f%% of original)\n",
			l.name, buf.Len(), float64(buf.Len())/float64(len(testData))*100)
	}

	// 1.4 خواندن مستقیم از فایل Gzip
	fmt.Println("\n--- 1.4 Reading from Gzip File ---")

	// ایجاد فایل Gzip تست
	gzipFile, err := os.CreateTemp("", "test*.gz")
	if err != nil {
		fmt.Printf("  Error creating temp file: %v\n", err)
		return
	}
	defer os.Remove(gzipFile.Name())

	gzWriter := gzip.NewWriter(gzipFile)
	gzWriter.Write([]byte("Data in gzip file"))
	gzWriter.Close()
	gzipFile.Close()

	// خواندن از فایل Gzip
	gzFileReader, _ := os.Open(gzipFile.Name())
	gzReader, _ := gzip.NewReader(gzFileReader)

	content, _ := io.ReadAll(gzReader)
	fmt.Printf("  Content from gzip file: %q\n", string(content))

	gzReader.Close()
	gzFileReader.Close()

	// ============================================================================
	// بخش 2: Zlib - فرمت استاندارد برای پروتکل‌ها
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📦 SECTION 2: ZLIB COMPRESSION")
	fmt.Println(strings.Repeat("=", 80))

	// 2.1 فشرده‌سازی با Zlib
	fmt.Println("\n--- 2.1 Zlib Compression ---")

	zlibData := []byte("Data to compress with zlib format. " +
		"Zlib is commonly used in TLS/SSL and PNG images.")

	var zlibBuf bytes.Buffer

	// ایجاد Zlib writer (سطوح: BestSpeed, BestCompression, DefaultCompression, NoCompression)
	zlibWriter, err := zlib.NewWriterLevel(&zlibBuf, zlib.BestCompression)
	if err != nil {
		fmt.Printf("  Error creating zlib writer: %v\n", err)
		return
	}

	zlibWriter.Write(zlibData)
	zlibWriter.Close()

	fmt.Printf("  Original size: %d bytes\n", len(zlibData))
	fmt.Printf("  Zlib compressed: %d bytes (%.1f%%)\n",
		zlibBuf.Len(), float64(zlibBuf.Len())/float64(len(zlibData))*100)

	// 2.2 از حالت فشرده خارج کردن Zlib
	fmt.Println("\n--- 2.2 Zlib Decompression ---")

	zlibReader, err := zlib.NewReader(&zlibBuf)
	if err != nil {
		fmt.Printf("  Error creating zlib reader: %v\n", err)
		return
	}
	defer zlibReader.Close()

	zlibDecompressed, err := io.ReadAll(zlibReader)
	if err != nil {
		fmt.Printf("  Error reading decompressed data: %v\n", err)
		return
	}

	fmt.Printf("  Decompressed size: %d bytes\n", len(zlibDecompressed))
	fmt.Printf("  Data matches: %v\n", bytes.Equal(zlibData, zlibDecompressed))

	// 2.3 Zlib با تنظیمات Dictionary (برای بهبود فشرده‌سازی)
	fmt.Println("\n--- 2.3 Zlib with Dictionary ---")

	// داده‌های تکراری
	repetitiveData := []byte("HTTP/1.1 200 OK\r\nContent-Type: text/html\r\n\r\n")
	htmlData := []byte("HTTP/1.1 200 OK\r\nContent-Type: text/html\r\n\r\n<html><body>Hello</body></html>")

	var dictBuf bytes.Buffer
	dictWriter, _ := zlib.NewWriterLevelDict(&dictBuf, zlib.BestCompression, repetitiveData)
	dictWriter.Write(htmlData)
	dictWriter.Close()

	fmt.Printf("  With dictionary: %d bytes\n", dictBuf.Len())

	var noDictBuf bytes.Buffer
	noDictWriter, _ := zlib.NewWriterLevel(&noDictBuf, zlib.BestCompression)
	noDictWriter.Write(htmlData)
	noDictWriter.Close()

	fmt.Printf("  Without dictionary: %d bytes\n", noDictBuf.Len())
	fmt.Printf("  Improvement: %.1f%%\n",
		(1-float64(dictBuf.Len())/float64(noDictBuf.Len()))*100)

	// ============================================================================
	// بخش 3: Bzip2 - فشرده‌سازی قوی (فقط خواندن)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📦 SECTION 3: BZIP2 COMPRESSION (Read-only)")
	fmt.Println(strings.Repeat("=", 80))

	// 3.1 خواندن فایل Bzip2 (نوشتن مستقیم در Go پشتیبانی نمی‌شود)
	fmt.Println("\n--- 3.1 Bzip2 Decompression ---")

	// ایجاد داده Bzip2 با استفاده از ابزار خارجی (در مثال، از داده شبیه‌سازی استفاده می‌کنیم)
	// در عمل از فایل‌های bz2 موجود استفاده کنید
	bzip2Data := []byte("BZ2 data would be here - in practice, read from .bz2 files")

	// شبیه‌سازی خواندن Bzip2 (در واقعیت از bzip2.NewReader استفاده کنید)
	fmt.Println("  Note: bzip2 only supports decompression, not compression in standard library")
	fmt.Printf("  To read bzip2: reader := bzip2.NewReader(file)\n")

	// مثال واقعی (کامنت شده):
	/*
		file, _ := os.Open("data.bz2")
		reader := bzip2.NewReader(file)
		decompressed, _ := io.ReadAll(reader)
		fmt.Printf("  Decompressed: %s\n", decompressed)
	*/

	// ============================================================================
	// بخش 4: Flate - پیاده‌سازی پایه DEFLATE
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📦 SECTION 4: FLATE COMPRESSION (DEFLATE)")
	fmt.Println(strings.Repeat("=", 80))

	// 4.1 فشرده‌سازی با Flate
	fmt.Println("\n--- 4.1 Flate Compression ---")

	flateData := []byte("Data to compress using DEFLATE algorithm. " +
		"This is the same algorithm used by gzip and zlib.")

	var flateBuf bytes.Buffer

	// ایجاد Flate writer با سطوح مختلف
	flateWriter, err := flate.NewWriter(&flateBuf, flate.BestCompression)
	if err != nil {
		fmt.Printf("  Error creating flate writer: %v\n", err)
		return
	}

	flateWriter.Write(flateData)
	flateWriter.Close()

	fmt.Printf("  Original size: %d bytes\n", len(flateData))
	fmt.Printf("  Flate compressed: %d bytes (%.1f%%)\n",
		flateBuf.Len(), float64(flateBuf.Len())/float64(len(flateData))*100)

	// 4.2 از حالت فشرده خارج کردن Flate
	fmt.Println("\n--- 4.2 Flate Decompression ---")

	flateReader := flate.NewReader(&flateBuf)
	defer flateReader.Close()

	flateDecompressed, err := io.ReadAll(flateReader)
	if err != nil {
		fmt.Printf("  Error reading decompressed data: %v\n", err)
		return
	}

	fmt.Printf("  Decompressed size: %d bytes\n", len(flateDecompressed))
	fmt.Printf("  Data matches: %v\n", bytes.Equal(flateData, flateDecompressed))

	// 4.3 Flate سطوح مختلف فشرده‌سازی
	fmt.Println("\n--- 4.3 Flate Compression Levels ---")

	flateLevels := []struct {
		level int
		name  string
	}{
		{flate.NoCompression, "No Compression"},
		{flate.BestSpeed, "Best Speed"},
		{flate.BestCompression, "Best Compression"},
		{flate.DefaultCompression, "Default"},
	}

	flateTestData := strings.Repeat("1234567890", 100)

	for _, l := range flateLevels {
		var buf bytes.Buffer
		writer, _ := flate.NewWriter(&buf, l.level)
		writer.Write([]byte(flateTestData))
		writer.Close()

		fmt.Printf("  %s: %d bytes\n", l.name, buf.Len())
	}

	// ============================================================================
	// بخش 5: LZW - فشرده‌سازی برای GIF و PDF
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📦 SECTION 5: LZW COMPRESSION (GIF, PDF, TIFF)")
	fmt.Println(strings.Repeat("=", 80))

	// 5.1 فشرده‌سازی با LZW
	fmt.Println("\n--- 5.1 LZW Compression ---")

	lzwData := []byte("LZW compression is used in GIF images and PDF files. " +
		"It works by building a dictionary of repeated sequences.")

	var lzwBuf bytes.Buffer

	// ایجاد LZW writer
	// order: LSB (Least Significant Byte first) or MSB
	// litWidth: number of literal bits (2-8)
	lzwWriter := lzw.NewWriter(&lzwBuf, lzw.LSB, 8)

	lzwWriter.Write(lzwData)
	lzwWriter.Close()

	fmt.Printf("  Original size: %d bytes\n", len(lzwData))
	fmt.Printf("  LZW compressed: %d bytes (%.1f%%)\n",
		lzwBuf.Len(), float64(lzwBuf.Len())/float64(len(lzwData))*100)

	// 5.2 از حالت فشرده خارج کردن LZW
	fmt.Println("\n--- 5.2 LZW Decompression ---")

	lzwReader := lzw.NewReader(&lzwBuf, lzw.LSB, 8)
	defer lzwReader.Close()

	lzwDecompressed, err := io.ReadAll(lzwReader)
	if err != nil {
		fmt.Printf("  Error reading decompressed data: %v\n", err)
		return
	}

	fmt.Printf("  Decompressed size: %d bytes\n", len(lzwDecompressed))
	fmt.Printf("  Data matches: %v\n", bytes.Equal(lzwData, lzwDecompressed))

	// 5.3 LZW با ترتیب بایت مختلف
	fmt.Println("\n--- 5.3 LZW Byte Order ---")

	original := []byte("Test data for LZW with different byte orders")

	// LSB (Least Significant Byte first)
	var lzwLSB bytes.Buffer
	writerLSB := lzw.NewWriter(&lzwLSB, lzw.LSB, 8)
	writerLSB.Write(original)
	writerLSB.Close()

	// MSB (Most Significant Byte first)
	var lzwMSB bytes.Buffer
	writerMSB := lzw.NewWriter(&lzwMSB, lzw.MSB, 8)
	writerMSB.Write(original)
	writerMSB.Close()

	fmt.Printf("  LSB compressed size: %d bytes\n", lzwLSB.Len())
	fmt.Printf("  MSB compressed size: %d bytes\n", lzwMSB.Len())
	fmt.Printf("  Files are different: %v\n", !bytes.Equal(lzwLSB.Bytes(), lzwMSB.Bytes()))

	// ============================================================================
	// بخش 6: مقایسه فرمت‌های مختلف فشرده‌سازی
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📊 SECTION 6: COMPRESSION FORMAT COMPARISON")
	fmt.Println(strings.Repeat("=", 80))

	// داده تست
	testString := `Lorem ipsum dolor sit amet, consectetur adipiscing elit. 
	Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. 
	Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris. 
	Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur.
	Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.`
	testData := []byte(strings.Repeat(testString, 10))

	fmt.Println("\n--- 6.1 Compression Ratio Comparison ---")
	fmt.Printf("  Original size: %d bytes\n", len(testData))

	// Gzip
	var gzipCompare bytes.Buffer
	gzipCmp := gzip.NewWriter(&gzipCompare)
	gzipCmp.Write(testData)
	gzipCmp.Close()
	fmt.Printf("  Gzip: %d bytes (%.1f%%)\n",
		gzipCompare.Len(), float64(gzipCompare.Len())/float64(len(testData))*100)

	// Zlib
	var zlibCompare bytes.Buffer
	zlibCmp, _ := zlib.NewWriterLevel(&zlibCompare, zlib.BestCompression)
	zlibCmp.Write(testData)
	zlibCmp.Close()
	fmt.Printf("  Zlib: %d bytes (%.1f%%)\n",
		zlibCompare.Len(), float64(zlibCompare.Len())/float64(len(testData))*100)

	// Flate
	var flateCompare bytes.Buffer
	flateCmp, _ := flate.NewWriter(&flateCompare, flate.BestCompression)
	flateCmp.Write(testData)
	flateCmp.Close()
	fmt.Printf("  Flate: %d bytes (%.1f%%)\n",
		flateCompare.Len(), float64(flateCompare.Len())/float64(len(testData))*100)

	// LZW
	var lzwCompare bytes.Buffer
	lzwCmp := lzw.NewWriter(&lzwCompare, lzw.LSB, 8)
	lzwCmp.Write(testData)
	lzwCmp.Close()
	fmt.Printf("  LZW: %d bytes (%.1f%%)\n",
		lzwCompare.Len(), float64(lzwCompare.Len())/float64(len(testData))*100)

	// ============================================================================
	// بخش 7: کاربردهای عملی
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("💼 SECTION 7: PRACTICAL APPLICATIONS")
	fmt.Println(strings.Repeat("=", 80))

	// 7.1 فشرده‌سازی HTTP پاسخ
	fmt.Println("\n--- 7.1 HTTP Response Compression ---")
	fmt.Println("  In HTTP servers, you can compress responses with gzip:")
	fmt.Println("  w.Header().Set(\"Content-Encoding\", \"gzip\")")
	fmt.Println("  gw := gzip.NewWriter(w)")
	fmt.Println("  gw.Write(responseData)")
	fmt.Println("  gw.Close()")

	// 7.2 فشرده‌سازی فایل لاگ
	fmt.Println("\n--- 7.2 Log File Compression ---")

	// شبیه‌سازی فشرده‌سازی لاگ
	logData := strings.Repeat("2024-01-15 10:00:00 INFO: Application started\n", 1000)

	var compressedLog bytes.Buffer
	logGzip := gzip.NewWriter(&compressedLog)
	logGzip.Write([]byte(logData))
	logGzip.Close()

	fmt.Printf("  Log size: %d bytes\n", len(logData))
	fmt.Printf("  Compressed: %d bytes (%.1f%%)\n",
		compressedLog.Len(), float64(compressedLog.Len())/float64(len(logData))*100)

	// 7.3 فشرده‌سازی در حافظه با بافر
	fmt.Println("\n--- 7.3 In-Memory Compression ---")

	type CompressedBuffer struct {
		buf    bytes.Buffer
		writer *gzip.Writer
	}

	cb := &CompressedBuffer{}
	cb.writer = gzip.NewWriter(&cb.buf)

	// نوشتن داده
	cb.writer.Write([]byte("First chunk of data\n"))
	cb.writer.Write([]byte("Second chunk of data\n"))
	cb.writer.Write([]byte("Third chunk of data\n"))
	cb.writer.Close()

	fmt.Printf("  Compressed buffer size: %d bytes\n", cb.buf.Len())

	// خواندن داده
	reader, _ := gzip.NewReader(&cb.buf)
	decompressed, _ := io.ReadAll(reader)
	fmt.Printf("  Decompressed: %q\n", string(decompressed[:50])+"...")

	// 7.4 تشخیص نوع فشرده‌سازی از روی magic bytes
	fmt.Println("\n--- 7.4 Compression Type Detection ---")

	detectCompression := func(data []byte) string {
		if len(data) < 2 {
			return "Unknown"
		}
		switch {
		case data[0] == 0x1F && data[1] == 0x8B:
			return "Gzip"
		case data[0] == 0x78 && (data[1] == 0x01 || data[1] == 0x5E || data[1] == 0x9C || data[1] == 0xDA):
			return "Zlib"
		case data[0] == 0x42 && data[1] == 0x5A:
			return "Bzip2"
		default:
			return "Unknown"
		}
	}

	// تست magic bytes
	fmt.Printf("  Gzip magic: %s\n", detectCompression(gzipCompare.Bytes()))
	fmt.Printf("  Zlib magic: %s\n", detectCompression(zlibCompare.Bytes()))
	fmt.Printf("  Plain text: %s\n", detectCompression([]byte("plain text")))

	// ============================================================================
	// بخش 8: اشتباهات رایج
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("❌ SECTION 8: COMMON MISTAKES")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println("\n❌ Mistake 1: Forgetting to Close writer")
	fmt.Println("   gz := gzip.NewWriter(&buf)")
	fmt.Println("   gz.Write(data)  // data not flushed!")
	fmt.Println("   ✅ defer gz.Close() or call gz.Close()")

	fmt.Println("\n❌ Mistake 2: Not checking Close errors")
	fmt.Println("   gz.Close()  // error ignored")
	fmt.Println("   ✅ Always check Close error (important for data integrity)")

	fmt.Println("\n❌ Mistake 3: Using gzip for small data")
	fmt.Println("   gzip adds header (18 bytes overhead)")
	fmt.Println("   ✅ Use for data > 200 bytes, else overhead may exceed benefit")

	fmt.Println("\n❌ Mistake 4: Not closing reader")
	fmt.Println("   gz, _ := gzip.NewReader(r)")
	fmt.Println("   io.Copy(w, gz)  // resources not released")
	fmt.Println("   ✅ defer gz.Close()")

	fmt.Println("\n❌ Mistake 5: Assuming bzip2 supports writing")
	fmt.Println("   compress/bzip2 only supports reading!")
	fmt.Println("   ✅ Use external tool or other formats for writing")

	fmt.Println("\n❌ Mistake 6: Incorrect LZW parameters")
	fmt.Println("   lzw.NewWriter(w, lzw.LSB, 8)  // wrong for GIF")
	fmt.Println("   ✅ GIF uses LSB with litWidth=8")
	fmt.Println("   ✅ PDF uses MSB with litWidth=8")

	fmt.Println("\n❌ Mistake 7: Not handling dictionary in zlib")
	fmt.Println("   // can't decompress without same dictionary")
	fmt.Println("   ✅ Use NewWriterLevelDict and NewReaderDict with same dict")

	// ============================================================================
	// بخش 9: جمع‌بندی و جدول مرجع
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📚 SECTION 9: QUICK REFERENCE")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ PACKAGE      │ USE CASE                      │ SUPPORT         │")
	fmt.Println("├──────────────┼───────────────────────────────┼─────────────────┤")
	fmt.Println("│ compress/gzip│ General purpose, web, files   │ Read & Write    │")
	fmt.Println("│ compress/zlib│ TLS/SSL, PNG, protocols       │ Read & Write    │")
	fmt.Println("│ compress/bzip2│ High compression ratio       │ Read only       │")
	fmt.Println("│ compress/flate│ Base DEFLATE                 │ Read & Write    │")
	fmt.Println("│ compress/lzw │ GIF, PDF, TIFF               │ Read & Write    │")
	fmt.Println("└──────────────┴───────────────────────────────┴─────────────────┘")

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ GZIP Compression Levels                                        │")
	fmt.Println("├─────────────────────────────────────────────────────────────────┤")
	fmt.Println("│ gzip.NoCompression      - No compression (fastest)            │")
	fmt.Println("│ gzip.BestSpeed          - Fastest compression                 │")
	fmt.Println("│ gzip.BestCompression    - Best compression (slowest)          │")
	fmt.Println("│ gzip.DefaultCompression - Balance between speed and size      │")
	fmt.Println("└─────────────────────────────────────────────────────────────────┘")

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ MAGIC BYTES FOR COMPRESSION DETECTION                          │")
	fmt.Println("├─────────────────────────────────────────────────────────────────┤")
	fmt.Println("│ Gzip:  0x1F, 0x8B                                             │")
	fmt.Println("│ Zlib:  0x78, 0x01/0x5E/0x9C/0xDA                               │")
	fmt.Println("│ Bzip2: 0x42, 0x5A                                             │")
	fmt.Println("│ Flate: no fixed magic (raw DEFLATE)                           │")
	fmt.Println("│ LZW:   no fixed magic                                         │")
	fmt.Println("└─────────────────────────────────────────────────────────────────┘")

	fmt.Println("\n💡 GOLDEN RULES:")
	fmt.Println("  1. Always Close writers to flush remaining data")
	fmt.Println("  2. Check Close errors for data integrity")
	fmt.Println("  3. Use gzip for general purpose compression")
	fmt.Println("  4. Use zlib for protocol-level compression")
	fmt.Println("  5. Defer Close on readers to prevent resource leaks")
	fmt.Println("  6. Choose compression level based on your needs")
	fmt.Println("  7. Small data (<200 bytes) may not benefit from compression")
	fmt.Println("  8. LZW parameters must match the format (GIF vs PDF)")
	fmt.Println("  9. bzip2 in standard library is read-only")
	fmt.Println("  10. Test compression ratio before using in production")

	fmt.Println("\n🎯 PERFORMANCE TIPS:")
	fmt.Println("  • Use BestSpeed for real-time compression")
	fmt.Println("  • Use BestCompression for storage/archives")
	fmt.Println("  • Reuse compressors with Reset() for multiple data blocks")
	fmt.Println("  • Pool compressors with sync.Pool for high concurrency")
	fmt.Println("  • Compress in goroutines for large datasets")
	fmt.Println("  • Use io.Pipe for streaming compression")
}

// تابع کمکی برای تکرار رشته
func stringsRepeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}
