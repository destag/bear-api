package main

import (
	"github.com/destag/bear-api/controller"
	"github.com/destag/bear-api/database"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func initDatabase() *database.Database {
	db := database.Connect(":memory:")
	db.Migrate()
	return db
}

func main() {
	db := initDatabase()

	r := gin.Default()

	pingController := controller.PingController{}
	bearController := controller.BearController{DB: db}

	r.GET("/ping", pingController.Ping)

	r.GET("/bears", bearController.ListBears)
	r.POST("/bears", bearController.CreateBear)

	r.Run()
}
