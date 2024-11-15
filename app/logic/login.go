package logic

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"tushuguanli/app/model"
	"tushuguanli/app/tools"
)

type Admin struct {
	Name         string `json:"name" form:"name"`
	Password     string `json:"password" form:"password"`
	CaptchaId    string `json:"captcha_id" form:"captcha_id"`
	CaptchaValue string `json:"captcha_value" form:"captcha_value"`
}
type User struct {
	Name         string `json:"name" form:"name"`
	Password     string `json:"password" form:"password"`
	CaptchaId    string `json:"captcha_id" form:"captcha_id"`
	CaptchaValue string `json:"captcha_value" form:"captcha_value"`
}

func GetAdminLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "login.tmpl", nil)
}
func GetUserLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "userlogin.tmpl", nil)
}

// @Summary 管理员登录
// @Description 执行管理员登录操作
// @Tags Admin
// @Accept json

// @Produce json
// @Param admin body Admin true "管理员登录信息"
// @Success 200 {object} ECode "登录成功"
// @Failure 200 {object} ECode "登录失败"
// @Router /admin/login [post]
func DoAdminLogin(c *gin.Context) {
	var admin Admin
	if err := c.ShouldBind(&admin); err != nil {
		c.JSON(http.StatusOK, tools.ECode{
			Message: err.Error(), //这里有风险
		})
	}

	if !tools.CaptchaVerify(tools.CaptchaData{
		CaptchaId: admin.CaptchaId,
		Data:      admin.CaptchaValue,
	}) {
		c.JSON(http.StatusOK, tools.ECode{
			Code:    10010,
			Message: "验证码校验失败!",
		})
		return
	}

	ret := model.GetAdmin(admin.Name)
	if ret.Id < 1 /*|| ret.Password != encryptV1(user.Password)*/ {
		c.JSON(http.StatusOK, tools.UserErr)
		return
	}

	//context.SetCookie("name", user.Name, 3600, "/", "", true, false)
	//context.SetCookie("Id", fmt.Sprint(ret.Id), 3600, "/", "", true, false)

	_ = model.SetSession(c, admin.Name, ret.Id)
	c.JSON(http.StatusOK, tools.ECode{
		Message: "登录成功",
	})
}

func DoUserLogin(c *gin.Context) {

	var user User
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusOK, tools.ECode{
			Code:    10001,
			Message: err.Error(), //这里有风险
		})
	}

	if !tools.CaptchaVerify(tools.CaptchaData{
		CaptchaId: user.CaptchaId,
		Data:      user.CaptchaValue,
	}) {
		c.JSON(http.StatusOK, tools.ECode{
			Code:    10002,
			Message: "验证码校验失败!",
		})
		return
	}

	ret := model.GetUser(user.Name)
	//if ret.Id < 1 || ret.Password != encryptV1(user.Password) {
	//	c.JSON(http.StatusOK, tools.ECode{
	//		Code:    10001,
	//		Message: "账号密码错误",
	//	})
	//	return
	//}

	if ret.Id < 1 || ret.Password != user.Password {
		c.JSON(http.StatusOK, tools.ECode{
			Code:    10001,
			Message: "账号密码错误",
		})
		return
	}
	//context.SetCookie("name", user.Name, 3600, "/", "", true, false)
	//context.SetCookie("Id", fmt.Sprint(ret.Id), 3600, "/", "", true, false)

	token, _ := model.GetJwt(ret.Id, user.Name)
	c.SetCookie("auth", fmt.Sprint(ret.Id), 3600, "/", "", true, false)
	c.JSON(http.StatusOK, tools.ECode{
		Message: "登录成功",
		Data:    token,
	})
	return
}

func UserLogout(c *gin.Context) {
	var name string
	model.ClearJWTMap(name)
	c.JSON(200, "退出")
}
