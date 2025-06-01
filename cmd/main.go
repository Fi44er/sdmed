package main

import (
	_ "github.com/Fi44er/sdmed/docs"
	"github.com/Fi44er/sdmed/internal/app"
	_ "github.com/lib/pq"
)

// @title				sdmedik API
// @version		1.0
// @description	Swagger docs from sdmedik backend
// @host			127.0.0.1:8080
// @BasePath		/api/
func main() {
	a := app.NewApp()
	err := a.Run()
	if err != nil {
		panic(err)
	}
}
