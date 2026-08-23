// ============================================================================
// FILE: pointers_complete_guide.go
// TITLE: راهنمای کامل اشاره‌گرها (Pointers) در Go - از صفر تا صد
// HOW TO RUN: go run pointers_complete_guide.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - اشاره‌گر چیست و چرا در Go مهم است؟
// ============================================================================
//
// اشاره‌گر (Pointer): متغیری که آدرس حافظه متغیر دیگر را نگهداری می‌کند
//
// چرا اشاره‌گرها در Go مهم هستند؟
// 1. کارایی: به جای کپی کردن داده‌های حجیم، فقط آدرس آن را منتقل می‌کنیم
// 2. تغییر داده‌ها: توابع می‌توانند متغیرهای خارج از scope خود را تغییر دهند
// 3. ساختارهای داده: linked list, tree, graph بدون اشاره‌گر غیرممکن است
// 4. متدهای receiver: انتخاب بین value receiver و pointer receiver
// 5. nil: نشان‌دهنده "هیچ مقداری" برای انواع reference
//
// مفاهیم کلیدی:
// - & عملگر: آدرس یک متغیر را می‌گیرد (address-of)
// - * عملگر: مقدار در آدرس را می‌خواند (dereference)
// - نوع *T: اشاره‌گر به نوع T (مثلاً *int)
// - nil: مقدار صفر برای اشاره‌گرها (هیچ آدرسی)
//
// قانون طلایی:
// "از اشاره‌گرها برای تغییر داده‌ها و جلوگیری از کپی استفاده کن،
//  اما از آن‌ها زیاد استفاده نکن - سادگی Go را حفظ کن"
// ============================================================================

package main

import (
	"fmt"
	"unsafe"
)

// ============================================================================
// بخش 1: مبانی اشاره‌گرها - تعریف، عملگرها، مقداردهی
// ============================================================================

func demonstratePointerBasics() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🔍 POINTER BASICS - Definition, Operators, Initialization")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 1.1 تعریف اشاره‌گر و گرفتن آدرس
	// ============================================
	fmt.Println("\n--- 1.1 Declaring Pointers and Address-of Operator (&) ---")

	var x int = 42  // یک متغیر معمولی
	var p *int = &x // p اشاره‌گر به x (آدرس x را نگه می‌دارد)

	fmt.Printf("Value of x: %d\n", x)
	fmt.Printf("Address of x: %p\n", &x)
	fmt.Printf("Value of p (address): %p\n", p)
	fmt.Printf("Type of p: %T\n", p)

	// ============================================
	// 1.2 عملگر Dereference (*)
	// ============================================
	fmt.Println("\n--- 1.2 Dereference Operator (*) ---")

	fmt.Printf("Value at address p: %d\n", *p) // خواندن مقدار

	// تغییر مقدار از طریق اشاره‌گر
	*p = 100
	fmt.Printf("After *p = 100, x is now: %d\n", x)

	// تغییر مستقیم x
	x = 200
	fmt.Printf("After x = 200, *p is now: %d\n", *p)

	// ============================================
	// 1.3 مقدار صفر اشاره‌گر (nil)
	// ============================================
	fmt.Println("\n--- 1.3 Zero Value (nil) ---")

	var q *int // اعلان اشاره‌گر بدون مقداردهی
	fmt.Printf("Declared pointer q: %v (nil)\n", q)

	if q == nil {
		fmt.Println("  q is nil - points to nothing")
	}

	// ❌ این خط کد panic می‌کند (dereference nil pointer)
	//fmt.Println(*q)

	// مقداردهی به اشاره‌گر nil
	var r int = 99
	q = &r
	fmt.Printf("After q = &r, q points to: %d\n", *q)

	// ============================================
	// 1.4 اشاره‌گر به اشاره‌گر (Double Pointer)
	// ============================================
	fmt.Println("\n--- 1.4 Pointer to Pointer (**T) ---")

	var a int = 10
	var ptr *int = &a
	var ptrPtr **int = &ptr

	fmt.Printf("a = %d\n", a)
	fmt.Printf("ptr = &a = %p, *ptr = %d\n", ptr, *ptr)
	fmt.Printf("ptrPtr = &ptr = %p, *ptrPtr = %p, **ptrPtr = %d\n",
		ptrPtr, *ptrPtr, **ptrPtr)

	// تغییر a از طریق double pointer
	**ptrPtr = 999
	fmt.Printf("After **ptrPtr = 999, a = %d\n", a)
}

