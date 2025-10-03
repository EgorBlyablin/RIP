package handler

import (
	"net/http"
	"rip/internal/app/ds"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (h *TurbineHandler) GetTurbines(ctx *gin.Context) {
	turbines := &[]ds.Turbine{}
	var err error

	turbineTitleQuery := ctx.Query("turbine-title-query")

	if turbineTitleQuery == "" { // если поле поиска пусто, то просто получаем из репозитория все записи
		turbines, err = h.TurbinesRepository.GetTurbines()
	} else {
		turbines, err = h.TurbinesRepository.GetTurbinesByTitle(turbineTitleQuery) // в ином случае ищем заказ по заголовку
	}
	if err != nil {
		log.Error(err)
	}

	ctx.HTML(http.StatusOK, "turbines-list", gin.H{
		"turbines":               turbines,
		"turbineTitleQuery":      turbineTitleQuery,
		"generationRequestItems": h.TurbinesRepository.GetGenerationRequestsCount(userId),
	})
}

func (h *TurbineHandler) GetTurbine(ctx *gin.Context) {
	turbineIdStr := ctx.Param("turbineId")
	turbineIdSigned, err := strconv.Atoi(turbineIdStr)
	turbineId := uint(turbineIdSigned)

	if err != nil {
		log.Error(err)
	}

	turbine, err := h.TurbinesRepository.GetTurbine(turbineId)
	if err != nil {
		log.Error(err)
	}

	ctx.HTML(http.StatusOK, "turbine-details", gin.H{
		"turbine": turbine,
	})
}
