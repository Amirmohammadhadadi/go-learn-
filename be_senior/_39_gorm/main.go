// ============================================================================
// FILE: gorm_guide.go
// TITLE: راهنمای کامل GORM در Go - ORM مدرن و قدرتمند
// HOW TO RUN: go run gorm_guide.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - GORM چیست و چه مزایایی دارد؟
// ============================================================================
//
// GORM یک ORM (Object Relational Mapping) کامل و قدرتمند برای Go است.
//
// مزایای استفاده از GORM:
// 1. کاهش کد تکراری (CRUD ساده)
// 2. Auto-migration (ایجاد خودکار جدول‌ها)
// 3. ارتباطات پیشرفته (Associations: BelongsTo, HasOne, HasMany, ManyToMany)
// 4. Callbacks و Hooks (BeforeCreate, AfterUpdate, etc.)
// 5. Scopes (بازیافت queryهای رایج)
// 6. Soft Delete (حذف منطقی)
// 7. Transactionها
// 8. Preloading (Eager Loading برای جلوگیری از N+1)
// 9. Raw SQL (زمانی که نیاز به کنترل کامل دارید)
//
// چه زمانی از GORM استفاده کنیم؟
// ✅ پروژه‌های با CRUD ساده
// ✅ زمان محدود (راه‌اندازی سریع)
// ✅ تیمی که با ORM راحت‌تر است
//
// چه زمانی از Raw SQL استفاده کنیم؟
// ❌ کوئری‌های بسیار پیچیده
// ❌ نیاز به حداکثر performance
// ❌ کنترل کامل روی کوئری
//
// قانون طلایی:
// "برای CRUD معمولی از GORM استفاده کن.
//  برای کوئری‌های پیچیده و performance-critical از Raw SQL استفاده کن.
//  همیشه Associations را با Preload بارگذاری کن تا از N+1 جلوگیری شود."
// ============================================================================

package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ============================================================================
// بخش 1: نصب و راه‌اندازی
// ============================================================================

/*
نصب:

# GORM core
$ go get -u gorm.io/gorm

# درایور PostgreSQL
$ go get -u gorm.io/driver/postgres

# درایور MySQL (در صورت نیاز)
$ go get -u gorm.io/driver/mysql

# درایور SQLite (در صورت نیاز)
$ go get -u gorm.io/driver/sqlite
*/

// ============================================================================
// بخش 2: مدل‌های GORM با تگ‌ها و Associations
// ============================================================================

// User مدل کاربر (GORM Model)
type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Name     string `gorm:"size:100;not null" json:"name"`
	Email    string `gorm:"size:255;uniqueIndex;not null" json:"email"`
	Age      int    `gorm:"default:0;check:age >= 0 AND age <= 150" json:"age"`
	IsActive bool   `gorm:"default:true" json:"is_active"`
	Password string `gorm:"-" json:"-"` // ignore در GORM
	Role     string `gorm:"default:user;index" json:"role"`

	// Associations
	Profile   Profile    `gorm:"foreignKey:UserID" json:"profile,omitempty"`           // One-to-One
	Posts     []Post     `gorm:"foreignKey:UserID" json:"posts,omitempty"`             // One-to-Many
	Orders    []Order    `gorm:"foreignKey:UserID" json:"orders,omitempty"`            // One-to-Many
	Companies []Company  `gorm:"many2many:user_companies;" json:"companies,omitempty"` // Many-to-Many
	Languages []Language `gorm:"many2many:user_languages;" json:"languages,omitempty"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // Soft Delete
}

// Profile مدل پروفایل (One-to-One با User)
type Profile struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `gorm:"uniqueIndex;not null"` // foreign key
	Bio       string `gorm:"type:text;size:500"`
	Avatar    string `gorm:"size:255"`
	Phone     string `gorm:"size:20"`
	Address   string `gorm:"size:255"`
	City      string `gorm:"size:100"`
	Country   string `gorm:"size:100"`
	ZipCode   string `gorm:"size:20"`
	BirthDate time.Time
}

// Post مدل پست (One-to-Many با User)
type Post struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"index;not null"` // foreign key
	Title     string    `gorm:"size:200;not null;index"`
	Content   string    `gorm:"type:text;not null"`
	Published bool      `gorm:"default:false;index"`
	ViewCount int       `gorm:"default:0"`
	Tags      []Tag     `gorm:"many2many:post_tags;" json:"tags,omitempty"` // Many-to-Many
	Comments  []Comment `gorm:"foreignKey:PostID" json:"comments,omitempty"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Comment مدل کامنت (One-to-Many با Post)
type Comment struct {
	ID        uint   `gorm:"primaryKey"`
	PostID    uint   `gorm:"index;not null"`
	UserID    uint   `gorm:"index;not null"`
	Content   string `gorm:"type:text;not null"`
	Approved  bool   `gorm:"default:false;index"`
	CreatedAt time.Time
}

// Tag مدل تگ (Many-to-Many با Post)
type Tag struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"size:50;uniqueIndex;not null"`
}

