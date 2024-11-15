package model

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

//数据库操作

var Conn *gorm.DB

var Rdb *redis.Client
var Mdb *mongo.Client

func NewMysql() {
	my := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", "root", "123456", "127.0.0.1:3306", "library")
	conn, err := gorm.Open(mysql.Open(my), &gorm.Config{})

	if err != nil {
		fmt.Printf("err:%s\n", err)
		panic(err)
	}
	Conn = conn
}

func NewRdb() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "192.168.189.11:6379",
		Password: "", //未使用密码
		DB:       0,  //数据库编号  使用默认数据库编号0
	})

	//通过调用 redisstore.NewRedisStore 方法初始化会话存储。
	//该方法接收一个上下文对象和一个 Redis 客户端实例作为参数，返回一个会话存储对象和一个错误对象。
	//在这里，会话存储对象赋值给全局变量 store

	Rdb = rdb
	//初始化 session
	//store, _ = redisstore.NewRedisStore(context.TODO(), Rdb)
	return
}

func NewMongoDB() {
	// 设置MongoDB连接选项
	clientOptions := options.Client().ApplyURI("mongodb://192.168.189.11:27017")

	// 连接到MongoDB
	mdb, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	// 检查连接
	err = mdb.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB!")

	// 将MongoDB客户端赋值给全局变量
	Mdb = mdb
}

// 关闭数据库连接和redis连接

func Close() {
	db, _ := Conn.DB()
	_ = db.Close()
	_ = Rdb.Close()

	//关闭 MongoDB 数据库
	if Mdb != nil {
		err := Mdb.Disconnect(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		log.Println("MongoDB connection closed.")
	}
}
