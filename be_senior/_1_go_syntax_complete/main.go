// ============================================================================
// FILE: go_syntax_complete_guide.go
// TITLE: راهنمای کامل Syntax زبان Go - متغیرها، حلقه‌ها، شرط‌ها، توابع
// HOW TO RUN: go run go_syntax_complete_guide.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - فلسفه طراحی Syntax در Go
// ============================================================================
//
// Go زبانی است با syntax ساده و حداقلی. فلسفه طراحی:
// 1. "کمتر، بیشتر است" - تنها یک راه برای انجام هر کار وجود دارد
// 2. خوانایی مهم‌تر از اختصار است
// 3. بدون قابلیت‌های پیچیده و غیرضروری (مثل ternary operator)
// 4. کامپایلر سریع با syntax بدون ابهام
//
// در این فایل: همه چیز از صفر تا صد با مثال‌های اجرایی
// ============================================================================

package main

import (
	"fmt"
	"math"
)

// ============================================================================
// بخش 1: متغیرها (Variables) - اعلان، مقداردهی، انواع
// ============================================================================

func demonstrateVariables() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📦 VARIABLES - Declaration, Initialization, Types")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 1.1 اعلان با var (صریح)
	// ============================================
	fmt.Println("\n--- 1.1 Var Declaration ---")

	var name string               // اعلان با مقدار پیش‌فرض (zero value)
	var age int = 30              // اعلان با مقداردهی اولیه
	var salary float64 = 55000.50 // اعلان با نوع صریح
	var isActive bool = true      // boolean

	fmt.Printf("name: %q (zero value: empty string)\n", name)
	fmt.Printf("age: %d\n", age)
	fmt.Printf("salary: %.2f\n", salary)
	fmt.Printf("isActive: %t\n", isActive)

	// ============================================
	// 1.2 اعلان گروهی (Block declaration)
	// ============================================
	fmt.Println("\n--- 1.2 Block Declaration ---")

	var (
		firstName string = "Ali"
		lastName         = "Rezaei" // نوع از روی مقدار推断 می‌شود
		country          = "Iran"
		zipCode          = 12345
	)

	fmt.Printf("First: %s, Last: %s, Country: %s, Zip: %d\n",
		firstName, lastName, country, zipCode)

	// ============================================
	// 1.3 Short Declaration (:=) - فقط داخل توابع
	// ============================================
	fmt.Println("\n--- 1.3 Short Declaration (:=) ---")

	city := "Tehran"      // string
	population := 8000000 // int
	temperature := 25.5   // float64
	isRaining := false    // bool

	fmt.Printf("city: %s (%T)\n", city, city)
	fmt.Printf("population: %d (%T)\n", population, population)
	fmt.Printf("temperature: %.1f (%T)\n", temperature, temperature)
	fmt.Printf("isRaining: %t (%T)\n", isRaining, isRaining)

	// اعلان چند متغیره در یک خط
	x, y, z := 10, 20.5, "hello"
	fmt.Printf("x=%d, y=%.1f, z=%s\n", x, y, z)

	// ============================================
	// 1.4 Zero Values (مقادیر پیش‌فرض)
	// ============================================
	fmt.Println("\n--- 1.4 Zero Values ---")

	var (
		i   int     // 0
		f   float64 // 0.0
		b   bool    // false
		s   string  // "" (empty string)
		p   *int    // nil
		arr [3]int  // [0,0,0]
	)

	fmt.Printf("int: %d\n", i)
	fmt.Printf("float64: %.1f\n", f)
	fmt.Printf("bool: %t\n", b)
	fmt.Printf("string: %q\n", s)
	fmt.Printf("pointer: %v\n", p)
	fmt.Printf("array: %v\n", arr)

	// ============================================
	// 1.5 Type Conversion (تبدیل نوع)
	// ============================================
	fmt.Println("\n--- 1.5 Type Conversion ---")

	var integer int = 42
	var floatNum float64 = float64(integer) // تبدیل صریح
	var unsigned uint = uint(integer)       // تبدیل به unsigned

	fmt.Printf("int: %d -> float64: %.1f\n", integer, floatNum)
	fmt.Printf("int: %d -> uint: %d\n", integer, unsigned)

	// تبدیل بین اعداد و رشته‌ها (نیاز به strconv)
	// numStr := strconv.Itoa(42)    // "42"
	// num, _ := strconv.Atoi("42")  // 42

	// ============================================
	// 1.6 Constants (ثابت‌ها)
	// ============================================
	fmt.Println("\n--- 1.6 Constants ---")

	const pi = 3.14159
	const appName string = "MyApp"
	const (
		statusOK       = 200
		statusNotFound = 404
	)

	// iota - شمارنده خودکار برای ثابت‌ها
	const (
		Sunday    = iota // 0
		Monday           // 1
		Tuesday          // 2
		Wednesday        // 3
		Thursday         // 4
		Friday           // 5
		Saturday         // 6
	)

	fmt.Printf("pi: %.5f\n", pi)
	fmt.Printf("appName: %s\n", appName)
	fmt.Printf("Status: OK=%d, NotFound=%d\n", statusOK, statusNotFound)
	fmt.Printf("Days: Sun=%d, Mon=%d, Tue=%d\n", Sunday, Monday, Tuesday)
}

