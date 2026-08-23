// ============================================================================
// FILE: fmt_complete_guide.go
// TITLE: راهنمای کامل پکیج fmt در Go - ورودی/خروجی فرمت شده
// HOW TO RUN: go run fmt_complete_guide.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - پکیج fmt چیست و چرا مهم است؟
// ============================================================================
//
// پکیج fmt امکانات کامل برای:
// 1. چاپ فرمت شده (Printf, Println, Print)
// 2. خواندن فرمت شده (Scanf, Scanln, Scan)
// 3. تبدیل به رشته (Sprintf, Sprint, Sprintln)
// 4. نوشتن در رشته (Fprintf, Fprint, Fprintln)
// 5. ورودی از رشته (Sscanf, Sscan, Sscanln)
//
// قانون طلایی:
// "از %v برای چاپ مقدار پیش‌فرض هر نوع، از %+v برای نمایش فیلدهای struct،
//  از %#v برای نمایش literal Go، و از %T برای نمایش نوع استفاده کن."
// ============================================================================

package __internal_packages

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

// ============================================================================
// بخش 1: انواع داده برای تست
// ============================================================================

// انواع مختلف برای نمایش
type Person struct {
	Name string
	Age  int
	City string
}

// پیاده‌سازی Stringer interface
func (p Person) String() string {
	return fmt.Sprintf("%s (%d) from %s", p.Name, p.Age, p.City)
}

// نوع سفارشی
type Color string

const (
	Red   Color = "red"
	Blue  Color = "blue"
	Green Color = "green"
)

// ============================================================================
// بخش 2: چاپ به stdout (Print, Println, Printf)
// ============================================================================

func demonstratePrint() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🖨️ PRINT TO STDOUT - Print, Println, Printf")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 2.1 Print - چاپ بدون خط جدید
	// ============================================
	fmt.Println("\n--- 2.1 Print (no newline) ---")

	fmt.Print("Hello")
	fmt.Print(" ")
	fmt.Print("World")
	fmt.Print("\n") // خط جدید دستی

	// Print با چند آرگومان
	fmt.Print("Value1:", 42, ", Value2:", 3.14, "\n")

	// ============================================
	// 2.2 Println - چاپ با خط جدید
	// ============================================
	fmt.Println("\n--- 2.2 Println (with newline) ---")

	fmt.Println("Hello World")
	fmt.Println(42, 3.14, true, "text")

	// بین آرگومان‌ها space اضافه می‌کند
	fmt.Println("First", "Second", "Third")

	// ============================================
	// 2.3 Printf - چاپ فرمت شده
	// ============================================
	fmt.Println("\n--- 2.3 Printf (formatted) ---")

	// Verbs عمومی
	name := "Ali"
	age := 30
	height := 1.85
	isActive := true

	fmt.Printf("Name: %s, Age: %d, Height: %.2f, Active: %t\n",
		name, age, height, isActive)

	// %v (default format)
	fmt.Printf("Default: %v, %v, %v, %v\n", name, age, height, isActive)

	// %+v (با نام فیلدها برای struct)
	person := Person{Name: "Sara", Age: 28, City: "Tehran"}
	fmt.Printf("Struct with %%+v: %+v\n", person)

	// %#v (Go syntax representation)
	fmt.Printf("Go syntax: %#v\n", person)

	// %T (نوع)
	fmt.Printf("Types: %T, %T, %T\n", name, age, person)

	// عرض و دقت
	fmt.Printf("Width 10: |%10s|, |%10d|\n", name, age)
	fmt.Printf("Width -10 (left align): |%-10s|, |%-10d|\n", name, age)
	fmt.Printf("Precision: %.2f, %.4f\n", height, 3.14159265)

	// اعداد با پایه‌های مختلف
	num := 255
	fmt.Printf("Decimal: %d, Binary: %b, Octal: %o, Hex: %x, Hex upper: %X\n",
		num, num, num, num, num)

	// pointers
	ptr := &age
	fmt.Printf("Pointer: %p, Value: %d\n", ptr, *ptr)

	// Unicode و rune
	runeVal := '世'
	fmt.Printf("Rune: %c, Unicode: %U, Value: %d\n", runeVal, runeVal, runeVal)
}

