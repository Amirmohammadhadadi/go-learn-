// ============================================================================
// FILE: redis_complete_guide.go
// TITLE: راهنمای کامل Redis در Go - کش، Session، Rate Limiting، Pub/Sub
// HOW TO RUN: go run redis_complete_guide.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - Redis چیست و چه کاربردی دارد؟
// ============================================================================
//
// Redis (REmote DIctionary Server) یک دیتابیس in-memory کلید-مقدار است
// که به دلیل سرعت بالا و پشتیبانی از ساختارهای داده متنوع محبوب است.
//
// کاربردهای اصلی Redis:
// 1. کش کردن (Caching): کاهش بار دیتابیس، افزایش سرعت پاسخ
// 2. Session Storage: ذخیره نشست‌های کاربری
// 3. Rate Limiting: محدود کردن نرخ درخواست‌ها
// 4. Pub/Sub: پیام‌رسانی real-time بین سرویس‌ها
// 5. Queue: صف‌های پیام با List/Stream
// 6. Leaderboard: رتبه‌بندی با Sorted Set
// 7. Counter: شمارنده‌های با کارایی بالا
// 8. Lock: قفل‌های توزیع‌شده
//
// انواع داده در Redis:
// - String: ساده‌ترین نوع (شمارنده، کش)
// - Hash: شبیه به map (ذخیره object)
// - List: لیست مرتب (صف، stack)
// - Set: مجموعه بدون ترتیب (tags, unique items)
// - Sorted Set: مجموعه مرتب با نمره (leaderboard)
// - Bitmap: عملیات بیتی
// - HyperLogLog: تخمین تعداد unique
// - Geo: موقعیت جغرافیایی
//
// قانون طلایی:
// "از Redis برای داده‌های موقت و پرسرعت استفاده کن.
//  همیشه TTL تنظیم کن تا حافظه آزاد شود.
//  از Pipeline برای عملیات دسته‌ای استفاده کن.
//  هرگز داده حیاتی را فقط در Redis ذخیره نکن."
// ============================================================================

package __redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

// ============================================================================
// بخش 1: نصب و راه‌اندازی
// ============================================================================

/*
نصب:
$ go get github.com/go-redis/redis/v8

# اجرای Redis با Docker
$ docker run -d --name redis -p 6379:6379 redis:alpine
$ docker run -d --name redis -p 6379:6379 redis:7-alpine
*/

// RedisClient ساختار کلاینت Redis
type RedisClient struct {
	client *redis.Client
	ctx    context.Context
}

// NewRedisClient ایجاد کلاینت جدید
func NewRedisClient(addr, password string, db int) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           db,
		PoolSize:     10, // حداکثر اتصالات همزمان
		MinIdleConns: 5,  // حداقل اتصالات idle
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolTimeout:  4 * time.Second,
	})

	ctx := context.Background()

	// تست اتصال
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Println("Redis connected successfully")
	return &RedisClient{client: client, ctx: ctx}, nil
}

// Close بستن اتصال
func (r *RedisClient) Close() error {
	return r.client.Close()
}

// ============================================================================
// بخش 2: کش کردن (Caching)
// ============================================================================

// CacheService سرویس کش
type CacheService struct {
	client *RedisClient
}

// NewCacheService ایجاد سرویس کش
func NewCacheService(client *RedisClient) *CacheService {
	return &CacheService{client: client}
}

// Set کش کردن داده با TTL
func (c *CacheService) Set(key string, value interface{}, ttl time.Duration) error {
	// تبدیل به JSON
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	return c.client.client.Set(c.client.ctx, key, data, ttl).Err()
}

// Get دریافت داده از کش
func (c *CacheService) Get(key string, dest interface{}) error {
	data, err := c.client.client.Get(c.client.ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("cache miss: key not found")
		}
		return fmt.Errorf("failed to get from cache: %w", err)
	}

	if err := json.Unmarshal(data, dest); err != nil {
		return fmt.Errorf("failed to unmarshal data: %w", err)
	}

	return nil
}

// Delete حذف از کش
func (c *CacheService) Delete(key string) error {
	return c.client.client.Del(c.client.ctx, key).Err()
}

