package api

import (
	"cloudpan/internal/model"
	"cloudpan/internal/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
)

// Download 处理文件下载
func Download(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid file id!"})
		return
	}

	session, _ := model.GetSession(c.Request)
	user := session.Values["user"].(*model.User)

	// 获取文件信息
	file := model.GetFileByID(uint(id))
	if file.OwnerUID != user.ID {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Unauthorized access!"})
		return
	}

	// 从磁盘中读取加密文件
	encrypted, err := readFile(file.Filename, file.Size)
	if err != nil {
		util.Log().Warning("read file error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "read file error!"})
		return
	}
	util.Log().Info(fmt.Sprintf("Before decryption: size %d, md5sum %s", len(encrypted), util.MD5(encrypted)))

	// 对文件内容进行解密
	content, err := util.DecryptFile(encrypted, user.Key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "decrypt file error!"})
	}
	util.Log().Info(fmt.Sprintf("After decryption: size %d, md5sum %s", len(content), util.MD5(content)))

	// 返回解密结果
	c.Writer.WriteHeader(http.StatusOK)
	c.Header("Content-Disposition", "attachment; filename="+file.Filename)
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Accept-Length", fmt.Sprintf("%d", file.Size))
	c.Writer.Write(content)
}

func readFile(filename string, size int64) ([]byte, error) {
	content := make([]byte, size+16)

	src, err := os.Open(os.Getenv("UPLOAD_DIR") + "/" + filename)
	defer src.Close()
	if err != nil {
		return nil, err
	}

	_, err = src.Read(content)
	if err != nil {
		return nil, err
	}
	return content, nil
}
