# Utils - Go工具包集合

一个精心设计的Go工具包集合，专注于简洁实用，让开发更高效。

## 🚀 包含工具

### 📝 Log - 日志工具
- **简洁配置**: 基于YAML的配置文件
- **自动分类**: 按日期和级别自动组织日志文件
- **多种级别**: Debug、Info、Warn、Error
- **格式化支持**: 支持格式化日志输出

### 🌐 HttpUtil - HTTP工具
- **极简API**: 类似Java HttpUtil的调用方式
- **表单支持**: 自动处理表单数据编码
- **JSON处理**: 内置JSON编码/解码
- **全局配置**: 支持全局认证和头部设置

### 🔐 JWT - 令牌工具
- **安全可靠**: 使用HMAC-SHA256签名算法
- **灵活配置**: 支持全局配置和单次配置
- **丰富字段**: 支持标准字段和自定义字段
- **令牌刷新**: 内置令牌刷新功能

### ⚙️ Config - 配置管理工具
- **多格式支持**: 支持YAML、JSON、TOML等格式
- **环境变量**: 自动映射和覆盖配置
- **结构体绑定**: 类型安全的配置绑定
- **配置验证**: 内置配置验证功能
- **热重载**: 支持配置文件热重载

### 🔐 Crypto - 加密工具
- **AES加密**: 支持AES-128/192/256加密解密
- **RSA加密**: 支持RSA公钥/私钥加密解密和数字签名
- **哈希算法**: 支持MD5、SHA1、SHA256、SHA512
- **密码哈希**: 支持bcrypt密码加盐哈希和强度检查
- **工具函数**: 随机数生成、Base64/Hex编码等

## 📦 安装

```bash
go get github.com/fastgox/utils
```

## 🎯 快速开始

### Log 使用示例

```go
package main

import (
    "github.com/fastgox/utils/Log"
)

func main() {
    // 初始化日志
    err := Log.InitDefault()
    if err != nil {
        panic(err)
    }

    // 使用日志
    Log.Info("应用启动成功")
    Log.Debugf("处理用户: %s", "helwd")
    Log.Warn("磁盘空间不足")
    Log.Error("数据库连接失败")
}
```

### HttpUtil 使用示例

```go
package main

import (
    "github.com/fastgox/utils/HttpUtil"
)

func main() {
    // GET请求
    result, err := HttpUtil.Get("https://api.example.com/users")
    if err != nil {
        panic(err)
    }

    // POST表单数据（推荐用法）
    paramMap := map[string]interface{}{
        "city": "北京",
        "name": "helwd",
    }
    result, err = HttpUtil.Post("https://api.example.com/search", paramMap)

    // POST JSON数据
    data := map[string]interface{}{
        "user": "helwd",
        "message": "Hello World",
    }
    result, err = HttpUtil.PostJSON("https://api.example.com/messages", data)
}
```

### Config 使用示例

```go
package main

import (
    "github.com/fastgox/utils/config"
)

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
}

func main() {
    // 初始化配置
    err := config.Init("config.yaml")
    if err != nil {
        panic(err)
    }

    // 获取配置值
    appName := config.GetString("app.name")
    serverPort := config.GetInt("server.port")

    // 结构体绑定
    var cfg AppConfig
    err = config.Unmarshal(&cfg)
    if err != nil {
        panic(err)
    }

    // 配置验证
    err = config.ValidateStruct(&cfg)
    if err != nil {
        panic(err)
    }

    // 环境变量覆盖
    config.SetEnvPrefix("MYAPP")
    config.BindEnv("server.port") // 对应 MYAPP_SERVER_PORT
}
```

### Crypto 使用示例

```go
package main

import (
    "fmt"
    "github.com/fastgox/utils/crypto"
)

func main() {
    // AES加密解密
    plaintext := "Hello, World!"
    password := "my-secure-password"

    encrypted, err := crypto.QuickEncrypt(plaintext, password)
    if err != nil {
        panic(err)
    }

    decrypted, err := crypto.QuickDecrypt(encrypted, password)
    if err != nil {
        panic(err)
    }

    // RSA加密解密
    privateKey, publicKey, err := crypto.GenerateKeyPair()
    if err != nil {
        panic(err)
    }

    rsaEncrypted, err := crypto.RSAEncrypt("Hello, RSA!", publicKey)
    if err != nil {
        panic(err)
    }

    rsaDecrypted, err := crypto.RSADecrypt(rsaEncrypted, privateKey)
    if err != nil {
        panic(err)
    }

    // 哈希算法
    md5Hash := crypto.MD5("Hello, Hash!")
    sha256Hash := crypto.SHA256("Hello, Hash!")
    hmacHash := crypto.HMACSHA256("data", "secret-key")

    // 密码哈希
    hashedPassword, err := crypto.HashPassword("my-password")
    if err != nil {
        panic(err)
    }

    isValid := crypto.CheckPassword("my-password", hashedPassword)
    fmt.Printf("密码验证: %v\n", isValid)

    // 生成强密码
    strongPassword, err := crypto.GenerateStrongPassword(16)
    if err != nil {
        panic(err)
    }

    strength := crypto.CheckPasswordStrength(strongPassword)
    fmt.Printf("生成的强密码: %s (强度: %s)\n", strongPassword, strength.String())
}
```

## 📚 详细文档

- [Log 工具文档](./log/README.md)
- [HttpUtil 工具文档](./http/README.md)
- [JWT 工具文档](./jwt/README.md)
- [Config 工具文档](./config/README.md)
- [Crypto 工具文档](./crypto/README.md)

## 🎮 运行示例

```bash
# 克隆项目
git clone https://github.com/fastgox/utils.git
cd utils

# 运行示例
go run main.go
```

这将演示所有工具的功能，并生成示例配置文件。

## 🌟 特色

- **🎯 专注实用**: 只包含最常用的功能，避免过度设计
- **📝 简洁API**: 链式调用，代码简洁优雅
- **⚡ 开箱即用**: 无需复杂配置，直接使用
- **🔧 灵活配置**: 支持配置文件和代码配置
- **💡 最佳实践**: 遵循Go语言最佳实践

## 🤝 贡献

欢迎提交Issue和Pull Request！

## 📄 许可证

MIT License
