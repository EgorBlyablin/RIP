package handler

import (
	"net/http"
	"rip/internal/app/repository"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type GenerationCalculationRequestHandler struct {
	TurbinesRepository *repository.TurbinesRepository
}

func NewGenerationCalculationRequestHandler(r *repository.TurbinesRepository) *GenerationCalculationRequestHandler {
	return &GenerationCalculationRequestHandler{TurbinesRepository: r}
}

func (r *GenerationCalculationRequestHandler) GetRequest(ctx *gin.Context) {
	idParam := ctx.Param("id")

	idSigned, err := strconv.Atoi(idParam)
	if err != nil {
		logrus.Error(err)
	}

	request, err := r.TurbinesRepository.GetRequest(uint(idSigned))
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "request", gin.H{
		"request":  request,
		"terrains": repository.Terrains,
	})
}
