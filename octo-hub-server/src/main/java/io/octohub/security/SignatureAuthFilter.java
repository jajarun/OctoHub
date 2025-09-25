package io.octohub.security;

import com.fasterxml.jackson.databind.ObjectMapper;
import io.octohub.annotation.SignatureAuth;
import io.octohub.dto.ApiResponse;
import io.octohub.enums.ErrorCode;
import io.octohub.util.ResponseUtil;
import io.octohub.util.SignatureUtils;
import jakarta.servlet.FilterChain;
import jakarta.servlet.ServletException;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.MediaType;
import org.springframework.lang.NonNull;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.authority.SimpleGrantedAuthority;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.util.StringUtils;
import org.springframework.web.filter.OncePerRequestFilter;
import org.springframework.web.method.HandlerMethod;
import org.springframework.web.servlet.HandlerExecutionChain;
import org.springframework.web.servlet.mvc.method.annotation.RequestMappingHandlerMapping;

import java.io.IOException;
import java.util.Collections;
import java.util.HashMap;
import java.util.Map;

public class SignatureAuthFilter extends OncePerRequestFilter {
    
    private static final Logger logger = LoggerFactory.getLogger(SignatureAuthFilter.class);
    
    @Autowired
    private SignatureUtils signatureUtils;
    
    @Autowired
    private RequestMappingHandlerMapping handlerMapping;
    
    private final ObjectMapper objectMapper = new ObjectMapper();
    
    @Override
    protected void doFilterInternal(@NonNull HttpServletRequest request, @NonNull HttpServletResponse response, 
                                  @NonNull FilterChain filterChain) throws ServletException, IOException {
        try {
            // 检查当前请求是否需要签名验证
            if (requiresSignatureAuth(request)) {
                logger.debug("Request requires signature authentication: {}", request.getRequestURI());
                
                if (!validateSignature(request)) {
                    logger.warn("Signature validation failed for request: {}", request.getRequestURI());
                    sendErrorResponse(response, ErrorCode.SIGNATURE_INVALID);
                    return;
                }
                
                // 签名验证成功，设置一个匿名认证，表示已通过签名验证
                UsernamePasswordAuthenticationToken authentication = 
                    new UsernamePasswordAuthenticationToken("signature_auth", null, 
                        Collections.singletonList(new SimpleGrantedAuthority("SIGNATURE_AUTH")));
                SecurityContextHolder.getContext().setAuthentication(authentication);
                logger.debug("Signature authentication successful for request: {}", request.getRequestURI());
            }
        } catch (Exception e) {
            logger.error("Error during signature authentication: {}", e.getMessage(), e);
            sendErrorResponse(response, ErrorCode.SIGNATURE_INVALID);
            return;
        }
        
        filterChain.doFilter(request, response);
    }
    
    /**
     * 检查请求是否需要签名验证
     */
    private boolean requiresSignatureAuth(HttpServletRequest request) {
        try {
            HandlerExecutionChain handlerChain = handlerMapping.getHandler(request);
            if (handlerChain != null && handlerChain.getHandler() instanceof HandlerMethod) {
                HandlerMethod handlerMethod = (HandlerMethod) handlerChain.getHandler();
                
                // 检查方法级别的注解
                SignatureAuth methodAnnotation = handlerMethod.getMethodAnnotation(SignatureAuth.class);
                if (methodAnnotation != null) {
                    return methodAnnotation.required();
                }
                
                // 检查类级别的注解
                SignatureAuth classAnnotation = handlerMethod.getBeanType().getAnnotation(SignatureAuth.class);
                if (classAnnotation != null) {
                    return classAnnotation.required();
                }
            }
        } catch (Exception e) {
            logger.debug("Could not determine handler for request: {}", request.getRequestURI());
        }
        
        return false;
    }
    
    /**
     * 验证签名
     */
    private boolean validateSignature(HttpServletRequest request) {
        String signature = request.getHeader("X-Signature");
        String timestamp = request.getHeader("X-Timestamp");
        String nonce = request.getHeader("X-Nonce");
        
        if (!StringUtils.hasText(signature) || !StringUtils.hasText(timestamp) || !StringUtils.hasText(nonce)) {
            logger.debug("Missing signature headers");
            return false;
        }
        
        // 获取请求参数
        Map<String, String> params = new HashMap<>();
        request.getParameterMap().forEach((key, values) -> {
            if (values.length > 0) {
                params.put(key, values[0]);
            }
        });
        
        String method = request.getMethod();
        String uri = request.getRequestURI();
        
        return signatureUtils.validateSignature(method, uri, params, timestamp, nonce, signature);
    }
    
    /**
     * 发送错误响应
     */
    private void sendErrorResponse(HttpServletResponse response, ErrorCode errorCode) throws IOException {
        response.setStatus(HttpServletResponse.SC_UNAUTHORIZED);
        response.setContentType(MediaType.APPLICATION_JSON_VALUE);
        response.setCharacterEncoding("UTF-8");
        
        ApiResponse<Object> apiResponse = ResponseUtil.error(errorCode).getBody();
        String jsonResponse = objectMapper.writeValueAsString(apiResponse);
        
        response.getWriter().write(jsonResponse);
        response.getWriter().flush();
    }
}
