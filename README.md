# Utils - 快速调用工具库

这是一个基于Go语言开发的快速调用工具项目，旨在提供高效、简洁的功能模块，方便开发者快速集成和使用。

## 安装

```bash
go get github.com/fastgox/utils@latest
```

## 模块

### Logger 日志模块

支持多级别日志记录，自动按日期和类型组织日志文件。

```go
import "github.com/fastgox/utils/logger"

// 初始化
logger.InitDefault()

// 使用
logger.Info("应用启动")
logger.Error("错误信息")
logger.Debugf("用户 %s 登录", "helwd")
```

## 示例

运行示例程序：

```bash
go run main.go
```
