package orm_test

import (
	"testing"
	"time"

	"github.com/fastgox/utils/orm"
)

// SimpleUser 简单用户模型（用于接口测试）
type SimpleUser struct {
	ID        uint      `orm:"id,primary,auto_increment" json:"id"`
	Name      string    `orm:"name,size:100,not_null" json:"name"`
	Email     string    `orm:"email,size:255,unique" json:"email"`
	Age       int       `orm:"age" json:"age"`
	IsActive  bool      `orm:"is_active,default:true" json:"is_active"`
	CreatedAt time.Time `orm:"created_at" json:"created_at"`
	UpdatedAt time.Time `orm:"updated_at" json:"updated_at"`
}

// TableName 自定义表名
func (SimpleUser) TableName() string {
	return "simple_users"
}

// TestORMInterfaces 测试ORM接口和基本功能（不需要实际数据库连接）
func TestORMInterfaces(t *testing.T) {
	t.Log("=== ORM接口测试 ===")

	// 1. 测试配置创建
	t.Log("1. 测试配置创建")
	config := &orm.Config{
		Type:     orm.MySQL,
		Host:     "localhost",
		Port:     3306,
		Username: "test",
		Password: "test",
		Database: "testdb",
		Charset:  "utf8mb4",
	}

	if config.Type != orm.MySQL {
		t.Errorf("期望数据库类型为 MySQL，实际为 %s", config.Type)
	}

	// 2. 测试默认配置
	t.Log("2. 测试默认配置")
	defaultConfig := orm.DefaultConfig()
	if defaultConfig.Type != orm.MySQL {
		t.Errorf("期望默认数据库类型为 MySQL，实际为 %s", defaultConfig.Type)
	}
	if defaultConfig.MaxOpenConns != 100 {
		t.Errorf("期望默认最大连接数为 100，实际为 %d", defaultConfig.MaxOpenConns)
	}

	// 3. 测试ORM实例创建
	t.Log("3. 测试ORM实例创建")
	ormInstance := orm.New(config)
	if ormInstance == nil {
		t.Fatal("ORM实例创建失败")
	}

	// 4. 测试模型管理器
	t.Log("4. 测试模型管理器")
	modelManager := orm.NewModelManager(ormInstance)
	if modelManager == nil {
		t.Fatal("模型管理器创建失败")
	}

	// 测试获取表信息
	tableInfo := modelManager.GetTableInfo(&SimpleUser{})
	if tableInfo == nil {
		t.Fatal("获取表信息失败")
	}

	if tableInfo.Name != "simple_users" {
		t.Errorf("期望表名为 'simple_users'，实际为 '%s'", tableInfo.Name)
	}

	if len(tableInfo.Columns) == 0 {
		t.Error("期望有列信息，但实际为空")
	}

	// 检查主键列
	primaryKey := tableInfo.GetPrimaryKey()
	if primaryKey == nil {
		t.Error("期望有主键列，但未找到")
	} else if primaryKey.Name != "id" {
		t.Errorf("期望主键列名为 'id'，实际为 '%s'", primaryKey.Name)
	}

	// 5. 测试数据库管理器
	t.Log("5. 测试数据库管理器")
	dbManager := orm.NewDatabaseManager(ormInstance)
	if dbManager == nil {
		t.Fatal("数据库管理器创建失败")
	}

	// 测试方言获取
	dialect := dbManager.GetDialect()
	if dialect == nil {
		t.Fatal("获取数据库方言失败")
	}

	// 测试SQL生成
	columns := []orm.ColumnDefinition{
		{
			Name:          "id",
			Type:          "BIGINT",
			Primary:       true,
			AutoIncrement: true,
			NotNull:       true,
		},
		{
			Name:    "name",
			Type:    "VARCHAR(100)",
			NotNull: true,
		},
	}

	createSQL := dialect.CreateTableSQL("test_table", columns)
	if createSQL == "" {
		t.Error("生成的CREATE TABLE SQL为空")
	}
	t.Logf("生成的CREATE TABLE SQL: %s", createSQL)

	dropSQL := dialect.DropTableSQL("test_table")
	if dropSQL == "" {
		t.Error("生成的DROP TABLE SQL为空")
	}
	t.Logf("生成的DROP TABLE SQL: %s", dropSQL)

	// 6. 测试事务管理器
	t.Log("6. 测试事务管理器")
	txManager := orm.NewTransactionManager(ormInstance)
	if txManager == nil {
		t.Fatal("事务管理器创建失败")
	}

	// 7. 测试迁移管理器
	t.Log("7. 测试迁移管理器")
	migrationManager := orm.NewMigrationManager(ormInstance)
	if migrationManager == nil {
		t.Fatal("迁移管理器创建失败")
	}

	// 8. 测试Schema管理器
	t.Log("8. 测试Schema管理器")
	schema := orm.NewSchema(ormInstance)
	if schema == nil {
		t.Fatal("Schema管理器创建失败")
	}

	t.Log("=== ORM接口测试完成 ===")
}

