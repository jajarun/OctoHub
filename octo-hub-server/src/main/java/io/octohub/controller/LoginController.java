package io.octohub.controller;

import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;


@RestController
@RequestMapping("/login")
public class LoginController {

    @RequestMapping("/test")
	String test() {
		return "OctoHub Service";
	}

}