// Exists بررسی وجود کلید
func (c *CacheService) Exists(key string) (bool, error) {
	result, err := c.client.client.Exists(c.client.ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

// GetOrSet دریافت یا ذخیره (اگر不存在)
func (c *CacheService) GetOrSet(key string, ttl time.Duration, loader func() (interface{}, error), dest interface{}) error {
	// تلاش برای دریافت از کش
	err := c.Get(key, dest)
	if err == nil {
		return nil // کش hit
	}

	if err.Error() != "cache miss: key not found" {
		return err
	}

	// کش miss - بارگذاری از دیتابیس
	data, err := loader()
	if err != nil {
		return fmt.Errorf("failed to load data: %w", err)
	}

	// ذخیره در کش
	if err := c.Set(key, data, ttl); err != nil {
		log.Printf("failed to set cache: %v", err)
	}

	// کپی داده به destination
	dataBytes, _ := json.Marshal(data)
	return json.Unmarshal(dataBytes, dest)
}

// Increment افزایش شمارنده (برای کش کردن شمارنده‌ها)
func (c *CacheService) Increment(key string) (int64, error) {
	return c.client.client.Incr(c.client.ctx, key).Result()
}

// Decrement کاهش شمارنده
func (c *CacheService) Decrement(key string) (int64, error) {
	return c.client.client.Decr(c.client.ctx, key).Result()
}

// IncrementBy افزایش با مقدار مشخص
func (c *CacheService) IncrementBy(key string, value int64) (int64, error) {
	return c.client.client.IncrBy(c.client.ctx, key, value).Result()
}

// ============================================================================
// بخش 3: Session Storage
// ============================================================================

// Session ساختار نشست کاربری
type Session struct {
	ID        string                 `json:"id"`
	UserID    int                    `json:"user_id"`
	Username  string                 `json:"username"`
	Role      string                 `json:"role"`
	CreatedAt time.Time              `json:"created_at"`
	ExpiresAt time.Time              `json:"expires_at"`
	Data      map[string]interface{} `json:"data"`
}

// SessionService سرویس مدیریت نشست
type SessionService struct {
	client     *RedisClient
	prefix     string
	sessionTTL time.Duration
}

// NewSessionService ایجاد سرویس نشست
func NewSessionService(client *RedisClient, sessionTTL time.Duration) *SessionService {
	return &SessionService{
		client:     client,
		prefix:     "session:",
		sessionTTL: sessionTTL,
	}
}

// CreateSession ایجاد نشست جدید
func (s *SessionService) CreateSession(userID int, username, role string) (*Session, error) {
	sessionID := generateSessionID()
	now := time.Now()

	session := &Session{
		ID:        sessionID,
		UserID:    userID,
		Username:  username,
		Role:      role,
		CreatedAt: now,
		ExpiresAt: now.Add(s.sessionTTL),
		Data:      make(map[string]interface{}),
	}

	key := s.prefix + sessionID
	data, err := json.Marshal(session)
	if err != nil {
		return nil, err
	}

	if err := s.client.client.Set(s.client.ctx, key, data, s.sessionTTL).Err(); err != nil {
		return nil, err
	}

	return session, nil
}

// GetSession دریافت نشست
func (s *SessionService) GetSession(sessionID string) (*Session, error) {
	key := s.prefix + sessionID
	data, err := s.client.client.Get(s.client.ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("session not found")
		}
		return nil, err
	}

	var session Session
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, err
	}

	// بررسی انقضا
	if time.Now().After(session.ExpiresAt) {
		s.DeleteSession(sessionID)
		return nil, fmt.Errorf("session expired")
	}

	// تمدید TTL
	s.client.client.Expire(s.client.ctx, key, s.sessionTTL)

	return &session, nil
}

// UpdateSession به‌روزرسانی نشست
func (s *SessionService) UpdateSession(session *Session) error {
	key := s.prefix + session.ID
	data, err := json.Marshal(session)
	if err != nil {
		return err
	}

	return s.client.client.Set(s.client.ctx, key, data, s.sessionTTL).Err()
}

// DeleteSession حذف نشست
func (s *SessionService) DeleteSession(sessionID string) error {
	key := s.prefix + sessionID
	return s.client.client.Del(s.client.ctx, key).Err()
}

