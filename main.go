package main

import (
	"fmt"
	"os"
	"time"

	"github.com/fastgox/utils/logger"
)

func main() {
	fmt.Println("=== Utils Library Demo ===")
	fmt.Println("Demonstrating logger functionality...")

	// 创建示例配置文件
	createExampleConfig()

	// 初始化logger
	err := logger.InitDefault()
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		return
	}

	// 演示各种日志级别
	fmt.Println("\n1. Testing different log levels:")
	logger.Debug("This is a debug message")
	logger.Info("Application started successfully")
	logger.Warn("This is a warning message")
	logger.Error("This is an error message")

	// 演示格式化日志
	fmt.Println("\n2. Testing formatted logging:")
	username := "helwd"
	userID := 12345
	logger.Infof("User %s (ID: %d) logged in at %s", username, userID, time.Now().Format("15:04:05"))
	logger.Debugf("Processing request for user: %s", username)
	logger.Warnf("User %s attempted %d failed logins", username, 3)
	logger.Errorf("Failed to process payment for user %s: insufficient funds", username)

	// 演示不同的初始化方法
	fmt.Println("\n3. Testing different initialization methods:")
	
	// 使用环境变量初始化
	os.Setenv("LOG_CONFIG", "log.yml")
	err = logger.InitFromEnv("LOG_CONFIG", "default.yml")
	if err != nil {
		logger.Errorf("Failed to init from env: %v", err)
	} else {
		logger.Info("Successfully initialized from environment variable")
	}

	// 使用自定义路径初始化
	err = logger.InitWithPath("log.yml")
	if err != nil {
		logger.Errorf("Failed to init with custom path: %v", err)
	} else {
		logger.Info("Successfully initialized with custom path")
	}

	fmt.Println("\n=== Demo completed ===")
	fmt.Println("Check the 'logs' directory for generated log files!")
	fmt.Println("Log files are organized by date and log level.")
}

// createExampleConfig 创建示例配置文件
func createExampleConfig() {
	configContent := `base_dir: "logs"
log_type: "demo"`

	err := os.WriteFile("log.yml", []byte(configContent), 0644)
	if err != nil {
		fmt.Printf("Warning: Could not create example config file: %v\n", err)
	}
}
