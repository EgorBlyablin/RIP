package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (r *TurbineHandler) GetGenerationRequest(ctx *gin.Context) {
	generationRequest, err := r.TurbinesRepository.GetGenerationRequest(userId)
	if err != nil {
		log.Error(err)
	}

	ctx.HTML(http.StatusOK, "generation-request", gin.H{
		"generationRequest": generationRequest,
	})
}

func (r *TurbineHandler) AddTurbineToGenerationRequest(ctx *gin.Context) {
	turbineId, err := strconv.Atoi(ctx.Request.FormValue("turbineId"))
	if err != nil {
		log.Error(err)
	}

	if err := r.TurbinesRepository.AddTurbineToGenerationRequest(userId, uint(turbineId)); err != nil {
		log.Error(err)
	}

	ctx.Redirect(http.StatusSeeOther, "/")
}

func (r *TurbineHandler) DeleteGenerationRequest(ctx *gin.Context) {
	if err := r.TurbinesRepository.DeleteGenerationRequest(userId); err != nil {
		log.WithError(err).Error("Failed to delete generation request")
	}

	ctx.Redirect(http.StatusSeeOther, "/")
}
