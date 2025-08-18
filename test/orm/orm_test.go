package orm_test

import (
	"testing"
	"time"

	"github.com/fastgox/utils/orm"
	_ "github.com/mattn/go-sqlite3" // SQLite驱动
)

// TestUser 测试用户模型
type TestUser struct {
	ID        uint      `orm:"id,primary,auto_increment" json:"id"`
	Name      string    `orm:"name,size:100,not_null" json:"name"`
	Email     string    `orm:"email,size:255,unique" json:"email"`
	Age       int       `orm:"age" json:"age"`
	IsActive  bool      `orm:"is_active,default:true" json:"is_active"`
	CreatedAt time.Time `orm:"created_at" json:"created_at"`
	UpdatedAt time.Time `orm:"updated_at" json:"updated_at"`
}

// TableName 自定义表名
func (TestUser) TableName() string {
	return "test_users"
}

// TestORMBasic 测试ORM基本功能
func TestORMBasic(t *testing.T) {
	// 配置SQLite数据库（用于测试）
	config := &orm.Config{
		Type:     orm.SQLite,
		Database: ":memory:", // 内存数据库
	}

	// 初始化ORM
	if err := orm.Init(config); err != nil {
		t.Fatalf("初始化ORM失败: %v", err)
	}
	defer orm.Close()

	// 自动迁移
	if err := orm.AutoMigrate(&TestUser{}); err != nil {
		t.Fatalf("自动迁移失败: %v", err)
	}

	t.Log("ORM基本功能测试通过")
}

// TestORMCRUD 测试CRUD操作
func TestORMCRUD(t *testing.T) {
	// 配置SQLite数据库（用于测试）
	config := &orm.Config{
		Type:     orm.SQLite,
		Database: ":memory:",
	}

	if err := orm.Init(config); err != nil {
		t.Fatalf("初始化ORM失败: %v", err)
	}
	defer orm.Close()

	// 自动迁移
	if err := orm.AutoMigrate(&TestUser{}); err != nil {
		t.Fatalf("自动迁移失败: %v", err)
	}

	// 测试插入
	user := &TestUser{
		Name:     "测试用户",
		Email:    "test@example.com",
		Age:      25,
		IsActive: true,
	}

	if err := orm.Model(&TestUser{}).Insert(user); err != nil {
		t.Fatalf("插入用户失败: %v", err)
	}

	// 测试查询
	var foundUser TestUser
	if err := orm.Model(&TestUser{}).Where("email = ?", "test@example.com").First(&foundUser); err != nil {
		t.Fatalf("查询用户失败: %v", err)
	}

	if foundUser.Name != "测试用户" {
		t.Errorf("期望用户名为 '测试用户'，实际为 '%s'", foundUser.Name)
	}

	// 测试更新
	if err := orm.Model(&TestUser{}).Where("id = ?", foundUser.ID).UpdateColumns(map[string]interface{}{
		"name": "更新的用户",
		"age":  26,
	}); err != nil {
		t.Fatalf("更新用户失败: %v", err)
	}

	// 验证更新
	var updatedUser TestUser
	if err := orm.Model(&TestUser{}).Where("id = ?", foundUser.ID).First(&updatedUser); err != nil {
		t.Fatalf("查询更新后的用户失败: %v", err)
	}

	if updatedUser.Name != "更新的用户" {
		t.Errorf("期望更新后的用户名为 '更新的用户'，实际为 '%s'", updatedUser.Name)
	}

	// 测试统计
	count, err := orm.Model(&TestUser{}).Count()
	if err != nil {
		t.Fatalf("统计用户数量失败: %v", err)
	}

	if count != 1 {
		t.Errorf("期望用户数量为 1，实际为 %d", count)
	}

	// 测试删除
	if err := orm.Model(&TestUser{}).Where("id = ?", foundUser.ID).Delete(); err != nil {
		t.Fatalf("删除用户失败: %v", err)
	}

	// 验证删除
	count, err = orm.Model(&TestUser{}).Count()
	if err != nil {
		t.Fatalf("统计删除后的用户数量失败: %v", err)
	}

	if count != 0 {
		t.Errorf("期望删除后用户数量为 0，实际为 %d", count)
	}

	t.Log("CRUD操作测试通过")
}

