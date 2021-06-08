package api

import (
	"cloudpan/internal/model"
	"cloudpan/internal/util"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
)

// Upload 处理上传文件请求
func Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		util.Log().Warning("upload file error", err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": "upload file error!"})
		return
	}

	// 读取用户信息
	session, _ := model.GetSession(c.Request)
	user := session.Values["user"].(*model.User)

	// 读取上传文件的内容
	content := make([]byte, file.Size)
	src, err := file.Open()
	if err != nil {
		util.Log().Warning("open uploaded file error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "open uploaded file error!"})
		return
	}
	_, err = io.ReadFull(src, content)
	if err != nil {
		util.Log().Warning("read uploaded file error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "read uploaded file error!"})
		return
	}
	_ = src.Close()
	util.Log().Info(fmt.Sprintf("Before encryption: size %d, md5sum %s", len(content), util.MD5(content)))

	// 对文件内容进行加密
	encrypted, err := util.EncryptFile(content, user.Key)
	if err != nil {
		util.Log().Warning("encrypt file error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "encrypt file error!"})
		return
	}
	util.Log().Info(fmt.Sprintf("After encryption: size %d, md5sum %s", len(encrypted), util.MD5(encrypted)))

	// 将文件加密结果保存到磁盘文件中
	name := strings.Split(path.Base(file.Filename), ".")[0]
	filename := name + "_" + util.MD5(content) + path.Ext(file.Filename)
	if err = saveFile(encrypted, filename); err != nil {
		util.Log().Warning("save file error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "save file error!"})
		return
	}

	// 创建文件信息记录
	record := model.NewFile()
	record.Filename = filename
	record.Size = file.Size
	record.OwnerUID = user.ID
	if err = model.CreateFile(&record); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "create file record error!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": fmt.Sprintf("uploaded file has been saved as %s", filename)})
}

func saveFile(content []byte, filename string) error {
	dst, err := os.Create(os.Getenv("UPLOAD_DIR") + "/" + filename)
	if err != nil {
		return err
	}
	_, err = dst.Write(content)
	return err
}
