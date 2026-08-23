// ============================================================================
// FILE: encoding_json_complete_guide.go
// TITLE: راهنمای کامل پکیج encoding/json در Go - Marshaling و Unmarshaling
// HOW TO RUN: go run encoding_json_complete_guide.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - JSON در Go
// ============================================================================
//
// JSON (JavaScript Object Notation) محبوب‌ترین فرمت تبادل داده در وب است.
// پکیج encoding/json امکانات کامل برای:
// 1. Marshaling: تبدیل داده‌های Go به JSON (encoding)
// 2. Unmarshaling: تبدیل JSON به داده‌های Go (decoding)
// 3. کار با استریم‌ها (Encoder/Decoder)
// 4. کنترل دقیق با تگ‌ها (tags)
// 5. فیلدهای اختیاری و سفارشی‌سازی
//
// قانون طلایی:
// "از تگ‌های json برای کنترل نام فیلدها، omitempty برای فیلدهای خالی،
//  و string برای اعداد درون رشته استفاده کن. هرگز از panic در Marshal/Unmarshal استفاده نکن."
// ============================================================================

package __internal_packages

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// ============================================================================
// بخش 1: ساختارهای پایه با تگ‌های JSON
// ============================================================================

// User یک مدل ساده با تگ‌های JSON
type User struct {
	ID        int       `json:"id"`              // نام فیلد در JSON: "id"
	Name      string    `json:"name"`            // نام فیلد در JSON: "name"
	Email     string    `json:"email,omitempty"` // اگر خالی باشد، در JSON نمی‌آید
	Age       int       `json:"age,string"`      // به عنوان رشته در JSON ذخیره می‌شود
	IsActive  bool      `json:"is_active"`       // تبدیل به boolean
	CreatedAt time.Time `json:"created_at"`      // زمان به فرمت RFC3339
	Password  string    `json:"-"`               // این فیلد نادیده گرفته می‌شود
	Role      string    `json:"role,omitempty"`  // اختیاری
}

// Product ساختار با فیلدهای تو در تو
type Product struct {
	ID       int                    `json:"id"`
	Name     string                 `json:"name"`
	Price    float64                `json:"price"`
	InStock  bool                   `json:"in_stock"`
	Tags     []string               `json:"tags,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// Order ساختار با nested struct
type Order struct {
	ID        string    `json:"order_id"`
	UserID    int       `json:"user_id"`
	Products  []Product `json:"products"`
	Total     float64   `json:"total"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	Shipping  *Address  `json:"shipping,omitempty"` // اشاره‌گر (می‌تواند nil باشد)
	Billing   *Address  `json:"billing,omitempty"`
}

// Address آدرس برای سفارش
type Address struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	Country string `json:"country"`
	ZipCode string `json:"zip_code"`
}

// ============================================================================
// بخش 2: Marshal - تبدیل Go به JSON
// ============================================================================

