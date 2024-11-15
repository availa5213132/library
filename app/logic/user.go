package logic

import (
	"context"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/alipay"
	"github.com/goccy/go-json"
	"github.com/golang-jwt/jwt/v5"

	"github.com/pkg/errors"
	"github.com/skip2/go-qrcode"
	"gorm.io/gorm"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"time"
	"tushuguanli/app/config"
	"tushuguanli/app/model"
	"tushuguanli/app/tools"
)

type CUser struct {
	Name         string `json:"name" form:"name"`
	Password     string `json:"password" form:"password"`
	Password2    string `json:"password_2" form:"password_2"`
	Phone        string `json:"phone" form:"phone"`
	CaptchaId    string `json:"captcha_id" form:"captcha_id"`
	CaptchaValue string `json:"captcha_value" form:"captcha_value"`
}

func CreatUser(context *gin.Context) {
	var user CUser
	if err := context.ShouldBind(&user); err != nil {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10001,
			Message: err.Error(),
		})
		return
	}

	//encrypt(user.Password)
	//encryptV1(user.Password)
	//encryptV2(user.Password)
	//return

	if user.Name == "" || user.Password == "" || user.Password2 == "" {

		context.JSON(http.StatusOK, tools.ParamErr)
		return
	}
	//校验密码
	if user.Password != user.Password2 {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10003,
			Message: "两次密码不相同！",
		})
		return
	}

	//校验用户是否存在，这种写法十分不安全，有并发安全
	if oldUser := model.GetUser(user.Name); oldUser.Id > 0 {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10004,
			Message: "用户名已存在",
		})
		return
	}

	nameLen := len(user.Name)
	passwordLen := len(user.Password)
	if nameLen > 16 || nameLen < 8 || passwordLen > 16 || passwordLen < 8 {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10005,
			Message: "账号或密码大于8小于16！",
		})
		return
	}
	//regexp  go语言自带正则表达式的包
	regex := regexp.MustCompile(`^[0-9]+$`)
	if regex.MatchString(user.Password) {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10006,
			Message: "密码不能为纯数字",
		})
		return
	}

	newUser := model.User{
		Name:        user.Name,
		Password:    tools.EncryptV1(user.Password),
		Phone:       tools.EncryptV1(user.Phone),
		CreatedTime: time.Now(),
		UpdatedTime: time.Now(),
		Uid:         tools.GetUid(),
	}
	if err := model.CreatUser(&newUser); err != nil {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10007,
			Message: "新用户创建失败",
		})
		return
	}
	context.JSON(http.StatusOK, tools.OK)
	return
}

func checkXYZ(c *gin.Context) bool {
	//拿到IP和UA
	ip := c.ClientIP() //用户的ip
	ua := c.GetHeader("user-agent")
	fmt.Printf("ip:%s\nua:%s\n", ip, ua)
	//转下MD5
	hash := md5.New()
	hash.Write([]byte(ip + ua))
	hashBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)
	//校验是否被ban
	flag, _ := model.Rdb.Get(c, "ban-"+hashString).Bool()
	if flag {
		return false
	}

	i, _ := model.Rdb.Get(c, "xyz-"+hashString).Int()
	fmt.Printf("i:%d\n", i)
	if i > 5 {
		model.Rdb.SetEx(c, "ban-"+hashString, true, 30*time.Minute)
		return false
	}

	model.Rdb.Incr(c, "xyz-"+hashString)
	fmt.Println(i)
	model.Rdb.Expire(c, "xyz-"+hashString, 50*time.Minute)
	return true
}

func GetCaptcha(c *gin.Context) {
	if !checkXYZ(c) {
		c.JSON(http.StatusOK, tools.ECode{
			Code:    10005,
			Message: "您的手速真的是太快了！",
		})
		return
	}
	captcha, err := tools.CaptchaGenerate()
	if err != nil {
		c.JSON(http.StatusOK, tools.ECode{
			Code:    10005,
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, tools.ECode{

		Data: captcha,
	})
}

func VerifyCaptcha(context *gin.Context) {
	var param tools.CaptchaData
	if err := context.ShouldBind(&param); err != nil {
		context.JSON(http.StatusOK, tools.ParamErr)
		return
	}
	fmt.Printf("参数为：%+v", param)
	if !tools.CaptchaVerify(param) {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10008,
			Message: "验证失败",
		})
		return
	}
	context.JSON(http.StatusOK, tools.OK)
}

