package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

// Config JWT配置
type Config struct {
	Secret     string        // 签名密钥
	Issuer     string        // 签发者
	Expiration time.Duration // 过期时间，0表示永不过期
}

// Claims JWT载荷
type Claims struct {
	UserID    interface{}            `json:"user_id,omitempty"`  // 用户ID
	Username  string                 `json:"username,omitempty"` // 用户名
	Role      string                 `json:"role,omitempty"`     // 角色
	Email     string                 `json:"email,omitempty"`    // 邮箱
	Issuer    string                 `json:"iss,omitempty"`      // 签发者
	Subject   string                 `json:"sub,omitempty"`      // 主题
	Audience  string                 `json:"aud,omitempty"`      // 受众
	IssuedAt  int64                  `json:"iat,omitempty"`      // 签发时间
	ExpireAt  int64                  `json:"exp,omitempty"`      // 过期时间
	NotBefore int64                  `json:"nbf,omitempty"`      // 生效时间
	Custom    map[string]interface{} `json:"-"`                  // 自定义字段
}

// Header JWT头部
type Header struct {
	Type      string `json:"typ"`
	Algorithm string `json:"alg"`
}

var (
	// 全局配置
	globalConfig = &Config{
		Secret:     "helwd-jwt-secret",
		Issuer:     "helwd-app",
		Expiration: 24 * time.Hour, // 默认24小时
	}
)

// Init 初始化JWT全局配置
func Init(secret, issuer string, expiration time.Duration) {
	globalConfig.Secret = secret
	globalConfig.Issuer = issuer
	globalConfig.Expiration = expiration
}

// InitDefault 使用默认配置初始化
func InitDefault() {
	Init("helwd-jwt-secret", "helwd-app", 24*time.Hour)
}

// SetSecret 设置全局密钥
func SetSecret(secret string) {
	globalConfig.Secret = secret
}

// SetIssuer 设置全局签发者
func SetIssuer(issuer string) {
	globalConfig.Issuer = issuer
}

// SetExpiration 设置全局过期时间
func SetExpiration(expiration time.Duration) {
	globalConfig.Expiration = expiration
}

// Generate 生成JWT令牌
func Generate(claims *Claims) (string, error) {
	return GenerateWithConfig(claims, nil)
}

// GenerateWithConfig 使用自定义配置生成JWT令牌
func GenerateWithConfig(claims *Claims, config *Config) (string, error) {
	// 确定使用的配置
	cfg := globalConfig
	if config != nil {
		cfg = config
	}

	// 设置默认值
	now := time.Now()
	if claims.IssuedAt == 0 {
		claims.IssuedAt = now.Unix()
	}
	if claims.Issuer == "" {
		claims.Issuer = cfg.Issuer
	}
	// 只有在ExpireAt为0且配置了过期时间时才自动设置
	if claims.ExpireAt == 0 && cfg.Expiration > 0 {
		claims.ExpireAt = now.Add(cfg.Expiration).Unix()
	}

	// 创建头部
	header := &Header{
		Type:      "JWT",
		Algorithm: "HS256",
	}

	// 编码头部
	headerBytes, err := json.Marshal(header)
	if err != nil {
		return "", fmt.Errorf("编码头部失败: %w", err)
	}
	headerEncoded := base64.RawURLEncoding.EncodeToString(headerBytes)

	// 合并自定义字段到Claims
	claimsMap := make(map[string]interface{})
	claimsBytes, _ := json.Marshal(claims)
	json.Unmarshal(claimsBytes, &claimsMap)

	// 添加自定义字段
	if claims.Custom != nil {
		for k, v := range claims.Custom {
			claimsMap[k] = v
		}
	}

	// 编码载荷
	payloadBytes, err := json.Marshal(claimsMap)
	if err != nil {
		return "", fmt.Errorf("编码载荷失败: %w", err)
	}
	payloadEncoded := base64.RawURLEncoding.EncodeToString(payloadBytes)

	// 创建签名
	message := headerEncoded + "." + payloadEncoded
	signature := createSignature(message, cfg.Secret)

	// 组合最终令牌
	token := message + "." + signature
	return token, nil
}

// Parse 解析JWT令牌
func Parse(token string) (*Claims, error) {
	return ParseWithConfig(token, nil)
}

