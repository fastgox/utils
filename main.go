package main

import (
	"fmt"
	"time"

	"github.com/fastgox/utils/config"
	"github.com/fastgox/utils/crypto"
	"github.com/fastgox/utils/jwt"
	"github.com/fastgox/utils/logger"
)

// AppConfig 应用配置结构体
type AppConfig struct {
	App struct {
		Name        string `config:"name" json:"name" validate:"required"`
		Version     string `config:"version" json:"version" validate:"required"`
		Debug       bool   `config:"debug" json:"debug"`
		Environment string `config:"environment" json:"environment" validate:"oneof=development testing production"`
	} `config:"app" json:"app"`

	Server struct {
		Host           string        `config:"host" json:"host" validate:"required"`
		Port           int           `config:"port" json:"port" validate:"min=1,max=65535"`
		Timeout        time.Duration `config:"timeout" json:"timeout"`
		MaxConnections int           `config:"max_connections" json:"max_connections" validate:"min=1"`
	} `config:"server" json:"server"`

	Database struct {
		Host           string        `config:"host" json:"host" validate:"required"`
		Port           int           `config:"port" json:"port" validate:"min=1,max=65535"`
		Username       string        `config:"username" json:"username" validate:"required"`
		Password       string        `config:"password" json:"password" validate:"required"`
		DBName         string        `config:"dbname" json:"dbname" validate:"required"`
		MaxConnections int           `config:"max_connections" json:"max_connections" validate:"min=1"`
		Timeout        time.Duration `config:"timeout" json:"timeout"`
	} `config:"database" json:"database"`

	JWT struct {
		Secret     string        `config:"secret" json:"secret" validate:"required"`
		Issuer     string        `config:"issuer" json:"issuer" validate:"required"`
		Expiration time.Duration `config:"expiration" json:"expiration"`
	} `config:"jwt" json:"jwt"`

	Log struct {
		Level   string `config:"level" json:"level" validate:"oneof=debug info warn error"`
		Format  string `config:"format" json:"format" validate:"oneof=json text"`
		Output  string `config:"output" json:"output" validate:"required"`
		MaxSize string `config:"max_size" json:"max_size"`
		MaxAge  string `config:"max_age" json:"max_age"`
	} `config:"log" json:"log"`
}

func main() {
	fmt.Println("=== Utils Library Demo ===")

	// 演示Config工具
	demoConfig()

	// 演示Crypto工具
	demoCrypto()

	fmt.Println("\n=== Demo completed ===")
	fmt.Println("Check the 'logs' directory for generated log files!")
}

func demoConfig() {
	fmt.Println("\n🔧 Config 工具演示:")

	// 初始化配置
	err := config.Init("config.yaml")
	if err != nil {
		fmt.Printf("❌ 配置初始化失败: %v\n", err)
		return
	}

	fmt.Println("✅ 配置文件加载成功")

	// 演示基本配置获取
	appName := config.GetString("app.name")
	serverPort := config.GetInt("server.port")
	debugMode := config.GetBool("app.debug")
	serverTimeout := config.GetDuration("server.timeout")

	fmt.Printf("📋 基本配置:")
	fmt.Printf("  应用名称: %s\n", appName)
	fmt.Printf("  服务端口: %d\n", serverPort)
	fmt.Printf("  调试模式: %v\n", debugMode)
	fmt.Printf("  服务超时: %v\n", serverTimeout)

	// 演示结构体绑定
	var cfg AppConfig

	err = config.Unmarshal(&cfg)
	if err != nil {
		fmt.Printf("❌ 结构体绑定失败: %v\n", err)
		return
	}

	fmt.Printf("\n📦 结构体绑定:")
	fmt.Printf("  应用: %s v%s (%s)\n", cfg.App.Name, cfg.App.Version, cfg.App.Environment)
	fmt.Printf("  服务器: %s:%d (超时: %v, 最大连接: %d)\n", cfg.Server.Host, cfg.Server.Port, cfg.Server.Timeout, cfg.Server.MaxConnections)
	fmt.Printf("  数据库: %s@%s:%d/%s (最大连接: %d)\n", cfg.Database.Username, cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName, cfg.Database.MaxConnections)

	// 演示配置验证
	err = config.ValidateStruct(&cfg)
	if err != nil {
		fmt.Printf("❌ 配置验证失败: %v\n", err)
	} else {
		fmt.Println("✅ 配置验证通过")
	}

	// 演示环境变量覆盖
	config.SetEnvPrefix("HELWD")
	config.BindEnv("app.debug")
	config.BindEnv("server.port")

	fmt.Printf("\n🌍 环境变量支持:")
	fmt.Printf("  设置环境变量 HELWD_APP_DEBUG=false 可覆盖 app.debug\n")
	fmt.Printf("  设置环境变量 HELWD_SERVER_PORT=9090 可覆盖 server.port\n")

	// 演示默认值
	config.SetDefault("features.new_feature", true)
	newFeature := config.GetBool("features.new_feature")
	fmt.Printf("\n⚙️  默认值设置:")
	fmt.Printf("  新功能开关: %v (默认值)\n", newFeature)

	// 初始化其他工具使用配置
	initOtherTools(&cfg)
}