// ============================================================================
// بخش 2: حلقه‌ها (Loops) - فقط for در Go
// ============================================================================

func demonstrateLoops() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🔄 LOOPS - Only 'for' keyword in Go")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 2.1 حلقه کلاسیک (C-style for)
	// ============================================
	fmt.Println("\n--- 2.1 Classic For Loop ---")

	for i := 0; i < 5; i++ {
		fmt.Printf("  Iteration %d\n", i)
	}

	// ============================================
	// 2.2 حلقه while-style (فقط شرط)
	// ============================================
	fmt.Println("\n--- 2.2 While-style Loop ---")

	count := 0
	for count < 3 {
		fmt.Printf("  Count: %d\n", count)
		count++
	}

	// ============================================
	// 2.3 حلقه بی‌نهایت (Infinite loop)
	// ============================================
	fmt.Println("\n--- 2.3 Infinite Loop with break ---")

	counter := 0
	for {
		fmt.Printf("  Loop iteration: %d\n", counter)
		counter++
		if counter >= 3 {
			break
		}
	}

	// ============================================
	// 2.4 حلقه با continue (پرش از تکرار)
	// ============================================
	fmt.Println("\n--- 2.4 Continue Example ---")

	for i := 0; i < 10; i++ {
		if i%2 == 0 {
			continue // اعداد زوج را رد کن
		}
		fmt.Printf("  Odd number: %d\n", i)
	}

	// ============================================
	// 2.5 حلقه روی اسلایس (range)
	// ============================================
	fmt.Println("\n--- 2.5 Range Over Slice ---")

	names := []string{"Ali", "Reza", "Sara", "Mina"}

	// با اندیس و مقدار
	for index, name := range names {
		fmt.Printf("  index=%d, value=%s\n", index, name)
	}

	// فقط مقدار (با ignoring index)
	fmt.Println("\n  Only values:")
	for _, name := range names {
		fmt.Printf("    %s\n", name)
	}

	// فقط اندیس
	fmt.Println("\n  Only indices:")
	for index := range names {
		fmt.Printf("    %d\n", index)
	}

	// ============================================
	// 2.6 حلقه روی مپ (map)
	// ============================================
	fmt.Println("\n--- 2.6 Range Over Map ---")

	ages := map[string]int{
		"Ali":  30,
		"Reza": 25,
		"Sara": 28,
	}

	for name, age := range ages {
		fmt.Printf("  %s is %d years old\n", name, age)
	}

	// ============================================
	// 2.7 حلقه روی رشته (string) - Unicode aware
	// ============================================
	fmt.Println("\n--- 2.7 Range Over String (Unicode) ---")

	text := "سلام Go"
	fmt.Printf("String: %s\n", text)

	fmt.Println("  Bytes (not recommended for Unicode):")
	for i := 0; i < len(text); i++ {
		fmt.Printf("    %d: %c (%d)\n", i, text[i], text[i])
	}

	fmt.Println("  Runes (correct for Unicode):")
	for index, runeValue := range text {
		fmt.Printf("    %d: %c (U+%04X)\n", index, runeValue, runeValue)
	}

	// ============================================
	// 2.8 حلقه با break به برچسب (Labeled break)
	// ============================================
	fmt.Println("\n--- 2.8 Labeled Break ---")

outerLoop:
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if i == 1 && j == 1 {
				fmt.Printf("  Breaking at i=%d, j=%d\n", i, j)
				break outerLoop // از هر دو حلقه خارج می‌شود
			}
			fmt.Printf("  i=%d, j=%d\n", i, j)
		}
	}
}

