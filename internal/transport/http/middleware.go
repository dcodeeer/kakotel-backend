package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
)

func (h *handler) AuthenticateMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")

		user, err := h.useCase.Users.GetByToken(token)
		if err != nil {
			return c.String(401, "Unauthorized")
		}

		c.Set("userId", user.ID)

		return next(c)
	}
}

func (h *handler) GetUserID(c echo.Context) int {
	value := c.Get("userId")
	if value == nil {
		return 0
	}

	return value.(int)
}

func (h *handler) CORS(ctx *gin.Context) {
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	if ctx.Request.Method == "OPTIONS" {
		ctx.Writer.WriteHeader(http.StatusOK)
		return
	}
	ctx.Next()
}