// ============================================================================
// بخش 2: اشاره‌گرها در توابع (پاس دادن با value vs reference)
// ============================================================================

// تابع با value receiver (کپی می‌کند)
func changeValueCopy(val int) {
	val = 100 // فقط کپی محلی را تغییر می‌دهد
}

// تابع با pointer receiver (اصل را تغییر می‌دهد)
func changeValuePointer(val *int) {
	*val = 100 // مقدار اصلی را تغییر می‌دهد
}

// تابع با value receiver برای struct
type Person struct {
	Name string
	Age  int
}

func updatePersonCopy(p Person) {
	p.Name = "Changed"
	p.Age = 99
}

func updatePersonPointer(p *Person) {
	p.Name = "Changed" // Go اجازه می‌دهد بدون (*p).Name بنویسیم
	p.Age = 99
}

func demonstratePointersInFunctions() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📤 POINTERS IN FUNCTIONS - Pass by Value vs Reference")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 2.1 پاس دادن value vs pointer برای انواع پایه
	// ============================================
	fmt.Println("\n--- 2.1 Basic Types ---")

	num := 10
	fmt.Printf("Before changeValueCopy: %d\n", num)
	changeValueCopy(num)
	fmt.Printf("After changeValueCopy: %d (unchanged)\n", num)

	fmt.Printf("\nBefore changeValuePointer: %d\n", num)
	changeValuePointer(&num)
	fmt.Printf("After changeValuePointer: %d (changed!)\n", num)

	// ============================================
	// 2.2 پاس دادن struct با value vs pointer
	// ============================================
	fmt.Println("\n--- 2.2 Struct Types ---")

	p1 := Person{Name: "Ali", Age: 30}
	fmt.Printf("Before updatePersonCopy: %+v\n", p1)
	updatePersonCopy(p1)
	fmt.Printf("After updatePersonCopy: %+v (unchanged)\n", p1)

	p2 := &Person{Name: "Reza", Age: 25}
	fmt.Printf("\nBefore updatePersonPointer: %+v\n", p2)
	updatePersonPointer(p2)
	fmt.Printf("After updatePersonPointer: %+v (changed!)\n", p2)

	// ============================================
	// 2.3 برگرداندن اشاره‌گر از تابع
	// ============================================
	fmt.Println("\n--- 2.3 Returning Pointers from Functions ---")

	newPerson := createPerson("Sara", 28)
	fmt.Printf("Created person: %+v\n", newPerson)
}

// تابعی که اشاره‌گر برمی‌گرداند
func createPerson(name string, age int) *Person {
	// Go در heap allocation تصمیم می‌گیرد (escape analysis)
	return &Person{Name: name, Age: age}
}

// ============================================================================
// بخش 3: اشاره‌گرها در ساختارهای داده (struct)
// ============================================================================

type Node struct {
	Value int
	Next  *Node // اشاره‌گر به Node بعدی (برای linked list)
}

type LinkedList struct {
	Head *Node
}

func (l *LinkedList) Insert(value int) {
	newNode := &Node{Value: value}

	if l.Head == nil {
		l.Head = newNode
		return
	}

	current := l.Head
	for current.Next != nil {
		current = current.Next
	}
	current.Next = newNode
}

func (l *LinkedList) Display() {
	current := l.Head
	for current != nil {
		fmt.Printf("%d -> ", current.Value)
		current = current.Next
	}
	fmt.Println("nil")
}

func (l *LinkedList) Search(value int) bool {
	current := l.Head
	for current != nil {
		if current.Value == value {
			return true
		}
		current = current.Next
	}
	return false
}

func demonstratePointersInDataStructures() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🔗 POINTERS IN DATA STRUCTURES - Linked List Example")
	fmt.Println(stringsRepeat("=", 80))

	list := &LinkedList{}

	fmt.Println("Inserting values: 10, 20, 30, 40, 50")
	list.Insert(10)
	list.Insert(20)
	list.Insert(30)
	list.Insert(40)
	list.Insert(50)

	fmt.Print("Linked list: ")
	list.Display()

	fmt.Printf("Search 30: %v\n", list.Search(30))
	fmt.Printf("Search 99: %v\n", list.Search(99))
}

