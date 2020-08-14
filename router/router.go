package router

import (
	"database/sql"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/mensurowary/juno/auth"
	"github.com/mensurowary/juno/commons"
	"github.com/mensurowary/juno/config"
	"github.com/mensurowary/juno/resources"
	"github.com/mensurowary/juno/resources/download"
	"github.com/mensurowary/juno/resources/interactions"
	"github.com/mensurowary/juno/resources/upload"
	"github.com/sirupsen/logrus"
	"net/http"
)

func Initialize(db *sql.DB) *gin.Engine {
	// dependencies
	uploadService := upload.NewService(db)
	downloadService := download.NewService(db)
	interactionService := interactions.NewService(db, downloadService)
	//
	engine := gin.Default()

	authMiddleware := auth.JwtMiddleware()

	engine.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		logrus.Printf("NoRoute claims: %#v", claims)
		c.JSON(http.StatusNotFound, commons.MakeFailureResponse(
			"Page not found", http.StatusNotFound,
		))
	})

	versioning := engine.Group(config.Config.ApiVersion)
	{
		versioning.POST("/auth/login", authMiddleware.LoginHandler)
		versioning.POST("/auth/refresh_token", authMiddleware.RefreshHandler)
		versioning.POST("/auth/logout", authMiddleware.LogoutHandler)

		resourcesGroup := versioning.Group("/resources")
		resourcesGroup.Use(authMiddleware.MiddlewareFunc())
		{
			resourcesGroup.Handle(http.MethodGet, "", resources.GetAppResourcesInformation(downloadService))
			resourcesGroup.Handle(http.MethodPost, "/upload", resources.Upload(uploadService))
			resourcesGroup.Handle(http.MethodGet, "/:id", resources.DownloadSingleAppResource(downloadService))
			resourcesGroup.Handle(http.MethodDelete, "/:id", resources.DeleteSingleAppResource(interactionService))
		}
	}

	return engine
}
