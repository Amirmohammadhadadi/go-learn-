// ============================================================================
// FILE: time_complete_guide.go
// TITLE: راهنمای کامل پکیج time در Go - زمان، تاریخ، تایمر، ددلاین
// HOW TO RUN: go run time_complete_guide.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - پکیج time چیست و چرا مهم است؟
// ============================================================================
//
// پکیج time امکانات کامل برای:
// 1. کار با زمان و تاریخ (Time struct)
// 2. اندازه‌گیری زمان (Duration)
// 3. تایمرها و تیکرها (Timer, Ticker)
// 4. فرمت‌سازی و parsing زمان
// 5. محاسبات زمانی (افزودن، تفریق، مقایسه)
// 6. منطقه زمانی (Time Zone)
//
// قانون طلایی:
// "از time.Time برای ذخیره زمان، از time.Duration برای اندازه‌گیری فاصله،
//  از time.After و time.Ticker برای عملیات تایم‌اوت و تکراری استفاده کن.
//  همیشه از time.Time استفاده کن، نه از int64 یا string."
// ============================================================================

package __internal_packages

import (
	"fmt"
	"time"
)

// ============================================================================
// بخش 1: Time ساختار پایه - ایجاد و دستکاری زمان
// ============================================================================

func demonstrateTimeBasics() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("⏰ TIME BASICS - Creating and Manipulating Time")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 1.1 گرفتن زمان فعلی
	// ============================================
	fmt.Println("\n--- 1.1 Current Time ---")

	now := time.Now()
	fmt.Printf("  Current time: %v\n", now)
	fmt.Printf("  Unix timestamp: %d\n", now.Unix())
	fmt.Printf("  Unix nano: %d\n", now.UnixNano())
	fmt.Printf("  Year: %d, Month: %s, Day: %d\n", now.Year(), now.Month(), now.Day())
	fmt.Printf("  Hour: %d, Minute: %d, Second: %d\n", now.Hour(), now.Minute(), now.Second())
	fmt.Printf("  Weekday: %s\n", now.Weekday())
	fmt.Printf("  YearDay: %d\n", now.YearDay())

	// ============================================
	// 1.2 ایجاد زمان مشخص
	// ============================================
	fmt.Println("\n--- 1.2 Creating Specific Time ---")

	// با تاریخ و زمان مشخص
	t1 := time.Date(2024, time.January, 15, 14, 30, 0, 0, time.UTC)
	fmt.Printf("  Specific time (UTC): %v\n", t1)

	// با زمان محلی
	t2 := time.Date(2024, time.December, 25, 10, 0, 0, 0, time.Local)
	fmt.Printf("  Specific time (Local): %v\n", t2)

	// ============================================
	// 1.3 از روی Unix timestamp
	// ============================================
	fmt.Println("\n--- 1.3 From Unix Timestamp ---")

	timestamp := int64(1705312200)
	t3 := time.Unix(timestamp, 0)
	fmt.Printf("  From seconds: %v\n", t3)

	// با نانوثانیه
	t4 := time.Unix(0, timestamp*1000000000)
	fmt.Printf("  From nanoseconds: %v\n", t4)

	// ============================================
	// 1.4 از روی Parsing
	// ============================================
	fmt.Println("\n--- 1.4 From Parsing ---")

	// Parsing با فرمت استاندارد
	t5, _ := time.Parse(time.RFC3339, "2024-01-15T14:30:00Z")
	fmt.Printf("  Parsed RFC3339: %v\n", t5)

	// ============================================
	// 1.5 اجزای زمان
	// ============================================
	fmt.Println("\n--- 1.5 Time Components ---")

	now = time.Now()
	fmt.Printf("  Date: %d-%02d-%02d\n", now.Year(), now.Month(), now.Day())
	fmt.Printf("  Time: %02d:%02d:%02d\n", now.Hour(), now.Minute(), now.Second())
	fmt.Printf("  Nanosecond: %d\n", now.Nanosecond())
	fmt.Printf("  Location: %s\n", now.Location())
	fmt.Printf("  Zone: %s, offset: %d\n", now.Zone())
	fmt.Printf("  ISO Week: %d\n", now.ISOWeek())
}

