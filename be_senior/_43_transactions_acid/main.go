// ============================================================================
// FILE: transactions_acid_guide.go
// TITLE: راهنمای کامل تراکنش‌ها و ACID در PostgreSQL با Go
// HOW TO RUN: go run transactions_acid_guide.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - ACID چیست؟
// ============================================================================
//
// ACID مجموعه‌ای از خصوصیات است که تراکنش‌های دیتابیس را قابل اعتماد می‌کند:
//
// Atomicity (اتمی بودن):
// - تراکنش یا به طور کامل اجرا می‌شود یا اصلاً اجرا نمی‌شود
// - در صورت خطا، همه تغییرات برگردانده می‌شوند (ROLLBACK)
//
// Consistency (سازگاری):
// - تراکنش دیتابیس را از یک state معتبر به state معتبر دیگر می‌برد
// - قوانین (constraints, triggers) رعایت می‌شوند
//
// Isolation (انزوا):
// - تراکنش‌های همزمان از یکدیگر جدا هستند
// - سطوح مختلف isolation: Read Uncommitted, Read Committed, Repeatable Read, Serializable
//
// Durability (دوام):
// - پس از commit، تغییرات دائمی هستند
// - حتی در صورت قطع برق یا crash، داده از دست نمی‌رود
//
// قانون طلایی:
// "از تراکنش برای عملیات‌هایی که چندین جدول را تغییر می‌دهند استفاده کن.
//  تراکنش‌ها را کوتاه نگه دار (تا حداقل زمان ممکن).
//  همیشه از defer tx.Rollback() استفاده کن و در صورت موفقیت Commit کن."
// ============================================================================

package __transactions_acid

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// ============================================================================
// بخش 1: مفاهیم پایه تراکنش
// ============================================================================

