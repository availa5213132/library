package logic

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tushuguanli/app/model"
	"tushuguanli/app/tools"
)

func GetEmail(c *gin.Context) {
	c.HTML(200, "email.tmpl", nil)
}

func SendEmailCode(c *gin.Context) {
	email := c.PostForm("email")
	emailUtil := model.NewEmailUtil()
	emailUtil.LoadConfig()
	_, err := emailUtil.SendVerificationCode(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "发送验证码失败！"})
		return
	}
	c.JSON(http.StatusOK, tools.ECode{
		Code:    10001,
		Message: "发送验证码成功！",
	})
}

func VerifyEmailCode(c *gin.Context) {
	email := c.PostForm("email")
	key := model.GenerateVerificationCodeKey(email)
	code := c.PostForm("code")
	emailUtil := model.NewEmailUtil()
	emailUtil.LoadConfig()
	if emailUtil.VerifyEmailCode(email, key, code) == false {
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