// Category مدل دسته‌بندی
type Category struct {
	ID       uint       `gorm:"primaryKey"`
	Name     string     `gorm:"size:100;uniqueIndex;not null"`
	Slug     string     `gorm:"size:100;uniqueIndex;not null"`
	ParentID *uint      `gorm:"index"` // Self-referential
	Parent   *Category  `gorm:"foreignKey:ParentID"`
	Children []Category `gorm:"foreignKey:ParentID"`
}

// Order مدل سفارش
type Order struct {
	ID          uint        `gorm:"primaryKey"`
	UserID      uint        `gorm:"index;not null"`
	OrderNumber string      `gorm:"size:50;uniqueIndex;not null"`
	Total       float64     `gorm:"type:decimal(10,2);not null"`
	Status      string      `gorm:"size:50;default:pending;index"`
	Items       []OrderItem `gorm:"foreignKey:OrderID"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// OrderItem مدل آیتم سفارش
type OrderItem struct {
	ID        uint    `gorm:"primaryKey"`
	OrderID   uint    `gorm:"index;not null"`
	ProductID uint    `gorm:"index;not null"`
	Quantity  int     `gorm:"not null;check:quantity > 0"`
	Price     float64 `gorm:"type:decimal(10,2);not null"`
	Product   Product `gorm:"foreignKey:ProductID"`
}

// Product مدل محصول
type Product struct {
	ID          uint    `gorm:"primaryKey"`
	Name        string  `gorm:"size:200;not null;index"`
	Description string  `gorm:"type:text"`
	Price       float64 `gorm:"type:decimal(10,2);not null;index"`
	Stock       int     `gorm:"not null;default:0"`
	CategoryID  uint    `gorm:"index"`
	Category    Category
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Company مدل شرکت (Many-to-Many با User)
type Company struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"size:100;uniqueIndex;not null"`
	Code string `gorm:"size:20;uniqueIndex;not null"`
}

