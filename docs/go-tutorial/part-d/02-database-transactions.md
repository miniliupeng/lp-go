# 02 数据库访问与事务

`database/sql` 是连接池和 API，不是数据库。隔离级别、锁、SQL 方言都由驱动和数据库决定。驱动不支持取消时，`QueryContext` 会等到查询结束才返回。

本仓库没有第三方驱动。 [sqltutorial](../examples/partd/sqltutorial/driver.go) 注册了一个只为教程存在的驱动 `tutorialkv`，这样 `database/sql` 的 `Begin` / `Commit` / `Rollback` 能在 `go test` 里跑完。

## 这个驱动实际保证什么

语句只有两条，参数走 `database/sql` 的占位参数，不是拼进字符串：

- `upsert`：键 `string`，值 `int64`
- `get`：键 `string`。没有这个键时返回空结果，`QueryRow` 得到 `sql.ErrNoRows`

事务开始时复制一份读快照。本事务里先写后读，能读到自己的写。提交时**只把本事务写过的键**写回共享表。因此：

- 两个事务写不同的键，两个键都在。测试 `TestDifferentKeysBothStick`。
- 两个事务写同一个键，后提交的覆盖先提交的，没有冲突错误。测试 `TestSameKeyLastCommitWins`。
- `Rollback` 丢掉写集合，外面读不到。测试 `TestCommitAndRollback`。

这不是 MySQL 的 RR 或 SI。它没有行锁，也没有「提交时发现版本变了就失败」。讲生产事务时把这一点说在前面，避免把教学驱动的丢更新说成数据库的正常隔离。

已取消的 `context` 调用 `BeginTx`，返回 `context.Canceled`。`database/sql` 在进驱动之前就会看 context。

## 连接池

`sql.Open` 不一定马上连上。池的参数：

- `SetMaxOpenConns`：同时打开的连接上限。0 表示不限制，生产上要设。
- `SetMaxIdleConns`：池里保留的空闲连接。
- `SetConnMaxLifetime` / `SetConnMaxIdleTime`：连接活太久或闲太久就关掉，用来躲开服务端单方面关掉的连接。

`Rows` 必须 `Close`，否则连接回不了池。循环之后还要看 `rows.Err()`。`QueryRow` 找不到行是 `sql.ErrNoRows`，用 `errors.Is` 判断，不要比字符串。

## 事务的写法

```go
tx, err := db.BeginTx(ctx, nil)
if err != nil {
    return err
}
defer tx.Rollback()

if _, err = tx.ExecContext(ctx, "upsert", key, value); err != nil {
    return err
}
return tx.Commit()
```

`Rollback` 在 `Commit` 成功之后再调用，驱动应允许它变成空操作。教学驱动就是这样做的。`defer Rollback` 然后 `return Commit()` 时，成功路径上的 Rollback 是无害的；失败路径上它丢掉未提交的写。

`TxOptions.Isolation` 只是转交给驱动。教学驱动忽略它。生产代码要写明数据库实际提供的级别，而不是假设 `sql.LevelSerializable` 在所有驱动里都生效。

预备语句由池在每条连接上分别准备。不要把 `*sql.Stmt` 的并发使用当成「一条数据库连接」。

## 面试常问

- `*sql.DB` 是池，可以并发用。`*sql.Tx` 绑定在一条连接上，不要并发使用同一个 `Tx`。
- 为什么 `defer Rollback` 仍要检查 `Commit` 的错误。
- 丢更新：教学驱动故意最后提交获胜。真正的数据库要靠事务隔离、唯一约束或 `UPDATE ... WHERE version = ?`。

## 运行

```bash
go test ./docs/go-tutorial/examples/partd/sqltutorial -count=1
```

## 练习

在测试里让事务 A 写入 `a=1` 并提交，事务 B 在 A 提交前只做了 `get a`。写出 B 的快照里有没有 `a`。然后跑测试验证你的判断。
