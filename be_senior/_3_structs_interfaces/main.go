// ============================================================================
// FILE: structs_interfaces_guide.go
// TITLE: راهنمای کامل ساختارها (structs) و اینترفیس‌ها (interfaces) در Go
// HOW TO RUN: go run structs_interfaces_guide.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - struct و interface چه هستند و چرا مهمند؟
// ============================================================================
//
// Struct (ساختار):
// - یک نوع داده ترکیبی که گروهی از فیلدها (با انواع مختلف) را کنار هم جمع می‌کند
// - معادل کلاس در زبان‌های دیگر، اما بدون ارث‌بری
// - فقط داده را نگه می‌دارد (بدون متد - متدها جدا تعریف می‌شوند)
//
// Interface (اینترفیس):
// - مجموعه‌ای از متدها که یک نوع باید پیاده‌سازی کند
// - قلب برنامه‌نویسی شیء‌گرا در Go (اما بدون ارث‌بری)
// - Go از استنتاج نوع (structural typing) استفاده می‌کند (نه nominal)
//
// تفاوت مهم Go با زبان‌های دیگر:
// 1. ارث‌بری وجود ندارد (به جای آن از composition استفاده می‌کنیم)
// 2. اینترفیس‌ها به صورت implicit پیاده‌سازی می‌شوند (بدون کلمه implements)
// 3. هر نوعی که متدهای یک اینترفیس را داشته باشد، آن اینترفیس را پیاده‌سازی می‌کند
//
// قانون طلایی:
// "با structها داده را سازماندهی کن، با interfaceها رفتار را تعریف کن"
// ============================================================================

package __structs_interfaces

import (
	"fmt"
	"math"
)

// ============================================================================
// بخش 1: ساختارها (Structs) - تعریف، فیلدها، مقداردهی
// ============================================================================

// تعریف یک struct ساده
type Person struct {
	Name    string
	Age     int
	Email   string
	Address string
}

// struct با تگ (tags) - برای JSON، XML، Validation و غیره
type User struct {
	ID       int    `json:"id" xml:"id"`
	Username string `json:"username" validate:"required,min=3"`
	Email    string `json:"email" validate:"email"`
	Password string `json:"-"` // این فیلد در JSON نادیده گرفته می‌شود
}

// struct تو در تو (embedded struct)
type Employee struct {
	Person     // embedding (ارث‌بری نیست، composition است)
	EmployeeID string
	Department string
	Salary     float64
	Manager    *Employee // اشاره‌گر به خودش (برای ساختارهای سلسله‌مراتبی)
}

// struct با فیلدهای ناشناس (anonymous fields)
type AnonymousStruct struct {
	string  // فیلد بدون نام - نوع به عنوان نام استفاده می‌شود
	int
	bool
}