func demonstrateMarshal() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📤 MARSHAL - Converting Go to JSON")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 2.1 Marshal ساده
	// ============================================
	fmt.Println("\n--- 2.1 Basic Marshal ---")

	user := User{
		ID:        1,
		Name:      "Ali Rezaei",
		Email:     "ali@example.com",
		Age:       30,
		IsActive:  true,
		CreatedAt: time.Now(),
		Password:  "secret123", // این فیلد نادیده گرفته می‌شود
		Role:      "admin",
	}

	jsonData, err := json.Marshal(user)
	if err != nil {
		log.Printf("Marshal error: %v", err)
		return
	}

	fmt.Printf("  Original: %+v\n", user)
	fmt.Printf("  JSON: %s\n", jsonData)
	fmt.Printf("  Pretty JSON:\n%s\n", prettyPrint(jsonData))

	// ============================================
	// 2.2 Marshal با omitempty
	// ============================================
	fmt.Println("\n--- 2.2 Marshal with omitempty ---")

	user2 := User{
		ID:        2,
		Name:      "Sara",
		Email:     "", // خالی - در JSON نمی‌آید
		Age:       25,
		IsActive:  false, // مقدار صفر - در JSON می‌آید (false)
		CreatedAt: time.Now(),
		Role:      "", // خالی - در JSON نمی‌آید
	}

	jsonData2, _ := json.Marshal(user2)
	fmt.Printf("  JSON (empty fields omitted): %s\n", jsonData2)

	// ============================================
	// 2.3 Marshal با انواع مختلف
	// ============================================
	fmt.Println("\n--- 2.3 Marshal Different Types ---")

	// Slice
	users := []User{
		{ID: 1, Name: "Ali", Age: 30},
		{ID: 2, Name: "Sara", Age: 25},
	}
	jsonUsers, _ := json.Marshal(users)
	fmt.Printf("  Slice of users: %s\n", jsonUsers)

	// Map
	data := map[string]interface{}{
		"name":   "Test",
		"value":  42,
		"active": true,
		"tags":   []string{"a", "b", "c"},
	}
	jsonMap, _ := json.Marshal(data)
	fmt.Printf("  Map: %s\n", jsonMap)

	// Array
	arr := [3]int{1, 2, 3}
	jsonArr, _ := json.Marshal(arr)
	fmt.Printf("  Array: %s\n", jsonArr)

	// Primitive types
	jsonInt, _ := json.Marshal(42)
	jsonString, _ := json.Marshal("hello")
	jsonBool, _ := json.Marshal(true)
	fmt.Printf("  Int: %s, String: %s, Bool: %s\n", jsonInt, jsonString, jsonBool)

	// ============================================
	// 2.4 MarshalIndent - JSON زیبا
	// ============================================
	fmt.Println("\n--- 2.4 MarshalIndent (Pretty Print) ---")

	product := Product{
		ID:      1,
		Name:    "Laptop",
		Price:   999.99,
		InStock: true,
		Tags:    []string{"electronics", "computer"},
		Metadata: map[string]interface{}{
			"brand": "Dell",
			"color": "black",
		},
	}

	pretty, _ := json.MarshalIndent(product, "", "  ")
	fmt.Printf("  Pretty printed:\n%s\n", pretty)
}

// ============================================================================
// بخش 3: Unmarshal - تبدیل JSON به Go
// ============================================================================

func demonstrateUnmarshal() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📥 UNMARSHAL - Converting JSON to Go")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 3.1 Unmarshal ساده
	// ============================================
	fmt.Println("\n--- 3.1 Basic Unmarshal ---")

	jsonStr := `{"id":1,"name":"Ali Rezaei","email":"ali@example.com","age":"30","is_active":true,"created_at":"2024-01-15T10:30:00Z"}`

	var user User
	err := json.Unmarshal([]byte(jsonStr), &user)
	if err != nil {
		log.Printf("Unmarshal error: %v", err)
		return
	}

	fmt.Printf("  JSON: %s\n", jsonStr)
	fmt.Printf("  Parsed: %+v\n", user)

	// ============================================
	// 3.2 Unmarshal به انواع مختلف
	// ============================================
	fmt.Println("\n--- 3.2 Unmarshal Different Types ---")

	// به Slice
	jsonUsers := `[{"id":1,"name":"Ali"},{"id":2,"name":"Sara"}]`
	var users []User
	json.Unmarshal([]byte(jsonUsers), &users)
	fmt.Printf("  Slice: %+v\n", users)

	// به Map
	jsonMap := `{"name":"Product","price":99.99,"in_stock":true}`
	var result map[string]interface{}
	json.Unmarshal([]byte(jsonMap), &result)
	fmt.Printf("  Map: %+v\n", result)

	// دسترسی به مقادیر map
	if name, ok := result["name"]; ok {
		fmt.Printf("    name = %v (type: %T)\n", name, name)
	}
	if price, ok := result["price"]; ok {
		fmt.Printf("    price = %v (type: %T)\n", price, price)
	}

	// به Array
	jsonArr := `[1,2,3,4,5]`
	var numbers [5]int
	json.Unmarshal([]byte(jsonArr), &numbers)
	fmt.Printf("  Array: %v\n", numbers)

	// ============================================
	// 3.3 Unmarshal با داده‌های نامعتبر
	// ============================================
	fmt.Println("\n--- 3.3 Error Handling ---")

	invalidJSON := `{"id":1,"name":"Ali"` // ناقص

	var u User
	err = json.Unmarshal([]byte(invalidJSON), &u)
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	}
}

