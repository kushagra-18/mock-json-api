package com.mock_json.api;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.scheduling.annotation.EnableAsync;

@SpringBootApplication
@EnableAsync
public class MockJsonApplication {

	public static void main(String[] args) {
		SpringApplication.run(MockJsonApplication.class, args);
	}

}
