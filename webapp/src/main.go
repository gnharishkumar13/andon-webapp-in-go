package main

import (
	"html/template"
	"log"
	"os"
)

func main() {

	//Pipelines
	data := 42                            //Pipeline can be struct or function too
	content := `this is a template {{.}}` //Actions{{}} and . means pipelines that have the entire data
	t := template.New("")
	_, err := t.Parse(content)
	if err != nil {
		log.Fatal(err)
	}
	err = t.Execute(os.Stdout, data)
	if err != nil {
		log.Fatal(err)
	}
}