// ============================================================================
// بخش 2: Duration - اندازه‌گیری فاصله زمانی
// ============================================================================

func demonstrateDuration() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("⏱️ DURATION - Measuring Time Intervals")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 2.1 ایجاد Duration
	// ============================================
	fmt.Println("\n--- 2.1 Creating Durations ---")

	d1 := 10 * time.Second
	d2 := 5 * time.Minute
	d3 := 2 * time.Hour
	d4 := 100 * time.Millisecond
	d5 := 500 * time.Microsecond
	d6 := 1000 * time.Nanosecond

	fmt.Printf("  Seconds: %v\n", d1)
	fmt.Printf("  Minutes: %v\n", d2)
	fmt.Printf("  Hours: %v\n", d3)
	fmt.Printf("  Milliseconds: %v\n", d4)
	fmt.Printf("  Microseconds: %v\n", d5)
	fmt.Printf("  Nanoseconds: %v\n", d6)

	// ============================================
	// 2.2 محاسبه Duration بین دو زمان
	// ============================================
	fmt.Println("\n--- 2.2 Calculating Duration ---")

	start := time.Now()
	time.Sleep(100 * time.Millisecond)
	elapsed := time.Since(start)

	fmt.Printf("  Elapsed: %v\n", elapsed)
	fmt.Printf("  Milliseconds: %d\n", elapsed.Milliseconds())
	fmt.Printf("  Microseconds: %d\n", elapsed.Microseconds())
	fmt.Printf("  Nanoseconds: %d\n", elapsed.Nanoseconds())
	fmt.Printf("  Seconds (float): %.3f\n", elapsed.Seconds())
	fmt.Printf("  Minutes (float): %.3f\n", elapsed.Minutes())
	fmt.Printf("  Hours (float): %.3f\n", elapsed.Hours())

	// ============================================
	// 2.3 عملیات روی Duration
	// ============================================
	fmt.Println("\n--- 2.3 Duration Operations ---")

	dur1 := 10 * time.Second
	dur2 := 30 * time.Second

	fmt.Printf("  dur1: %v, dur2: %v\n", dur1, dur2)
	fmt.Printf("  dur1 + dur2: %v\n", dur1+dur2)
	fmt.Printf("  dur2 - dur1: %v\n", dur2-dur1)
	fmt.Printf("  dur1 * 3: %v\n", dur1*3)
	fmt.Printf("  dur2 / 2: %v\n", dur2/2)
	fmt.Printf("  dur1 < dur2: %v\n", dur1 < dur2)
	fmt.Printf("  dur1 == dur2: %v\n", dur1 == dur2)

	// ============================================
	// 2.4 Rounding و Truncation
	// ============================================
	fmt.Println("\n--- 2.4 Rounding and Truncation ---")

	dur := 2*time.Second + 350*time.Millisecond
	fmt.Printf("  Original: %v\n", dur)
	fmt.Printf("  Round to second: %v\n", dur.Round(time.Second))
	fmt.Printf("  Round to 500ms: %v\n", dur.Round(500*time.Millisecond))
	fmt.Printf("  Truncate to second: %v\n", dur.Truncate(time.Second))
}

// ============================================================================
// بخش 3: عملیات روی زمان - افزودن، تفریق، مقایسه
// ============================================================================

