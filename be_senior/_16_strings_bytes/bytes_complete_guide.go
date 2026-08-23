package __internal_packages

// ============================================================================
// FILE: bytes_complete_guide.go
// TITLE: راهنمای کامل پکیج bytes در Go - تمام توابع با مثال
// HOW TO RUN: go run bytes_complete_guide.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - پکیج bytes چیست؟
// ============================================================================
//
// پکیج bytes عملیات روی slice های بایت ([]byte) را ارائه می‌دهد.
// این پکیج بسیار شبیه به پکیج strings است اما به جای رشته روی بایت‌ها کار می‌کند.
//
// کاربردهای اصلی:
// 1. پردازش بافر و داده‌های باینری
// 2. دستکاری کارآمد رشته‌ها (بدون کپی)
// 3. خواندن و نوشتن از بافر
// 4. عملیات جستجو و جایگزینی روی داده‌های باینری
//
// قانون طلایی:
// "هر وقت نیاز به عملیات روی داده‌های باینری یا بافر داری، از پکیج bytes استفاده کن.
//  این پکیج کارآمدترین روش برای دستکاری بایت‌ها در Go است."
// ============================================================================

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("🎯 COMPLETE bytes PACKAGE GUIDE IN GO")
	fmt.Println("All functions with practical examples")
	fmt.Println(strings.Repeat("=", 80))

	// ============================================================================
	// بخش 1: توابع مقایسه (Comparison)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("⚖️ SECTION 1: COMPARISON FUNCTIONS")
	fmt.Println(strings.Repeat("=", 80))

	// 1.1 Compare - مقایسه دو slice بایت
	fmt.Println("\n--- 1.1 bytes.Compare ---")
	// Compare returns an integer comparing two byte slices lexicographically.
	a := []byte("abc")
	b := []byte("abd")
	c := []byte("abc")

	fmt.Printf("  Compare(%q, %q): %d\n", a, b, bytes.Compare(a, b))
	fmt.Printf("  Compare(%q, %q): %d\n", b, a, bytes.Compare(b, a))
	fmt.Printf("  Compare(%q, %q): %d\n", a, c, bytes.Compare(a, c))

	// 1.2 Equal - بررسی برابری دو slice
	fmt.Println("\n--- 1.2 bytes.Equal ---")
	// Equal reports whether a and b are the same length and contain the same bytes.
	d := []byte("hello")
	e := []byte("hello")
	f := []byte("world")

	fmt.Printf("  Equal(%q, %q): %v\n", d, e, bytes.Equal(d, e))
	fmt.Printf("  Equal(%q, %q): %v\n", d, f, bytes.Equal(d, f))

	// 1.3 EqualFold - مقایسه بدون در نظر گرفتن حروف بزرگ/کوچک
	fmt.Println("\n--- 1.3 bytes.EqualFold ---")
	// EqualFold reports whether s and t, interpreted as UTF-8 strings, are equal under Unicode case-folding.
	g := []byte("GoLang")
	h := []byte("golang")

	fmt.Printf("  EqualFold(%q, %q): %v\n", g, h, bytes.EqualFold(g, h))

	// ============================================================================
	// بخش 2: توابع جستجو (Search)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🔍 SECTION 2: SEARCH FUNCTIONS")
	fmt.Println(strings.Repeat("=", 80))

	data := []byte("The quick brown fox jumps over the lazy dog")

	// 2.1 Contains - آیا slice شامل زیرslice است؟
	fmt.Println("\n--- 2.1 bytes.Contains ---")
	// Contains reports whether subslice is within b.
	fmt.Printf("  Contains(%q, %q): %v\n", data, []byte("fox"), bytes.Contains(data, []byte("fox")))
	fmt.Printf("  Contains(%q, %q): %v\n", data, []byte("cat"), bytes.Contains(data, []byte("cat")))

	// 2.2 ContainsAny - آیا شامل هر یک از کاراکترهاست؟
	fmt.Println("\n--- 2.2 bytes.ContainsAny ---")
	// ContainsAny reports whether any of the UTF-8-encoded code points in chars are within b.
	fmt.Printf("  ContainsAny(%q, %q): %v\n", data, "xyz", bytes.ContainsAny(data, "xyz"))
	fmt.Printf("  ContainsAny(%q, %q): %v\n", data, "123", bytes.ContainsAny(data, "123"))

	// 2.3 ContainsRune - آیا شامل رون خاص است؟
	fmt.Println("\n--- 2.3 bytes.ContainsRune ---")
	// ContainsRune reports whether the rune is contained in the UTF-8-encoded byte slice b.
	fmt.Printf("  ContainsRune(%q, 'e'): %v\n", data, 'e', bytes.ContainsRune(data, 'e'))
	fmt.Printf("  ContainsRune(%q, '世'): %v\n", data, '世', bytes.ContainsRune(data, '世'))

	// 2.4 Count - تعداد دفعات تکرار زیرslice
	fmt.Println("\n--- 2.4 bytes.Count ---")
	// Count counts the number of non-overlapping instances of sep in s.
	banana := []byte("banana")
	fmt.Printf("  Count(%q, %q): %d\n", banana, []byte("na"), bytes.Count(banana, []byte("na")))
	fmt.Printf("  Count(%q, %q): %d\n", banana, []byte("a"), bytes.Count(banana, []byte("a")))

	// 2.5 Index - اولین موقعیت زیرslice
	fmt.Println("\n--- 2.5 bytes.Index ---")
	// Index returns the index of the first instance of sep in s, or -1 if sep is not present in s.
	fmt.Printf("  Index(%q, %q): %d\n", data, []byte("fox"), bytes.Index(data, []byte("fox")))
	fmt.Printf("  Index(%q, %q): %d\n", data, []byte("cat"), bytes.Index(data, []byte("cat")))

	// 2.6 IndexAny - اولین موقعیت هر یک از کاراکترها
	fmt.Println("\n--- 2.6 bytes.IndexAny ---")
	// IndexAny interprets s as a sequence of UTF-8-encoded code points.
	// It returns the byte index of the first occurrence in s of any of the Unicode code points in chars.
	fmt.Printf("  IndexAny(%q, %q): %d\n", data, "aeiou", bytes.IndexAny(data, "aeiou"))

	// 2.7 IndexByte - موقعیت یک بایت خاص
	fmt.Println("\n--- 2.7 bytes.IndexByte ---")
	// IndexByte returns the index of the first instance of c in b, or -1 if c is not present in b.
	fmt.Printf("  IndexByte(%q, 'o'): %d\n", data, 'o', bytes.IndexByte(data, 'o'))
	fmt.Printf("  IndexByte(%q, 'z'): %d\n", data, 'z', bytes.IndexByte(data, 'z'))

	// 2.8 IndexFunc - موقعیت کاراکتری که شرط را برآورده می‌کند
	fmt.Println("\n--- 2.8 bytes.IndexFunc ---")
	// IndexFunc interprets s as a sequence of UTF-8-encoded code points.
	// It returns the byte index of the first Unicode code point in s that satisfies f(c), or -1 if none do.
	isUpper := func(r rune) bool { return r >= 'A' && r <= 'Z' }
	fmt.Printf("  IndexFunc(%q, isUpper): %d (اولین حرف بزرگ)\n", data, bytes.IndexFunc(data, isUpper))

	// 2.9 IndexRune - موقعیت یک رون خاص
	fmt.Println("\n--- 2.9 bytes.IndexRune ---")
	// IndexRune interprets s as a sequence of UTF-8-encoded code points.
	// It returns the byte index of the first occurrence in s of the given rune.
	fmt.Printf("  IndexRune(%q, '世'): %d\n", []byte("Hello世World"), '世', bytes.IndexRune([]byte("Hello世World"), '世'))

	// 2.10 LastIndex - آخرین موقعیت زیرslice
	fmt.Println("\n--- 2.10 bytes.LastIndex ---")
	// LastIndex returns the index of the last instance of sep in s, or -1 if sep is not present in s.
	fmt.Printf("  LastIndex(%q, %q): %d\n", banana, []byte("na"), bytes.LastIndex(banana, []byte("na")))

	// 2.11 LastIndexAny - آخرین موقعیت هر یک از کاراکترها
	fmt.Println("\n--- 2.11 bytes.LastIndexAny ---")
	// LastIndexAny interprets s as a sequence of UTF-8-encoded code points.
	// It returns the byte index of the last occurrence in s of any of the Unicode code points in chars.
	fmt.Printf("  LastIndexAny(%q, %q): %d\n", data, "aeiou", bytes.LastIndexAny(data, "aeiou"))

	// 2.12 LastIndexByte - آخرین موقعیت یک بایت خاص
	fmt.Println("\n--- 2.12 bytes.LastIndexByte ---")
	// LastIndexByte returns the index of the last instance of c in b, or -1 if c is not present in b.
	fmt.Printf("  LastIndexByte(%q, 'a'): %d\n", banana, 'a', bytes.LastIndexByte(banana, 'a'))

	// 2.13 LastIndexFunc - آخرین موقعیت کاراکتری که شرط را برآورده می‌کند
	fmt.Println("\n--- 2.13 bytes.LastIndexFunc ---")
	// LastIndexFunc interprets s as a sequence of UTF-8-encoded code points.
	// It returns the byte index of the last Unicode code point in s that satisfies f(c), or -1 if none do.
	fmt.Printf("  LastIndexFunc(%q, isUpper): %d\n", data, bytes.LastIndexFunc(data, isUpper))

	// ============================================================================
	// بخش 3: توابع جایگزینی (Replacement)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🔄 SECTION 3: REPLACEMENT FUNCTIONS")
	fmt.Println(strings.Repeat("=", 80))

	// 3.1 Replace - جایگزینی با تعداد مشخص
	fmt.Println("\n--- 3.1 bytes.Replace ---")
	// Replace returns a copy of the slice s with the first n non-overlapping instances of old replaced by new.
	orig := []byte("oink oink oink")
	replaced := bytes.Replace(orig, []byte("oink"), []byte("moo"), 2)
	fmt.Printf("  Original: %q\n", orig)
	fmt.Printf("  Replace first 2: %q\n", replaced)

	// با n = -1 همه را جایگزین کن
	replacedAll := bytes.Replace(orig, []byte("oink"), []byte("moo"), -1)
	fmt.Printf("  Replace all: %q\n", replacedAll)

	// 3.2 ReplaceAll - جایگزینی همه موارد
	fmt.Println("\n--- 3.2 bytes.ReplaceAll ---")
	// ReplaceAll returns a copy of the slice s with all non-overlapping instances of old replaced by new.
	hello := []byte("hello world hello")
	replacedAll2 := bytes.ReplaceAll(hello, []byte("hello"), []byte("hi"))
	fmt.Printf("  ReplaceAll(%q, %q, %q): %q\n", hello, "hello", "hi", replacedAll2)

	// 3.3 NewReplacer - جایگزینی چندتایی با Replacer
	fmt.Println("\n--- 3.3 bytes.NewReplacer ---")
	// NewReplacer returns a new Replacer from a list of old, new string pairs.
	replacer := bytes.NewReplacer(
		[]byte("hello"), []byte("hi"),
		[]byte("world"), []byte("earth"),
		[]byte("good"), []byte("bad"),
	)
	text := []byte("hello world good morning")
	result := replacer.Replace(text)
	fmt.Printf("  Original: %q\n", text)
	fmt.Printf("  Replaced: %q\n", result)

	// ============================================================================
	// بخش 4: توابع تقسیم و اتصال (Split & Join)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("✂️ SECTION 4: SPLIT & JOIN FUNCTIONS")
	fmt.Println(strings.Repeat("=", 80))

	// 4.1 Split - تقسیم با جداکننده
	fmt.Println("\n--- 4.1 bytes.Split ---")
	// Split slices s into all subslices separated by sep and returns a slice of the subslices between those separators.
	csv := []byte("a,b,c,d,e")
	parts := bytes.Split(csv, []byte(","))
	fmt.Printf("  Split(%q, ','): %q\n", csv, parts)

	// 4.2 SplitN - تقسیم با حداکثر n قسمت
	fmt.Println("\n--- 4.2 bytes.SplitN ---")
	// SplitN slices s into subslices separated by sep and returns a slice of the subslices between those separators.
	partsN := bytes.SplitN(csv, []byte(","), 3)
	fmt.Printf("  SplitN(%q, ',', 3): %q\n", csv, partsN)

	// 4.3 SplitAfter - تقسیم با حفظ جداکننده
	fmt.Println("\n--- 4.3 bytes.SplitAfter ---")
	// SplitAfter slices s into all subslices after each instance of sep and returns a slice of those subslices.
	splitAfter := bytes.SplitAfter(csv, []byte(","))
	fmt.Printf("  SplitAfter(%q, ','): %q\n", csv, splitAfter)

	// 4.4 SplitAfterN - تقسیم با حداکثر n قسمت (حفظ جداکننده)
	fmt.Println("\n--- 4.4 bytes.SplitAfterN ---")
	// SplitAfterN slices s into subslices after each instance of sep and returns a slice of those subslices.
	splitAfterN := bytes.SplitAfterN(csv, []byte(","), 3)
	fmt.Printf("  SplitAfterN(%q, ',', 3): %q\n", csv, splitAfterN)

	// 4.5 Fields - تقسیم بر اساس whitespace
	fmt.Println("\n--- 4.5 bytes.Fields ---")
	// Fields interprets s as a sequence of UTF-8-encoded code points.
	// It splits the slice s around each instance of one or more consecutive white space characters.
	whitespace := []byte("hello   world  from   go")
	fields := bytes.Fields(whitespace)
	fmt.Printf("  Fields(%q): %q\n", whitespace, fields)

	// 4.6 FieldsFunc - تقسیم با تابع شرط
	fmt.Println("\n--- 4.6 bytes.FieldsFunc ---")
	// FieldsFunc interprets s as a sequence of UTF-8-encoded code points.
	// It splits the slice s at each run of code points c satisfying f(c) and returns a slice of subslices of s.
	isPunct := func(r rune) bool {
		return r == ',' || r == '.' || r == '!' || r == '?'
	}
	punctuation := []byte("hello,world!how.are?you")
	fieldsFunc := bytes.FieldsFunc(punctuation, isPunct)
	fmt.Printf("  FieldsFunc(%q, isPunct): %q\n", punctuation, fieldsFunc)

	// 4.7 Join - اتصال slice از بایت‌ها با جداکننده
	fmt.Println("\n--- 4.7 bytes.Join ---")
	// Join concatenates the elements of s to create a new byte slice. The separator sep is placed between elements.
	words := [][]byte{[]byte("hello"), []byte("world"), []byte("from"), []byte("go")}
	joined := bytes.Join(words, []byte(" "))
	fmt.Printf("  Join(%q, ' '): %q\n", words, joined)

	// ============================================================================
	// بخش 5: توابع حذف فاصله و کاراکترهای خاص (Trim)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("✂️ SECTION 5: TRIM FUNCTIONS")
	fmt.Println(strings.Repeat("=", 80))

	// 5.1 Trim - حذف کاراکترهای مشخص از دو طرف
	fmt.Println("\n--- 5.1 bytes.Trim ---")
	// Trim returns a subslice of s by slicing off all leading and trailing UTF-8-encoded code points contained in cutset.
	exclam := []byte("!!!Hello!!!")
	trimmed := bytes.Trim(exclam, "!")
	fmt.Printf("  Trim(%q, '!'): %q\n", exclam, trimmed)

	// 5.2 TrimSpace - حذف whitespace از دو طرف
	fmt.Println("\n--- 5.2 bytes.TrimSpace ---")
	// TrimSpace returns a subslice of s by slicing off all leading and trailing white space.
	spacey := []byte("  \t\nHello World\n\t  ")
	trimmedSpace := bytes.TrimSpace(spacey)
	fmt.Printf("  TrimSpace(%q): %q\n", spacey, trimmedSpace)

	// 5.3 TrimLeft - حذف کاراکترها از سمت چپ
	fmt.Println("\n--- 5.3 bytes.TrimLeft ---")
	// TrimLeft returns a subslice of s by slicing off all leading UTF-8-encoded code points contained in cutset.
	trimmedLeft := bytes.TrimLeft(exclam, "!")
	fmt.Printf("  TrimLeft(%q, '!'): %q\n", exclam, trimmedLeft)

	// 5.4 TrimRight - حذف کاراکترها از سمت راست
	fmt.Println("\n--- 5.4 bytes.TrimRight ---")
	// TrimRight returns a subslice of s by slicing off all trailing UTF-8-encoded code points that are contained in cutset.
	trimmedRight := bytes.TrimRight(exclam, "!")
	fmt.Printf("  TrimRight(%q, '!'): %q\n", exclam, trimmedRight)

	// 5.5 TrimPrefix - حذف پیشوند
	fmt.Println("\n--- 5.5 bytes.TrimPrefix ---")
	// TrimPrefix returns s without the provided leading prefix string. If s doesn't start with prefix, s is returned unchanged.
	prefix := []byte("Hello")
	withPrefix := []byte("HelloWorld")
	trimmedPrefix := bytes.TrimPrefix(withPrefix, prefix)
	fmt.Printf("  TrimPrefix(%q, %q): %q\n", withPrefix, prefix, trimmedPrefix)

	// 5.6 TrimSuffix - حذف پسوند
	fmt.Println("\n--- 5.6 bytes.TrimSuffix ---")
	// TrimSuffix returns s without the provided trailing suffix string. If s doesn't end with suffix, s is returned unchanged.
	suffix := []byte("World")
	withSuffix := []byte("HelloWorld")
	trimmedSuffix := bytes.TrimSuffix(withSuffix, suffix)
	fmt.Printf("  TrimSuffix(%q, %q): %q\n", withSuffix, suffix, trimmedSuffix)

	// 5.7 TrimFunc - حذف کاراکترهایی که شرط را برآورده می‌کنند
	fmt.Println("\n--- 5.7 bytes.TrimFunc ---")
	// TrimFunc returns a subslice of s by slicing off all leading and trailing UTF-8-encoded code points c that satisfy f(c).
	isExclam := func(r rune) bool { return r == '!' }
	trimmedFunc := bytes.TrimFunc(exclam, isExclam)
	fmt.Printf("  TrimFunc(%q, isExclam): %q\n", exclam, trimmedFunc)

	// 5.8 TrimLeftFunc - حذف از چپ با شرط
	fmt.Println("\n--- 5.8 bytes.TrimLeftFunc ---")
	// TrimLeftFunc returns a subslice of s by slicing off all leading UTF-8-encoded code points c that satisfy f(c).
	trimmedLeftFunc := bytes.TrimLeftFunc(exclam, isExclam)
	fmt.Printf("  TrimLeftFunc(%q, isExclam): %q\n", exclam, trimmedLeftFunc)

	// 5.9 TrimRightFunc - حذف از راست با شرط
	fmt.Println("\n--- 5.9 bytes.TrimRightFunc ---")
	// TrimRightFunc returns a subslice of s by slicing off all trailing UTF-8-encoded code points c that satisfy f(c).
	trimmedRightFunc := bytes.TrimRightFunc(exclam, isExclam)
	fmt.Printf("  TrimRightFunc(%q, isExclam): %q\n", exclam, trimmedRightFunc)

	// ============================================================================
	// بخش 6: توابع تبدیل حروف (Case Conversion)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🔠 SECTION 6: CASE CONVERSION FUNCTIONS")
	fmt.Println(strings.Repeat("=", 80))

	// 6.1 ToLower - تبدیل به حروف کوچک
	fmt.Println("\n--- 6.1 bytes.ToLower ---")
	// ToLower returns a copy of the byte slice s with all Unicode letters mapped to their lower case.
	upper := []byte("HELLO WORLD")
	lower := bytes.ToLower(upper)
	fmt.Printf("  ToLower(%q): %q\n", upper, lower)

	// 6.2 ToUpper - تبدیل به حروف بزرگ
	fmt.Println("\n--- 6.2 bytes.ToUpper ---")
	// ToUpper returns a copy of the byte slice s with all Unicode letters mapped to their upper case.
	lower2 := []byte("hello world")
	upper2 := bytes.ToUpper(lower2)
	fmt.Printf("  ToUpper(%q): %q\n", lower2, upper2)

	// 6.3 ToTitle - تبدیل به حروف عنوان
	fmt.Println("\n--- 6.3 bytes.ToTitle ---")
	// ToTitle returns a copy of the byte slice s with all Unicode letters mapped to their title case.
	title := []byte("hello world")
	titleCase := bytes.ToTitle(title)
	fmt.Printf("  ToTitle(%q): %q\n", title, titleCase)

	// 6.4 ToLowerSpecial - تبدیل با قوانین خاص
	fmt.Println("\n--- 6.4 bytes.ToLowerSpecial ---")
	// ToLowerSpecial returns a copy of the byte slice s with all Unicode letters mapped to their lower case using the case mapping specified by c.
	// (نیاز به unicode.SpecialCase دارد)

	// 6.5 ToUpperSpecial - تبدیل با قوانین خاص
	fmt.Println("\n--- 6.5 bytes.ToUpperSpecial ---")
	// ToUpperSpecial returns a copy of the byte slice s with all Unicode letters mapped to their upper case using the case mapping specified by c.

	// ============================================================================
	// بخش 7: توابع بررسی پیشوند و پسوند (Prefix & Suffix)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📌 SECTION 7: PREFIX & SUFFIX FUNCTIONS")
	fmt.Println(strings.Repeat("=", 80))

	// 7.1 HasPrefix - آیا با پیشوند مشخص شروع می‌شود؟
	fmt.Println("\n--- 7.1 bytes.HasPrefix ---")
	// HasPrefix tests whether the byte slice s begins with prefix.
	helloWorld := []byte("HelloWorld")
	fmt.Printf("  HasPrefix(%q, %q): %v\n", helloWorld, []byte("Hello"), bytes.HasPrefix(helloWorld, []byte("Hello")))
	fmt.Printf("  HasPrefix(%q, %q): %v\n", helloWorld, []byte("World"), bytes.HasPrefix(helloWorld, []byte("World")))

	// 7.2 HasSuffix - آیا با پسوند مشخص تمام می‌شود؟
	fmt.Println("\n--- 7.2 bytes.HasSuffix ---")
	// HasSuffix tests whether the byte slice s ends with suffix.
	fmt.Printf("  HasSuffix(%q, %q): %v\n", helloWorld, []byte("World"), bytes.HasSuffix(helloWorld, []byte("World")))
	fmt.Printf("  HasSuffix(%q, %q): %v\n", helloWorld, []byte("Hello"), bytes.HasSuffix(helloWorld, []byte("Hello")))

	// ============================================================================
	// بخش 8: بافر (bytes.Buffer) - کارآمدترین روش کار با بایت‌ها
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📦 SECTION 8: bytes.Buffer - Efficient Byte Buffer")
	fmt.Println(strings.Repeat("=", 80))

	// 8.1 ایجاد و نوشتن در Buffer
	fmt.Println("\n--- 8.1 Creating and Writing to Buffer ---")

	var buf bytes.Buffer

	// نوشتن رشته
	buf.WriteString("Hello ")
	buf.WriteString("World")

	// نوشتن بایت
	buf.WriteByte('!')

	// نوشتن slice بایت
	buf.Write([]byte(" Welcome"))

	// نوشتن رون
	buf.WriteRune('🎉')

	fmt.Printf("  Buffer content: %q\n", buf.String())
	fmt.Printf("  Buffer length: %d, capacity: %d\n", buf.Len(), buf.Cap())

	// 8.2 Grow - پیش‌تخصیص ظرفیت
	fmt.Println("\n--- 8.2 Buffer Grow (Pre-allocation) ---")

	var buf2 bytes.Buffer
	fmt.Printf("  Before Grow: len=%d, cap=%d\n", buf2.Len(), buf2.Cap())

	buf2.Grow(1024) // پیش‌تخصیص 1KB
	fmt.Printf("  After Grow: len=%d, cap=%d\n", buf2.Len(), buf2.Cap())

	// 8.3 خواندن از Buffer
	fmt.Println("\n--- 8.3 Reading from Buffer ---")

	buf3 := bytes.NewBufferString("Hello World")

	// ReadByte - خواندن یک بایت
	b, _ := buf3.ReadByte()
	fmt.Printf("  ReadByte: %c\n", b)

	// ReadString - خواندن تا جداکننده
	line, _ := buf3.ReadString(' ')
	fmt.Printf("  ReadString(' '): %q\n", line)

	// Read - خواندن به slice
	remaining := make([]byte, 10)
	n, _ := buf3.Read(remaining)
	fmt.Printf("  Read remaining: %q (%d bytes)\n", remaining[:n], n)

	// 8.4 Reset - پاک کردن Buffer
	fmt.Println("\n--- 8.4 Buffer Reset ---")

	buf4 := bytes.NewBufferString("Data to reset")
	fmt.Printf("  Before Reset: %q\n", buf4.String())

	buf4.Reset()
	fmt.Printf("  After Reset: %q (empty)\n", buf4.String())
	fmt.Printf("  Note: Reset doesn't free memory, just sets length to 0\n")

	// 8.5 Next - دریافت n بایت بعدی بدون حذف
	fmt.Println("\n--- 8.5 Buffer Next ---")

	buf5 := bytes.NewBufferString("Hello World")
	next2 := buf5.Next(5)
	fmt.Printf("  Next(5): %q\n", next2)
	fmt.Printf("  Remaining: %q\n", buf5.String())

	// 8.6 Bytes - دریافت slice بایت (بدون کپی)
	fmt.Println("\n--- 8.6 Buffer Bytes ---")

	buf6 := bytes.NewBufferString("data")
	bytes6 := buf6.Bytes()
	fmt.Printf("  Bytes(): %q (direct reference)\n", bytes6)

	// 8.7 String - دریافت رشته (بدون کپی)
	fmt.Println("\n--- 8.7 Buffer String ---")

	buf7 := bytes.NewBufferString("string data")
	str7 := buf7.String()
	fmt.Printf("  String(): %q\n", str7)

	// 8.8 WriteTo - نوشتن محتوا به Writer
	fmt.Println("\n--- 8.8 Buffer WriteTo ---")

	buf8 := bytes.NewBufferString("Content to write")
	var dest bytes.Buffer

	written, _ := buf8.WriteTo(&dest)
	fmt.Printf("  WriteTo wrote %d bytes: %q\n", written, dest.String())

	// 8.9 ReadFrom - خواندن از Reader
	fmt.Println("\n--- 8.9 Buffer ReadFrom ---")

	source := bytes.NewBufferString("Source content")
	var dest2 bytes.Buffer

	read, _ := dest2.ReadFrom(source)
	fmt.Printf("  ReadFrom read %d bytes: %q\n", read, dest2.String())

	// ============================================================================
	// بخش 9: توابع کاربردی دیگر
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🔧 SECTION 9: OTHER USEFUL FUNCTIONS")
	fmt.Println(strings.Repeat("=", 80))

	// 9.1 Repeat - تکرار slice بایت
	fmt.Println("\n--- 9.1 bytes.Repeat ---")
	// Repeat returns a new byte slice consisting of count copies of b.
	star := []byte("*")
	repeated := bytes.Repeat(star, 10)
	fmt.Printf("  Repeat(%q, 10): %q\n", star, repeated)

	// 9.2 Runes - تبدیل به slice رون
	fmt.Println("\n--- 9.2 bytes.Runes ---")
	// Runes interprets s as a UTF-8-encoded sequence and returns a slice of runes.
	utf8 := []byte("Hello 世界")
	runes := bytes.Runes(utf8)
	fmt.Printf("  Runes(%q): %U\n", utf8, runes)

	// 9.3 NewBuffer - ایجاد Buffer از slice بایت
	fmt.Println("\n--- 9.3 bytes.NewBuffer ---")
	// NewBuffer creates and initializes a new Buffer using buf as its initial contents.
	buf9 := bytes.NewBuffer([]byte{72, 101, 108, 108, 111})
	fmt.Printf("  NewBuffer from bytes: %q\n", buf9.String())

	// 9.4 NewBufferString - ایجاد Buffer از رشته
	fmt.Println("\n--- 9.4 bytes.NewBufferString ---")
	// NewBufferString creates and initializes a new Buffer using string s as its initial contents.
	buf10 := bytes.NewBufferString("Hello from string")
	fmt.Printf("  NewBufferString: %q\n", buf10.String())

	// 9.5 Cut - برش slice بایت به دو قسمت (Go 1.18+)
	fmt.Println("\n--- 9.5 bytes.Cut ---")
	// Cut slices s around the first instance of sep, returning the text before and after sep.
	cutData := []byte("hello-world")
	before, after, found := bytes.Cut(cutData, []byte("-"))
	fmt.Printf("  Cut(%q, '-'): before=%q, after=%q, found=%v\n", cutData, before, after, found)

	// 9.6 Clone - کپی slice بایت
	fmt.Println("\n--- 9.6 bytes.Clone ---")
	// Clone returns a copy of b. The copy is always distinct from b.
	originalBytes := []byte("original")
	clonedBytes := bytes.Clone(originalBytes)
	clonedBytes[0] = 'X'
	fmt.Printf("  Original: %q, Cloned: %q (independent)\n", originalBytes, clonedBytes)

	// 9.7 Map - تبدیل هر بایت با تابع
	fmt.Println("\n--- 9.7 bytes.Map ---")
	// Map returns a copy of the byte slice s with all its characters modified according to the mapping function.
	toUpperMap := func(r rune) rune {
		if r >= 'a' && r <= 'z' {
			return r - 32
		}
		return r
	}
	mapped := bytes.Map(toUpperMap, []byte("hello world"))
	fmt.Printf("  Map(toUpper, %q): %q\n", "hello world", mapped)

	// ============================================================================
	// بخش 10: کاربردهای عملی (Practical Examples)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("💡 SECTION 10: PRACTICAL EXAMPLES")
	fmt.Println(strings.Repeat("=", 80))

	// 10.1 ساخت کارآمد رشته با Buffer
	fmt.Println("\n--- 10.1 Efficient String Building ---")

	var builder bytes.Buffer
	for i := 0; i < 10; i++ {
		builder.WriteString(fmt.Sprintf("Item %d, ", i))
	}
	result2 := builder.String()
	fmt.Printf("  Built string: %s\n", result2)

	// 10.2 خواندن خط به خط از بافر
	fmt.Println("\n--- 10.2 Line by Line Reading ---")

	multiLine := []byte("Line 1\nLine 2\nLine 3\nLine 4")
	reader := bytes.NewReader(multiLine)

	lineNum := 1
	for {
		line, err := reader.ReadBytes('\n')
		if err == io.EOF {
			if len(line) > 0 {
				fmt.Printf("  Line %d: %q", lineNum, line)
			}
			break
		}
		if err != nil {
			break
		}
		fmt.Printf("  Line %d: %q", lineNum, line)
		lineNum++
	}
	fmt.Println()

	// 10.3 تشخیص نوع فایل با بررسی magic bytes
	fmt.Println("\n--- 10.3 File Type Detection (Magic Bytes) ---")

	detectFileType := func(data []byte) string {
		if len(data) < 4 {
			return "unknown"
		}
		if bytes.Equal(data[:4], []byte{0x89, 0x50, 0x4E, 0x47}) {
			return "PNG"
		}
		if bytes.Equal(data[:4], []byte{0xFF, 0xD8, 0xFF, 0xE0}) ||
			bytes.Equal(data[:4], []byte{0xFF, 0xD8, 0xFF, 0xE1}) {
			return "JPEG"
		}
		if bytes.Equal(data[:4], []byte{0x25, 0x50, 0x44, 0x46}) {
			return "PDF"
		}
		if bytes.HasPrefix(data, []byte("GIF")) {
			return "GIF"
		}
		return "unknown"
	}

	pngHeader := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	fmt.Printf("  PNG header: %s\n", detectFileType(pngHeader))

	pdfHeader := []byte("%PDF-1.4")
	fmt.Printf("  PDF header: %s\n", detectFileType(pdfHeader))

	// 10.4 کپی کردن بین بافرها
	fmt.Println("\n--- 10.4 Copy Between Buffers ---")

	srcBuffer := bytes.NewBufferString("Source data")
	dstBuffer := bytes.NewBuffer(nil)

	written, _ = io.Copy(dstBuffer, srcBuffer)
	fmt.Printf("  Copied %d bytes: %q\n", written, dstBuffer.String())

	// 10.5 جستجو و جایگزینی در فایل (شبیه‌سازی)
	fmt.Println("\n--- 10.5 Search and Replace in Buffer ---")

	content := []byte("The quick brown fox jumps over the lazy dog")
	oldWord := []byte("fox")
	newWord := []byte("cat")

	replacedContent := bytes.ReplaceAll(content, oldWord, newWord)
	fmt.Printf("  Original: %q\n", content)
	fmt.Printf("  Replaced: %q\n", replacedContent)

	// ============================================================================
	// بخش 11: اشتباهات رایج (Common Mistakes)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("❌ SECTION 11: COMMON MISTAKES")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println("\n❌ Mistake 1: Modifying bytes returned by Buffer.Bytes()")
	fmt.Println("   b := buf.Bytes()")
	fmt.Println("   b[0] = 'X'  // Modifies buffer internals!")
	fmt.Println("   ✅ Use copy() for safe modification")

	fmt.Println("\n❌ Mistake 2: Not resetting Buffer between uses")
	fmt.Println("   var buf bytes.Buffer")
	fmt.Println("   buf.WriteString(\"data\")")
	fmt.Println("   // buf not reset before next use")
	fmt.Println("   ✅ Use buf.Reset() or create new buffer")

	fmt.Println("\n❌ Mistake 3: Assuming bytes.Equal works like string equality")
	fmt.Println("   bytes.Equal(nil, []byte{})  // returns true")
	fmt.Println("   ✅ Check nil explicitly if needed")

	fmt.Println("\n❌ Mistake 4: Using bytes.Buffer without Grow for large data")
	fmt.Println("   var buf bytes.Buffer")
	fmt.Println("   for i := 0; i < 1000000; i++ { buf.WriteString(\"x\") }")
	fmt.Println("   ✅ buf.Grow(1000000) for pre-allocation")

	fmt.Println("\n❌ Mistake 5: Ignoring return values of Read functions")
	fmt.Println("   n, _ := buf.Read(data)  // ignoring n")
	fmt.Println("   ✅ Always use the returned n")

	fmt.Println("\n❌ Mistake 6: Using bytes.Split with empty separator")
	fmt.Println("   bytes.Split(data, []byte{})  // splits every byte!")
	fmt.Println("   ✅ Check for empty separator before splitting")

	fmt.Println("\n❌ Mistake 7: Forgetting that Replace returns new slice")
	fmt.Println("   bytes.Replace(data, old, new, -1)  // data unchanged")
	fmt.Println("   ✅ data = bytes.Replace(data, old, new, -1)")

	// ============================================================================
	// بخش 12: جمع‌بندی و جدول مرجع سریع
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📚 SECTION 12: QUICK REFERENCE")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ CATEGORY          │ FUNCTIONS                                    │")
	fmt.Println("├───────────────────┼──────────────────────────────────────────────┤")
	fmt.Println("│ Comparison        │ Compare, Equal, EqualFold                   │")
	fmt.Println("│ Search            │ Contains, ContainsAny, ContainsRune, Count,  │")
	fmt.Println("│                   │ Index, IndexAny, IndexByte, IndexFunc,       │")
	fmt.Println("│                   │ IndexRune, LastIndex, LastIndexAny,          │")
	fmt.Println("│                   │ LastIndexByte, LastIndexFunc                 │")
	fmt.Println("│ Replacement       │ Replace, ReplaceAll, NewReplacer             │")
	fmt.Println("│ Split/Join        │ Split, SplitN, SplitAfter, SplitAfterN,      │")
	fmt.Println("│                   │ Fields, FieldsFunc, Join                     │")
	fmt.Println("│ Trim              │ Trim, TrimSpace, TrimLeft, TrimRight,        │")
	fmt.Println("│                   │ TrimPrefix, TrimSuffix, TrimFunc,            │")
	fmt.Println("│                   │ TrimLeftFunc, TrimRightFunc                  │")
	fmt.Println("│ Case Conversion   │ ToLower, ToUpper, ToTitle                    │")
	fmt.Println("│ Prefix/Suffix     │ HasPrefix, HasSuffix                         │")
	fmt.Println("│ Buffer            │ NewBuffer, NewBufferString, Write, WriteString│")
	fmt.Println("│                   │ WriteByte, WriteRune, Read, ReadByte,        │")
	fmt.Println("│                   │ ReadString, ReadBytes, Next, Bytes, String,  │")
	fmt.Println("│                   │ Len, Cap, Grow, Reset, Truncate, WriteTo,    │")
	fmt.Println("│                   │ ReadFrom, Next                               │")
	fmt.Println("│ Other             │ Repeat, Runes, Cut, Clone, Map               │")
	fmt.Println("└───────────────────┴──────────────────────────────────────────────┘")

	fmt.Println("\n💡 GOLDEN RULES:")
	fmt.Println("  1. Use bytes.Buffer for efficient string/buffer concatenation")
	fmt.Println("  2. Pre-allocate buffer capacity with Grow() when size is known")
	fmt.Println("  3. Use bytes.Equal for comparing byte slices (not bytes.Compare)")
	fmt.Println("  4. Always check returned n from Read operations")
	fmt.Println("  5. Be careful modifying bytes returned by Buffer.Bytes()")
	fmt.Println("  6. Use bytes.ReplaceAll for simple replacements")
	fmt.Println("  7. Use bytes.Fields for splitting on whitespace")
	fmt.Println("  8. Use bytes.TrimSpace for removing all whitespace")
	fmt.Println("  9. bytes.Clone for independent copy of byte slice")
	fmt.Println("  10. Remember: bytes are not strings - use conversion when needed")

	fmt.Println("\n🎯 BUFFER vs STRING BUILDER:")
	fmt.Println("  • bytes.Buffer  → Works with bytes ([]byte)")
	fmt.Println("  • strings.Builder → Works with strings (string)")
	fmt.Println("  • Choose based on what you need to build")
}

// تابع کمکی برای تکرار رشته
func stringsRepeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}
