package main

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

// ============================================================================
// 1. 令牌桶限流器 (Token Bucket RateLimiter)
// ============================================================================

// TokenBucketLimiter 生产级令牌桶限流器
type TokenBucketLimiter struct {
	rate       float64    // 放入令牌速率 (tokens/sec)
	capacity   float64    // 桶容量
	tokens     float64    // 当前剩余令牌
	lastUpdate time.Time  // 上次补充令牌时间
	mu         sync.Mutex // 保护互斥锁
}

func NewTokenBucketLimiter(rate, capacity float64) *TokenBucketLimiter {
	return &TokenBucketLimiter{
		rate:       rate,
		capacity:   capacity,
		tokens:     capacity, // 初始装满
		lastUpdate: time.Now(),
	}
}

// Allow 判断是否允许一个请求通过 (非阻塞)
func (tb *TokenBucketLimiter) Allow() bool {
	return tb.AllowN(time.Now(), 1)
}

func (tb *TokenBucketLimiter) AllowN(now time.Time, n float64) bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	// 计算时间差并补充令牌
	elapsed := now.Sub(tb.lastUpdate).Seconds()
	tb.lastUpdate = now
	tb.tokens = math.Min(tb.capacity, tb.tokens+elapsed*tb.rate)

	if tb.tokens >= n {
		tb.tokens -= n
		return true
	}
	return false
}

// ============================================================================
// 2. 自适应熔断器 (Circuit Breaker - 状态机模式)
// ============================================================================

type State int

const (
	StateClosed   State = iota // 正常放行
	StateHalfOpen              // 半开探测
	StateOpen                  // 熔断拦截
)

func (s State) String() string {
	switch s {
	case StateClosed:
		return "CLOSED"
	case StateHalfOpen:
		return "HALF-OPEN"
	case StateOpen:
		return "OPEN"
	default:
		return "UNKNOWN"
	}
}

var ErrCircuitBreakerOpen = errors.New("circuit breaker is open: request rejected")

type CircuitBreaker struct {
	state             State
	failureThreshold  int           // 熔断错误阈值
	consecutiveFails  int           // 当前连续错误数
	openTimeout       time.Duration // 熔断后多久进入半开探测
	lastStateChange   time.Time
	halfOpenSuccesses int
	mu                sync.Mutex
}

func NewCircuitBreaker(failureThreshold int, openTimeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		state:            StateClosed,
		failureThreshold: failureThreshold,
		openTimeout:      openTimeout,
		lastStateChange:  time.Now(),
	}
}

// Execute 封装受保护的调用
func (cb *CircuitBreaker) Execute(fn func() error) error {
	cb.mu.Lock()
	now := time.Now()

	// 状态转移检查：若处于 OPEN 状态且超过 openTimeout，转为 HALF-OPEN
	if cb.state == StateOpen && now.Sub(cb.lastStateChange) > cb.openTimeout {
		cb.state = StateHalfOpen
		cb.lastStateChange = now
		cb.halfOpenSuccesses = 0
	}

	if cb.state == StateOpen {
		cb.mu.Unlock()
		return ErrCircuitBreakerOpen
	}

	cb.mu.Unlock()

	// 执行受保护的操作
	err := fn()

	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err != nil {
		// 调用失败
		cb.consecutiveFails++
		if cb.state == StateHalfOpen || cb.consecutiveFails >= cb.failureThreshold {
			cb.state = StateOpen
			cb.lastStateChange = time.Now()
		}
		return err
	}

	// 调用成功
	if cb.state == StateHalfOpen {
		cb.halfOpenSuccesses++
		if cb.halfOpenSuccesses >= 2 { // 连续两次探测成功即转为 CLOSED
			cb.state = StateClosed
			cb.consecutiveFails = 0
			cb.lastStateChange = time.Now()
		}
	} else if cb.state == StateClosed {
		cb.consecutiveFails = 0
	}
	return nil
}

func (cb *CircuitBreaker) CurrentState() State {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.state
}

// ============================================================================
// 3. 业务幂等防重机制 (Idempotency Token Manager)
// ============================================================================

type IdempotencyRecord struct {
	CreatedAt time.Time
	Result    string
	Completed bool
}

type IdempotencyManager struct {
	store sync.Map
	ttl   time.Duration
}

func NewIdempotencyManager(ttl time.Duration) *IdempotencyManager {
	return &IdempotencyManager{ttl: ttl}
}

var ErrConcurrentDuplicate = errors.New("request already in progress")

// ExecuteIdempotent 执行幂等保护调用
func (im *IdempotencyManager) ExecuteIdempotent(idempotencyKey string, action func() (string, error)) (string, error) {
	record := &IdempotencyRecord{
		CreatedAt: time.Now(),
		Completed: false,
	}

	// 尝试抢占 key
	actual, loaded := im.store.LoadOrStore(idempotencyKey, record)
	if loaded {
		existing := actual.(*IdempotencyRecord)
		if !existing.Completed {
			return "", ErrConcurrentDuplicate
		}
		// 返回已持久化的既往结果
		return existing.Result, nil
	}

	// 抢占成功，执行业务操作
	res, err := action()
	if err != nil {
		im.store.Delete(idempotencyKey) // 失败允许重试
		return "", err
	}

	record.Result = res
	record.Completed = true
	return res, nil
}

func main() {
	fmt.Println("=== 阶段三 专题02：微服务高可用三板斧（限流、熔断、幂等） ===")

	// 1. 令牌桶限流演示
	limiter := NewTokenBucketLimiter(10, 5) // 10/s, 桶容量 5
	fmt.Println("[1. 令牌桶限流实测]")
	passed := 0
	for i := 0; i < 8; i++ {
		if limiter.Allow() {
			passed++
		}
	}
	fmt.Printf("瞬时并发 8 次请求，成功放行: %d，限流丢弃: %d\n", passed, 8-passed)

	// 2. 熔断器状态机演示
	cb := NewCircuitBreaker(3, 100*time.Millisecond)
	fmt.Println("\n[2. 熔断降级实测]")
	// 连续触发 3 次失败
	for i := 0; i < 3; i++ {
		_ = cb.Execute(func() error {
			return errors.New("upstream service error")
		})
	}
	fmt.Printf("连续 3 次失败后，熔断器状态: %s\n", cb.CurrentState())

	// 熔断期拒绝请求
	errReject := cb.Execute(func() error { return nil })
	fmt.Printf("熔断拦截结果: %v\n", errReject)

	// 等待探测窗口
	time.Sleep(120 * time.Millisecond)
	fmt.Println("等待 120ms 超时恢复，发起半开探测...")
	_ = cb.Execute(func() error { return nil })
	_ = cb.Execute(func() error { return nil })
	fmt.Printf("连续探测成功，熔断器状态重置为: %s\n", cb.CurrentState())

	// 3. 业务幂等性演示
	im := NewIdempotencyManager(5 * time.Minute)
	fmt.Println("\n[3. 业务幂等防重实测]")
	token := "ORDER_REQ_TOKEN_9988"
	var orderCounter atomic.Int32

	createOrder := func() (string, error) {
		cnt := orderCounter.Add(1)
		return fmt.Sprintf("OrderCreated_No_%d", cnt), nil
	}

	res1, _ := im.ExecuteIdempotent(token, createOrder)
	res2, _ := im.ExecuteIdempotent(token, createOrder)
	fmt.Printf("首次提交返回: %s\n", res1)
	fmt.Printf("重复提交返回: %s (底层执行次数计数器: %d)\n", res2, orderCounter.Load())
	_ = context.Background()
}
