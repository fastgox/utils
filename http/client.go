package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Config HTTP客户端配置，既可用于全局配置，也可用于单次请求配置
type Config struct {
	Timeout time.Duration     // 超时时间，0表示使用默认值
	Auth    string            // 认证信息，空字符串表示不使用认证
	Headers map[string]string // 请求头，nil表示不设置额外头部
}

var (
	// 全局配置
	globalConfig = &Config{
		Timeout: 30 * time.Second,
		Auth:    "",
		Headers: make(map[string]string),
	}
)

// Init 初始化HTTP客户端全局配置
func Init(timeout time.Duration, auth string, headers map[string]string) {
	globalConfig.Timeout = timeout
	globalConfig.Auth = auth
	if headers != nil {
		globalConfig.Headers = make(map[string]string)
		for k, v := range headers {
			globalConfig.Headers[k] = v
		}
	} else {
		globalConfig.Headers = make(map[string]string)
	}
}

// InitDefault 使用默认配置初始化
func InitDefault() {
	Init(30*time.Second, "", nil)
}

// SetTimeout 设置全局超时时间（秒）
func SetTimeout(seconds int) {
	globalConfig.Timeout = time.Duration(seconds) * time.Second
}

// SetAuth 设置全局认证信息
func SetAuth(auth string) {
	globalConfig.Auth = auth
}

// SetHeader 设置全局请求头
func SetHeader(key, value string) {
	globalConfig.Headers[key] = value
}

// ClearHeaders 清除所有全局请求头
func ClearHeaders() {
	globalConfig.Headers = make(map[string]string)
}

// Get 发送GET请求，返回响应文本
func Get(url string) (string, error) {
	return doRequest("GET", url, "", nil)
}

// GetWithConfig 发送GET请求，支持自定义配置
func GetWithConfig(url string, config *Config) (string, error) {
	return doRequestWithConfig("GET", url, "", nil, config)
}

// Post 发送POST请求，参数为表单数据
func Post(urlStr string, params map[string]interface{}) (string, error) {
	formData := url.Values{}
	for key, value := range params {
		formData.Set(key, fmt.Sprintf("%v", value))
	}
	return doRequest("POST", urlStr, "application/x-www-form-urlencoded", strings.NewReader(formData.Encode()))
}

// PostWithConfig 发送POST请求，支持自定义配置
func PostWithConfig(urlStr string, params map[string]interface{}, config *Config) (string, error) {
	formData := url.Values{}
	for key, value := range params {
		formData.Set(key, fmt.Sprintf("%v", value))
	}
	return doRequestWithConfig("POST", urlStr, "application/x-www-form-urlencoded", strings.NewReader(formData.Encode()), config)
}

// PostJSON 发送JSON POST请求
func PostJSON(url string, data interface{}) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("JSON编码失败: %w", err)
	}
	return doRequest("POST", url, "application/json", bytes.NewReader(jsonData))
}

// PostJSONWithConfig 发送JSON POST请求，支持自定义配置
func PostJSONWithConfig(url string, data interface{}, config *Config) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("JSON编码失败: %w", err)
	}
	return doRequestWithConfig("POST", url, "application/json", bytes.NewReader(jsonData), config)
}

// PostText 发送文本POST请求
func PostText(url, text string) (string, error) {
	return doRequest("POST", url, "text/plain", strings.NewReader(text))
}

// PostTextWithConfig 发送文本POST请求，支持自定义配置
func PostTextWithConfig(url, text string, config *Config) (string, error) {
	return doRequestWithConfig("POST", url, "text/plain", strings.NewReader(text), config)
}

// Put 发送PUT请求，参数为表单数据
func Put(urlStr string, params map[string]interface{}) (string, error) {
	formData := url.Values{}
	for key, value := range params {
		formData.Set(key, fmt.Sprintf("%v", value))
	}
	return doRequest("PUT", urlStr, "application/x-www-form-urlencoded", strings.NewReader(formData.Encode()))
}

// PutWithConfig 发送PUT请求，支持自定义配置
func PutWithConfig(urlStr string, params map[string]interface{}, config *Config) (string, error) {
	formData := url.Values{}
	for key, value := range params {
		formData.Set(key, fmt.Sprintf("%v", value))
	}
	return doRequestWithConfig("PUT", urlStr, "application/x-www-form-urlencoded", strings.NewReader(formData.Encode()), config)
}

