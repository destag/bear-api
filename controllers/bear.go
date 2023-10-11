package controllers

import (
	"net/http"

	"github.com/destag/bear-api/models"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type BearController struct {
	DB *sqlx.DB
}

type createBearInput struct {
	Name string `json:"name"`
}

func (bc *BearController) ListBears(c *gin.Context) {
	var bears []models.Bear

	err := bc.DB.SelectContext(c.Request.Context(), &bears, "SELECT * FROM bears")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, bears)
}
