package tournament

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"net/http"
	apierrors "pinman/internal/app/api/errors"
	"pinman/internal/app/generated"
	"pinman/internal/models"
	"pinman/internal/utils"
	"strings"
)

type Controller struct {
	DB *gorm.DB
}

func NewController(db *gorm.DB) *Controller {
	return &Controller{
		DB: db,
	}
}

// CreateTournament creates a new tournament
func (c *Controller) CreateTournament(ctx *gin.Context) {
	payload := &generated.TournamentCreate{}

	if err := ctx.ShouldBindJSON(payload); err != nil {
		apierrors.AbortWithError(http.StatusBadRequest, err.Error(), ctx)
		return
	}

	if err := validateSettingsPayload(payload); err != nil {
		apierrors.AbortWithError(http.StatusBadRequest, err.Error(), ctx)
		return
	}

	league := models.League{}
	if err := c.DB.First(&league, "id = ?", payload.LeagueId).Error; err != nil {
		apierrors.AbortWithError(http.StatusBadRequest, fmt.Sprintf("invalid league: %s", err.Error()), ctx)
		return
	}

	location := models.Location{}
	if err := c.DB.First(&location, "id = ?", payload.LocationId).Error; err != nil {
		apierrors.AbortWithError(http.StatusBadRequest, fmt.Sprintf("invalid location: %s", err.Error()), ctx)
		return
	}

	settings, err := payload.Settings.MarshalJSON()
	if err != nil {
		apierrors.AbortWithError(http.StatusInternalServerError, err.Error(), ctx)
		return
	}

	tournament := &models.Tournament{
		Name:       payload.Name,
		Slug:       payload.Slug,
		Type:       payload.Type,
		Settings:   settings,
		LocationID: location.ID,
		LeagueID:   league.ID,
	}

	result := c.DB.Create(tournament)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key") {
			apierrors.AbortWithError(http.StatusBadRequest, "tournament with slug already exists", ctx)
			return
		} else {
			log.Error().Err(result.Error).Msg("failed to create tournament")
			apierrors.AbortWithError(http.StatusInternalServerError, "failed to create tournament", ctx)
			return
		}
	}

	ctx.JSON(http.StatusCreated, generated.TournamentResponse{
		Tournament: generated.Tournament{
			Id:   tournament.ID.String(),
			Name: tournament.Name,
			Slug: tournament.Slug,
			Type: tournament.Type,
			// Use the original payload settings to avoid having to unmarshal/marshal to the generated type
			Settings:   payload.Settings,
			LocationId: tournament.LocationID.String(),
			LeagueId:   tournament.LeagueID.String(),
			CreatedAt:  utils.FormatTime(tournament.CreatedAt),
			UpdatedAt:  utils.FormatTime(tournament.UpdatedAt),
		},
	})

}

func validateSettingsPayload(payload *generated.TournamentCreate) error {
	switch payload.Type {
	case generated.MultiRoundTournament:
		settings, err := payload.Settings.AsMultiRoundTournamentSettings()
		if err != nil {
			return err
		}
		err = binding.Validator.ValidateStruct(settings)
		if err != nil {
			return err
		}
	}

	return nil
}

// ListTournaments lists all tournaments
func (c *Controller) ListTournaments(ctx *gin.Context) {
	var tournaments []models.Tournament
	result := c.DB.Find(&tournaments)
	if result.Error != nil {
		log.Error().Err(result.Error).Msg("failed to list tournaments")
		apierrors.AbortWithError(http.StatusInternalServerError, "failed to list tournaments", ctx)
		return
	}

	response := generated.TournamentListResponse{
		Tournaments: []generated.Tournament{},
	}

	for _, tournament := range tournaments {
		settings, err := tournament.GetSettings()
		if err != nil {
			log.Error().Err(err).Msg("failed to read tournament settings")
			apierrors.AbortWithError(http.StatusInternalServerError, "failed to read tournament settings", ctx)
			return
		}
		response.Tournaments = append(response.Tournaments, generated.Tournament{
			Id:         tournament.ID.String(),
			Name:       tournament.Name,
			Slug:       tournament.Slug,
			Type:       tournament.Type,
			Settings:   *settings,
			LocationId: tournament.LocationID.String(),
			LeagueId:   tournament.LeagueID.String(),
			CreatedAt:  utils.FormatTime(tournament.CreatedAt),
			UpdatedAt:  utils.FormatTime(tournament.UpdatedAt),
		})
	}

	ctx.JSON(http.StatusOK, response)
}
