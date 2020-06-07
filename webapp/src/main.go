package main

import (
	"html/template"
	"log"
	"os"
)

func main() {

	content := `{{- /** This is a comment */ -}} this is a template`
	t := template.New("my first template")
	t, err := t.Parse(content)
	if err != nil {
		log.Fatal(err)
	}
	err = t.Execute(os.Stdout, nil)
	if err != nil {
		log.Fatal(err)
	}
}
