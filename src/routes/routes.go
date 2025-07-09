package routes

import (
	"github.com/cakezero/go-server/src/controllers"
	"net/http"
)

func Routes() {
	http.HandleFunc("/", controllers.Home)
	http.HandleFunc("/login", controllers.Login)
	http.HandleFunc("/register", controllers.Register)
}