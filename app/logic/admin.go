package logic

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"tushuguanli/app/model"
	"tushuguanli/app/tools"
)

type CBooks struct {
	Uid         int64     `gorm:"column:uid;type:bigint(20);primary_key;AUTO_INCREMENT" json:"uid" form:"uid"`
	BookName    string    `gorm:"column:book_name;type:varchar(200);comment:书名" json:"book_name" form:"book_name"`
	Cate        string    `gorm:"column:cate;type:varchar(50);comment:图书类别" json:"cate" form:"cate"`
	Num         int       `gorm:"column:num;type:int(20);comment:书的数量" json:"num" form:"num"`
	Price       float64   `gorm:"column:price;type:double;default:1;comment:价格" json:"price" form:"price"`
	CreatedTime time.Time `gorm:"column:created_time;type:datetime" json:"created_time" form:"created_time"`
	UpdatedTime time.Time `gorm:"column:updated_time;type:datetime" json:"updated_time" form:"updated_time"`
}

func AddBooks(c *gin.Context) {
	var book CBooks
	if err := c.ShouldBind(&book); err != nil {
		c.JSON(200, tools.ParamErr)
		return
	}
	fmt.Println(book)
	if oldBook := model.GetBookV1(book.BookName); oldBook.Id > 0 {
		c.JSON(http.StatusOK, tools.ECode{
			Code:    10004,
			Message: "书籍已存在！",
		})
		return
	}
	newBook := model.BookInfo{
		Uid:         tools.GetUid(),
		BookName:    book.BookName,
		Cate:        book.Cate,
		Num:         book.Num,
		Price:       book.Price,
		CreatedTime: time.Now(),
		UpdatedTime: time.Now(),
	}

	err := model.CreatBooks(&newBook)
	if err != nil {
		c.JSON(200, tools.ECode{
			Message: "添加失败!",
		})
		return
	}
	c.JSON(200, tools.ECode{
		Message: "添加成功！",
	})
}

func GetBook(c *gin.Context) {
	idStr := c.Query("name")
	ret := model.GetBookV1(idStr)
	if ret.Id <= 0 {
		c.JSON(404, ret)
		return
	}
	c.JSON(200, gin.H{"data": ret})
}
