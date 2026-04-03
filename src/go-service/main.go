// @title           Items API (Go Gin)
// @version         1.0.0
// @description     REST API для управления каталогом товаров. Часть лабораторной работы №10 — сравнение FastAPI (Python) и Gin (Go).
// @contact.name    Евланичев Максим Юрьевич
// @contact.email   storhetmax@mail.ru
// @host            localhost:8080
// @BasePath        /
package main

import (
	"go-service/app"
	_ "go-service/docs"
)

func main() {
	r := app.SetupRouter()
	r.Run(":8080")
}
