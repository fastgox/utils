package crypto

import (
	"fmt"
	"io"
	"os"
)

// Init 初始化加密工具
func Init() {
	// 初始化全局配置
	globalConfig = &Config{
		DefaultAESKey:     "",
		DefaultBcryptCost: DefaultBcryptCost,
		DefaultRSAKeySize: RSA2048KeySize,
	}
}

// InitWithConfig 使用自定义配置初始化
func InitWithConfig(config *Config) {
	if config != nil {
		globalConfig = config
	} else {
		Init()
	}
}

// EncryptFile 加密文件
func EncryptFile(inputFile, outputFile, password string) error {
	return EncryptFileWithOptions(inputFile, outputFile, password, DefaultFileEncryptionOptions())
}

// EncryptFileWithOptions 使用选项加密文件
func EncryptFileWithOptions(inputFile, outputFile, password string, options *FileEncryptionOptions) error {
	if options == nil {
		options = DefaultFileEncryptionOptions()
	}
	
	// 读取输入文件
	inputData, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("读取输入文件失败: %w", err)
	}
	
	// 生成密钥
	key, err := AESKeyFromPassword(password, "file-salt", options.KeySize)
	if err != nil {
		return fmt.Errorf("生成密钥失败: %w", err)
	}
	
	// 加密数据
	encryptedData, err := AESEncryptBytes(inputData, key)
	if err != nil {
		return fmt.Errorf("加密数据失败: %w", err)
	}
	
	// 写入输出文件
	err = os.WriteFile(outputFile, encryptedData, 0644)
	if err != nil {
		return fmt.Errorf("写入输出文件失败: %w", err)
	}
	
	return nil
}

// DecryptFile 解密文件
func DecryptFile(inputFile, outputFile, password string) error {
	return DecryptFileWithOptions(inputFile, outputFile, password, DefaultFileEncryptionOptions())
}

// DecryptFileWithOptions 使用选项解密文件
func DecryptFileWithOptions(inputFile, outputFile, password string, options *FileEncryptionOptions) error {
	if options == nil {
		options = DefaultFileEncryptionOptions()
	}
	
	// 读取输入文件
	encryptedData, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("读取输入文件失败: %w", err)
	}
	
	// 生成密钥
	key, err := AESKeyFromPassword(password, "file-salt", options.KeySize)
	if err != nil {
		return fmt.Errorf("生成密钥失败: %w", err)
	}
	
	// 解密数据
	decryptedData, err := AESDecryptBytes(encryptedData, key)
	if err != nil {
		return fmt.Errorf("解密数据失败: %w", err)
	}
	
	// 写入输出文件
	err = os.WriteFile(outputFile, decryptedData, 0644)
	if err != nil {
		return fmt.Errorf("写入输出文件失败: %w", err)
	}
	
	return nil
}

// EncryptStream 加密数据流
func EncryptStream(reader io.Reader, writer io.Writer, password string) error {
	// 生成密钥
	key, err := AESKeyFromPassword(password, "stream-salt", AES256KeySize)
	if err != nil {
		return fmt.Errorf("生成密钥失败: %w", err)
	}
	
	// 读取所有数据（简化处理）
	data, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("读取数据失败: %w", err)
	}
	
	// 加密数据
	encryptedData, err := AESEncryptBytes(data, key)
	if err != nil {
		return fmt.Errorf("加密数据失败: %w", err)
	}
	
	// 写入加密数据
	_, err = writer.Write(encryptedData)
	if err != nil {
		return fmt.Errorf("写入数据失败: %w", err)
	}
	
	return nil
}

// DecryptStream 解密数据流
func DecryptStream(reader io.Reader, writer io.Writer, password string) error {
	// 生成密钥
	key, err := AESKeyFromPassword(password, "stream-salt", AES256KeySize)
	if err != nil {
		return fmt.Errorf("生成密钥失败: %w", err)
	}
	
	// 读取所有数据（简化处理）
	encryptedData, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("读取数据失败: %w", err)
	}
	
	// 解密数据
	decryptedData, err := AESDecryptBytes(encryptedData, key)
	if err != nil {
		return fmt.Errorf("解密数据失败: %w", err)
	}
	
	// 写入解密数据
	_, err = writer.Write(decryptedData)
	if err != nil {
		return fmt.Errorf("写入数据失败: %w", err)
	}
	
	return nil
}

// QuickEncrypt 快速加密（使用默认设置）
func QuickEncrypt(data, password string) (string, error) {
	return AESEncryptWithPassword(data, password)
}

// QuickDecrypt 快速解密（使用默认设置）
func QuickDecrypt(encryptedData, password string) (string, error) {
	return AESDecryptWithPassword(encryptedData, password)
}

// QuickHash 快速哈希（使用SHA256）
func QuickHash(data string) string {
	return SHA256(data)
}

// QuickHMAC 快速HMAC（使用SHA256）
func QuickHMAC(data, key string) string {
	return HMACSHA256(data, key)
}

// QuickSign 快速签名（使用RSA）
func QuickSign(data, privateKey string) (string, error) {
	return RSASign(data, privateKey)
}

