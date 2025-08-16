package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// Config YAML配置结构
type Config struct {
	BaseDir string `yaml:"base_dir"` // 基础目录，如 "logs"
	LogType string `yaml:"log_type"` // 日志类型，如 "app"
}

var (
	config      *Config
	debugLogger *log.Logger
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
)

// Init 从YAML文件初始化logger
func Init(configFile string) error {
	// 读取配置文件
	data, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 解析YAML
	config = &Config{}
	if err := yaml.Unmarshal(data, config); err != nil {
		return fmt.Errorf("解析YAML配置失败: %w", err)
	}

	// 设置默认值
	if config.BaseDir == "" {
		config.BaseDir = "logs"
	}
	if config.LogType == "" {
		config.LogType = "app"
	}

	// 初始化各级别的logger
	debugLogger = createLogger("debug")
	infoLogger = createLogger("info")
	warnLogger = createLogger("warn")
	errorLogger = createLogger("error")

	return nil
}

// createLogger 创建指定级别的logger
func createLogger(level string) *log.Logger {
	writer := getWriter(level)
	return log.New(writer, "", log.LstdFlags)
}

// getWriter 获取指定级别的文件写入器
func getWriter(level string) io.Writer {
	// 构建文件路径: baseDir/日期/logType/level.log
	today := time.Now().Format("2006-01-02")
	logDir := filepath.Join(config.BaseDir, today, config.LogType)
	logFile := filepath.Join(logDir, level+".log")

	// 确保目录存在
	if err := os.MkdirAll(logDir, 0755); err != nil {
		// 如果创建目录失败，返回标准输出
		return os.Stdout
	}

	// 打开文件
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		// 如果打开文件失败，返回标准输出
		return os.Stdout
	}

	return file
}

// Debug 调试日志
func Debug(v ...interface{}) {
	if debugLogger != nil {
		debugLogger.Print(v...)
	}
}

// Debugf 格式化调试日志
func Debugf(format string, v ...interface{}) {
	if debugLogger != nil {
		debugLogger.Printf(format, v...)
	}
}

// Info 信息日志
func Info(v ...interface{}) {
	if infoLogger != nil {
		infoLogger.Print(v...)
	}
}

// Infof 格式化信息日志
func Infof(format string, v ...interface{}) {
	if infoLogger != nil {
		infoLogger.Printf(format, v...)
	}
}

// Warn 警告日志
func Warn(v ...interface{}) {
	if warnLogger != nil {
		warnLogger.Print(v...)
	}
}

// Warnf 格式化警告日志
func Warnf(format string, v ...interface{}) {
	if warnLogger != nil {
		warnLogger.Printf(format, v...)
	}
}

// Error 错误日志
func Error(v ...interface{}) {
	if errorLogger != nil {
		errorLogger.Print(v...)
	}
}

// Errorf 格式化错误日志
func Errorf(format string, v ...interface{}) {
	if errorLogger != nil {
		errorLogger.Printf(format, v...)
	}
}

// InitDefault 便捷函数，自动从log.yml初始化
func InitDefault() error {
	return Init("log.yml")
}

// InitWithPath 便捷函数，支持自定义配置文件路径
func InitWithPath(configPath string) error {
	if configPath == "" {
		configPath = "log.yml" // 默认路径
	}
	return Init(configPath)
}

// InitFromEnv 从环境变量获取配置文件路径，支持默认值
func InitFromEnv(envKey string, defaultPath string) error {
	configPath := os.Getenv(envKey)
	if configPath == "" {
		if defaultPath == "" {
			defaultPath = "log.yml"
		}
		configPath = defaultPath
	}
	return Init(configPath)
}
