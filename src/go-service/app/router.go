package app

import (
	"net/http"
	"runtime"
	"strconv"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go-service/middleware"
)

// Item представляет товар в каталоге.
type Item struct {
	ID    int     `json:"id"    example:"1"`
	Name  string  `json:"name"  example:"Apple"`
	Price float64 `json:"price" example:"1.5"`
}

// ErrorResponse описывает ответ при ошибке.
type ErrorResponse struct {
	Error string `json:"error" example:"item not found"`
}

var defaultItems = []Item{
	{ID: 1, Name: "Apple", Price: 1.5},
	{ID: 2, Name: "Banana", Price: 0.75},
}

// SetupRouter создаёт и возвращает настроенный gin.Engine.
// Использует gin.New() + кастомный Logger middleware (не gin.Default()).
func SetupRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Logger())

	r.GET("/ping", pingHandler)
	r.GET("/items", getItemsHandler)
	r.GET("/items/:id", getItemByIDHandler)
	r.GET("/memory", memoryHandler)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}

// pingHandler godoc
// @Summary     Проверка работоспособности
// @Description Возвращает pong — используется для health-check
// @Tags        health
// @Produce     json
// @Success     200 {object} map[string]string
// @Router      /ping [get]
func pingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

// getItemsHandler godoc
// @Summary     Список всех товаров
// @Description Возвращает полный список товаров в каталоге
// @Tags        items
// @Produce     json
// @Success     200 {array} Item
// @Router      /items [get]
func getItemsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, defaultItems)
}

// MemoryStats содержит метрики потребления памяти Go-рантайма.
type MemoryStats struct {
	AllocMB      float64 `json:"alloc_mb"       example:"1.2"`
	TotalAllocMB float64 `json:"total_alloc_mb" example:"3.4"`
	SysMB        float64 `json:"sys_mb"         example:"8.0"`
	NumGC        uint32  `json:"num_gc"         example:"5"`
}

// memoryHandler godoc
// @Summary     Потребление памяти
// @Description Возвращает метрики памяти Go-рантайма (runtime.MemStats)
// @Tags        metrics
// @Produce     json
// @Success     200 {object} MemoryStats
// @Router      /memory [get]
func memoryHandler(c *gin.Context) {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	c.JSON(http.StatusOK, MemoryStats{
		AllocMB:      float64(ms.Alloc) / 1024 / 1024,
		TotalAllocMB: float64(ms.TotalAlloc) / 1024 / 1024,
		SysMB:        float64(ms.Sys) / 1024 / 1024,
		NumGC:        ms.NumGC,
	})
}

// getItemByIDHandler godoc
// @Summary     Товар по ID
// @Description Возвращает товар по его числовому идентификатору
// @Tags        items
// @Produce     json
// @Param       id  path     int  true  "ID товара"
// @Success     200 {object} Item
// @Failure     400 {object} ErrorResponse
// @Failure     404 {object} ErrorResponse
// @Router      /items/{id} [get]
func getItemByIDHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "id must be a number"})
		return
	}
	for _, item := range defaultItems {
		if item.ID == id {
			c.JSON(http.StatusOK, item)
			return
		}
	}
	c.JSON(http.StatusNotFound, ErrorResponse{Error: "item not found"})
}
