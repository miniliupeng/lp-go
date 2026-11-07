// Package sqltutorial 是一个只为教程服务的 database/sql 驱动。
//
// 方言只有两条语句，参数用 database/sql 的占位符传递：
//
//	upsert    参数：key string, value int64
//	get       参数：key string，结果一列 v；没有这个键时返回空行
//
// 事务在 Begin 时复制一份读快照。提交时只把本事务写过的键写回共享表。
// 不同的键可以一起留下。同一个键后提交的覆盖先提交的，驱动不做冲突检测。
// 这不是数据库的隔离级别，只让 database/sql 的 Begin/Commit/Rollback 能在没有第三方驱动时跑起来。
package sqltutorial

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"io"
	"maps"
	"sync"
)

func init() {
	sql.Register("tutorialkv", drv{})
}

type drv struct{}

func (drv) Open(name string) (driver.Conn, error) {
	return &memConn{st: storeFor(name)}, nil
}

type memStore struct {
	mu sync.Mutex
	kv map[string]int64
}

var (
	storesMu sync.Mutex
	stores   = map[string]*memStore{}
)

func storeFor(dsn string) *memStore {
	storesMu.Lock()
	defer storesMu.Unlock()
	st, ok := stores[dsn]
	if !ok {
		st = &memStore{kv: map[string]int64{}}
		stores[dsn] = st
	}
	return st
}

type memConn struct {
	mu sync.Mutex
	st *memStore
	tx *memTx
}

func (c *memConn) Prepare(query string) (driver.Stmt, error) {
	return &memStmt{c: c, q: query}, nil
}

func (c *memConn) Close() error { return nil }

func (c *memConn) Begin() (driver.Tx, error) {
	return c.BeginTx(context.Background(), driver.TxOptions{})
}

func (c *memConn) BeginTx(ctx context.Context, _ driver.TxOptions) (driver.Tx, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.tx != nil {
		return nil, errors.New("tutorialkv: already in a transaction")
	}
	c.st.mu.Lock()
	base := maps.Clone(c.st.kv)
	c.st.mu.Unlock()
	c.tx = &memTx{
		c:       c,
		base:    base,
		overlay: map[string]int64{},
		dirty:   map[string]bool{},
	}
	return c.tx, nil
}

func (c *memConn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.execLocked(query, namedValues(args))
}

func (c *memConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.queryLocked(query, namedValues(args))
}

func (c *memConn) execLocked(query string, args []driver.Value) (driver.Result, error) {
	if query != "upsert" {
		return nil, fmt.Errorf("tutorialkv: unknown statement %q", query)
	}
	if len(args) != 2 {
		return nil, errors.New("tutorialkv: upsert wants 2 args")
	}
	key, err := asString(args[0])
	if err != nil {
		return nil, err
	}
	val, err := asInt64(args[1])
	if err != nil {
		return nil, err
	}
	if c.tx != nil {
		c.tx.overlay[key] = val
		c.tx.dirty[key] = true
		return result(1), nil
	}
	c.st.mu.Lock()
	c.st.kv[key] = val
	c.st.mu.Unlock()
	return result(1), nil
}

func (c *memConn) queryLocked(query string, args []driver.Value) (driver.Rows, error) {
	if query != "get" {
		return nil, fmt.Errorf("tutorialkv: unknown statement %q", query)
	}
	if len(args) != 1 {
		return nil, errors.New("tutorialkv: get wants 1 arg")
	}
	key, err := asString(args[0])
	if err != nil {
		return nil, err
	}
	rows := &memRows{cols: []string{"v"}}
	if c.tx != nil {
		if c.tx.dirty[key] {
			rows.vals = []driver.Value{c.tx.overlay[key]}
			return rows, nil
		}
		if v, ok := c.tx.base[key]; ok {
			rows.vals = []driver.Value{v}
		}
		return rows, nil
	}
	c.st.mu.Lock()
	v, ok := c.st.kv[key]
	c.st.mu.Unlock()
	if ok {
		rows.vals = []driver.Value{v}
	}
	return rows, nil
}

type memTx struct {
	c       *memConn
	base    map[string]int64
	overlay map[string]int64
	dirty   map[string]bool
	done    bool
}

func (t *memTx) Commit() error {
	t.c.mu.Lock()
	defer t.c.mu.Unlock()
	if t.done || t.c.tx != t {
		return errors.New("tutorialkv: transaction is not open")
	}
	t.c.st.mu.Lock()
	for key, dirty := range t.dirty {
		if dirty {
			t.c.st.kv[key] = t.overlay[key]
		}
	}
	t.c.st.mu.Unlock()
	t.done = true
	t.c.tx = nil
	return nil
}

func (t *memTx) Rollback() error {
	t.c.mu.Lock()
	defer t.c.mu.Unlock()
	if t.done || t.c.tx != t {
		return nil
	}
	t.done = true
	t.c.tx = nil
	return nil
}

type memStmt struct {
	c *memConn
	q string
}

func (s *memStmt) Close() error  { return nil }
func (s *memStmt) NumInput() int { return -1 }

func (s *memStmt) Exec(args []driver.Value) (driver.Result, error) {
	s.c.mu.Lock()
	defer s.c.mu.Unlock()
	return s.c.execLocked(s.q, args)
}

func (s *memStmt) Query(args []driver.Value) (driver.Rows, error) {
	s.c.mu.Lock()
	defer s.c.mu.Unlock()
	return s.c.queryLocked(s.q, args)
}

type memRows struct {
	cols []string
	vals []driver.Value
	i    int
}

func (r *memRows) Columns() []string { return r.cols }
func (r *memRows) Close() error      { return nil }

func (r *memRows) Next(dest []driver.Value) error {
	if r.i >= len(r.vals) {
		return io.EOF
	}
	dest[0] = r.vals[r.i]
	r.i++
	return nil
}

type memResult int64

func result(n int64) driver.Result { return memResult(n) }

func (r memResult) LastInsertId() (int64, error) {
	return 0, errors.New("tutorialkv: no insert id")
}
func (r memResult) RowsAffected() (int64, error) { return int64(r), nil }

func namedValues(args []driver.NamedValue) []driver.Value {
	out := make([]driver.Value, len(args))
	for i, arg := range args {
		out[i] = arg.Value
	}
	return out
}

func asString(v driver.Value) (string, error) {
	s, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("tutorialkv: want string, got %T", v)
	}
	return s, nil
}

func asInt64(v driver.Value) (int64, error) {
	n, ok := v.(int64)
	if !ok {
		return 0, fmt.Errorf("tutorialkv: want int64, got %T", v)
	}
	return n, nil
}
