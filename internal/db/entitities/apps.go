package entitities

type AppScope string
type AppIntegration string

const (
	googleIntegration AppIntegration = "google"
	githubIntegration AppIntegration = "github"
	facebookIntegration AppIntegration = "facebook"
	linkedinIntegration AppIntegration = "linkedin"
	slackIntegration AppIntegration = "slack"
	teamsIntegration AppIntegration = "teams"
	zoomIntegration AppIntegration = "zoom"
	whatsappIntegration AppIntegration = "whatsapp"
	telegramIntegration AppIntegration = "telegram"
	discordIntegration AppIntegration = "discord"
	spotifyIntegration AppIntegration = "spotify"
	twitchIntegration AppIntegration = "twitch"
	redditIntegration AppIntegration = "reddit"
	tiktokIntegration AppIntegration = "tiktok"
	instagramIntegration AppIntegration = "instagram"
	twitterIntegration AppIntegration = "twitter"
)

const (
	authScope           AppScope = "auth"
	auditLogScope       AppScope = "audit_log"
	userManagementScope AppScope = "user_management"
)



type App struct {
	BaseEntity
	Name           string       `gorm:"not null"`
	ClientID       string       `gorm:"not null"`
	ClientSecret   string       `gorm:"not null"`
	RedirectURI    string       `gorm:"not null"`
	Scopes []AppScope `gorm:"type:app_scope[];not null; default:{auth,audit_log}"`
	Integrations []AppIntegration `gorm:"type:app_integration[];not null; default:{google,github,facebook,linkedin}"`
	MfaEnabled     bool         `gorm:"not null"`
	OrganizationID string       `gorm:"not null"`
	Organization   Organization `gorm:"foreignKey:OrganizationID"`
}

func (App) TableName() string {
	return "apps"
}
