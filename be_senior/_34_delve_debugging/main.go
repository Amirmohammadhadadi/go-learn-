// ============================================================================
// FILE: delve_debugging_guide.go
// TITLE: راهنمای کامل دیباگ با Delve در Go - از نصب تا دیباگ حرفه‌ای
// HOW TO RUN: go run delve_debugging_guide.go (این فایل فقط توضیحات است)
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - Delve چیست و چرا نیاز است؟
// ============================================================================
//
// Delve یک دیباگر مخصوص زبان Go است [citation:1][citation:3].
// اهداف اصلی Delve [citation:3]:
// 1. سادگی در استفاده (Easy to invoke and easy to use)
// 2. قابلیت‌های کامل دیباگ (Breakpoints, variable inspection, stack traces)
// 3. پشتیبانی از همزمانی (Goroutines, channels)
// 4. حالت Headless برای اتصال از راه دور
//
// چرا Delve به جای GDB؟
// - درک کامل از Goroutines و Channelها
// - پشتیبانی بهتر از انواع داده Go (slices, maps, interfaces)
// - سینتکس ساده‌تر و اختصاصی Go
// - دیباگ بدون نیاز به اطلاعات DWARF اضافی
//
// قانون طلایی:
// "برای دیباگ برنامه‌های Go حتماً از Delve استفاده کن.
//  برنامه را با -gcflags='all=-N -l' کامپایل کن تا بهینه‌سازی‌ها غیرفعال شوند [citation:5].
//  در VSCode نیازی به تنظیمات دستی نیست، Go Extension این کار را انجام می‌دهد [citation:2]."
// ============================================================================

package __delve_debugging

import (
	"fmt"
	"strings"
)

// این فایل فقط یک فایل توضیحی است.
// تمام دستورات باید در ترمینال اجرا شوند.

