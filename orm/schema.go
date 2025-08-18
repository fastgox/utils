package orm

import (
	"fmt"
)

// schema 表结构实现
type schema struct {
	orm *ORM
}

// NewSchema 创建表结构管理器
func NewSchema(orm *ORM) Schema {
	return &schema{orm: orm}
}

// CreateTable 创建表
func (s *schema) CreateTable(tableName string, callback func(TableInterface)) error {
	table := NewTableBuilder(tableName, s.orm)
	callback(table)

	sql := table.(*tableBuilder).ToSQL()
	_, err := s.orm.Exec(sql)
	return err
}

// DropTable 删除表
func (s *schema) DropTable(tableName string) error {
	dialect := NewDatabaseManager(s.orm).GetDialect()
	sql := dialect.DropTableSQL(tableName)
	_, err := s.orm.Exec(sql)
	return err
}

// AlterTable 修改表
func (s *schema) AlterTable(tableName string, callback func(TableInterface)) error {
	table := NewTableBuilder(tableName, s.orm)
	table.(*tableBuilder).SetAlterMode(true)
	callback(table)

	sqls := table.(*tableBuilder).ToAlterSQLs()
	for _, sql := range sqls {
		if _, err := s.orm.Exec(sql); err != nil {
			return err
		}
	}
	return nil
}

// HasTable 检查表是否存在
func (s *schema) HasTable(tableName string) (bool, error) {
	var sql string
	switch s.orm.config.Type {
	case MySQL:
		sql = "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?"
	case PostgreSQL:
		sql = "SELECT COUNT(*) FROM information_schema.tables WHERE table_name = ?"
	case SQLite:
		sql = "SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name = ?"
	case SQLServer:
		sql = "SELECT COUNT(*) FROM information_schema.tables WHERE table_name = ?"
	default:
		return false, fmt.Errorf("不支持的数据库类型")
	}

	var count int
	err := s.orm.QueryRow(sql, tableName).Scan(&count)
	return count > 0, err
}

// HasColumn 检查列是否存在
func (s *schema) HasColumn(tableName, columnName string) (bool, error) {
	var sql string
	switch s.orm.config.Type {
	case MySQL:
		sql = "SELECT COUNT(*) FROM information_schema.columns WHERE table_schema = DATABASE() AND table_name = ? AND column_name = ?"
	case PostgreSQL:
		sql = "SELECT COUNT(*) FROM information_schema.columns WHERE table_name = ? AND column_name = ?"
	case SQLite:
		sql = "SELECT COUNT(*) FROM pragma_table_info(?) WHERE name = ?"
	case SQLServer:
		sql = "SELECT COUNT(*) FROM information_schema.columns WHERE table_name = ? AND column_name = ?"
	default:
		return false, fmt.Errorf("不支持的数据库类型")
	}

	var count int
	err := s.orm.QueryRow(sql, tableName, columnName).Scan(&count)
	return count > 0, err
}

// tableBuilder 表构建器实现
type tableBuilder struct {
	tableName string
	orm       *ORM
	columns   []ColumnDefinition
	indexes   []IndexDefinition
	alterMode bool
	alterOps  []AlterOperation
}

// IndexDefinition 索引定义
type IndexDefinition struct {
	Name    string
	Columns []string
	Unique  bool
}

// AlterOperation 修改操作
type AlterOperation struct {
	Type string // ADD_COLUMN, DROP_COLUMN, ADD_INDEX, DROP_INDEX
	Data interface{}
}

// NewTableBuilder 创建表构建器
func NewTableBuilder(tableName string, orm *ORM) TableInterface {
	return &tableBuilder{
		tableName: tableName,
		orm:       orm,
		columns:   make([]ColumnDefinition, 0),
		indexes:   make([]IndexDefinition, 0),
		alterOps:  make([]AlterOperation, 0),
	}
}

// SetAlterMode 设置修改模式
func (tb *tableBuilder) SetAlterMode(alter bool) {
	tb.alterMode = alter
}

// ID 添加ID主键列
func (tb *tableBuilder) ID() TableInterface {
	tb.columns = append(tb.columns, ColumnDefinition{
		Name:          "id",
		Type:          "BIGINT",
		Primary:       true,
		AutoIncrement: true,
		NotNull:       true,
	})
	return tb
}

// String 添加字符串列
func (tb *tableBuilder) String(name string, length ...int) TableInterface {
	size := 255
	if len(length) > 0 {
		size = length[0]
	}

	columnType := "VARCHAR"
	if size > 0 {
		columnType = fmt.Sprintf("VARCHAR(%d)", size)
	}

	col := ColumnDefinition{
		Name: name,
		Type: columnType,
		Size: size,
	}

	if tb.alterMode {
		tb.alterOps = append(tb.alterOps, AlterOperation{
			Type: "ADD_COLUMN",
			Data: col,
		})
	} else {
		tb.columns = append(tb.columns, col)
	}

	return tb
}

