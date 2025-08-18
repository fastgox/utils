package logger_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/fastgox/utils/logger"
)

func TestLogger(t *testing.T) {
	// 初始化日志系统
	err := logger.InitWithPath("test_logs")
	if err != nil {
		t.Fatalf("InitWithPath失败: %v", err)
	}

	// 获取不同事件类型的logger
	userLogger, err := logger.GetLogger("user")
	if err != nil {
		t.Fatalf("获取user logger失败: %v", err)
	}

	apiLogger, err := logger.GetLogger("api")
	if err != nil {
		t.Fatalf("获取api logger失败: %v", err)
	}

	// 记录不同级别的日志
	userLogger.Info("用户登录: userID=%d", 12345)
	userLogger.Debug("用户详情: %s", "调试信息")
	apiLogger.Warn("API响应慢: duration=%dms", 1500)
	apiLogger.Error("API错误: %s", "连接超时")

	// 验证文件创建
	today := time.Now().Format("2006-01-02")
	userInfoFile := filepath.Join("test_logs", today, "user", "info.log")
	apiWarnFile := filepath.Join("test_logs", today, "api", "warn.log")

	if _, err := os.Stat(userInfoFile); os.IsNotExist(err) {
		t.Error("应该创建user/info.log文件")
	}

	if _, err := os.Stat(apiWarnFile); os.IsNotExist(err) {
		t.Error("应该创建api/warn.log文件")
	}

	// 关闭所有logger
	logger.CloseAll()
}
func TestLogger2(t *testing.T) {
	logger.InitWithPath("test_logs")
	// 初始化日志系统
	logger.Info("测试日志")
}
