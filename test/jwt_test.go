package test

import (
	"testing"
	"time"

	jwt "github.com/fastgox/utils/jwt"
)

func TestJWT(t *testing.T) {
	// 初始化JWT配置
	jwt.Init("test-secret-key", "test-app", 1*time.Hour)

	// 测试基本令牌生成和解析
	t.Run("基本功能测试", func(t *testing.T) {
		claims := &jwt.Claims{
			UserID:   12345,
			Username: "helwd",
			Role:     "admin",
			Email:    "helwd@example.com",
		}

		// 生成令牌
		token, err := jwt.Generate(claims)
		if err != nil {
			t.Fatalf("生成令牌失败: %v", err)
		}

		if token == "" {
			t.Error("生成的令牌为空")
		}

		// 解析令牌
		parsedClaims, err := jwt.Parse(token)
		if err != nil {
			t.Fatalf("解析令牌失败: %v", err)
		}

		// 验证解析结果
		if parsedClaims.UserID != claims.UserID {
			t.Errorf("UserID不匹配，期望: %v, 实际: %v", claims.UserID, parsedClaims.UserID)
		}

		if parsedClaims.Username != claims.Username {
			t.Errorf("Username不匹配，期望: %s, 实际: %s", claims.Username, parsedClaims.Username)
		}

		if parsedClaims.Role != claims.Role {
			t.Errorf("Role不匹配，期望: %s, 实际: %s", claims.Role, parsedClaims.Role)
		}

		if parsedClaims.Email != claims.Email {
			t.Errorf("Email不匹配，期望: %s, 实际: %s", claims.Email, parsedClaims.Email)
		}
	})

	// 测试自定义字段
	t.Run("自定义字段测试", func(t *testing.T) {
		claims := &jwt.Claims{
			UserID:   67890,
			Username: "testuser",
			Custom: map[string]interface{}{
				"department":  "技术部",
				"level":       5,
				"permissions": []interface{}{"read", "write", "admin"},
			},
		}

		token, err := jwt.Generate(claims)
		if err != nil {
			t.Fatalf("生成令牌失败: %v", err)
		}

		parsedClaims, err := jwt.Parse(token)
		if err != nil {
			t.Fatalf("解析令牌失败: %v", err)
		}

		// 验证自定义字段
		if parsedClaims.Custom["department"] != "技术部" {
			t.Errorf("自定义字段department不匹配")
		}

		if parsedClaims.Custom["level"] != float64(5) { // JSON解析数字为float64
			t.Errorf("自定义字段level不匹配")
		}
	})

	// 测试令牌验证
	t.Run("令牌验证测试", func(t *testing.T) {
		claims := &jwt.Claims{
			UserID:   11111,
			Username: "verifyuser",
		}

		token, err := jwt.Generate(claims)
		if err != nil {
			t.Fatalf("生成令牌失败: %v", err)
		}

		// 验证有效令牌
		if err := jwt.Verify(token); err != nil {
			t.Errorf("有效令牌验证失败: %v", err)
		}

		// 测试无效令牌
		invalidToken := token + "invalid"
		if err := jwt.Verify(invalidToken); err == nil {
			t.Error("无效令牌应该验证失败")
		}
	})

	// 测试自定义配置
	t.Run("自定义配置测试", func(t *testing.T) {
		customConfig := &jwt.Config{
			Secret:     "custom-secret",
			Issuer:     "custom-app",
			Expiration: 30 * time.Minute,
		}

		claims := &jwt.Claims{
			UserID:   22222,
			Username: "customuser",
		}

		// 使用自定义配置生成令牌
		token, err := jwt.GenerateWithConfig(claims, customConfig)
		if err != nil {
			t.Fatalf("使用自定义配置生成令牌失败: %v", err)
		}

		// 使用自定义配置解析令牌
		parsedClaims, err := jwt.ParseWithConfig(token, customConfig)
		if err != nil {
			t.Fatalf("使用自定义配置解析令牌失败: %v", err)
		}

		if parsedClaims.Issuer != "custom-app" {
			t.Errorf("Issuer不匹配，期望: custom-app, 实际: %s", parsedClaims.Issuer)
		}

		// 使用错误配置应该解析失败
		wrongConfig := &jwt.Config{
			Secret:     "wrong-secret",
			Issuer:     "wrong-app",
			Expiration: 1 * time.Hour,
		}

		_, err = jwt.ParseWithConfig(token, wrongConfig)
		if err == nil {
			t.Error("使用错误密钥应该解析失败")
		}
	})

	// 测试过期时间
	t.Run("过期时间测试", func(t *testing.T) {
		// 创建一个已经过期的令牌
		claims := &jwt.Claims{
			UserID:   33333,
			Username: "expireuser",
			ExpireAt: time.Now().Add(-1 * time.Hour).Unix(), // 1小时前过期
		}

		shortConfig := &jwt.Config{
			Secret:     "short-secret",
			Issuer:     "short-app",
			Expiration: 0, // 不自动设置过期时间，使用Claims中的
		}

		token, err := jwt.GenerateWithConfig(claims, shortConfig)
		if err != nil {
			t.Fatalf("生成过期令牌失败: %v", err)
		}

		// 验证过期令牌
		if err := jwt.VerifyWithConfig(token, shortConfig); err == nil {
			t.Error("过期令牌应该验证失败")
		}

		// 测试IsExpired函数
		if !jwt.IsExpired(token) {
			t.Error("IsExpired应该返回true")
		}
	})

	// 测试令牌刷新
	t.Run("令牌刷新测试", func(t *testing.T) {
		claims := &jwt.Claims{
			UserID:   44444,
			Username: "refreshuser",
		}

		originalToken, err := jwt.Generate(claims)
		if err != nil {
			t.Fatalf("生成原始令牌失败: %v", err)
		}

		// 等待一小段时间确保时间戳不同
		time.Sleep(1 * time.Second)

		// 刷新令牌
		newToken, err := jwt.Refresh(originalToken)
		if err != nil {
			t.Fatalf("刷新令牌失败: %v", err)
		}

		if newToken == originalToken {
			t.Logf("原令牌: %s", originalToken)
			t.Logf("新令牌: %s", newToken)
			t.Error("刷新后的令牌应该与原令牌不同")
		}

		// 验证新令牌
		newClaims, err := jwt.Parse(newToken)
		if err != nil {
			t.Fatalf("解析新令牌失败: %v", err)
		}

		if newClaims.UserID != claims.UserID {
			t.Errorf("刷新后UserID不匹配")
		}

		if newClaims.Username != claims.Username {
			t.Errorf("刷新后Username不匹配")
		}
	})

	// 测试GetClaims函数
	t.Run("GetClaims测试", func(t *testing.T) {
		claims := &jwt.Claims{
			UserID:   55555,
			Username: "getclaimsuser",
		}

		token, err := jwt.Generate(claims)
		if err != nil {
			t.Fatalf("生成令牌失败: %v", err)
		}

		// 获取Claims（不验证签名）
		getClaims, err := jwt.GetClaims(token)
		if err != nil {
			t.Fatalf("获取Claims失败: %v", err)
		}

		if getClaims.UserID != claims.UserID {
			t.Errorf("GetClaims UserID不匹配")
		}

		if getClaims.Username != claims.Username {
			t.Errorf("GetClaims Username不匹配")
		}
	})

	// 测试配置函数
	t.Run("配置函数测试", func(t *testing.T) {
		// 测试设置函数
		jwt.SetSecret("new-secret")
		jwt.SetIssuer("new-issuer")
		jwt.SetExpiration(2 * time.Hour)

		claims := &jwt.Claims{
			UserID:   66666,
			Username: "configuser",
		}

		token, err := jwt.Generate(claims)
		if err != nil {
			t.Fatalf("生成令牌失败: %v", err)
		}

		parsedClaims, err := jwt.Parse(token)
		if err != nil {
			t.Fatalf("解析令牌失败: %v", err)
		}

		if parsedClaims.Issuer != "new-issuer" {
			t.Errorf("新Issuer不匹配，期望: new-issuer, 实际: %s", parsedClaims.Issuer)
		}

		// 重置为测试配置
		jwt.Init("test-secret-key", "test-app", 1*time.Hour)
	})
}
