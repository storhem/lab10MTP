package app

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go-service/middleware"
)

// Item представляет товар в каталоге.
type Item struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
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

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	r.GET("/items", func(c *gin.Context) {
		c.JSON(http.StatusOK, defaultItems)
	})

	r.GET("/items/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		for _, item := range defaultItems {
			if fmt.Sprintf("%d", item.ID) == idStr {
				c.JSON(http.StatusOK, item)
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
	})

	return r
}
