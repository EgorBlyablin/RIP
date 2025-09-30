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

func StartGenerationCalculationServer() {
	log.Println("Starting server")

	repo, err := repository.NewTurbineRepository()
	if err != nil {
		logrus.Error("ошибка инициализации репозитория")
	}
	turbineHandler := handler.NewTurbineHandler(repo)
	generationCalculationRequestHandler := handler.NewGenerationCalculationRequestHandler(repo)

	r := gin.Default()

	r.SetFuncMap(template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
	})

	generation_calculation_templates_files, _ := filepath.Glob("templates/*.*")
	generation_calculation_templates_subfiles, _ := filepath.Glob("templates/**/*.*")
	generation_calculation_templates_files = append(generation_calculation_templates_files, generation_calculation_templates_subfiles...)

	r.LoadHTMLFiles(generation_calculation_templates_files...) // сразу все файлы рекурсивно
	r.Static("/static", "./static")

	r.GET("/", turbineHandler.GetTurbines)
	r.GET("/turbines", turbineHandler.GetTurbines)
	r.GET("/turbines/:turbineId", turbineHandler.GetTurbine)
	r.GET("/requests/:generationCalculationId", generationCalculationRequestHandler.GetGenerationCalculationRequest)

	r.Run()
	log.Println("Server down")
}
