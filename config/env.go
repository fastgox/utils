package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// EnvManager 环境变量管理器
type EnvManager struct {
	config *Config
}

// NewEnvManager 创建环境变量管理器
func NewEnvManager(config *Config) *EnvManager {
	return &EnvManager{
		config: config,
	}
}

// BindEnv 绑定环境变量
func (e *EnvManager) BindEnv(key string) error {
	envKey := e.keyToEnvVar(key)
	if e.config.envBindings == nil {
		e.config.envBindings = make(map[string]string)
	}
	e.config.envBindings[key] = envKey
	return nil
}

// LoadEnvVars 加载环境变量
func (e *EnvManager) LoadEnvVars() {
	if e.config.data == nil {
		e.config.data = make(map[string]interface{})
	}

	// 加载绑定的环境变量
	for key, envKey := range e.config.envBindings {
		if value := os.Getenv(envKey); value != "" {
			e.setConfigValue(key, value)
		}
	}

	// 如果启用了自动环境变量，扫描所有环境变量
	if e.config.automaticEnv {
		e.loadAutomaticEnvVars()
	}
}

// loadAutomaticEnvVars 自动加载环境变量
func (e *EnvManager) loadAutomaticEnvVars() {
	prefix := e.config.envPrefix
	if prefix != "" && !strings.HasSuffix(prefix, "_") {
		prefix += "_"
	}

	// 遍历所有环境变量
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 {
			continue
		}

		envKey := parts[0]
		envValue := parts[1]

		// 检查是否匹配前缀
		if prefix != "" && !strings.HasPrefix(envKey, prefix) {
			continue
		}

		// 转换环境变量名为配置键
		configKey := e.envVarToKey(envKey)
		if configKey != "" {
			e.setConfigValue(configKey, envValue)
		}
	}
}

// keyToEnvVar 将配置键转换为环境变量名
func (e *EnvManager) keyToEnvVar(key string) string {
	// 将点号替换为下划线，转换为大写
	envKey := strings.ToUpper(strings.ReplaceAll(key, ".", "_"))
	
	// 添加前缀
	if e.config.envPrefix != "" {
		prefix := strings.ToUpper(e.config.envPrefix)
		if !strings.HasSuffix(prefix, "_") {
			prefix += "_"
		}
		envKey = prefix + envKey
	}
	
	return envKey
}

// envVarToKey 将环境变量名转换为配置键
func (e *EnvManager) envVarToKey(envVar string) string {
	key := envVar
	
	// 移除前缀
	if e.config.envPrefix != "" {
		prefix := strings.ToUpper(e.config.envPrefix)
		if !strings.HasSuffix(prefix, "_") {
			prefix += "_"
		}
		if strings.HasPrefix(key, prefix) {
			key = key[len(prefix):]
		} else {
			return "" // 不匹配前缀，忽略
		}
	}
	
	// 转换为小写，下划线替换为点号
	key = strings.ToLower(key)
	key = strings.ReplaceAll(key, "_", ".")
	
	return key
}

// setConfigValue 设置配置值，自动类型转换
func (e *EnvManager) setConfigValue(key, value string) {
	// 尝试类型转换
	convertedValue := e.convertValue(value)
	
	// 设置到配置中
	e.setNestedValue(e.config.data, key, convertedValue)
}

// convertValue 转换字符串值为合适的类型
func (e *EnvManager) convertValue(value string) interface{} {
	// 尝试转换为布尔值
	if boolVal, err := strconv.ParseBool(value); err == nil {
		return boolVal
	}
	
	// 尝试转换为整数
	if intVal, err := strconv.ParseInt(value, 10, 64); err == nil {
		// 如果值在int范围内，返回int，否则返回int64
		if intVal >= int64(int(^uint(0)>>1)*-1) && intVal <= int64(int(^uint(0)>>1)) {
			return int(intVal)
		}
		return intVal
	}
	
	// 尝试转换为浮点数
	if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
		return floatVal
	}
	
	// 尝试转换为时间间隔
	if duration, err := time.ParseDuration(value); err == nil {
		return duration
	}
	
	// 检查是否为数组格式（逗号分隔）
	if strings.Contains(value, ",") {
		parts := strings.Split(value, ",")
		result := make([]string, len(parts))
		for i, part := range parts {
			result[i] = strings.TrimSpace(part)
		}
		return result
	}
	
	// 默认返回字符串
	return value
}

// setNestedValue 设置嵌套值
func (e *EnvManager) setNestedValue(data map[string]interface{}, key string, value interface{}) {
	keys := strings.Split(key, ".")
	current := data

	// 遍历到倒数第二层
	for i := 0; i < len(keys)-1; i++ {
		k := keys[i]
		if _, exists := current[k]; !exists {
			current[k] = make(map[string]interface{})
		}
		if nextMap, ok := current[k].(map[string]interface{}); ok {
			current = nextMap
		} else {
			// 如果不是map，创建新的map覆盖
			current[k] = make(map[string]interface{})
			current = current[k].(map[string]interface{})
		}
	}

	// 设置最终值
	current[keys[len(keys)-1]] = value
}

// GetEnvVar 获取环境变量值
func (e *EnvManager) GetEnvVar(key string) string {
	envKey := e.keyToEnvVar(key)
	return os.Getenv(envKey)
}

// SetEnvVar 设置环境变量值
func (e *EnvManager) SetEnvVar(key, value string) error {
	envKey := e.keyToEnvVar(key)
	return os.Setenv(envKey, value)
}

// HasEnvVar 检查环境变量是否存在
func (e *EnvManager) HasEnvVar(key string) bool {
	envKey := e.keyToEnvVar(key)
	_, exists := os.LookupEnv(envKey)
	return exists
}

// ListEnvVars 列出所有相关的环境变量
func (e *EnvManager) ListEnvVars() map[string]string {
	result := make(map[string]string)
	prefix := e.config.envPrefix
	if prefix != "" && !strings.HasSuffix(prefix, "_") {
		prefix += "_"
	}

	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 {
			continue
		}

		envKey := parts[0]
		envValue := parts[1]

		// 检查是否匹配前缀
		if prefix != "" && !strings.HasPrefix(envKey, prefix) {
			continue
		}

		// 转换为配置键
		configKey := e.envVarToKey(envKey)
		if configKey != "" {
			result[configKey] = envValue
		}
	}

	return result
}

// ClearEnvVars 清除相关的环境变量
func (e *EnvManager) ClearEnvVars() error {
	prefix := e.config.envPrefix
	if prefix != "" && !strings.HasSuffix(prefix, "_") {
		prefix += "_"
	}

	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 {
			continue
		}

		envKey := parts[0]

		// 检查是否匹配前缀
		if prefix != "" && strings.HasPrefix(envKey, prefix) {
			if err := os.Unsetenv(envKey); err != nil {
				return err
			}
		}
	}

	return nil
}