// ============================================================================
// بخش 3: شرط‌ها (Conditionals) - if, switch
// ============================================================================

func demonstrateConditionals() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🎯 CONDITIONALS - if, if-else, switch")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 3.1 if-else ساده
	// ============================================
	fmt.Println("\n--- 3.1 Simple if-else ---")

	score := 85

	if score >= 90 {
		fmt.Println("  Grade: A")
	} else if score >= 80 {
		fmt.Println("  Grade: B")
	} else if score >= 70 {
		fmt.Println("  Grade: C")
	} else {
		fmt.Println("  Grade: F")
	}

	// ============================================
	// 3.2 if با statement (متغیر موقت)
	// ============================================
	fmt.Println("\n--- 3.2 If with Statement ---")

	// متغیر فقط در scope if و else موجود است
	if value := math.Pow(2, 3); value > 10 {
		fmt.Printf("  Value %.0f is greater than 10\n", value)
	} else {
		fmt.Printf("  Value %.0f is not greater than 10\n", value)
	}
	// value اینجا دیگر قابل دسترس نیست

	// ============================================
	// 3.3 بررسی خطا (idiomatic error checking)
	// ============================================
	fmt.Println("\n--- 3.3 Error Checking Pattern ---")

	divide := func(a, b float64) (float64, error) {
		if b == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return a / b, nil
	}

	if result, err := divide(10, 2); err != nil {
		fmt.Printf("  Error: %v\n", err)
	} else {
		fmt.Printf("  Result: %.2f\n", result)
	}

	if result, err := divide(10, 0); err != nil {
		fmt.Printf("  Error: %v\n", err)
	} else {
		fmt.Printf("  Result: %.2f\n", result)
	}

	// ============================================
	// 3.4 Switch - پایه
	// ============================================
	fmt.Println("\n--- 3.4 Basic Switch ---")

	day := 3

	switch day {
	case 1:
		fmt.Println("  Saturday")
	case 2:
		fmt.Println("  Sunday")
	case 3:
		fmt.Println("  Monday")
	default:
		fmt.Println("  Other day")
	}

	// ============================================
	// 3.5 Switch با چند مقدار در هر case
	// ============================================
	fmt.Println("\n--- 3.5 Multiple Values in Case ---")

	month := "Jan"

	switch month {
	case "Jan", "Feb", "Dec":
		fmt.Println("  Winter")
	case "Mar", "Apr", "May":
		fmt.Println("  Spring")
	case "Jun", "Jul", "Aug":
		fmt.Println("  Summer")
	case "Sep", "Oct", "Nov":
		fmt.Println("  Fall")
	default:
		fmt.Println("  Unknown")
	}

	// ============================================
	// 3.6 Switch بدون expression (به جای if-else chain)
	// ============================================
	fmt.Println("\n--- 3.6 Tagless Switch ---")

	number := 15

	switch {
	case number < 0:
		fmt.Println("  Negative")
	case number == 0:
		fmt.Println("  Zero")
	case number > 0 && number < 10:
		fmt.Println("  Small positive")
	case number >= 10 && number < 20:
		fmt.Println("  Medium positive")
	default:
		fmt.Println("  Large positive")
	}

	// ============================================
	// 3.7 Switch با fallthrough
	// ============================================
	fmt.Println("\n--- 3.7 Fallthrough Example ---")

	num := 2

	switch num {
	case 1:
		fmt.Println("  One")
	case 2:
		fmt.Println("  Two")
		fallthrough // به case بعدی هم می‌رود
	case 3:
		fmt.Println("  Three (from fallthrough)")
	case 4:
		fmt.Println("  Four")
	}

	// ============================================
	// 3.8 Switch با break (break زودتر از موعد)
	// ============================================
	fmt.Println("\n--- 3.8 Break in Switch ---")

	for i := 0; i < 5; i++ {
		switch {
		case i == 2:
			fmt.Printf("  Breaking at i=%d\n", i)
			break // فقط از switch خارج می‌شود، نه از حلقه
		default:
			fmt.Printf("  i=%d\n", i)
		}
	}

	// برای خروج از حلقه از label استفاده کن
	fmt.Println("\n  Breaking out of loop:")
loop:
	for i := 0; i < 5; i++ {
		switch {
		case i == 2:
			fmt.Printf("    Breaking loop at i=%d\n", i)
			break loop
		default:
			fmt.Printf("    i=%d\n", i)
		}
	}
}

