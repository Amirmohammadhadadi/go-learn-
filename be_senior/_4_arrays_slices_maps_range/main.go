// ============================================================================
// FILE: arrays_slices_maps_range_guide.go
// TITLE: راهنمای کامل آرایه، اسلایس، مپ و range در Go
// HOW TO RUN: go run arrays_slices_maps_range_guide.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - تفاوت آرایه، اسلایس و مپ
// ============================================================================
//
// آرایه (Array):
// - طول ثابت (تعیین شده در زمان کامپایل)
// - مقدارها در حافظه پشت سر هم (contiguous)
// - value type (وقتی به تابع می‌دهید، کل آرایه کپی می‌شود)
// - کاربرد: زمانی که اندازه مشخص و ثابت است
//
// اسلایس (Slice):
// - طول پویا (می‌تواند رشد کند)
// - ویو (view) روی یک آرایه زیرین
// - reference type (اشاره‌گر به آرایه زیرین)
// - کاربرد: 99% مواقع از اسلایس استفاده می‌کنیم
//
// مپ (Map):
// - جدول هش (key-value)
// - unordered collection
// - reference type
// - کاربرد: جستجوی سریع، داده‌های کلید-مقدار
//
// قانون طلایی:
// "در Go، تقریباً همیشه از اسلایس استفاده کن، نه آرایه.
//  مپ برای جستجو و گروه‌بندی، range برای پیمایش"
// ============================================================================

package main

import (
	"fmt"
	"sort"
)

// ============================================================================
// بخش 1: آرایه‌ها (Arrays) - طول ثابت، value type
// ============================================================================

func demonstrateArrays() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📊 ARRAYS - Fixed Length, Value Type")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 1.1 تعریف و مقداردهی آرایه
	// ============================================
	fmt.Println("\n--- 1.1 Array Declaration and Initialization ---")

	// روش 1: با اندازه و مقداردهی
	var arr1 [5]int
	fmt.Printf("arr1: %v (zero values)\n", arr1)

	// روش 2: با مقداردهی اولیه
	arr2 := [5]int{1, 2, 3, 4, 5}
	fmt.Printf("arr2: %v\n", arr2)

	// روش 3: با سه نقطه (اندازه از روی مقدارها مشخص می‌شود)
	arr3 := [...]int{10, 20, 30, 40}
	fmt.Printf("arr3: %v, length: %d\n", arr3, len(arr3))

	// روش 4: مقداردهی به اندیس‌های خاص
	arr4 := [5]int{0: 100, 2: 200, 4: 300}
	fmt.Printf("arr4: %v\n", arr4)

	// آرایه دو بعدی
	matrix := [3][3]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}
	fmt.Printf("Matrix: %v\n", matrix)

	// ============================================
	// 1.2 دسترسی و تغییر عناصر
	// ============================================
	fmt.Println("\n--- 1.2 Access and Modify ---")

	nums := [5]int{10, 20, 30, 40, 50}
	fmt.Printf("Original: %v\n", nums)

	// دسترسی با اندیس
	fmt.Printf("nums[0]: %d, nums[2]: %d\n", nums[0], nums[2])

	// تغییر مقدار
	nums[1] = 200
	nums[4] = 500
	fmt.Printf("After modification: %v\n", nums)

	// len() و cap() برای آرایه
	fmt.Printf("Length: %d, Capacity: %d\n", len(nums), cap(nums))

	// ============================================
	// 1.3 آرایه value type (کپی می‌شود)
	// ============================================
	fmt.Println("\n--- 1.3 Array is Value Type (Copied) ---")

	original := [3]int{1, 2, 3}
	copy := original // کپی کامل

	copy[0] = 999

	fmt.Printf("Original: %v\n", original)
	fmt.Printf("Copy: %v\n", copy)

	// تابعی که آرایه می‌گیرد (کپی می‌شود)
	modifyArray := func(arr [3]int) {
		arr[0] = 1000
	}

	testArr := [3]int{1, 2, 3}
	modifyArray(testArr)
	fmt.Printf("After function call: %v (unchanged)\n", testArr)

	// ============================================
	// 1.4 اشاره‌گر به آرایه (اگر می‌خواهید تغییر کند)
	// ============================================
	fmt.Println("\n--- 1.4 Pointer to Array ---")

	modifyArrayPtr := func(arr *[3]int) {
		arr[0] = 1000 // Go allows this without (*arr)[0]
	}

	testArr2 := [3]int{1, 2, 3}
	modifyArrayPtr(&testArr2)
	fmt.Printf("After pointer function: %v (changed!)\n", testArr2)
}

