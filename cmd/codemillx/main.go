package main

import (
	"flag"
	"fmt"

	"github.com/hudangwei/codemillx"
	"github.com/hudangwei/codemillx/codemill"
)

var moduleName string
var customizeCodeQLAction bool
var output string

func init() {
	flag.StringVar(&moduleName, "module", "Customizations", "module name")
	flag.BoolVar(&customizeCodeQLAction, "customizeCodeQLAction", false, "customize CodeQL Action")
	flag.StringVar(&output, "output", "Customizations.qll", "output file name")
}

func main() {
	flag.Parse()

	pkgs, err := codemillx.LoadProject(flag.Args())
	if err != nil {
		fmt.Println("LoadProject with error:", err)
		return
	}

	codeqlModuleSpec := codemillx.ExtractCodeqlModuleSpec(moduleName, pkgs)

	if customizeCodeQLAction {
		if customizeFile := codemillx.FindCustomizationsQllFile(); len(customizeFile) > 0 {
			output = customizeFile
		}
	}

	if err := codemill.GenerateCodeQL(codeqlModuleSpec, output); err != nil {
		fmt.Println("GenerateCodeQL with error:", err)
		return
	}
}
