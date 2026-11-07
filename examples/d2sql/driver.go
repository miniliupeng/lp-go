// Package sqltut 是一个只为讲事务边界的内存驱动。
// 它不是数据库隔离级别：同一键后提交的事务覆盖先提交的事务；不同键都保留。
// driver.TxOptions 被忽略。
package sqltut

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"maps"
	"sync"
)

func init() {
	sql.Register("sqltut", drv{})
}

type drv struct{}

func (drv) Open(name string) (driver.Conn, error) {
	return &conn{db: dbFor(name)}, nil
}

var (
	regMu sync.Mutex
	reg   = map[string]*mem{}
)

type mem struct {
	mu sync.Mutex
	kv map[string]int64
}

func dbFor(name string) *mem {
	regMu.Lock()
	defer regMu.Unlock()
	m, ok := reg[name]
	if !ok {
		m = &mem{kv: map[string]int64{}}
		reg[name] = m
	}
	return m
}

type conn struct {
	db *mem
	mu sync.Mutex
	tx *tx
}

func (c *conn) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("sqltut: prepare is not supported")
}

func (c *conn) Close() error { return nil }

func (c *conn) Begin() (driver.Tx, error) {
	return c.BeginTx(context.Background(), driver.TxOptions{})
}

func (c *conn) BeginTx(ctx context.Context, _ driver.TxOptions) (driver.Tx, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.tx != nil {
		return nil, errors.New("sqltut: transaction already open")
	}
	c.db.mu.Lock()
	base := maps.Clone(c.db.kv)
	c.db.mu.Unlock()
	c.tx = &tx{c: c, base: base, dirty: map[string]int64{}}
	return c.tx, nil
}

// database/sql 的 Tx.ExecContext / QueryContext 调用的是这条连接上的方法，
// 不会调用 driver.Tx 上的同名方法。连接自己记住当前事务，才能看见未提交的写入。
func (c *conn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	c.mu.Lock()
	tx := c.tx
	c.mu.Unlock()
	if tx != nil {
		return execMap(nil, tx.dirty, query, args)
	}
	return execMap(c.db, nil, query, args)
}

func (c *conn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	c.mu.Lock()
	tx := c.tx
	c.mu.Unlock()
	if tx != nil {
		return queryMap(tx.base, tx.dirty, query, args)
	}
	c.db.mu.Lock()
	defer c.db.mu.Unlock()
	return queryMap(c.db.kv, nil, query, args)
}

type tx struct {
	c     *conn
	base  map[string]int64
	dirty map[string]int64
}

func (t *tx) Commit() error {
	t.c.mu.Lock()
	defer t.c.mu.Unlock()
	t.c.db.mu.Lock()
	defer t.c.db.mu.Unlock()
	for k, v := range t.dirty {
		t.c.db.kv[k] = v
	}
	t.c.tx = nil
	return nil
}

func (t *tx) Rollback() error {
	t.c.mu.Lock()
	defer t.c.mu.Unlock()
	t.c.tx = nil
	return nil
}

type result struct{}

func (result) LastInsertId() (int64, error) { return 0, errors.New("sqltut: no insert id") }
func (result) RowsAffected() (int64, error) { return 1, nil }

func execMap(db *mem, dirty map[string]int64, query string, args []driver.NamedValue) (driver.Result, error) {
	if query != "upsert" {
		return nil, errors.New("sqltut: unknown exec")
	}
	k, v, err := keyVal(args)
	if err != nil {
		return nil, err
	}
	if dirty != nil {
		dirty[k] = v
		return result{}, nil
	}
	db.mu.Lock()
	db.kv[k] = v
	db.mu.Unlock()
	return result{}, nil
}

func queryMap(base, dirty map[string]int64, query string, args []driver.NamedValue) (driver.Rows, error) {
	if query != "get" || len(args) != 1 {
		return nil, errors.New("sqltut: get wants one key")
	}
	k, ok := args[0].Value.(string)
	if !ok {
		return nil, errors.New("sqltut: key must be string")
	}
	if dirty != nil {
		if v, ok := dirty[k]; ok {
			return &rows{val: v, ok: true}, nil
		}
	}
	v, ok := base[k]
	return &rows{val: v, ok: ok}, nil
}

func keyVal(args []driver.NamedValue) (string, int64, error) {
	if len(args) != 2 {
		return "", 0, errors.New("sqltut: upsert wants key and int64")
	}
	k, ok := args[0].Value.(string)
	if !ok {
		return "", 0, errors.New("sqltut: key must be string")
	}
	v, ok := args[1].Value.(int64)
	if !ok {
		return "", 0, errors.New("sqltut: value must be int64")
	}
	return k, v, nil
}

type rows struct {
	val  int64
	ok   bool
	done bool
}

func (r *rows) Columns() []string { return []string{"v"} }
func (r *rows) Close() error      { return nil }
func (r *rows) Next(dest []driver.Value) error {
	if r.done || !r.ok {
		return io.EOF
	}
	dest[0] = r.val
	r.done = true
	return nil
}
