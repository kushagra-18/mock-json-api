package models

// Project represents a project within a team
type Project struct {
	BaseModel
	Name                 string        `gorm:"not null" json:"name"`
	Slug                 string        `gorm:"unique;not null" json:"slug"`
	ChannelID            string        `gorm:"not null" json:"channel_id"` // Assuming this is a Slack channel ID or similar
	Description          string        `json:"description"`
	IsForwardProxyActive bool          `gorm:"default:false" json:"is_forward_proxy_active"`
	TeamID               uint          `json:"team_id"`                                             // Foreign key for Team
	Team                 Team          `json:"team,omitempty"`                                      // Belongs to Team
	ForwardProxy         *ForwardProxy `gorm:"foreignKey:ProjectID" json:"forward_proxy,omitempty"` // Has one ForwardProxy, use pointer
}