func demonstrateTimeOperations() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("➕ TIME OPERATIONS - Add, Subtract, Compare")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 3.1 افزودن زمان (Add)
	// ============================================
	fmt.Println("\n--- 3.1 Adding Time (Add) ---")

	now := time.Now()
	fmt.Printf("  Now: %v\n", now)

	// افزودن Duration
	tomorrow := now.Add(24 * time.Hour)
	fmt.Printf("  Tomorrow: %v\n", tomorrow)

	nextHour := now.Add(1 * time.Hour)
	fmt.Printf("  Next hour: %v\n", nextHour)

	nextWeek := now.AddDate(0, 0, 7)
	fmt.Printf("  Next week: %v\n", nextWeek)

	nextMonth := now.AddDate(0, 1, 0)
	fmt.Printf("  Next month: %v\n", nextMonth)

	nextYear := now.AddDate(1, 0, 0)
	fmt.Printf("  Next year: %v\n", nextYear)

	// ============================================
	// 3.2 تفریق زمان (Sub)
	// ============================================
	fmt.Println("\n--- 3.2 Subtracting Time (Sub) ---")

	past := now.Add(-48 * time.Hour)
	duration := now.Sub(past)
	fmt.Printf("  Past: %v\n", past)
	fmt.Printf("  Difference: %v (%.1f hours)\n", duration, duration.Hours())

	// ============================================
	// 3.3 مقایسه زمان‌ها
	// ============================================
	fmt.Println("\n--- 3.3 Comparing Times ---")

	t1 := time.Now()
	t2 := t1.Add(1 * time.Hour)
	t3 := t1.Add(-1 * time.Hour)

	fmt.Printf("  t1: %v\n", t1)
	fmt.Printf("  t2: %v\n", t2)
	fmt.Printf("  t3: %v\n", t3)

	fmt.Printf("  t1.Before(t2): %v\n", t1.Before(t2))
	fmt.Printf("  t1.After(t2): %v\n", t1.After(t2))
	fmt.Printf("  t1.Equal(t2): %v\n", t1.Equal(t2))
	fmt.Printf("  t1.Equal(t1): %v\n", t1.Equal(t1))

	// ============================================
	// 3.4 UTC و Local تبدیل
	// ============================================
	fmt.Println("\n--- 3.4 UTC and Local Conversion ---")

	local := time.Now()
	utc := local.UTC()

	fmt.Printf("  Local: %v\n", local)
	fmt.Printf("  UTC: %v\n", utc)
	fmt.Printf("  Local zone: %s\n", local.Location())
	fmt.Printf("  UTC zone: %s\n", utc.Location())

	// تبدیل به منطقه زمانی دیگر
	location, _ := time.LoadLocation("America/New_York")
	nyTime := local.In(location)
	fmt.Printf("  New York time: %v\n", nyTime)
}

// ============================================================================
// بخش 4: فرمت‌سازی و Parsing زمان
// ============================================================================

func demonstrateTimeFormatting() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📝 TIME FORMATTING AND PARSING")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 4.1 فرمت‌های استاندارد
	// ============================================
	fmt.Println("\n--- 4.1 Standard Formats ---")

	now := time.Now()

	formats := []struct {
		name   string
		layout string
	}{
		{"RFC3339", time.RFC3339},
		{"RFC3339Nano", time.RFC3339Nano},
		{"RFC1123", time.RFC1123},
		{"RFC1123Z", time.RFC1123Z},
		{"RFC822", time.RFC822},
		{"RFC822Z", time.RFC822Z},
		{"RFC850", time.RFC850},
		{"Kitchen", time.Kitchen},
		{"ANSIC", time.ANSIC},
		{"UnixDate", time.UnixDate},
		{"RubyDate", time.RubyDate},
	}

	for _, f := range formats {
		fmt.Printf("  %s: %s\n", f.name, now.Format(f.layout))
	}

	// ============================================
	// 4.2 فرمت سفارشی
	// ============================================
	fmt.Println("\n--- 4.2 Custom Formats ---")

	// نکته مهم: Go از یک زمان مرجع خاص برای فرمت استفاده می‌کند:
	// Mon Jan 2 15:04:05 MST 2006 (01/02 03:04:05PM '06 -0700)
	// 1  2   3  4  5   6   7   8   9   10  11  12  13  14

	customFormats := []struct {
		name   string
		layout string
	}{
		{"YYYY-MM-DD", "2006-01-02"},
		{"DD/MM/YYYY", "02/01/2006"},
		{"Month DD, YYYY", "January 02, 2006"},
		{"HH:MM:SS", "15:04:05"},
		{"HH:MM AM/PM", "03:04 PM"},
		{"Full custom", "2006-01-02 15:04:05 (Monday)"},
		{"With timezone", "2006-01-02 15:04:05 MST"},
		{"Weekday, Date", "Monday, January 2"},
	}

	for _, f := range customFormats {
		fmt.Printf("  %s: %s\n", f.name, now.Format(f.layout))
	}

	// ============================================
	// 4.3 Parsing زمان
	// ============================================
	fmt.Println("\n--- 4.3 Parsing Time ---")

	parseExamples := []struct {
		layout string
		value  string
	}{
		{time.RFC3339, "2024-01-15T14:30:00Z"},
		{"2006-01-02", "2024-01-15"},
		{"02/01/2006", "15/01/2024"},
		{"January 2, 2006", "January 15, 2024"},
		{"2006-01-02 15:04:05", "2024-01-15 14:30:00"},
	}

	for _, ex := range parseExamples {
		parsed, err := time.Parse(ex.layout, ex.value)
		if err == nil {
			fmt.Printf("  Parsed '%s' → %v\n", ex.value, parsed)
		} else {
			fmt.Printf("  Error parsing '%s': %v\n", ex.value, err)
		}
	}

	// ============================================
	// 4.4 Parsing با منطقه زمانی
	// ============================================
	fmt.Println("\n--- 4.4 Parsing with Timezone ---")

	// با location خاص
	loc, _ := time.LoadLocation("Asia/Tehran")
	parsed, _ := time.ParseInLocation("2006-01-02 15:04:05", "2024-01-15 14:30:00", loc)
	fmt.Printf("  Parsed in Tehran: %v\n", parsed)
	fmt.Printf("  UTC equivalent: %v\n", parsed.UTC())
}

