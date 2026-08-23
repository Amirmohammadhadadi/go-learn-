// ============================================================================
// FILE: soft_skills_guide.go
// TITLE: راهنمای مهارت‌های نرم برای توسعه‌دهندگان ارشد Go
// HOW TO RUN: go run soft_skills_guide.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - مهارت‌های نرم برای Senior Developer
// ============================================================================
//
// مهارت‌های فنی به تنهایی برای یک Senior Developer کافی نیست.
// مهارت‌های نرم (Soft Skills) تفاوت بین یک برنامه‌نویس خوب و یک رهبر فنی را مشخص می‌کنند.
//
// مهارت‌های کلیدی:
// 1. Code Review مؤثر
//    - نقد سازنده، یادگیری جمعی، بهبود کیفیت کد
//
// 2. مستندسازی (Documentation)
//    - API documentation با Swagger/OpenAPI
//    - README, Architecture docs, Runbooks
//
// 3. منتورینگ (Mentoring)
//    - آموزش اعضای تیم، انتقال دانش، رشد حرفه‌ای
//
// 4. برآورد زمان و تسک‌ها
//    - تخمین دقیق، مدیریت انتظارات، شکستن تسک‌ها
//
// 5. ارتباط با تیم محصول و فنی
//    - ترجمه نیازمندی‌ها، مدیریت ذی‌نفعان، تصمیم‌گیری
//
// قانون طلایی:
// "کد را نقد کن، نه برنامه‌نویس را.
//  مستندات را به اندازه کد جدی بگیر.
//  منتورینگ بهترین راه برای یادگیری خودت است.
//  تخمین‌ها همیشه عدم قطعیت دارند - شفاف باش.
//  ارتباط مؤثر 80% کار یک سنیور است."
// ============================================================================

package ______soft_skills

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/swaggo/http-swagger"
)

// برای swaggo - این کامنت‌ها در فایل واقعی استفاده می‌شوند
// در این فایل فقط توضیحات هستند

// ============================================================================
// بخش 1: Code Review مؤثر
// ============================================================================

