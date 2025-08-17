package test

import (
	"strings"
	"testing"

	"github.com/fastgox/utils/crypto"
)

func TestCryptoAES(t *testing.T) {
	t.Run("AES基本加密解密", func(t *testing.T) {
		plaintext := "Hello, World! 这是一个测试消息。"
		key := "12345678901234567890123456789012" // 32字节密钥

		// 加密
		encrypted, err := crypto.AESEncrypt(plaintext, key)
		if err != nil {
			t.Fatalf("AES加密失败: %v", err)
		}

		if encrypted == "" {
			t.Fatal("加密结果为空")
		}

		// 解密
		decrypted, err := crypto.AESDecrypt(encrypted, key)
		if err != nil {
			t.Fatalf("AES解密失败: %v", err)
		}

		if decrypted != plaintext {
			t.Fatalf("解密结果不匹配: 期望 %s, 得到 %s", plaintext, decrypted)
		}

		t.Logf("AES加密解密测试通过")
	})

	t.Run("AES密码加密", func(t *testing.T) {
		plaintext := "Secret message"
		password := "my-password"

		// 使用密码加密
		encrypted, err := crypto.AESEncryptWithPassword(plaintext, password)
		if err != nil {
			t.Fatalf("密码加密失败: %v", err)
		}

		// 使用密码解密
		decrypted, err := crypto.AESDecryptWithPassword(encrypted, password)
		if err != nil {
			t.Fatalf("密码解密失败: %v", err)
		}

		if decrypted != plaintext {
			t.Fatalf("解密结果不匹配: 期望 %s, 得到 %s", plaintext, decrypted)
		}

		t.Logf("AES密码加密解密测试通过")
	})

	t.Run("AES密钥生成", func(t *testing.T) {
		// 测试不同长度的密钥生成
		keySizes := []int{16, 24, 32}

		for _, size := range keySizes {
			key, err := crypto.GenerateAESKey(size)
			if err != nil {
				t.Fatalf("生成%d字节AES密钥失败: %v", size, err)
			}

			if len(key) != size {
				t.Fatalf("密钥长度不正确: 期望 %d, 得到 %d", size, len(key))
			}
		}

		t.Logf("AES密钥生成测试通过")
	})
}

func TestCryptoRSA(t *testing.T) {
	t.Run("RSA密钥生成", func(t *testing.T) {
		privateKey, publicKey, err := crypto.GenerateRSAKeyPair(2048)
		if err != nil {
			t.Fatalf("RSA密钥生成失败: %v", err)
		}

		if privateKey == "" || publicKey == "" {
			t.Fatal("生成的密钥为空")
		}

		if !strings.Contains(privateKey, "PRIVATE KEY") {
			t.Fatal("私钥格式不正确")
		}

		if !strings.Contains(publicKey, "PUBLIC KEY") {
			t.Fatal("公钥格式不正确")
		}

		t.Logf("RSA密钥生成测试通过")
	})

	t.Run("RSA加密解密", func(t *testing.T) {
		// 生成密钥对
		privateKey, publicKey, err := crypto.GenerateRSAKeyPair(2048)
		if err != nil {
			t.Fatalf("RSA密钥生成失败: %v", err)
		}

		plaintext := "Hello, RSA!"

		// 公钥加密
		encrypted, err := crypto.RSAEncrypt(plaintext, publicKey)
		if err != nil {
			t.Fatalf("RSA加密失败: %v", err)
		}

		if encrypted == "" {
			t.Fatal("加密结果为空")
		}

		// 私钥解密
		decrypted, err := crypto.RSADecrypt(encrypted, privateKey)
		if err != nil {
			t.Fatalf("RSA解密失败: %v", err)
		}

		if decrypted != plaintext {
			t.Fatalf("解密结果不匹配: 期望 %s, 得到 %s", plaintext, decrypted)
		}

		t.Logf("RSA加密解密测试通过")
	})

	t.Run("RSA签名验证", func(t *testing.T) {
		// 生成密钥对
		privateKey, publicKey, err := crypto.GenerateRSAKeyPair(2048)
		if err != nil {
			t.Fatalf("RSA密钥生成失败: %v", err)
		}

		data := "Hello, Signature!"

		// 私钥签名
		signature, err := crypto.RSASign(data, privateKey)
		if err != nil {
			t.Fatalf("RSA签名失败: %v", err)
		}

		if signature == "" {
			t.Fatal("签名结果为空")
		}

		// 公钥验证
		isValid, err := crypto.RSAVerify(data, signature, publicKey)
		if err != nil {
			t.Fatalf("RSA验证失败: %v", err)
		}

		if !isValid {
			t.Fatal("签名验证失败")
		}

		// 测试错误数据的验证
		isValid, err = crypto.RSAVerify("wrong data", signature, publicKey)
		if err != nil {
			t.Fatalf("RSA验证失败: %v", err)
		}

		if isValid {
			t.Fatal("错误数据的签名验证应该失败")
		}

		t.Logf("RSA签名验证测试通过")
	})
}

