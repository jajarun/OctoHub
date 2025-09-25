package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"time"
)

// SignatureValidator 签名验证器
type SignatureValidator struct {
	secretKey string
}

// NewSignatureValidator 创建新的签名验证器
func NewSignatureValidator(secretKey string) *SignatureValidator {
	return &SignatureValidator{
		secretKey: secretKey,
	}
}

// ValidateSignature 验证签名
func (sv *SignatureValidator) ValidateSignature(id, timestamp, signature string, timeoutSeconds int) bool {
	// 生成期望的签名
	message := id + timestamp
	mac := hmac.New(sha256.New, []byte(sv.secretKey))
	mac.Write([]byte(message))
	expectedSignature := hex.EncodeToString(mac.Sum(nil))

	// 验证签名是否匹配
	if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
		return false
	}

	// 验证时间戳
	ts, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return false
	}

	now := time.Now().Unix()
	if now-ts > int64(timeoutSeconds) {
		return false
	}

	return true
}

// GenerateSignature 生成签名（用于测试）
func (sv *SignatureValidator) GenerateSignature(id string) (string, string) {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	message := id + timestamp

	mac := hmac.New(sha256.New, []byte(sv.secretKey))
	mac.Write([]byte(message))
	signature := hex.EncodeToString(mac.Sum(nil))

	return timestamp, signature
}
