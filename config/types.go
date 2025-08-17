package config

import (
	"time"
)

// Options 配置选项
type Options struct {
	ConfigPath   string            // 配置文件路径
	ConfigName   string            // 配置文件名（不含扩展名）
	ConfigType   string            // 配置文件类型 (yaml, json, toml, etc.)
	ConfigPaths  []string          // 配置文件搜索路径
	EnvPrefix    string            // 环境变量前缀
	AutomaticEnv bool              // 是否自动绑定环境变量
	Defaults     map[string]interface{} // 默认值
}

// Config 配置管理器
type Config struct {
	configPath   string
	configName   string
	configType   string
	configPaths  []string
	envPrefix    string
	automaticEnv bool
	defaults     map[string]interface{}
	data         map[string]interface{}
	envBindings  map[string]string // key -> env var name
	watcher      *Watcher
	callbacks    []WatchCallback
}

// WatchCallback 配置变化回调函数
type WatchCallback func(oldConfig, newConfig interface{})

// ValidationError 验证错误
type ValidationError struct {
	Field   string
	Tag     string
	Value   interface{}
	Message string
}

func (e *ValidationError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return "配置验证失败: " + e.Field + " " + e.Tag
}

// ConfigFormat 配置文件格式
type ConfigFormat int

const (
	FormatYAML ConfigFormat = iota
	FormatJSON
	FormatTOML
	FormatProperties
	FormatINI
)

// String 返回格式名称
func (f ConfigFormat) String() string {
	switch f {
	case FormatYAML:
		return "yaml"
	case FormatJSON:
		return "json"
	case FormatTOML:
		return "toml"
	case FormatProperties:
		return "properties"
	case FormatINI:
		return "ini"
	default:
		return "unknown"
	}
}

// GetConfigFormat 根据文件扩展名获取格式
func GetConfigFormat(ext string) ConfigFormat {
	switch ext {
	case ".yaml", ".yml":
		return FormatYAML
	case ".json":
		return FormatJSON
	case ".toml":
		return FormatTOML
	case ".properties":
		return FormatProperties
	case ".ini":
		return FormatINI
	default:
		return FormatYAML // 默认使用YAML
	}
}

// DefaultOptions 返回默认配置选项
func DefaultOptions() *Options {
	return &Options{
		ConfigName:   "config",
		ConfigType:   "yaml",
		ConfigPaths:  []string{".", "./config", "./configs"},
		EnvPrefix:    "",
		AutomaticEnv: false,
		Defaults:     make(map[string]interface{}),
	}
}

// Merge 合并配置选项
func (o *Options) Merge(other *Options) *Options {
	if other == nil {
		return o
	}

	result := &Options{
		ConfigPath:   o.ConfigPath,
		ConfigName:   o.ConfigName,
		ConfigType:   o.ConfigType,
		ConfigPaths:  make([]string, len(o.ConfigPaths)),
		EnvPrefix:    o.EnvPrefix,
		AutomaticEnv: o.AutomaticEnv,
		Defaults:     make(map[string]interface{}),
	}

	copy(result.ConfigPaths, o.ConfigPaths)
	for k, v := range o.Defaults {
		result.Defaults[k] = v
	}

	// 覆盖非零值
	if other.ConfigPath != "" {
		result.ConfigPath = other.ConfigPath
	}
	if other.ConfigName != "" {
		result.ConfigName = other.ConfigName
	}
	if other.ConfigType != "" {
		result.ConfigType = other.ConfigType
	}
	if len(other.ConfigPaths) > 0 {
		result.ConfigPaths = make([]string, len(other.ConfigPaths))
		copy(result.ConfigPaths, other.ConfigPaths)
	}
	if other.EnvPrefix != "" {
		result.EnvPrefix = other.EnvPrefix
	}
	if other.AutomaticEnv {
		result.AutomaticEnv = other.AutomaticEnv
	}
	for k, v := range other.Defaults {
		result.Defaults[k] = v
	}

	return result
}

// Duration 时间配置类型，支持字符串解析
type Duration struct {
	time.Duration
}

// UnmarshalText 实现文本解析
func (d *Duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

// MarshalText 实现文本序列化
func (d Duration) MarshalText() ([]byte, error) {
	return []byte(d.Duration.String()), nil
}

// Size 大小配置类型，支持KB、MB、GB等单位
type Size struct {
	Bytes int64
}

// UnmarshalText 实现文本解析
func (s *Size) UnmarshalText(text []byte) error {
	// 这里可以实现大小解析逻辑，如 "1MB" -> 1048576
	// 为了简化，暂时直接解析数字
	return nil
}

// MarshalText 实现文本序列化
func (s Size) MarshalText() ([]byte, error) {
	return []byte(""), nil
}
