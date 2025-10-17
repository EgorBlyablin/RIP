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

func (api *UsersApi) GetCurrentUser(ctx *gin.Context) {
	userId := middlewares.GetUserId()

	user, err := api.usersService.GetUser(userId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

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

func (api *UsersApi) Login(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Logged in successfully"})
}

func (api *UsersApi) Logout(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
