package com.mock_json.api.models;

import jakarta.persistence.*;
import lombok.*;
import lombok.Builder.Default;

import org.hibernate.annotations.SoftDelete;
import org.hibernate.annotations.UpdateTimestamp;

import com.fasterxml.jackson.annotation.JsonBackReference;
import com.fasterxml.jackson.annotation.JsonProperty;

import org.hibernate.annotations.CreationTimestamp;

import java.time.LocalDateTime;

@Entity
@Table(name = "urls")
@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
@Builder
@SoftDelete
public class Url {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @Builder.Default 
    @Column(nullable = false)
    private Long latency = 0L; 

    @Column(nullable = true)
    private String description;

    @Column(nullable = false)
    private String name;

    @Lob
    @JsonProperty("json_data") 
    @Column(nullable = false,columnDefinition = "LONGTEXT")
    private String jsonData;

    @Column(nullable = false, unique = true)
    private String url;

    @CreationTimestamp
    @Column(nullable = false, updatable = false)
    private LocalDateTime createdAt;

    @UpdateTimestamp
    @Column(nullable = false)
    private LocalDateTime updatedAt;

    @ManyToOne(fetch = FetchType.EAGER)
    @JoinColumn(name = "project_id")
    @JsonBackReference
    private Project project;

}
