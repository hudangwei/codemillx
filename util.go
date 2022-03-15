package codemillx

import (
	"os"
	"path/filepath"
)

func FindCustomizationsQllFile() string {
	matches, err := filepath.Glob("/opt/hostedtoolcache/CodeQL/*/x64/codeql/qlpacks/codeql/go-all/*")
	if err != nil || len(matches) == 0 {
		return ""
	}
	for _, v := range matches {
		f := filepath.Join(v, "Customizations.qll")
		if FileExists(f) {
			return f
		}
	}
	return ""
}

func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
