package main

import (
	"fmt"
	"testing"
)

// TestKafkaZeroLossAndIdempotence 验证零丢失策略与消费端幂等防重拦截
func TestKafkaZeroLossAndIdempotence(t *testing.T) {
	broker := NewMockBroker(2)
	msg := ProducerMessage{Key: "k", Value: "v", MessageID: "m1"}

	// 1. 副本不足拒绝
	if err := broker.Append(msg, 1); err == nil {
		t.Fatal("expected error when active replicas < min.insync.replicas")
	}

	// 2. 副本满足成功
	if err := broker.Append(msg, 2); err != nil {
		t.Fatalf("expected success, got: %v", err)
	}

	// 3. 消费端幂等性
	consumer := NewIdempotentConsumer()
	isDup, _ := consumer.ProcessMessage(msg)
	if isDup {
		t.Fatal("first process should not be duplicate")
	}

	isDup, _ = consumer.ProcessMessage(msg)
	if !isDup {
		t.Fatal("second process must be detected as duplicate")
	}

	// 4. 哈希路由确定性
	if HashRoutePartition("test_order", 10) != HashRoutePartition("test_order", 10) {
		t.Fatal("hash routing must be deterministic")
	}
}

// 1. 消费端幂等表查询基准测试
func BenchmarkIdempotentConsumer(b *testing.B) {
	consumer := NewIdempotentConsumer()
	msg := ProducerMessage{Key: "order", Value: "val", MessageID: "unique_id_100"}
	_, _ = consumer.ProcessMessage(msg) // 预先插入

	b.ReportAllocs()
	for b.Loop() {
		_, _ = consumer.ProcessMessage(msg)
	}
}

// 2. 顺序分区路由哈希计算基准测试
func BenchmarkHashRoutePartition(b *testing.B) {
	keys := make([]string, 100)
	for i := 0; i < 100; i++ {
		keys[i] = fmt.Sprintf("order_key_%d", i)
	}

	b.ReportAllocs()
	for b.Loop() {
		for _, k := range keys {
			_ = HashRoutePartition(k, 16)
		}
	}
}
