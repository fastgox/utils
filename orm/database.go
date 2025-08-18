package orm

import (
	"fmt"
	"reflect"
	"strings"
)

// DatabaseManager 数据库管理器
type DatabaseManager struct {
	orm *ORM
}

// NewDatabaseManager 创建数据库管理器
func NewDatabaseManager(orm *ORM) *DatabaseManager {
	return &DatabaseManager{orm: orm}
}

// GetDialect 获取数据库方言
func (dm *DatabaseManager) GetDialect() Dialect {
	switch dm.orm.config.Type {
	case MySQL:
		return &MySQLDialect{}
	case PostgreSQL:
		return &PostgreSQLDialect{}
	case SQLite:
		return &SQLiteDialect{}
	case SQLServer:
		return &SQLServerDialect{}
	default:
		return &MySQLDialect{} // 默认使用MySQL方言
	}
}

// Dialect 数据库方言接口
type Dialect interface {
	Quote(name string) string
	QuoteString(s string) string
	DataType(fieldType reflect.Type, size int) string
	AutoIncrement() string
	PrimaryKey() string
	CreateTableSQL(tableName string, columns []ColumnDefinition) string
	DropTableSQL(tableName string) string
	AddColumnSQL(tableName, columnName string, definition ColumnDefinition) string
	DropColumnSQL(tableName, columnName string) string
	CreateIndexSQL(tableName, indexName string, columns []string, unique bool) string
	DropIndexSQL(tableName, indexName string) string
}

// ColumnDefinition 列定义
type ColumnDefinition struct {
	Name          string
	Type          string
	Size          int
	Precision     int
	Scale         int
	NotNull       bool
	Primary       bool
	AutoIncrement bool
	Unique        bool
	Default       interface{}
	Comment       string
	ForeignKey    string
	References    string
}

// MySQLDialect MySQL方言
type MySQLDialect struct{}

func (d *MySQLDialect) Quote(name string) string {
	return "`" + name + "`"
}

func (d *MySQLDialect) QuoteString(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "''") + "'"
}

func (d *MySQLDialect) DataType(fieldType reflect.Type, size int) string {
	switch fieldType.Kind() {
	case reflect.Bool:
		return "BOOLEAN"
	case reflect.Int, reflect.Int32:
		return "INT"
	case reflect.Int64:
		return "BIGINT"
	case reflect.Float32:
		return "FLOAT"
	case reflect.Float64:
		return "DOUBLE"
	case reflect.String:
		if size > 0 && size <= 255 {
			return fmt.Sprintf("VARCHAR(%d)", size)
		}
		return "TEXT"
	default:
		if fieldType.String() == "time.Time" {
			return "DATETIME"
		}
		return "TEXT"
	}
}

func (d *MySQLDialect) AutoIncrement() string {
	return "AUTO_INCREMENT"
}

func (d *MySQLDialect) PrimaryKey() string {
	return "PRIMARY KEY"
}

func (d *MySQLDialect) CreateTableSQL(tableName string, columns []ColumnDefinition) string {
	var parts []string
	var primaryKeys []string

	for _, col := range columns {
		part := d.Quote(col.Name) + " " + col.Type

		if col.NotNull {
			part += " NOT NULL"
		}

		if col.AutoIncrement {
			part += " " + d.AutoIncrement()
		}

		if col.Default != nil {
			part += " DEFAULT " + fmt.Sprintf("%v", col.Default)
		}

		if col.Comment != "" {
			part += " COMMENT " + d.QuoteString(col.Comment)
		}

		parts = append(parts, part)

		if col.Primary {
			primaryKeys = append(primaryKeys, d.Quote(col.Name))
		}
	}

	if len(primaryKeys) > 0 {
		parts = append(parts, d.PrimaryKey()+" ("+strings.Join(primaryKeys, ", ")+")")
	}

	return fmt.Sprintf("CREATE TABLE %s (%s)", d.Quote(tableName), strings.Join(parts, ", "))
}

func (d *MySQLDialect) DropTableSQL(tableName string) string {
	return fmt.Sprintf("DROP TABLE IF EXISTS %s", d.Quote(tableName))
}

func (d *MySQLDialect) AddColumnSQL(tableName, columnName string, definition ColumnDefinition) string {
	return fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s",
		d.Quote(tableName), d.Quote(columnName), definition.Type)
}

func (d *MySQLDialect) DropColumnSQL(tableName, columnName string) string {
	return fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s", d.Quote(tableName), d.Quote(columnName))
}

func (d *MySQLDialect) CreateIndexSQL(tableName, indexName string, columns []string, unique bool) string {
	indexType := "INDEX"
	if unique {
		indexType = "UNIQUE INDEX"
	}

	quotedColumns := make([]string, len(columns))
	for i, col := range columns {
		quotedColumns[i] = d.Quote(col)
	}

	return fmt.Sprintf("CREATE %s %s ON %s (%s)",
		indexType, d.Quote(indexName), d.Quote(tableName), strings.Join(quotedColumns, ", "))
}

