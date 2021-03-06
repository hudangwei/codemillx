package codemill

import (
	"fmt"
	"os"

	cqljen "github.com/gagliardetto/cqlgen/jen"
	"github.com/gagliardetto/utilz"
)

func GenerateCodeQL(module CodeqlModuleSpec, output string) error {
	cqlFile := cqljen.NewFile()
	cqlFile.Import("go")
	cqlFile.Private().Module().Id(utilz.ToCamel(module.ModuleName)).BlockFunc(func(moduleGroup *cqljen.Group) {
		if module.UntrustedFlowSourceSpec != nil && !module.UntrustedFlowSourceSpec.IsEmpty() {
			if err := module.UntrustedFlowSourceSpec.GenerateCodeQL(moduleGroup); err != nil {
				fmt.Println(err)
			}
		}
		if module.TaintTrackingSpec != nil && !module.TaintTrackingSpec.IsEmpty() {
			if err := module.TaintTrackingSpec.GenerateCodeQL(moduleGroup); err != nil {
				fmt.Println(err)
			}
		}
		if module.SQLQueryStringSinkSpec != nil && !module.SQLQueryStringSinkSpec.IsEmpty() {
			if err := module.SQLQueryStringSinkSpec.GenerateCodeQL(moduleGroup); err != nil {
				fmt.Println(err)
			}
		}
		if module.LoggerCallSpec != nil && !module.LoggerCallSpec.IsEmpty() {
			if err := module.LoggerCallSpec.GenerateCodeQL(moduleGroup); err != nil {
				fmt.Println(err)
			}
		}
		if module.HTTPRedirectSpec != nil && !module.HTTPRedirectSpec.IsEmpty() {
			if err := module.HTTPRedirectSpec.GenerateCodeQL(moduleGroup); err != nil {
				fmt.Println(err)
			}
		}
	})

	codeqlFile, err := os.Create(output)
	if err != nil {
		return err
	}
	defer codeqlFile.Close()
	return cqlFile.Render(codeqlFile)
}

func GetFuncQualifierCodeElements(qual *FuncQualifier) []cqljen.Code {
	codeElements := make([]cqljen.Code, 0)
	inpCodeElements := GenFunctionInputOutput("inp", qual.Inp)
	outCodeElements := GenFunctionInputOutput("out", qual.Out)
	codeElements = append(codeElements,
		cqljen.Parens(
			cqljen.Join(
				cqljen.Or(),
				inpCodeElements...,
			),
		).
			And().
			Parens(
				cqljen.Join(
					cqljen.Or(),
					outCodeElements...,
				),
			),
	)
	return codeElements
}

func GenFunctionInputOutput(idName string, funcDecl FuncDeclMetaData) []cqljen.Code {
	codeElements := make([]cqljen.Code, 0)

	if funcDecl.IsReceiver {
		codeElements = append(codeElements,
			cqljen.Id(idName).Dot("isReceiver").Call(),
		)
	}

	if len(funcDecl.Parameters) > 0 {
		codeElements = append(codeElements,
			cqljen.Id(idName).Dot("isParameter").Call(cqljen.IntsToSetOrLit(funcDecl.Parameters...)),
		)
	}

	if len(funcDecl.Results) > 0 {
		codeElements = append(codeElements,
			cqljen.Id(idName).Dot("isResult").Call(cqljen.IntsToSetOrLit(funcDecl.Results...)),
		)
	}

	return codeElements
}

func GenFunctionParam(idName string, funcDecl FuncDeclMetaData) []cqljen.Code {
	codeElements := make([]cqljen.Code, 0)

	if len(funcDecl.Parameters) > 0 {
		codeElements = append(codeElements,
			cqljen.Id(idName).Eq().Add(cqljen.IntsToSetOrLit(funcDecl.Parameters...)),
		)
	}

	return codeElements
}

func CqlFormatPackagePath(path string) cqljen.Code {
	return cqljen.Id("package").Call(cqljen.List(cqljen.Lit(path), cqljen.Lit("")))
}
