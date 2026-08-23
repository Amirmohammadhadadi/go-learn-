// ============================================================================
// FILE: os_complete_guide.go
// TITLE: راهنمای کامل پکیج os در Go - سیستم عامل، فایل‌ها، فرآیندها، محیط
// HOW TO RUN: go run os_complete_guide.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - پکیج os چیست و چه قابلیت‌هایی دارد؟
// ============================================================================
//
// پکیج os امکانات مستقل از سیستم عامل برای:
// 1. کار با فایل‌ها و دایرکتوری‌ها (ایجاد، خواندن، نوشتن، حذف)
// 2. متغیرهای محیطی (Environment Variables)
// 3. آرگومان‌های خط فرمان (Command Line Arguments)
// 4. اطلاعات فرآیند (Process Info)
// 5. سیگنال‌ها (Signals)
// 6. دستورات سیستم (Exec)
// 7. کار با stdin/stdout/stderr
//
// قانون طلایی:
// "برای کار با فایل‌ها، از os.Open و os.Create استفاده کن.
//  همیشه فایل‌ها را با defer close کن.
//  از path/filepath برای کار با مسیرها استفاده کن."
// ============================================================================

package __internal_packages

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// ============================================================================
// بخش 1: متغیرهای محیطی (Environment Variables)
// ============================================================================

func demonstrateEnvVariables() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🌍 ENVIRONMENT VARIABLES")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 1.1 گرفتن و تنظیم متغیرهای محیطی
	// ============================================
	fmt.Println("\n--- 1.1 Get and Set Environment Variables ---")

	// گرفتن یک متغیر
	path := os.Getenv("PATH")
	fmt.Printf("PATH length: %d characters\n", len(path))

	// گرفتن با مقدار پیش‌فرض
	home := os.Getenv("HOME")
	if home == "" {
		home = "/home/user"
	}
	fmt.Printf("HOME: %s\n", home)

	// تنظیم متغیر (فقط برای این برنامه)
	os.Setenv("MY_APP_MODE", "development")
	mode := os.Getenv("MY_APP_MODE")
	fmt.Printf("MY_APP_MODE: %s\n", mode)

	// بررسی وجود متغیر
	if val, exists := os.LookupEnv("GOPATH"); exists {
		fmt.Printf("GOPATH: %s\n", val)
	} else {
		fmt.Println("GOPATH not set")
	}

	// حذف متغیر
	os.Unsetenv("MY_APP_MODE")
	if _, exists := os.LookupEnv("MY_APP_MODE"); !exists {
		fmt.Println("MY_APP_MODE unset successfully")
	}

	// ============================================
	// 1.2 گرفتن همه متغیرهای محیطی
	// ============================================
	fmt.Println("\n--- 1.2 All Environment Variables ---")

	allEnv := os.Environ()
	fmt.Printf("Total environment variables: %d\n", len(allEnv))

	// نمایش چند نمونه
	for i, env := range allEnv {
		if i >= 5 {
			fmt.Printf("  ... and %d more\n", len(allEnv)-5)
			break
		}
		fmt.Printf("  %s\n", env)
	}

	// تبدیل به مپ
	envMap := make(map[string]string)
	for _, env := range allEnv {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 2 {
			envMap[parts[0]] = parts[1]
		}
	}
	fmt.Printf("Environment map size: %d\n", len(envMap))
}

// ============================================================================
// بخش 2: آرگومان‌های خط فرمان (Command Line Arguments)
// ============================================================================