// ============================================================================
// بخش 4: اشاره‌گرها در آرایه‌ها و اسلایس‌ها
// ============================================================================

func demonstratePointersWithArraysSlices() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📊 POINTERS WITH ARRAYS AND SLICES")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 4.1 اشاره‌گر به آرایه
	// ============================================
	fmt.Println("\n--- 4.1 Pointer to Array ---")

	arr := [5]int{1, 2, 3, 4, 5}
	var ptrToArr *[5]int = &arr

	fmt.Printf("Original array: %v\n", arr)

	// تغییر از طریق اشاره‌گر
	ptrToArr[0] = 100
	(*ptrToArr)[1] = 200
	fmt.Printf("After pointer modification: %v\n", arr)

	// ============================================
	// 4.2 اشاره‌گر به اسلایس (کمتر رایج)
	// ============================================
	fmt.Println("\n--- 4.2 Pointer to Slice ---")

	slice := []int{1, 2, 3}
	var ptrToSlice *[]int = &slice

	fmt.Printf("Original slice: %v\n", slice)

	// تغییر از طریق اشاره‌گر
	*ptrToSlice = append(*ptrToSlice, 4, 5)
	fmt.Printf("After pointer modification: %v\n", slice)

	// ============================================
	// 4.3 تفاوت آرایه و اسلایس در اشاره‌گرها
	// ============================================
	fmt.Println("\n--- 4.3 Array vs Slice with Pointers ---")

	modifyArray := func(arr *[3]int) {
		(*arr)[0] = 999
	}

	modifySlice := func(slice *[]int) {
		(*slice)[0] = 999
	}

	arr2 := [3]int{1, 2, 3}
	slice2 := []int{1, 2, 3}

	fmt.Printf("Before modifyArray: %v\n", arr2)
	modifyArray(&arr2)
	fmt.Printf("After modifyArray: %v\n", arr2)

	fmt.Printf("\nBefore modifySlice: %v\n", slice2)
	modifySlice(&slice2)
	fmt.Printf("After modifySlice: %v\n", slice2)
}

// ============================================================================
// بخش 5: اشاره‌گرها در متدها (Value Receiver vs Pointer Receiver)
// ============================================================================

type Counter struct {
	value int
}

// Value receiver - کپی می‌کند
func (c Counter) IncrementCopy() {
	c.value++
	fmt.Printf("  Inside IncrementCopy: %d (local copy)\n", c.value)
}

// Pointer receiver - اصل را تغییر می‌دهد
func (c *Counter) IncrementPointer() {
	c.value++
	fmt.Printf("  Inside IncrementPointer: %d\n", c.value)
}

// Value receiver برای خواندن (بدون تغییر)
func (c Counter) GetValue() int {
	return c.value
}

type LargeStruct struct {
	Data [1024]byte // 1KB داده
	ID   int
}

// Value receiver - کل struct کپی می‌شود (expensive)
func (l LargeStruct) ProcessValue() int {
	return l.ID
}

// Pointer receiver - فقط اشاره‌گر کپی می‌شود (cheap)
func (l *LargeStruct) ProcessPointer() int {
	return l.ID
}