// SetSessionData ذخیره داده در نشست
func (s *SessionService) SetSessionData(sessionID, key string, value interface{}) error {
	session, err := s.GetSession(sessionID)
	if err != nil {
		return err
	}

	session.Data[key] = value
	return s.UpdateSession(session)
}

// GetSessionData دریافت داده از نشست
func (s *SessionService) GetSessionData(sessionID, key string) (interface{}, error) {
	session, err := s.GetSession(sessionID)
	if err != nil {
		return nil, err
	}

	value, ok := session.Data[key]
	if !ok {
		return nil, fmt.Errorf("key not found in session")
	}
	return value, nil
}

// GetUserSessions دریافت همه نشست‌های کاربر (با الگو)
func (s *SessionService) GetUserSessions(userID int) ([]Session, error) {
	pattern := s.prefix + "*"
	keys, err := s.client.client.Keys(s.client.ctx, pattern).Result()
	if err != nil {
		return nil, err
	}

	var sessions []Session
	for _, key := range keys {
		data, err := s.client.client.Get(s.client.ctx, key).Bytes()
		if err != nil {
			continue
		}

		var session Session
		if err := json.Unmarshal(data, &session); err != nil {
			continue
		}

		if session.UserID == userID {
			sessions = append(sessions, session)
		}
	}

	return sessions, nil
}

// ClearExpiredSessions پاک کردن نشست‌های منقضی (Redis خودکار انجام می‌دهد)
func (s *SessionService) ClearExpiredSessions() error {
	// Redis به صورت خودکار TTL را مدیریت می‌کند
	return nil
}

// ============================================================================
// بخش 4: Rate Limiting
// ============================================================================

// RateLimiter ساختار محدودکننده نرخ
type RateLimiter struct {
	client *RedisClient
}

// NewRateLimiter ایجاد محدودکننده نرخ جدید
func NewRateLimiter(client *RedisClient) *RateLimiter {
	return &RateLimiter{client: client}
}

// FixedWindowLimiter محدودیت با پنجره ثابت
func (r *RateLimiter) FixedWindowLimiter(key string, limit int, window time.Duration) (bool, int64, error) {
	// کلید منحصر به فرد برای پنجره زمانی
	windowKey := fmt.Sprintf("rate_limit:%s:%d", key, time.Now().Unix()/int64(window.Seconds()))

	// افزایش شمارنده
	count, err := r.client.client.Incr(r.client.ctx, windowKey).Result()
	if err != nil {
		return false, 0, err
	}

	// تنظیم TTL برای اولین بار
	if count == 1 {
		r.client.client.Expire(r.client.ctx, windowKey, window)
	}

	allowed := count <= int64(limit)
	remaining := int64(limit) - count
	if remaining < 0 {
		remaining = 0
	}

	return allowed, remaining, nil
}

// SlidingWindowLimiter محدودیت با پنجره لغزنده
func (r *RateLimiter) SlidingWindowLimiter(key string, limit int, window time.Duration) (bool, int64, error) {
	now := time.Now().UnixNano()
	windowKey := fmt.Sprintf("rate_limit:%s", key)

	// افزودن درخواست فعلی به Sorted Set با نمره = زمان فعلی
	member := fmt.Sprintf("%d", now)
	err := r.client.client.ZAdd(r.client.ctx, windowKey, &redis.Z{
		Score:  float64(now),
		Member: member,
	}).Err()
	if err != nil {
		return false, 0, err
	}

	// حذف درخواست‌های خارج از پنجره زمانی
	minScore := float64(now - window.Nanoseconds())
	err = r.client.client.ZRemRangeByScore(r.client.ctx, windowKey, "0", fmt.Sprintf("%f", minScore)).Err()
	if err != nil {
		return false, 0, err
	}

	// شمارش درخواست‌های داخل پنجره
	count, err := r.client.client.ZCard(r.client.ctx, windowKey).Result()
	if err != nil {
		return false, 0, err
	}

	// تنظیم TTL برای کلید (برای جلوگیری از堆积)
	r.client.client.Expire(r.client.ctx, windowKey, window)

	allowed := count <= int64(limit)
	remaining := int64(limit) - count
	if remaining < 0 {
		remaining = 0
	}

	return allowed, remaining, nil
}

