package crypto

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
)

// MD5 计算MD5哈希
func MD5(data string) string {
	return hex.EncodeToString(MD5Bytes([]byte(data)))
}

// MD5Bytes 计算MD5哈希（字节）
func MD5Bytes(data []byte) []byte {
	hash := md5.Sum(data)
	return hash[:]
}

// SHA1 计算SHA1哈希
func SHA1(data string) string {
	return hex.EncodeToString(SHA1Bytes([]byte(data)))
}

// SHA1Bytes 计算SHA1哈希（字节）
func SHA1Bytes(data []byte) []byte {
	hash := sha1.Sum(data)
	return hash[:]
}

// SHA256 计算SHA256哈希
func SHA256(data string) string {
	return hex.EncodeToString(SHA256Bytes([]byte(data)))
}

// SHA256Bytes 计算SHA256哈希（字节）
func SHA256Bytes(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

// SHA512 计算SHA512哈希
func SHA512(data string) string {
	return hex.EncodeToString(SHA512Bytes([]byte(data)))
}

// SHA512Bytes 计算SHA512哈希（字节）
func SHA512Bytes(data []byte) []byte {
	hash := sha512.Sum512(data)
	return hash[:]
}

// HMACSHA256 计算HMAC-SHA256
func HMACSHA256(data, key string) string {
	return hex.EncodeToString(HMACSHA256Bytes([]byte(data), []byte(key)))
}

// HMACSHA256Bytes 计算HMAC-SHA256（字节）
func HMACSHA256Bytes(data, key []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}

// HMACSHA512 计算HMAC-SHA512
func HMACSHA512(data, key string) string {
	return hex.EncodeToString(HMACSHA512Bytes([]byte(data), []byte(key)))
}

// HMACSHA512Bytes 计算HMAC-SHA512（字节）
func HMACSHA512Bytes(data, key []byte) []byte {
	h := hmac.New(sha512.New, key)
	h.Write(data)
	return h.Sum(nil)
}

// HMACMD5 计算HMAC-MD5
func HMACMD5(data, key string) string {
	return hex.EncodeToString(HMACMD5Bytes([]byte(data), []byte(key)))
}

// HMACMD5Bytes 计算HMAC-MD5（字节）
func HMACMD5Bytes(data, key []byte) []byte {
	h := hmac.New(md5.New, key)
	h.Write(data)
	return h.Sum(nil)
}

// Hash 通用哈希函数
func Hash(data []byte, algorithm HashAlgorithm) []byte {
	var h hash.Hash

	switch algorithm {
	case HashMD5:
		h = md5.New()
	case HashSHA1:
		h = sha1.New()
	case HashSHA256:
		h = sha256.New()
	case HashSHA512:
		h = sha512.New()
	default:
		return nil
	}

	h.Write(data)
	return h.Sum(nil)
}

// HashString 通用哈希函数（字符串）
func HashString(data string, algorithm HashAlgorithm) string {
	hashBytes := Hash([]byte(data), algorithm)
	if hashBytes == nil {
		return ""
	}
	return hex.EncodeToString(hashBytes)
}

// HMAC 通用HMAC函数
func HMAC(data, key []byte, algorithm HashAlgorithm) []byte {
	var h func() hash.Hash

	switch algorithm {
	case HashMD5:
		h = md5.New
	case HashSHA1:
		h = sha1.New
	case HashSHA256:
		h = sha256.New
	case HashSHA512:
		h = sha512.New
	default:
		return nil
	}

	mac := hmac.New(h, key)
	mac.Write(data)
	return mac.Sum(nil)
}

// HMACString 通用HMAC函数（字符串）
func HMACString(data, key string, algorithm HashAlgorithm) string {
	macBytes := HMAC([]byte(data), []byte(key), algorithm)
	if macBytes == nil {
		return ""
	}
	return hex.EncodeToString(macBytes)
}

// VerifyHMAC 验证HMAC
func VerifyHMAC(data, key []byte, expectedMAC []byte, algorithm HashAlgorithm) bool {
	computedMAC := HMAC(data, key, algorithm)
	return hmac.Equal(expectedMAC, computedMAC)
}

// VerifyHMACString 验证HMAC（字符串）
func VerifyHMACString(data, key, expectedMAC string, algorithm HashAlgorithm) bool {
	expectedMACBytes, err := hex.DecodeString(expectedMAC)
	if err != nil {
		return false
	}
	return VerifyHMAC([]byte(data), []byte(key), expectedMACBytes, algorithm)
}

// FileHash 计算文件哈希
func FileHash(filename string, algorithm HashAlgorithm) (string, error) {
	data, err := readFile(filename)
	if err != nil {
		return "", fmt.Errorf("读取文件失败: %w", err)
	}

	hashBytes := Hash(data, algorithm)
	if hashBytes == nil {
		return "", fmt.Errorf("不支持的哈希算法: %s", algorithm.String())
	}

	return hex.EncodeToString(hashBytes), nil
}

// FileMD5 计算文件MD5
func FileMD5(filename string) (string, error) {
	return FileHash(filename, HashMD5)
}

// FileSHA256 计算文件SHA256
func FileSHA256(filename string) (string, error) {
	return FileHash(filename, HashSHA256)
}

// FileSHA512 计算文件SHA512
func FileSHA512(filename string) (string, error) {
	return FileHash(filename, HashSHA512)
}

// CompareHash 比较哈希值
func CompareHash(hash1, hash2 string) bool {
	return hash1 == hash2
}

// ValidateHash 验证哈希格式
func ValidateHash(hashStr string, algorithm HashAlgorithm) bool {
	var expectedLength int

	switch algorithm {
	case HashMD5:
		expectedLength = 32
	case HashSHA1:
		expectedLength = 40
	case HashSHA256:
		expectedLength = 64
	case HashSHA512:
		expectedLength = 128
	default:
		return false
	}

	if len(hashStr) != expectedLength {
		return false
	}

	// 检查是否为有效的十六进制字符串
	_, err := hex.DecodeString(hashStr)
	return err == nil
}

// HashMultiple 计算多个数据的组合哈希
func HashMultiple(data [][]byte, algorithm HashAlgorithm) []byte {
	var h hash.Hash

	switch algorithm {
	case HashMD5:
		h = md5.New()
	case HashSHA1:
		h = sha1.New()
	case HashSHA256:
		h = sha256.New()
	case HashSHA512:
		h = sha512.New()
	default:
		return nil
	}

	for _, d := range data {
		h.Write(d)
	}

	return h.Sum(nil)
}

// HashMultipleString 计算多个字符串的组合哈希
func HashMultipleString(data []string, algorithm HashAlgorithm) string {
	byteData := make([][]byte, len(data))
	for i, s := range data {
		byteData[i] = []byte(s)
	}

	hashBytes := HashMultiple(byteData, algorithm)
	if hashBytes == nil {
		return ""
	}

	return hex.EncodeToString(hashBytes)
}

// readFile 读取文件内容（简化版，实际应该使用io包进行流式读取）
func readFile(filename string) ([]byte, error) {
	// 这里应该实现文件读取逻辑
	// 为了简化，暂时返回错误
	return nil, fmt.Errorf("文件读取功能待实现")
}