// QuickVerify 快速验证签名（使用RSA）
func QuickVerify(data, signature, publicKey string) (bool, error) {
	return RSAVerify(data, signature, publicKey)
}

// GenerateKeyPair 生成密钥对（默认RSA-2048）
func GenerateKeyPair() (privateKey, publicKey string, err error) {
	return GenerateRSAKeyPair(globalConfig.DefaultRSAKeySize)
}

// GenerateSecretKey 生成密钥（默认AES-256）
func GenerateSecretKey() ([]byte, error) {
	return GenerateAESKey(AES256KeySize)
}

// GenerateSecretKeyString 生成密钥字符串（默认AES-256）
func GenerateSecretKeyString() (string, error) {
	key, err := GenerateSecretKey()
	if err != nil {
		return "", err
	}
	return Base64Encode(key), nil
}

// Benchmark 性能测试
func Benchmark() {
	fmt.Println("=== Crypto 性能测试 ===")
	
	// AES加密性能测试
	fmt.Println("\n🔐 AES加密性能测试:")
	testData := "Hello, World! This is a test message for encryption benchmark."
	testKey := "my-secret-key-32-bytes-long!!"
	
	// 测试AES加密
	encrypted, err := AESEncrypt(testData, testKey)
	if err != nil {
		fmt.Printf("❌ AES加密失败: %v\n", err)
	} else {
		fmt.Printf("✅ AES加密成功，密文长度: %d\n", len(encrypted))
	}
	
	// 测试AES解密
	decrypted, err := AESDecrypt(encrypted, testKey)
	if err != nil {
		fmt.Printf("❌ AES解密失败: %v\n", err)
	} else if decrypted == testData {
		fmt.Printf("✅ AES解密成功，数据一致\n")
	} else {
		fmt.Printf("❌ AES解密数据不一致\n")
	}
	
	// RSA加密性能测试
	fmt.Println("\n🔑 RSA加密性能测试:")
	privateKey, publicKey, err := GenerateRSAKeyPair(2048)
	if err != nil {
		fmt.Printf("❌ RSA密钥生成失败: %v\n", err)
		return
	}
	fmt.Printf("✅ RSA密钥生成成功\n")
	
	// 测试RSA加密
	rsaTestData := "Hello, RSA!"
	rsaEncrypted, err := RSAEncrypt(rsaTestData, publicKey)
	if err != nil {
		fmt.Printf("❌ RSA加密失败: %v\n", err)
	} else {
		fmt.Printf("✅ RSA加密成功，密文长度: %d\n", len(rsaEncrypted))
	}
	
	// 测试RSA解密
	rsaDecrypted, err := RSADecrypt(rsaEncrypted, privateKey)
	if err != nil {
		fmt.Printf("❌ RSA解密失败: %v\n", err)
	} else if rsaDecrypted == rsaTestData {
		fmt.Printf("✅ RSA解密成功，数据一致\n")
	} else {
		fmt.Printf("❌ RSA解密数据不一致\n")
	}
	
	// 哈希性能测试
	fmt.Println("\n🔒 哈希性能测试:")
	hashTestData := "Hello, Hash!"
	
	md5Hash := MD5(hashTestData)
	sha256Hash := SHA256(hashTestData)
	sha512Hash := SHA512(hashTestData)
	
	fmt.Printf("✅ MD5: %s\n", md5Hash)
	fmt.Printf("✅ SHA256: %s\n", sha256Hash)
	fmt.Printf("✅ SHA512: %s\n", sha512Hash)
	
	// 密码哈希性能测试
	fmt.Println("\n🛡️ 密码哈希性能测试:")
	password := "test-password"
	
	hashedPassword, err := HashPassword(password)
	if err != nil {
		fmt.Printf("❌ 密码哈希失败: %v\n", err)
	} else {
		fmt.Printf("✅ 密码哈希成功，长度: %d\n", len(hashedPassword))
	}
	
	isValid := CheckPassword(password, hashedPassword)
	if isValid {
		fmt.Printf("✅ 密码验证成功\n")
	} else {
		fmt.Printf("❌ 密码验证失败\n")
	}
	
	fmt.Println("\n=== 性能测试完成 ===")
}

// Version 返回版本信息
func Version() string {
	return "Crypto v1.0.0 - helwd工具包加密模块"
}

// Info 返回工具信息
func Info() map[string]interface{} {
	return map[string]interface{}{
		"name":        "Crypto",
		"version":     "1.0.0",
		"description": "Go加密工具包",
		"features": []string{
			"AES加密解密",
			"RSA加密解密",
			"哈希算法",
			"密码哈希",
			"数字签名",
			"文件加密",
		},
		"algorithms": map[string][]string{
			"symmetric":  {"AES-128", "AES-192", "AES-256"},
			"asymmetric": {"RSA-1024", "RSA-2048", "RSA-3072", "RSA-4096"},
			"hash":       {"MD5", "SHA1", "SHA256", "SHA512"},
			"password":   {"bcrypt"},
		},
	}
}
