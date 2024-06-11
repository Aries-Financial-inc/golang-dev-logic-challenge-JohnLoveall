package main

import (
	"github.com/Aries-Financial-inc/golang-dev-logic-challenge-JohnLoveall/routes"
)

func main() {
	router := routes.SetupRouter()
	router.Run(":8080")
}
