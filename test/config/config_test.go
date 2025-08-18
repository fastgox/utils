package config_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/fastgox/utils/config"
)

// TestConfig 配置结构体
type TestConfig struct {
	App struct {
		Name    string `config:"name" validate:"required"`
		Version string `config:"version" validate:"required"`
		Debug   bool   `config:"debug"`
	} `config:"app"`

	Server struct {
		Host    string        `config:"host" validate:"required"`
		Port    int           `config:"port" validate:"min=1,max=65535"`
		Timeout time.Duration `config:"timeout"`
	} `config:"server"`

	Database struct {
		Host     string `config:"host" validate:"required"`
		Port     int    `config:"port" validate:"min=1,max=65535"`
		Username string `config:"username" validate:"required"`
		Password string `config:"password" validate:"required"`
		DBName   string `config:"dbname" validate:"required"`
	} `config:"database"`

	Redis struct {
		Host     string `config:"host"`
		Port     int    `config:"port"`
		Password string `config:"password"`
		DB       int    `config:"db"`
	} `config:"redis"`
}

func TestConfigBasic(t *testing.T) {
	// 重置全局配置
	config.Reset()

	// 创建测试配置文件
	configContent := `
app:
  name: "test-app"
  version: "1.0.0"
  debug: true

server:
  host: "localhost"
  port: 8080
  timeout: "30s"

database:
  host: "localhost"
  port: 3306
  username: "root"
  password: "password"
  dbname: "testdb"

redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0
`

	// 创建临时配置文件
	configPath := filepath.Join("test_configs", "config.yaml")
	err := os.MkdirAll(filepath.Dir(configPath), 0755)
	if err != nil {
		t.Fatalf("创建配置目录失败: %v", err)
	}

	err = os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("创建配置文件失败: %v", err)
	}
	defer os.RemoveAll("test_configs")

	// 初始化配置
	err = config.Init(configPath)
	if err != nil {
		t.Fatalf("初始化配置失败: %v", err)
	}

	// 测试基本配置获取
	appName := config.GetString("app.name")
	if appName != "test-app" {
		t.Errorf("期望 app.name = 'test-app', 实际得到: %s", appName)
	}

	serverPort := config.GetInt("server.port")
	if serverPort != 8080 {
		t.Errorf("期望 server.port = 8080, 实际得到: %d", serverPort)
	}

	debugMode := config.GetBool("app.debug")
	if !debugMode {
		t.Errorf("期望 app.debug = true, 实际得到: %v", debugMode)
	}

	timeout := config.GetDuration("server.timeout")
	if timeout != 30*time.Second {
		t.Errorf("期望 server.timeout = 30s, 实际得到: %v", timeout)
	}

	t.Logf("基本配置测试通过")
}

func TestConfigStructBinding(t *testing.T) {
	// 重置全局配置
	config.Reset()

	// 创建测试配置文件
	configContent := `
app:
  name: "struct-test-app"
  version: "2.0.0"
  debug: false

server:
  host: "0.0.0.0"
  port: 9090
  timeout: "60s"

database:
  host: "db.example.com"
  port: 5432
  username: "dbuser"
  password: "dbpass"
  dbname: "myapp"

redis:
  host: "redis.example.com"
  port: 6380
  password: "redispass"
  db: 1
`

	configPath := filepath.Join("test_configs", "struct_config.yaml")
	err := os.MkdirAll(filepath.Dir(configPath), 0755)
	if err != nil {
		t.Fatalf("创建配置目录失败: %v", err)
	}

	err = os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("创建配置文件失败: %v", err)
	}
	defer os.RemoveAll("test_configs")

	// 初始化配置
	err = config.Init(configPath)
	if err != nil {
		t.Fatalf("初始化配置失败: %v", err)
	}

	// 测试结构体绑定
	var cfg TestConfig
	err = config.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("结构体绑定失败: %v", err)
	}

	// 验证绑定结果
	if cfg.App.Name != "struct-test-app" {
		t.Errorf("期望 App.Name = 'struct-test-app', 实际得到: %s", cfg.App.Name)
	}

	if cfg.Server.Port != 9090 {
		t.Errorf("期望 Server.Port = 9090, 实际得到: %d", cfg.Server.Port)
	}

	if cfg.Database.Host != "db.example.com" {
		t.Errorf("期望 Database.Host = 'db.example.com', 实际得到: %s", cfg.Database.Host)
	}

	if cfg.Redis.DB != 1 {
		t.Errorf("期望 Redis.DB = 1, 实际得到: %d", cfg.Redis.DB)
	}

	t.Logf("结构体绑定测试通过")
}

