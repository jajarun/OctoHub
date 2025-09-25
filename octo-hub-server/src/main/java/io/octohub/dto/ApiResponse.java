package io.octohub.dto;

import com.fasterxml.jackson.annotation.JsonInclude;
import io.octohub.enums.ErrorCode;

/**
 * 统一API响应格式
 * @param <T> 数据类型
 */
@JsonInclude(JsonInclude.Include.NON_NULL)
public class ApiResponse<T> {
    
    /**
     * 错误码，0表示成功，非0表示失败
     */
    private Integer errcode;
    
    /**
     * 错误信息，成功时可为空
     */
    private String errmsg;
    
    /**
     * 响应数据
     */
    private T data;
    
    /**
     * 私有构造函数
     */
    private ApiResponse() {}
    
    /**
     * 私有构造函数
     */
    private ApiResponse(Integer errcode, String errmsg, T data) {
        this.errcode = errcode;
        this.errmsg = errmsg;
        this.data = data;
    }
    
    /**
     * 成功响应，带数据
     */
    public static <T> ApiResponse<T> success(T data) {
        return new ApiResponse<>(0, null, data);
    }
    
    /**
     * 成功响应，无数据
     */
    public static <T> ApiResponse<T> success() {
        return new ApiResponse<>(0, null, null);
    }
    
    /**
     * 成功响应，带消息和数据
     */
    public static <T> ApiResponse<T> success(String message, T data) {
        return new ApiResponse<>(0, message, data);
    }
    
    /**
     * 失败响应
     */
    public static <T> ApiResponse<T> error(Integer errcode, String errmsg) {
        return new ApiResponse<>(errcode, errmsg, null);
    }
    
    /**
     * 失败响应，带数据
     */
    public static <T> ApiResponse<T> error(Integer errcode, String errmsg, T data) {
        return new ApiResponse<>(errcode, errmsg, data);
    }
    
    /**
     * 使用错误码枚举的失败响应
     */
    public static <T> ApiResponse<T> error(ErrorCode errorCode) {
        return new ApiResponse<>(errorCode.getCode(), errorCode.getMessage(), null);
    }
    
    /**
     * 使用错误码枚举的失败响应，带数据
     */
    public static <T> ApiResponse<T> error(ErrorCode errorCode, T data) {
        return new ApiResponse<>(errorCode.getCode(), errorCode.getMessage(), data);
    }
    
    /**
     * 使用错误码枚举的失败响应，自定义消息
     */
    public static <T> ApiResponse<T> error(ErrorCode errorCode, String customMessage) {
        return new ApiResponse<>(errorCode.getCode(), customMessage, null);
    }
    
    // Getters and Setters
    public Integer getErrcode() {
        return errcode;
    }
    
    public void setErrcode(Integer errcode) {
        this.errcode = errcode;
    }
    
    public String getErrmsg() {
        return errmsg;
    }
    
    public void setErrmsg(String errmsg) {
        this.errmsg = errmsg;
    }
    
    public T getData() {
        return data;
    }
    
    public void setData(T data) {
        this.data = data;
    }
    
    @Override
    public String toString() {
        return "ApiResponse{" +
                "errcode=" + errcode +
                ", errmsg='" + errmsg + '\'' +
                ", data=" + data +
                '}';
    }
}