func demonstrateCommandLineArgs() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📝 COMMAND LINE ARGUMENTS")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 2.1 خواندن آرگومان‌ها
	// ============================================
	fmt.Println("\n--- 2.1 Reading Arguments ---")

	// os.Args[0] نام برنامه است
	fmt.Printf("Program name: %s\n", os.Args[0])

	// آرگومان‌های ارسالی
	if len(os.Args) > 1 {
		fmt.Printf("Arguments received: %v\n", os.Args[1:])
	} else {
		fmt.Println("No arguments provided")
		fmt.Println("Try: go run os_complete_guide.go arg1 arg2 arg3")
	}

	// ============================================
	// 2.2 پردازش آرگومان‌ها (مثال ساده)
	// ============================================
	fmt.Println("\n--- 2.2 Simple Argument Parser ---")

	// پردازش flags ساده (برای پروژه‌های کوچک)
	// برای پروژه‌های بزرگ از پکیج flag استفاده کن
	var (
		verbose bool
		name    string
		count   int
	)

	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		switch arg {
		case "-v", "--verbose":
			verbose = true
		case "-n", "--name":
			if i+1 < len(os.Args) {
				name = os.Args[i+1]
				i++
			}
		case "-c", "--count":
			if i+1 < len(os.Args) {
				count, _ = strconv.Atoi(os.Args[i+1])
				i++
			}
		}
	}

	if verbose {
		fmt.Printf("  Verbose mode: ON\n")
	}
	if name != "" {
		fmt.Printf("  Name: %s\n", name)
	}
	if count > 0 {
		fmt.Printf("  Count: %d\n", count)
	}
}

// ============================================================================
// بخش 3: کار با فایل‌ها - عملیات پایه
// ============================================================================

func demonstrateBasicFileOperations() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📁 BASIC FILE OPERATIONS")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 3.1 ایجاد فایل
	// ============================================
	fmt.Println("\n--- 3.1 Creating Files ---")

	// os.Create - ایجاد یا خالی کردن فایل (اگر وجود داشته باشد)
	file1, err := os.Create("example1.txt")
	if err != nil {
		log.Printf("Error creating file: %v", err)
		return
	}
	defer file1.Close()
	defer os.Remove("example1.txt")

	file1.WriteString("Hello from os.Create\n")
	fmt.Println("  Created example1.txt")

	// ============================================
	// 3.2 باز کردن فایل
	// ============================================
	fmt.Println("\n--- 3.2 Opening Files ---")

	// os.Open - فقط خواندن
	file2, err := os.Open("example1.txt")
	if err != nil {
		log.Printf("Error opening file: %v", err)
		return
	}
	defer file2.Close()

	data, _ := io.ReadAll(file2)
	fmt.Printf("  Read from file: %s", string(data))

	// os.OpenFile - با权限‌های مشخص
	file3, err := os.OpenFile("example2.txt",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	defer file3.Close()
	defer os.Remove("example2.txt")

	file3.WriteString("Appended line\n")
	fmt.Println("  Created example2.txt with append mode")

	// ============================================
	// 3.3 نوشتن در فایل
	// ============================================
	fmt.Println("\n--- 3.3 Writing to Files ---")

	file4, _ := os.Create("write_example.txt")
	defer os.Remove("write_example.txt")
	defer file4.Close()

	// نوشتن با WriteString
	n1, _ := file4.WriteString("String write\n")
	fmt.Printf("  WriteString: %d bytes\n", n1)

	// نوشتن با Write ([]byte)
	n2, _ := file4.Write([]byte("Byte slice write\n"))
	fmt.Printf("  Write: %d bytes\n", n2)

	// نوشتن با WriteAt (نوشتن در موقعیت خاص)
	n3, _ := file4.WriteAt([]byte("OVERWRITE"), 0)
	fmt.Printf("  WriteAt: %d bytes (overwrote beginning)\n", n3)

	// دریافت اطلاعات فایل
	info, _ := file4.Stat()
	fmt.Printf("  File size: %d bytes\n", info.Size())

	// ============================================
	// 3.4 خواندن از فایل
	// ============================================
	fmt.Println("\n--- 3.4 Reading from Files ---")

	// آماده‌سازی فایل
	readFile, _ := os.Create("read_example.txt")
	readFile.WriteString("Line 1\nLine 2\nLine 3\nLine 4\nLine 5\n")
	readFile.Close()
	defer os.Remove("read_example.txt")

	// باز کردن برای خواندن
	readFile, _ = os.Open("read_example.txt")
	defer readFile.Close()

	// روش 1: خواندن بایت به بایت
	buf := make([]byte, 8)
	readFile.Seek(0, 0)
	fmt.Println("  Reading byte by byte:")
	for {
		n, err := readFile.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("    Error: %v\n", err)
			break
		}
		fmt.Printf("    %s", buf[:n])
	}

	// روش 2: خواندن کل فایل با ReadFile
	data, _ = os.ReadFile("read_example.txt")
	fmt.Printf("\n  ReadFile: %s", string(data))
}

