package main

import (
	"html/template"
	"log"
	"os"
)

func main() {

	childContent := `this is a child template`
	parentContent := `
	this is a parent template first line
	{{ template "child_template" }}
	this is a parent template second line
	`
	parent := template.New("parent_template")
	child := parent.New("child_template") //Template composition

	_, err := parent.Parse(parentContent)

	if err != nil {
		log.Fatal(err)
	}

	_, err = child.Parse(childContent)

	if err != nil {
		log.Fatal(err)
	}

	err = parent.Execute(os.Stdout, nil)

}