// ============================================================================
// بخش 5: تایمرها (Timer)
// ============================================================================

func demonstrateTimers() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("⏲️ TIMERS - One-time Delayed Execution")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 5.1 Timer ساده
	// ============================================
	fmt.Println("\n--- 5.1 Basic Timer ---")

	timer1 := time.NewTimer(1 * time.Second)
	fmt.Println("  Timer started, waiting...")

	<-timer1.C
	fmt.Println("  Timer expired!")

	// ============================================
	// 5.2 Timer با After
	// ============================================
	fmt.Println("\n--- 5.2 Timer with After ---")

	fmt.Println("  Waiting with After...")
	<-time.After(500 * time.Millisecond)
	fmt.Println("  After expired!")

	// ============================================
	// 5.3 Stop کردن Timer
	// ============================================
	fmt.Println("\n--- 5.3 Stopping Timer ---")

	timer2 := time.NewTimer(2 * time.Second)

	stopped := timer2.Stop()
	if stopped {
		fmt.Println("  Timer stopped before expiration")
	} else {
		fmt.Println("  Timer already expired")
	}

	// ============================================
	// 5.4 Reset کردن Timer
	// ============================================
	fmt.Println("\n--- 5.4 Resetting Timer ---")

	timer3 := time.NewTimer(1 * time.Second)

	go func() {
		<-timer3.C
		fmt.Println("  Timer expired (first time)")
	}()

	time.Sleep(500 * time.Millisecond)
	timer3.Reset(1 * time.Second)
	fmt.Println("  Timer reset, waiting again...")

	time.Sleep(1 * time.Second)
	fmt.Println("  Timer expired (second time)")

	// ============================================
	// 5.5 Timer با AfterFunc
	// ============================================
	fmt.Println("\n--- 5.5 AfterFunc ---")

	done := make(chan bool)

	time.AfterFunc(1*time.Second, func() {
		fmt.Println("  AfterFunc executed!")
		done <- true
	})

	<-done
}

// ============================================================================
// بخش 6: تیکرها (Ticker)
// ============================================================================

func demonstrateTickers() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🔄 TICKERS - Periodic Execution")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 6.1 Ticker ساده
	// ============================================
	fmt.Println("\n--- 6.1 Basic Ticker ---")

	ticker := time.NewTicker(500 * time.Millisecond)

	count := 0
	done := make(chan bool)

	go func() {
		for range ticker.C {
			count++
			fmt.Printf("  Tick %d\n", count)
			if count >= 3 {
				done <- true
			}
		}
	}()

	<-done
	ticker.Stop()
	fmt.Println("  Ticker stopped")

	// ============================================
	// 6.2 Ticker با range
	// ============================================
	fmt.Println("\n--- 6.2 Ticker with Range ---")

	ticker2 := time.NewTicker(200 * time.Millisecond)

	go func() {
		for tick := range ticker2.C {
			fmt.Printf("  Tick at: %v\n", tick.Format("15:04:05.000"))
		}
	}()

	time.Sleep(1 * time.Second)
	ticker2.Stop()
	fmt.Println("  Ticker stopped")

	// ============================================
	// 6.3 Ticker با NewTicker vs Tick
	// ============================================
	fmt.Println("\n--- 6.3 NewTicker vs Tick ---")

	// Tick (ساده‌تر اما نمی‌توان Stop کرد)
	tickChan := time.Tick(300 * time.Millisecond)

	go func() {
		for i := 0; i < 2; i++ {
			<-tickChan
			fmt.Println("  Tick from time.Tick()")
		}
	}()

	time.Sleep(1 * time.Second)
}

