// ============================================================================
// FILE: slices_complete_guide.go
// TITLE: راهنمای کامل پکیج slices در Go - تمام توابع با مثال
// HOW TO RUN: go run slices_complete_guide.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - پکیج slices چیست؟
// ============================================================================
//
// پکیج slices (اضافه شده در Go 1.21) توابع عمومی برای کار با اسلایس‌های هر نوعی ارائه می‌دهد.
// با استفاده از Generics، این توابع برای همه انواع اسلایس قابل استفاده هستند.
//
// قانون طلایی:
// "هر زمان که نیاز به عملیات روی اسلایس داری، اول ببین آیا تابعی در پکیج slices هست.
//  این پکیج بیشتر نیازهای روزمره را پوشش می‌دهد و کد را خواناتر می‌کند."
// ============================================================================

package _5_internal_packages_slice_map

import (
	"cmp"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

func main() {
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("🎯 COMPLETE slices PACKAGE GUIDE IN GO")
	fmt.Println("All functions with practical examples")
	fmt.Println(strings.Repeat("=", 80))

	// ============================================================================
	// بخش 1: توابع جستجو (Search)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🔍 SECTION 1: SEARCH FUNCTIONS")
	fmt.Println(strings.Repeat("=", 80))

	// 1.1 BinarySearch - جستجوی دودویی در اسلایس مرتب شده
	fmt.Println("\n--- 1.1 slices.BinarySearch ---")
	// BinarySearch searches for target in a sorted slice and returns the earliest position
	// where target is found, or the position where target would appear in the sort order;
	// it also returns a bool saying whether the target is really found in the slice [citation:1].
	numbers := []int{10, 20, 30, 40, 50}

	idx, found := slices.BinarySearch(numbers, 30)
	fmt.Printf("  BinarySearch([10,20,30,40,50], 30): index=%d, found=%v\n", idx, found)

	idx, found = slices.BinarySearch(numbers, 35)
	fmt.Printf("  BinarySearch([10,20,30,40,50], 35): index=%d, found=%v\n", idx, found)

	// 1.2 BinarySearchFunc - جستجوی دودویی با تابع مقایسه سفارشی
	fmt.Println("\n--- 1.2 slices.BinarySearchFunc ---")
	// BinarySearchFunc works like BinarySearch, but uses a custom comparison function [citation:1].
	type Person struct {
		Name string
		Age  int
	}
	people := []Person{
		{"Alice", 25},
		{"Bob", 30},
		{"Charlie", 35},
	}

	idx, found = slices.BinarySearchFunc(people, Person{Name: "Bob", Age: 0},
		func(a, b Person) int {
			return strings.Compare(a.Name, b.Name)
		})
	fmt.Printf("  BinarySearchFunc by name 'Bob': index=%d, found=%v\n", idx, found)

	// 1.3 Contains - بررسی وجود عنصر در اسلایس
	fmt.Println("\n--- 1.3 slices.Contains ---")
	// Contains reports whether v is present in s [citation:1].
	colors := []string{"red", "green", "blue"}
	fmt.Printf("  Contains([red,green,blue], green): %v\n", slices.Contains(colors, "green"))
	fmt.Printf("  Contains([red,green,blue], yellow): %v\n", slices.Contains(colors, "yellow"))

	// 1.4 ContainsFunc - بررسی وجود عنصر با تابع شرط
	fmt.Println("\n--- 1.4 slices.ContainsFunc ---")
	// ContainsFunc reports whether at least one element e of s satisfies f(e) [citation:1].
	hasEven := slices.ContainsFunc([]int{1, 3, 5, 7, 8}, func(n int) bool {
		return n%2 == 0
	})
	fmt.Printf("  ContainsFunc (has even): %v\n", hasEven)

	// 1.5 Index - پیدا کردن اولین موقعیت عنصر
	fmt.Println("\n--- 1.5 slices.Index ---")
	// Index returns the index of the first occurrence of v in s, or -1 if not present [citation:1].
	fruits := []string{"apple", "banana", "cherry", "banana"}
	fmt.Printf("  Index([apple,banana,cherry,banana], banana): %d\n", slices.Index(fruits, "banana"))
	fmt.Printf("  Index([apple,banana,cherry,banana], orange): %d\n", slices.Index(fruits, "orange"))

	// 1.6 IndexFunc - پیدا کردن اولین موقعیت با تابع شرط
	fmt.Println("\n--- 1.6 slices.IndexFunc ---")
	// IndexFunc returns the first index i satisfying f(s[i]), or -1 if none do [citation:1].
	idx = slices.IndexFunc([]int{10, 20, 30, 35, 40}, func(n int) bool {
		return n > 30
	})
	fmt.Printf("  IndexFunc (first > 30): %d\n", idx)

	// ============================================================================
	// بخش 2: توابع مقایسه (Comparison)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("⚖️ SECTION 2: COMPARISON FUNCTIONS")
	fmt.Println(strings.Repeat("=", 80))

	// 2.1 Equal - بررسی برابری دو اسلایس
	fmt.Println("\n--- 2.1 slices.Equal ---")
	// Equal reports whether two slices are equal: the same length and all elements equal [citation:1].
	a := []int{1, 2, 3}
	b := []int{1, 2, 3}
	c := []int{1, 2, 4}
	fmt.Printf("  Equal([1,2,3], [1,2,3]): %v\n", slices.Equal(a, b))
	fmt.Printf("  Equal([1,2,3], [1,2,4]): %v\n", slices.Equal(a, c))

	// 2.2 EqualFunc - بررسی برابری با تابع مقایسه سفارشی
	fmt.Println("\n--- 2.2 slices.EqualFunc ---")
	// EqualFunc reports whether two slices are equal using an equality function on each pair of elements [citation:1].
	nums := []int{1, 2, 3}
	strs := []string{"1", "2", "3"}
	equal := slices.EqualFunc(nums, strs, func(n int, s string) bool {
		return strconv.Itoa(n) == s
	})
	fmt.Printf("  EqualFunc([1,2,3], [\"1\",\"2\",\"3\"]): %v\n", equal)

	// 2.3 Compare - مقایسه lexicographical دو اسلایس
	fmt.Println("\n--- 2.3 slices.Compare ---")
	// Compare compares the elements of s1 and s2, using cmp.Compare on each pair of elements [citation:1].
	// Returns -1 if s1 < s2, 0 if s1 == s2, +1 if s1 > s2
	fmt.Printf("  Compare([1,2,3], [1,2,3]): %d\n", slices.Compare([]int{1, 2, 3}, []int{1, 2, 3}))
	fmt.Printf("  Compare([1,2,3], [1,2,4]): %d\n", slices.Compare([]int{1, 2, 3}, []int{1, 2, 4}))
	fmt.Printf("  Compare([1,2,3], [1,2]): %d\n", slices.Compare([]int{1, 2, 3}, []int{1, 2}))

	// 2.4 CompareFunc - مقایسه با تابع سفارشی
	fmt.Println("\n--- 2.4 slices.CompareFunc ---")
	// CompareFunc is like Compare but uses a custom comparison function on each pair of elements [citation:1].
	result := slices.CompareFunc([]int{1, 2, 3}, []string{"1", "2", "3"},
		func(n int, s string) int {
			sn, _ := strconv.Atoi(s)
			return cmp.Compare(n, sn)
		})
	fmt.Printf("  CompareFunc result: %d\n", result)

	// 2.5 IsSorted - بررسی مرتب بودن اسلایس
	fmt.Println("\n--- 2.5 slices.IsSorted ---")
	// IsSorted reports whether x is sorted in ascending order [citation:1].
	sorted := []int{1, 2, 3, 4, 5}
	unsorted := []int{1, 3, 2, 4, 5}
	fmt.Printf("  IsSorted([1,2,3,4,5]): %v\n", slices.IsSorted(sorted))
	fmt.Printf("  IsSorted([1,3,2,4,5]): %v\n", slices.IsSorted(unsorted))

	// 2.6 IsSortedFunc - بررسی مرتب بودن با تابع سفارشی
	fmt.Println("\n--- 2.6 slices.IsSortedFunc ---")
	// IsSortedFunc reports whether x is sorted in ascending order, with cmp as the comparison function [citation:1].
	isSortedDesc := slices.IsSortedFunc([]int{5, 4, 3, 2, 1}, func(a, b int) int {
		return cmp.Compare(b, a) // descending order
	})
	fmt.Printf("  IsSortedFunc (descending): %v\n", isSortedDesc)

	// ============================================================================
	// بخش 3: توابع کپی و کلون (Copy & Clone)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📋 SECTION 3: COPY & CLONE FUNCTIONS")
	fmt.Println(strings.Repeat("=", 80))

	// 3.1 Clone - کپی سطحی (shallow copy) از اسلایس
	fmt.Println("\n--- 3.1 slices.Clone ---")
	// Clone returns a copy of the slice. The elements are copied using assignment, so this is a shallow clone [citation:1].
	original := []int{1, 2, 3, 4, 5}
	cloned := slices.Clone(original)
	cloned[0] = 100
	fmt.Printf("  Original: %v\n", original)
	fmt.Printf("  Cloned (modified): %v\n", cloned)
	fmt.Printf("  Note: Clone creates independent copy\n")

	// 3.2 Clip - حذف ظرفیت اضافی اسلایس
	fmt.Println("\n--- 3.2 slices.Clip ---")
	// Clip removes unused capacity from the slice, returning s[:len(s):len(s)] [citation:1].
	arr := [...]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	s := arr[:4:10]
	fmt.Printf("  Before Clip: len=%d, cap=%d, value=%v\n", len(s), cap(s), s)
	clipped := slices.Clip(s)
	fmt.Printf("  After Clip: len=%d, cap=%d, value=%v\n", len(clipped), cap(clipped), clipped)

	// 3.3 Grow - افزایش ظرفیت اسلایس
	fmt.Println("\n--- 3.3 slices.Grow ---")
	// Grow increases the slice's capacity, if necessary, to guarantee space for another n elements [citation:1].
	nums2 := []int{1, 2, 3}
	fmt.Printf("  Before Grow: len=%d, cap=%d\n", len(nums2), cap(nums2))
	nums2 = slices.Grow(nums2, 10)
	fmt.Printf("  After Grow: len=%d, cap=%d\n", len(nums2), cap(nums2))

	// ============================================================================
	// بخش 4: توابع حذف (Delete)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🗑️ SECTION 4: DELETE FUNCTIONS")
	fmt.Println(strings.Repeat("=", 80))

	// 4.1 Delete - حذف بازه‌ای از عناصر
	fmt.Println("\n--- 4.1 slices.Delete ---")
	// Delete removes the elements s[i:j] from s, returning the modified slice [citation:1].
	letters := []string{"a", "b", "c", "d", "e"}
	fmt.Printf("  Before Delete: %v\n", letters)
	letters = slices.Delete(letters, 1, 4)
	fmt.Printf("  After Delete indices 1-4: %v\n", letters)

	// 4.2 DeleteFunc - حذف عناصر بر اساس شرط
	fmt.Println("\n--- 4.2 slices.DeleteFunc ---")
	// DeleteFunc removes any elements from s for which del returns true, returning the modified slice [citation:1].
	numbers2 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	fmt.Printf("  Before DeleteFunc: %v\n", numbers2)
	numbers2 = slices.DeleteFunc(numbers2, func(n int) bool {
		return n%2 == 0 // حذف اعداد زوج
	})
	fmt.Printf("  After DeleteFunc (remove evens): %v\n", numbers2)

	// 4.3 Compact - حذف تکرارهای متوالی
	fmt.Println("\n--- 4.3 slices.Compact ---")
	// Compact replaces consecutive runs of equal elements with a single copy [citation:1].
	duplicates := []int{1, 1, 2, 2, 2, 3, 4, 4, 5}
	fmt.Printf("  Before Compact: %v\n", duplicates)
	duplicates = slices.Compact(duplicates)
	fmt.Printf("  After Compact: %v\n", duplicates)

	// 4.4 CompactFunc - حذف تکرارهای متوالی با تابع مقایسه
	fmt.Println("\n--- 4.4 slices.CompactFunc ---")
	// CompactFunc is like Compact but uses an equality function to compare elements [citation:1].
	names := []string{"bob", "Bob", "alice", "Vera", "VERA"}
	fmt.Printf("  Before CompactFunc: %v\n", names)
	names = slices.CompactFunc(names, strings.EqualFold)
	fmt.Printf("  After CompactFunc (case-insensitive): %v\n", names)

	// ============================================================================
	// بخش 5: توابع درج (Insert)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("➕ SECTION 5: INSERT FUNCTIONS")
	fmt.Println(strings.Repeat("=", 80))

	// 5.1 Insert - درج عناصر در موقعیت مشخص
	fmt.Println("\n--- 5.1 slices.Insert ---")
	// Insert inserts the values v... into s at index i, returning the modified slice [citation:1].
	fruits2 := []string{"apple", "orange"}
	fmt.Printf("  Before Insert: %v\n", fruits2)
	fruits2 = slices.Insert(fruits2, 1, "banana", "grape")
	fmt.Printf("  After Insert at index 1: %v\n", fruits2)

	// ============================================================================
	// بخش 6: توابع الحاق (Append & Concat)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🔗 SECTION 6: APPEND & CONCAT FUNCTIONS")
	fmt.Println(strings.Repeat("=", 80))

	// 6.1 Concat - اتصال چند اسلایس
	fmt.Println("\n--- 6.1 slices.Concat ---")
	// Concat returns a new slice concatenating the passed in slices [citation:1].
	s1 := []int{1, 2, 3}
	s2 := []int{4, 5, 6}
	s3 := []int{7, 8, 9}
	concatenated := slices.Concat(s1, s2, s3)
	fmt.Printf("  Concat([1,2,3], [4,5,6], [7,8,9]): %v\n", concatenated)

	// ============================================================================
	// بخش 7: توابع یافتن ماکزیمم و مینیمم (Max & Min)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📊 SECTION 7: MAX & MIN FUNCTIONS")
	fmt.Println(strings.Repeat("=", 80))

	// 7.1 Max - یافتن بزرگترین مقدار
	fmt.Println("\n--- 7.1 slices.Max ---")
	// Max returns the maximal value in x. It panics if x is empty [citation:1].
	values := []int{42, 17, 89, 3, 56}
	maxVal := slices.Max(values)
	fmt.Printf("  Max([42,17,89,3,56]): %d\n", maxVal)

	// 7.2 MaxFunc - یافتن بزرگترین مقدار با تابع مقایسه
	fmt.Println("\n--- 7.2 slices.MaxFunc ---")
	// MaxFunc returns the maximal value in x using a custom comparison function [citation:1].
	type Product struct {
		Name  string
		Price int
	}
	products := []Product{
		{"Laptop", 1000},
		{"Mouse", 20},
		{"Keyboard", 50},
		{"Monitor", 300},
	}
	maxPrice := slices.MaxFunc(products, func(a, b Product) int {
		return cmp.Compare(a.Price, b.Price)
	})
	fmt.Printf("  MaxFunc (by price): %+v\n", maxPrice)

	// 7.3 Min - یافتن کوچکترین مقدار
	fmt.Println("\n--- 7.3 slices.Min ---")
	// Min returns the minimal value in x. It panics if x is empty [citation:1].
	minVal := slices.Min(values)
	fmt.Printf("  Min([42,17,89,3,56]): %d\n", minVal)

	// 7.4 MinFunc - یافتن کوچکترین مقدار با تابع مقایسه
	fmt.Println("\n--- 7.4 slices.MinFunc ---")
	minPrice := slices.MinFunc(products, func(a, b Product) int {
		return cmp.Compare(a.Price, b.Price)
	})
	fmt.Printf("  MinFunc (by price): %+v\n", minPrice)

	// ============================================================================
	// بخش 8: توابع مرتب‌سازی (Sort)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📊 SECTION 8: SORT FUNCTIONS")
	fmt.Println(strings.Repeat("=", 80))

	// 8.1 Sort - مرتب‌سازی اسلایس
	fmt.Println("\n--- 8.1 slices.Sort ---")
	// Sort sorts a slice of any ordered type in increasing order [citation:1].
	unsortedNums := []int{5, 2, 8, 1, 9, 3}
	fmt.Printf("  Before Sort: %v\n", unsortedNums)
	slices.Sort(unsortedNums)
	fmt.Printf("  After Sort: %v\n", unsortedNums)

	// 8.2 SortFunc - مرتب‌سازی با تابع مقایسه سفارشی
	fmt.Println("\n--- 8.2 slices.SortFunc ---")
	// SortFunc sorts a slice using a custom comparison function [citation:1].
	people2 := []Person{
		{"Charlie", 35},
		{"Alice", 25},
		{"Bob", 30},
	}
	fmt.Printf("  Before SortFunc: %+v\n", people2)
	slices.SortFunc(people2, func(a, b Person) int {
		return cmp.Compare(a.Name, b.Name)
	})
	fmt.Printf("  After SortFunc (by name): %+v\n", people2)

	// 8.3 SortStableFunc - مرتب‌سازی پایدار با تابع مقایسه
	fmt.Println("\n--- 8.3 slices.SortStableFunc ---")
	// SortStableFunc sorts the slice x while keeping the original order of equal elements [citation:1].
	people3 := []Person{
		{"Alice", 30},
		{"Bob", 25},
		{"Alice", 20},
	}
	fmt.Printf("  Before SortStableFunc: %+v\n", people3)
	slices.SortStableFunc(people3, func(a, b Person) int {
		return cmp.Compare(a.Name, b.Name)
	})
	fmt.Printf("  After SortStableFunc: %+v (order of 'Alice' preserved)\n", people3)

	// ============================================================================
	// بخش 9: توابع معکوس‌سازی (Reverse)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🔄 SECTION 9: REVERSE FUNCTION")
	fmt.Println(strings.Repeat("=", 80))

	// 9.1 Reverse - معکوس‌سازی اسلایس درجا
	fmt.Println("\n--- 9.1 slices.Reverse ---")
	// Reverse reverses the elements of the slice in place [citation:1].
	toReverse := []int{1, 2, 3, 4, 5}
	fmt.Printf("  Before Reverse: %v\n", toReverse)
	slices.Reverse(toReverse)
	fmt.Printf("  After Reverse: %v\n", toReverse)

	// ============================================================================
	// بخش 10: توابع جایگزینی (Replace)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🔄 SECTION 10: REPLACE FUNCTION")
	fmt.Println(strings.Repeat("=", 80))

	// 10.1 Replace - جایگزینی بازه‌ای از عناصر
	fmt.Println("\n--- 10.1 slices.Replace ---")
	// Replace replaces the elements s[i:j] with the given v, returning the modified slice [citation:1].
	original2 := []string{"a", "b", "c", "d", "e"}
	fmt.Printf("  Before Replace: %v\n", original2)
	replaced := slices.Replace(original2, 1, 4, "x", "y")
	fmt.Printf("  After Replace indices 1-4 with [x,y]: %v\n", replaced)

	// ============================================================================
	// بخش 11: توابع تکراری (Iterators - Go 1.23+)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🔄 SECTION 11: ITERATOR FUNCTIONS (Go 1.23+)")
	fmt.Println(strings.Repeat("=", 80))

	// 11.1 All - بازگرداندن iterator روی جفت‌های (index, value)
	fmt.Println("\n--- 11.1 slices.All ---")
	// All returns an iterator over index-value pairs in the slice in the usual order [citation:1].
	colors2 := []string{"red", "green", "blue"}
	fmt.Print("  All iterator: ")
	for i, v := range slices.All(colors2) {
		fmt.Printf("(%d,%s) ", i, v)
	}
	fmt.Println()

	// 11.2 Backward - بازگرداندن iterator به صورت معکوس
	fmt.Println("\n--- 11.2 slices.Backward ---")
	// Backward returns an iterator over index-value pairs in the slice, traversing backward [citation:1].
	fmt.Print("  Backward iterator: ")
	for i, v := range slices.Backward(colors2) {
		fmt.Printf("(%d,%s) ", i, v)
	}
	fmt.Println()

	// 11.3 Chunk - تقسیم اسلایس به تکه‌های n تایی
	fmt.Println("\n--- 11.3 slices.Chunk ---")
	// Chunk returns an iterator over consecutive sub-slices of up to n elements of s [citation:1].
	nums3 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	fmt.Print("  Chunk (size 3): ")
	for chunk := range slices.Chunk(nums3, 3) {
		fmt.Printf("%v ", chunk)
	}
	fmt.Println()

	// ============================================================================
	// بخش 12: مثال‌های کاربردی (Practical Examples)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("💡 SECTION 12: PRACTICAL EXAMPLES")
	fmt.Println(strings.Repeat("=", 80))

	// 12.1 حذف تکراری‌ها با استفاده از Compact
	fmt.Println("\n--- 12.1 Remove Duplicates ---")
	duplicateStrings := []string{"apple", "banana", "apple", "cherry", "banana", "apple"}
	fmt.Printf("  Original: %v\n", duplicateStrings)

	// ابتدا مرتب می‌کنیم، سپس Compact می‌کنیم
	slices.Sort(duplicateStrings)
	unique := slices.Compact(duplicateStrings)
	fmt.Printf("  Unique: %v\n", unique)

	// 12.2 فیلتر کردن (Filter) - پیاده‌سازی دستی چون در پکیج slices وجود ندارد
	fmt.Println("\n--- 12.2 Filter Pattern ---")
	// Note: slices package does NOT have Filter function (as of Go 1.23) [citation:9]
	numbers3 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// پیاده‌سازی Filter
	filterEven := func(n int) bool { return n%2 == 0 }
	filtered := make([]int, 0)
	for _, v := range numbers3 {
		if filterEven(v) {
			filtered = append(filtered, v)
		}
	}
	fmt.Printf("  Original: %v\n", numbers3)
	fmt.Printf("  Filtered (evens): %v\n", filtered)

	// 12.3 Map - پیاده‌سازی دستی (در پکیج slices وجود ندارد)
	fmt.Println("\n--- 12.3 Map Pattern ---")
	// Map pattern (not in standard slices package) [citation:7]
	square := func(n int) int { return n * n }
	mapped := make([]int, len(numbers3))
	for i, v := range numbers3 {
		mapped[i] = square(v)
	}
	fmt.Printf("  Original: %v\n", numbers3)
	fmt.Printf("  Mapped (squares): %v\n", mapped)

	// 12.4 بررسی اسلایس خالی
	fmt.Println("\n--- 12.4 Empty Slice Check ---")
	var emptySlice []int
	nilSlice := []int{}
	nonEmpty := []int{1, 2, 3}

	fmt.Printf("  emptySlice == nil: %v\n", emptySlice == nil)
	fmt.Printf("  len(emptySlice) == 0: %v\n", len(emptySlice) == 0)
	fmt.Printf("  len(nilSlice) == 0: %v\n", len(nilSlice) == 0)
	fmt.Printf("  slices.Equal(emptySlice, nilSlice): %v\n", slices.Equal(emptySlice, nilSlice))

	// 12.5 تبدیل اسلایس به مپ
	fmt.Println("\n--- 12.5 Slice to Map Conversion ---")
	keys := []string{"a", "b", "c"}
	values2 := []int{1, 2, 3}

	resultMap := make(map[string]int)
	for i, k := range keys {
		if i < len(values2) {
			resultMap[k] = values2[i]
		}
	}
	fmt.Printf("  Keys: %v\n", keys)
	fmt.Printf("  Values: %v\n", values2)
	fmt.Printf("  Map: %v\n", resultMap)

	// ============================================================================
	// بخش 13: نکات عملکردی (Performance Tips)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("⚡ SECTION 13: PERFORMANCE TIPS")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println("\n💡 Important notes about slice functions:")
	fmt.Println("  1. Delete/Insert/Replace functions return a NEW slice - always assign the result")
	fmt.Println("  2. Sort/Reverse modify the slice in place - no return value needed")
	fmt.Println("  3. Compact and DeleteFunc zero out obsolete elements to prevent memory leaks [citation:4]")
	fmt.Println("  4. Clone creates a shallow copy - for deep copy, implement manually")
	fmt.Println("  5. Use Grow when you know the final size to reduce allocations")
	fmt.Println("  6. BinarySearch only works on SORTED slices - check IsSorted first")

	// ============================================================================
	// بخش 14: اشتباهات رایج (Common Mistakes)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("❌ SECTION 14: COMMON MISTAKES")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println("\n❌ Mistake 1: Ignoring return value of Delete/Insert")
	fmt.Println("   slices.Delete(s, i, j)  // WRONG - s is not modified")
	fmt.Println("   ✅ s = slices.Delete(s, i, j)")

	fmt.Println("\n❌ Mistake 2: Using BinarySearch on unsorted slice")
	fmt.Println("   slices.BinarySearch(unsorted, 5)  // WRONG - undefined behavior")
	fmt.Println("   ✅ slices.Sort(s); slices.BinarySearch(s, 5)")

	fmt.Println("\n❌ Mistake 3: Assuming nil and empty slices are different in Equal")
	fmt.Println("   slices.Equal(nil, []int{}) returns true")
	fmt.Println("   ✅ They are considered equal - check nil separately if needed")

	fmt.Println("\n❌ Mistake 4: Using Max/Min on empty slice (panics)")
	fmt.Println("   slices.Max([]int{})  // PANIC!")
	fmt.Println("   ✅ Check len(s) > 0 before calling Max/Min")

	fmt.Println("\n❌ Mistake 5: Forgetting that Compact removes ONLY consecutive duplicates")
	fmt.Println("   slices.Compact([]int{1,2,1,2})  // returns [1,2,1,2] unchanged")
	fmt.Println("   ✅ Sort first: slices.Sort(s); slices.Compact(s)")

	// ============================================================================
	// بخش 15: جمع‌بندی و جدول مرجع سریع
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📚 SECTION 15: QUICK REFERENCE")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ CATEGORY          │ FUNCTIONS                                    │")
	fmt.Println("├───────────────────┼──────────────────────────────────────────────┤")
	fmt.Println("│ Search            │ BinarySearch, BinarySearchFunc, Contains,    │")
	fmt.Println("│                   │ ContainsFunc, Index, IndexFunc               │")
	fmt.Println("│ Comparison        │ Equal, EqualFunc, Compare, CompareFunc,      │")
	fmt.Println("│                   │ IsSorted, IsSortedFunc                       │")
	fmt.Println("│ Copy/Clone        │ Clone, Clip, Grow                            │")
	fmt.Println("│ Delete            │ Delete, DeleteFunc, Compact, CompactFunc     │")
	fmt.Println("│ Insert/Replace    │ Insert, Replace                              │")
	fmt.Println("│ Append/Concat     │ Concat                                       │")
	fmt.Println("│ Min/Max           │ Min, Max, MinFunc, MaxFunc                   │")
	fmt.Println("│ Sort              │ Sort, SortFunc, SortStableFunc, Reverse      │")
	fmt.Println("│ Iterators (1.23+) │ All, Backward, Chunk                         │")
	fmt.Println("└───────────────────┴──────────────────────────────────────────────┘")

	fmt.Println("\n💡 GOLDEN RULES:")
	fmt.Println("  1. Always assign result of functions that may change slice length")
	fmt.Println("  2. Sort before using BinarySearch")
	fmt.Println("  3. Sort before Compact to remove all duplicates (not just consecutive)")
	fmt.Println("  4. Check length before calling Max/Min (they panic on empty slices)")
	fmt.Println("  5. nil and empty slices are considered equal by Equal function")
	fmt.Println("  6. Use Clone for independent copy, not just assignment")
	fmt.Println("  7. Use Grow to pre-allocate when final size is known")
	fmt.Println("  8. Delete/Insert/Replace modify underlying array - be careful with shared slices")
	fmt.Println("  9. Compact/DeleteFunc zero out tail elements to prevent memory leaks [citation:4]")
	fmt.Println("  10. For Filter/Map operations, implement manually (not in standard library) [citation:7]")
}

// تابع کمکی برای تکرار رشته (در خود پکیج وجود دارد، اینجا فقط برای مثال)
// func stringsRepeat(s string, count int) string {
//     return strings.Repeat(s, count)
// }
