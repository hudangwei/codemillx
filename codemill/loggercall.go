package codemill

import (
	"sort"

	. "github.com/gagliardetto/cqlgen/jen"
)

func NewLoggerCallSpec() *LoggerCallSpec {
	return &LoggerCallSpec{
		Funcs:            make(map[string][]*FuncQualifier),
		Methods:          make(map[string]map[string][]*FuncQualifier),
		InterfaceMethods: make(map[string]map[string][]*FuncQualifier),
	}
}

func (spec *LoggerCallSpec) IsEmpty() bool {
	return len(spec.Funcs) == 0 &&
		len(spec.Methods) == 0 &&
		len(spec.InterfaceMethods) == 0
}

func (spec *LoggerCallSpec) GenerateCodeQL(rootModuleGroup *Group) error {
	addedCount := 0
	tmp := DoGroup(func(tempFuncsModel *Group) {
		tempFuncsModel.Doc("Models LoggerCall functions.")
		tempFuncsModel.Private().Class().Id("LoggerCallFunctions").Extends().List(Qual("LoggerCall", "Range"), Qual("DataFlow", "CallNode")).BlockFunc(
			func(funcModelsClassGroup *Group) {
				funcModelsClassGroup.Id("LoggerCallFunctions").Call().BlockFunc(
					func(funcModelsSelfMethodGroup *Group) {
						funcModelsSelfMethodGroup.DoGroup(
							func(groupCase *Group) {
								for k, v := range spec.Funcs {
									pathCodez := make([]Code, 0)
									fnNames := make([]string, 0)
									for _, funcQual := range v {
										fnNames = append(fnNames, funcQual.FunctionName)
									}
									sort.Strings(fnNames)
									pathCodez = append(pathCodez,
										ParensFunc(
											func(par *Group) {
												par.This().Dot("getTarget").Call().Dot("hasQualifiedName").Call(CqlFormatPackagePath(k), StringsToSetOrLit(fnNames...))
											},
										),
									)

									if len(pathCodez) > 0 {
										if addedCount > 0 {
											groupCase.Or()
										}
										groupCase.Commentf("LoggerCall models for package: %s", k).Parens(
											Join(
												Or(),
												pathCodez...,
											),
										)
										addedCount++
									}
								}
							})
					})

				funcModelsClassGroup.Override().Qual("DataFlow", "Node").Id("getAMessageComponent").Call().BlockFunc(
					func(overrideBlockGroup *Group) {
						overrideBlockGroup.Id("result").Eq().This().Dot("getAnArgument").Call()
					})
			})
	})
	if addedCount > 0 {
		rootModuleGroup.Add(tmp)
	}

	addedCount = 0
	tmp = DoGroup(func(tempMethodsModel *Group) {
		tempMethodsModel.Doc("Models LoggerCall methods.")
		tempMethodsModel.Private().Class().Id("LoggerCallMethods").Extends().List(Qual("LoggerCall", "Range"), Qual("DataFlow", "MethodCallNode")).BlockFunc(
			func(methodModelsClassGroup *Group) {
				methodModelsClassGroup.Id("LoggerCallMethods").Call().BlockFunc(func(methodModelsSelfMethodGroup *Group) {
					methodModelsSelfMethodGroup.DoGroup(func(groupCase *Group) {
						for k, v := range spec.Methods {
							pathCodez := make([]Code, 0)
							for receiverTypeID, methodQualifiers := range v {
								codez := DoGroup(func(mtdGroup *Group) {
									mtdGroup.Commentf("Receiver type: %s", receiverTypeID)
									mtdGroup.ParensFunc(func(parMethods *Group) {
										methNames := make([]string, 0)
										for _, methodQual := range methodQualifiers {
											methNames = append(methNames, methodQual.FunctionName)
										}
										sort.Strings(methNames)
										parMethods.ParensFunc(func(par *Group) {
											par.This().Dot("getTarget").Call().Dot("hasQualifiedName").Call(CqlFormatPackagePath(k), Lit(receiverTypeID), StringsToSetOrLit(methNames...))
										})
									})
								})
								pathCodez = append(pathCodez, codez)
							}
							if len(pathCodez) > 0 {
								if addedCount > 0 {
									groupCase.Or()
								}
								groupCase.Commentf("LoggerCall models for package: %s", k).Parens(
									Join(
										Or(),
										pathCodez...,
									),
								)
								addedCount++
							}
						}
						for k, v := range spec.InterfaceMethods {
							pathCodez := make([]Code, 0)
							for receiverTypeID, methodQualifiers := range v {
								codez := DoGroup(func(mtdGroup *Group) {
									mtdGroup.Commentf("Interface type: %s", receiverTypeID)
									mtdGroup.ParensFunc(func(parMethods *Group) {
										methNames := make([]string, 0)
										for _, methodQual := range methodQualifiers {
											methNames = append(methNames, methodQual.FunctionName)
										}
										sort.Strings(methNames)
										parMethods.ParensFunc(func(par *Group) {
											par.This().Dot("getTarget").Call().Dot("implements").Call(CqlFormatPackagePath(k), Lit(receiverTypeID), StringsToSetOrLit(methNames...))
										})
									})
								})
								pathCodez = append(pathCodez, codez)
							}
							if len(pathCodez) > 0 {
								if addedCount > 0 {
									groupCase.Or()
								}
								groupCase.Commentf("LoggerCall models for package: %s", k).Parens(
									Join(
										Or(),
										pathCodez...,
									),
								)
								addedCount++
							}
						}
					})
					methodModelsClassGroup.Override().Qual("DataFlow", "Node").Id("getAMessageComponent").Call().BlockFunc(
						func(overrideBlockGroup *Group) {
							overrideBlockGroup.Id("result").Eq().This().Dot("getAnArgument").Call()
						})
				})
			})
	})
	if addedCount > 0 {
		rootModuleGroup.Add(tmp)
	}

	return nil
}
