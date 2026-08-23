// ============================================================================
// FILE: bson_complete_guide.go
// TITLE: راهنمای کامل پکیج bson در Go - کار با داده‌های باینری MongoDB
// HOW TO RUN: go run bson_complete_guide.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - BSON چیست و چه تفاوتی با JSON دارد؟
// ============================================================================
//
// BSON (Binary JSON) یک فرمت سریال‌سازی باینری است که توسط MongoDB استفاده می‌شود.
// برخلاف JSON که متنی است، BSON باینری است و قابلیت‌های بیشتری دارد.
//
// تفاوت‌های کلیدی BSON با JSON:
// 1. BSON انواع داده بیشتری دارد (ObjectID, Date, Binary, etc.)
// 2. BSON حجم کمتری برای اعداد بزرگ دارد (fixed 8 bytes)
// 3. BSON سرعت parsing بالاتری دارد (به دلیل فرمت باینری)
// 4. BSON ترتیب فیلدها را حفظ می‌کند (با استفاده از نوع D)
// 5. BSON جستجو در داده‌ها را بدون unmarshal کامل ممکن می‌کند (bson.Raw)
//
// قانون طلایی:
// "برای کار با MongoDB همیشه از bson استفاده کن. برای ذخیره‌سازی، نوع D را برای
//  حفظ ترتیب، نوع M را برای دسترسی آسان، و structها را برای مدل‌های مشخص استفاده کن."
// ============================================================================

package __internal_packages

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ============================================================================
// بخش 1: انواع پایه BSON (D, M, A, E)
// ============================================================================

func demonstrateBasicTypes() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📦 BASIC BSON TYPES - D, M, A, E")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 1.1 نوع D (Document - Ordered Slice)
	// ============================================
	fmt.Println("\n--- 1.1 Type D (Ordered Document) ---")

	// D یک اسلایس از Eهاست که ترتیب را حفظ می‌کند
	// برای دستورات MongoDB که به ترتیب اهمیت دارد مناسب است
	doc := bson.D{
		{Key: "name", Value: "Ali"},
		{Key: "age", Value: 30},
		{Key: "city", Value: "Tehran"},
	}
	fmt.Printf("  D document: %+v\n", doc)
	fmt.Printf("  First field: %s = %v\n", doc[0].Key, doc[0].Value)

	// استفاده در فیلتر کوئری
	filter := bson.D{
		{"age", bson.D{{"$gt", 18}}},
		{"city", "Tehran"},
	}
	fmt.Printf("  Query filter: %+v\n", filter)

	// ============================================
	// 1.2 نوع M (Map - Unordered Document)
	// ============================================
	fmt.Println("\n--- 1.2 Type M (Unordered Map) ---")

	// M یک map است، ترتیب不重要
	// برای خوانایی و دسترسی آسان مناسب است
	docMap := bson.M{
		"name": "Sara",
		"age":  25,
		"city": "Isfahan",
	}
	fmt.Printf("  M document: %+v\n", docMap)
	fmt.Printf("  Access by key: name = %v\n", docMap["name"])

	// استفاده در آپدیت
	update := bson.M{
		"$set": bson.M{
			"age":  26,
			"city": "Shiraz",
		},
	}
	fmt.Printf("  Update document: %+v\n", update)

	// ============================================
	// 1.3 نوع A (Array)
	// ============================================
	fmt.Println("\n--- 1.3 Type A (Array) ---")

	// A برای نمایش آرایه‌های BSON
	array := bson.A{"apple", "banana", "cherry"}
	fmt.Printf("  A array: %+v\n", array)

	// آرایه در کوئری
	filterArray := bson.M{
		"tags": bson.M{"$in": bson.A{"fruit", "organic"}},
	}
	fmt.Printf("  Array in query: %+v\n", filterArray)

	// ============================================
	// 1.4 نوع E (Element)
	// ============================================
	fmt.Println("\n--- 1.4 Type E (Element) ---")

	// E یک عنصر تکی در D است
	elements := []bson.E{
		{Key: "name", Value: "Reza"},
		{Key: "age", Value: 35},
	}
	fmt.Printf("  E elements: %+v\n", elements)
}

