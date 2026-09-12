package main

import (
	"testing"
)

// TestMVCCVisibilityRules 验证 MVCC ReadView 的判定完整性
func TestMVCCVisibilityRules(t *testing.T) {
	rv := ReadView{
		CreatorTrxID: 100,
		MIDs:         []int64{80, 90},
		MinTrxID:     80,
		MaxTrxID:     101,
	}

	// 1. 自身事务修改必定可见
	if !rv.IsVisible(100) {
		t.Fatal("creator trx should be visible")
	}

	// 2. < MinTrxID 必定可见
	if !rv.IsVisible(70) {
		t.Fatal("trx < min_trx_id should be visible")
	}

	// 3. >= MaxTrxID 必定不可见
	if rv.IsVisible(101) || rv.IsVisible(200) {
		t.Fatal("trx >= max_trx_id should not be visible")
	}

	// 4. MIDs 中的未提交事务不可见
	if rv.IsVisible(80) || rv.IsVisible(90) {
		t.Fatal("uncommitted active trx should not be visible")
	}

	// 5. 在区间内但不在 MIDs 中 (已提交) 可见
	if !rv.IsVisible(85) {
		t.Fatal("committed trx in range should be visible")
	}
}

// 1. MVCC 快照读判定基准测试 (超高吞吐验证)
func BenchmarkMVCCVisibilityCheck(b *testing.B) {
	rv := ReadView{
		CreatorTrxID: 500,
		MIDs:         []int64{200, 300, 400},
		MinTrxID:     200,
		MaxTrxID:     600,
	}

	b.ReportAllocs()
	for b.Loop() {
		// 覆盖可见与不可见判定分支
		_ = rv.IsVisible(150)
		_ = rv.IsVisible(300)
		_ = rv.IsVisible(500)
		_ = rv.IsVisible(700)
	}
}

// 2. 连接池参数合法性测试
func TestDBPoolConfigValidity(t *testing.T) {
	cfg := getProductionDBPoolConfig()
	if cfg.MaxIdleConns > cfg.MaxOpenConns {
		t.Fatal("MaxIdleConns should not exceed MaxOpenConns")
	}
	if cfg.ConnMaxLifetime <= 0 {
		t.Fatal("ConnMaxLifetime must be positive")
	}
}
