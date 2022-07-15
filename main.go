package main

import (
	"RMS/database"
	"RMS/server"
	"log"
)

func main() {

	dbConnection := database.ConnectAndMigrate("localhost", "5433", "rms", "local", "local", "disable")
	if dbConnection != nil {
		log.Printf("Main : Error in database connection")
		panic(dbConnection)
	}

	log.Printf("Connected")

	srv := server.SetUpRoutes()
	start := srv.Run(":8080")
	if start != nil {
		log.Printf("main : Error in listenig to the requests.")
		log.Fatal(start)
	}

}
