package main

import (
	"html/template"

	"github.com/Jawadh-Salih/go-web-analyzer/internal/router"
)

var tmpl = template.Must(template.ParseFiles(
	"web/index.html",
))

func main() {
	engine := router.Init(":8080")

	err := engine.Run()
	if err != nil {
		// handle this gracefully later
		panic(err)
	}

}
