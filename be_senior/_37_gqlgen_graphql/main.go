// ============================================================================
// FILE: gqlgen_complete_guide.go
// TITLE: راهنمای کامل GraphQL با gqlgen در Go
// HOW TO RUN: go run gqlgen_complete_guide.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - GraphQL و gqlgen چیست؟
// ============================================================================
//
// GraphQL یک زبان کوئری برای APIها است که توسط فیسبوک ساخته شده است.
// برخلاف REST که چندین endpoint دارد، GraphQL یک endpoint واحد با کوئری‌های قابل تنظیم دارد.
//
// مزایای GraphQL نسبت به REST:
// 1. جلوگیری از over-fetching (دریافت داده اضافی)
// 2. جلوگیری از under-fetching (نیاز به چندین درخواست)
// 3. تایپ قوی با Schema Definition Language
// 4. یک endpoint برای تمام عملیات
// 5. نسخه‌بندی خودکار (با اضافه کردن فیلدهای جدید)
//
// gqlgen کتابخانه‌ای برای ساخت GraphQL سرور در Go است با ویژگی‌های:
// - Schema first: ابتدا Schema را تعریف می‌کنی، سپس کد生成 می‌شود
// - Type safe: هرگز map[string]interface{} نمی‌بینی
// - Code generation: کدهای تکراری自动生成 می‌شوند [citation:4]
//
// قانون طلایی:
// "اول Schema را بنویس، سپس gqlgen generate را اجرا کن.
//  فقط resolverها را دستی بنویس - بقیه کدها生成 می‌شوند.
//  برای کنترل دقیق تر، از custom models و resolverهای explicit استفاده کن."
// ============================================================================

package model

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/99designs/gqlgen/graphql"
)

// ============================================================================
// بخش 1: نصب و راه‌اندازی پروژه
// ============================================================================

/*
نصب:

1. ایجاد پروژه جدید:
   $ mkdir graphql-demo
   $ cd graphql-demo
   $ go mod init graphql-demo

2. ایجاد فایل tools.go برای ردیابی وابستگی (workaround):
*/
//go:build tools

//package tools

import (
_ "github.com/99designs/gqlgen"
_ "github.com/99designs/gqlgen/graphql/introspection"
)

/*
   3. نصب وابستگی‌ها:
      $ go mod tidy

   4. مقداردهی اولیه پروژه gqlgen:
      $ go run github.com/99designs/gqlgen init

   پس از اجرا، ساختار پوشه به این صورت خواهد بود:
   ├── go.mod
   ├── go.sum
   ├── gqlgen.yml
   ├── server.go
   ├── graph/
   │   ├── generated.go
   │   ├── resolver.go
   │   ├── schema.graphqls
   │   └── schema.resolvers.go
   └── graph/model/
       └── models_gen.go

   فایل‌های مهم:
   - schema.graphqls: تعریف Schema GraphQL
   - resolver.go: مدیریت state و وابستگی‌ها
   - schema.resolvers.go: پیاده‌سازی منطق resolverها [citation:1]
   - gqlgen.yml: فایل تنظیمات gqlgen
*/

// ============================================================================
// بخش 2: تعریف Schema (schema.graphqls)
// ============================================================================