// ============================================================================
// بخش 3: Stringer اینترفیس (چاپ سفارشی)
// ============================================================================

func demonstrateStringer() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🎨 STRINGER INTERFACE - Custom String Representation")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 3.1 استفاده از Stringer
	// ============================================
	fmt.Println("\n--- 3.1 Using Stringer ---")

	p := Person{Name: "Ali", Age: 30, City: "Tehran"}

	// بدون Stringer: {Ali 30 Tehran}
	// با Stringer: Ali (30) from Tehran
	fmt.Printf("Person: %v\n", p)
	fmt.Printf("Person with %%s: %s\n", p)
	fmt.Printf("Person with %%v: %v\n", p)

	// ============================================
	// 3.2 Stringer برای انواع سفارشی
	// ============================================
	fmt.Println("\n--- 3.2 Stringer for Custom Types ---")

	// تعریف نوع جدید با Stringer
	type Temperature float64

	func (t Temperature) String() string {
		return fmt.Sprintf("%.1f°C", t)
	}

	temp := Temperature(23.5)
	fmt.Printf("Temperature: %v\n", temp)
	fmt.Printf("Temperature with fmt: %s\n", temp)

	// نوع Color
	func (c Color) String() string {
		switch c {
	case Red:
		return "🔴 Red"
	case Blue:
		return "🔵 Blue"
	case Green:
		return "🟢 Green"
	default:
		return "⚪ Unknown"
	}
	}

	fmt.Printf("Color: %v\n", Red)
	fmt.Printf("Color: %v\n", Blue)
	fmt.Printf("Color: %v\n", Green)
}

// ============================================================================
// بخش 4: Formatter اینترفیس (کنترل پیشرفته‌تر)
// ============================================================================

type CustomInt int

// پیاده‌سازی Formatter برای کنترل دقیق فرمت
func (c CustomInt) Format(f fmt.State, verb rune) {
	switch verb {
	case 'v':
		if f.Flag('+') {
			fmt.Fprintf(f, "CustomInt(%d)", int(c))
			return
		}
		fallthrough
	case 'd':
		fmt.Fprintf(f, "%d", int(c))
	case 'x':
		fmt.Fprintf(f, "%x", int(c))
	case 'b':
		fmt.Fprintf(f, "%b", int(c))
	default:
		fmt.Fprintf(f, "%%!%c(%T=%d)", verb, c, int(c))
	}
}

func demonstrateFormatter() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🔧 FORMATTER INTERFACE - Advanced Format Control")
	fmt.Println(stringsRepeat("=", 80))

	num := CustomInt(42)

	fmt.Printf("Default: %v\n", num)
	fmt.Printf("With + flag: %+v\n", num)
	fmt.Printf("Decimal: %d\n", num)
	fmt.Printf("Hex: %x\n", num)
	fmt.Printf("Binary: %b\n", num)
}

// ============================================================================
// بخش 5: نوشتن در رشته (Sprint, Sprintln, Sprintf)
// ============================================================================

func demonstrateSprint() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📝 SPRINT - Print to String")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 5.1 Sprint - بدون خط جدید
	// ============================================
	fmt.Println("\n--- 5.1 Sprint (no newline) ---")

	str1 := fmt.Sprint("Hello", " ", "World")
	fmt.Printf("Sprint: %q\n", str1)

	str2 := fmt.Sprint(42, 3.14, true)
	fmt.Printf("Sprint multiple: %q\n", str2)

	// ============================================
	// 5.2 Sprintln - با خط جدید
	// ============================================
	fmt.Println("\n--- 5.2 Sprintln (with newline) ---")

	str3 := fmt.Sprintln("Hello", "World")
	fmt.Printf("Sprintln: %q\n", str3)

	// ============================================
	// 5.3 Sprintf - فرمت شده
	// ============================================
	fmt.Println("\n--- 5.3 Sprintf (formatted) ---")

	name := "Ali"
	age := 30
	str4 := fmt.Sprintf("Name: %s, Age: %d", name, age)
	fmt.Printf("Sprintf: %s\n", str4)

	// کاربرد: ساخت پیام خطا
	errMsg := fmt.Sprintf("Error: user %s not found", name)
	fmt.Printf("Error message: %s\n", errMsg)

	// کاربرد: ساخت query
	query := fmt.Sprintf("SELECT * FROM users WHERE name='%s' AND age=%d", name, age)
	fmt.Printf("SQL Query: %s\n", query)
}

