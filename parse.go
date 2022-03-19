package codemillx

import (
	"fmt"
	"go/ast"
	"go/token"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/hudangwei/codemillx/codemill"
	"golang.org/x/tools/go/packages"
)

func LoadProject(patterns, flags []string) ([]*packages.Package, error) {
	config := &packages.Config{
		Mode:       packages.LoadSyntax,
		BuildFlags: flags,
	}
	pkgs, err := packages.Load(config, patterns...)
	if err != nil {
		return nil, fmt.Errorf("packages.Load with error: %s", err)
	}
	if packages.PrintErrors(pkgs) > 0 {
		return nil, fmt.Errorf("packages contain errors")
	}
	return pkgs, nil
}

func ExtractCodeqlModuleSpec(moduleName string, pkgs []*packages.Package) codemill.CodeqlModuleSpec {
	untrustedFlowSourceSpec := codemill.NewUntrustedFlowSourceSpec()
	taintTrackingSpec := codemill.NewTaintTrackingSpec()
	sqlquerystringsinkSpec := codemill.NewSQLQueryStringSinkSpec()
	loggercallSpec := codemill.NewLoggerCallSpec()
	httpredirectSpec := codemill.NewHTTPRedirect()
	for _, v := range pkgs {
		models := make(map[string][]*codemill.Selector)
		for _, f := range v.Syntax {
			parseComment(v.PkgPath, f, models)
		}
		ExtractUntrustedFlowSourceSpec(models[codemill.UntrustedFlowSourceKind], untrustedFlowSourceSpec)
		ExtractTaintTrackingSpec(models[codemill.TaintTrackingKind], taintTrackingSpec)
		ExtractSQLQueryStringSinkSpec(models[codemill.SQLQueryStringSinkKind], sqlquerystringsinkSpec)
		ExtractLoggerCallSpec(models[codemill.LoggerCallKind], loggercallSpec)
		ExtractHTTPRedirectSpec(models[codemill.HTTPRedirectKind], httpredirectSpec)
	}
	return codemill.CodeqlModuleSpec{
		ModuleName:              moduleName,
		UntrustedFlowSourceSpec: untrustedFlowSourceSpec,
		TaintTrackingSpec:       taintTrackingSpec,
		SQLQueryStringSinkSpec:  sqlquerystringsinkSpec,
		LoggerCallSpec:          loggercallSpec,
		HTTPRedirectSpec:        httpredirectSpec,
	}
}

func addFunc(funcs map[string][]*codemill.FuncQualifier, name string, fn *codemill.FuncQualifier) {
	if v, ok := funcs[name]; ok {
		v = append(v, fn)
		funcs[name] = v
	} else {
		funcs[name] = []*codemill.FuncQualifier{fn}
	}
}

func addMethod(methods map[string]map[string][]*codemill.FuncQualifier, fn *codemill.FuncQualifier) {
	if v, ok := methods[fn.PkgPath]; ok {
		addFunc(v, fn.Receiver, fn)
	} else {
		vv := make(map[string][]*codemill.FuncQualifier)
		addFunc(vv, fn.Receiver, fn)
		methods[fn.PkgPath] = vv
	}
}

func addInterface(methods map[string]map[string][]*codemill.FuncQualifier, fn *codemill.FuncQualifier) {
	if v, ok := methods[fn.PkgPath]; ok {
		addFunc(v, fn.Interface, fn)
	} else {
		vv := make(map[string][]*codemill.FuncQualifier)
		addFunc(vv, fn.Interface, fn)
		methods[fn.PkgPath] = vv
	}
}

func addType(types map[string][]*codemill.TypeQualifier, typ *codemill.TypeQualifier) {
	if v, ok := types[typ.PkgPath]; ok {
		v = append(v, typ)
		types[typ.PkgPath] = v
	} else {
		types[typ.PkgPath] = []*codemill.TypeQualifier{typ}
	}
}

func addStruct(sqs map[string][]*codemill.StructQualifier, sq *codemill.StructQualifier) {
	if v, ok := sqs[sq.PkgPath]; ok {
		v = append(v, sq)
		sqs[sq.PkgPath] = v
	} else {
		sqs[sq.PkgPath] = []*codemill.StructQualifier{sq}
	}
}

