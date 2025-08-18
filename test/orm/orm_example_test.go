package orm_test

import (
	"testing"
	"time"

	"github.com/fastgox/utils/orm"
	_ "github.com/mattn/go-sqlite3" // SQLite驱动
)

// User 用户模型
type User struct {
	ID        uint      `orm:"id,primary,auto_increment" json:"id"`
	Name      string    `orm:"name,size:100,not_null" json:"name"`
	Email     string    `orm:"email,size:255,unique" json:"email"`
	Age       int       `orm:"age" json:"age"`
	IsActive  bool      `orm:"is_active,default:true" json:"is_active"`
	CreatedAt time.Time `orm:"created_at" json:"created_at"`
	UpdatedAt time.Time `orm:"updated_at" json:"updated_at"`
}

// TableName 自定义表名
func (User) TableName() string {
	return "users"
}

// Product 产品模型
type Product struct {
	ID          uint      `orm:"id,primary,auto_increment" json:"id"`
	Name        string    `orm:"name,size:200,not_null" json:"name"`
	Price       float64   `orm:"price,not_null" json:"price"`
	Description string    `orm:"description,type:text" json:"description"`
	CreatedAt   time.Time `orm:"created_at" json:"created_at"`
	UpdatedAt   time.Time `orm:"updated_at" json:"updated_at"`
}

// Order 订单模型
type Order struct {
	ID        uint      `orm:"id,primary,auto_increment" json:"id"`
	UserID    uint      `orm:"user_id" json:"user_id"`
	ProductID uint      `orm:"product_id" json:"product_id"`
	Quantity  int       `orm:"quantity" json:"quantity"`
	Amount    float64   `orm:"amount" json:"amount"`
	Status    string    `orm:"status,size:50" json:"status"`
	CreatedAt time.Time `orm:"created_at" json:"created_at"`
	UpdatedAt time.Time `orm:"updated_at" json:"updated_at"`
}

// setupTestDB 设置测试数据库
func setupTestDB(t *testing.T) {
	config := &orm.Config{
		Type:     orm.SQLite,
		Database: ":memory:", // 使用内存数据库进行测试
	}

	if err := orm.Init(config); err != nil {
		t.Fatalf("初始化ORM失败: %v", err)
	}

	// 自动迁移
	if err := orm.AutoMigrate(&User{}, &Product{}, &Order{}); err != nil {
		t.Fatalf("自动迁移失败: %v", err)
	}
}

