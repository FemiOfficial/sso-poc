package auth

import (
	"sso-poc/cmd/lib/auth"
	"sso-poc/internal/db/entitities"

	"github.com/gin-gonic/gin"
)

type AuthService struct {
	authLib *auth.AuthLib
}

func CreateAuthService(authLib *auth.AuthLib) *AuthService {
	return &AuthService{authLib: authLib}
}

func (s *AuthService) InitiateAuthSession(ctx *gin.Context) (*string, error, int, gin.H) {
	app := ctx.MustGet("app").(*entitities.App)
	providers := []string{ctx.Query("provider")}

	message, err, statusCode, data := s.authLib.InitiateAuthSession(ctx, app, providers)
	if err != nil {
		return message, err, statusCode, nil
	}
	return message, nil, statusCode, data
}

// func (s *AuthService) LoginUser(ctx *gin.Context) {
// 	// create new auth with goth
// 	// begin auth session or complete it if that is the case
// 	// begin auth will calls .BeginAuth from goth api

// 	var loginUserRequest types.LoginUserRequest
// 	if err := ctx.ShouldBindJSON(&loginUserRequest); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	provider := loginUserRequest.Provider
// 	sessionId := loginUserRequest.SessionID

// 	app := ctx.MustGet("app").(*entitities.App)

// 	authRequest := &entitities.AuthRequest{}
// 	err := s.redis.Get(ctx, sessionId).Scan(authRequest)
// 	if err != nil {
// 		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session id, session not found"})
// 		return
// 	}

// 	if authRequest.State.Status != "initiated" {
// 		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session id, session is not yet initiated"})
// 		return
// 	}

// 	if !slices.Contains(authRequest.Providers, provider) {
// 		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid provider, provider is not valid for this session"})
// 		return
// 	}

// 	callbackURL := fmt.Sprintf("%s/auth/%s/%s/callback", os.Getenv("APP_URL"), provider, sessionId)

// 	appIdentityProvider := &entitities.AppIdentityProvider{}
// 	s.db.DB.Joins("IdentityProvider").Where("app_id = ?", app.ID).Where("identity_provider.name = ?", provider).First(appIdentityProvider)
// 	providerInstance, err := CreateProvider(appIdentityProvider, s.vaultEncrypt, callbackURL)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	providerInstance.SetName(provider)

// 	goth.UseProviders(providerInstance)

// 	gothic.BeginAuthHandler(ctx.Writer, ctx.Request)
// }

// func (s *AuthService) Callback(ctx *gin.Context) {
// 	provider := ctx.Param("provider")
// 	sessionId := ctx.Query("session_id")

// 	query := ctx.Request.URL.Query()
// 	query.Set("provider", provider)
// 	query.Set("session_id", sessionId)
// 	ctx.Request.URL.RawQuery = query.Encode()

// 	user, err := gothic.CompleteUserAuth(ctx.Writer, ctx.Request)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	fmt.Println(user)
// 	http.Redirect(ctx.Writer, ctx.Request, "/api/auth/profile", http.StatusTemporaryRedirect)
// }

// func (s *AuthService) GetAuthProfileData(ctx *gin.Context) {
// 	sessionId := ctx.Query("session_id")
// 	dbConnection := s.db.DB

// 	// Get app from context (set by middleware)
// 	app := ctx.MustGet("app").(*entitities.App)

// 	authRequest := &entitities.AuthRequest{}
// 	err := s.redis.Get(ctx, sessionId).Scan(authRequest)

// 	if err != nil {
// 		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session id, session not found"})
// 		return
// 	}

// 	if authRequest.AppID != app.ID {
// 		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session id, session is not for this app"})
// 		return
// 	}

// 	if authRequest.State.Status != "auth_completed" {
// 		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session id, session is not yet completed"})
// 		return
// 	}

// 	user := &entitities.User{}

// 	email := authRequest.State.Data["email"].Email

// 	if email == nil {
// 		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session id, email not found"})
// 		return
// 	}

// 	dbConnection.Where("email = ?", authRequest.State.Data["email"].Email).First(user)
// 	if user.ID == "" {
// 		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, user)
// }
