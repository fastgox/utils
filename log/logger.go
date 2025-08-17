package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// 配置结构
type config struct {
	BaseDir string // 基础目录，如 "logs"
}

// Logger 实例，每个事件类型一个独立的logger
type Logger struct {
	cfg       *config
	eventType string                    // 事件类型
	loggers   map[string]*log.Logger    // 按需创建的logger
	writers   map[string]io.WriteCloser // 管理文件句柄
	mu        sync.RWMutex              // 保护并发访问
}

var (
	defaultLogger *Logger
	loggerMap     = make(map[string]*Logger) // 存储不同事件类型的logger
	mapMu         sync.RWMutex               // 保护loggerMap
)

// NewLogger 创建新的Logger实例（使用路径类型）
func NewLogger(basePath string) (*Logger, error) {
	// 清理路径
	if basePath == "" {
		basePath = "logs"
	}
	basePath = filepath.Clean(basePath)

	cfg := &config{
		BaseDir: basePath,
	}

	// 创建Logger实例（默认logger没有特定事件类型）
	logger := &Logger{
		cfg:       cfg,
		eventType: "app", // 默认事件类型
		loggers:   make(map[string]*log.Logger),
		writers:   make(map[string]io.WriteCloser),
	}

	return logger, nil
}

// InitWithPath 使用路径类型初始化默认logger
func InitWithPath(basePath string) error {
	logger, err := NewLogger(basePath)
	if err != nil {
		return err
	}
	defaultLogger = logger
	return nil
}

// GetLoggerWithBaseDir 获取指定事件类型和基础目录的logger
func GetLoggerWithBaseDir(eventType, baseDir string) (*Logger, error) {
	if baseDir == "" {
		baseDir = "logs"
	}

	loggerKey := fmt.Sprintf("%s_%s", baseDir, eventType)

	mapMu.RLock()
	if logger, exists := loggerMap[loggerKey]; exists {
		mapMu.RUnlock()
		return logger, nil
	}
	mapMu.RUnlock()

	mapMu.Lock()
	defer mapMu.Unlock()

	// 双重检查
	if logger, exists := loggerMap[loggerKey]; exists {
		return logger, nil
	}

	// 创建新的logger配置
	cfg := &config{
		BaseDir: baseDir,
	}

	logger := &Logger{
		cfg:       cfg,
		eventType: eventType,
		loggers:   make(map[string]*log.Logger),
		writers:   make(map[string]io.WriteCloser),
	}

	loggerMap[loggerKey] = logger
	return logger, nil
}

// GetLogger 获取指定事件类型的logger，如果不存在则创建（使用默认baseDir）
func GetLogger(eventType string) (*Logger, error) {
	baseDir := "logs" // 默认值
	if defaultLogger != nil {
		baseDir = defaultLogger.cfg.BaseDir
	}
	return GetLoggerWithBaseDir(eventType, baseDir)
}

// createLogger 创建指定级别的logger（假设已经持有锁）
func (l *Logger) createLogger(level string) *log.Logger {
	writer := l.getWriterUnsafe(level)
	return log.New(writer, "", log.LstdFlags)
}

// getWriterUnsafe 获取指定级别的文件写入器（不加锁，内部使用）
func (l *Logger) getWriterUnsafe(level string) io.Writer {
	// 构建文件路径: baseDir/日期/eventType/level.log
	today := time.Now().Format("2006-01-02")
	logDir := filepath.Join(l.cfg.BaseDir, today, l.eventType)
	logFile := filepath.Join(logDir, level+".log")

	// 检查是否已经有这个文件的writer
	writerKey := fmt.Sprintf("%s_%s", today, level)
	if writer, exists := l.writers[writerKey]; exists {
		return writer
	}

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

	// 存储writer以便复用和后续关闭
	l.writers[writerKey] = file
	return file
}

// Close 关闭所有文件句柄
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	var lastErr error
	for key, writer := range l.writers {
		if err := writer.Close(); err != nil {
			lastErr = err
		}
		delete(l.writers, key)
	}
	return lastErr
}

// getOrCreateLogger 懒加载获取指定level的logger
func (l *Logger) getOrCreateLogger(level string) *log.Logger {
	l.mu.RLock()
	if logger, exists := l.loggers[level]; exists {
		l.mu.RUnlock()
		return logger
	}
	l.mu.RUnlock()

	l.mu.Lock()
	defer l.mu.Unlock()

	// 双重检查
	if logger, exists := l.loggers[level]; exists {
		return logger
	}

	// 创建新的logger
	logger := l.createLogger(level)
	l.loggers[level] = logger
	return logger
}

// Logger实例方法
func (l *Logger) Debug(format string, v ...interface{}) {
	logger := l.getOrCreateLogger("debug")
	logger.Printf(format, v...)
}

func (l *Logger) Info(format string, v ...interface{}) {
	logger := l.getOrCreateLogger("info")
	logger.Printf(format, v...)
}

func (l *Logger) Warn(format string, v ...interface{}) {
	logger := l.getOrCreateLogger("warn")
	logger.Printf(format, v...)
}

func (l *Logger) Error(format string, v ...interface{}) {
	logger := l.getOrCreateLogger("error")
	logger.Printf(format, v...)
}

// 全局方法（使用默认logger，保持向后兼容）
func Debug(format string, v ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Debug(format, v...)
	}
}

func Info(format string, v ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Info(format, v...)
	}
}

func Warn(format string, v ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Warn(format, v...)
	}
}

func Error(format string, v ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Error(format, v...)
	}
}

// InitDefault 便捷函数，使用默认目录初始化
func InitDefault() error {
	return InitWithPath("logs")
}

// CloseAll 关闭所有logger的文件句柄
func CloseAll() {
	mapMu.Lock()
	defer mapMu.Unlock()

	for _, logger := range loggerMap {
		logger.Close()
	}

	if defaultLogger != nil {
		defaultLogger.Close()
	}
}