func demonstrateStructBasics() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📦 STRUCT BASICS - Definition, Fields, Initialization")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 1.1 روش‌های مختلف مقداردهی struct
	// ============================================
	fmt.Println("\n--- 1.1 Initialization Methods ---")

	// روش 1: با نام فیلدها (recommended)
	p1 := Person{
		Name:    "Ali Rezaei",
		Age:     30,
		Email:   "ali@example.com",
		Address: "Tehran, Iran",
	}
	fmt.Printf("Method 1 (named fields): %+v\n", p1)

	// روش 2: به ترتیب فیلدها (not recommended - fragile)
	p2 := Person{"Sara Mohammadi", 28, "sara@example.com", "Isfahan, Iran"}
	fmt.Printf("Method 2 (order): %+v\n", p2)

	// روش 3: مقداردهی جداگانه
	var p3 Person
	p3.Name = "Reza Karimi"
	p3.Age = 35
	p3.Email = "reza@example.com"
	fmt.Printf("Method 3 (separate): %+v\n", p3)

	// روش 4: با new (برمی‌گرداند *Person)
	p4 := new(Person)
	p4.Name = "Mina Hosseini"
	p4.Age = 32
	fmt.Printf("Method 4 (new): %+v\n", p4)

	// روش 5: اشاره‌گر با &
	p5 := &Person{Name: "Hassan Nasiri", Age: 40}
	fmt.Printf("Method 5 (&struct): %+v\n", p5)

	// ============================================
	// 1.2 دسترسی به فیلدها (با .)
	// ============================================
	fmt.Println("\n--- 1.2 Accessing Fields ---")

	p := Person{Name: "Test User", Age: 25}
	fmt.Printf("Name: %s, Age: %d\n", p.Name, p.Age)

	// تغییر مقدار فیلد
	p.Age = 26
	fmt.Printf("After change: Age = %d\n", p.Age)

	// ============================================
	// 1.3 Struct با فیلدهای ناشناس
	// ============================================
	fmt.Println("\n--- 1.3 Anonymous Fields ---")

	anon := AnonymousStruct{
		string: "Hello",
		int:    42,
		bool:   true,
	}
	fmt.Printf("Anonymous struct: %+v\n", anon)
	fmt.Printf("  string field: %s\n", anon.string)
	fmt.Printf("  int field: %d\n", anon.int)
	fmt.Printf("  bool field: %t\n", anon.bool)

	// ============================================
	// 1.4 Struct embedding (تو در تو)
	// ============================================
	fmt.Println("\n--- 1.4 Struct Embedding ---")

	emp := Employee{
		Person: Person{
			Name:  "Ali Rezaei",
			Age:   30,
			Email: "ali@company.com",
		},
		EmployeeID: "EMP001",
		Department: "Engineering",
		Salary:     7500000,
	}

	// دسترسی مستقیم به فیلدهای Person (promoted fields)
	fmt.Printf("Employee: %s (ID: %s), Age: %d\n",
		emp.Name, emp.EmployeeID, emp.Age)
	fmt.Printf("Full struct: %+v\n", emp)
}

// ============================================================================
// بخش 2: متدها روی Structها (Value vs Pointer Receivers)
// ============================================================================

// متد با value receiver (کپی می‌کند)
func (p Person) Greet() string {
	return fmt.Sprintf("Hello, I'm %s, %d years old", p.Name, p.Age)
}

// متد با pointer receiver (اصل را تغییر می‌دهد)
func (p *Person) HaveBirthday() {
	p.Age++
	fmt.Printf("  Happy birthday %s! Now %d years old\n", p.Name, p.Age)
}

// متد برای تغییر ایمیل
func (p *Person) UpdateEmail(newEmail string) {
	p.Email = newEmail
}

// متد برای Employee
func (e Employee) Work() string {
	return fmt.Sprintf("%s is working in %s department", e.Name, e.Department)
}

func (e *Employee) GiveRaise(percent float64) {
	e.Salary *= (1 + percent/100)
}

type Rectangle struct {
	Width, Height float64
}

// متد با value receiver (خواندن)
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

// متد با value receiver (خواندن)
func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

// متد با pointer receiver (تغییر)
func (r *Rectangle) Scale(factor float64) {
	r.Width *= factor
	r.Height *= factor
}

func demonstrateStructMethods() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🔧 STRUCT METHODS - Value vs Pointer Receivers")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 2.1 متدها روی Person
	// ============================================
	fmt.Println("\n--- 2.1 Person Methods ---")

	p := Person{Name: "Sara", Age: 28, Email: "sara@example.com"}

	// Value receiver (کپی)
	fmt.Println(p.Greet())

	// Pointer receiver (تغییر اصل)
	p.HaveBirthday()
	p.UpdateEmail("sara.new@example.com")
	fmt.Printf("Updated email: %s\n", p.Email)

	// ============================================
	// 2.2 متدها روی Rectangle
	// ============================================
	fmt.Println("\n--- 2.2 Rectangle Methods ---")

	r := Rectangle{Width: 10, Height: 5}
	fmt.Printf("Rectangle: width=%.1f, height=%.1f\n", r.Width, r.Height)
	fmt.Printf("Area: %.1f\n", r.Area())
	fmt.Printf("Perimeter: %.1f\n", r.Perimeter())

	r.Scale(2)
	fmt.Printf("\nAfter scaling by 2: width=%.1f, height=%.1f\n", r.Width, r.Height)
	fmt.Printf("New area: %.1f\n", r.Area())

	// ============================================
	// 2.3 متدها روی Employee
	// ============================================
	fmt.Println("\n--- 2.3 Employee Methods ---")

	emp := Employee{
		Person:     Person{Name: "Reza", Age: 35},
		EmployeeID: "EMP002",
		Department: "Sales",
		Salary:     5000000,
	}

	fmt.Println(emp.Work())
	fmt.Printf("Salary before raise: %.0f\n", emp.Salary)
	emp.GiveRaise(10)
	fmt.Printf("Salary after 10%% raise: %.0f\n", emp.Salary)
}

