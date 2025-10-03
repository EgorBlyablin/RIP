package handler

import (
	"net/http"
	"path/filepath"
	"rip/internal/app/repository"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type TurbineHandler struct {
	TurbinesRepository *repository.TurbinesRepository
}

func NewTurbineHandler(r *repository.TurbinesRepository) *TurbineHandler {
	return &TurbineHandler{TurbinesRepository: r}
}

const userId = 1

// RegisterHandler Функция, в которой мы отдельно регистрируем маршруты, чтобы не писать все в одном месте
func (h *TurbineHandler) RegisterHandler(router *gin.Engine) {
	router.GET("/", h.GetTurbines)
	router.GET("/turbines", h.GetTurbines)
	router.GET("/turbines/:turbineId", h.GetTurbine)

	router.GET("/generation-request/:generationRequestId", h.GetGenerationRequest)
	router.POST("/add-to-generation-request", h.AddTurbineToGenerationRequest)
	router.POST("/delete-generation-request/:generationRequestId", h.DeleteGenerationRequest)
}

// RegisterStatic То же самое, что и с маршрутами, регистрируем статику
func (h *TurbineHandler) RegisterStatic(router *gin.Engine) {
	generation_templates_files, _ := filepath.Glob("templates/*.*")
	generation_templates_subfiles, _ := filepath.Glob("templates/**/*.*")
	generation_templates_files = append(generation_templates_files, generation_templates_subfiles...)

	router.LoadHTMLFiles(generation_templates_files...) // сразу все файлы рекурсивно
	router.Static("/static", "./static")
}

// errorHandler для более удобного вывода ошибок
func (h *TurbineHandler) errorHandler(ctx *gin.Context, err error) {
	log.Error(err.Error())
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}