// ============================================================================
// بخش 2: اسلایس‌ها (Slices) - قلب کار با داده‌های دنباله‌ای در Go
// ============================================================================

func demonstrateSlices() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🔪 SLICES - Dynamic Arrays, Heart of Go Collections")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 2.1 ایجاد اسلایس
	// ============================================
	fmt.Println("\n--- 2.1 Creating Slices ---")

	// روش 1: با make (نوع، طول، ظرفیت)
	slice1 := make([]int, 5)     // طول 5، ظرفیت 5
	slice2 := make([]int, 5, 10) // طول 5، ظرفیت 10
	fmt.Printf("slice1: len=%d, cap=%d, val=%v\n", len(slice1), cap(slice1), slice1)
	fmt.Printf("slice2: len=%d, cap=%d, val=%v\n", len(slice2), cap(slice2), slice2)

	// روش 2: با literal
	slice3 := []int{1, 2, 3, 4, 5}
	fmt.Printf("slice3: %v\n", slice3)

	// روش 3: از آرایه
	arr := [5]int{10, 20, 30, 40, 50}
	slice4 := arr[1:4] // از اندیس 1 تا 3 (قبل از 4)
	fmt.Printf("slice4 (arr[1:4]): %v\n", slice4)

	// روش 4: از اسلایس دیگر
	slice5 := slice3[1:3]
	fmt.Printf("slice5 (slice3[1:3]): %v\n", slice5)

	// روش 5: اسلایس خالی
	var emptySlice []int
	fmt.Printf("emptySlice: %v, len=%d, cap=%d, is nil? %v\n",
		emptySlice, len(emptySlice), cap(emptySlice), emptySlice == nil)

	// ============================================
	// 2.2 append - افزودن عنصر به اسلایس
	// ============================================
	fmt.Println("\n--- 2.2 Append Operation ---")

	nums := []int{1, 2, 3}
	fmt.Printf("Initial: %v (len=%d, cap=%d)\n", nums, len(nums), cap(nums))

	// append یک عنصر
	nums = append(nums, 4)
	fmt.Printf("After append 4: %v (len=%d, cap=%d)\n", nums, len(nums), cap(nums))

	// append چند عنصر
	nums = append(nums, 5, 6, 7)
	fmt.Printf("After append 5,6,7: %v (len=%d, cap=%d)\n", nums, len(nums), cap(nums))

	// append یک اسلایس دیگر (با ...)
	more := []int{8, 9, 10}
	nums = append(nums, more...)
	fmt.Printf("After append slice: %v (len=%d, cap=%d)\n", nums, len(nums), cap(nums))

	// ============================================
	// 2.3 برش (slicing) - ایجاد ویو روی آرایه زیرین
	// ============================================
	fmt.Println("\n--- 2.3 Slicing (Views) ---")

	original := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// برش‌های مختلف
	slice10 := original[2:7] // اندیس 2 تا 6
	slice11 := original[:5]  // 0 تا 4
	slice12 := original[5:]  // 5 تا آخر
	slice13 := original[:]   // کل اسلایس

	fmt.Printf("original: %v\n", original)
	fmt.Printf("original[2:7]: %v\n", slice10)
	fmt.Printf("original[:5]: %v\n", slice11)
	fmt.Printf("original[5:]: %v\n", slice12)
	fmt.Printf("original[:]: %v\n", slice13)

	// تغییر در برش، آرایه زیرین را تغییر می‌دهد
	slice10[0] = 999
	fmt.Printf("\nAfter modifying slice[2:7]: original = %v\n", original)

	// ============================================
	// 2.4 copy - کپی کردن اسلایس
	// ============================================
	fmt.Println("\n--- 2.4 Copy Operation ---")

	src := []int{1, 2, 3, 4, 5}
	dst := make([]int, len(src))

	copied := copy(dst, src)
	fmt.Printf("src: %v\n", src)
	fmt.Printf("dst: %v (copied %d elements)\n", dst, copied)

	// کپی با اندازه‌های مختلف
	src2 := []int{1, 2, 3, 4, 5, 6}
	dst2 := make([]int, 3)
	copy(dst2, src2)
	fmt.Printf("src2: %v, dst2: %v (only first 3 copied)\n", src2, dst2)

	// کپی به وسط اسلایس
	dst3 := []int{0, 0, 0, 0, 0}
	copy(dst3[2:], src2[:3])
	fmt.Printf("dst3 after copy to middle: %v\n", dst3)

	// ============================================
	// 2.5 ظرفیت (capacity) و رشد اسلایس
	// ============================================
	fmt.Println("\n--- 2.5 Capacity and Growth ---")

	// ایجاد با ظرفیت اولیه
	sl := make([]int, 0, 3)
	fmt.Printf("Initial: len=%d, cap=%d\n", len(sl), cap(sl))

	// افزودن تا ظرفیت
	for i := 1; i <= 5; i++ {
		sl = append(sl, i)
		fmt.Printf("After append %d: len=%d, cap=%d\n", i, len(sl), cap(sl))
	}

	// نکته: وقتی ظرفیت پر می‌شود، Go ظرفیت را دو برابر می‌کند

	// ============================================
	// 2.6 اسلایس multidimentional
	// ============================================
	fmt.Println("\n--- 2.6 Multidimensional Slices ---")

	// ماتریس 3x3 با اسلایس
	matrix := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}

	fmt.Println("Matrix:")
	for i, row := range matrix {
		fmt.Printf("  Row %d: %v\n", i, row)
	}

	// افزودن سطر جدید
	matrix = append(matrix, []int{10, 11, 12})
	fmt.Printf("After adding row: %v\n", matrix)
}

