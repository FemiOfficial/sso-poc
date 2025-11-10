package auth

import (
	"net/http"
	"sso-poc/cmd/lib/auth"
	authTypes "sso-poc/cmd/lib/auth/types"
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

func (s *AuthService) ResolveSession(ctx *gin.Context) (*string, error, int, gin.H) {
	sessionId := ctx.Param("sessionId")

	authRequest, err := s.authLib.ResolveSession(sessionId)
	if err != nil {
		message := "Session not found"
		return &message, err, http.StatusNotFound, nil
	}

	message := "Session resolved successfully"
	data := gin.H{"session": authRequest}
	return &message, nil, http.StatusOK, data
}

// func (s *AuthService) LoginUser(ctx *gin.Context) (*string, error, int, gin.H) {
// 	app := ctx.MustGet("app").(*entitities.App)
// 	provider := ctx.Query("provider")
// 	sessionId := ctx.Query("session_id")

// 	message, err, statusCode, _ := s.authLib.LoginUser(ctx, app, providerObject, sessionId)
// 	if err != nil {
// 		return message, err, statusCode, nil
// 	}
// 	return message, nil, statusCode, nil
// }

func (s *AuthService) LoginUserWithSession(ctx *gin.Context) (*string, error, int, gin.H) {
	sessionId := ctx.Param("sessionId")
	provider := ctx.Query("provider")

	// Resolve session to get app
	authRequest, err := s.authLib.ResolveSession(sessionId)
	if err != nil {
		message := "Session not found"
		return &message, err, 404, nil
	}

	if provider == "" {
		message := "Provider is required"
		return &message, nil, 400, nil
	}

	providerObject := &authTypes.SessionProviders{}
	for _, _provider := range authRequest.Providers {
		if _provider.Name == provider {
			providerObject = &_provider
			break
		}
	}

	if providerObject == nil {
		message := "Invalid provider"
		return &message, nil, http.StatusBadRequest, nil
	}

	// Get app from authRequest
	app := &entitities.App{}
	err = s.authLib.GetDB().DB.Where("id = ?", authRequest.AppID).First(app).Error
	if err != nil {
		message := "App not found"
		return &message, err, http.StatusNotFound, nil
	}

	message, err, statusCode, _ := s.authLib.LoginUser(ctx, app, providerObject, sessionId)
	if err != nil {
		return message, err, statusCode, nil
	}
	return message, nil, statusCode, nil
}

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