func TestCryptoHash(t *testing.T) {
	t.Run("基本哈希算法", func(t *testing.T) {
		data := "Hello, Hash!"

		// 测试各种哈希算法
		md5Hash := crypto.MD5(data)
		sha1Hash := crypto.SHA1(data)
		sha256Hash := crypto.SHA256(data)
		sha512Hash := crypto.SHA512(data)

		// 验证哈希长度
		if len(md5Hash) != 32 {
			t.Fatalf("MD5哈希长度不正确: 期望 32, 得到 %d", len(md5Hash))
		}

		if len(sha1Hash) != 40 {
			t.Fatalf("SHA1哈希长度不正确: 期望 40, 得到 %d", len(sha1Hash))
		}

		if len(sha256Hash) != 64 {
			t.Fatalf("SHA256哈希长度不正确: 期望 64, 得到 %d", len(sha256Hash))
		}

		if len(sha512Hash) != 128 {
			t.Fatalf("SHA512哈希长度不正确: 期望 128, 得到 %d", len(sha512Hash))
		}

		// 验证哈希一致性
		md5Hash2 := crypto.MD5(data)
		if md5Hash != md5Hash2 {
			t.Fatal("相同数据的MD5哈希结果不一致")
		}

		t.Logf("基本哈希算法测试通过")
		t.Logf("MD5: %s", md5Hash)
		t.Logf("SHA256: %s", sha256Hash)
	})

	t.Run("HMAC算法", func(t *testing.T) {
		data := "Hello, HMAC!"
		key := "secret-key"

		// 测试HMAC算法
		hmacSHA256 := crypto.HMACSHA256(data, key)
		hmacSHA512 := crypto.HMACSHA512(data, key)

		if len(hmacSHA256) != 64 {
			t.Fatalf("HMAC-SHA256长度不正确: 期望 64, 得到 %d", len(hmacSHA256))
		}

		if len(hmacSHA512) != 128 {
			t.Fatalf("HMAC-SHA512长度不正确: 期望 128, 得到 %d", len(hmacSHA512))
		}

		// 验证HMAC一致性
		hmacSHA256_2 := crypto.HMACSHA256(data, key)
		if hmacSHA256 != hmacSHA256_2 {
			t.Fatal("相同数据和密钥的HMAC结果不一致")
		}

		t.Logf("HMAC算法测试通过")
	})
}