func TestConfigEnvOverride(t *testing.T) {
	// 重置全局配置
	config.Reset()

	// 创建基础配置文件
	configContent := `
app:
  name: "env-test-app"
  debug: false

server:
  host: "localhost"
  port: 8080
`

	configPath := filepath.Join("test_configs", "env_config.yaml")
	err := os.MkdirAll(filepath.Dir(configPath), 0755)
	if err != nil {
		t.Fatalf("创建配置目录失败: %v", err)
	}

	err = os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("创建配置文件失败: %v", err)
	}
	defer os.RemoveAll("test_configs")

	// 设置环境变量
	os.Setenv("TESTAPP_APP_DEBUG", "true")
	os.Setenv("TESTAPP_SERVER_PORT", "9000")
	defer func() {
		os.Unsetenv("TESTAPP_APP_DEBUG")
		os.Unsetenv("TESTAPP_SERVER_PORT")
	}()

	// 初始化配置
	opts := config.DefaultOptions()
	opts.ConfigPath = configPath
	opts.EnvPrefix = "TESTAPP"
	err = config.InitWithOptions(opts)
	if err != nil {
		t.Fatalf("初始化配置失败: %v", err)
	}

	// 绑定环境变量
	config.BindEnv("app.debug")
	config.BindEnv("server.port")

	// 重新加载环境变量
	config.AutomaticEnv()

	// 验证环境变量覆盖
	debug := config.GetBool("app.debug")
	if !debug {
		t.Errorf("期望环境变量覆盖 app.debug = true, 实际得到: %v", debug)
	}

	port := config.GetInt("server.port")
	if port != 9000 {
		t.Errorf("期望环境变量覆盖 server.port = 9000, 实际得到: %d", port)
	}

	t.Logf("环境变量覆盖测试通过")
}

func TestConfigDefaults(t *testing.T) {
	// 重置全局配置
	config.Reset()

	// 不创建配置文件，只使用默认值
	opts := config.DefaultOptions()
	opts.ConfigName = "nonexistent"
	opts.ConfigPaths = []string{"./nonexistent"}

	// 在选项中设置默认值
	opts.Defaults["app.name"] = "default-app"
	opts.Defaults["server.port"] = 3000
	opts.Defaults["app.debug"] = true

	// 尝试初始化（会失败，但默认值应该可用）
	config.InitWithOptions(opts)

	// 验证默认值
	appName := config.GetString("app.name")
	if appName != "default-app" {
		t.Errorf("期望默认值 app.name = 'default-app', 实际得到: %s", appName)
	}

	serverPort := config.GetInt("server.port")
	if serverPort != 3000 {
		t.Errorf("期望默认值 server.port = 3000, 实际得到: %d", serverPort)
	}

	debug := config.GetBool("app.debug")
	if !debug {
		t.Errorf("期望默认值 app.debug = true, 实际得到: %v", debug)
	}

	t.Logf("默认值测试通过")
}

func TestConfigValidation(t *testing.T) {
	// 重置全局配置
	config.Reset()

	// 创建无效配置
	configContent := `
app:
  name: ""  # 违反required规则
  version: "1.0.0"

server:
  host: "localhost"
  port: 70000  # 违反max=65535规则

database:
  host: "localhost"
  port: 3306
  username: "root"
  password: ""  # 违反required规则
  dbname: "test"
`

	configPath := filepath.Join("test_configs", "invalid_config.yaml")
	err := os.MkdirAll(filepath.Dir(configPath), 0755)
	if err != nil {
		t.Fatalf("创建配置目录失败: %v", err)
	}

	err = os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("创建配置文件失败: %v", err)
	}
	defer os.RemoveAll("test_configs")

	// 初始化配置
	err = config.Init(configPath)
	if err != nil {
		t.Fatalf("初始化配置失败: %v", err)
	}

	// 测试结构体验证
	var cfg TestConfig
	err = config.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("结构体绑定失败: %v", err)
	}

	// 验证配置（应该失败）
	err = config.ValidateStruct(&cfg)
	if err == nil {
		t.Errorf("期望验证失败，但验证通过了")
	} else {
		t.Logf("验证失败（符合预期）: %v", err)
	}

	t.Logf("配置验证测试通过")
}