// ============================================================================
// بخش 3: اینترفیس‌ها (Interfaces) - قلب برنامه‌نویسی Go
// ============================================================================

// تعریف یک اینترفیس ساده
type Greeter interface {
	Greet() string
}

// اینترفیس برای شکل‌های هندسی
type Shape interface {
	Area() float64
	Perimeter() float64
}

// اینترفیس برای موجوداتی که صدا تولید می‌کنند
type Animal interface {
	Speak() string
	Move() string
}

// نوع Dog اینترفیس Animal را پیاده‌سازی می‌کند
type Dog struct {
	Name string
	Breed string
}

func (d Dog) Speak() string {
	return "Woof!"
}

func (d Dog) Move() string {
	return "Running on 4 legs"
}

// نوع Cat اینترفیس Animal را پیاده‌سازی می‌کند
type Cat struct {
	Name string
	Color string
}

func (c Cat) Speak() string {
	return "Meow!"
}

func (c Cat) Move() string {
	return "Walking gracefully"
}

// نوع Bird اینترفیس Animal را پیاده‌سازی می‌کند
type Bird struct {
	Name string
	CanFly bool
}

func (b Bird) Speak() string {
	return "Chirp!"
}

func (b Bird) Move() string {
	if b.CanFly {
		return "Flying in the sky"
	}
	return "Walking on ground"
}

// Circle پیاده‌سازی Shape
type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

// Square پیاده‌سازی Shape
type Square struct {
	Side float64
}

func (s Square) Area() float64 {
	return s.Side * s.Side
}

func (s Square) Perimeter() float64 {
	return 4 * s.Side
}

func demonstrateInterfaces() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🎯 INTERFACES - Defining Behavior")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 3.1 استفاده از اینترفیس به عنوان نوع
	// ============================================
	fmt.Println("\n--- 3.1 Using Interface as Type ---")

	var animal Animal

	animal = Dog{Name: "Max", Breed: "German Shepherd"}
	fmt.Printf("Dog: %s - %s, %s\n", animal.(Dog).Name, animal.Speak(), animal.Move())

	animal = Cat{Name: "Luna", Color: "Black"}
	fmt.Printf("Cat: %s - %s, %s\n", animal.(Cat).Name, animal.Speak(), animal.Move())

	animal = Bird{Name: "Tweety", CanFly: true}
	fmt.Printf("Bird: %s - %s, %s\n", animal.(Bird).Name, animal.Speak(), animal.Move())

	// ============================================
	// 3.2 اینترفیس به عنوان پارامتر تابع
	// ============================================
	fmt.Println("\n--- 3.2 Interface as Function Parameter ---")

	printAnimalInfo := func(a Animal) {
		fmt.Printf("  Animal says: %s, moves: %s\n", a.Speak(), a.Move())
	}

	printAnimalInfo(Dog{Name: "Buddy", Breed: "Labrador"})
	printAnimalInfo(Cat{Name: "Whiskers", Color: "White"})

	// ============================================
	// 3.3 اینترفیس به عنوان خروجی تابع
	// ============================================
	fmt.Println("\n--- 3.3 Interface as Return Type ---")

	getAnimal := func(animalType string) Animal {
		switch animalType {
		case "dog":
			return Dog{Name: "Generic Dog", Breed: "Mixed"}
		case "cat":
			return Cat{Name: "Generic Cat", Color: "Orange"}
		default:
			return Bird{Name: "Generic Bird", CanFly: true}
		}
	}

	animal1 := getAnimal("dog")
	animal2 := getAnimal("cat")
	fmt.Printf("Got: %T and %T\n", animal1, animal2)

	// ============================================
	// 3.4 اینترفیس Shape با انواع مختلف
	// ============================================
	fmt.Println("\n--- 3.4 Shape Interface ---")

	shapes := []Shape{
		Circle{Radius: 5},
		Square{Side: 4},
		Circle{Radius: 3},
		Square{Side: 6},
	}

	for i, shape := range shapes {
		fmt.Printf("Shape %d: Area=%.2f, Perimeter=%.2f\n",
			i+1, shape.Area(), shape.Perimeter())
	}
}

