package model

import (
	"fmt"
	"github.com/spf13/viper"
	"math/rand"
	"net/mail"
	"net/smtp"
	"strconv"
	"time"
)

type EmailUtil struct {
	SMTPHost string
	SMTPPort int
	Email    string
	Password string
}

func NewEmailUtil() *EmailUtil {
	return &EmailUtil{}
}

//获取配置信息，如不用配置文件，这里直接赋值即可

func (eu *EmailUtil) LoadConfig() {
	viper.SetConfigFile("./app/view/config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Failed to read config file: %v", err)
	}
	eu.SMTPHost = viper.GetString("smtp.host")
	eu.SMTPPort = viper.GetInt("smtp.port")
	eu.Email = viper.GetString("smtp.email")
	eu.Password = viper.GetString("smtp.password")

}

func (eu *EmailUtil) SendVerificationCode(email string) (string, error) {
	code := GenerateVerificationCode()
	from := mail.Address{Name: "Notification", Address: eu.Email}

	// 2. 设置SMTP服务器相关信息
	auth := smtp.PlainAuth("", eu.Email, eu.Password, eu.SMTPHost)
	addr := fmt.Sprintf("%s:%d", eu.SMTPHost, eu.SMTPPort)

	// 3. 组装邮件内容
	header := make(map[string]string)
	header["From"] = from.String()
	header["To"] = email
	header["Subject"] = "Email verification code"
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/plain; charset=\"utf-8\""
	content := "你的验证码是: " + code
	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + content
	// 4. 发送邮件
	err := smtp.SendMail(addr, auth, eu.Email, []string{email}, []byte(message))

	if err != nil {
		return "", err
	}

	// 存储验证码到 Redis
	key := GenerateVerificationCodeKey(email)
	err = SetRedis(key, code)
	if err != nil {
		return "", err
	}
	return key, nil
}

// 验证

func (eu *EmailUtil) VerifyEmailCode(email, key, code string) bool {
	// 1. 从Redis获取存储的验证码
	storedCode, err := Getsession(key)
	if err != nil {
		return false
	}
	if email != key[18:] {
		return false
	}
	// 2. 比较验证码
	if storedCode != code {
		return false
	}
	// 3. 验证码匹配,可以删除Redis中的验证码
	_, err = DelRedis(key)
	if err != nil {
		fmt.Printf("Error deleting code: %v", err)
	}
	// 4. 返回成功
	return true
}

// 生成6位数字验证码

func GenerateVerificationCode() string {
	rand.Seed(time.Now().UnixNano())
	return strconv.Itoa(rand.Intn(999999))
}

// 设置redis-key

func GenerateVerificationCodeKey(email string) string {
	return fmt.Sprintf("verification_code:%s", email)
}
