package configs

import (
	"fmt"
	"time"
)

// Config 定义企业级网关全量运行配置
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Security SecurityConfig `yaml:"security"`
	Limiter  LimiterConfig  `yaml:"limiter"`
	Breaker  BreakerConfig  `yaml:"breaker"`
	RAG      RAGConfig      `yaml:"rag"`
	Upstream UpstreamConfig `yaml:"upstream"`
}

type ServerConfig struct {
	Port         int           `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
	DrainTimeout time.Duration `yaml:"drain_timeout"` // 优雅退出排空超时
}

type SecurityConfig struct {
	APIKeys []string `yaml:"api_keys"` // 授权 API-Key 列表
}

type LimiterConfig struct {
	Enabled  bool    `yaml:"enabled"`
	Rate     float64 `yaml:"rate"`     // 令牌填充速率 (token/s)
	Capacity float64 `yaml:"capacity"` // 令牌桶容量
}

type BreakerConfig struct {
	Enabled          bool          `yaml:"enabled"`
	FailureThreshold int           `yaml:"failure_threshold"` // 触发熔断的连续失败次数
	OpenTimeout      time.Duration `yaml:"open_timeout"`       // 熔断后冷却进入半开状态的时长
}

type RAGConfig struct {
	Enabled bool    `yaml:"enabled"`
	TopK    int     `yaml:"top_k"`
	RRFK    float32 `yaml:"rrf_k"` // RRF 平滑系数
}

type UpstreamConfig struct {
	DefaultModel string        `yaml:"default_model"`
	Timeout      time.Duration `yaml:"timeout"`
	MockDelay    time.Duration `yaml:"mock_delay"` // 模拟推理 Token 吐出时延
}

// DefaultConfig 提供经过生产调优的默认配置
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         8090,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 60 * time.Second, // SSE 流式长连接需要较长 WriteTimeout
			IdleTimeout:  120 * time.Second,
			DrainTimeout: 10 * time.Second,
		},
		Security: SecurityConfig{
			APIKeys: []string{"sk-production-ai-gateway-key-9527", "sk-test-client-key-8888"},
		},
		Limiter: LimiterConfig{
			Enabled:  true,
			Rate:     100.0, // 每秒 100 个请求
			Capacity: 200.0, // 允许突发 200 个请求
		},
		Breaker: BreakerConfig{
			Enabled:          true,
			FailureThreshold: 5,
			OpenTimeout:      5 * time.Second,
		},
		RAG: RAGConfig{
			Enabled: true,
			TopK:    3,
			RRFK:    60.0,
		},
		Upstream: UpstreamConfig{
			DefaultModel: "deepseek-r1-chat",
			Timeout:      30 * time.Second,
			MockDelay:    15 * time.Millisecond,
		},
	}
}

// Validate 校验配置合法性
func (c *Config) Validate() error {
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}
	if c.Limiter.Enabled && (c.Limiter.Rate <= 0 || c.Limiter.Capacity <= 0) {
		return fmt.Errorf("invalid rate limiter params: rate=%.2f, capacity=%.2f", c.Limiter.Rate, c.Limiter.Capacity)
	}
	if c.Breaker.Enabled && c.Breaker.FailureThreshold <= 0 {
		return fmt.Errorf("invalid circuit breaker threshold: %d", c.Breaker.FailureThreshold)
	}
	return nil
}