// ============================================================================
// بخش 4: اینترفیس خالی (empty interface) - interface{}
// ============================================================================

// interface{} می‌تواند هر نوعی را نگه دارد (مثل any در Go 1.18+)
func printAnything(v interface{}) {
	fmt.Printf("  Type: %T, Value: %v\n", v, v)
}

// تابع با پارامترهای variadic از نوع interface{}
func printAll(items ...interface{}) {
	for i, item := range items {
		fmt.Printf("  [%d] %T: %v\n", i, item, item)
	}
}

// استفاده از interface{} برای ساخت container عمومی
type Container struct {
	items []interface{}
}

func (c *Container) Add(item interface{}) {
	c.items = append(c.items, item)
}

func (c *Container) Get(index int) interface{} {
	if index < 0 || index >= len(c.items) {
		return nil
	}
	return c.items[index]
}

func (c *Container) Length() int {
	return len(c.items)
}

func demonstrateEmptyInterface() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📦 EMPTY INTERFACE (interface{}) - Any Type")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 4.1 interface{} می‌تواند هر چیزی باشد
	// ============================================
	fmt.Println("\n--- 4.1 Empty Interface Can Hold Anything ---")

	var anything interface{}

	anything = 42
	fmt.Printf("anything = %v (type: %T)\n", anything, anything)

	anything = "hello"
	fmt.Printf("anything = %v (type: %T)\n", anything, anything)

	anything = 3.14
	fmt.Printf("anything = %v (type: %T)\n", anything, anything)

	anything = Dog{Name: "Rex", Breed: "Husky"}
	fmt.Printf("anything = %v (type: %T)\n", anything, anything)

	// ============================================
	// 4.2 استفاده از interface{} در توابع
	// ============================================
	fmt.Println("\n--- 4.2 Using interface{} in Functions ---")

	printAnything(42)
	printAnything("text")
	printAnything(3.14)
	printAnything(Dog{Name: "Max"})

	fmt.Println("\n  Variadic interface{}:")
	printAll(1, "two", 3.0, true, Dog{Name: "Buddy"})

	// ============================================
	// 4.3 Generic Container با interface{}
	// ============================================
	fmt.Println("\n--- 4.3 Generic Container ---")

	container := &Container{}
	container.Add(10)
	container.Add("hello")
	container.Add(3.14)
	container.Add(Dog{Name: "Rex"})

	fmt.Printf("Container has %d items\n", container.Length())
	for i := 0; i < container.Length(); i++ {
		fmt.Printf("  Item %d: %v (type: %T)\n", i, container.Get(i), container.Get(i))
	}
}

// ============================================================================
// بخش 5: Type Assertion و Type Switch
// ============================================================================