// ============================================================================
// بخش 4: فیلدهای اختیاری و سفارشی (Custom Marshaling)
// ============================================================================

// CustomTime زمان سفارشی با فرمت دلخواه
type CustomTime struct {
	time.Time
}

// فرمت سفارشی برای JSON
const customTimeFormat = "2006-01-02"

// MarshalJSON پیاده‌سازی سفارشی Marshal
func (ct CustomTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + ct.Format(customTimeFormat) + `"`), nil
}

// UnmarshalJSON پیاده‌سازی سفارشی Unmarshal
func (ct *CustomTime) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), `"`)
	parsed, err := time.Parse(customTimeFormat, str)
	if err != nil {
		return err
	}
	ct.Time = parsed
	return nil
}

// Event رویداد با فیلدهای سفارشی
type Event struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	Date      CustomTime `json:"date"`
	Timestamp time.Time  `json:"timestamp"`
}

// ProductWithCustom ساختار با فیلدهای نادیده گرفته شده
type ProductWithCustom struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price,omitempty"`
	Internal string  `json:"-"`                // همیشه نادیده گرفته می‌شود
	Secret   string  `json:"secret,omitempty"` // اگر خالی باشد نادیده گرفته می‌شود
}

func demonstrateCustomMarshaling() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🔧 CUSTOM MARSHALING - Custom JSON Handling")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 4.1 Custom Time Marshal/Unmarshal
	// ============================================
	fmt.Println("\n--- 4.1 Custom Time Format ---")

	event := Event{
		ID:        1,
		Name:      "Conference",
		Date:      CustomTime{Time: time.Now()},
		Timestamp: time.Now(),
	}

	jsonData, _ := json.Marshal(event)
	fmt.Printf("  Marshaled: %s\n", jsonData)

	var event2 Event
	json.Unmarshal(jsonData, &event2)
	fmt.Printf("  Unmarshaled: %+v\n", event2)

	// ============================================
	// 4.2 Omitted Fields
	// ============================================
	fmt.Println("\n--- 4.2 Omitted Fields ---")

	prod := ProductWithCustom{
		ID:       1,
		Name:     "Laptop",
		Price:    0,        // omitempty - حذف می‌شود
		Internal: "secret", // - حذف می‌شود
		Secret:   "",       // omitempty - حذف می‌شود
	}

	jsonProd, _ := json.Marshal(prod)
	fmt.Printf("  Fields omitted: %s\n", jsonProd)
}

// ============================================================================
// بخش 5: کار با Encoder و Decoder (استریم‌ها)
// ============================================================================

