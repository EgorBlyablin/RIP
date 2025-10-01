package handler

import (
	"net/http"
	"rip/internal/app/ds"
	"rip/internal/app/repository"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type TurbineHandler struct {
	Repository *repository.TurbinesRepository
}

func NewTurbineHandler(r *repository.TurbinesRepository) *TurbineHandler {
	return &TurbineHandler{Repository: r}
}

func (h *TurbineHandler) GetTurbines(ctx *gin.Context) {
	var turbines []ds.Turbine
	var err error

	turbineTitleQuery := ctx.Query("turbine-title-query")

	if turbineTitleQuery == "" { // если поле поиска пусто, то просто получаем из репозитория все записи
		turbines, err = h.Repository.GetTurbines()
	} else {
		turbines, err = h.Repository.GetTurbinesByTitle(turbineTitleQuery) // в ином случае ищем заказ по заголовку
	}
	if err != nil {
		logrus.Error(err)
	}

	request, err := h.Repository.GetCalculationGenerationRequest(1)
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "turbines-list", gin.H{
		"turbines":                turbines,
		"turbineTitleQuery":       turbineTitleQuery,
		"calculationRequestItems": len(request.SelectedTurbines),
	})
}

func (h *TurbineHandler) GetTurbine(ctx *gin.Context) {
	turbineIdStr := ctx.Param("turbineId")
	turbineIdSigned, err := strconv.Atoi(turbineIdStr)
	turbineId := uint(turbineIdSigned)

	if err != nil {
		logrus.Error(err)
	}

	turbine, err := h.Repository.GetTurbine(turbineId)
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "turbine-details", gin.H{
		"turbine": turbine,
	})
}
