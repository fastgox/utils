package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

// AESEncrypt AES加密（字符串）
func AESEncrypt(plaintext, key string) (string, error) {
	keyBytes := []byte(key)
	plaintextBytes := []byte(plaintext)
	
	ciphertext, err := AESEncryptBytes(plaintextBytes, keyBytes)
	if err != nil {
		return "", err
	}
	
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// AESDecrypt AES解密（字符串）
func AESDecrypt(ciphertext, key string) (string, error) {
	keyBytes := []byte(key)
	
	ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("base64解码失败: %w", err)
	}
	
	plaintext, err := AESDecryptBytes(ciphertextBytes, keyBytes)
	if err != nil {
		return "", err
	}
	
	return string(plaintext), nil
}

// AESEncryptBytes AES加密（字节）
func AESEncryptBytes(plaintext, key []byte) ([]byte, error) {
	// 验证密钥长度
	if err := ValidateAESKeySize(len(key)); err != nil {
		return nil, err
	}
	
	// 创建AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("创建AES cipher失败: %w", err)
	}
	
	// 使用GCM模式
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("创建GCM失败: %w", err)
	}
	
	// 生成随机nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("生成nonce失败: %w", err)
	}
	
	// 加密
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// AESDecryptBytes AES解密（字节）
func AESDecryptBytes(ciphertext, key []byte) ([]byte, error) {
	// 验证密钥长度
	if err := ValidateAESKeySize(len(key)); err != nil {
		return nil, err
	}
	
	// 创建AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("创建AES cipher失败: %w", err)
	}
	
	// 使用GCM模式
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("创建GCM失败: %w", err)
	}
	
	// 检查密文长度
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, ErrInvalidCiphertext
	}
	
	// 提取nonce和密文
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	
	// 解密
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("AES解密失败: %w", err)
	}
	
	return plaintext, nil
}

// AESEncryptCBC AES-CBC模式加密
func AESEncryptCBC(plaintext, key []byte) ([]byte, error) {
	// 验证密钥长度
	if err := ValidateAESKeySize(len(key)); err != nil {
		return nil, err
	}
	
	// 创建AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("创建AES cipher失败: %w", err)
	}
	
	// PKCS7填充
	plaintext = pkcs7Padding(plaintext, aes.BlockSize)
	
	// 生成随机IV
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, fmt.Errorf("生成IV失败: %w", err)
	}
	
	// 创建CBC模式
	mode := cipher.NewCBCEncrypter(block, iv)
	
	// 加密
	ciphertext := make([]byte, len(plaintext))
	mode.CryptBlocks(ciphertext, plaintext)
	
	// 将IV添加到密文前面
	result := make([]byte, len(iv)+len(ciphertext))
	copy(result[:len(iv)], iv)
	copy(result[len(iv):], ciphertext)
	
	return result, nil
}

// AESDecryptCBC AES-CBC模式解密
func AESDecryptCBC(ciphertext, key []byte) ([]byte, error) {
	// 验证密钥长度
	if err := ValidateAESKeySize(len(key)); err != nil {
		return nil, err
	}
	
	// 检查密文长度
	if len(ciphertext) < aes.BlockSize {
		return nil, ErrInvalidCiphertext
	}
	
	// 创建AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("创建AES cipher失败: %w", err)
	}
	
	// 提取IV和密文
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	
	// 检查密文长度是否为块大小的倍数
	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, ErrInvalidCiphertext
	}
	
	// 创建CBC模式
	mode := cipher.NewCBCDecrypter(block, iv)
	
	// 解密
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)
	
	// 去除PKCS7填充
	plaintext, err = pkcs7UnPadding(plaintext)
	if err != nil {
		return nil, fmt.Errorf("去除填充失败: %w", err)
	}
	
	return plaintext, nil
}

// GenerateAESKey 生成AES密钥
func GenerateAESKey(keySize int) ([]byte, error) {
	if err := ValidateAESKeySize(keySize); err != nil {
		return nil, err
	}
	
	key := make([]byte, keySize)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, fmt.Errorf("生成AES密钥失败: %w", err)
	}
	
	return key, nil
}

// AESEncryptDefault 使用默认密钥加密
func AESEncryptDefault(plaintext string) (string, error) {
	if globalConfig.DefaultAESKey == "" {
		return "", fmt.Errorf("未设置默认AES密钥")
	}
	return AESEncrypt(plaintext, globalConfig.DefaultAESKey)
}

// AESDecryptDefault 使用默认密钥解密
func AESDecryptDefault(ciphertext string) (string, error) {
	if globalConfig.DefaultAESKey == "" {
		return "", fmt.Errorf("未设置默认AES密钥")
	}
	return AESDecrypt(ciphertext, globalConfig.DefaultAESKey)
}

// pkcs7Padding PKCS7填充
func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := make([]byte, padding)
	for i := range padtext {
		padtext[i] = byte(padding)
	}
	return append(data, padtext...)
}

// pkcs7UnPadding 去除PKCS7填充
func pkcs7UnPadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, fmt.Errorf("数据为空")
	}
	
	unpadding := int(data[length-1])
	if unpadding > length {
		return nil, fmt.Errorf("无效的填充")
	}
	
	return data[:(length - unpadding)], nil
}

// AESKeyFromPassword 从密码生成AES密钥
func AESKeyFromPassword(password, salt string, keySize int) ([]byte, error) {
	if err := ValidateAESKeySize(keySize); err != nil {
		return nil, err
	}
	
	// 使用PBKDF2生成密钥
	return PBKDF2([]byte(password), []byte(salt), 10000, keySize, SHA256Bytes), nil
}

// AESEncryptWithPassword 使用密码加密
func AESEncryptWithPassword(plaintext, password string) (string, error) {
	// 生成随机盐
	salt, err := GenerateRandomBytes(16)
	if err != nil {
		return "", fmt.Errorf("生成盐失败: %w", err)
	}
	
	// 从密码生成密钥
	key, err := AESKeyFromPassword(password, string(salt), AES256KeySize)
	if err != nil {
		return "", err
	}
	
	// 加密
	ciphertext, err := AESEncryptBytes([]byte(plaintext), key)
	if err != nil {
		return "", err
	}
	
	// 将盐和密文组合
	result := append(salt, ciphertext...)
	return base64.StdEncoding.EncodeToString(result), nil
}

// AESDecryptWithPassword 使用密码解密
func AESDecryptWithPassword(ciphertext, password string) (string, error) {
	// Base64解码
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("base64解码失败: %w", err)
	}
	
	// 检查数据长度
	if len(data) < 16 {
		return "", ErrInvalidCiphertext
	}
	
	// 提取盐和密文
	salt := data[:16]
	ciphertextBytes := data[16:]
	
	// 从密码生成密钥
	key, err := AESKeyFromPassword(password, string(salt), AES256KeySize)
	if err != nil {
		return "", err
	}
	
	// 解密
	plaintext, err := AESDecryptBytes(ciphertextBytes, key)
	if err != nil {
		return "", err
	}
	
	return string(plaintext), nil
}