// Text 添加文本列
func (tb *tableBuilder) Text(name string) TableInterface {
	col := ColumnDefinition{
		Name: name,
		Type: "TEXT",
	}

	if tb.alterMode {
		tb.alterOps = append(tb.alterOps, AlterOperation{
			Type: "ADD_COLUMN",
			Data: col,
		})
	} else {
		tb.columns = append(tb.columns, col)
	}

	return tb
}

// Integer 添加整数列
func (tb *tableBuilder) Integer(name string) TableInterface {
	col := ColumnDefinition{
		Name: name,
		Type: "INT",
	}

	if tb.alterMode {
		tb.alterOps = append(tb.alterOps, AlterOperation{
			Type: "ADD_COLUMN",
			Data: col,
		})
	} else {
		tb.columns = append(tb.columns, col)
	}

	return tb
}

// BigInteger 添加大整数列
func (tb *tableBuilder) BigInteger(name string) TableInterface {
	col := ColumnDefinition{
		Name: name,
		Type: "BIGINT",
	}

	if tb.alterMode {
		tb.alterOps = append(tb.alterOps, AlterOperation{
			Type: "ADD_COLUMN",
			Data: col,
		})
	} else {
		tb.columns = append(tb.columns, col)
	}

	return tb
}

// Float 添加浮点数列
func (tb *tableBuilder) Float(name string, precision, scale int) TableInterface {
	col := ColumnDefinition{
		Name:      name,
		Type:      fmt.Sprintf("FLOAT(%d,%d)", precision, scale),
		Precision: precision,
		Scale:     scale,
	}

	if tb.alterMode {
		tb.alterOps = append(tb.alterOps, AlterOperation{
			Type: "ADD_COLUMN",
			Data: col,
		})
	} else {
		tb.columns = append(tb.columns, col)
	}

	return tb
}

// Double 添加双精度浮点数列
func (tb *tableBuilder) Double(name string) TableInterface {
	col := ColumnDefinition{
		Name: name,
		Type: "DOUBLE",
	}

	if tb.alterMode {
		tb.alterOps = append(tb.alterOps, AlterOperation{
			Type: "ADD_COLUMN",
			Data: col,
		})
	} else {
		tb.columns = append(tb.columns, col)
	}

	return tb
}

// Decimal 添加精确小数列
func (tb *tableBuilder) Decimal(name string, precision, scale int) TableInterface {
	col := ColumnDefinition{
		Name:      name,
		Type:      fmt.Sprintf("DECIMAL(%d,%d)", precision, scale),
		Precision: precision,
		Scale:     scale,
	}

	if tb.alterMode {
		tb.alterOps = append(tb.alterOps, AlterOperation{
			Type: "ADD_COLUMN",
			Data: col,
		})
	} else {
		tb.columns = append(tb.columns, col)
	}

	return tb
}

// Boolean 添加布尔列
func (tb *tableBuilder) Boolean(name string) TableInterface {
	col := ColumnDefinition{
		Name: name,
		Type: "BOOLEAN",
	}

	if tb.alterMode {
		tb.alterOps = append(tb.alterOps, AlterOperation{
			Type: "ADD_COLUMN",
			Data: col,
		})
	} else {
		tb.columns = append(tb.columns, col)
	}

	return tb
}

// Date 添加日期列
func (tb *tableBuilder) Date(name string) TableInterface {
	col := ColumnDefinition{
		Name: name,
		Type: "DATE",
	}

	if tb.alterMode {
		tb.alterOps = append(tb.alterOps, AlterOperation{
			Type: "ADD_COLUMN",
			Data: col,
		})
	} else {
		tb.columns = append(tb.columns, col)
	}

	return tb
}

// DateTime 添加日期时间列
func (tb *tableBuilder) DateTime(name string) TableInterface {
	col := ColumnDefinition{
		Name: name,
		Type: "DATETIME",
	}

	if tb.alterMode {
		tb.alterOps = append(tb.alterOps, AlterOperation{
			Type: "ADD_COLUMN",
			Data: col,
		})
	} else {
		tb.columns = append(tb.columns, col)
	}

	return tb
}

// Timestamp 添加时间戳列
func (tb *tableBuilder) Timestamp(name string) TableInterface {
	col := ColumnDefinition{
		Name: name,
		Type: "TIMESTAMP",
	}

	if tb.alterMode {
		tb.alterOps = append(tb.alterOps, AlterOperation{
			Type: "ADD_COLUMN",
			Data: col,
		})
	} else {
		tb.columns = append(tb.columns, col)
	}

	return tb
}

// JSON 添加JSON列
func (tb *tableBuilder) JSON(name string) TableInterface {
	col := ColumnDefinition{
		Name: name,
		Type: "JSON",
	}

	if tb.alterMode {
		tb.alterOps = append(tb.alterOps, AlterOperation{
			Type: "ADD_COLUMN",
			Data: col,
		})
	} else {
		tb.columns = append(tb.columns, col)
	}

	return tb
}