// Language مدل زبان (Many-to-Many با User)
type Language struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"size:50;uniqueIndex;not null"`
	Code string `gorm:"size:10;uniqueIndex;not null"`
}

// ============================================================================
// بخش 3: Connection و تنظیمات
// ============================================================================

// Config تنظیمات اتصال
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// NewGormConnection ایجاد اتصال GORM
func NewGormConnection(config *DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Tehran",
		config.Host, config.User, config.Password, config.DBName, config.Port, config.SSLMode,
	)

	// تنظیمات logger
	newLogger := logger.Default.LogMode(logger.Info)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                                   newLogger,
		SkipDefaultTransaction:                   true, // بهبود performance
		PrepareStmt:                              true, // prepared statements
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	// گرفتن connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// تنظیمات connection pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

// ============================================================================
// بخش 4: Auto Migration
// ============================================================================

func runMigration(db *gorm.DB) error {
	// Auto migrate
	err := db.AutoMigrate(
		&User{},
		&Profile{},
		&Post{},
		&Comment{},
		&Tag{},
		&Category{},
		&Order{},
		&OrderItem{},
		&Product{},
		&Company{},
		&Language{},
	)
	if err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	// افزودن constraintهای سفارشی
	// db.Exec("CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_name ON users(name)")

	log.Println("Migration completed successfully")
	return nil
}

// ============================================================================
// بخش 5: Basic CRUD Operations
// ============================================================================

// UserService سرویس کاربران
type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

// Create ایجاد کاربر جدید
func (s *UserService) Create(user *User) error {
	result := s.db.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// CreateBatch ایجاد دسته‌ای کاربران
func (s *UserService) CreateBatch(users []User) error {
	return s.db.CreateInBatches(users, 100).Error
}

// GetByID دریافت کاربر با ID
func (s *UserService) GetByID(id uint) (*User, error) {
	var user User
	result := s.db.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// GetByEmail دریافت کاربر با ایمیل
func (s *UserService) GetByEmail(email string) (*User, error) {
	var user User
	result := s.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// GetAll دریافت همه کاربران با pagination
func (s *UserService) GetAll(page, pageSize int) ([]User, int64, error) {
	var users []User
	var total int64

	offset := (page - 1) * pageSize

	result := s.db.Model(&User{}).Count(&total)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	result = s.db.Offset(offset).Limit(pageSize).Find(&users)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return users, total, nil
}

// Update به‌روزرسانی کاربر
func (s *UserService) Update(user *User) error {
	// روش 1: Save (همه فیلدها را ذخیره می‌کند)
	// return s.db.Save(user).Error

	// روش 2: Update (فقط فیلدهای مشخص)
	return s.db.Model(user).Updates(map[string]interface{}{
		"name":      user.Name,
		"email":     user.Email,
		"age":       user.Age,
		"is_active": user.IsActive,
	}).Error
}

// UpdatePartial به‌روزرسانی جزئی
func (s *UserService) UpdatePartial(id uint, updates map[string]interface{}) error {
	return s.db.Model(&User{}).Where("id = ?", id).Updates(updates).Error
}

// Delete حذف کاربر (Soft Delete)
func (s *UserService) Delete(id uint) error {
	return s.db.Delete(&User{}, id).Error
}

// DeletePermanent حذف فیزیکی
func (s *UserService) DeletePermanent(id uint) error {
	return s.db.Unscoped().Delete(&User{}, id).Error
}

// ============================================================================
// بخش 6: Advanced Queries (شرط‌ها، Joins، گروه‌بندی)
// ============================================================================

// FindActiveUsers یافتن کاربران فعال
func (s *UserService) FindActiveUsers() ([]User, error) {
	var users []User
	result := s.db.Where("is_active = ?", true).Find(&users)
	return users, result.Error
}

// FindUsersByRole یافتن کاربران با نقش مشخص
func (s *UserService) FindUsersByRole(role string) ([]User, error) {
	var users []User
	result := s.db.Where("role = ?", role).Find(&users)
	return users, result.Error
}

// SearchUsers جستجوی کاربران
func (s *UserService) SearchUsers(keyword string) ([]User, error) {
	var users []User
	query := s.db.Where("name ILIKE ?", "%"+keyword+"%").
		Or("email ILIKE ?", "%"+keyword+"%")
	result := query.Find(&users)
	return users, result.Error
}

// FindUsersWithConditions شرط‌های پیچیده
func (s *UserService) FindUsersWithConditions(minAge, maxAge int, isActive bool, role string) ([]User, error) {
	var users []User
	result := s.db.Where("age BETWEEN ? AND ?", minAge, maxAge).
		Where("is_active = ?", isActive).
		Where("role = ?", role).
		Find(&users)
	return users, result.Error
}

// GetUserStats آمار کاربران
func (s *UserService) GetUserStats() (map[string]interface{}, error) {
	var stats struct {
		Total    int64
		Active   int64
		Inactive int64
		AvgAge   float64
		MaxAge   int
		MinAge   int
	}

	db := s.db.Model(&User{})

	db.Count(&stats.Total)
	db.Where("is_active = ?", true).Count(&stats.Active)
	db.Where("is_active = ?", false).Count(&stats.Inactive)
	db.Select("AVG(age)").Row().Scan(&stats.AvgAge)
	db.Select("MAX(age)").Row().Scan(&stats.MaxAge)
	db.Select("MIN(age)").Row().Scan(&stats.MinAge)

	return map[string]interface{}{
		"total":    stats.Total,
		"active":   stats.Active,
		"inactive": stats.Inactive,
		"avg_age":  stats.AvgAge,
		"max_age":  stats.MaxAge,
		"min_age":  stats.MinAge,
	}, nil
}

// GetUsersByAgeGroup گروه‌بندی سنی
func (s *UserService) GetUsersByAgeGroup() ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	result := s.db.Model(&User{}).
		Select("CASE " +
			"WHEN age < 18 THEN 'Under 18' " +
			"WHEN age BETWEEN 18 AND 30 THEN '18-30' " +
			"WHEN age BETWEEN 31 AND 50 THEN '31-50' " +
			"ELSE '50+' END as age_group, COUNT(*) as count").
		Group("age_group").
		Order("age_group").
		Find(&results)

	return results, result.Error
}

// ============================================================================
// بخش 7: Associations (One-to-One, One-to-Many, Many-to-Many)
// ============================================================================

// CreateUserWithProfile ایجاد کاربر همراه با پروفایل
func (s *UserService) CreateUserWithProfile(user *User, profile *Profile) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		profile.UserID = user.ID
		if err := tx.Create(profile).Error; err != nil {
			return err
		}
		return nil
	})
}

// GetUserWithProfile دریافت کاربر با پروفایل (Preload)
func (s *UserService) GetUserWithProfile(id uint) (*User, error) {
	var user User
	result := s.db.Preload("Profile").First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// GetUserWithAllAssociations دریافت کاربر با همه ارتباطات
func (s *UserService) GetUserWithAllAssociations(id uint) (*User, error) {
	var user User
	result := s.db.
		Preload("Profile").
		Preload("Posts").
		Preload("Posts.Comments").
		Preload("Posts.Tags").
		Preload("Orders").
		Preload("Orders.Items").
		Preload("Orders.Items.Product").
		Preload("Companies").
		Preload("Languages").
		First(&user, id)

	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// AddUserToCompany اضافه کردن کاربر به شرکت (Many-to-Many)
func (s *UserService) AddUserToCompany(userID, companyID uint) error {
	return s.db.Model(&User{ID: userID}).Association("Companies").Append(&Company{ID: companyID})
}

// GetUserCompanies دریافت شرکت‌های کاربر
func (s *UserService) GetUserCompanies(userID uint) ([]Company, error) {
	var user User
	var companies []Company

	result := s.db.First(&user, userID)
	if result.Error != nil {
		return nil, result.Error
	}

	result = s.db.Model(&user).Association("Companies").Find(&companies)
	return companies, result.Error
}

// ============================================================================
// بخش 8: Scopes (بازیافت کوئری‌ها)
// ============================================================================

// ActiveUserScope scope برای کاربران فعال
func ActiveUserScope(db *gorm.DB) *gorm.DB {
	return db.Where("is_active = ?", true)
}

// AdultUserScope scope برای کاربران بالای 18 سال
func AdultUserScope(db *gorm.DB) *gorm.DB {
	return db.Where("age >= ?", 18)
}

// RoleScope scope برای نقش خاص
func RoleScope(role string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("role = ?", role)
	}
}

// GetActiveAdultUsers دریافت کاربران فعال بالای 18 سال
func (s *UserService) GetActiveAdultUsers() ([]User, error) {
	var users []User
	result := s.db.Scopes(ActiveUserScope, AdultUserScope).Find(&users)
	return users, result.Error
}

// GetAdminUsers دریافت کاربران ادمین
func (s *UserService) GetAdminUsers() ([]User, error) {
	var users []User
	result := s.db.Scopes(RoleScope("admin")).Find(&users)
	return users, result.Error
}

// ============================================================================
// بخش 9: Transactions
// ============================================================================

// CreateOrderWithItems ایجاد سفارش با آیتم‌ها (با تراکنش)
func CreateOrderWithItems(db *gorm.DB, order *Order, items []OrderItem) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// ایجاد سفارش
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		// ایجاد آیتم‌ها
		for i := range items {
			items[i].OrderID = order.ID
			if err := tx.Create(&items[i]).Error; err != nil {
				return err
			}

			// کاهش موجودی
			if err := tx.Model(&Product{}).Where("id = ?", items[i].ProductID).
				Update("stock", gorm.Expr("stock - ?", items[i].Quantity)).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// TransferUserData انتقال داده‌های کاربر (مثال پیشرفته)
func TransferUserData(db *gorm.DB, fromUserID, toUserID uint) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// به‌روزرسانی پست‌ها
		if err := tx.Model(&Post{}).Where("user_id = ?", fromUserID).
			Update("user_id", toUserID).Error; err != nil {
			return err
		}

		// به‌روزرسانی سفارشات
		if err := tx.Model(&Order{}).Where("user_id = ?", fromUserID).
			Update("user_id", toUserID).Error; err != nil {
			return err
		}

		// حذف کاربر قدیمی
		if err := tx.Delete(&User{}, fromUserID).Error; err != nil {
			return err
		}

		return nil
	})
}

// ============================================================================
// بخش 10: Raw SQL و SQL Builder
// ============================================================================

// GetUsersWithRawSQL کوئری با SQL خام
func GetUsersWithRawSQL(db *gorm.DB) ([]User, error) {
	var users []User

	// روش 1: Raw
	result := db.Raw("SELECT * FROM users WHERE age > ? AND is_active = ?", 18, true).Scan(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}

// ExecuteRawUpdate اجرای آپدیت با SQL خام
func ExecuteRawUpdate(db *gorm.DB) error {
	return db.Exec("UPDATE users SET is_active = ? WHERE age < ?", false, 18).Error
}

// UsingSQLBuilder استفاده از SQL Builder
func UsingSQLBuilder(db *gorm.DB) ([]User, error) {
	var users []User

	result := db.Where("age > ?", 18).
		Where("is_active = ?", true).
		Order("created_at DESC").
		Limit(10).
		Find(&users)

	return users, result.Error
}

// ============================================================================
// بخش 11: Hooks و Callbacks
// ============================================================================

// BeforeCreate hook قبل از ایجاد
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}

	// تنظیم مقدار پیش‌فرض
	if u.Role == "" {
		u.Role = "user"
	}

	return nil
}

// AfterCreate hook بعد از ایجاد
func (u *User) AfterCreate(tx *gorm.DB) error {
	// لاگ کردن یا ارسال نوتیفیکیشن
	log.Printf("User created: %s (%s)", u.Name, u.Email)
	return nil
}

// BeforeUpdate hook قبل از به‌روزرسانی
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	// اعتبارسنجی
	if u.Age < 0 || u.Age > 150 {
		return fmt.Errorf("invalid age")
	}
	return nil
}

// AfterFind hook بعد از پیدا کردن
func (u *User) AfterFind(tx *gorm.DB) error {
	// تبدیل داده‌ها یا لاگ
	return nil
}

// BeforeDelete hook قبل از حذف
func (u *User) BeforeDelete(tx *gorm.DB) error {
	// جلوگیری از حذف کاربران ادمین
	if u.Role == "admin" {
		return fmt.Errorf("cannot delete admin user")
	}
	return nil
}

// ============================================================================
// بخش 12: Soft Delete و بازیابی
// ============================================================================

// RestoreDeletedUser بازیابی کاربر حذف شده
func RestoreDeletedUser(db *gorm.DB, id uint) error {
	return db.Unscoped().Model(&User{}).Where("id = ?", id).Update("deleted_at", nil).Error
}

// GetDeletedUsers دریافت کاربران حذف شده
func GetDeletedUsers(db *gorm.DB) ([]User, error) {
	var users []User
	result := db.Unscoped().Where("deleted_at IS NOT NULL").Find(&users)
	return users, result.Error
}

// ============================================================================
// بخش 13: Best Practices
// ============================================================================

func bestPractices() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("💡 GORM BEST PRACTICES")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────┐
│ PERFORMANCE OPTIMIZATIONS                                     │
├─────────────────────────────────────────────────────────────────┤
│ • Use Select to specify fields                                 │
│   db.Select("name", "email").Find(&users)                     │
│                                                                │
│ • Use Pluck for single column                                  │
│   db.Model(&User{}).Pluck("name", &names)                     │
│                                                                │
│ • Use FindInBatches for large datasets                         │
│   db.FindInBatches(&users, 1000, func(tx *gorm.DB, batch int) error│
│                                                                │
│ • Use Prepared Statements                                      │
│   db.PrepareStmt = true                                        │
│                                                                │
│ • Disable Default Transaction for read operations              │
│   db.SkipDefaultTransaction = true                            │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ ASSOCIATIONS                                                   │
├─────────────────────────────────────────────────────────────────┤
│ • Use Preload to avoid N+1 problem                            │
│   db.Preload("Profile").Find(&users)                          │
│                                                                │
│ • Use Joins for better performance (when possible)             │
│   db.Joins("Profile").Find(&users)                            │
│                                                                │
│ • Use Association mode for CRUD on relations                   │
│   db.Model(&user).Association("Posts").Find(&posts)           │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ ERROR HANDLING                                                │
├─────────────────────────────────────────────────────────────────┤
│ • Always check result.Error                                   │
│ • Use errors.Is for specific errors                           │
│   errors.Is(result.Error, gorm.ErrRecordNotFound)             │
│ • Use Transaction for multiple operations                     │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ MIGRATION                                                      │
├─────────────────────────────────────────────────────────────────┤
│ • Use AutoMigrate for development only                        │
│ • Use proper migration tools for production                   │
│ • Add indexes for frequently queried fields                   │
│ • Use foreign keys for data integrity                         │
└─────────────────────────────────────────────────────────────────┘
`)
}

