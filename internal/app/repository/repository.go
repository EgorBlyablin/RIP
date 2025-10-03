package repository

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type TurbinesRepository struct {
	db *gorm.DB
}

func NewTurbineRepository() (*TurbinesRepository, error) {
	return &TurbinesRepository{}, nil
}

func New(dsn string) (*TurbinesRepository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // standard Go logger
			logger.Config{
				IgnoreRecordNotFoundError: true,
			},
		),
	}) // подключаемся к БД
	if err != nil {
		return nil, err
	}

	// Возвращаем объект TurbinesRepository с подключенной базой данных
	return &TurbinesRepository{
		db: db,
	}, nil
}
