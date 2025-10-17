package pkg

import (
	"fmt"

	"rip/internal/app/api"
	"rip/internal/app/config"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TurbinesApplication struct{}

func (app *TurbinesApplication) Run(config *config.Config, db *gorm.DB) {
	logrus.Info("Server start up")

	engine := gin.Default()
	turbinesApiRouter := engine.Group("/api")

	turbinesApi := api.NewTurbinesAppApi(config, db)
	turbinesApi.RegisterEndpoints(turbinesApiRouter)

	serverAddress := fmt.Sprintf("%s:%d", config.Service.Host, config.Service.Port)
	if err := engine.Run(serverAddress); err != nil {
		logrus.Fatal(err)
	}

	logrus.Info("Server down")
}
