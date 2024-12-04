package com.mock_json.mock_api.models;

import jakarta.persistence.*;
import lombok.*;

import org.hibernate.annotations.SoftDelete;
import org.hibernate.annotations.UpdateTimestamp;

import com.fasterxml.jackson.annotation.JsonBackReference;
import com.fasterxml.jackson.annotation.JsonProperty;

import org.hibernate.annotations.CreationTimestamp;

import java.time.LocalDateTime;

@Entity
@Table(name = "projects")
@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
@Builder
@SoftDelete
public class Project {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @Column(nullable = false)
    private String name;

    @Column(nullable = false, unique = true)
    private String slug;

    @CreationTimestamp
    @JsonProperty("created_at")
    @Column(nullable = false, updatable = false)
    private LocalDateTime createdAt;

    @Column(nullable = false)
    @JsonProperty("channel_id")
    private String ChannelId;

    @Column(nullable = true,columnDefinition = "TEXT")
    private String description;

    @UpdateTimestamp
    @JsonProperty("updated_at")
    @Column(nullable = false)
    private LocalDateTime updatedAt;

    @ManyToOne(fetch = FetchType.EAGER)
    @JoinColumn(name = "team_id")
    @JsonBackReference
    private Team team;


    @OneToOne(fetch = FetchType.EAGER)
    @JoinColumn(name = "project_id")
    @JsonBackReference
    private ForwardProxy forwardProxy;

    @Column(name = "is_forward_proxy_active", columnDefinition = "BOOLEAN DEFAULT FALSE")
    @JsonProperty("is_forward_proxy_active")
    @Builder.Default
    private Boolean isForwardProxyActive = false;
}
