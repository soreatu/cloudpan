package api

import (
	model "cloudpan/internal/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
)

// GetFiles 获取当前用户的所有文件
func GetFiles(c *gin.Context) {
	session, _ := model.GetSession(c.Request)
	user := session.Values["user"].(*model.User)

	files := model.GetFilesByUID(user.ID)
	c.JSON(http.StatusOK, gin.H{"data": files})
}

// DeleteFile 删除指定id的文件
func DeleteFile(c *gin.Context) {
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

	// 删除磁盘中的文件
	if err = os.Remove(os.Getenv("UPLOAD_DIR") + "/" + file.Filename); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Delete file error!"})
		return
	}
	// 删除文件信息记录
	model.DeleteFile(uint(id))

	c.JSON(http.StatusOK, gin.H{"msg": "Delete success!"})
}
