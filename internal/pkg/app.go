package pkg

import (
	"fmt"

	"rip/internal/app/api"
	"rip/internal/app/config"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"

	_ "rip/docs"
)

type TurbinesApplication struct{}

func (a *TurbinesApplication) Run(config *config.Config, db *gorm.DB) {
	logrus.Info("Server start up")

	engine := gin.Default()
	router := engine.Group("/api")

	turbinesApi := api.NewTurbinesAppApi(config, db)
	turbinesApi.RegisterEndpoints(router)

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	serverAddress := fmt.Sprintf("%s:%d", config.Service.Host, config.Service.Port)
	if err := engine.Run(serverAddress); err != nil {
		logrus.Fatal(err)
	}

	logrus.Info("Server down")
}