/*
   # فایل: graph/schema.graphqls

   # 2.1 انواع پایه (Types)
   type Todo {
       id: ID!
       text: String!
       done: Boolean!
       user: User!
       createdAt: Time!
   }

   type User {
       id: ID!
       name: String!
       email: String!
       todos: [Todo!]!
   }

   # 2.2 Query (عملیات خواندن)
   type Query {
       todos: [Todo!]!
       todo(id: ID!): Todo
       users: [User!]!
       user(id: ID!): User
   }

   # 2.3 Mutation (عملیات نوشتن)
   type Mutation {
       createTodo(input: NewTodo!): Todo!
       updateTodo(id: ID!, input: UpdateTodo!): Todo!
       deleteTodo(id: ID!): Boolean!
       createUser(input: NewUser!): User!
   }

   # 2.4 Input Types (برای mutations)
   input NewTodo {
       text: String!
       userId: ID!
   }

   input UpdateTodo {
       text: String
       done: Boolean
   }

   input NewUser {
       name: String!
       email: String!
   }

   # 2.5 Scalarهای سفارشی
   scalar Time
   scalar JSON
   scalar Upload

   # 2.6 Enumها
   enum TodoStatus {
       PENDING
       IN_PROGRESS
       COMPLETED
   }

   # 2.7 Interfaceها
   interface Node {
       id: ID!
   }

   # 2.8 Unionها
   union SearchResult = User | Todo

   # 2.9 Directiveها (برای authorization)
   directive @auth(requires: Role!) on OBJECT | FIELD_DEFINITION

   enum Role {
       ADMIN
       USER
       GUEST
   }

   # 2.10 انواع خروجی با pagination
   type TodoConnection {
       edges: [TodoEdge!]!
       pageInfo: PageInfo!
       totalCount: Int!
   }

   type TodoEdge {
       node: Todo!
       cursor: String!
   }

   type PageInfo {
       hasNextPage: Boolean!
       hasPreviousPage: Boolean!
       startCursor: String
       endCursor: String
   }
*/

// ============================================================================
// بخش 3: مدل‌های Go (مدل‌های دستی برای کنترل بیشتر)
// ============================================================================

// Model تعریف مدل‌های سفارشی (برای جلوگیری از استفاده از models_gen.go)
// می‌توان در gqlgen.yml این مدل‌ها را معرفی کرد تا gqlgen از آن‌ها استفاده کند

package model

import (
"time"
)