// ============================================================================
// بخش 4: کار با فایل‌ها - عملیات پیشرفته
// ============================================================================

func demonstrateAdvancedFileOperations() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🔧 ADVANCED FILE OPERATIONS")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 4.1 کپی، جابجایی، حذف
	// ============================================
	fmt.Println("\n--- 4.1 Copy, Move, Delete ---")

	// ایجاد فایل منبع
	src, _ := os.Create("source.txt")
	src.WriteString("Source content")
	src.Close()
	defer os.Remove("source.txt")

	// کپی فایل (روش دستی)
	dst, _ := os.Create("destination.txt")
	src, _ = os.Open("source.txt")
	io.Copy(dst, src)
	src.Close()
	dst.Close()
	defer os.Remove("destination.txt")
	fmt.Println("  Copied source.txt to destination.txt")

	// جابجایی (rename)
	os.Rename("destination.txt", "moved.txt")
	defer os.Remove("moved.txt")
	fmt.Println("  Renamed to moved.txt")

	// حذف
	os.Remove("moved.txt")
	fmt.Println("  Deleted moved.txt")

	// ============================================
	// 4.2 اطلاعات فایل (FileInfo)
	// ============================================
	fmt.Println("\n--- 4.2 File Information (FileInfo) ---")

	// ایجاد فایل تست
	testFile, _ := os.Create("info_test.txt")
	testFile.WriteString("Test content")
	testFile.Close()
	defer os.Remove("info_test.txt")

	info, err := os.Stat("info_test.txt")
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
		return
	}

	fmt.Printf("  Name: %s\n", info.Name())
	fmt.Printf("  Size: %d bytes\n", info.Size())
	fmt.Printf("  Mode: %s\n", info.Mode())
	fmt.Printf("  ModTime: %s\n", info.ModTime().Format(time.RFC3339))
	fmt.Printf("  IsDir: %v\n", info.IsDir())

	// بررسی وجود فایل
	if _, err := os.Stat("nonexistent.txt"); os.IsNotExist(err) {
		fmt.Println("  File does not exist (as expected)")
	}

	// ============================================
	// 4.3 تغییر permissions و ownership
	// ============================================
	fmt.Println("\n--- 4.3 Changing Permissions ---")

	permFile, _ := os.Create("perm_test.txt")
	permFile.Close()
	defer os.Remove("perm_test.txt")

	// تغییر permissions (chmod)
	err = os.Chmod("perm_test.txt", 0644)
	if err == nil {
		fmt.Println("  Changed permissions to 0644")
	}

	// تغییر owner (فقط در Unix/Linux)
	// os.Chown("perm_test.txt", uid, gid)

	// تغییر زمان آخرین دسترسی/تغییر
	now := time.Now()
	os.Chtimes("perm_test.txt", now, now)
	fmt.Println("  Updated file timestamps")
}

// ============================================================================
// بخش 5: کار با دایرکتوری‌ها
// ============================================================================