// ============================================================================
// بخش 7: Timeout و Deadline با Context (مرور)
// ============================================================================

func demonstrateTimeouts() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("⏰ TIMEOUT AND DEADLINE PATTERNS")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 7.1 Timeout با select و time.After
	// ============================================
	fmt.Println("\n--- 7.1 Timeout with select and time.After ---")

	ch := make(chan string)

	go func() {
		time.Sleep(2 * time.Second)
		ch <- "result"
	}()

	select {
	case res := <-ch:
		fmt.Printf("  Received: %s\n", res)
	case <-time.After(1 * time.Second):
		fmt.Println("  Timeout! Operation took too long")
	}

	// ============================================
	// 7.2 Timeout در حلقه
	// ============================================
	fmt.Println("\n--- 7.2 Timeout in Loop ---")

	work := make(chan int)

	go func() {
		for i := 1; i <= 5; i++ {
			time.Sleep(300 * time.Millisecond)
			work <- i
		}
		close(work)
	}()

	for {
		select {
		case v, ok := <-work:
			if !ok {
				fmt.Println("  Work channel closed")
				return
			}
			fmt.Printf("  Received: %d\n", v)
		case <-time.After(1 * time.Second):
			fmt.Println("  Timeout! No activity")
			return
		}
	}
}

// ============================================================================
// بخش 8: Sleep و عملیات مسدودکننده
// ============================================================================

func demonstrateSleep() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("💤 SLEEP AND BLOCKING OPERATIONS")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 8.1 Sleep ساده
	// ============================================
	fmt.Println("\n--- 8.1 Basic Sleep ---")

	fmt.Println("  Sleeping for 500ms...")
	start := time.Now()
	time.Sleep(500 * time.Millisecond)
	fmt.Printf("  Awake after %v\n", time.Since(start))

	// ============================================
	// 8.2 Sleep قابل قطع
	// ============================================
	fmt.Println("\n--- 8.2 Interruptible Sleep ---")

	cancel := make(chan bool)

	go func() {
		time.Sleep(2 * time.Second)
		cancel <- true
	}()

	select {
	case <-time.After(1 * time.Second):
		fmt.Println("  Sleep completed")
	case <-cancel:
		fmt.Println("  Sleep interrupted")
	}
}

// ============================================================================
// بخش 9: منطقه زمانی (Time Zone)
// ============================================================================

func demonstrateTimeZones() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🌍 TIME ZONES - Working with Different Time Zones")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 9.1 مناطق زمانی موجود
	// ============================================
	fmt.Println("\n--- 9.1 Available Time Zones ---")

	zones := []string{
		"UTC",
		"Local",
		"America/New_York",
		"Europe/London",
		"Asia/Tehran",
		"Asia/Tokyo",
		"Australia/Sydney",
	}

	for _, zone := range zones {
		if loc, err := time.LoadLocation(zone); err == nil {
			fmt.Printf("  %s: %s\n", zone, loc)
		}
	}

	// ============================================
	// 9.2 تبدیل بین مناطق زمانی
	// ============================================
	fmt.Println("\n--- 9.2 Converting Between Time Zones ---")

	now := time.Now()

	zonesToShow := []string{
		"UTC",
		"America/New_York",
		"Europe/London",
		"Asia/Tehran",
		"Asia/Tokyo",
	}

	fmt.Printf("  Original (Local): %v\n", now)

	for _, zone := range zonesToShow {
		loc, _ := time.LoadLocation(zone)
		converted := now.In(loc)
		fmt.Printf("  %s: %v\n", zone, converted)
	}

	// ============================================
	// 9.3 ایجاد زمان در منطقه خاص
	// ============================================
	fmt.Println("\n--- 9.3 Creating Time in Specific Zone ---")

	tehranLoc, _ := time.LoadLocation("Asia/Tehran")
	tehranTime := time.Date(2024, time.January, 15, 14, 30, 0, 0, tehranLoc)

	fmt.Printf("  Time in Tehran: %v\n", tehranTime)
	fmt.Printf("  UTC equivalent: %v\n", tehranTime.UTC())
	fmt.Printf("  New York time: %v\n", tehranTime.In(time.UTC))
}

