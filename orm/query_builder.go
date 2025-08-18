package orm

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

// queryBuilder 查询构建器实现
type queryBuilder struct {
	orm        *ORM
	tx         Tx
	tableName  string
	selectCols []string
	conditions []QueryCondition
	joins      []JoinClause
	orders     []OrderClause
	groups     []string
	havings    []HavingClause
	limitNum   int
	offsetNum  int
}

// NewQueryBuilder 创建新的查询构建器
func NewQueryBuilder(orm *ORM, tableName string) QueryBuilder {
	return &queryBuilder{
		orm:       orm,
		tableName: tableName,
	}
}

// NewTransactionQueryBuilder 创建事务查询构建器
func NewTransactionQueryBuilder(tx Tx, tableName string) QueryBuilder {
	return &queryBuilder{
		tx:        tx,
		tableName: tableName,
	}
}

// Select 选择字段
func (qb *queryBuilder) Select(columns ...string) QueryBuilder {
	qb.selectCols = columns
	return qb
}

// From 设置表名
func (qb *queryBuilder) From(table string) QueryBuilder {
	qb.tableName = table
	return qb
}

// Where 添加WHERE条件
func (qb *queryBuilder) Where(condition string, args ...interface{}) QueryBuilder {
	qb.conditions = append(qb.conditions, QueryCondition{
		Column:   condition,
		Operator: "=",
		Value:    args,
		Logic:    "AND",
	})
	return qb
}

// WhereIn 添加IN条件
func (qb *queryBuilder) WhereIn(column string, values ...interface{}) QueryBuilder {
	qb.conditions = append(qb.conditions, QueryCondition{
		Column:   column,
		Operator: "IN",
		Values:   values,
		Logic:    "AND",
	})
	return qb
}

// WhereNotIn 添加NOT IN条件
func (qb *queryBuilder) WhereNotIn(column string, values ...interface{}) QueryBuilder {
	qb.conditions = append(qb.conditions, QueryCondition{
		Column:   column,
		Operator: "NOT IN",
		Values:   values,
		Logic:    "AND",
	})
	return qb
}

// WhereBetween 添加BETWEEN条件
func (qb *queryBuilder) WhereBetween(column string, start, end interface{}) QueryBuilder {
	qb.conditions = append(qb.conditions, QueryCondition{
		Column:   column,
		Operator: "BETWEEN",
		Values:   []interface{}{start, end},
		Logic:    "AND",
	})
	return qb
}

// WhereNull 添加IS NULL条件
func (qb *queryBuilder) WhereNull(column string) QueryBuilder {
	qb.conditions = append(qb.conditions, QueryCondition{
		Column:   column,
		Operator: "IS NULL",
		Logic:    "AND",
	})
	return qb
}

// WhereNotNull 添加IS NOT NULL条件
func (qb *queryBuilder) WhereNotNull(column string) QueryBuilder {
	qb.conditions = append(qb.conditions, QueryCondition{
		Column:   column,
		Operator: "IS NOT NULL",
		Logic:    "AND",
	})
	return qb
}

// OrderBy 添加排序
func (qb *queryBuilder) OrderBy(column string, direction ...string) QueryBuilder {
	dir := "ASC"
	if len(direction) > 0 {
		dir = strings.ToUpper(direction[0])
	}
	qb.orders = append(qb.orders, OrderClause{
		Column:    column,
		Direction: dir,
	})
	return qb
}

// GroupBy 添加分组
func (qb *queryBuilder) GroupBy(columns ...string) QueryBuilder {
	qb.groups = append(qb.groups, columns...)
	return qb
}

// Having 添加HAVING条件
func (qb *queryBuilder) Having(condition string, args ...interface{}) QueryBuilder {
	qb.havings = append(qb.havings, HavingClause{
		Condition: condition,
		Args:      args,
	})
	return qb
}

// Limit 设置限制数量
func (qb *queryBuilder) Limit(limit int) QueryBuilder {
	qb.limitNum = limit
	return qb
}

// Offset 设置偏移量
func (qb *queryBuilder) Offset(offset int) QueryBuilder {
	qb.offsetNum = offset
	return qb
}