// TokenBucketLimiter محدودیت با سطل توکن
func (r *RateLimiter) TokenBucketLimiter(key string, capacity int, refillRate float64) (bool, error) {
	luaScript := redis.NewScript(`
        local key = KEYS[1]
        local capacity = tonumber(ARGV[1])
        local refill_rate = tonumber(ARGV[2])
        local now = tonumber(ARGV[3])
        
        local bucket = redis.call('HMGET', key, 'tokens', 'last_refill')
        local tokens = tonumber(bucket[1]) or capacity
        local last_refill = tonumber(bucket[2]) or now
        
        local refill = (now - last_refill) * refill_rate
        tokens = math.min(capacity, tokens + refill)
        
        local allowed = tokens >= 1
        if allowed then
            tokens = tokens - 1
        end
        
        redis.call('HMSET', key, 'tokens', tokens, 'last_refill', now)
        redis.call('EXPIRE', key, 60)
        
        return allowed and 1 or 0
    `)

	result, err := luaScript.Run(r.client.ctx, r.client.client, []string{key}, capacity, refillRate, time.Now().Unix()).Int()
	if err != nil {
		return false, err
	}

	return result == 1, nil
}

// RateLimitMiddleware میدلور rate limiting برای HTTP
func (r *RateLimiter) RateLimitMiddleware(next http.HandlerFunc, limit int, window time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// استفاده از IP کلاینت به عنوان کلید
		key := req.RemoteAddr

		allowed, remaining, err := r.FixedWindowLimiter(key, limit, window)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", limit))
		w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))

		if !allowed {
			w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(window).Unix()))
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		next(w, req)
	}
}

// ============================================================================
// بخش 5: Pub/Sub (Publish/Subscribe)
// ============================================================================

// Message ساختار پیام
type Message struct {
	ID        string      `json:"id"`
	Channel   string      `json:"channel"`
	Type      string      `json:"type"`
	Payload   interface{} `json:"payload"`
	Timestamp time.Time   `json:"timestamp"`
	Sender    string      `json:"sender"`
}

// PubSubService سرویس Pub/Sub
type PubSubService struct {
	client      *RedisClient
	pubsub      *redis.PubSub
	subscribers map[string]chan Message
	mu          sync.RWMutex
}

// NewPubSubService ایجاد سرویس Pub/Sub جدید
func NewPubSubService(client *RedisClient) *PubSubService {
	return &PubSubService{
		client:      client,
		subscribers: make(map[string]chan Message),
	}
}

// Publish ارسال پیام به کانال
func (p *PubSubService) Publish(channel string, msg Message) error {
	msg.Timestamp = time.Now()
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return p.client.client.Publish(p.client.ctx, channel, data).Err()
}

// Subscribe اشتراک در کانال
func (p *PubSubService) Subscribe(channel string) (<-chan Message, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// اگر اولین مشترک است، PubSub را ایجاد کن
	if p.pubsub == nil {
		p.pubsub = p.client.client.Subscribe(p.client.ctx)
	}

	// اشتراک در کانال
	if err := p.pubsub.Subscribe(p.client.ctx, channel); err != nil {
		return nil, err
	}

	// کانال برای ارسال پیام‌ها
	msgChan := make(chan Message, 100)
	p.subscribers[channel] = msgChan

	// شروع listening در گوروتین جدا
	go p.listen(channel)

	return msgChan, nil
}

// listen گوش دادن به پیام‌های کانال
func (p *PubSubService) listen(channel string) {
	ch := p.pubsub.Channel()

	for msg := range ch {
		if msg.Channel == channel {
			var message Message
			if err := json.Unmarshal([]byte(msg.Payload), &message); err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				continue
			}

			p.mu.RLock()
			if ch, ok := p.subscribers[channel]; ok {
				select {
				case ch <- message:
				default:
					log.Printf("Channel %s is full, dropping message", channel)
				}
			}
			p.mu.RUnlock()
		}
	}
}

// Unsubscribe لغو اشتراک
func (p *PubSubService) Unsubscribe(channel string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if err := p.pubsub.Unsubscribe(p.client.ctx, channel); err != nil {
		return err
	}

	if ch, ok := p.subscribers[channel]; ok {
		close(ch)
		delete(p.subscribers, channel)
	}

	return nil
}