func basicTransactionConcepts() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📚 ACID PROPERTIES")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ PROPERTY      │ DESCRIPTION                                                │
├───────────────┼─────────────────────────────────────────────────────────────┤
│ Atomicity     │ All or nothing - transaction要么全部成功要么全部失败       │
│ Consistency   │ Database remains in valid state                            │
│ Isolation     │ Concurrent transactions don't interfere                    │
│ Durability    │ Committed changes persist even after crash                 │
└───────────────┴─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ ISOLATION LEVELS                                                          │
├───────────────┬─────────────────────────────────────────────────────────────┤
│ LEVEL         │ PHENOMENA PREVENTED                                        │
├───────────────┼─────────────────────────────────────────────────────────────┤
│ Read Uncommitted │ Dirty reads (not in PostgreSQL)                         │
│ Read Committed   │ Dirty reads prevented (PostgreSQL default)              │
│ Repeatable Read  │ Dirty reads, non-repeatable reads prevented             │
│ Serializable     │ All phenomena prevented (highest isolation)             │
└───────────────┴─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ ANOMALIES                                                                 │
├───────────────┬─────────────────────────────────────────────────────────────┤
│ Dirty Read    │ Reading uncommitted data from another transaction         │
│ Non-repeatable│同一トランザクション内で同じ行を読むたびに値が変わる         │
│ Phantom Read  │同一クエリを複数回実行すると異なる行セットが返る             │
│ Serialization│ Transactions that would not occur in serial execution      │
│ Anomaly       │                                                           │
└───────────────┴─────────────────────────────────────────────────────────────┘
`)
}

// ============================================================================
// بخش 2: تراکنش پایه با database/sql
// ============================================================================

// TransactionManager مدیریت تراکنش‌ها
type TransactionManager struct {
	db *sql.DB
}

// NewTransactionManager ایجاد مدیر تراکنش
func NewTransactionManager(db *sql.DB) *TransactionManager {
	return &TransactionManager{db: db}
}

// SimpleTransaction مثال تراکنش ساده
func (tm *TransactionManager) SimpleTransaction() error {
	// شروع تراکنش
	tx, err := tm.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	// همیشه Rollback در defer - اگر Commit نشود، Rollback می‌شود
	defer tx.Rollback()

	// عملیات 1: ایجاد کاربر
	_, err = tx.Exec(`
		INSERT INTO users (username, email, age) 
		VALUES ($1, $2, $3)
	`, "john_doe", "john@example.com", 25)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	// عملیات 2: ایجاد پروفایل
	_, err = tx.Exec(`
		INSERT INTO profiles (user_id, bio, avatar) 
		VALUES (currval('users_id_seq'), $1, $2)
	`, "Software developer", "avatar.jpg")
	if err != nil {
		return fmt.Errorf("failed to create profile: %w", err)
	}

	// اگر همه چیز موفق بود، Commit کن
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Println("Transaction committed successfully")
	return nil
}

// TransferMoney انتقال وجه بین دو حساب (مثال کلاسیک)
func (tm *TransactionManager) TransferMoney(fromAccountID, toAccountID int, amount float64) error {
	tx, err := tm.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// بررسی موجودی حساب مبدأ
	var balance float64
	err = tx.QueryRow(`
		SELECT balance FROM accounts WHERE id = $1 FOR UPDATE
	`, fromAccountID).Scan(&balance)
	if err != nil {
		return fmt.Errorf("failed to get balance: %w", err)
	}

	if balance < amount {
		return fmt.Errorf("insufficient funds: balance %.2f, amount %.2f", balance, amount)
	}

	// کاهش از حساب مبدأ
	_, err = tx.Exec(`
		UPDATE accounts SET balance = balance - $1, updated_at = NOW()
		WHERE id = $2
	`, amount, fromAccountID)
	if err != nil {
		return fmt.Errorf("failed to debit from account: %w", err)
	}

	// افزایش به حساب مقصد
	_, err = tx.Exec(`
		UPDATE accounts SET balance = balance + $1, updated_at = NOW()
		WHERE id = $2
	`, amount, toAccountID)
	if err != nil {
		return fmt.Errorf("failed to credit to account: %w", err)
	}

	// ثبت تراکنش در جدول history
	_, err = tx.Exec(`
		INSERT INTO transfer_history (from_account, to_account, amount, created_at)
		VALUES ($1, $2, $3, NOW())
	`, fromAccountID, toAccountID, amount)
	if err != nil {
		return fmt.Errorf("failed to record transfer: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transfer: %w", err)
	}

	log.Printf("Transferred %.2f from account %d to account %d", amount, fromAccountID, toAccountID)
	return nil
}

// ============================================================================
// بخش 3: سطوح مختلف Isolation
// ============================================================================

// IsolationExample مثال‌های سطوح انزوا
type IsolationExample struct {
	db *sql.DB
}

// ReadCommittedExample مثال Read Committed (پیش‌فرض PostgreSQL)
func (ie *IsolationExample) ReadCommittedExample() error {
	tx, err := ie.db.BeginTx(context.Background(), &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// در سطح Read Committed:
	// - فقط داده‌های committed خوانده می‌شوند
	// - Non-repeatable read ممکن است رخ دهد

	var balance float64
	err = tx.QueryRow(`SELECT balance FROM accounts WHERE id = $1`, 1).Scan(&balance)
	if err != nil {
		return err
	}

	// اگر در این لحظه تراکنش دیگری balance را تغییر دهد،
	// خواندن دوباره مقدار متفاوتی برمی‌گرداند (non-repeatable read)

	return tx.Commit()
}

// RepeatableReadExample مثال Repeatable Read
func (ie *IsolationExample) RepeatableReadExample() error {
	tx, err := ie.db.BeginTx(context.Background(), &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  false,
	})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// در سطح Repeatable Read:
	// - Non-repeatable reads جلوگیری می‌شود
	// - Phantom reads ممکن است رخ دهد

	var balance1, balance2 float64
	err = tx.QueryRow(`SELECT balance FROM accounts WHERE id = $1`, 1).Scan(&balance1)
	if err != nil {
		return err
	}

	// در اینجا تراکنش دیگر نمی‌تواند balance را تغییر دهد
	// حتی اگر تراکنش دیگری آن را تغییر دهد، این تراکنش مقدار قبلی را می‌بیند

	err = tx.QueryRow(`SELECT balance FROM accounts WHERE id = $1`, 1).Scan(&balance2)
	if err != nil {
		return err
	}

	// balance1 و balance2 همیشه برابر هستند
	if balance1 != balance2 {
		return fmt.Errorf("repeatable read violated")
	}

	return tx.Commit()
}

// SerializableExample مثال Serializable (بالاترین سطح)
func (ie *IsolationExample) SerializableExample() error {
	tx, err := ie.db.BeginTx(context.Background(), &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// در سطح Serializable:
	// - همه anomalies جلوگیری می‌شود
	// - ممکن است serialization failure رخ دهد (باید retry کرد)

	// مثال: رزرو صندلی
	var available int
	err = tx.QueryRow(`
		SELECT available FROM seats WHERE id = $1
	`, 1).Scan(&available)
	if err != nil {
		return err
	}

	if available <= 0 {
		return fmt.Errorf("seat not available")
	}

	_, err = tx.Exec(`
		UPDATE seats SET available = available - 1, reserved_by = $1
		WHERE id = $2 AND available > 0
	`, "user123", 1)
	if err != nil {
		return err
	}

	// در سطح Serializable، اگر دو تراکنش همزمان سعی کنند
	// آخرین صندلی را رزرو کنند، یکی از آنها Serialization failure می‌خورد

	if err := tx.Commit(); err != nil {
		// اگر خطای serialization بود، می‌توان retry کرد
		if errors.Is(err, &sql.ErrTxSerialization{}) {
			return fmt.Errorf("serialization failure, please retry")
		}
		return err
	}

	return nil
}

// ============================================================================
// بخش 4: Deadlock Detection و Prevention
// ============================================================================

// DeadlockExample مثال deadlock
type DeadlockExample struct {
	db *sql.DB
}

// PotentialDeadlock عملیاتی که ممکن است deadlock ایجاد کند
func (de *DeadlockExample) PotentialDeadlock(userID1, userID2 int) error {
	tx, err := de.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// ❌ خطرناک: قفل کردن به ترتیب مختلف در تراکنش‌های مختلف
	// اگر تراکنش اول user1 را قفل کند و تراکنش دوم user2 را،
	// هر دو منتظر هم می‌مانند (deadlock)

	// قفل کردن user1
	_, err = tx.Exec(`UPDATE users SET updated_at = NOW() WHERE id = $1`, userID1)
	if err != nil {
		return err
	}

	// شبیه‌سازی کار
	time.Sleep(100 * time.Millisecond)

	// قفل کردن user2
	_, err = tx.Exec(`UPDATE users SET updated_at = NOW() WHERE id = $1`, userID2)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// SafeDeadlockPrevention جلوگیری از deadlock با ترتیب ثابت
func (de *DeadlockExample) SafeDeadlockPrevention(id1, id2 int) error {
	// ✅ راه حل: همیشه IDs را به ترتیب ثابت قفل کن
	if id1 > id2 {
		id1, id2 = id2, id1 // swap to maintain order
	}

	tx, err := de.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// قفل کردن با ترتیب ثابت (صعودی)
	_, err = tx.Exec(`UPDATE users SET updated_at = NOW() WHERE id = $1`, id1)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`UPDATE users SET updated_at = NOW() WHERE id = $1`, id2)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// RetryOnDeadlock خودکار retry در صورت deadlock
func (de *DeadlockExample) RetryOnDeadlock(fn func() error, maxRetries int) error {
	var err error
	for i := 0; i < maxRetries; i++ {
		err = fn()
		if err == nil {
			return nil
		}

		// بررسی خطای deadlock
		if isDeadlockError(err) {
			log.Printf("Deadlock detected, retrying (%d/%d)", i+1, maxRetries)
			time.Sleep(time.Duration(100*(i+1)) * time.Millisecond) // exponential backoff
			continue
		}
		return err
	}
	return fmt.Errorf("failed after %d retries: %w", maxRetries, err)
}

func isDeadlockError(err error) bool {
	return err != nil && (stringsContains(err.Error(), "deadlock detected") ||
		stringsContains(err.Error(), "could not serialize access"))
}

// ============================================================================
// بخش 5: Savepoints (نقاط ذخیره)
// ============================================================================

// SavepointExample مثال استفاده از Savepoints
type SavepointExample struct {
	db *sql.DB
}

// ProcessWithSavepoints پردازش با نقاط ذخیره
func (se *SavepointExample) ProcessWithSavepoints(records []Record) error {
	tx, err := se.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for i, record := range records {
		// ایجاد savepoint قبل از هر رکورد
		savepointName := fmt.Sprintf("sp_%d", i)
		_, err := tx.Exec(fmt.Sprintf("SAVEPOINT %s", savepointName))
		if err != nil {
			return err
		}

		// پردازش رکورد
		err = se.processRecord(tx, record)
		if err != nil {
			// در صورت خطا، فقط به savepoint برگرد
			_, rbErr := tx.Exec(fmt.Sprintf("ROLLBACK TO SAVEPOINT %s", savepointName))
			if rbErr != nil {
				return rbErr
			}
			log.Printf("Record %d failed, rolled back to savepoint", i)
			continue // ادامه با رکورد بعدی
		}

		// اگر موفق بود، savepoint را آزاد کن
		_, err = tx.Exec(fmt.Sprintf("RELEASE SAVEPOINT %s", savepointName))
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (se *SavepointExample) processRecord(tx *sql.Tx, record Record) error {
	// پردازش رکورد
	_, err := tx.Exec(`
		INSERT INTO processed_records (data, status, processed_at)
		VALUES ($1, $2, NOW())
	`, record.Data, "success")
	return err
}

type Record struct {
	Data string
}

// ============================================================================
// بخش 6: Optimistic Locking (قفل خوش‌بینانه)
// ============================================================================

// OptimisticLock مثال قفل خوش‌بینانه
type OptimisticLock struct {
	db *sql.DB
}

// UpdateWithVersion به‌روزرسانی با استفاده از version (optimistic locking)
func (ol *OptimisticLock) UpdateWithVersion(id int, name string, version int) error {
	result, err := ol.db.Exec(`
		UPDATE users 
		SET name = $1, version = version + 1, updated_at = NOW()
		WHERE id = $2 AND version = $3
	`, name, id, version)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("record was modified by another transaction, please retry")
	}

	return nil
}

// PessimisticLock مثال قفل بدبینانه (با FOR UPDATE)
type PessimisticLock struct {
	db *sql.DB
}

// UpdateWithLock به‌روزرسانی با قفل بدبینانه
func (pl *PessimisticLock) UpdateWithLock(id int, name string) error {
	tx, err := pl.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// قفل کردن ردیف (تا پایان تراکنش)
	var currentName string
	err = tx.QueryRow(`
		SELECT name FROM users WHERE id = $1 FOR UPDATE
	`, id).Scan(&currentName)
	if err != nil {
		return err
	}

	// به‌روزرسانی
	_, err = tx.Exec(`
		UPDATE users SET name = $1, updated_at = NOW()
		WHERE id = $2
	`, name, id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// ============================================================================
// بخش 7: Two-Phase Commit (2PC) برای تراکنش‌های توزیع‌شده
// ============================================================================

// TwoPhaseCommit مثال Two-Phase Commit
type TwoPhaseCommit struct {
	db1 *sql.DB
	db2 *sql.DB
}

// PrepareTwoPhase آماده‌سازی تراکنش دو فازی
func (tpc *TwoPhaseCommit) PrepareTwoPhase() error {
	// مرحله 1: PREPARE - آماده‌سازی روی هر دو دیتابیس
	tx1, err := tpc.db1.Begin()
	if err != nil {
		return err
	}

	tx2, err := tpc.db2.Begin()
	if err != nil {
		tx1.Rollback()
		return err
	}

	// انجام عملیات روی هر دو دیتابیس
	if err := tpc.operationOnDB1(tx1); err != nil {
		tx1.Rollback()
		tx2.Rollback()
		return err
	}

	if err := tpc.operationOnDB2(tx2); err != nil {
		tx1.Rollback()
		tx2.Rollback()
		return err
	}

	// مرحله 2: COMMIT - commit هر دو
	if err := tx1.Commit(); err != nil {
		tx2.Rollback()
		return err
	}

	if err := tx2.Commit(); err != nil {
		// در صورت خطا در دومی، اولی قبلاً commit شده - inconsistent state
		// نیاز به compensating transaction
		return err
	}

	return nil
}

func (tpc *TwoPhaseCommit) operationOnDB1(tx *sql.Tx) error {
	_, err := tx.Exec("UPDATE accounts SET balance = balance - 100 WHERE id = 1")
	return err
}

func (tpc *TwoPhaseCommit) operationOnDB2(tx *sql.Tx) error {
	_, err := tx.Exec("UPDATE accounts SET balance = balance + 100 WHERE id = 2")
	return err
}

// ============================================================================
// بخش 8: Distributed Transaction Pattern (Saga)
// ============================================================================

// SagaPattern الگوی Saga برای تراکنش‌های توزیع‌شده
type SagaPattern struct {
	db *sql.DB
}

// SagaStep یک مرحله در Saga
type SagaStep struct {
	Name       string
	Execute    func() error
	Compensate func() error
}

// ExecuteSaga اجرای Saga با قابلیت جبران
func (s *SagaPattern) ExecuteSaga(steps []SagaStep) error {
	executedSteps := make([]int, 0)

	for i, step := range steps {
		log.Printf("Executing step: %s", step.Name)

		if err := step.Execute(); err != nil {
			log.Printf("Step %s failed: %v", step.Name, err)

			// جبران مراحل قبلی (به ترتیب معکوس)
			for j := len(executedSteps) - 1; j >= 0; j-- {
				stepIdx := executedSteps[j]
				log.Printf("Compensating step: %s", steps[stepIdx].Name)
				if err := steps[stepIdx].Compensate(); err != nil {
					log.Printf("Failed to compensate step %s: %v", steps[stepIdx].Name, err)
				}
			}
			return fmt.Errorf("saga failed at step %s: %w", step.Name, err)
		}

		executedSteps = append(executedSteps, i)
	}

	log.Println("Saga completed successfully")
	return nil
}

// ExampleOrderSaga مثال سفارش با Saga
func (s *SagaPattern) ExampleOrderSaga(orderID int) error {
	steps := []SagaStep{
		{
			Name: "Reserve Stock",
			Execute: func() error {
				_, err := s.db.Exec(`
					UPDATE products SET stock = stock - 1 
					WHERE id = $1 AND stock > 0
				`, orderID)
				return err
			},
			Compensate: func() error {
				_, err := s.db.Exec(`
					UPDATE products SET stock = stock + 1 
					WHERE id = $1
				`, orderID)
				return err
			},
		},
		{
			Name: "Process Payment",
			Execute: func() error {
				// درخواست به درگاه پرداخت
				return nil
			},
			Compensate: func() error {
				// برگرداندن وجه
				return nil
			},
		},
		{
			Name: "Create Order",
			Execute: func() error {
				_, err := s.db.Exec(`
					INSERT INTO orders (id, status, created_at)
					VALUES ($1, 'confirmed', NOW())
				`, orderID)
				return err
			},
			Compensate: func() error {
				_, err := s.db.Exec(`
					UPDATE orders SET status = 'cancelled' 
					WHERE id = $1
				`, orderID)
				return err
			},
		},
	}

	return s.ExecuteSaga(steps)
}

// ============================================================================
// بخش 9: Best Practices
// ============================================================================

func bestPractices() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("💡 TRANSACTION BEST PRACTICES")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────┐
│ TRANSACTION DESIGN                                            │
├─────────────────────────────────────────────────────────────────┤
│ • Keep transactions SHORT                                     │
│   - Don't include user interaction                             │
│   - Don't do heavy processing                                  │
│   - Don't call external APIs                                   │
│                                                                 │
│ • Always use defer tx.Rollback()                              │
│   - Prevents leaked transactions                               │
│   - Safe even if Commit succeeds                               │
│                                                                 │
│ • Choose appropriate isolation level                          │
│   - Read Committed for most cases                              │
│   - Repeatable Read for consistent reads                       │
│   - Serializable when needed (with retry)                      │
│                                                                 │
│ • Handle retries for serialization failures                    │
│   - Implement exponential backoff                              │
│   - Set max retry limit                                        │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ LOCKING STRATEGIES                                            │
├─────────────────────────────────────────────────────────────────┤
│ • Pessimistic Locking (FOR UPDATE)                            │
│   - Good for high contention                                   │
│   - Can cause deadlocks                                        │
│   - Lock until transaction ends                                │
│                                                                 │
│ • Optimistic Locking (version column)                         │
│   - Good for low contention                                    │
│   - No deadlocks                                               │
│   - Retry on conflict                                          │
│                                                                 │
│ • Lock Ordering                                                │
│   - Always acquire locks in same order                        │
│   - Prevents deadlocks                                         │
│   - Document lock order                                        │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ ERROR HANDLING                                                │
├─────────────────────────────────────────────────────────────────┤
│ • Always check for serialization failures                     │
│ • Implement retry logic                                        │
│ • Log transaction duration                                     │
│ • Monitor deadlock frequency                                   │
│ • Use timeouts for long-running transactions                   │
└─────────────────────────────────────────────────────────────────┘
`)
}

