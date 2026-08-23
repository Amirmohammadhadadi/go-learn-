// ============================================================================
// FILE: go_tools_guide.go
// TITLE: راهنمای کامل ابزارهای خط فرمان Go - go fmt, go vet, go mod, go generate
// HOW TO RUN: go run go_tools_guide.go (این فایل فقط توضیحات است)
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - ابزارهای خط فرمان Go
// ============================================================================
//
// Go همراه با مجموعه‌ای از ابزارهای خط فرمان قدرتمند است:
//
// 1. go fmt   - فرمت کردن خودکار کد (indentation, spacing, etc.)
// 2. go vet   - آنالیز استاتیک کد (یافتن باگ‌های احتمالی)
// 3. go mod   - مدیریت ماژول‌ها و وابستگی‌ها
// 4. go generate - تولید خودکار کد از الگوها
//
// قانون طلایی:
// "قبل از commit، همیشه go fmt و go vet را اجرا کن.
//  از go mod برای مدیریت وابستگی‌ها استفاده کن.
//  از go generate برای تولید کدهای تکراری (مثل mock) استفاده کن."
// ============================================================================

package __go_tools_guide

import (
	"fmt"
	"strings"
)

// این فایل فقط یک فایل توضیحی است.
// تمام دستورات باید در ترمینال اجرا شوند.

func main() {
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("🎯 GO COMMAND LINE TOOLS GUIDE")
	fmt.Println("go fmt | go vet | go mod | go generate")
	fmt.Println(strings.Repeat("=", 80))

	// ============================================================================
	// بخش 1: go fmt - قالب‌بندی خودکار کد
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📝 SECTION 1: go fmt - Code Formatting")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────┐
│ go fmt COMMANDS                                                │
├─────────────────────────────────────────────────────────────────┤
│ go fmt              - فرمت فایل‌های текущего دایرکتوری        │
│ go fmt ./...        - فرمت همه فایل‌های پروژه (recursive)      │
│ go fmt -n ./...     - نمایش دستورات بدون اجرا (dry run)        │
│ go fmt -x ./...     - نمایش دستورات در حین اجرا                │
│ go fmt -l ./...     - نمایش نام فایل‌هایی که تغییر می‌کنند     │
│ go fmt -d ./...     - نمایش diff تغییرات (بدون اعمال)          │
└─────────────────────────────────────────────────────────────────┘

💡 RULES OF go fmt:
   • Indentation: tabs for indentation, spaces for alignment
   • No extra spaces: removes trailing whitespace
   • Consistent braces: always on same line
   • Imports: groups and sorts imports alphabetically
   • Line wrapping: no hard limit, but avoid long lines

📌 BEST PRACTICE:
   • Run go fmt before every commit
   • Use editor integration (save on format)
   • CI/CD should check: go fmt -l . | wc -l (should be 0)
`)

	// ============================================================================
	// بخش 2: go vet - آنالیز استاتیک
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🔍 SECTION 2: go vet - Static Analysis")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────┐
│ go vet COMMANDS                                                │
├─────────────────────────────────────────────────────────────────┤
│ go vet              - آنالیز текущего پکیج                     │
│ go vet ./...        - آنالیز همه پکیج‌ها                       │
│ go vet -v ./...     - با خروجی verbose                         │
│ go vet -tags test   - با build tags خاص                        │
└─────────────────────────────────────────────────────────────────┘

🔍 WHAT go vet DETECTS:

   1. Unreachable code (کدهای غیرقابل دسترسی)
      func example() {
          return
          fmt.Println("never called") // go vet detects this
      }

   2. Suspicious assignments
      var x int
      x = x // self-assignment

   3. Wrong format in printf
      fmt.Printf("%d", "string") // wrong type

   4. Unused struct fields
      type User struct {
          Name string
          age  int // unexported, unused
      }

   5. Nil function calls
      var fn func()
      fn() // nil pointer dereference

   6. Lock misuse
      var mu sync.Mutex
      mu.Unlock() // unlock without lock

   7. Suspicious shifts
      x := 1 << 100 // shift too large

   8. Unkeyed struct literals
      type Point struct { X, Y int }
      p := Point{10, 20} // should use Point{X:10, Y:20}

📌 BEST PRACTICE:
   • Run go vet in CI/CD pipeline
   • Fix all warnings (even if code "works")
   • Use golangci-lint for more checks
`)

	// ============================================================================
	// بخش 3: go mod - مدیریت ماژول
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📦 SECTION 3: go mod - Module Management")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────┐
│ INITIALIZATION                                                 │
├─────────────────────────────────────────────────────────────────┤
│ go mod init <module-name>   - ایجاد ماژول جدید                 │
│ go mod init example.com/myproject                              │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ DEPENDENCY MANAGEMENT                                          │
├─────────────────────────────────────────────────────────────────┤
│ go get <pkg>                    - افزودن وابستگی               │
│ go get github.com/gorilla/mux@v1.8.0                           │
│ go get -u ./...                 - به‌روزرسانی همه وابستگی‌ها    │
│ go get -u=patch ./...           - به‌روزرسانی فقط patch        │
│ go mod tidy                     - حذف وابستگی‌های بی‌استفاده    │
│ go mod vendor                   - کپی وابستگی‌ها در vendor/     │
│ go mod download                 - دانلود وابستگی‌ها             │
│ go mod verify                   - بررسی یکپارچگی وابستگی‌ها     │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ INSPECTION                                                     │
├─────────────────────────────────────────────────────────────────┤
│ go list -m all                 - لیست همه وابستگی‌ها            │
│ go list -m -u all              - لیست با نسخه‌های جدیدتر        │
│ go mod graph                   - نمایش درخت وابستگی             │
│ go mod why <pkg>               - دلیل وجود وابستگی              │
│ go mod edit -json              - نمایش go.mod به صورت JSON      │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ TROUBLESHOOTING                                                │
├─────────────────────────────────────────────────────────────────┤
│ go clean -modcache            - پاک کردن کش ماژول‌ها            │
│ go mod edit -replace <old>=<new> - جایگزینی ماژول              │
│ go mod edit -dropreplace <pkg>  - حذف جایگزینی                 │
│ go mod edit -exclude <pkg>@<ver> - حذف نسخه خاص                │
│ go mod edit -retract <ver>      -撤回 کردن نسخه                 │
└─────────────────────────────────────────────────────────────────┘