func TestCryptoPassword(t *testing.T) {
	t.Run("密码哈希和验证", func(t *testing.T) {
		password := "my-secure-password"

		// 哈希密码
		hashedPassword, err := crypto.HashPassword(password)
		if err != nil {
			t.Fatalf("密码哈希失败: %v", err)
		}

		if hashedPassword == "" {
			t.Fatal("哈希密码为空")
		}

		// 验证正确密码
		isValid := crypto.CheckPassword(password, hashedPassword)
		if !isValid {
			t.Fatal("正确密码验证失败")
		}

		// 验证错误密码
		isValid = crypto.CheckPassword("wrong-password", hashedPassword)
		if isValid {
			t.Fatal("错误密码验证应该失败")
		}

		t.Logf("密码哈希和验证测试通过")
	})

	t.Run("密码强度检查", func(t *testing.T) {
		testCases := []struct {
			password string
			expected crypto.PasswordStrength
		}{
			{"123", crypto.Weak},
			{"password", crypto.Weak},    // 只有小写字母，强度为弱
			{"Password123", crypto.Fair}, // 有大小写字母和数字，但长度不够长
			{"Password123!", crypto.Strong},
		}

		for _, tc := range testCases {
			strength := crypto.CheckPasswordStrength(tc.password)
			if strength != tc.expected {
				t.Fatalf("密码 %s 强度检查失败: 期望 %s, 得到 %s",
					tc.password, tc.expected.String(), strength.String())
			}
		}

		t.Logf("密码强度检查测试通过")
	})

	t.Run("密码生成", func(t *testing.T) {
		// 生成普通密码
		password, err := crypto.GeneratePassword(12, false)
		if err != nil {
			t.Fatalf("生成密码失败: %v", err)
		}

		if len(password) != 12 {
			t.Fatalf("生成密码长度不正确: 期望 12, 得到 %d", len(password))
		}

		// 生成强密码
		strongPassword, err := crypto.GenerateStrongPassword(16)
		if err != nil {
			t.Fatalf("生成强密码失败: %v", err)
		}

		if len(strongPassword) != 16 {
			t.Fatalf("生成强密码长度不正确: 期望 16, 得到 %d", len(strongPassword))
		}

		// 检查强密码强度
		strength := crypto.CheckPasswordStrength(strongPassword)
		if strength < crypto.Strong {
			t.Fatalf("生成的强密码强度不够: %s", strength.String())
		}

		t.Logf("密码生成测试通过")
		t.Logf("生成的密码: %s", password)
		t.Logf("生成的强密码: %s (强度: %s)", strongPassword, strength.String())
	})
}

func TestCryptoUtils(t *testing.T) {
	t.Run("随机数生成", func(t *testing.T) {
		// 生成随机字节
		randomBytes, err := crypto.GenerateRandomBytes(32)
		if err != nil {
			t.Fatalf("生成随机字节失败: %v", err)
		}

		if len(randomBytes) != 32 {
			t.Fatalf("随机字节长度不正确: 期望 32, 得到 %d", len(randomBytes))
		}

		// 生成随机字符串
		randomString, err := crypto.GenerateRandomString(16)
		if err != nil {
			t.Fatalf("生成随机字符串失败: %v", err)
		}

		if len(randomString) != 16 {
			t.Fatalf("随机字符串长度不正确: 期望 16, 得到 %d", len(randomString))
		}

		t.Logf("随机数生成测试通过")
		t.Logf("随机字符串: %s", randomString)
	})

	t.Run("编码解码", func(t *testing.T) {
		data := []byte("Hello, Encoding!")

		// Base64编码解码
		base64Encoded := crypto.Base64Encode(data)
		base64Decoded, err := crypto.Base64Decode(base64Encoded)
		if err != nil {
			t.Fatalf("Base64解码失败: %v", err)
		}

		if string(base64Decoded) != string(data) {
			t.Fatal("Base64编码解码结果不一致")
		}

		// Hex编码解码
		hexEncoded := crypto.HexEncode(data)
		hexDecoded, err := crypto.HexDecode(hexEncoded)
		if err != nil {
			t.Fatalf("Hex解码失败: %v", err)
		}

		if string(hexDecoded) != string(data) {
			t.Fatal("Hex编码解码结果不一致")
		}

		t.Logf("编码解码测试通过")
	})
}