func ExtractUntrustedFlowSourceSpec(sels []*codemill.Selector, spec *codemill.UntrustedFlowSourceSpec) {
	for _, sel := range sels {
		if sel.Kind == codemill.SelectorKindType {
			addType(spec.Types, sel.Qualifier.(*codemill.TypeQualifier))
		} else if sel.Kind == codemill.SelectorKindStruct {
			addStruct(spec.StructFieldsmap, sel.Qualifier.(*codemill.StructQualifier))
		} else if sel.Kind == codemill.SelectorKindFunc {
			if fn, ok := sel.Qualifier.(*codemill.FuncQualifier); ok {
				if len(fn.Interface) == 0 && len(fn.Receiver) == 0 {
					addFunc(spec.Funcs, fn.PkgPath, fn)
					continue
				}
				if len(fn.Receiver) > 0 {
					addMethod(spec.Methods, fn)
					continue
				}
				if len(fn.Interface) > 0 {
					addInterface(spec.InterfaceMethods, fn)
				}
			}
		}
	}
}

func ExtractTaintTrackingSpec(sels []*codemill.Selector, spec *codemill.TaintTrackingSpec) {
	for _, sel := range sels {
		if fn, ok := sel.Qualifier.(*codemill.FuncQualifier); ok {
			//func
			if len(fn.Interface) == 0 && len(fn.Receiver) == 0 {
				addFunc(spec.Funcs, fn.PkgPath, fn)
				continue
			}
			if len(fn.Receiver) > 0 {
				addMethod(spec.Methods, fn)
				continue
			}
			if len(fn.Interface) > 0 {
				addInterface(spec.Methods, fn)
			}
		}
	}
}

func ExtractSQLQueryStringSinkSpec(sels []*codemill.Selector, spec *codemill.SQLQueryStringSinkSpec) {
	for _, sel := range sels {
		if sel.Kind == codemill.SelectorKindFunc {
			if fn, ok := sel.Qualifier.(*codemill.FuncQualifier); ok {
				if len(fn.Interface) == 0 && len(fn.Receiver) == 0 {
					addFunc(spec.Funcs, fn.PkgPath, fn)
					continue
				}
				if len(fn.Receiver) > 0 {
					addMethod(spec.Methods, fn)
					continue
				}
				if len(fn.Interface) > 0 {
					addInterface(spec.InterfaceMethods, fn)
				}
			}
		}
	}
}

func ExtractLoggerCallSpec(sels []*codemill.Selector, spec *codemill.LoggerCallSpec) {
	for _, sel := range sels {
		if sel.Kind == codemill.SelectorKindFunc {
			if fn, ok := sel.Qualifier.(*codemill.FuncQualifier); ok {
				if len(fn.Interface) == 0 && len(fn.Receiver) == 0 {
					addFunc(spec.Funcs, fn.PkgPath, fn)
					continue
				}
				if len(fn.Receiver) > 0 {
					addMethod(spec.Methods, fn)
					continue
				}
				if len(fn.Interface) > 0 {
					addInterface(spec.InterfaceMethods, fn)
				}
			}
		}
	}
}

func ExtractHTTPRedirectSpec(sels []*codemill.Selector, spec *codemill.HTTPRedirectSpec) {
	for _, sel := range sels {
		if sel.Kind == codemill.SelectorKindFunc {
			if fn, ok := sel.Qualifier.(*codemill.FuncQualifier); ok {
				if len(fn.Interface) == 0 && len(fn.Receiver) == 0 {
					addFunc(spec.Funcs, fn.PkgPath, fn)
					continue
				}
				if len(fn.Receiver) > 0 {
					addMethod(spec.Methods, fn)
					continue
				}
				if len(fn.Interface) > 0 {
					addInterface(spec.Methods, fn)
				}
			}
		}
	}
}

func addModel(models map[string][]*codemill.Selector, kind string, sel *codemill.Selector) {
	if v, ok := models[kind]; ok {
		v = append(v, sel)
		models[kind] = v
	} else {
		models[kind] = []*codemill.Selector{sel}
	}
}