func demonstrateDirectories() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📂 DIRECTORY OPERATIONS")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 5.1 ایجاد، حذف و خواندن دایرکتوری
	// ============================================
	fmt.Println("\n--- 5.1 Create, Remove, Read Directory ---")

	// ایجاد دایرکتوری
	err := os.Mkdir("test_dir", 0755)
	if err == nil {
		fmt.Println("  Created directory: test_dir")
		defer os.RemoveAll("test_dir")
	} else if os.IsExist(err) {
		fmt.Println("  Directory already exists")
	}

	// ایجاد دایرکتوری با parentها (مانند mkdir -p)
	err = os.MkdirAll("test_dir/subdir/subsubdir", 0755)
	if err == nil {
		fmt.Println("  Created nested directories")
	}

	// خواندن محتویات دایرکتوری
	entries, _ := os.ReadDir("test_dir")
	fmt.Printf("  Contents of test_dir: %d entries\n", len(entries))

	// ایجاد چند فایل برای نمایش
	os.Create("test_dir/file1.txt")
	os.Create("test_dir/file2.txt")
	os.Mkdir("test_dir/empty_dir", 0755)

	// خواندن با جزئیات
	entries2, _ := os.ReadDir("test_dir")
	fmt.Println("  Directory contents:")
	for _, entry := range entries2 {
		info, _ := entry.Info()
		fmt.Printf("    %s (size: %d, isDir: %v)\n",
			entry.Name(), info.Size(), entry.IsDir())
	}

	// حذف دایرکتوری (فقط اگر خالی باشد)
	err = os.Remove("test_dir/empty_dir")
	if err == nil {
		fmt.Println("  Removed empty_dir")
	}

	// حذف دایرکتوری و همه محتویات
	err = os.RemoveAll("test_dir")
	if err == nil {
		fmt.Println("  Removed test_dir and all contents")
	}

	// ============================================
	// 5.2 مسیر جاری و تغییر مسیر
	// ============================================
	fmt.Println("\n--- 5.2 Working Directory ---")

	// گرفتن مسیر جاری
	wd, err := os.Getwd()
	if err == nil {
		fmt.Printf("  Current working directory: %s\n", wd)
	}

	// تغییر مسیر (موقت)
	oldWd, _ := os.Getwd()
	os.Chdir("/tmp")
	newWd, _ := os.Getwd()
	fmt.Printf("  Changed to: %s\n", newWd)
	os.Chdir(oldWd)
	fmt.Printf("  Changed back to: %s\n", oldWd)

	// ============================================
	// 5.3 پیمایش دایرکتوری (Walk)
	// ============================================
	fmt.Println("\n--- 5.3 Directory Walk ---")

	// ایجاد ساختار تست
	os.MkdirAll("walk_test/sub1", 0755)
	os.MkdirAll("walk_test/sub2/deep", 0755)
	os.Create("walk_test/file1.txt")
	os.Create("walk_test/sub1/file2.txt")
	os.Create("walk_test/sub2/deep/file3.txt")
	defer os.RemoveAll("walk_test")

	// پیمایش با WalkDir (Go 1.16+)
	fmt.Println("  Walking through walk_test:")
	err = filepath.WalkDir("walk_test", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		fmt.Printf("    %s (isDir: %v)\n", path, d.IsDir())
		return nil
	})
	if err != nil {
		fmt.Printf("  Walk error: %v\n", err)
	}
}

// ============================================================================
// بخش 6: مسیرها (Paths) - کار با filepath
// ============================================================================

func demonstratePaths() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🛤️ PATH OPERATIONS (filepath package)")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 6.1 عملیات پایه روی مسیر
	// ============================================
	fmt.Println("\n--- 6.1 Basic Path Operations ---")

	paths := []string{
		"/home/user/file.txt",
		"relative/path/to/file.go",
		"/usr/local/bin/",
		"file.txt",
	}

	for _, p := range paths {
		fmt.Printf("  Path: %s\n", p)
		fmt.Printf("    Dir: %s\n", filepath.Dir(p))
		fmt.Printf("    Base: %s\n", filepath.Base(p))
		fmt.Printf("    Ext: %s\n", filepath.Ext(p))
		fmt.Println()
	}

	// ============================================
	// 6.2 Join و Split مسیرها
	// ============================================
	fmt.Println("\n--- 6.2 Join and Split ---")

	// Join (اتصال مسیرها با جداکننده مناسب OS)
	joined := filepath.Join("home", "user", "documents", "file.txt")
	fmt.Printf("  Joined: %s\n", joined)

	// Split (جداسازی آخرین عنصر)
	dir, file := filepath.Split("/home/user/file.txt")
	fmt.Printf("  Split: dir=%s, file=%s\n", dir, file)

	// ============================================
	// 6.3 مسیر مطلق و نسبی
	// ============================================
	fmt.Println("\n--- 6.3 Absolute and Relative Paths ---")

	// تبدیل به مسیر مطلق
	abs, _ := filepath.Abs("relative/path")
	fmt.Printf("  Absolute: %s\n", abs)

	// تبدیل به مسیر نسبی
	rel, _ := filepath.Rel("/home/user", "/home/user/documents/file.txt")
	fmt.Printf("  Relative: %s\n", rel)

	// تمیز کردن مسیر
	clean := filepath.Clean("/home//user/../user/./file.txt")
	fmt.Printf("  Cleaned: %s\n", clean)

	// ============================================
	// 6.4 تطبیق الگو (Glob)
	// ============================================
	fmt.Println("\n--- 6.4 Glob Pattern Matching ---")

	// ایجاد فایل‌های تست
	os.MkdirAll("glob_test", 0755)
	os.Create("glob_test/file1.txt")
	os.Create("glob_test/file2.txt")
	os.Create("glob_test/file3.log")
	os.Create("glob_test/data.json")
	defer os.RemoveAll("glob_test")

	// الگوهای مختلف
	patterns := []string{
		"glob_test/*.txt",
		"glob_test/*.log",
		"glob_test/*.*",
	}

	for _, pattern := range patterns {
		matches, _ := filepath.Glob(pattern)
		fmt.Printf("  Pattern %s: %v\n", pattern, matches)
	}
}

