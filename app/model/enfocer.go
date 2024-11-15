package model

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/sirupsen/logrus"
)

var (
	e   *casbin.Enforcer
	err error
)

func init() {
	e, err = casbin.NewEnforcer("./app/tools/model.conf", "./app/tools/policy.csv")
	if err != nil {
		//logrus.Fatal("load file failed, %v", err.Error())
	}
}
func CheckPermission(c *gin.Context, sub, obj, act string) {
	logrus.Infof("sub = %s obj = %s act = %s", sub, obj, act)
	ok, err := e.Enforce(sub, obj, act)
	if err != nil {
		logrus.Print("enforce failed %s", err.Error())
		c.String(http.StatusInternalServerError, "内部服务器错误")
		return
	}
	if !ok {
		logrus.Println("权限验证不通过")
		c.String(http.StatusOK, "权限验证不通过")
		return
	}
	logrus.Println("权限验证通过")
	c.String(http.StatusOK, "权限验证通过")
}