func newSelector(commentFuncs []CommentFunc, sel *codemill.Selector) *codemill.Selector {
	var isDependsPkg, isFunc, isType, isField bool
	for _, v := range commentFuncs {
		if v.Name == "Package" {
			isDependsPkg = true
		}
		if v.Name == "Func" || v.Name == "Method" || v.Name == "Interface" {
			isFunc = true
		}
		if v.Name == "Type" {
			isType = true
		}
		if v.Name == "Field" {
			isField = true
		}
	}
	if isDependsPkg && isType {
		return &codemill.Selector{
			Kind:      codemill.SelectorKindType,
			Qualifier: &codemill.TypeQualifier{},
		}
	}
	if isDependsPkg && isField {
		return &codemill.Selector{
			Kind:      codemill.SelectorKindStruct,
			Qualifier: &codemill.StructQualifier{},
		}
	}
	if isDependsPkg && isFunc {
		return &codemill.Selector{
			Kind:      codemill.SelectorKindFunc,
			Qualifier: &codemill.FuncQualifier{},
		}
	}

	return sel
}

func parseComment(pkgPath string, f *ast.File, models map[string][]*codemill.Selector) {
	for _, d := range f.Decls {
		switch specDecl := d.(type) {
		case *ast.GenDecl:
			if specDecl.Tok == token.TYPE {
				for _, s := range specDecl.Specs {
					specName := s.(*ast.TypeSpec).Name.String()
					if specDecl.Doc != nil {
						if metaData, err := extractCommentGroupMetaData(specDecl.Doc, "type"); err == nil && len(metaData) > 0 {
							for _, v := range metaData {
								sel := newSelector(v.CommentFuncs, &codemill.Selector{
									Kind: codemill.SelectorKindType,
									Qualifier: &codemill.TypeQualifier{
										PkgPath:  pkgPath,
										TypeName: specName,
									},
								})
								for _, vv := range v.CommentFuncs {
									if _, err := funcs.Call(vv.Name, v.ModelKind, sel, vv.Params...); err != nil {
										fmt.Printf("call %s with error:%v\n", vv.Name, err)
									}
								}
								addModel(models, v.ModelKind, sel)
							}
						}
					}
					switch tp := s.(*ast.TypeSpec).Type.(type) {
					case *ast.StructType:
						var kind string
						var fields []string
						for _, fd := range tp.Fields.List {
							if len(fd.Names) == 0 {
								continue
							}
							if fd.Comment != nil {
								if metaData, err := extractCommentGroupMetaData(fd.Comment, "field"); err == nil && len(metaData) == 1 {
									fields = append(fields, fd.Names[0].Name)
									kind = metaData[0].ModelKind
								}
							}
						}
						if len(fields) > 0 && len(kind) > 0 {
							sel := &codemill.Selector{
								Kind: codemill.SelectorKindStruct,
								Qualifier: &codemill.StructQualifier{
									PkgPath:    pkgPath,
									StructName: specName,
									Fields:     fields,
								},
							}
							addModel(models, kind, sel)
						}
					case *ast.InterfaceType:
						for _, meth := range tp.Methods.List {
							if len(meth.Names) == 0 {
								continue
							}
							if meth.Comment != nil {
								if metaData, err := extractCommentGroupMetaData(meth.Comment, "method"); err == nil && len(metaData) > 0 {
									for _, v := range metaData {
										sel := &codemill.Selector{
											Kind: codemill.SelectorKindFunc,
											Qualifier: &codemill.FuncQualifier{
												PkgPath:      pkgPath,
												FunctionName: meth.Names[0].Name,
												Interface:    specName,
											},
										}
										for _, vv := range v.CommentFuncs {
											if _, err := funcs.Call(vv.Name, v.ModelKind, sel, vv.Params...); err != nil {
												fmt.Printf("call %s with error:%v\n", vv.Name, err)
											}
										}
										addModel(models, v.ModelKind, sel)
									}
								}
							}
						}
					}
				}
			}
		case *ast.FuncDecl:
			if specDecl.Doc != nil {
				if metaData, err := extractCommentGroupMetaData(specDecl.Doc, "method"); err == nil && len(metaData) > 0 {
					rec := ""
					if specDecl.Recv != nil {
						if specDecl.Recv != nil && len(specDecl.Recv.List) > 0 {
							if star, ok := specDecl.Recv.List[0].Type.(*ast.StarExpr); ok {
								rec = star.X.(*ast.Ident).Name
							}
							if ident, ok := specDecl.Recv.List[0].Type.(*ast.Ident); ok {
								rec = ident.Name
							}
						}
					}

					for _, v := range metaData {
						sel := newSelector(v.CommentFuncs, &codemill.Selector{
							Kind: codemill.SelectorKindFunc,
							Qualifier: &codemill.FuncQualifier{
								PkgPath:      pkgPath,
								FunctionName: specDecl.Name.String(),
								Receiver:     rec,
							},
						})
						for _, vv := range v.CommentFuncs {
							if _, err := funcs.Call(vv.Name, v.ModelKind, sel, vv.Params...); err != nil {
								fmt.Printf("call %s with error:%v\n", vv.Name, err)
							}
						}
						addModel(models, v.ModelKind, sel)
					}
				} else if err != nil {
					fmt.Println(err)
				}
			}
		}
	}
}

