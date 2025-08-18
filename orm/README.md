# ORM - 对象关系映射工具

一个功能强大、简单易用的Go ORM工具，支持所有主流数据库。

## ✨ 特性

- 🚀 **多数据库支持**: MySQL、PostgreSQL、SQLite、SQL Server
- 🔗 **链式查询**: 流畅的查询构建器API
- 🏗️ **模型映射**: 结构体到数据库表的自动映射
- 🔄 **事务支持**: 完整的事务管理功能
- 📊 **迁移工具**: 数据库表结构迁移和版本管理
- 🎯 **类型安全**: 编译时类型检查
- ⚡ **高性能**: 优化的SQL生成和执行

## 📦 安装

```bash
go get github.com/fastgox/utils/orm
```

## 🎯 快速开始

### 1. 配置数据库连接

```go
package main

import (
    "log"
    "github.com/fastgox/utils/orm"
)

func main() {
    // 配置数据库连接
    config := &orm.Config{
        Type:     orm.MySQL,
        Host:     "localhost",
        Port:     3306,
        Username: "root",
        Password: "password",
        Database: "testdb",
        Charset:  "utf8mb4",
    }

    // 初始化ORM
    if err := orm.Init(config); err != nil {
        log.Fatal("初始化ORM失败:", err)
    }
    defer orm.Close()
}
```

### 2. 定义模型

```go
type User struct {
    ID        uint      `orm:"id,primary,auto_increment" json:"id"`
    Name      string    `orm:"name,size:100,not_null" json:"name"`
    Email     string    `orm:"email,size:255,unique" json:"email"`
    Age       int       `orm:"age" json:"age"`
    IsActive  bool      `orm:"is_active,default:true" json:"is_active"`
    CreatedAt time.Time `orm:"created_at" json:"created_at"`
    UpdatedAt time.Time `orm:"updated_at" json:"updated_at"`
}

// 自定义表名（可选）
func (User) TableName() string {
    return "users"
}
```

### 3. 自动迁移

```go
// 自动创建表
if err := orm.AutoMigrate(&User{}); err != nil {
    panic(err)
}
```

### 4. 基本操作

#### 创建记录

```go
user := &User{
    Name:     "张三",
    Email:    "zhangsan@example.com",
    Age:      25,
    IsActive: true,
}

// 插入单条记录
if err := orm.Model(&User{}).Insert(user); err != nil {
    panic(err)
}

// 批量插入
users := []*User{
    {Name: "李四", Email: "lisi@example.com", Age: 30},
    {Name: "王五", Email: "wangwu@example.com", Age: 28},
}
if err := orm.Model(&User{}).InsertBatch(users); err != nil {
    panic(err)
}
```

#### 查询记录

```go
// 查询所有记录
var users []User
if err := orm.Model(&User{}).Find(&users); err != nil {
    panic(err)
}

// 条件查询
var user User
if err := orm.Model(&User{}).Where("email = ?", "zhangsan@example.com").First(&user); err != nil {
    panic(err)
}

// 复杂查询
var activeUsers []User
err := orm.Model(&User{}).
    Where("is_active = ?", true).
    Where("age > ?", 18).
    OrderBy("created_at", "DESC").
    Limit(10).
    Find(&activeUsers)
```

#### 更新记录

```go
// 更新单个字段
err := orm.Model(&User{}).
    Where("id = ?", 1).
    UpdateColumns(map[string]interface{}{
        "name": "新名字",
        "age":  26,
    })

// 更新整个结构体
user.Name = "更新的名字"
err := orm.Model(&User{}).Where("id = ?", user.ID).Update(&user)
```

#### 删除记录

```go
// 删除指定记录
err := orm.Model(&User{}).Where("id = ?", 1).Delete()

// 批量删除
err := orm.Model(&User{}).Where("is_active = ?", false).Delete()
```

### 5. 高级查询

#### JOIN查询

```go
type Order struct {
    ID        uint      `orm:"id,primary,auto_increment" json:"id"`
    UserID    uint      `orm:"user_id" json:"user_id"`
    Amount    float64   `orm:"amount" json:"amount"`
    CreatedAt time.Time `orm:"created_at" json:"created_at"`
    UpdatedAt time.Time `orm:"updated_at" json:"updated_at"`
}

// LEFT JOIN查询
var results []struct {
    User
    OrderAmount float64 `orm:"amount"`
}

err := orm.Table("users").
    Select("users.*, orders.amount").
    LeftJoin("orders", "users.id = orders.user_id").
    Where("users.is_active = ?", true).
    Find(&results)
```

#### 聚合查询

```go
// 统计记录数
count, err := orm.Model(&User{}).Where("is_active = ?", true).Count()

// 检查记录是否存在
exists, err := orm.Model(&User{}).Where("email = ?", "test@example.com").Exists()
```

