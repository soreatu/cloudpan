package model

import (
	"github.com/gorilla/sessions"
	"net/http"
	"os"
)

func GetSession(r *http.Request) (*sessions.Session, error) {
	return store.Get(r, os.Getenv("SESSION_ID"))
}