// Close بستن PubSub
func (p *PubSubService) Close() error {
	if p.pubsub != nil {
		return p.pubsub.Close()
	}
	return nil
}

// ============================================================================
// بخش 6: کاربردهای پیشرفته
// ============================================================================

// 6.1 Distributed Lock (قفل توزیع‌شده)
type DistributedLock struct {
	client *RedisClient
	key    string
	value  string
	ttl    time.Duration
}

// NewDistributedLock ایجاد قفل توزیع‌شده
func NewDistributedLock(client *RedisClient, key string, ttl time.Duration) *DistributedLock {
	return &DistributedLock{
		client: client,
		key:    "lock:" + key,
		value:  generateLockValue(),
		ttl:    ttl,
	}
}

// Acquire تلاش برای گرفتن قفل
func (l *DistributedLock) Acquire() (bool, error) {
	// SET NX (Not eXists) + EX (Expire)
	ok, err := l.client.client.SetNX(l.client.ctx, l.key, l.value, l.ttl).Result()
	if err != nil {
		return false, err
	}
	return ok, nil
}

// Release آزاد کردن قفل
func (l *DistributedLock) Release() error {
	luaScript := redis.NewScript(`
        if redis.call("get", KEYS[1]) == ARGV[1] then
            return redis.call("del", KEYS[1])
        else
            return 0
        end
    `)

	result, err := luaScript.Run(l.client.ctx, l.client.client, []string{l.key}, l.value).Int()
	if err != nil {
		return err
	}
	if result == 0 {
		return fmt.Errorf("lock not owned")
	}
	return nil
}

// 6.2 Leaderboard (رتبه‌بندی با Sorted Set)
type Leaderboard struct {
	client *RedisClient
	key    string
}

func NewLeaderboard(client *RedisClient, name string) *Leaderboard {
	return &Leaderboard{
		client: client,
		key:    "leaderboard:" + name,
	}
}

// AddScore افزودن امتیاز
func (l *Leaderboard) AddScore(userID string, score float64) error {
	return l.client.client.ZIncrBy(l.client.ctx, l.key, score, userID).Err()
}

// GetTopUsers دریافت کاربران برتر
func (l *Leaderboard) GetTopUsers(limit int64) ([]redis.Z, error) {
	return l.client.client.ZRevRangeWithScores(l.client.ctx, l.key, 0, limit-1).Result()
}

// GetUserRank دریافت رتبه کاربر
func (l *Leaderboard) GetUserRank(userID string) (int64, error) {
	return l.client.client.ZRevRank(l.client.ctx, l.key, userID).Result()
}

// GetUserScore دریافت امتیاز کاربر
func (l *Leaderboard) GetUserScore(userID string) (float64, error) {
	return l.client.client.ZScore(l.client.ctx, l.key, userID).Result()
}

// 6.3 Queue (صف با List)
type Queue struct {
	client *RedisClient
	key    string
}

func NewQueue(client *RedisClient, name string) *Queue {
	return &Queue{
		client: client,
		key:    "queue:" + name,
	}
}

// Push افزودن به صف
func (q *Queue) Push(value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return q.client.client.LPush(q.client.ctx, q.key, data).Err()
}

// Pop برداشتن از صف (بلوکینگ)
func (q *Queue) Pop(timeout time.Duration) ([]byte, error) {
	result, err := q.client.client.BRPop(q.client.ctx, timeout, q.key).Result()
	if err != nil {
		return nil, err
	}
	if len(result) < 2 {
		return nil, fmt.Errorf("no data")
	}
	return []byte(result[1]), nil
}

// Length طول صف
func (q *Queue) Length() (int64, error) {
	return q.client.client.LLen(q.client.ctx, q.key).Result()
}

// ============================================================================
// بخش 7: Helper Functions
// ============================================================================

func generateSessionID() string {
	return fmt.Sprintf("%d_%d", time.Now().UnixNano(), rand.Int63())
}

func generateLockValue() string {
	return fmt.Sprintf("%d_%d", time.Now().UnixNano(), rand.Int63())
}

// ============================================================================
// بخش 8: Best Practices
// ============================================================================

