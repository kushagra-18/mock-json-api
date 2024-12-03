package com.mock_json.mock_api.models;

import lombok.*;
import org.springframework.data.annotation.CreatedDate;
import org.springframework.data.annotation.Id;
import org.springframework.data.jpa.domain.support.AuditingEntityListener;
import org.springframework.data.mongodb.core.mapping.Document;
import org.springframework.data.mongodb.core.mapping.Field;

import jakarta.persistence.EntityListeners;

import java.time.LocalDateTime;

@Document(collection = "request_logs") 
@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
@Builder
@EntityListeners(AuditingEntityListener.class)
public class RequestLog{

    @Id
    private String id;

    private String ip;

    @CreatedDate 
    @Field("created_at") 
    private LocalDateTime createdAt;

    @Field("url_id") 
    private Long urlId;

    @Field("project_id") 
    private Long projectId;

    @Field("method") 
    private String method;

    @Field("status")    
    private int status; 

    @Field("url") 
    private String url;
}