// ============================================================================
// بخش 7: اطلاعات سیستم و فرآیند
// ============================================================================

func demonstrateSystemAndProcess() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("💻 SYSTEM AND PROCESS INFORMATION")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 7.1 اطلاعات سیستم
	// ============================================
	fmt.Println("\n--- 7.1 System Information ---")

	// Hostname
	hostname, _ := os.Hostname()
	fmt.Printf("  Hostname: %s\n", hostname)

	// Process ID
	fmt.Printf("  Process ID (PID): %d\n", os.Getpid())

	// Parent Process ID
	fmt.Printf("  Parent PID (PPID): %d\n", os.Getppid())

	// User ID (Unix/Linux)
	fmt.Printf("  User ID (UID): %d\n", os.Getuid())

	// Group ID (Unix/Linux)
	fmt.Printf("  Group ID (GID): %d\n", os.Getgid())

	// ============================================
	// 7.2 اطلاعات فرآیند
	// ============================================
	fmt.Println("\n--- 7.2 Process Information ---")

	// آرگومان‌ها
	fmt.Printf("  Arguments: %v\n", os.Args)

	// مسیر اجرایی
	execPath, _ := os.Executable()
	fmt.Printf("  Executable path: %s\n", execPath)

	// صفحات محیطی
	fmt.Printf("  Pagesize: %d\n", os.Getpagesize())

	// ============================================
	// 7.3 متغیرهای موقت
	// ============================================
	fmt.Println("\n--- 7.3 Temporary Files and Directories ---")

	// دایرکتوری موقت سیستم
	tempDir := os.TempDir()
	fmt.Printf("  System temp dir: %s\n", tempDir)

	// ایجاد فایل موقت
	tempFile, err := os.CreateTemp("", "example_*.txt")
	if err == nil {
		fmt.Printf("  Temporary file: %s\n", tempFile.Name())
		tempFile.WriteString("Temporary data")
		tempFile.Close()
		os.Remove(tempFile.Name())
	}

	// ایجاد دایرکتوری موقت
	tempDir2, err := os.MkdirTemp("", "example_dir_*")
	if err == nil {
		fmt.Printf("  Temporary directory: %s\n", tempDir2)
		os.RemoveAll(tempDir2)
	}
}

// ============================================================================
// بخش 8: استاندارد ورودی/خروجی/خطا (stdin, stdout, stderr)
// ============================================================================

