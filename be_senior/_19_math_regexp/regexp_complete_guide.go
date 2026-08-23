// ============================================================================
// FILE: regexp_complete_guide.go
// TITLE: راهنمای کامل پکیج regexp در Go - عبارات منظم
// HOW TO RUN: go run regexp_complete_guide.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - پکیج regexp چیست؟
// ============================================================================
//
// پکیج regexp پیاده‌سازی سریع و امن عبارات منظم (Regular Expressions) در Go است.
// از سینتکس RE2 استفاده می‌کند که در برابر حملات ReDoS امن است.
//
// کاربردهای اصلی:
// 1. اعتبارسنجی (Validation): ایمیل، شماره تلفن، کد پستی
// 2. جستجو و استخراج (Search/Extract): پیدا کردن الگوها در متن
// 3. جایگزینی (Replace): تغییر فرمت داده‌ها
// 4. تقسیم (Split): جدا کردن رشته بر اساس الگو
//
// قانون طلایی:
// "برای ساده‌ترین عملیات روی رشته (مثل Contains, HasPrefix) از توابع strings استفاده کن.
//  از regexp فقط زمانی استفاده کن که الگوی جستجو پیچیده است.
//  همیشه regexp را با MustCompile کامپایل کن و به عنوان متغیر سراسری ذخیره کن."
// ============================================================================

package __internal_packages

import (
	"fmt"
	"regexp"
	"strings"
)

