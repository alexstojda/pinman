package user

import (
	"github.com/rs/zerolog/log"
	"net/http"
	"pinman/internal/app/api/auth"
	"pinman/internal/app/api/errors"
	"pinman/internal/app/generated"
	"pinman/internal/models"
	"pinman/internal/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Controller struct {
	DB *gorm.DB
}

func NewController(DB *gorm.DB) *Controller {
	return &Controller{DB}
}

func (c *Controller) GetMe(ctx *gin.Context) {
	currentUser, err := auth.GetUser(ctx)
	if err != nil {
		errors.AbortWithError(http.StatusForbidden, err.Error(), ctx)
	}

	ctx.JSON(http.StatusOK, generated.UserResponse{
		User: &generated.User{
			Id:        currentUser.ID.String(),
			Name:      currentUser.Name,
			Email:     currentUser.Email,
			Role:      currentUser.Role,
			CreatedAt: utils.FormatTime(currentUser.CreatedAt),
			UpdatedAt: utils.FormatTime(currentUser.UpdatedAt),
		},
	})
}

func (c *Controller) SignUpUser(ctx *gin.Context) {
	payload := &generated.UserRegister{}

	if err := ctx.ShouldBindJSON(payload); err != nil {
		errors.AbortWithError(http.StatusBadRequest, err.Error(), ctx)
		return
	}

	if payload.Password != payload.PasswordConfirm {
		errors.AbortWithError(http.StatusBadRequest, "passwords don't match", ctx)
		return
	}

	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		errors.AbortWithError(http.StatusInternalServerError, err.Error(), ctx)
		return
	}

	now := time.Now()
	newUser := models.User{
		Name:      payload.Name,
		Email:     strings.ToLower(payload.Email),
		Password:  hashedPassword,
		Role:      "user",
		Verified:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	result := c.DB.Create(&newUser)

	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key value violates unique") {
			errors.AbortWithError(http.StatusConflict, "user with that email already exists", ctx)
			return
		} else {
			log.Err(result.Error).Msg("failed to create user")
			errors.AbortWithError(http.StatusInternalServerError, "failed to create user", ctx)
			return
		}
	}

	ctx.JSON(http.StatusCreated, generated.UserResponse{
		User: &generated.User{
			Id:        newUser.ID.String(),
			Name:      newUser.Name,
			Email:     newUser.Email,
			Role:      newUser.Role,
			CreatedAt: utils.FormatTime(newUser.CreatedAt),
			UpdatedAt: utils.FormatTime(newUser.UpdatedAt),
		},
	})
}
