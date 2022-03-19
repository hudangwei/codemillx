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
var buildFlags []string

func init() {
	flag.StringVar(&moduleName, "module", "Customizations", "module name")
	flag.BoolVar(&customizeCodeQLAction, "customizeCodeQLAction", false, "customize CodeQL Action")
	flag.StringVar(&output, "output", "Customizations.qll", "output file name")
	flag.Var((*BuildFlags)(&buildFlags), "buildFlags", "build flags")
}

func main() {
	flag.Parse()

	pkgs, err := codemillx.LoadProject(flag.Args(), buildFlags)
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

type BuildFlags []string

func (v *BuildFlags) Set(s string) error {
	var err error
	*v, err = splitQuotedFields(s)
	if *v == nil {
		*v = []string{}
	}
	return err
}

func (v *BuildFlags) Get() interface{} { return *v }

func splitQuotedFields(s string) ([]string, error) {
	// Split fields allowing '' or "" around elements.
	// Quotes further inside the string do not count.
	var f []string
	for len(s) > 0 {
		for len(s) > 0 && isSpaceByte(s[0]) {
			s = s[1:]
		}
		if len(s) == 0 {
			break
		}
		// Accepted quoted string. No unescaping inside.
		if s[0] == '"' || s[0] == '\'' {
			quote := s[0]
			s = s[1:]
			i := 0
			for i < len(s) && s[i] != quote {
				i++
			}
			if i >= len(s) {
				return nil, fmt.Errorf("unterminated %c string", quote)
			}
			f = append(f, s[:i])
			s = s[i+1:]
			continue
		}
		i := 0
		for i < len(s) && !isSpaceByte(s[i]) {
			i++
		}
		f = append(f, s[:i])
		s = s[i:]
	}
	return f, nil
}

func (v *BuildFlags) String() string {
	return "<BuildFlags>"
}

func isSpaceByte(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}
