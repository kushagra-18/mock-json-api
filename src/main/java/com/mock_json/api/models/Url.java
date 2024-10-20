package com.mock_json.api.models;

import jakarta.persistence.*;
import lombok.*;

import org.hibernate.annotations.SoftDelete;
import org.hibernate.annotations.UpdateTimestamp;

import com.fasterxml.jackson.annotation.JsonBackReference;
import com.fasterxml.jackson.annotation.JsonIgnore;
import com.fasterxml.jackson.annotation.JsonManagedReference;
import com.fasterxml.jackson.annotation.JsonProperty;

import org.hibernate.annotations.CreationTimestamp;

import java.time.LocalDateTime;
import java.util.List;

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

    @Column(nullable = true)
    private String description;

    @Column(nullable = false)
    private String name;

    // for mocking rate limits
    @Column(nullable = true)
    private  Integer requests;

    @Column(nullable = true)
    private  Long time;

    @Column(nullable = false, unique = true)
    private String url;

    // @Column(nullable = false)
    // @Enumerated(EnumType.STRING) 
    // private StatusCode status; 

    @CreationTimestamp
    @Column(nullable = false, updatable = false)
    @JsonIgnore
    private LocalDateTime createdAt;

    @UpdateTimestamp
    @Column(nullable = false)
    @JsonIgnore
    private LocalDateTime updatedAt;

    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "project_id")
    @JsonBackReference
    private Project project;

    @OneToMany(mappedBy = "urlId", cascade = CascadeType.ALL, fetch = FetchType.EAGER)
    @JsonProperty("mock_content_list")
    @JsonManagedReference
    private List<MockContent> mockContentList;

}