📂 FILES STRUCTURE:

   go.mod      - تعریف ماژول و وابستگی‌ها
   go.sum      - checksum وابستگی‌ها (برای یکپارچگی)
   vendor/     - کپی وابستگی‌ها (اختیاری)

📝 EXAMPLE go.mod:

   module github.com/myuser/myproject

   go 1.21

   require (
       github.com/gorilla/mux v1.8.0
       github.com/stretchr/testify v1.8.4
   )

   replace github.com/old/pkg => github.com/new/pkg v1.0.0

   exclude github.com/bad/pkg v1.2.3

📌 BEST PRACTICE:
   • Commit go.mod and go.sum to version control
   • Do NOT commit vendor/ unless necessary
   • Run go mod tidy before commit
   • Use semantic versioning for modules
`)

	// ============================================================================
	// بخش 4: go generate - تولید خودکار کد
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("⚙️ SECTION 4: go generate - Code Generation")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────┐
│ go generate COMMANDS                                           │
├─────────────────────────────────────────────────────────────────┤
│ go generate             - اجرا در текущем دایرکتوری            │
│ go generate ./...       - اجرا در همه دایرکتوری‌ها             │
│ go generate -x ./...    - نمایش دستورات در حین اجرا            │
│ go generate -n ./...    - نمایش دستورات بدون اجرا              │
│ go generate -run pattern - فقط دستورات matching pattern       │
│ go generate -v ./...    - verbose output                       │
└─────────────────────────────────────────────────────────────────┘

📝 SYNTAX (در فایل‌های Go):

   //go:generate <command> <arguments>

📚 COMMON USE CASES:

   1. Generate mocks
      //go:generate mockgen -source=user.go -destination=mock_user.go

   2. Generate stringer
      //go:generate stringer -type=Status

   3. Generate protobuf
      //go:generate protoc --go_out=. *.proto

   4. Generate embedded files
      //go:generate go run gen_assets.go

   5. Generate SQL code
      //go:generate sqlc generate

   6. Generate documentation
      //go:generate go run docs/generate.go

   7. Generate errors
      //go:generate go run golang.org/x/tools/cmd/stringer -type=ErrorCode

💡 BEST PRACTICE:
   • Keep generation commands simple and repeatable
   • Commit generated code (for reproducibility)
   • Run go generate in CI/CD
   • Use //go:generate comments near related types
`)
}