// Primary 设置主键
func (tb *tableBuilder) Primary(columns ...string) TableInterface {
	// 标记指定列为主键
	for _, colName := range columns {
		for i := range tb.columns {
			if tb.columns[i].Name == colName {
				tb.columns[i].Primary = true
			}
		}
	}
	return tb
}

// Index 添加索引
func (tb *tableBuilder) Index(name string, columns ...string) TableInterface {
	index := IndexDefinition{
		Name:    name,
		Columns: columns,
		Unique:  false,
	}

	if tb.alterMode {
		tb.alterOps = append(tb.alterOps, AlterOperation{
			Type: "ADD_INDEX",
			Data: index,
		})
	} else {
		tb.indexes = append(tb.indexes, index)
	}

	return tb
}

// Unique 添加唯一索引
func (tb *tableBuilder) Unique(name string, columns ...string) TableInterface {
	index := IndexDefinition{
		Name:    name,
		Columns: columns,
		Unique:  true,
	}

	if tb.alterMode {
		tb.alterOps = append(tb.alterOps, AlterOperation{
			Type: "ADD_INDEX",
			Data: index,
		})
	} else {
		tb.indexes = append(tb.indexes, index)
	}

	return tb
}

// Foreign 添加外键约束
func (tb *tableBuilder) Foreign(column, references string) TableInterface {
	// 找到对应的列并设置外键
	for i := range tb.columns {
		if tb.columns[i].Name == column {
			tb.columns[i].ForeignKey = column
			tb.columns[i].References = references
		}
	}
	return tb
}

// Nullable 设置列可为空
func (tb *tableBuilder) Nullable() TableInterface {
	if len(tb.columns) > 0 {
		tb.columns[len(tb.columns)-1].NotNull = false
	}
	return tb
}

// NotNull 设置列不可为空
func (tb *tableBuilder) NotNull() TableInterface {
	if len(tb.columns) > 0 {
		tb.columns[len(tb.columns)-1].NotNull = true
	}
	return tb
}

// Default 设置默认值
func (tb *tableBuilder) Default(value interface{}) TableInterface {
	if len(tb.columns) > 0 {
		tb.columns[len(tb.columns)-1].Default = value
	}
	return tb
}

// Comment 设置注释
func (tb *tableBuilder) Comment(comment string) TableInterface {
	if len(tb.columns) > 0 {
		tb.columns[len(tb.columns)-1].Comment = comment
	}
	return tb
}

// AutoIncrement 设置自增
func (tb *tableBuilder) AutoIncrement() TableInterface {
	if len(tb.columns) > 0 {
		tb.columns[len(tb.columns)-1].AutoIncrement = true
	}
	return tb
}

// ToSQL 生成创建表SQL
func (tb *tableBuilder) ToSQL() string {
	dialect := NewDatabaseManager(tb.orm).GetDialect()
	return dialect.CreateTableSQL(tb.tableName, tb.columns)
}

// ToAlterSQLs 生成修改表SQL
func (tb *tableBuilder) ToAlterSQLs() []string {
	var sqls []string
	dialect := NewDatabaseManager(tb.orm).GetDialect()

	for _, op := range tb.alterOps {
		switch op.Type {
		case "ADD_COLUMN":
			if col, ok := op.Data.(ColumnDefinition); ok {
				sql := dialect.AddColumnSQL(tb.tableName, col.Name, col)
				sqls = append(sqls, sql)
			}
		case "DROP_COLUMN":
			if colName, ok := op.Data.(string); ok {
				sql := dialect.DropColumnSQL(tb.tableName, colName)
				sqls = append(sqls, sql)
			}
		case "ADD_INDEX":
			if index, ok := op.Data.(IndexDefinition); ok {
				sql := dialect.CreateIndexSQL(tb.tableName, index.Name, index.Columns, index.Unique)
				sqls = append(sqls, sql)
			}
		case "DROP_INDEX":
			if indexName, ok := op.Data.(string); ok {
				sql := dialect.DropIndexSQL(tb.tableName, indexName)
				sqls = append(sqls, sql)
			}
		}
	}

	return sqls
}

// DropColumn 删除列（用于ALTER模式）
func (tb *tableBuilder) DropColumn(columnName string) TableInterface {
	if tb.alterMode {
		tb.alterOps = append(tb.alterOps, AlterOperation{
			Type: "DROP_COLUMN",
			Data: columnName,
		})
	}
	return tb
}

// DropIndex 删除索引（用于ALTER模式）
func (tb *tableBuilder) DropIndex(indexName string) TableInterface {
	if tb.alterMode {
		tb.alterOps = append(tb.alterOps, AlterOperation{
			Type: "DROP_INDEX",
			Data: indexName,
		})
	}
	return tb
}