// Join 添加JOIN
func (qb *queryBuilder) Join(table, condition string) QueryBuilder {
	qb.joins = append(qb.joins, JoinClause{
		Type:      "INNER",
		Table:     table,
		Condition: condition,
	})
	return qb
}

// LeftJoin 添加LEFT JOIN
func (qb *queryBuilder) LeftJoin(table, condition string) QueryBuilder {
	qb.joins = append(qb.joins, JoinClause{
		Type:      "LEFT",
		Table:     table,
		Condition: condition,
	})
	return qb
}

// RightJoin 添加RIGHT JOIN
func (qb *queryBuilder) RightJoin(table, condition string) QueryBuilder {
	qb.joins = append(qb.joins, JoinClause{
		Type:      "RIGHT",
		Table:     table,
		Condition: condition,
	})
	return qb
}

// InnerJoin 添加INNER JOIN
func (qb *queryBuilder) InnerJoin(table, condition string) QueryBuilder {
	qb.joins = append(qb.joins, JoinClause{
		Type:      "INNER",
		Table:     table,
		Condition: condition,
	})
	return qb
}

// Get 获取多条记录
func (qb *queryBuilder) Get(dest interface{}) error {
	query, args := qb.buildSelectSQL()

	var rows *sql.Rows
	var err error

	if qb.tx != nil {
		rows, err = qb.tx.Query(query, args...)
	} else {
		rows, err = qb.orm.Query(query, args...)
	}

	if err != nil {
		return err
	}
	defer rows.Close()

	return scanRows(rows, dest)
}

// First 获取第一条记录
func (qb *queryBuilder) First(dest interface{}) error {
	qb.limitNum = 1
	query, args := qb.buildSelectSQL()

	var row *sql.Row

	if qb.tx != nil {
		row = qb.tx.QueryRow(query, args...)
	} else {
		row = qb.orm.QueryRow(query, args...)
	}

	return scanRow(row, dest)
}

// Find 查找记录（别名）
func (qb *queryBuilder) Find(dest interface{}) error {
	return qb.Get(dest)
}

// Count 统计记录数
func (qb *queryBuilder) Count() (int64, error) {
	query, args := qb.buildCountSQL()

	var row *sql.Row

	if qb.tx != nil {
		row = qb.tx.QueryRow(query, args...)
	} else {
		row = qb.orm.QueryRow(query, args...)
	}

	var count int64
	err := row.Scan(&count)
	return count, err
}

// Exists 检查记录是否存在
func (qb *queryBuilder) Exists() (bool, error) {
	count, err := qb.Count()
	return count > 0, err
}

// Insert 插入记录
func (qb *queryBuilder) Insert(data interface{}) error {
	query, args := qb.buildInsertSQL(data)

	if qb.tx != nil {
		_, err := qb.tx.Exec(query, args...)
		return err
	} else {
		_, err := qb.orm.Exec(query, args...)
		return err
	}
}

// InsertBatch 批量插入记录
func (qb *queryBuilder) InsertBatch(data interface{}) error {
	query, args := qb.buildBatchInsertSQL(data)

	if qb.tx != nil {
		_, err := qb.tx.Exec(query, args...)
		return err
	} else {
		_, err := qb.orm.Exec(query, args...)
		return err
	}
}

// Update 更新记录
func (qb *queryBuilder) Update(data interface{}) error {
	query, args := qb.buildUpdateSQL(data)

	if qb.tx != nil {
		_, err := qb.tx.Exec(query, args...)
		return err
	} else {
		_, err := qb.orm.Exec(query, args...)
		return err
	}
}

// UpdateColumns 更新指定列
func (qb *queryBuilder) UpdateColumns(columns map[string]interface{}) error {
	query, args := qb.buildUpdateColumnsSQL(columns)

	if qb.tx != nil {
		_, err := qb.tx.Exec(query, args...)
		return err
	} else {
		_, err := qb.orm.Exec(query, args...)
		return err
	}
}

// Delete 删除记录
func (qb *queryBuilder) Delete() error {
	query, args := qb.buildDeleteSQL()

	if qb.tx != nil {
		_, err := qb.tx.Exec(query, args...)
		return err
	} else {
		_, err := qb.orm.Exec(query, args...)
		return err
	}
}

