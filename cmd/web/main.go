package main

import (
	_ "Hackathon/docs"
	"Hackathon/internal/app"
)

//	@title			Hackathon API
//	@version		1.0
//	@description	API Server for Hackathon

//	@BasePath	/
func main() {
	app.Run()
}
