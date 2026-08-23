package __internal_packages

// ============================================================================
// FILE: reflect_complete_guide.go
// TITLE: راهنمای کامل پکیج reflect در Go - بازتاب و متاپروگرامینگ
// HOW TO RUN: go run reflect_complete_guide.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - پکیج reflect چیست و چه کاربردی دارد؟
// ============================================================================
//
// پکیج reflect (بازتاب) به برنامه اجازه می‌دهد در زمان اجرا (runtime) به نوع و مقدار
// متغیرها دسترسی پیدا کند و حتی آن‌ها را تغییر دهد.
//
// کاربردهای اصلی:
// 1. بررسی نوع و مقدار متغیرها در زمان اجرا
// 2. فراخوانی داینامیک توابع و متدها
// 3. ساخت توابع عمومی (generic functions) قبل از Go 1.18
// 4. سریالایز/دی سریالایز کردن داده‌ها (JSON, XML, BSON)
// 5. نوشتن ابزارهای دیباگ و تست
// 6. Dependency Injection frameworks
//
// قانون طلایی:
// "reflect قدرتمند است اما هزینه دارد. هر زمان که می‌توانی بدون reflect کار کن،
//  از آن استفاده نکن. فقط زمانی از reflect استفاده کن که راه دیگری وجود ندارد."
// ============================================================================


import (
"fmt"
"reflect"
"strings"
)

