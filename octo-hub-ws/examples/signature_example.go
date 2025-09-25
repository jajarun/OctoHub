package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"
)

// 生成签名的示例函数
func generateSignature(id, secretKey string) (string, string, string) {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	message := id + timestamp

	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(message))
	signature := hex.EncodeToString(mac.Sum(nil))

	return timestamp, signature, message
}

func main() {
	secretKey := "your-secret-key-here"

	// 生成用户连接签名示例
	userID := "user123"
	timestamp, signature, message := generateSignature(userID, secretKey)

	fmt.Println("=== 用户连接签名示例 ===")
	fmt.Printf("用户ID: %s\n", userID)
	fmt.Printf("时间戳: %s\n", timestamp)
	fmt.Printf("签名消息: %s\n", message)
	fmt.Printf("签名: %s\n", signature)
	fmt.Printf("用户连接URL: ws://localhost:8080/ws/user?user_id=%s&timestamp=%s&signature=%s\n\n", userID, timestamp, signature)

	// 生成PC连接签名示例
	pcID := "pc456"
	timestamp2, signature2, message2 := generateSignature(pcID, secretKey)

	fmt.Println("=== PC连接签名示例 ===")
	fmt.Printf("PC ID: %s\n", pcID)
	fmt.Printf("时间戳: %s\n", timestamp2)
	fmt.Printf("签名消息: %s\n", message2)
	fmt.Printf("签名: %s\n", signature2)
	fmt.Printf("PC连接URL: ws://localhost:8080/ws/pc?pc_id=%s&timestamp=%s&signature=%s\n", pcID, timestamp2, signature2)
}
