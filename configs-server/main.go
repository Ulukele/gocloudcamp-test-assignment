package main

import (
	"configs-server/data"
	"configs-server/server"
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
	log.Fatal(s.StartServer())
}
