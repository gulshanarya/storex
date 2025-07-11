package main

import (
	"log"
	"storex/db"
	"storex/routes"
)

func main() {
	err := db.InitDB()
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer func() {
		err = db.DB.Close()
		if err != nil {
			log.Fatalf("DB close failed: %v", err)
		}
	}()

	//routes setup here
	routes.Setup()

}