func (d *MySQLDialect) DropIndexSQL(tableName, indexName string) string {
	return fmt.Sprintf("DROP INDEX %s ON %s", d.Quote(indexName), d.Quote(tableName))
}

// PostgreSQLDialect PostgreSQL方言
type PostgreSQLDialect struct{}

func (d *PostgreSQLDialect) Quote(name string) string {
	return `"` + name + `"`
}

func (d *PostgreSQLDialect) QuoteString(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "''") + "'"
}

func (d *PostgreSQLDialect) DataType(fieldType reflect.Type, size int) string {
	switch fieldType.Kind() {
	case reflect.Bool:
		return "BOOLEAN"
	case reflect.Int, reflect.Int32:
		return "INTEGER"
	case reflect.Int64:
		return "BIGINT"
	case reflect.Float32:
		return "REAL"
	case reflect.Float64:
		return "DOUBLE PRECISION"
	case reflect.String:
		if size > 0 {
			return fmt.Sprintf("VARCHAR(%d)", size)
		}
		return "TEXT"
	default:
		if fieldType.String() == "time.Time" {
			return "TIMESTAMP"
		}
		return "TEXT"
	}
}

func (d *PostgreSQLDialect) AutoIncrement() string {
	return "SERIAL"
}

func (d *PostgreSQLDialect) PrimaryKey() string {
	return "PRIMARY KEY"
}

func (d *PostgreSQLDialect) CreateTableSQL(tableName string, columns []ColumnDefinition) string {
	var parts []string
	var primaryKeys []string

	for _, col := range columns {
		part := d.Quote(col.Name) + " " + col.Type

		if col.NotNull {
			part += " NOT NULL"
		}

		if col.Default != nil {
			part += " DEFAULT " + fmt.Sprintf("%v", col.Default)
		}

		parts = append(parts, part)

		if col.Primary {
			primaryKeys = append(primaryKeys, d.Quote(col.Name))
		}
	}

	if len(primaryKeys) > 0 {
		parts = append(parts, d.PrimaryKey()+" ("+strings.Join(primaryKeys, ", ")+")")
	}

	return fmt.Sprintf("CREATE TABLE %s (%s)", d.Quote(tableName), strings.Join(parts, ", "))
}

func (d *PostgreSQLDialect) DropTableSQL(tableName string) string {
	return fmt.Sprintf("DROP TABLE IF EXISTS %s", d.Quote(tableName))
}

func (d *PostgreSQLDialect) AddColumnSQL(tableName, columnName string, definition ColumnDefinition) string {
	return fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s",
		d.Quote(tableName), d.Quote(columnName), definition.Type)
}

func (d *PostgreSQLDialect) DropColumnSQL(tableName, columnName string) string {
	return fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s", d.Quote(tableName), d.Quote(columnName))
}

func (d *PostgreSQLDialect) CreateIndexSQL(tableName, indexName string, columns []string, unique bool) string {
	indexType := ""
	if unique {
		indexType = "UNIQUE "
	}

	quotedColumns := make([]string, len(columns))
	for i, col := range columns {
		quotedColumns[i] = d.Quote(col)
	}

	return fmt.Sprintf("CREATE %sINDEX %s ON %s (%s)",
		indexType, d.Quote(indexName), d.Quote(tableName), strings.Join(quotedColumns, ", "))
}

func (d *PostgreSQLDialect) DropIndexSQL(tableName, indexName string) string {
	return fmt.Sprintf("DROP INDEX IF EXISTS %s", d.Quote(indexName))
}

// SQLiteDialect SQLite方言
type SQLiteDialect struct{}

func (d *SQLiteDialect) Quote(name string) string {
	return "`" + name + "`"
}

func (d *SQLiteDialect) QuoteString(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "''") + "'"
}

func (d *SQLiteDialect) DataType(fieldType reflect.Type, size int) string {
	switch fieldType.Kind() {
	case reflect.Bool:
		return "INTEGER"
	case reflect.Int, reflect.Int32, reflect.Int64:
		return "INTEGER"
	case reflect.Float32, reflect.Float64:
		return "REAL"
	case reflect.String:
		return "TEXT"
	default:
		if fieldType.String() == "time.Time" {
			return "DATETIME"
		}
		return "TEXT"
	}
}

func (d *SQLiteDialect) AutoIncrement() string {
	return "AUTOINCREMENT"
}

func (d *SQLiteDialect) PrimaryKey() string {
	return "PRIMARY KEY"
}

func (d *SQLiteDialect) CreateTableSQL(tableName string, columns []ColumnDefinition) string {
	var parts []string

	for _, col := range columns {
		part := d.Quote(col.Name) + " " + col.Type

		if col.Primary {
			part += " " + d.PrimaryKey()
		}

		if col.AutoIncrement {
			part += " " + d.AutoIncrement()
		}

		if col.NotNull {
			part += " NOT NULL"
		}

		if col.Default != nil {
			part += " DEFAULT " + fmt.Sprintf("%v", col.Default)
		}

		parts = append(parts, part)
	}

	return fmt.Sprintf("CREATE TABLE %s (%s)", d.Quote(tableName), strings.Join(parts, ", "))
}