// TestORMTransaction 测试事务
func TestORMTransaction(t *testing.T) {
	config := &orm.Config{
		Type:     orm.SQLite,
		Database: ":memory:",
	}

	if err := orm.Init(config); err != nil {
		t.Fatalf("初始化ORM失败: %v", err)
	}
	defer orm.Close()

	if err := orm.AutoMigrate(&TestUser{}); err != nil {
		t.Fatalf("自动迁移失败: %v", err)
	}

	// 测试事务提交
	err := orm.WithTransaction(func(tx orm.Tx) error {
		user1 := &TestUser{Name: "用户1", Email: "user1@example.com", Age: 20}
		if err := tx.Table("users").Insert(user1); err != nil {
			return err
		}

		user2 := &TestUser{Name: "用户2", Email: "user2@example.com", Age: 22}
		if err := tx.Table("users").Insert(user2); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		t.Fatalf("事务执行失败: %v", err)
	}

	// 验证事务提交
	count, err := orm.Model(&TestUser{}).Count()
	if err != nil {
		t.Fatalf("统计用户数量失败: %v", err)
	}

	if count != 2 {
		t.Errorf("期望用户数量为 2，实际为 %d", count)
	}

	t.Log("事务测试通过")
}

// TestORMQueryBuilder 测试查询构建器
func TestORMQueryBuilder(t *testing.T) {
	config := &orm.Config{
		Type:     orm.SQLite,
		Database: ":memory:",
	}

	if err := orm.Init(config); err != nil {
		t.Fatalf("初始化ORM失败: %v", err)
	}
	defer orm.Close()

	if err := orm.AutoMigrate(&TestUser{}); err != nil {
		t.Fatalf("自动迁移失败: %v", err)
	}

	// 插入测试数据
	users := []*TestUser{
		{Name: "张三", Email: "zhangsan@example.com", Age: 25, IsActive: true},
		{Name: "李四", Email: "lisi@example.com", Age: 30, IsActive: true},
		{Name: "王五", Email: "wangwu@example.com", Age: 28, IsActive: false},
	}

	for _, user := range users {
		if err := orm.Model(&TestUser{}).Insert(user); err != nil {
			t.Fatalf("插入用户失败: %v", err)
		}
	}

	// 测试条件查询
	var activeUsers []TestUser
	err := orm.Model(&TestUser{}).
		Where("is_active = ?", true).
		Where("age > ?", 20).
		OrderBy("age", "ASC").
		Find(&activeUsers)

	if err != nil {
		t.Fatalf("条件查询失败: %v", err)
	}

	if len(activeUsers) != 2 {
		t.Errorf("期望查询到 2 个活跃用户，实际为 %d", len(activeUsers))
	}

	// 测试IN查询
	var users1 []TestUser
	err = orm.Model(&TestUser{}).
		WhereIn("name", "张三", "李四").
		Find(&users1)

	if err != nil {
		t.Fatalf("IN查询失败: %v", err)
	}

	if len(users1) != 2 {
		t.Errorf("期望IN查询到 2 个用户，实际为 %d", len(users1))
	}

	// 测试限制查询
	var limitedUsers []TestUser
	err = orm.Model(&TestUser{}).
		Limit(2).
		OrderBy("id", "ASC").
		Find(&limitedUsers)

	if err != nil {
		t.Fatalf("限制查询失败: %v", err)
	}

	if len(limitedUsers) != 2 {
		t.Errorf("期望限制查询到 2 个用户，实际为 %d", len(limitedUsers))
	}

	t.Log("查询构建器测试通过")
}

// TestORMModelMapping 测试模型映射
func TestORMModelMapping(t *testing.T) {
	config := &orm.Config{
		Type:     orm.SQLite,
		Database: ":memory:",
	}

	if err := orm.Init(config); err != nil {
		t.Fatalf("初始化ORM失败: %v", err)
	}
	defer orm.Close()

	// 测试表是否存在
	exists, err := orm.HasTable(&TestUser{})
	if err != nil {
		t.Fatalf("检查表是否存在失败: %v", err)
	}

	if exists {
		t.Error("期望表不存在，但实际存在")
	}

	// 创建表
	if err := orm.CreateTable(&TestUser{}); err != nil {
		t.Fatalf("创建表失败: %v", err)
	}

	// 再次检查表是否存在
	exists, err = orm.HasTable(&TestUser{})
	if err != nil {
		t.Fatalf("检查表是否存在失败: %v", err)
	}

	if !exists {
		t.Error("期望表存在，但实际不存在")
	}

	t.Log("模型映射测试通过")
}

// BenchmarkORMInsert 基准测试插入性能
func BenchmarkORMInsert(b *testing.B) {
	config := &orm.Config{
		Type:     orm.SQLite,
		Database: ":memory:",
	}

	if err := orm.Init(config); err != nil {
		b.Fatalf("初始化ORM失败: %v", err)
	}
	defer orm.Close()

	if err := orm.AutoMigrate(&TestUser{}); err != nil {
		b.Fatalf("自动迁移失败: %v", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		user := &TestUser{
			Name:     "基准测试用户",
			Email:    "benchmark@example.com",
			Age:      25,
			IsActive: true,
		}

		if err := orm.Model(&TestUser{}).Insert(user); err != nil {
			b.Fatalf("插入用户失败: %v", err)
		}
	}
}
