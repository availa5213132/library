package logic

import (
	"github.com/gin-gonic/gin"
	"tushuguanli/app/model"
)

func Check(c *gin.Context) {
	sub := c.Query("username")
	obj := c.Request.URL.Path
	act := c.Request.Method
	model.CheckPermission(c, sub, obj, act)
}
