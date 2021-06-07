package model

import (
	util "cloudpan/internal/util"
	"errors"
	"github.com/jinzhu/gorm"
)

// User 定义一个用户
type User struct {
	gorm.Model

	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Key      []byte
}

// NewUser 返回一个空的User对象
func NewUser() User {
	return User{}
}

// CreateUser 在数据库中创建一个新的User记录
func CreateUser(u *User) error {
	db = db.Create(u)
	if db.Error != nil {
		return db.Error
	}
	if db.RowsAffected != 1 {
		return errors.New("affected rows != 1")
	}
	return nil
}

// GetUserByName 通过用户名获取User对象
func GetUserByName(username string) User {
	var user User

	db.Model(&User{}).Where(&User{Username: username}).First(&user)

	return user
}

// GetUserKey 从数据库中获取指定id用户的加密密钥
func GetUserKey(uid int) []byte {
	var user User

	db.First(&user, uid)

	return user.Key
}

// BeforeCreate 将用户的明文密钥进行加密存储
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	//util.Log().Info(fmt.Sprintf("Plain key: %v", u.Key))
	return util.EncryptUserKey(u.Key)
	//util.Log().Info(fmt.Sprintf("Cipher key: %v", u.Key))
}

// AfterFind 将取出来的被加密用户密钥进行解密
func (u *User) AfterFind(tx *gorm.DB) (err error) {
	return util.DecryptUserKey(u.Key)
}
