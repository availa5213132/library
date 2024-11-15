package model

import "fmt"

func CreatBooks(bookinfo *BookInfo) error {
	if err := Conn.Create(bookinfo).Error; err != nil {
		fmt.Printf("err%s", err.Error())
		return err
	}
	return nil
}
