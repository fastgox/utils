# 测试文件

这里包含了所有工具包的测试文件，按功能模块分类组织。

## 📁 目录结构

```
test/
├── README.md           # 测试说明文档
├── config/            # 配置工具测试
│   └── config_test.go
├── crypto/            # 加密工具测试
│   └── crypto_test.go
├── http/              # HTTP客户端测试
│   └── http_test.go
├── jwt/               # JWT工具测试
│   └── jwt_test.go
├── logger/            # 日志工具测试
│   └── logger_test.go
├── orm/               # ORM工具测试
│   ├── orm_test.go           # 基础功能测试
│   ├── orm_interface_test.go # 接口测试
│   └── orm_example_test.go   # 完整示例测试
├── string/            # 字符串工具测试
│   └── string_test.go
└── test_logs/         # 测试日志输出目录
```

## 🧪 运行测试

### 运行所有测试
```bash
go test ./test/...
```

### 运行特定模块测试
```bash
# ORM模块测试
go test ./test/orm -v

# 加密模块测试
go test ./test/crypto -v

# HTTP客户端测试
go test ./test/http -v

# JWT模块测试
go test ./test/jwt -v

# 日志模块测试
go test ./test/logger -v

# 配置模块测试
go test ./test/config -v
```

### 运行特定测试函数
```bash
# 运行ORM接口测试
go test ./test/orm -run TestORMInterfaces -v

# 运行加密基础测试
go test ./test/crypto -run TestCryptoBasic -v
```

## 📝 测试规范

1. **包名约定**: 每个测试包使用 `模块名_test` 作为包名
2. **文件命名**: 测试文件以 `_test.go` 结尾
3. **函数命名**: 测试函数以 `Test` 开头，使用驼峰命名
4. **分类组织**: 按功能模块分目录存放测试文件
5. **文档完整**: 每个测试函数都有清晰的注释说明

## 🎯 测试覆盖

- ✅ **配置工具**: JSON/YAML配置文件读取和解析
- ✅ **加密工具**: AES、RSA、哈希等加密功能
- ✅ **HTTP客户端**: GET/POST请求、参数处理、错误处理
- ✅ **JWT工具**: Token生成、验证、Claims处理
- ✅ **日志工具**: 多级别日志、文件输出、格式化
- ✅ **ORM工具**: 接口测试、查询构建器、数据库方言

## 🔧 测试环境

- **Go版本**: 1.19+
- **测试数据库**: SQLite（内存数据库，用于ORM测试）
- **依赖管理**: Go Modules
- **CI/CD**: 支持自动化测试流水线
