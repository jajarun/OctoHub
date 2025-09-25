package io.octohub.util;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;

import javax.crypto.Mac;
import javax.crypto.spec.SecretKeySpec;
import java.nio.charset.StandardCharsets;
import java.security.InvalidKeyException;
import java.security.NoSuchAlgorithmException;

@Component
public class WebSocketSignatureUtils {
    
    private static final String HMAC_SHA256 = "HmacSHA256";

    @Value("${websocket.signature.secret.key:your-secret-key-here}")
    private String secretKey;

    /**
     * 生成WebSocket连接签名
     * @param id 用户ID或PC ID
     * @param timestamp 时间戳
     * @return 签名字符串 (hex编码)
     */
    public String generateSignature(String id, String timestamp) {
        String message = id + timestamp;
        
        try {
            Mac mac = Mac.getInstance(HMAC_SHA256);
            SecretKeySpec secretKeySpec = new SecretKeySpec(secretKey.getBytes(StandardCharsets.UTF_8), HMAC_SHA256);
            mac.init(secretKeySpec);
            
            byte[] hash = mac.doFinal(message.getBytes(StandardCharsets.UTF_8));
            return bytesToHex(hash);
        } catch (NoSuchAlgorithmException | InvalidKeyException e) {
            throw new RuntimeException("WebSocket签名生成失败", e);
        }
    }

    /**
     * 验证WebSocket签名
     * @param id 用户ID或PC ID
     * @param timestamp 时间戳
     * @param signature 签名
     * @return 验证结果
     */
    public boolean validateSignature(String id, String timestamp, String signature) {
        try {
            // 检查时间戳，防止重放攻击（5分钟内有效）
            long currentTime = System.currentTimeMillis() / 1000;
            long requestTime = Long.parseLong(timestamp);
            if (Math.abs(currentTime - requestTime) > 300) { // 5分钟
                return false;
            }
            
            String expectedSignature = generateSignature(id, timestamp);
            return signature.equals(expectedSignature);
        } catch (Exception e) {
            return false;
        }
    }

    /**
     * 字节数组转十六进制字符串
     */
    private String bytesToHex(byte[] bytes) {
        StringBuilder result = new StringBuilder();
        for (byte b : bytes) {
            result.append(String.format("%02x", b));
        }
        return result.toString();
    }
}
