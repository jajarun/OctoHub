package io.octohub.service;

import io.octohub.dto.WebSocketConnectionDto;
import io.octohub.util.WebSocketSignatureUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;

@Service
public class WebSocketService {

    @Autowired
    private WebSocketSignatureUtils signatureUtils;

    @Value("${websocket.server.host:localhost}")
    private String wsHost;

    @Value("${websocket.server.port:8000}")
    private String wsPort;

    @Value("${websocket.server.protocol:ws}")
    private String wsProtocol;

    /**
     * 生成用户WebSocket连接地址
     * @param userId 用户ID
     * @return WebSocket连接信息
     */
    public WebSocketConnectionDto generateUserConnectionUrl(String userId) {
        String timestamp = String.valueOf(System.currentTimeMillis() / 1000);
        String signature = signatureUtils.generateSignature(userId, timestamp);
        
        String wsUrl = String.format("%s://%s:%s/ws/user?user_id=%s&timestamp=%s&signature=%s",
                wsProtocol, wsHost, wsPort, userId, timestamp, signature);
        
        return new WebSocketConnectionDto(wsUrl);
    }

    /**
     * 验证WebSocket连接签名
     * @param id 用户ID或PC ID
     * @param timestamp 时间戳
     * @param signature 签名
     * @return 验证结果
     */
    public boolean validateConnectionSignature(String id, String timestamp, String signature) {
        return signatureUtils.validateSignature(id, timestamp, signature);
    }
}
