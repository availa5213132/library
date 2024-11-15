package model

import "time"

type Admin struct {
	Id          int64     `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT" json:"id"`
	Name        string    `gorm:"column:name;type:varchar(50)" json:"name"`
	Password    string    `gorm:"column:password;type:varchar(50)" json:"password"`
	CreatedTime time.Time `gorm:"column:created_time;type:datetime" json:"created_time"`
	UpdatedTime time.Time `gorm:"column:updated_time;type:datetime" json:"updated_time"`
}

func (m *Admin) TableName() string {
	return "admin"
}

type Book struct {
	Id          int64     `gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL" json:"id" form:"id"`
	Uid         int64     `gorm:"column:uid;default:NULL" json:"uid"form:"uid"`
	Name        string    `gorm:"column:name;default:NULL" json:"name"form:"name"`
	Cate        string    `gorm:"column:cate;default:NULL" json:"cate"form:"cate"`
	Status      int64     `gorm:"column:status;default:NULL" json:"status"form:"status"`
	Num         int64     `gorm:"column:num;default:NULL" json:"num"form:"num"`
	Price       float64   `gorm:"column:price;type:bigint(20)" json:"price"form:"price"`
	CreatedTime time.Time `gorm:"column:created_time;default:NULL" json:"created_time"form:"created_time"`
	UpdatedTime time.Time `gorm:"column:updated_time;default:NULL" json:"updated_time"form:"updated_time"`
	Timestmpe   time.Time `gorm:"column:timestmpe;default:NULL"`
}

func (b *Book) TableName() string {
	return "book"
}

type BookInfo struct {
	Id                 uint      `gorm:"column:id;type:int(11) unsigned;comment:书的id" json:"id" form:"id"`
	Uid                int64     `gorm:"column:uid;type:bigint(20);primary_key;AUTO_INCREMENT" json:"uid" form:"uid"`
	BookName           string    `gorm:"column:book_name;type:varchar(200);comment:书名" json:"book_name" form:"book_name"`
	Author             string    `gorm:"column:author;type:varchar(50);comment:作者" json:"author" form:"author"`
	PublishingHouse    string    `gorm:"column:publishing_house;type:varchar(50);comment:出版社" json:"publishing_house" form:"publishing_house"`
	Translator         string    `gorm:"column:translator;type:varchar(50);comment:译者" json:"translator" form:"translator"`
	PublishDate        time.Time `gorm:"column:publish_date;type:date;comment:出版时间" json:"publish_date" form:"publish_date"`
	Pages              int       `gorm:"column:pages;type:int(10);default:100;comment:页数" json:"pages" form:"pages"`
	Num                int       `gorm:"column:num;type:int(20);comment:书的数量" json:"num" form:"num"`
	ISBN               string    `gorm:"column:ISBN;type:varchar(20);comment:ISBN号码" json:"ISBN" form:"ISBN"`
	Price              float64   `gorm:"column:price;type:double;default:1;comment:价格" json:"price" form:"price"`
	BriefIntroduction  string    `gorm:"column:brief_introduction;type:varchar(15000);comment:内容简介" json:"brief_introduction" form:"brief_introduction"`
	AuthorIntroduction string    `gorm:"column:author_introduction;type:varchar(5000);comment:作者简介" json:"author_introduction" form:"author_introduction"`
	ImgUrl             string    `gorm:"column:img_url;type:varchar(200);comment:封面地址" json:"img_url" form:"img_url"`
	DelFlg             int       `gorm:"column:del_flg;type:int(1);default:0;comment:删除标识" json:"del_flg" form:"del_flg"`
	Cate               string    `gorm:"column:cate;type:varchar(50);comment:图书类别" json:"cate" form:"cate"`
	CreatedTime        time.Time `gorm:"column:created_time;type:datetime" json:"created_time" form:"created_time"`
	UpdatedTime        time.Time `gorm:"column:updated_time;type:datetime" json:"updated_time" form:"updated_time"`
}

func (m *BookInfo) TableName() string {
	return "book_info"
}

type User struct {
	Id          int64     `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT" json:"id" form:"id"`
	Uid         int64     `gorm:"column:uid;type:bigint(20)" json:"uid" form:"uid"`
	Name        string    `gorm:"column:name;type:varchar(50)" json:"name" form:"name"`
	Password    string    `gorm:"column:password;type:varchar(50)" json:"password" form:"password"`
	Phone       string    `gorm:"column:phone;type:varchar(20)" json:"phone"`
	CreatedTime time.Time `gorm:"column:created_time;type:datetime" json:"created_time" form:"created_time"`
	UpdatedTime time.Time `gorm:"column:updated_time;type:datetime" json:"updated_time" form:"updated_time"`
	UserName    string
	UserID      interface{}
	Identity    string
}

func (m *User) TableName() string {
	return "user"
}

type BookUser struct {
	Id          int64     `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT" json:"id" form:"id"`
	UserId      int64     `gorm:"column:user_id;type:bigint(20)" json:"user_id"form:"user_id"`
	BookId      int64     `gorm:"column:book_id;type:bigint(20)" json:"book_id"form:"book_id"`
	Status      int64     `gorm:"column:status;type:bigint(20)" json:"status"form:"status"`
	BuyNumber   int64     `gorm:"column:buy_number;type:int(11)" json:"buy_number"form:"buy_number"`
	Price       int64     `gorm:"column:price;type:bigint(20)" json:"price"form:"price"`
	Time        int64     `gorm:"column:time;type:bigint(20)" json:"time"form:"time"`
	CreatedTime time.Time `gorm:"column:created_time;type:datetime" json:"created_time"form:"created_time"`
	UpdatedTime time.Time `gorm:"column:updated_time;type:datetime" json:"updated_time"form:"updated_time"`
}

func (m *BookUser) TableName() string {
	return "book_user"
}
