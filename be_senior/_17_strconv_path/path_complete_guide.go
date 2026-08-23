// ============================================================================
// FILE: path_complete_guide.go
// TITLE: راهنمای کامل پکیج path و path/filepath در Go - کار با مسیرها
// HOW TO RUN: go run path_complete_guide.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - پکیج‌های path و path/filepath چیستند؟
// ============================================================================
//
// Go دو پکیج برای کار با مسیرها ارائه می‌دهد:
//
// 1. path: برای مسیرهای با جداکننده slash (/)
//    - مناسب برای URL paths, Unix/Linux paths
//    - مستقل از سیستم عامل
//    - همیشه از / به عنوان جداکننده استفاده می‌کند
//
// 2. path/filepath: برای مسیرهای فایل سیستم
//    - قابل حمل بین سیستم‌عامل‌های مختلف
//    - از جداکننده مناسب سیستم عامل استفاده می‌کند (/ در Unix, \ در Windows)
//    - توابع بیشتر و قابلیت‌های پیشرفته‌تر
//
// قانون طلایی:
// "برای مسیرهای URL و مسیرهای عمومی از path استفاده کن.
//  برای مسیرهای فایل سیستم حتماً از path/filepath استفاده کن.
//  هرگز مسیرها را با + یا fmt.Sprintf نساز - همیشه از filepath.Join استفاده کن."
// ============================================================================

package __internal_packages

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"
)

