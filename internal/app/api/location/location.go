package location

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"net/http"
	apierrors "pinman/internal/app/api/errors"
	"pinman/internal/app/generated"
	"pinman/internal/clients/pinballmap"
	"pinman/internal/models"
	"pinman/internal/utils"
	"strings"
)

type Controller struct {
	DB       *gorm.DB
	pmClient pinballmap.ClientInterface
}

func NewController(db *gorm.DB) *Controller {
	return &Controller{
		DB:       db,
		pmClient: pinballmap.NewClient(),
	}
}

func NewControllerWithClient(db *gorm.DB, client pinballmap.ClientInterface) *Controller {
	return &Controller{
		DB:       db,
		pmClient: client,
	}
}

func (c *Controller) CreateLocation(ctx *gin.Context) {
	payload := &generated.LocationCreate{}

	if err := ctx.ShouldBindJSON(payload); err != nil {
		apierrors.AbortWithError(http.StatusBadRequest, err.Error(), ctx)
		return
	}

	pinballMapLocation, err := c.pmClient.GetLocation(payload.PinballMapId)
	if err != nil {
		log.Error().Err(err).Msgf("failed to get location with id '%d' from pinball map API", payload.PinballMapId)
		apierrors.AbortWithError(http.StatusInternalServerError, err.Error(), ctx)
		return
	}

	location := models.Location{
		Name:         pinballMapLocation.Name,
		Slug:         utils.Slugify(pinballMapLocation.Name, 20),
		PinballMapID: pinballMapLocation.ID,
		Address: fmt.Sprintf(
			"%s, %s, %s, %s",
			pinballMapLocation.Street, pinballMapLocation.City, pinballMapLocation.State, pinballMapLocation.Country,
		),
	}

	result := c.DB.Create(&location)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key") {
			// TODO: Add some retry logic to generate a different slug if this happens
			apierrors.AbortWithError(http.StatusConflict, "location with slug already exists", ctx)
			return
		} else {
			log.Error().Err(result.Error).Msg("failed to create location")
			apierrors.AbortWithError(http.StatusInternalServerError, "failed to create location", ctx)
			return
		}
	}

	ctx.JSON(http.StatusCreated, generated.LocationResponse{
		Location: generated.Location{
			Id:           location.ID.String(),
			Name:         location.Name,
			Slug:         location.Slug,
			Address:      location.Address,
			PinballMapId: location.PinballMapID,
			CreatedAt:    utils.FormatTime(location.CreatedAt),
			UpdatedAt:    utils.FormatTime(location.UpdatedAt),
		},
	})
}

func (c *Controller) ListLocations(ctx *gin.Context) {
	var dbLocations []models.Location
	result := c.DB.Find(&dbLocations)
	if result.Error != nil {
		log.Err(result.Error).Msg("failed to list locations")
		apierrors.AbortWithError(http.StatusInternalServerError, "failed to list locations", ctx)
		return
	}

	locations := make([]generated.Location, len(dbLocations))
	for i, location := range dbLocations {
		locations[i] = generated.Location{
			Id:           location.ID.String(),
			Name:         location.Name,
			Slug:         location.Slug,
			Address:      location.Address,
			PinballMapId: location.PinballMapID,
			CreatedAt:    utils.FormatTime(location.CreatedAt),
			UpdatedAt:    utils.FormatTime(location.UpdatedAt),
		}
	}

	ctx.JSON(http.StatusOK, generated.LocationListResponse{
		Locations: locations,
	})
}

func (c *Controller) GetLocationWithSlug(ctx *gin.Context, slug string) {
	var location models.Location
	result := c.DB.Where("slug = ?", slug).First(&location)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "not found") {
			apierrors.AbortWithError(http.StatusNotFound, "location not found", ctx)
			return
		} else {
			log.Error().Err(result.Error).Msg("failed to get location")
			apierrors.AbortWithError(http.StatusInternalServerError, "failed to get location", ctx)
			return
		}
	}

	ctx.JSON(http.StatusOK, generated.LocationResponse{
		Location: generated.Location{
			Id:           location.ID.String(),
			Name:         location.Name,
			Slug:         location.Slug,
			Address:      location.Address,
			PinballMapId: location.PinballMapID,
			CreatedAt:    utils.FormatTime(location.CreatedAt),
			UpdatedAt:    utils.FormatTime(location.UpdatedAt),
		},
	})
}
