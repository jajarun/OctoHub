package io.octohub.controller;

import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RestController;
import io.octohub.dto.LoginRequest;
import io.octohub.dto.ApiResponse;
import io.octohub.dto.JwtResponse;
import io.octohub.util.ResponseUtil;
import org.springframework.http.ResponseEntity;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.authentication.AuthenticationManager;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import io.octohub.service.UserDetailsServiceImpl.UserPrincipal;
import io.octohub.util.JwtUtils;
import io.octohub.enums.ErrorCode;

@RestController
@RequestMapping("/login")
public class LoginController {

	@Autowired
    private AuthenticationManager authenticationManager;

	@Autowired
	private JwtUtils jwtUtils;

	@PostMapping("/login")
	ResponseEntity<ApiResponse<JwtResponse>> login(@RequestBody LoginRequest loginRequest) {
		try {
            Authentication authentication = authenticationManager
                    .authenticate(new UsernamePasswordAuthenticationToken(
                            loginRequest.getEmail(),
                            loginRequest.getPassword()));
            
            SecurityContextHolder.getContext().setAuthentication(authentication);
            
            UserPrincipal userDetails = (UserPrincipal) authentication.getPrincipal();
            
			String jwt = jwtUtils.generateJwtToken(userDetails.getId());

            return ResponseUtil.success(new JwtResponse(jwt,
                    userDetails.getId()));
        } catch (Exception e) {
            return ResponseUtil.error(ErrorCode.LOGIN_FAILED);
        }
	}

}
