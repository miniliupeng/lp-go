package sqltutorial

import (
	"context"
	"database/sql"
	"errors"
	"testing"
)

func openDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("tutorialkv", t.Name())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func get(t *testing.T, db *sql.DB, key string) (int64, error) {
	t.Helper()
	var v int64
	err := db.QueryRow("get", key).Scan(&v)
	return v, err
}

func TestCommitAndRollback(t *testing.T) {
	db := openDB(t)

	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	if _, err = tx.Exec("upsert", "alice", int64(10)); err != nil {
		t.Fatal(err)
	}
	if err = tx.Rollback(); err != nil {
		t.Fatal(err)
	}
	_, err = get(t, db, "alice")
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("after rollback: %v", err)
	}

	tx, err = db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	if _, err = tx.Exec("upsert", "alice", int64(10)); err != nil {
		t.Fatal(err)
	}
	var own int64
	if err = tx.QueryRow("get", "alice").Scan(&own); err != nil {
		t.Fatal(err)
	}
	if own != 10 {
		t.Fatalf("read your write = %d", own)
	}
	if err = tx.Commit(); err != nil {
		t.Fatal(err)
	}
	got, err := get(t, db, "alice")
	if err != nil || got != 10 {
		t.Fatalf("committed = %d, %v", got, err)
	}
}

func TestDifferentKeysBothStick(t *testing.T) {
	db := openDB(t)
	tx1, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	tx2, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	if _, err = tx1.Exec("upsert", "a", int64(1)); err != nil {
		t.Fatal(err)
	}
	if _, err = tx2.Exec("upsert", "b", int64(2)); err != nil {
		t.Fatal(err)
	}
	if err = tx1.Commit(); err != nil {
		t.Fatal(err)
	}
	if err = tx2.Commit(); err != nil {
		t.Fatal(err)
	}
	a, err := get(t, db, "a")
	if err != nil || a != 1 {
		t.Fatalf("a=%d %v", a, err)
	}
	b, err := get(t, db, "b")
	if err != nil || b != 2 {
		t.Fatalf("b=%d %v", b, err)
	}
}

func TestSameKeyLastCommitWins(t *testing.T) {
	db := openDB(t)
	tx1, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	tx2, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	if _, err = tx1.Exec("upsert", "k", int64(1)); err != nil {
		t.Fatal(err)
	}
	if _, err = tx2.Exec("upsert", "k", int64(2)); err != nil {
		t.Fatal(err)
	}
	if err = tx1.Commit(); err != nil {
		t.Fatal(err)
	}
	if err = tx2.Commit(); err != nil {
		t.Fatal(err)
	}
	got, err := get(t, db, "k")
	if err != nil || got != 2 {
		t.Fatalf("last commit = %d, %v", got, err)
	}
}

func TestCanceledBegin(t *testing.T) {
	db := openDB(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := db.BeginTx(ctx, nil)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("begin: %v", err)
	}
}
