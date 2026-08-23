// ============================================================================
// FILE: request_parsing_guide.go
// TITLE: راهنمای کامل Parsing درخواست‌ها - JSON, XML, Form-Data
// HOW TO RUN: go run request_parsing_guide.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - انواع داده‌های درخواست
// ============================================================================
//
// سه فرمت اصلی برای ارسال داده در درخواست‌های HTTP:
//
// 1. JSON (application/json)
//    - محبوب‌ترین فرمت برای APIها
//    - ساختار key-value
//    - پشتیبانی از انواع داده پیچیده
//
// 2. XML (application/xml)
//    - فرمت قدیمی‌تر، still used in SOAP, RSS, etc.
//    - ساختار tree-based با تگ‌ها
//
// 3. Form-Data (application/x-www-form-urlencoded, multipart/form-data)
//    - فرمت سنتی HTML forms
//    - urlencoded برای داده ساده
//    - multipart برای فایل و داده باینری
//
// قانون طلایی:
// "برای APIهای جدید از JSON استفاده کن.
//  برای فایل‌ها از multipart/form-data استفاده کن.
//  برای سازگاری با سیستم‌های قدیمی از XML استفاده کن."
// ============================================================================

package __request_parsing

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// ============================================================================
// بخش 1: مدل‌های داده مشترک
// ============================================================================

// User مدل کاربر برای JSON/XML
type User struct {
	ID        int       `json:"id" xml:"id"`
	Name      string    `json:"name" xml:"name"`
	Email     string    `json:"email" xml:"email"`
	Age       int       `json:"age" xml:"age"`
	IsActive  bool      `json:"is_active" xml:"is_active"`
	CreatedAt time.Time `json:"created_at" xml:"created_at"`
	Tags      []string  `json:"tags" xml:"tags>tag"`
}

// Product مدل محصول
type Product struct {
	ID    int     `json:"id" xml:"id"`
	Name  string  `json:"name" xml:"name"`
	Price float64 `json:"price" xml:"price"`
}

// Response پاسخ استاندارد API
type Response struct {
	Status  string      `json:"status" xml:"status"`
	Message string      `json:"message" xml:"message"`
	Data    interface{} `json:"data,omitempty" xml:"data,omitempty"`
}

// ============================================================================
// بخش 2: Parsing JSON (با net/http خالص)
// ============================================================================