// ============================================================================
// بخش 5: مثال‌های عملی (با کامنت‌های go:generate)
// ============================================================================

// 5.1 مثال: Stringer برای enum
type Status int

const (
	Pending Status = iota
	Active
	Inactive
	Deleted
)

//go:generate stringer -type=Status -output=status_string.go

// 5.2 مثال: Mock generation
type UserRepository interface {
	GetUser(id int) (*User, error)
	SaveUser(user *User) error
}

//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE

// 5.3 مثال: Embedded files
//go:generate go run embed_gen.go

// 5.4 مثال: Generate from template
//go:generate go run gen/template_gen.go -input=data.json -output=generated.go

// ============================================================================
// بخش 6: ترکیب ابزارها (Workflow)
// ============================================================================

func demonstrateWorkflow() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🔄 COMPLETE WORKFLOW")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────┐
│ TYPICAL DEVELOPMENT WORKFLOW                                   │
├─────────────────────────────────────────────────────────────────┤
│                                                                │
│  1. Initialize project                                        │
│     $ go mod init github.com/user/project                     │
│                                                                │
│  2. Add dependencies                                           │
│     $ go get github.com/gorilla/mux                           │
│     $ go get github.com/stretchr/testify                      │
│                                                                │
│  3. Write code with //go:generate directives                   │
│     //go:generate stringer -type=MyType                       │
│                                                                │
│  4. Generate code                                              │
│     $ go generate ./...                                       │
│                                                                │
│  5. Format code                                                │
│     $ go fmt ./...                                            │
│                                                                │
│  6. Run static analysis                                        │
│     $ go vet ./...                                            │
│                                                                │
│  7. Run tests                                                  │
│     $ go test -race -cover ./...                              │
│                                                                │
│  8. Tidy dependencies                                          │
│     $ go mod tidy                                             │
│                                                                │
│  9. Build                                                      │
│     $ go build -o bin/app .                                   │
│                                                                │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ CI/CD PIPELINE CHECKLIST                                       │
├─────────────────────────────────────────────────────────────────┤
│                                                                │
│  ✅ go mod verify                                              │
│  ✅ go mod download                                            │
│  ✅ go generate ./...                                          │
│  ✅ go fmt -l . | wc -l (should be 0)                         │
│  ✅ go vet ./...                                               │
│  ✅ go test -race -cover ./...                                 │
│  ✅ go build ./...                                             │
│                                                                │
└─────────────────────────────────────────────────────────────────┘
`)
}

// ============================================================================
// بخش 7: اشتباهات رایج
// ============================================================================

func demonstrateCommonMistakes() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("❌ COMMON MISTAKES")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
❌ Mistake 1: Committing unformatted code
   ✅ Run go fmt before commit

❌ Mistake 2: Ignoring go vet warnings
   ✅ Fix all warnings (they indicate real bugs)

❌ Mistake 3: Committing vendor/ unnecessarily
   ✅ Only commit vendor/ for production builds

❌ Mistake 4: Not running go mod tidy
   ✅ Always tidy before commit

❌ Mistake 5: Hardcoding versions in go.mod
   ✅ Use go get -u to update

❌ Mistake 6: Not committing go.sum
   ✅ go.sum is essential for security

❌ Mistake 7: Forgetting to run go generate
   ✅ Add to CI/CD pipeline

❌ Mistake 8: Using //go:generate for manual tasks
   ✅ Keep commands automated and repeatable
`)
}

// ============================================================================
// بخش 8: نکات پیشرفته
// ============================================================================