// ============================================================================
// بخش 10: کاربردهای عملی
// ============================================================================

func demonstratePracticalUses() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("💡 PRACTICAL USE CASES")
	fmt.Println(stringsRepeat("=", 80))

	// ============================================
	// 10.1 اندازه‌گیری زمان اجرا
	// ============================================
	fmt.Println("\n--- 10.1 Measuring Execution Time ---")

	measureFunc := func(name string, fn func()) {
		start := time.Now()
		fn()
		elapsed := time.Since(start)
		fmt.Printf("  %s took %v\n", name, elapsed)
	}

	measureFunc("Sleep 100ms", func() {
		time.Sleep(100 * time.Millisecond)
	})

	// ============================================
	// 10.2 Retry با Backoff
	// ============================================
	fmt.Println("\n--- 10.2 Retry with Backoff ---")

	retryWithBackoff := func(maxRetries int, fn func() error) error {
		var err error
		for i := 0; i < maxRetries; i++ {
			err = fn()
			if err == nil {
				return nil
			}

			// exponential backoff
			backoff := time.Duration(1<<uint(i)) * 100 * time.Millisecond
			fmt.Printf("  Retry %d after %v\n", i+1, backoff)
			time.Sleep(backoff)
		}
		return err
	}

	attempts := 0
	retryWithBackoff(3, func() error {
		attempts++
		if attempts < 2 {
			return fmt.Errorf("simulated error")
		}
		fmt.Printf("  Success on attempt %d\n", attempts)
		return nil
	})

	// ============================================
	// 10.3 Rate Limiter ساده
	// ============================================
	fmt.Println("\n--- 10.3 Simple Rate Limiter ---")

	type RateLimiter struct {
		ticker *time.Ticker
	}

	func NewRateLimiter(rate time.Duration) *RateLimiter {
		return &RateLimiter{
		ticker: time.NewTicker(rate),
	}
	}

	func (rl *RateLimiter) Allow() bool {
		select {
	case <-rl.ticker.C:
		return true
	default:
		return false
	}
	}

	limiter := NewRateLimiter(200 * time.Millisecond)
	for i := 0; i < 5; i++ {
		if limiter.Allow() {
			fmt.Printf("  Request %d allowed\n", i+1)
		} else {
			fmt.Printf("  Request %d rate limited\n", i+1)
		}
		time.Sleep(100 * time.Millisecond)
	}

	// ============================================
	// 10.4 Cron-like Scheduler
	// ============================================
	fmt.Println("\n--- 10.4 Simple Scheduler ---")

	scheduleAt := func(targetHour, targetMinute int, fn func()) {
		now := time.Now()
		target := time.Date(now.Year(), now.Month(), now.Day(),
			targetHour, targetMinute, 0, 0, now.Location())

		if target.Before(now) {
			target = target.Add(24 * time.Hour)
		}

		duration := target.Sub(now)
		fmt.Printf("  Scheduled in %v\n", duration)

		time.AfterFunc(duration, fn)
	}

	scheduleAt(15, 30, func() {
		fmt.Println("  Scheduled task executed!")
	})

	// ============================================
	// 10.5 Parse duration از string
	// ============================================
	fmt.Println("\n--- 10.5 Parsing Duration from String ---")

	durationStrings := []string{"1h30m", "2m15s", "500ms", "1.5h"}

	for _, ds := range durationStrings {
		d, err := time.ParseDuration(ds)
		if err == nil {
			fmt.Printf("  '%s' → %v\n", ds, d)
		}
	}
}