// ============================================================================
// بخش 6: نوشتن در Writer (Fprint, Fprintln, Fprintf)
// ============================================================================

func demonstrateFprint() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("✍️ FPRINT - Print to io.Writer")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 6.1 نوشتن در فایل
	// ============================================
	fmt.Println("\n--- 6.1 Writing to File ---")

	// ایجاد فایل موقت
	file, err := os.CreateTemp("", "fmt_example_*.txt")
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer os.Remove(file.Name())
	defer file.Close()

	// نوشتن در فایل
	fmt.Fprint(file, "Hello from Fprint\n")
	fmt.Fprintln(file, "Hello from Fprintln")
	fmt.Fprintf(file, "Hello from Fprintf: %s\n", "formatted")

	// خواندن و نمایش محتوا
	file.Seek(0, 0)
	data := make([]byte, 1024)
	n, _ := file.Read(data)
	fmt.Printf("File content:\n%s", data[:n])

	// ============================================
	// 6.2 نوشتن در bytes.Buffer
	// ============================================
	fmt.Println("\n--- 6.2 Writing to bytes.Buffer ---")

	var buf bytes.Buffer

	fmt.Fprint(&buf, "Line 1\n")
	fmt.Fprintf(&buf, "Line %d\n", 2)
	fmt.Fprintln(&buf, "Line 3")

	fmt.Printf("Buffer content:\n%s", buf.String())

	// ============================================
	// 6.3 نوشتن در strings.Builder
	// ============================================
	fmt.Println("\n--- 6.3 Writing to strings.Builder ---")

	var builder strings.Builder

	fmt.Fprint(&builder, "Building ")
	fmt.Fprintf(&builder, "a %s ", "string")
	fmt.Fprintln(&builder, "efficiently")

	fmt.Printf("Builder result: %s\n", builder.String())

	// کارایی: pre-allocation
	builder.Grow(1024)
	fmt.Fprint(&builder, "Pre-allocated builder")
}

// ============================================================================
// بخش 7: خواندن از stdin (Scan, Scanln, Scanf)
// ============================================================================

func demonstrateScan() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📖 SCAN - Reading from stdin")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 7.1 مثال (غیرتعاملی - نمایش نحوه استفاده)
	// ============================================
	fmt.Println("\n--- 7.1 Scan Examples (non-interactive) ---")

	// این بخش فقط نحوه استفاده را نشان می‌دهد،
	// در اجرای واقعی نیاز به ورودی کاربر دارد

	fmt.Println("  // Example 1: Reading a single value")
	fmt.Println("  var name string")
	fmt.Println("  fmt.Print(\"Enter name: \")")
	fmt.Println("  fmt.Scan(&name)")

	fmt.Println("\n  // Example 2: Reading multiple values")
	fmt.Println("  var age int")
	fmt.Println("  var height float64")
	fmt.Println("  fmt.Scan(&age, &height)")

	fmt.Println("\n  // Example 3: Formatted reading")
	fmt.Println("  var first, last string")
	fmt.Println("  var year int")
	fmt.Println("  fmt.Scanf(\"%s %s %d\", &first, &last, &year)")

	// ============================================
	// 7.2 Scan با رشته شبیه‌سازی شده
	// ============================================
	fmt.Println("\n--- 7.2 Simulated Scan with String ---")

	// شبیه‌سازی ورودی
	input := "Ali 30 1.85"
	reader := strings.NewReader(input)

	var simName string
	var simAge int
	var simHeight float64

	fmt.Fscan(reader, &simName, &simAge, &simHeight)
	fmt.Printf("  Scanned: name=%s, age=%d, height=%.2f\n",
		simName, simAge, simHeight)
}

// ============================================================================
// بخش 8: خواندن از رشته (Sscan, Sscanln, Sscanf)
// ============================================================================

