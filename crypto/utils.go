package crypto

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"math/big"
	"os"

	"golang.org/x/crypto/pbkdf2"
)

// GenerateRandomBytes 生成随机字节
func GenerateRandomBytes(length int) ([]byte, error) {
	if length <= 0 {
		return nil, fmt.Errorf("长度必须大于0")
	}

	bytes := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, bytes); err != nil {
		return nil, fmt.Errorf("生成随机字节失败: %w", err)
	}

	return bytes, nil
}

// GenerateRandomString 生成随机字符串（字母数字）
func GenerateRandomString(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	return GenerateRandomStringFromChars(length, charset)
}

// GenerateRandomStringFromChars 从指定字符集生成随机字符串
func GenerateRandomStringFromChars(length int, charset string) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("长度必须大于0")
	}

	if len(charset) == 0 {
		return "", fmt.Errorf("字符集不能为空")
	}

	result := make([]byte, length)
	charsetLen := big.NewInt(int64(len(charset)))

	for i := 0; i < length; i++ {
		randomIndex, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", fmt.Errorf("生成随机索引失败: %w", err)
		}
		result[i] = charset[randomIndex.Int64()]
	}

	return string(result), nil
}

// GenerateRandomHex 生成随机十六进制字符串
func GenerateRandomHex(length int) (string, error) {
	bytes, err := GenerateRandomBytes(length)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// GenerateRandomBase64 生成随机Base64字符串
func GenerateRandomBase64(length int) (string, error) {
	bytes, err := GenerateRandomBytes(length)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}

// Base64Encode Base64编码
func Base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// Base64Decode Base64解码
func Base64Decode(data string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(data)
}

// Base64URLEncode Base64 URL安全编码
func Base64URLEncode(data []byte) string {
	return base64.URLEncoding.EncodeToString(data)
}

// Base64URLDecode Base64 URL安全解码
func Base64URLDecode(data string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(data)
}

// HexEncode 十六进制编码
func HexEncode(data []byte) string {
	return hex.EncodeToString(data)
}

// HexDecode 十六进制解码
func HexDecode(data string) ([]byte, error) {
	return hex.DecodeString(data)
}

// SecureCompare 安全比较两个字节切片（防止时序攻击）
func SecureCompare(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	var result byte
	for i := 0; i < len(a); i++ {
		result |= a[i] ^ b[i]
	}

	return result == 0
}

// SecureCompareString 安全比较两个字符串（防止时序攻击）
func SecureCompareString(a, b string) bool {
	return SecureCompare([]byte(a), []byte(b))
}

// ZeroBytes 安全清零字节切片
func ZeroBytes(data []byte) {
	for i := range data {
		data[i] = 0
	}
}

// ZeroString 安全清零字符串（通过字节切片）
func ZeroString(s *string) {
	if s != nil {
		data := []byte(*s)
		ZeroBytes(data)
		*s = ""
	}
}

// IsValidHex 检查是否为有效的十六进制字符串
func IsValidHex(s string) bool {
	_, err := hex.DecodeString(s)
	return err == nil
}

// IsValidBase64 检查是否为有效的Base64字符串
func IsValidBase64(s string) bool {
	_, err := base64.StdEncoding.DecodeString(s)
	return err == nil
}

// PadBytes 字节填充到指定长度
func PadBytes(data []byte, length int, padByte byte) []byte {
	if len(data) >= length {
		return data
	}

	padded := make([]byte, length)
	copy(padded, data)

	for i := len(data); i < length; i++ {
		padded[i] = padByte
	}

	return padded
}

// UnpadBytes 移除字节填充
func UnpadBytes(data []byte, padByte byte) []byte {
	for i := len(data) - 1; i >= 0; i-- {
		if data[i] != padByte {
			return data[:i+1]
		}
	}
	return []byte{}
}

// XORBytes 异或操作
func XORBytes(a, b []byte) []byte {
	if len(a) != len(b) {
		return nil
	}

	result := make([]byte, len(a))
	for i := 0; i < len(a); i++ {
		result[i] = a[i] ^ b[i]
	}

	return result
}

// RotateLeft 循环左移
func RotateLeft(data []byte, positions int) []byte {
	if len(data) == 0 {
		return data
	}

	positions = positions % len(data)
	if positions == 0 {
		return data
	}

	result := make([]byte, len(data))
	copy(result, data[positions:])
	copy(result[len(data)-positions:], data[:positions])

	return result
}

// RotateRight 循环右移
func RotateRight(data []byte, positions int) []byte {
	if len(data) == 0 {
		return data
	}

	positions = positions % len(data)
	if positions == 0 {
		return data
	}

	return RotateLeft(data, len(data)-positions)
}

// FileExists 检查文件是否存在
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

// SecureDeleteFile 安全删除文件（多次覆写）
func SecureDeleteFile(filename string) error {
	// 检查文件是否存在
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return nil // 文件不存在，认为删除成功
	}
	if err != nil {
		return fmt.Errorf("获取文件信息失败: %w", err)
	}

	// 打开文件进行覆写
	file, err := os.OpenFile(filename, os.O_WRONLY, 0)
	if err != nil {
		return fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	fileSize := info.Size()

	// 多次覆写文件内容
	patterns := []byte{0x00, 0xFF, 0xAA, 0x55}
	for _, pattern := range patterns {
		// 移动到文件开头
		if _, err := file.Seek(0, 0); err != nil {
			return fmt.Errorf("文件定位失败: %w", err)
		}

		// 用指定模式覆写整个文件
		buffer := make([]byte, 4096)
		for i := range buffer {
			buffer[i] = pattern
		}

		remaining := fileSize
		for remaining > 0 {
			writeSize := int64(len(buffer))
			if remaining < writeSize {
				writeSize = remaining
			}

			if _, err := file.Write(buffer[:writeSize]); err != nil {
				return fmt.Errorf("文件覆写失败: %w", err)
			}

			remaining -= writeSize
		}

		// 强制写入磁盘
		if err := file.Sync(); err != nil {
			return fmt.Errorf("文件同步失败: %w", err)
		}
	}

	// 关闭文件
	file.Close()

	// 删除文件
	if err := os.Remove(filename); err != nil {
		return fmt.Errorf("删除文件失败: %w", err)
	}

	return nil
}

// GenerateUUID 生成简单的UUID（基于随机数）
func GenerateUUID() (string, error) {
	bytes, err := GenerateRandomBytes(16)
	if err != nil {
		return "", err
	}

	// 设置版本和变体位
	bytes[6] = (bytes[6] & 0x0f) | 0x40 // 版本4
	bytes[8] = (bytes[8] & 0x3f) | 0x80 // 变体10

	return fmt.Sprintf("%x-%x-%x-%x-%x",
		bytes[0:4], bytes[4:6], bytes[6:8], bytes[8:10], bytes[10:16]), nil
}

// SplitBytes 将字节切片分割成指定大小的块
func SplitBytes(data []byte, chunkSize int) [][]byte {
	if chunkSize <= 0 {
		return nil
	}

	var chunks [][]byte
	for i := 0; i < len(data); i += chunkSize {
		end := i + chunkSize
		if end > len(data) {
			end = len(data)
		}
		chunks = append(chunks, data[i:end])
	}

	return chunks
}

// JoinBytes 连接多个字节切片
func JoinBytes(chunks [][]byte) []byte {
	totalLen := 0
	for _, chunk := range chunks {
		totalLen += len(chunk)
	}

	result := make([]byte, totalLen)
	offset := 0
	for _, chunk := range chunks {
		copy(result[offset:], chunk)
		offset += len(chunk)
	}

	return result
}

// PBKDF2 密钥派生函数
func PBKDF2(password, salt []byte, iterations, keyLength int, hashFunc func([]byte) []byte) []byte {
	// 简化处理，默认使用SHA256
	return pbkdf2.Key(password, salt, iterations, keyLength, sha256.New)
}