// ToSQL 构建SQL语句
func (qb *queryBuilder) ToSQL() (string, []interface{}) {
	return qb.buildSelectSQL()
}

// buildSelectSQL 构建SELECT SQL
func (qb *queryBuilder) buildSelectSQL() (string, []interface{}) {
	var parts []string
	var args []interface{}

	// SELECT子句
	if len(qb.selectCols) > 0 {
		parts = append(parts, "SELECT "+strings.Join(qb.selectCols, ", "))
	} else {
		parts = append(parts, "SELECT *")
	}

	// FROM子句
	parts = append(parts, "FROM "+qb.tableName)

	// JOIN子句
	for _, join := range qb.joins {
		parts = append(parts, fmt.Sprintf("%s JOIN %s ON %s", join.Type, join.Table, join.Condition))
	}

	// WHERE子句
	if len(qb.conditions) > 0 {
		whereClause, whereArgs := qb.buildWhereClause()
		parts = append(parts, "WHERE "+whereClause)
		args = append(args, whereArgs...)
	}

	// GROUP BY子句
	if len(qb.groups) > 0 {
		parts = append(parts, "GROUP BY "+strings.Join(qb.groups, ", "))
	}

	// HAVING子句
	if len(qb.havings) > 0 {
		havingClause, havingArgs := qb.buildHavingClause()
		parts = append(parts, "HAVING "+havingClause)
		args = append(args, havingArgs...)
	}

	// ORDER BY子句
	if len(qb.orders) > 0 {
		var orderParts []string
		for _, order := range qb.orders {
			orderParts = append(orderParts, order.Column+" "+order.Direction)
		}
		parts = append(parts, "ORDER BY "+strings.Join(orderParts, ", "))
	}

	// LIMIT子句
	if qb.limitNum > 0 {
		parts = append(parts, fmt.Sprintf("LIMIT %d", qb.limitNum))
	}

	// OFFSET子句
	if qb.offsetNum > 0 {
		parts = append(parts, fmt.Sprintf("OFFSET %d", qb.offsetNum))
	}

	return strings.Join(parts, " "), args
}

// buildCountSQL 构建COUNT SQL
func (qb *queryBuilder) buildCountSQL() (string, []interface{}) {
	var parts []string
	var args []interface{}

	parts = append(parts, "SELECT COUNT(*)")
	parts = append(parts, "FROM "+qb.tableName)

	// JOIN子句
	for _, join := range qb.joins {
		parts = append(parts, fmt.Sprintf("%s JOIN %s ON %s", join.Type, join.Table, join.Condition))
	}

	// WHERE子句
	if len(qb.conditions) > 0 {
		whereClause, whereArgs := qb.buildWhereClause()
		parts = append(parts, "WHERE "+whereClause)
		args = append(args, whereArgs...)
	}

	return strings.Join(parts, " "), args
}