// ============================================================================
// بخش 3: مپ‌ها (Maps) - کلید-مقدار، جستجوی سریع
// ============================================================================

func demonstrateMaps() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🗺️ MAPS - Key-Value Stores, Fast Lookup")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 3.1 ایجاد و مقداردهی مپ
	// ============================================
	fmt.Println("\n--- 3.1 Creating Maps ---")

	// روش 1: با make
	m1 := make(map[string]int)
	fmt.Printf("m1: %v, len=%d\n", m1, len(m1))

	// روش 2: با literal
	m2 := map[string]int{
		"apple":  5,
		"banana": 3,
		"cherry": 7,
	}
	fmt.Printf("m2: %v, len=%d\n", m2, len(m2))

	// روش 3: مپ خالی (nil map)
	var m3 map[string]int
	fmt.Printf("m3: %v, is nil? %v\n", m3, m3 == nil)

	// ⚠️ نمی‌توان به m3 nil اضافه کرد (panic)
	// m3["key"] = 1  // panic!

	// ✅ باید با make ایجاد کرد
	m4 := make(map[string]int)
	m4["key"] = 1 // OK

	// ============================================
	// 3.2 عملیات اصلی روی مپ
	// ============================================
	fmt.Println("\n--- 3.2 Basic Operations ---")

	scores := make(map[string]int)

	// افزودن/به‌روزرسانی
	scores["Ali"] = 95
	scores["Reza"] = 87
	scores["Sara"] = 92
	fmt.Printf("After adding: %v\n", scores)

	// خواندن
	fmt.Printf("Ali's score: %d\n", scores["Ali"])

	// بررسی وجود کلید (comma ok idiom)
	val, ok := scores["Mohammad"]
	if ok {
		fmt.Printf("Mohammad's score: %d\n", val)
	} else {
		fmt.Println("Mohammad not found")
	}

	// حذف
	delete(scores, "Reza")
	fmt.Printf("After deleting Reza: %v\n", scores)

	// طول مپ
	fmt.Printf("Length: %d\n", len(scores))

	// ============================================
	// 3.3 مپ با انواع مختلف کلید
	// ============================================
	fmt.Println("\n--- 3.3 Different Key Types ---")

	// کلید از نوع int
	intMap := map[int]string{
		1: "one",
		2: "two",
		3: "three",
	}
	fmt.Printf("int map: %v\n", intMap)

	// کلید از نوع float64 (کمی risky)
	floatMap := map[float64]string{
		1.0: "one",
		2.0: "two",
	}
	fmt.Printf("float map: %v\n", floatMap)

	// کلید از نوع struct
	type Point struct {
		X, Y int
	}
	pointMap := map[Point]string{
		{X: 0, Y: 0}: "origin",
		{X: 1, Y: 0}: "x-axis",
	}
	fmt.Printf("struct map: %v\n", pointMap)

	// کلید از نوع array (نه slice)
	arrayMap := map[[3]int]string{
		{1, 2, 3}: "first",
		{4, 5, 6}: "second",
	}
	fmt.Printf("array map: %v\n", arrayMap)

	// ============================================
	// 3.4 مپ با مقدار از نوع اسلایس یا مپ
	// ============================================
	fmt.Println("\n--- 3.4 Maps with Slice/Map Values ---")

	// مپ به اسلایس
	students := make(map[string][]string)
	students["Ali"] = []string{"Math", "Physics"}
	students["Sara"] = []string{"English", "Art", "Music"}
	fmt.Printf("Student courses: %v\n", students)

	// افزودن به اسلایس
	students["Ali"] = append(students["Ali"], "Chemistry")
	fmt.Printf("After adding course: %v\n", students)

	// مپ به مپ (nested map)
	grades := make(map[string]map[string]int)
	grades["Ali"] = make(map[string]int)
	grades["Ali"]["Math"] = 95
	grades["Ali"]["Physics"] = 88
	fmt.Printf("Nested map: %v\n", grades)

	// ============================================
	// 3.5 پیمایش مپ با range
	// ============================================
	fmt.Println("\n--- 3.5 Iterating Over Maps ---")

	colors := map[string]string{
		"red":   "#FF0000",
		"green": "#00FF00",
		"blue":  "#0000FF",
	}

	fmt.Println("All colors:")
	for key, value := range colors {
		fmt.Printf("  %s -> %s\n", key, value)
	}

	// فقط کلیدها
	fmt.Println("\nKeys only:")
	for key := range colors {
		fmt.Printf("  %s\n", key)
	}

	// فقط مقادیر
	fmt.Println("\nValues only:")
	for _, value := range colors {
		fmt.Printf("  %s\n", value)
	}

	// ⚠️ توجه: ترتیب پیمایش مپ تصادفی است
	fmt.Println("\n⚠️ Map iteration order is random!")

	// ============================================
	// 3.6 مرتب‌سازی کلیدهای مپ
	// ============================================
	fmt.Println("\n--- 3.6 Sorting Map Keys ---")

	ages := map[string]int{
		"Reza":  30,
		"Ali":   25,
		"Sara":  28,
		"Zahra": 35,
	}

	// گرفتن کلیدها
	keys := make([]string, 0, len(ages))
	for k := range ages {
		keys = append(keys, k)
	}

	// مرتب‌سازی
	sort.Strings(keys)

	// پیمایش مرتب
	fmt.Println("Sorted by key:")
	for _, k := range keys {
		fmt.Printf("  %s: %d\n", k, ages[k])
	}
}

