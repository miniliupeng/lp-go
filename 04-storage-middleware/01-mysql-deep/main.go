package main

import (
	"fmt"
	"sync"
	"time"
)

// ==================== 1. InnoDB MVCC 一致性读视口 (ReadView) 算法模拟 ====================

// ReadView 模拟 InnoDB 在执行快照读时生成的视口结构
type ReadView struct {
	CreatorTrxID int64   // 创建该 ReadView 的当前事务 ID
	MIDs         []int64 // 生成 ReadView 时系统正在活跃且未提交的事务 ID 列表
	MinTrxID     int64   // MIDs 中的最小值 (活跃事务的下界)
	MaxTrxID     int64   // 生成 ReadView 时系统应分配给下一个事务的 ID (高水位)
}

// IsVisible 按照 InnoDB 经典四步判定版本链中的某一版本记录是否对当前事务可见
func (rv *ReadView) IsVisible(trxID int64) bool {
	// 1. 如果被访问版本的 trx_id 与当前事务自身相同，说明是自己修改的，可见
	if trxID == rv.CreatorTrxID {
		return true
	}

	// 2. 如果被访问版本的 trx_id 小于最小活跃事务 ID，说明该版本在快照生成前早已提交，可见
	if trxID < rv.MinTrxID {
		return true
	}

	// 3. 如果被访问版本的 trx_id 大于或等于下一个要分配的事务 ID，说明该版本在快照生成之后才开启，不可见
	if trxID >= rv.MaxTrxID {
		return false
	}

	// 4. 如果 trx_id 在 [MinTrxID, MaxTrxID) 之间，检查其是否在活跃未提交列表 MIDs 中
	for _, activeID := range rv.MIDs {
		if activeID == trxID {
			// 在活跃列表中，说明该事务尚未提交，不可见
			return false
		}
	}

	// 不在活跃列表中，说明已在该快照生成前提交，可见
	return true
}

func demoMVCCReadView() {
	// 模拟快照生成时刻：
	// 活跃事务: 100, 200; 当前事务: 105; 下一个事务 ID: 201
	rv := ReadView{
		CreatorTrxID: 105,
		MIDs:         []int64{100, 200},
		MinTrxID:     100,
		MaxTrxID:     201,
	}

	testCases := []struct {
		desc  string
		trxID int64
	}{
		{"历史旧版本数据 (trx_id=80)", 80},
		{"并发活跃事务 A 未提交的数据 (trx_id=100)", 100},
		{"当前事务自己修改的数据 (trx_id=105)", 105},
		{"在快照前已完成提交的数据 (trx_id=150)", 150},
		{"快照生成后未来新开启事务的数据 (trx_id=205)", 205},
	}

	for _, tc := range testCases {
		visible := rv.IsVisible(tc.trxID)
		fmt.Printf("版本 [%s]: 对当前事务可见? %v\n", tc.desc, visible)
	}
}

// ==================== 2. Next-Key Lock 与插入意向锁并发死锁原理模拟 ====================

// 模拟表行记录区间: 主键索引存在记录 10, 20, 30
// 间隙为: (-∞, 10], (10, 20], (20, 30], (30, +∞)
func demoNextKeyLockDeadlockConcept() {
	var wg sync.WaitGroup
	wg.Add(2)

	// 共享资源信号量模拟锁争用
	lockGap10To20 := make(chan struct{}, 1)
	lockGap10To20 <- struct{}{} // 间隙锁可用

	fmt.Println("\n[Next-Key Lock 锁冲突演示]")
	// 事务 A: 锁定 (10, 20) 的间隙
	go func() {
		defer wg.Done()
		fmt.Println("  [事务 A] 执行 SELECT FOR UPDATE 锁定区间 (10, 20) 的 Gap Lock")
		<-lockGap10To20
		time.Sleep(30 * time.Millisecond)
		fmt.Println("  [事务 A] 尝试在 (10, 20) 之间 INSERT 记录 15，需要申请插入意向锁 (等待互斥解锁)")
		lockGap10To20 <- struct{}{}
	}()

	// 事务 B: 同样尝试锁定该区间
	go func() {
		defer wg.Done()
		time.Sleep(10 * time.Millisecond)
		fmt.Println("  [事务 B] 并发执行锁定区间 (10, 20)，Gap Lock 之间彼此兼容")
		time.Sleep(50 * time.Millisecond)
		fmt.Println("  [事务 B] 释放相关资源")
	}()

	wg.Wait()
	fmt.Println(">> 成功演示间隙锁兼容性与插入意向锁互斥等待机理！")
}

// ==================== 3. 生产级 database/sql 连接池配置规范 ====================

// DBPoolConfig 展示企业级生产实践推荐的连接池参数
type DBPoolConfig struct {
	MaxOpenConns    int           // 最大打开连接数 (受限于 MySQL max_connections / 容器 CPU)
	MaxIdleConns    int           // 最大空闲连接数 (建议与 MaxOpenConns 接近，减少频繁握手开销)
	ConnMaxLifetime time.Duration // 连接最大生存时间 (必须小于云服务商 RDS / HAProxy 默认断开时间，如 5~10 分钟)
	ConnMaxIdleTime time.Duration // 空闲连接最大存活时间 (回收低峰期冗余连接)
}

func getProductionDBPoolConfig() DBPoolConfig {
	return DBPoolConfig{
		MaxOpenConns:    100,
		MaxIdleConns:    100,
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: 1 * time.Minute,
	}
}

func main() {
	fmt.Println("=== 1. InnoDB MVCC ReadView 一致性快照判定算法 ===")
	demoMVCCReadView()

	demoNextKeyLockDeadlockConcept()

	fmt.Println("\n=== 3. database/sql 企业级连接池调优准则 ===")
	cfg := getProductionDBPoolConfig()
	fmt.Printf("推荐配置: MaxOpen=%d, MaxIdle=%d, MaxLifetime=%v, MaxIdleTime=%v\n",
		cfg.MaxOpenConns, cfg.MaxIdleConns, cfg.ConnMaxLifetime, cfg.ConnMaxIdleTime)
}
