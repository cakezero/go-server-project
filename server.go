package main

import (
	"fmt"
	"github.com/cakezero/go-server/src/routes"
	"github.com/cakezero/go-server/src/utils"
	"net/http"
	"log"
)

func init () {
	loadEnvErr := utils.LoadEnv();
	if loadEnvErr != nil {
		panic(loadEnvErr)
	}

	dbConnectErr := utils.DB()
	if dbConnectErr != nil {
		panic(dbConnectErr)
	}

	fmt.Println("DB connected")
}


func main() {
	router := routes.Routes()

	port := utils.PORT

	fmt.Printf("Server is running on port %s\n", port)

	err := http.ListenAndServe(port, router)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