// Todo مدل Todo
type Todo struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	Done      bool      `json:"done"`
	UserID    string    `json:"userId"`
	User      *User     `json:"user"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// User مدل User
type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
	Todos     []*Todo   `json:"todos"`
}

// NewTodoInput برای mutation createTodo
type NewTodoInput struct {
	Text   string `json:"text"`
	UserID string `json:"userId"`
}

// UpdateTodoInput برای mutation updateTodo
type UpdateTodoInput struct {
	Text *string `json:"text"`
	Done *bool   `json:"done"`
}

// NewUserInput برای mutation createUser
type NewUserInput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// PageInfo برای pagination
type PageInfo struct {
	HasNextPage     bool    `json:"hasNextPage"`
	HasPreviousPage bool    `json:"hasPreviousPage"`
	StartCursor     *string `json:"startCursor"`
	EndCursor       *string `json:"endCursor"`
}

// ============================================================================
// بخش 4: Resolver (مدیریت State و وابستگی‌ها)
// ============================================================================

// Resolver ساختار اصلی که state برنامه را نگهداری می‌کند
// اینجا می‌توان database connection، cache، logger و غیره را اضافه کرد

package graph

import (
"context"
"crypto/rand"
"encoding/base64"
"fmt"
"graphql-demo/graph/model"
"sync"
"time"
)

// Resolver struct برای مدیریت state
type Resolver struct {
	todos     map[string]*model.Todo
	users     map[string]*model.User
	mu        sync.RWMutex
	// می‌توان وابستگی‌های دیگر را اضافه کرد:
	// db     *sql.DB
	// redis  *redis.Client
	// logger *log.Logger
}

// NewResolver سازنده Resolver
func NewResolver() *Resolver {
	r := &Resolver{
		todos: make(map[string]*model.Todo),
		users: make(map[string]*model.User),
	}
	// اضافه کردن داده نمونه
	r.seedData()
	return r
}

// seedData افزودن داده نمونه
func (r *Resolver) seedData() {
	user1 := &model.User{
		ID:        r.generateID(),
		Name:      "Ali Rezaei",
		Email:     "ali@example.com",
		CreatedAt: time.Now(),
	}
	user2 := &model.User{
		ID:        r.generateID(),
		Name:      "Sara Mohammadi",
		Email:     "sara@example.com",
		CreatedAt: time.Now(),
	}
	r.users[user1.ID] = user1
	r.users[user2.ID] = user2

	todo1 := &model.Todo{
		ID:        r.generateID(),
		Text:      "Learn GraphQL",
		Done:      false,
		UserID:    user1.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	todo2 := &model.Todo{
		ID:        r.generateID(),
		Text:      "Build API with gqlgen",
		Done:      false,
		UserID:    user1.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	r.todos[todo1.ID] = todo1
	r.todos[todo2.ID] = todo2
}

// generateID تولید ID یکتا
func (r *Resolver) generateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// ============================================================================
// بخش 5: پیاده‌سازی Query Resolverها
// ============================================================================

// این توابع توسط gqlgen generate در فایل schema.resolvers.go قرار می‌گیرند
// اما برای نمایش، اینجا پیاده‌سازی می‌شوند

// Todos بازگرداندن لیست همه Todoها
func (r *Resolver) Todos(ctx context.Context) ([]*model.Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	todos := make([]*model.Todo, 0, len(r.todos))
	for _, todo := range r.todos {
		todos = append(todos, todo)
	}
	return todos, nil
}

// Todo بازگرداندن یک Todo بر اساس ID
func (r *Resolver) Todo(ctx context.Context, id string) (*model.Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	todo, ok := r.todos[id]
	if !ok {
		return nil, fmt.Errorf("todo not found")
	}
	return todo, nil
}

// Users بازگرداندن لیست همه Userها
func (r *Resolver) Users(ctx context.Context) ([]*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]*model.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}
	return users, nil
}

// User بازگرداندن یک User بر اساس ID
func (r *Resolver) User(ctx context.Context, id string) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.users[id]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

// ============================================================================
// بخش 6: پیاده‌سازی Mutation Resolverها
// ============================================================================

// CreateTodo ایجاد Todo جدید
func (r *Resolver) CreateTodo(ctx context.Context, input model.NewTodoInput) (*model.Todo, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// بررسی وجود User
	user, ok := r.users[input.UserID]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}

	// ایجاد Todo جدید
	todo := &model.Todo{
		ID:        r.generateID(),
		Text:      input.Text,
		Done:      false,
		UserID:    input.UserID,
		User:      user,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	r.todos[todo.ID] = todo
	return todo, nil
}

// UpdateTodo به‌روزرسانی Todo
func (r *Resolver) UpdateTodo(ctx context.Context, id string, input model.UpdateTodoInput) (*model.Todo, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	todo, ok := r.todos[id]
	if !ok {
		return nil, fmt.Errorf("todo not found")
	}

	if input.Text != nil {
		todo.Text = *input.Text
	}
	if input.Done != nil {
		todo.Done = *input.Done
	}
	todo.UpdatedAt = time.Now()

	return todo, nil
}

// DeleteTodo حذف Todo
func (r *Resolver) DeleteTodo(ctx context.Context, id string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.todos[id]; !ok {
		return false, fmt.Errorf("todo not found")
	}

	delete(r.todos, id)
	return true, nil
}

// CreateUser ایجاد User جدید
func (r *Resolver) CreateUser(ctx context.Context, input model.NewUserInput) (*model.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// بررسی عدم تکراری بودن ایمیل
	for _, u := range r.users {
		if u.Email == input.Email {
			return nil, fmt.Errorf("email already exists")
		}
	}

	user := &model.User{
		ID:        r.generateID(),
		Name:      input.Name,
		Email:     input.Email,
		CreatedAt: time.Now(),
	}

	r.users[user.ID] = user
	return user, nil
}

// ============================================================================
// بخش 7: Field Resolverها (برای روابط)
// ============================================================================

// این resolverها برای فیلدهایی که در struct نیستند استفاده می‌شوند
// مثلاً User.Todos در مدل User وجود ندارد، باید به صورت جداگانه resolver شود

// UserTodos resolver برای فیلد todos در User
func (r *Resolver) UserTodos(ctx context.Context, obj *model.User) ([]*model.Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	todos := make([]*model.Todo, 0)
	for _, todo := range r.todos {
		if todo.UserID == obj.ID {
			todos = append(todos, todo)
		}
	}
	return todos, nil
}

// TodoUser resolver برای فیلد user در Todo
func (r *Resolver) TodoUser(ctx context.Context, obj *model.Todo) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.users[obj.UserID]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

// ============================================================================
// بخش 8: جلوگیری از N+1 Problem با DataLoader
// ============================================================================

/*
DataLoader برای جلوگیری از مشکل N+1 در resolverهای تو در تو استفاده می‌شود.
مثال: درخواست زیر 1+N کوئری دیتابیس می‌زند بدون DataLoader
query {
  todos {
    user { name }  // برای هر todo یک کوئری جداگانه
  }
}

با DataLoader، درخواست‌های مشابه batch می‌شوند.
*/

// UserLoader برای batch کردن درخواست‌های User
type UserLoader struct {
	mu    sync.Mutex
	batch []string
	fetch func(ids []string) ([]*model.User, []error)
	cache map[string]*model.User
}

// NewUserLoader ایجاد UserLoader جدید
func NewUserLoader(fetch func(ids []string) ([]*model.User, []error)) *UserLoader {
	return &UserLoader{
		fetch: fetch,
		cache: make(map[string]*model.User),
	}
}

// Load بارگذاری یک User
func (l *UserLoader) Load(ctx context.Context, id string) (*model.User, error) {
	l.mu.Lock()
	l.batch = append(l.batch, id)
	l.mu.Unlock()

	// در عمل باید از timer و channel استفاده کرد
	return l.fetchSingle(id)
}

func (l *UserLoader) fetchSingle(id string) (*model.User, error) {
	if user, ok := l.cache[id]; ok {
		return user, nil
	}
	// در عمل باید batch را اجرا کرد
	return nil, nil
}

// ============================================================================
// بخش 9: Custom Scalarها
// ============================================================================

// Custom Scalarها باید در gqlgen.yml معرفی شوند

// TimeScalar scalar سفارشی برای Time
type TimeScalar struct{}

// MarshalTime مارشال کردن Time به JSON
func MarshalTime(t time.Time) interface{} {
	return t.Format(time.RFC3339Nano)
}

// UnmarshalTime آنمارشال کردن JSON به Time
func UnmarshalTime(v interface{}) (time.Time, error) {
	switch v := v.(type) {
	case string:
		return time.Parse(time.RFC3339Nano, v)
	case time.Time:
		return v, nil
	default:
		return time.Time{}, fmt.Errorf("invalid time format")
	}
}

// JSONScalar scalar سفارشی برای JSON
type JSONScalar map[string]interface{}

func MarshalJSON(j JSONScalar) interface{} {
	return j
}

func UnmarshalJSON(v interface{}) (JSONScalar, error) {
	switch v := v.(type) {
	case map[string]interface{}:
		return JSONScalar(v), nil
	default:
		return nil, fmt.Errorf("invalid JSON format")
	}
}

// ============================================================================
// بخش 10: Middleware و Authentication
// ============================================================================

// AuthContext کلید برای ذخیره user در context
type contextKey string

const UserContextKey contextKey = "user"

// AuthMiddleware میدلور احراز هویت
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// استخراج token از هدر
		token := r.Header.Get("Authorization")

		// اعتبارسنجی token و یافتن user
		user, err := validateToken(token)
		if err == nil && user != nil {
			ctx := context.WithValue(r.Context(), UserContextKey, user)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}

func validateToken(token string) (*model.User, error) {
	// در عمل token را validate کن
	return nil, nil
}

// Directive احراز هویت در GraphQL
func AuthDirective(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	user := ctx.Value(UserContextKey)
	if user == nil {
		return nil, fmt.Errorf("authentication required")
	}
	return next(ctx)
}

// ============================================================================
// بخش 11: Error Handling
// ============================================================================

// GraphQLError خطای سفارشی GraphQL
type GraphQLError struct {
	Message    string                 `json:"message"`
	Code       string                 `json:"code"`
	Path       []interface{}          `json:"path,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
}

