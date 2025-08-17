# Utils - Go工具包集合

一个精心设计的Go工具包集合，专注于简洁实用，让开发更高效。

## 📦 安装

```bash
go get github.com/fastgox/utils
```

## 🎯 快速开始

查看各工具的详细使用方法：
- [Logger 日志工具](./log/README.md) - 结构化日志记录
- [Config 配置工具](./config/README.md) - 配置文件管理和环境变量
- [Crypto 加密工具](./crypto/README.md) - AES/RSA加密、哈希和密码处理
- [HTTP 客户端工具](./http/README.md) - HTTP请求封装和客户端
- [JWT 令牌工具](./jwt/README.md) - JSON Web Token生成和验证
- [测试示例](./test/README.md) - 所有工具的测试用例和使用示例

### 基本使用示例

```go
package main

import (
    "fmt"
    "github.com/fastgox/utils/log"
    "github.com/fastgox/utils/config"
    "github.com/fastgox/utils/crypto"
)

func main() {
    // 日志记录
    log.Info("应用启动")

    // 配置管理
    cfg := config.Load("config.yaml")

    // 加密解密
    encrypted, _ := crypto.AESEncrypt("hello world", "your-secret-key")
    fmt.Println("加密结果:", encrypted)
}
```

## ✅ 已实现工具

- **Logger** - 结构化日志记录，支持多级别和格式化输出
- **Config** - 配置文件管理，支持YAML、环境变量和热重载
- **Crypto** - 加密工具集，包含AES、RSA、哈希和密码处理
- **HTTP** - HTTP客户端封装，支持GET/POST/PUT/DELETE和超时控制
- **JWT** - JSON Web Token工具，支持生成、验证和刷新

## 📋 工具开发计划 (TODO)

### 📊 Database - 数据库工具 (计划中)
- [ ] MySQL 连接池管理
- [ ] Redis 客户端封装
- [ ] 事务处理
- [ ] 查询构建器
- [ ] 数据迁移工具

### 🔧 Validator - 验证工具 (计划中)
- [ ] 结构体字段验证
- [ ] 自定义验证规则
- [ ] 错误信息国际化
- [ ] 嵌套结构验证

### 📧 Email - 邮件工具 (计划中)
- [ ] SMTP 邮件发送
- [ ] HTML/文本邮件支持
- [ ] 附件处理
- [ ] 邮件模板

### 🕒 Time - 时间工具 (计划中)
- [ ] 时间格式化
- [ ] 时区转换
- [ ] 时间计算
- [ ] 定时任务

### 📁 File - 文件工具 (计划中)
- [ ] 文件上传下载
- [ ] 文件压缩解压
- [ ] 文件类型检测
- [ ] 目录操作

### 🔄 Cache - 缓存工具 (计划中)
- [ ] 内存缓存
- [ ] Redis 缓存
- [ ] 缓存策略
- [ ] 过期管理




## 🤝 贡献

欢迎提交Issue和Pull Request！

## 📄 许可证

MIT License
