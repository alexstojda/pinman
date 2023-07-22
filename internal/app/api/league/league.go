package league

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"net/http"
	"pinman/internal/app/api/auth"
	"pinman/internal/app/api/errors"
	"pinman/internal/app/generated"
	"pinman/internal/models"
	"pinman/internal/utils"
	"strings"
	"time"
)

type Controller struct {
	DB *gorm.DB
}

func NewController(db *gorm.DB) *Controller {
	return &Controller{
		DB: db,
	}
}

func (c *Controller) CreateLeague(ctx *gin.Context) {
	payload := &generated.LeagueCreate{}

	currentUser, err := auth.GetUser(ctx)
	if err != nil {
		errors.AbortWithError(http.StatusForbidden, err.Error(), ctx)
		return
	}

	if err := ctx.ShouldBindJSON(payload); err != nil {
		errors.AbortWithError(http.StatusBadRequest, err.Error(), ctx)
		return
	}

	now := time.Now()
	league := models.League{
		Name:      payload.Name,
		OwnerID:   currentUser.ID,
		Slug:      payload.Slug,
		Location:  payload.Location,
		CreatedAt: now,
		UpdatedAt: now,
	}

	result := c.DB.Create(&league)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key value violates unique") {
			errors.AbortWithError(http.StatusConflict, "league with slug already exists", ctx)
			return
		} else {
			log.Err(result.Error).Msg("failed to create league")
			errors.AbortWithError(http.StatusInternalServerError, "failed to create league", ctx)
			return
		}
	}

	ctx.JSON(http.StatusCreated, generated.LeagueResponse{
		League: &generated.League{
			Id:        league.ID.String(),
			Name:      league.Name,
			Slug:      league.Slug,
			Location:  league.Location,
			OwnerId:   league.OwnerID.String(),
			CreatedAt: utils.FormatTime(league.CreatedAt),
			UpdatedAt: utils.FormatTime(league.UpdatedAt),
		},
	})
}

func (c *Controller) ListLeagues(ctx *gin.Context) {
	var dbResults []models.League
	result := c.DB.Find(&dbResults)
	if result.Error != nil {
		log.Err(result.Error).Msg("failed to list leagues")
		errors.AbortWithError(http.StatusInternalServerError, "failed to list leagues", ctx)
		return
	}

	var leagues []generated.League
	for _, league := range dbResults {
		leagues = append(leagues, generated.League{
			Id:        league.ID.String(),
			Name:      league.Name,
			Slug:      league.Slug,
			Location:  league.Location,
			OwnerId:   league.Owner.ID.String(),
			CreatedAt: utils.FormatTime(league.CreatedAt),
			UpdatedAt: utils.FormatTime(league.UpdatedAt),
		})
	}

	ctx.JSON(http.StatusOK, generated.LeagueListResponse{
		Leagues: leagues,
	})
}