// TestORMQueryBuilderInterface 测试查询构建器接口（不需要实际数据库连接）
func TestORMQueryBuilderInterface(t *testing.T) {
	t.Log("=== 查询构建器测试 ===")

	config := orm.DefaultConfig()
	ormInstance := orm.New(config)

	// 创建查询构建器
	qb := orm.NewQueryBuilder(ormInstance, "users")
	if qb == nil {
		t.Fatal("查询构建器创建失败")
	}

	// 测试链式调用
	qb = qb.Select("id", "name", "email").
		Where("age > ?", 18).
		Where("is_active = ?", true).
		OrderBy("created_at", "DESC").
		Limit(10).
		Offset(5)

	// 测试SQL生成
	sql, args := qb.ToSQL()
	if sql == "" {
		t.Error("生成的SQL为空")
	}
	if len(args) == 0 {
		t.Error("期望有查询参数，但实际为空")
	}

	t.Logf("生成的SQL: %s", sql)
	t.Logf("查询参数: %v", args)

	// 验证SQL包含期望的部分
	expectedParts := []string{"SELECT", "FROM users", "WHERE", "ORDER BY", "LIMIT", "OFFSET"}
	for _, part := range expectedParts {
		if !contains(sql, part) {
			t.Errorf("生成的SQL不包含期望的部分: %s", part)
		}
	}

	t.Log("=== 查询构建器测试完成 ===")
}

// TestORMDifferentDialects 测试不同数据库方言
func TestORMDifferentDialects(t *testing.T) {
	t.Log("=== 数据库方言测试 ===")

	databases := []struct {
		name   string
		dbType orm.DatabaseType
	}{
		{"MySQL", orm.MySQL},
		{"PostgreSQL", orm.PostgreSQL},
		{"SQLite", orm.SQLite},
		{"SQL Server", orm.SQLServer},
	}

	for _, db := range databases {
		t.Run(db.name, func(t *testing.T) {
			config := &orm.Config{Type: db.dbType}
			ormInstance := orm.New(config)
			dbManager := orm.NewDatabaseManager(ormInstance)
			dialect := dbManager.GetDialect()

			// 测试基本SQL生成
			columns := []orm.ColumnDefinition{
				{Name: "id", Type: "BIGINT", Primary: true, AutoIncrement: true},
				{Name: "name", Type: "VARCHAR(100)", NotNull: true},
			}

			createSQL := dialect.CreateTableSQL("test_table", columns)
			if createSQL == "" {
				t.Errorf("%s: 生成的CREATE TABLE SQL为空", db.name)
			}

			dropSQL := dialect.DropTableSQL("test_table")
			if dropSQL == "" {
				t.Errorf("%s: 生成的DROP TABLE SQL为空", db.name)
			}

			indexSQL := dialect.CreateIndexSQL("test_table", "idx_name", []string{"name"}, false)
			if indexSQL == "" {
				t.Errorf("%s: 生成的CREATE INDEX SQL为空", db.name)
			}

			t.Logf("%s - CREATE TABLE: %s", db.name, createSQL)
			t.Logf("%s - DROP TABLE: %s", db.name, dropSQL)
			t.Logf("%s - CREATE INDEX: %s", db.name, indexSQL)
		})
	}

	t.Log("=== 数据库方言测试完成 ===")
}

// contains 检查字符串是否包含子字符串（忽略大小写）
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			len(s) > len(substr) &&
				(s[:len(substr)] == substr ||
					s[len(s)-len(substr):] == substr ||
					containsInMiddle(s, substr)))
}

// containsInMiddle 检查字符串中间是否包含子字符串
func containsInMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
