package main

import (
	"pb-dev-be/db"
	"pb-dev-be/routes"
)

func main() {
	db.Init()

	e := routes.Init()

	e.Logger.Fatal(e.Start(":1234"))
	// e.Logger.Fatal(e.StartAutoTLS(":1234"))
	// fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
}
