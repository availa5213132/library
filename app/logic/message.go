package logic

import (
	"github.com/gin-gonic/gin"
	"net/http"
	_ "net/http"
	_ "time"
	"tushuguanli/app/model"
	"tushuguanli/app/tools"
)

func VerifyPhoneCode(c *gin.Context) {
	phone := c.PostForm("phone")
	key := model.GenerateVerificationCodeKey(phone)
	code := c.PostForm("code")
	if model.VerifyPhoneCode(phone, key, code) == false {
		c.JSON(http.StatusInternalServerError, tools.ECode{
			Code:    10001,
			Message: "验证失败！",
		})
	} else {
		c.JSON(200, tools.ECode{
			Code:    10001,
			Message: "验证成功！",
		})
	}

}