func extractCommentGroupMetaData(f *ast.CommentGroup, typ string) (ret []CommentGroupMetaData, err error) {
	for _, d := range f.List {
		t := strings.TrimSpace(strings.TrimPrefix(d.Text, "//"))
		content := strings.Split(t, " ")
		if len(content) > 1 && content[0] == "@codeql" {
			var cfs []CommentFunc
			for i := 2; i < len(content); i++ {
				var cf CommentFunc
				cf, err = parseFunc(content[i])
				if err != nil {
					return
				}
				cfs = append(cfs, cf)
			}
			ret = append(ret, CommentGroupMetaData{
				ModelKind:    content[1],
				CommentFuncs: cfs,
			})
		}
	}
	return
}

func parseFunc(vfunc string) (cf CommentFunc, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	vfunc = strings.TrimSpace(vfunc)
	start := strings.Index(vfunc, "(")
	var num int
	if start == -1 {
		if num, err = numIn(vfunc); err != nil {
			return
		}
		if num != 0 {
			err = fmt.Errorf("%s require %d parameters", vfunc, num)
			return
		}
		cf.Name = vfunc
		return
	}

	end := strings.Index(vfunc, ")")
	if end == -1 {
		err = fmt.Errorf("invalid valid function")
		return
	}

	name := strings.TrimSpace(vfunc[:start])
	params := strings.Split(vfunc[start+1:end], ",")
	tParams, err := trim(name, params)
	if err != nil {
		return
	}
	cf.Name = name
	cf.Params = tParams
	return
}

func numIn(name string) (num int, err error) {
	fn, ok := funcs[name]
	if !ok {
		err = fmt.Errorf("doesn't exists %s valid function", name)
		return
	}
	//去除前两个默认参数，比如：func Result(modelKind string, sel *codemill.Selector, args []int)
	num = fn.Type().NumIn() - 2
	return
}

func trim(name string, s []string) ([]interface{}, error) {
	fn, ok := funcs[name]
	if !ok {
		return nil, fmt.Errorf("doesn't exists %s valid function", name)
	}
	if len(s) > 0 && fn.Type().In(2).Kind() == reflect.Slice && fn.Type().In(2).Elem().String() == "int" {
		var param []int
		for _, v := range s {
			if n, err := strconv.Atoi(v); err == nil {
				param = append(param, n)
			}
		}
		return []interface{}{param}, nil
	}
	if len(s) > 0 && fn.Type().NumIn() == 4 && fn.Type().In(2).Kind() == reflect.String && fn.Type().In(3).Kind() == reflect.Slice && fn.Type().In(3).Elem().String() == "string" {
		return []interface{}{s[0], s[1:]}, nil
	}
	var ret []interface{}
	for i := 0; i < len(s); i++ {
		var param interface{}
		param, err := parseParam(fn.Type().In(i+2), strings.TrimSpace(s[i]))
		if err != nil {
			return nil, err
		}
		ret = append(ret, param)
	}
	return ret, nil
}

func parseParam(t reflect.Type, s string) (i interface{}, err error) {
	switch t.Kind() {
	case reflect.Int:
		i, err = strconv.Atoi(s)
	case reflect.Int64:
		i, err = strconv.ParseInt(s, 10, 64)
	case reflect.Int32:
		var v int64
		v, err = strconv.ParseInt(s, 10, 32)
		if err == nil {
			i = int32(v)
		}
	case reflect.Int16:
		var v int64
		v, err = strconv.ParseInt(s, 10, 16)
		if err == nil {
			i = int16(v)
		}
	case reflect.Int8:
		var v int64
		v, err = strconv.ParseInt(s, 10, 8)
		if err == nil {
			i = int8(v)
		}
	case reflect.String:
		i = s
	case reflect.Ptr:
		if t.Elem().String() != "regexp.Regexp" {
			err = fmt.Errorf("not support %s", t.Elem().String())
			return
		}
		i, err = regexp.Compile(s)
	default:
		err = fmt.Errorf("not support %s", t.Kind().String())
	}
	return
}
