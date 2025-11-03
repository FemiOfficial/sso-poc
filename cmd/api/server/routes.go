package server

import (
	"net/http"
	"sso-poc/cmd/api/server/auth"
	"sso-poc/cmd/api/server/dashboard/middlewares"
	miscTypes "sso-poc/cmd/api/server/dashboard/misc/types"
	organisationTypes "sso-poc/cmd/api/server/dashboard/organisation/types"
	appTypes "sso-poc/cmd/api/server/dashboard/app/types"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	routes := gin.Default()

	// routes.Use(middlewares.ErrorHandlerMiddleware())

	routes.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Allow all origins for development
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "Client-Id", "Client-Secret"},
		AllowCredentials: false, // Set to false when using wildcard origins
	}))

	api := routes.Group("/api")
	{
		api.GET("/health", s.healthHandler)
	}

	protectedAPI := routes.Group("/api")
	protectedAPI.Use(auth.ClientAuthMiddleware(s.db))
	{
		protectedAPI.POST("/auth/initiate", s.authController.InitiateAuthSession)
		// protectedAPI.GET("/auth/profile", s.authController.GetAuthProfileData)
		protectedAPI.POST("/auth/login", s.authController.LoginUser)
	}

	dashboardAPI := routes.Group("/api/dashboard")
	{
		dashboardAPI.POST("/organisation/create",
			middlewares.ValidateRequestBody[organisationTypes.CreateOrganizationRequest](),
			s.organizationController.CreateOrganization)

		dashboardAPI.POST("/organisation/verifification/email",
			middlewares.ValidateRequestBody[organisationTypes.VerifyEmailRequest](),
			s.organizationController.VerifyOrganizationEmail)

		dashboardAPI.POST("/organisation/signin",
			middlewares.ValidateRequestBody[organisationTypes.LoginOrganizationRequest](),
			s.organizationController.LoginOrganization)

		dashboardAPI.POST("/organisation/verification/email/resend",
			middlewares.ValidateRequestBody[organisationTypes.ResendEmailVerificationOtpRequest](),
			s.organizationController.ResendEmailVerificationOtp)
	}

	protectedDashboardAPIMisc := routes.Group("/api/dashboard/misc")
	protectedDashboardAPIMisc.Use(middlewares.JwtMiddleware(s.db))
	{
		protectedDashboardAPIMisc.GET("/identity-providers",
			middlewares.ValidateRequestQuery[miscTypes.GetIDPRequest](),
			s.miscController.GetIdentityProviders)
	}

	protectedDashboardApps := routes.Group("/api/dashboard/app")
	protectedDashboardApps.Use(middlewares.JwtMiddleware(s.db))
	{
		protectedDashboardApps.POST("/create",
			middlewares.ValidateRequestBody[appTypes.CreateAppRequest](),
			s.appController.CreateApp)
	}
	// lib := routes.Group("/lib")
	// {
	// 	lib.GET("/auth/:provider", s.authHandler)
	// 	lib.GET("/auth/:provider/callback", s.callBackHandler)
	// 	lib.GET("/auth/logout", s.logoutHandler)
	// }

	return routes
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

// func (s *Server) callBackHandler(c *gin.Context) {
// 	provider := c.Param("provider")
// 	q := c.Request.URL.Query()
// 	q.Set("provider", provider)
// 	c.Request.URL.RawQuery = q.Encode()

// 	user, err := s.auth.CompleteAuth(c.Writer, c.Request)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	fmt.Println(user)
// 	http.Redirect(c.Writer, c.Request, "/", http.StatusTemporaryRedirect)
// }

// func (s *Server) authHandler(c *gin.Context) {
// 	provider := c.Param("provider")
// 	q := c.Request.URL.Query()
// 	q.Set("provider", provider)
// 	c.Request.URL.RawQuery = q.Encode()

// 	if gothUser, err := s.auth.CompleteAuth(c.Writer, c.Request); err == nil {
// 		c.JSON(http.StatusOK, gothUser)
// 	} else {
// 		s.auth.BeginAuth(c.Writer, c.Request)
// 	}
// }

// func (s *Server) logoutHandler(c *gin.Context) {
// 	s.auth.Logout(c.Writer, c.Request)
// 	c.Header("Location", "/")
// 	http.Redirect(c.Writer, c.Request, "/", http.StatusTemporaryRedirect)
// }
