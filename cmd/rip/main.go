package main

import (
	"fmt"

	"rip/internal/app/config"
	"rip/internal/app/dsn"
	"rip/internal/app/repositories"
	"rip/internal/pkg"

	"github.com/sirupsen/logrus"
)

// @title Turbines API
// @version 0.2.0
// @description API для работы с сервисом расчета генерации электроэнергии ветрогенераторами
// @license.name MIT License
// @contact.name Егор Бляблин
// @host localhost:8000
// @schemes http

func main() {
	config, err := config.NewConfig()
	if err != nil {
		logrus.Fatalf("error loading config: %v", err)
	}

	postgresDsn := dsn.FromEnv()
	fmt.Println(postgresDsn)

	db, errRep := repositories.NewTurbinesDB(postgresDsn)
	if errRep != nil {
		logrus.Fatalf("error initializing repository: %v", errRep)
	}

	app := &pkg.TurbinesApplication{}
	app.Run(config, db)
}