func demonstrateStdIO() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📡 STANDARD I/O (stdin, stdout, stderr)")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 8.1 نوشتن در stdout و stderr
	// ============================================
	fmt.Println("\n--- 8.1 Writing to stdout/stderr ---")

	// نوشتن در stdout (os.Stdout)
	os.Stdout.WriteString("This goes to stdout\n")

	// نوشتن در stderr (os.Stderr)
	os.Stderr.WriteString("This goes to stderr\n")

	// استفاده از fmt (که پیش‌فرض می‌رود به stdout)
	fmt.Println("fmt.Println goes to stdout")

	// ============================================
	// 8.2 خواندن از stdin
	// ============================================
	fmt.Println("\n--- 8.2 Reading from stdin ---")
	fmt.Println("  (Skipping interactive read in this demo)")

	// مثال (در صورت نیاز)
	// reader := bufio.NewReader(os.Stdin)
	// fmt.Print("Enter text: ")
	// text, _ := reader.ReadString('\n')
	// fmt.Printf("You entered: %s", text)

	// ============================================
	// 8.3 بازگردانی stdout/stderr
	// ============================================
	fmt.Println("\n--- 8.3 Redirecting stdout/stderr ---")

	// ذخیره stdout اصلی
	oldStdout := os.Stdout

	// ایجاد فایل برای redirect
	r, w, _ := os.Pipe()
	os.Stdout = w

	// چاپ (به فایل می‌رود)
	fmt.Println("This line is redirected to pipe")

	// بازگردانی
	w.Close()
	os.Stdout = oldStdout

	// خواندن از pipe
	var buf bytes.Buffer
	io.Copy(&buf, r)
	fmt.Printf("  Redirected output: %q\n", buf.String())
}

// ============================================================================
// بخش 9: سیگنال‌ها (Signals)
// ============================================================================

func demonstrateSignals() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📡 SIGNAL HANDLING")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 9.1 دریافت سیگنال‌ها
	// ============================================
	fmt.Println("\n--- 9.1 Receiving Signals ---")

	// ایجاد کانال برای سیگنال‌ها
	sigChan := make(chan os.Signal, 1)

	// ثبت برای دریافت سیگنال‌های خاص
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// مثال غیرفعال (در دمو اجرا نمی‌شود)
	fmt.Println("  Signal handling configured for SIGINT and SIGTERM")
	fmt.Println("  (Press Ctrl+C to see in real execution)")

	// در برنامه واقعی:
	// go func() {
	//     sig := <-sigChan
	//     fmt.Printf("\nReceived signal: %v\n", sig)
	//     // انجام clean up
	//     os.Exit(0)
	// }()

	// متوقف کردن دریافت سیگنال‌ها
	signal.Stop(sigChan)
}

// ============================================================================
// بخش 10: اجرای دستورات سیستمی (os/exec)
// ============================================================================

func demonstrateExecCommands() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("⚙️ EXECUTING SYSTEM COMMANDS (os/exec)")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 10.1 اجرای دستورات ساده
	// ============================================
	fmt.Println("\n--- 10.1 Simple Commands ---")

	// اجرا و گرفتن خروجی
	out, err := exec.Command("echo", "Hello from exec").Output()
	if err == nil {
		fmt.Printf("  echo output: %s", out)
	}

	// ============================================
	// 10.2 اجرا با ورودی و خروجی
	// ============================================
	fmt.Println("\n--- 10.2 Commands with Input/Output ---")

	cmd := exec.Command("grep", "go")
	cmd.Stdin = strings.NewReader("go language\npython\njava\ngolang\n")

	output, _ := cmd.Output()
	fmt.Printf("  grep 'go' output: %s", output)

	// ============================================
	// 10.3 اجرای دستورات طولانی
	// ============================================
	fmt.Println("\n--- 10.3 Long-running Commands ---")

	cmd2 := exec.Command("sleep", "2")
	fmt.Println("  Running sleep 2...")

	// اجرا با timeout
	start := time.Now()
	err = cmd2.Run()
	elapsed := time.Since(start)

	if err == nil {
		fmt.Printf("  Command completed in %v\n", elapsed)
	}

	// ============================================
	// 10.4 بررسی وجود دستور
	// ============================================
	fmt.Println("\n--- 10.4 Checking Command Existence ---")

	commands := []string{"go", "python", "node", "nonexistent"}
	for _, cmdName := range commands {
		_, err := exec.LookPath(cmdName)
		if err == nil {
			fmt.Printf("  %s found in PATH\n", cmdName)
		} else {
			fmt.Printf("  %s not found\n", cmdName)
		}
	}
}

// ============================================================================
// بخش 11: اشتباهات رایج و نکات مهم
// ============================================================================

