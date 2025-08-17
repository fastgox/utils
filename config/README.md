# Config - 配置管理工具

一个功能强大的Go配置管理工具，支持多种配置格式和环境变量覆盖。

## 🚀 特性

- **🎯 多格式支持**: 支持YAML、JSON、TOML、Properties、INI格式
- **🌍 环境变量**: 自动映射环境变量，支持配置覆盖
- **📋 结构体绑定**: 类型安全的配置绑定到Go结构体
- **✅ 配置验证**: 内置配置验证功能
- **🔄 热重载**: 支持配置文件变化监听和热重载
- **🏗️ 多环境**: 支持开发、测试、生产环境配置
- **⚡ 高性能**: 配置缓存，避免重复解析

## 📦 安装

```bash
go get github.com/fastgox/utils/config
```

## 🎯 快速开始

### 基础使用

```go
package main

import (
    "fmt"
    "github.com/fastgox/utils/config"
)

func main() {
    // 初始化配置
    err := config.Init("config.yaml")
    if err != nil {
        panic(err)
    }
    
    // 获取配置值
    dbHost := config.GetString("database.host")
    dbPort := config.GetInt("database.port")
    debug := config.GetBool("app.debug")
    
    fmt.Printf("数据库: %s:%d\n", dbHost, dbPort)
    fmt.Printf("调试模式: %v\n", debug)
}
```

### 结构体绑定

```go
type AppConfig struct {
    App struct {
        Name    string `config:"name" validate:"required"`
        Version string `config:"version" validate:"required"`
        Debug   bool   `config:"debug"`
    } `config:"app"`
    
    Server struct {
        Host string `config:"host" validate:"required"`
        Port int    `config:"port" validate:"min=1,max=65535"`
    } `config:"server"`
    
    Database struct {
        Host     string `config:"host" validate:"required"`
        Port     int    `config:"port" validate:"min=1,max=65535"`
        Username string `config:"username" validate:"required"`
        Password string `config:"password" validate:"required"`
    } `config:"database"`
}

func main() {
    config.Init("config.yaml")
    
    var cfg AppConfig
    err := config.Unmarshal(&cfg)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("应用: %s v%s\n", cfg.App.Name, cfg.App.Version)
    fmt.Printf("服务器: %s:%d\n", cfg.Server.Host, cfg.Server.Port)
}
```

### 环境变量覆盖

```go
func main() {
    // 设置环境变量前缀
    config.SetEnvPrefix("MYAPP")
    
    // 绑定环境变量
    config.BindEnv("database.host")  // 对应 MYAPP_DATABASE_HOST
    config.BindEnv("database.port")  // 对应 MYAPP_DATABASE_PORT
    
    config.Init("config.yaml")
    
    // 环境变量会自动覆盖配置文件中的值
    dbHost := config.GetString("database.host")
    fmt.Printf("数据库主机: %s\n", dbHost)
}
```

### 配置热重载

```go
func main() {
    config.Init("config.yaml")
    
    // 监听配置变化
    err := config.Watch(func(oldConfig, newConfig interface{}) {
        fmt.Println("配置文件已更新，重新加载应用配置")
        // 这里可以重新初始化数据库连接等
    })
    if err != nil {
        panic(err)
    }
    
    // 应用主逻辑...
    select {} // 保持程序运行
}
```

## 📁 配置文件示例

### config.yaml
```yaml
app:
  name: "helwd-app"
  version: "1.0.0"
  debug: true

server:
  host: "localhost"
  port: 8080
  timeout: 30s

database:
  host: "localhost"
  port: 3306
  username: "root"
  password: "password"
  dbname: "myapp"
  max_connections: 100

redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0

log:
  level: "info"
  format: "json"
  output: "logs/app.log"
```

## 📚 API 文档

### 初始化函数

```go
// 使用配置文件初始化
config.Init(configPath string) error

// 使用选项初始化
config.InitWithOptions(opts *Options) error

// 使用默认配置初始化
config.InitDefault() error
```

### 配置获取

```go
// 获取原始值
config.Get(key string) interface{}

// 获取特定类型值
config.GetString(key string) string
config.GetInt(key string) int
config.GetBool(key string) bool
config.GetFloat64(key string) float64
config.GetStringSlice(key string) []string
config.GetDuration(key string) time.Duration

// 带默认值获取
config.GetStringDefault(key, defaultValue string) string
config.GetIntDefault(key string, defaultValue int) int
```

### 结构体绑定

```go
// 绑定整个配置到结构体
config.Unmarshal(v interface{}) error

// 绑定指定键的配置到结构体
config.UnmarshalKey(key string, v interface{}) error
```

### 环境变量

```go
// 设置环境变量前缀
config.SetEnvPrefix(prefix string)

// 绑定环境变量
config.BindEnv(key string) error

// 自动绑定环境变量
config.AutomaticEnv()
```

### 配置监听

```go
// 监听配置变化
config.Watch(callback func(oldConfig, newConfig interface{})) error

// 停止监听
config.StopWatch()
```

## 🔧 高级功能

### 多环境支持

```go
// 设置环境
config.SetEnvironment("dev") // dev, test, prod

// 加载环境特定配置
config.LoadEnvironmentConfig("prod") // 会加载 config.prod.yaml
```

### 配置验证

```go
// 验证当前配置
err := config.Validate()

// 验证结构体
err := config.ValidateStruct(&cfg)
```

### 默认值设置

```go
// 设置默认值
config.SetDefault("server.port", 8080)
config.SetDefault("app.debug", false)
```

## 🌟 最佳实践

1. **配置文件命名**: 使用 `config.yaml` 作为主配置文件
2. **环境变量**: 使用应用名作为环境变量前缀
3. **结构体验证**: 为重要配置添加验证标签
4. **敏感信息**: 密码等敏感信息通过环境变量传入
5. **配置分层**: 将配置按功能模块分组

## 🤝 贡献

欢迎提交Issue和Pull Request！

## 📄 许可证

MIT License
