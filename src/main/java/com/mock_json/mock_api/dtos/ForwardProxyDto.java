package com.mock_json.mock_api.dtos;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.*;

import java.time.LocalDateTime;

@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
@Builder
public class ForwardProxyDto {

    private Long id;

    @JsonProperty("project_id")
    private Long projectId;

    private String domain;

    @JsonProperty("created_at")
    private LocalDateTime createdAt;

    @JsonProperty("updated_at")
    private LocalDateTime updatedAt;

    @JsonProperty("is_active")
    private Boolean isActive;

}