func (d *SQLiteDialect) DropTableSQL(tableName string) string {
	return fmt.Sprintf("DROP TABLE IF EXISTS %s", d.Quote(tableName))
}

func (d *SQLiteDialect) AddColumnSQL(tableName, columnName string, definition ColumnDefinition) string {
	return fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s",
		d.Quote(tableName), d.Quote(columnName), definition.Type)
}

func (d *SQLiteDialect) DropColumnSQL(tableName, columnName string) string {
	// SQLite不直接支持删除列，需要重建表
	return fmt.Sprintf("-- SQLite不支持直接删除列: %s.%s", tableName, columnName)
}

func (d *SQLiteDialect) CreateIndexSQL(tableName, indexName string, columns []string, unique bool) string {
	indexType := ""
	if unique {
		indexType = "UNIQUE "
	}

	quotedColumns := make([]string, len(columns))
	for i, col := range columns {
		quotedColumns[i] = d.Quote(col)
	}

	return fmt.Sprintf("CREATE %sINDEX %s ON %s (%s)",
		indexType, d.Quote(indexName), d.Quote(tableName), strings.Join(quotedColumns, ", "))
}

func (d *SQLiteDialect) DropIndexSQL(tableName, indexName string) string {
	return fmt.Sprintf("DROP INDEX IF EXISTS %s", d.Quote(indexName))
}

// SQLServerDialect SQL Server方言
type SQLServerDialect struct{}

func (d *SQLServerDialect) Quote(name string) string {
	return "[" + name + "]"
}

func (d *SQLServerDialect) QuoteString(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "''") + "'"
}

func (d *SQLServerDialect) DataType(fieldType reflect.Type, size int) string {
	switch fieldType.Kind() {
	case reflect.Bool:
		return "BIT"
	case reflect.Int, reflect.Int32:
		return "INT"
	case reflect.Int64:
		return "BIGINT"
	case reflect.Float32:
		return "REAL"
	case reflect.Float64:
		return "FLOAT"
	case reflect.String:
		if size > 0 && size <= 4000 {
			return fmt.Sprintf("NVARCHAR(%d)", size)
		}
		return "NTEXT"
	default:
		if fieldType.String() == "time.Time" {
			return "DATETIME2"
		}
		return "NTEXT"
	}
}

func (d *SQLServerDialect) AutoIncrement() string {
	return "IDENTITY(1,1)"
}

func (d *SQLServerDialect) PrimaryKey() string {
	return "PRIMARY KEY"
}

func (d *SQLServerDialect) CreateTableSQL(tableName string, columns []ColumnDefinition) string {
	var parts []string
	var primaryKeys []string

	for _, col := range columns {
		part := d.Quote(col.Name) + " " + col.Type

		if col.AutoIncrement {
			part += " " + d.AutoIncrement()
		}

		if col.NotNull {
			part += " NOT NULL"
		}

		if col.Default != nil {
			part += " DEFAULT " + fmt.Sprintf("%v", col.Default)
		}

		parts = append(parts, part)

		if col.Primary {
			primaryKeys = append(primaryKeys, d.Quote(col.Name))
		}
	}

	if len(primaryKeys) > 0 {
		parts = append(parts, d.PrimaryKey()+" ("+strings.Join(primaryKeys, ", ")+")")
	}

	return fmt.Sprintf("CREATE TABLE %s (%s)", d.Quote(tableName), strings.Join(parts, ", "))
}

func (d *SQLServerDialect) DropTableSQL(tableName string) string {
	return fmt.Sprintf("DROP TABLE IF EXISTS %s", d.Quote(tableName))
}

func (d *SQLServerDialect) AddColumnSQL(tableName, columnName string, definition ColumnDefinition) string {
	return fmt.Sprintf("ALTER TABLE %s ADD %s %s",
		d.Quote(tableName), d.Quote(columnName), definition.Type)
}

func (d *SQLServerDialect) DropColumnSQL(tableName, columnName string) string {
	return fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s", d.Quote(tableName), d.Quote(columnName))
}

func (d *SQLServerDialect) CreateIndexSQL(tableName, indexName string, columns []string, unique bool) string {
	indexType := "NONCLUSTERED INDEX"
	if unique {
		indexType = "UNIQUE NONCLUSTERED INDEX"
	}

	quotedColumns := make([]string, len(columns))
	for i, col := range columns {
		quotedColumns[i] = d.Quote(col)
	}

	return fmt.Sprintf("CREATE %s %s ON %s (%s)",
		indexType, d.Quote(indexName), d.Quote(tableName), strings.Join(quotedColumns, ", "))
}

func (d *SQLServerDialect) DropIndexSQL(tableName, indexName string) string {
	return fmt.Sprintf("DROP INDEX %s ON %s", d.Quote(indexName), d.Quote(tableName))
}