func main() {
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("🎯 COMPLETE path & path/filepath PACKAGES GUIDE IN GO")
	fmt.Println("Working with file paths and URL paths")
	fmt.Println(strings.Repeat("=", 80))

	// ============================================================================
	// بخش 1: پکیج path - برای مسیرهای با slash (/)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📁 SECTION 1: path PACKAGE (for slash-separated paths)")
	fmt.Println(strings.Repeat("=", 80))

	// 1.1 Clean - پاکسازی و نرمال‌سازی مسیر
	fmt.Println("\n--- 1.1 path.Clean ---")
	// Clean returns the shortest path name equivalent to path by purely lexical processing.
	paths := []string{
		"a/b/../c",
		"a//b///c",
		"a/b/c/..",
		"/a/b/../c/./d",
		".",
		"..",
		"",
	}
	for _, p := range paths {
		cleaned := path.Clean(p)
		fmt.Printf("  Clean(%q) = %q\n", p, cleaned)
	}

	// 1.2 Join - اتصال مسیرها
	fmt.Println("\n--- 1.2 path.Join ---")
	// Join joins any number of path elements into a single path, separating them with slashes.
	joinExamples := [][]string{
		{"a", "b", "c"},
		{"a", "b/c", "d"},
		{"a", "../b", "c"},
		{"/a", "b", "c"},
		{"a", "b/", "/c"},
	}
	for _, elems := range joinExamples {
		joined := path.Join(elems...)
		fmt.Printf("  Join(%q) = %q\n", elems, joined)
	}

	// 1.3 Split - تقسیم مسیر به دایرکتوری و فایل
	fmt.Println("\n--- 1.3 path.Split ---")
	// Split splits path immediately following the final slash, separating it into a directory and file name component.
	splitExamples := []string{
		"/a/b/c.txt",
		"a/b/c.txt",
		"c.txt",
		"/",
		"",
	}
	for _, p := range splitExamples {
		dir, file := path.Split(p)
		fmt.Printf("  Split(%q) = dir: %q, file: %q\n", p, dir, file)
	}

	// 1.4 Dir - دایرکتوری مسیر
	fmt.Println("\n--- 1.4 path.Dir ---")
	// Dir returns all but the last element of path, typically the path's directory.
	dirExamples := []string{
		"/a/b/c.txt",
		"a/b/c.txt",
		"/a/b/c/",
		"c.txt",
		"/",
		".",
		"..",
	}
	for _, p := range dirExamples {
		dir := path.Dir(p)
		fmt.Printf("  Dir(%q) = %q\n", p, dir)
	}

	// 1.5 Base - آخرین عنصر مسیر (نام فایل)
	fmt.Println("\n--- 1.5 path.Base ---")
	// Base returns the last element of path.
	baseExamples := []string{
		"/a/b/c.txt",
		"a/b/c.txt",
		"/a/b/c/",
		"c.txt",
		"/",
		".",
		"..",
	}
	for _, p := range baseExamples {
		base := path.Base(p)
		fmt.Printf("  Base(%q) = %q\n", p, base)
	}

	// 1.6 Ext - پسوند فایل
	fmt.Println("\n--- 1.6 path.Ext ---")
	// Ext returns the file name extension used by path.
	extExamples := []string{
		"file.txt",
		"file.go",
		"file.tar.gz",
		"file",
		".hidden",
		"file.",
		"/a/b/file.jpg",
	}
	for _, p := range extExamples {
		ext := path.Ext(p)
		fmt.Printf("  Ext(%q) = %q\n", p, ext)
	}

	// 1.7 IsAbs - بررسی مطلق بودن مسیر
	fmt.Println("\n--- 1.7 path.IsAbs ---")
	// IsAbs reports whether the path is absolute.
	absExamples := []string{
		"/a/b/c",
		"a/b/c",
		"/",
		"",
	}
	for _, p := range absExamples {
		isAbs := path.IsAbs(p)
		fmt.Printf("  IsAbs(%q) = %v\n", p, isAbs)
	}

	// 1.8 Match - تطبیق مسیر با الگو
	fmt.Println("\n--- 1.8 path.Match ---")
	// Match reports whether name matches the shell file name pattern.
	matchExamples := []struct {
		pattern string
		name    string
	}{
		{"*.go", "main.go"},
		{"*.go", "test.txt"},
		{"*_test.go", "math_test.go"},
		{"*_test.go", "math.go"},
		{"a/*/c", "a/b/c"},
		{"a/*/c", "a/b/c/d"},
	}
	for _, ex := range matchExamples {
		matched, _ := path.Match(ex.pattern, ex.name)
		fmt.Printf("  Match(%q, %q) = %v\n", ex.pattern, ex.name, matched)
	}

	// ============================================================================
	// بخش 2: پکیج path/filepath - برای مسیرهای فایل سیستم
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📂 SECTION 2: path/filepath PACKAGE (OS-agnostic file paths)")
	fmt.Println(strings.Repeat("=", 80))

	// 2.1 Clean - پاکسازی مسیر (نسخه filepath)
	fmt.Println("\n--- 2.1 filepath.Clean ---")
	fileCleanExamples := []string{
		"a/b/../c",
		"a//b///c",
		"a/b/c/..",
		"/a/b/../c/./d",
		".",
		"..",
		"",
	}
	for _, p := range fileCleanExamples {
		cleaned := filepath.Clean(p)
		fmt.Printf("  Clean(%q) = %q\n", p, cleaned)
	}

	// 2.2 Join - اتصال مسیرها (نسخه filepath)
	fmt.Println("\n--- 2.2 filepath.Join ---")
	// Join joins any number of path elements into a single path, using OS-specific separator.
	fileJoinExamples := [][]string{
		{"home", "user", "file.txt"},
		{"home", "user", "..", "file.txt"},
		{"", "home", "user"},
		{"home", "", "file.txt"},
		{"C:", "Users", "file.txt"}, // Windows example
	}
	for _, elems := range fileJoinExamples {
		joined := filepath.Join(elems...)
		fmt.Printf("  Join(%q) = %q\n", elems, joined)
	}

	// 2.3 Split - تقسیم مسیر (نسخه filepath)
	fmt.Println("\n--- 2.3 filepath.Split ---")
	fileSplitExamples := []string{
		"/home/user/file.txt",
		"home/user/file.txt",
		"file.txt",
		"/",
		"",
		"C:\\Users\\file.txt", // Windows example
	}
	for _, p := range fileSplitExamples {
		dir, file := filepath.Split(p)
		fmt.Printf("  Split(%q) = dir: %q, file: %q\n", p, dir, file)
	}

	// 2.4 Dir - دایرکتوری (نسخه filepath)
	fmt.Println("\n--- 2.4 filepath.Dir ---")
	fileDirExamples := []string{
		"/home/user/file.txt",
		"home/user/file.txt",
		"/home/user/",
		"file.txt",
		"/",
		".",
		"..",
	}
	for _, p := range fileDirExamples {
		dir := filepath.Dir(p)
		fmt.Printf("  Dir(%q) = %q\n", p, dir)
	}

	// 2.5 Base - آخرین عنصر (نام فایل) (نسخه filepath)
	fmt.Println("\n--- 2.5 filepath.Base ---")
	fileBaseExamples := []string{
		"/home/user/file.txt",
		"home/user/file.txt",
		"/home/user/",
		"file.txt",
		"/",
		".",
		"..",
	}
	for _, p := range fileBaseExamples {
		base := filepath.Base(p)
		fmt.Printf("  Base(%q) = %q\n", p, base)
	}

	// 2.6 Ext - پسوند فایل (نسخه filepath)
	fmt.Println("\n--- 2.6 filepath.Ext ---")
	fileExtExamples := []string{
		"file.txt",
		"file.go",
		"file.tar.gz",
		"file",
		".hidden",
		"file.",
		"/home/user/file.jpg",
	}
	for _, p := range fileExtExamples {
		ext := filepath.Ext(p)
		fmt.Printf("  Ext(%q) = %q\n", p, ext)
	}

	// 2.7 IsAbs - بررسی مطلق بودن مسیر (نسخه filepath)
	fmt.Println("\n--- 2.7 filepath.IsAbs ---")
	fileAbsExamples := []string{
		"/home/user",
		"home/user",
		"/",
		"",
		"C:\\Users", // Windows absolute
		"\\",        // Windows root
	}
	for _, p := range fileAbsExamples {
		isAbs := filepath.IsAbs(p)
		fmt.Printf("  IsAbs(%q) = %v\n", p, isAbs)
	}

	// 2.8 Rel - مسیر نسبی بین دو مسیر
	fmt.Println("\n--- 2.8 filepath.Rel ---")
	// Rel returns a relative path that is lexically equivalent to targpath when joined to basepath.
	relExamples := []struct {
		base string
		targ string
	}{
		{"/a/b", "/a/b/c/d"},
		{"/a/b", "/a/b"},
		{"/a/b", "/a/c"},
		{"/a/b", "/a/b/../c"},
		{"/a/b", "/c/d"},
	}
	for _, ex := range relExamples {
		rel, err := filepath.Rel(ex.base, ex.targ)
		if err != nil {
			fmt.Printf("  Rel(%q, %q) = error: %v\n", ex.base, ex.targ, err)
		} else {
			fmt.Printf("  Rel(%q, %q) = %q\n", ex.base, ex.targ, rel)
		}
	}

	// 2.9 Abs - مسیر مطلق
	fmt.Println("\n--- 2.9 filepath.Abs ---")
	// Abs returns an absolute representation of path.
	// Note: In real code, check error. Here we're just showing the concept.
	absPath, _ := filepath.Abs(".")
	fmt.Printf("  Abs(\".\") = %q\n", absPath)

	// 2.10 VolumeName - نام حجم (در Windows)
	fmt.Println("\n--- 2.10 filepath.VolumeName ---")
	// VolumeName returns the volume name (e.g., "C:" on Windows).
	volumeExamples := []string{
		"C:\\Users\\file.txt",
		"D:\\Data",
		"/home/user",
		"\\\\server\\share\\file", // UNC path
	}
	for _, p := range volumeExamples {
		volume := filepath.VolumeName(p)
		fmt.Printf("  VolumeName(%q) = %q\n", p, volume)
	}

	// 2.11 ToSlash - تبدیل به slash (/)
	fmt.Println("\n--- 2.11 filepath.ToSlash ---")
	// ToSlash returns the result of replacing each separator character in path with a slash ('/').
	toSlashExamples := []string{
		"a\\b\\c",
		"a/b/c",
		"C:\\Users\\file.txt",
	}
	for _, p := range toSlashExamples {
		slashed := filepath.ToSlash(p)
		fmt.Printf("  ToSlash(%q) = %q\n", p, slashed)
	}

	// 2.12 FromSlash - تبدیل از slash به جداکننده سیستم
	fmt.Println("\n--- 2.12 filepath.FromSlash ---")
	// FromSlash returns the result of replacing each slash ('/') character in path with a separator.
	fromSlashExamples := []string{
		"a/b/c",
		"/home/user/file.txt",
		"http://example.com/path",
	}
	for _, p := range fromSlashExamples {
		fromSlashed := filepath.FromSlash(p)
		fmt.Printf("  FromSlash(%q) = %q\n", p, fromSlashed)
	}

	// ============================================================================
	// بخش 3: تطبیق الگو و Glob
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🔍 SECTION 3: PATTERN MATCHING AND GLOB")
	fmt.Println(strings.Repeat("=", 80))

	// 3.1 Match - تطبیق الگو (نسخه filepath)
	fmt.Println("\n--- 3.1 filepath.Match ---")
	fileMatchExamples := []struct {
		pattern string
		name    string
	}{
		{"*.go", "main.go"},
		{"*.go", "test.txt"},
		{"*_test.go", "math_test.go"},
		{"*_test.go", "math.go"},
		{"a/*/c", "a/b/c"},
		{"a/*/c", "a/b/c/d"},
		{"[a-z]*", "hello"},
		{"[a-z]*", "Hello"},
	}
	for _, ex := range fileMatchExamples {
		matched, err := filepath.Match(ex.pattern, ex.name)
		if err != nil {
			fmt.Printf("  Match(%q, %q) = error: %v\n", ex.pattern, ex.name, err)
		} else {
			fmt.Printf("  Match(%q, %q) = %v\n", ex.pattern, ex.name, matched)
		}
	}

	// 3.2 Glob - یافتن فایل‌های تطبیق‌دهنده با الگو
	fmt.Println("\n--- 3.2 filepath.Glob ---")
	// Glob returns the names of all files matching pattern or nil if there is no matching file.

	// ایجاد فایل‌های تست (در حافظه - شبیه‌سازی)
	fmt.Println("  (Creating test files in memory for demonstration)")

	// در عمل، Glob با فایل‌های واقعی کار می‌کند
	// matches, _ := filepath.Glob("*.go")
	// fmt.Printf("  Glob(\"*.go\") = %v\n", matches)

	// مثال‌های الگو
	patternExamples := []string{
		"*.txt",
		"*_test.go",
		"a/*/c",
		"*/*.go",
	}
	fmt.Println("  Pattern examples:")
	for _, p := range patternExamples {
		fmt.Printf("    %q - matches files with this pattern\n", p)
	}

	// ============================================================================
	// بخش 4: پیمایش دایرکتوری (Walk)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🚶 SECTION 4: DIRECTORY WALKING")
	fmt.Println(strings.Repeat("=", 80))

	// 4.1 WalkDir - پیمایش دایرکتوری (Go 1.16+)
	fmt.Println("\n--- 4.1 filepath.WalkDir ---")
	// WalkDir walks the file tree rooted at root, calling fn for each file or directory.

	// شبیه‌سازی ساختار دایرکتوری
	fmt.Println("  Example directory structure:")
	fmt.Println("    .")
	fmt.Println("    ├── main.go")
	fmt.Println("    ├── utils")
	fmt.Println("    │   ├── helper.go")
	fmt.Println("    │   └── validator.go")
	fmt.Println("    └── tests")
	fmt.Println("        └── main_test.go")

	// در عمل، WalkDir به این صورت استفاده می‌شود:
	/*
		err := filepath.WalkDir(".", func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			fmt.Printf("  %s (isDir: %v)\n", path, d.IsDir())
			return nil
		})
	*/

	fmt.Println("  (WalkDir traverses directory tree recursively)")

	// ============================================================================
	// بخش 5: توابع کاربردی دیگر
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🔧 SECTION 5: OTHER USEFUL FUNCTIONS")
	fmt.Println(strings.Repeat("=", 80))

	// 5.1 EvalSymlinks - ارزیابی symlink
	fmt.Println("\n--- 5.1 filepath.EvalSymlinks ---")
	// EvalSymlinks returns the path name after the evaluation of any symbolic links.
	// (در عمل برای دنبال کردن symlink‌ها استفاده می‌شود)
	fmt.Println("  EvalSymlinks resolves symbolic links to their target")

	// 5.2 HasPrefix - بررسی پیشوند (در path و filepath وجود ندارد)
	// می‌توانیم خودمان بسازیم
	fmt.Println("\n--- 5.2 Path Prefix Check ---")
	hasPrefix := func(p, prefix string) bool {
		// ساده‌ترین روش
		return strings.HasPrefix(p, prefix)
	}
	fmt.Printf("  HasPrefix(\"/a/b/c\", \"/a\") = %v\n", hasPrefix("/a/b/c", "/a"))
	fmt.Printf("  HasPrefix(\"/a/b/c\", \"/b\") = %v\n", hasPrefix("/a/b/c", "/b"))

	// ============================================================================
	// بخش 6: کاربردهای عملی
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("💡 SECTION 6: PRACTICAL EXAMPLES")
	fmt.Println(strings.Repeat("=", 80))

	// 6.1 استخراج نام فایل بدون پسوند
	fmt.Println("\n--- 6.1 Extract Filename Without Extension ---")
	filename := "document.pdf"
	nameWithoutExt := strings.TrimSuffix(filename, filepath.Ext(filename))
	fmt.Printf("  Original: %q, Without extension: %q\n", filename, nameWithoutExt)

	// 6.2 تغییر پسوند فایل
	fmt.Println("\n--- 6.2 Change File Extension ---")
	changeExt := func(path, newExt string) string {
		ext := filepath.Ext(path)
		return strings.TrimSuffix(path, ext) + newExt
	}
	fmt.Printf("  changeExt(\"image.jpg\", \".png\") = %q\n", changeExt("image.jpg", ".png"))
	fmt.Printf("  changeExt(\"data.tar.gz\", \".tgz\") = %q\n", changeExt("data.tar.gz", ".tgz"))

	// 6.3 ساخت مسیر امن
	fmt.Println("\n--- 6.3 Safe Path Construction ---")
	baseDir := "/var/log"
	userInput := "../../etc/passwd"

	// هرگز از concatenation ساده استفاده نکنید!
	unsafePath := baseDir + "/" + userInput
	fmt.Printf("  Unsafe: %q (path traversal vulnerability)\n", unsafePath)

	// استفاده از Clean برای نرمال‌سازی
	cleanedPath := filepath.Clean(unsafePath)
	fmt.Printf("  Cleaned: %q\n", cleanedPath)

	// برای امنیت بیشتر، بررسی کنید که مسیر در baseDir باشد
	safePath := filepath.Join(baseDir, userInput)
	fmt.Printf("  Join: %q\n", safePath)

	// 6.4 استخراج دایرکتوری والد
	fmt.Println("\n--- 6.4 Get Parent Directory ---")
	parent := func(p string) string {
		return filepath.Dir(p)
	}
	fmt.Printf("  Parent of \"/a/b/c/d.txt\" = %q\n", parent("/a/b/c/d.txt"))
	fmt.Printf("  Parent of \"/a/b/c/\" = %q\n", parent("/a/b/c/"))

	// 6.5 بررسی نوع فایل با پسوند
	fmt.Println("\n--- 6.5 File Type by Extension ---")
	fileType := func(filename string) string {
		switch filepath.Ext(filename) {
		case ".go":
			return "Go source file"
		case ".txt", ".md":
			return "Text file"
		case ".jpg", ".jpeg", ".png", ".gif":
			return "Image file"
		case ".mp3", ".wav":
			return "Audio file"
		case ".mp4", ".avi":
			return "Video file"
		default:
			return "Unknown type"
		}
	}
	fmt.Printf("  main.go: %s\n", fileType("main.go"))
	fmt.Printf("  image.jpg: %s\n", fileType("image.jpg"))
	fmt.Printf("  song.mp3: %s\n", fileType("song.mp3"))

	// 6.6 ساخت مسیر با تاریخ
	fmt.Println("\n--- 6.6 Path with Timestamp ---")
	// در عمل از time.Now() استفاده می‌شود
	logFile := filepath.Join("logs", "app-2024-01-15.log")
	fmt.Printf("  Log path: %q\n", logFile)

	// ============================================================================
	// بخش 7: تفاوت path و filepath
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("⚖️ SECTION 7: path vs filepath COMPARISON")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ Feature              │ path           │ path/filepath           │")
	fmt.Println("├──────────────────────┼────────────────┼─────────────────────────┤")
	fmt.Println("│ Separator            │ Always /       │ OS-specific (/, \\)     │")
	fmt.Println("│ Use case             │ URL paths,     │ File system paths       │")
	fmt.Println("│                      │ Unix paths     │                         │")
	fmt.Println("│ Platform independent │ Yes            │ Yes (adapts to OS)      │")
	fmt.Println("│ Volume names         │ No             │ Yes (C:, D:, etc.)      │")
	fmt.Println("│ Walk/Glob            │ No             │ Yes                     │")
	fmt.Println("│ Absolute path check  │ Simple         │ Full OS-aware          │")
	fmt.Println("│ Relative path        │ No             │ Yes (Rel function)      │")
	fmt.Println("│ Symlink resolution   │ No             │ Yes (EvalSymlinks)      │")
	fmt.Println("└──────────────────────┴────────────────┴─────────────────────────┘")

	// ============================================================================
	// بخش 8: اشتباهات رایج
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("❌ SECTION 8: COMMON MISTAKES")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println("\n❌ Mistake 1: Using path for file system paths")
	fmt.Println("   path.Join(\"home\", \"user\", \"file.txt\")  // uses /")
	fmt.Println("   ✅ Use filepath.Join for file system paths")

	fmt.Println("\n❌ Mistake 2: Manual path concatenation")
	fmt.Println("   dir + \"/\" + file  // wrong separator on Windows")
	fmt.Println("   ✅ filepath.Join(dir, file)")

	fmt.Println("\n❌ Mistake 3: Assuming path separators")
	fmt.Println("   strings.Split(path, \"/\")  // fails on Windows")
	fmt.Println("   ✅ filepath.SplitList or filepath.ToSlash first")

	fmt.Println("\n❌ Mistake 4: Not cleaning user input paths")
	fmt.Println("   filepath.Join(base, userInput)  // may contain ../")
	fmt.Println("   ✅ Clean after join or validate")

	fmt.Println("\n❌ Mistake 5: Using filepath.Ext for multiple extensions")
	fmt.Println("   filepath.Ext(\"file.tar.gz\") = \".gz\"  // only last")
	fmt.Println("   ✅ Handle multiple extensions manually")

	fmt.Println("\n❌ Mistake 6: Assuming Abs returns error for relative paths")
	fmt.Println("   abs, _ := filepath.Abs(\"file.txt\")  // works")
	fmt.Println("   ✅ But check error for absolute paths")

	fmt.Println("\n❌ Mistake 7: Using Glob without checking error")
	fmt.Println("   matches, _ := filepath.Glob(pattern)  // ignore error")
	fmt.Println("   ✅ Always check error for permission issues")

	// ============================================================================
	// بخش 9: جمع‌بندی و جدول مرجع
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📚 SECTION 9: QUICK REFERENCE")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ path FUNCTIONS                                                │")
	fmt.Println("├─────────────────────────────────────────────────────────────────┤")
	fmt.Println("│ path.Clean(p)      - Normalize path                           │")
	fmt.Println("│ path.Join(elems...) - Join path elements                      │")
	fmt.Println("│ path.Split(p)      - Split into dir and file                  │")
	fmt.Println("│ path.Dir(p)        - Get directory portion                    │")
	fmt.Println("│ path.Base(p)       - Get last element (file name)             │")
	fmt.Println("│ path.Ext(p)        - Get file extension                       │")
	fmt.Println("│ path.IsAbs(p)      - Check if absolute                        │")
	fmt.Println("│ path.Match(p, name) - Match pattern                           │")
	fmt.Println("└─────────────────────────────────────────────────────────────────┘")

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ filepath FUNCTIONS                                             │")
	fmt.Println("├─────────────────────────────────────────────────────────────────┤")
	fmt.Println("│ filepath.Clean(p)      - Normalize path                        │")
	fmt.Println("│ filepath.Join(elems...) - Join path elements (OS-aware)        │")
	fmt.Println("│ filepath.Split(p)      - Split into dir and file               │")
	fmt.Println("│ filepath.Dir(p)        - Get directory portion                 │")
	fmt.Println("│ filepath.Base(p)       - Get last element (file name)          │")
	fmt.Println("│ filepath.Ext(p)        - Get file extension                    │")
	fmt.Println("│ filepath.IsAbs(p)      - Check if absolute                     │")
	fmt.Println("│ filepath.Rel(base, targ) - Get relative path                   │")
	fmt.Println("│ filepath.Abs(p)        - Get absolute path                     │")
	fmt.Println("│ filepath.VolumeName(p) - Get volume name (Windows)             │")
	fmt.Println("│ filepath.ToSlash(p)    - Convert to slash (/)                  │")
	fmt.Println("│ filepath.FromSlash(p)  - Convert from slash to OS separator    │")
	fmt.Println("│ filepath.Match(p, name) - Match pattern                        │")
	fmt.Println("│ filepath.Glob(p)       - Find matching files                   │")
	fmt.Println("│ filepath.WalkDir(root, fn) - Walk directory tree               │")
	fmt.Println("│ filepath.EvalSymlinks(p) - Resolve symbolic links              │")
	fmt.Println("└─────────────────────────────────────────────────────────────────┘")

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ GLOB PATTERNS                                                 │")
	fmt.Println("├─────────────────────────────────────────────────────────────────┤")
	fmt.Println("│ *        - matches any sequence of non-separator characters   │")
	fmt.Println("│ ?        - matches any single non-separator character         │")
	fmt.Println("│ [abc]    - matches any character in the set                    │")
	fmt.Println("│ [a-z]    - matches any character in the range                  │")
	fmt.Println("│ {a,b}    - matches any of the alternatives (not in standard)  │")
	fmt.Println("│ \\        - escape character                                    │")
	fmt.Println("└─────────────────────────────────────────────────────────────────┘")

	fmt.Println("\n💡 GOLDEN RULES:")
	fmt.Println("  1. Use path for URL/Unix paths, filepath for file system paths")
	fmt.Println("  2. Always use filepath.Join for constructing file paths")
	fmt.Println("  3. Never concatenate paths with + or fmt.Sprintf")
	fmt.Println("  4. Clean user input paths to prevent path traversal")
	fmt.Println("  5. Use filepath.ToSlash for consistent path representation")
	fmt.Println("  6. Check errors from Glob, Abs, and Rel functions")
	fmt.Println("  7. Use filepath.VolumeName for Windows compatibility")
	fmt.Println("  8. Use filepath.WalkDir for directory traversal")
	fmt.Println("  9. Be careful with filepath.Ext for double extensions")
	fmt.Println("  10. Use filepath.Rel for relative path calculations")

	fmt.Println("\n🎯 PLATFORM-SPECIFIC BEHAVIOR:")
	fmt.Println("  • Windows: Separator is \\, volume names (C:), UNC paths (\\\\server\\share)")
	fmt.Println("  • Unix/Linux: Separator is /, no volume names")
	fmt.Println("  • macOS: Same as Unix")
	fmt.Println("  • filepath handles these differences automatically")
}

// تابع کمکی برای تکرار رشته
func stringsRepeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}
