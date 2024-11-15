package tools

import (
	"github.com/mojocn/base64Captcha"
)

type CaptchaData struct {
	CaptchaId string `json:"captcha_id"`
	Data      string `json:"data"`
}

type driverString struct {
	Id            string
	CaptchaType   string
	VerifyValue   string
	DriverString  *base64Captcha.DriverString  //字符串
	DriverChinese *base64Captcha.DriverChinese //中文
	DriverMath    *base64Captcha.DriverMath    //数学
	DriverDigit   *base64Captcha.DriverDigit   //数字
}

// 数字驱动
var digitDriver = base64Captcha.DriverDigit{
	Height:   50,  //生成图片高度
	Width:    150, //生成图片宽度
	Length:   1,   //验证码长度
	MaxSkew:  0,   //文字的倾斜度 越大倾斜度越狠  越不容易看懂
	DotCount: 0,   //背景的点数，越大，字体越模糊
}

// 使用内存驱动，相关数据会存在内存空间里
var store = base64Captcha.DefaultMemStore

func CaptchaGenerate() (CaptchaData, error) {
	var ret CaptchaData
	//直接使用digitDriver会报错，必须传一个指针
	c := base64Captcha.NewCaptcha(&digitDriver, store)
	id, b64s, _, _ := c.Generate()

	ret.CaptchaId = id
	ret.Data = b64s
	return ret, nil
}
func CaptchaVerify(data CaptchaData) bool {
	return store.Verify(data.CaptchaId, data.Data, true)
}
