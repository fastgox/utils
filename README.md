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

## 📚 详细文档

- [Log 工具文档](./Log/README.md)
- [HttpUtil 工具文档](./HttpUtil/README.md)

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
