package io.octohub.enums;

/**
 * 错误码枚举
 * 统一管理所有系统错误码
 */
public enum ErrorCode {
    
    // 成功
    SUCCESS(0, "成功"),
    
    //通用错误
    ERROR_LOGIN(401, "未登录"),
    LOGIN_FAILED(1001, "登录失败"),
    SYSTEM_ERROR(1000, "系统内部错误");
    
    
    private final Integer code;
    private final String message;
    
    ErrorCode(Integer code, String message) {
        this.code = code;
        this.message = message;
    }
    
    public Integer getCode() {
        return code;
    }
    
    public String getMessage() {
        return message;
    }
    
    /**
     * 根据错误码获取枚举
     */
    public static ErrorCode fromCode(Integer code) {
        for (ErrorCode errorCode : ErrorCode.values()) {
            if (errorCode.getCode().equals(code)) {
                return errorCode;
            }
        }
        return SYSTEM_ERROR;
    }
    
    @Override
    public String toString() {
        return "ErrorCode{" +
                "code=" + code +
                ", message='" + message + '\'' +
                '}';
    }
}
