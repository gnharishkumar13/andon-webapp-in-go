package main

import (
	"fmt"
	"html/template"
	"log"
	"os"
)

type data struct {
	Message string
	Answer  int
}

func (d data) PrintMessage(message string, answer int) string {
	return fmt.Sprintf("%v : %v", message, answer)
}

func main() {

	//Pipelines                          //Pipeline can be struct or function too
	content := `{{.PrintMessage .Message .Answer}}` //Actions{{}} and . means pipelines that have the entire data
	t := template.New("")
	_, err := t.Parse(content)
	if err != nil {
		log.Fatal(err)
	}
	err = t.Execute(os.Stdout, data{
		Message: "The answer is",
		Answer:  42,
	})
	if err != nil {
		log.Fatal(err)
	}
}