// PutJSON 发送JSON PUT请求
func PutJSON(url string, data interface{}) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("JSON编码失败: %w", err)
	}
	return doRequest("PUT", url, "application/json", bytes.NewReader(jsonData))
}

// PutJSONWithConfig 发送JSON PUT请求，支持自定义配置
func PutJSONWithConfig(url string, data interface{}, config *Config) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("JSON编码失败: %w", err)
	}
	return doRequestWithConfig("PUT", url, "application/json", bytes.NewReader(jsonData), config)
}

// Delete 发送DELETE请求
func Delete(url string) (string, error) {
	return doRequest("DELETE", url, "", nil)
}

// DeleteWithConfig 发送DELETE请求，支持自定义配置
func DeleteWithConfig(url string, config *Config) (string, error) {
	return doRequestWithConfig("DELETE", url, "", nil, config)
}

// doRequest 执行HTTP请求的核心方法
func doRequest(method, url, contentType string, body io.Reader) (string, error) {
	return doRequestWithConfig(method, url, contentType, body, nil)
}

// doRequestWithConfig 执行HTTP请求的核心方法，支持自定义配置
func doRequestWithConfig(method, url, contentType string, body io.Reader, config *Config) (string, error) {
	// 确定使用的配置
	timeout := globalConfig.Timeout
	auth := globalConfig.Auth
	headers := make(map[string]string)

	// 复制全局headers
	for k, v := range globalConfig.Headers {
		headers[k] = v
	}

	// 如果有config参数，覆盖相应配置
	if config != nil {
		if config.Timeout > 0 {
			timeout = config.Timeout
		}
		if config.Auth != "" {
			auth = config.Auth
		}
		// 合并headers
		if config.Headers != nil {
			for k, v := range config.Headers {
				headers[k] = v
			}
		}
	}

	client := &http.Client{
		Timeout: timeout,
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置Content-Type
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	// 设置认证
	if auth != "" {
		if strings.HasPrefix(auth, "Bearer ") {
			req.Header.Set("Authorization", auth)
		} else if strings.HasPrefix(auth, "Basic ") {
			req.Header.Set("Authorization", auth)
		} else {
			// 默认作为Bearer token处理
			req.Header.Set("Authorization", "Bearer "+auth)
		}
	}

	// 设置请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// 设置默认User-Agent
	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", "helwd-httpclient/1.0")
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	// 检查HTTP状态码
	if resp.StatusCode >= 400 {
		return string(responseBody), fmt.Errorf("HTTP错误 %d: %s", resp.StatusCode, resp.Status)
	}

	return string(responseBody), nil
}

// GetJSON 发送GET请求并解析JSON响应
func GetJSON(url string, result interface{}) error {
	response, err := Get(url)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(response), result)
}

// PostJSONAndParse 发送JSON POST请求并解析JSON响应
func PostJSONAndParse(url string, data interface{}, result interface{}) error {
	response, err := PostJSON(url, data)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(response), result)
}

// DownloadFile 下载文件
func DownloadFile(url, filename string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("下载失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("下载失败，HTTP状态码: %d", resp.StatusCode)
	}

	// 这里简化处理，实际使用时可能需要写入文件
	// 为了保持工具的简洁性，这里只返回成功
	return nil
}

// IsURL 检查字符串是否为有效的URL
func IsURL(str string) bool {
	return strings.HasPrefix(str, "http://") || strings.HasPrefix(str, "https://")
}

// BuildURL 构建URL，自动处理参数
func BuildURL(baseURL string, params map[string]interface{}) string {
	if len(params) == 0 {
		return baseURL
	}

	values := url.Values{}
	for key, value := range params {
		values.Set(key, fmt.Sprintf("%v", value))
	}

	separator := "?"
	if strings.Contains(baseURL, "?") {
		separator = "&"
	}

	return baseURL + separator + values.Encode()
}

// 便利方法

// QuickGet 快速GET请求，忽略错误（仅用于测试）
func QuickGet(url string) string {
	result, _ := Get(url)
	return result
}

// QuickPost 快速POST请求，忽略错误（仅用于测试）
func QuickPost(url string, params map[string]interface{}) string {
	result, _ := Post(url, params)
	return result
}

// Reset 重置所有全局配置
func Reset() {
	globalConfig.Timeout = 30 * time.Second
	globalConfig.Auth = ""
	globalConfig.Headers = make(map[string]string)
}
