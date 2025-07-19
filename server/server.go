package main

import (
	"fmt"
	"github.com/cakezero/go-server/src/routes"
	"github.com/cakezero/go-server/src/utils"
	"net/http"
	"github.com/rs/cors"
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

	newCors := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "https://go-cal.vercel.app"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	corsHandler := newCors.Handler(router)

	port := utils.PORT

	fmt.Printf("Server is running on port %s\n", port)

	err := http.ListenAndServe(port, corsHandler)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
