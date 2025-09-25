package io.octohub.util;

import java.util.HashMap;
import java.util.Map;

/**
 * 签名验证客户端示例
 * 用于生成请求签名的工具类
 */
public class SignatureClient {
    
    private final SignatureUtils signatureUtils;
    
    public SignatureClient() {
        this.signatureUtils = new SignatureUtils();
        this.signatureUtils.setSecretKey("mySignatureKey123456789012345");
    }
    
    /**
     * 生成请求头信息，包含签名验证所需的所有头部
     * @param method HTTP方法
     * @param uri 请求URI
     * @param params 请求参数
     * @return 包含签名信息的请求头Map
     */
    public Map<String, String> generateHeaders(String method, String uri, Map<String, String> params) {
        Map<String, String> headers = new HashMap<>();
        
        // 生成时间戳和随机数
        String timestamp = String.valueOf(System.currentTimeMillis() / 1000);
        String nonce = generateNonce();
        
        // 生成签名
        String signature = signatureUtils.generateSignature(method, uri, params, timestamp, nonce);
        
        // 设置请求头
        headers.put("X-Signature", signature);
        headers.put("X-Timestamp", timestamp);
        headers.put("X-Nonce", nonce);
        headers.put("Content-Type", "application/json");
        
        return headers;
    }
    
    /**
     * 生成随机数
     */
    private String generateNonce() {
        return String.valueOf(System.currentTimeMillis() + (int)(Math.random() * 1000));
    }
    
    /**
     * 使用示例
     */
    public static void main(String[] args) {
        SignatureClient client = new SignatureClient();
        
        // 示例：为GET /user/public-info请求生成签名头
        Map<String, String> params = new HashMap<>();
        // 如果有查询参数，在这里添加
        // params.put("param1", "value1");
        
        Map<String, String> headers = client.generateHeaders("GET", "/user/public-info", params);
        
        System.out.println("请求头信息：");
        headers.forEach((key, value) -> System.out.println(key + ": " + value));
        
        System.out.println("\n使用curl测试命令：");
        StringBuilder curlCommand = new StringBuilder("curl -X GET");
        headers.forEach((key, value) -> curlCommand.append(" -H \"").append(key).append(": ").append(value).append("\""));
        curlCommand.append(" http://localhost:8080/user/public-info");
        System.out.println(curlCommand.toString());
    }
}
