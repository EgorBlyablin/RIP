package api

import (
	"errors"
	"net/http"
	"rip/internal/app/ds"
	"rip/internal/app/middlewares"
	"rip/internal/app/repositories"
	"rip/internal/app/services"

	"github.com/gin-gonic/gin"
)

type UsersApi struct {
	usersService *services.UsersService
}

func NewUsersApi(usersService *services.UsersService) *UsersApi {
	return &UsersApi{
		usersService: usersService,
	}
}

func (api *UsersApi) RegisterEndpoints(router *gin.RouterGroup) {
	router.POST("/", api.RegisterUser)
	router.GET("/", api.GetCurrentUser)
	router.PUT("/", api.UpdateCurrentUser)
	router.POST("/login", api.Login)
	router.POST("/logout", api.Logout)
}

// @Summary Регистрация нового пользователя
// @Description Создает нового пользователя с указанными учетными данными
// @Tags Пользователи
// @Accept json
// @Produce json
// @Param user body ds.CreateUser true "Данные нового пользователя"
// @Success 201 {object} ds.User "Созданный пользователь"
// @Failure 400 {object} map[string]string "Некорректный запрос"
// @Failure 409 {object} map[string]string "Логин уже занят"
// @Router /api/users/ [post]
func (api *UsersApi) RegisterUser(ctx *gin.Context) {
	user := ds.CreateUser{}
	if err := ctx.Bind(&user); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	createdUser, err := api.usersService.RegisterUser(user)
	if err != nil {
		if errors.Is(err, repositories.ErrorLoginIsTaken) {
			ctx.AbortWithStatusJSON(http.StatusConflict, gin.H{"message": err.Error()})
			return
		}

		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, createdUser)
}

// @Summary Получение информации о текущем пользователе
// @Description Возвращает информацию о текущем аутентифицированном пользователе
// @Tags Пользователи
// @Accept json
// @Produce json
// @Success 200 {object} ds.User "Информация о пользователе"
// @Failure 404 {object} map[string]string "Пользователь не найден"
// @Router /api/users/ [get]
func (api *UsersApi) GetCurrentUser(ctx *gin.Context) {
	userId := middlewares.GetUserId()

	user, err := api.usersService.GetUser(userId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

// @Summary Обновление информации о текущем пользователе
// @Description Обновляет информацию о текущем аутентифицированном пользователе
// @Tags Пользователи
// @Accept json
// @Produce json
// @Param user body ds.UpdateUser true "Данные для обновления"
// @Success 200 {object} ds.User "Обновленный пользователь"
// @Failure 400 {object} map[string]string "Некорректный запрос"
// @Failure 404 {object} map[string]string "Пользователь не найден"
// @Failure 409 {object} map[string]string "Логин уже занят"
// @Router /api/users/ [put]
func (api *UsersApi) UpdateCurrentUser(ctx *gin.Context) {
	userId := middlewares.GetUserId()

	user := ds.UpdateUser{}
	if err := ctx.Bind(&user); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	updatedUser, err := api.usersService.UpdateUser(userId, user)
	if err != nil {
		if errors.Is(err, repositories.ErrorLoginIsTaken) {
			ctx.AbortWithStatusJSON(http.StatusConflict, gin.H{"message": err.Error()})
			return
		}

		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, updatedUser)
}

// @Summary Аутентификация пользователя
// @Description Аутентифицирует пользователя и создает сессию
// @Tags Пользователи
// @Accept json
// @Produce json
// @Param credentials body ds.CreateUser true "Учетные данные пользователя (логин и пароль)"
// @Success 200 {object} map[string]string "Сообщение об успешной аутентификации"
// @Failure 400 {object} map[string]string "Некорректные учетные данные"
// @Router /api/users/login/ [post]
func (api *UsersApi) Login(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Logged in successfully"})
}

// @Summary Деавторизация пользователя
// @Description Завершает сессию текущего аутентифицированного пользователя
// @Tags Пользователи
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "Сообщение об успешной деавторизации"
// @Router /api/users/logout/ [post]
func (api *UsersApi) Logout(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