func main() {
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("🎯 COMPLETE DELVE DEBUGGING GUIDE")
	fmt.Println("Installation | Commands | Headless Mode | VS Code Integration")
	fmt.Println(strings.Repeat("=", 80))

	// ============================================================================
	// بخش 1: نصب Delve
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📦 SECTION 1: INSTALLATION")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────┐
│ INSTALLATION METHODS                                           │
├─────────────────────────────────────────────────────────────────┤
│                                                                │
│  1. Using go install (RECOMMENDED) [citation:5][citation:8]   │
│     $ go install github.com/go-delve/delve/cmd/dlv@latest     │
│                                                                │
│  2. Using VS Code Go Extension [citation:2][citation:9]       │
│     - Open Command Palette (Ctrl+Shift+P)                     │
│     - Select "Go: Install/Update Tools"                       │
│     - Select "dlv" and click OK                               │
│                                                                │
│  3. From source (for development)                             │
│     $ git clone https://github.com/go-delve/delve.git         │
│     $ cd delve                                                 │
│     $ go install cmd/dlv                                      │
│                                                                │
│  4. Verify installation                                       │
│     $ dlv version                                             │
│     # Output: Delve Debugger Version: 1.25.0... [citation:7]  │
│                                                                │
└─────────────────────────────────────────────────────────────────┘

💡 PATH SETUP:
   Add $GOPATH/bin or $HOME/go/bin to your PATH [citation:8]:
   $ echo 'export PATH=$PATH:$HOME/go/bin' >> ~/.bashrc
   $ source ~/.bashrc
`)
	// ============================================================================
	// بخش 2: کامپایل برنامه برای دیباگ
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🔧 SECTION 2: COMPILING FOR DEBUGGING")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────┐
│ IMPORTANT COMPILATION FLAGS                                    │
├─────────────────────────────────────────────────────────────────┤
│                                                                │
│  go build -gcflags="all=-N -l" [citation:5]                   │
│                                                                │
│  Flags explanation:                                            │
│  • -N  : Disable optimizations                                 │
│  • -l  : Disable inlining                                      │
│                                                                │
│  WHY? Optimizations can cause:                                 │
│  • Variables to be optimized away (not visible)                │
│  • Line numbers to be incorrect                                │
│  • Step-by-step execution to jump around                       │
│                                                                │
│  Delve automatically adds these flags when using dlv debug [citation:9]│
│  For pre-built binaries, ensure they were built with these flags│
│                                                                │
└─────────────────────────────────────────────────────────────────┘
`)
	// ============================================================================
	// بخش 3: حالت‌های اجرای Delve
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🚀 SECTION 3: DELVE MODES")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────┐
│ MODE 1: dlv debug - Build and Debug                           │
├─────────────────────────────────────────────────────────────────┤
│                                                                │
│  $ dlv debug [package]                                         │
│                                                                │
│  • Builds the program with debug flags                         │
│  • Starts the program and attaches debugger                    │
│  • Stops at program start (before main) [citation:10]         │
│                                                                │
│  Examples:                                                     │
│  $ dlv debug                          # current directory     │
│  $ dlv debug main.go                  # specific file         │
│  $ dlv debug github.com/user/repo     # specific package      │
│  $ dlv debug -- -arg1 value           # pass args to program  │
│                                                                │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ MODE 2: dlv exec - Debug Existing Binary                      │
├─────────────────────────────────────────────────────────────────┤
│                                                                │
│  $ dlv exec ./binary                                           │
│                                                                │
│  • Debug a pre-compiled binary                                 │
│  • Binary must be built with -gcflags="all=-N -l"             │
│                                                                │
│  Example:                                                      │
│  $ go build -gcflags="all=-N -l" -o myapp                     │
│  $ dlv exec ./myapp                                            │
│  $ dlv exec ./myapp -- --config config.yaml                   │
│                                                                │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ MODE 3: dlv test - Debug Tests                                │
├─────────────────────────────────────────────────────────────────┤
│                                                                │
│  $ dlv test [package]                                          │
│                                                                │
│  • Builds and debugs test binary [citation:10]                │
│  • Stops at test execution start                               │
│                                                                │
│  Examples:                                                     │
│  $ dlv test                          # current package tests  │
│  $ dlv test ./...                    # all tests              │
│  $ dlv test github.com/user/repo     # specific package       │
│  $ dlv test -- -test.run TestName    # run specific test      │
│                                                                │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ MODE 4: Headless Mode (Remote Debugging) [citation:7][citation:8]│
├─────────────────────────────────────────────────────────────────┤
│                                                                │
│  $ dlv debug --headless --listen=:2345 --api-version=2        │
│                                                                │
│  • Starts Delve as a server waiting for connections           │
│  • IDE connects to this port (VS Code, GoLand, etc.)          │
│  • Allows debugging of containerized/remote applications      │
│                                                                │
│  Common flags [citation:1]:                                   │
│  • --headless         : Run without interactive terminal      │
│  • --listen=:2345     : Listen on port 2345                   │
│  • --api-version=2    : Use API v2 (required for DAP)         │
│  • --accept-multiclient: Allow multiple client connections    │
│  • --log              : Enable logging                        │
│  • --continue         : Start program immediately             │
│                                                                │
│  Example with all flags:                                      │
│  $ dlv debug --headless --listen=:2345 --api-version=2 \      │
│    --accept-multiclient --log --continue                      │
│                                                                │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ MODE 5: dlv attach - Attach to Running Process               │
├─────────────────────────────────────────────────────────────────┤
│                                                                │
│  $ dlv attach <PID>                                            │
│                                                                │
│  • Attaches to an already running Go process                   │
│  • Process must have been built with debug symbols             │
│                                                                │
│  Example:                                                      │
│  $ ps aux | grep myapp        # find process ID               │
│  $ dlv attach 12345                                           │
│                                                                │
└─────────────────────────────────────────────────────────────────┘
`)
	// ============================================================================
	// بخش 4: دستورات اصلی Delve
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("⌨️ SECTION 4: DELVE COMMANDS")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────┐
│ BREAKPOINT COMMANDS [citation:5][citation:8]                  │
├─────────────────────────────────────────────────────────────────┤
│                                                                │
│  break (b) <location>        - Set breakpoint                  │
│  break main.go:25            - Break at line 25 in main.go    │
│  break main.main             - Break at main function         │
│  break fmt.Println           - Break at function              │
│  break +5                    - Break 5 lines from current     │
│  break -3                    - Break 3 lines before current   │
│                                                                │
│  cond <breakpoint> <expr>    - Conditional breakpoint         │
│  cond 1 i == 10              - Break when i equals 10         │
│                                                                │
│  clear <breakpoint>          - Remove breakpoint by number    │
│  clear main.go:25            - Remove breakpoint at location  │
│  clearall                    - Remove all breakpoints         │
│                                                                │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ EXECUTION CONTROL COMMANDS [citation:5][citation:8]           │
├─────────────────────────────────────────────────────────────────┤
│                                                                │
│  continue (c)                - Continue execution until break │
│  next (n)                    - Step over (don't enter func)   │
│  step (s)                    - Step into function             │
│  stepout (so)                - Step out of current function   │
│  restart (r)                 - Restart the program            │
│  exit (q)                    - Exit Delve                     │
│                                                                │
│  Example workflow:                                             │
│  (dlv) break main.main                                        │
│  (dlv) continue                                               │
│  (dlv) next                                                   │
│  (dlv) step                                                   │
│  (dlv) print result                                           │
│                                                                │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ INFORMATION DISPLAY COMMANDS [citation:5][citation:8]         │
├─────────────────────────────────────────────────────────────────┤
│                                                                │
│  print (p) <expr>            - Print variable value           │
│  p myVar                     - Print variable                 │
│  p &myVar                    - Print address                  │
│  p *myPointer                - Dereference pointer            │
│  p myStruct.Field            - Print struct field             │
│  p mySlice[2:5]              - Print slice subrange           │
│                                                                │
│  whatis <expr>               - Print type of expression       │
│                                                                │
│  locals                      - Print local variables [citation:5]│
│  args                        - Print function arguments [citation:5]│
│  vars                        - Print global variables         │
│                                                                │
│  stack (bt)                  - Print stack trace [citation:8] │
│  stack 10                    - Print 10 stack frames          │
│  stack -1                    - Print all frames               │
│                                                                │
│  list (ls)                   - Show source code [citation:8]  │
│  list                        - Show around current line       │
│  list main.go:20             - Show around line 20            │
│                                                                │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ GOROUTINE COMMANDS [citation:5]                               │
├─────────────────────────────────────────────────────────────────┤
│                                                                │
│  goroutines (grs)            - List all goroutines            │
│  goroutine <id>              - Switch to specific goroutine   │
│  goroutine 5                 - Switch to goroutine 5          │
│  goroutine 5 stack           - Show stack of goroutine 5      │
│                                                                │
│  Example:                                                      │
│  (dlv) goroutines                                             │
│  (dlv) goroutine 3                                            │
│  (dlv) print myVar                                            │
│                                                                │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ CONFIGURATION COMMANDS                                        │
├─────────────────────────────────────────────────────────────────┤
│                                                                │
│  config max-string-len <n>    - Max string length to print    │
│  config max-array-values <n>  - Max array elements to print   │
│  config max-variable-recurse <n> - Recursion depth for types  │
│                                                                │
│  Example:                                                      │
│  (dlv) config max-string-len 256                              │
│  (dlv) config max-array-values 100                            │
│                                                                │
└─────────────────────────────────────────────────────────────────┘
`)
	// ============================================================================
	// بخش 5: مثال کامل دیباگ
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("💡 SECTION 5: COMPLETE DEBUGGING EXAMPLE")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────┐
│ STEP 1: Create a sample program (main.go)                     │
├─────────────────────────────────────────────────────────────────┤
│                                                                │
│  package main                                                  │
│                                                                │
│  import "fmt"                                                  │
│                                                                │
│  func calculate(x, y int) int {                               │
│      result := x * y                                          │
│      return result                                             │
│  }                                                            │
│                                                                │
│  func main() {                                                │
│      a := 10                                                   │
│      b := 20                                                   │
│      c := calculate(a, b)                                     │
│      fmt.Printf("Result: %d\n", c)                            │
│  }                                                            │
│                                                                │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ STEP 2: Start debugging                                       │
├─────────────────────────────────────────────────────────────────┤
│                                                                │
│  $ dlv debug main.go                                          │
│  Type 'help' for list of commands.                            │
│  (dlv)                                                        │
│                                                                │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ STEP 3: Set breakpoints and debug                             │
├─────────────────────────────────────────────────────────────────┤
│                                                                │
│  (dlv) break main.calculate                                   │
│  Breakpoint 1 set at 0x... for main.calculate()               │
│                                                                │
│  (dlv) break main.go:14                                       │
│  Breakpoint 2 set at 0x... for main.main() line 14            │
│                                                                │
│  (dlv) continue                                               │
│  > main.main() ./main.go:10 (hits goroutine(1):1 total:1)    │
│     9: func main() {                                          │
│ => 10:     a := 10                                            │
│    11:     b := 20                                            │
│    12:     c := calculate(a, b)                               │
│    13:     fmt.Printf("Result: %d\n", c)                      │
│                                                                │
│  (dlv) next                                                   │
│  > main.main() ./main.go:11 (hits goroutine(1):1 total:1)    │
│    10:     a := 10                                            │
│ => 11:     b := 20                                            │
│    12:     c := calculate(a, b)                               │
│    13:     fmt.Printf("Result: %d\n", c)                      │
│                                                                │
│  (dlv) print a                                                │
│  10                                                           │
│                                                                │
│  (dlv) next                                                   │
│  > main.calculate() ./main.go:6 (hits goroutine(1):1 total:1)│
│     5: func calculate(x, y int) int {                         │
│ =>  6:     result := x * y                                    │
│     7:     return result                                      │
│     8: }                                                      │
│                                                                │
│  (dlv) print x, y                                             │
│  10 20                                                        │
│                                                                │
│  (dlv) next                                                   │
│  > main.calculate() ./main.go:7 (hits goroutine(1):1 total:1)│
│     6:     result := x * y                                    │
│ =>  7:     return result                                      │
│     8: }                                                      │
│                                                                │
│  (dlv) print result                                           │
│  200                                                          │
│                                                                │
│  (dlv) continue                                               │
│  Result: 200                                                  │
│  Process 12345 has exited with status 0                       │
│  (dlv) exit                                                   │
│                                                                │
└─────────────────────────────────────────────────────────────────┘
`)
	// ============================================================================
	// بخش 6: VS Code Integration
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🖥️ SECTION 6: VS CODE INTEGRATION")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────┐
│ AUTOMATIC SETUP (No configuration needed!) [citation:2]       │
├─────────────────────────────────────────────────────────────────┤
│                                                                │
│  1. Install Go Extension (golang.go)                          │
│  2. Open a Go file                                            │
│  3. Press F5 or go to Run → Start Debugging                   │
│  4. VS Code automatically uses Delve [citation:9]            │
│                                                                │
│  The Go extension will prompt to install Delve if missing [citation:2]│
│                                                                │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ MANUAL launch.json CONFIGURATION [citation:2][citation:9]     │
├─────────────────────────────────────────────────────────────────┤
│                                                                │
│  Create .vscode/launch.json:                                  │
│                                                                │
│  {                                                            │
│    "version": "0.2.0",                                        │
│    "configurations": [                                        │
│      {                                                        │
│        "name": "Launch Package",                              │
│        "type": "go",                                          │
│        "request": "launch",                                   │
│        "mode": "auto",                                        │
│        "program": "${workspaceFolder}",                       │
│        "env": {},                                             │
│        "args": []                                             │
│      },                                                       │
│      {                                                        │
│        "name": "Launch File",                                 │
│        "type": "go",                                          │
│        "request": "launch",                                   │
│        "mode": "auto",                                        │
│        "program": "${file}"                                   │
│      },                                                       │
│      {                                                        │
│        "name": "Launch Test Function",                        │
│        "type": "go",                                          │
│        "request": "launch",                                   │
│        "mode": "test",                                        │
│        "program": "${workspaceFolder}",                       │
│        "args": ["-test.run", "MyTestFunction"]                │
│      },                                                       │
│      {                                                        │
│        "name": "Attach to Process",                           │
│        "type": "go",                                          │
│        "request": "attach",                                   │
│        "mode": "local",                                       │
│        "processId": "${command:pickGoProcess}"                │
│      },                                                       │
│      {                                                        │
│        "name": "Connect to Remote",                           │
│        "type": "go",                                          │
│        "request": "attach",                                   │
│        "mode": "remote",                                      │
│        "remotePath": "${workspaceFolder}",                    │
│        "port": 2345,                                          │
│        "host": "127.0.0.1"                                    │
│      }                                                       │
│    ]                                                          │
│  }                                                            │
│                                                                │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ DELVE CONFIGURATION in VS Code [citation:2][citation:9]       │
├─────────────────────────────────────────────────────────────────┤
│                                                                │
│  In settings.json:                                            │
│                                                                │
│  "go.delveConfig": {                                          │
│    "debugAdapter": "dlv-dap",  // Use DAP protocol (default)  │
│    "apiVersion": 2,                                           │
│    "dlvLoadConfig": {                                         │
│      "maxStringLen": 256,          // Max string length       │
│      "maxArrayValues": 100,        // Max array elements      │
│      "maxStructFields": -1,        // All struct fields       │
│      "maxVariableRecurse": 3,      // Recursion depth         │
│      "followPointers": true        // Auto-dereference        │
│    },                                                         │
│    "showGlobalVariables": true     // Show globals in debug   │
│  }                                                            │
│                                                                │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ DEBUGGING SHORTCUTS in VS Code [citation:2]                   │
├─────────────────────────────────────────────────────────────────┤
│                                                                │
│  F5              - Start / Continue                           │
│  F10             - Step Over (next)                           │
│  F11             - Step Into (step)                           │
│  Shift+F11       - Step Out (stepout)                         │
│  Ctrl+Shift+F5   - Restart                                    │
│  Shift+F5        - Stop                                       │
│  F9              - Toggle breakpoint                          │
│                                                                │
└─────────────────────────────────────────────────────────────────┘
`)
	// ============================================================================
	// بخش 7: Docker Container Debugging
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🐳 SECTION 7: DOCKER CONTAINER DEBUGGING")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────┐
│ DOCKERFILE with Delve [citation:4]                            │
├─────────────────────────────────────────────────────────────────┤
│                                                                │
│  FROM golang:1.24                                             │
│  WORKDIR /app                                                  │
│                                                                │
│  # Copy go.mod and go.sum first (for caching)                 │
│  COPY go.mod go.sum ./                                        │
│  RUN go mod download                                          │
│                                                                │
│  # Copy source code                                           │
│  COPY . ./                                                     │
│                                                                │
│  # Install Delve                                              │
│  RUN go install github.com/go-delve/delve/cmd/dlv@latest      │
│                                                                │
│  # Expose ports                                               │
│  EXPOSE 8080 40000    # 8080: app, 40000: delve              │
│                                                                │
│  # Run Delve in headless mode [citation:4]                    │
│  CMD ["/go/bin/dlv", "debug", "--headless",                   │
│       "--listen=0.0.0.0:40000",                               │
│       "--api-version=2",                                      │
│       "--accept-multiclient",                                 │
│       "--log",                                                │
│       "--continue"]                                           │
│                                                                │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ VS CODE launch.json for Docker [citation:4]                   │
├─────────────────────────────────────────────────────────────────┤
│                                                                │
│  {                                                            │
│    "name": "Attach to Docker",                                │
│    "type": "go",                                              │
│    "request": "attach",                                       │
│    "mode": "remote",                                          │
│    "remotePath": "/app",                                      │
│    "port": 40000,                                             │
│    "host": "127.0.0.1",                                       │
│    "substitutePath": [                                        │
│      {                                                        │
│        "from": "\${workspaceFolder}",                         │
│        "to": "/app"                                           │
│      }                                                        │
│    ]                                                          │
│  }                                                            │
│                                                                │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ DOCKER COMMANDS                                               │
├─────────────────────────────────────────────────────────────────┤
│                                                                │
│  # Build image                                                │
│  $ docker build -t myapp:debug .                              │
│                                                                │
│  # Run container with port mapping                            │
│  $ docker run -d -p 8080:8080 -p 40000:40000 --name myapp-debug myapp:debug│
│                                                                │
│  # Check logs                                                 │
│  $ docker logs myapp-debug                                    │
│                                                                │
│  # Stop and remove                                            │
│  $ docker stop myapp-debug && docker rm myapp-debug           │
│                                                                │
└─────────────────────────────────────────────────────────────────┘
`)
	// ============================================================================
	// بخش 8: عیب‌یابی
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🔧 SECTION 8: TROUBLESHOOTING")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────┐
│ COMMON ISSUES AND SOLUTIONS                                   │
├─────────────────────────────────────────────────────────────────┤
│                                                                │
│  Issue 1: "could not find file" in breakpoints [citation:4]   │
│  • Problem: Path mismatch between local and remote            │
│  • Solution: Check substitutePath in launch.json              │
│                                                                │
│  Issue 2: Variables show as "unavailable"                     │
│  • Problem: Optimizations enabled [citation:5]                │
│  • Solution: Rebuild with -gcflags="all=-N -l"                │
│                                                                │
│  Issue 3: Connection refused on port 40000 [citation:4]       │
│  • Problem: Container not running or port not mapped          │
│  • Solution: docker ps and check port mapping                 │
│                                                                │
│  Issue 4: Breakpoints not hitting                             │
│  • Problem: Wrong file path or optimized build                │
│  • Solution: Verify substitutePath and build flags            │
│                                                                │
│  Issue 5: "Version of Delve is too old" [citation:2]          │
│  • Problem: Outdated dlv version                              │
│  • Solution: Update with go install ...@latest                │
│                                                                │
│  Issue 6: Cannot find debug info for package                  │
│  • Problem: Module cache or build issues                      │
│  • Solution: Run go mod tidy and rebuild                      │
│                                                                │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ DEBUGGING TIPS                                                │
├─────────────────────────────────────────────────────────────────┤
│                                                                │
│  1. Use logOutput for verbose debugging [citation:9]          │
│     "logOutput": "debugger,dap"                               │
│                                                                │
│  2. Check VS Code debug logs [citation:4]                     │
│     Windows: %TEMP%\vscode-go-debug.txt                       │
│     Linux: /tmp/vscode-go-debug.txt                           │
│                                                                │
│  3. Use dlv trace for function call tracing                   │
│     $ dlv trace ./... myFunction                              │
│                                                                │
│  4. Use dlv replay for recorded debug sessions                │
│                                                                │
│  5. For race condition debugging, combine with -race flag     │
│     $ dlv debug -- --race                                     │
│                                                                │
└─────────────────────────────────────────────────────────────────┘
`)
	// ============================================================================
	// بخش 9: جمع‌بندی
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📚 SECTION 9: QUICK REFERENCE")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────┐
│ COMMAND CHEAT SHEET                                           │
├─────────────────────────────────────────────────────────────────┤
│                                                                │
│  Start Debugging:                                             │
│  • dlv debug               - Build and debug current package  │
│  • dlv exec ./binary       - Debug existing binary            │
│  • dlv test                - Debug tests                      │
│  • dlv attach <PID>        - Attach to running process        │
│  • dlv debug --headless    - Remote debugging mode            │
│                                                                │
│  Breakpoints:                                                 │
│  • b main.go:10            - Set breakpoint at line           │
│  • b main.main             - Set breakpoint at function       │
│  • cond 1 i == 10          - Conditional breakpoint           │
│  • clear 1                 - Remove breakpoint                │
│  • clearall                - Remove all breakpoints           │
│                                                                │
│  Execution:                                                   │
│  • c                       - Continue                         │
│  • n                       - Step over (next)                 │
│  • s                       - Step into                        │
│  • so                      - Step out                         │
│  • r                       - Restart                          │
│  • q                       - Quit                             │
│                                                                │
│  Inspection:                                                  │
│  • p <var>                 - Print variable                   │
│  • whatis <var>            - Show type                        │
│  • locals                  - Show local variables             │
│  • args                    - Show arguments                   │
│  • bt                      - Show stack trace                 │
│  • goroutines              - List goroutines                  │
│  • config max-string-len 256 - Increase string limit          │
│                                                                │
└─────────────────────────────────────────────────────────────────┘

