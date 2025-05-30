package dtos

// ForwardProxyDTO is used for creating or updating forward proxy settings.
type ForwardProxyDTO struct {
	ProjectID uint   `json:"projectId" binding:"required"`
	Domain    string `json:"domain" binding:"required,url|fqdn_rfc1123"` // Validate as URL or FQDN
}

// UpdateForwardProxyStatusDTO is used for updating the active status of a forward proxy.
type UpdateForwardProxyStatusDTO struct {
	IsActive bool `json:"isActive"` // No binding:"required", as default false is acceptable if missing
}
