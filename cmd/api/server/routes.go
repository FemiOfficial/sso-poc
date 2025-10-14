package server

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	routes := gin.Default()

	routes.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true, // Enable cookies/auth
	}))

	api := routes.Group("/api")
	{
		api.GET("/health", s.healthHandler)
		api.POST("/auth/initiate", s.authController.InitiateAuthSession)
		api.GET()
	}

	lib := routes.Group("/lib")
	{
		lib.GET("/auth/:provider", s.authHandler)
		lib.GET("/auth/:provider/callback", s.callBackHandler)
		lib.GET("/auth/logout", s.logoutHandler)
	}

	return routes
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) callBackHandler(c *gin.Context) {
	provider := c.Param("provider")
	q := c.Request.URL.Query()
	q.Set("provider", provider)
	c.Request.URL.RawQuery = q.Encode()

	user, err := s.auth.CompleteAuth(c.Writer, c.Request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Println(user)
	http.Redirect(c.Writer, c.Request, "/", http.StatusTemporaryRedirect)
}

func (s *Server) authHandler(c *gin.Context) {
	provider := c.Param("provider")
	q := c.Request.URL.Query()
	q.Set("provider", provider)
	c.Request.URL.RawQuery = q.Encode()

	if gothUser, err := s.auth.CompleteAuth(c.Writer, c.Request); err == nil {
		c.JSON(http.StatusOK, gothUser)
	} else {
		s.auth.BeginAuth(c.Writer, c.Request)
	}
}

func (s *Server) logoutHandler(c *gin.Context) {
	s.auth.Logout(c.Writer, c.Request)
	c.Header("Location", "/")
	http.Redirect(c.Writer, c.Request, "/", http.StatusTemporaryRedirect)
}
