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
	"strings"
	"time"
)

const (
	IdentityKey = "user"
)

func CreateJWTMiddleware(db *gorm.DB) (*jwt.GinJWTMiddleware, error) {
	// the jwt middleware
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "pinman",
		Key:         []byte("secret key"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: IdentityKey,
		// Callback function that will be called during login.
		// Using this function it is possible to add additional payload data to the webtoken.
		// The data is then made available during requests via c.Get("JWT_PAYLOAD").
		// Note that the payload is not encrypted.
		// The attributes mentioned on jwt.io can't be used as keys for the map.
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*models.User); ok {
				return jwt.MapClaims{
					IdentityKey: v.ID.String(),
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)

			user := &models.User{}
			result := db.First(user, "id = ?", claims[IdentityKey])
			if result.Error != nil {
				return nil
			}

			return user
		},
		// Callback function that should perform the authentication of the user based on login info.
		// Must return user data as user identifier, it will be stored in Claim Array. Required.
		// Check error (e) to determine the appropriate error message.
		Authenticator: func(ctx *gin.Context) (interface{}, error) {
			payload := &generated.UserLogin{}

			if ctx.ContentType() == "application/json" {
				if err := ctx.ShouldBind(payload); err != nil {
					return nil, err
				}
			} else if ctx.ContentType() == "application/x-www-form-urlencoded" {
				payload.Username = utils.PtrString(ctx.PostForm("username"))
				payload.Password = utils.PtrString(ctx.PostForm("password"))
			} else {
				errors.AbortWithError(http.StatusBadRequest, "content-type not supported", ctx)
			}

			user := &models.User{}
			result := db.First(user, "email = ?", strings.ToLower(*payload.Username))
			if result.Error != nil {
				log.Err(result.Error).Msg("failed to authenticate user - database query failed")
				return nil, jwt.ErrFailedAuthentication
			}

			if err := utils.VerifyPassword(user.Password, *payload.Password); err != nil {
				log.Err(jwt.ErrFailedAuthentication).Str("userId", user.ID.String()).Msg("invalid password")
				return nil, jwt.ErrFailedAuthentication
			}

			return user, nil
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			errors.AbortWithError(code, message, c)
			return
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if _, ok := c.Get(generated.PinmanAuthScopes); ok {
				if _, ok := data.(*models.User); ok {
					//TODO: Validate authorization
				}
			}

			return true
		},
		LoginResponse: func(c *gin.Context, code int, token string, expire time.Time) {
			c.JSON(http.StatusOK, generated.TokenResponse{
				AccessToken: utils.PtrString(token),
				Expire:      utils.PtrString(expire.Format(time.RFC3339)),
			})
		},
		RefreshResponse: func(c *gin.Context, code int, token string, expire time.Time) {
			c.JSON(http.StatusOK, generated.TokenResponse{
				AccessToken: utils.PtrString(token),
				Expire:      utils.PtrString(expire.Format(time.RFC3339)),
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create middleware: %v", err)
	}

	return authMiddleware, nil
}

func GetUser(ctx *gin.Context) *models.User {
	return ctx.MustGet(IdentityKey).(*models.User)
}
