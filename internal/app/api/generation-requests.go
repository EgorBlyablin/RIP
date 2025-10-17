package api

import (
	"errors"
	"net/http"
	"rip/internal/app/ds"
	"rip/internal/app/middlewares"
	"rip/internal/app/repositories"
	"rip/internal/app/services"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type GenerationRequestsApi struct {
	generationRequestsService *services.GenerationRequestsService
}

func NewGenerationRequestsApi(generationRequestsService *services.GenerationRequestsService) *GenerationRequestsApi {
	return &GenerationRequestsApi{
		generationRequestsService: generationRequestsService,
	}
}

func (api *GenerationRequestsApi) RegisterEndpoints(router *gin.RouterGroup) {
	router.GET("/", api.GetSentGenerationRequests)
	router.GET("/:generationRequestId", api.GetGenerationRequest)
	router.PUT("/:generationRequestId/close", api.CloseGenerationRequest)

	router.GET("/draft", api.GetDraftBriefInfo)
	router.POST("/draft/:turbineId", api.AddTurbineToDraft)
	router.PUT("/draft", api.UpdateDraftGenerationRequest)
	router.PUT("/draft/:turbineId", api.UpdateTurbineInDraft)
	router.DELETE("/draft/:turbineId", api.RemoveTurbineFromDraft)
	router.PUT("/draft/submit", api.SubmitDraftGenerationRequest)
	router.DELETE("/draft", api.DeleteDraftGenerationRequest)
}

func (api *GenerationRequestsApi) GetSentGenerationRequests(ctx *gin.Context) {
	userId := middlewares.GetUserId()

	generationRequestsFilter := repositories.GenerationRequestsFilter{}
	if err := ctx.BindQuery(&generationRequestsFilter); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(generationRequestsFilter); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	generationRequests, err := api.generationRequestsService.GetGenerationRequests(userId, generationRequestsFilter)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, generationRequests)
}

func (api *GenerationRequestsApi) GetGenerationRequest(ctx *gin.Context) {
	userId := middlewares.GetUserId()

	generationRequestsId, err := strconv.Atoi(ctx.Param("generationRequestId"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	generationRequest, err := api.generationRequestsService.GetGenerationRequest(uint(generationRequestsId))
	if err != nil {
		if errors.Is(err, repositories.ErrorGenerationRequestNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	if generationRequest.CreatedByID != userId {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	ctx.JSON(http.StatusOK, generationRequest)
}

func (api *GenerationRequestsApi) CloseGenerationRequest(ctx *gin.Context) {
	userId := middlewares.GetUserId()

	generationRequestsId, err := strconv.Atoi(ctx.Param("generationRequestId"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	type req struct {
		Status string `json:"status" binding:"required,oneof=completed rejected"`
	}
	closeGenerationRequest := req{}
	if err := ctx.Bind(&closeGenerationRequest); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	closedGenerationRequest, err := api.generationRequestsService.CloseGenerationRequest(uint(generationRequestsId), userId, closeGenerationRequest.Status)
	if err != nil {
		switch err {
		case services.ErrorGenerationRequestIncorrectStatus:
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		case services.ErrorUserHasNoAccess:
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": err.Error()})
			return
		case repositories.ErrorGenerationRequestNotFound:
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		case repositories.ErrorGenerationRequestCannotBeClosed:
			ctx.AbortWithStatusJSON(http.StatusConflict, gin.H{"message": err.Error()})
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, closedGenerationRequest)
}

func (api *GenerationRequestsApi) GetDraftBriefInfo(ctx *gin.Context) {
	userId := middlewares.GetUserId()

	generationRequestDraftBriefInfo, err := api.generationRequestsService.GetDraftBriefInfo(userId)
	if err != nil {
		if errors.Is(err, repositories.ErrorGenerationRequestNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, generationRequestDraftBriefInfo)
}

func (api *GenerationRequestsApi) UpdateDraftGenerationRequest(ctx *gin.Context) {
	userId := middlewares.GetUserId()

	generationRequestUpdates := ds.UpdateGenerationRequest{}
	if err := ctx.Bind(&generationRequestUpdates); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	updatedGenerationRequest, err := api.generationRequestsService.UpdateDraftGenerationRequest(userId, generationRequestUpdates)
	if err != nil {
		if errors.Is(err, services.ErrorUserHasNoAccess) {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": err.Error()})
			return
		}
		if errors.Is(err, repositories.ErrorGenerationRequestNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updatedGenerationRequest)
}

func (api *GenerationRequestsApi) AddTurbineToDraft(ctx *gin.Context) {
	userId := middlewares.GetUserId()

	turbinesId, err := strconv.Atoi(ctx.Param("turbineId"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = api.generationRequestsService.AddTurbineToDraft(userId, uint(turbinesId))
	if err != nil {
		if errors.Is(err, repositories.ErrorGenerationRequestTurbineAlreadyInDraft) {
			ctx.AbortWithStatusJSON(http.StatusNotModified, gin.H{"message": err.Error()})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

func (api *GenerationRequestsApi) UpdateTurbineInDraft(ctx *gin.Context) {
	userId := middlewares.GetUserId()

	turbineId, err := strconv.Atoi(ctx.Param("turbineId"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	turbineUpdates := ds.UpdateTurbineGenerationRequest{}
	if err := ctx.Bind(&turbineUpdates); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = api.generationRequestsService.UpdateTurbineInDraft(userId, uint(turbineId), turbineUpdates)
	if errors.Is(err, repositories.ErrorGenerationRequestTurbineNotFound) {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

func (api *GenerationRequestsApi) RemoveTurbineFromDraft(ctx *gin.Context) {
	userId := middlewares.GetUserId()

	turbineId, err := strconv.Atoi(ctx.Param("turbineId"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = api.generationRequestsService.RemoveTurbineFromDraft(userId, uint(turbineId))
	if errors.Is(err, repositories.ErrorGenerationRequestTurbineNotFound) {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	if errors.Is(err, repositories.ErrorGenerationRequestNotFound) {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

func (api *GenerationRequestsApi) SubmitDraftGenerationRequest(ctx *gin.Context) {
	userId := middlewares.GetUserId()

	err := api.generationRequestsService.SubmitDraftGenerationRequest(userId)
	if errors.Is(err, repositories.ErrorGenerationRequestNotFound) {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

func (api *GenerationRequestsApi) DeleteDraftGenerationRequest(ctx *gin.Context) {
	userId := middlewares.GetUserId()

	err := api.generationRequestsService.DeleteDraftGenerationRequest(userId)
	if err != nil {
		if errors.Is(err, repositories.ErrorGenerationRequestNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}
