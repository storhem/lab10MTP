package main

import "go-service/app"

func main() {
	r := app.SetupRouter()
	r.Run(":8080")
}