func BorrowBook(c *gin.Context) {
	//获取用户信息
	uidStr := c.Query("uid")
	//获取图书ID
	idStr := c.Query("id")
	if idStr == "" || idStr == "0" {
		c.JSON(http.StatusOK, tools.ParamErr)
		return
	}
	id, _ := strconv.ParseInt(idStr, 10, 64)
	uid, _ := strconv.ParseInt(uidStr, 10, 64)
	//执行借书逻辑
	err := model.BorrowBook(uid, id)
	if err != nil {
		c.JSON(http.StatusOK, tools.ECode{
			Code:    10002,
			Message: err.Error(),
		})
		return
	}
	//返回成功
	c.JSON(http.StatusOK, tools.OK)
}

func ReturnBook(c *gin.Context) {
	//获取用户信息
	uidStr := c.Query("uid")
	//获取图书ID
	idStr := c.Query("id")
	if idStr == "" || idStr == "0" {
		c.JSON(http.StatusOK, tools.ParamErr)
		return
	}
	id, _ := strconv.ParseInt(idStr, 10, 64)
	uid, _ := strconv.ParseInt(uidStr, 10, 64)
	//执行借书逻辑
	err := model.ReturnBook(uid, id)
	if err != nil {
		c.JSON(http.StatusOK, tools.ECode{
			Code:    10002,
			Message: err.Error(),
		})
		return
	}
	//返回成功
	c.JSON(http.StatusOK, tools.OK)
}

func BuyBook(c *gin.Context) {
	userId, _ := strconv.ParseInt(c.PostForm("id"), 10, 64)
	bookId, _ := strconv.ParseInt(c.PostForm("uid"), 10, 64)
	num, _ := strconv.ParseInt(c.PostForm("num"), 10, 64)
	err := model.BuyBook(userId, bookId, num)
	if err != nil {
		c.JSON(200, tools.ECode{Message: err.Error()})
		return
	}
	c.JSON(200, tools.ECode{
		Message: "购买成功！",
	})
}

func GetBuyBook(c *gin.Context) {
	c.HTML(200, "buy_book.tmpl", nil)
}
func BuyBooks(c *gin.Context) {
	tx := model.Conn.Begin()
	userId, _ := strconv.ParseInt(c.PostForm("id"), 10, 64)
	bookId, _ := strconv.ParseInt(c.PostForm("uid"), 10, 64)
	num, _ := strconv.ParseInt(c.PostForm("num"), 10, 64)
	var book model.Book
	if err := tx.Table("book").Where("uid=?", bookId).Find(&book).Error; err != nil {
		tx.Rollback()
	}
	price := book.Price
	price = price * float64(num)
	tx.Commit()
	// 获取url进行支付
	client, err := alipay.NewClient(config.AppId, config.PrivateKey, config.IsProduction)
	if err != nil {
		log.Println("支付宝初始化错误")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "支付宝初始化错误"})
		return
	}
	client.SetCharset("utf-8").SetSignType(alipay.RSA2).SetNotifyUrl(config.NotifyURL).SetReturnUrl(config.ReturnURL)

	ts := time.Now().UnixMilli()
	outTradeNo := fmt.Sprintf("%d", ts)
	bm := make(gopay.BodyMap)
	bm.Set("subject", "这里是小陈的支付页面")
	bm.Set("out_trade_no", outTradeNo)
	bm.Set("total_amount", price)
	bm.Set("product_code", config.ProductCode)
	payUrl, err := client.TradePagePay(context.Background(), bm)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "支付链接生成失败"})
		return
	}

	// 更新购买逻辑，例如生成订单、更新库存等
	err = model.BuyBook(userId, bookId, num)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": err.Error()})

		return
	}

	// 返回购买成功的 JSON 响应，并返回支付链接
	c.JSON(http.StatusOK, gin.H{"message": "恭喜您，买书成功", "payUrl": payUrl})
}

// TOKEN 假设您在Go代码中定义了一个名为TOKEN的常量，用于存储您的令牌值
const TOKEN = "123"

// 配置公众号的token
func CheckSignature(c *gin.Context) {
	// 获取查询参数中的签名、时间戳和随机数
	signature := c.Query("signature")
	timestamp := c.Query("timestamp")
	nonce := c.Query("nonce")
	echostr := c.Query("echostr")
	// 创建包含令牌、时间戳和随机数的字符串切片
	tmpArr := []string{TOKEN, timestamp, nonce}
	// 对切片进行字典排序
	sort.Strings(tmpArr)
	// 将排序后的元素拼接成单个字符串
	tmpStr := ""
	for _, v := range tmpArr {
		tmpStr += v
	}
	// 对字符串进行SHA-1哈希计算
	tmpHash := sha1.New()
	tmpHash.Write([]byte(tmpStr))
	tmpStr = fmt.Sprintf("%x", tmpHash.Sum(nil))
	fmt.Println(tmpStr)
	fmt.Println(signature)
	// 将计算得到的签名与请求中提供的签名进行比较，并根据结果发送相应的响应
	if tmpStr == signature {
		c.String(200, echostr)
		model.Rdb.Set(context.Background(), "library:token", tmpStr, 7*24*time.Hour)
	} else {
		c.String(403, "签名验证失败 "+timestamp)
	}
}