func demonstrateTypeAssertion() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🔍 TYPE ASSERTION - Extracting Concrete Type from Interface")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 5.1 Type Assertion پایه
	// ============================================
	fmt.Println("\n--- 5.1 Basic Type Assertion ---")

	var i interface{} = "hello"

	// Safe type assertion with check
	s, ok := i.(string)
	if ok {
		fmt.Printf("  String value: %s\n", s)
	}

	// Unsafe type assertion (panics if wrong type)
	s2 := i.(string)
	fmt.Printf("  Unsafe assertion: %s\n", s2)

	// Type assertion that fails
	f, ok := i.(float64)
	if !ok {
		fmt.Printf("  Not a float64: %v\n", f)
	}

	// ============================================
	// 5.2 Type Switch (بهتر از چندین type assertion)
	// ============================================
	fmt.Println("\n--- 5.2 Type Switch ---")

	checkType := func(v interface{}) {
		switch v := v.(type) {
		case int:
			fmt.Printf("  Integer: %d\n", v)
		case string:
			fmt.Printf("  String: %q\n", v)
		case float64:
			fmt.Printf("  Float64: %.2f\n", v)
		case bool:
			fmt.Printf("  Boolean: %t\n", v)
		case Dog:
			fmt.Printf("  Dog: %s (%s)\n", v.Name, v.Breed)
		case Cat:
			fmt.Printf("  Cat: %s (%s)\n", v.Name, v.Color)
		default:
			fmt.Printf("  Unknown type: %T\n", v)
		}
	}

	checkType(42)
	checkType("hello")
	checkType(3.14)
	checkType(true)
	checkType(Dog{Name: "Max", Breed: "Shepherd"})
	checkType(Cat{Name: "Luna", Color: "Black"})

	// ============================================
	// 5.3 Type Switch در کار با اینترفیس‌ها
	// ============================================
	fmt.Println("\n--- 5.3 Type Switch with Interface Methods ---")

	describeAnimal := func(a Animal) {
		switch a.(type) {
		case Dog:
			fmt.Printf("  This is a Dog: %s\n", a.Speak())
		case Cat:
			fmt.Printf("  This is a Cat: %s\n", a.Speak())
		case Bird:
			fmt.Printf("  This is a Bird: %s\n", a.Speak())
		default:
			fmt.Println("  Unknown animal")
		}
	}

	describeAnimal(Dog{})
	describeAnimal(Cat{})
	describeAnimal(Bird{})
}

// ============================================================================
// بخش 6: اینترفیس‌های Embedding (ترکیب اینترفیس‌ها)
// ============================================================================

// اینترفیس پایه
type Reader interface {
	Read(p []byte) (n int, err error)
}

type Writer interface {
	Write(p []byte) (n int, err error)
}

type Closer interface {
	Close() error
}

// اینترفیس ترکیبی (embedding)
type ReadWriter interface {
	Reader
	Writer
}

type ReadWriteCloser interface {
	Reader
	Writer
	Closer
}

// پیاده‌سازی ساده
type File struct {
	Name string
}

func (f File) Read(p []byte) (n int, err error) {
	fmt.Printf("  Reading from file %s\n", f.Name)
	return len(p), nil
}

func (f File) Write(p []byte) (n int, err error) {
	fmt.Printf("  Writing to file %s\n", f.Name)
	return len(p), nil
}

func (f File) Close() error {
	fmt.Printf("  Closing file %s\n", f.Name)
	return nil
}

func demonstrateInterfaceEmbedding() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🔗 INTERFACE EMBEDDING - Composing Interfaces")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 6.1 Embedding ساده
	// ============================================
	fmt.Println("\n--- 6.1 Basic Embedding ---")

	var rw ReadWriter = File{Name: "data.txt"}
	rw.Read([]byte{})
	rw.Write([]byte{})

	var rwc ReadWriteCloser = File{Name: "log.txt"}
	rwc.Read([]byte{})
	rwc.Write([]byte{})
	rwc.Close()

	// ============================================
	// 6.2 Embedding در اینترفیس‌های سفارشی
	// ============================================
	fmt.Println("\n--- 6.2 Custom Interface Embedding ---")

	type Logger interface {
		Log(message string)
	}

	type Formatter interface {
		Format() string
	}

	type LoggerFormatter interface {
		Logger
		Formatter
	}

	type AppLogger struct {
		prefix string
	}

	func (l AppLogger) Log(message string) {
		fmt.Printf("[%s] %s\n", l.prefix, message)
	}

	func (l AppLogger) Format() string {
		return l.prefix
	}

	var lf LoggerFormatter = AppLogger{prefix: "APP"}
	lf.Log("Starting application")
	fmt.Printf("Logger format: %s\n", lf.Format())
}

// ============================================================================
// بخش 7: اینترفیس‌های معروف در استاندارد Go
// ============================================================================

// Stringer اینترفیس (fmt.Stringer)
type Person2 struct {
	FirstName string
	LastName  string
	Age       int
}

// پیاده‌سازی Stringer برای کنترل چاپ
func (p Person2) String() string {
	return fmt.Sprintf("%s %s (%d years old)", p.FirstName, p.LastName, p.Age)
}

// Error اینترفیس (برای خطاهای سفارشی)
type ValidationError struct {
	Field string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error on %s: %s", e.Field, e.Message)
}

