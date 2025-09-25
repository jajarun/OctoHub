package io.octohub.util;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;
import org.springframework.util.StringUtils;

import javax.crypto.Mac;
import javax.crypto.spec.SecretKeySpec;
import java.nio.charset.StandardCharsets;
import java.security.InvalidKeyException;
import java.security.NoSuchAlgorithmException;
import java.util.Base64;
import java.util.Map;
import java.util.TreeMap;

@Component
public class SignatureUtils {
    
    private static final String HMAC_SHA256 = "HmacSHA256";

    @Value("${signature.secret.key}")
    private String secretKey; // 实际使用时应该从配置文件读取

    public void setSecretKey(String secretKey) {
        this.secretKey = secretKey;
    }
    
    /**
     * 生成签名
     * @param method HTTP方法
     * @param uri 请求URI
     * @param params 请求参数
     * @param timestamp 时间戳
     * @param nonce 随机数
     * @return 签名字符串
     */
    public String generateSignature(String method, String uri, Map<String, String> params, 
                                  String timestamp, String nonce) {
        // 构建签名字符串
        String signString = buildSignString(method, uri, params, timestamp, nonce);
        
        try {
            Mac mac = Mac.getInstance(HMAC_SHA256);
            SecretKeySpec secretKeySpec = new SecretKeySpec(secretKey.getBytes(StandardCharsets.UTF_8), HMAC_SHA256);
            mac.init(secretKeySpec);
            
            byte[] hash = mac.doFinal(signString.getBytes(StandardCharsets.UTF_8));
            return Base64.getEncoder().encodeToString(hash);
        } catch (NoSuchAlgorithmException | InvalidKeyException e) {
            throw new RuntimeException("签名生成失败", e);
        }
    }
    
    /**
     * 验证签名
     * @param method HTTP方法
     * @param uri 请求URI
     * @param params 请求参数
     * @param timestamp 时间戳
     * @param nonce 随机数
     * @param signature 待验证的签名
     * @return 验证结果
     */
    public boolean validateSignature(String method, String uri, Map<String, String> params, 
                                   String timestamp, String nonce, String signature) {
        if (!StringUtils.hasText(signature) || !StringUtils.hasText(timestamp) || !StringUtils.hasText(nonce)) {
            return false;
        }
        
        // 检查时间戳，防止重放攻击（5分钟内有效）
        long currentTime = System.currentTimeMillis() / 1000;
        long requestTime = Long.parseLong(timestamp);
        if (Math.abs(currentTime - requestTime) > 300) { // 5分钟
            return false;
        }
        
        String expectedSignature = generateSignature(method, uri, params, timestamp, nonce);
        return signature.equals(expectedSignature);
    }
    
    /**
     * 构建签名字符串
     * 格式: METHOD&URI&PARAMS&TIMESTAMP&NONCE
     */
    private String buildSignString(String method, String uri, Map<String, String> params, 
                                 String timestamp, String nonce) {
        StringBuilder sb = new StringBuilder();
        
        // HTTP方法
        sb.append(method.toUpperCase()).append("&");
        
        // URI
        sb.append(uri).append("&");
        
        // 参数按key排序后拼接
        if (params != null && !params.isEmpty()) {
            TreeMap<String, String> sortedParams = new TreeMap<>(params);
            for (Map.Entry<String, String> entry : sortedParams.entrySet()) {
                sb.append(entry.getKey()).append("=").append(entry.getValue()).append("&");
            }
            // 移除最后一个&
            sb.deleteCharAt(sb.length() - 1);
        }
        
        sb.append("&");
        
        // 时间戳
        sb.append(timestamp).append("&");
        
        // 随机数
        sb.append(nonce);
        
        return sb.toString();
    }
}
