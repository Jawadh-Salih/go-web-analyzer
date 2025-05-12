package main

import (
	"github.com/Jawadh-Salih/go-web-analyzer/internal/router"
)

func main() {
	engine := router.Init(":8080")

	err := engine.Run()
	if err != nil {
		// handle this gracefully later
		panic(err)
	}

}