func main() {
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("🎯 COMPLETE reflect PACKAGE GUIDE IN GO")
	fmt.Println("Reflection and Metaprogramming")
	fmt.Println(strings.Repeat("=", 80))

	// ============================================================================
	// بخش 1: انواع و مقادیر (Type and Value)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📊 SECTION 1: reflect.Type AND reflect.Value")
	fmt.Println(strings.Repeat("=", 80))

	// 1.1 TypeOf - گرفتن نوع متغیر
	fmt.Println("\n--- 1.1 reflect.TypeOf ---")
	// TypeOf returns the reflection Type that represents the dynamic type of i.
	var x int = 42
	var y string = "hello"
	var z float64 = 3.14

	tx := reflect.TypeOf(x)
	ty := reflect.TypeOf(y)
	tz := reflect.TypeOf(z)

	fmt.Printf("  Type of x (%v): %v, Kind: %v\n", x, tx, tx.Kind())
	fmt.Printf("  Type of y (%v): %v, Kind: %v\n", y, ty, ty.Kind())
	fmt.Printf("  Type of z (%v): %v, Kind: %v\n", z, tz, tz.Kind())

	// 1.2 ValueOf - گرفتن مقدار متغیر
	fmt.Println("\n--- 1.2 reflect.ValueOf ---")
	// ValueOf returns a new Value initialized to the concrete value stored in the interface i.
	vx := reflect.ValueOf(x)
	vy := reflect.ValueOf(y)
	vz := reflect.ValueOf(z)

	fmt.Printf("  Value of x: %v, Type: %v, Kind: %v\n", vx, vx.Type(), vx.Kind())
	fmt.Printf("  Value of y: %v, Type: %v, Kind: %v\n", vy, vy.Type(), vy.Kind())
	fmt.Printf("  Value of z: %v, Type: %v, Kind: %v\n", vz, vz.Type(), vz.Kind())

	// 1.3 Kind - انواع پایه
	fmt.Println("\n--- 1.3 reflect.Kind ---")
	// Kind represents the specific kind of type that a Type represents.
	var (
		intVal    int     = 42
		floatVal  float64 = 3.14
		boolVal   bool    = true
		stringVal string  = "hello"
		sliceVal  []int   = []int{1, 2, 3}
		mapVal    map[string]int
		structVal struct{ Name string }
		funcVal   func()
		ptrVal    *int
	)

	printKind := func(v interface{}) {
		t := reflect.TypeOf(v)
		fmt.Printf("  Value: %v, Type: %v, Kind: %v\n", v, t, t.Kind())
	}

	printKind(intVal)
	printKind(floatVal)
	printKind(boolVal)
	printKind(stringVal)
	printKind(sliceVal)
	printKind(mapVal)
	printKind(structVal)
	printKind(funcVal)
	printKind(ptrVal)

	// ============================================================================
	// بخش 2: کار با Structها (بازتاب فیلدها و تگ‌ها)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🏗️ SECTION 2: STRUCT REFLECTION")
	fmt.Println(strings.Repeat("=", 80))

	// تعریف یک struct نمونه با تگ‌ها
	type User struct {
		ID        int    `json:"id" db:"user_id" validate:"required"`
		Name      string `json:"name" db:"user_name" validate:"required,min=3"`
		Email     string `json:"email,omitempty" db:"email" validate:"email"`
		Password  string `json:"-" db:"password_hash"`
		IsActive  bool   `json:"is_active" db:"is_active"`
		CreatedAt string `json:"created_at" db:"created_at"`
	}

	// 2.1 بررسی فیلدهای Struct
	fmt.Println("\n--- 2.1 Inspecting Struct Fields ---")

	user := User{
		ID:        1,
		Name:      "Ali Rezaei",
		Email:     "ali@example.com",
		Password:  "secret",
		IsActive:  true,
		CreatedAt: "2024-01-15",
	}

	t := reflect.TypeOf(user)
	v := reflect.ValueOf(user)

	fmt.Printf("  Struct name: %s\n", t.Name())
	fmt.Printf("  Number of fields: %d\n", t.NumField())

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		fmt.Printf("\n  Field %d:\n", i)
		fmt.Printf("    Name: %s\n", field.Name)
		fmt.Printf("    Type: %v\n", field.Type)
		fmt.Printf("    Value: %v\n", value.Interface())
		fmt.Printf("    Tags: %v\n", field.Tag)

		// خواندن تگ‌های خاص
		jsonTag := field.Tag.Get("json")
		dbTag := field.Tag.Get("db")
		validateTag := field.Tag.Get("validate")

		if jsonTag != "" {
			fmt.Printf("    json tag: %s\n", jsonTag)
		}
		if dbTag != "" {
			fmt.Printf("    db tag: %s\n", dbTag)
		}
		if validateTag != "" {
			fmt.Printf("    validate tag: %s\n", validateTag)
		}
	}

	// 2.2 اصلاح فیلدهای Struct (با اشاره‌گر)
	fmt.Println("\n--- 2.2 Modifying Struct Fields (with pointer) ---")

	userPtr := &User{Name: "Original Name", IsActive: false}
	vPtr := reflect.ValueOf(userPtr).Elem()

	fmt.Printf("  Before modification: %+v\n", userPtr)

	// پیدا کردن فیلد Name و تغییر آن
	nameField := vPtr.FieldByName("Name")
	if nameField.IsValid() && nameField.CanSet() {
		nameField.SetString("New Name")
	}

	// پیدا کردن فیلد IsActive و تغییر آن
	activeField := vPtr.FieldByName("IsActive")
	if activeField.IsValid() && activeField.CanSet() {
		activeField.SetBool(true)
	}

	fmt.Printf("  After modification: %+v\n", userPtr)

	// 2.3 پیمایش فیلدها با FieldByName
	fmt.Println("\n--- 2.3 FieldByName ---")

	getFieldValue := func(s interface{}, fieldName string) (interface{}, bool) {
		v := reflect.ValueOf(s)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		field := v.FieldByName(fieldName)
		if !field.IsValid() {
			return nil, false
		}
		return field.Interface(), true
	}

	if val, ok := getFieldValue(user, "Email"); ok {
		fmt.Printf("  Email field value: %v\n", val)
	}
	if val, ok := getFieldValue(user, "NonExistent"); ok {
		fmt.Printf("  NonExistent: %v\n", val)
	} else {
		fmt.Println("  NonExistent field not found")
	}

	// ============================================================================
	// بخش 3: کار با Slice و Array
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📊 SECTION 3: SLICE AND ARRAY REFLECTION")
	fmt.Println(strings.Repeat("=", 80))

	// 3.1 بررسی Slice
	fmt.Println("\n--- 3.1 Inspecting Slice ---")

	slice := []int{10, 20, 30, 40, 50}
	vSlice := reflect.ValueOf(slice)

	fmt.Printf("  Slice: %v\n", slice)
	fmt.Printf("  Type: %v\n", vSlice.Type())
	fmt.Printf("  Kind: %v\n", vSlice.Kind())
	fmt.Printf("  Length: %d\n", vSlice.Len())
	fmt.Printf("  Capacity: %d\n", vSlice.Cap())

	// پیمایش عناصر Slice
	fmt.Println("  Elements:")
	for i := 0; i < vSlice.Len(); i++ {
		elem := vSlice.Index(i)
		fmt.Printf("    [%d]: %v (type: %v)\n", i, elem.Interface(), elem.Type())
	}

	// 3.2 ایجاد و اصلاح Slice
	fmt.Println("\n--- 3.2 Creating and Modifying Slice ---")

	// ایجاد Slice جدید با reflect
	sliceType := reflect.TypeOf([]int{})
	newSlice := reflect.MakeSlice(sliceType, 0, 5)
	fmt.Printf("  New slice: %v (len=%d, cap=%d)\n", newSlice, newSlice.Len(), newSlice.Cap())

	// افزودن عناصر به Slice
	newSlice = reflect.Append(newSlice, reflect.ValueOf(100))
	newSlice = reflect.Append(newSlice, reflect.ValueOf(200))
	newSlice = reflect.Append(newSlice, reflect.ValueOf(300))
	fmt.Printf("  After appending: %v\n", newSlice)

	// تغییر عنصر
	newSlice.Index(1).SetInt(999)
	fmt.Printf("  After modification: %v\n", newSlice)

	// 3.3 بررسی Array
	fmt.Println("\n--- 3.3 Inspecting Array ---")

	arr := [5]int{1, 2, 3, 4, 5}
	vArr := reflect.ValueOf(arr)

	fmt.Printf("  Array: %v\n", arr)
	fmt.Printf("  Type: %v\n", vArr.Type())
	fmt.Printf("  Length: %d\n", vArr.Len())

	// پیمایش عناصر Array
	fmt.Println("  Elements:")
	for i := 0; i < vArr.Len(); i++ {
		fmt.Printf("    [%d]: %v\n", i, vArr.Index(i).Interface())
	}

	// ============================================================================
	// بخش 4: کار با Map
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🗺️ SECTION 4: MAP REFLECTION")
	fmt.Println(strings.Repeat("=", 80))

	// 4.1 بررسی Map
	fmt.Println("\n--- 4.1 Inspecting Map ---")

	m := map[string]int{
		"apple":  5,
		"banana": 3,
		"cherry": 7,
	}

	vMap := reflect.ValueOf(m)

	fmt.Printf("  Map: %v\n", m)
	fmt.Printf("  Type: %v\n", vMap.Type())
	fmt.Printf("  Kind: %v\n", vMap.Kind())
	fmt.Printf("  Length: %d\n", vMap.Len())

	// پیمایش Map
	fmt.Println("  Key-Value pairs:")
	iter := vMap.MapRange()
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()
		fmt.Printf("    %v: %v\n", key.Interface(), value.Interface())
	}

	// 4.2 ایجاد و اصلاح Map
	fmt.Println("\n--- 4.2 Creating and Modifying Map ---")

	// ایجاد Map جدید
	mapType := reflect.TypeOf(map[string]int{})
	newMap := reflect.MakeMap(mapType)

	// افزودن key-value
	newMap.SetMapIndex(reflect.ValueOf("x"), reflect.ValueOf(10))
	newMap.SetMapIndex(reflect.ValueOf("y"), reflect.ValueOf(20))
	newMap.SetMapIndex(reflect.ValueOf("z"), reflect.ValueOf(30))

	fmt.Printf("  New map: %v\n", newMap)

	// حذف key
	newMap.SetMapIndex(reflect.ValueOf("y"), reflect.Value{})
	fmt.Printf("  After deleting 'y': %v\n", newMap)

	// ============================================================================
	// بخش 5: فراخوانی توابع و متدها با reflect
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🔧 SECTION 5: CALLING FUNCTIONS AND METHODS")
	fmt.Println(strings.Repeat("=", 80))

	// 5.1 فراخوانی توابع معمولی
	fmt.Println("\n--- 5.1 Calling Regular Functions ---")

	// تعریف توابع نمونه
	add := func(a, b int) int {
		return a + b
	}

	greet := func(name string) string {
		return "Hello, " + name
	}

	// فراخوانی add
	vAdd := reflect.ValueOf(add)
	args := []reflect.Value{reflect.ValueOf(10), reflect.ValueOf(20)}
	results := vAdd.Call(args)
	fmt.Printf("  add(10, 20) = %v\n", results[0].Interface())

	// فراخوانی greet
	vGreet := reflect.ValueOf(greet)
	args2 := []reflect.Value{reflect.ValueOf("Ali")}
	results2 := vGreet.Call(args2)
	fmt.Printf("  greet(\"Ali\") = %v\n", results2[0].Interface())

	// 5.2 فراخوانی توابع با تعداد متغیر آرگومان (Variadic)
	fmt.Println("\n--- 5.2 Calling Variadic Functions ---")

	sum := func(nums ...int) int {
		total := 0
		for _, n := range nums {
			total += n
		}
		return total
	}

	vSum := reflect.ValueOf(sum)
	args3 := []reflect.Value{reflect.ValueOf([]int{1, 2, 3, 4, 5})}
	results3 := vSum.Call(args3)
	fmt.Printf("  sum(1,2,3,4,5) = %v\n", results3[0].Interface())

	// 5.3 فراخوانی متدهای Struct
	fmt.Println("\n--- 5.3 Calling Struct Methods ---")

	type Calculator struct {
		Value int
	}

	func (c *Calculator) Add(n int) int {
		c.Value += n
		return c.Value
	}

	func (c Calculator) Multiply(n int) int {
		return c.Value * n
	}

	calc := &Calculator{Value: 10}
	vCalc := reflect.ValueOf(calc)

	// فراخوانی متد Add (pointer receiver)
	addMethod := vCalc.MethodByName("Add")
	if addMethod.IsValid() {
		args4 := []reflect.Value{reflect.ValueOf(5)}
		result4 := addMethod.Call(args4)
		fmt.Printf("  Add(5) called via reflect: %v\n", result4[0].Interface())
		fmt.Printf("  Calculator.Value after: %d\n", calc.Value)
	}

	// فراخوانی متد Multiply (value receiver)
	multiplyMethod := vCalc.MethodByName("Multiply")
	if multiplyMethod.IsValid() {
		args5 := []reflect.Value{reflect.ValueOf(3)}
		result5 := multiplyMethod.Call(args5)
		fmt.Printf("  Multiply(3) called via reflect: %v\n", result5[0].Interface())
	}

	// 5.4 بررسی وجود متد
	fmt.Println("\n--- 5.4 Checking Method Existence ---")

	hasMethod := func(s interface{}, methodName string) bool {
		v := reflect.ValueOf(s)
		method := v.MethodByName(methodName)
		return method.IsValid()
	}

	fmt.Printf("  Calculator has Add method: %v\n", hasMethod(calc, "Add"))
	fmt.Printf("  Calculator has Subtract method: %v\n", hasMethod(calc, "Subtract"))

	// ============================================================================
	// بخش 6: ایجاد انواع جدید در زمان اجرا
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🔨 SECTION 6: CREATING NEW TYPES AT RUNTIME")
	fmt.Println(strings.Repeat("=", 80))

	// 6.1 ایجاد Struct داینامیک
	fmt.Println("\n--- 6.1 Creating Dynamic Struct ---")

	fields := []reflect.StructField{
		{
			Name: "Name",
			Type: reflect.TypeOf(""),
			Tag:  `json:"name"`,
		},
		{
			Name: "Age",
			Type: reflect.TypeOf(0),
			Tag:  `json:"age"`,
		},
		{
			Name: "Active",
			Type: reflect.TypeOf(true),
			Tag:  `json:"active"`,
		},
	}

	dynamicType := reflect.StructOf(fields)
	dynamicValue := reflect.New(dynamicType).Elem()

	dynamicValue.FieldByName("Name").SetString("Dynamic Person")
	dynamicValue.FieldByName("Age").SetInt(30)
	dynamicValue.FieldByName("Active").SetBool(true)

	fmt.Printf("  Dynamic struct type: %v\n", dynamicType)
	fmt.Printf("  Dynamic struct value: %+v\n", dynamicValue.Interface())

	// تبدیل به interface معمولی
	var dynamic interface{} = dynamicValue.Interface()
	fmt.Printf("  As interface: %+v\n", dynamic)

	// 6.2 ایجاد Slice از نوع داینامیک
	fmt.Println("\n--- 6.2 Creating Dynamic Slice ---")

	elemType := reflect.TypeOf(0) // int
	sliceType2 := reflect.SliceOf(elemType)
	sliceValue := reflect.MakeSlice(sliceType2, 0, 5)

	for i := 0; i < 5; i++ {
		sliceValue = reflect.Append(sliceValue, reflect.ValueOf(i*10))
	}

	fmt.Printf("  Dynamic slice: %v\n", sliceValue)

	// 6.3 ایجاد Map داینامیک
	fmt.Println("\n--- 6.3 Creating Dynamic Map ---")

	keyType := reflect.TypeOf("")
	valueType := reflect.TypeOf(0)
	mapType2 := reflect.MapOf(keyType, valueType)
	mapValue := reflect.MakeMap(mapType2)

	mapValue.SetMapIndex(reflect.ValueOf("a"), reflect.ValueOf(1))
	mapValue.SetMapIndex(reflect.ValueOf("b"), reflect.ValueOf(2))
	mapValue.SetMapIndex(reflect.ValueOf("c"), reflect.ValueOf(3))

	fmt.Printf("  Dynamic map: %v\n", mapValue)

	// ============================================================================
	// بخش 7: بررسی و مدیریت خطاها
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("⚠️ SECTION 7: ERROR HANDLING AND VALIDATION")
	fmt.Println(strings.Repeat("=", 80))

	// 7.1 اعتبارسنجی Struct با تگ‌ها
	fmt.Println("\n--- 7.1 Struct Validation with Tags ---")

	type ValidatedUser struct {
		Name  string `validate:"required,min=3,max=50"`
		Age   int    `validate:"min=0,max=150"`
		Email string `validate:"required,email"`
	}

	validate := func(s interface{}) []string {
		var errors []string
		v := reflect.ValueOf(s)
		t := reflect.TypeOf(s)

		if v.Kind() == reflect.Ptr {
			v = v.Elem()
			t = t.Elem()
		}

		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			value := v.Field(i)
			tag := field.Tag.Get("validate")

			if tag == "" {
				continue
			}

			// بررسی required
			if strings.Contains(tag, "required") && value.IsZero() {
				errors = append(errors, fmt.Sprintf("%s is required", field.Name))
			}

			// بررسی min length برای رشته
			if strings.Contains(tag, "min") && value.Kind() == reflect.String {
				// در عمل باید tag را parse کرد
				if len(value.String()) < 3 {
					errors = append(errors, fmt.Sprintf("%s length must be at least 3", field.Name))
				}
			}

			// بررسی min/max برای اعداد
			if strings.Contains(tag, "min") && value.Kind() == reflect.Int {
				if value.Int() < 0 {
					errors = append(errors, fmt.Sprintf("%s must be >= 0", field.Name))
				}
			}
		}

		return errors
	}

	validUser := ValidatedUser{Name: "Ali", Age: 30, Email: "ali@test.com"}
	invalidUser := ValidatedUser{Name: "A", Age: -5, Email: ""}

	fmt.Printf("  Valid user errors: %v\n", validate(validUser))
	fmt.Printf("  Invalid user errors: %v\n", validate(invalidUser))

	// 7.2 بررسی nil بودن
	fmt.Println("\n--- 7.2 Nil Check ---")

	var ptr *int
	var sliceNil []int
	var mapNil map[string]int

	vPtr := reflect.ValueOf(ptr)
	vSliceNil := reflect.ValueOf(sliceNil)
	vMapNil := reflect.ValueOf(mapNil)

	fmt.Printf("  ptr is nil: %v\n", vPtr.IsNil())
	fmt.Printf("  sliceNil is nil: %v\n", vSliceNil.IsNil())
	fmt.Printf("  mapNil is nil: %v\n", vMapNil.IsNil())

	// ============================================================================
	// بخش 8: تبدیل انواع (Type Conversion)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🔄 SECTION 8: TYPE CONVERSION")
	fmt.Println(strings.Repeat("=", 80))

	// 8.1 تبدیل به انواع پایه
	fmt.Println("\n--- 8.1 Converting to Basic Types ---")

	val := reflect.ValueOf(42)

	fmt.Printf("  Int: %d\n", val.Int())
	fmt.Printf("  Float: %f\n", val.Float())
	fmt.Printf("  Interface: %v\n", val.Interface())

	strVal := reflect.ValueOf("hello")
	fmt.Printf("  String: %s\n", strVal.String())

	// 8.2 تبدیل بین انواع مختلف
	fmt.Println("\n--- 8.2 Converting Between Types ---")

	convert := func(v interface{}, targetKind reflect.Kind) interface{} {
		val := reflect.ValueOf(v)
		switch targetKind {
		case reflect.String:
			return fmt.Sprintf("%v", v)
		case reflect.Int:
			return int(val.Int())
		case reflect.Float64:
			return val.Float()
		default:
			return v
		}
	}

	fmt.Printf("  Convert 42 to string: %v (type: %T)\n", convert(42, reflect.String), convert(42, reflect.String))
	fmt.Printf("  Convert 3.14 to int: %v (type: %T)\n", convert(3.14, reflect.Int), convert(3.14, reflect.Int))

	// ============================================================================
	// بخش 9: بررسی قابلیت‌های Value (CanSet, CanAddr, etc.)
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🔍 SECTION 9: VALUE CAPABILITIES")
	fmt.Println(strings.Repeat("=", 80))

	// 9.1 بررسی قابلیت تغییر
	fmt.Println("\n--- 9.1 CanSet and CanAddr ---")

	num := 42
	ptrNum := &num

	vNum := reflect.ValueOf(num)
	vPtrNum := reflect.ValueOf(ptrNum)
	vElem := vPtrNum.Elem()

	fmt.Printf("  Value of num: canSet=%v, canAddr=%v\n", vNum.CanSet(), vNum.CanAddr())
	fmt.Printf("  Value of ptr: canSet=%v, canAddr=%v\n", vPtrNum.CanSet(), vPtrNum.CanAddr())
	fmt.Printf("  Value of elem: canSet=%v, canAddr=%v\n", vElem.CanSet(), vElem.CanAddr())

	// تغییر مقدار از طریق اشاره‌گر
	if vElem.CanSet() {
		vElem.SetInt(100)
		fmt.Printf("  After modification: num = %d\n", num)
	}

	// 9.2 بررسی نوع‌های مختلف
	fmt.Println("\n--- 9.2 Type Checks ---")

	checkValue := func(v interface{}) {
		val := reflect.ValueOf(v)
		fmt.Printf("  Value: %v, Kind: %v\n", v, val.Kind())
		fmt.Printf("    IsValid: %v\n", val.IsValid())
		fmt.Printf("    IsZero: %v\n", val.IsZero())
		fmt.Printf("    IsNil (if applicable): %v\n", (val.Kind() == reflect.Ptr ||
			val.Kind() == reflect.Slice || val.Kind() == reflect.Map) && val.IsNil())
	}

	checkValue(42)
	checkValue("hello")
	checkValue(nil)
	checkValue([]int(nil))

	// ============================================================================
	// بخش 10: کاربردهای عملی
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("💡 SECTION 10: PRACTICAL EXAMPLES")
	fmt.Println(strings.Repeat("=", 80))

	// 10.1 کپی عمیق (Deep Copy)
	fmt.Println("\n--- 10.1 Deep Copy ---")

	deepCopy := func(src interface{}) interface{} {
		srcVal := reflect.ValueOf(src)
		if srcVal.Kind() == reflect.Ptr {
			srcVal = srcVal.Elem()
		}

		dstVal := reflect.New(srcVal.Type()).Elem()
		deepCopyRecursive(srcVal, dstVal)
		return dstVal.Interface()
	}

	deepCopyRecursive := func(src, dst reflect.Value) {
		switch src.Kind() {
		case reflect.Struct:
			for i := 0; i < src.NumField(); i++ {
				deepCopyRecursive(src.Field(i), dst.Field(i))
			}
		case reflect.Slice:
			dst.Set(reflect.MakeSlice(src.Type(), src.Len(), src.Cap()))
			for i := 0; i < src.Len(); i++ {
				deepCopyRecursive(src.Index(i), dst.Index(i))
			}
		case reflect.Map:
			dst.Set(reflect.MakeMap(src.Type()))
			iter := src.MapRange()
			for iter.Next() {
				key := reflect.New(iter.Key().Type()).Elem()
				value := reflect.New(iter.Value().Type()).Elem()
				deepCopyRecursive(iter.Key(), key)
				deepCopyRecursive(iter.Value(), value)
				dst.SetMapIndex(key, value)
			}
		default:
			dst.Set(src)
		}
	}

	originalStruct := struct {
		Name string
		Age  int
		Tags []string
	}{
		Name: "Original",
		Age:  30,
		Tags: []string{"a", "b", "c"},
	}

	copied := deepCopy(originalStruct)
	fmt.Printf("  Original: %+v\n", originalStruct)
	fmt.Printf("  Copied: %+v\n", copied)

	// 10.2 تبدیل Struct به Map
	fmt.Println("\n--- 10.2 Struct to Map Conversion ---")

	structToMap := func(s interface{}) map[string]interface{} {
		result := make(map[string]interface{})
		v := reflect.ValueOf(s)
		t := reflect.TypeOf(s)

		if v.Kind() == reflect.Ptr {
			v = v.Elem()
			t = t.Elem()
		}

		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			value := v.Field(i)
			jsonTag := field.Tag.Get("json")

			key := field.Name
			if jsonTag != "" && jsonTag != "-" {
				key = strings.Split(jsonTag, ",")[0]
			}

			result[key] = value.Interface()
		}

		return result
	}

	userMap := structToMap(user)
	fmt.Printf("  Struct to map: %+v\n", userMap)

	// 10.3 Map به Struct
	fmt.Println("\n--- 10.3 Map to Struct Conversion ---")

	mapToStruct := func(m map[string]interface{}, s interface{}) error {
		v := reflect.ValueOf(s)
		if v.Kind() != reflect.Ptr || v.IsNil() {
			return fmt.Errorf("s must be a non-nil pointer")
		}

		v = v.Elem()
		t := v.Type()

		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			jsonTag := field.Tag.Get("json")

			key := field.Name
			if jsonTag != "" && jsonTag != "-" {
				key = strings.Split(jsonTag, ",")[0]
			}

			if val, ok := m[key]; ok {
				fieldVal := v.Field(i)
				if fieldVal.CanSet() {
					fieldVal.Set(reflect.ValueOf(val))
				}
			}
		}

		return nil
	}

	newUser := &User{}
	mapToStruct(userMap, newUser)
	fmt.Printf("  Map to struct: %+v\n", newUser)

	// ============================================================================
	// بخش 11: اشتباهات رایج
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("❌ SECTION 11: COMMON MISTAKES")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println("\n❌ Mistake 1: Calling Set on unexported field")
	fmt.Println("   v.Field(i).Set(newValue)  // panics if field is unexported")
	fmt.Println("   ✅ Check CanSet() before calling Set()")

	fmt.Println("\n❌ Mistake 2: Using reflect on non-pointer when modification needed")
	fmt.Println("   v := reflect.ValueOf(s)  // s is not pointer")
	fmt.Println("   v.Set(newValue)  // panics")
	fmt.Println("   ✅ Use reflect.ValueOf(&s).Elem()")

	fmt.Println("\n❌ Mistake 3: Comparing reflect.Value with == directly")
	fmt.Println("   v1 == v2  // compares reflect.Value structs, not underlying values")
	fmt.Println("   ✅ Use v1.Interface() == v2.Interface() or v1.Equal(v2)")

	fmt.Println("\n❌ Mistake 4: Assuming reflect.DeepEqual works for all types")
	fmt.Println("   reflect.DeepEqual(func1, func2)  // always false for functions")
	fmt.Println("   ✅ Use reflect.Value.Pointer() for function comparison")

	fmt.Println("\n❌ Mistake 5: Performance issues in hot paths")
	fmt.Println("   reflect operations are slow (10-100x slower)")
	fmt.Println("   ✅ Use type switches or generics when possible")

	fmt.Println("\n❌ Mistake 6: Using reflect on nil interface")
	fmt.Println("   var i interface{} = nil")
	fmt.Println("   reflect.ValueOf(i)  // invalid Value")
	fmt.Println("   ✅ Check IsValid() before using Value")

	// ============================================================================
	// بخش 12: جمع‌بندی و جدول مرجع
	// ============================================================================
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📚 SECTION 12: QUICK REFERENCE")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ FUNCTION                │ DESCRIPTION                          │")
	fmt.Println("├─────────────────────────┼──────────────────────────────────────┤")
	fmt.Println("│ reflect.TypeOf(i)       │ Returns the reflection Type of i     │")
	fmt.Println("│ reflect.ValueOf(i)      │ Returns the reflection Value of i    │")
	fmt.Println("│ reflect.New(t)          │ Returns pointer to new zero value    │")
	fmt.Println("│ reflect.Zero(t)         │ Returns zero value of type t         │")
	fmt.Println("│ reflect.DeepEqual(a,b)  │ Deep equality comparison             │")
	fmt.Println("│ reflect.SliceOf(t)      │ Returns slice type of element t      │")
	fmt.Println("│ reflect.MapOf(k,v)      │ Returns map type with key k, value v │")
	fmt.Println("│ reflect.StructOf(fields)│ Returns struct type from field list  │")
	fmt.Println("└─────────────────────────┴──────────────────────────────────────┘")

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ reflect.Value METHODS                                        │")
	fmt.Println("├─────────────────────────────────────────────────────────────────┤")
	fmt.Println("│ v.Kind()              - Returns the kind of the value         │")
	fmt.Println("│ v.Type()              - Returns the type of the value         │")
	fmt.Println("│ v.Interface()         - Returns value as interface{}          │")
	fmt.Println("│ v.IsNil()             - Reports whether v is nil              │")
	fmt.Println("│ v.IsValid()           - Reports whether v represents a value  │")
	fmt.Println("│ v.IsZero()            - Reports whether v is zero value       │")
	fmt.Println("│ v.CanSet()            - Reports whether value can be changed  │")
	fmt.Println("│ v.CanAddr()           - Reports whether value can be addressed│")
	fmt.Println("│ v.Elem()              - Returns value that v points to        │")
	fmt.Println("│ v.NumField()          - Returns number of struct fields       │")
	fmt.Println("│ v.Field(i)            - Returns struct field i                │")
	fmt.Println("│ v.FieldByName(name)   - Returns struct field by name          │")
	fmt.Println("│ v.NumMethod()         - Returns number of methods             │")
	fmt.Println("│ v.Method(i)           - Returns method i                      │")
	fmt.Println("│ v.MethodByName(name)  - Returns method by name                │")
	fmt.Println("│ v.Call(args)          - Calls function with arguments         │")
	fmt.Println("│ v.Len()               - Returns length (for array/slice/map)  │")
	fmt.Println("│ v.Cap()               - Returns capacity (for array/slice)    │")
	fmt.Println("│ v.Index(i)            - Returns element i (for array/slice)   │")
	fmt.Println("│ v.MapRange()          - Returns iterator for map              │")
	fmt.Println("│ v.SetMapIndex(k, v)   - Sets map key to value                 │")
	fmt.Println("└─────────────────────────────────────────────────────────────────┘")

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ reflect.Kind VALUES                                          │")
	fmt.Println("├─────────────────────────────────────────────────────────────────┤")
	fmt.Println("│ Bool, Int, Int8, Int16, Int32, Int64, Uint, Uint8, Uint16,   │")
	fmt.Println("│ Uint32, Uint64, Uintptr, Float32, Float64, Complex64,        │")
	fmt.Println("│ Complex128, Array, Chan, Func, Interface, Map, Ptr, Slice,   │")
	fmt.Println("│ String, Struct, UnsafePointer                                 │")
	fmt.Println("└─────────────────────────────────────────────────────────────────┘")

	fmt.Println("\n💡 GOLDEN RULES:")
	fmt.Println("  1. Always check CanSet() before modifying values")
	fmt.Println("  2. Use pointer when you need to modify original value")
	fmt.Println("  3. reflect.DeepEqual is convenient but slower than manual comparison")
	fmt.Println("  4. Avoid reflect in performance-critical code paths")
	fmt.Println("  5. Use type switches when types are known at compile time")
	fmt.Println("  6. Check IsValid() before using Value methods")
	fmt.Println("  7. Unexported struct fields cannot be accessed or modified")
	fmt.Println("  8. reflect.Call is powerful but expensive")
	fmt.Println("  9. Use reflect.MakeSlice/Map for dynamic collection creation")
	fmt.Println("  10. When in doubt, avoid reflect and use interfaces")

	fmt.Println("\n🎯 WHEN TO USE REFLECT:")
	fmt.Println("  • JSON/XML/BSON serialization/deserialization")
	fmt.Println("  • ORM (Object-Relational Mapping) libraries")
	fmt.Println("  • Dependency injection frameworks")
	fmt.Println("  • Testing and mocking frameworks")
	fmt.Println("  • Configuration parsers")
	fmt.Println("  • Generic algorithms (pre-Go 1.18)")
	fmt.Println("  • Debugging and inspection tools")

	fmt.Println("\n🎯 WHEN NOT TO USE REFLECT:")
	fmt.Println("  • Simple type conversions")
	fmt.Println("  • Performance-critical code")
	fmt.Println("  • When type is known at compile time")
	fmt.Println("  • When generics can solve the problem (Go 1.18+)")
	fmt.Println("  • For basic CRUD operations")
}

// تابع کمکی برای تکرار رشته
func stringsRepeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}