func demonstrateValueVsPointerReceivers() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🎯 VALUE RECEIVER vs POINTER RECEIVER")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 5.1 تفاوت در تغییر مقدار
	// ============================================
	fmt.Println("\n--- 5.1 Modifying State ---")

	c := Counter{value: 10}

	fmt.Printf("Initial value: %d\n", c.GetValue())

	c.IncrementCopy()
	fmt.Printf("After IncrementCopy: %d (unchanged!)\n", c.GetValue())

	c.IncrementPointer()
	fmt.Printf("After IncrementPointer: %d (changed!)\n", c.GetValue())

	// ============================================
	// 5.2 عملکرد (performance)
	// ============================================
	fmt.Println("\n--- 5.2 Performance Considerations ---")

	large := LargeStruct{ID: 42}

	// Value receiver: کل 1KB کپی می‌شود
	_ = large.ProcessValue()

	// Pointer receiver: فقط 8 بایت (آدرس) کپی می‌شود
	_ = large.ProcessPointer()

	fmt.Println("  LargeStruct size: 1KB+")
	fmt.Println("  Value receiver: copies entire struct")
	fmt.Println("  Pointer receiver: copies only pointer (8 bytes)")

	// ============================================
	// 5.3 چه زمانی از کدام استفاده کنیم؟
	// ============================================
	fmt.Println("\n--- 5.3 When to Use Which ---")

	fmt.Println("\n✅ Use Pointer Receiver when:")
	fmt.Println("  • You need to modify the receiver")
	fmt.Println("  • The struct is large (avoid copying)")
	fmt.Println("  • The struct contains mutex (mutex can't be copied)")
	fmt.Println("  • You want consistency (all methods use pointer)")

	fmt.Println("\n✅ Use Value Receiver when:")
	fmt.Println("  • The method doesn't modify the receiver")
	fmt.Println("  • The struct is small (a few fields)")
	fmt.Println("  • The type is a basic type (int, string, etc.)")
	fmt.Println("  • You want immutability")
}

// ============================================================================
// بخش 6: اشاره‌گرها و nil (مدیریت حافظه)
// ============================================================================

type SafeResource struct {
	data *int
}

func (r *SafeResource) Initialize(value int) {
	r.data = new(int) // تخصیص حافظه جدید
	*r.data = value
}

func (r *SafeResource) GetValue() (int, error) {
	if r.data == nil {
		return 0, fmt.Errorf("resource not initialized")
	}
	return *r.data, nil
}

func (r *SafeResource) Cleanup() {
	if r.data != nil {
		r.data = nil // آزاد کردن اشاره‌گر (GC بعداً جمع می‌کند)
	}
}

func demonstrateNilPointers() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("⚠️ NIL POINTERS - Safety and Best Practices")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 6.1 بررسی nil قبل از استفاده
	// ============================================
	fmt.Println("\n--- 6.1 Checking for nil ---")

	var p *int

	if p == nil {
		fmt.Println("  p is nil - safe to check")
	}

	// همیشه قبل از dereference بررسی کن
	if p != nil {
		fmt.Printf("  Value: %d\n", *p)
	} else {
		fmt.Println("  Cannot dereference nil pointer")
	}

	// ============================================
	// 6.2 الگوی Safe Resource
	// ============================================
	fmt.Println("\n--- 6.2 Safe Resource Pattern ---")

	var res SafeResource

	// تلاش برای استفاده قبل از initialization
	if val, err := res.GetValue(); err != nil {
		fmt.Printf("  Error: %v\n", err)
	} else {
		fmt.Printf("  Value: %d\n", val)
	}

	// Initialize
	res.Initialize(42)
	if val, err := res.GetValue(); err != nil {
		fmt.Printf("  Error: %v\n", err)
	} else {
		fmt.Printf("  Value after init: %d\n", val)
	}

	// Cleanup
	res.Cleanup()
	if _, err := res.GetValue(); err != nil {
		fmt.Printf("  After cleanup: %v\n", err)
	}

	// ============================================
	// 6.3 new() vs &T{}
	// ============================================
	fmt.Println("\n--- 6.3 new() vs &T{} ---")

	// new(T): تخصیص حافظه و برگرداندن اشاره‌گر به مقدار صفر
	p1 := new(int)
	fmt.Printf("new(int): p1=%p, *p1=%d\n", p1, *p1)

	// هر دو معادل هستند
	var p3 *int = new(int)
	*p3 = 100
	fmt.Printf("p3=%d, p4=%d (equivalent)\n", *p3)
}

// ============================================================================
// بخش 7: اشاره‌گرها و escape analysis (تخصیص heap vs stack)
// ============================================================================

// این تابع ممکن است باعث escape به heap شود
func createPointer() *int {
	x := 42
	return &x // x به heap می‌رود (escape)
}

// این تابع در stack می‌ماند
func createValue() int {
	x := 42
	return x
}

