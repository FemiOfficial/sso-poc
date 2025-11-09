package types


type AuthRequestState struct {
	Status string `json:"status" gorm:"not null;enum:initiated,auth_id_opened,auth_failed,auth_completed"`
	Data   map[string]any `json:"data" gorm:"type:jsonb"`
}

type SessionProviders struct {
	ID string `gorm:"not null" json:"id"`
	Name string `gorm:"not null" json:"name"`
	Logo string `gorm:"not null" json:"logo"`
	Status string `gorm:"not null" json:"status"`
	Scopes []string `gorm:"not null" json:"scopes"`
	DisplayName string `gorm:"not null" json:"display_name"`
	CredentialFields []string `gorm:"type:text[];not null" json:"credential_fields"`
}

type ResolveSessionResponse struct {
	SessionID string `gorm:"not null" json:"session_id"`
	AppID     string `gorm:"not null" json:"app_id"`
	Status     string   `gorm:"not null" json:"status"`
	State      AuthRequestState `gorm:"not null" json:"state"`
	Providers  []SessionProviders `gorm:"not null" json:"providers"`
}