func bestPractices() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("💡 REDIS BEST PRACTICES")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────┐
│ CACHING BEST PRACTICES                                        │
├─────────────────────────────────────────────────────────────────┤
│ • Always set TTL for cache keys                               │
│ • Use Cache-Aside pattern                                     │
│ • Implement cache invalidation strategy                       │
│ • Use connection pooling                                      │
│ • Monitor memory usage                                        │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ SESSION MANAGEMENT                                            │
├─────────────────────────────────────────────────────────────────┤
│ • Use short TTL for sessions (15-60 minutes)                  │
│ • Implement refresh mechanism                                 │
│ • Store minimal data in session                               │
│ • Use Redis Cluster for high availability                     │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ RATE LIMITING                                                 │
├─────────────────────────────────────────────────────────────────┤
│ • Use appropriate algorithm (Fixed Window, Sliding Window,    │
│   Token Bucket)                                               │
│ • Add rate limit headers to responses                         │
│ • Implement different limits for different endpoints          │
│ • Use Redis Lua scripts for atomicity                         │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ PUB/SUB                                                        │
├─────────────────────────────────────────────────────────────────┤
│ • Use channels for different message types                    │
│ • Implement retry mechanism for lost messages                 │
│ • Use Redis Streams for persistent messaging                  │
│ • Monitor subscriber count                                    │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ GENERAL                                                       │
├─────────────────────────────────────────────────────────────────┤
│ • Use Pipeline for batch operations                           │
│ • Use Lua scripts for atomic operations                       │
│ • Monitor slow queries                                        │
│ • Set appropriate maxmemory policy                            │
│ • Use Redis Sentinel or Cluster for production                │
└─────────────────────────────────────────────────────────────────┘
`)
}

// ============================================================================
// بخش 9: Main
// ============================================================================

func main() {
	fmt.Println(stringsRepeat("=", 80))
	fmt.Println("🎯 COMPLETE REDIS GUIDE")
	fmt.Println("Caching | Session | Rate Limiting | Pub/Sub")
	fmt.Println(stringsRepeat("=", 80))

	bestPractices()

	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📝 EXAMPLE USAGE")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println(`
// مثال کامل استفاده از Redis
func main() {
    // اتصال به Redis
    client, err := NewRedisClient("localhost:6379", "", 0)
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // 1. کش کردن
    cache := NewCacheService(client)
    
    type User struct {
        ID   int    ` + "`json:\"id\"`" + `
        Name string ` + "`json:\"name\"`" + `
    }
    
    // ذخیره در کش
    user := User{ID: 1, Name: "Ali"}
    cache.Set("user:1", user, 10*time.Minute)
    
    // دریافت از کش
    var cachedUser User
    if err := cache.Get("user:1", &cachedUser); err == nil {
        fmt.Printf("Cached user: %+v\n", cachedUser)
    }
    
    // GetOrSet
    var user2 User
    cache.GetOrSet("user:2", 10*time.Minute, func() (interface{}, error) {
        // بارگذاری از دیتابیس
        return User{ID: 2, Name: "Sara"}, nil
    }, &user2)

    // 2. Session Storage
    sessionService := NewSessionService(client, 30*time.Minute)
    
    // ایجاد نشست
    session, err := sessionService.CreateSession(1, "ali", "user")
    if err != nil {
        log.Fatal(err)
    }
    
    // ذخیره داده در نشست
    sessionService.SetSessionData(session.ID, "cart", []int{1,2,3})

    // 3. Rate Limiting
    limiter := NewRateLimiter(client)
    
    // Fixed Window
    allowed, remaining, _ := limiter.FixedWindowLimiter("api:login", 10, time.Minute)
    
    // Token Bucket
    allowed, _ = limiter.TokenBucketLimiter("api:search", 100, 10)

    // 4. Pub/Sub
    pubsub := NewPubSubService(client)
    
    // اشتراک
    msgChan, _ := pubsub.Subscribe("notifications")
    go func() {
        for msg := range msgChan {
            fmt.Printf("Received: %+v\n", msg)
        }
    }()
    
    // انتشار
    pubsub.Publish("notifications", Message{
        Type:    "user_joined",
        Payload: map[string]string{"user": "Ali"},
    })
}
`)

	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🎯 REDIS GUIDE - COMPLETE")
	fmt.Println("Ready to build high-performance applications with Redis!")
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