// Custom error creator
func validateAge(age int) error {
	if age < 0 {
		return ValidationError{Field: "age", Message: "cannot be negative"}
	}
	if age > 150 {
		return ValidationError{Field: "age", Message: "cannot be more than 150"}
	}
	return nil
}

func demonstrateStandardInterfaces() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("⭐ STANDARD INTERFACES - Stringer, Error, etc.")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 7.1 fmt.Stringer اینترفیس
	// ============================================
	fmt.Println("\n--- 7.1 fmt.Stringer Interface ---")

	p := Person2{
		FirstName: "Ali",
		LastName:  "Rezaei",
		Age:       30,
	}

	// بدون Stringer: {Ali Rezaei 30}
	// با Stringer: Ali Rezaei (30 years old)
	fmt.Printf("Person: %s\n", p.String())
	fmt.Printf("Person with fmt: %v\n", p)  // fmt از Stringer استفاده می‌کند
	fmt.Printf("Person with fmt: %+v\n", p)

	// ============================================
	// 7.2 error اینترفیس
	// ============================================
	fmt.Println("\n--- 7.2 error Interface ---")

	if err := validateAge(200); err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	if err := validateAge(-5); err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	if err := validateAge(30); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Age is valid")
	}
}

// ============================================================================
// بخش 8: Composition vs Inheritance (ترکیب به جای ارث‌بری)
// ============================================================================

// مثال ترکیب (Composition) - روش Go-way
type Engine struct {
	Horsepower int
	Type       string
}

func (e Engine) Start() string {
	return fmt.Sprintf("Engine started (%d HP %s)", e.Horsepower, e.Type)
}

type Wheels struct {
	Count int
	Type  string
}

func (w Wheels) Rotate() string {
	return fmt.Sprintf("%d wheels rotating", w.Count)
}

type Car struct {
	Brand  string
	Model  string
	Engine // embedding (ترکیب)
	Wheels // embedding
}

func (c Car) Drive() string {
	return fmt.Sprintf("Driving %s %s", c.Brand, c.Model)
}

func demonstrateComposition() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🔨 COMPOSITION - Go's Alternative to Inheritance")
	fmt.Println(stringsRepeat("=", 80))

	car := Car{
		Brand:  "Tesla",
		Model:  "Model 3",
		Engine: Engine{Horsepower: 283, Type: "Electric"},
		Wheels: Wheels{Count: 4, Type: "Alloy"},
	}

	fmt.Println(car.Drive())
	fmt.Println(car.Engine.Start())
	fmt.Println(car.Wheels.Rotate())

	// دسترسی مستقیم به فیلدهای embedded (promoted)
	fmt.Printf("Car specs: %d HP, %d wheels\n", car.Horsepower, car.Count)
}

// ============================================================================
// بخش 9: الگوهای پیشرفته با Struct و Interface
// ============================================================================

// 1. الگوی Constructor (تابع سازنده)
func NewPerson(name string, age int) *Person {
	return &Person{
		Name: name,
		Age:  age,
	}
}

// 2. الگوی Builder Pattern
type PersonBuilder struct {
	person Person
}

func NewPersonBuilder() *PersonBuilder {
	return &PersonBuilder{person: Person{}}
}

func (b *PersonBuilder) WithName(name string) *PersonBuilder {
	b.person.Name = name
	return b
}

func (b *PersonBuilder) WithAge(age int) *PersonBuilder {
	b.person.Age = age
	return b
}

func (b *PersonBuilder) WithEmail(email string) *PersonBuilder {
	b.person.Email = email
	return b
}

func (b *PersonBuilder) Build() Person {
	return b.person
}

// 3. الگوی Options Pattern (با استفاده از closure)
type Server struct {
	host    string
	port    int
	timeout int
	maxConn int
}

type ServerOption func(*Server)

func WithHost(host string) ServerOption {
	return func(s *Server) {
		s.host = host
	}
}

func WithPort(port int) ServerOption {
	return func(s *Server) {
		s.port = port
	}
}

func WithTimeout(timeout int) ServerOption {
	return func(s *Server) {
		s.timeout = timeout
	}
}

