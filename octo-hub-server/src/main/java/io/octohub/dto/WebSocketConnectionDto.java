package io.octohub.dto;

public class WebSocketConnectionDto {
    private String wsUrl;

    public WebSocketConnectionDto() {}

    public WebSocketConnectionDto(String wsUrl) {
        this.wsUrl = wsUrl;
    }

    public String getWsUrl() {
        return wsUrl;
    }

    public void setWsUrl(String wsUrl) {
        this.wsUrl = wsUrl;
    }

    @Override
    public String toString() {
        return "WebSocketConnectionDto{" +
                "wsUrl='" + wsUrl + '\'' +
                '}';
    }
}
