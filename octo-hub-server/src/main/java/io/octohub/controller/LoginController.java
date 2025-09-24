package io.octohub.controller;

import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RestController;
import io.octohub.dto.LoginRequest;
import java.util.Map;
import java.util.HashMap;
import org.springframework.http.ResponseEntity;


@RestController
@RequestMapping("/api/login")
public class LoginController {

	@PostMapping("/login")
	ResponseEntity<?> login(@RequestBody LoginRequest request) {
		Map<String, Object> response = new HashMap<>();
		response.put("errcode", 0);
		return ResponseEntity.ok(response);
	}

    @RequestMapping("/test")
	String test() {
		return "OctoHub Service";
	}

}