func demonstrateEncoderDecoder() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📡 ENCODER/DECODER - Streaming JSON")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 5.1 Encoder - نوشتن مستقیم به Writer
	// ============================================
	fmt.Println("\n--- 5.1 Encoder (Stream Writing) ---")

	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetIndent("", "  ") // pretty print

	products := []Product{
		{ID: 1, Name: "Laptop", Price: 999.99, InStock: true},
		{ID: 2, Name: "Mouse", Price: 19.99, InStock: true},
		{ID: 3, Name: "Keyboard", Price: 49.99, InStock: false},
	}

	for _, p := range products {
		if err := encoder.Encode(p); err != nil {
			log.Printf("Encode error: %v", err)
		}
	}

	fmt.Printf("  Encoded stream:\n%s", buf.String())

	// ============================================
	// 5.2 Decoder - خواندن مستقیم از Reader
	// ============================================
	fmt.Println("\n--- 5.2 Decoder (Stream Reading) ---")

	jsonStream := `{"id":1,"name":"Laptop","price":999.99}
{"id":2,"name":"Mouse","price":19.99}
{"id":3,"name":"Keyboard","price":49.99}`

	decoder := json.NewDecoder(strings.NewReader(jsonStream))

	var products2 []Product
	for decoder.More() {
		var p Product
		if err := decoder.Decode(&p); err != nil {
			log.Printf("Decode error: %v", err)
			break
		}
		products2 = append(products2, p)
	}

	fmt.Printf("  Decoded products: %+v\n", products2)

	// ============================================
	// 5.3 Decoder با JSON بزرگ (تکه تکه)
	// ============================================
	fmt.Println("\n--- 5.3 Streaming Large JSON ---")

	largeJSON := `[
		{"id":1,"name":"Item1"},
		{"id":2,"name":"Item2"},
		{"id":3,"name":"Item3"}
	]`

	decoder = json.NewDecoder(strings.NewReader(largeJSON))

	// خواندن token اول (آغاز آرایه)
	token, err := decoder.Token()
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	fmt.Printf("  First token: %v\n", token)

	// خواندن هر عنصر به صورت جداگانه
	for decoder.More() {
		var item struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}
		if err := decoder.Decode(&item); err != nil {
			break
		}
		fmt.Printf("  Item: %+v\n", item)
	}

	// خواندن token آخر (پایان آرایه)
	token, _ = decoder.Token()
	fmt.Printf("  Last token: %v\n", token)
}

// ============================================================================
// بخش 6: RawMessage - پردازش دیرهنگام JSON
// ============================================================================

type FlexibleEvent struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"` // ذخیره خام برای پردازش بعدی
}

type UserPayload struct {
	UserID int    `json:"user_id"`
	Name   string `json:"name"`
}

type OrderPayload struct {
	OrderID string  `json:"order_id"`
	Amount  float64 `json:"amount"`
}

func demonstrateRawMessage() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📦 RAW MESSAGE - Delayed JSON Processing")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 6.1 استفاده از RawMessage
	// ============================================
	fmt.Println("\n--- 6.1 Using RawMessage ---")

	events := []string{
		`{"type":"user","payload":{"user_id":1,"name":"Ali"}}`,
		`{"type":"order","payload":{"order_id":"ORD-001","amount":99.99}}`,
	}

	for _, eventStr := range events {
		var event FlexibleEvent
		json.Unmarshal([]byte(eventStr), &event)

		fmt.Printf("  Event type: %s\n", event.Type)
		fmt.Printf("  Raw payload: %s\n", event.Payload)

		switch event.Type {
		case "user":
			var user UserPayload
			json.Unmarshal(event.Payload, &user)
			fmt.Printf("    User: %+v\n", user)
		case "order":
			var order OrderPayload
			json.Unmarshal(event.Payload, &order)
			fmt.Printf("    Order: %+v\n", order)
		}
	}
}

// ============================================================================
// بخش 7: MarshalJSON و UnmarshalJSON (اینترفیس‌های سفارشی)
// ============================================================================

type Status string

const (
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
	StatusPending  Status = "pending"
)

// MarshalJSON برای Status
func (s Status) MarshalJSON() ([]byte, error) {
	return []byte(`"` + string(s) + `"`), nil
}

// UnmarshalJSON برای Status
func (s *Status) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), `"`)
	switch str {
	case "active":
		*s = StatusActive
	case "inactive":
		*s = StatusInactive
	case "pending":
		*s = StatusPending
	default:
		return fmt.Errorf("invalid status: %s", str)
	}
	return nil
}

type Task struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status Status `json:"status"`
}

