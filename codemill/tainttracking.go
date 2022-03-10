package codemill

import (
	. "github.com/gagliardetto/cqlgen/jen"
)

func NewTaintTrackingSpec() *TaintTrackingSpec {
	return &TaintTrackingSpec{
		Funcs:   make(map[string][]*FuncQualifier),
		Methods: make(map[string]map[string][]*FuncQualifier),
	}
}

func (spec *TaintTrackingSpec) IsEmpty() bool {
	return len(spec.Funcs) == 0 && len(spec.Methods) == 0
}

func (spec *TaintTrackingSpec) GenerateCodeQL(rootModuleGroup *Group) error {
	addedCount := 0
	tmp := DoGroup(func(tempFuncsModel *Group) {
		tempFuncsModel.Doc("Models taint-tracking through functions.")
		tempFuncsModel.Private().Class().Id("TaintTrackingFunctionModels").Extends().Qual("TaintTracking", "FunctionModel").BlockFunc(
			func(funcModelsClassGroup *Group) {
				funcModelsClassGroup.Id("FunctionInput").Id("inp").Semicolon().Line()
				funcModelsClassGroup.Id("FunctionOutput").Id("out").Semicolon().Line()

				funcModelsClassGroup.Id("TaintTrackingFunctionModels").Call().BlockFunc(
					func(funcModelsSelfMethodGroup *Group) {
						funcModelsSelfMethodGroup.DoGroup(
							func(groupCase *Group) {
								for k, v := range spec.Funcs {
									pathCodez := make([]Code, 0)
									for _, funcQual := range v {
										codeElements := GetFuncQualifierCodeElements(funcQual)
										pathCodez = append(pathCodez,
											ParensFunc(
												func(par *Group) {
													par.This().Dot("hasQualifiedName").Call(CqlFormatPackagePath(funcQual.PkgPath), Lit(funcQual.FunctionName))
													par.And()

													joined := Join(
														Or(),
														codeElements...,
													)
													if len(codeElements) > 1 {
														par.Parens(
															joined,
														)
													} else {
														par.Add(joined)
													}
												},
											),
										)
									}
									if len(pathCodez) > 0 {
										if addedCount > 0 {
											groupCase.Or()
										}
										groupCase.Commentf("Taint-tracking models for package: %s", k).Parens(
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

				funcModelsClassGroup.Override().Predicate().Id("hasTaintFlow").Call(Id("FunctionInput").Id("input"), Id("FunctionOutput").Id("output")).BlockFunc(
					func(overrideBlockGroup *Group) {
						overrideBlockGroup.Id("input").Eq().Id("inp").And().Id("output").Eq().Id("out")
					})
			})
	})
	if addedCount > 0 {
		rootModuleGroup.Add(tmp)
	}

	addedCount = 0
	tmp = DoGroup(func(tempMethodsModel *Group) {
		tempMethodsModel.Doc("Models taint-tracking through method calls.")
		tempMethodsModel.Private().Class().Id("TaintTrackingMethodModels").Extends().List(Qual("TaintTracking", "FunctionModel"), Id("Method")).BlockFunc(
			func(methodModelsClassGroup *Group) {
				methodModelsClassGroup.Id("FunctionInput").Id("inp").Semicolon().Line()
				methodModelsClassGroup.Id("FunctionOutput").Id("out").Semicolon().Line()

				methodModelsClassGroup.Id("TaintTrackingMethodModels").Call().BlockFunc(func(methodModelsSelfMethodGroup *Group) {
					methodModelsSelfMethodGroup.DoGroup(func(groupCase *Group) {
						for k, v := range spec.Methods {
							pathCodez := make([]Code, 0)
							for receiverTypeID, methodQualifiers := range v {
								codez := DoGroup(func(mtdGroup *Group) {
									mtdGroup.Commentf("Receiver type: %s", receiverTypeID)
									methodIndex := 0
									mtdGroup.ParensFunc(func(parMethods *Group) {
										for _, methodQual := range methodQualifiers {
											if methodIndex > 0 {
												parMethods.Or()
											}
											methodIndex++

											codeElements := GetFuncQualifierCodeElements(methodQual)
											parMethods.ParensFunc(func(par *Group) {
												if len(methodQual.Interface) > 0 {
													par.This().Dot("implements").Call(CqlFormatPackagePath(methodQual.PkgPath), Lit(receiverTypeID), Lit(methodQual.FunctionName))
												} else {
													par.This().Dot("hasQualifiedName").Call(CqlFormatPackagePath(methodQual.PkgPath), Lit(receiverTypeID), Lit(methodQual.FunctionName))
												}
												par.And()

												joined := Join(
													Or(),
													codeElements...,
												)
												if len(codeElements) > 1 {
													par.Parens(
														joined,
													)
												} else {
													par.Add(joined)
												}
											})
										}
									})
								})
								pathCodez = append(pathCodez, codez)

							}
							if len(pathCodez) > 0 {
								if addedCount > 0 {
									groupCase.Or()
								}
								groupCase.Commentf("Taint-tracking models for package: %s", k).Parens(
									Join(
										Or(),
										pathCodez...,
									),
								)
								addedCount++
							}
						}
					})
					methodModelsClassGroup.Override().Predicate().Id("hasTaintFlow").Call(Id("FunctionInput").Id("input"), Id("FunctionOutput").Id("output")).BlockFunc(func(overrideBlockGroup *Group) {
						overrideBlockGroup.Id("input").Eq().Id("inp").And().Id("output").Eq().Id("out")
					})
				})
			})
	})
	if addedCount > 0 {
		rootModuleGroup.Add(tmp)
	}

	return nil
}
