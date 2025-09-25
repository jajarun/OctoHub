package io.octohub.controller;

import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.http.ResponseEntity;
import org.springframework.beans.factory.annotation.Autowired;
import io.octohub.dto.ApiResponse;
import io.octohub.entity.User;
import io.octohub.util.ResponseUtil;
import io.octohub.service.UserService;

@RestController
@RequestMapping("/user")
public class UserController {

    @Autowired
    private UserService userService;

    @GetMapping("/info")
    public ResponseEntity<ApiResponse<User>> getUserInfo() {
        return ResponseUtil.success(userService.getUserInfo());
    }
    
}
