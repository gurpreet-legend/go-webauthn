package config

// import (
// 	"time"

// 	"github.com/gofiber/fiber/v2/middleware/session"
// 	"github.com/gofiber/storage/mysql"
// )

// var store *session.Store

// func CreateSession() {
// 	storage := mysql.New(mysql.Config{
// 		ConnectionURI: "root:lostworld2701@tcp(127.0.0.1:3306)/webauthn-users",
// 		Reset:         false,
// 		GCInterval:    10 * time.Second,
// 	})
// 	// storage := sqlite3.New()
// 	store = session.New(session.Config{
// 		Storage: storage,
// 	})
// }

// func GetSession() *session.Store {
// 	return store
// }
