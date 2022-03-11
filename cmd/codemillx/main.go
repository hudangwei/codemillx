package main

import (
	"flag"
	"fmt"

	"github.com/hudangwei/codemillx"
	"github.com/hudangwei/codemillx/codemill"
)

var moduleName string

func init() {
	flag.StringVar(&moduleName, "module", "", "module name")
}

func main() {
	flag.Parse()

	pkgs, err := codemillx.LoadProject(flag.Args())
	if err != nil {
		fmt.Println("LoadProject with error:", err)
		return
	}

	codeqlModuleSpec := codemillx.ExtractCodeqlModuleSpec(moduleName, pkgs)

	ouputPath, err := codemill.GenerateCodeQL(codeqlModuleSpec)
	if err != nil {
		fmt.Println("GenerateCodeQL with error:", err)
		return
	}
	fmt.Println(ouputPath)
}
