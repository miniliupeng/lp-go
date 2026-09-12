package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

// ============================================================================
// 1. 跨进程 RPC 协议与编解码对比 (Protobuf/Binary vs JSON)
// ============================================================================

// RPCRequest 表示 RPC 请求报文
type RPCRequest struct {
	SeqID   uint64
	Method  string
	Payload []byte
}

// RPCResponse 表示 RPC 响应报文
type RPCResponse struct {
	SeqID uint64
	Err   string
	Data  []byte
}

// BinaryCodec 模拟高性能二进制序列化 (Protobuf/Kitex FastCodec 理念：固定定长头 + 变长数据)
type BinaryCodec struct{}

func (b *BinaryCodec) EncodeRequest(req *RPCRequest) ([]byte, error) {
	buf := new(bytes.Buffer)
	// 协议头：Magic(2B) + SeqID(8B) + MethodLen(2B) + PayloadLen(4B)
	if err := binary.Write(buf, binary.BigEndian, uint16(0xCAFE)); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, req.SeqID); err != nil {
		return nil, err
	}
	methodBytes := []byte(req.Method)
	if err := binary.Write(buf, binary.BigEndian, uint16(len(methodBytes))); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, uint32(len(req.Payload))); err != nil {
		return nil, err
	}
	buf.Write(methodBytes)
	buf.Write(req.Payload)
	return buf.Bytes(), nil
}

func (b *BinaryCodec) DecodeRequest(r io.Reader) (*RPCRequest, error) {
	var magic uint16
	if err := binary.Read(r, binary.BigEndian, &magic); err != nil {
		return nil, err
	}
	if magic != 0xCAFE {
		return nil, errors.New("invalid magic number")
	}

	var req RPCRequest
	if err := binary.Read(r, binary.BigEndian, &req.SeqID); err != nil {
		return nil, err
	}

	var methodLen uint16
	if err := binary.Read(r, binary.BigEndian, &methodLen); err != nil {
		return nil, err
	}

	var payloadLen uint32
	if err := binary.Read(r, binary.BigEndian, &payloadLen); err != nil {
		return nil, err
	}

	methodBuf := make([]byte, methodLen)
	if _, err := io.ReadFull(r, methodBuf); err != nil {
		return nil, err
	}
	req.Method = string(methodBuf)

	req.Payload = make([]byte, payloadLen)
	if _, err := io.ReadFull(r, req.Payload); err != nil {
		return nil, err
	}

	return &req, nil
}

// JSONCodec 经典文本协议编解码 (作为性能对比基准)
type JSONCodec struct{}

func (j *JSONCodec) EncodeRequest(req *RPCRequest) ([]byte, error) {
	return json.Marshal(req)
}

func (j *JSONCodec) DecodeRequest(data []byte) (*RPCRequest, error) {
	var req RPCRequest
	err := json.Unmarshal(data, &req)
	return &req, err
}

// ============================================================================
// 2. 模拟 Kitex Netpoll 多路复用连接池与调用客户端
// ============================================================================

type UserQuery struct {
	UserID int64  `json:"user_id"`
	Region string `json:"region"`
}

type UserInfo struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Level    string `json:"level"`
}

// RPCClient 模拟支持连接复用与超时的客户端
type RPCClient struct {
	conn   net.Conn
	codec  *BinaryCodec
	seq    atomic.Uint64
	mu     sync.Mutex
	pend   map[uint64]chan *RPCResponse
	closed atomic.Bool
}

func NewRPCClient(conn net.Conn) *RPCClient {
	c := &RPCClient{
		conn:  conn,
		codec: &BinaryCodec{},
		pend:  make(map[uint64]chan *RPCResponse),
	}
	go c.readLoop()
	return c
}

func (c *RPCClient) readLoop() {
	for !c.closed.Load() {
		var magic uint16
		if err := binary.Read(c.conn, binary.BigEndian, &magic); err != nil {
			break
		}
		var seqID uint64
		_ = binary.Read(c.conn, binary.BigEndian, &seqID)
		var errLen uint16
		_ = binary.Read(c.conn, binary.BigEndian, &errLen)
		var dataLen uint32
		_ = binary.Read(c.conn, binary.BigEndian, &dataLen)

		errBytes := make([]byte, errLen)
		_, _ = io.ReadFull(c.conn, errBytes)

		dataBytes := make([]byte, dataLen)
		_, _ = io.ReadFull(c.conn, dataBytes)

		resp := &RPCResponse{
			SeqID: seqID,
			Err:   string(errBytes),
			Data:  dataBytes,
		}

		c.mu.Lock()
		ch, ok := c.pend[seqID]
		if ok {
			delete(c.pend, seqID)
		}
		c.mu.Unlock()

		if ok && ch != nil {
			ch <- resp
		}
	}
}