func demonstrateSscan() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🔍 SSCAN - Reading from String")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 8.1 Sscan - خواندن ساده
	// ============================================
	fmt.Println("\n--- 8.1 Sscan (simple) ---")

	data := "Alice 25 1.68"
	var name string
	var age int
	var height float64

	n, err := fmt.Sscan(data, &name, &age, &height)
	if err == nil {
		fmt.Printf("  Read %d items: name=%s, age=%d, height=%.2f\n",
			n, name, age, height)
	}

	// ============================================
	// 8.2 Sscanln - تا خط جدید
	// ============================================
	fmt.Println("\n--- 8.2 Sscanln (until newline) ---")

	multiLine := "Bob 30\nCharlie 25\n"

	var name2 string
	var age2 int

	// فقط خط اول را می‌خواند
	n, _ = fmt.Sscanln(multiLine, &name2, &age2)
	fmt.Printf("  Sscanln read: name=%s, age=%d\n", name2, age2)

	// ============================================
	// 8.3 Sscanf - فرمت شده
	// ============================================
	fmt.Println("\n--- 8.3 Sscanf (formatted) ---")

	logLine := "[ERROR] 2024-01-15: Connection failed"
	var level, date, message string

	n, err = fmt.Sscanf(logLine, "[%s] %s: %s", &level, &date, &message)
	if err == nil {
		fmt.Printf("  Parsed %d items: level=%s, date=%s, message=%s\n",
			n, level, date, message)
	}

	// مثال با اعداد
	coordinates := "Point: (10, 20)"
	var x, y int
	fmt.Sscanf(coordinates, "Point: (%d, %d)", &x, &y)
	fmt.Printf("  Coordinates: x=%d, y=%d\n", x, y)
}

// ============================================================================
// بخش 9: Errorf - ساخت خطاهای فرمت شده
// ============================================================================

func demonstrateErrorf() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("⚠️ ERRORF - Creating Formatted Errors")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 9.1 ایجاد خطا با Errorf
	// ============================================
	fmt.Println("\n--- 9.1 Creating Errors ---")

	userID := 123
	err1 := fmt.Errorf("user %d not found", userID)
	fmt.Printf("  Error: %v\n", err1)

	// خطا با چند پارامتر
	err2 := fmt.Errorf("failed to process: user=%s, action=%s, code=%d",
		"ali", "login", 401)
	fmt.Printf("  Error: %v\n", err2)

	// ============================================
	// 9.2 Wrapping errors (با %w)
	// ============================================
	fmt.Println("\n--- 9.2 Wrapping Errors ---")

	baseErr := fmt.Errorf("connection refused")
	wrappedErr := fmt.Errorf("database error: %w", baseErr)

	fmt.Printf("  Base error: %v\n", baseErr)
	fmt.Printf("  Wrapped error: %v\n", wrappedErr)

	// بررسی با errors.Is (نیاز به import errors)
	// if errors.Is(wrappedErr, baseErr) {
	//     fmt.Println("  Error is wrapped")
	// }
}

// ============================================================================
// بخش 10: GoStringer اینترفیس (برای %#v)
// ============================================================================

type Product struct {
	ID    int
	Name  string
	Price float64
}

// GoStringer برای نمایش Go syntax
func (p Product) GoString() string {
	return fmt.Sprintf("&Product{ID:%d, Name:%q, Price:%.2f}",
		p.ID, p.Name, p.Price)
}

func demonstrateGoStringer() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🔧 GoStringer INTERFACE - For %#v verb")
	fmt.Println(stringsRepeat("=", 80))

	p := Product{ID: 1, Name: "Laptop", Price: 999.99}

	fmt.Printf("Default %%v: %v\n", p)
	fmt.Printf("Default %%+v: %+v\n", p)
	fmt.Printf("Go syntax %%#v: %#v\n", p)
}

// ============================================================================
// بخش 11: State و Flag در فرمت‌دهی پیشرفته
// ============================================================================

type HexColor int