func (e *GraphQLError) Error() string {
	return e.Message
}

// خطاهای از پیش تعریف شده
var (
	ErrNotFound     = &GraphQLError{Code: "NOT_FOUND", Message: "Resource not found"}
	ErrUnauthorized = &GraphQLError{Code: "UNAUTHORIZED", Message: "Unauthorized"}
	ErrBadRequest   = &GraphQLError{Code: "BAD_REQUEST", Message: "Invalid request"}
	ErrConflict     = &GraphQLError{Code: "CONFLICT", Message: "Resource already exists"}
)

// ============================================================================
// بخش 12: فایل تنظیمات gqlgen.yml
// ============================================================================

/*
# فایل: gqlgen.yml
# تنظیمات اصلی gqlgen

schema:
  - graph/schema.graphqls

exec:
  filename: graph/generated.go
  package: graph

model:
  filename: graph/model/models_gen.go
  package: model

resolver:
  layout: follow-schema
  dir: graph
  package: graph
  filename_template: "{name}.resolvers.go"

# مدل‌های سفارشی (برای استفاده از structهای خودمان به جای مدل‌های generated)
models:
  ID:
    model:
      - github.com/99designs/gqlgen/graphql.String
  Time:
    model:
      - graphql-demo/graph/model.Time
  User:
    model: graphql-demo/graph/model.User
    fields:
      todos:
        resolver: true
  Todo:
    model: graphql-demo/graph/model.Todo
    fields:
      user:
        resolver: true

# اتوماتیک binding برای types
autobind:
  - "graphql-demo/graph/model"

# تنظیمات introspection
introspection:
  enabled: true

# تنظیمات directive
directives:
  auth:
    skip_runtime: false
*/

