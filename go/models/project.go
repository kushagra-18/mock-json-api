package models

// Project represents a project within a team
type Project struct {
	BaseModel
	Name                 string `gorm:"not null"`
	Slug                 string `gorm:"unique;not null"`
	ChannelID            string `gorm:"not null"` // Assuming this is a Slack channel ID or similar
	Description          string
	IsForwardProxyActive bool   `gorm:"default:false"`
	TeamID               uint   // Foreign key for Team
	Team                 Team   // Belongs to Team
	ForwardProxy         *ForwardProxy `gorm:"foreignKey:ProjectID"` // Has one ForwardProxy, use pointer
}
