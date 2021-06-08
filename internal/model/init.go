package model

import (
	"encoding/gob"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"net/http"
	"os"
)

var (
	// db 表示数据库对象
	db *gorm.DB
	// store 表示session store对象
	store *sessions.FilesystemStore
)

// Init 初始化数据库和session
func Init() {
	// 连接到数据库
	dbName := os.Getenv("DB_NAME")
	if dbName == "sqlite3" {
		sqlite3Conn(os.Getenv("DB_PATH"))
	}

	// 更新schema
	db.AutoMigrate(&User{}, &File{})

	// 设置session store
	SetupSession()

	// 创建upload和session目录
	dirs := []string{os.Getenv("UPLOAD_DIR"), os.Getenv("SESSION_DIR")}
	for _, dir := range dirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err = os.MkdirAll(dir, 0755)
			if err != nil {
				panic(err)
			}
		}
	}
}

// sqlite3Conn 连接到sqlite3数据库
func sqlite3Conn(path string) {
	var err error
	db, err = gorm.Open("sqlite3", path)
	if err != nil {
		panic(err)
	}
}

// SetupSession 初始化session store
func SetupSession() {
	store = sessions.NewFilesystemStore(os.Getenv("SESSION_DIR"), securecookie.GenerateRandomKey(16))
	store.Options = &sessions.Options{
		MaxAge:   24 * 60 * 60, // 1 day
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode, // Set cookies for cross origin requests
	}

	gob.Register(&User{})
}
