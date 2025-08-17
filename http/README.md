# HTTP Client 工具

## 使用方法

```go
package main

import (
    "time"
    "github.com/fastgox/utils/http"
)

func main() {
    // 初始化全局配置
    headers := map[string]string{
        "User-Agent": "MyApp/1.0",
        "Accept":     "application/json",
    }
    client.Init(10*time.Second, "Bearer your-token", headers)

    // 发送请求
    response, err := client.Get("https://api.example.com/users")
    if err != nil {
        fmt.Printf("请求失败: %v\n", err)
    }

    // POST JSON请求
    jsonData := map[string]interface{}{
        "name": "张三",
        "age":  25,
    }
    response, err = client.PostJSON("https://api.example.com/users", jsonData)
}
```