// ============================================================================
// بخش 13: سرور GraphQL
// ============================================================================

/*
// فایل: server.go
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"graphql-demo/graph"
)

func main() {
	// ایجاد Resolver با state
	resolver := graph.NewResolver()

	// ایجاد GraphQL server
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	// Routes
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	// شروع سرور
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("GraphQL server running on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
*/

// ============================================================================
// بخش 14: Performance Tips (نکات عملکردی)
// ============================================================================

func performanceTips() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("⚡ GQLGEN PERFORMANCE TIPS")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────┐
│ PERFORMANCE OPTIMIZATIONS                                     │
├─────────────────────────────────────────────────────────────────┤
│                                                                │
│ 1. استفاده از DataLoader برای جلوگیری از N+1                  │
│    - Batch کردن درخواست‌های دیتابیس                           │
│    - کاهش تعداد کوئری‌ها از O(N) به O(1)                       │
│                                                                │
│ 2. Complex Query Limiting                                      │
│    - محدود کردن عمق کوئری‌ها                                   │
│    - محدود کردن تعداد فیلدها                                   │
│    - جلوگیری از حملات DoS                                     │
│                                                                │
│ 3. Caching Strategy                                            │
│    - کش کردن Queryهای تکراری                                   │
│    - استفاده از CDN برای persisted queries                     │
│    - کش کردن سطح resolver                                     │
│                                                                │
│ 4. Connection Pooling                                          │
│    - استفاده از connection pool برای دیتابیس                   │
│    - تنظیم مناسب max open connections                         │
│                                                                │
│ 5. Query Cost Analysis                                         │
│    - محاسبه cost هر query                                      │
│    - رد کردن queryهای expensive                                │
│    - پیاده‌سازی rate limiting                                 │
│                                                                │
│ 6. Persisted Queries                                           │
│    - ذخیره queryهای پرکاربرد در سرور                          │
│    - کاهش سایز payload                                         │
│    - افزایش امنیت                                              │
│                                                                │
│ 7. gzip Compression                                            │
│    - فشرده‌سازی پاسخ‌ها                                        │
│    - کاهش پهنای باند مصرفی                                     │
│                                                                │
│ 8. Benchmarking                                                │
│    - استفاده از benchmarkهای داخلی gqlgen                     │
│    - تست با ابزارهایی مثل k6                                  │
│                                                                │
└─────────────────────────────────────────────────────────────────┘

