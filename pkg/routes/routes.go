package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/remaster/webauthn/pkg/controllers"
)

var router *mux.Router

func setupRoutes() {
	router = mux.NewRouter()
	// router.Handle("/login/begin/{username}", http.HandlerFunc(controllers.BeginLogin)).Methods("GET")
	// router.Handle("/login/finish/{username}", http.HandlerFunc(controllers.FinishLogin)).Methods("POST")
	router.Handle("/register/begin/{username}", http.HandlerFunc(controllers.BeginRegistration)).Methods("GET")
	router.Handle("/register/finish/{username}", http.HandlerFunc(controllers.FinishRegistration)).Methods("POST")
}

func GetRouter() *mux.Router {
	setupRoutes()
	return router
}