func demonstrateEscapeAnalysis() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🏃 ESCAPE ANALYSIS - Stack vs Heap Allocation")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println("\n--- Escape Analysis Example ---")

	// کامپایلر تصمیم می‌گیرد کجا تخصیص دهد
	val := createValue()   // احتمالاً در stack
	ptr := createPointer() // حتماً در heap (escape)

	fmt.Printf("Value: %d (stack)\n", val)
	fmt.Printf("Pointer: %d (heap - escaped)\n", *ptr)

	fmt.Println("\n💡 Check escape analysis with:")
	fmt.Println("   go build -gcflags=\"-m\" pointers_complete_guide.go")
	fmt.Println("\n📌 Rules of thumb:")
	fmt.Println("  • Returning pointer to local variable → escapes to heap")
	fmt.Println("  • Storing pointer in struct that escapes → escapes")
	fmt.Println("  • Interface values → often escape")
	fmt.Println("  • Large objects (>64KB) → always heap")
}

// ============================================================================
// بخش 8: اشاره‌گرهای unsafe (برای موارد خاص)
// ============================================================================

func demonstrateUnsafePointers() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("⚠️ UNSAFE POINTERS - Advanced (Use with caution)")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 8.1 تبدیل بین انواع با unsafe.Pointer
	// ============================================
	fmt.Println("\n--- 8.1 Type Punning with unsafe.Pointer ---")

	var f float64 = 3.14159

	// تبدیل float64 به uint64 (برای دیدن بیت‌ها)
	bits := *(*uint64)(unsafe.Pointer(&f))
	fmt.Printf("float64: %f\n", f)
	fmt.Printf("bits: 0x%X\n", bits)

	// تبدیل back
	f2 := *(*float64)(unsafe.Pointer(&bits))
	fmt.Printf("Back to float64: %f\n", f2)

	// ============================================
	// 8.2 محاسبه offset فیلدهای struct
	// ============================================
	fmt.Println("\n--- 8.2 Field Offsets ---")

	type MyStruct struct {
		A int32
		B int64
		C byte
	}

	var s MyStruct
	offsetA := unsafe.Offsetof(s.A)
	offsetB := unsafe.Offsetof(s.B)
	offsetC := unsafe.Offsetof(s.C)

	fmt.Printf("Offset of A: %d bytes\n", offsetA)
	fmt.Printf("Offset of B: %d bytes\n", offsetB)
	fmt.Printf("Offset of C: %d bytes\n", offsetC)
	fmt.Printf("Total size: %d bytes\n", unsafe.Sizeof(s))

	fmt.Println("\n⚠️  WARNING: unsafe package is called 'unsafe' for a reason!")
	fmt.Println("  • Not guaranteed to work across Go versions")
	fmt.Println("  • Can crash your program")
	fmt.Println("  • Only use when absolutely necessary")
}

// ============================================================================
// بخش 9: الگوهای رایج و بهترین تمرین‌ها
// ============================================================================

// الگوی 1: Optional parameters (پارامترهای اختیاری)
type Config struct {
	Host string
	Port int
}

type ConfigOption func(*Config)

func WithHost(host string) ConfigOption {
	return func(c *Config) {
		c.Host = host
	}
}

func WithPort(port int) ConfigOption {
	return func(c *Config) {
		c.Port = port
	}
}