📊 BENCHMARK RESULTS (gqlgen vs REST):
   • Query response time: 25-35% faster than REST [citation:3]
   • Data transfer: 50-70% reduction [citation:3]
   • Memory usage: ~40% less than REST [citation:3]
   • Concurrent requests: Supports thousands of concurrent requests [citation:3]
`)
}

// ============================================================================
// بخش 15: Best Practices
// ============================================================================

func bestPractices() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("💡 GQLGEN BEST PRACTICES")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────┐
│ BEST PRACTICES FOR GQLGEN                                     │
├─────────────────────────────────────────────────────────────────┤
│                                                                │
│ 1. Schema Design                                               │
│    - استفاده از naming conventions یکسان                       │
│    - افزودن توضیحات برای هر type و field                      │
│    - استفاده از interface و union برای چندریختی              │
│                                                                │
│ 2. Error Handling                                              │
│    - برگرداندن errorهای معنادار                               │
│    - استفاده از extensions برای اطلاعات اضافی                 │
│    - لاگ کردن خطاهای داخلی                                     │
│                                                                │
│ 3. Security                                                    │
│    - اعتبارسنجی query depth                                   │
│    - محدود کردن introspection در production                   │
│    - استفاده از persisted queries                            │
│    - rate limiting بر اساس complexity                         │
│                                                                │
│ 4. Federation (برای میکروسرویس‌ها)                            │
│    - استفاده از @key directive                                │
│    - تقسیم Schema بین سرویس‌های مختلف                        │
│    - استفاده از Apollo Federation                            │
│                                                                │
│ 5. Testing                                                     │
│    - تست unit برای هر resolver                                │
│    - تست integration با دیتابیس واقعی                        │
│    - تست query complexity                                     │
│                                                                │
│ 6. Monitoring                                                  │
│    - metrics با Prometheus                                    │
│    - tracing با OpenTelemetry                                │
│    - logging درخواست‌ها                                       │
│                                                                │
│ 7. Versioning                                                  │
│    - GraphQL نسخه‌بندی نیاز ندارد!                           │
│    - فقط فیلدهای جدید اضافه کن                                │
│    - از @deprecated directive استفاده کن                     │
│                                                                │
│ 8. Real-time with Subscriptions                               │
│    - استفاده از WebSocket برای subscriptions                  │
│    - مدیریت اتصالات                                           │
│    - پیاده‌سازی pub/sub                                      │
│                                                                │
└─────────────────────────────────────────────────────────────────┘

📋 COMMANDS QUICK REFERENCE:

   # Initialize project
   $ go run github.com/99designs/gqlgen init

   # Generate code after schema changes
   $ go run github.com/99designs/gqlgen generate

   # Run server
   $ go run server.go

   # Update gqlgen version
   $ go get -u github.com/99designs/gqlgen

🎯 WHEN TO USE GRAPHQL:

   ✅ Complex data requirements (multiple related entities)
   ✅ Rapid client iteration (mobile apps, multiple frontends)
   ✅ Real-time applications (with subscriptions)
   ✅ Aggregating multiple data sources
   ✅ When over-fetching/under-fetching is a problem

   ❌ Simple CRUD with single entity
   ❌ Serverless with cold start concerns
   ❌ Public APIs with unknown clients (requires careful security)
`)
}

