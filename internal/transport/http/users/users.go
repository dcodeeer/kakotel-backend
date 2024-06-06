package users

import (
	"api/internal/application"
	"api/internal/transport/websocket"
	"io"
	"log"
	"strconv"

	"github.com/labstack/echo/v4"
)

type handler struct {
	useCase application.IUsers
	ws      *websocket.WebSocket
}

func InitRoutes(api *echo.Group, middleware echo.MiddlewareFunc, useCase application.IUsers, ws *websocket.WebSocket) {
	h := &handler{useCase: useCase, ws: ws}

	router := api.Group("/users")
	router.GET("", h.getOne)
	router.POST("/signup", h.signUp)
	router.POST("/login", h.signIn)

	{
		private := router.Group("", middleware)
		private.GET("/getMe", h.getMe)
		private.PATCH("/update", h.update)
		private.PATCH("/changePassword", h.changePassword)
		private.PATCH("/uploadPhoto", h.uploadPhoto)
	}
}

func (h *handler) signUp(c echo.Context) error {
	var dto signUpDTO
	if err := c.Bind(&dto); err != nil {
		return c.String(400, "Bad Request")
	}
	if err := dto.Validate(); err != nil {
		return c.String(400, "Bad Request")
	}

	token, err := h.useCase.SignUp(dto.Email, dto.Password)
	if err != nil {
		return c.String(400, "Bad Request")
	}

	return c.JSON(200, map[string]string{`token`: token})
}

func (h *handler) signIn(c echo.Context) error {
	var dto signUpDTO
	if err := c.Bind(&dto); err != nil {
		return c.String(400, "Bad Request")
	}

	token, err := h.useCase.SignIn(dto.Email, dto.Password)
	if err != nil {
		return c.String(400, "Bad Request")
	}

	return c.JSON(200, map[string]string{`token`: token})
}

func (h *handler) getMe(c echo.Context) error {
	userId := c.Get("userId").(int)

	user, err := h.useCase.GetOneById(userId)
	if err != nil {
		return c.String(404, "BadRequest")
	}

	return c.JSON(200, user)
}

func (h *handler) getOne(c echo.Context) error {
	userIdQuery := c.QueryParam("id")

	userId, err := strconv.Atoi(userIdQuery)
	if err != nil {
		return c.NoContent(400)
	}

	user, err := h.useCase.GetOneInfo(userId)
	if err != nil {
		return c.NoContent(400)
	}

	user["last_seen"] = h.ws.IsUserOnline(userId)

	return c.JSON(200, user)
}

func (h *handler) update(c echo.Context) error {
	var dto updateDTO
	if err := c.Bind(&dto); err != nil {
		return c.String(400, "Bad Request")
	}

	userId := c.Get("userId").(int)

	user := dto.ToUser()
	user.ID = userId

	if err := h.useCase.Update(user); err != nil {
		return c.NoContent(400)
	}

	return c.NoContent(200)
}

func (h *handler) uploadPhoto(c echo.Context) error {
	userId := c.Get("userId").(int)

	bytes, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.NoContent(400)
	}

	path, err := h.useCase.UpdatePhoto(userId, bytes)
	if err != nil {
		log.Println(err)
		return c.NoContent(400)
	}

	return c.JSON(200, map[string]string{"path": path})
}

func (h *handler) changePassword(c echo.Context) error {
	var dto changePasswordDTO
	if err := c.Bind(&dto); err != nil {
		return c.NoContent(400)
	}

	userId := c.Get("userId").(int)

	err := h.useCase.ChangePassword(userId, dto.OldPassword, dto.NewPassword)
	if err != nil {
		log.Println(err)
		return c.NoContent(400)
	}

	return c.NoContent(200)
}