func initOtherTools(cfg *AppConfig) {
	fmt.Printf("\n🔗 使用配置初始化其他工具:")

	// 初始化日志工具
	err := logger.InitWithPath("logs")
	if err != nil {
		fmt.Printf("❌ 日志工具初始化失败: %v\n", err)
	} else {
		fmt.Printf("✅ 日志工具初始化成功 (输出目录: %s)\n", cfg.Log.Output)

		// 记录一些日志
		appLogger, _ := logger.GetLogger("app")
		appLogger.Info("应用启动: %s v%s", cfg.App.Name, cfg.App.Version)
		appLogger.Debug("调试模式: %v", cfg.App.Debug)
	}

	// 初始化JWT工具
	jwt.Init(cfg.JWT.Secret, cfg.JWT.Issuer, cfg.JWT.Expiration)
	fmt.Printf("✅ JWT工具初始化成功 (签发者: %s, 过期时间: %v)\n", cfg.JWT.Issuer, cfg.JWT.Expiration)

	// 生成一个示例JWT令牌
	claims := &jwt.Claims{
		UserID:   12345,
		Username: "demo-user",
		Role:     "admin",
		Email:    "demo@example.com",
	}

	token, err := jwt.Generate(claims)
	if err != nil {
		fmt.Printf("❌ JWT令牌生成失败: %v\n", err)
	} else {
		fmt.Printf("🎫 生成JWT令牌成功 (长度: %d)\n", len(token))

		// 验证令牌
		if err := jwt.Verify(token); err != nil {
			fmt.Printf("❌ JWT令牌验证失败: %v\n", err)
		} else {
			fmt.Printf("✅ JWT令牌验证成功\n")
		}
	}
}

