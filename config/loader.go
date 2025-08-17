package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Loader 配置加载器
type Loader struct {
	config *Config
}

// NewLoader 创建新的配置加载器
func NewLoader(config *Config) *Loader {
	return &Loader{
		config: config,
	}
}

// LoadFromFile 从文件加载配置
func (l *Loader) LoadFromFile(filePath string) error {
	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("配置文件不存在: %s", filePath)
	}

	// 读取文件内容
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 根据文件扩展名确定格式
	ext := strings.ToLower(filepath.Ext(filePath))
	format := GetConfigFormat(ext)

	// 解析配置
	configData, err := l.parseConfig(data, format)
	if err != nil {
		return fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 合并到现有配置
	l.mergeConfig(configData)

	return nil
}

// LoadFromPath 从路径搜索并加载配置文件
func (l *Loader) LoadFromPath() error {
	var filePath string
	var err error

	// 如果指定了完整路径，直接使用
	if l.config.configPath != "" {
		filePath = l.config.configPath
	} else {
		// 在搜索路径中查找配置文件
		filePath, err = l.FindConfigFile()
		if err != nil {
			return err
		}
	}

	return l.LoadFromFile(filePath)
}

// FindConfigFile 在搜索路径中查找配置文件
func (l *Loader) FindConfigFile() (string, error) {
	// 支持的扩展名
	extensions := []string{".yaml", ".yml", ".json", ".toml", ".properties", ".ini"}

	// 如果指定了配置类型，优先使用对应扩展名
	if l.config.configType != "" {
		switch l.config.configType {
		case "yaml":
			extensions = []string{".yaml", ".yml"}
		case "json":
			extensions = []string{".json"}
		case "toml":
			extensions = []string{".toml"}
		case "properties":
			extensions = []string{".properties"}
		case "ini":
			extensions = []string{".ini"}
		}
	}

	// 在每个搜索路径中查找
	for _, searchPath := range l.config.configPaths {
		for _, ext := range extensions {
			filePath := filepath.Join(searchPath, l.config.configName+ext)
			if _, err := os.Stat(filePath); err == nil {
				return filePath, nil
			}
		}
	}

	return "", fmt.Errorf("未找到配置文件: %s", l.config.configName)
}

// parseConfig 解析配置数据
func (l *Loader) parseConfig(data []byte, format ConfigFormat) (map[string]interface{}, error) {
	var result map[string]interface{}

	switch format {
	case FormatYAML:
		err := yaml.Unmarshal(data, &result)
		if err != nil {
			return nil, fmt.Errorf("解析YAML失败: %w", err)
		}
	case FormatJSON:
		err := json.Unmarshal(data, &result)
		if err != nil {
			return nil, fmt.Errorf("解析JSON失败: %w", err)
		}
	case FormatTOML:
		// TODO: 实现TOML解析
		return nil, fmt.Errorf("TOML格式暂未支持")
	case FormatProperties:
		// TODO: 实现Properties解析
		return nil, fmt.Errorf("Properties格式暂未支持")
	case FormatINI:
		// TODO: 实现INI解析
		return nil, fmt.Errorf("INI格式暂未支持")
	default:
		return nil, fmt.Errorf("不支持的配置格式: %s", format.String())
	}

	return result, nil
}

// mergeConfig 合并配置数据
func (l *Loader) mergeConfig(newData map[string]interface{}) {
	if l.config.data == nil {
		l.config.data = make(map[string]interface{})
	}

	l.deepMerge(l.config.data, newData)
}

// deepMerge 深度合并map
func (l *Loader) deepMerge(dst, src map[string]interface{}) {
	for key, srcVal := range src {
		if dstVal, exists := dst[key]; exists {
			// 如果两个值都是map，递归合并
			if dstMap, dstOk := dstVal.(map[string]interface{}); dstOk {
				if srcMap, srcOk := srcVal.(map[string]interface{}); srcOk {
					l.deepMerge(dstMap, srcMap)
					continue
				}
			}
		}
		// 否则直接覆盖
		dst[key] = srcVal
	}
}

// LoadDefaults 加载默认值
func (l *Loader) LoadDefaults() {
	if l.config.data == nil {
		l.config.data = make(map[string]interface{})
	}

	// 将默认值合并到配置中（不覆盖已存在的值）
	for key, value := range l.config.defaults {
		if !l.hasKey(key) {
			l.setNestedValue(l.config.data, key, value)
		}
	}
}

// hasKey 检查是否存在指定键
func (l *Loader) hasKey(key string) bool {
	_, exists := l.getNestedValue(l.config.data, key)
	return exists
}

// setNestedValue 设置嵌套值
func (l *Loader) setNestedValue(data map[string]interface{}, key string, value interface{}) {
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

// getNestedValue 获取嵌套值
func (l *Loader) getNestedValue(data map[string]interface{}, key string) (interface{}, bool) {
	keys := strings.Split(key, ".")
	current := data

	for i, k := range keys {
		if val, exists := current[k]; exists {
			if i == len(keys)-1 {
				// 最后一个键，返回值
				return val, true
			}
			// 不是最后一个键，继续向下查找
			if nextMap, ok := val.(map[string]interface{}); ok {
				current = nextMap
			} else {
				// 不是map，无法继续查找
				return nil, false
			}
		} else {
			// 键不存在
			return nil, false
		}
	}

	return nil, false
}

// SaveToFile 保存配置到文件
func (l *Loader) SaveToFile(filePath string) error {
	// 根据文件扩展名确定格式
	ext := strings.ToLower(filepath.Ext(filePath))
	format := GetConfigFormat(ext)

	var data []byte
	var err error

	switch format {
	case FormatYAML:
		data, err = yaml.Marshal(l.config.data)
		if err != nil {
			return fmt.Errorf("序列化YAML失败: %w", err)
		}
	case FormatJSON:
		data, err = json.MarshalIndent(l.config.data, "", "  ")
		if err != nil {
			return fmt.Errorf("序列化JSON失败: %w", err)
		}
	default:
		return fmt.Errorf("不支持保存格式: %s", format.String())
	}

	// 确保目录存在
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	// 写入文件
	err = ioutil.WriteFile(filePath, data, 0644)
	if err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}

	return nil
}
