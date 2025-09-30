package io.octohub.controller;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.RestController;
import io.octohub.annotation.SignatureAuth;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.http.ResponseEntity;
import io.octohub.dto.ApiResponse;
import io.octohub.dto.WebSocketConnectionDto;
import io.octohub.enums.ErrorCode;
import io.octohub.util.ResponseUtil;
import io.octohub.service.WebSocketService;
import org.springframework.web.bind.annotation.RequestParam;

@RestController 
@RequestMapping("/node")
@SignatureAuth
public class NodeController {

    @Autowired
    private WebSocketService webSocketService;
    
    /**
     * 获取Node节点WebSocket连接地址
     * @param pcId Node节点ID
     * @return WebSocket连接信息
     */
    @GetMapping("/ws")
    public ResponseEntity<ApiResponse<WebSocketConnectionDto>> getWsAddress(
        @RequestParam("pc_id") String pcId
    ) {
        
        try {
            WebSocketConnectionDto connectionInfo;
            connectionInfo = webSocketService.generateNodeConnectionUrl(pcId);
            
            return ResponseUtil.success(connectionInfo);
        } catch (Exception e) {
            return ResponseUtil.error(ErrorCode.SYSTEM_ERROR, "生成WebSocket连接地址失败: " + e.getMessage());
        }
    }

}
