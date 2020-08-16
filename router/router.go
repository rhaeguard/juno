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
	engine := gin.Default()

	// dependencies init
	us := upload.NewService(db)
	dr := download.NewRepository(db)
	ds := download.NewService(dr)
	ir := interactions.NewRepository(db)
	is := interactions.NewService(ir, ds)
	// dependencies init end

	authMiddleware := auth.JwtMiddleware()

	engine.NoRoute(authMiddleware.MiddlewareFunc(), NoRouteHandler())

	versioning := engine.Group(config.Config.ApiVersion)
	{
		versioning.POST("/auth/login", authMiddleware.LoginHandler)
		versioning.POST("/auth/refresh_token", authMiddleware.RefreshHandler)
		versioning.POST("/auth/logout", authMiddleware.LogoutHandler)

		resourcesGroup := versioning.Group("/resources")
		resourcesGroup.Use(authMiddleware.MiddlewareFunc())
		{
			resourcesGroup.Handle(http.MethodGet, "", resources.GetAppResourcesInformation(ds))
			resourcesGroup.Handle(http.MethodPost, "/upload", resources.Upload(us))
			resourcesGroup.Handle(http.MethodGet, "/:id", resources.DownloadSingleAppResource(ds))
			resourcesGroup.Handle(http.MethodDelete, "/:id", resources.DeleteSingleAppResource(is))
		}
	}

	return engine
}

func NoRouteHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		logrus.Printf("NoRoute claims: %#v", claims)
		c.JSON(http.StatusNotFound, commons.MakeFailureResponse(
			"Page not found", http.StatusNotFound,
		))
	}
}