func NewConfig(opts ...ConfigOption) *Config {
	cfg := &Config{
		Host: "localhost",
		Port: 8080,
	}
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

// الگوی 2: Singleflight با اشاره‌گر (قبلاً دیدیم)
type LazyInit struct {
	value *int
	init  func() int
}

func (l *LazyInit) Get() int {
	if l.value == nil {
		l.value = new(int)
		*l.value = l.init()
	}
	return *l.value
}

// الگوی 3: Object Pool با اشاره‌گر
type Pool struct {
	objects chan *LargeStruct
}

func NewPool(size int) *Pool {
	p := &Pool{
		objects: make(chan *LargeStruct, size),
	}
	for i := 0; i < size; i++ {
		p.objects <- &LargeStruct{ID: i}
	}
	return p
}

func (p *Pool) Get() *LargeStruct {
	select {
	case obj := <-p.objects:
		return obj
	default:
		return &LargeStruct{} // جدید
	}
}

func (p *Pool) Put(obj *LargeStruct) {
	select {
	case p.objects <- obj:
	default:
		// drop if full
	}
}

func demonstratePointerPatterns() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🎨 COMMON POINTER PATTERNS")
	fmt.Println(stringsRepeat("=", 80))

	// الگوی 1: Options Pattern
	fmt.Println("\n--- Pattern 1: Options Pattern ---")
	cfg1 := NewConfig()
	cfg2 := NewConfig(WithHost("0.0.0.0"), WithPort(9000))
	fmt.Printf("  Default config: %+v\n", cfg1)
	fmt.Printf("  Custom config: %+v\n", cfg2)

	// الگوی 2: Lazy Initialization
	fmt.Println("\n--- Pattern 2: Lazy Initialization ---")
	lazy := &LazyInit{
		init: func() int {
			fmt.Println("    Initializing expensive resource...")
			return 42
		},
	}
	fmt.Printf("  First get: %d\n", lazy.Get())
	fmt.Printf("  Second get: %d (cached)\n", lazy.Get())

	// الگوی 3: Object Pool
	fmt.Println("\n--- Pattern 3: Object Pool ---")
	pool := NewPool(3)
	obj1 := pool.Get()
	obj2 := pool.Get()
	fmt.Printf("  Got objects: ID=%d, ID=%d\n", obj1.ID, obj2.ID)
	pool.Put(obj1)
	obj3 := pool.Get()
	fmt.Printf("  Reused object ID: %d\n", obj3.ID)
}

// ============================================================================
// بخش 10: اشتباهات رایج و نکات مهم
// ============================================================================

func demonstrateCommonMistakes() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("❌ COMMON POINTER MISTAKES AND HOW TO AVOID THEM")
	fmt.Println(stringsRepeat("=", 80))

	// اشتباه 1: Dereference nil pointer
	fmt.Println("\n❌ Mistake 1: Dereferencing nil pointer")
	fmt.Println("   var p *int")
	fmt.Println("   *p = 42  // PANIC!")
	fmt.Println("   ✅ Always check: if p != nil { *p = 42 }")

	// اشتباه 2: برگرداندن اشاره‌گر به متغیر محلی (گاهی مشکل‌ساز)
	fmt.Println("\n❌ Mistake 2: Returning pointer to local variable")
	fmt.Println("   func bad() *int { x := 42; return &x } // escapes to heap")
	fmt.Println("   ✅ It's actually safe (escape analysis), but be aware")

	// اشتباه 3: استفاده از اشاره‌گر به map
	fmt.Println("\n❌ Mistake 3: Pointer to map (unnecessary)")
	fmt.Println("   var m *map[string]int  // rarely needed")
	fmt.Println("   ✅ Maps are already references: var m map[string]int")

	// اشتباه 4: استفاده از اشاره‌گر به channel
	fmt.Println("\n❌ Mistake 4: Pointer to channel (unnecessary)")
	fmt.Println("   var ch *chan int  // rarely needed")
	fmt.Println("   ✅ Channels are already references: var ch chan int")

	// اشتباه 5: کپی کردن struct با mutex
	fmt.Println("\n❌ Mistake 5: Copying struct with mutex")
	fmt.Println("   type Counter struct { mu sync.Mutex; count int }")
	fmt.Println("   c2 := c1  // ❌ copies mutex!")
	fmt.Println("   ✅ Use pointer: c2 := &c1")

	// اشتباه 6: فراموش کردن & در فراخوانی متد pointer receiver
	fmt.Println("\n❌ Mistake 6: Forgetting & for pointer receiver")
	fmt.Println("   type Counter struct { count int }")
	fmt.Println("   func (c *Counter) Inc() { c.count++ }")
	fmt.Println("   var c Counter")
	fmt.Println("   c.Inc()  // ✅ Go automatically takes address!")
	fmt.Println("   ✅ Go handles this automatically for methods")
}

// ============================================================================
// بخش 11: جمع‌بندی و جدول مرجع سریع
// ============================================================================