func demonstrateMarshalJSONInterface() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🔌 MarshalJSON/UnmarshalJSON Interfaces")
	fmt.Println(stringsRepeat("=", 80))

	task := Task{
		ID:     1,
		Title:  "Complete project",
		Status: StatusActive,
	}

	jsonData, _ := json.Marshal(task)
	fmt.Printf("  Marshaled: %s\n", jsonData)

	var task2 Task
	json.Unmarshal(jsonData, &task2)
	fmt.Printf("  Unmarshaled: %+v\n", task2)

	// تست با مقدار نامعتبر
	invalid := `{"id":2,"title":"Test","status":"invalid"}`
	var task3 Task
	err := json.Unmarshal([]byte(invalid), &task3)
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	}
}

// ============================================================================
// بخش 8: کار با فایل‌های JSON
// ============================================================================

func demonstrateFileJSON() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📁 WORKING WITH JSON FILES")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 8.1 نوشتن JSON در فایل
	// ============================================
	fmt.Println("\n--- 8.1 Writing JSON to File ---")

	users := []User{
		{ID: 1, Name: "Ali", Email: "ali@test.com", Age: 30, IsActive: true},
		{ID: 2, Name: "Sara", Email: "sara@test.com", Age: 25, IsActive: true},
	}

	// روش 1: Marshal و WriteFile
	data, _ := json.MarshalIndent(users, "", "  ")
	err := os.WriteFile("users.json", data, 0644)
	if err == nil {
		fmt.Println("  Saved users to users.json")
		defer os.Remove("users.json")
	}

	// روش 2: Encoder
	file, _ := os.Create("products.json")
	defer file.Close()
	defer os.Remove("products.json")

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	products := []Product{
		{ID: 1, Name: "Laptop", Price: 999.99, InStock: true},
		{ID: 2, Name: "Mouse", Price: 19.99, InStock: true},
	}
	encoder.Encode(products)
	fmt.Println("  Saved products to products.json")

	// ============================================
	// 8.2 خواندن JSON از فایل
	// ============================================
	fmt.Println("\n--- 8.2 Reading JSON from File ---")

	// روش 1: ReadFile و Unmarshal
	data2, _ := os.ReadFile("users.json")
	var loadedUsers []User
	json.Unmarshal(data2, &loadedUsers)
	fmt.Printf("  Loaded users: %+v\n", loadedUsers)

	// روش 2: Decoder
	file2, _ := os.Open("products.json")
	defer file2.Close()

	var loadedProducts []Product
	decoder := json.NewDecoder(file2)
	decoder.Decode(&loadedProducts)
	fmt.Printf("  Loaded products: %+v\n", loadedProducts)
}

// ============================================================================
// بخش 9: JSON با انواع داینامیک (map[string]interface{})
// ============================================================================

func demonstrateDynamicJSON() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🎭 DYNAMIC JSON - Working with map[string]interface{}")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 9.1 Unmarshal به map
	// ============================================
	fmt.Println("\n--- 9.1 Unmarshal to Map ---")

	jsonStr := `{
		"name": "Product X",
		"price": 49.99,
		"in_stock": true,
		"tags": ["new", "sale"],
		"metadata": {
			"weight": 1.5,
			"color": "black"
		}
	}`

	var data map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Println("  Parsed dynamic data:")
	for k, v := range data {
		fmt.Printf("    %s: %v (type: %T)\n", k, v, v)
	}

	// ============================================
	// 9.2 دسترسی به nested data
	// ============================================
	fmt.Println("\n--- 9.2 Accessing Nested Data ---")

	if metadata, ok := data["metadata"].(map[string]interface{}); ok {
		if weight, ok := metadata["weight"].(float64); ok {
			fmt.Printf("  Weight: %.2f\n", weight)
		}
	}

	if tags, ok := data["tags"].([]interface{}); ok {
		fmt.Printf("  Tags: %v\n", tags)
	}

	// ============================================
	// 9.3 Marshal از map
	// ============================================
	fmt.Println("\n--- 9.3 Marshal from Map ---")

	newData := map[string]interface{}{
		"id":     100,
		"name":   "New Item",
		"price":  29.99,
		"active": true,
		"tags":   []string{"featured", "hot"},
		"metadata": map[string]interface{}{
			"created": "2024-01-15",
		},
	}

	jsonOutput, _ := json.MarshalIndent(newData, "", "  ")
	fmt.Printf("  Marshaled map:\n%s\n", jsonOutput)
}

