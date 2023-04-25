package config

import (
	"github.com/gofiber/fiber/v2/middleware/session"
)

var store *session.Store

func CreateSession() {
	store = session.New()
}

func GetSession() *session.Store {
	return store
}