func codeReviewGuide() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📝 CODE REVIEW BEST PRACTICES")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ WHAT TO LOOK FOR IN CODE REVIEW                                           │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  1. CORRECTNESS                                                            │
│     • Does the code do what it's supposed to do?                           │
│     • Are edge cases handled?                                              │
│     • Are errors handled properly?                                         │
│     • Is there proper validation?                                          │
│                                                                             │
│  2. READABILITY                                                            │
│     • Are names meaningful?                                                │
│     • Is the code easy to understand?                                      │
│     • Are comments useful (not obvious)?                                   │
│     • Does it follow project conventions?                                  │
│                                                                             │
│  3. MAINTAINABILITY                                                        │
│     • Is the code modular?                                                 │
│     • Are there duplicate code?                                            │
│     • Are dependencies properly managed?                                   │
│     • Is it testable?                                                      │
│                                                                             │
│  4. PERFORMANCE                                                            │
│     • Are there obvious performance issues?                                │
│     • Is there unnecessary allocation?                                     │
│     • Are there N+1 queries?                                               │
│                                                                             │
│  5. SECURITY                                                               │
│     • Is user input validated?                                             │
│     • Are there SQL injection risks?                                       │
│     • Are secrets hardcoded?                                               │
│     • Is proper authentication/authorization in place?                     │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ HOW TO GIVE FEEDBACK                                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ✅ GOOD FEEDBACK:                                                         │
│     • "This function is quite long (120 lines).                           │
│        Could we break it into smaller helper functions?"                   │
│     • "I see a potential race condition here when accessing this map.     │
│        Consider using sync.RWMutex."                                       │
│     • "Great approach! The error handling is very clean."                  │
│                                                                             │
│  ❌ BAD FEEDBACK:                                                          │
│     • "This code is terrible."                                             │
│     • "You should know better."                                            │
│     • "Rewrite everything."                                                │
│     • "This is wrong." (without explanation)                               │
│                                                                             │
│  TIPS FOR EFFECTIVE REVIEWS:                                              │
│     • Be specific and constructive                                        │
│     • Explain the "why" behind suggestions                                 │
│     • Praise good code, not just criticize                                 │
│     • Ask questions instead of making demands                             │
│     • Suggest alternatives, not just problems                              │
│     • Respect different approaches                                        │
│     • Focus on the code, not the person                                    │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ CODE REVIEW CHECKLIST FOR GO                                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ☐ Error handling: all errors are checked                                  │
│  ☐ Goroutines: have proper cleanup/context                                 │
│  ☐ Channels: closed properly                                               │
│  ☐ Mutexes: no copying, proper locking/unlocking                          │
│  ☐ Context: passed as first parameter                                      │
│  ☐ Interfaces: small and focused                                           │
│  ☐ Package naming: clear, lowercase                                        │
│  ☐ Tests: covers edge cases, uses table-driven                            │
│  ☐ Documentation: exported symbols documented                              │
│  ☐ No race conditions: go test -race passes                                │
│  ☐ No deadlocks: proper channel/mutex usage                                │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
`)
}

// ============================================================================
// بخش 2: مستندسازی با Swaggo (OpenAPI)
// ============================================================================

// User مدل کاربر برای مستندسازی
type User struct {
	ID        int       `json:"id" example:"1" description:"User ID"`
	Name      string    `json:"name" example:"Ali Rezaei" description:"Full name"`
	Email     string    `json:"email" example:"ali@example.com" description:"Email address"`
	Age       int       `json:"age" example:"30" description:"User age"`
	CreatedAt time.Time `json:"created_at" example:"2024-01-01T00:00:00Z" description:"Creation timestamp"`
}

// APIResponse پاسخ استاندارد API
type APIResponse struct {
	Success bool        `json:"success" example:"true" description:"Request success status"`
	Message string      `json:"message" example:"Operation completed" description:"Response message"`
	Data    interface{} `json:"data,omitempty" description:"Response data"`
}

// UserHandler هندلر کاربران با مستندات Swagger
type UserHandler struct {
	// می‌توانید وابستگی‌ها را اینجا اضافه کنید
}

// @Summary      Create a new user
// @Description  Creates a new user in the system
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user body User true "User object to create"
// @Success      201  {object}  APIResponse{data=User}
// @Failure      400  {object}  APIResponse
// @Failure      409  {object}  APIResponse
// @Router       /users [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	// implementation
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(APIResponse{
		Success: true,
		Message: "User created",
	})
}

// @Summary      Get user by ID
// @Description  Returns a single user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  APIResponse{data=User}
// @Failure      404  {object}  APIResponse
// @Router       /users/{id} [get]
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	// implementation
	id := chi.URLParam(r, "id")
	user := User{
		ID:        1,
		Name:      "Ali Rezaei",
		Email:     "ali@example.com",
		Age:       30,
		CreatedAt: time.Now(),
	}
	json.NewEncoder(w).Encode(APIResponse{
		Success: true,
		Data:    user,
	})
}

// @Summary      Update user
// @Description  Updates an existing user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Param        user body User true "User object to update"
// @Success      200  {object}  APIResponse{data=User}
// @Failure      400  {object}  APIResponse
// @Failure      404  {object}  APIResponse
// @Router       /users/{id} [put]
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// implementation
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(APIResponse{
		Success: true,
		Message: "User updated",
	})
}

// @Summary      Delete user
// @Description  Deletes a user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      204  {object}  nil
// @Failure      404  {object}  APIResponse
// @Router       /users/{id} [delete]
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// implementation
	w.WriteHeader(http.StatusNoContent)
}

// @Summary      List users
// @Description  Returns a list of users with pagination
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        page     query     int  false  "Page number"  default(1)  minimum(1)
// @Param        pageSize query     int  false  "Items per page"  default(10)  minimum(1)  maximum(100)
// @Param        search   query     string false "Search term"
// @Success      200  {object}  APIResponse{data=[]User}
// @Router       /users [get]
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	// implementation
	users := []User{
		{ID: 1, Name: "Ali", Email: "ali@example.com", Age: 30},
		{ID: 2, Name: "Sara", Email: "sara@example.com", Age: 25},
	}
	json.NewEncoder(w).Encode(APIResponse{
		Success: true,
		Data:    users,
	})
}

// HealthCheck health check endpoint
// @Summary      Health check
// @Description  Returns the health status of the service
// @Tags         system
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /health [get]
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
		"time":   time.Now().String(),
	})
}

// setupSwagger راه‌اندازی Swagger UI
func setupSwagger() {
	// در فایل main.go:
	// r := chi.NewRouter()
	//
	// // Swagger documentation
	// r.Get("/swagger/*", httpSwagger.Handler(
	//     httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	// ))
	//
	// // یا با Swaggo CLI:
	// // $ swag init
	// // $ go run main.go
}

// Documentation best practices
func documentationBestPractices() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📚 DOCUMENTATION BEST PRACTICES")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ SWAGGO/SWAGGER SETUP                                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  # Install swaggo                                                          │
│  $ go get -u github.com/swaggo/swag/cmd/swag                               │
│  $ go get -u github.com/swaggo/http-swagger                                │
│                                                                             │
│  # Generate docs                                                           │
│  $ swag init                                                               │
│                                                                             │
│  # Run with live reload                                                    │
│  $ swag init --parseDependency --parseInternal                             │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ MAIN DOCUMENTATION TYPES                                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  1. API Documentation (Swagger/OpenAPI)                                    │
│     • Endpoints, parameters, responses                                      │
│     • Authentication requirements                                          │
│     • Rate limits                                                          │
│                                                                             │
│  2. README.md                                                              │
│     • Project overview                                                     │
│     • Setup instructions                                                   │
│     • Dependencies                                                         │
│     • Environment variables                                                │
│     • How to run                                                           │
│     • How to test                                                          │
│                                                                             │
│  3. Architecture Documentation                                             │
│     • System design                                                        │
│     • Component diagrams                                                   │
│     • Data flow                                                            │
│     • Technology choices                                                   │
│                                                                             │
│  4. Runbooks                                                               │
│     • Deployment procedures                                                │
│     • Backup/restore                                                       │
│     • Incident response                                                    │
│     • Monitoring and alerts                                                │
│                                                                             │
│  5. Code Documentation                                                     │
│     • Package comments                                                     │
│     • Exported types/functions                                             │
│     • Complex logic explanations                                           │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
`)
}