// ============================================================================
// بخش 4: توابع (Functions) - اعلان، پارامترها، برگشتی‌ها
// ============================================================================

// ============================================
// 4.1 تابع ساده
// ============================================
func add(a int, b int) int {
	return a + b
}

// ============================================
// 4.2 پارامترهای همنوع (شورت‌نویسی)
// ============================================
func multiply(x, y int) int {
	return x * y
}

// ============================================
// 4.3 توابع با چند خروجی (multiple return values)
// ============================================
func divideWithRemainder(dividend, divisor int) (int, int) {
	quotient := dividend / divisor
	remainder := dividend % divisor
	return quotient, remainder
}

// ============================================
// 4.4 Named return values (خروجی‌های نام‌گذاری شده)
// ============================================
func rectangleAreaAndPerimeter(width, height float64) (area, perimeter float64) {
	area = width * height
	perimeter = 2 * (width + height)
	return // naked return (فقط return بنویس)
}

// ============================================
// 4.5 Variadic functions (توابع با تعداد متغیر پارامتر)
// ============================================
func sum(numbers ...int) int {
	total := 0
	for _, n := range numbers {
		total += n
	}
	return total
}

// ============================================
// 4.6 تابع بازگشتی (recursive)
// ============================================
func factorial(n int) int {
	if n <= 1 {
		return 1
	}
	return n * factorial(n-1)
}

// ============================================
// 4.7 تابع به عنوان مقدار (first-class function)
// ============================================
func applyOperation(a, b int, operation func(int, int) int) int {
	return operation(a, b)
}

// ============================================
// 4.8 Closure (بستار) - تابع درون تابع
// ============================================
func counter() func() int {
	count := 0
	return func() int {
		count++
		return count
	}
}

// ============================================
// 4.9 تابع با defer (تأخیر در اجرا)
// ============================================
func demonstrateDefer() {
	fmt.Println("\n  Starting function")

	defer fmt.Println("  This runs last (deferred)")
	defer fmt.Println("  This runs second last (LIFO)")

	fmt.Println("  Normal execution")
	// ترتیب اجرا: Normal -> This runs second last -> This runs last
}

// ============================================
// 4.10 تابع با panic و recover
// ============================================
func safeDivide(a, b int) (result int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic recovered: %v", r)
		}
	}()

	if b == 0 {
		panic("division by zero")
	}

	return a / b, nil
}

func demonstrateFunctions() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🔧 FUNCTIONS - Declaration, Parameters, Returns")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// فراخوانی توابع مختلف
	// ============================================

	fmt.Println("\n--- 4.1 Simple Function ---")
	fmt.Printf("  add(5, 3) = %d\n", add(5, 3))

	fmt.Println("\n--- 4.2 Same Type Parameters ---")
	fmt.Printf("  multiply(4, 6) = %d\n", multiply(4, 6))

	fmt.Println("\n--- 4.3 Multiple Return Values ---")
	q, r := divideWithRemainder(17, 5)
	fmt.Printf("  divideWithRemainder(17, 5) = quotient=%d, remainder=%d\n", q, r)

	fmt.Println("\n--- 4.4 Named Return Values ---")
	area, perimeter := rectangleAreaAndPerimeter(5, 3)
	fmt.Printf("  Rectangle: width=5, height=3 → area=%.1f, perimeter=%.1f\n", area, perimeter)

	fmt.Println("\n--- 4.5 Variadic Function ---")
	fmt.Printf("  sum(1,2,3,4,5) = %d\n", sum(1, 2, 3, 4, 5))
	fmt.Printf("  sum(10,20) = %d\n", sum(10, 20))

	fmt.Println("\n--- 4.6 Recursive Function ---")
	fmt.Printf("  factorial(5) = %d\n", factorial(5))

	fmt.Println("\n--- 4.7 Function as Value ---")
	double := func(a, b int) int { return a * b }
	result := applyOperation(6, 7, double)
	fmt.Printf("  applyOperation(6, 7, double) = %d\n", result)

	fmt.Println("\n--- 4.8 Closure ---")
	c := counter()
	fmt.Printf("  counter(): %d\n", c())
	fmt.Printf("  counter(): %d\n", c())
	fmt.Printf("  counter(): %d\n", c())

	fmt.Println("\n--- 4.9 Defer Example ---")
	demonstrateDefer()

	fmt.Println("\n--- 4.10 Panic/Recover Example ---")
	if res, err := safeDivide(10, 2); err != nil {
		fmt.Printf("  safeDivide(10, 2): error=%v\n", err)
	} else {
		fmt.Printf("  safeDivide(10, 2): %d\n", res)
	}

	if res, err := safeDivide(10, 0); err != nil {
		fmt.Printf("  safeDivide(10, 0): error=%v\n", err)
	} else {
		fmt.Printf("  safeDivide(10, 0): %d\n", res)
	}
}

