package logic

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
	"tushuguanli/app/model"
	"tushuguanli/app/tools"
)

func GetBooks(c *gin.Context) {
	var id int64
	idStr := c.Query("id")
	id, _ = strconv.ParseInt(idStr, 10, 64)
	ret := model.GetBook(id)
	if ret.Id <= 0 {
		c.JSON(404, ret)
		return
	}

	c.JSON(200, gin.H{"data": ret})
	fmt.Printf("ret%s", ret)
}

type CBook struct {
	Name  string  `gorm:"column:name;default:NULL" json:"name"form:"name"`
	Cate  string  `gorm:"column:cate;default:NULL" json:"cate"form:"cate"`
	Num   int64   `gorm:"column:num;default:NULL" json:"num"form:"num"`
	Price float64 `gorm:"column:price;type:double;default:1;comment:价格" json:"price" form:"price"`
}

func AddBook(c *gin.Context) {
	var book CBook
	if err := c.ShouldBind(&book); err != nil {
		c.JSON(200, tools.ParamErr)
		return
	}
	if oldBook := model.GetUser(book.Name); oldBook.Id > 0 {
		c.JSON(http.StatusOK, tools.ECode{
			Code:    10004,
			Message: "书籍已存在！",
		})
		return
	}
	newBook := model.Book{
		Uid:         tools.GetUid(),
		Name:        book.Name,
		Cate:        book.Cate,
		Status:      0,
		Num:         book.Num,
		Price:       book.Price,
		CreatedTime: time.Now(),
		UpdatedTime: time.Now(),
	}
	err := model.CreatBook(&newBook)
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

func DelBook(c *gin.Context) {
	var id int64
	idStr := c.Query("id")
	id, _ = strconv.ParseInt(idStr, 10, 64)
	book := model.GetBook(id)
	if book.Id <= 0 {
		c.JSON(200, tools.OK)
		return
	}
	if err := model.DelBook(id); err != nil {
		c.JSON(http.StatusOK, tools.ECode{
			Code:    10006,
			Message: "删除失败",
		})
		return
	}
	c.JSON(200, tools.ECode{
		Message: "删除成功",
	})
}

func UpdateBook(c *gin.Context) {

}
