package middlewares

import (
	"net/http"
	"rip/internal/app/ds"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

const jwtPrefix = "Bearer "

type UserMiddlewares struct {
	Token string
}

func NewUserMiddlewares(token string) UserMiddlewares {
	return UserMiddlewares{
		Token: token,
	}
}

func (u *UserMiddlewares) WithAuth(ctx *gin.Context) {
	jwtStr := ctx.GetHeader("Authorization")
	if !strings.HasPrefix(jwtStr, jwtPrefix) {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	jwtStr = jwtStr[len(jwtPrefix):]

	jwt, err := jwt.ParseWithClaims(jwtStr, &ds.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(u.Token), nil
	})
	if err != nil {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	jwtData := jwt.Claims.(*ds.JWTClaims)
	ctx.Set(ds.UserIDKey, jwtData.UserID)
	ctx.Next()
}

func (u *UserMiddlewares) WithModeratorAccess(ctx *gin.Context) {
	jwtStr := ctx.GetHeader("Authorization")
	if !strings.HasPrefix(jwtStr, jwtPrefix) {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	jwtStr = jwtStr[len(jwtPrefix):]

	jwt, _ := jwt.ParseWithClaims(jwtStr, &ds.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(u.Token), nil
	})

	jwtData := jwt.Claims.(*ds.JWTClaims)
	if !jwtData.IsModerator {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	ctx.Set(ds.UserIDKey, jwtData.UserID)
	ctx.Set(ds.UserIsModeratorKey, true)
	ctx.Next()
}
