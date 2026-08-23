// ============================================================================
// FILE: strings_complete_guide.go
// TITLE: راهنمای کامل پکیج strings در Go - تمام توابع با مثال
// HOW TO RUN: go run strings_complete_guide.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - پکیج strings چیست؟
// ============================================================================
//
// پکیج strings توابع متعددی برای کار با رشته‌ها (string) ارائه می‌دهد:
// - مقایسه، جستجو، جایگزینی
// - تقسیم، اتصال، برش
// - تبدیل حروف بزرگ/کوچک
// - حذف فاصله‌ها و کاراکترهای خاص
// - بررسی پیشوند/پسوند
// - تکرار، padding
// - و بسیاری دیگر...
//
// قانون طلایی:
// "هر زمان که با رشته کار می‌کنی، اول ببین آیا تابعی در پکیج strings هست
//  که کارت را انجام دهد. این پکیج بسیار کامل است و نیازی به re-invent کردن wheel نیست."
// ============================================================================

package __internal_packages

import (
	"fmt"
	"strings"
	"unicode"
)

func main() {
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("🎯 COMPLETE strings PACKAGE GUIDE IN GO")
	fmt.Println("All functions with practical examples")
	fmt.Println(strings.Repeat("=", 80))

	// ============================================================================
	// بخش 1: مقایسه رشته‌ها (Comparison)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📊 SECTION 1: COMPARISON FUNCTIONS")
	fmt.Println(strings.Repeat("=", 80))

	// 1.1 Compare - مقایسه دو رشته (حساس به حروف بزرگ/کوچک)
	fmt.Println("\n--- 1.1 strings.Compare(a, b string) int ---")
	fmt.Println("  returns -1 if a < b, 0 if a == b, 1 if a > b")
	fmt.Printf("  Compare(\"a\", \"b\"): %d\n", strings.Compare("a", "b"))
	fmt.Printf("  Compare(\"b\", \"a\"): %d\n", strings.Compare("b", "a"))
	fmt.Printf("  Compare(\"a\", \"a\"): %d\n", strings.Compare("a", "a"))
	fmt.Printf("  Compare(\"abc\", \"abd\"): %d\n", strings.Compare("abc", "abd"))

	// 1.2 EqualFold - مقایسه بدون در نظر گرفتن حروف بزرگ/کوچک
	fmt.Println("\n--- 1.2 strings.EqualFold(s, t string) bool ---")
	fmt.Printf("  EqualFold(\"Go\", \"go\"): %v\n", strings.EqualFold("Go", "go"))
	fmt.Printf("  EqualFold(\"Hello\", \"HELLO\"): %v\n", strings.EqualFold("Hello", "HELLO"))
	fmt.Printf("  EqualFold(\"abc\", \"def\"): %v\n", strings.EqualFold("abc", "def"))

	// ============================================================================
	// بخش 2: جستجو در رشته (Search)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🔍 SECTION 2: SEARCH FUNCTIONS")
	fmt.Println(strings.Repeat("=", 80))

	sample := "The quick brown fox jumps over the lazy dog"

	// 2.1 Contains - آیا رشته شامل زیررشته است؟
	fmt.Println("\n--- 2.1 strings.Contains(s, substr string) bool ---")
	fmt.Printf("  Contains(\"%s\", \"fox\"): %v\n", sample, strings.Contains(sample, "fox"))
	fmt.Printf("  Contains(\"%s\", \"cat\"): %v\n", sample, strings.Contains(sample, "cat"))

	// 2.2 ContainsAny - آیا شامل هر یک از کاراکترهاست؟
	fmt.Println("\n--- 2.2 strings.ContainsAny(s, chars string) bool ---")
	fmt.Printf("  ContainsAny(\"%s\", \"xyz\"): %v\n", sample, strings.ContainsAny(sample, "xyz"))
	fmt.Printf("  ContainsAny(\"%s\", \"123\"): %v\n", sample, strings.ContainsAny(sample, "123"))

	// 2.3 ContainsRune - آیا شامل رون خاص است؟
	fmt.Println("\n--- 2.3 strings.ContainsRune(s string, r rune) bool ---")
	fmt.Printf("  ContainsRune(\"Hello\", 'e'): %v\n", strings.ContainsRune("Hello", 'e'))
	fmt.Printf("  ContainsRune(\"Hello\", 'x'): %v\n", strings.ContainsRune("Hello", 'x'))

	// 2.4 Count - تعداد دفعات تکرار زیررشته
	fmt.Println("\n--- 2.4 strings.Count(s, substr string) int ---")
	fmt.Printf("  Count(\"banana\", \"na\"): %d\n", strings.Count("banana", "na"))
	fmt.Printf("  Count(\"aaaaaa\", \"aa\"): %d\n", strings.Count("aaaaaa", "aa"))
	fmt.Printf("  Count(\"hello\", \"x\"): %d\n", strings.Count("hello", "x"))

	// 2.5 Index - اولین موقعیت زیررشته
	fmt.Println("\n--- 2.5 strings.Index(s, substr string) int ---")
	fmt.Printf("  Index(\"hello world\", \"world\"): %d\n", strings.Index("hello world", "world"))
	fmt.Printf("  Index(\"hello world\", \"xyz\"): %d\n", strings.Index("hello world", "xyz"))

	// 2.6 IndexAny - اولین موقعیت هر یک از کاراکترها
	fmt.Println("\n--- 2.6 strings.IndexAny(s, chars string) int ---")
	fmt.Printf("  IndexAny(\"hello\", \"aeiou\"): %d (اولین مصوت)\n", strings.IndexAny("hello", "aeiou"))
	fmt.Printf("  IndexAny(\"hello\", \"xyz\"): %d\n", strings.IndexAny("hello", "xyz"))

	// 2.7 IndexByte - موقعیت یک بایت خاص
	fmt.Println("\n--- 2.7 strings.IndexByte(s string, c byte) int ---")
	fmt.Printf("  IndexByte(\"hello\", 'e'): %d\n", strings.IndexByte("hello", 'e'))
	fmt.Printf("  IndexByte(\"hello\", 'z'): %d\n", strings.IndexByte("hello", 'z'))

	// 2.8 IndexFunc - موقعیت کاراکتری که شرط را برآورده می‌کند
	fmt.Println("\n--- 2.8 strings.IndexFunc(s string, f func(rune) bool) int ---")
	isDigit := func(r rune) bool { return r >= '0' && r <= '9' }
	isUpper := func(r rune) bool { return r >= 'A' && r <= 'Z' }
	fmt.Printf("  IndexFunc(\"abc123\", isDigit): %d (اولین رقم)\n", strings.IndexFunc("abc123", isDigit))
	fmt.Printf("  IndexFunc(\"hello\", isUpper): %d\n", strings.IndexFunc("hello", isUpper))

	// 2.9 IndexRune - موقعیت یک رون خاص
	fmt.Println("\n--- 2.9 strings.IndexRune(s string, r rune) int ---")
	fmt.Printf("  IndexRune(\"Hello\", 'e'): %d\n", strings.IndexRune("Hello", 'e'))
	fmt.Printf("  IndexRune(\"Hello\", '世'): %d (فارسی/چینی)\n", strings.IndexRune("Hello世", '世'))

	// 2.10 LastIndex - آخرین موقعیت زیررشته
	fmt.Println("\n--- 2.10 strings.LastIndex(s, substr string) int ---")
	fmt.Printf("  LastIndex(\"banana\", \"na\"): %d\n", strings.LastIndex("banana", "na"))
	fmt.Printf("  LastIndex(\"hello hello\", \"hello\"): %d\n", strings.LastIndex("hello hello", "hello"))

	// 2.11 LastIndexAny - آخرین موقعیت هر یک از کاراکترها
	fmt.Println("\n--- 2.11 strings.LastIndexAny(s, chars string) int ---")
	fmt.Printf("  LastIndexAny(\"hello world\", \"aeiou\"): %d\n", strings.LastIndexAny("hello world", "aeiou"))

	// 2.12 LastIndexByte - آخرین موقعیت یک بایت خاص
	fmt.Println("\n--- 2.12 strings.LastIndexByte(s string, c byte) int ---")
	fmt.Printf("  LastIndexByte(\"banana\", 'a'): %d\n", strings.LastIndexByte("banana", 'a'))

	// 2.13 LastIndexFunc - آخرین موقعیت کاراکتری که شرط را برآورده می‌کند
	fmt.Println("\n--- 2.13 strings.LastIndexFunc(s string, f func(rune) bool) int ---")
	fmt.Printf("  LastIndexFunc(\"abc123def456\", isDigit): %d (آخرین رقم)\n", strings.LastIndexFunc("abc123def456", isDigit))

	// ============================================================================
	// بخش 3: جایگزینی (Replacement)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🔄 SECTION 3: REPLACEMENT FUNCTIONS")
	fmt.Println(strings.Repeat("=", 80))

	// 3.1 Replace - جایگزینی با تعداد مشخص
	fmt.Println("\n--- 3.1 strings.Replace(s, old, new string, n int) string ---")
	fmt.Printf("  Replace(\"oink oink oink\", \"oink\", \"moo\", 2): %s\n",
		strings.Replace("oink oink oink", "oink", "moo", 2))
	fmt.Printf("  Replace(\"oink oink oink\", \"oink\", \"moo\", -1): %s (همه)\n",
		strings.Replace("oink oink oink", "oink", "moo", -1))

	// 3.2 ReplaceAll - جایگزینی همه موارد
	fmt.Println("\n--- 3.2 strings.ReplaceAll(s, old, new string) string ---")
	fmt.Printf("  ReplaceAll(\"hello world hello\", \"hello\", \"hi\"): %s\n",
		strings.ReplaceAll("hello world hello", "hello", "hi"))

	// 3.3 NewReplacer - جایگزینی چندتایی با Replacer
	fmt.Println("\n--- 3.3 strings.NewReplacer(oldnew ...string) *strings.Replacer ---")
	replacer := strings.NewReplacer(
		"hello", "hi",
		"world", "earth",
		"good", "bad",
	)
	fmt.Printf("  Replacer: %s\n", replacer.Replace("hello world good morning"))

	// ============================================================================
	// بخش 4: تقسیم و اتصال (Split & Join)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("✂️ SECTION 4: SPLIT & JOIN FUNCTIONS")
	fmt.Println(strings.Repeat("=", 80))

	// 4.1 Split - تقسیم با جداکننده
	fmt.Println("\n--- 4.1 strings.Split(s, sep string) []string ---")
	fmt.Printf("  Split(\"a,b,c\", \",\"): %q\n", strings.Split("a,b,c", ","))
	fmt.Printf("  Split(\"a,,c\", \",\"): %q (توجه: رشته خالی)\n", strings.Split("a,,c", ","))

	// 4.2 SplitN - تقسیم با حداکثر n قسمت
	fmt.Println("\n--- 4.2 strings.SplitN(s, sep string, n int) []string ---")
	fmt.Printf("  SplitN(\"a,b,c,d\", \",\", 2): %q\n", strings.SplitN("a,b,c,d", ",", 2))
	fmt.Printf("  SplitN(\"a,b,c,d\", \",\", 3): %q\n", strings.SplitN("a,b,c,d", ",", 3))

	// 4.3 SplitAfter - تقسیم با حفظ جداکننده
	fmt.Println("\n--- 4.3 strings.SplitAfter(s, sep string) []string ---")
	fmt.Printf("  SplitAfter(\"a,b,c\", \",\"): %q\n", strings.SplitAfter("a,b,c", ","))

	// 4.4 SplitAfterN - تقسیم با حداکثر n قسمت (حفظ جداکننده)
	fmt.Println("\n--- 4.4 strings.SplitAfterN(s, sep string, n int) []string ---")
	fmt.Printf("  SplitAfterN(\"a,b,c,d\", \",\", 2): %q\n", strings.SplitAfterN("a,b,c,d", ",", 2))

	// 4.5 Fields - تقسیم بر اساس whitespace
	fmt.Println("\n--- 4.5 strings.Fields(s string) []string ---")
	fmt.Printf("  Fields(\"hello world   from  go\"): %q\n",
		strings.Fields("hello world   from  go"))

	// 4.6 FieldsFunc - تقسیم با تابع شرط
	fmt.Println("\n--- 4.6 strings.FieldsFunc(s string, f func(rune) bool) []string ---")
	isNotLetter := func(r rune) bool { return (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') }
	fmt.Printf("  FieldsFunc(\"hello123world456go\", isNotLetter): %q\n",
		strings.FieldsFunc("hello123world456go", isNotLetter))

	// 4.7 Join - اتصال اسلایس رشته با جداکننده
	fmt.Println("\n--- 4.7 strings.Join(elems []string, sep string) string ---")
	words := []string{"hello", "world", "from", "go"}
	fmt.Printf("  Join(%q, \" \"): %s\n", words, strings.Join(words, " "))
	fmt.Printf("  Join(%q, \", \"): %s\n", words, strings.Join(words, ", "))

	// ============================================================================
	// بخش 5: حذف فاصله و کاراکترهای خاص (Trim)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("✂️ SECTION 5: TRIM FUNCTIONS")
	fmt.Println(strings.Repeat("=", 80))

	// 5.1 Trim - حذف کاراکترهای مشخص از دو طرف
	fmt.Println("\n--- 5.1 strings.Trim(s, cutset string) string ---")
	fmt.Printf("  Trim(\"!!!Hello!!!\", \"!\"): %q\n", strings.Trim("!!!Hello!!!", "!"))
	fmt.Printf("  Trim(\"   Hello   \", \" \"): %q\n", strings.Trim("   Hello   ", " "))

	// 5.2 TrimSpace - حذف whitespace از دو طرف
	fmt.Println("\n--- 5.2 strings.TrimSpace(s string) string ---")
	fmt.Printf("  TrimSpace(\"  \\t\\nHello World\\n\\t  \"): %q\n",
		strings.TrimSpace("  \t\nHello World\n\t  "))

	// 5.3 TrimLeft - حذف کاراکترها از سمت چپ
	fmt.Println("\n--- 5.3 strings.TrimLeft(s, cutset string) string ---")
	fmt.Printf("  TrimLeft(\"!!!Hello!!!\", \"!\"): %q\n", strings.TrimLeft("!!!Hello!!!", "!"))

	// 5.4 TrimRight - حذف کاراکترها از سمت راست
	fmt.Println("\n--- 5.4 strings.TrimRight(s, cutset string) string ---")
	fmt.Printf("  TrimRight(\"!!!Hello!!!\", \"!\"): %q\n", strings.TrimRight("!!!Hello!!!", "!"))

	// 5.5 TrimPrefix - حذف پیشوند
	fmt.Println("\n--- 5.5 strings.TrimPrefix(s, prefix string) string ---")
	fmt.Printf("  TrimPrefix(\"HelloWorld\", \"Hello\"): %q\n",
		strings.TrimPrefix("HelloWorld", "Hello"))
	fmt.Printf("  TrimPrefix(\"HelloWorld\", \"World\"): %q (تغییری نمی‌کند)\n",
		strings.TrimPrefix("HelloWorld", "World"))

	// 5.6 TrimSuffix - حذف پسوند
	fmt.Println("\n--- 5.6 strings.TrimSuffix(s, suffix string) string ---")
	fmt.Printf("  TrimSuffix(\"HelloWorld\", \"World\"): %q\n",
		strings.TrimSuffix("HelloWorld", "World"))

	// 5.7 TrimFunc - حذف کاراکترهایی که شرط را برآورده می‌کنند
	fmt.Println("\n--- 5.7 strings.TrimFunc(s string, f func(rune) bool) string ---")
	isPunct := func(r rune) bool { return r == '!' || r == '?' || r == '.' }
	fmt.Printf("  TrimFunc(\"!!!Hello!!!\", isPunct): %q\n",
		strings.TrimFunc("!!!Hello!!!", isPunct))

	// 5.8 TrimLeftFunc - حذف از چپ با شرط
	fmt.Println("\n--- 5.8 strings.TrimLeftFunc(s string, f func(rune) bool) string ---")
	fmt.Printf("  TrimLeftFunc(\"!!!Hello!!!\", isPunct): %q\n",
		strings.TrimLeftFunc("!!!Hello!!!", isPunct))

	// 5.9 TrimRightFunc - حذف از راست با شرط
	fmt.Println("\n--- 5.9 strings.TrimRightFunc(s string, f func(rune) bool) string ---")
	fmt.Printf("  TrimRightFunc(\"!!!Hello!!!\", isPunct): %q\n",
		strings.TrimRightFunc("!!!Hello!!!", isPunct))

	// ============================================================================
	// بخش 6: تبدیل حروف (Case Conversion)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🔠 SECTION 6: CASE CONVERSION FUNCTIONS")
	fmt.Println(strings.Repeat("=", 80))

	// 6.1 ToLower - تبدیل به حروف کوچک
	fmt.Println("\n--- 6.1 strings.ToLower(s string) string ---")
	fmt.Printf("  ToLower(\"HELLO World\"): %s\n", strings.ToLower("HELLO World"))
	fmt.Printf("  ToLower(\"Αθήνα\"): %s (یونانی)\n", strings.ToLower("Αθήνα"))

	// 6.2 ToUpper - تبدیل به حروف بزرگ
	fmt.Println("\n--- 6.2 strings.ToUpper(s string) string ---")
	fmt.Printf("  ToUpper(\"Hello World\"): %s\n", strings.ToUpper("Hello World"))

	// 6.3 ToTitle - تبدیل به حروف عنوان (حرف اول هر کلمه بزرگ)
	fmt.Println("\n--- 6.3 strings.ToTitle(s string) string ---")
	fmt.Printf("  ToTitle(\"hello world\"): %s\n", strings.ToTitle("hello world"))
	fmt.Printf("  ToTitle(\"hElLo wOrLd\"): %s\n", strings.ToTitle("hElLo wOrLd"))

	// 6.4 ToTitleSpecial - تبدیل با قوانین خاص (مثل ترکی)
	fmt.Println("\n--- 6.4 strings.ToTitleSpecial(c unicode.SpecialCase, s string) string ---")
	fmt.Printf("  ToTitleSpecial(unicode.TurkishCase, \"i\"): %s (ترکی)\n",
		strings.ToTitleSpecial(unicode.TurkishCase, "i"))

	// 6.5 ToLowerSpecial - تبدیل با قوانین خاص
	fmt.Println("\n--- 6.5 strings.ToLowerSpecial(c unicode.SpecialCase, s string) string ---")
	fmt.Printf("  ToLowerSpecial(unicode.TurkishCase, \"İ\"): %s\n",
		strings.ToLowerSpecial(unicode.TurkishCase, "İ"))

	// 6.6 ToUpperSpecial - تبدیل با قوانین خاص
	fmt.Println("\n--- 6.6 strings.ToUpperSpecial(c unicode.SpecialCase, s string) string ---")
	fmt.Printf("  ToUpperSpecial(unicode.TurkishCase, \"i\"): %s\n",
		strings.ToUpperSpecial(unicode.TurkishCase, "i"))

	// ============================================================================
	// بخش 7: بررسی پیشوند و پسوند (Prefix & Suffix)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📌 SECTION 7: PREFIX & SUFFIX FUNCTIONS")
	fmt.Println(strings.Repeat("=", 80))

	// 7.1 HasPrefix - آیا با پیشوند مشخص شروع می‌شود؟
	fmt.Println("\n--- 7.1 strings.HasPrefix(s, prefix string) bool ---")
	fmt.Printf("  HasPrefix(\"HelloWorld\", \"Hello\"): %v\n", strings.HasPrefix("HelloWorld", "Hello"))
	fmt.Printf("  HasPrefix(\"HelloWorld\", \"World\"): %v\n", strings.HasPrefix("HelloWorld", "World"))

	// 7.2 HasSuffix - آیا با پسوند مشخص تمام می‌شود؟
	fmt.Println("\n--- 7.2 strings.HasSuffix(s, suffix string) bool ---")
	fmt.Printf("  HasSuffix(\"HelloWorld\", \"World\"): %v\n", strings.HasSuffix("HelloWorld", "World"))
	fmt.Printf("  HasSuffix(\"HelloWorld\", \"Hello\"): %v\n", strings.HasSuffix("HelloWorld", "Hello"))

	// ============================================================================
	// بخش 8: تکرار و ساخت رشته (Repeat & Builder)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🏗️ SECTION 8: REPEAT & BUILDER FUNCTIONS")
	fmt.Println(strings.Repeat("=", 80))

	// 8.1 Repeat - تکرار رشته
	fmt.Println("\n--- 8.1 strings.Repeat(s string, count int) string ---")
	fmt.Printf("  Repeat(\"Go\", 3): %s\n", strings.Repeat("Go", 3))
	fmt.Printf("  Repeat(\"-\", 10): %s\n", strings.Repeat("-", 10))

	// 8.2 Builder - ساخت کارآمد رشته (بدون allocation زیاد)
	fmt.Println("\n--- 8.2 strings.Builder (کارآمدترین روش ساخت رشته) ---")
	var builder strings.Builder

	// پیش‌تخصیص حافظه برای بهبود کارایی
	builder.Grow(100)

	// افزودن رشته‌ها
	builder.WriteString("Hello")
	builder.WriteString(" ")
	builder.WriteString("World")
	builder.WriteString(" ")
	builder.WriteString("from")
	builder.WriteString(" ")
	builder.WriteString("Go")

	// افزودن بایت
	builder.WriteByte('!')

	// افزودن رون
	builder.WriteRune('🎉')

	fmt.Printf("  Builder result: %s\n", builder.String())
	fmt.Printf("  Builder length: %d, capacity: %d\n", builder.Len(), builder.Cap())

	// Reset کردن builder
	builder.Reset()
	builder.WriteString("After reset")
	fmt.Printf("  After reset: %s\n", builder.String())

	// ============================================================================
	// بخش 9: توابع کمکی دیگر (Other Helpers)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🔧 SECTION 9: OTHER HELPER FUNCTIONS")
	fmt.Println(strings.Repeat("=", 80))

	// 9.1 Clone - کپی کردن رشته
	fmt.Println("\n--- 9.1 strings.Clone(s string) string ---")
	original := "hello"
	cloned := strings.Clone(original)
	fmt.Printf("  Clone(\"%s\"): %s (آدرس متفاوت)\n", original, cloned)

	// 9.2 Cut - برش رشته به دو قسمت بر اساس جداکننده (Go 1.18+)
	fmt.Println("\n--- 9.2 strings.Cut(s, sep string) (before, after string, found bool) ---")
	before, after, found := strings.Cut("hello-world", "-")
	fmt.Printf("  Cut(\"hello-world\", \"-\"): before=%q, after=%q, found=%v\n",
		before, after, found)

	before2, after2, found2 := strings.Cut("hello-world", "+")
	fmt.Printf("  Cut(\"hello-world\", \"+\"): before=%q, after=%q, found=%v\n",
		before2, after2, found2)

	// 9.3 CutPrefix - حذف پیشوند اگر وجود داشته باشد (Go 1.20+)
	fmt.Println("\n--- 9.3 strings.CutPrefix(s, prefix string) (after string, found bool) ---")
	after3, found3 := strings.CutPrefix("HelloWorld", "Hello")
	fmt.Printf("  CutPrefix(\"HelloWorld\", \"Hello\"): after=%q, found=%v\n", after3, found3)

	// 9.4 CutSuffix - حذف پسوند اگر وجود داشته باشد (Go 1.20+)
	fmt.Println("\n--- 9.4 strings.CutSuffix(s, suffix string) (before string, found bool) ---")
	before4, found4 := strings.CutSuffix("HelloWorld", "World")
	fmt.Printf("  CutSuffix(\"HelloWorld\", \"World\"): before=%q, found=%v\n", before4, found4)

	// 9.5 Map - تبدیل هر کاراکتر با تابع
	fmt.Println("\n--- 9.5 strings.Map(mapping func(rune) rune, s string) string ---")
	toUpper := func(r rune) rune {
		if r >= 'a' && r <= 'z' {
			return r - 32
		}
		return r
	}
	fmt.Printf("  Map(toUpper, \"hello\"): %s\n", strings.Map(toUpper, "hello"))

	// 9.6 NewReader - ایجاد Reader از رشته
	fmt.Println("\n--- 9.6 strings.NewReader(s string) *strings.Reader ---")
	reader := strings.NewReader("Hello from reader")
	buf := make([]byte, 10)
	n, _ := reader.Read(buf)
	fmt.Printf("  NewReader read %d bytes: %s\n", n, buf[:n])

	// ============================================================================
	// بخش 10: مثال‌های کاربردی (Practical Examples)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("💡 SECTION 10: PRACTICAL EXAMPLES")
	fmt.Println(strings.Repeat("=", 80))

	// 10.1 اعتبارسنجی ایمیل (ساده)
	fmt.Println("\n--- 10.1 Email Validation ---")
	email := "user@example.com"
	if strings.Contains(email, "@") && strings.Contains(email, ".") {
		fmt.Printf("  %s is valid email\n", email)
	}

	// 10.2 حذف فاصله‌های اضافی
	fmt.Println("\n--- 10.2 Remove Extra Spaces ---")
	messy := "  hello   world  from   go  "
	cleaned := strings.Join(strings.Fields(messy), " ")
	fmt.Printf("  Original: %q\n", messy)
	fmt.Printf("  Cleaned: %q\n", cleaned)

	// 10.3 بررسی پسوند فایل
	fmt.Println("\n--- 10.3 File Extension Check ---")
	filename := "document.pdf"
	if strings.HasSuffix(filename, ".pdf") {
		fmt.Printf("  %s is a PDF file\n", filename)
	}

	// 10.4 ساخت رشته با جداکننده
	fmt.Println("\n--- 10.4 Build CSV Line ---")
	fields := []string{"Ali", "30", "Tehran"}
	csvLine := strings.Join(fields, ",")
	fmt.Printf("  CSV line: %s\n", csvLine)

	// 10.5 تشخیص کلمات تکراری
	fmt.Println("\n--- 10.5 Count Word Frequency ---")
	text := "go is great and go is powerful"
	word := "go"
	count := strings.Count(text, word)
	fmt.Printf("  '%s' appears %d times in '%s'\n", word, count, text)

	// 10.6 مخفی کردن بخشی از رشته
	fmt.Println("\n--- 10.6 Mask String ---")
	creditCard := "1234567812345678"
	masked := creditCard[:4] + strings.Repeat("*", len(creditCard)-8) + creditCard[len(creditCard)-4:]
	fmt.Printf("  Original: %s\n", creditCard)
	fmt.Printf("  Masked: %s\n", masked)

	// ============================================================================
	// بخش 11: اشتباهات رایج
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("❌ SECTION 11: COMMON MISTAKES")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println("\n❌ Mistake 1: Compare با == برای حساسیت به حروف")
	fmt.Println("   \"Go\" == \"go\"  // false")
	fmt.Println("   ✅ Use strings.EqualFold(\"Go\", \"go\") // true")

	fmt.Println("\n❌ Mistake 2: Split با جداکننده خالی")
	fmt.Println("   strings.Split(\"hello\", \"\")  // [h e l l o]")
	fmt.Println("   ✅ Use strings.Split after validation")

	fmt.Println("\n❌ Mistake 3: فراموش کردن فضای خالی در Trim")
	fmt.Println("   strings.Trim(\"  hello  \", \" \")  // \"hello\"")
	fmt.Println("   ✅ Use strings.TrimSpace() for all whitespace")

	fmt.Println("\n❌ Mistake 4: اصلاح مستقیم رشته (غیرممکن)")
	fmt.Println("   s[0] = 'H'  // compile error (strings are immutable)")
	fmt.Println("   ✅ Convert to []byte first or use builder")

	// ============================================================================
	// بخش 12: جمع‌بندی و جدول مرجع سریع
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📚 SECTION 12: QUICK REFERENCE")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ CATEGORY          │ FUNCTIONS                                    │")
	fmt.Println("├───────────────────┼──────────────────────────────────────────────┤")
	fmt.Println("│ Comparison        │ Compare, EqualFold                           │")
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
	fmt.Println("│ Case Conversion   │ ToLower, ToUpper, ToTitle, ToLowerSpecial,   │")
	fmt.Println("│                   │ ToUpperSpecial, ToTitleSpecial               │")
	fmt.Println("│ Prefix/Suffix     │ HasPrefix, HasSuffix                         │")
	fmt.Println("│ Build/Repeat      │ Repeat, Builder (WriteString, WriteByte,     │")
	fmt.Println("│                   │ WriteRune, Grow, Len, Cap, Reset)            │")
	fmt.Println("│ Other             │ Clone, Cut, CutPrefix, CutSuffix, Map,       │")
	fmt.Println("│                   │ NewReader                                     │")
	fmt.Println("└───────────────────┴──────────────────────────────────────────────┘")

	fmt.Println("\n💡 GOLDEN RULES:")
	fmt.Println("  1. Strings are immutable - always assign result to new variable")
	fmt.Println("  2. Use strings.Builder for efficient string concatenation")
	fmt.Println("  3. Pre-allocate builder capacity with Grow() when size is known")
	fmt.Println("  4. Use EqualFold for case-insensitive comparison")
	fmt.Println("  5. Use Fields for splitting on whitespace")
	fmt.Println("  6. Use TrimSpace for removing all whitespace")
	fmt.Println("  7. Use Contains instead of Index > -1 for readability")
	fmt.Println("  8. Use Cut for simple two-part splits (Go 1.18+)")
	fmt.Println("  9. Clone when you need a copy to avoid reference issues")
	fmt.Println("  10. Never modify strings directly - convert to []byte first")

	fmt.Println("\n🎯 PERFORMANCE TIPS:")
	fmt.Println("  • strings.Builder is fastest for many concatenations")
	fmt.Println("  • Pre-allocate builder capacity: builder.Grow(estimatedSize)")
	fmt.Println("  • Join is faster than manual concatenation with +")
	fmt.Println("  • For single concatenation, + is fine")
	fmt.Println("  • Use ContainsRune for single character search")
	fmt.Println("  • Fields is optimized for whitespace splitting")
}

// تابع کمکی برای تکرار رشته (در خود پکیج وجود دارد، اینجا فقط برای مثال)
// func stringsRepeat(s string, count int) string {
//     return strings.Repeat(s, count)
// }
