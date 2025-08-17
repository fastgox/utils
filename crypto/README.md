# Crypto - 加密工具

一个功能强大的Go加密工具包，提供常用的加密、解密、哈希和签名功能。

## 🚀 特性

- **🔐 AES加密**: 支持AES-128/192/256加密解密
- **🔑 RSA加密**: 支持RSA公钥/私钥加密解密
- **🔒 哈希算法**: 支持MD5、SHA1、SHA256、SHA512
- **🛡️ 密码哈希**: 支持bcrypt密码加盐哈希
- **📝 数字签名**: 支持RSA/ECDSA数字签名
- **🎯 简洁API**: 类似其他helwd工具的简洁设计
- **⚡ 高性能**: 优化的加密算法实现

## 📦 安装

```bash
go get github.com/fastgox/utils/crypto
```

## 🎯 快速开始

### AES加密

```go
package main

import (
    "fmt"
    "github.com/fastgox/utils/crypto"
)

func main() {
    // AES加密
    plaintext := "Hello, World!"
    key := "my-secret-key-32-bytes-long!!"
    
    encrypted, err := crypto.AESEncrypt(plaintext, key)
    if err != nil {
        panic(err)
    }
    fmt.Printf("加密结果: %s\n", encrypted)
    
    // AES解密
    decrypted, err := crypto.AESDecrypt(encrypted, key)
    if err != nil {
        panic(err)
    }
    fmt.Printf("解密结果: %s\n", decrypted)
}
```

### RSA加密

```go
func main() {
    // 生成RSA密钥对
    privateKey, publicKey, err := crypto.GenerateRSAKeyPair(2048)
    if err != nil {
        panic(err)
    }
    
    // RSA加密
    plaintext := "Hello, RSA!"
    encrypted, err := crypto.RSAEncrypt(plaintext, publicKey)
    if err != nil {
        panic(err)
    }
    
    // RSA解密
    decrypted, err := crypto.RSADecrypt(encrypted, privateKey)
    if err != nil {
        panic(err)
    }
    fmt.Printf("解密结果: %s\n", decrypted)
}
```

### 哈希算法

```go
func main() {
    data := "Hello, Hash!"
    
    // MD5哈希
    md5Hash := crypto.MD5(data)
    fmt.Printf("MD5: %s\n", md5Hash)
    
    // SHA256哈希
    sha256Hash := crypto.SHA256(data)
    fmt.Printf("SHA256: %s\n", sha256Hash)
    
    // SHA512哈希
    sha512Hash := crypto.SHA512(data)
    fmt.Printf("SHA512: %s\n", sha512Hash)
}
```

### 密码哈希

```go
func main() {
    password := "my-password"
    
    // 生成密码哈希
    hashedPassword, err := crypto.HashPassword(password)
    if err != nil {
        panic(err)
    }
    fmt.Printf("密码哈希: %s\n", hashedPassword)
    
    // 验证密码
    isValid := crypto.CheckPassword(password, hashedPassword)
    fmt.Printf("密码验证: %v\n", isValid)
}
```

### 数字签名

```go
func main() {
    // 生成密钥对
    privateKey, publicKey, err := crypto.GenerateRSAKeyPair(2048)
    if err != nil {
        panic(err)
    }
    
    data := "Hello, Signature!"
    
    // 生成签名
    signature, err := crypto.RSASign(data, privateKey)
    if err != nil {
        panic(err)
    }
    
    // 验证签名
    isValid, err := crypto.RSAVerify(data, signature, publicKey)
    if err != nil {
        panic(err)
    }
    fmt.Printf("签名验证: %v\n", isValid)
}
```

## 📚 API 文档

### AES加密函数

```go
// AES加密解密
crypto.AESEncrypt(plaintext, key string) (string, error)
crypto.AESDecrypt(ciphertext, key string) (string, error)

// AES加密解密（字节）
crypto.AESEncryptBytes(plaintext, key []byte) ([]byte, error)
crypto.AESDecryptBytes(ciphertext, key []byte) ([]byte, error)

// 生成AES密钥
crypto.GenerateAESKey(keySize int) ([]byte, error) // 16, 24, 32
```

### RSA加密函数

```go
// RSA密钥生成
crypto.GenerateRSAKeyPair(keySize int) (privateKey, publicKey string, err error)
crypto.GenerateRSAKeyPairToFile(keySize int, privateKeyFile, publicKeyFile string) error

// RSA加密解密
crypto.RSAEncrypt(plaintext, publicKey string) (string, error)
crypto.RSADecrypt(ciphertext, privateKey string) (string, error)

// RSA签名验证
crypto.RSASign(data, privateKey string) (string, error)
crypto.RSAVerify(data, signature, publicKey string) (bool, error)
```

### 哈希函数

```go
// 基本哈希
crypto.MD5(data string) string
crypto.SHA1(data string) string
crypto.SHA256(data string) string
crypto.SHA512(data string) string

// 字节哈希
crypto.MD5Bytes(data []byte) []byte
crypto.SHA256Bytes(data []byte) []byte

// HMAC
crypto.HMACSHA256(data, key string) string
crypto.HMACSHA512(data, key string) string
```

### 密码哈希函数

```go
// bcrypt密码哈希
crypto.HashPassword(password string) (string, error)
crypto.CheckPassword(password, hashedPassword string) bool

// 自定义成本
crypto.HashPasswordWithCost(password string, cost int) (string, error)
```

### 工具函数

```go
// 随机数生成
crypto.GenerateRandomBytes(length int) ([]byte, error)
crypto.GenerateRandomString(length int) (string, error)

// Base64编码
crypto.Base64Encode(data []byte) string
crypto.Base64Decode(data string) ([]byte, error)

// Hex编码
crypto.HexEncode(data []byte) string
crypto.HexDecode(data string) ([]byte, error)
```

## 🔧 高级功能

### 配置选项

```go
// 设置默认AES密钥
crypto.SetDefaultAESKey("your-default-key-32-bytes!!")

// 使用默认密钥加密
encrypted, err := crypto.AESEncryptDefault("Hello, World!")

// 设置默认bcrypt成本
crypto.SetDefaultBcryptCost(12)
```

### 文件加密

```go
// 加密文件
err := crypto.EncryptFile("input.txt", "output.enc", "my-key")

// 解密文件
err := crypto.DecryptFile("output.enc", "decrypted.txt", "my-key")
```

## 🛡️ 安全建议

1. **密钥管理**: 不要在代码中硬编码密钥，使用环境变量或配置文件
2. **密钥长度**: AES使用32字节密钥，RSA使用至少2048位
3. **随机性**: 使用加密安全的随机数生成器
4. **密码哈希**: 使用bcrypt等安全的密码哈希算法
5. **定期更新**: 定期更新密钥和算法

## 🎮 运行示例

```bash
# 克隆项目
git clone https://github.com/fastgox/utils.git
cd utils

# 运行加密演示
go run main.go
```

## 🌟 特色

- **🎯 专注实用**: 只包含最常用的加密功能
- **📝 简洁API**: 链式调用，代码简洁优雅
- **⚡ 开箱即用**: 无需复杂配置，直接使用
- **🔧 灵活配置**: 支持自定义参数和选项
- **💡 最佳实践**: 遵循加密安全最佳实践

## 🤝 贡献

欢迎提交Issue和Pull Request！

## 📄 许可证

MIT License