// ============================================================================
// بخش 3: منتورینگ (Mentoring)
// ============================================================================

func mentoringGuide() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🎓 MENTORING BEST PRACTICES")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ EFFECTIVE MENTORING STRATEGIES                                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  1. SET EXPECTATIONS                                                       │
│     • Define goals and timeline                                            │
│     • Establish regular meeting schedule                                   │
│     • Clarify availability and boundaries                                  │
│     • Agree on communication channels                                      │
│                                                                             │
│  2. ASK QUESTIONS, DON'T JUST ANSWER                                       │
│     • "What have you tried?"                                               │
│     • "What do you think the issue might be?"                              │
│     • "How would you approach this?"                                       │
│     • Guide, don't dictate                                                 │
│                                                                             │
│  3. PROVIDE ACTIONABLE FEEDBACK                                            │
│     • Be specific and timely                                               │
│     • Focus on behavior, not personality                                   │
│     • Balance positive and constructive feedback                           │
│     • Follow up on progress                                                │
│                                                                             │
│  4. SHARE CONTEXT, NOT JUST CODE                                           │
│     • Explain why decisions were made                                      │
│     • Share past mistakes and learnings                                    │
│     • Discuss trade-offs                                                   │
│     • Connect technical decisions to business outcomes                     │
│                                                                             │
│  5. ENCOURAGE INDEPENDENCE                                                 │
│     • Gradually reduce hands-on help                                       │
│     • Assign challenging tasks                                             │
│     • Let them make mistakes (safely)                                      │
│     • Celebrate growth and achievements                                    │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ MENTORING TOPICS FOR GO DEVELOPERS                                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Beginner → Intermediate                                                   │
│  • Concurrency (goroutines, channels)                                      │
│  • Error handling patterns                                                 │
│  • Testing (table-driven, mocking)                                         │
│  • Context usage                                                           │
│  • Package organization                                                    │
│                                                                             │
│  Intermediate → Senior                                                     │
│  • Performance optimization                                                │
│  • Design patterns in Go                                                   │
│  • Architecture decisions                                                  │
│  • Production debugging                                                    │
│  • Code review leadership                                                  │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ MENTORING CHECKLIST                                                       │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ☐ Have regular 1-on-1 meetings (weekly)                                   │
│  ☐ Set clear, measurable goals                                             │
│  ☐ Review code together                                                    │
│  ☐ Pair program on complex tasks                                           │
│  ☐ Share resources (books, articles, talks)                                │
│  ☐ Introduce to the wider team                                             │
│  ☐ Give credit for their work                                              │
│  ☐ Advocate for their growth                                               │
│  ☐ Ask for feedback on your mentoring                                      │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
`)
}

// ============================================================================
// بخش 4: برآورد زمان و تسک‌ها (Estimation)
// ============================================================================

func estimationGuide() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("⏱️ TASK ESTIMATION BEST PRACTICES")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ ESTIMATION TECHNIQUES                                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  1. BREAK DOWN TASKS                                                       │
│     • Split into smaller chunks (< 8 hours)                                │
│     • Identify dependencies                                                │
│     • Separate known vs unknown                                            │
│                                                                             │
│  2. USE RELATIVE ESTIMATION                                                │
│     • Fibonacci sequence (1, 2, 3, 5, 8, 13, 21)                          │
│     • T-shirt sizing (XS, S, M, L, XL)                                     │
│     • Planning poker (team estimation)                                     │
│                                                                             │
│  3. ACCOUNT FOR OVERHEAD                                                   │
│     • Meetings (15-20%)                                                    │
│     • Code review (10-15%)                                                 │
│     • Documentation (5-10%)                                                │
│     • Testing (20-30%)                                                     │
│     • Bug fixes (10-15%)                                                   │
│                                                                             │
│  4. ADD BUFFERS                                                            │
│     • Multiply by 1.5-2x for uncertainty                                   │
│     • Add buffer for unexpected issues                                     │
│     • Communicate range, not single number                                 │
│                                                                             │
│  5. TRACK AND LEARN                                                        │
│     • Compare actual vs estimate                                           │
│     • Identify patterns of over/under estimation                           │
│     • Adjust future estimates accordingly                                  │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ COMMON ESTIMATION PITFALLS                                                │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ❌ HOFTAD ESTIMATION                                                       │
│     • Giving an estimate immediately                                       │
│     • "This looks easy, probably 2 hours"                                  │
│     ✅ Take time to understand first                                       │
│                                                                             │
│  ❌ OPTIMISM BIAS                                                           │
│     • Assuming everything will go perfectly                                │
│     • "There should be no issues"                                          │
│     ✅ Always add buffer for unknowns                                      │
│                                                                             │
│  ❌ PRESSURE-BASED ESTIMATION                                               │
│     • "We need this done by Friday"                                        │
│     • Estimating to please stakeholders                                    │
│     ✅ Be honest about what's realistic                                    │
│                                                                             │
│  ❌ COMPARISON FALLACY                                                      │
│     • "We did something similar in 3 days"                                 │
│     • Every task is unique                                                 │
│     ✅ Use historical data as reference, not absolute                        │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ TASK BREAKDOWN EXAMPLE                                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Task: Implement user authentication with JWT                             │
│                                                                             │
│  Sub-tasks:                                                               │
│  1. Research JWT libraries (1-2 hours)                                    │
│  2. Setup database tables (1 hour)                                        │
│  3. Implement registration endpoint (2-3 hours)                           │
│  4. Implement login endpoint (2-3 hours)                                  │
│  5. Implement token refresh (2 hours)                                     │
│  6. Add middleware for protected routes (2 hours)                         │
│  7. Write unit tests (3-4 hours)                                          │
│  8. Write integration tests (2-3 hours)                                   │
│  9. Documentation (1 hour)                                                 │
│  10. Code review and fixes (2 hours)                                      │
│                                                                             │
│  Total: 18-24 hours (Development)                                         │
│  With buffers: 27-36 hours (~3.5-4.5 days)                                │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
`)
}

