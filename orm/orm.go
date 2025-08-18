package orm

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"sync"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var (
	// 全局ORM实例
	globalORM *ORM
	once      sync.Once
)

// ORM 主要结构体
type ORM struct {
	config *Config
	db     *sql.DB
	mu     sync.RWMutex
}

// New 创建新的ORM实例
func New(config *Config) *ORM {
	if config == nil {
		config = DefaultConfig()
	}

	return &ORM{
		config: config,
	}
}

// Init 初始化全局ORM实例
func Init(config *Config) error {
	var err error
	once.Do(func() {
		globalORM = New(config)
		err = globalORM.Connect()
	})
	return err
}

// GetGlobalORM 获取全局ORM实例
func GetGlobalORM() *ORM {
	if globalORM == nil {
		panic("ORM未初始化，请先调用Init()方法")
	}
	return globalORM
}

// Connect 连接数据库
func (o *ORM) Connect() error {
	o.mu.Lock()
	defer o.mu.Unlock()

	dsn, err := o.buildDSN()
	if err != nil {
		return fmt.Errorf("构建DSN失败: %w", err)
	}

	db, err := sql.Open(string(o.config.Type), dsn)
	if err != nil {
		return fmt.Errorf("打开数据库连接失败: %w", err)
	}

	// 设置连接池参数
	db.SetMaxOpenConns(o.config.MaxOpenConns)
	db.SetMaxIdleConns(o.config.MaxIdleConns)
	db.SetConnMaxLifetime(o.config.MaxLifetime)

	// 测试连接
	if err := db.Ping(); err != nil {
		db.Close()
		return fmt.Errorf("数据库连接测试失败: %w", err)
	}

	o.db = db
	return nil
}

// Close 关闭数据库连接
func (o *ORM) Close() error {
	o.mu.Lock()
	defer o.mu.Unlock()

	if o.db != nil {
		return o.db.Close()
	}
	return nil
}

// Ping 测试数据库连接
func (o *ORM) Ping() error {
	o.mu.RLock()
	defer o.mu.RUnlock()

	if o.db == nil {
		return fmt.Errorf("数据库未连接")
	}
	return o.db.Ping()
}

// Query 执行查询
func (o *ORM) Query(query string, args ...interface{}) (*sql.Rows, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	if o.db == nil {
		return nil, fmt.Errorf("数据库未连接")
	}
	return o.db.Query(query, args...)
}

// QueryRow 执行单行查询
func (o *ORM) QueryRow(query string, args ...interface{}) *sql.Row {
	o.mu.RLock()
	defer o.mu.RUnlock()

	if o.db == nil {
		panic("数据库未连接")
	}
	return o.db.QueryRow(query, args...)
}

// Exec 执行SQL语句
func (o *ORM) Exec(query string, args ...interface{}) (sql.Result, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	if o.db == nil {
		return nil, fmt.Errorf("数据库未连接")
	}
	return o.db.Exec(query, args...)
}

// Begin 开始事务
func (o *ORM) Begin() (Tx, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	if o.db == nil {
		return nil, fmt.Errorf("数据库未连接")
	}

	tx, err := o.db.Begin()
	if err != nil {
		return nil, err
	}

	return &transaction{tx: tx}, nil
}

// BeginTx 开始带选项的事务
func (o *ORM) BeginTx(ctx context.Context, opts *sql.TxOptions) (Tx, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	if o.db == nil {
		return nil, fmt.Errorf("数据库未连接")
	}

	tx, err := o.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}

	return &transaction{tx: tx}, nil
}

// Raw 获取原始数据库连接
func (o *ORM) Raw() *sql.DB {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.db
}

// Table 创建查询构建器
func (o *ORM) Table(tableName string) QueryBuilder {
	return NewQueryBuilder(o, tableName)
}

// Model 基于模型创建查询构建器
func (o *ORM) Model(model interface{}) QueryBuilder {
	tableName := o.getTableName(model)
	return NewQueryBuilder(o, tableName)
}

// buildDSN 构建数据源名称
func (o *ORM) buildDSN() (string, error) {
	switch o.config.Type {
	case MySQL:
		return o.buildMySQLDSN(), nil
	case PostgreSQL:
		return o.buildPostgreSQLDSN(), nil
	case SQLite:
		return o.buildSQLiteDSN(), nil
	case SQLServer:
		return o.buildSQLServerDSN(), nil
	default:
		return "", fmt.Errorf("不支持的数据库类型: %s", o.config.Type)
	}
}

// buildMySQLDSN 构建MySQL DSN
func (o *ORM) buildMySQLDSN() string {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		o.config.Username,
		o.config.Password,
		o.config.Host,
		o.config.Port,
		o.config.Database,
	)

	params := []string{}
	if o.config.Charset != "" {
		params = append(params, "charset="+o.config.Charset)
	}
	if o.config.Timezone != "" {
		params = append(params, "loc="+o.config.Timezone)
	}

	if len(params) > 0 {
		dsn += "?" + strings.Join(params, "&")
	}

	return dsn
}

// buildPostgreSQLDSN 构建PostgreSQL DSN
func (o *ORM) buildPostgreSQLDSN() string {
	params := []string{
		fmt.Sprintf("host=%s", o.config.Host),
		fmt.Sprintf("port=%d", o.config.Port),
		fmt.Sprintf("user=%s", o.config.Username),
		fmt.Sprintf("password=%s", o.config.Password),
		fmt.Sprintf("dbname=%s", o.config.Database),
	}

	if o.config.SSLMode != "" {
		params = append(params, "sslmode="+o.config.SSLMode)
	}
	if o.config.Timezone != "" {
		params = append(params, "timezone="+o.config.Timezone)
	}

	return strings.Join(params, " ")
}

// buildSQLiteDSN 构建SQLite DSN
func (o *ORM) buildSQLiteDSN() string {
	return o.config.Database
}

// buildSQLServerDSN 构建SQL Server DSN
func (o *ORM) buildSQLServerDSN() string {
	return fmt.Sprintf("server=%s;port=%d;user id=%s;password=%s;database=%s",
		o.config.Host,
		o.config.Port,
		o.config.Username,
		o.config.Password,
		o.config.Database,
	)
}

// getTableName 获取表名
func (o *ORM) getTableName(model interface{}) string {
	if m, ok := model.(ModelInterface); ok {
		return m.TableName()
	}

	// 使用反射获取结构体名称
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// 将驼峰命名转换为下划线命名
	name := t.Name()
	return camelToSnake(name)
}

// camelToSnake 驼峰命名转下划线命名
func camelToSnake(s string) string {
	var result []rune
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '_')
		}
		result = append(result, r)
	}
	return strings.ToLower(string(result))
}

// 全局便捷方法

// Connect 连接数据库
func Connect() error {
	return GetGlobalORM().Connect()
}

// Close 关闭数据库连接
func Close() error {
	return GetGlobalORM().Close()
}

// Table 创建查询构建器
func Table(tableName string) QueryBuilder {
	return GetGlobalORM().Table(tableName)
}

// Model 基于模型创建查询构建器
func Model(model interface{}) QueryBuilder {
	return GetGlobalORM().Model(model)
}

// Begin 开始事务
func Begin() (Tx, error) {
	return GetGlobalORM().Begin()
}

// Raw 获取原始数据库连接
func Raw() *sql.DB {
	return GetGlobalORM().Raw()
}
