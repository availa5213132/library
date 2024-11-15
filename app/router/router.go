package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"tushuguanli/app/logic"
	"tushuguanli/app/middleware"
	"tushuguanli/app/model"
	"tushuguanli/app/tools"
)

func New() {
	r := gin.Default()
	r.LoadHTMLGlob("app/view/*")

	//管理员对 book_info 表进行操作   前端页面
	index := r.Group("")
	//index.Use(checkAdmin)
	{
		index.GET("/index", logic.Index)                           //页面
		index.GET("/book/info", logic.GetBookInfo)                 //显示图书详情
		index.GET("/book/list", logic.GetBook)                     //对书进行查询
		index.POST("/book/list", logic.AddBooks)                   //添加书
		index.StaticFS("/static", http.Dir("./app/static/images")) //引入图书封面
	}

	//图片上传下载
	img := r.Group("")
	{
		img.GET("/img", middleware.UploadHand)
		img.POST("/img", middleware.UploadHandler)
		img.POST("/mongo", middleware.UploadHandlerV1)

	}

	// 管理员登录模块
	admin := r.Group("/admin")
	{
		admin.GET("/login", logic.GetAdminLogin)
		admin.POST("/login", logic.DoAdminLogin)
	}

	//用户操作模块
	user := r.Group("/user")
	//user.Use(middleware.CheckUser)
	{
		user.POST("/create", logic.CreatUser)
		user.GET("/login", logic.GetUserLogin)
		user.POST("/login", logic.DoUserLogin)
		user.GET("/logout", logic.UserLogout)
		user.POST("/book/borrow", logic.BorrowBook)
		user.POST("/book/return", logic.ReturnBook)
		user.POST("/book/buy", logic.BuyBook)
		user.GET("/books/buy", logic.GetBuyBook)
		user.POST("/books/buy", logic.BuyBooks)
		user.GET("/users", logic.Check)
		user.POST("/users", logic.Check)
		user.GET("/wechat", logic.CheckSignature) //微信扫码登录
		user.GET("/wechat/login", logic.Redirect)
		user.GET("/wechat/callback", logic.Callback)
		user.GET("/wechat/vx", logic.VChat)

	}

	//对 book 表进行增删改查
	book := r.Group("/book")
	{
		book.GET("/book", logic.GetBooks)
		book.POST("/book", logic.AddBook)
		book.DELETE("/book", logic.DelBook)
		book.PUT("/book", logic.UpdateBook)
	}

	//验证码相关操作
	captcha := r.Group("")
	{
		//图片验证
		captcha.GET("/captcha", logic.GetCaptcha)
		captcha.POST("/captcha/verify", logic.VerifyCaptcha)

		//邮箱验证
		captcha.GET("/email", logic.GetEmail)
		captcha.POST("/email", logic.SendEmailCode)
		captcha.POST("/check/email", logic.VerifyEmailCode)

		//手机验证
		captcha.POST("/meg", logic.SendPhoneCode)
		captcha.POST("/check/meg", logic.VerifyPhoneCode)

	}

	if err := r.Run(":8080"); err != nil {
		fmt.Print("启动失败")
	}
}

func checkAdmin(context *gin.Context) {
	var name string
	var id int64 //TODO 存在一个bug
	values := model.GetSession(context)

	if v, ok := values["name"]; ok {
		name = v.(string)
	}
	if v, ok := values["id"]; ok {
		id = v.(int64)
	}
	if name == "" || id <= 0 {
		context.JSON(http.StatusUnauthorized, tools.NotLogin)
		context.Abort()
	}

	context.Next()
}