func main() {
	fmt.Println(stringsRepeat("=", 80))
	fmt.Println("🎯 COMPLETE POINTERS GUIDE IN GO")
	fmt.Println("From Basics to Advanced Patterns")
	fmt.Println(stringsRepeat("=", 80))

	// بخش 1: مبانی
	demonstratePointerBasics()

	// بخش 2: اشاره‌گرها در توابع
	demonstratePointersInFunctions()

	// بخش 3: اشاره‌گرها در ساختارهای داده
	demonstratePointersInDataStructures()

	// بخش 4: اشاره‌گرها با آرایه و اسلایس
	demonstratePointersWithArraysSlices()

	// بخش 5: Value vs Pointer Receiver
	demonstrateValueVsPointerReceivers()

	// بخش 6: Nil Pointers
	demonstrateNilPointers()

	// بخش 7: Escape Analysis
	demonstrateEscapeAnalysis()

	// بخش 8: Unsafe Pointers
	demonstrateUnsafePointers()

	// بخش 9: الگوهای رایج
	demonstratePointerPatterns()

	// بخش 10: اشتباهات رایج
	demonstrateCommonMistakes()

	// بخش 11: جمع‌بندی
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📚 POINTERS QUICK REFERENCE")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ OPERATOR  │ MEANING                          │ EXAMPLE          │")
	fmt.Println("├───────────┼──────────────────────────────────┼──────────────────┤")
	fmt.Println("│ &x        │ Address of x                     │ p := &x          │")
	fmt.Println("│ *p        │ Value at address p (dereference) │ y := *p          │")
	fmt.Println("│ *T        │ Pointer to type T                │ var p *int       │")
	fmt.Println("│ new(T)    │ Allocate zero T, return pointer  │ p := new(int)    │")
	fmt.Println("│ &T{...}   │ Allocate and initialize T        │ p := &Point{1,2} │")
	fmt.Println("└───────────┴──────────────────────────────────┴──────────────────┘")

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ USE CASE                     │ RECOMMENDATION                  │")
	fmt.Println("├──────────────────────────────┼──────────────────────────────────┤")
	fmt.Println("│ Modify function parameter    │ Use pointer                     │")
	fmt.Println("│ Large struct (>64 bytes)     │ Use pointer (avoid copy)        │")
	fmt.Println("│ Method modifies receiver     │ Pointer receiver                │")
	fmt.Println("│ Struct with mutex            │ Always pointer                  │")
	fmt.Println("│ Small immutable data         │ Value is fine                   │")
	fmt.Println("│ Map, slice, channel          │ Already references              │")
	fmt.Println("│ JSON/API response            │ Pointer for optional fields     │")
	fmt.Println("└──────────────────────────────┴──────────────────────────────────┘")

	fmt.Println("\n💡 GOLDEN RULES:")
	fmt.Println("  1. Use pointers to modify data across function boundaries")
	fmt.Println("  2. Always check for nil before dereferencing")
	fmt.Println("  3. Prefer value receivers for small, immutable types")
	fmt.Println("  4. Use pointer receivers for methods that modify state")
	fmt.Println("  5. Maps, slices, channels are already references")
	fmt.Println("  6. Don't copy structs that contain mutexes")
	fmt.Println("  7. Use new(T) or &T{} for allocating pointers")
	fmt.Println("  8. Escape analysis decides stack vs heap")
	fmt.Println("  9. unsafe package is truly unsafe - avoid if possible")
	fmt.Println("  10. When in doubt, use value (simpler and safer)")

	fmt.Println("\n🎯 WHEN TO USE POINTERS:")
	fmt.Println("  • You need to modify the original value")
	fmt.Println("  • The struct is large (avoid expensive copies)")
	fmt.Println("  • You need nil to represent 'no value'")
	fmt.Println("  • Working with sync.Mutex or other non-copyable types")
	fmt.Println("  • Implementing linked structures (trees, lists)")
	fmt.Println("  • Performance optimization (after profiling)")

	fmt.Println("\n🎯 WHEN NOT TO USE POINTERS:")
	fmt.Println("  • Small primitive types (int, bool, float64)")
	fmt.Println("  • Strings (already small and immutable)")
	fmt.Println("  • Maps, slices, channels (already references)")
	fmt.Println("  • When you don't need to modify the value")
	fmt.Println("  • For simple data passing (keep it simple)")
}

// ============================================================================
// بخش 12: توابع کمکی
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
go run pointers_complete_guide.go

# بررسی escape analysis
go build -gcflags="-m" pointers_complete_guide.go

# اجرا با race detector
go run -race pointers_complete_guide.go
*/
