package routes

import (
	"net/http"

	"github.com/cakezero/go-server/src/controllers"
	"github.com/cakezero/go-server/src/middlewares"
	"github.com/gorilla/mux"
)

func Routes() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", controllers.Home)
	router.HandleFunc("/login", controllers.Login).Methods("POST")
	router.HandleFunc("/register", controllers.Register).Methods("POST")
	router.Handle("/perform-action", middlewares.AuthMiddleware(http.HandlerFunc(controllers.PerformAction))).Methods("POST")
	router.Handle("/get-history", middlewares.AuthMiddleware(http.HandlerFunc(controllers.ArithmeticHistory))).Methods("GET")

	return router
}