// TestORMCompleteExample 完整的ORM使用示例
func TestORMCompleteExample(t *testing.T) {
	setupTestDB(t)
	defer orm.Close()

	t.Log("=== ORM完整示例测试 ===")

	// 1. 创建用户
	t.Log("1. 创建用户")
	user := &User{
		Name:     "张三",
		Email:    "zhangsan@example.com",
		Age:      25,
		IsActive: true,
	}

	if err := orm.Model(&User{}).Insert(user); err != nil {
		t.Fatalf("创建用户失败: %v", err)
	}
	t.Logf("创建用户成功: ID=%d, Name=%s", user.ID, user.Name)

	// 2. 批量创建产品
	t.Log("2. 批量创建产品")
	products := []*Product{
		{Name: "笔记本电脑", Price: 5999.99, Description: "高性能笔记本电脑"},
		{Name: "无线鼠标", Price: 99.99, Description: "人体工学无线鼠标"},
		{Name: "机械键盘", Price: 299.99, Description: "青轴机械键盘"},
	}

	if err := orm.Model(&Product{}).InsertBatch(products); err != nil {
		t.Fatalf("批量创建产品失败: %v", err)
	}
	t.Logf("批量创建了 %d 个产品", len(products))

	// 3. 查询和验证
	t.Log("3. 查询和验证")
	var foundUser User
	if err := orm.Model(&User{}).Where("email = ?", "zhangsan@example.com").First(&foundUser); err != nil {
		t.Fatalf("查询用户失败: %v", err)
	}
	t.Logf("查询到用户: ID=%d, Name=%s, Email=%s", foundUser.ID, foundUser.Name, foundUser.Email)

	// 4. 条件查询产品
	t.Log("4. 条件查询产品")
	var expensiveProducts []Product
	err := orm.Model(&Product{}).
		Where("price > ?", 200).
		OrderBy("price", "DESC").
		Find(&expensiveProducts)
	if err != nil {
		t.Fatalf("查询产品失败: %v", err)
	}

	t.Logf("价格大于200的产品有 %d 个", len(expensiveProducts))
	for _, product := range expensiveProducts {
		t.Logf("  - %s: ¥%.2f", product.Name, product.Price)
	}

	// 5. 更新操作
	t.Log("5. 更新操作")
	if err := orm.Model(&User{}).Where("id = ?", foundUser.ID).UpdateColumns(map[string]interface{}{
		"name": "张三（已更新）",
		"age":  26,
	}); err != nil {
		t.Fatalf("更新用户失败: %v", err)
	}
	t.Log("用户更新成功")

	// 6. 事务操作
	t.Log("6. 事务操作")
	err = orm.WithTransaction(func(tx orm.Tx) error {
		// 创建订单
		order := &Order{
			UserID:    foundUser.ID,
			ProductID: products[0].ID,
			Quantity:  2,
			Amount:    products[0].Price * 2,
			Status:    "pending",
		}
		if err := tx.Model(&Order{}).Insert(order); err != nil {
			return err
		}

		// 更新产品库存（这里只是示例）
		if err := tx.Model(&Product{}).Where("id = ?", products[0].ID).UpdateColumns(map[string]interface{}{
			"description": "高性能笔记本电脑（库存已更新）",
		}); err != nil {
			return err
		}

		t.Logf("事务中创建了订单: ID=%d, Amount=%.2f", order.ID, order.Amount)
		return nil
	})

	if err != nil {
		t.Fatalf("事务执行失败: %v", err)
	}
	t.Log("事务执行成功")

	// 7. 统计查询
	t.Log("7. 统计查询")
	userCount, err := orm.Model(&User{}).Count()
	if err != nil {
		t.Fatalf("统计用户数量失败: %v", err)
	}
	t.Logf("用户总数: %d", userCount)

	productCount, err := orm.Model(&Product{}).Count()
	if err != nil {
		t.Fatalf("统计产品数量失败: %v", err)
	}
	t.Logf("产品总数: %d", productCount)

	orderCount, err := orm.Model(&Order{}).Count()
	if err != nil {
		t.Fatalf("统计订单数量失败: %v", err)
	}
	t.Logf("订单总数: %d", orderCount)

	// 8. 复杂查询
	t.Log("8. 复杂查询")
	var activeUsers []User
	err = orm.Model(&User{}).
		Where("is_active = ?", true).
		Where("age >= ?", 20).
		OrderBy("created_at", "DESC").
		Limit(10).
		Find(&activeUsers)

	if err != nil {
		t.Fatalf("复杂查询失败: %v", err)
	}

	t.Logf("活跃用户（年龄>=20）有 %d 个", len(activeUsers))
	for _, u := range activeUsers {
		t.Logf("  - %s (年龄: %d)", u.Name, u.Age)
	}

	// 9. IN查询
	t.Log("9. IN查询")
	var selectedProducts []Product
	err = orm.Model(&Product{}).
		WhereIn("name", "笔记本电脑", "机械键盘").
		Find(&selectedProducts)

	if err != nil {
		t.Fatalf("IN查询失败: %v", err)
	}

	t.Logf("指定产品有 %d 个", len(selectedProducts))
	for _, p := range selectedProducts {
		t.Logf("  - %s: ¥%.2f", p.Name, p.Price)
	}

	// 10. 存在性检查
	t.Log("10. 存在性检查")
	exists, err := orm.Model(&User{}).Where("email = ?", "zhangsan@example.com").Exists()
	if err != nil {
		t.Fatalf("检查用户存在性失败: %v", err)
	}
	t.Logf("用户是否存在: %v", exists)

	// 11. 删除操作
	t.Log("11. 删除操作")
	// 删除订单
	if err := orm.Model(&Order{}).Where("status = ?", "pending").Delete(); err != nil {
		t.Fatalf("删除订单失败: %v", err)
	}
	t.Log("删除待处理订单成功")

	// 验证删除
	remainingOrders, err := orm.Model(&Order{}).Count()
	if err != nil {
		t.Fatalf("统计剩余订单失败: %v", err)
	}
	t.Logf("剩余订单数: %d", remainingOrders)

	t.Log("=== ORM完整示例测试完成 ===")
}

// TestORMDifferentDatabases 测试不同数据库类型
func TestORMDifferentDatabases(t *testing.T) {
	// 这里只测试SQLite，其他数据库需要相应的环境
	databases := []struct {
		name   string
		config *orm.Config
	}{
		{
			name: "SQLite",
			config: &orm.Config{
				Type:     orm.SQLite,
				Database: ":memory:",
			},
		},
		// 可以添加其他数据库的测试配置
		// {
		//     name: "MySQL",
		//     config: &orm.Config{
		//         Type:     orm.MySQL,
		//         Host:     "localhost",
		//         Port:     3306,
		//         Username: "test",
		//         Password: "test",
		//         Database: "test",
		//     },
		// },
	}

	for _, db := range databases {
		t.Run(db.name, func(t *testing.T) {
			if err := orm.Init(db.config); err != nil {
				t.Skipf("跳过 %s 测试，初始化失败: %v", db.name, err)
				return
			}
			defer orm.Close()

			// 自动迁移
			if err := orm.AutoMigrate(&User{}); err != nil {
				t.Fatalf("%s 自动迁移失败: %v", db.name, err)
			}

			// 基本CRUD测试
			user := &User{
				Name:     "测试用户",
				Email:    "test@example.com",
				Age:      30,
				IsActive: true,
			}

			// 创建
			if err := orm.Model(&User{}).Insert(user); err != nil {
				t.Fatalf("%s 创建用户失败: %v", db.name, err)
			}

			// 查询
			var foundUser User
			if err := orm.Model(&User{}).Where("email = ?", "test@example.com").First(&foundUser); err != nil {
				t.Fatalf("%s 查询用户失败: %v", db.name, err)
			}

			if foundUser.Name != "测试用户" {
				t.Errorf("%s 用户名不匹配，期望: 测试用户，实际: %s", db.name, foundUser.Name)
			}

			t.Logf("%s 数据库测试通过", db.name)
		})
	}
}