// ============================================================================
// بخش 2: Marshal و Unmarshal - تبدیل به BSON
// ============================================================================

// Person مدل نمونه برای Marshal/Unmarshal
type Person struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	Age       int                `bson:"age"`
	Email     string             `bson:"email,omitempty"`
	CreatedAt time.Time          `bson:"created_at"`
	Tags      []string           `bson:"tags,omitempty"`
	Address   Address            `bson:"address"`
}

type Address struct {
	Street  string `bson:"street"`
	City    string `bson:"city"`
	Country string `bson:"country"`
}

func demonstrateMarshalUnmarshal() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🔄 MARSHAL & UNMARSHAL - Go to BSON and Back")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 2.1 Marshal - تبدیل Go به BSON
	// ============================================
	fmt.Println("\n--- 2.1 Marshal (Go to BSON) ---")

	person := Person{
		Name:      "Ali Rezaei",
		Age:       30,
		Email:     "ali@example.com",
		CreatedAt: time.Now(),
		Tags:      []string{"developer", "golang"},
		Address: Address{
			Street:  "Valiasr St",
			City:    "Tehran",
			Country: "Iran",
		},
	}

	// تبدیل به BSON bytes
	bsonData, err := bson.Marshal(person)
	if err != nil {
		log.Printf("  Marshal error: %v", err)
		return
	}
	fmt.Printf("  Marshaled BSON length: %d bytes\n", len(bsonData))
	fmt.Printf("  BSON bytes (hex): %x...\n", bsonData[:20])

	// ============================================
	// 2.2 Unmarshal - تبدیل BSON به Go
	// ============================================
	fmt.Println("\n--- 2.2 Unmarshal (BSON to Go) ---")

	var decoded Person
	err = bson.Unmarshal(bsonData, &decoded)
	if err != nil {
		log.Printf("  Unmarshal error: %v", err)
		return
	}
	fmt.Printf("  Decoded person: %+v\n", decoded)
	fmt.Printf("  Name: %s, Age: %d, City: %s\n",
		decoded.Name, decoded.Age, decoded.Address.City)

	// ============================================
	// 2.3 Marshal به D و M
	// ============================================
	fmt.Println("\n--- 2.3 Marshal to D and M ---")

	// تبدیل مستقیم به D
	doc := bson.D{
		{"name", "Test"},
		{"value", 42},
	}
	bsonDoc, _ := bson.Marshal(doc)
	fmt.Printf("  Marshaled D: %x...\n", bsonDoc[:20])

	// تبدیل مستقیم به M
	docMap := bson.M{
		"name":  "Test",
		"value": 42,
	}
	bsonMap, _ := bson.Marshal(docMap)
	fmt.Printf("  Marshaled M: %x...\n", bsonMap[:20])
}

// ============================================================================
// بخش 3: تگ‌های BSON (Struct Tags)
// ============================================================================

