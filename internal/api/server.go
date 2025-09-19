package api

import (
	"html/template"
	"log"
	"path/filepath"
	"rip/internal/app/handler"
	"rip/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func StartServer() {
	log.Println("Starting server")

	repo, err := repository.NewRepository()
	if err != nil {
		logrus.Error("ошибка инициализации репозитория")
	}
	turbineHandler := handler.NewTurbineHandler(repo)
	generationCalculationRequestHandler := handler.NewGenerationCalculationRequestHandler(repo)

	r := gin.Default()
	// добавляем наш html/шаблон

	r.SetFuncMap(template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
	})

	files, _ := filepath.Glob("templates/*.*")
	subfiles, _ := filepath.Glob("templates/**/*.*")
	files = append(files, subfiles...)

	r.LoadHTMLFiles(files...) // сразу все файлы рекурсивно
	r.Static("/static", "./static")
	// слева название папки, в которую выгрузится наша статика
	// справа путь к папке, в которой лежит статика

	r.GET("/", turbineHandler.GetTurbines)
	r.GET("/turbines", turbineHandler.GetTurbines)
	r.GET("/turbines/:id", turbineHandler.GetTurbine)
	r.GET("/requests/:id", generationCalculationRequestHandler.GetRequest)

	r.Run()
	log.Println("Server down")
}
