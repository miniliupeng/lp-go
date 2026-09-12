package main

import (
	"bytes"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
)

// 1. 模拟与探究 Channel 底层 hchan 的行为特征
// 展示：从已关闭通道无锁读取剩余数据与零值的 fast path 机制
func demoChannelInternalFastPath() {
	ch := make(chan int, 3)
	ch <- 10
	ch <- 20
	close(ch)

	// fast path: 缓冲区有数据，仍能正常读取
	v1, ok1 := <-ch
	v2, ok2 := <-ch
	// 缓冲区排空后读取，立即返回零值与 false (无需阻塞)
	v3, ok3 := <-ch

	fmt.Printf("[Channel fast path] 读取1: %d (ok=%v), 读取2: %d (ok=%v), 读取3: %d (ok=%v)\n",
		v1, ok1, v2, ok2, v3, ok3)
}

// 2. 模拟 sync.Mutex 正常模式与饥饿模式的竞争场景
// 当大量新协程持续竞争锁时，若无饥饿模式，尾部排队协程会被饿死
func demoMutexStarvationConcept() {
	var mu sync.Mutex
	var counter int64

	// 启动一个长耗时持有锁的协程，模拟持续排队
	var wg sync.WaitGroup
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				mu.Lock()
				atomic.AddInt64(&counter, 1)
				// 极短计算模拟临界区
				runtime.Gosched()
				mu.Unlock()
			}
		}(i)
	}

	wg.Wait()
	fmt.Printf("[Mutex 竞争防护] 10 个协程并发完成累加，总计计数: %d (无死锁与数据竞争)\n", counter)
}

// 3. sync.Pool 高并发对象复用与双层缓存机制
// 定义一个全局字节缓冲池
var bufferPool = sync.Pool{
	New: func() any {
		// New 函数在池中无可用对象时兜底创建
		return new(bytes.Buffer)
	},
}

func demoSyncPoolLifecycle() {
	// 1. 从 Pool 获取对象 (命中 localPool 或 New)
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	buf.WriteString("Hello sync.Pool")
	fmt.Printf("[sync.Pool 获取] 数据: %s, 初始容量: %d\n", buf.String(), buf.Cap())

	// 2. 用完后放回 Pool (供本 P 或通过 victim cache 复用)
	buf.Reset()
	bufferPool.Put(buf)

	// 3. 再次获取，验证复用到了同一个指针
	buf2 := bufferPool.Get().(*bytes.Buffer)
	fmt.Printf("[sync.Pool 复用] 再次获取成功，指针地址是否相同: %v\n", buf == buf2)
	bufferPool.Put(buf2)
}

// 4. sync.Once 现代双重检查锁 (DCL) 与 Fast Path 无锁读取
type DatabaseConn struct {
	Endpoint string
}

var (
	dbInstance *DatabaseConn
	once       sync.Once
	initCount  int64
)

func GetDatabaseInstance() *DatabaseConn {
	// once.Do 底层: if atomic.LoadUint32(&o.done) == 0 { o.doSlow(f) }
	// 一旦初始化完成，后续所有并发读取全走原子加载 Fast Path，0 锁损耗！
	once.Do(func() {
		atomic.AddInt64(&initCount, 1)
		dbInstance = &DatabaseConn{Endpoint: "10.0.0.1:3306"}
	})
	return dbInstance
}

func demoSyncOnceFastPath() {
	var wg sync.WaitGroup
	// 启动 50 个并发协程获取单例
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = GetDatabaseInstance()
		}()
	}
	wg.Wait()
	fmt.Printf("[sync.Once 单例] 50 个并发协程请求完成，初始化函数实际执行次数: %d (确保且仅执行 1 次)\n", initCount)
}

// 5. sync.Map 读写分离架构 (read map + dirty map)
func demoSyncMapReadDirty() {
	var m sync.Map

	// 1. 写入键值对 (首次写入 dirty map)
	m.Store("config_env", "production")
	m.Store("max_conns", 5000)

	// 2. 高并发无锁读取 (命中了 read map 后的原子指针读取，无需任何 Mutex)
	val, ok := m.Load("config_env")
	fmt.Printf("[sync.Map 读写分离] Load 读取 config_env: %v (命中: %v)\n", val, ok)

	// 3. LoadOrStore 原子存在性检查与设值
	actual, loaded := m.LoadOrStore("max_conns", 9999)
	fmt.Printf("[sync.Map LoadOrStore] 读取已有值: %v (是否已存在: %v)\n", actual, loaded)

	// 4. 遍历与安全删除
	m.Delete("max_conns")
	_, exists := m.Load("max_conns")
	fmt.Printf("[sync.Map 删除验证] max_conns 删除后是否存在: %v\n", exists)
}

func main() {
	fmt.Println("=== 1. Channel 状态与 hchan 缓冲区读取验证 ===")
	demoChannelInternalFastPath()

	fmt.Println("\n=== 2. sync.Mutex 状态流转与并发临界区安全 ===")
	demoMutexStarvationConcept()

	fmt.Println("\n=== 3. sync.Pool 高并发对象复用生命周期 ===")
	demoSyncPoolLifecycle()

	fmt.Println("\n=== 4. sync.Once 现代双重检查锁与零锁 Fast Path ===")
	demoSyncOnceFastPath()

	fmt.Println("\n=== 5. sync.Map 读写分离与高并发无锁读取 ===")
	demoSyncMapReadDirty()
}