// ============================================================================
// بخش 5: توابع پیشرفته (Function Advanced Patterns)
// ============================================================================

// 5.1 Function Options Pattern (الگوی تنظیمات اختیاری)
type ServerConfig struct {
	Host    string
	Port    int
	Timeout int
}

type ServerOption func(*ServerConfig)

func WithHost(host string) ServerOption {
	return func(c *ServerConfig) {
		c.Host = host
	}
}

func WithPort(port int) ServerOption {
	return func(c *ServerConfig) {
		c.Port = port
	}
}

func WithTimeout(timeout int) ServerOption {
	return func(c *ServerConfig) {
		c.Timeout = timeout
	}
}

func NewServer(options ...ServerOption) *ServerConfig {
	// مقادیر پیش‌فرض
	config := &ServerConfig{
		Host:    "localhost",
		Port:    8080,
		Timeout: 30,
	}

	// اعمال options
	for _, option := range options {
		option(config)
	}

	return config
}

// 5.2 Generic Function (از Go 1.18+)
func genericMax[T int | float64](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// 5.3 تابع با نوع‌های مختلف با interface{}
func printAnything(v interface{}) {
	fmt.Printf("  Type: %T, Value: %v\n", v, v)
}

func demonstrateAdvancedFunctions() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("⚡ ADVANCED FUNCTIONS - Options Pattern, Generics")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println("\n--- 5.1 Function Options Pattern ---")
	server1 := NewServer()
	server2 := NewServer(WithHost("0.0.0.0"), WithPort(3000))
	server3 := NewServer(WithTimeout(60), WithPort(9000))

	fmt.Printf("  Default: %+v\n", server1)
	fmt.Printf("  Custom host+port: %+v\n", server2)
	fmt.Printf("  Custom timeout+port: %+v\n", server3)

	fmt.Println("\n--- 5.2 Generic Function (Go 1.18+) ---")
	fmt.Printf("  genericMax(10, 20): %d\n", genericMax(10, 20))
	fmt.Printf("  genericMax(3.14, 2.71): %.2f\n", genericMax(3.14, 2.71))

	fmt.Println("\n--- 5.3 Print Anything with interface{} ---")
	printAnything(42)
	printAnything("hello")
	printAnything(3.14)
	printAnything(true)
}

// ============================================================================
// بخش 6: نکات مهم و اشتباهات رایج
// ============================================================================

func demonstrateCommonMistakes() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("⚠️ COMMON MISTAKES AND BEST PRACTICES")
	fmt.Println(stringsRepeat("=", 80))

	// اشتباه 1: استفاده از := در خارج از تابع
	fmt.Println("\n❌ Mistake 1: Using := outside function")
	fmt.Println("   // This is invalid:")
	fmt.Println("   // x := 10  // syntax error")
	fmt.Println("   ✅ Use var x = 10 instead")

	// اشتباه 2: اعلان متغیر و استفاده نکردن
	fmt.Println("\n❌ Mistake 2: Declared but not used")
	fmt.Println("   var unused int = 5  // compilation error")
	fmt.Println("   ✅ Use _ = unused to ignore, or remove it")

	// اشتباه 3: فراموش کردن نوع در short declaration
	fmt.Println("\n❌ Mistake 3: Short declaration with existing variable")
	fmt.Println("   x := 10")
	fmt.Println("   x := 20  // error: no new variables")
	fmt.Println("   ✅ Use x = 20 for assignment")

	// اشتباه 4: استفاده از break در switch بدون label
	fmt.Println("\n❌ Mistake 4: Break in switch doesn't break loop")
	fmt.Println("   for i := 0; i < 5; i++ {")
	fmt.Println("       switch { case i == 2: break } // breaks switch only")
	fmt.Println("   }")
	fmt.Println("   ✅ Use labeled break: break loopLabel")

	// اشتباه 5: shadowing متغیرها
	fmt.Println("\n❌ Mistake 5: Variable shadowing")
	fmt.Println("   x := 10")
	fmt.Println("   if true {")
	fmt.Println("       x := 20  // new variable, outer x unchanged")
	fmt.Println("   }")
	fmt.Println("   // x is still 10")
}