// Redirect 微信扫码登录
// @Summary 用户登录接口3
// @Description 通过微信扫码登录，手机进行登录验证
// @Tags 公开
// @Accept json
// @Produce application/json
// @Param Url query string true "内网穿透地址"
// @Router /api/v1/wechat/login [get]
func Redirect(c *gin.Context) {
	//path := c.Query("Url")
	//防止跨站请求伪造攻击 增加安全性
	redirectURL := url.QueryEscape("http://yji2ai.natappfree.cc/user/wechat/vx") //userinfo,
	wechatLoginURL := fmt.Sprintf("https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&state=%s&scope=snsapi_userinfo#wechat_redirect", "wx068062b4c8b47f31", redirectURL, "state")
	wechatLoginURL, _ = url.QueryUnescape(wechatLoginURL)
	// 生成二维码
	qrCode, err := qrcode.Encode(wechatLoginURL, qrcode.Medium, 256)
	if err != nil {
		// 错误处理
		c.String(http.StatusInternalServerError, "Error generating QR code")
		return
	}
	// 将二维码图片作为响应返回给用户
	c.Header("Content-Type", "image/png")
	c.Writer.Write(qrCode)
}

func Callback(c *gin.Context) {
	// 获取微信返回的授权码
	code := c.Query("code")
	// 向微信服务器发送请求，获取access_token和openid
	tokenResp, err := http.Get(fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code", "wx068062b4c8b47f31", "appsecret", code))
	if err != nil {
		fmt.Println(err)
		resp := &tools.ECode{
			Data:    nil,
			Message: "error,获取token失败",
			Code:    10001,
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	// 解析响应中的access_token和openid
	var tokenData struct {
		AccessToken  string `json:"access_token"`
		ExpiresIn    int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
		OpenID       string `json:"openid"`
		Scope        string `json:"scope"`
	}
	if err1 := json.NewDecoder(tokenResp.Body).Decode(&tokenData); err1 != nil {
		resp := &tools.ECode{
			Data:    nil,
			Message: "error,获取token失败",
			Code:    10002,
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	userInfoURL := fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s", tokenData.AccessToken, tokenData.OpenID)
	userInfoResp, err := http.Get(userInfoURL)
	if err != nil {
		// 错误处理
		//zap.L().Error("获取失败")
		fmt.Println(err)
		return
	}
	defer userInfoResp.Body.Close()

	var userData struct {
		OpenID   string `json:"openid"`
		Nickname string `json:"nickname"`
	}
	if err1 := json.NewDecoder(userInfoResp.Body).Decode(&userData); err1 != nil {
		// 错误处理
		//zap.L().Error("获取用户信息失败")
		fmt.Println(err1)
		return
	}
	//用户的名字
	var user1 model.User
	nickname := userData.Nickname
	if err2 := model.Conn.Where("user_name=?", nickname).First(&user1).Error; err2 != nil {
		if errors.Is(err2, gorm.ErrRecordNotFound) {
			user1.UserName = nickname
			//user1.UserID, _ = snowflake.GetID()
			user1.UserID = tools.GetUid()
			user1.Identity = "普通用户"
		} else {
			//zap.L().Error("验证登录信息过程中出错")
			//ResponseError(c, CodeServerBusy)
			return
		}
	}
	////添加jwt验证
	//atoken, rtoken, err3 := GetToken(user1.UserID, user1.UserName, user1.Identity)
	//
	//if err3 != nil {
	//	//zap.L().Error("生成JWT令牌失败")
	//	return
	//}
	//c.Header("Authorization", atoken)
	////发送成功响应
	//ResponseSuccess(c, &LoginData{
	//	AccessToken:  atoken,
	//	RefreshToken: rtoken,
	//})
	////zap.L().Info("登录成功")
	//return
}
func GetToken(userID uint64, userName, identity string) (string, string, error) {
	// 创建jwt claims
	claims := jwt.MapClaims{
		"user_id":   userID,
		"user_name": userName,
		"identity":  identity,
		"exp":       time.Now().Add(time.Hour * 24).Unix(), // 过期时间为24小时
	}

	// 创建jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte("your_secret_key"))
	if err != nil {
		return "", "", err
	}

	// 创建刷新令牌
	refreshClaims := jwt.MapClaims{
		"user_id":   userID,
		"user_name": userName,
		"identity":  identity,
		"exp":       time.Now().Add(time.Hour * 720).Unix(), // 过期时间为30天
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte("your_refresh_secret_key"))
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshTokenString, nil
}
func ResponseSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": 10003,
		"msg":  "success",
		"data": data,
	})
}

type LoginData struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