func demonstrateCommonMistakes() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("❌ COMMON MISTAKES WITH os PACKAGE")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println("\n❌ Mistake 1: Not closing files")
	fmt.Println("   file, _ := os.Open(\"file.txt\")")
	fmt.Println("   // missing defer file.Close()")
	fmt.Println("   ✅ defer file.Close()")

	fmt.Println("\n❌ Mistake 2: Assuming files always open successfully")
	fmt.Println("   file, _ := os.Open(\"missing.txt\")  // ignoring error")
	fmt.Println("   ✅ Always check errors")

	fmt.Println("\n❌ Mistake 3: Using hardcoded path separators")
	fmt.Println("   path := \"home/user/file.txt\"  // doesn't work on Windows")
	fmt.Println("   ✅ Use filepath.Join(\"home\", \"user\", \"file.txt\")")

	fmt.Println("\n❌ Mistake 4: Not handling permission errors")
	fmt.Println("   os.Mkdir(\"/root/dir\", 0755)  // may fail without sudo")
	fmt.Println("   ✅ Check error and handle appropriately")

	fmt.Println("\n❌ Mistake 5: Using Remove on non-empty directory")
	fmt.Println("   os.Remove(\"dir_with_files\")  // error")
	fmt.Println("   ✅ Use os.RemoveAll() for directories with content")

	fmt.Println("\n❌ Mistake 6: Not using absolute paths for consistency")
	fmt.Println("   os.Chdir(\"/some/dir\"); os.Open(\"file.txt\")  // fragile")
	fmt.Println("   ✅ Use filepath.Abs() or full paths")

	fmt.Println("\n❌ Mistake 7: Reading entire large files")
	fmt.Println("   data, _ := os.ReadFile(\"huge.log\")  // memory issue")
	fmt.Println("   ✅ Use bufio.Scanner or io.Copy with chunks")
}

// ============================================================================
// بخش 12: جمع‌بندی و جدول مرجع
// ============================================================================

