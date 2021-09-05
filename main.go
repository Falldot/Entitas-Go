package main

import (
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"

	"github.com/Falldot/Entitas-Go/generator"
)

func main() {
	fset := token.NewFileSet()
	inFileName := os.Getenv("GOFILE")

	src, err := ioutil.ReadFile(inFileName)
	if err != nil {
		panic(err)
	}

	f, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	components := generator.FindComponents(f)

	generator.CreateEntitasLibFile()

	generator.InitContext(inFileName[:len(inFileName)-3])

	generator.CreateEntitasContextFile(inFileName, components, src)
}
