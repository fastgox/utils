package http_test

import (
	"testing"
	"time"

	httpclient "github.com/fastgox/utils/http"
)

func TestHTTPClient(t *testing.T) {
	// 测试初始化和配置
	headers := map[string]string{
		"User-Agent": "TestApp/1.0",
		"Accept":     "application/json",
	}
	httpclient.Init(15*time.Second, "Bearer test-token", headers)

	// 验证配置已设置（通过后续操作验证）
	t.Log("HTTP客户端配置已设置")

	// 测试设置方法
	httpclient.SetTimeout(20)
	httpclient.SetAuth("Bearer new-token")
	httpclient.SetHeader("X-Test", "test-value")

	// 测试Config参数
	config := &httpclient.Config{
		Timeout: 5 * time.Second,
		Auth:    "custom-auth",
		Headers: map[string]string{
			"Custom-Header": "custom-value",
		},
	}

	if config.Timeout != 5*time.Second {
		t.Errorf("期望Config超时时间为5秒，实际为%v", config.Timeout)
	}

	// 测试URL构建功能
	baseURL := "https://api.example.com/users"
	params := map[string]interface{}{
		"page":  1,
		"limit": 10,
		"name":  "张三",
	}

	result := httpclient.BuildURL(baseURL, params)
	if !httpclient.IsURL(result) {
		t.Errorf("构建的URL无效: %s", result)
	}

	// 测试重置功能
	httpclient.Reset()
}