// ============================================================================
// بخش 10: Validation و فیلتر کردن فیلدها
// ============================================================================

type ValidatedUser struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
	XX    int    `bson:"age"`
}

// Validate اعتبارسنجی بعد از Unmarshal
func (u *ValidatedUser) Validate() error {
	if u.Name == "" {
		return fmt.Errorf("name is required")
	}
	if u.Email == "" {
		return fmt.Errorf("email is required")
	}
	if u.Age < 0 || u.Age > 150 {
		return fmt.Errorf("age must be between 0 and 150")
	}
	return nil
}

// Custom unmarshal with validation
func (u *ValidatedUser) UnmarshalJSON(data []byte) error {
	type Alias ValidatedUser
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(u),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	return u.Validate()
}

func demonstrateValidation() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("✅ VALIDATION - Post-Unmarshal Validation")
	fmt.Println(stringsRepeat("=", 80))

	validJSON := `{"id":1,"name":"Ali","email":"ali@test.com","age":30}`
	invalidJSON := `{"id":2,"name":"","email":"","age":200}`

	var user1 ValidatedUser
	if err := json.Unmarshal([]byte(validJSON), &user1); err != nil {
		fmt.Printf("  Error: %v\n", err)
	} else {
		fmt.Printf("  Valid user: %+v\n", user1)
	}

	var user2 ValidatedUser
	if err := json.Unmarshal([]byte(invalidJSON), &user2); err != nil {
		fmt.Printf("  Invalid user error: %v\n", err)
	}
}

// ============================================================================
// بخش 11: اشتباهات رایج
// ============================================================================

func demonstrateCommonMistakes() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("❌ COMMON MISTAKES WITH encoding/json")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println("\n❌ Mistake 1: Passing non-pointer to Unmarshal")
	fmt.Println("   var user User")
	fmt.Println("   json.Unmarshal(data, user)  // ❌ user is not pointer")
	fmt.Println("   ✅ json.Unmarshal(data, &user)")

	fmt.Println("\n❌ Mistake 2: Ignoring errors")
	fmt.Println("   json.Unmarshal(data, &user)  // ignoring error")
	fmt.Println("   ✅ Always check returned error")

	fmt.Println("\n❌ Mistake 3: Capitalization mismatch")
	fmt.Println("   JSON: {\"user_id\":1}  // lowercase")
	fmt.Println("   Go:   UserID int      // uppercase")
	fmt.Println("   ✅ Use tags: `json:\"user_id\"`")

	fmt.Println("\n❌ Mistake 4: Assuming map order")
	fmt.Println("   Maps in JSON have no order")
	fmt.Println("   ✅ Don't rely on iteration order")

	fmt.Println("\n❌ Mistake 5: Not handling null values")
	fmt.Println("   JSON null becomes nil pointer in Go")
	fmt.Println("   ✅ Check for nil before dereferencing")

	fmt.Println("\n❌ Mistake 6: Using int for large numbers")
	fmt.Println("   JSON numbers can be large")
	fmt.Println("   ✅ Use int64 or json.Number")

	fmt.Println("\n❌ Mistake 7: Timezone issues")
	fmt.Println("   time.Time uses UTC by default")
	fmt.Println("   ✅ Use custom marshaling for specific timezone")
}

// ============================================================================
// بخش 12: جمع‌بندی و جدول مرجع
// ============================================================================

func prettyPrint(data []byte) string {
	var pretty bytes.Buffer
	json.Indent(&pretty, data, "  ", "  ")
	return pretty.String()
}

