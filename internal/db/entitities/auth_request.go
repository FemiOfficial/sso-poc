package entitities

type AuthRequestState string

const (
	AuthRequestStateInitiated   AuthRequestState = "initiated"
	AuthRequestStateAuthIDOpened AuthRequestState = "auth_id_opened"
	AuthRequestStateFailed      AuthRequestState = "auth_failed"
	AuthRequestStateCompleted   AuthRequestState = "auth_completed"
)

type AuthRequest struct {
	BaseEntity
	SessionID string            `gorm:"not null"`
	AppID     string            `gorm:"not null"`
	App       App               `gorm:"foreignKey:AppID"`
	Providers []string          `gorm:"not null;type:text[]"`
	State     AuthRequestState `gorm:"not null"`
}

func (AuthRequest) TableName() string {
	return "auth_requests"
}
