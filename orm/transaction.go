package orm

import (
	"context"
	"database/sql"
)

// transaction 事务实现
type transaction struct {
	tx *sql.Tx
}

// Query 执行查询
func (t *transaction) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return t.tx.Query(query, args...)
}

// QueryRow 执行单行查询
func (t *transaction) QueryRow(query string, args ...interface{}) *sql.Row {
	return t.tx.QueryRow(query, args...)
}

// Exec 执行SQL语句
func (t *transaction) Exec(query string, args ...interface{}) (sql.Result, error) {
	return t.tx.Exec(query, args...)
}

// Commit 提交事务
func (t *transaction) Commit() error {
	return t.tx.Commit()
}

// Rollback 回滚事务
func (t *transaction) Rollback() error {
	return t.tx.Rollback()
}

// Table 在事务中创建查询构建器
func (t *transaction) Table(tableName string) QueryBuilder {
	return NewTransactionQueryBuilder(t, tableName)
}

// Model 在事务中基于模型创建查询构建器
func (t *transaction) Model(model interface{}) QueryBuilder {
	tableName := getTableNameFromModel(model)
	return NewTransactionQueryBuilder(t, tableName)
}

// getTableNameFromModel 从模型获取表名
func getTableNameFromModel(model interface{}) string {
	if m, ok := model.(ModelInterface); ok {
		return m.TableName()
	}

	// 使用反射获取结构体名称并转换为表名
	return camelToSnake(getStructName(model))
}

// TransactionManager 事务管理器
type TransactionManager struct {
	orm *ORM
}

// NewTransactionManager 创建事务管理器
func NewTransactionManager(orm *ORM) *TransactionManager {
	return &TransactionManager{orm: orm}
}

// WithTransaction 在事务中执行函数
func (tm *TransactionManager) WithTransaction(fn func(tx Tx) error) error {
	tx, err := tm.orm.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// WithTransactionContext 在带上下文的事务中执行函数
func (tm *TransactionManager) WithTransactionContext(ctx context.Context, opts *sql.TxOptions, fn func(tx Tx) error) error {
	tx, err := tm.orm.BeginTx(ctx, opts)
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// SavePoint 保存点
type SavePoint struct {
	name string
	tx   Tx
}

// NewSavePoint 创建保存点
func NewSavePoint(tx Tx, name string) (*SavePoint, error) {
	// 注意：不是所有数据库都支持保存点
	_, err := tx.Exec("SAVEPOINT " + name)
	if err != nil {
		return nil, err
	}

	return &SavePoint{
		name: name,
		tx:   tx,
	}, nil
}

// Rollback 回滚到保存点
func (sp *SavePoint) Rollback() error {
	_, err := sp.tx.Exec("ROLLBACK TO SAVEPOINT " + sp.name)
	return err
}

// Release 释放保存点
func (sp *SavePoint) Release() error {
	_, err := sp.tx.Exec("RELEASE SAVEPOINT " + sp.name)
	return err
}

// 全局便捷方法

// WithTransaction 在事务中执行函数
func WithTransaction(fn func(tx Tx) error) error {
	tm := NewTransactionManager(GetGlobalORM())
	return tm.WithTransaction(fn)
}

// WithTransactionContext 在带上下文的事务中执行函数
func WithTransactionContext(ctx context.Context, opts *sql.TxOptions, fn func(tx Tx) error) error {
	tm := NewTransactionManager(GetGlobalORM())
	return tm.WithTransactionContext(ctx, opts, fn)
}