func main() {
	fmt.Println(stringsRepeat("=", 80))
	fmt.Println("🎯 COMPLETE encoding/json GUIDE IN GO")
	fmt.Println(stringsRepeat("=", 80))

	// بخش 1: Marshal
	demonstrateMarshal()

	// بخش 2: Unmarshal
	demonstrateUnmarshal()

	// بخش 3: Custom Marshaling
	demonstrateCustomMarshaling()

	// بخش 4: Encoder/Decoder
	demonstrateEncoderDecoder()

	// بخش 5: RawMessage
	demonstrateRawMessage()

	// بخش 6: MarshalJSON/UnmarshalJSON
	demonstrateMarshalJSONInterface()

	// بخش 7: File operations
	demonstrateFileJSON()

	// بخش 8: Dynamic JSON
	demonstrateDynamicJSON()

	// بخش 9: Validation
	demonstrateValidation()

	// بخش 10: Common mistakes
	demonstrateCommonMistakes()

	// بخش 11: Quick reference
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📚 encoding/json QUICK REFERENCE")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ BASIC OPERATIONS                                              │")
	fmt.Println("├─────────────────────────────────────────────────────────────────┤")
	fmt.Println("│ json.Marshal(v)           - Convert Go to JSON                 │")
	fmt.Println("│ json.MarshalIndent(v, prefix, indent) - Pretty JSON           │")
	fmt.Println("│ json.Unmarshal(data, &v)  - Convert JSON to Go                 │")
	fmt.Println("│ json.NewEncoder(w).Encode(v) - Stream encode                   │")
	fmt.Println("│ json.NewDecoder(r).Decode(&v) - Stream decode                  │")
	fmt.Println("└─────────────────────────────────────────────────────────────────┘")

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ JSON TAGS                                                     │")
	fmt.Println("├─────────────────────────────────────────────────────────────────┤")
	fmt.Println("│ `json:\"field_name\"`        - Rename field                     │")
	fmt.Println("│ `json:\"field,omitempty\"`   - Omit if empty                    │")
	fmt.Println("│ `json:\"-\"`                - Ignore field                      │")
	fmt.Println("│ `json:\"field,string\"`      - Force string conversion          │")
	fmt.Println("│ `json:\",omitempty\"`        - Use original name, omit if empty │")
	fmt.Println("└─────────────────────────────────────────────────────────────────┘")

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ CUSTOM INTERFACES                                             │")
	fmt.Println("├─────────────────────────────────────────────────────────────────┤")
	fmt.Println("│ MarshalJSON() ([]byte, error)   - Custom encoding              │")
	fmt.Println("│ UnmarshalJSON(data []byte) error - Custom decoding            │")
	fmt.Println("│ json.RawMessage                - Raw JSON storage             │")
	fmt.Println("│ json.Number                    - Preserve number precision    │")
	fmt.Println("└─────────────────────────────────────────────────────────────────┘")

	fmt.Println("\n💡 GOLDEN RULES:")
	fmt.Println("  1. Always pass pointer to Unmarshal")
	fmt.Println("  2. Always check errors from Marshal/Unmarshal")
	fmt.Println("  3. Use tags to match JSON field names")
	fmt.Println("  4. Use omitempty to skip empty fields")
	fmt.Println("  5. Use RawMessage for delayed parsing")
	fmt.Println("  6. Use map[string]interface{} for dynamic JSON")
	fmt.Println("  7. Use Encoder/Decoder for streams and files")
	fmt.Println("  8. Validate data after Unmarshal")
	fmt.Println("  9. Use int64 for large numbers")
	fmt.Println("  10. Never panic in MarshalJSON/UnmarshalJSON")

	fmt.Println("\n🎯 TYPES AND THEIR JSON REPRESENTATION:")
	fmt.Println("  bool      → true/false")
	fmt.Println("  int/float → number")
	fmt.Println("  string    → \"string\"")
	fmt.Println("  []byte    → base64 string")
	fmt.Println("  struct    → object")
	fmt.Println("  slice     → array")
	fmt.Println("  map       → object")
	fmt.Println("  pointer   → value or null")
	fmt.Println("  interface → any JSON type")
	fmt.Println("  time.Time → RFC3339 string")
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
go run encoding_json_complete_guide.go

# فایل‌های ایجاد شده را می‌توانید ببینید
ls -la users.json products.json
*/
