package generator

import (
	"io/ioutil"
	"os"
	"strings"
)

const contextsTemplate = `
const (
	{First} = iota
	//NEXT1
)

func CreateContexts() Contexts {
	return SetContexts(
		//NEXT2
	)
}`

const getContextTemplate = `
func (c *Contexts) {name}() EntityBase {
	return (*c)[{name}].(EntityBase)
}`

func InitContext(context string) {
	if _, err := os.Stat("./Entitas/Contexts.go"); os.IsNotExist(err) {
		file, _ := os.Create("./Entitas/Contexts.go")
		defer file.Close()

		getter := strings.Replace(getContextTemplate, "{name}", context, -1)
		body := strings.Replace(contextsTemplate, "//NEXT2", "CreateEntityBase("+context+"ComponentTotal),\n//NEXT2", -1)
		body = strings.Replace(body, "{First}", context, -1)
		contextData := header + body + getter

		file.WriteString(contextData)
	} else {
		data, _ := ioutil.ReadFile("./Entitas/Contexts.go")

		strData := string(data)

		getter := strings.Replace(getContextTemplate, "{name}", context, -1)
		strData = strings.Replace(strData, "//NEXT2", "CreateEntityBase("+context+"ComponentTotal),\n//NEXT2", -1)
		strData = strings.Replace(strData, "//NEXT1", context+"\n//NEXT2", -1)
		strData += getter

		os.Remove("./Entitas/Contexts.go")
		file, _ := os.Create("./Entitas/Contexts.go")
		defer file.Close()

		file.WriteString(strData)
	}

}