func (c HexColor) Format(f fmt.State, verb rune) {
	switch verb {
	case 'v', 'x', 'X':
		if f.Flag('#') {
			fmt.Fprintf(f, "#%06X", int(c))
		} else {
			fmt.Fprintf(f, "%X", int(c))
		}
	default:
		fmt.Fprintf(f, "%%!%c(%T=%d)", verb, c, int(c))
	}
}

func demonstrateFormatFlags() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🎨 FORMAT FLAGS - Advanced Formatting")
	fmt.Println(stringsRepeat("=", 80))

	color := HexColor(0xFF5733)

	fmt.Printf("Default: %v\n", color)
	fmt.Printf("With # flag: %#v\n", color)
	fmt.Printf("Hex: %x\n", color)
	fmt.Printf("Hex upper: %X\n", color)

	// فلگ‌های مختلف
	num := 123
	fmt.Printf("Default: |%v|\n", num)
	fmt.Printf("Width 10: |%10v|\n", num)
	fmt.Printf("Width -10 (left): |%-10v|\n", num)
	fmt.Printf("Zero padding: |%010v|\n", num)
	fmt.Printf("Space for positive: |% d|\n", num)
	fmt.Printf("Sign for positive: |%+d|\n", num)
	fmt.Printf("Alternate format: |%#o|, |%#x|\n", num, num)
}

// ============================================================================
// بخش 12: کاربردهای عملی و الگوها
// ============================================================================

func demonstratePracticalUses() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("💡 PRACTICAL USES AND PATTERNS")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 12.1 لاگینگ فرمت شده
	// ============================================
	fmt.Println("\n--- 12.1 Formatted Logging ---")

	logMessage := func(level, msg string, args ...interface{}) {
		prefix := fmt.Sprintf("[%s] ", level)
		format := prefix + msg + "\n"
		fmt.Printf(format, args...)
	}

	logMessage("INFO", "User %s logged in", "Ali")
	logMessage("ERROR", "Failed to connect to %s: %v", "database", "timeout")
	logMessage("DEBUG", "Value: %d, Status: %s", 42, "ok")

	// ============================================
	// 12.2 جدول‌سازی (Table formatting)
	// ============================================
	fmt.Println("\n--- 12.2 Table Formatting ---")

	type User struct {
		Name  string
		Email string
		Age   int
	}

	users := []User{
		{"Ali", "ali@test.com", 30},
		{"Sara Mohammadi", "sara@test.com", 25},
		{"Reza", "reza@test.com", 35},
	}

	// هدر جدول
	fmt.Printf("%-20s | %-25s | %5s\n", "Name", "Email", "Age")
	fmt.Println(stringsRepeat("-", 20+2+25+2+5))

	// ردیف‌ها
	for _, u := range users {
		fmt.Printf("%-20s | %-25s | %5d\n", u.Name, u.Email, u.Age)
	}

	// ============================================
	// 12.3 نمایش پیشرفت (Progress display)
	// ============================================
	fmt.Println("\n--- 12.3 Progress Display ---")

	for i := 0; i <= 100; i += 10 {
		fmt.Printf("\rProgress: %3d%% [%s]", i, stringsRepeat("=", i/2))
		// time.Sleep(100 * time.Millisecond) // در اجرای واقعی
	}
	fmt.Println() // خط جدید

	// ============================================
	// 12.4 ساخت query parameters
	// ============================================
	fmt.Println("\n--- 12.4 Building Query Parameters ---")

	buildQuery := func(params map[string]interface{}) string {
		var parts []string
		for k, v := range params {
			parts = append(parts, fmt.Sprintf("%s=%v", k, v))
		}
		return strings.Join(parts, "&")
	}

	query := buildQuery(map[string]interface{}{
		"q":     "golang",
		"page":  1,
		"limit": 10,
	})
	fmt.Printf("Query string: %s\n", query)
}

// ============================================================================
// بخش 13: اشتباهات رایج
// ============================================================================