// ParseWithConfig 使用自定义配置解析JWT令牌
func ParseWithConfig(token string, config *Config) (*Claims, error) {
	// 确定使用的配置
	cfg := globalConfig
	if config != nil {
		cfg = config
	}

	// 分割令牌
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errors.New("无效的JWT格式")
	}

	headerEncoded, payloadEncoded, signatureEncoded := parts[0], parts[1], parts[2]

	// 验证签名
	message := headerEncoded + "." + payloadEncoded
	expectedSignature := createSignature(message, cfg.Secret)
	if signatureEncoded != expectedSignature {
		return nil, errors.New("JWT签名验证失败")
	}

	// 解码载荷
	payloadBytes, err := base64.RawURLEncoding.DecodeString(payloadEncoded)
	if err != nil {
		return nil, fmt.Errorf("解码载荷失败: %w", err)
	}

	// 解析Claims
	var claimsMap map[string]interface{}
	if err := json.Unmarshal(payloadBytes, &claimsMap); err != nil {
		return nil, fmt.Errorf("解析载荷失败: %w", err)
	}

	claims := &Claims{
		Custom: make(map[string]interface{}),
	}

	// 提取标准字段
	if v, ok := claimsMap["user_id"]; ok {
		// 处理数字类型转换
		if f, ok := v.(float64); ok {
			claims.UserID = int(f)
		} else {
			claims.UserID = v
		}
	}
	if v, ok := claimsMap["username"].(string); ok {
		claims.Username = v
	}
	if v, ok := claimsMap["role"].(string); ok {
		claims.Role = v
	}
	if v, ok := claimsMap["email"].(string); ok {
		claims.Email = v
	}
	if v, ok := claimsMap["iss"].(string); ok {
		claims.Issuer = v
	}
	if v, ok := claimsMap["sub"].(string); ok {
		claims.Subject = v
	}
	if v, ok := claimsMap["aud"].(string); ok {
		claims.Audience = v
	}
	if v, ok := claimsMap["iat"].(float64); ok {
		claims.IssuedAt = int64(v)
	}
	if v, ok := claimsMap["exp"].(float64); ok {
		claims.ExpireAt = int64(v)
	}
	if v, ok := claimsMap["nbf"].(float64); ok {
		claims.NotBefore = int64(v)
	}

	// 提取自定义字段
	standardFields := map[string]bool{
		"user_id": true, "username": true, "role": true, "email": true,
		"iss": true, "sub": true, "aud": true, "iat": true, "exp": true, "nbf": true,
	}
	for k, v := range claimsMap {
		if !standardFields[k] {
			claims.Custom[k] = v
		}
	}

	return claims, nil
}

// Verify 验证JWT令牌有效性
func Verify(token string) error {
	return VerifyWithConfig(token, nil)
}

// VerifyWithConfig 使用自定义配置验证JWT令牌
func VerifyWithConfig(token string, config *Config) error {
	claims, err := ParseWithConfig(token, config)
	if err != nil {
		return err
	}

	now := time.Now().Unix()

	// 检查是否已过期
	if claims.ExpireAt > 0 && now > claims.ExpireAt {
		return errors.New("JWT令牌已过期")
	}

	// 检查是否还未生效
	if claims.NotBefore > 0 && now < claims.NotBefore {
		return errors.New("JWT令牌还未生效")
	}

	return nil
}

// createSignature 创建HMAC-SHA256签名
func createSignature(message, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(message))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}

// IsExpired 检查令牌是否过期
func IsExpired(token string) bool {
	claims, err := Parse(token)
	if err != nil {
		return true
	}

	if claims.ExpireAt == 0 {
		return false // 永不过期
	}

	return time.Now().Unix() > claims.ExpireAt
}

// GetClaims 获取令牌中的Claims（不验证签名）
func GetClaims(token string) (*Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errors.New("无效的JWT格式")
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("解码载荷失败: %w", err)
	}

	var claimsMap map[string]interface{}
	if err := json.Unmarshal(payloadBytes, &claimsMap); err != nil {
		return nil, fmt.Errorf("解析载荷失败: %w", err)
	}

	claims := &Claims{Custom: make(map[string]interface{})}

	// 简化版本，只提取基本字段
	if v, ok := claimsMap["user_id"]; ok {
		// 处理数字类型转换
		if f, ok := v.(float64); ok {
			claims.UserID = int(f)
		} else {
			claims.UserID = v
		}
	}
	if v, ok := claimsMap["username"].(string); ok {
		claims.Username = v
	}
	if v, ok := claimsMap["exp"].(float64); ok {
		claims.ExpireAt = int64(v)
	}

	return claims, nil
}

// Refresh 刷新令牌（重新生成过期时间）
func Refresh(token string) (string, error) {
	return RefreshWithConfig(token, nil)
}

// RefreshWithConfig 使用自定义配置刷新令牌
func RefreshWithConfig(token string, config *Config) (string, error) {
	claims, err := ParseWithConfig(token, config)
	if err != nil {
		return "", err
	}

	// 重置时间字段，确保令牌会发生变化
	now := time.Now()
	// 强制更新IssuedAt，确保与原令牌不同
	claims.IssuedAt = now.Unix()

	cfg := globalConfig
	if config != nil {
		cfg = config
	}

	// 重新设置过期时间
	if cfg.Expiration > 0 {
		claims.ExpireAt = now.Add(cfg.Expiration).Unix()
	}
	// 注意：如果原来有过期时间但配置中没有，保持原有的过期时间逻辑

	return GenerateWithConfig(claims, config)
}
