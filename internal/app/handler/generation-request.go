package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (r *TurbineHandler) GetGenerationRequest(ctx *gin.Context) {
	generationRequestId, err := strconv.Atoi(ctx.Param("generationRequestId"))
	if err != nil {
		r.errorHandler(ctx, err)
		return
	}

	generationRequest, err := r.TurbinesRepository.GetGenerationRequest(uint(generationRequestId))
	if err != nil {
		r.errorHandler(ctx, err)
		return
	}

	ctx.HTML(http.StatusOK, "generation-request", gin.H{
		"generationRequest": generationRequest,
	})
}

func (r *TurbineHandler) AddTurbineToGenerationRequest(ctx *gin.Context) {
	turbineId, err := strconv.Atoi(ctx.Request.FormValue("turbineId"))
	if err != nil {
		r.errorHandler(ctx, err)
		return
	}

	generationRequest, err := r.TurbinesRepository.GetOrCreateUserDraftGenerationRequest(userId)
	if err != nil {
		r.errorHandler(ctx, err)
		return
	}

	if err := r.TurbinesRepository.AddTurbineToDraftGenerationRequest(uint(generationRequest.ID), uint(turbineId)); err != nil {
		r.errorHandler(ctx, err)
		return
	}

	ctx.Redirect(http.StatusSeeOther, "/")
}

func (r *TurbineHandler) DeleteGenerationRequest(ctx *gin.Context) {
	generationRequestId, err := strconv.Atoi(ctx.Param("generationRequestId"))
	if err != nil {
		r.errorHandler(ctx, err)
		return
	}

	if err := r.TurbinesRepository.DeleteDraftGenerationRequest(uint(generationRequestId)); err != nil {
		r.errorHandler(ctx, err)
		return
	}

	ctx.Redirect(http.StatusSeeOther, "/")
}
