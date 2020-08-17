package auth

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/mensurowary/juno/commons"
	"github.com/mensurowary/juno/config"
	log "github.com/sirupsen/logrus"
	"time"
)

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

var identityKey = "app_id"

// Application represents the application/user of juno
type Application struct {
	ID string
}

// JwtMiddleware is the middleware that handles authentication and authorization
func JwtMiddleware() *jwt.GinJWTMiddleware {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:           config.Config.JwtRealm,
		Key:             []byte(config.Config.JwtSecret),
		Timeout:         15 * time.Minute,
		MaxRefresh:      15 * time.Minute,
		IdentityKey:     identityKey,
		PayloadFunc:     payloadHandler(),
		IdentityHandler: identityHandler(),
		Authenticator:   authenticator(),
		Authorizator:    authorizer(),
		Unauthorized:    unauthorized(),
		TokenLookup:     "header: Authorization",
		TokenHeadName:   "Bearer",
		TimeFunc:        time.Now,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	return authMiddleware

}

func authorizer() func(data interface{}, c *gin.Context) bool {
	return func(data interface{}, c *gin.Context) bool {
		if v, ok := data.(*Application); ok && v.ID == "admin" {
			return true
		}

		return false
	}
}

func authenticator() func(c *gin.Context) (interface{}, error) {
	return func(c *gin.Context) (interface{}, error) {
		var loginValues login
		if err := c.ShouldBind(&loginValues); err != nil {
			return "", jwt.ErrMissingLoginValues
		}
		userID := loginValues.Username
		password := loginValues.Password

		if (userID == "admin" && password == "admin") || (userID == "test" && password == "test") {
			return &Application{
				ID: userID,
			}, nil
		}

		return nil, jwt.ErrFailedAuthentication
	}
}

func identityHandler() func(c *gin.Context) interface{} {
	return func(c *gin.Context) interface{} {
		claims := jwt.ExtractClaims(c)
		return &Application{
			ID: claims[identityKey].(string),
		}
	}
}

func payloadHandler() func(data interface{}) jwt.MapClaims {
	return func(data interface{}) jwt.MapClaims {
		if v, ok := data.(*Application); ok {
			return jwt.MapClaims{
				identityKey: v.ID,
			}
		}
		return jwt.MapClaims{}
	}
}

func unauthorized() func(c *gin.Context, code int, message string) {
	return func(c *gin.Context, code int, message string) {
		c.JSON(code, commons.MakeFailureResponse(
			"Unauthorized", uint16(code),
		))
	}
}
