package entitities

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type AuthRequestStateData struct {
	Email      *string `json:"email,omitempty"`
	FirstName  *string `json:"first_name,omitempty"`
	LastName   *string `json:"last_name,omitempty"`
	AvatarURL  *string `json:"avatar_url,omitempty"`
	ProviderID *string `json:"provider_id,omitempty"`
}

type AuthRequestState struct {
	Status string                          `json:"status" gorm:"not null;enum:initiated,auth_id_opened,auth_failed,auth_completed"`
	Data   map[string]AuthRequestStateData `json:"data" gorm:"type:jsonb"`
}

// Value implements the driver.Valuer interface for database storage
func (a AuthRequestState) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan implements the sql.Scanner interface for database retrieval
func (a *AuthRequestState) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("cannot scan %T into AuthRequestState", value)
	}

	return json.Unmarshal(bytes, a)
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
