package com.mock_json.mock_api.models;

import jakarta.persistence.*;
import lombok.*;

import org.hibernate.annotations.SoftDelete;
import org.hibernate.annotations.UpdateTimestamp;

import com.fasterxml.jackson.annotation.JsonBackReference;
import com.fasterxml.jackson.annotation.JsonIgnore;
import com.fasterxml.jackson.annotation.JsonProperty;

import org.hibernate.annotations.CreationTimestamp;

import java.time.LocalDateTime;

@Entity
@Table(name = "mock_content")
@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
@Builder
@SoftDelete
public class MockContent {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @ManyToOne
    @JsonBackReference
    @JoinColumn(name = "url_id", nullable = false)
    private Url urlId;

    @Builder.Default
    @Column(nullable = false)
    private Long randomness = 0L;

    @Builder.Default
    @Column(nullable = false)
    private Long latency = 0L;

    @Column(nullable = true)
    private String description;

    @Column(nullable = false)
    private String name;

    @Lob
    @JsonProperty("data")
    @Column(nullable = false, columnDefinition = "LONGTEXT")
    private String data;

    @CreationTimestamp
    @Column(nullable = false, updatable = false)
    @JsonIgnore
    private LocalDateTime createdAt;

    @UpdateTimestamp
    @Column(nullable = false)
    @JsonIgnore
    private LocalDateTime updatedAt;

}
