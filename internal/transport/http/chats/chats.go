package chats

import (
	"api/internal/application"
	"log"
	"strconv"

	"github.com/labstack/echo/v4"
)

type handler struct {
	useCase application.IChats
}

func InitRoutes(api *echo.Group, middleware echo.MiddlewareFunc, useCase application.IChats) {
	h := &handler{useCase: useCase}

	router := api.Group("/chats", middleware)
	router.GET("", h.getAll)
	router.GET("/messages", h.getMessages)
}

func (h *handler) getAll(c echo.Context) error {
	userId := c.Get("userId").(int)

	chats, err := h.useCase.GetAll(userId)
	if err != nil {
		log.Println(err)
		return c.NoContent(400)
	}

	return c.JSON(200, chats)
}

func (h *handler) getMessages(c echo.Context) error {
	userId := c.Get("userId").(int)

	chatId, err := strconv.Atoi(c.QueryParam("chatId"))
	if err != nil {
		return c.NoContent(400)
	}

	messages, err := h.useCase.GetMessages(userId, chatId)
	if err != nil {
		log.Println(err)
		return c.NoContent(400)
	}

	return c.JSON(200, messages)
}