// ============================================================================
// بخش 4: Range - پیمایش آرایه، اسلایس، مپ، رشته
// ============================================================================

func demonstrateRange() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🔄 RANGE - Iterating Over Arrays, Slices, Maps, Strings")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 4.1 range روی اسلایس
	// ============================================
	fmt.Println("\n--- 4.1 Range Over Slice ---")

	nums := []int{10, 20, 30, 40, 50}

	// با اندیس و مقدار
	fmt.Println("With index and value:")
	for i, v := range nums {
		fmt.Printf("  index=%d, value=%d\n", i, v)
	}

	// فقط مقدار (با ignore index)
	fmt.Println("\nOnly values:")
	for _, v := range nums {
		fmt.Printf("  %d\n", v)
	}

	// فقط اندیس
	fmt.Println("\nOnly indices:")
	for i := range nums {
		fmt.Printf("  %d\n", i)
	}

	// ============================================
	// 4.2 range روی آرایه
	// ============================================
	fmt.Println("\n--- 4.2 Range Over Array ---")

	arr := [5]int{100, 200, 300, 400, 500}

	for i, v := range arr {
		fmt.Printf("  arr[%d] = %d\n", i, v)
	}

	// ============================================
	// 4.3 range روی مپ
	// ============================================
	fmt.Println("\n--- 4.3 Range Over Map ---")

	population := map[string]int{
		"Tehran":  8000000,
		"Mashhad": 3000000,
		"Isfahan": 2000000,
	}

	for city, pop := range population {
		fmt.Printf("  %s: %d\n", city, pop)
	}

	// ============================================
	// 4.4 range روی رشته (Unicode-aware)
	// ============================================
	fmt.Println("\n--- 4.4 Range Over String (Unicode) ---")

	text := "Hello, 世界"

	fmt.Println("Bytes (not recommended for Unicode):")
	for i := 0; i < len(text); i++ {
		fmt.Printf("  %d: %c (%d)\n", i, text[i], text[i])
	}

	fmt.Println("\nRunes (correct for Unicode):")
	for i, r := range text {
		fmt.Printf("  %d: %c (U+%04X)\n", i, r, r)
	}

	// مثال با فارسی
	persian := "سلام دنیا"
	fmt.Println("\nPersian string:")
	for i, r := range persian {
		fmt.Printf("  %d: %c (U+%04X)\n", i, r, r)
	}

	// ============================================
	// 4.5 range روی کانال
	// ============================================
	fmt.Println("\n--- 4.5 Range Over Channel ---")

	ch := make(chan int, 5)
	go func() {
		for i := 1; i <= 5; i++ {
			ch <- i
		}
		close(ch)
	}()

	for v := range ch {
		fmt.Printf("  Received from channel: %d\n", v)
	}
}

