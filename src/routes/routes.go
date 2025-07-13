package routes

import (
	"github.com/cakezero/go-server/src/controllers"
	"github.com/gorilla/mux"
)

func Routes() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", controllers.Home)
	router.HandleFunc("/login", controllers.Login).Methods("POST")
	router.HandleFunc("/register", controllers.Register).Methods("POST")
	router.HandleFunc("/perform-action", controllers.PerformAction).Methods("POST")
	router.HandleFunc("/get-history", controllers.ArithmeticHistory).Methods("GET")

	return router
}
