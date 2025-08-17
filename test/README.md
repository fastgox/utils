# 测试文件

这个目录包含所有工具包的测试文件。

## 运行测试

```bash
# 运行所有测试
go test -v ./test

# 运行特定测试
go test -v ./test -run TestLogger
go test -v ./test -run TestHTTPClient
```

## 测试文件

- `logger_test.go` - 日志工具测试
- `http_test.go` - HTTP客户端工具测试