type Product struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`        // نام فیلد در BSON
	Name        string             `bson:"name"`                 // نام معمولی
	Price       float64            `bson:"price,omitempty"`      // اگر صفر باشد حذف می‌شود
	Quantity    int                `bson:"quantity,minsize"`     // اگر در int32 جا شد، به عنوان int32 ذخیره کن
	Description string             `bson:"description,truncate"` // اعشار را حذف کن (برای اعداد)
	Metadata    map[string]string  `bson:"metadata,inline"`      // فیلدها را flat کن
	Internal    string             `bson:"-"`                    // نادیده گرفته می‌شود
	Secret      string             `bson:"secret,omitempty"`     // اختیاری
}

func demonstrateStructTags() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🏷️ STRUCT TAGS - Customizing BSON Mapping")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 3.1 استفاده از تگ‌ها
	// ============================================
	fmt.Println("\n--- 3.1 Using Struct Tags ---")

	product := Product{
		Name:     "Laptop",
		Price:    999.99,
		Quantity: 10,
		Metadata: map[string]string{
			"brand": "Dell",
			"color": "black",
		},
		Internal: "secret value", // این فیلد نادیده گرفته می‌شود
	}

	bsonData, _ := bson.Marshal(product)
	fmt.Printf("  Marshaled with tags: %x...\n", bsonData[:20])

	// نمایش فیلدهایی که در BSON هستند
	var raw bson.M
	bson.Unmarshal(bsonData, &raw)
	fmt.Printf("  Fields in BSON: %v\n", getMapKeys(raw))

	// ============================================
	// 3.2 توضیح تگ‌ها
	// ============================================
	fmt.Println("\n--- 3.2 Tag Explanations ---")
	fmt.Println("  • omitempty  - اگر مقدار صفر باشد، فیلد حذف می‌شود")
	fmt.Println("  • minsize    - اعداد را در کوچکترین نوع ممکن ذخیره کن")
	fmt.Println("  • truncate   - اعشار اعداد را حذف کن (برای انواع غیر-float)")
	fmt.Println("  • inline     - فیلدهای struct/map را flat کن")
	fmt.Println("  • -          - فیلد را نادیده بگیر (مثل json:\"-\")")
}

// ============================================================================
// بخش 4: انواع داده‌های ویژه BSON (primitive)
// ============================================================================

func demonstratePrimitiveTypes() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🔷 PRIMITIVE TYPES - ObjectID, DateTime, Binary, etc.")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 4.1 ObjectID
	// ============================================
	fmt.Println("\n--- 4.1 ObjectID ---")

	// ایجاد ObjectID جدید
	oid := primitive.NewObjectID()
	fmt.Printf("  New ObjectID: %s\n", oid.Hex())
	fmt.Printf("  Timestamp: %v\n", oid.Timestamp())

	// تبدیل از hex
	oidFromHex, err := primitive.ObjectIDFromHex("507f1f77bcf86cd799439011")
	if err == nil {
		fmt.Printf("  From hex: %s\n", oidFromHex.Hex())
	}

	// ObjectID خالی
	var emptyOid primitive.ObjectID
	fmt.Printf("  Empty ObjectID: %s (is zero: %v)\n", emptyOid.Hex(), emptyOid.IsZero())

	// ============================================
	// 4.2 DateTime
	// ============================================
	fmt.Println("\n--- 4.2 DateTime ---")

	// DateTime در BSON میلی‌ثانیه از epoch است
	now := time.Now()
	dateTime := primitive.NewDateTimeFromTime(now)
	fmt.Printf("  Current time: %v\n", now)
	fmt.Printf("  BSON DateTime: %d\n", dateTime)

	// تبدیل به time.Time
	backToTime := dateTime.Time()
	fmt.Printf("  Back to time.Time: %v\n", backToTime)

	// ============================================
	// 4.3 Binary Data
	// ============================================
	fmt.Println("\n--- 4.3 Binary Data ---")

	binaryData := primitive.Binary{
		Subtype: 0x00, // Generic binary subtype
		Data:    []byte("hello world"),
	}
	fmt.Printf("  Binary data: %x (len=%d)\n", binaryData.Data, len(binaryData.Data))

	// ============================================
	// 4.4 Decimal128 (برای اعداد دقیق)
	// ============================================
	fmt.Println("\n--- 4.4 Decimal128 ---")

	// برای اعداد مالی که به دقت بالا نیاز دارند
	// decimal, _ := primitive.ParseDecimal128("123.456")
	// fmt.Printf("  Decimal128: %s\n", decimal.String())

	// ============================================
	// 4.5 Timestamp (برای replication)
	// ============================================
	fmt.Println("\n--- 4.5 Timestamp ---")

	timestamp := primitive.Timestamp{T: 1234567890, I: 1}
	fmt.Printf("  Timestamp: T=%d, I=%d\n", timestamp.T, timestamp.I)

	// ============================================
	// 4.6 سایر انواع
	// ============================================
	fmt.Println("\n--- 4.6 Other Primitive Types ---")
	fmt.Println("  • primitive.Regex    - Regular expression")
	fmt.Println("  • primitive.JavaScript - JavaScript code")
	fmt.Println("  • primitive.Symbol   - Symbol")
	fmt.Println("  • primitive.MinKey   - Minimum possible key")
	fmt.Println("  • primitive.MaxKey   - Maximum possible key")
	fmt.Println("  • primitive.Undefined - Undefined value")
	fmt.Println("  • primitive.Null     - Null value")
}

// ============================================================================
// بخش 5: BSON Raw - کار با داده خام BSON
// ============================================================================

func demonstrateRawBSON() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📦 RAW BSON - Working with Raw BSON Data")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 5.1 ایجاد و اعتبارسنجی Raw BSON
	// ============================================
	fmt.Println("\n--- 5.1 Creating and Validating Raw BSON ---")

	// ایجاد یک سند BSON
	doc := bson.D{
		{"name", "Ali"},
		{"age", 30},
		{"city", "Tehran"},
	}

	bsonBytes, _ := bson.Marshal(doc)
	raw := bson.Raw(bsonBytes)

	// اعتبارسنجی
	err := raw.Validate()
	if err == nil {
		fmt.Println("  Raw BSON is valid")
	}

	// ============================================
	// 5.2 جستجو در Raw BSON
	// ============================================
	fmt.Println("\n--- 5.2 Looking Up Values ---")

	// جستجوی فیلد
	value := raw.Lookup("name")
	fmt.Printf("  Lookup 'name': %s\n", value.String())

	// جستجو با type assertion
	if name, ok := raw.Lookup("name").StringValueOK(); ok {
		fmt.Printf("  Name value: %s\n", name)
	}

	// جستجوی فیلد تو در تو (در صورت وجود)
	// nested := raw.Lookup("address", "city")

	// ============================================
	// 5.3 عناصر Raw BSON
	// ============================================
	fmt.Println("\n--- 5.3 Raw Elements ---")

	elements, _ := raw.Elements()
	for _, elem := range elements {
		key := elem.Key()
		value := elem.Value()
		fmt.Printf("  Element: %s = %v\n", key, value)
	}
}

// ============================================================================
// بخش 6: BSON Options - تنظیمات پیشرفته
// ============================================================================

func demonstrateBSONOptions() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("⚙️ BSON OPTIONS - Advanced Configuration")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println("\n--- BSON Options Overview ---")
	fmt.Println(`
