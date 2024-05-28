package estates

import (
	"api/internal/application"
	"api/internal/core"
	"io"
	"log"
	"strconv"

	"github.com/labstack/echo/v4"
)

type handler struct {
	useCase application.IEstates
}

func InitRoutes(api *echo.Group, authMiddleware echo.MiddlewareFunc, useCase application.IEstates) {
	h := &handler{useCase: useCase}

	router := api.Group("/estates")
	router.GET("", h.getAll)
	router.GET("/getOne", h.getOne)
	router.GET("/amenities", h.getAmenities)
	router.GET("/categories", h.getCategories)
	{
		private := router.Group("", authMiddleware)
		private.POST("", h.add)
		private.POST("/tempImage", h.uploadTempImage)
	}
}

func (h *handler) add(c echo.Context) error {
	userId := c.Get("userId").(int)

	var dto CreateEstateDto
	if err := c.Bind(&dto); err != nil {
		log.Println(err)
		return c.NoContent(400)
	}

	images, err := h.useCase.GetTempImages(dto.Images)
	if err != nil {
		log.Println(err)
		return c.NoContent(400)
	}

	estate := &core.Estate{
		Description: dto.Description,
		Images:      images,
		Amenities:   dto.Amenities,
		OwnerId:     userId,
		PriceNight:  dto.PriceNight,
		PriceWeek:   dto.PriceWeek,
		Area:        dto.Area,
		Rooms:       dto.Rooms,
		Showers:     dto.Showers,
		BabyRooms:   dto.BabyRooms,
		CategoryId:  dto.CategoryId,
		Address: core.Address{
			Number:   dto.Address.Number,
			Street:   dto.Address.Street,
			City:     dto.Address.City,
			District: dto.Address.District,
		},
	}

	id, err := h.useCase.Add(estate)
	if err != nil {
		log.Println(err)
		return c.NoContent(400)
	}

	return c.JSON(200, map[string]any{"id": id})
}

func (h *handler) uploadTempImage(c echo.Context) error {
	bytes, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.NoContent(400)
	}

	imageId, err := h.useCase.AddTempImage(bytes)
	if err != nil {
		log.Println(err)
		return c.NoContent(400)
	}

	return c.JSON(200, imageId)
}

func (h *handler) getAmenities(c echo.Context) error {
	items, err := h.useCase.GetAmenities()
	if err != nil {
		log.Println(err)
		return c.NoContent(400)
	}

	return c.JSON(200, items)
}
func (h *handler) getCategories(c echo.Context) error {
	items, err := h.useCase.GetCategories()
	if err != nil {
		log.Println(err)
		return c.NoContent(400)
	}

	return c.JSON(200, items)
}

func (h *handler) getAll(c echo.Context) error {
	estates, err := h.useCase.GetAll()
	if err != nil {
		log.Println(err)
		return c.NoContent(400)
	}

	return c.JSON(200, estates)
}

func (h *handler) getOne(c echo.Context) error {
	estateIdStr := c.QueryParam("id")

	estateId, err := strconv.Atoi(estateIdStr)
	if err != nil {
		return c.NoContent(400)
	}

	estate, err := h.useCase.GetById(estateId)
	if err != nil {
		log.Println(err)
		return c.NoContent(404)
	}

	return c.JSON(200, estate)
}