func demonstrateCommonMistakes() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("❌ COMMON MISTAKES WITH fmt")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println("\n❌ Mistake 1: Wrong verb for type")
	fmt.Println("   var num int = 42")
	fmt.Println("   fmt.Printf(\"%s\", num)  // %!s(int=42)")
	fmt.Println("   ✅ Use %d for integers")

	fmt.Println("\n❌ Mistake 2: Not enough arguments")
	fmt.Println("   fmt.Printf(\"%s %s\", \"hello\")  // missing argument")
	fmt.Println("   ✅ Match number of verbs with arguments")

	fmt.Println("\n❌ Mistake 3: Too many arguments")
	fmt.Println("   fmt.Printf(\"%s\", \"hello\", \"world\")  // extra arg ignored")
	fmt.Println("   ✅ Use correct number of arguments")

	fmt.Println("\n❌ Mistake 4: Forgetting newline in Println")
	fmt.Println("   fmt.Print(\"Hello\")  // no newline")
	fmt.Println("   ✅ Use fmt.Println() or add \\n")

	fmt.Println("\n❌ Mistake 5: Not checking scan errors")
	fmt.Println("   var x int")
	fmt.Println("   fmt.Scan(&x)  // ignores error on invalid input")
	fmt.Println("   ✅ Check returned error")

	fmt.Println("\n❌ Mistake 6: Using %v for everything")
	fmt.Println("   // %v is convenient but sometimes too verbose")
	fmt.Println("   ✅ Use specific verbs for specific needs")

	fmt.Println("\n❌ Mistake 7: Not handling pointer formatting")
	fmt.Println("   p := &Person{Name:\"Ali\"}")
	fmt.Println("   fmt.Println(p)  // prints &{Ali 0 }")
	fmt.Println("   ✅ Use *p or implement Stringer")
}

// ============================================================================
// بخش 14: جمع‌بندی و جدول مرجع
// ============================================================================

