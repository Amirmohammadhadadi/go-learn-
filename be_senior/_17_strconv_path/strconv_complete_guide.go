// ============================================================================
// FILE: strconv_complete_guide.go
// TITLE: راهنمای کامل پکیج strconv در Go - تبدیل رشته به انواع داده و بالعکس
// HOW TO RUN: go run strconv_complete_guide.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - پکیج strconv چیست؟
// ============================================================================
//
// پکیج strconv (String Conversion) توابع تبدیل بین رشته (string) و انواع پایه را ارائه می‌دهد:
// 1. تبدیل رشته به bool، int، float، uint
// 2. تبدیل bool، int، float، uint به رشته
// 3. افزودن/حذف علامت نقل قول (quoting/unquoting)
// 4. تبدیل رشته به rune (Atoi, Itoa)
// 5. کار با اعداد در پایه‌های مختلف (binary, octal, hex)
//
// قانون طلایی:
// "برای تبدیل رشته به عدد از Parse استفاده کن، برای تبدیل عدد به رشته از Format استفاده کن.
//  همیشه خطاهای Parse را بررسی کن. برای اعداد صحیح ساده، Atoi و Itoa کافی هستند."
// ============================================================================

package __internal_packages

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("🎯 COMPLETE strconv PACKAGE GUIDE IN GO")
	fmt.Println("String conversions for all basic types")
	fmt.Println(strings.Repeat("=", 80))

	// ============================================================================
	// بخش 1: تبدیل به رشته (Format functions)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📝 SECTION 1: FORMAT FUNCTIONS (Convert to string)")
	fmt.Println(strings.Repeat("=", 80))

	// 1.1 FormatBool - تبدیل bool به رشته
	fmt.Println("\n--- 1.1 strconv.FormatBool ---")
	fmt.Printf("  FormatBool(true) = %q\n", strconv.FormatBool(true))
	fmt.Printf("  FormatBool(false) = %q\n", strconv.FormatBool(false))

	// 1.2 FormatInt - تبدیل int64 به رشته با پایه دلخواه
	fmt.Println("\n--- 1.2 strconv.FormatInt ---")
	// FormatInt returns the string representation of i in the given base, 2 <= base <= 36.
	number := int64(42)
	bases := []int{2, 8, 10, 16, 36}
	for _, base := range bases {
		result := strconv.FormatInt(number, base)
		fmt.Printf("  FormatInt(%d, %d) = %q (base %d)\n", number, base, result, base)
	}

	// اعداد منفی
	negNumber := int64(-42)
	for _, base := range bases[:3] {
		result := strconv.FormatInt(negNumber, base)
		fmt.Printf("  FormatInt(%d, %d) = %q\n", negNumber, base, result)
	}

	// 1.3 FormatUint - تبدیل uint64 به رشته
	fmt.Println("\n--- 1.3 strconv.FormatUint ---")
	unumber := uint64(42)
	for _, base := range bases {
		result := strconv.FormatUint(unumber, base)
		fmt.Printf("  FormatUint(%d, %d) = %q (base %d)\n", unumber, base, result, base)
	}

	// 1.4 FormatFloat - تبدیل float64 به رشته
	fmt.Println("\n--- 1.4 strconv.FormatFloat ---")
	// FormatFloat converts the floating-point number f to a string
	// fmt: 'f' (-ddd.dddd), 'e' (-d.dddde±dd), 'E' (-d.ddddE±dd), 'g' ('e' for large exponents, 'f' otherwise), 'G' ('E' for large exponents, 'f' otherwise)
	// prec: number of digits (except for 'f' and 'g')
	floatNum := 123.456789
	formats := []struct {
		fmt  byte
		prec int
		desc string
	}{
		{'f', 2, "fixed point, 2 digits"},
		{'f', 5, "fixed point, 5 digits"},
		{'e', 5, "scientific notation"},
		{'E', 5, "scientific notation (uppercase)"},
		{'g', 5, "compact representation"},
		{'g', -1, "smallest representation"},
	}
	for _, f := range formats {
		result := strconv.FormatFloat(floatNum, f.fmt, f.prec, 64)
		fmt.Printf("  FormatFloat(%.6f, '%c', %d, 64) = %q (%s)\n",
			floatNum, f.fmt, f.prec, result, f.desc)
	}

	// 1.5 FormatFloat با اعداد خاص
	fmt.Println("\n--- 1.5 FormatFloat Special Values ---")
	specialVals := []float64{
		math_.NaN(),
		math_.Inf(1),
		math_.Inf(-1),
		0.0,
		-0.0,
	}
	for _, v := range specialVals {
		result := strconv.FormatFloat(v, 'g', -1, 64)
		fmt.Printf("  FormatFloat(%v, 'g', -1, 64) = %q\n", v, result)
	}

	// 1.6 Itoa - تبدیل int به رشته (ساده‌ترین روش)
	fmt.Println("\n--- 1.6 strconv.Itoa (Integer to ASCII) ---")
	// Itoa is shorthand for FormatInt(int64(i), 10)
	intVals := []int{-100, -1, 0, 1, 42, 100, 999}
	for _, v := range intVals {
		result := strconv.Itoa(v)
		fmt.Printf("  Itoa(%d) = %q\n", v, result)
	}

	// 1.7 Append - افزودن به slice بایت
	fmt.Println("\n--- 1.7 Append Functions ---")
	// Append functions append the string representation to a byte slice
	buf := make([]byte, 0, 100)

	buf = strconv.AppendBool(buf, true)
	buf = strconv.AppendInt(buf, 42, 10)
	buf = strconv.AppendFloat(buf, 3.14159, 'f', 5, 64)
	buf = strconv.AppendQuote(buf, "hello")

	fmt.Printf("  Appended result: %q\n", buf)

	// ============================================================================
	// بخش 2: تبدیل از رشته (Parse functions)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📖 SECTION 2: PARSE FUNCTIONS (Convert from string)")
	fmt.Println(strings.Repeat("=", 80))

	// 2.1 ParseBool - تبدیل رشته به bool
	fmt.Println("\n--- 2.1 strconv.ParseBool ---")
	// ParseBool accepts: 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False
	boolStrings := []string{"true", "TRUE", "True", "1", "t", "false", "FALSE", "False", "0", "f", "invalid"}
	for _, s := range boolStrings {
		result, err := strconv.ParseBool(s)
		if err != nil {
			fmt.Printf("  ParseBool(%q) = error: %v\n", s, err)
		} else {
			fmt.Printf("  ParseBool(%q) = %v\n", s, result)
		}
	}

	// 2.2 ParseInt - تبدیل رشته به int64
	fmt.Println("\n--- 2.2 strconv.ParseInt ---")
	// ParseInt interprets a string s in the given base (0, 2-36) and bitSize (0-64)
	intStrings := []string{
		"42", "-42", "0x2A", "0b101010", "052", "FF", "100",
	}
	for _, s := range intStrings {
		// base 0 = auto-detect (0x for hex, 0 for octal, otherwise decimal)
		result, err := strconv.ParseInt(s, 0, 64)
		if err != nil {
			fmt.Printf("  ParseInt(%q, 0, 64) = error: %v\n", s, err)
		} else {
			fmt.Printf("  ParseInt(%q, 0, 64) = %d\n", s, result)
		}
	}

	// با base مشخص
	fmt.Println("\n  With specific base:")
	hexStr := "FF"
	result, _ := strconv.ParseInt(hexStr, 16, 64)
	fmt.Printf("  ParseInt(%q, 16, 64) = %d\n", hexStr, result)

	binaryStr := "101010"
	result, _ = strconv.ParseInt(binaryStr, 2, 64)
	fmt.Printf("  ParseInt(%q, 2, 64) = %d\n", binaryStr, result)

	// 2.3 ParseUint - تبدیل رشته به uint64
	fmt.Println("\n--- 2.3 strconv.ParseUint ---")
	uintStrings := []string{"42", "0x2A", "0b101010", "052", "FF", "18446744073709551615"}
	for _, s := range uintStrings {
		result, err := strconv.ParseUint(s, 0, 64)
		if err != nil {
			fmt.Printf("  ParseUint(%q, 0, 64) = error: %v\n", s, err)
		} else {
			fmt.Printf("  ParseUint(%q, 0, 64) = %d\n", s, result)
		}
	}

	// 2.4 ParseFloat - تبدیل رشته به float64
	fmt.Println("\n--- 2.4 strconv.ParseFloat ---")
	floatStrings := []string{
		"3.14159",
		"-2.71828",
		"1.23e-4",
		"1.23E+4",
		"NaN",
		"+Inf",
		"-Inf",
		"invalid",
	}
	for _, s := range floatStrings {
		result, err := strconv.ParseFloat(s, 64)
		if err != nil {
			fmt.Printf("  ParseFloat(%q, 64) = error: %v\n", s, err)
		} else {
			fmt.Printf("  ParseFloat(%q, 64) = %v\n", s, result)
		}
	}

	// 2.5 Atoi - تبدیل رشته به int (ساده‌ترین روش)
	fmt.Println("\n--- 2.5 strconv.Atoi (ASCII to Integer) ---")
	// Atoi is shorthand for ParseInt(s, 10, 0)
	atoiStrings := []string{"123", "-456", "0", "9999999999", "abc", "12.34"}
	for _, s := range atoiStrings {
		result, err := strconv.Atoi(s)
		if err != nil {
			fmt.Printf("  Atoi(%q) = error: %v\n", s, err)
		} else {
			fmt.Printf("  Atoi(%q) = %d\n", s, result)
		}
	}

	// ============================================================================
	// بخش 3: نقل قول (Quoting and Unquoting)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📝 SECTION 3: QUOTING AND UNQUOTING")
	fmt.Println(strings.Repeat("=", 80))

	// 3.1 Quote - افزودن علامت نقل قول به رشته
	fmt.Println("\n--- 3.1 strconv.Quote ---")
	strings := []string{
		"hello",
		"hello\nworld",
		"hello\tworld",
		`hello "world"`,
		"hello \\ world",
	}
	for _, s := range strings {
		quoted := strconv.Quote(s)
		fmt.Printf("  Quote(%q) = %q\n", s, quoted)
	}

	// 3.2 QuoteToASCII - افزودن نقل قول با escape کاراکترهای غیر ASCII
	fmt.Println("\n--- 3.2 strconv.QuoteToASCII ---")
	unicodeStr := "Hello 世界"
	quotedASCII := strconv.QuoteToASCII(unicodeStr)
	fmt.Printf("  QuoteToASCII(%q) = %q\n", unicodeStr, quotedASCII)

	// 3.3 QuoteRune - نقل قول برای rune
	fmt.Println("\n--- 3.3 strconv.QuoteRune ---")
	runes := []rune{'a', '世', '\n', '\t', '\''}
	for _, r := range runes {
		quoted := strconv.QuoteRune(r)
		fmt.Printf("  QuoteRune(%q) = %s\n", r, quoted)
	}

	// 3.4 QuoteRuneToASCII - نقل قول rune با escape
	fmt.Println("\n--- 3.4 strconv.QuoteRuneToASCII ---")
	quotedRuneASCII := strconv.QuoteRuneToASCII('世')
	fmt.Printf("  QuoteRuneToASCII('世') = %s\n", quotedRuneASCII)

	// 3.5 Unquote - حذف علامت نقل قول
	fmt.Println("\n--- 3.5 strconv.Unquote ---")
	quotedStrings := []string{
		`"hello"`,
		"`hello`",
		`"hello\nworld"`,
		`"hello \"world\""`,
	}
	for _, s := range quotedStrings {
		unquoted, err := strconv.Unquote(s)
		if err != nil {
			fmt.Printf("  Unquote(%q) = error: %v\n", s, err)
		} else {
			fmt.Printf("  Unquote(%q) = %q\n", s, unquoted)
		}
	}

	// 3.6 UnquoteChar - حذف نقل قول از یک کاراکتر
	fmt.Println("\n--- 3.6 strconv.UnquoteChar ---")
	// UnquoteChar unquotes the first character or escape sequence in s
	quotedChar := `"\u4E16"`
	value, _, _, err := strconv.UnquoteChar(quotedChar, '"')
	if err == nil {
		fmt.Printf("  UnquoteChar(%q) = %c\n", quotedChar, value)
	}

	// 3.7 CanBackquote - بررسی امکان استفاده از backquote
	fmt.Println("\n--- 3.7 strconv.CanBackquote ---")
	backquoteTests := []string{
		"hello",
		"hello world",
		"hello\nworld",
		"hello`world",
		"hello\tworld",
	}
	for _, s := range backquoteTests {
		can := strconv.CanBackquote(s)
		fmt.Printf("  CanBackquote(%q) = %v\n", s, can)
	}

	// ============================================================================
	// بخش 4: Append Functions (پیشرفته)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🔧 SECTION 4: APPEND FUNCTIONS (Advanced)")
	fmt.Println(strings.Repeat("=", 80))

	// 4.1 AppendQuote - افزودن نقل قول به slice
	fmt.Println("\n--- 4.1 strconv.AppendQuote ---")
	slice := make([]byte, 0, 100)
	slice = strconv.AppendQuote(slice, "hello world")
	fmt.Printf("  AppendQuote: %q\n", slice)

	// 4.2 AppendQuoteToASCII - افزودن نقل قول ASCII
	fmt.Println("\n--- 4.2 strconv.AppendQuoteToASCII ---")
	slice2 := make([]byte, 0, 100)
	slice2 = strconv.AppendQuoteToASCII(slice2, "Hello 世界")
	fmt.Printf("  AppendQuoteToASCII: %q\n", slice2)

	// 4.3 AppendInt - افزودن عدد به slice
	fmt.Println("\n--- 4.3 strconv.AppendInt ---")
	slice3 := make([]byte, 0, 100)
	slice3 = strconv.AppendInt(slice3, 42, 10)
	slice3 = strconv.AppendInt(slice3, 255, 16)
	slice3 = strconv.AppendInt(slice3, 0b101010, 2)
	fmt.Printf("  AppendInt results: %q\n", slice3)

	// ============================================================================
	// بخش 5: نکات و ترفندها (Tips and Tricks)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("💡 SECTION 5: TIPS AND TRICKS")
	fmt.Println(strings.Repeat("=", 80))

	// 5.1 تبدیل با پیش‌تخصیص ظرفیت
	fmt.Println("\n--- 5.1 Pre-allocating Capacity ---")

	// بدون پیش‌تخصیص
	noPrealloc := make([]byte, 0)
	for i := 0; i < 1000; i++ {
		noPrealloc = strconv.AppendInt(noPrealloc, int64(i), 10)
		noPrealloc = append(noPrealloc, ',')
	}
	fmt.Printf("  Without preallocation: len=%d, cap=%d\n", len(noPrealloc), cap(noPrealloc))

	// با پیش‌تخصیص
	withPrealloc := make([]byte, 0, 5000)
	for i := 0; i < 1000; i++ {
		withPrealloc = strconv.AppendInt(withPrealloc, int64(i), 10)
		withPrealloc = append(withPrealloc, ',')
	}
	fmt.Printf("  With preallocation: len=%d, cap=%d\n", len(withPrealloc), cap(withPrealloc))

	// 5.2 تبدیل اعداد با ارور هندلینگ دقیق
	fmt.Println("\n--- 5.2 Detailed Error Handling ---")

	parseWithDetail := func(s string) {
		result, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			if e, ok := err.(*strconv.NumError); ok {
				fmt.Printf("  ParseInt(%q): %v (Func: %s, Num: %s)\n",
					s, e.Err, e.Func, e.Num)
			}
		} else {
			fmt.Printf("  ParseInt(%q) = %d\n", s, result)
		}
	}

	parseWithDetail("123")
	parseWithDetail("9999999999999999999") // overflow
	parseWithDetail("abc")                 // invalid syntax

	// 5.3 تبدیل اعداد با پایه‌های مختلف
	fmt.Println("\n--- 5.3 Converting with Different Bases ---")

	convertBase := func(s string, fromBase, toBase int) string {
		// Parse from fromBase
		num, err := strconv.ParseInt(s, fromBase, 64)
		if err != nil {
			return "error"
		}
		// Format to toBase
		return strconv.FormatInt(num, toBase)
	}

	fmt.Printf("  Convert \"FF\" from hex to decimal: %s\n", convertBase("FF", 16, 10))
	fmt.Printf("  Convert \"101010\" from binary to hex: %s\n", convertBase("101010", 2, 16))
	fmt.Printf("  Convert \"52\" from octal to binary: %s\n", convertBase("52", 8, 2))

	// 5.4 تبدیل رشته به انواع مختلف در یک خط
	fmt.Println("\n--- 5.4 One-line Conversions with Default Values ---")

	atoiDefault := func(s string, defaultVal int) int {
		if val, err := strconv.Atoi(s); err == nil {
			return val
		}
		return defaultVal
	}

	parseFloatDefault := func(s string, defaultVal float64) float64 {
		if val, err := strconv.ParseFloat(s, 64); err == nil {
			return val
		}
		return defaultVal
	}

	fmt.Printf("  atoiDefault(\"123\", 0) = %d\n", atoiDefault("123", 0))
	fmt.Printf("  atoiDefault(\"abc\", 0) = %d\n", atoiDefault("abc", 0))
	fmt.Printf("  parseFloatDefault(\"3.14\", 0) = %.2f\n", parseFloatDefault("3.14", 0))
	fmt.Printf("  parseFloatDefault(\"invalid\", 0) = %.2f\n", parseFloatDefault("invalid", 0))

	// ============================================================================
	// بخش 6: کاربردهای عملی
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("💼 SECTION 6: PRACTICAL EXAMPLES")
	fmt.Println(strings.Repeat("=", 80))

	// 6.1 پارس کردن CSV اعداد
	fmt.Println("\n--- 6.1 Parsing CSV Numbers ---")

	csvNumbers := "10,20,30,40,50"
	parts := strings.Split(csvNumbers, ",")
	numbers := make([]int, 0, len(parts))

	for _, part := range parts {
		if num, err := strconv.Atoi(part); err == nil {
			numbers = append(numbers, num)
		}
	}
	fmt.Printf("  CSV: %q -> %v\n", csvNumbers, numbers)

	// 6.2 اعتبارسنجی ورودی
	fmt.Println("\n--- 6.2 Input Validation ---")

	validateAge := func(ageStr string) (int, error) {
		age, err := strconv.Atoi(ageStr)
		if err != nil {
			return 0, fmt.Errorf("invalid age format: %w", err)
		}
		if age < 0 || age > 150 {
			return 0, fmt.Errorf("age must be between 0 and 150")
		}
		return age, nil
	}

	ages := []string{"25", "-5", "200", "abc"}
	for _, ageStr := range ages {
		if age, err := validateAge(ageStr); err != nil {
			fmt.Printf("  Age %q: error - %v\n", ageStr, err)
		} else {
			fmt.Printf("  Age %q: valid - %d\n", ageStr, age)
		}
	}

	// 6.3 فرمت کردن اعداد با هزارگان (thousands separator)
	fmt.Println("\n--- 6.3 Format Numbers with Thousands Separator ---")

	formatWithCommas := func(n int) string {
		str := strconv.Itoa(n)
		if len(str) <= 3 {
			return str
		}
		result := make([]byte, 0, len(str)+len(str)/3)
		for i, c := range []byte(str) {
			if i > 0 && (len(str)-i)%3 == 0 {
				result = append(result, ',')
			}
			result = append(result, c)
		}
		return string(result)
	}

	testNums := []int{123, 1234, 12345, 1234567, 1234567890}
	for _, n := range testNums {
		fmt.Printf("  %d -> %s\n", n, formatWithCommas(n))
	}

	// 6.4 تبدیل hex رنگ به RGB
	fmt.Println("\n--- 6.4 Hex Color to RGB ---")

	hexToRGB := func(hex string) (int, int, int, error) {
		if len(hex) != 7 || hex[0] != '#' {
			return 0, 0, 0, fmt.Errorf("invalid hex format")
		}
		r, err := strconv.ParseInt(hex[1:3], 16, 64)
		if err != nil {
			return 0, 0, 0, err
		}
		g, err := strconv.ParseInt(hex[3:5], 16, 64)
		if err != nil {
			return 0, 0, 0, err
		}
		b, err := strconv.ParseInt(hex[5:7], 16, 64)
		if err != nil {
			return 0, 0, 0, err
		}
		return int(r), int(g), int(b), nil
	}

	colors := []string{"#FF0000", "#00FF00", "#0000FF", "#FFFFFF"}
	for _, hex := range colors {
		r, g, b, err := hexToRGB(hex)
		if err != nil {
			fmt.Printf("  %s: error - %v\n", hex, err)
		} else {
			fmt.Printf("  %s -> RGB(%d, %d, %d)\n", hex, r, g, b)
		}
	}

	// 6.5 تبدیل انواع مختلف در query string
	fmt.Println("\n--- 6.5 Query String Parsing ---")

	type QueryParams struct {
		Page  int
		Limit int
		Sort  string
	}

	parseQuery := func(query map[string][]string) QueryParams {
		params := QueryParams{Page: 1, Limit: 10, Sort: "asc"}

		if pageVals, ok := query["page"]; ok && len(pageVals) > 0 {
			if page, err := strconv.Atoi(pageVals[0]); err == nil && page > 0 {
				params.Page = page
			}
		}

		if limitVals, ok := query["limit"]; ok && len(limitVals) > 0 {
			if limit, err := strconv.Atoi(limitVals[0]); err == nil && limit > 0 && limit <= 100 {
				params.Limit = limit
			}
		}

		if sortVals, ok := query["sort"]; ok && len(sortVals) > 0 {
			if sortVals[0] == "asc" || sortVals[0] == "desc" {
				params.Sort = sortVals[0]
			}
		}

		return params
	}

	query := map[string][]string{
		"page":  {"2"},
		"limit": {"20"},
		"sort":  {"desc"},
	}
	params := parseQuery(query)
	fmt.Printf("  Parsed query params: %+v\n", params)

	// ============================================================================
	// بخش 7: اشتباهات رایج
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("❌ SECTION 7: COMMON MISTAKES")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println("\n❌ Mistake 1: Not checking Parse errors")
	fmt.Println("   num, _ := strconv.Atoi(\"abc\")  // num = 0, error ignored")
	fmt.Println("   ✅ Always check the error return value")

	fmt.Println("\n❌ Mistake 2: Assuming ParseInt with base 0 always works")
	fmt.Println("   strconv.ParseInt(\"09\", 0, 64)  // error: invalid octal")
	fmt.Println("   ✅ Use base 10 explicitly if you want decimal")

	fmt.Println("\n❌ Mistake 3: Overflow in ParseInt/ParseUint")
	fmt.Println("   strconv.ParseInt(\"9999999999999999999\", 10, 32)  // overflow")
	fmt.Println("   ✅ Use larger bitSize or check error")

	fmt.Println("\n❌ Mistake 4: Using FormatFloat with too much precision")
	fmt.Println("   strconv.FormatFloat(1.2, 'f', 100, 64)  // huge string")
	fmt.Println("   ✅ Use reasonable precision (e.g., 10-15 digits)")

	fmt.Println("\n❌ Mistake 5: Forgetting that Atoi returns int (not int64)")
	fmt.Println("   var x int64 = strconv.Atoi(\"42\")  // compile error")
	fmt.Println("   ✅ Convert: x := int64(strconv.Atoi(\"42\"))")

	fmt.Println("\n❌ Mistake 6: Using Quote instead of QuoteToASCII for JSON")
	fmt.Println("   json.Marshal(\"Hello 世界\")  // works fine")
	fmt.Println("   ✅ strconv.Quote is fine, QuoteToASCII escapes non-ASCII")

	fmt.Println("\n❌ Mistake 7: Not handling NaN and Inf in ParseFloat")
	fmt.Println("   strconv.ParseFloat(\"NaN\", 64)  // returns NaN")
	fmt.Println("   ✅ Check with math_.IsNaN/IsInf if needed")

	fmt.Println("\n❌ Mistake 8: Performance issues in loops")
	fmt.Println("   for i := 0; i < 1000000; i++ { strconv.Itoa(i) }")
	fmt.Println("   ✅ Pre-allocate buffer with AppendInt")

	// ============================================================================
	// بخش 8: جمع‌بندی و جدول مرجع
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📚 SECTION 8: QUICK REFERENCE")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ FORMAT FUNCTIONS (Convert to string)                           │")
	fmt.Println("├─────────────────────────────────────────────────────────────────┤")
	fmt.Println("│ strconv.FormatBool(b)        - bool to string                  │")
	fmt.Println("│ strconv.FormatInt(i, base)   - int64 to string (base 2-36)     │")
	fmt.Println("│ strconv.FormatUint(i, base)  - uint64 to string (base 2-36)    │")
	fmt.Println("│ strconv.FormatFloat(f, fmt, prec, bits) - float to string      │")
	fmt.Println("│ strconv.Itoa(i)              - int to string (base 10)         │")
	fmt.Println("│ strconv.AppendXxx()          - Append to byte slice            │")
	fmt.Println("└─────────────────────────────────────────────────────────────────┘")

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ PARSE FUNCTIONS (Convert from string)                          │")
	fmt.Println("├─────────────────────────────────────────────────────────────────┤")
	fmt.Println("│ strconv.ParseBool(s)         - string to bool                  │")
	fmt.Println("│ strconv.ParseInt(s, base, bits) - string to int64              │")
	fmt.Println("│ strconv.ParseUint(s, base, bits) - string to uint64            │")
	fmt.Println("│ strconv.ParseFloat(s, bits)  - string to float64               │")
	fmt.Println("│ strconv.Atoi(s)              - string to int (base 10)         │")
	fmt.Println("└─────────────────────────────────────────────────────────────────┘")

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ QUOTE FUNCTIONS                                               │")
	fmt.Println("├─────────────────────────────────────────────────────────────────┤")
	fmt.Println("│ strconv.Quote(s)             - Add quotes to string           │")
	fmt.Println("│ strconv.QuoteToASCII(s)      - Quote with ASCII escaping      │")
	fmt.Println("│ strconv.QuoteRune(r)         - Quote a rune                   │")
	fmt.Println("│ strconv.QuoteRuneToASCII(r)  - Quote rune with ASCII escaping │")
	fmt.Println("│ strconv.Unquote(s)           - Remove quotes from string      │")
	fmt.Println("│ strconv.CanBackquote(s)      - Check if backquote works       │")
	fmt.Println("└─────────────────────────────────────────────────────────────────┘")

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ FORMAT SPECIFIERS FOR FormatFloat                             │")
	fmt.Println("├─────────────────────────────────────────────────────────────────┤")
	fmt.Println("│ 'f'  - decimal: -ddd.dddd                                    │")
	fmt.Println("│ 'e'  - scientific: -d.dddde±dd                               │")
	fmt.Println("│ 'E'  - scientific (uppercase): -d.ddddE±dd                   │")
	fmt.Println("│ 'g'  - compact: 'e' for large exponents, 'f' otherwise       │")
	fmt.Println("│ 'G'  - compact (uppercase): 'E' for large, 'f' otherwise     │")
	fmt.Println("│ 'b'  - binary exponent (rare)                                │")
	fmt.Println("└─────────────────────────────────────────────────────────────────┘")

	fmt.Println("\n💡 GOLDEN RULES:")
	fmt.Println("  1. Always check errors from Parse functions")
	fmt.Println("  2. Use Atoi/Itoa for base-10 integer conversion")
	fmt.Println("  3. Use ParseInt with base 0 for auto-detection")
	fmt.Println("  4. Use FormatFloat with 'g' for compact representation")
	fmt.Println("  5. Quote strings before inserting into SQL/JSON")
	fmt.Println("  6. Pre-allocate buffers with Append functions for performance")
	fmt.Println("  7. Handle overflow errors from Parse functions")
	fmt.Println("  8. Use CanBackquote before using backticks for strings")
	fmt.Println("  9. Be careful with ParseInt on octal strings (leading zero)")
	fmt.Println("  10. Use proper bitSize to avoid unnecessary conversions")

	fmt.Println("\n🎯 COMMON USE CASES:")
	fmt.Println("  • Atoi/Itoa → Simple integer conversion")
	fmt.Println("  • ParseFloat → Parsing decimal numbers from user input")
	fmt.Println("  • FormatFloat → Displaying numbers to users")
	fmt.Println("  • Quote/Unquote → Working with string literals")
	fmt.Println("  • ParseInt with base 16 → Hex color parsing")
	fmt.Println("  • ParseInt with base 2 → Binary string parsing")
}

// تابع کمکی برای تکرار رشته (برای math_.NaN و math_.Inf نیاز است)
// در برنامه واقعی باید math_ را import کنید
var math_ = struct {
	NaN func() float64
	Inf func(sign int) float64
}{
	NaN: func() float64 { return 0.0 / 0.0 },
	Inf: func(sign int) float64 {
		if sign >= 0 {
			return 1.0 / 0.0
		}
		return -1.0 / 0.0
	},
}

func stringsRepeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}