func demonstrateAdvancedTips() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("⚡ ADVANCED TIPS")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────┐
│ ADVANCED go mod                                                │
├─────────────────────────────────────────────────────────────────┤
│ go mod edit -go=1.22           - به‌روزرسانی نسخه Go           │
│ go mod edit -toolchain=go1.22.0 - تنظیم toolchain              │
│ go mod edit -retract=v1.0.0    -撤回 کردن نسخه                  │
│ go mod graph | grep <pkg>      - جستجو در درخت وابستگی         │
│ GOSUMDB=off go get <pkg>       - غیرفعال کردن checksum         │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ ADVANCED go generate                                           │
├─────────────────────────────────────────────────────────────────┤
│ //go:generate -command mycmd go run ./cmd/mycmd               │
│ //go:generate mycmd --input file.txt --output out.go          │
│                                                                │
│ //go:generate go run golang.org/x/tools/cmd/stringer \        │
│ //go:generate   -type=MyType \                                │
│ //go:generate   -output=mytype_string.go                      │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ ENVIRONMENT VARIABLES                                          │
├─────────────────────────────────────────────────────────────────┤
│ GO111MODULE=on/off/auto    - فعال/غیرفعال کردن ماژول‌ها         │
│ GOPROXY=https://proxy.golang.org - پروکسی ماژول‌ها             │
│ GOPRIVATE=github.com/myorg - ماژول‌های خصوصی                   │
│ GOSUMDB=sum.golang.org     - checksum database                 │
│ GONOSUMDB=github.com/myorg - غیرفعال کردن checksum برای خصوصی  │
│ GOTOOLCHAIN=auto           - مدیریت خودکار toolchain           │
└─────────────────────────────────────────────────────────────────┘
`)
}

// ============================================================================
// بخش 9: جمع‌بندی
// ============================================================================

func main() {
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("🎯 GO COMMAND LINE TOOLS - COMPLETE GUIDE")
	fmt.Println(strings.Repeat("=", 80))

	// این فایل فقط توضیحات است و کد واقعی ندارد.
	// تمام دستورات باید در ترمینال اجرا شوند.

	fmt.Println(`
╔═════════════════════════════════════════════════════════════════╗
║                    QUICK REFERENCE CARD                         ║
╠═════════════════════════════════════════════════════════════════╣
║                                                                 ║
║  go fmt ./...        → Format all code                         ║
║  go vet ./...        → Static analysis                         ║
║  go mod init <name>   → Create module                          ║
║  go get <pkg>        → Add dependency                          ║
║  go mod tidy         → Clean dependencies                      ║
║  go generate ./...    → Generate code                          ║
║  go test ./...       → Run tests                               ║
║  go build ./...      → Build                                   ║
║  go install <pkg>    → Install tool                            ║
║                                                                 ║
╚═════════════════════════════════════════════════════════════════╝

💡 GOLDEN RULES:

   1. go fmt: Run before every commit
   2. go vet: Fix all warnings (no exceptions)
   3. go mod: Commit go.mod AND go.sum
   4. go generate: Use for repetitive code
   5. CI/CD: Include all tools in pipeline
   6. Never ignore tool warnings
   7. Keep tools updated (go install <tool>@latest)
   8. Use //go:generate for mocks, stringers, etc.
   9. Run go mod tidy regularly
   10. Use -race flag with tests

📚 USEFUL COMMANDS COMBINATIONS:

   # Full check before commit
   go fmt ./... && go vet ./... && go test -race -cover ./...

   # Update all dependencies
   go get -u ./... && go mod tidy

   # Clean and rebuild
   go clean -modcache && go mod download && go build ./...

   # Generate and test
   go generate ./... && go test -v ./...
`)
}

// تابع کمکی برای تکرار رشته
func stringsRepeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}

/*
خلاصه دستورات مهم
ابزار	دستور	توضیح
go fmt	go fmt ./...	فرمت کردن همه فایل‌ها
go fmt -d ./...	نمایش diff بدون اعمال
go fmt -l ./...	نمایش فایل‌هایی که تغییر می‌کنند
go vet	go vet ./...	آنالیز استاتیک
go vet -v ./...	با خروجی verbose
go mod	go mod init <name>	ایجاد ماژول جدید
go get <pkg>	افزودن وابستگی
go mod tidy	پاکسازی وابستگی‌ها
go mod vendor	کپی وابستگی‌ها در vendor
go list -m all	لیست همه وابستگی‌ها
go generate	go generate ./...	اجرای همه دستورات generate
go generate -x ./...	نمایش دستورات در حین اجرا
go generate -run pattern	اجرای دستورات matching pattern

*/
