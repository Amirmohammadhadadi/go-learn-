package _5_internal_packages_slice_map

// ============================================================================
// FILE: maps_complete_guide.go
// TITLE: راهنمای کامل پکیج maps در Go - تمام توابع با مثال
// HOW TO RUN: go run maps_complete_guide.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - پکیج maps چیست؟
// ============================================================================
//
// پکیج maps (اضافه شده در Go 1.21) توابع عمومی برای کار با مپ‌های هر نوعی ارائه می‌دهد.
// با استفاده از Generics، این توابع برای همه انواع مپ قابل استفاده هستند.
//
// قانون طلایی:
// "هر زمان که نیاز به عملیات روی مپ داری، اول ببین آیا تابعی در پکیج maps هست.
//  این پکیج بیشتر نیازهای روزمره مانند کپی، حذف شرطی، و مقایسه را پوشش می‌دهد."
// ============================================================================

import (
	"fmt"
	"maps"
	"slices"
	"strings"
)

func main() {
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("🎯 COMPLETE maps PACKAGE GUIDE IN GO")
	fmt.Println("All functions with practical examples")
	fmt.Println(strings.Repeat("=", 80))

	// ============================================================================
	// بخش 1: توابع کپی و کلون (Copy & Clone)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📋 SECTION 1: COPY & CLONE FUNCTIONS")
	fmt.Println(strings.Repeat("=", 80))

	// 1.1 Clone - کپی عمیق (shallow copy) از مپ
	fmt.Println("\n--- 1.1 maps.Clone ---")
	// Clone returns a copy of m. This is a shallow clone: the new map contains the same keys
	// and values as the original map (using assignment).
	original := map[string]int{
		"apple":  5,
		"banana": 3,
		"cherry": 7,
	}

	cloned := maps.Clone(original)
	cloned["apple"] = 100

	fmt.Printf("  Original map: %v\n", original)
	fmt.Printf("  Cloned map (modified): %v\n", cloned)
	fmt.Printf("  Note: Clone creates independent copy (shallow)\n")

	// 1.2 Copy - کپی کردن تمام key-value از مپ مبدأ به مقصد
	fmt.Println("\n--- 1.2 maps.Copy ---")
	// Copy copies all key/value pairs from src to dst, overwriting any existing keys in dst.
	dst := map[string]int{
		"apple": 1,
		"grape": 10,
	}
	src := map[string]int{
		"apple":  5,
		"banana": 3,
		"cherry": 7,
	}

	fmt.Printf("  Before Copy - dst: %v\n", dst)
	fmt.Printf("  Before Copy - src: %v\n", src)

	maps.Copy(dst, src)

	fmt.Printf("  After Copy - dst: %v (overwrote apple, added banana, cherry)\n", dst)
	fmt.Printf("  After Copy - src: %v (unchanged)\n", src)

	// ============================================================================
	// بخش 2: توابع حذف (Delete)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🗑️ SECTION 2: DELETE FUNCTIONS")
	fmt.Println(strings.Repeat("=", 80))

	// 2.1 DeleteFunc - حذف عناصر بر اساس شرط
	fmt.Println("\n--- 2.1 maps.DeleteFunc ---")
	// DeleteFunc deletes any key/value pairs from m for which del returns true.
	scores := map[string]int{
		"Alice":   95,
		"Bob":     45,
		"Charlie": 85,
		"Diana":   30,
		"Eve":     60,
	}

	fmt.Printf("  Before DeleteFunc: %v\n", scores)

	// حذف افرادی که نمره کمتر از 50 دارند
	maps.DeleteFunc(scores, func(key string, value int) bool {
		return value < 50
	})

	fmt.Printf("  After DeleteFunc (remove scores < 50): %v\n", scores)

	// ============================================================================
	// بخش 3: توابع مقایسه (Comparison)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("⚖️ SECTION 3: COMPARISON FUNCTIONS")
	fmt.Println(strings.Repeat("=", 80))

	// 3.1 Equal - بررسی برابری دو مپ
	fmt.Println("\n--- 3.1 maps.Equal ---")
	// Equal reports whether two maps contain the same key/value pairs.
	map1 := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}
	map2 := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}
	map3 := map[string]int{
		"a": 1,
		"b": 2,
		"d": 4,
	}

	fmt.Printf("  Equal(map1, map2): %v\n", maps.Equal(map1, map2))
	fmt.Printf("  Equal(map1, map3): %v\n", maps.Equal(map1, map3))

	// 3.2 EqualFunc - بررسی برابری با تابع مقایسه سفارشی
	fmt.Println("\n--- 3.2 maps.EqualFunc ---")
	// EqualFunc is like Equal, but compares values using eq. Keys are still compared normally.
	type Person struct {
		Name string
		Age  int
	}

	mapA := map[int]Person{
		1: {"Alice", 25},
		2: {"Bob", 30},
	}
	mapB := map[int]Person{
		1: {"ALICE", 25},
		2: {"BOB", 30},
	}

	equal := maps.EqualFunc(mapA, mapB, func(p1, p2 Person) bool {
		return p1.Age == p2.Age && strings.EqualFold(p1.Name, p2.Name)
	})

	fmt.Printf("  EqualFunc (case-insensitive name match): %v\n", equal)

	// ============================================================================
	// بخش 4: کار با کلیدها و مقادیر (Keys & Values) - استفاده از slices
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🔑 SECTION 4: KEYS & VALUES (using slices package)")
	fmt.Println(strings.Repeat("=", 80))

	// 4.1 دریافت کلیدهای مپ (با استفاده از slices.Collect و maps.Keys)
	fmt.Println("\n--- 4.1 Getting Keys ---")
	// Note: maps.Keys returns an iterator, use slices.Collect to get slice
	colors := map[string]string{
		"red":   "#FF0000",
		"green": "#00FF00",
		"blue":  "#0000FF",
	}

	keys := slices.Collect(maps.Keys(colors))
	fmt.Printf("  Keys: %v\n", keys)

	// 4.2 دریافت مقادیر مپ
	fmt.Println("\n--- 4.2 Getting Values ---")
	values := slices.Collect(maps.Values(colors))
	fmt.Printf("  Values: %v\n", values)

	// 4.3 مرتب‌سازی کلیدها
	fmt.Println("\n--- 4.3 Sorted Keys ---")
	sortedKeys := slices.Collect(maps.Keys(colors))
	slices.Sort(sortedKeys)
	fmt.Printf("  Sorted keys: %v\n", sortedKeys)

	// چاپ مپ به صورت مرتب
	fmt.Println("  Map in order:")
	for _, k := range sortedKeys {
		fmt.Printf("    %s -> %s\n", k, colors[k])
	}

	// ============================================================================
	// بخش 5: توابع ساخته نشده در پکیج maps (که باید خودمان بسازیم)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🔧 SECTION 5: CUSTOM FUNCTIONS (Not in standard maps package)")
	fmt.Println(strings.Repeat("=", 80))

	// 5.1 Filter - فیلتر کردن مپ بر اساس شرط (نیاز به ساخت دستی)
	fmt.Println("\n--- 5.1 Filter (Custom Implementation) ---")

	filterMap := func(m map[string]int, predicate func(string, int) bool) map[string]int {
		result := make(map[string]int)
		for k, v := range m {
			if predicate(k, v) {
				result[k] = v
			}
		}
		return result
	}

	ages := map[string]int{
		"Alice":   25,
		"Bob":     17,
		"Charlie": 30,
		"Diana":   16,
		"Eve":     22,
	}

	// فیلتر افراد بالای 18 سال
	adults := filterMap(ages, func(name string, age int) bool {
		return age >= 18
	})

	fmt.Printf("  Original ages: %v\n", ages)
	fmt.Printf("  Adults (age >= 18): %v\n", adults)

	// 5.2 Map (Transform) - تبدیل مقادیر مپ
	fmt.Println("\n--- 5.2 Map (Transform Values) ---")

	transformValues := func(m map[string]int, transform func(int) int) map[string]int {
		result := make(map[string]int)
		for k, v := range m {
			result[k] = transform(v)
		}
		return result
	}

	// افزایش 5 سال به همه
	aged := transformValues(ages, func(age int) int {
		return age + 5
	})

	fmt.Printf("  Original: %v\n", ages)
	fmt.Printf("  After +5 years: %v\n", aged)

	// 5.3 Merge - ادغام چند مپ
	fmt.Println("\n--- 5.3 Merge Multiple Maps ---")

	mergeMaps := func(maps ...map[string]int) map[string]int {
		result := make(map[string]int)
		for _, m := range maps {
			for k, v := range m {
				result[k] = v
			}
		}
		return result
	}

	mapA := map[string]int{"a": 1, "b": 2}
	mapB := map[string]int{"c": 3, "d": 4}
	mapC := map[string]int{"e": 5, "f": 6}

	merged := mergeMaps(mapA, mapB, mapC)
	fmt.Printf("  Merged: %v\n", merged)

	// 5.4 Invert - معکوس کردن مپ (key و value جابجا می‌شوند)
	fmt.Println("\n--- 5.4 Invert Map (Swap Key/Value) ---")

	invertMap := func(m map[string]int) map[int]string {
		result := make(map[int]string)
		for k, v := range m {
			result[v] = k
		}
		return result
	}

	original2 := map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	inverted := invertMap(original2)
	fmt.Printf("  Original: %v\n", original2)
	fmt.Printf("  Inverted: %v\n", inverted)

	// 5.5 KeyExists - بررسی وجود کلید با مقدار پیش‌فرض
	fmt.Println("\n--- 5.5 KeyExists with Default Value ---")

	getValue := func(m map[string]int, key string, defaultValue int) int {
		if value, ok := m[key]; ok {
			return value
		}
		return defaultValue
	}

	// استفاده از تابع
	value := getValue(ages, "Frank", 0)
	fmt.Printf("  Value for 'Frank' (not exists): %d\n", value)
	value = getValue(ages, "Alice", 0)
	fmt.Printf("  Value for 'Alice' (exists): %d\n", value)

	// 5.6 Intersection - اشتراک دو مپ
	fmt.Println("\n--- 5.6 Map Intersection ---")

	intersection := func(m1, m2 map[string]int) map[string]int {
		result := make(map[string]int)
		for k, v := range m1 {
			if val2, ok := m2[k]; ok && val2 == v {
				result[k] = v
			}
		}
		return result
	}

	m1 := map[string]int{"a": 1, "b": 2, "c": 3}
	m2 := map[string]int{"b": 2, "c": 4, "d": 5}

	common := intersection(m1, m2)
	fmt.Printf("  Map1: %v\n", m1)
	fmt.Printf("  Map2: %v\n", m2)
	fmt.Printf("  Intersection: %v\n", common)

	// ============================================================================
	// بخش 6: مثال‌های کاربردی (Practical Examples)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("💡 SECTION 6: PRACTICAL EXAMPLES")
	fmt.Println(strings.Repeat("=", 80))

	// 6.1 شمارش کلمات (Word Count)
	fmt.Println("\n--- 6.1 Word Count ---")

	text := "go is great and go is powerful and go is simple"
	words := strings.Fields(text)

	wordCount := make(map[string]int)
	for _, word := range words {
		wordCount[word]++
	}

	fmt.Printf("  Text: %s\n", text)
	fmt.Printf("  Word count: %v\n", wordCount)

	// 6.2 گروه‌بندی داده‌ها
	fmt.Println("\n--- 6.2 Grouping Data ---")

	type Person2 struct {
		Name string
		City string
	}

	people := []Person2{
		{"Alice", "Tehran"},
		{"Bob", "Shiraz"},
		{"Charlie", "Tehran"},
		{"Diana", "Isfahan"},
		{"Eve", "Shiraz"},
	}

	// گروه‌بندی بر اساس شهر
	cityGroups := make(map[string][]string)
	for _, p := range people {
		cityGroups[p.City] = append(cityGroups[p.City], p.Name)
	}

	fmt.Println("  People grouped by city:")
	for city, names := range cityGroups {
		fmt.Printf("    %s: %v\n", city, names)
	}

	// 6.3 Cache با زمان انقضا (ساده)
	fmt.Println("\n--- 6.3 Simple Cache with TTL ---")

	type CacheItem struct {
		Value      string
		Expiration int64
	}

	// در یک مثال واقعی از time.Now().Unix() استفاده می‌شود
	// cache := make(map[string]CacheItem)

	fmt.Println("  (Cache implementation would use time package for TTL)")

	// 6.4 تبدیل اسلایس به مپ
	fmt.Println("\n--- 6.4 Slice to Map Conversion ---")

	ids := []int{101, 102, 103, 104}
	names2 := []string{"Alice", "Bob", "Charlie", "Diana"}

	idToName := make(map[int]string)
	for i, id := range ids {
		if i < len(names2) {
			idToName[id] = names2[i]
		}
	}

	fmt.Printf("  IDs: %v\n", ids)
	fmt.Printf("  Names: %v\n", names2)
	fmt.Printf("  ID to Name map: %v\n", idToName)

	// 6.5 معکوس کردن مپ یک به یک
	fmt.Println("\n--- 6.5 Inverse of One-to-One Map ---")

	countryCode := map[string]string{
		"Iran":    "IR",
		"USA":     "US",
		"UK":      "GB",
		"Germany": "DE",
	}

	codeToCountry := make(map[string]string)
	for country, code := range countryCode {
		codeToCountry[code] = country
	}

	fmt.Printf("  Country to Code: %v\n", countryCode)
	fmt.Printf("  Code to Country: %v\n", codeToCountry)

	// 6.6 محاسبه جمع مقادیر
	fmt.Println("\n--- 6.6 Sum of Values ---")

	prices := map[string]float64{
		"apple":  0.5,
		"banana": 0.3,
		"cherry": 0.8,
		"date":   1.2,
	}

	var total float64
	for _, price := range prices {
		total += price
	}

	fmt.Printf("  Prices: %v\n", prices)
	fmt.Printf("  Total: %.2f\n", total)

	// ============================================================================
	// بخش 7: نکات عملکردی (Performance Tips)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("⚡ SECTION 7: PERFORMANCE TIPS")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println("\n💡 Important notes about map operations:")
	fmt.Println("  1. Pre-allocate map capacity with make(map[K]V, size) when size is known")
	fmt.Println("  2. maps.Clone creates a new map with same capacity as original")
	fmt.Println("  3. maps.Copy modifies the destination map in place - no new allocation")
	fmt.Println("  4. For large maps, DeleteFunc iterates over entire map O(n)")
	fmt.Println("  5. Equal and EqualFunc also iterate over entire map O(n)")
	fmt.Println("  6. Use range directly for simple operations (fastest)")
	fmt.Println("  7. maps.Keys and maps.Values return iterators - use slices.Collect to get slice")

	// 7.1 مثال پیش‌تخصیص ظرفیت
	fmt.Println("\n--- 7.1 Capacity Pre-allocation ---")

	size := 1000
	// بدون پیش‌تخصیص (slower, more allocations)
	withoutCap := make(map[int]int)
	for i := 0; i < size; i++ {
		withoutCap[i] = i
	}

	// با پیش‌تخصیص (faster, fewer allocations)
	withCap := make(map[int]int, size)
	for i := 0; i < size; i++ {
		withCap[i] = i
	}

	fmt.Printf("  Pre-allocating capacity reduces rehashing overhead\n")

	// ============================================================================
	// بخش 8: اشتباهات رایج (Common Mistakes)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("❌ SECTION 8: COMMON MISTAKES")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println("\n❌ Mistake 1: Reading from nil map")
	fmt.Println("   var m map[string]int")
	fmt.Println("   val := m[\"key\"]  // OK - returns zero value")
	fmt.Println("   m[\"key\"] = 1     // PANIC! assignment to nil map")
	fmt.Println("   ✅ Use make(map[string]int) before writing")

	fmt.Println("\n❌ Mistake 2: Assuming map iteration order")
	fmt.Println("   for k, v := range m { ... }  // order is random")
	fmt.Println("   ✅ Sort keys first if order matters")

	fmt.Println("\n❌ Mistake 3: Modifying map while iterating")
	fmt.Println("   for k, v := range m {")
	fmt.Println("       delete(m, k)  // safe in Go")
	fmt.Println("       m[\"new\"] = 1  // may work but not guaranteed")
	fmt.Println("   }")
	fmt.Println("   ✅ Collect keys first for modification")

	fmt.Println("\n❌ Mistake 4: Using map as a set for non-comparable types")
	fmt.Println("   type Point struct { X, Y int }")
	fmt.Println("   set := make(map[Point]bool)  // works because Point is comparable")
	fmt.Println("   // For slices: use map[string]bool with fmt.Sprint as key")

	fmt.Println("\n❌ Mistake 5: maps.Equal vs reflect.DeepEqual")
	fmt.Println("   maps.Equal only compares map contents, not additional fields")
	fmt.Println("   ✅ Use maps.Equal for maps, reflect.DeepEqual for complex structs")

	fmt.Println("\n❌ Mistake 6: Copying map in loop (inefficient)")
	fmt.Println("   for k, v := range src { dst[k] = v }")
	fmt.Println("   ✅ Use maps.Copy(dst, src) - optimized")

	fmt.Println("\n❌ Mistake 7: maps.Clone on map with pointer values")
	fmt.Println("   cloned := maps.Clone(original)  // shallow copy only")
	fmt.Println("   ✅ For deep copy, implement manually or use encoding/gob")

	// ============================================================================
	// بخش 9: توابع ساخته نشده که مفید هستند
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🔧 SECTION 9: USEFUL CUSTOM FUNCTIONS (Not in standard maps)")
	fmt.Println(strings.Repeat("=", 80))

	// 9.1 MapDifference - تفاوت دو مپ
	fmt.Println("\n--- 9.1 Map Difference ---")

	mapDifference := func(m1, m2 map[string]int) (onlyInM1, onlyInM2, common map[string]int) {
		onlyInM1 = make(map[string]int)
		onlyInM2 = make(map[string]int)
		common = make(map[string]int)

		for k, v := range m1 {
			if val2, ok := m2[k]; ok {
				if val2 == v {
					common[k] = v
				} else {
					onlyInM1[k] = v
				}
			} else {
				onlyInM1[k] = v
			}
		}

		for k, v := range m2 {
			if _, ok := m1[k]; !ok {
				onlyInM2[k] = v
			}
		}

		return onlyInM1, onlyInM2, common
	}

	mA := map[string]int{"a": 1, "b": 2, "c": 3}
	mB := map[string]int{"b": 2, "c": 4, "d": 5}

	onlyA, onlyB, common2 := mapDifference(mA, mB)

	fmt.Printf("  Map A: %v\n", mA)
	fmt.Printf("  Map B: %v\n", mB)
	fmt.Printf("  Only in A: %v\n", onlyA)
	fmt.Printf("  Only in B: %v\n", onlyB)
	fmt.Printf("  Common: %v\n", common2)

	// 9.2 Merge with Conflict Resolution
	fmt.Println("\n--- 9.2 Merge with Conflict Resolution ---")

	mergeWithConflict := func(m1, m2 map[string]int, resolve func(int, int) int) map[string]int {
		result := make(map[string]int)

		// کپی m1
		for k, v := range m1 {
			result[k] = v
		}

		// اضافه کردن m2 با حل اختلاف
		for k, v := range m2 {
			if existing, ok := result[k]; ok {
				result[k] = resolve(existing, v)
			} else {
				result[k] = v
			}
		}

		return result
	}

	mX := map[string]int{"a": 1, "b": 2}
	mY := map[string]int{"b": 10, "c": 3}

	// حل اختلاف: جمع مقادیر
	merged2 := mergeWithConflict(mX, mY, func(a, b int) int {
		return a + b
	})

	fmt.Printf("  Map X: %v\n", mX)
	fmt.Printf("  Map Y: %v\n", mY)
	fmt.Printf("  Merged (sum conflicts): %v\n", merged2)

	// 9.3 GroupBy - گروه‌بندی بر اساس تابع
	fmt.Println("\n--- 9.3 GroupBy (Custom) ---")

	groupBy := func(values []string, keyFunc func(string) string) map[string][]string {
		result := make(map[string][]string)
		for _, v := range values {
			key := keyFunc(v)
			result[key] = append(result[key], v)
		}
		return result
	}

	fruits3 := []string{"apple", "apricot", "banana", "blueberry", "cherry"}

	// گروه‌بندی بر اساس حرف اول
	byFirstLetter := groupBy(fruits3, func(s string) string {
		return string(s[0])
	})

	fmt.Printf("  Fruits: %v\n", fruits3)
	fmt.Printf("  Grouped by first letter: %v\n", byFirstLetter)

	// ============================================================================
	// بخش 10: جمع‌بندی و جدول مرجع سریع
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📚 SECTION 10: QUICK REFERENCE")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ FUNCTION              │ DESCRIPTION                              │")
	fmt.Println("├───────────────────────┼─────────────────────────────────────────┤")
	fmt.Println("│ maps.Clone(m)         │ Returns a shallow copy of map m         │")
	fmt.Println("│ maps.Copy(dst, src)   │ Copies all keys/values from src to dst  │")
	fmt.Println("│ maps.DeleteFunc(m, f) │ Deletes entries where f(key,val) is true│")
	fmt.Println("│ maps.Equal(m1, m2)    │ Checks if two maps contain same entries │")
	fmt.Println("│ maps.EqualFunc(m1,m2,f)│ Checks equality with custom value func │")
	fmt.Println("│ maps.Keys(m)          │ Returns iterator over map keys          │")
	fmt.Println("│ maps.Values(m)        │ Returns iterator over map values        │")
	fmt.Println("└───────────────────────┴─────────────────────────────────────────┘")

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ CUSTOM FUNCTION       │ DESCRIPTION                              │")
	fmt.Println("├───────────────────────┼─────────────────────────────────────────┤")
	fmt.Println("│ Filter(m, f)          │ Returns new map with filtered entries    │")
	fmt.Println("│ MapValues(m, f)       │ Returns new map with transformed values  │")
	fmt.Println("│ Merge(maps...)        │ Merges multiple maps                     │")
	fmt.Println("│ Invert(m)             │ Swaps keys and values                    │")
	fmt.Println("│ Intersection(m1,m2)   │ Returns common entries                   │")
	fmt.Println("│ Difference(m1,m2)     │ Returns entries only in one map          │")
	fmt.Println("└───────────────────────┴─────────────────────────────────────────┘")

	fmt.Println("\n💡 GOLDEN RULES:")
	fmt.Println("  1. Always use make() before writing to a map")
	fmt.Println("  2. Never assume iteration order - sort keys if needed")
	fmt.Println("  3. Use maps.Copy instead of manual loops for efficiency")
	fmt.Println("  4. Check 'ok' when reading from maps: val, ok := m[key]")
	fmt.Println("  5. Pre-allocate capacity with make(map[K]V, size) when size known")
	fmt.Println("  6. maps.Equal compares nil and empty maps as equal")
	fmt.Println("  7. maps.Clone creates shallow copy - pointer values are shared")
	fmt.Println("  8. Use maps.Keys and slices.Collect for sorted iteration")
	fmt.Println("  9. For concurrent access, use sync.Map or mutex")
	fmt.Println("  10. nil maps are safe to read but not to write - will panic")

	fmt.Println("\n🎯 QUICK COMPARISON:")
	fmt.Println("  • maps.Equal(m1, m2)    → Compares values using ==")
	fmt.Println("  • maps.EqualFunc(m1,m2,f) → Compares values using custom function")
	fmt.Println("  • reflect.DeepEqual(m1,m2) → Deep comparison (slower)")
	fmt.Println("  • manual loop           → Most flexible but verbose")

	fmt.Println("\n📦 INSTALLATION:")
	fmt.Println("  No installation needed - maps package is part of standard library since Go 1.21")
}

// تابع کمکی برای تکرار رشته
func stringsRepeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}
