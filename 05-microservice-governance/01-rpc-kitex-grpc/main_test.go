package main

import (
	"context"
	"net"
	"testing"
	"time"
)

func TestRPCClientServer(t *testing.T) {
	server := NewRPCServer()
	server.Register("Ping", func(payload []byte) ([]byte, error) {
		return []byte("Pong:" + string(payload)), nil
	})

	if err := server.Start("127.0.0.1:0"); err != nil {
		t.Fatalf("failed to start server: %v", err)
	}
	defer server.Close()

	conn, err := net.Dial("tcp", server.Addr())
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}
	client := NewRPCClient(conn)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	resp, err := client.Call(ctx, "Ping", []byte("Hello"))
	if err != nil {
		t.Fatalf("call failed: %v", err)
	}

	if string(resp) != "Pong:Hello" {
		t.Fatalf("unexpected resp: %s", string(resp))
	}
}

func BenchmarkBinaryCodec(b *testing.B) {
	codec := &BinaryCodec{}
	req := &RPCRequest{
		SeqID:   1001,
		Method:  "UserService.GetUserInfo",
		Payload: []byte(`{"user_id":123456,"region":"ap-southeast"}`),
	}

	b.ResetTimer()
	for b.Loop() {
		data, err := codec.EncodeRequest(req)
		if err != nil {
			b.Fatal(err)
		}
		if len(data) == 0 {
			b.Fatal("empty data")
		}
	}
}

func BenchmarkJSONCodec(b *testing.B) {
	codec := &JSONCodec{}
	req := &RPCRequest{
		SeqID:   1001,
		Method:  "UserService.GetUserInfo",
		Payload: []byte(`{"user_id":123456,"region":"ap-southeast"}`),
	}

	b.ResetTimer()
	for b.Loop() {
		data, err := codec.EncodeRequest(req)
		if err != nil {
			b.Fatal(err)
		}
		if len(data) == 0 {
			b.Fatal("empty data")
		}
	}
}

func BenchmarkEndToEndRPCCall(b *testing.B) {
	server := NewRPCServer()
	server.Register("Echo", func(payload []byte) ([]byte, error) {
		return payload, nil
	})
	if err := server.Start("127.0.0.1:0"); err != nil {
		b.Fatalf("failed to start server: %v", err)
	}
	defer server.Close()

	conn, err := net.Dial("tcp", server.Addr())
	if err != nil {
		b.Fatalf("failed to dial: %v", err)
	}
	client := NewRPCClient(conn)
	defer client.Close()

	ctx := context.Background()
	payload := []byte("BenchmarkPayloadData")

	b.ResetTimer()
	for b.Loop() {
		_, err := client.Call(ctx, "Echo", payload)
		if err != nil {
			b.Fatal(err)
		}
	}
}