// ============================================================================
// بخش 16: جمع‌بندی
// ============================================================================

func main() {
	fmt.Println(stringsRepeat("=", 80))
	fmt.Println("🎯 COMPLETE GQLGEN GUIDE")
	fmt.Println("GraphQL Server with Go - Schema First, Type Safe, Code Generation")
	fmt.Println(stringsRepeat("=", 80))

	performanceTips()
	bestPractices()

	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📚 QUICK REFERENCE")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────┐
│ WORKFLOW SUMMARY                                              │
├─────────────────────────────────────────────────────────────────┤
│                                                                │
│  1. Define schema.graphqls                                    │
│     - Types, Queries, Mutations, Subscriptions                │
│                                                                │
│  2. Configure gqlgen.yml                                      │
│     - Models, resolvers, directives                           │
│                                                                │
│  3. Run code generation                                       │
│     $ go run github.com/99designs/gqlgen generate            │
│                                                                │
│  4. Implement resolvers                                       │
│     - Business logic in schema.resolvers.go                   │
│     - State management in resolver.go                         │
│                                                                │
│  5. Add DataLoaders                                           │
│     - Prevent N+1 problems                                    │
│     - Batch database requests                                 │
│                                                                │
│  6. Add middleware                                            │
│     - Authentication, logging, rate limiting                  │
│                                                                │
│  7. Run server                                                │
│     $ go run server.go                                        │
│                                                                │
└─────────────────────────────────────────────────────────────────┘

💡 GOLDEN RULES:

   1. Schema first: always start with schema definition
   2. Never edit generated files (resolver implementation is safe)
   3. Use custom models for complex business logic
   4. Use DataLoaders for N+1 prevention
   5. Implement query complexity limits for production
   6. Use @deprecated instead of removing fields
   7. Test resolvers in isolation
   8. Monitor query performance
   9. Use persisted queries for production APIs
   10. Keep schema.graphqls documented

🔗 USEFUL LINKS:
   • gqlgen docs: https://gqlgen.com/
   • GraphQL spec: https://graphql.org/
   • Apollo Federation: https://www.apollographql.com/docs/federation/
   • GraphQL playground: http://localhost:8080/
`)

	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🎯 GQLGEN GUIDE - COMPLETE")
	fmt.Println("Ready to build efficient GraphQL APIs with Go!")
	fmt.Println(stringsRepeat("=", 80))
}

// تابع کمکی برای تکرار رشته
func stringsRepeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}

/*
# مقداردهی اولیه پروژه
go run github.com/99designs/gqlgen init

# باز生成 کد بعد از تغییر Schema
go run github.com/99designs/gqlgen generate

# اجرای سرور
go run server.go
 */

/*
خلاصه gqlgen
مفهوم	توضیح
Schema First	ابتدا Schema GraphQL را تعریف می‌کنی، سپس کد生成 می‌شود
Code Generation	gqlgen کدهای تکراری (مثل Marshal/Unmarshal) را自动生成 می‌کند
Type Safety	هرگز با map[string]interface{} سروکار نداری - همه چیز type safe است
Resolver	منطق business logic در resolverها پیاده‌سازی می‌شود
DataLoader	برای جلوگیری از مشکل N+1 در resolverهای تو در تو
*/
/*
فایل‌های اصلی
فایل	توضیح
schema.graphqls	تعریف Schema GraphQL
resolver.go	مدیریت state و وابستگی‌ها
schema.resolvers.go	پیاده‌سازی منطق resolverها
generated.go	کد生成 شده (دستکاری نکن)
models_gen.go	مدل‌های生成 شده (می‌توان override کرد)
gqlgen.yml	فایل تنظیمات

*/