// ============================================================================
// بخش 5: ارتباط با تیم محصول و فنی
// ============================================================================

func communicationGuide() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🗣️ TEAM COMMUNICATION BEST PRACTICES")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ COMMUNICATION WITH PRODUCT TEAM                                           │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  1. TRANSLATE REQUIREMENTS                                                 │
│     • Understand business goals, not just features                         │
│     • Ask "why" to clarify intent                                          │
│     • Break down into technical tasks                                      │
│     • Identify hidden complexity                                           │
│                                                                             │
│  2. MANAGE EXPECTATIONS                                                    │
│     • Be transparent about trade-offs                                      │
│     • Explain technical constraints                                        │
│     • Propose alternatives                                                 │
│     • Set realistic deadlines                                              │
│                                                                             │
│  3. PROVIDE PROGRESS UPDATES                                               │
│     • Regular status updates (daily/weekly)                                │
│     • Flag blockers early                                                  │
│     • Celebrate wins                                                       │
│     • Be honest about delays                                               │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ TECHNICAL COMMUNICATION                                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  1. EFFECTIVE MEETINGS                                                     │
│     • Have clear agenda                                                    │
│     • Invite only necessary people                                         │
│     • Start and end on time                                                │
│     • Document decisions                                                   │
│     • Assign action items                                                  │
│                                                                             │
│  2. ASYNCHRONOUS COMMUNICATION                                              │
│     • Use appropriate channels (Slack, email, Jira)                       │
│     • Write clear, concise messages                                        │
│     • Use threads to keep context                                          │
│     • Tag relevant people                                                  │
│                                                                             │
│  3. CONFLICT RESOLUTION                                                    │
│     • Focus on the problem, not the person                                 │
│     • Listen actively                                                      │
│     • Find common ground                                                   │
│     • Escalate when necessary                                              │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ COMMUNICATION TEMPLATES                                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  STATUS UPDATE:                                                           │
│  • What I accomplished yesterday                                           │
│  • What I plan to do today                                                 │
│  • Blockers (if any)                                                       │
│                                                                             │
│  INCIDENT REPORT:                                                          │
│  • What happened?                                                          │
│  • Why did it happen?                                                      │
│  • What was the impact?                                                    │
│  • What is the fix?                                                        │
│  • How to prevent recurrence?                                              │
│                                                                             │
│  TECH DEBT PROPOSAL:                                                       │
│  • Current situation                                                       │
│  • Proposed solution                                                       │
│  • Benefits                                                                │
│  • Effort estimate                                                         │
│  • Risks                                                                   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
`)
}

// ============================================================================
// بخش 6: Senior Developer Responsibilities
// ============================================================================

func seniorResponsibilities() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("👔 SENIOR DEVELOPER RESPONSIBILITIES")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ TECHNICAL RESPONSIBILITIES                                                │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  • Architecture decisions                                                  │
│  • Code review leadership                                                  │
│  • Technical debt management                                               │
│  • Performance optimization                                                │
│  • Security best practices                                                 │
│  • Production incident response                                            │
│  • Documentation standards                                                 │
│  • Tool and library evaluation                                             │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ TEAM RESPONSIBILITIES                                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  • Mentoring junior developers                                             │
│  • Onboarding new team members                                             │
│  • Facilitating technical discussions                                      │
│  • Improving team processes                                                │
│  • Knowledge sharing (lunch & learns, docs)                                │
│  • Interviewing candidates                                                 │
│  • Running retrospectives                                                  │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ leadership RESPONSIBILITIES                                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  • Aligning technical decisions with business goals                        │
│  • Communicating with stakeholders                                         │
│  • Project estimation and planning                                         │
│  • Risk identification and mitigation                                      │
│  • Driving technical initiatives                                           │
│  • Building consensus                                                      │
│  • Representing the engineering team                                       │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
`)
}

