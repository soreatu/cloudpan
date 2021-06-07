package model

import (
	"errors"
	"time"
)

type File struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	Filename  string     `gorm:"unique" json:"filename"`
	Size      int64      `json:"size"`
	OwnerUID  uint       `json:"uid"`
}

// NewFile 返回一个空的File对象
func NewFile() File {
	return File{}
}

// CreateFile 在数据库中创建一个新的User记录
func CreateFile(f *File) (err error) {
	db = db.Create(f)
	if db.Error != nil {
		return db.Error
	}
	if db.RowsAffected != 1 {
		return errors.New("affected rows != 1")
	}
	return nil
}

// GetFilesByUID 获取指定用户id的所有文件信息
func GetFilesByUID(uid uint) (files []File) {
	db.Model(&File{}).Where(&File{OwnerUID: uid}).Find(&files)
	return
}

// GetFileByID 获取指定id的文件信息
func GetFileByID(id uint) (file File) {
	db.First(&file, id)
	return
}

// DeleteFile 删除指定id的文件信息
func DeleteFile(id uint) {
	db.Delete(&File{}, id)
}
