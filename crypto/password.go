package crypto

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword 使用bcrypt哈希密码
func HashPassword(password string) (string, error) {
	return HashPasswordWithCost(password, globalConfig.DefaultBcryptCost)
}

// HashPasswordWithCost 使用指定成本哈希密码
func HashPasswordWithCost(password string, cost int) (string, error) {
	if err := ValidateBcryptCost(cost); err != nil {
		return "", err
	}
	
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", fmt.Errorf("密码哈希失败: %w", err)
	}
	
	return string(hashedPassword), nil
}

// CheckPassword 验证密码
func CheckPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// CheckPasswordWithError 验证密码（返回错误信息）
func CheckPasswordWithError(password, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return fmt.Errorf("密码验证失败: %w", err)
	}
	return nil
}

// GetPasswordHashCost 获取密码哈希的成本
func GetPasswordHashCost(hashedPassword string) (int, error) {
	cost, err := bcrypt.Cost([]byte(hashedPassword))
	if err != nil {
		return 0, fmt.Errorf("获取密码哈希成本失败: %w", err)
	}
	return cost, nil
}

// IsValidPasswordHash 检查是否为有效的bcrypt哈希
func IsValidPasswordHash(hashedPassword string) bool {
	_, err := bcrypt.Cost([]byte(hashedPassword))
	return err == nil
}

// PasswordStrength 密码强度评估
type PasswordStrength int

const (
	Weak PasswordStrength = iota
	Fair
	Good
	Strong
	VeryStrong
)

// String 返回密码强度描述
func (p PasswordStrength) String() string {
	switch p {
	case Weak:
		return "弱"
	case Fair:
		return "一般"
	case Good:
		return "良好"
	case Strong:
		return "强"
	case VeryStrong:
		return "很强"
	default:
		return "未知"
	}
}

// CheckPasswordStrength 检查密码强度
func CheckPasswordStrength(password string) PasswordStrength {
	score := 0
	length := len(password)
	
	// 长度评分
	if length >= 8 {
		score++
	}
	if length >= 12 {
		score++
	}
	if length >= 16 {
		score++
	}
	
	// 字符类型评分
	hasLower := false
	hasUpper := false
	hasDigit := false
	hasSpecial := false
	
	for _, char := range password {
		switch {
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= '0' && char <= '9':
			hasDigit = true
		default:
			hasSpecial = true
		}
	}
	
	if hasLower {
		score++
	}
	if hasUpper {
		score++
	}
	if hasDigit {
		score++
	}
	if hasSpecial {
		score++
	}
	
	// 根据评分返回强度
	switch {
	case score <= 2:
		return Weak
	case score <= 4:
		return Fair
	case score <= 5:
		return Good
	case score <= 6:
		return Strong
	default:
		return VeryStrong
	}
}

// GeneratePassword 生成随机密码
func GeneratePassword(length int, includeSymbols bool) (string, error) {
	if length < 4 {
		return "", fmt.Errorf("密码长度至少为4位")
	}
	
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	if includeSymbols {
		chars += "!@#$%^&*()_+-=[]{}|;:,.<>?"
	}
	
	return GenerateRandomStringFromChars(length, chars)
}

// GenerateStrongPassword 生成强密码
func GenerateStrongPassword(length int) (string, error) {
	if length < 8 {
		return "", fmt.Errorf("强密码长度至少为8位")
	}
	
	// 确保包含各种字符类型
	lower := "abcdefghijklmnopqrstuvwxyz"
	upper := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits := "0123456789"
	symbols := "!@#$%^&*()_+-=[]{}|;:,.<>?"
	
	password := ""
	
	// 至少包含一个小写字母
	char, err := GenerateRandomStringFromChars(1, lower)
	if err != nil {
		return "", err
	}
	password += char
	
	// 至少包含一个大写字母
	char, err = GenerateRandomStringFromChars(1, upper)
	if err != nil {
		return "", err
	}
	password += char
	
	// 至少包含一个数字
	char, err = GenerateRandomStringFromChars(1, digits)
	if err != nil {
		return "", err
	}
	password += char
	
	// 至少包含一个符号
	char, err = GenerateRandomStringFromChars(1, symbols)
	if err != nil {
		return "", err
	}
	password += char
	
	// 填充剩余长度
	if length > 4 {
		allChars := lower + upper + digits + symbols
		remaining, err := GenerateRandomStringFromChars(length-4, allChars)
		if err != nil {
			return "", err
		}
		password += remaining
	}
	
	// 打乱密码字符顺序
	return shuffleString(password), nil
}

// shuffleString 打乱字符串顺序
func shuffleString(s string) string {
	runes := []rune(s)
	
	// 简单的Fisher-Yates洗牌算法
	for i := len(runes) - 1; i > 0; i-- {
		// 生成0到i之间的随机数
		randomBytes, err := GenerateRandomBytes(1)
		if err != nil {
			return s // 如果出错，返回原字符串
		}
		j := int(randomBytes[0]) % (i + 1)
		runes[i], runes[j] = runes[j], runes[i]
	}
	
	return string(runes)
}

// ValidatePasswordPolicy 验证密码策略
type PasswordPolicy struct {
	MinLength      int  // 最小长度
	RequireLower   bool // 需要小写字母
	RequireUpper   bool // 需要大写字母
	RequireDigit   bool // 需要数字
	RequireSpecial bool // 需要特殊字符
}

// DefaultPasswordPolicy 返回默认密码策略
func DefaultPasswordPolicy() *PasswordPolicy {
	return &PasswordPolicy{
		MinLength:      8,
		RequireLower:   true,
		RequireUpper:   true,
		RequireDigit:   true,
		RequireSpecial: false,
	}
}

// ValidatePassword 根据策略验证密码
func ValidatePassword(password string, policy *PasswordPolicy) error {
	if policy == nil {
		policy = DefaultPasswordPolicy()
	}
	
	// 检查长度
	if len(password) < policy.MinLength {
		return fmt.Errorf("密码长度至少为%d位", policy.MinLength)
	}
	
	hasLower := false
	hasUpper := false
	hasDigit := false
	hasSpecial := false
	
	for _, char := range password {
		switch {
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= '0' && char <= '9':
			hasDigit = true
		default:
			hasSpecial = true
		}
	}
	
	if policy.RequireLower && !hasLower {
		return fmt.Errorf("密码必须包含小写字母")
	}
	
	if policy.RequireUpper && !hasUpper {
		return fmt.Errorf("密码必须包含大写字母")
	}
	
	if policy.RequireDigit && !hasDigit {
		return fmt.Errorf("密码必须包含数字")
	}
	
	if policy.RequireSpecial && !hasSpecial {
		return fmt.Errorf("密码必须包含特殊字符")
	}
	
	return nil
}
