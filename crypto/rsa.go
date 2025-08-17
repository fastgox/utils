package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"
)

// GenerateRSAKeyPair 生成RSA密钥对（返回PEM格式字符串）
func GenerateRSAKeyPair(keySize int) (privateKey, publicKey string, err error) {
	if err := ValidateRSAKeySize(keySize); err != nil {
		return "", "", err
	}

	// 生成私钥
	privKey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return "", "", fmt.Errorf("生成RSA私钥失败: %w", err)
	}

	// 私钥转PEM格式
	privKeyBytes, err := x509.MarshalPKCS8PrivateKey(privKey)
	if err != nil {
		return "", "", fmt.Errorf("序列化私钥失败: %w", err)
	}

	privKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privKeyBytes,
	})

	// 公钥转PEM格式
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(&privKey.PublicKey)
	if err != nil {
		return "", "", fmt.Errorf("序列化公钥失败: %w", err)
	}

	pubKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubKeyBytes,
	})

	return string(privKeyPEM), string(pubKeyPEM), nil
}

// GenerateRSAKeyPairToFile 生成RSA密钥对并保存到文件
func GenerateRSAKeyPairToFile(keySize int, privateKeyFile, publicKeyFile string) error {
	privateKey, publicKey, err := GenerateRSAKeyPair(keySize)
	if err != nil {
		return err
	}

	// 保存私钥
	err = os.WriteFile(privateKeyFile, []byte(privateKey), 0600)
	if err != nil {
		return fmt.Errorf("保存私钥文件失败: %w", err)
	}

	// 保存公钥
	err = os.WriteFile(publicKeyFile, []byte(publicKey), 0644)
	if err != nil {
		return fmt.Errorf("保存公钥文件失败: %w", err)
	}

	return nil
}

// RSAEncrypt RSA公钥加密
func RSAEncrypt(plaintext, publicKeyPEM string) (string, error) {
	// 解析公钥
	pubKey, err := parsePublicKey(publicKeyPEM)
	if err != nil {
		return "", err
	}

	// 加密
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pubKey, []byte(plaintext), nil)
	if err != nil {
		return "", fmt.Errorf("RSA加密失败: %w", err)
	}

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// RSADecrypt RSA私钥解密
func RSADecrypt(ciphertext, privateKeyPEM string) (string, error) {
	// 解析私钥
	privKey, err := parsePrivateKey(privateKeyPEM)
	if err != nil {
		return "", err
	}

	// Base64解码
	ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("base64解码失败: %w", err)
	}

	// 解密
	plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privKey, ciphertextBytes, nil)
	if err != nil {
		return "", fmt.Errorf("RSA解密失败: %w", err)
	}

	return string(plaintext), nil
}

// RSASign RSA私钥签名
func RSASign(data, privateKeyPEM string) (string, error) {
	// 解析私钥
	privKey, err := parsePrivateKey(privateKeyPEM)
	if err != nil {
		return "", err
	}

	// 计算哈希
	hash := sha256.Sum256([]byte(data))

	// 签名
	signature, err := rsa.SignPKCS1v15(rand.Reader, privKey, crypto.SHA256, hash[:])
	if err != nil {
		return "", fmt.Errorf("RSA签名失败: %w", err)
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

// RSAVerify RSA公钥验证签名
func RSAVerify(data, signature, publicKeyPEM string) (bool, error) {
	// 解析公钥
	pubKey, err := parsePublicKey(publicKeyPEM)
	if err != nil {
		return false, err
	}

	// Base64解码签名
	signatureBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false, fmt.Errorf("base64解码失败: %w", err)
	}

	// 计算哈希
	hash := sha256.Sum256([]byte(data))

	// 验证签名
	err = rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hash[:], signatureBytes)
	if err != nil {
		return false, nil // 签名无效，但不是错误
	}

	return true, nil
}

// RSAEncryptBytes RSA公钥加密（字节）
func RSAEncryptBytes(plaintext []byte, publicKeyPEM string) ([]byte, error) {
	// 解析公钥
	pubKey, err := parsePublicKey(publicKeyPEM)
	if err != nil {
		return nil, err
	}

	// 计算最大加密长度
	keySize := pubKey.Size()
	maxPlaintextSize := keySize - 2*sha256.Size - 2

	if len(plaintext) > maxPlaintextSize {
		return nil, fmt.Errorf("明文长度超过限制: %d > %d", len(plaintext), maxPlaintextSize)
	}

	// 加密
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pubKey, plaintext, nil)
	if err != nil {
		return nil, fmt.Errorf("RSA加密失败: %w", err)
	}

	return ciphertext, nil
}

// RSADecryptBytes RSA私钥解密（字节）
func RSADecryptBytes(ciphertext []byte, privateKeyPEM string) ([]byte, error) {
	// 解析私钥
	privKey, err := parsePrivateKey(privateKeyPEM)
	if err != nil {
		return nil, err
	}

	// 解密
	plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privKey, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("RSA解密失败: %w", err)
	}

	return plaintext, nil
}

// LoadRSAPrivateKeyFromFile 从文件加载RSA私钥
func LoadRSAPrivateKeyFromFile(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("读取私钥文件失败: %w", err)
	}
	return string(data), nil
}

// LoadRSAPublicKeyFromFile 从文件加载RSA公钥
func LoadRSAPublicKeyFromFile(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("读取公钥文件失败: %w", err)
	}
	return string(data), nil
}

// GetRSAPublicKeyFromPrivate 从私钥提取公钥
func GetRSAPublicKeyFromPrivate(privateKeyPEM string) (string, error) {
	// 解析私钥
	privKey, err := parsePrivateKey(privateKeyPEM)
	if err != nil {
		return "", err
	}

	// 提取公钥
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(&privKey.PublicKey)
	if err != nil {
		return "", fmt.Errorf("序列化公钥失败: %w", err)
	}

	pubKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubKeyBytes,
	})

	return string(pubKeyPEM), nil
}

// parsePrivateKey 解析私钥PEM
func parsePrivateKey(privateKeyPEM string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return nil, fmt.Errorf("无效的PEM格式私钥")
	}

	// 尝试解析PKCS8格式
	privKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		// 尝试解析PKCS1格式
		privKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("解析私钥失败: %w", err)
		}
	}

	rsaPrivKey, ok := privKey.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("不是RSA私钥")
	}

	return rsaPrivKey, nil
}

// parsePublicKey 解析公钥PEM
func parsePublicKey(publicKeyPEM string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return nil, fmt.Errorf("无效的PEM格式公钥")
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("解析公钥失败: %w", err)
	}

	rsaPubKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("不是RSA公钥")
	}

	return rsaPubKey, nil
}

// RSAKeyInfo 获取RSA密钥信息
func RSAKeyInfo(keyPEM string) (keySize int, keyType string, err error) {
	// 尝试解析为私钥
	if privKey, err := parsePrivateKey(keyPEM); err == nil {
		return privKey.Size() * 8, "private", nil
	}

	// 尝试解析为公钥
	if pubKey, err := parsePublicKey(keyPEM); err == nil {
		return pubKey.Size() * 8, "public", nil
	}

	return 0, "", fmt.Errorf("无法解析RSA密钥")
}
