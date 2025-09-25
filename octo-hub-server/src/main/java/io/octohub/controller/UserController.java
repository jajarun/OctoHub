package io.octohub.controller;

import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.http.ResponseEntity;
import org.springframework.beans.factory.annotation.Autowired;
import io.octohub.dto.ApiResponse;
import io.octohub.dto.WebSocketConnectionDto;
import io.octohub.entity.User;
import io.octohub.enums.ErrorCode;
import io.octohub.util.ResponseUtil;
import io.octohub.service.UserService;
import io.octohub.service.WebSocketService;

@RestController
@RequestMapping("/user")
public class UserController {

    @Autowired
    private UserService userService;

    @Autowired
    private WebSocketService webSocketService;

    @GetMapping("/info")
    public ResponseEntity<ApiResponse<User>> getUserInfo() {
        return ResponseUtil.success(userService.getUserInfo());
    }
    
    /**
     * 获取用户WebSocket连接地址
     * @param userId 用户ID，如果不传则使用默认值
     * @param type 连接类型：user(默认) 或 pc
     * @return WebSocket连接信息
     */
    @GetMapping("/ws")
    public ResponseEntity<ApiResponse<WebSocketConnectionDto>> getWsAddress() {
        
        try {
            WebSocketConnectionDto connectionInfo;
            User user = userService.getUserInfo();
            connectionInfo = webSocketService.generateUserConnectionUrl(user.getId().toString());
            
            return ResponseUtil.success(connectionInfo);
            
        } catch (Exception e) {
            return ResponseUtil.error(ErrorCode.SYSTEM_ERROR, "生成WebSocket连接地址失败: " + e.getMessage());
        }
    }
    
}
