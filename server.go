package main

import (
	"fmt"
	"github.com/cakezero/go-server/src/routes"
	// "github.com/cakezero/go-server/src/utils"
	"net/http"
	"log"
)


func main() {
	routes.Routes()

	// utils.GetRedisClient()

	fmt.Println("Server is running on port: 9000")

	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
