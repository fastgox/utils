# JWT - JSON Web Token 工具

一个简洁高效的JWT工具包，支持令牌生成、解析、验证和刷新。

## 🚀 特性

- **🎯 简洁API**: 类似其他helwd工具的简洁设计
- **🔐 安全可靠**: 使用HMAC-SHA256签名算法
- **⚙️ 灵活配置**: 支持全局配置和单次配置
- **📝 丰富字段**: 支持标准字段和自定义字段
- **🔄 令牌刷新**: 内置令牌刷新功能
- **⏰ 时间验证**: 自动处理过期时间和生效时间

## 📦 安装

```bash
go get github.com/fastgox/utils/jwt
```

## 🎯 快速开始

### 基础使用

```go
package main

import (
    "fmt"
    "github.com/fastgox/utils/jwt"
)

func main() {
    // 初始化JWT配置
    jwt.Init("my-secret-key", "my-app", 24*time.Hour)
    
    // 创建Claims
    claims := &jwt.Claims{
        UserID:   12345,
        Username: "helwd",
        Role:     "admin",
        Email:    "helwd@example.com",
    }
    
    // 生成令牌
    token, err := jwt.Generate(claims)
    if err != nil {
        panic(err)
    }
    fmt.Println("生成的令牌:", token)
    
    // 解析令牌
    parsedClaims, err := jwt.Parse(token)
    if err != nil {
        panic(err)
    }
    fmt.Printf("用户ID: %v\n", parsedClaims.UserID)
    fmt.Printf("用户名: %s\n", parsedClaims.Username)
    
    // 验证令牌
    if err := jwt.Verify(token); err != nil {
        fmt.Println("令牌无效:", err)
    } else {
        fmt.Println("令牌有效")
    }
}
```

### 自定义字段

```go
claims := &jwt.Claims{
    UserID:   12345,
    Username: "helwd",
    Custom: map[string]interface{}{
        "department": "技术部",
        "level":      5,
        "permissions": []string{"read", "write", "admin"},
    },
}

token, err := jwt.Generate(claims)
```

### 使用自定义配置

```go
// 为特定操作使用不同配置
config := &jwt.Config{
    Secret:     "special-secret",
    Issuer:     "special-app",
    Expiration: 1 * time.Hour, // 1小时过期
}

token, err := jwt.GenerateWithConfig(claims, config)
parsedClaims, err := jwt.ParseWithConfig(token, config)
```

### 令牌刷新

```go
// 刷新令牌（重新生成过期时间）
newToken, err := jwt.Refresh(oldToken)
if err != nil {
    fmt.Println("刷新失败:", err)
} else {
    fmt.Println("新令牌:", newToken)
}
```

## 📚 API 文档

### 配置函数

```go
// 初始化全局配置
jwt.Init(secret, issuer string, expiration time.Duration)

// 使用默认配置初始化
jwt.InitDefault()

// 设置全局密钥
jwt.SetSecret(secret string)

// 设置全局签发者
jwt.SetIssuer(issuer string)

// 设置全局过期时间
jwt.SetExpiration(expiration time.Duration)
```

### 核心函数

```go
// 生成令牌
jwt.Generate(claims *Claims) (string, error)
jwt.GenerateWithConfig(claims *Claims, config *Config) (string, error)

// 解析令牌
jwt.Parse(token string) (*Claims, error)
jwt.ParseWithConfig(token string, config *Config) (*Claims, error)

// 验证令牌
jwt.Verify(token string) error
jwt.VerifyWithConfig(token string, config *Config) error

// 刷新令牌
jwt.Refresh(token string) (string, error)
jwt.RefreshWithConfig(token string, config *Config) (string, error)
```

### 工具函数

```go
// 检查是否过期
jwt.IsExpired(token string) bool

// 获取Claims（不验证签名）
jwt.GetClaims(token string) (*Claims, error)
```

## 🏗️ 数据结构

### Claims 结构

```go
type Claims struct {
    UserID    interface{}            `json:"user_id,omitempty"`
    Username  string                 `json:"username,omitempty"`
    Role      string                 `json:"role,omitempty"`
    Email     string                 `json:"email,omitempty"`
    Issuer    string                 `json:"iss,omitempty"`
    Subject   string                 `json:"sub,omitempty"`
    Audience  string                 `json:"aud,omitempty"`
    IssuedAt  int64                  `json:"iat,omitempty"`
    ExpireAt  int64                  `json:"exp,omitempty"`
    NotBefore int64                  `json:"nbf,omitempty"`
    Custom    map[string]interface{} `json:"-"`
}
```

### Config 结构

```go
type Config struct {
    Secret     string        // 签名密钥
    Issuer     string        // 签发者
    Expiration time.Duration // 过期时间，0表示永不过期
}
```

## 🔒 安全建议

1. **密钥管理**: 使用足够复杂的密钥，建议从环境变量读取
2. **过期时间**: 根据应用场景设置合适的过期时间
3. **HTTPS**: 在生产环境中始终使用HTTPS传输令牌
4. **存储安全**: 客户端安全存储令牌，避免XSS攻击

## 🎮 完整示例

```go
package main

import (
    "fmt"
    "time"
    "github.com/fastgox/utils/jwt"
)

func main() {
    // 初始化
    jwt.Init("my-super-secret-key", "helwd-app", 2*time.Hour)
    
    // 用户登录，生成令牌
    loginClaims := &jwt.Claims{
        UserID:   12345,
        Username: "helwd",
        Role:     "admin",
        Email:    "helwd@example.com",
        Custom: map[string]interface{}{
            "login_ip": "192.168.1.100",
            "device":   "mobile",
        },
    }
    
    token, err := jwt.Generate(loginClaims)
    if err != nil {
        panic(err)
    }
    
    fmt.Println("🎫 生成令牌成功")
    fmt.Println("令牌:", token)
    
    // 验证令牌
    if err := jwt.Verify(token); err != nil {
        fmt.Println("❌ 令牌验证失败:", err)
        return
    }
    
    fmt.Println("✅ 令牌验证成功")
    
    // 解析用户信息
    claims, err := jwt.Parse(token)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("👤 用户信息:\n")
    fmt.Printf("   ID: %v\n", claims.UserID)
    fmt.Printf("   用户名: %s\n", claims.Username)
    fmt.Printf("   角色: %s\n", claims.Role)
    fmt.Printf("   邮箱: %s\n", claims.Email)
    fmt.Printf("   登录IP: %v\n", claims.Custom["login_ip"])
    fmt.Printf("   设备: %v\n", claims.Custom["device"])
    
    // 检查过期时间
    if jwt.IsExpired(token) {
        fmt.Println("⏰ 令牌已过期")
    } else {
        expireTime := time.Unix(claims.ExpireAt, 0)
        fmt.Printf("⏰ 令牌将于 %s 过期\n", expireTime.Format("2006-01-02 15:04:05"))
    }
    
    // 刷新令牌
    newToken, err := jwt.Refresh(token)
    if err != nil {
        fmt.Println("🔄 令牌刷新失败:", err)
    } else {
        fmt.Println("🔄 令牌刷新成功")
        fmt.Println("新令牌:", newToken)
    }
}
```

## 📄 许可证

MIT License