💡 GOLDEN RULES:

   1. Always compile with -gcflags="all=-N -l" when debugging [citation:5]
   2. Use VS Code Go extension - it handles Delve automatically [citation:2]
   3. For Docker debugging, use headless mode with port mapping [citation:4]
   4. Check substitutePath when breakpoints don't hit [citation:4]
   5. Use conditional breakpoints for complex scenarios
   6. Increase max-string-len for long string inspection
   7. Use goroutines command to debug concurrent code
   8. Check VS Code debug logs when things don't work
   9. Keep Delve updated: go install ...@latest
   10. Use dlv trace for function call tracing

🎯 QUICK START:

   # Install Delve
   $ go install github.com/go-delve/delve/cmd/dlv@latest

   # Debug your program
   $ cd /path/to/your/project
   $ dlv debug

   # In Delve prompt
   (dlv) break main.main
   (dlv) continue
   (dlv) next
   (dlv) print myVariable
   (dlv) quit

   # Or use VS Code (easier!)
   - Open project in VS Code
   - Press F5
   - Set breakpoints by clicking line numbers
   - Use debugging toolbar
`)

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🎯 DELVE DEBUGGING - COMPLETE GUIDE")
	fmt.Println("Ready to debug your Go programs!")
	fmt.Println(strings.Repeat("=", 80))
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
خلاصه دستورات مهم Delve
دسته	دستور	توضیح
شروع دیباگ	dlv debug	Build و دیباگ پکیج текущего
dlv exec ./binary	دیباگ باینری موجود
dlv test	دیباگ تست‌ها
dlv attach <PID>	اتصال به فرآیند در حال اجرا
Breakpoint	b main.go:10	تنظیم breakpoint در خط مشخص
b main.main	تنظیم breakpoint در تابع
cond 1 i == 10	breakpoint شرطی
clear 1	حذف breakpoint
اجرا	c	ادامه اجرا
n	Step over (وارد تابع نشو)
s	Step into (وارد تابع شو)
so	Step out (خروج از تابع текущего)
مشاهده	p <var>	چاپ مقدار متغیر
locals	نمایش متغیرهای محلی
args	نمایش آرگومان‌های تابع
bt	نمایش stack trace
goroutines	لیست گوروتین‌ها
تنظیمات	config max-string-len 256	افزایش حد رشته
خروج	q	خروج از Delve

*/
