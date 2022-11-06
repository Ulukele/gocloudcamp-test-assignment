package main

import (
	"configs-server/data"
	_ "configs-server/docs"
	"configs-server/server"
	"fmt"
	"log"
	"os"
)

func main() {

	// Configure mongo connection
	mongoConfig := data.MongoConnectionConfig{}
	mongoConfig.User = os.Getenv("MONGO_USER")
	mongoConfig.Password = os.Getenv("MONGO_PASSWORD")
	mongoConfig.Host = os.Getenv("MONGO_HOST")

	// Create data manager
	manager, err := data.NewManager(mongoConfig)
	if err != nil {
		panic("Unable to create data manager")
	}

	// Create server
	s := server.NewServer(manager)

	// Start server
	serverPort := os.Getenv("PORT")
	log.Fatal(s.StartServer(fmt.Sprintf(":%s", serverPort)))
}
