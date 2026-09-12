package domain

// 遵循 OpenAI /v1/chat/completions 工业标准 API 协议规范

type ChatRole string

const (
	RoleSystem    ChatRole = "system"
	RoleUser      ChatRole = "user"
	RoleAssistant ChatRole = "assistant"
)

type ChatMessage struct {
	Role    ChatRole `json:"role"`
	Content string   `json:"content"`
}

// ChatCompletionRequest 聊天补全请求
type ChatCompletionRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	Stream      bool          `json:"stream"`
	Temperature float32       `json:"temperature,omitempty"`
	EnableRAG   bool          `json:"enable_rag,omitempty"` // 是否开启企业私域 RAG 检索增强
}

// ChatChoice 非流式响应结果
type ChatChoice struct {
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message"`
	FinishReason string      `json:"finish_reason"`
}

// ChatCompletionResponse 非流式完整响应
type ChatCompletionResponse struct {
	ID      string       `json:"id"`
	Object  string       `json:"object"`
	Created int64        `json:"created"`
	Model   string       `json:"model"`
	Choices []ChatChoice `json:"choices"`
}

// StreamDelta 流式增量内容
type StreamDelta struct {
	Role    ChatRole `json:"role,omitempty"`
	Content string   `json:"content,omitempty"`
}

// StreamChoice 流式单条增量
type StreamChoice struct {
	Index        int         `json:"index"`
	Delta        StreamDelta `json:"delta"`
	FinishReason *string     `json:"finish_reason"`
}

// ChatCompletionChunk 标准 SSE 流式 Chunk 报文
type ChatCompletionChunk struct {
	ID      string         `json:"id"`
	Object  string         `json:"object"`
	Created int64          `json:"created"`
	Model   string         `json:"model"`
	Choices []StreamChoice `json:"choices"`
}
