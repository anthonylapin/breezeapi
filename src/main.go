package main

import (
	"breezeapi/src/server"
	"flag"
	"fmt"
)

func parseCLIArgs() map[string]string {
	directory := flag.String("directory", "", "Path to the directory")
	flag.Parse()

	result := map[string]string{}

	if directory != nil {
		result["fileDirectory"] = *directory
	}

	return result
}

func main() {
	cliArgs := parseCLIArgs()
	fmt.Println(cliArgs)

	controller := RequestController{CLIArgs: cliArgs}

	httpServer := server.NewServer()

	router := server.NewRouter()

	router.Get("/", controller.ping)

	router.Get("/echo/{str}", controller.echo)

	router.Get("/user-agent", controller.userAgent)

	router.Get("/files/{fileName}", controller.getFile)

	router.Post("/files/{fileName}", controller.writeFile)

	httpServer.AddRouter(router)

	const PORT = 4221
	httpServer.Listen(PORT)
}