تنظیمات قابل اعمال در ClientOptions:

1. UseJSONStructTags
   - اگر تگ bson وجود نداشت، از تگ json استفاده کن
   - مفید برای کدهایی که قبلاً با json کار می‌کردند

2. NilSliceAsEmpty
   - nil sliceها را به عنوان آرایه خالی marshal کن (به جای null)
   - برای consistency با سایر سیستم‌ها

3. OmitEmpty
   - به صورت سراسری فیلدهای خالی را حذف کن
   - معادل omitempty روی همه فیلدها

4. DiscardUnknown
   - فیلدهای ناشناخته را نادیده بگیر (به جای خطا)

5. JSONFallbackStructTagParser
   - از تگ json به عنوان fallback استفاده کن
`)

	// مثال استفاده (در کد واقعی)
	// bsonOpts := &options.BSONOptions{
	//     UseJSONStructTags: true,
	//     NilSliceAsEmpty:   true,
	//     OmitEmpty:         true,
	// }
	// clientOpts := options.Client().ApplyURI(uri).SetBSONOptions(bsonOpts)
}

// ============================================================================
// بخش 7: BSON در مقابل JSON - مقایسه عملی
// ============================================================================

type CompareStruct struct {
	ID       int     `bson:"id" json:"id"`
	Name     string  `bson:"name" json:"name"`
	Price    float64 `bson:"price" json:"price"`
	InStock  bool    `bson:"in_stock" json:"in_stock"`
	LongNum  int64   `bson:"long_num" json:"long_num"`
	LargeNum int64   `bson:"large_num" json:"large_num"`
}

func demonstrateBSONvsJSON() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("⚖️ BSON vs JSON - Comparison")
	fmt.Println(stringsRepeat("=", 80))

	testData := CompareStruct{
		ID:       123,
		Name:     "Test Product",
		Price:    99.99,
		InStock:  true,
		LongNum:  42,
		LargeNum: 9223372036854775807, // max int64
	}

	// Marshal به BSON
	bsonData, _ := bson.Marshal(testData)

	// Marshal به JSON
	jsonData, _ := json.Marshal(testData)

	fmt.Printf("\n--- Size Comparison ---\n")
	fmt.Printf("  BSON size: %d bytes\n", len(bsonData))
	fmt.Printf("  JSON size: %d bytes\n", len(jsonData))

	fmt.Printf("\n--- BSON vs JSON Characteristics ---\n")
	fmt.Println("  +------------------+------------------+------------------+")
	fmt.Println("  | Feature          | BSON             | JSON             |")
	fmt.Println("  +------------------+------------------+------------------+")
	fmt.Println("  | Format           | Binary           | Text             |")
	fmt.Println("  | Human readable   | ❌ No            | ✅ Yes           |")
	fmt.Println("  | Data types       | ✅ Rich          | ❌ Limited       |")
	fmt.Println("  | ObjectID support | ✅ Native        | ❌ String only   |")
	fmt.Println("  | Date support     | ✅ Native        | ❌ String only   |")
	fmt.Println("  | Traversal speed  | ✅ Fast          | ⚠️ Medium        |")
	fmt.Println("  | Parsing speed    | ✅ Fast          | ⚠️ Medium        |")
	fmt.Println("  | Binary data      | ✅ Native        | ❌ Base64        |")
	fmt.Println("  | Order preservation| ✅ D type        | ❌ Not preserved|")
	fmt.Println("  +------------------+------------------+------------------+")
}

// ============================================================================
// بخش 8: کاربردهای عملی در MongoDB
// ============================================================================

func demonstratePracticalUseCases() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("💡 PRACTICAL USE CASES - MongoDB Operations")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 8.1 فیلترهای کوئری
	// ============================================
	fmt.Println("\n--- 8.1 Query Filters ---")

	// فیلتر ساده
	filter1 := bson.M{"age": bson.M{"$gt": 18}}
	fmt.Printf("  Age > 18: %+v\n", filter1)

	// فیلتر ترکیبی
	filter2 := bson.M{
		"$and": []bson.M{
			{"age": bson.M{"$gte": 18}},
			{"age": bson.M{"$lte": 65}},
			{"city": "Tehran"},
		},
	}
	fmt.Printf("  Complex filter: %+v\n", filter2)

	// فیلتر با ObjectID
	oid := primitive.NewObjectID()
	filter3 := bson.M{"_id": oid}
	fmt.Printf("  By ObjectID: %+v\n", filter3)

	// ============================================
	// 8.2 آپدیت‌ها
	// ============================================
	fmt.Println("\n--- 8.2 Update Operations ---")

	// $set
	update1 := bson.M{
		"$set": bson.M{
			"name": "New Name",
			"age":  31,
		},
	}
	fmt.Printf("  \$set update: %+v\n", update1)

	// $inc
	update2 := bson.M{
		"$inc": bson.M{"views": 1},
	}
	fmt.Printf("  \$inc update: %+v\n", update2)

	// $push (به آرایه)
	update3 := bson.M{
		"$push": bson.M{"tags": "new-tag"},
	}
	fmt.Printf("  \$push update: %+v\n", update3)

	// ============================================
	// 8.3 Aggregation Pipeline
	// ============================================
	fmt.Println("\n--- 8.3 Aggregation Pipeline ---")

	pipeline := bson.A{
		bson.M{"$match": bson.M{"status": "active"}},
		bson.M{"$group": bson.M{
			"_id":   "$city",
			"count": bson.M{"$sum": 1},
		}},
		bson.M{"$sort": bson.M{"count": -1}},
		bson.M{"$limit": 10},
	}
	fmt.Printf("  Pipeline stages: %d\n", len(pipeline))
}

// ============================================================================
// بخش 9: تبدیل BSON به JSON و بالعکس
// ============================================================================

func demonstrateBSONToJSON() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🔄 BSON TO JSON CONVERSION")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 9.1 BSON به JSON
	// ============================================
	fmt.Println("\n--- 9.1 BSON to JSON ---")

	// ایجاد سند BSON
	doc := bson.D{
		{"name", "Ali"},
		{"age", 30},
		{"active", true},
		{"score", 98.5},
		{"tags", bson.A{"go", "mongodb"}},
	}

	// Marshal به BSON
	bsonData, _ := bson.Marshal(doc)

	// روش 1: Unmarshal به map و سپس Marshal به JSON
	var temp bson.M
	bson.Unmarshal(bsonData, &temp)
	jsonData, _ := json.MarshalIndent(temp, "", "  ")
	fmt.Printf("  Converted to JSON:\n%s\n", jsonData)

	// ============================================
	// 9.2 JSON به BSON
	// ============================================
	fmt.Println("\n--- 9.2 JSON to BSON ---")

	jsonStr := `{"name":"Sara","age":25,"city":"Shiraz"}`
	var jsonMap map[string]interface{}
	json.Unmarshal([]byte(jsonStr), &jsonMap)

	// تبدیل به BSON (با استفاده از Marshal)
	bsonFromJSON, _ := bson.Marshal(jsonMap)
	fmt.Printf("  JSON string: %s\n", jsonStr)
	fmt.Printf("  Converted to BSON: %x...\n", bsonFromJSON[:20])
}

// ============================================================================
// بخش 10: اشتباهات رایج
// ============================================================================

func demonstrateCommonMistakes() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("❌ COMMON MISTAKES WITH BSON")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println("\n❌ Mistake 1: Using D when order doesn't matter")
	fmt.Println("   // D is slower for lookups")
	fmt.Println("   ✅ Use M for unordered access")

	fmt.Println("\n❌ Mistake 2: Forgetting omitempty for optional fields")
	fmt.Println("   type User struct {")
	fmt.Println("       Name string `bson:\"name\"`  // always present")
	fmt.Println("   }")
	fmt.Println("   ✅ Use `bson:\"name,omitempty\"`")

	fmt.Println("\n❌ Mistake 3: Not handling ObjectID zero value")
	fmt.Println("   var id primitive.ObjectID")
	fmt.Println("   filter := bson.M{\"_id\": id}  // matches nothing?")
	fmt.Println("   ✅ Check id.IsZero() before using")

	fmt.Println("\n❌ Mistake 4: Using struct with unexported fields")
	fmt.Println("   type User struct {")
	fmt.Println("       name string  // unexported")
	fmt.Println("   }")
	fmt.Println("   ✅ Unexported fields are ignored")

	fmt.Println("\n❌ Mistake 5: Not validating bson.Raw before use")
	fmt.Println("   var raw bson.Raw")
	fmt.Println("   raw.Lookup(\"field\")  // may panic")
	fmt.Println("   ✅ raw.Validate() first")

	fmt.Println("\n❌ Mistake 6: Assuming map iteration order")
	fmt.Println("   for k, v := range bson.M{...}  // order undefined")
	fmt.Println("   ✅ Use bson.D if order matters")
}

