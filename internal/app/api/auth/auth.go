package auth

import (
	"fmt"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"net/http"
	"pinman/internal/app/api/errors"
	"pinman/internal/app/generated"
	"pinman/internal/models"
	"pinman/internal/utils"
	"reflect"
	"strings"
	"time"
)

const (
	IdentityKey = "user"
)

func CreateJWTMiddleware(config *utils.Config, db *gorm.DB) (*jwt.GinJWTMiddleware, error) {
	// the jwt middleware
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:           "pinman",
		Key:             []byte(config.TokenSecretKey),
		Timeout:         config.TokenExpiresAfter,
		MaxRefresh:      config.TokenExpiresAfter,
		PubKeyBytes:     []byte(config.TokenPublicKey),
		PrivKeyBytes:    []byte(config.TokenPrivateKey),
		IdentityKey:     IdentityKey,
		PayloadFunc:     payloadFunc,
		IdentityHandler: getIdentityHandlerFunc(db),
		Authenticator:   getAuthenticatorFunc(db),
		Unauthorized:    unauthorizedFunc,
		Authorizator:    authorizationFunc,
		LoginResponse:   loginResponseFunc,
		RefreshResponse: refreshResponseFunc,
		TokenLookup:     "header:Authorization",
		TokenHeadName:   "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create middleware: %v", err)
	}

	return authMiddleware, nil
}

func GetUser(ctx *gin.Context) (*models.User, error) {
	user, ok := ctx.Get(IdentityKey)
	if !ok && reflect.TypeOf(user).AssignableTo(reflect.TypeOf(&models.User{})) {
		return nil, fmt.Errorf("user obj not populated in ctx")
	}

	return user.(*models.User), nil
}

func GetAuthMiddlewareFunc(mw *jwt.GinJWTMiddleware) func(c *gin.Context) {
	return func(c *gin.Context) {
		// Only run the JWT middleware if auth scopes are required for the endpoint
		if _, ok := c.Get(generated.PinmanAuthScopes); ok {
			mw.MiddlewareFunc()(c)
		}
		c.Next()
		return
	}
}

func payloadFunc(data interface{}) jwt.MapClaims {
	if v, ok := data.(*models.User); ok {
		return jwt.MapClaims{
			IdentityKey: v.ID.String(),
		}
	}
	return jwt.MapClaims{}
}

func unauthorizedFunc(c *gin.Context, code int, message string) {
	errors.AbortWithError(code, message, c)
	return
}

func authorizationFunc(data interface{}, c *gin.Context) bool {
	if scopes, ok := c.Get(generated.PinmanAuthScopes); ok {
		if user, ok := data.(*models.User); ok {
			for _, scope := range scopes.([]string) {
				if user.Role == scope {
					return true
				}
			}
			c.Abort()
			return false
		} else {
			log.Err(fmt.Errorf("data var was not *models.User"))
			return false
		}
	}

	return true
}

func loginResponseFunc(c *gin.Context, _ int, token string, expire time.Time) {
	c.JSON(http.StatusOK, generated.TokenResponse{
		AccessToken: token,
		Expire:      expire.Format(time.RFC3339),
	})
}

func refreshResponseFunc(c *gin.Context, _ int, token string, expire time.Time) {
	c.JSON(http.StatusOK, generated.TokenResponse{
		AccessToken: token,
		Expire:      expire.Format(time.RFC3339),
	})
}

func getIdentityHandlerFunc(db *gorm.DB) func(c *gin.Context) interface{} {
	return func(c *gin.Context) interface{} {
		claims := jwt.ExtractClaims(c)

		user := &models.User{}
		result := db.First(user, "id = ?", claims[IdentityKey])
		if result.Error != nil {
			return nil
		}

		return user
	}
}

func getAuthenticatorFunc(db *gorm.DB) func(ctx *gin.Context) (interface{}, error) {
	return func(ctx *gin.Context) (interface{}, error) {
		payload := &generated.UserLogin{}

		if ctx.ContentType() == "application/json" {
			if err := ctx.ShouldBind(payload); err != nil {
				return nil, err
			}
		} else if ctx.ContentType() == "application/x-www-form-urlencoded" {
			payload.Username = ctx.PostForm("username")
			payload.Password = ctx.PostForm("password")
		} else {
			return nil, fmt.Errorf("unsuported content-type")
		}

		user := &models.User{}
		result := db.First(user, "email = ?", strings.ToLower(payload.Username))
		if result.Error != nil {
			log.Err(result.Error).Msg("failed to authenticate user - database query failed")
			return nil, jwt.ErrFailedAuthentication
		}

		if err := utils.VerifyPassword(user.Password, payload.Password); err != nil {
			log.Err(jwt.ErrFailedAuthentication).Str("userId", user.ID.String()).Msg("invalid password")
			return nil, jwt.ErrFailedAuthentication
		}

		return user, nil
	}
}
