package main

import (
	"context"
	"crypto/sha256"
	"fmt"
	"sync"
	"time"
)

// ==================== 1. 消息端到端零丢失：生产者与 Broker 交互模拟 ====================

// ProducerMessage 模拟一条经过幂等哈希包装的 Kafka 消息
type ProducerMessage struct {
	Key       string
	Value     string
	MessageID string // 业务唯一生成的去重 ID (UUID / 订单号)
}

// MockBroker 模拟 Kafka Broker 端副本同步与提交日志
type MockBroker struct {
	mu             sync.Mutex
	committedLog   []ProducerMessage
	minInsyncCount int
}

func NewMockBroker(minInsync int) *MockBroker {
	return &MockBroker{
		minInsyncCount: minInsync,
	}
}

// Append 模拟 acks=all 写入：必须达到最小同步副本数后才返回成功
func (b *MockBroker) Append(msg ProducerMessage, activeReplicas int) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if activeReplicas < b.minInsyncCount {
		return fmt.Errorf("写入失败: 当前活跃 ISR 副本数 (%d) 低于 min.insync.replicas (%d)",
			activeReplicas, b.minInsyncCount)
	}

	b.committedLog = append(b.committedLog, msg)
	return nil
}

func demoProducerZeroLoss() {
	broker := NewMockBroker(2) // 生产级要求：min.insync.replicas = 2

	msg := ProducerMessage{
		Key:       "order_1001",
		Value:     "PAYMENT_SUCCESS",
		MessageID: "msg_uuid_998877",
	}

	// 1. 模拟网络分区或节点宕机：当前 ISR 副本数仅为 1，不满足最小同步要求
	err := broker.Append(msg, 1)
	fmt.Printf("[零丢失策略实证] ISR不足时 (acks=all) 写入结果: %v\n", err)

	// 2. 集群正常：当前 ISR 副本数为 2，满足要求
	err = broker.Append(msg, 2)
	fmt.Printf("[零丢失策略实证] ISR充足时 (acks=all) 写入结果: %v (消息安全持久化)\n", err == nil)
}

// ==================== 2. 消费端幂等去重状态机 (Idempotent Consumer) ====================

// IdempotentConsumer 模拟支持去重的消费端
type IdempotentConsumer struct {
	mu             sync.Mutex
	processedKeyDB map[string]struct{} // 模拟 Redis/MySQL 幂等唯一键表
}

func NewIdempotentConsumer() *IdempotentConsumer {
	return &IdempotentConsumer{
		processedKeyDB: make(map[string]struct{}),
	}
}

// ProcessMessage 执行幂等消费：同一消息重试投递时自动拦截
func (c *IdempotentConsumer) ProcessMessage(msg ProducerMessage) (isDuplicate bool, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 1. 检查去重表中是否已存在该消息唯一 ID
	if _, exists := c.processedKeyDB[msg.MessageID]; exists {
		// 已消费过，直接确认并跳过，防止重复扣款或发货
		return true, nil
	}

	// 2. 执行实际业务逻辑 (如修改订单状态)
	// 3. 将消息 ID 记录进幂等表中 (生产中与业务更新处于同一数据库事务中)
	c.processedKeyDB[msg.MessageID] = struct{}{}
	return false, nil
}

func demoConsumerIdempotence() {
	consumer := NewIdempotentConsumer()

	msg := ProducerMessage{
		Key:       "order_1001",
		Value:     "CHARGE_50_DOLLARS",
		MessageID: "tx_order_charge_unique_id_001",
	}

	fmt.Println("\n=== 2. 消费端幂等防重表实战 (防止 At-least-once 重复投递) ===")
	// 第 1 次正常消费
	isDup1, _ := consumer.ProcessMessage(msg)
	fmt.Printf("首次消费消息 (ID: %s): 是否重复? %v (成功扣款)\n", msg.MessageID, isDup1)

	// 模拟由于网络抖动，ACK 丢失导致 Broker 重新投递了同一条消息
	isDup2, _ := consumer.ProcessMessage(msg)
	fmt.Printf("重复收到相同消息 (ID: %s): 是否重复? %v (成功拦截，避免重复扣款！)\n", msg.MessageID, isDup2)
}

// ==================== 3. 顺序消息哈希路由算法 (Partition Routing) ====================

// HashRoutePartition 根据业务 Key 计算目标分区 (确保相同订单的消息绝对进入同一分区)
func HashRoutePartition(key string, partitionCount int) int {
	h := sha256.Sum256([]byte(key))
	// 取前 4 字节转整数
	hashVal := int(h[0])<<24 | int(h[1])<<16 | int(h[2])<<8 | int(h[3])
	if hashVal < 0 {
		hashVal = -hashVal
	}
	return hashVal % partitionCount
}

func demoHashRouting() {
	fmt.Println("\n=== 3. 严格顺序消息：基于 Key 的一致性分区路由 ===")
	partitionCount := 8
	orderID := "order_2026_0908"

	p1 := HashRoutePartition(orderID, partitionCount)
	p2 := HashRoutePartition(orderID, partitionCount)

	fmt.Printf("订单 [%s] 创建事件路由至分区: %d\n", orderID, p1)
	fmt.Printf("订单 [%s] 支付事件路由至分区: %d (两次路由一致: %v，保证单分区内严格 FIFO 有序！)\n",
		orderID, p2, p1 == p2)
}

func main() {
	_ = context.Background()
	_ = time.Second

	fmt.Println("=== 1. Kafka 生产端 acks=all 与 min.insync.replicas 零丢失机制 ===")
	demoProducerZeroLoss()

	demoConsumerIdempotence()

	demoHashRouting()
}