func NewServer(opts ...ServerOption) *Server {
	s := &Server{
		host:    "localhost",
		port:    8080,
		timeout: 30,
		maxConn: 100,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// 4. الگوی Interface-based (برای تست و mock)
type UserRepository interface {
	GetUser(id int) (*Person, error)
	SaveUser(user *Person) error
}

type InMemoryUserRepository struct {
	users map[int]Person
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: make(map[int]Person),
	}
}

func (r *InMemoryUserRepository) GetUser(id int) (*Person, error) {
	if user, ok := r.users[id]; ok {
		return &user, nil
	}
	return nil, fmt.Errorf("user not found")
}

func (r *InMemoryUserRepository) SaveUser(user *Person) error {
	// in real code, you'd have an ID
	return nil
}

func demonstrateAdvancedPatterns() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🎨 ADVANCED PATTERNS - Builder, Options, Repository")
	fmt.Println(stringsRepeat("=", 80))

	// Builder Pattern
	fmt.Println("\n--- Builder Pattern ---")
	person := NewPersonBuilder().
		WithName("Ali Rezaei").
		WithAge(30).
		WithEmail("ali@example.com").
		Build()
	fmt.Printf("Built person: %+v\n", person)

	// Options Pattern
	fmt.Println("\n--- Options Pattern ---")
	server1 := NewServer()
	server2 := NewServer(WithHost("0.0.0.0"), WithPort(9090), WithTimeout(60))
	fmt.Printf("Default server: %+v\n", server1)
	fmt.Printf("Custom server: %+v\n", server2)

	// Repository Pattern
	fmt.Println("\n--- Repository Pattern ---")
	repo := NewInMemoryUserRepository()
	_ = repo.SaveUser(&person)
	if user, err := repo.GetUser(1); err == nil {
		fmt.Printf("Found user: %+v\n", user)
	} else {
		fmt.Printf("User not found: %v\n", err)
	}
}

// ============================================================================
// بخش 10: اشتباهات رایج با Struct و Interface
// ============================================================================

func demonstrateCommonMistakes() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("❌ COMMON MISTAKES WITH STRUCTS AND INTERFACES")
	fmt.Println(stringsRepeat("=", 80))

	// اشتباه 1: کپی کردن struct با mutex
	fmt.Println("\n❌ Mistake 1: Copying struct with mutex")
	fmt.Println("   type Counter struct { mu sync.Mutex; count int }")
	fmt.Println("   c2 := c1  // copies mutex!")
	fmt.Println("   ✅ Use pointer: c2 := &c1")

	// اشتباه 2: تعریف اینترفیس از سمت پیاده‌سازی
	fmt.Println("\n❌ Mistake 2: Defining interface on the implementer side")
	fmt.Println("   // Wrong: type Dog interface { Bark() }")
	fmt.Println("   // Define interface where it's USED, not where it's IMPLEMENTED")
	fmt.Println("   ✅ Define interfaces in the consumer package")

	// اشتباه 3: اینترفیس‌های خیلی بزرگ
	fmt.Println("\n❌ Mistake 3: Too large interfaces")
	fmt.Println("   type BigInterface interface { method1(); method2(); ... method20() }")
	fmt.Println("   ✅ Keep interfaces small (1-3 methods) - 'Do one thing'")

	// اشتباه 4: استفاده از interface{} بدون دلیل
	fmt.Println("\n❌ Mistake 4: Using interface{} unnecessarily")
	fmt.Println("   func process(data interface{}) // vague")
	fmt.Println("   ✅ Use specific types or generics (Go 1.18+)")

	// اشتباه 5: فراموش کردن pointer receiver در متدها
	fmt.Println("\n❌ Mistake 5: Forgetting pointer receiver when needed")
	fmt.Println("   type Counter struct { count int }")
	fmt.Println("   func (c Counter) Inc() { c.count++ } // doesn't work!")
	fmt.Println("   ✅ Use pointer: func (c *Counter) Inc()")

	// اشتباه 6: nil interface vs nil concrete type
	fmt.Println("\n❌ Mistake 6: nil interface vs nil concrete type")
	fmt.Println("   var p *Person = nil")
	fmt.Println("   var i interface{} = p")
	fmt.Println("   fmt.Println(i == nil) // false!")
	fmt.Println("   ✅ Check both: i == nil || reflect.ValueOf(i).IsNil()")
}