// ============================================================================
// بخش 11: اشتباهات رایج
// ============================================================================

func demonstrateCommonMistakes() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("❌ COMMON MISTAKES WITH time PACKAGE")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println("\n❌ Mistake 1: Using time.Sleep for synchronization")
	fmt.Println("   time.Sleep(1 * time.Second)  // unpredictable")
	fmt.Println("   ✅ Use channels or sync.WaitGroup")

	fmt.Println("\n❌ Mistake 2: Comparing times with == or !=")
	fmt.Println("   t1 == t2  // only works for exact same time")
	fmt.Println("   ✅ Use t1.Equal(t2) for correct comparison")

	fmt.Println("\n❌ Mistake 3: Not handling timezone")
	fmt.Println("   time.Parse(\"2006-01-02\", \"2024-01-15\")  // UTC")
	fmt.Println("   ✅ Use ParseInLocation for local time")

	fmt.Println("\n❌ Mistake 4: Forgetting to stop Ticker")
	fmt.Println("   ticker := time.NewTicker(1*time.Second)")
	fmt.Println("   // missing ticker.Stop() → goroutine leak")
	fmt.Println("   ✅ defer ticker.Stop()")

	fmt.Println("\n❌ Mistake 5: Using time.After in loop (memory leak)")
	fmt.Println("   for { select { case <-time.After(1s): } }")
	fmt.Println("   ✅ Use time.NewTicker for repeated operations")

	fmt.Println("\n❌ Mistake 6: Not resetting Timer after expiration")
	fmt.Println("   timer := time.NewTimer(1s)")
	fmt.Println("   <-timer.C")
	fmt.Println("   timer.Reset(1s)  // works but careful")
	fmt.Println("   ✅ Check timer.Stop() before Reset")

	fmt.Println("\n❌ Mistake 7: Assuming monotonic clock for all operations")
	fmt.Println("   time.Now() includes wall clock and monotonic")
	fmt.Println("   ✅ Use time.Since() for measuring intervals")
}

// ============================================================================
// بخش 12: جمع‌بندی و جدول مرجع
// ============================================================================

