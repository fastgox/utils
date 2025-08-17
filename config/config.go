package config

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	// 全局配置实例
	globalConfig *Config
)

// Init 使用配置文件路径初始化
func Init(configPath string) error {
	opts := DefaultOptions()
	opts.ConfigPath = configPath
	return InitWithOptions(opts)
}

// InitWithOptions 使用选项初始化
func InitWithOptions(opts *Options) error {
	if opts == nil {
		opts = DefaultOptions()
	}

	// 创建配置实例
	config := &Config{
		configPath:   opts.ConfigPath,
		configName:   opts.ConfigName,
		configType:   opts.ConfigType,
		configPaths:  opts.ConfigPaths,
		envPrefix:    opts.EnvPrefix,
		automaticEnv: opts.AutomaticEnv,
		defaults:     make(map[string]interface{}),
		data:         make(map[string]interface{}),
		envBindings:  make(map[string]string),
	}

	// 复制默认值
	for k, v := range opts.Defaults {
		config.defaults[k] = v
	}

	// 加载默认值
	loader := NewLoader(config)
	loader.LoadDefaults()

	// 尝试加载配置文件（如果失败，只使用默认值）
	err := loader.LoadFromPath()
	if err != nil {
		// 如果没有指定配置路径，或者文件不存在，只使用默认值
		if opts.ConfigPath == "" {
			// 这是正常情况，只使用默认值
		} else {
			return fmt.Errorf("加载配置文件失败: %w", err)
		}
	}

	// 加载环境变量
	envManager := NewEnvManager(config)
	envManager.LoadEnvVars()

	// 设置全局配置
	globalConfig = config

	return nil
}

// InitDefault 使用默认配置初始化
func InitDefault() error {
	return InitWithOptions(DefaultOptions())
}

// SetDefault 设置默认值
func SetDefault(key string, value interface{}) {
	ensureGlobalConfig()
	if globalConfig.defaults == nil {
		globalConfig.defaults = make(map[string]interface{})
	}
	globalConfig.defaults[key] = value

	// 如果配置中还没有这个值，设置它
	if !hasKey(key) {
		setNestedValue(globalConfig.data, key, value)
	}
}

// Get 获取配置值
func Get(key string) interface{} {
	ensureGlobalConfig()
	value, _ := getNestedValue(globalConfig.data, key)
	return value
}

// GetString 获取字符串值
func GetString(key string) string {
	value := Get(key)
	if value == nil {
		return ""
	}
	return fmt.Sprintf("%v", value)
}

