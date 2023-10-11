package main

import (
	"fmt"
	"net/http"

	"github.com/destag/bear-api/database"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Bear struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func initDatabase() *sqlx.DB {
	db := database.Connect(":memory:")
	database.Migrate(db)
	return db
}

func main() {
	db := initDatabase()

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/bears", func(c *gin.Context) {
		bears := []Bear{}
		fmt.Printf("%#v\n", bears)
		fmt.Println(bears)
		err := db.Select(&bears, "SELECT * FROM bears;")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err)
			return
		}

		fmt.Println(bears)
		fmt.Println([]Bear{})
		fmt.Printf("%T\n", bears)
		fmt.Printf("%+v\n", bears)
		fmt.Printf("%#v\n", bears)
		fmt.Printf("%#v\n", []Bear{})

		c.JSON(http.StatusOK, bears)
	})

	r.POST("/bears", func(c *gin.Context) {
		input := Bear{}
		if err := c.BindJSON(&input); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err)
			return
		}

		fmt.Printf("input: %+v\n", input)

		res, err := db.NamedExec("INSERT INTO bears (name) VALUES (:name)", input)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err)
			return
		}

		id, err := res.LastInsertId()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err)
			return
		}

		input.ID = uint(id)

		c.JSON(http.StatusOK, input)
	})
	r.Run()
}
