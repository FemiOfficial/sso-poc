package entitities

type AuthRequestStateData struct {
	Email      *string `json:"email,omitempty"`
	FirstName  *string `json:"first_name,omitempty"`
	LastName   *string `json:"last_name,omitempty"`
	AvatarURL  *string `json:"avatar_url,omitempty"`
	ProviderID *string `json:"provider_id,omitempty"`
}

type AuthRequestState struct {
	Status string                          `gorm:"not null;enum:initiated,auth_id_opened,auth_failed,auth_completed"`
	Data   map[string]AuthRequestStateData `gorm:"type:jsonb;serializer:json"`
}

type AuthRequest struct {
	BaseEntity
	SessionID string           `gorm:"not null"`
	AppID     string           `gorm:"not null"`
	App       App              `gorm:"foreignKey:AppID"`
	Providers []string         `gorm:"not null;type:text[]"`
	State     AuthRequestState `gorm:"not null"`
}

func (AuthRequest) TableName() string {
	return "auth_requests"
}
