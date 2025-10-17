package api

import (
	"rip/internal/app/config"
	"rip/internal/app/repositories"
	"rip/internal/app/services"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TurbinesAppApi struct {
	turbinesAppConfig *config.Config
	turbinesAppDB     *gorm.DB
}

func NewTurbinesAppApi(turbinesAppConfig *config.Config, turbinesAppDB *gorm.DB) *TurbinesAppApi {
	return &TurbinesAppApi{
		turbinesAppConfig: turbinesAppConfig,
		turbinesAppDB:     turbinesAppDB,
	}
}

func (turbinesAppApi *TurbinesAppApi) RegisterEndpoints(router *gin.RouterGroup) {
	generationRequestsService := services.NewGenerationRequestsService(turbinesAppApi.turbinesAppDB)
	generationRequestApi := NewGenerationRequestsApi(generationRequestsService)
	generationRequestApi.RegisterEndpoints(router.Group("/generation-requests"))

	turbinesImagesS3, err := repositories.NewS3Repository(
		turbinesAppApi.turbinesAppConfig.S3.Host,
		turbinesAppApi.turbinesAppConfig.S3.Port,
		"turbines",
	)
	if err != nil {
		log.WithError(err).Error("Failed to initialize S3 repository for turbines images")
	}
	turbinesService := services.NewTurbinesService(turbinesAppApi.turbinesAppDB, turbinesImagesS3)
	turbinesApi := NewTurbinesApi(turbinesService)
	turbinesApi.RegisterEndpoints(router.Group("/turbines"))

	usersService := services.NewUsersService(turbinesAppApi.turbinesAppDB)
	usersApi := NewUsersApi(usersService)
	usersApi.RegisterEndpoints(router.Group("/users"))
}