func main() {
	fmt.Println(stringsRepeat("=", 80))
	fmt.Println("🎯 COMPLETE fmt PACKAGE GUIDE IN GO")
	fmt.Println(stringsRepeat("=", 80))

	// بخش 1: Print functions
	demonstratePrint()

	// بخش 2: Stringer
	demonstrateStringer()

	// بخش 3: Formatter
	demonstrateFormatter()

	// بخش 4: Sprint functions
	demonstrateSprint()

	// بخش 5: Fprint functions
	demonstrateFprint()

	// بخش 6: Scan functions
	demonstrateScan()

	// بخش 7: Sscan functions
	demonstrateSscan()

	// بخش 8: Errorf
	demonstrateErrorf()

	// بخش 9: GoStringer
	demonstrateGoStringer()

	// بخش 10: Format flags
	demonstrateFormatFlags()

	// بخش 11: Practical uses
	demonstratePracticalUses()

	// بخش 12: Common mistakes
	demonstrateCommonMistakes()

	// بخش 13: Quick reference
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📚 fmt PACKAGE QUICK REFERENCE")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ PRINT FUNCTIONS                                                │")
	fmt.Println("├─────────────────────────────────────────────────────────────────┤")
	fmt.Println("│ fmt.Print(v...)        - Print without newline                 │")
	fmt.Println("│ fmt.Println(v...)      - Print with newline and spaces         │")
	fmt.Println("│ fmt.Printf(format, v...) - Formatted print                     │")
	fmt.Println("│ fmt.Sprint(v...)       - Return string without newline         │")
	fmt.Println("│ fmt.Sprintln(v...)     - Return string with newline            │")
	fmt.Println("│ fmt.Sprintf(format, v...) - Return formatted string            │")
	fmt.Println("│ fmt.Fprint(w, v...)    - Print to io.Writer                    │")
	fmt.Println("│ fmt.Fprintln(w, v...)  - Print to io.Writer with newline       │")
	fmt.Println("│ fmt.Fprintf(w, format, v...) - Formatted print to io.Writer    │")
	fmt.Println("└─────────────────────────────────────────────────────────────────┘")

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ SCAN FUNCTIONS                                                 │")
	fmt.Println("├─────────────────────────────────────────────────────────────────┤")
	fmt.Println("│ fmt.Scan(a...)          - Read from stdin                      │")
	fmt.Println("│ fmt.Scanln(a...)        - Read until newline                   │")
	fmt.Println("│ fmt.Scanf(format, a...) - Formatted read from stdin            │")
	fmt.Println("│ fmt.Sscan(str, a...)    - Read from string                     │")
	fmt.Println("│ fmt.Sscanln(str, a...)  - Read from string until newline       │")
	fmt.Println("│ fmt.Sscanf(str, format, a...) - Formatted read from string     │")
	fmt.Println("│ fmt.Fscan(r, a...)      - Read from io.Reader                  │")
	fmt.Println("│ fmt.Fscanln(r, a...)    - Read from io.Reader until newline    │")
	fmt.Println("│ fmt.Fscanf(r, format, a...) - Formatted read from io.Reader    │")
	fmt.Println("└─────────────────────────────────────────────────────────────────┘")

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ VERBS (Format Verbs)                                           │")
	fmt.Println("├─────────────────────────────────────────────────────────────────┤")
	fmt.Println("│ %v      - Default format                                       │")
	fmt.Println("│ %+v     - Default with field names (for structs)               │")
	fmt.Println("│ %#v     - Go syntax representation                             │")
	fmt.Println("│ %T      - Type of value                                        │")
	fmt.Println("│ %%      - Literal percent sign                                 │")
	fmt.Println("│ %t      - Boolean (true or false)                              │")
	fmt.Println("│ %d      - Decimal integer                                      │")
	fmt.Println("│ %b      - Binary integer                                       │")
	fmt.Println("│ %o      - Octal integer                                        │")
	fmt.Println("│ %x/%X   - Hexadecimal integer (lower/upper)                    │")
	fmt.Println("│ %c      - Character (rune)                                     │")
	fmt.Println("│ %q      - Quoted string or char                                │")
	fmt.Println("│ %s      - String                                               │")
	fmt.Println("│ %f/%F   - Floating point                                       │")
	fmt.Println("│ %e/%E   - Scientific notation                                  │")
	fmt.Println("│ %g/%G   - Compact floating point                               │")
	fmt.Println("│ %p      - Pointer address                                      │")
	fmt.Println("└─────────────────────────────────────────────────────────────────┘")

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ FLAGS                                                          │")
	fmt.Println("├─────────────────────────────────────────────────────────────────┤")
	fmt.Println("│ +     - Always print sign for numeric values                   │")
	fmt.Println("│ -     - Left-justify (default right)                           │")
	fmt.Println("│ #     - Alternate format (e.g., 0x for hex)                    │")
	fmt.Println("│ ' '   - Leave space for positive numbers                       │")
	fmt.Println("│ 0     - Pad with zeros instead of spaces                       │")
	fmt.Println("│ width - Minimum width (e.g., %10s)                             │")
	fmt.Println("│ .prec - Precision (e.g., %.2f)                                 │")
	fmt.Println("└─────────────────────────────────────────────────────────────────┘")

	fmt.Println("\n💡 GOLDEN RULES:")
	fmt.Println("  1. Use %v for quick debugging, %+v for structs")
	fmt.Println("  2. Implement Stringer for custom types")
	fmt.Println("  3. Use Errorf for creating formatted errors")
	fmt.Println("  4. Always check errors from Scan functions")
	fmt.Println("  5. Use Fprint for logging to files")
	fmt.Println("  6. Use Sprint for building strings without printing")
	fmt.Println("  7. Match verb count with argument count")
	fmt.Println("  8. Use specific verbs for specific types")
	fmt.Println("  9. Use width and precision for aligned output")
	fmt.Println("  10. Prefer Println over Print when you need newline")

	fmt.Println("\n🎯 PERFORMANCE TIPS:")
	fmt.Println("  • Sprint is faster than string concatenation")
	fmt.Println("  • strings.Builder is faster than Sprint for many concatenations")
	fmt.Println("  • Pre-allocate buffers when possible")
	fmt.Println("  • Avoid Printf when Println/Print works")
	fmt.Println("  • Use Fprint with bytes.Buffer for efficient building")
}

// ============================================================================
// بخش 15: توابع کمکی
// ============================================================================

func stringsRepeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}
/*
# اجرای کامل برنامه
go run fmt_complete_guide.go

# برای تست Scan (نیاز به ورودی دارد)
echo "Ali 30" | go run fmt_complete_guide.go
 */