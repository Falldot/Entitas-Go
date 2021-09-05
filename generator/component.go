package generator

import (
	"go/ast"
	"unicode"
)

type Component struct {
	Name   string
	Fields map[string]string
	Ident  bool
}

func FindComponents(f *ast.File) []*Component {
	var components []*Component
	ast.Inspect(f, func(n ast.Node) bool {
		comp := &Component{Fields: make(map[string]string)}
		switch t := n.(type) {
		case *ast.TypeSpec:
			switch t.Type.(type) {
			case *ast.StructType:
				if unicode.IsUpper(rune(t.Name.Name[0])) {
					comp.Name = t.Name.Name
				}
			case *ast.Ident:
				if unicode.IsUpper(rune(t.Name.Name[0])) {
					comp.Ident = true
					comp.Name = t.Name.Name
					comp.Fields[comp.Name] = t.Type.(*ast.Ident).String()
				}
			}
		}
		if comp.Name != "" {
			components = append(components, comp)
		}
		return true
	})
	return components
}

const componentConstansTemplate = `
const {context}ComponentTotal = {componentCount}

const (
	{const} = iota //next
)
`

const componentTemplate = `
func (e *Entity) Add{name}({argsWithType}) {
	var c *{name}Component
	if comp, ok := e.Create({name}); ok {
		c = comp.(*{name}Component)
	} else {
		c = new({name}Component)
	}
	{result}
	e.components[{name}] = c
	e.onComponentAdd.Execute(e, {name}, c)
}

func (e *Entity) Replace{name}({argsWithType}) {
	if e.components[{name}] != nil {
		c := e.components[{name}].(*{name}Component)
		{result}
		e.onComponentReplace.Execute(e, {name}, c)
	} else {
		e.Add{name}({args})
	}
}

func (e *Entity) On{name}() {
	e.onComponentAdd.Execute(e, {name}, e.components[{name}])
}

func (e *Entity) Off{name}() {
	e.onComponentOff.Execute(e, {name}, e.components[{name}])
}

func (e *Entity) Remove{name}() {
	e.onComponentRemove.Execute(e, {name}, e.components[{name}])
}
`

const componentTemplateGetMethodStruct = `
func (e *Entity) Get{name}() *{name}Component {
	return e.components[{name}].(*{name}Component)
}
`

const componentTemplateGetMethodSingleType = `
func (e *Entity) Get{name}() {type} {
	return e.components[{name}].({type})
}
`