func demoCrypto() {
	fmt.Println("\n🔐 Crypto 工具演示:")

	// 演示AES加密
	fmt.Println("\n🔒 AES加密演示:")
	plaintext := "Hello, Crypto World! 这是一个加密测试消息。"
	password := "my-secure-password"

	encrypted, err := crypto.QuickEncrypt(plaintext, password)
	if err != nil {
		fmt.Printf("❌ AES加密失败: %v\n", err)
		return
	}
	fmt.Printf("✅ 原文: %s\n", plaintext)
	fmt.Printf("✅ 密文: %s\n", encrypted[:50]+"...")

	decrypted, err := crypto.QuickDecrypt(encrypted, password)
	if err != nil {
		fmt.Printf("❌ AES解密失败: %v\n", err)
		return
	}
	fmt.Printf("✅ 解密: %s\n", decrypted)

	// 演示RSA加密
	fmt.Println("\n🔑 RSA加密演示:")
	privateKey, publicKey, err := crypto.GenerateKeyPair()
	if err != nil {
		fmt.Printf("❌ RSA密钥生成失败: %v\n", err)
		return
	}
	fmt.Printf("✅ RSA密钥对生成成功\n")

	rsaPlaintext := "Hello, RSA!"
	rsaEncrypted, err := crypto.RSAEncrypt(rsaPlaintext, publicKey)
	if err != nil {
		fmt.Printf("❌ RSA加密失败: %v\n", err)
		return
	}
	fmt.Printf("✅ RSA加密成功，密文长度: %d\n", len(rsaEncrypted))

	rsaDecrypted, err := crypto.RSADecrypt(rsaEncrypted, privateKey)
	if err != nil {
		fmt.Printf("❌ RSA解密失败: %v\n", err)
		return
	}
	fmt.Printf("✅ RSA解密成功: %s\n", rsaDecrypted)

	// 演示数字签名
	fmt.Println("\n📝 数字签名演示:")
	signData := "重要文档内容"
	signature, err := crypto.QuickSign(signData, privateKey)
	if err != nil {
		fmt.Printf("❌ 数字签名失败: %v\n", err)
		return
	}
	fmt.Printf("✅ 数字签名成功，签名长度: %d\n", len(signature))

	isValid, err := crypto.QuickVerify(signData, signature, publicKey)
	if err != nil {
		fmt.Printf("❌ 签名验证失败: %v\n", err)
		return
	}
	fmt.Printf("✅ 签名验证结果: %v\n", isValid)

	// 演示哈希算法
	fmt.Println("\n🔒 哈希算法演示:")
	hashData := "Hello, Hash World!"

	md5Hash := crypto.MD5(hashData)
	sha256Hash := crypto.SHA256(hashData)
	hmacHash := crypto.HMACSHA256(hashData, "secret-key")

	fmt.Printf("✅ 原文: %s\n", hashData)
	fmt.Printf("✅ MD5: %s\n", md5Hash)
	fmt.Printf("✅ SHA256: %s\n", sha256Hash)
	fmt.Printf("✅ HMAC-SHA256: %s\n", hmacHash)

	// 演示密码哈希
	fmt.Println("\n🛡️ 密码哈希演示:")
	userPassword := "user-password-123"

	hashedPassword, err := crypto.HashPassword(userPassword)
	if err != nil {
		fmt.Printf("❌ 密码哈希失败: %v\n", err)
		return
	}
	fmt.Printf("✅ 密码哈希成功，长度: %d\n", len(hashedPassword))

	isPasswordValid := crypto.CheckPassword(userPassword, hashedPassword)
	fmt.Printf("✅ 密码验证结果: %v\n", isPasswordValid)

	// 演示密码生成
	fmt.Println("\n🎲 密码生成演示:")
	randomPassword, err := crypto.GeneratePassword(12, false)
	if err != nil {
		fmt.Printf("❌ 密码生成失败: %v\n", err)
		return
	}

	strongPassword, err := crypto.GenerateStrongPassword(16)
	if err != nil {
		fmt.Printf("❌ 强密码生成失败: %v\n", err)
		return
	}

	strength := crypto.CheckPasswordStrength(strongPassword)
	fmt.Printf("✅ 随机密码: %s\n", randomPassword)
	fmt.Printf("✅ 强密码: %s (强度: %s)\n", strongPassword, strength.String())

	// 演示工具函数
	fmt.Println("\n🔧 工具函数演示:")
	randomBytes, err := crypto.GenerateRandomBytes(16)
	if err != nil {
		fmt.Printf("❌ 随机字节生成失败: %v\n", err)
		return
	}

	base64Encoded := crypto.Base64Encode(randomBytes)
	hexEncoded := crypto.HexEncode(randomBytes)

	fmt.Printf("✅ 随机字节 (Base64): %s\n", base64Encoded)
	fmt.Printf("✅ 随机字节 (Hex): %s\n", hexEncoded)

	fmt.Printf("\n🎯 Crypto工具演示完成！")
}
