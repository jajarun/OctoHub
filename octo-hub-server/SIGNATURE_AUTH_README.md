# 签名验证功能使用说明

## 概述

本系统现在支持两种认证方式：
1. **JWT Token认证**：适用于需要用户身份验证的接口
2. **签名验证**：适用于开放API或第三方系统调用的接口

## 签名验证机制

### 签名算法
- 使用 HMAC-SHA256 算法
- 签名字符串格式：`METHOD&URI&PARAMS&TIMESTAMP&NONCE`
- 参数按key字典序排序

### 请求头要求
客户端需要在请求头中包含以下信息：
```
X-Signature: 签名值（Base64编码）
X-Timestamp: 时间戳（秒级，用于防重放攻击）
X-Nonce: 随机数（防重放攻击）
```

### 安全特性
- **时间窗口验证**：请求时间戳必须在5分钟内有效
- **防重放攻击**：使用时间戳和随机数
- **参数完整性**：所有参数参与签名计算

## 使用方法

### 1. 在Controller中标记需要签名验证的接口

```java
@RestController
public class ApiController {
    
    // 使用签名验证的接口
    @GetMapping("/api/public-data")
    @SignatureAuth
    public ResponseEntity<ApiResponse<String>> getPublicData() {
        return ResponseUtil.success("公开数据");
    }
    
    // 整个Controller都使用签名验证
    @SignatureAuth
    @RestController
    public class PublicApiController {
        // 所有接口都需要签名验证
    }
}
```

### 2. 客户端生成签名

使用提供的 `SignatureClient` 工具类：

```java
SignatureClient client = new SignatureClient();
Map<String, String> params = new HashMap<>();
// 添加查询参数（如果有）
params.put("param1", "value1");

Map<String, String> headers = client.generateHeaders("GET", "/api/public-data", params);
// 使用headers发送HTTP请求
```

### 3. 测试示例

运行 `SignatureClient.main()` 方法可以生成测试用的curl命令：

```bash
curl -X GET \
  -H "X-Signature: xxx" \
  -H "X-Timestamp: 1234567890" \
  -H "X-Nonce: xxx" \
  -H "Content-Type: application/json" \
  http://localhost:8080/user/public-info
```

## 配置说明

### 密钥配置
当前密钥硬编码在 `SignatureUtils` 中，生产环境应该：
1. 将密钥移至配置文件 `application.properties`
2. 使用环境变量或配置中心管理密钥
3. 支持多个密钥和密钥轮换

```properties
# application.properties
signature.secret.key=your-production-secret-key
```

### 时间窗口配置
默认时间窗口为5分钟，可以根据需要调整：

```java
// 在SignatureUtils中修改
private static final long TIME_WINDOW = 300; // 秒
```

## 错误码

| 错误码 | 错误信息 | 说明 |
|--------|----------|------|
| 1002 | 签名验证失败 | 签名不正确或签名已过期 |
| 1003 | 签名已过期 | 请求时间戳超出允许范围 |
| 1004 | 缺少签名信息 | 缺少必要的签名头信息 |

## 接口分类

- **JWT Token接口**：`/user/info` 等需要用户身份的接口
- **签名验证接口**：`/user/public-info` 等开放API接口
- **公开接口**：`/login/**` 等不需要任何验证的接口

## 注意事项

1. **密钥安全**：生产环境必须使用安全的密钥管理方案
2. **时间同步**：客户端和服务端时间必须同步
3. **HTTPS**：生产环境建议使用HTTPS传输
4. **日志记录**：建议记录签名验证失败的请求用于安全审计