// ============================================================================
// بخش 11: جمع‌بندی و جدول مرجع سریع
// ============================================================================

func main() {
	fmt.Println(stringsRepeat("=", 80))
	fmt.Println("🎯 COMPLETE STRUCTS & INTERFACES GUIDE IN GO")
	fmt.Println(stringsRepeat("=", 80))

	// بخش 1: Struct مبانی
	demonstrateStructBasics()

	// بخش 2: متدها روی Struct
	demonstrateStructMethods()

	// بخش 3: اینترفیس‌ها
	demonstrateInterfaces()

	// بخش 4: Empty Interface
	demonstrateEmptyInterface()

	// بخش 5: Type Assertion & Type Switch
	demonstrateTypeAssertion()

	// بخش 6: Interface Embedding
	demonstrateInterfaceEmbedding()

	// بخش 7: Standard Interfaces
	demonstrateStandardInterfaces()

	// بخش 8: Composition
	demonstrateComposition()

	// بخش 9: الگوهای پیشرفته
	demonstrateAdvancedPatterns()

	// بخش 10: اشتباهات رایج
	demonstrateCommonMistakes()

	// بخش 11: جمع‌بندی
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📚 STRUCTS & INTERFACES QUICK REFERENCE")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ STRUCTS                                                         │")
	fmt.Println("├─────────────────────────────────────────────────────────────────┤")
	fmt.Println("│ type Person struct { Name string; Age int }                     │")
	fmt.Println("│ p := Person{Name: \"Ali\", Age: 30}     // initialization      │")
	fmt.Println("│ p := &Person{Name: \"Ali\"}             // pointer             │")
	fmt.Println("│ p.Name = \"Reza\"                        // access field        │")
	fmt.Println("│ func (p Person) Greet() string {...}   // value receiver       │")
	fmt.Println("│ func (p *Person) SetAge(a int) {...}   // pointer receiver     │")
	fmt.Println("└─────────────────────────────────────────────────────────────────┘")

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ INTERFACES                                                      │")
	fmt.Println("├─────────────────────────────────────────────────────────────────┤")
	fmt.Println("│ type Greeter interface { Greet() string }                       │")
	fmt.Println("│ // Automatic implementation (no 'implements' keyword)          │")
	fmt.Println("│ var g Greeter = Dog{}          // Dog implements Greeter       │")
	fmt.Println("│ var empty interface{}          // can hold any type            │")
	fmt.Println("│ s, ok := i.(string)            // type assertion               │")
	fmt.Println("│ switch v := i.(type) {         // type switch                  │")
	fmt.Println("│ type ReadWriter interface { Reader; Writer }  // embedding     │")
	fmt.Println("└─────────────────────────────────────────────────────────────────┘")

	fmt.Println("\n💡 GOLDEN RULES:")
	fmt.Println("  1. Use structs for data, interfaces for behavior")
	fmt.Println("  2. Keep interfaces small (1-3 methods)")
	fmt.Println("  3. Define interfaces where they are USED, not implemented")
	fmt.Println("  4. Use pointer receivers for methods that modify state")
	fmt.Println("  5. Use value receivers for small, immutable types")
	fmt.Println("  6. Prefer composition over inheritance (embedding)")
	fmt.Println("  7. Empty interface (interface{}) is NOT 'any' - use with care")
	fmt.Println("  8. Type assertions should always check ok (safe assertion)")
	fmt.Println("  9. Interfaces are satisfied implicitly (don't force it)")
	fmt.Println("  10. Accept interfaces, return structs (common idiom)")

	fmt.Println("\n🎯 WHEN TO USE WHAT:")
	fmt.Println("  • Struct: Group related data together")
	fmt.Println("  • Methods: Attach behavior to structs")
	fmt.Println("  • Interface: Define contracts and enable polymorphism")
	fmt.Println("  • Embedding: Compose behaviors (not inheritance)")
	fmt.Println("  • Type Assertion: Extract concrete type from interface")
	fmt.Println("  • Type Switch: Handle multiple types gracefully")
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