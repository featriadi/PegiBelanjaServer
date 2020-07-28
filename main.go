package main

import (
	"pb-dev-be/db"
	"pb-dev-be/routes"
)

func main() {
	db.Init()

	e := routes.Init()

	// e.Logger.Fatal(e.Start(":4646"))
	e.Logger.Fatal(e.Start(":1234"))
}
