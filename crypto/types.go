package crypto

import (
	"crypto/rsa"
	"errors"
)

// 常用的密钥长度
const (
	AES128KeySize = 16 // AES-128
	AES192KeySize = 24 // AES-192
	AES256KeySize = 32 // AES-256

	RSA1024KeySize = 1024 // RSA-1024 (不推荐)
	RSA2048KeySize = 2048 // RSA-2048 (推荐)
	RSA3072KeySize = 3072 // RSA-3072
	RSA4096KeySize = 4096 // RSA-4096

	DefaultBcryptCost = 12 // bcrypt默认成本
)

// 常见错误
var (
	ErrInvalidKeySize      = errors.New("无效的密钥长度")
	ErrInvalidKey          = errors.New("无效的密钥")
	ErrInvalidCiphertext   = errors.New("无效的密文")
	ErrInvalidPlaintext    = errors.New("无效的明文")
	ErrInvalidSignature    = errors.New("无效的签名")
	ErrKeyGenerationFailed = errors.New("密钥生成失败")
	ErrEncryptionFailed    = errors.New("加密失败")
	ErrDecryptionFailed    = errors.New("解密失败")
	ErrSigningFailed       = errors.New("签名失败")
	ErrVerificationFailed  = errors.New("验证失败")
)

// Config 加密工具配置
type Config struct {
	DefaultAESKey     string // 默认AES密钥
	DefaultBcryptCost int    // 默认bcrypt成本
	DefaultRSAKeySize int    // 默认RSA密钥长度
}

// RSAKeyPair RSA密钥对
type RSAKeyPair struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
	PrivatePEM string // PEM格式私钥
	PublicPEM  string // PEM格式公钥
}

// HashAlgorithm 哈希算法类型
type HashAlgorithm int

const (
	HashMD5 HashAlgorithm = iota
	HashSHA1
	HashSHA224
	HashSHA256
	HashSHA384
	HashSHA512
)

// String 返回哈希算法名称
func (h HashAlgorithm) String() string {
	switch h {
	case HashMD5:
		return "MD5"
	case HashSHA1:
		return "SHA1"
	case HashSHA224:
		return "SHA224"
	case HashSHA256:
		return "SHA256"
	case HashSHA384:
		return "SHA384"
	case HashSHA512:
		return "SHA512"
	default:
		return "Unknown"
	}
}

// EncryptionMode 加密模式
type EncryptionMode int

const (
	CBC EncryptionMode = iota // CBC模式
	GCM                       // GCM模式
	CFB                       // CFB模式
	OFB                       // OFB模式
)

// String 返回加密模式名称
func (e EncryptionMode) String() string {
	switch e {
	case CBC:
		return "CBC"
	case GCM:
		return "GCM"
	case CFB:
		return "CFB"
	case OFB:
		return "OFB"
	default:
		return "Unknown"
	}
}

// SignatureAlgorithm 签名算法类型
type SignatureAlgorithm int

const (
	RSA_PKCS1v15 SignatureAlgorithm = iota // RSA PKCS#1 v1.5
	RSA_PSS                                // RSA PSS
	ECDSA_P256                             // ECDSA P-256
	ECDSA_P384                             // ECDSA P-384
	ECDSA_P521                             // ECDSA P-521
)

// String 返回签名算法名称
func (s SignatureAlgorithm) String() string {
	switch s {
	case RSA_PKCS1v15:
		return "RSA-PKCS1v15"
	case RSA_PSS:
		return "RSA-PSS"
	case ECDSA_P256:
		return "ECDSA-P256"
	case ECDSA_P384:
		return "ECDSA-P384"
	case ECDSA_P521:
		return "ECDSA-P521"
	default:
		return "Unknown"
	}
}

// FileEncryptionOptions 文件加密选项
type FileEncryptionOptions struct {
	Algorithm     string         // 加密算法 (AES)
	Mode          EncryptionMode // 加密模式
	KeySize       int            // 密钥长度
	BufferSize    int            // 缓冲区大小
	Compress      bool           // 是否压缩
	IncludeHeader bool           // 是否包含文件头
}

// DefaultFileEncryptionOptions 返回默认文件加密选项
func DefaultFileEncryptionOptions() *FileEncryptionOptions {
	return &FileEncryptionOptions{
		Algorithm:     "AES",
		Mode:          GCM,
		KeySize:       AES256KeySize,
		BufferSize:    64 * 1024, // 64KB
		Compress:      false,
		IncludeHeader: true,
	}
}

// PasswordHashOptions 密码哈希选项
type PasswordHashOptions struct {
	Algorithm string // 哈希算法 (bcrypt, scrypt, argon2)
	Cost      int    // 成本参数
	SaltSize  int    // 盐长度
}

// DefaultPasswordHashOptions 返回默认密码哈希选项
func DefaultPasswordHashOptions() *PasswordHashOptions {
	return &PasswordHashOptions{
		Algorithm: "bcrypt",
		Cost:      DefaultBcryptCost,
		SaltSize:  16,
	}
}

// RandomOptions 随机数生成选项
type RandomOptions struct {
	Length      int    // 长度
	UseNumbers  bool   // 使用数字
	UseLetters  bool   // 使用字母
	UseSymbols  bool   // 使用符号
	CustomChars string // 自定义字符集
}

// DefaultRandomOptions 返回默认随机数选项
func DefaultRandomOptions() *RandomOptions {
	return &RandomOptions{
		Length:     32,
		UseNumbers: true,
		UseLetters: true,
		UseSymbols: false,
	}
}

var (
	// 全局配置
	globalConfig = &Config{
		DefaultAESKey:     "",
		DefaultBcryptCost: DefaultBcryptCost,
		DefaultRSAKeySize: RSA2048KeySize,
	}
)

// GetGlobalConfig 获取全局配置
func GetGlobalConfig() *Config {
	return globalConfig
}

// SetGlobalConfig 设置全局配置
func SetGlobalConfig(config *Config) {
	if config != nil {
		globalConfig = config
	}
}

// SetDefaultAESKey 设置默认AES密钥
func SetDefaultAESKey(key string) {
	globalConfig.DefaultAESKey = key
}

// SetDefaultBcryptCost 设置默认bcrypt成本
func SetDefaultBcryptCost(cost int) {
	if cost >= 4 && cost <= 31 {
		globalConfig.DefaultBcryptCost = cost
	}
}

// SetDefaultRSAKeySize 设置默认RSA密钥长度
func SetDefaultRSAKeySize(keySize int) {
	if keySize >= RSA1024KeySize {
		globalConfig.DefaultRSAKeySize = keySize
	}
}

// ValidateAESKeySize 验证AES密钥长度
func ValidateAESKeySize(keySize int) error {
	if keySize != AES128KeySize && keySize != AES192KeySize && keySize != AES256KeySize {
		return ErrInvalidKeySize
	}
	return nil
}

// ValidateRSAKeySize 验证RSA密钥长度
func ValidateRSAKeySize(keySize int) error {
	if keySize < RSA1024KeySize || keySize%8 != 0 {
		return ErrInvalidKeySize
	}
	return nil
}

// ValidateBcryptCost 验证bcrypt成本
func ValidateBcryptCost(cost int) error {
	if cost < 4 || cost > 31 {
		return errors.New("bcrypt成本必须在4-31之间")
	}
	return nil
}
