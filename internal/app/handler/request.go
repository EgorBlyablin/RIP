package handler

import (
	"net/http"
	"rip/internal/app/repository"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type RequestHandler struct {
	Repository *repository.Repository
}

func NewRequestHandler(r *repository.Repository) *RequestHandler {
	return &RequestHandler{Repository: r}
}

func (r *RequestHandler) GetRequest(ctx *gin.Context) {
	idParam := ctx.Param("id")

	idSigned, err := strconv.Atoi(idParam)
	if err != nil {
		logrus.Error(err)
	}

	request, err := r.Repository.GetRequest(uint(idSigned))
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "request", gin.H{
		"request":  request,
		"terrains": repository.Terrains,
	})
}