// ============================================================================
// بخش 7: Career Growth Plan
// ============================================================================

func careerGrowthPlan() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📈 SENIOR DEVELOPER CAREER GROWTH")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ 6-MONTH GROWTH PLAN                                                       │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  MONTHS 1-2:                                                              │
│  • Master code review skills                                               │
│  • Document at least 3 internal systems                                    │
│  • Mentor one junior developer                                             │
│  • Improve estimation accuracy                                             │
│                                                                             │
│  MONTHS 3-4:                                                              │
│  • Lead a medium-sized feature                                             │
│  • Create technical design docs                                            │
│  • Interview 3+ candidates                                                 │
│  • Present at team meeting                                                 │
│                                                                             │
│  MONTHS 5-6:                                                              │
│  • Propose and implement process improvement                               │
│  • Write a blog post or internal wiki page                                 │
│  • Take ownership of a service                                             │
│  • Run incident post-mortem                                                │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ SKILLS TO DEVELOP                                                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  TECHNICAL:                                                               │
│  • Advanced Go patterns                                                    │
│  • System design                                                           │
│  • Performance optimization                                                │
│  • Security best practices                                                 │
│                                                                             │
│  SOFT SKILLS:                                                             │
│  • Public speaking                                                         │
│  • Writing (technical docs, proposals)                                    │
│  • Conflict resolution                                                     │
│  • Facilitation                                                            │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
`)
}

// ============================================================================
// بخش 8: جمع‌بندی
// ============================================================================

func summary() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📚 SUMMARY: PATH TO SENIOR DEVELOPER")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ KEY TAKEAWAYS                                                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  1. CODE REVIEW                                                            │
│     • Be constructive and specific                                         │
│     • Focus on the code, not the person                                    │
│     • Learn from reviewing code                                            │
│                                                                             │
│  2. DOCUMENTATION                                                          │
│     • Use Swagger for API docs                                             │
│     • Document as you code                                                 │
│     • Keep docs up to date                                                 │
│                                                                             │
│  3. MENTORING                                                              │
│     • Teaching is the best way to learn                                    │
│     • Set clear expectations                                               │
│     • Celebrate progress                                                   │
│                                                                             │
│  4. ESTIMATION                                                             │
│     • Break down tasks                                                     │
│     • Add buffers for unknowns                                             │
│     • Learn from past estimates                                            │
│                                                                             │
│  5. COMMUNICATION                                                          │
│     • Be clear and concise                                                 │
│     • Listen actively                                                      │
│     • Build relationships                                                  │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

💡 GOLDEN RULES:

   1. Always review code with empathy
   2. Document everything that's not obvious
   3. Mentor others to grow yourself
   4. Estimate ranges, not points
   5. Communicate clearly and often
   6. Lead by example
   7. Take ownership of outcomes
   8. Be curious and always learning
   9. Help others succeed
   10. Balance technical excellence with business value
`)
}

// ============================================================================
// بخش 9: Main
// ============================================================================

func main() {
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("🎯 SOFT SKILLS FOR SENIOR GO DEVELOPERS")
	fmt.Println("Code Review | Documentation | Mentoring | Estimation | Communication")
	fmt.Println(strings.Repeat("=", 80))

	codeReviewGuide()
	documentationBestPractices()
	mentoringGuide()
	estimationGuide()
	communicationGuide()
	seniorResponsibilities()
	careerGrowthPlan()
	summary()

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🎯 SOFT SKILLS - COMPLETE")
	fmt.Println("Ready to level up your career!")
	fmt.Println(strings.Repeat("=", 80))
}

// تابع کمکی
func stringsRepeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}
