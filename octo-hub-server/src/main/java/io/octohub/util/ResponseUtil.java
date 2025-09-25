package io.octohub.util;

import io.octohub.dto.ApiResponse;
import io.octohub.enums.ErrorCode;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;

/**
 * 响应工具类
 * 提供便捷的ResponseEntity创建方法
 */
public class ResponseUtil {
    
    /**
     * 成功响应，带数据
     */
    public static <T> ResponseEntity<ApiResponse<T>> success(T data) {
        return ResponseEntity.ok(ApiResponse.success(data));
    }
    
    /**
     * 成功响应，无数据
     */
    public static <T> ResponseEntity<ApiResponse<T>> success() {
        return ResponseEntity.ok(ApiResponse.success());
    }
    
    /**
     * 成功响应，带消息和数据
     */
    public static <T> ResponseEntity<ApiResponse<T>> success(String message, T data) {
        return ResponseEntity.ok(ApiResponse.success(message, data));
    }
    
    /**
     * 使用错误码枚举的失败响应
     */
    public static <T> ResponseEntity<ApiResponse<T>> error(ErrorCode errorCode) {
        return ResponseEntity.badRequest().body(ApiResponse.error(errorCode));
    }
    
    /**
     * 使用错误码枚举的失败响应，自定义HTTP状态码
     */
    public static <T> ResponseEntity<ApiResponse<T>> error(HttpStatus status, ErrorCode errorCode) {
        return ResponseEntity.status(status).body(ApiResponse.error(errorCode));
    }
    
    /**
     * 使用错误码枚举的失败响应，带数据
     */
    public static <T> ResponseEntity<ApiResponse<T>> error(ErrorCode errorCode, T data) {
        return ResponseEntity.badRequest().body(ApiResponse.error(errorCode, data));
    }
    
    /**
     * 使用错误码枚举的失败响应，自定义消息
     */
    public static <T> ResponseEntity<ApiResponse<T>> error(ErrorCode errorCode, String customMessage) {
        return ResponseEntity.badRequest().body(ApiResponse.error(errorCode, customMessage));
    }
    
    /**
     * 使用错误码枚举的失败响应，自定义HTTP状态码和消息
     */
    public static <T> ResponseEntity<ApiResponse<T>> error(HttpStatus status, ErrorCode errorCode, String customMessage) {
        return ResponseEntity.status(status).body(ApiResponse.error(errorCode, customMessage));
    }
}