// Call 同步RPC调用（支持Context超时控制）
func (c *RPCClient) Call(ctx context.Context, method string, payload []byte) ([]byte, error) {
	seq := c.seq.Add(1)
	req := &RPCRequest{
		SeqID:   seq,
		Method:  method,
		Payload: payload,
	}

	encoded, err := c.codec.EncodeRequest(req)
	if err != nil {
		return nil, err
	}

	respCh := make(chan *RPCResponse, 1)
	c.mu.Lock()
	c.pend[seq] = respCh
	c.mu.Unlock()

	// 发送数据
	if _, err := c.conn.Write(encoded); err != nil {
		c.mu.Lock()
		delete(c.pend, seq)
		c.mu.Unlock()
		return nil, err
	}

	select {
	case <-ctx.Done():
		c.mu.Lock()
		delete(c.pend, seq)
		c.mu.Unlock()
		return nil, ctx.Err()
	case resp := <-respCh:
		if resp.Err != "" {
			return nil, errors.New(resp.Err)
		}
		return resp.Data, nil
	}
}

func (c *RPCClient) Close() error {
	c.closed.Store(true)
	return c.conn.Close()
}

// ============================================================================
// 3. RPC 服务端实现
// ============================================================================

type RPCServer struct {
	listener net.Listener
	codec    *BinaryCodec
	handlers map[string]func([]byte) ([]byte, error)
	closed   atomic.Bool
}

func NewRPCServer() *RPCServer {
	return &RPCServer{
		codec:    &BinaryCodec{},
		handlers: make(map[string]func([]byte) ([]byte, error)),
	}
}

func (s *RPCServer) Register(method string, handler func([]byte) ([]byte, error)) {
	s.handlers[method] = handler
}

func (s *RPCServer) Start(addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	s.listener = l

	go func() {
		for !s.closed.Load() {
			conn, err := s.listener.Accept()
			if err != nil {
				if s.closed.Load() {
					return
				}
				continue
			}
			go s.handleConn(conn)
		}
	}()
	return nil
}

func (s *RPCServer) handleConn(conn net.Conn) {
	defer conn.Close()
	for {
		req, err := s.codec.DecodeRequest(conn)
		if err != nil {
			return
		}

		handler, ok := s.handlers[req.Method]
		var respData []byte
		var errMsg string
		if !ok {
			errMsg = fmt.Sprintf("method %s not found", req.Method)
		} else {
			respData, err = handler(req.Payload)
			if err != nil {
				errMsg = err.Error()
			}
		}

		// 发送响应报文: Magic(2B) + SeqID(8B) + ErrLen(2B) + DataLen(4B) + Err + Data
		buf := new(bytes.Buffer)
		_ = binary.Write(buf, binary.BigEndian, uint16(0xCAFE))
		_ = binary.Write(buf, binary.BigEndian, req.SeqID)
		_ = binary.Write(buf, binary.BigEndian, uint16(len(errMsg)))
		_ = binary.Write(buf, binary.BigEndian, uint32(len(respData)))
		buf.WriteString(errMsg)
		buf.Write(respData)

		if _, err := conn.Write(buf.Bytes()); err != nil {
			return
		}
	}
}

func (s *RPCServer) Close() error {
	s.closed.Store(true)
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func (s *RPCServer) Addr() string {
	if s.listener != nil {
		return s.listener.Addr().String()
	}
	return ""
}

func main() {
	fmt.Println("=== 阶段三 专题01：微服务 RPC 与编解码机制 (Kitex vs gRPC) ===")

	// 1. 启动轻量二进制 RPC 服务
	server := NewRPCServer()
	server.Register("UserService.GetUserInfo", func(payload []byte) ([]byte, error) {
		var query UserQuery
		if err := json.Unmarshal(payload, &query); err != nil {
			return nil, err
		}
		info := UserInfo{
			UserID:   query.UserID,
			Username: fmt.Sprintf("Engineer_%d", query.UserID),
			Level:    "Staff Engineer (T2)",
		}
		return json.Marshal(info)
	})

	if err := server.Start("127.0.0.1:0"); err != nil {
		panic(err)
	}
	defer server.Close()

	// 2. 建立客户端连接
	conn, err := net.Dial("tcp", server.Addr())
	if err != nil {
		panic(err)
	}
	client := NewRPCClient(conn)
	defer client.Close()

	// 3. 执行 RPC 调用
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	queryData, _ := json.Marshal(UserQuery{UserID: 8888, Region: "CN-North"})
	data, err := client.Call(ctx, "UserService.GetUserInfo", queryData)
	if err != nil {
		fmt.Printf("RPC Call Failed: %v\n", err)
		return
	}

	var user UserInfo
	_ = json.Unmarshal(data, &user)
	fmt.Printf("[Success] RPC 响应成功: UserID=%d, Username=%s, Level=%s\n", user.UserID, user.Username, user.Level)

	// 4. 超时验证
	ctxTimeout, cancelTimeout := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancelTimeout()
	_, errTimeout := client.Call(ctxTimeout, "UserService.GetUserInfo", queryData)
	fmt.Printf("[Timeout Check] 预期超时拦截: %v\n", errTimeout)
}