// ============================================================================
// بخش 11: جمع‌بندی و جدول مرجع
// ============================================================================

func getMapKeys(m bson.M) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func main() {
	fmt.Println(stringsRepeat("=", 80))
	fmt.Println("🎯 COMPLETE bson PACKAGE GUIDE IN GO")
	fmt.Println("Binary JSON for MongoDB")
	fmt.Println(stringsRepeat("=", 80))

	// توجه: برای اجرای این کد نیاز به نصب پکیج دارید:
	// go get go.mongodb.org/mongo-driver/bson

	// بخش 1: انواع پایه
	demonstrateBasicTypes()

	// بخش 2: Marshal/Unmarshal
	demonstrateMarshalUnmarshal()

	// بخش 3: Struct Tags
	demonstrateStructTags()

	// بخش 4: Primitive Types
	demonstratePrimitiveTypes()

	// بخش 5: Raw BSON
	demonstrateRawBSON()

	// بخش 6: BSON Options
	demonstrateBSONOptions()

	// بخش 7: BSON vs JSON
	demonstrateBSONvsJSON()

	// بخش 8: Practical Use Cases
	demonstratePracticalUseCases()

	// بخش 9: BSON to JSON Conversion
	demonstrateBSONToJSON()

	// بخش 10: Common Mistakes
	demonstrateCommonMistakes()

	// بخش 11: Quick Reference
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📚 bson PACKAGE QUICK REFERENCE")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ BASIC TYPES                                                   │")
	fmt.Println("├─────────────────────────────────────────────────────────────────┤")
	fmt.Println("│ bson.D         - Ordered document ([]bson.E)                  │")
	fmt.Println("│ bson.M         - Unordered document (map[string]interface{})  │")
	fmt.Println("│ bson.A         - BSON array                                   │")
	fmt.Println("│ bson.E         - Single element (key, value)                  │")
	fmt.Println("│ bson.Raw       - Raw BSON bytes (no unmarshal)                │")
	fmt.Println("└─────────────────────────────────────────────────────────────────┘")

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ PRIMITIVE TYPES                                               │")
	fmt.Println("├─────────────────────────────────────────────────────────────────┤")
	fmt.Println("│ primitive.ObjectID    - MongoDB ObjectID (12 bytes)           │")
	fmt.Println("│ primitive.DateTime     - BSON datetime (milliseconds since epoch)│")
	fmt.Println("│ primitive.Binary       - Binary data                          │")
	fmt.Println("│ primitive.Decimal128   - 128-bit decimal (for finance)        │")
	fmt.Println("│ primitive.Timestamp    - MongoDB timestamp                    │")
	fmt.Println("│ primitive.Regex        - Regular expression                   │")
	fmt.Println("│ primitive.JavaScript   - JavaScript code                      │")
	fmt.Println("│ primitive.MinKey/MaxKey - Special comparison values           │")
	fmt.Println("└─────────────────────────────────────────────────────────────────┘")

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ MARSHAL/UNMARSHAL                                             │")
	fmt.Println("├─────────────────────────────────────────────────────────────────┤")
	fmt.Println("│ bson.Marshal(v)          - Go to BSON bytes                   │")
	fmt.Println("│ bson.Unmarshal(data, &v) - BSON bytes to Go                   │")
	fmt.Println("│ bson.MarshalExtJSON(v)   - Go to Extended JSON                │")
	fmt.Println("│ bson.UnmarshalExtJSON(data, &v) - Extended JSON to Go         │")
	fmt.Println("└─────────────────────────────────────────────────────────────────┘")

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ STRUCT TAGS                                                   │")
	fmt.Println("├─────────────────────────────────────────────────────────────────┤")
	fmt.Println("│ `bson:\"field_name\"`      - Specify BSON field name           │")
	fmt.Println("│ `bson:\"field,omitempty\"` - Omit if zero value                │")
	fmt.Println("│ `bson:\"field,minsize\"`   - Use smallest integer type         │")
	fmt.Println("│ `bson:\"field,truncate\"`  - Truncate decimal values           │")
	fmt.Println("│ `bson:\"field,inline\"`    - Flatten embedded struct/map       │")
	fmt.Println("│ `bson:\"-\"`               - Ignore field                      │")
	fmt.Println("└─────────────────────────────────────────────────────────────────┘")

	fmt.Println("\n💡 GOLDEN RULES:")
	fmt.Println("  1. Use bson.D when order matters (MongoDB commands)")
	fmt.Println("  2. Use bson.M for filters and updates (readability)")
	fmt.Println("  3. Always use omitempty for optional fields")
	fmt.Println("  4. Use primitive.ObjectID for MongoDB _id")
	fmt.Println("  5. Check ObjectID.IsZero() before using as filter")
	fmt.Println("  6. Use bson.Raw for partial document access (performance)")
	fmt.Println("  7. Validate bson.Raw with Validate() before Lookup")
	fmt.Println("  8. Unexported struct fields are ignored")
	fmt.Println("  9. Use minsize to save space for small numbers")
	fmt.Println("  10. Prefer structs over D/M for known schemas")

	fmt.Println("\n🎯 WHEN TO USE WHAT:")
	fmt.Println("  • struct with tags → Known schema, type safety")
	fmt.Println("  • bson.D → MongoDB commands, ordered operations")
	fmt.Println("  • bson.M → Filters, updates, projections")
	fmt.Println("  • bson.A → Array values in queries")
	fmt.Println("  • bson.Raw → High-performance partial reads")
	fmt.Println("  • primitive.ObjectID → MongoDB _id fields")

	fmt.Println("\n📦 INSTALLATION:")
	fmt.Println("  go get go.mongodb.org/mongo-driver/bson")
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
# نصب پکیج bson
go get go.mongodb.org/mongo-driver/bson

# اجرای برنامه
go run bson_complete_guide.go
*/
