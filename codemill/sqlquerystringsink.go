package codemill

import (
	. "github.com/gagliardetto/cqlgen/jen"
)

func NewSQLQueryStringSinkSpec() *SQLQueryStringSinkSpec {
	return &SQLQueryStringSinkSpec{
		Funcs:            make(map[string][]*FuncQualifier),
		Methods:          make(map[string]map[string][]*FuncQualifier),
		InterfaceMethods: make(map[string]map[string][]*FuncQualifier),
	}
}

func (spec *SQLQueryStringSinkSpec) IsEmpty() bool {
	return len(spec.Funcs) == 0 &&
		len(spec.Methods) == 0 &&
		len(spec.InterfaceMethods) == 0
}

func (spec *SQLQueryStringSinkSpec) GenerateCodeQL(moduleGroup *Group) error {
	moduleGroup.Doc("Provides models of sql querystring sink.")
	moduleGroup.Private().Class().Id("SQLQueryStringSink").Extends().List(Qual("SQL::QueryString", "Range")).
		BlockFunc(func(classGr *Group) {
			classGr.Id("SQLQueryStringSink").Call().BlockFunc(func(metGr *Group) {

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
							Int().Id("n"),
						),
						DoGroup(func(st *Group) {
							st.This().Eq().Id("fn").Dot("getACall").Call().Dot("getArgument").Call(Id("n"))
						}),
						DoGroup(func(st *Group) {
							for i, funcQual := range v {
								if i > 0 {
									st.Or()
								}

								codeElements := GenFunctionParam("n", funcQual.FuncDeclMetaData)
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
							Id("Method").Id("mtd"),
							String().Id("receiverName"),
							String().Id("methodName"),
							Int().Id("n"),
						),
						DoGroup(func(st *Group) {
							st.This().Eq().Id("mtd").Dot("getACall").Call().Dot("getArgument").Call(Id("n"))

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
											codeElements := GenFunctionParam("n", methodQual.FuncDeclMetaData)
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
							Id("Method").Id("mtd"),
							String().Id("interfaceName"),
							String().Id("methodName"),
							Int().Id("n"),
						),
						DoGroup(func(st *Group) {
							st.This().Eq().Id("mtd").Dot("getACall").Call().Dot("getArgument").Call(Id("n"))

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

											codeElements := GenFunctionParam("n", methodQual.FuncDeclMetaData)

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
			})
		})

	return nil
}
