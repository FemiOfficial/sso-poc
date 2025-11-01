package entitities

type AppScope string
type AppIntegration string

const (
	googleIntegration    AppIntegration = "google"
	githubIntegration    AppIntegration = "github"
	facebookIntegration  AppIntegration = "facebook"
	linkedinIntegration  AppIntegration = "linkedin"
	slackIntegration     AppIntegration = "slack"
	teamsIntegration     AppIntegration = "teams"
	zoomIntegration      AppIntegration = "zoom"
	whatsappIntegration  AppIntegration = "whatsapp"
	telegramIntegration  AppIntegration = "telegram"
	discordIntegration   AppIntegration = "discord"
	spotifyIntegration   AppIntegration = "spotify"
	twitchIntegration    AppIntegration = "twitch"
	redditIntegration    AppIntegration = "reddit"
	tiktokIntegration    AppIntegration = "tiktok"
	instagramIntegration AppIntegration = "instagram"
	twitterIntegration   AppIntegration = "twitter"
)

const (
	AuthScope           AppScope = "auth"
	AuditLogScope       AppScope = "audit_log"
	UserManagementScope AppScope = "user_management"
)

type App struct {
	BaseEntity
	Name              string                `gorm:"not null"`
	ClientID          string                `gorm:"not null"`
	ClientSecret      string                `gorm:"not null"`
	RedirectURI       string                `gorm:"not null"`
	Live              bool                  `gorm:"not null; default:false"`
	Scopes            []AppScope            `gorm:"type:text[];not null; default:'{auth,audit_log}'"`
	IdentityProviders []AppIdentityProvider `gorm:"many2many:app_identity_providers;"`
	MfaEnabled        bool                  `gorm:"not null"`
	OrganizationID    string                `gorm:"not null"`
	Organization      Organization          `gorm:"foreignKey:OrganizationID"`
}

func (App) TableName() string {
	return "apps"
}