// ============================================================================
// بخش 5: توابع کاربردی برای اسلایس و مپ
// ============================================================================

// تابع حذف عنصر از اسلایس (با حفظ ترتیب)
func removeElement(slice []int, index int) []int {
	if index < 0 || index >= len(slice) {
		return slice
	}
	return append(slice[:index], slice[index+1:]...)
}

// تابع حذف عنصر بدون حفظ ترتیب (سریعتر)
func removeElementFast(slice []int, index int) []int {
	if index < 0 || index >= len(slice) {
		return slice
	}
	slice[index] = slice[len(slice)-1]
	return slice[:len(slice)-1]
}

// تابع فیلتر اسلایس
func filter(slice []int, predicate func(int) bool) []int {
	result := make([]int, 0)
	for _, v := range slice {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}

// تابع مپ (transform)
func mapSlice(slice []int, transform func(int) int) []int {
	result := make([]int, len(slice))
	for i, v := range slice {
		result[i] = transform(v)
	}
	return result
}

func demonstrateUtilityFunctions() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🛠️ UTILITY FUNCTIONS for Slices and Maps")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 5.1 حذف عنصر از اسلایس
	// ============================================
	fmt.Println("\n--- 5.1 Remove Element from Slice ---")

	nums := []int{10, 20, 30, 40, 50}
	fmt.Printf("Original: %v\n", nums)

	nums = removeElement(nums, 2)
	fmt.Printf("After removing index 2: %v\n", nums)

	// حذف سریع
	nums2 := []int{1, 2, 3, 4, 5}
	nums2 = removeElementFast(nums2, 2)
	fmt.Printf("Fast remove index 2: %v\n", nums2)

	// ============================================
	// 5.2 فیلتر اسلایس
	// ============================================
	fmt.Println("\n--- 5.2 Filter Slice ---")

	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// فیلتر اعداد زوج
	evens := filter(numbers, func(n int) bool {
		return n%2 == 0
	})
	fmt.Printf("Evens: %v\n", evens)

	// فیلتر اعداد بزرگتر از 5
	greater := filter(numbers, func(n int) bool {
		return n > 5
	})
	fmt.Printf("Greater than 5: %v\n", greater)

	// ============================================
	// 5.3 تبدیل (map) اسلایس
	// ============================================
	fmt.Println("\n--- 5.3 Map (Transform) Slice ---")

	nums3 := []int{1, 2, 3, 4, 5}

	// ضرب در 2
	doubled := mapSlice(nums3, func(n int) int {
		return n * 2
	})
	fmt.Printf("Doubled: %v\n", doubled)

	// به توان 2
	squared := mapSlice(nums3, func(n int) int {
		return n * n
	})
	fmt.Printf("Squared: %v\n", squared)
}

