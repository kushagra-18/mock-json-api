package com.mock_json.mock_api.models;

import lombok.*;
import org.springframework.data.annotation.Id;
import org.springframework.data.jpa.domain.support.AuditingEntityListener;
import org.springframework.data.mongodb.core.mapping.Document;
import org.springframework.data.mongodb.core.mapping.Field;

import jakarta.persistence.EntityListeners;

@Document(collection = "random_words") 
@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
@Builder
@EntityListeners(AuditingEntityListener.class)
public class RandomWords{

    @Id
    private String id;

    @Field("word") 
    private String word;

}
