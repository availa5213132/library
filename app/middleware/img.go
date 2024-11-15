package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"time"
	"tushuguanli/app/model"
)

func UploadHand(c *gin.Context) {
	c.HTML(http.StatusOK, "img.tmpl", nil)
}
func UploadHandler(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, "文件上传错误！")
		return
	}
	files := form.File["image"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, "上传失败！")
		return
	}
	saveDir := "img"
	err = os.MkdirAll(saveDir, os.ModePerm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "创建存储文件夹出错！")
		return
	}
	file := files[0]
	err = c.SaveUploadedFile(file, saveDir+"/"+file.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "文件保存错误！")
		return
	}
	c.JSON(http.StatusOK, "文件上传成功！")
}

type Image struct {
	ID       string    `bson:"_id,omitempty"`
	Filename string    `bson:"filename"`
	FilePath string    `bson:"filepath"`
	UploadAt time.Time `bson:"upload_at"`
}

func UploadHandlerV1(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.String(http.StatusBadRequest, "Error retrieving the file")
		return
	}

	files := form.File["image"]
	if len(files) == 0 {
		c.String(http.StatusBadRequest, "No file uploaded")
		return
	}

	file := files[0]

	// 获取对应的集合
	collection := model.Mdb.Database("Image").Collection("images")

	// 构造要插入的文档
	image := Image{
		Filename: file.Filename,
		FilePath: "./img/" + file.Filename, // 这里将文件名作为图片地址存储
		UploadAt: time.Now(),
	}
	// 构建目标文件路径
	targetPath := "./img/" + file.Filename

	// 将文件保存到目标路径
	err = c.SaveUploadedFile(file, targetPath)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error saving the file")
		return
	}
	// 插入文档到 MongoDB
	_, err = collection.InsertOne(context.Background(), image)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error saving the file information")
		return
	}

	c.String(http.StatusOK, "File uploaded successfully")
}