// 2.1 Parsing JSON از Body
func parseJSONHandler(w http.ResponseWriter, r *http.Request) {
	// بررسی Content-Type
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	// محدود کردن حجم Body (1MB)
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var user User
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // جلوگیری از فیلدهای ناشناخته

	if err := decoder.Decode(&user); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// اعتبارسنجی ساده
	if user.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	response := Response{
		Status:  "success",
		Message: "User created",
		Data:    user,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// 2.2 Parsing JSON از Query String
func parseJSONQueryHandler(w http.ResponseWriter, r *http.Request) {
	// خواندن JSON از query parameter ?data={"name":"Ali"}
	jsonStr := r.URL.Query().Get("data")
	if jsonStr == "" {
		http.Error(w, "Missing data parameter", http.StatusBadRequest)
		return
	}

	var user User
	if err := json.Unmarshal([]byte(jsonStr), &user); err != nil {
		http.Error(w, "Invalid JSON in query: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// 2.3 Parsing JSON Stream (برای داده بزرگ)
func parseJSONStreamHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	// خواندن token اول (باز شدن آرایه)
	token, err := decoder.Token()
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if delim, ok := token.(json.Delim); !ok || delim != '[' {
		http.Error(w, "Expected array", http.StatusBadRequest)
		return
	}

	var users []User
	for decoder.More() {
		var user User
		if err := decoder.Decode(&user); err != nil {
			http.Error(w, "Error parsing user: "+err.Error(), http.StatusBadRequest)
			return
		}
		users = append(users, user)
	}

	// خواندن token آخر (بستن آرایه)
	decoder.Token()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// 2.4 Parsing JSON به map (برای داده داینامیک)
func parseJSONMapHandler(w http.ResponseWriter, r *http.Request) {
	var data map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// دسترسی به فیلدها
	name, _ := data["name"].(string)
	age, _ := data["age"].(float64)

	response := map[string]interface{}{
		"received_name": name,
		"received_age":  age,
		"all_fields":    data,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ============================================================================
// بخش 3: Parsing XML (با net/http خالص)
// ============================================================================

// 3.1 Parsing XML از Body
func parseXMLHandler(w http.ResponseWriter, r *http.Request) {
	// بررسی Content-Type
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/xml" && contentType != "text/xml" {
		http.Error(w, "Content-Type must be application/xml", http.StatusUnsupportedMediaType)
		return
	}

	var user User
	if err := xml.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid XML: "+err.Error(), http.StatusBadRequest)
		return
	}

	response := Response{
		Status:  "success",
		Message: "XML parsed successfully",
		Data:    user,
	}

	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
	xml.NewEncoder(w).Encode(response)
}

// 3.2 Parsing XML با تگ‌های سفارشی
type CustomXMLUser struct {
	XMLName xml.Name `xml:"user"`
	ID      int      `xml:"id,attr"`       // attribute
	Name    string   `xml:"name"`          // element
	Email   string   `xml:"email>address"` // nested
	Age     int      `xml:"age,omitempty"` // omit if zero
}

func parseCustomXMLHandler(w http.ResponseWriter, r *http.Request) {
	var user CustomXMLUser
	if err := xml.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid XML: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	xml.NewEncoder(w).Encode(user)
}

// ============================================================================
// بخش 4: Parsing Form-Data (x-www-form-urlencoded)
// ============================================================================

// 4.1 Parsing Form URL Encoded
func parseFormURLEncodedHandler(w http.ResponseWriter, r *http.Request) {
	// بررسی Content-Type
	if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		http.Error(w, "Content-Type must be application/x-www-form-urlencoded", http.StatusUnsupportedMediaType)
		return
	}

	// Parsing form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form: "+err.Error(), http.StatusBadRequest)
		return
	}

	// خواندن مقادیر
	name := r.FormValue("name")
	email := r.FormValue("email")
	age, _ := strconv.Atoi(r.FormValue("age"))

	// خواندن چندین مقدار (مثلاً checkbox)
	hobbies := r.Form["hobbies"]

	response := map[string]interface{}{
		"name":    name,
		"email":   email,
		"age":     age,
		"hobbies": hobbies,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// 4.2 Parsing Form با Validation
func parseFormWithValidationHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Validation
	errors := make(map[string]string)

	name := strings.TrimSpace(r.FormValue("name"))
	if name == "" {
		errors["name"] = "Name is required"
	}

	email := strings.TrimSpace(r.FormValue("email"))
	if email == "" {
		errors["email"] = "Email is required"
	} else if !strings.Contains(email, "@") {
		errors["email"] = "Invalid email format"
	}

	ageStr := r.FormValue("age")
	age := 0
	if ageStr != "" {
		var err error
		age, err = strconv.Atoi(ageStr)
		if err != nil || age < 0 || age > 150 {
			errors["age"] = "Age must be between 0 and 150"
		}
	}

	if len(errors) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors)
		return
	}

	response := map[string]interface{}{
		"status": "success",
		"user": map[string]interface{}{
			"name":  name,
			"email": email,
			"age":   age,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ============================================================================
// بخش 5: Parsing Multipart Form-Data (برای فایل‌ها)
// ============================================================================

// 5.1 Parsing Multipart Form (با فایل)
func parseMultipartHandler(w http.ResponseWriter, r *http.Request) {
	// محدودیت حجم (10 MB)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Error parsing multipart form: "+err.Error(), http.StatusBadRequest)
		return
	}

	// خواندن فیلدهای متنی
	name := r.FormValue("name")
	description := r.FormValue("description")

	// خواندن فایل
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error reading file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// خواندن محتویات فایل (نمایش)
	content, _ := io.ReadAll(file)
	contentPreview := string(content)
	if len(contentPreview) > 100 {
		contentPreview = contentPreview[:100] + "..."
	}

	response := map[string]interface{}{
		"name":         name,
		"description":  description,
		"filename":     header.Filename,
		"file_size":    header.Size,
		"file_preview": contentPreview,
		"content_type": header.Header.Get("Content-Type"),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// 5.2 آپلود فایل با ذخیره‌سازی
func uploadFileHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error reading file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// ایجاد دایرکتوری uploads اگر وجود ندارد
	os.MkdirAll("./uploads", 0755)

	// ذخیره فایل
	dst, err := os.Create(fmt.Sprintf("./uploads/%d_%s", time.Now().Unix(), header.Filename))
	if err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Error copying file", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"status":   "success",
		"message":  "File uploaded successfully",
		"filename": header.Filename,
		"size":     header.Size,
		"path":     dst.Name(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// 5.3 آپلود چند فایل همزمان
func uploadMultipleFilesHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(50 << 20); err != nil { // 50 MB
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["files"]
	var uploadedFiles []map[string]interface{}

	os.MkdirAll("./uploads", 0755)

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			continue
		}

		dst, err := os.Create(fmt.Sprintf("./uploads/%d_%s", time.Now().UnixNano(), fileHeader.Filename))
		if err != nil {
			file.Close()
			continue
		}

		io.Copy(dst, file)
		dst.Close()
		file.Close()

		uploadedFiles = append(uploadedFiles, map[string]interface{}{
			"filename": fileHeader.Filename,
			"size":     fileHeader.Size,
		})
	}

	response := map[string]interface{}{
		"status":   "success",
		"uploaded": len(uploadedFiles),
		"files":    uploadedFiles,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ============================================================================
// بخش 6: Parsing ترکیبی (JSON + Form + Query)
// ============================================================================

func parseCombinedHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Query Parameters
	queryParam := r.URL.Query().Get("param")

	// 2. URL Parameters (با Chi/Mux)
	// id := chi.URLParam(r, "id")

	// 3. Headers
	userAgent := r.Header.Get("User-Agent")

	// 4. Body (JSON یا Form)
	var bodyData map[string]interface{}
	contentType := r.Header.Get("Content-Type")

	if strings.Contains(contentType, "application/json") {
		json.NewDecoder(r.Body).Decode(&bodyData)
	} else if strings.Contains(contentType, "application/x-www-form-urlencoded") {
		r.ParseForm()
		bodyData = make(map[string]interface{})
		for k, v := range r.Form {
			if len(v) == 1 {
				bodyData[k] = v[0]
			} else {
				bodyData[k] = v
			}
		}
	}

	response := map[string]interface{}{
		"query_param": queryParam,
		"user_agent":  userAgent,
		"body":        bodyData,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ============================================================================
// بخش 7: Chi Router Examples
// ============================================================================

/*
// Chi Router با JSON Parsing
func chiJSONHandler(c *chi.Router) {
	r := chi.NewRouter()

	// JSON
	r.Post("/users", func(w http.ResponseWriter, r *http.Request) {
		var user User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	})

	// XML
	r.Post("/users/xml", func(w http.ResponseWriter, r *http.Request) {
		var user User
		if err := xml.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/xml")
		xml.NewEncoder(w).Encode(user)
	})

	// Form
	r.Post("/form", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		name := r.FormValue("name")
		w.Write([]byte("Name: " + name))
	})

	// Multipart
	r.Post("/upload", func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(10 << 20)
		file, _, _ := r.FormFile("file")
		defer file.Close()
		// process file
		w.Write([]byte("File uploaded"))
	})
}
*/

// ============================================================================
// بخش 8: Gin Framework Examples
// ============================================================================

/*
// Gin Router با Parsing خودکار
func ginExamples() {
	r := gin.Default()

	// 1. JSON Binding (خودکار)
	r.POST("/users", func(c *gin.Context) {
		var user User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)
	})

	// 2. XML Binding
	r.POST("/users/xml", func(c *gin.Context) {
		var user User
		if err := c.ShouldBindXML(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.XML(http.StatusOK, user)
	})

	// 3. Form Binding
	r.POST("/form", func(c *gin.Context) {
		var form struct {
			Name  string `form:"name" binding:"required"`
			Email string `form:"email" binding:"required,email"`
			Age   int    `form:"age" binding:"min=0,max=150"`
		}
		if err := c.ShouldBind(&form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, form)
	})

	// 4. Multipart Form
	r.POST("/upload", func(c *gin.Context) {
		// Text fields
		name := c.PostForm("name")
		description := c.PostForm("description")

		// File
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Save file
		c.SaveUploadedFile(file, "./uploads/"+file.Filename)

		c.JSON(http.StatusOK, gin.H{
			"name":        name,
			"description": description,
			"filename":    file.Filename,
			"size":        file.Size,
		})
	})

	// 5. Multiple files
	r.POST("/upload-multiple", func(c *gin.Context) {
		form, _ := c.MultipartForm()
		files := form.File["files"]

		var uploaded []string
		for _, file := range files {
			c.SaveUploadedFile(file, "./uploads/"+file.Filename)
			uploaded = append(uploaded, file.Filename)
		}

		c.JSON(http.StatusOK, gin.H{"uploaded": uploaded})
	})

	// 6. Query + Form + JSON ترکیبی
	r.POST("/combined", func(c *gin.Context) {
		queryParam := c.Query("param")
		headerValue := c.GetHeader("X-Custom")

		var jsonBody map[string]interface{}
		c.ShouldBindJSON(&jsonBody)

		c.JSON(http.StatusOK, gin.H{
			"query":  queryParam,
			"header": headerValue,
			"body":   jsonBody,
		})
	})

	// 7. ShouldBind با تشخیص خودکار Content-Type
	r.POST("/auto-bind", func(c *gin.Context) {
		var user User
		// Gin به طور خودکار بر اساس Content-Type تصمیم می‌گیرد
		if err := c.ShouldBind(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)
	})

	r.Run(":8080")
}
*/

// ============================================================================
// بخش 9: Validation و Error Handling
// ============================================================================

// 9.1 Validation Helper
func validateUser(user User) map[string]string {
	errors := make(map[string]string)

	if user.Name == "" {
		errors["name"] = "Name is required"
	} else if len(user.Name) < 3 {
		errors["name"] = "Name must be at least 3 characters"
	}

	if user.Email == "" {
		errors["email"] = "Email is required"
	} else if !strings.Contains(user.Email, "@") {
		errors["email"] = "Invalid email format"
	}

	if user.Age < 0 || user.Age > 150 {
		errors["age"] = "Age must be between 0 and 150"
	}

	return errors
}

// 9.2 Generic Parse Function
func parseJSON(r *http.Request, target interface{}) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("read body: %w", err)
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, target); err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}
	return nil
}

// ============================================================================
// بخش 10: سرور کامل با تمام نمونه‌ها
// ============================================================================

func runFullServer() {
	mux := http.NewServeMux()

	// JSON endpoints
	mux.HandleFunc("/json", parseJSONHandler)
	mux.HandleFunc("/json/query", parseJSONQueryHandler)
	mux.HandleFunc("/json/stream", parseJSONStreamHandler)
	mux.HandleFunc("/json/map", parseJSONMapHandler)

	// XML endpoints
	mux.HandleFunc("/xml", parseXMLHandler)
	mux.HandleFunc("/xml/custom", parseCustomXMLHandler)

	// Form endpoints
	mux.HandleFunc("/form/urlencoded", parseFormURLEncodedHandler)
	mux.HandleFunc("/form/validate", parseFormWithValidationHandler)

	// Multipart endpoints
	mux.HandleFunc("/upload", uploadFileHandler)
	mux.HandleFunc("/upload/multiple", uploadMultipleFilesHandler)
	mux.HandleFunc("/multipart", parseMultipartHandler)

	// Combined
	mux.HandleFunc("/combined", parseCombinedHandler)

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

// ============================================================================
// بخش 11: Client Examples (ارسال درخواست)
// ============================================================================

// 11.1 Send JSON Request
func sendJSONRequest() {
	user := User{
		Name:  "Ali Rezaei",
		Email: "ali@example.com",
		Age:   30,
	}

	jsonData, _ := json.Marshal(user)

	resp, err := http.Post("http://localhost:8080/json", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	defer resp.Body.Close()

	var response Response
	json.NewDecoder(resp.Body).Decode(&response)
	log.Printf("Response: %+v", response)
}

// 11.2 Send XML Request
func sendXMLRequest() {
	user := User{
		Name:  "Sara Mohammadi",
		Email: "sara@example.com",
		Age:   25,
	}

	xmlData, _ := xml.Marshal(user)

	resp, err := http.Post("http://localhost:8080/xml", "application/xml", bytes.NewBuffer(xmlData))
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	log.Printf("Response: %s", body)
}

// 11.3 Send Form URL Encoded
func sendFormRequest() {
	data := url.Values{}
	data.Set("name", "Ali Rezaei")
	data.Set("email", "ali@example.com")
	data.Set("age", "30")
	data.Set("hobbies", "coding")
	data.Add("hobbies", "reading")

	resp, err := http.PostForm("http://localhost:8080/form/urlencoded", data)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	log.Printf("Response: %s", body)
}

// 11.4 Send Multipart Form with File
func sendMultipartRequest() {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Text fields
	writer.WriteField("name", "Ali Rezaei")
	writer.WriteField("description", "Profile picture")

	// File
	file, _ := os.Open("test.txt")
	defer file.Close()

	part, _ := writer.CreateFormFile("file", "test.txt")
	io.Copy(part, file)

	writer.Close()

	resp, err := http.Post("http://localhost:8080/upload", writer.FormDataContentType(), body)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	defer resp.Body.Close()

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	log.Printf("Response: %+v", response)
}

// ============================================================================
// بخش 12: جدول مرجع سریع
// ============================================================================

func quickReference() {
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("📚 QUICK REFERENCE")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────┐
│ CONTENT-TYPE                    │ HOW TO PARSE                  │
├─────────────────────────────────┼───────────────────────────────┤
│ application/json                │ json.NewDecoder(r.Body)      │
│ application/xml                 │ xml.NewDecoder(r.Body)       │
│ application/x-www-form-urlencoded│ r.ParseForm() / r.FormValue()│
│ multipart/form-data             │ r.ParseMultipartForm()       │
└─────────────────────────────────┴───────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ LIBRARY    │ JSON              │ XML               │ Form       │
├────────────┼───────────────────┼───────────────────┼────────────┤
│ net/http   │ json.Decoder      │ xml.Decoder       │ ParseForm  │
│ Chi        │ Same as net/http  │ Same as net/http  │ Same       │
│ Gin        │ c.ShouldBindJSON()│ c.ShouldBindXML() │ c.ShouldBind│
└────────────┴───────────────────┴───────────────────┴────────────┘

💡 BEST PRACTICES:

   1. Always check Content-Type before parsing
   2. Set body size limits (http.MaxBytesReader)
   3. Use DisallowUnknownFields() for JSON to catch typos
   4. Close request body after reading (defer)
   5. Validate data after parsing
   6. Return meaningful error messages
   7. Use struct tags for mapping
   8. Handle multipart forms with appropriate size limits
   9. For large files, stream directly to disk
   10. Use ShouldBind in Gin for automatic detection
`)
}

// ============================================================================
// بخش 13: main
// ============================================================================

func main() {
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("🎯 REQUEST PARSING GUIDE")
	fmt.Println("JSON | XML | Form-Data | Multipart")
	fmt.Println(strings.Repeat("=", 80))

	quickReference()

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🚀 Starting server on :8080")
	fmt.Println("Test endpoints:")
	fmt.Println("  POST /json              - JSON data")
	fmt.Println("  POST /xml               - XML data")
	fmt.Println("  POST /form/urlencoded   - Form URL Encoded")
	fmt.Println("  POST /upload            - Multipart form with file")
	fmt.Println("  GET  /health            - Health check")
	fmt.Println(strings.Repeat("=", 80))

	runFullServer()
}

// تابع کمکی برای تکرار رشته
func stringsRepeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}