// ============================================================================
// بخش 7: جمع‌بندی و جدول مرجع سریع
// ============================================================================

func main() {
	fmt.Println(stringsRepeat("=", 80))
	fmt.Println("🎯 COMPLETE GO SYNTAX GUIDE")
	fmt.Println("Variables | Loops | Conditionals | Functions")
	fmt.Println(stringsRepeat("=", 80))

	// بخش 1: متغیرها
	demonstrateVariables()

	// بخش 2: حلقه‌ها
	demonstrateLoops()

	// بخش 3: شرط‌ها
	demonstrateConditionals()

	// بخش 4: توابع
	demonstrateFunctions()

	// بخش 5: توابع پیشرفته
	demonstrateAdvancedFunctions()

	// بخش 6: اشتباهات رایج
	demonstrateCommonMistakes()

	// بخش 7: جمع‌بندی
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📚 QUICK REFERENCE CARD")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ VARIABLES                                                       │")
	fmt.Println("├─────────────────────────────────────────────────────────────────┤")
	fmt.Println("│ var name string                    // declaration              │")
	fmt.Println("│ var age int = 30                   // declaration + init       │")
	fmt.Println("│ city := \"Tehran\"                    // short declaration        │")
	fmt.Println("│ var ( a int; b string )            // block declaration        │")
	fmt.Println("│ const pi = 3.14                    // constant                 │")
	fmt.Println("│ const ( A = iota; B )              // iota counter             │")
	fmt.Println("└─────────────────────────────────────────────────────────────────┘")

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ LOOPS (only 'for')                                              │")
	fmt.Println("├─────────────────────────────────────────────────────────────────┤")
	fmt.Println("│ for i := 0; i < 10; i++           // C-style loop              │")
	fmt.Println("│ for x < 10                        // while-style               │")
	fmt.Println("│ for { break }                     // infinite loop             │")
	fmt.Println("│ for i, v := range slice           // range over slice          │")
	fmt.Println("│ for k, v := range map             // range over map            │")
	fmt.Println("│ for i, r := range \"text\"          // range over string (rune)  │")
	fmt.Println("└─────────────────────────────────────────────────────────────────┘")

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ CONDITIONALS                                                    │")
	fmt.Println("├─────────────────────────────────────────────────────────────────┤")
	fmt.Println("│ if x > 0 { ... } else if x < 0 { ... } else { ... }            │")
	fmt.Println("│ if err := do(); err != nil { ... } // with statement            │")
	fmt.Println("│ switch x { case 1: ... default: ... }                          │")
	fmt.Println("│ switch { case x > 0: ... }        // tagless switch            │")
	fmt.Println("│ switch x { case 1,2: ... }        // multiple values           │")
	fmt.Println("└─────────────────────────────────────────────────────────────────┘")

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ FUNCTIONS                                                       │")
	fmt.Println("├─────────────────────────────────────────────────────────────────┤")
	fmt.Println("│ func add(a, b int) int { return a+b }                          │")
	fmt.Println("│ func div(a, b int) (int, error) { ... } // multiple returns    │")
	fmt.Println("│ func area(w, h float64) (a float64, p float64) { ... }         │")
	fmt.Println("│ func sum(nums ...int) int { ... }    // variadic               │")
	fmt.Println("│ func() { ... }()                    // anonymous + call        │")
	fmt.Println("│ defer fmt.Println(\"done\")            // defer execution        │")
	fmt.Println("└─────────────────────────────────────────────────────────────────┘")

	fmt.Println("\n💡 GOLDEN RULES:")
	fmt.Println("  1. Use short declaration (:=) inside functions")
	fmt.Println("  2. Use var for package-level variables")
	fmt.Println("  3. Go has only 'for' loop (no while, no do-while)")
	fmt.Println("  4. No ternary operator (?:) in Go")
	fmt.Println("  5. Use switch instead of long if-else chains")
	fmt.Println("  6. Functions can return multiple values")
	fmt.Println("  7. Use named returns for documentation")
	fmt.Println("  8. Defer runs LIFO (last-in-first-out)")
	fmt.Println("  9. All variables must be used (no unused variables)")
	fmt.Println("  10. Type conversions are explicit (no implicit conversion)")
}

// ============================================================================
// بخش 8: توابع کمکی
// ============================================================================

func stringsRepeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}
