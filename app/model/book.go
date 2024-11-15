package model

import (
	"errors"
	"fmt"
	"github.com/goccy/go-json"
	"golang.org/x/net/context"
	"gorm.io/gorm"
	"time"
)

func GetAdmin(name string) *Admin {
	var ret Admin
	if err := Conn.Table("admin").Where("name = ?", name).Find(&ret).Error; err != nil {
		fmt.Printf("err:%s", err.Error())
	}
	return &ret
}

func GetBook(id int64) *Book {
	var ret Book
	if err := Conn.Table("book").Where("id = ?", id).Find(&ret).Error; err != nil {
		fmt.Printf("err%s", err.Error())
	}
	return &ret
}
func GetBookV1(name string) *BookInfo {
	var ret BookInfo
	if err := Conn.Table("book_info").Where("book_name = ?", name).Find(&ret).Error; err != nil {
		fmt.Printf("err%s", err.Error())
	}
	return &ret
}
func GetBooks(pageSize int, pageNum int) []BookInfo {
	ret := make([]BookInfo, 0)
	if err := Conn.Table("book_info").Where("uid>?", pageNum).Limit(pageSize).Find(&ret).Error; err != nil {
		fmt.Printf("err%s", err.Error())
	}
	storeBooksInRedis(ret)
	return ret
}

// 把分页数据缓存到redis中
func storeBooksInRedis(books []BookInfo) {
	// 序列化 books 切片为 JSON 字符串
	data, err := json.Marshal(books)
	if err != nil {
		fmt.Println("序列化书籍失败!", err)
		return
	}
	// 设置 Redis 键的过期时间
	expiration := 24 * time.Hour
	// 将 books 存储到 Redis 中
	err = Rdb.Set(context.TODO(), "books", data, expiration).Err()
	if err != nil {
		fmt.Println("缓存失败!", err)
		return
	}
	fmt.Println("缓存成功!")
}

//对 book 表进行增删改查

func CreatBook(book *Book) error {
	if err := Conn.Create(book).Error; err != nil {
		fmt.Printf("err%s", err.Error())
		return err
	}
	return nil
}

func DelBook(id int64) error {
	if err := Conn.Table("book_info").Where("id = ?", id).Delete(&Book{}).Error; err != nil {
		fmt.Printf("err%s", err.Error())
		return err
	}
	return nil
}

func UpdateBook(book *Book) error {
	if err := Conn.Save(book).Error; err != nil {
		fmt.Printf("err%s", err.Error())
		return err
	}
	return nil
}

func BorrowBook(uid, id int64) error {
	tx := Conn.Begin()
	//查询用户是否存在并加悲观锁
	var user User
	tx.Where("uid = ?", uid).Set("gorm:query_option", "FOR UPDATE").Find(&user)
	if user.Id == 0 {
		tx.Rollback()
		return errors.New("用户信息不存在")
	}

	//查询图书是否存在，是否正常并加悲观锁
	var book Book
	//使用乐观锁 记录时间戳
	//tx.Where("id = ?", id).Find(&book)
	tx.Where("id = ?", id).Set("gorm:query_option", "FOR UPDATE").Find(&book)
	if book.Id == 0 || book.Num <= 0 {
		tx.Rollback()
		return errors.New("图书信息不存在或库存不足")
	}

	// 记录当前时间戳
	//currentTimestamp := book.Timestamp

	//创建借阅记录
	now := time.Now()
	bu := BookUser{
		UserId:      uid,
		BookId:      id,
		Status:      1,
		Time:        1,
		CreatedTime: now,
		UpdatedTime: now,
	}
	if tx.Create(&bu).Error != nil {
		tx.Rollback()
		return errors.New("创建一个借阅记录")
	}
	//扣减图书库存，同时更新时间戳
	book.Num = book.Num - 1
	//book.Timestamp = now  //更新时间戳
	if tx.Save(&book).Error != nil /*|| book.Timestamp != currentTimestamp */ {
		tx.Rollback()
		return errors.New("扣减图书库存")
	}

	tx.Commit()
	return nil
	//我们将 currentTimestamp 变量用于记录当前的时间戳，并在保存图书记录后与数据库中的时间戳进行比较。
	//如果两者不匹配，说明在此期间有其他事务修改了该记录，我们选择进行回滚
}

func ReturnBook(uid, id int64) error {
	tx := Conn.Begin()
	//查询用户是否存在
	var user User
	tx.Where("uid = ?", uid).First(&user)
	if user.Id == 0 {
		tx.Rollback()
		return errors.New("用户信息不存在")
	}

	//查询图书是否存在，是否正常
	var book Book
	tx.Where("id = ?", id).First(&book)
	if book.Id == 0 || book.Num <= 0 {
		tx.Rollback()
		return errors.New("图书信息不存在或库存不足")
	}

	//查询借书记录是否存在
	var bu BookUser
	tx.Where("user_id = ? and book_id = ?", uid, id).First(&bu)
	if bu.Id <= 0 {
		tx.Rollback()
		return errors.New("借阅记录不存在")
	}

	//更新借阅状态
	bu.Status = 1
	if err := tx.Save(&bu).Error; err != nil {
		tx.Rollback()
		return errors.New(fmt.Sprintf("修改借阅记录失败：%s", err.Error()))
	}

	//更新图书库存
	book.Num = book.Num + 1
	if err := tx.Save(&book).Error; err != nil {
		tx.Rollback()
		return errors.New(fmt.Sprintf("增加库存失败：%s", err.Error()))
	}
	tx.Commit()
	return nil
}

func BuyBook(userId, bookId, buyNum int64) error {
	tx := Conn.Begin()
	var user User
	if err := tx.Table("user").Where("id=?", userId).Find(&user).Error; err != nil {
		tx.Rollback()
		return errors.New("用户不存在")
	}

	var book Book
	if err := tx.Table("book").Where("uid=?", bookId).Find(&book).Error; err != nil {
		tx.Rollback()
		return errors.New("您查询的书本不存在")
	}
	if book.Num < 1 || book.Num-buyNum < 0 {
		tx.Rollback()
		return errors.New("库存不够了")
	}
	//更新book表中的库存
	if err := tx.Model(&Book{}).Where("uid=?", bookId).UpdateColumn("num", gorm.Expr("num - ?", buyNum)).Error; err != nil {
		tx.Rollback()
		return errors.New("库存更新失败")
	}
	//创建买书记录
	bookNum := BookUser{
		UserId:      userId,
		BookId:      bookId,
		BuyNumber:   buyNum,
		CreatedTime: time.Now(),
		UpdatedTime: time.Now(),
	}
	if err := tx.Table("book_user").Create(&bookNum).Error; err != nil {
		tx.Rollback()
		return errors.New("添加购买记录失败")
	}
	tx.Commit()
	return nil

}
