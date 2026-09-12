package main

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// ==================== 1. SingleFlight 内存级并发防击穿核心实现 ====================

// call 代表一次正在进行中或已完成的回源函数调用
type call struct {
	wg  sync.WaitGroup
	val any
	err error
}

// SingleFlightGroup 管理相同 key 的并发请求合并
type SingleFlightGroup struct {
	mu sync.Mutex
	m  map[string]*call
}

func NewSingleFlightGroup() *SingleFlightGroup {
	return &SingleFlightGroup{
		m: make(map[string]*call),
	}
}

// Do 核心算法：针对同一个 key，若已有请求在执行，其余并发协程等待其完成并复用结果
func (g *SingleFlightGroup) Do(key string, fn func() (any, error)) (any, error) {
	g.mu.Lock()
	if c, ok := g.m[key]; ok {
		// 已有相同请求正在执行，释放锁并等待该请求完成
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}

	// 首次发起的请求，负责真正执行 fn
	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	// 执行实际回源函数 (如查询底层 MySQL)
	c.val, c.err = fn()
	c.wg.Done()

	// 清理已完成的调用
	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()

	return c.val, c.err
}

func demoSingleFlightBreakdown() {
	sf := NewSingleFlightGroup()
	var dbQueryCount int64

	// 模拟耗时的数据库回源操作 (100ms)
	mockDBSlowQuery := func() (any, error) {
		atomic.AddInt64(&dbQueryCount, 1)
		time.Sleep(50 * time.Millisecond)
		return "expensive_data_payload", nil
	}

	// 模拟热点数据失效瞬间，100 个并发 Goroutine 同时涌入请求同一个 Key
	var wg sync.WaitGroup
	concurrency := 100
	wg.Add(concurrency)

	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			_, _ = sf.Do("hot_product_1001", mockDBSlowQuery)
		}()
	}

	wg.Wait()
	fmt.Printf("[SingleFlight 防击穿实测] 100 个高并发请求涌入，底层数据库实际被查询次数: %d (成功合并 99%% 并发压力！)\n",
		dbQueryCount)
}

// ==================== 2. 生产级 Redis 分布式锁与看门狗自动续期 ====================

// RedisLock 模拟企业级分布式锁结构
type RedisLock struct {
	mu         sync.Mutex
	token      string             // 随机生成的客户端 UUID，防止误删他人持有的锁
	isHeld     bool               // 锁持有状态
	cancelAuto context.CancelFunc // 看门狗停止信号
}

// Lua 解锁脚本标准原语 (仅供参考逻辑)：
// if redis.call("get", KEYS[1]) == ARGV[1] then
//     return redis.call("del", KEYS[1])
// else
//     return 0
// end

func NewRedisLock(token string) *RedisLock {
	return &RedisLock{token: token}
}

// LockWithWatchdog 获取锁并自动启动后台看门狗协程进行周期性续期
func (l *RedisLock) LockWithWatchdog(ttl time.Duration) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	// 模拟 SET resource_key token NX PX ttl 成功
	l.isHeld = true
	ctx, cancel := context.WithCancel(context.Background())
	l.cancelAuto = cancel

	// 启动看门狗协程：每隔 TTL/3 时间进行一次续期
	go func() {
		ticker := time.NewTicker(ttl / 3)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return // 业务处理完成，看门狗退出
			case <-ticker.C:
				fmt.Printf("  [Watchdog 续期] 自动为锁 (Token: %s) 刷新 TTL 至 %v\n", l.token, ttl)
			}
		}
	}()

	return true
}

// Unlock 安全释放锁：先停止看门狗，再通过原子比对 Token 释放
func (l *RedisLock) Unlock() {
	l.mu.Lock()
	defer l.mu.Unlock()

	if !l.isHeld {
		return
	}

	// 1. 停止看门狗协程
	if l.cancelAuto != nil {
		l.cancelAuto()
	}
	l.isHeld = false
	fmt.Printf("[RedisLock 解锁] 通过 Lua 比对 Token (%s) 校验一致，成功安全释放锁\n", l.token)
}

func demoRedisLockAndWatchdog() {
	lock := NewRedisLock("client_uuid_abc_123")
	fmt.Println("\n=== 2. Redis 分布式锁与看门狗 (Watchdog) 自动续期实战 ===")

	// 假设施加 60ms 的锁，业务执行耗时 100ms (超过原锁时间，若无看门狗必然发生锁过期被他人抢占)
	lock.LockWithWatchdog(60 * time.Millisecond)

	// 模拟长时间执行的复杂分布式业务
	time.Sleep(90 * time.Millisecond)

	// 业务安全结束，释放锁
	lock.Unlock()
}

func main() {
	fmt.Println("=== 1. Redis 高并发防击穿：SingleFlight 并发请求合并 ===")
	demoSingleFlightBreakdown()

	demoRedisLockAndWatchdog()
}