// buildInsertSQL 构建INSERT SQL
func (qb *queryBuilder) buildInsertSQL(data interface{}) (string, []interface{}) {
	columns, values := qb.extractColumnsAndValues(data)

	placeholders := make([]string, len(values))
	for i := range placeholders {
		placeholders[i] = "?"
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		qb.tableName,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "))

	return query, values
}

// buildBatchInsertSQL 构建批量INSERT SQL
func (qb *queryBuilder) buildBatchInsertSQL(data interface{}) (string, []interface{}) {
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Slice {
		return qb.buildInsertSQL(data)
	}

	if v.Len() == 0 {
		return "", nil
	}

	// 获取第一个元素的列名
	firstItem := v.Index(0).Interface()
	columns, _ := qb.extractColumnsAndValues(firstItem)

	var allValues []interface{}
	var valuePlaceholders []string

	for i := 0; i < v.Len(); i++ {
		item := v.Index(i).Interface()
		_, values := qb.extractColumnsAndValues(item)
		allValues = append(allValues, values...)

		placeholders := make([]string, len(values))
		for j := range placeholders {
			placeholders[j] = "?"
		}
		valuePlaceholders = append(valuePlaceholders, "("+strings.Join(placeholders, ", ")+")")
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s",
		qb.tableName,
		strings.Join(columns, ", "),
		strings.Join(valuePlaceholders, ", "))

	return query, allValues
}

// buildUpdateSQL 构建UPDATE SQL
func (qb *queryBuilder) buildUpdateSQL(data interface{}) (string, []interface{}) {
	columns, values := qb.extractColumnsAndValues(data)

	var setParts []string
	for _, col := range columns {
		setParts = append(setParts, col+" = ?")
	}

	var parts []string
	var args []interface{}

	parts = append(parts, "UPDATE "+qb.tableName)
	parts = append(parts, "SET "+strings.Join(setParts, ", "))
	args = append(args, values...)

	// WHERE子句
	if len(qb.conditions) > 0 {
		whereClause, whereArgs := qb.buildWhereClause()
		parts = append(parts, "WHERE "+whereClause)
		args = append(args, whereArgs...)
	}

	return strings.Join(parts, " "), args
}

// buildUpdateColumnsSQL 构建UPDATE指定列SQL
func (qb *queryBuilder) buildUpdateColumnsSQL(columns map[string]interface{}) (string, []interface{}) {
	var setParts []string
	var values []interface{}

	for col, val := range columns {
		setParts = append(setParts, col+" = ?")
		values = append(values, val)
	}

	var parts []string
	var args []interface{}

	parts = append(parts, "UPDATE "+qb.tableName)
	parts = append(parts, "SET "+strings.Join(setParts, ", "))
	args = append(args, values...)

	// WHERE子句
	if len(qb.conditions) > 0 {
		whereClause, whereArgs := qb.buildWhereClause()
		parts = append(parts, "WHERE "+whereClause)
		args = append(args, whereArgs...)
	}

	return strings.Join(parts, " "), args
}

// buildDeleteSQL 构建DELETE SQL
func (qb *queryBuilder) buildDeleteSQL() (string, []interface{}) {
	var parts []string
	var args []interface{}

	parts = append(parts, "DELETE FROM "+qb.tableName)

	// WHERE子句
	if len(qb.conditions) > 0 {
		whereClause, whereArgs := qb.buildWhereClause()
		parts = append(parts, "WHERE "+whereClause)
		args = append(args, whereArgs...)
	}

	return strings.Join(parts, " "), args
}

// buildWhereClause 构建WHERE子句
func (qb *queryBuilder) buildWhereClause() (string, []interface{}) {
	var parts []string
	var args []interface{}

	for i, condition := range qb.conditions {
		if i > 0 {
			parts = append(parts, condition.Logic)
		}

		switch condition.Operator {
		case "IN", "NOT IN":
			placeholders := make([]string, len(condition.Values))
			for j := range placeholders {
				placeholders[j] = "?"
			}
			parts = append(parts, fmt.Sprintf("%s %s (%s)",
				condition.Column, condition.Operator, strings.Join(placeholders, ", ")))
			args = append(args, condition.Values...)
		case "BETWEEN":
			parts = append(parts, fmt.Sprintf("%s BETWEEN ? AND ?", condition.Column))
			args = append(args, condition.Values...)
		case "IS NULL", "IS NOT NULL":
			parts = append(parts, fmt.Sprintf("%s %s", condition.Column, condition.Operator))
		default:
			if condition.Value != nil {
				if values, ok := condition.Value.([]interface{}); ok && len(values) > 0 {
					// 处理复杂条件，如 "name = ? AND age > ?"
					parts = append(parts, condition.Column)
					args = append(args, values...)
				} else {
					parts = append(parts, fmt.Sprintf("%s %s ?", condition.Column, condition.Operator))
					args = append(args, condition.Value)
				}
			}
		}
	}

	return strings.Join(parts, " "), args
}

// buildHavingClause 构建HAVING子句
func (qb *queryBuilder) buildHavingClause() (string, []interface{}) {
	var parts []string
	var args []interface{}

	for i, having := range qb.havings {
		if i > 0 {
			parts = append(parts, "AND")
		}
		parts = append(parts, having.Condition)
		args = append(args, having.Args...)
	}

	return strings.Join(parts, " "), args
}
