package main

import (
	"fmt"
	"log"

	"github.com/google/uuid"
)

func main() {
	// create new uiid
	userID, err := uuid.NewUUID()
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(userID)
}
