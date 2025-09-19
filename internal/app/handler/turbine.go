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
	Repository *repository.Repository
}

func NewTurbineHandler(r *repository.Repository) *TurbineHandler {
	return &TurbineHandler{Repository: r}
}

func (h *TurbineHandler) GetTurbines(ctx *gin.Context) {
	var turbines []ds.Turbine
	var err error

	searchQuery := ctx.Query("query") // получаем значение из поля поиска

	if searchQuery == "" { // если поле поиска пусто, то просто получаем из репозитория все записи
		turbines, err = h.Repository.GetTurbines()
	} else {
		turbines, err = h.Repository.GetTurbinesByTitle(searchQuery) // в ином случае ищем заказ по заголовку
	}
	if err != nil {
		logrus.Error(err)
	}

	request, err := h.Repository.GetRequest(1)
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "list", gin.H{
		"turbines":     turbines,
		"query":        searchQuery,
		"requestItems": len(request.SelectedTurbines),
	})
}

func (h *TurbineHandler) GetTurbine(ctx *gin.Context) {
	idStr := ctx.Param("id") // получаем id заказа из урла (то есть из /turbine/:id)
	// через двоеточие мы указываем параметры, которые потом сможем считать через функцию выше
	idSigned, err := strconv.Atoi(idStr) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
	id := uint(idSigned)

	if err != nil {
		logrus.Error(err)
	}

	turbine, err := h.Repository.GetTurbine(id)
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "details", gin.H{
		"turbine": turbine,
	})
}
