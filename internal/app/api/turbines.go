package api

import (
	"errors"
	"net/http"
	"rip/internal/app/ds"
	"rip/internal/app/repositories"
	"rip/internal/app/services"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TurbinesApi struct {
	turbinesService *services.TurbinesService
}

func NewTurbinesApi(turbinesService *services.TurbinesService) *TurbinesApi {
	return &TurbinesApi{
		turbinesService: turbinesService,
	}
}

func (api *TurbinesApi) RegisterEndpoints(router *gin.RouterGroup) {
	router.GET("/", api.GetTurbines)

	router.POST("/", api.CreateTurbine)
	router.GET("/:turbineId", api.GetTurbine)
	router.PUT("/:turbineId", api.UpdateTurbine)
	router.POST("/:turbineId/upload-image", api.UploadTurbineImage)
	router.DELETE("/:turbineId", api.DeleteTurbine)
}

func (api *TurbinesApi) GetTurbines(ctx *gin.Context) {
	titleFilter := func() *string {
		titleFilterValue := ctx.Query("turbineTitle")
		if titleFilterValue != "" {
			return &titleFilterValue
		}
		return nil
	}()

	turbines, err := api.turbinesService.GetActiveTurbines(titleFilter)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, turbines)
}

func (api *TurbinesApi) GetTurbine(ctx *gin.Context) {
	turbineId, err := strconv.Atoi(ctx.Param("turbineId"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	turbine, err := api.turbinesService.GetActiveTurbine(uint(turbineId))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, turbine)
}

func (api *TurbinesApi) CreateTurbine(ctx *gin.Context) {
	turbine := ds.CreateTurbine{}

	if err := ctx.Bind(&turbine); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	createdTurbine, err := api.turbinesService.CreateTurbine(turbine)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, createdTurbine)
}

func (api *TurbinesApi) UpdateTurbine(ctx *gin.Context) {
	turbine := ds.UpdateTurbine{}
	if err := ctx.Bind(&turbine); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	turbineId, err := strconv.Atoi(ctx.Param("turbineId"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	updatedTurbine, err := api.turbinesService.UpdateTurbine(uint(turbineId), turbine)
	if err != nil {
		if errors.Is(err, repositories.ErrorTurbineIsDeleted) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updatedTurbine)
}

func (api *TurbinesApi) UploadTurbineImage(ctx *gin.Context) {
	turbineId, err := strconv.Atoi(ctx.Param("turbineId"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	_, err = api.turbinesService.GetActiveTurbine(uint(turbineId))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	if ctx.Request.ContentLength <= 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	updatedTurbine, err := api.turbinesService.UpdateTurbineImage(
		uint(turbineId),
		ctx.Request.Body,
		ctx.Request.ContentLength,
		ctx.ContentType(),
	)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updatedTurbine)
}

func (api *TurbinesApi) DeleteTurbine(ctx *gin.Context) {
	turbineId, err := strconv.Atoi(ctx.Param("turbineId"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	if err := api.turbinesService.DeleteTurbine(uint(turbineId)); err != nil {
		if errors.Is(err, repositories.ErrorTurbineIsDeleted) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}