### 6. 事务处理

```go
// 使用事务
err := orm.WithTransaction(func(tx orm.Tx) error {
    // 在事务中执行操作
    user := &User{Name: "事务用户", Email: "tx@example.com"}
    if err := tx.Table("users").Insert(user); err != nil {
        return err // 自动回滚
    }
    
    order := &Order{UserID: user.ID, Amount: 100.0}
    if err := tx.Table("orders").Insert(order); err != nil {
        return err // 自动回滚
    }
    
    return nil // 自动提交
})
```

### 7. 数据库迁移

#### 创建迁移

```go
type CreateUsersTable struct {
    *orm.BaseMigration
}

func (m *CreateUsersTable) Up() error {
    return m.CreateTable("users", func(t orm.Table) {
        t.ID()
        t.String("name", 100).NotNull()
        t.String("email", 255).Unique()
        t.Integer("age").Default(0)
        t.Boolean("is_active").Default(true)
        t.Timestamp("created_at")
        t.Timestamp("updated_at")
        
        t.Index("idx_users_email", "email")
    })
}

func (m *CreateUsersTable) Down() error {
    return m.DropTable("users")
}

func (m *CreateUsersTable) Version() string {
    return "20240101_120000_create_users_table"
}
```

#### 运行迁移

```go
// 创建迁移实例
migration := &CreateUsersTable{
    BaseMigration: orm.NewBaseMigration("20240101_120000_create_users_table", orm.GetGlobalORM()),
}

// 运行迁移
if err := orm.Migrate(migration); err != nil {
    panic(err)
}

// 回滚迁移
if err := orm.RollbackMigration(1); err != nil {
    panic(err)
}

// 查看迁移状态
if err := orm.MigrationStatus(); err != nil {
    panic(err)
}
```

## 🏷️ 模型标签

支持以下ORM标签：

- `column`: 指定列名
- `type`: 指定数据类型
- `size`: 字符串长度
- `primary`: 主键
- `auto_increment`: 自增
- `not_null`: 非空
- `unique`: 唯一
- `default`: 默认值
- `comment`: 注释
- `index`: 索引名
- `-`: 忽略字段

```go
type User struct {
    ID       uint      `orm:"id,primary,auto_increment"`
    Name     string    `orm:"name,size:100,not_null,comment:用户名"`
    Email    string    `orm:"email,unique,index:idx_email"`
    Password string    `orm:"-"` // 忽略此字段
    Age      int       `orm:"age,default:0"`
    Status   bool      `orm:"status,default:true"`
}
```

## 🔧 配置选项

```go
config := &orm.Config{
    Type:         orm.MySQL,           // 数据库类型
    Host:         "localhost",         // 主机地址
    Port:         3306,               // 端口
    Username:     "root",             // 用户名
    Password:     "password",         // 密码
    Database:     "testdb",           // 数据库名
    Charset:      "utf8mb4",          // 字符集
    SSLMode:      "disable",          // SSL模式
    Timezone:     "Asia/Shanghai",    // 时区
    MaxOpenConns: 100,                // 最大打开连接数
    MaxIdleConns: 10,                 // 最大空闲连接数
    MaxLifetime:  time.Hour,          // 连接最大生存时间
}
```

## 🗄️ 支持的数据库

| 数据库 | 驱动 | 状态 |
|--------|------|------|
| MySQL | github.com/go-sql-driver/mysql | ✅ |
| PostgreSQL | github.com/lib/pq | ✅ |
| SQLite | github.com/mattn/go-sqlite3 | ✅ |
| SQL Server | github.com/denisenkom/go-mssqldb | ✅ |

## 📝 最佳实践

1. **标准字段**: 建议在模型中包含ID、CreatedAt、UpdatedAt等标准字段
2. **定义表名**: 实现TableName()方法自定义表名，否则使用结构体名的下划线形式
3. **使用事务**: 对于多表操作使用事务确保数据一致性
4. **索引优化**: 为经常查询的字段添加索引
5. **迁移管理**: 使用迁移工具管理数据库结构变更
6. **标签规范**: 合理使用ORM标签定义字段属性

## 🚀 完整示例

查看 [example.go](./example.go) 文件获取完整的使用示例，包括：

- 模型定义
- 数据库连接
- CRUD操作
- 事务处理
- 复杂查询
- 批量操作

运行示例：

```bash
cd orm
go run example.go
```

## 🚀 完整示例

查看 [test/orm/orm_example_test.go](../test/orm/orm_example_test.go) 文件获取完整的使用示例，包括：

- 模型定义
- 数据库连接
- CRUD操作
- 事务处理
- 复杂查询
- 批量操作

运行示例测试：

```bash
go test ./test/orm -run TestORMCompleteExample -v
```

## 🤝 贡献

欢迎提交Issue和Pull Request！

## 📄 许可证

MIT License