// ============================================================================
// بخش 10: Monitoring
// ============================================================================

// TransactionMonitor مانیتور تراکنش‌ها
type TransactionMonitor struct {
	db *sql.DB
}

// GetActiveTransactions دریافت تراکنش‌های فعال
func (tm *TransactionMonitor) GetActiveTransactions() ([]ActiveTransaction, error) {
	query := `
		SELECT 
			pid,
			usename,
			application_name,
			state,
			query,
			now() - xact_start as transaction_duration,
			now() - query_start as query_duration
		FROM pg_stat_activity
		WHERE xact_start IS NOT NULL
		AND state != 'idle'
		ORDER BY xact_start
	`

	rows, err := tm.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []ActiveTransaction
	for rows.Next() {
		var t ActiveTransaction
		err := rows.Scan(&t.PID, &t.Username, &t.AppName, &t.State, &t.Query,
			&t.TransactionDuration, &t.QueryDuration)
		if err != nil {
			continue
		}
		transactions = append(transactions, t)
	}

	return transactions, nil
}

type ActiveTransaction struct {
	PID                 int
	Username            string
	AppName             string
	State               string
	Query               string
	TransactionDuration interface{}
	QueryDuration       interface{}
}

// ============================================================================
// بخش 11: Types
// ============================================================================

