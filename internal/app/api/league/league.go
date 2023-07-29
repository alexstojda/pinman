package league

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"net/http"
	"pinman/internal/app/api/auth"
	apierrors "pinman/internal/app/api/errors"
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
		apierrors.AbortWithError(http.StatusForbidden, err.Error(), ctx)
		return
	}

	if err := ctx.ShouldBindJSON(payload); err != nil {
		apierrors.AbortWithError(http.StatusBadRequest, err.Error(), ctx)
		return
	}

	location := models.Location{}
	locationQueryResult := c.DB.Where("id = ?", payload.LocationId).First(&location)
	if locationQueryResult.Error != nil {
		if errors.Is(locationQueryResult.Error, gorm.ErrRecordNotFound) {
			apierrors.AbortWithError(http.StatusBadRequest, fmt.Sprintf("location with id %s does not exist", payload.LocationId), ctx)
			return
		} else {
			log.Err(locationQueryResult.Error).Msg("failed to get location")
			apierrors.AbortWithError(http.StatusInternalServerError, "failed to get location", ctx)
			return
		}
	}

	now := time.Now()
	league := models.League{
		Name:       payload.Name,
		OwnerID:    currentUser.ID,
		Slug:       payload.Slug,
		LocationID: location.ID,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	createResult := c.DB.Create(&league)
	if createResult.Error != nil {
		if strings.Contains(createResult.Error.Error(), "duplicate key value violates unique") {
			apierrors.AbortWithError(http.StatusConflict, "league with slug already exists", ctx)
			return
		} else {
			log.Err(createResult.Error).Msg("failed to create league")
			apierrors.AbortWithError(http.StatusInternalServerError, "failed to create league", ctx)
			return
		}
	}

	ctx.JSON(http.StatusCreated, generated.LeagueResponse{
		League: &generated.League{
			Id:         league.ID.String(),
			Name:       league.Name,
			Slug:       league.Slug,
			LocationId: location.ID.String(),
			OwnerId:    league.OwnerID.String(),
			CreatedAt:  utils.FormatTime(league.CreatedAt),
			UpdatedAt:  utils.FormatTime(league.UpdatedAt),
		},
	})
}

func (c *Controller) ListLeagues(ctx *gin.Context) {
	var dbResults []models.League
	result := c.DB.Find(&dbResults)
	if result.Error != nil {
		log.Err(result.Error).Msg("failed to list leagues")
		apierrors.AbortWithError(http.StatusInternalServerError, "failed to list leagues", ctx)
		return
	}

	var leagues []generated.League
	for _, league := range dbResults {
		leagues = append(leagues, generated.League{
			Id:         league.ID.String(),
			Name:       league.Name,
			Slug:       league.Slug,
			LocationId: league.LocationID.String(),
			OwnerId:    league.OwnerID.String(),
			CreatedAt:  utils.FormatTime(league.CreatedAt),
			UpdatedAt:  utils.FormatTime(league.UpdatedAt),
		})
	}

	ctx.JSON(http.StatusOK, generated.LeagueListResponse{
		Leagues: leagues,
	})
}

func (c *Controller) GetLeagueWithSlug(ctx *gin.Context, slug string) {
	var dbResult models.League
	result := c.DB.Where("slug = ?", slug).First(&dbResult)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			apierrors.AbortWithError(http.StatusNotFound, "league not found", ctx)
			return
		} else {
			log.Err(result.Error).Msg("failed to get league")
			apierrors.AbortWithError(http.StatusInternalServerError, "failed to get league", ctx)
			return
		}
	}

	ctx.JSON(http.StatusOK, generated.LeagueResponse{
		League: &generated.League{
			Id:         dbResult.ID.String(),
			Name:       dbResult.Name,
			Slug:       dbResult.Slug,
			LocationId: dbResult.LocationID.String(),
			OwnerId:    dbResult.OwnerID.String(),
			CreatedAt:  utils.FormatTime(dbResult.CreatedAt),
			UpdatedAt:  utils.FormatTime(dbResult.UpdatedAt),
		},
	})
}
