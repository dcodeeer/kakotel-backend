package http

import (
	"api/internal/application"
	"api/internal/transport/http/chats"
	"api/internal/transport/http/estates"
	"api/internal/transport/http/users"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type handler struct {
	useCase *application.UseCase
	router  *echo.Echo
}

func New(useCase *application.UseCase) *handler {
	return &handler{
		useCase: useCase,
		router:  echo.New(),
	}
}

func (h *handler) HandleWS(path string, callback http.Handler) {
	h.router.GET(path, func(c echo.Context) error {
		callback.ServeHTTP(c.Response(), c.Request())
		return nil
	})
}

func (h *handler) Run(addr string) error {
	srv := &http.Server{
		Addr:         addr,
		Handler:      h.getRouter(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	return srv.ListenAndServe()
}

func (h *handler) getRouter() http.Handler {
	h.router.Use(middleware.CORS())

	h.router.Static("/uploads", "./uploads")

	basePath := h.router.Group("/api/v1")

	users.InitRoutes(basePath, h.AuthenticateMiddleware, h.useCase.Users)
	chats.InitRoutes(basePath, h.AuthenticateMiddleware, h.useCase.Chats)
	estates.InitRoutes(basePath, h.AuthenticateMiddleware, h.useCase.Estates)

	return h.router
}
