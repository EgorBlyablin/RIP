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

func (r *GenerationCalculationRequestHandler) GetGenerationCalculationRequest(ctx *gin.Context) {
	generationCalculationIdParam := ctx.Param("generationCalculationId")

	generationCalculationIdSigned, err := strconv.Atoi(generationCalculationIdParam)
	if err != nil {
		logrus.Error(err)
	}

	generationCalculationRequest, err := r.TurbinesRepository.GetCalculationGenerationRequest(uint(generationCalculationIdSigned))
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "generation-calculation-request", gin.H{
		"generationCalculationRequest": generationCalculationRequest,
	})
}
