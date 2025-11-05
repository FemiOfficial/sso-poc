package entitities

type AppScope string
type AppIntegration string

const (
	AuthScope           AppScope = "auth"
	AuditLogScope       AppScope = "audit_log"
	UserManagementScope AppScope = "user_management"
)

type App struct {
	BaseEntity
	Name              string                `gorm:"not null" json:"name"`
	ClientID          string                `gorm:"not null" json:"client_id"`
	ClientSecret      string                `gorm:"not null" json:"client_secret"`
	RedirectURI       string                `gorm:"not null" json:"redirect_uri"`
	CallbackURI       string                `gorm:"null; default:null" json:"callback_uri"`
	Live              bool                  `gorm:"not null; default:false" json:"live"`
	Scopes            StringArray           `gorm:"type:text[];not null; default:'{auth,audit_log}'" json:"scopes"`
	AppIdentityProviders []AppIdentityProvider `gorm:"foreignKey:AppID" json:"app_identity_providers"`
	MfaEnabled        bool                  `gorm:"not null" json:"mfa_enabled"`
	OrganizationID    string                `gorm:"not null" json:"organization_id"`
	Organization      *Organization          `gorm:"foreignKey:OrganizationID;references:ID" json:"organization"`
}

func (App) TableName() string {
	return "apps"
}
