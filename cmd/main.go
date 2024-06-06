package main

import (
	"fmt"
	"gochat/cmd/api"
)

func main() {
	server := api.InitServer(":3001")

	if err := server.InitAPI(); err != nil {
		fmt.Println(err)
		fmt.Println("Error in init api")
		return
	}
}