func main() {
	fmt.Println(stringsRepeat("=", 80))
	fmt.Println("🎯 COMPLETE os PACKAGE GUIDE IN GO")
	fmt.Println(stringsRepeat("=", 80))

	// بخش 1: متغیرهای محیطی
	demonstrateEnvVariables()

	// بخش 2: آرگومان‌های خط فرمان
	demonstrateCommandLineArgs()

	// بخش 3: عملیات پایه فایل
	demonstrateBasicFileOperations()

	// بخش 4: عملیات پیشرفته فایل
	demonstrateAdvancedFileOperations()

	// بخش 5: دایرکتوری‌ها
	demonstrateDirectories()

	// بخش 6: مسیرها
	demonstratePaths()

	// بخش 7: اطلاعات سیستم و فرآیند
	demonstrateSystemAndProcess()

	// بخش 8: stdin/stdout/stderr
	demonstrateStdIO()

	// بخش 9: سیگنال‌ها
	demonstrateSignals()

	// بخش 10: اجرای دستورات
	demonstrateExecCommands()

	// بخش 11: اشتباهات رایج
	demonstrateCommonMistakes()

	// بخش 12: جمع‌بندی
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📚 os PACKAGE QUICK REFERENCE")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ FILE OPERATIONS                                                │")
	fmt.Println("├─────────────────────────────────────────────────────────────────┤")
	fmt.Println("│ os.Create(name)          - Create or truncate file            │")
	fmt.Println("│ os.Open(name)            - Open file read-only                │")
	fmt.Println("│ os.OpenFile(name, flag, perm) - Open with options             │")
	fmt.Println("│ os.ReadFile(name)        - Read entire file (Go 1.16+)        │")
	fmt.Println("│ os.WriteFile(name, data, perm) - Write entire file (Go 1.16+) │")
	fmt.Println("│ os.Remove(name)          - Delete file or empty directory     │")
	fmt.Println("│ os.RemoveAll(name)       - Delete recursively                 │")
	fmt.Println("│ os.Rename(old, new)      - Rename/move file                   │")
	fmt.Println("│ os.Stat(name)            - Get file info                      │")
	fmt.Println("└─────────────────────────────────────────────────────────────────┘")

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ DIRECTORY OPERATIONS                                           │")
	fmt.Println("├─────────────────────────────────────────────────────────────────┤")
	fmt.Println("│ os.Mkdir(name, perm)     - Create directory                   │")
	fmt.Println("│ os.MkdirAll(name, perm)  - Create directory with parents      │")
	fmt.Println("│ os.ReadDir(name)         - Read directory contents (Go 1.16+) │")
	fmt.Println("│ os.Getwd()               - Get current working directory      │")
	fmt.Println("│ os.Chdir(dir)            - Change working directory           │")
	fmt.Println("│ filepath.WalkDir(root, fn) - Walk directory tree              │")
	fmt.Println("└─────────────────────────────────────────────────────────────────┘")

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ ENVIRONMENT & PROCESS                                          │")
	fmt.Println("├─────────────────────────────────────────────────────────────────┤")
	fmt.Println("│ os.Getenv(key)           - Get environment variable           │")
	fmt.Println("│ os.Setenv(key, val)      - Set environment variable           │")
	fmt.Println("│ os.LookupEnv(key)        - Check if variable exists           │")
	fmt.Println("│ os.Environ()             - Get all environment variables      │")
	fmt.Println("│ os.Args                  - Command line arguments             │")
	fmt.Println("│ os.Getpid()              - Get process ID                     │")
	fmt.Println("│ os.Getppid()             - Get parent process ID              │")
	fmt.Println("│ os.Hostname()            - Get system hostname                │")
	fmt.Println("│ os.Exit(code)            - Exit program                       │")
	fmt.Println("└─────────────────────────────────────────────────────────────────┘")

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ PATH OPERATIONS (filepath)                                     │")
	fmt.Println("├─────────────────────────────────────────────────────────────────┤")
	fmt.Println("│ filepath.Join(parts...)  - Join path components               │")
	fmt.Println("│ filepath.Split(path)     - Split into dir and file            │")
	fmt.Println("│ filepath.Dir(path)       - Get directory part                 │")
	fmt.Println("│ filepath.Base(path)      - Get base name                      │")
	fmt.Println("│ filepath.Ext(path)       - Get file extension                 │")
	fmt.Println("│ filepath.Abs(path)       - Get absolute path                  │")
	fmt.Println("│ filepath.Rel(base, target) - Get relative path                │")
	fmt.Println("│ filepath.Clean(path)     - Clean up path                      │")
	fmt.Println("│ filepath.Glob(pattern)   - Match files with pattern           │")
	fmt.Println("└─────────────────────────────────────────────────────────────────┘")

	fmt.Println("\n💡 GOLDEN RULES:")
	fmt.Println("  1. Always check errors, especially for file operations")
	fmt.Println("  2. Always defer file.Close() after opening files")
	fmt.Println("  3. Use filepath package for cross-platform path handling")
	fmt.Println("  4. Use os.ReadFile/WriteFile for small files (Go 1.16+)")
	fmt.Println("  5. Use bufio.Scanner for reading large files line by line")
	fmt.Println("  6. Use os.CreateTemp for temporary files")
	fmt.Println("  7. Never ignore errors with _")
	fmt.Println("  8. Use os.RemoveAll for deleting directories with content")
	fmt.Println("  9. Handle signals for graceful shutdown")
	fmt.Println("  10. Use os.LookupEnv to check if env var exists")

	fmt.Println("\n🎯 FILE OPEN FLAGS:")
	fmt.Println("  os.O_RDONLY   - Read only")
	fmt.Println("  os.O_WRONLY   - Write only")
	fmt.Println("  os.O_RDWR     - Read and write")
	fmt.Println("  os.O_APPEND   - Append to end")
	fmt.Println("  os.O_CREATE   - Create if not exists")
	fmt.Println("  os.O_EXCL     - Used with O_CREATE, fail if exists")
	fmt.Println("  os.O_SYNC     - Synchronous writes")
	fmt.Println("  os.O_TRUNC    - Truncate on open")
}

// ============================================================================
// بخش 13: توابع کمکی
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
go run os_complete_guide.go

# اجرا با آرگومان‌های خط فرمان
go run os_complete_guide.go -v --name Ali --count 5

# اجرا و ارسال ورودی به stdin
echo "test input" | go run os_complete_guide.go
*/
