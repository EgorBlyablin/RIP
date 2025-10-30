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

// @Summary Список заявок с фильтрацией
// @Description Возвращает список заявок пользователя с фильтрацией по дате и статусу (кроме удаленных и черновиков)
// @Tags Заявки расчета выработки
// @Accept json
// @Produce json
// @Param filter path repositories.GenerationRequestsFilter false "Фильтр заявок"
// @Success 200 {array} ds.GenerationRequest "Список заявок"
// @Failure 400 {object} map[string]string "Некорректный запрос"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /api/generation-requests/ [get]
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

	for i := range generationRequests {
		grCount := uint(len(*generationRequests[i].TurbineGenerationRequests))
		generationRequests[i].TurbineGenerationRequestsCount = &grCount
		generationRequests[i].TurbineGenerationRequests = nil
	}

	ctx.JSON(http.StatusOK, generationRequests)
}

// @Summary Получение информации о заявке
// @Description Возвращает информацию о конкретной заявке с деталями и связанными турбинами
// @Tags Заявки расчета выработки
// @Accept json
// @Produce json
// @Param generationRequestId path int true "ID заявки"
// @Success 200 {object} ds.GenerationRequest "Информация о заявке"
// @Failure 403 {object} map[string]string "Доступ запрещен"
// @Failure 404 {object} map[string]string "Заявка не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /api/generation-requests/{generationRequestId}/ [get]
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

	grCount := uint(len(*generationRequest.TurbineGenerationRequests))
	generationRequest.TurbineGenerationRequestsCount = &grCount

	ctx.JSON(http.StatusOK, generationRequest)
}

// @Summary Завершение/отклонение заявки модератором
// @Description Модератор завершает или отклоняет заявку (только для модераторов)
// @Tags Заявки расчета выработки
// @Accept json
// @Produce json
// @Param generationRequestId path int true "ID заявки"
// @Param request body api.CloseGenerationRequest.req true "Статус: completed или rejected"
// @Success 200 {object} ds.GenerationRequest "Обновленная заявка"
// @Failure 400 {object} map[string]string "Некорректный запрос"
// @Failure 403 {object} map[string]string "Доступ запрещен"
// @Failure 404 {object} map[string]string "Заявка не найдена"
// @Failure 409 {object} map[string]string "Заявка не может быть закрыта"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /api/generation-requests/{generationRequestId}/close/ [put]
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

// @Summary Получение информации о черновике заявки
// @Description Возвращает информацию о черновике заявки текущего пользователя
// @Tags Заявки расчета выработки
// @Accept json
// @Produce json
// @Success 200 {object} ds.GenerationRequest "Информация о черновике"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /api/generation-requests/draft/ [get]
func (api *GenerationRequestsApi) GetDraftBriefInfo(ctx *gin.Context) {
	userId := middlewares.GetUserId()

	generationRequestDraftBriefInfo, err := api.generationRequestsService.GetDraftBriefInfo(userId)
	if err != nil {
		if errors.Is(err, repositories.ErrorGenerationRequestNotFound) {
			ctx.AbortWithStatusJSON(http.StatusOK, ds.DraftGenerationRequestsBriefInfo{
				GenerationRequestId: 0,
				TurbinesCount:       0,
			})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, generationRequestDraftBriefInfo)
}

// @Summary Обновление черновика заявки
// @Description Обновляет поля черновика заявки текущего пользователя
// @Tags Заявки расчета выработки
// @Accept json
// @Produce json
// @Param request body ds.UpdateGenerationRequest true "Данные для обновления"
// @Success 200 {object} ds.GenerationRequest "Обновленный черновик"
// @Failure 400 {object} map[string]string "Некорректный запрос"
// @Failure 403 {object} map[string]string "Доступ запрещен"
// @Failure 404 {object} map[string]string "Черновик не найден"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /api/generation-requests/draft/ [put]
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

// @Summary Добавление турбины в черновик
// @Description Добавляет турбину в черновик заявки текущего пользователя
// @Tags Заявки расчета выработки
// @Accept json
// @Produce json
// @Param turbineId path int true "ID турбины"
// @Success 200 {object} map[string]string "Турбина добавлена"
// @Failure 400 {object} map[string]string "Некорректный запрос"
// @Failure 304 {object} map[string]string "Турбина уже в черновике"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /api/generation-requests/draft/{turbineId}/ [post]
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

// @Summary Обновление параметров турбины в черновике
// @Description Обновляет параметры (avg_velocity, alpha) турбины в черновике
// @Tags Заявки расчета выработки
// @Accept json
// @Produce json
// @Param turbineId path int true "ID турбины"
// @Param updates body ds.UpdateTurbineGenerationRequest true "Параметры для обновления"
// @Success 200 {object} map[string]string "Параметры обновлены"
// @Failure 400 {object} map[string]string "Некорректный запрос"
// @Failure 404 {object} map[string]string "Турбина не найдена в черновике"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /api/generation-requests/draft/{turbineId}/ [put]
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

// @Summary Удаление турбины из черновика
// @Description Удаляет турбину из черновика заявки текущего пользователя
// @Tags Заявки расчета выработки
// @Accept json
// @Produce json
// @Param turbineId path int true "ID турбины"
// @Success 200 {object} map[string]string "Турбина удалена"
// @Failure 400 {object} map[string]string "Некорректный запрос"
// @Failure 404 {object} map[string]string "Турбина или черновик не найден"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /api/generation-requests/draft/{turbineId}/ [delete]
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

// @Summary Отправка черновика заявки
// @Description Формирует черновик заявки и отправляет на рассмотрение
// @Tags Заявки расчета выработки
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "Черновик отправлен"
// @Failure 404 {object} map[string]string "Черновик не найден"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /api/generation-requests/draft/submit/ [put]
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

// @Summary Удаление черновика заявки
// @Description Удаляет черновик заявки текущего пользователя
// @Tags Заявки расчета выработки
// @Accept json
// @Produce json
// @Success 204 "Черновик удален"
// @Failure 404 {object} map[string]string "Черновик не найден"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /api/generation-requests/draft/ [delete]
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
