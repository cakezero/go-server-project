package controllers

import (
	"fmt"
	"net/http"

	// "encoding/json"
)

func Login(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Login hit!")
}

func Register(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Register hit!")
}

func Home(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Home hit!")
}