func main() {
	fmt.Println(stringsRepeat("=", 80))
	fmt.Println("🎯 COMPLETE time PACKAGE GUIDE IN GO")
	fmt.Println(stringsRepeat("=", 80))

	// بخش 1: Time Basics
	demonstrateTimeBasics()

	// بخش 2: Duration
	demonstrateDuration()

	// بخش 3: Time Operations
	demonstrateTimeOperations()

	// بخش 4: Formatting
	demonstrateTimeFormatting()

	// بخش 5: Timers
	demonstrateTimers()

	// بخش 6: Tickers
	demonstrateTickers()

	// بخش 7: Timeouts
	demonstrateTimeouts()

	// بخش 8: Sleep
	demonstrateSleep()

	// بخش 9: Time Zones
	demonstrateTimeZones()

	// بخش 10: Practical Uses
	demonstratePracticalUses()

	// بخش 11: Common Mistakes
	demonstrateCommonMistakes()

	// بخش 12: Quick Reference
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📚 time PACKAGE QUICK REFERENCE")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ TIME CREATION                                                  │")
	fmt.Println("├─────────────────────────────────────────────────────────────────┤")
	fmt.Println("│ time.Now()                - Current time                       │")
	fmt.Println("│ time.Date(year, month, day, hour, min, sec, nsec, loc)        │")
	fmt.Println("│ time.Parse(layout, value) - Parse string to time               │")
	fmt.Println("│ time.ParseInLocation(layout, value, loc) - Parse with location │")
	fmt.Println("│ time.Unix(sec, nsec)      - From Unix timestamp                │")
	fmt.Println("└─────────────────────────────────────────────────────────────────┘")

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ TIME OPERATIONS                                                │")
	fmt.Println("├─────────────────────────────────────────────────────────────────┤")
	fmt.Println("│ t.Add(duration)           - Add duration                       │")
	fmt.Println("│ t.AddDate(y, m, d)        - Add date components                │")
	fmt.Println("│ t.Sub(t2)                 - Difference between times           │")
	fmt.Println("│ t.Before(t2)              - Check if before                    │")
	fmt.Println("│ t.After(t2)               - Check if after                     │")
	fmt.Println("│ t.Equal(t2)               - Check equality                     │")
	fmt.Println("│ t.UTC()                   - Convert to UTC                     │")
	fmt.Println("│ t.In(loc)                 - Convert to location                │")
	fmt.Println("└─────────────────────────────────────────────────────────────────┘")

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ DURATION                                                       │")
	fmt.Println("├─────────────────────────────────────────────────────────────────┤")
	fmt.Println("│ time.Duration               - Type for time intervals          │")
	fmt.Println("│ time.Second, Minute, Hour   - Constants                        │")
	fmt.Println("│ time.ParseDuration(s)       - Parse string to duration         │")
	fmt.Println("│ d.Nanoseconds() / Microseconds() / Milliseconds() / Seconds()  │")
	fmt.Println("│ d.Minutes() / Hours()       - Convert to units                 │")
	fmt.Println("│ d.Round(d2) / Truncate(d2)  - Round/truncate duration          │")
	fmt.Println("└─────────────────────────────────────────────────────────────────┘")

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ TIMERS & TICKERS                                               │")
	fmt.Println("├─────────────────────────────────────────────────────────────────┤")
	fmt.Println("│ time.Sleep(d)              - Block for duration                │")
	fmt.Println("│ time.After(d)              - Channel that receives after d     │")
	fmt.Println("│ time.NewTimer(d)           - Create new timer                  │")
	fmt.Println("│ timer.Stop()               - Stop timer                        │")
	fmt.Println("│ timer.Reset(d)             - Reset timer                       │")
	fmt.Println("│ time.NewTicker(d)          - Create new ticker                 │")
	fmt.Println("│ ticker.Stop()              - Stop ticker                       │")
	fmt.Println("│ time.Tick(d)               - Channel of ticks (no stop)        │")
	fmt.Println("└─────────────────────────────────────────────────────────────────┘")

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ FORMAT LAYOUT (Reference: Mon Jan 2 15:04:05 MST 2006)         │")
	fmt.Println("├─────────────────────────────────────────────────────────────────┤")
	fmt.Println("│ 2006  - Year                                                   │")
	fmt.Println("│ 01    - Month (01, Jan, January)                               │")
	fmt.Println("│ 02    - Day                                                    │")
	fmt.Println("│ 15    - Hour (24-hour)                                         │")
	fmt.Println("│ 03    - Hour (12-hour)                                         │")
	fmt.Println("│ 04    - Minute                                                 │")
	fmt.Println("│ 05    - Second                                                 │")
	fmt.Println("│ PM    - AM/PM                                                  │")
	fmt.Println("│ Mon   - Weekday (Mon, Monday)                                  │")
	fmt.Println("│ MST   - Timezone                                               │")
	fmt.Println("│ -0700 - Timezone offset                                        │")
	fmt.Println("└─────────────────────────────────────────────────────────────────┘")

	fmt.Println("\n💡 GOLDEN RULES:")
	fmt.Println("  1. Always use time.Time for timestamps, not int64")
	fmt.Println("  2. Use time.Equal() for comparing times")
	fmt.Println("  3. Always stop Tickers to avoid goroutine leaks")
	fmt.Println("  4. Use time.Since() for measuring intervals")
	fmt.Println("  5. Use time.AfterFunc() for delayed function calls")
	fmt.Println("  6. Never use time.Sleep for synchronization")
	fmt.Println("  7. Use time.ParseDuration() for user input")
	fmt.Println("  8. Store times in UTC for consistency")
	fmt.Println("  9. Use monotonic clock for duration measurements")
	fmt.Println("  10. Check error when parsing times")

	fmt.Println("\n🎯 PERFORMANCE TIPS:")
	fmt.Println("  • Cache time.Now() if you need it multiple times")
	fmt.Println("  • Avoid time.After in loops (use Ticker)")
	fmt.Println("  • Use time.NewTimer instead of time.After when you need Stop")
	fmt.Println("  • Pre-allocate time.Time when possible")
	fmt.Println("  • Use time.Since(start) instead of time.Now().Sub(start)")
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