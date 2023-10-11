package controller

import (
	"net/http"

	"github.com/destag/bear-api/model"
	"github.com/gin-gonic/gin"
)

type BearDatabase interface {
	ListBears() ([]*model.Bear, error)
	CreateBear(name string) (*model.Bear, error)
}

type BearController struct {
	DB BearDatabase
}

type BearParams struct {
	Name string `json:"name"`
}

func (c *BearController) ListBears(ctx *gin.Context) {
	bears, err := c.DB.ListBears()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, bears)
}

func (c *BearController) CreateBear(ctx *gin.Context) {
	var input BearParams
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	bear, err := c.DB.CreateBear(input.Name)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusOK, bear)
}
