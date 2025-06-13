package handler

import (
	"errors"
	"mfawzanid/warehouse-commerce/core/entity"
	"mfawzanid/warehouse-commerce/core/usecase"
	errorutil "mfawzanid/warehouse-commerce/utils/error"
	generalutil "mfawzanid/warehouse-commerce/utils/general"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthHandler interface {
	VerifyToken(next echo.HandlerFunc) echo.HandlerFunc
}

type authHandler struct {
	authUsecase usecase.AuthUsecaseInterface
}

func NewAuthHandler(authUsecase usecase.AuthUsecaseInterface) AuthHandler {
	return &authHandler{authUsecase}
}

func (h *authHandler) VerifyToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")

		if len(token) == 0 {
			return c.JSON(http.StatusUnauthorized, generalutil.MapAny{
				errorutil.Error: errorutil.CombineHTTPErrorMessage(http.StatusUnauthorized, errors.New("token is empty")),
			})
		}

		userId, err := h.authUsecase.VerifyToken(token)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, generalutil.MapAny{
				errorutil.Error: errorutil.CombineHTTPErrorMessage(http.StatusInternalServerError, err),
			})
		}

		c.Set(entity.ContextUserId, userId)

		return next(c)
	}
}