// ============================================================================
// بخش 12: Main
// ============================================================================

func main() {
	fmt.Println(stringsRepeat("=", 80))
	fmt.Println("🎯 TRANSACTIONS & ACID GUIDE")
	fmt.Println("PostgreSQL Transactions in Go")
	fmt.Println(stringsRepeat("=", 80))

	basicTransactionConcepts()
	bestPractices()

	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📚 QUICK REFERENCE")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ TRANSACTION PATTERNS                                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│ 1. Basic Transaction                                                       │
│    tx, _ := db.Begin()                                                     │
│    defer tx.Rollback()                                                     │
│    // operations                                                           │
│    tx.Commit()                                                             │
│                                                                             │
│ 2. Transaction with Context                                                │
│    tx, _ := db.BeginTx(ctx, &sql.TxOptions{                                │
│        Isolation: sql.LevelRepeatableRead,                                 │
│        ReadOnly: false,                                                    │
│    })                                                                      │
│                                                                             │
│ 3. Prepared Statement in Transaction                                       │
│    stmt, _ := tx.Prepare("UPDATE users SET name = $1 WHERE id = $2")      │
│    defer stmt.Close()                                                      │
│    stmt.Exec("Ali", 1)                                                     │
│                                                                             │
│ 4. FOR UPDATE (row-level lock)                                             │
│    tx.QueryRow("SELECT balance FROM accounts WHERE id = $1 FOR UPDATE", 1) │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

🔍 MONITORING QUERIES:

   # Active transactions
   SELECT * FROM pg_stat_activity WHERE xact_start IS NOT NULL;

   # Long-running transactions
   SELECT pid, now() - xact_start as duration, query 
   FROM pg_stat_activity 
   WHERE xact_start IS NOT NULL 
   AND now() - xact_start > interval '5 minutes';

   # Deadlocks
   SELECT * FROM pg_locks WHERE granted = false;

   # Transaction statistics
   SELECT datname, xact_commit, xact_rollback 
   FROM pg_stat_database;
`)

	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🎯 TRANSACTIONS & ACID - COMPLETE")
	fmt.Println("Ready to build reliable database applications!")
	fmt.Println(stringsRepeat("=", 80))
}

// تابع کمکی
func stringsContains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && (s[0:len(substr)] == substr ||
			(len(s) > len(substr) && s[len(s)-len(substr):] == substr) ||
			(len(s) > len(substr) && stringsContains(s[1:], substr)))))
}

func stringsRepeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}