func main() {
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("🎯 COMPLETE regexp PACKAGE GUIDE IN GO")
	fmt.Println("Regular Expressions - Match, Find, Replace, Split")
	fmt.Println(strings.Repeat("=", 80))

	// ============================================================================
	// بخش 1: کامپایل کردن الگو (Compile)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🔧 SECTION 1: COMPILING PATTERNS")
	fmt.Println(strings.Repeat("=", 80))

	// 1.1 Compile - کامپایل با احتمال خطا
	fmt.Println("\n--- 1.1 regexp.Compile ---")
	pattern := `^[a-z]+$`
	re, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Printf("  Compile error: %v\n", err)
	} else {
		fmt.Printf("  Compiled: %s\n", re.String())
	}

	// 1.2 MustCompile - کامپایل بدون خطا (panic در صورت خطا)
	fmt.Println("\n--- 1.2 regexp.MustCompile ---")
	re2 := regexp.MustCompile(`\d+`)
	fmt.Printf("  MustCompile: %s\n", re2.String())

	// این خط panic می‌کند (کامنت شده)
	// re3 := regexp.MustCompile(`[invalid`)

	// 1.3 بررسی الگوهای نامعتبر
	fmt.Println("\n--- 1.3 Invalid Patterns ---")
	invalidPatterns := []string{
		`[a-z`,      // unclosed character class
		`\d{1,3,5}`, // invalid repetition
		`(a`,        // unclosed group
	}
	for _, p := range invalidPatterns {
		_, err := regexp.Compile(p)
		fmt.Printf("  Pattern %q: error = %v\n", p, err)
	}

	// ============================================================================
	// بخش 2: توابع سطح بالا (Global Functions)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🌐 SECTION 2: GLOBAL FUNCTIONS (Without Pre-compilation)")
	fmt.Println(strings.Repeat("=", 80))

	// 2.1 MatchString - بررسی تطابق الگو با رشته
	fmt.Println("\n--- 2.1 regexp.MatchString ---")
	matched, _ := regexp.MatchString(`\d+`, "abc123def")
	fmt.Printf("  MatchString(`\\d+`, \"abc123def\"): %v\n", matched)

	matched, _ = regexp.MatchString(`^[A-Z]+$`, "Hello")
	fmt.Printf("  MatchString(`^[A-Z]+$`, \"Hello\"): %v\n", matched)

	// 2.2 Match - بررسی تطابق با []byte
	fmt.Println("\n--- 2.2 regexp.Match ---")
	matchedBytes, _ := regexp.Match(`\w+`, []byte("hello world"))
	fmt.Printf("  Match(`\\w+`, []byte(\"hello world\")): %v\n", matchedBytes)

	// 2.3 MatchReader - بررسی تطابق با io.RuneReader
	fmt.Println("\n--- 2.3 regexp.MatchReader ---")
	// (معمولاً کمتر استفاده می‌شود)

	// 2.4 QuoteMeta - escape کردن کاراکترهای خاص
	fmt.Println("\n--- 2.4 regexp.QuoteMeta ---")
	literal := regexp.QuoteMeta(`[a-z]+\d+`)
	fmt.Printf("  QuoteMeta(`[a-z]+\\d+`) = %q\n", literal)

	// ============================================================================
	// بخش 3: توابع تطابق (Match Functions on Regexp)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("✅ SECTION 3: MATCH FUNCTIONS")
	fmt.Println(strings.Repeat("=", 80))

	reMatch := regexp.MustCompile(`\b[A-Z][a-z]+\b`)

	// 3.1 MatchString - بررسی تطابق کامل
	fmt.Println("\n--- 3.1 MatchString ---")
	fmt.Printf("  MatchString(\"Hello world\"): %v\n", reMatch.MatchString("Hello world"))
	fmt.Printf("  MatchString(\"hello world\"): %v\n", reMatch.MatchString("hello world"))

	// 3.2 Match - بررسی تطابق با []byte
	fmt.Println("\n--- 3.2 Match ---")
	fmt.Printf("  Match([]byte(\"Hello\")): %v\n", reMatch.Match([]byte("Hello")))

	// ============================================================================
	// بخش 4: جستجو (Find Functions)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🔍 SECTION 4: FIND FUNCTIONS")
	fmt.Println(strings.Repeat("=", 80))

	text := "The quick brown fox jumps over the lazy dog. The fox is quick."
	reFind := regexp.MustCompile(`\b\w{4}\b`) // کلمات 4 حرفی

	// 4.1 FindString - اولین تطابق
	fmt.Println("\n--- 4.1 FindString ---")
	found := reFind.FindString(text)
	fmt.Printf("  FindString: %q\n", found)

	// 4.2 Find - اولین تطابق ([]byte)
	fmt.Println("\n--- 4.2 Find ---")
	foundBytes := reFind.Find([]byte(text))
	fmt.Printf("  Find: %q\n", foundBytes)

	// 4.3 FindStringIndex - اندیس اولین تطابق
	fmt.Println("\n--- 4.3 FindStringIndex ---")
	indices := reFind.FindStringIndex(text)
	fmt.Printf("  FindStringIndex: %v (%q)\n", indices, text[indices[0]:indices[1]])

	// 4.4 FindAllString - همه تطابق‌ها
	fmt.Println("\n--- 4.4 FindAllString ---")
	allMatches := reFind.FindAllString(text, -1)
	fmt.Printf("  FindAllString: %q\n", allMatches)

	// با محدودیت تعداد
	limitedMatches := reFind.FindAllString(text, 2)
	fmt.Printf("  FindAllString (limit 2): %q\n", limitedMatches)

	// 4.5 FindAllStringSubmatch - تطابق با گروه‌ها
	fmt.Println("\n--- 4.5 FindAllStringSubmatch ---")
	reSubmatch := regexp.MustCompile(`(\w+)@(\w+)\.(\w+)`)
	emailText := "Contact: ali@gmail.com and sara@yahoo.com"

	submatches := reSubmatch.FindAllStringSubmatch(emailText, -1)
	for i, match := range submatches {
		fmt.Printf("  Match %d: %q\n", i, match[0])
		fmt.Printf("    Username: %q\n", match[1])
		fmt.Printf("    Domain: %q\n", match[2])
		fmt.Printf("    TLD: %q\n", match[3])
	}

	// 4.6 FindStringSubmatchIndex - اندیس گروه‌ها
	fmt.Println("\n--- 4.6 FindStringSubmatchIndex ---")
	reIndex := regexp.MustCompile(`(\d{4})-(\d{2})-(\d{2})`)
	date := "2024-01-15"
	indices2 := reIndex.FindStringSubmatchIndex(date)
	fmt.Printf("  Pattern: %s\n", date)
	fmt.Printf("  Indices: %v\n", indices2)
	for i := 0; i < len(indices2); i += 2 {
		start, end := indices2[i], indices2[i+1]
		if start >= 0 {
			fmt.Printf("    Group %d: %q\n", i/2, date[start:end])
		}
	}

	// ============================================================================
	// بخش 5: جایگزینی (Replace Functions)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🔄 SECTION 5: REPLACE FUNCTIONS")
	fmt.Println(strings.Repeat("=", 80))

	textReplace := "Hello World! Hello Go! Hello Regexp!"
	reReplace := regexp.MustCompile(`Hello`)

	// 5.1 ReplaceAllString - جایگزینی همه
	fmt.Println("\n--- 5.1 ReplaceAllString ---")
	result := reReplace.ReplaceAllString(textReplace, "Hi")
	fmt.Printf("  ReplaceAllString: %q\n", result)

	// 5.2 ReplaceAllLiteralString - جایگزینی بدون تفسیر الگو
	fmt.Println("\n--- 5.2 ReplaceAllLiteralString ---")
	reLiteral := regexp.MustCompile(`\d+`)
	result2 := reLiteral.ReplaceAllLiteralString("Price: $100", "$50")
	fmt.Printf("  ReplaceAllLiteralString: %q\n", result2)

	// 5.3 ReplaceAllStringFunc - جایگزینی با تابع
	fmt.Println("\n--- 5.3 ReplaceAllStringFunc ---")
	reFunc := regexp.MustCompile(`\d+`)
	result3 := reFunc.ReplaceAllStringFunc("I have 5 apples and 3 oranges", func(s string) string {
		// تبدیل عدد به حروف
		numbers := map[string]string{"5": "five", "3": "three"}
		return numbers[s]
	})
	fmt.Printf("  ReplaceAllStringFunc: %q\n", result3)

	// 5.4 جایگزینی با گروه‌ها (با استفاده از $1, $2)
	fmt.Println("\n--- 5.4 Replace with Groups ---")
	reGroupReplace := regexp.MustCompile(`(\d{4})-(\d{2})-(\d{2})`)
	date2 := "2024-01-15"

	// تغییر فرمت از YYYY-MM-DD به DD/MM/YYYY
	reformatted := reGroupReplace.ReplaceAllString(date2, "$3/$2/$1")
	fmt.Printf("  Original: %s\n", date2)
	fmt.Printf("  Reformatted: %s\n", reformatted)

	// با نام گروه‌ها
	reNamed := regexp.MustCompile(`(?P<year>\d{4})-(?P<month>\d{2})-(?P<day>\d{2})`)
	reformatted2 := reNamed.ReplaceAllString(date2, "${day}/${month}/${year}")
	fmt.Printf("  With named groups: %s\n", reformatted2)

	// 5.5 ReplaceAll - کار با []byte
	fmt.Println("\n--- 5.5 ReplaceAll ---")
	reBytes := regexp.MustCompile(`foo`)
	resultBytes := reBytes.ReplaceAll([]byte("foo bar foo"), []byte("baz"))
	fmt.Printf("  ReplaceAll: %q\n", resultBytes)

	// ============================================================================
	// بخش 6: تقسیم (Split)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("✂️ SECTION 6: SPLIT FUNCTION")
	fmt.Println(strings.Repeat("=", 80))

	reSplit := regexp.MustCompile(`\s+`) // یک یا چند فاصله
	textSplit := "hello   world  from   go"

	// 6.1 Split - تقسیم بر اساس الگو
	fmt.Println("\n--- 6.1 Split ---")
	parts := reSplit.Split(textSplit, -1)
	fmt.Printf("  Split: %q\n", parts)

	// با محدودیت تعداد
	partsLimited := reSplit.Split(textSplit, 3)
	fmt.Printf("  Split (limit 3): %q\n", partsLimited)

	// 6.2 Split با الگوهای مختلف
	fmt.Println("\n--- 6.2 Split with Different Patterns ---")

	// تقسیم بر اساس کاما
	reComma := regexp.MustCompile(`,`)
	csv := "apple,banana,cherry,date"
	fruits := reComma.Split(csv, -1)
	fmt.Printf("  CSV split: %q\n", fruits)

	// تقسیم بر اساس اعداد
	reDigit := regexp.MustCompile(`\d+`)
	mixed := "a1b2c3d4e5"
	letters := reDigit.Split(mixed, -1)
	fmt.Printf("  Split by digits: %q\n", letters)

	// ============================================================================
	// بخش 7: متدهای مفید دیگر
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🔧 SECTION 7: OTHER USEFUL METHODS")
	fmt.Println(strings.Repeat("=", 80))

	reMisc := regexp.MustCompile(`\d+`)

	// 7.1 String - نمایش الگو
	fmt.Println("\n--- 7.1 String() ---")
	fmt.Printf("  Pattern: %s\n", reMisc.String())

	// 7.2 NumSubexp - تعداد گروه‌ها
	fmt.Println("\n--- 7.2 NumSubexp ---")
	reGroups := regexp.MustCompile(`(\w+)@(\w+)\.(\w+)`)
	fmt.Printf("  Number of subexpressions: %d\n", reGroups.NumSubexp())

	// 7.3 SubexpNames - نام گروه‌ها
	fmt.Println("\n--- 7.3 SubexpNames ---")
	reNamed2 := regexp.MustCompile(`(?P<user>\w+)@(?P<domain>\w+)\.(?P<tld>\w+)`)
	names := reNamed2.SubexpNames()
	fmt.Printf("  Group names: %v\n", names)

	// 7.4 LiteralPrefix - پیشوند ثابت الگو
	fmt.Println("\n--- 7.4 LiteralPrefix ---")
	rePrefix := regexp.MustCompile(`https?://[a-z]+\.com`)
	prefix, complete := rePrefix.LiteralPrefix()
	fmt.Printf("  Literal prefix: %q, complete: %v\n", prefix, complete)

	// 7.5 Longest - اولویت دادن به طولانی‌ترین تطابق
	fmt.Println("\n--- 7.5 Longest() ---")
	reLongest := regexp.MustCompile(`a|ab`)
	fmt.Printf("  Without Longest (a|ab): %q\n", reLongest.FindString("ab"))

	reLongest.Longest()
	fmt.Printf("  With Longest (a|ab): %q\n", reLongest.FindString("ab"))

	// ============================================================================
	// بخش 8: الگوهای پرکاربرد (Common Patterns)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📝 SECTION 8: COMMON REGEX PATTERNS")
	fmt.Println(strings.Repeat("=", 80))

	// 8.1 اعتبارسنجی ایمیل
	fmt.Println("\n--- 8.1 Email Validation ---")
	emailPattern := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	emails := []string{
		"user@example.com",
		"user.name@domain.co.uk",
		"invalid-email",
		"user@",
	}
	for _, email := range emails {
		fmt.Printf("  %s: %v\n", email, emailPattern.MatchString(email))
	}

	// 8.2 اعتبارسنجی شماره تلفن ایران
	fmt.Println("\n--- 8.2 Iran Phone Number Validation ---")
	phonePattern := regexp.MustCompile(`^09[0-9]{9}$`)
	phones := []string{
		"09123456789",
		"09301234567",
		"0912345678",   // 10 digits
		"091234567890", // 12 digits
	}
	for _, phone := range phones {
		fmt.Printf("  %s: %v\n", phone, phonePattern.MatchString(phone))
	}

	// 8.3 اعتبارسنجی کد پستی ایران
	fmt.Println("\n--- 8.3 Iran Postal Code Validation ---")
	postalPattern := regexp.MustCompile(`^\d{5}-\d{5}$|^\d{10}$`)
	postalCodes := []string{
		"12345-67890",
		"1234567890",
		"1234-56789",
	}
	for _, code := range postalCodes {
		fmt.Printf("  %s: %v\n", code, postalPattern.MatchString(code))
	}

	// 8.4 استخراج URL
	fmt.Println("\n--- 8.4 URL Extraction ---")
	urlPattern := regexp.MustCompile(`https?://[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}(/[a-zA-Z0-9/.-]*)?`)
	textUrl := "Visit https://golang.org and http://example.com/page for more info"
	urls := urlPattern.FindAllString(textUrl, -1)
	fmt.Printf("  Found URLs: %v\n", urls)

	// 8.5 استخراج هگزادسیمال
	fmt.Println("\n--- 8.5 Hex Color Extraction ---")
	hexPattern := regexp.MustCompile(`#[0-9A-Fa-f]{6}\b`)
	cssText := "color: #FF0000; background: #00FF00; border: #0000FF"
	colors2 := hexPattern.FindAllString(cssText, -1)
	fmt.Printf("  Hex colors: %v\n", colors2)

	// 8.6 استخراج اعداد با جداکننده هزارگان
	fmt.Println("\n--- 8.6 Number with Commas ---")
	numberPattern := regexp.MustCompile(`\d{1,3}(?:,\d{3})*`)
	numText := "The price is 1,234,567 and discount is 123,456"
	numbers2 := numberPattern.FindAllString(numText, -1)
	fmt.Printf("  Numbers: %v\n", numbers2)

	// ============================================================================
	// بخش 9: کاربردهای عملی
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("💼 SECTION 9: PRACTICAL APPLICATIONS")
	fmt.Println(strings.Repeat("=", 80))

	// 9.1 حذف HTML تگ‌ها
	fmt.Println("\n--- 9.1 Remove HTML Tags ---")
	htmlText := "<p>Hello <b>World</b>!</p>"
	htmlPattern := regexp.MustCompile(`<[^>]*>`)
	cleanText := htmlPattern.ReplaceAllString(htmlText, "")
	fmt.Printf("  Original: %s\n", htmlText)
	fmt.Printf("  Clean: %s\n", cleanText)

	// 9.2 تبدیل camelCase به snake_case
	fmt.Println("\n--- 9.2 CamelCase to snake_case ---")
	camelToSnake := regexp.MustCompile(`([A-Z])`)
	camelStr := "userFirstName"
	snakeStr := camelToSnake.ReplaceAllStringFunc(camelStr, func(s string) string {
		return "_" + strings.ToLower(s)
	})
	snakeStr = strings.TrimPrefix(snakeStr, "_")
	fmt.Printf("  CamelCase: %s\n", camelStr)
	fmt.Printf("  snake_case: %s\n", snakeStr)

	// 9.3 اعتبارسنجی رمز عبور قوی
	fmt.Println("\n--- 9.3 Strong Password Validation ---")
	// حداقل 8 کاراکتر، حداقل یک حرف بزرگ، یک حرف کوچک، یک عدد، یک کاراکتر خاص
	passwordPattern := regexp.MustCompile(`^.{8,}$`)
	hasUpper := regexp.MustCompile(`[A-Z]`)
	hasLower := regexp.MustCompile(`[a-z]`)
	hasDigit := regexp.MustCompile(`\d`)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`)

	validatePassword := func(pwd string) bool {
		return passwordPattern.MatchString(pwd) &&
			hasUpper.MatchString(pwd) &&
			hasLower.MatchString(pwd) &&
			hasDigit.MatchString(pwd) &&
			hasSpecial.MatchString(pwd)
	}

	passwords := []string{
		"weak",
		"WeakPassword",
		"WeakPass1",
		"Strong@Pass123",
	}
	for _, pwd := range passwords {
		fmt.Printf("  %s: %v\n", pwd, validatePassword(pwd))
	}

	// 9.4 استخراج نام فایل از مسیر
	fmt.Println("\n--- 9.4 Extract Filename from Path ---")
	pathPattern := regexp.MustCompile(`[^/\\]+$`)
	paths := []string{
		"/home/user/file.txt",
		"C:\\Users\\file.txt",
		"relative/path/document.pdf",
	}
	for _, p := range paths {
		filename := pathPattern.FindString(p)
		fmt.Printf("  %s -> %s\n", p, filename)
	}

	// 9.5 بررسی فرمت JSON
	fmt.Println("\n--- 9.5 JSON Format Validation (Basic) ---")
	jsonPattern := regexp.MustCompile(`^\s*(\{.*\}|\[.*\])\s*$`)
	jsonStrings := []string{
		`{"name":"Ali"}`,
		`[1,2,3]`,
		`not json`,
		`{"invalid"}`,
	}
	for _, js := range jsonStrings {
		fmt.Printf("  %s: %v\n", js, jsonPattern.MatchString(js))
	}

	// ============================================================================
	// بخش 10: اشتباهات رایج
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("❌ SECTION 10: COMMON MISTAKES")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println("\n❌ Mistake 1: Recompiling regexp in loops")
	fmt.Println("   for i := 0; i < 1000; i++ {")
	fmt.Println("       re, _ := regexp.Compile(pattern)  // SLOW!")
	fmt.Println("   }")
	fmt.Println("   ✅ Compile once globally: var re = regexp.MustCompile(pattern)")

	fmt.Println("\n❌ Mistake 2: Not escaping literal strings")
	fmt.Println("   re := regexp.MustCompile(`[a-z]`)  // matches a-z")
	fmt.Println("   ✅ Use QuoteMeta: regexp.QuoteMeta(`[a-z]`)")

	fmt.Println("\n❌ Mistake 3: Greedy vs Non-greedy confusion")
	fmt.Println("   `.*` matches as much as possible (greedy)")
	fmt.Println("   `.*?` matches as little as possible (non-greedy)")

	fmt.Println("\n❌ Mistake 4: Using regexp for simple string operations")
	fmt.Println("   regexp.MustCompile(`hello`).MatchString(s)")
	fmt.Println("   ✅ Use strings.Contains(s, \"hello\") - faster!")

	fmt.Println("\n❌ Mistake 5: Not handling Compile errors")
	fmt.Println("   re := regexp.MustCompile(pattern)  // panics on invalid")
	fmt.Println("   ✅ Use Compile for dynamic/user-provided patterns")

	fmt.Println("\n❌ Mistake 6: Assuming ^ and $ mean line boundaries")
	fmt.Println("   By default, ^ and $ match start/end of entire string")
	fmt.Println("   ✅ Use (?m) flag: (?m)^\\d+$ for multiline")

	fmt.Println("\n❌ Mistake 7: Catastrophic backtracking (ReDoS)")
	fmt.Println("   `(a+)+$` can be extremely slow on certain inputs")
	fmt.Println("   ✅ RE2 engine in Go is safe (no backtracking)")

	// ============================================================================
	// بخش 11: جمع‌بندی و جدول مرجع
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📚 SECTION 11: QUICK REFERENCE")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ COMPILE FUNCTIONS                                             │")
	fmt.Println("├─────────────────────────────────────────────────────────────────┤")
	fmt.Println("│ regexp.Compile(pattern)   - Compile with error return         │")
	fmt.Println("│ regexp.MustCompile(pattern) - Compile or panic                 │")
	fmt.Println("│ regexp.QuoteMeta(s)       - Escape special characters         │")
	fmt.Println("└─────────────────────────────────────────────────────────────────┘")

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ MATCH FUNCTIONS                                               │")
	fmt.Println("├─────────────────────────────────────────────────────────────────┤")
	fmt.Println("│ re.MatchString(s)   - Check if pattern matches string         │")
	fmt.Println("│ re.Match(b)         - Check if pattern matches []byte         │")
	fmt.Println("└─────────────────────────────────────────────────────────────────┘")

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ FIND FUNCTIONS                                                │")
	fmt.Println("├─────────────────────────────────────────────────────────────────┤")
	fmt.Println("│ re.FindString(s)           - First match                      │")
	fmt.Println("│ re.FindStringIndex(s)      - Indices of first match           │")
	fmt.Println("│ re.FindAllString(s, n)     - All matches (limit n)            │")
	fmt.Println("│ re.FindStringSubmatch(s)   - Match with groups                │")
	fmt.Println("│ re.FindAllStringSubmatch(s, n) - All matches with groups      │")
	fmt.Println("└─────────────────────────────────────────────────────────────────┘")

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ REPLACE FUNCTIONS                                             │")
	fmt.Println("├─────────────────────────────────────────────────────────────────┤")
	fmt.Println("│ re.ReplaceAllString(s, repl)   - Replace all matches          │")
	fmt.Println("│ re.ReplaceAllLiteralString(s, repl) - Literal replacement     │")
	fmt.Println("│ re.ReplaceAllStringFunc(s, fn) - Replace with function        │")
	fmt.Println("└─────────────────────────────────────────────────────────────────┘")

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ OTHER FUNCTIONS                                               │")
	fmt.Println("├─────────────────────────────────────────────────────────────────┤")
	fmt.Println("│ re.Split(s, n)          - Split string by pattern             │")
	fmt.Println("│ re.String()             - Return pattern string               │")
	fmt.Println("│ re.NumSubexp()          - Number of capturing groups          │")
	fmt.Println("│ re.SubexpNames()        - Names of capturing groups           │")
	fmt.Println("│ re.LiteralPrefix()      - Find literal prefix                 │")
	fmt.Println("│ re.Longest()            - Prefer longest match                │")
	fmt.Println("└─────────────────────────────────────────────────────────────────┘")

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ COMMON PATTERNS                                               │")
	fmt.Println("├─────────────────────────────────────────────────────────────────┤")
	fmt.Println("│ \\d          - Digit [0-9]                                     │")
	fmt.Println("│ \\w          - Word [a-zA-Z0-9_]                               │")
	fmt.Println("│ \\s          - Whitespace                                      │")
	fmt.Println("│ .           - Any character (except newline)                  │")
	fmt.Println("│ ^           - Start of string/line                            │")
	fmt.Println("│ $           - End of string/line                              │")
	fmt.Println("│ *           - Zero or more                                    │")
	fmt.Println("│ +           - One or more                                     │")
	fmt.Println("│ ?           - Zero or one                                     │")
	fmt.Println("│ {n}         - Exactly n times                                 │")
	fmt.Println("│ {n,m}       - n to m times                                    │")
	fmt.Println("│ [abc]       - Any of a, b, or c                               │")
	fmt.Println("│ [^abc]      - Anything except a, b, c                         │")
	fmt.Println("│ (abc)       - Capturing group                                 │")
	fmt.Println("│ (?:abc)     - Non-capturing group                             │")
	fmt.Println("│ (?P<name>...) - Named group                                   │")
	fmt.Println("│ a|b         - Alternation (a or b)                            │")
	fmt.Println("└─────────────────────────────────────────────────────────────────┘")

	fmt.Println("\n💡 GOLDEN RULES:")
	fmt.Println("  1. Always compile regexp once, reuse the compiled object")
	fmt.Println("  2. Use MustCompile for patterns that are valid at compile time")
	fmt.Println("  3. Use Compile for user-provided or dynamic patterns")
	fmt.Println("  4. Use strings functions for simple operations (faster)")
	fmt.Println("  5. Use non-greedy quantifiers (.*?) when needed")
	fmt.Println("  6. Use (?m) flag for multiline mode")
	fmt.Println("  7. Escape user input with QuoteMeta")
	fmt.Println("  8. Test regexp with regex101.com before implementing")
	fmt.Println("  9. RE2 engine doesn't support backreferences (\\1)")
	fmt.Println("  10. Keep patterns simple to avoid performance issues")

	fmt.Println("\n🎯 FLAGS (used as prefix in pattern):")
	fmt.Println("  (?i) - Case-insensitive")
	fmt.Println("  (?m) - Multiline (^ and $ match line boundaries)")
	fmt.Println("  (?s) - Dot matches newline (let . match \\n)")
	fmt.Println("  (?U) - Ungreedy (swap meaning of * and *?)")
	fmt.Println("  Example: (?i)abc matches \"ABC\", \"abc\", \"Abc\"")
}

// تابع کمکی برای تکرار رشته
func stringsRepeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}