// ============================================================================
// بخش 6: تله‌ها و اشتباهات رایج
// ============================================================================

func demonstrateCommonMistakes() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("⚠️ COMMON MISTAKES WITH ARRAYS, SLICES, MAPS")
	fmt.Println(stringsRepeat("=", 80))

	// اشتباه 1: append بدون اختصاص دادن نتیجه
	fmt.Println("\n❌ Mistake 1: append without assignment")
	fmt.Println("   s := []int{1,2,3}")
	fmt.Println("   append(s, 4)  // ❌ result ignored!")
	fmt.Println("   ✅ s = append(s, 4)")

	// اشتباه 2: اشتراک آرایه زیرین در اسلایس‌ها
	fmt.Println("\n❌ Mistake 2: Shared underlying array")
	fmt.Println("   s1 := []int{1,2,3,4,5}")
	fmt.Println("   s2 := s1[1:4]")
	fmt.Println("   s2[0] = 999  // ❌ modifies s1 too!")
	fmt.Println("   ✅ Use copy() for independent slice")

	// اشتباه 3: استفاده از مپ nil
	fmt.Println("\n❌ Mistake 3: Using nil map")
	fmt.Println("   var m map[string]int")
	fmt.Println("   m[\"key\"] = 1  // ❌ panic!")
	fmt.Println("   ✅ m := make(map[string]int)")

	// اشتباه 4: تغییر مپ در حین پیمایش
	fmt.Println("\n❌ Mistake 4: Modifying map during iteration")
	fmt.Println("   for k, v := range m {")
	fmt.Println("       delete(m, k)  // ❌ unpredictable!")
	fmt.Println("   }")
	fmt.Println("   ✅ Collect keys first, then delete")

	// اشتباه 5: فرض ترتیب در مپ
	fmt.Println("\n❌ Mistake 5: Assuming map order")
	fmt.Println("   Map iteration order is random!")
	fmt.Println("   ✅ Sort keys if you need order")

	// اشتباه 6: فراموش کردن ظرفیت در make
	fmt.Println("\n❌ Mistake 6: Forgetting capacity in make")
	fmt.Println("   s := make([]int, 0)  // cap=0")
	fmt.Println("   s = append(s, 1)    // allocates new array")
	fmt.Println("   ✅ s := make([]int, 0, 100)  // preallocate")
}

// ============================================================================
// بخش 7: جمع‌بندی و جدول مرجع سریع
// ============================================================================

