package logic

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"tushuguanli/app/model"
	"tushuguanli/app/tools"
)

var page int

func Index(c *gin.Context) {
	page, _ = strconv.Atoi(c.Query("page"))
	c.HTML(http.StatusOK, "indexV1.tmpl", nil)
}

func GetBookInfo(c *gin.Context) {
	pageSize := 10
	pageNum := page
	if pageNum == 0 {
		pageNum = -1
	}
	ret := model.GetBooks(pageSize, pageNum)
	c.JSON(http.StatusOK, tools.ECode{
		Data: ret,
	})

}
func VChat(c *gin.Context) {
	c.HTML(http.StatusOK, "vxsaoma.tmpl", nil)
	//c.HTML(http.StatusOK, "vx1.tmpl", nil)

}
