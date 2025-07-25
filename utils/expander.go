package utils

import (
	"fmt"
	"os"
)

type CustomExpander struct {
	varMap map[string]string
}

func NewCustomExpander() *CustomExpander {
	return &CustomExpander{
		varMap: make(map[string]string),
	}
}

func (ce *CustomExpander) GetVar(name string) string {
	return ce.varMap[name]
}

func (ce *CustomExpander) SetVar(name string, value string) {
	ce.varMap[name] = value
}

func (ce *CustomExpander) SetVars(varMap map[string]string) {
	for name, value := range varMap {
		ce.varMap[name] = value
	}
}

func (ce *CustomExpander) Expand(s string) string {
	return os.Expand(s, ce.GetVar)
}

func (ce *CustomExpander) Environ() []string {
	var envList []string
	for name, value := range ce.varMap {
		envList = append(envList, fmt.Sprintf("%s=%s", name, value))
	}
	return envList
}

func (ce *CustomExpander) Clone() *CustomExpander {
	varMap := make(map[string]string)
	for name, value := range ce.varMap {
		varMap[name] = value
	}
	return &CustomExpander{
		varMap: varMap,
	}
}
