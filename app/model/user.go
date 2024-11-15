package model

import "fmt"

func GetUser(name string) *User {
	var ret User
	if err := Conn.Table("user").Where("name = ?", name).Find(&ret).Error; err != nil {
		fmt.Printf("err:%s", err.Error())
	}
	return &ret
}

func CreatUser(user *User) error {
	if err := Conn.Create(user).Error; err != nil {
		fmt.Printf("err%s", err.Error())
		return err
	}
	return nil
}
