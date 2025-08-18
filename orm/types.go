package orm

import (
	"context"
	"database/sql"
	"time"
)

// DatabaseType 数据库类型
type DatabaseType string

const (
	MySQL      DatabaseType = "mysql"
	PostgreSQL DatabaseType = "postgres"
	SQLite     DatabaseType = "sqlite3"
	SQLServer  DatabaseType = "sqlserver"
)

// Config 数据库配置
type Config struct {
	Type         DatabaseType  `json:"type" yaml:"type"`
	Host         string        `json:"host" yaml:"host"`
	Port         int           `json:"port" yaml:"port"`
	Username     string        `json:"username" yaml:"username"`
	Password     string        `json:"password" yaml:"password"`
	Database     string        `json:"database" yaml:"database"`
	SSLMode      string        `json:"ssl_mode" yaml:"ssl_mode"`
	Charset      string        `json:"charset" yaml:"charset"`
	Timezone     string        `json:"timezone" yaml:"timezone"`
	MaxOpenConns int           `json:"max_open_conns" yaml:"max_open_conns"`
	MaxIdleConns int           `json:"max_idle_conns" yaml:"max_idle_conns"`
	MaxLifetime  time.Duration `json:"max_lifetime" yaml:"max_lifetime"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Type:         MySQL,
		Host:         "localhost",
		Port:         3306,
		Charset:      "utf8mb4",
		SSLMode:      "disable",
		MaxOpenConns: 100,
		MaxIdleConns: 10,
		MaxLifetime:  time.Hour,
	}
}

// DB 数据库接口
type DB interface {
	// 连接管理
	Connect() error
	Close() error
	Ping() error

	// 查询操作
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)

	// 事务操作
	Begin() (Tx, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (Tx, error)

	// 获取原始连接
	Raw() *sql.DB
}

// Tx 事务接口
type Tx interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
	Commit() error
	Rollback() error
	Table(tableName string) QueryBuilder
	Model(model interface{}) QueryBuilder
}

// ModelInterface 模型接口
type ModelInterface interface {
	TableName() string
}

// QueryBuilder 查询构建器接口
type QueryBuilder interface {
	// SELECT 操作
	Select(columns ...string) QueryBuilder
	From(table string) QueryBuilder
	Where(condition string, args ...interface{}) QueryBuilder
	WhereIn(column string, values ...interface{}) QueryBuilder
	WhereNotIn(column string, values ...interface{}) QueryBuilder
	WhereBetween(column string, start, end interface{}) QueryBuilder
	WhereNull(column string) QueryBuilder
	WhereNotNull(column string) QueryBuilder
	OrderBy(column string, direction ...string) QueryBuilder
	GroupBy(columns ...string) QueryBuilder
	Having(condition string, args ...interface{}) QueryBuilder
	Limit(limit int) QueryBuilder
	Offset(offset int) QueryBuilder
	Join(table, condition string) QueryBuilder
	LeftJoin(table, condition string) QueryBuilder
	RightJoin(table, condition string) QueryBuilder
	InnerJoin(table, condition string) QueryBuilder

	// 执行查询
	Get(dest interface{}) error
	First(dest interface{}) error
	Find(dest interface{}) error
	Count() (int64, error)
	Exists() (bool, error)

	// INSERT 操作
	Insert(data interface{}) error
	InsertBatch(data interface{}) error

	// UPDATE 操作
	Update(data interface{}) error
	UpdateColumns(columns map[string]interface{}) error

	// DELETE 操作
	Delete() error

	// 构建SQL
	ToSQL() (string, []interface{})
}

// Migration 迁移接口
type Migration interface {
	Up() error
	Down() error
	Version() string
}

// Schema 表结构接口
type Schema interface {
	CreateTable(tableName string, callback func(TableInterface)) error
	DropTable(tableName string) error
	AlterTable(tableName string, callback func(TableInterface)) error
	HasTable(tableName string) (bool, error)
	HasColumn(tableName, columnName string) (bool, error)
}

// TableInterface 表定义接口
type TableInterface interface {
	ID() TableInterface
	String(name string, length ...int) TableInterface
	Text(name string) TableInterface
	Integer(name string) TableInterface
	BigInteger(name string) TableInterface
	Float(name string, precision, scale int) TableInterface
	Double(name string) TableInterface
	Decimal(name string, precision, scale int) TableInterface
	Boolean(name string) TableInterface
	Date(name string) TableInterface
	DateTime(name string) TableInterface
	Timestamp(name string) TableInterface
	JSON(name string) TableInterface

	// 约束
	Primary(columns ...string) TableInterface
	Index(name string, columns ...string) TableInterface
	Unique(name string, columns ...string) TableInterface
	Foreign(column, references string) TableInterface

	// 修饰符
	Nullable() TableInterface
	NotNull() TableInterface
	Default(value interface{}) TableInterface
	Comment(comment string) TableInterface
	AutoIncrement() TableInterface
}

// FieldTag 字段标签
type FieldTag struct {
	Column        string `json:"column"`
	Type          string `json:"type"`
	Size          int    `json:"size"`
	Primary       bool   `json:"primary"`
	AutoIncrement bool   `json:"auto_increment"`
	NotNull       bool   `json:"not_null"`
	Unique        bool   `json:"unique"`
	Index         string `json:"index"`
	Default       string `json:"default"`
	Comment       string `json:"comment"`
	ForeignKey    string `json:"foreign_key"`
	References    string `json:"references"`
}

// QueryCondition 查询条件
type QueryCondition struct {
	Column   string        `json:"column"`
	Operator string        `json:"operator"`
	Value    interface{}   `json:"value"`
	Values   []interface{} `json:"values"`
	Logic    string        `json:"logic"` // AND, OR
}

// JoinClause JOIN子句
type JoinClause struct {
	Type      string `json:"type"` // INNER, LEFT, RIGHT, FULL
	Table     string `json:"table"`
	Condition string `json:"condition"`
}

// OrderClause 排序子句
type OrderClause struct {
	Column    string `json:"column"`
	Direction string `json:"direction"` // ASC, DESC
}

// GroupClause 分组子句
type GroupClause struct {
	Columns []string `json:"columns"`
}

// HavingClause HAVING子句
type HavingClause struct {
	Condition string        `json:"condition"`
	Args      []interface{} `json:"args"`
}

// LimitClause 限制子句
type LimitClause struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}
