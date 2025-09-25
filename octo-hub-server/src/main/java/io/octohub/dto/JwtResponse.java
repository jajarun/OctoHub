package io.octohub.dto;

public class JwtResponse {
    private String token;
    private String type = "Bearer";
    private Long id;
    
    public JwtResponse(String accessToken, Long id) {
        this.token = accessToken;
        this.id = id;
    }
    
    // Getters and Setters
    public String getToken() {
        return token;
    }
    
    public void setToken(String token) {
        this.token = token;
    }
    
    public String getType() {
        return type;
    }
    
    public void setType(String type) {
        this.type = type;
    }
    
    public Long getId() {
        return id;
    }
    
    public void setId(Long id) {
        this.id = id;
    }

} 