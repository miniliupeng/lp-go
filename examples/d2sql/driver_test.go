package sqltut

import (
	"context"
	"database/sql"
	"errors"
	"testing"
)

func open(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqltut", t.Name())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func TestCommitAndRollback(t *testing.T) {
	db := open(t)
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = tx.ExecContext(ctx, "upsert", "a", int64(1)); err != nil {
		t.Fatal(err)
	}
	var v int64
	if err = tx.QueryRowContext(ctx, "get", "a").Scan(&v); err != nil || v != 1 {
		t.Fatalf("own write %d %v", v, err)
	}
	if err = tx.Rollback(); err != nil {
		t.Fatal(err)
	}
	err = db.QueryRowContext(ctx, "get", "a").Scan(&v)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatal(err)
	}

	tx, err = db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = tx.ExecContext(ctx, "upsert", "a", int64(2)); err != nil {
		t.Fatal(err)
	}
	if err = tx.Commit(); err != nil {
		t.Fatal(err)
	}
	if err = db.QueryRowContext(ctx, "get", "a").Scan(&v); err != nil || v != 2 {
		t.Fatalf("committed %d %v", v, err)
	}
}

func TestSameKeyLastCommitWins(t *testing.T) {
	db := open(t)
	ctx := context.Background()
	tx1, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	tx2, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = tx1.ExecContext(ctx, "upsert", "k", int64(1)); err != nil {
		t.Fatal(err)
	}
	if _, err = tx2.ExecContext(ctx, "upsert", "k", int64(9)); err != nil {
		t.Fatal(err)
	}
	if err = tx1.Commit(); err != nil {
		t.Fatal(err)
	}
	if err = tx2.Commit(); err != nil {
		t.Fatal(err)
	}
	var v int64
	if err = db.QueryRowContext(ctx, "get", "k").Scan(&v); err != nil || v != 9 {
		t.Fatalf("%d %v", v, err)
	}
}

func TestDifferentKeysBothStick(t *testing.T) {
	db := open(t)
	ctx := context.Background()
	tx1, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	tx2, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = tx1.ExecContext(ctx, "upsert", "a", int64(1)); err != nil {
		t.Fatal(err)
	}
	if _, err = tx2.ExecContext(ctx, "upsert", "b", int64(2)); err != nil {
		t.Fatal(err)
	}
	if err = tx1.Commit(); err != nil || tx2.Commit() != nil {
		t.Fatal(err)
	}
	var a, b int64
	if err = db.QueryRowContext(ctx, "get", "a").Scan(&a); err != nil || a != 1 {
		t.Fatal(a, err)
	}
	if err = db.QueryRowContext(ctx, "get", "b").Scan(&b); err != nil || b != 2 {
		t.Fatal(b, err)
	}
}

func TestCanceledBegin(t *testing.T) {
	db := open(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := db.BeginTx(ctx, nil); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
}
