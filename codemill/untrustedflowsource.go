package codemill

import (
	"sort"

	. "github.com/gagliardetto/cqlgen/jen"
)

func NewUntrustedFlowSourceSpec() *UntrustedFlowSourceSpec {
	return &UntrustedFlowSourceSpec{
		Funcs:            make(map[string][]*FuncQualifier),
		Methods:          make(map[string]map[string][]*FuncQualifier),
		InterfaceMethods: make(map[string]map[string][]*FuncQualifier),
		StructFieldsmap:  make(map[string][]*StructQualifier),
		Types:            make(map[string][]*TypeQualifier),
	}
}

func (spec *UntrustedFlowSourceSpec) IsEmpty() bool {
	return len(spec.Funcs) == 0 &&
		len(spec.Methods) == 0 &&
		len(spec.InterfaceMethods) == 0 &&
		len(spec.StructFieldsmap) == 0 &&
		len(spec.Types) == 0
}

func (spec *UntrustedFlowSourceSpec) GenerateCodeQL(moduleGroup *Group) error {
	moduleGroup.Doc("Provides models of untrusted flow sources.")
	moduleGroup.Private().Class().Id("UntrustedSources").Extends().List(Qual("UntrustedFlowSource", "Range")).
		BlockFunc(func(classGr *Group) {
			classGr.Id("UntrustedSources").Call().BlockFunc(func(metGr *Group) {

				//funcs
				index := 0
				for k, v := range spec.Funcs {
					if index > 0 {
						metGr.Or()
					}
					index++

					metGr.Comment("Functions of package: " + k)
					metGr.Exists(
						List(
							Id("Function").Id("fn"),
							Id("FunctionOutput").Id("out"),
						),
						DoGroup(func(st *Group) {
							st.This().Eq().Id("out").Dot("getExitNode").Call(Id("fn").Dot("getACall").Call())
						}),
						DoGroup(func(st *Group) {
							for i, funcQual := range v {
								if i > 0 {
									st.Or()
								}

								codeElements := GenFunctionInputOutput("out", funcQual.FuncDeclMetaData)
								st.Id("fn").Dot("hasQualifiedName").Call(CqlFormatPackagePath(funcQual.PkgPath), Lit(funcQual.FunctionName)).
									And().
									Parens(
										Join(
											Or(),
											codeElements...,
										),
									)
							}
						}),
					)
				}

				if len(spec.Funcs) > 0 && len(spec.Methods) > 0 {
					metGr.Or()
				}

				//methods
				index = 0
				for k, v := range spec.Methods {
					if index > 0 {
						metGr.Or()
					}
					index++

					metGr.Comment("Methods on types of package: " + k)
					metGr.Exists(
						List(
							String().Id("receiverName"),
							String().Id("methodName"),
							Id("Method").Id("mtd"),
							Id("FunctionOutput").Id("out"),
						),
						DoGroup(func(st *Group) {
							st.This().Eq().Id("out").Dot("getExitNode").Call(Id("mtd").Dot("getACall").Call())

							st.And()

							st.Id("mtd").Dot("hasQualifiedName").Call(
								CqlFormatPackagePath(k),
								Id("receiverName"),
								Id("methodName"),
							)
						}),
						DoGroup(func(st *Group) {
							typeIndex := 0
							for receiverTypeID, methodQualifiers := range v {
								if typeIndex > 0 {
									st.Or()
								}
								typeIndex++

								st.Id("receiverName").Eq().Lit(receiverTypeID)
								st.And()

								st.ParensFunc(
									func(parMethods *Group) {
										for i, methodQual := range methodQualifiers {
											if i > 0 {
												parMethods.Or()
											}
											codeElements := GenFunctionInputOutput("out", methodQual.FuncDeclMetaData)
											parMethods.ParensFunc(
												func(par *Group) {
													par.Id("methodName").Eq().Lit(methodQual.FunctionName)
													par.And()
													par.Parens(
														Join(
															Or(),
															codeElements...,
														),
													)
												},
											)
										}
									},
								)
							}
						}),
					)
				}

				if (len(spec.Funcs) > 0 || len(spec.Methods) > 0) && len(spec.InterfaceMethods) > 0 {
					metGr.Or()
				}

				index = 0
				for k, v := range spec.InterfaceMethods {
					if index > 0 {
						metGr.Or()
					}
					index++

					metGr.Comment("Interfaces of package: " + k)
					metGr.Exists(
						List(
							String().Id("interfaceName"),
							String().Id("methodName"),
							Id("Method").Id("mtd"),
							Id("FunctionOutput").Id("out"),
						),
						DoGroup(func(st *Group) {
							st.This().Eq().Id("out").Dot("getExitNode").Call(Id("mtd").Dot("getACall").Call())

							st.And()

							st.Id("mtd").Dot("implements").Call(
								CqlFormatPackagePath(k),
								Id("interfaceName"),
								Id("methodName"),
							)
						}),
						DoGroup(func(st *Group) {
							typeIndex := 0
							for receiverTypeID, methodQualifiers := range v {
								if typeIndex > 0 {
									st.Or()
								}
								typeIndex++

								st.Id("interfaceName").Eq().Lit(receiverTypeID)

								st.And()

								st.ParensFunc(
									func(parMethods *Group) {
										for i, methodQual := range methodQualifiers {
											if i > 0 {
												parMethods.Or()
											}

											codeElements := GenFunctionInputOutput("out", methodQual.FuncDeclMetaData)

											parMethods.ParensFunc(
												func(par *Group) {
													par.Id("methodName").Eq().Lit(methodQual.FunctionName)
													par.And()
													par.Parens(
														Join(
															Or(),
															codeElements...,
														),
													)
												},
											)

										}
									},
								)
							}
						}),
					)
				}

				if (len(spec.Funcs) > 0 || len(spec.Methods) > 0 || len(spec.InterfaceMethods) > 0) && len(spec.StructFieldsmap) > 0 {
					metGr.Or()
				}

				index = 0
				for k, v := range spec.StructFieldsmap {
					if index > 0 {
						metGr.Or()
					}
					index++
					metGr.Comment("Structs of package: " + k)
					metGr.Exists(
						List(
							String().Id("structName"),
							String().Id("fields"),
							Qual("DataFlow", "Field").Id("fld"),
						),
						DoGroup(func(st *Group) {
							st.This().Eq().Id("fld").Dot("getARead").Call()

							st.And()

							st.Id("fld").Dot("hasQualifiedName").Call(
								CqlFormatPackagePath(k),
								Id("structName"),
								Id("fields"),
							)
						}),
						DoGroup(func(st *Group) {
							for qualIndex, qual := range v {
								if qualIndex > 0 {
									st.Or()
								}
								st.Id("structName").Eq().Lit(qual.StructName)
								st.And()
								st.Id("fields").Eq().Add(StringsToSetOrLit(qual.Fields...))
							}
						}),
					)
				}

				if (len(spec.Funcs) > 0 || len(spec.Methods) > 0 || len(spec.InterfaceMethods) > 0 || len(spec.StructFieldsmap) > 0) && len(spec.Types) > 0 {
					metGr.Or()
				}

				index = 0
				for k, v := range spec.Types {
					if index > 0 {
						metGr.Or()
					}
					index++
					metGr.Comment("Types of package: " + k)
					metGr.Exists(
						List(
							Id("ValueEntity").Id("v"),
						),
						DoGroup(func(st *Group) {
							var typeNames []string
							for _, qual := range v {
								typeNames = append(typeNames, qual.TypeName)
							}
							sort.Strings(typeNames)

							st.Id("v").Dot("getType").Call().Dot("hasQualifiedName").Call(
								CqlFormatPackagePath(k),
								StringsToSetOrLit(typeNames...),
							)
						}),
						DoGroup(func(st *Group) {
							st.This().Eq().Id("v").Dot("getARead").Call()
						}),
					)
				}

			})
		})

	return nil
}
