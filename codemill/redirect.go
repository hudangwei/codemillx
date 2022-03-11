package codemill

import (
	. "github.com/gagliardetto/cqlgen/jen"
)

func NewHTTPRedirect() *HTTPRedirectSpec {
	return &HTTPRedirectSpec{
		Funcs:   make(map[string][]*FuncQualifier),
		Methods: make(map[string]map[string][]*FuncQualifier),
	}
}

func (spec *HTTPRedirectSpec) IsEmpty() bool {
	return len(spec.Funcs) == 0 &&
		len(spec.Methods) == 0
}

func (spec *HTTPRedirectSpec) GenerateCodeQL(rootModuleGroup *Group) error {
	addedCount := 0
	tmp := DoGroup(func(tempFuncsModel *Group) {
		tempFuncsModel.Doc("Models HTTP redirects.")
		tempFuncsModel.Private().Class().Id("HttpRedirect").Extends().List(
			Id("HTTP::Redirect::Range"),
			Id("DataFlow::CallNode"),
		).BlockFunc(
			func(funcModelsClassGroup *Group) {
				funcModelsClassGroup.String().Id("package").Semicolon().Line()
				funcModelsClassGroup.Id("DataFlow::Node").Id("urlNode").Semicolon().Line()
				funcModelsClassGroup.Id("HttpRedirect").Call().BlockFunc(
					func(funcModelsSelfMethodGroup *Group) {
						{
							funcModelsSelfMethodGroup.DoGroup(
								func(groupCase *Group) {
									for k, v := range spec.Funcs {
										pathCodez := make([]Code, 0)
										for _, funcQual := range v {
											pathCodez = append(pathCodez,
												ParensFunc(
													func(par *Group) {
														par.This().
															Dot("getTarget").Call().
															Dot("hasQualifiedName").Call(
															Id("package"),
															Lit(funcQual.FunctionName),
														)

														par.And()

														par.Id("urlNode").Eq().Add(Id("this").Dot("getArgument").Call(IntsToSetOrLit(funcQual.Parameters...)))
													},
												),
											)
										}

										if len(pathCodez) > 0 {
											if addedCount > 0 {
												groupCase.Or()
											}
											groupCase.Commentf("HTTP redirect models for package: %s", k)
											groupCase.Id("package").Eq().Add(CqlFormatPackagePath(k)).And()

											groupCase.Parens(
												Join(
													Or(),
													pathCodez...,
												),
											)

											addedCount++
										}
									}
									for k, v := range spec.Methods {
										pathCodez := make([]Code, 0)
										for receiverTypeID, methodQualifiers := range v {
											codez := DoGroup(func(mtdGroup *Group) {
												methodIndex := 0
												mtdGroup.ParensFunc(
													func(parMethods *Group) {
														for _, methodQual := range methodQualifiers {
															if methodIndex > 0 {
																parMethods.Or()
															}
															methodIndex++

															parMethods.ParensFunc(
																func(par *Group) {
																	par.This().
																		Eq().
																		Any(
																			DoGroup(func(gr *Group) {
																				gr.Id("Method").Id("m")
																			}),
																			DoGroup(func(gr *Group) {
																				if len(methodQual.Interface) > 0 {
																					gr.Id("m").Dot("implements").Call(
																						Id("package"),
																						Lit(receiverTypeID),
																						Lit(methodQual.FunctionName),
																					)
																				} else {
																					gr.Id("m").Dot("hasQualifiedName").Call(
																						Id("package"),
																						Lit(receiverTypeID),
																						Lit(methodQual.FunctionName),
																					)
																				}
																			}),
																			nil,
																		).Dot("getACall").Call()

																	par.And()

																	par.Id("urlNode").Eq().Add(Id("this").Dot("getArgument").Call(IntsToSetOrLit(methodQual.Parameters...)))
																},
															)
														}
													},
												)
											})
											pathCodez = append(pathCodez, codez)
										}
										if len(pathCodez) > 0 {
											if addedCount > 0 {
												groupCase.Or()
											}
											groupCase.Commentf("HTTP redirect models for package: %s", k)
											groupCase.Id("package").Eq().Add(CqlFormatPackagePath(k)).And()

											groupCase.Parens(
												Join(
													Or(),
													pathCodez...,
												),
											)

											addedCount++
										}
									}
								})
						}
					})

				funcModelsClassGroup.Override().Id("DataFlow::Node").Id("getUrl").Call().BlockFunc(
					func(overrideBlockGroup *Group) {
						overrideBlockGroup.Id("result").Eq().Id("urlNode")
					})

				funcModelsClassGroup.Override().Id("HTTP::ResponseWriter").Id("getResponseWriter").Call().BlockFunc(
					func(overrideBlockGroup *Group) {
						overrideBlockGroup.Id("result").Dot("getANode").Call().Eq().This().Dot("getReceiver").Call()
					})
			})
	})
	if addedCount > 0 {
		rootModuleGroup.Add(tmp)
	}
	return nil
}
