# Logger 日志工具

## 使用方法

```go
package main

import (
    "github.com/fastgox/utils/log"
)

func main() {
    // 初始化日志系统
    err := logger.InitWithPath("logs")
    if err != nil {
        panic(err)
    }
    defer logger.CloseAll()

    // 获取不同事件类型的logger
    userLogger, _ := logger.GetLogger("user")
    orderLogger, _ := logger.GetLogger("order")

    // 记录日志
    userLogger.Info("用户登录: userID=%d", 12345)
    orderLogger.Error("订单创建失败: %s", "库存不足")
}
```