// GetStringDefault 获取字符串值，带默认值
func GetStringDefault(key, defaultValue string) string {
	value := GetString(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetInt 获取整数值
func GetInt(key string) int {
	value := Get(key)
	if value == nil {
		return 0
	}

	switch v := value.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	case string:
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return 0
}

// GetIntDefault 获取整数值，带默认值
func GetIntDefault(key string, defaultValue int) int {
	value := GetInt(key)
	if value == 0 {
		return defaultValue
	}
	return value
}

// GetBool 获取布尔值
func GetBool(key string) bool {
	value := Get(key)
	if value == nil {
		return false
	}

	switch v := value.(type) {
	case bool:
		return v
	case string:
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return false
}

// GetFloat64 获取浮点数值
func GetFloat64(key string) float64 {
	value := Get(key)
	if value == nil {
		return 0
	}

	switch v := value.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	case int64:
		return float64(v)
	case string:
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	}
	return 0
}

// GetStringSlice 获取字符串切片
func GetStringSlice(key string) []string {
	value := Get(key)
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case []string:
		return v
	case []interface{}:
		result := make([]string, len(v))
		for i, item := range v {
			result[i] = fmt.Sprintf("%v", item)
		}
		return result
	case string:
		// 尝试解析逗号分隔的字符串
		if strings.Contains(v, ",") {
			parts := strings.Split(v, ",")
			result := make([]string, len(parts))
			for i, part := range parts {
				result[i] = strings.TrimSpace(part)
			}
			return result
		}
		return []string{v}
	}
	return nil
}

// GetDuration 获取时间间隔
func GetDuration(key string) time.Duration {
	value := Get(key)
	if value == nil {
		return 0
	}

	switch v := value.(type) {
	case time.Duration:
		return v
	case string:
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	case int64:
		return time.Duration(v)
	case int:
		return time.Duration(v)
	}
	return 0
}

// Unmarshal 将配置绑定到结构体
func Unmarshal(v interface{}) error {
	ensureGlobalConfig()
	return unmarshalData(globalConfig.data, v)
}

// UnmarshalKey 将指定键的配置绑定到结构体
func UnmarshalKey(key string, v interface{}) error {
	data := Get(key)
	if data == nil {
		return fmt.Errorf("配置键不存在: %s", key)
	}
	return unmarshalData(data, v)
}

// unmarshalData 将数据绑定到结构体
func unmarshalData(data interface{}, v interface{}) error {
	// 预处理数据，处理特殊类型
	processedData := preprocessData(data)

	// 使用JSON作为中间格式进行转换
	jsonData, err := json.Marshal(processedData)
	if err != nil {
		return fmt.Errorf("序列化配置数据失败: %w", err)
	}

	err = json.Unmarshal(jsonData, v)
	if err != nil {
		return fmt.Errorf("反序列化到结构体失败: %w", err)
	}

	return nil
}

// preprocessData 预处理数据，转换特殊类型
func preprocessData(data interface{}) interface{} {
	switch v := data.(type) {
	case map[string]interface{}:
		result := make(map[string]interface{})
		for key, value := range v {
			result[key] = preprocessData(value)
		}
		return result
	case []interface{}:
		result := make([]interface{}, len(v))
		for i, item := range v {
			result[i] = preprocessData(item)
		}
		return result
	case time.Duration:
		// 将time.Duration转换为纳秒数
		return int64(v)
	case string:
		// 尝试解析时间间隔字符串
		if duration, err := time.ParseDuration(v); err == nil {
			return int64(duration)
		}
		return v
	default:
		return v
	}
}

// SetEnvPrefix 设置环境变量前缀
func SetEnvPrefix(prefix string) {
	ensureGlobalConfig()
	globalConfig.envPrefix = prefix
}

// BindEnv 绑定环境变量
func BindEnv(key string) error {
	ensureGlobalConfig()
	envManager := NewEnvManager(globalConfig)
	return envManager.BindEnv(key)
}

// AutomaticEnv 启用自动环境变量绑定
func AutomaticEnv() {
	ensureGlobalConfig()
	globalConfig.automaticEnv = true

	// 重新加载环境变量
	envManager := NewEnvManager(globalConfig)
	envManager.LoadEnvVars()
}

// Watch 监听配置文件变化
func Watch(callback WatchCallback) error {
	ensureGlobalConfig()

	if globalConfig.watcher == nil {
		watcher, err := NewWatcher(globalConfig)
		if err != nil {
			return err
		}
		globalConfig.watcher = watcher
	}

	globalConfig.watcher.AddCallback(callback)

	// 如果还没有开始监听，启动监听
	if !globalConfig.watcher.IsRunning() {
		configPath := globalConfig.configPath
		if configPath == "" {
			// 尝试找到配置文件路径
			loader := NewLoader(globalConfig)
			var err error
			configPath, err = loader.FindConfigFile()
			if err != nil {
				return fmt.Errorf("无法找到配置文件进行监听: %w", err)
			}
		}
		return globalConfig.watcher.Start(configPath)
	}

	return nil
}

// StopWatch 停止监听配置文件
func StopWatch() error {
	ensureGlobalConfig()
	if globalConfig.watcher != nil {
		return globalConfig.watcher.Stop()
	}
	return nil
}

// Validate 验证当前配置
func Validate() error {
	ensureGlobalConfig()
	validator := NewValidator(globalConfig)
	return validator.Validate()
}

// ValidateStruct 验证结构体
func ValidateStruct(v interface{}) error {
	ensureGlobalConfig()
	validator := NewValidator(globalConfig)
	return validator.ValidateStruct(v)
}

// WriteConfig 保存配置到原文件
func WriteConfig() error {
	ensureGlobalConfig()
	if globalConfig.configPath == "" {
		return fmt.Errorf("未指定配置文件路径")
	}
	loader := NewLoader(globalConfig)
	return loader.SaveToFile(globalConfig.configPath)
}

// WriteConfigAs 保存配置到指定文件
func WriteConfigAs(filename string) error {
	ensureGlobalConfig()
	loader := NewLoader(globalConfig)
	return loader.SaveToFile(filename)
}

// Reset 重置全局配置（主要用于测试）
func Reset() {
	if globalConfig != nil && globalConfig.watcher != nil {
		globalConfig.watcher.Stop()
	}
	globalConfig = nil
}

// 辅助函数

// ensureGlobalConfig 确保全局配置已初始化
func ensureGlobalConfig() {
	if globalConfig == nil {
		// 使用默认配置初始化
		InitDefault()
	}
}

// hasKey 检查是否存在指定键
func hasKey(key string) bool {
	ensureGlobalConfig()
	_, exists := getNestedValue(globalConfig.data, key)
	return exists
}

// getNestedValue 获取嵌套值
func getNestedValue(data map[string]interface{}, key string) (interface{}, bool) {
	keys := strings.Split(key, ".")
	current := data

	for i, k := range keys {
		if val, exists := current[k]; exists {
			if i == len(keys)-1 {
				return val, true
			}
			if nextMap, ok := val.(map[string]interface{}); ok {
				current = nextMap
			} else {
				return nil, false
			}
		} else {
			return nil, false
		}
	}

	return nil, false
}

// setNestedValue 设置嵌套值
func setNestedValue(data map[string]interface{}, key string, value interface{}) {
	keys := strings.Split(key, ".")
	current := data

	for i := 0; i < len(keys)-1; i++ {
		k := keys[i]
		if _, exists := current[k]; !exists {
			current[k] = make(map[string]interface{})
		}
		if nextMap, ok := current[k].(map[string]interface{}); ok {
			current = nextMap
		} else {
			current[k] = make(map[string]interface{})
			current = current[k].(map[string]interface{})
		}
	}

	current[keys[len(keys)-1]] = value
}
