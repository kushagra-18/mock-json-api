package com.mock_json;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.annotation.Configuration;
import org.springframework.data.mongodb.config.EnableMongoAuditing;
import org.springframework.scheduling.annotation.EnableAsync;

@SpringBootApplication
@EnableAsync
@Configuration
@EnableMongoAuditing
public class MockJsonApplication {

	public static void main(String[] args) {
		SpringApplication.run(MockJsonApplication.class, args);
	}

}