func main() {
	fmt.Println(stringsRepeat("=", 80))
	fmt.Println("🎯 COMPLETE GUIDE: ARRAYS | SLICES | MAPS | RANGE")
	fmt.Println(stringsRepeat("=", 80))

	// بخش 1: آرایه‌ها
	demonstrateArrays()

	// بخش 2: اسلایس‌ها
	demonstrateSlices()

	// بخش 3: مپ‌ها
	demonstrateMaps()

	// بخش 4: Range
	demonstrateRange()

	// بخش 5: توابع کاربردی
	demonstrateUtilityFunctions()

	// بخش 6: اشتباهات رایج
	demonstrateCommonMistakes()

	// بخش 7: جمع‌بندی
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📚 QUICK REFERENCE CARD")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ ARRAYS (Fixed Length)                                           │")
	fmt.Println("├─────────────────────────────────────────────────────────────────┤")
	fmt.Println("│ var arr [5]int                     // declaration              │")
	fmt.Println("│ arr := [5]int{1,2,3,4,5}          // initialization            │")
	fmt.Println("│ arr := [...]int{1,2,3}             // length inferred           │")
	fmt.Println("│ arr[0] = 10                        // set value                 │")
	fmt.Println("│ x := arr[2]                        // get value                 │")
	fmt.Println("│ len(arr), cap(arr)                 // length and capacity       │")
	fmt.Println("└─────────────────────────────────────────────────────────────────┘")

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ SLICES (Dynamic Arrays)                                         │")
	fmt.Println("├─────────────────────────────────────────────────────────────────┤")
	fmt.Println("│ s := make([]int, 5)                // len=5, cap=5             │")
	fmt.Println("│ s := make([]int, 5, 10)            // len=5, cap=10            │")
	fmt.Println("│ s := []int{1,2,3}                  // literal                   │")
	fmt.Println("│ s = append(s, 4)                   // add element               │")
	fmt.Println("│ s = append(s, 5,6,7)               // add multiple              │")
	fmt.Println("│ s = append(s, other...)            // append slice              │")
	fmt.Println("│ sub := s[1:4]                      // slicing (view)            │")
	fmt.Println("│ copy(dst, src)                     // copy slice                │")
	fmt.Println("└─────────────────────────────────────────────────────────────────┘")

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ MAPS (Key-Value)                                                │")
	fmt.Println("├─────────────────────────────────────────────────────────────────┤")
	fmt.Println("│ m := make(map[string]int)          // create                    │")
	fmt.Println("│ m := map[string]int{\"a\":1}         // literal                   │")
	fmt.Println("│ m[\"key\"] = 42                      // set value                 │")
	fmt.Println("│ x := m[\"key\"]                      // get value                 │")
	fmt.Println("│ x, ok := m[\"key\"]                  // check existence           │")
	fmt.Println("│ delete(m, \"key\")                   // delete key                │")
	fmt.Println("│ len(m)                             // number of keys            │")
	fmt.Println("└─────────────────────────────────────────────────────────────────┘")

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ RANGE (Iteration)                                               │")
	fmt.Println("├─────────────────────────────────────────────────────────────────┤")
	fmt.Println("│ for i, v := range slice            // index, value              │")
	fmt.Println("│ for i := range slice               // index only                │")
	fmt.Println("│ for _, v := range slice            // value only                │")
	fmt.Println("│ for k, v := range map              // key, value                │")
	fmt.Println("│ for i, r := range string           // rune index, rune value    │")
	fmt.Println("│ for v := range channel             // values from channel       │")
	fmt.Println("└─────────────────────────────────────────────────────────────────┘")

	fmt.Println("\n💡 GOLDEN RULES:")
	fmt.Println("  1. Use slices, not arrays (99% of the time)")
	fmt.Println("  2. Always assign result of append back to slice")
	fmt.Println("  3. Use make() to preallocate capacity for performance")
	fmt.Println("  4. Be careful with shared underlying arrays in slices")
	fmt.Println("  5. Check 'ok' when reading from maps")
	fmt.Println("  6. Never assume order when iterating maps")
	fmt.Println("  7. Use copy() for independent slice copies")
	fmt.Println("  8. Delete map keys during iteration is safe in Go 1.21+")
	fmt.Println("  9. Use range for strings (Unicode safe), not index")
	fmt.Println("  10. When in doubt, use slice + map combination")

	fmt.Println("\n🎯 PERFORMANCE TIPS:")
	fmt.Println("  • Preallocate slices with make() when size is known")
	fmt.Println("  • Use copy() for large slices instead of append")
	fmt.Println("  • Avoid unnecessary slice copying")
	fmt.Println("  • For maps, approximate size improves performance")
	fmt.Println("  • Use slice of structs instead of slice of pointers")
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