// ============================================================================
// بخش 14: Main
// ============================================================================

func main() {
	fmt.Println(stringsRepeat("=", 80))
	fmt.Println("🎯 COMPLETE GORM GUIDE")
	fmt.Println("ORM for Go - CRUD, Associations, Transactions, Hooks")
	fmt.Println(stringsRepeat("=", 80))

	bestPractices()

	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📝 EXAMPLE USAGE")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println(`
// مثال کامل استفاده از GORM
func main() {
    // اتصال به دیتابیس
    config := &DBConfig{
        Host:     "localhost",
        Port:     "5432",
        User:     "postgres",
        Password: "password",
        DBName:   "testdb",
        SSLMode:  "disable",
    }
    
    db, err := NewGormConnection(config)
    if err != nil {
        log.Fatal(err)
    }
    
    // Migration
    if err := runMigration(db); err != nil {
        log.Fatal(err)
    }
    
    // CRUD operations
    userService := NewUserService(db)
    
    // Create user
    user := &User{
        Name:  "Ali Rezaei",
        Email: "ali@example.com",
        Age:   30,
        Role:  "admin",
    }
    if err := userService.Create(user); err != nil {
        log.Fatal(err)
    }
    
    // Get user
    user, err := userService.GetByID(1)
    if err != nil {
        log.Fatal(err)
    }
    
    // Update user
    user.Name = "Ali Mohammadi"
    if err := userService.Update(user); err != nil {
        log.Fatal(err)
    }
    
    // Get all users with pagination
    users, total, err := userService.GetAll(1, 10)
    if err != nil {
        log.Fatal(err)
    }
    
    // Search users
    users, err := userService.SearchUsers("ali")
    if err != nil {
        log.Fatal(err)
    }
    
    // Delete user
    if err := userService.Delete(1); err != nil {
        log.Fatal(err)
    }
}
`)

	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🎯 GORM GUIDE - COMPLETE")
	fmt.Println("Ready to build efficient database applications with GORM!")
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
