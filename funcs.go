package codemillx

import (
	"fmt"
	"reflect"

	"github.com/hudangwei/codemillx/codemill"
)

var funcs = make(Funcs)

func init() {
	funcs["Result"] = reflect.ValueOf(Result)
	funcs["InResult"] = reflect.ValueOf(InResult)
	funcs["OutResult"] = reflect.ValueOf(OutResult)
	funcs["Param"] = reflect.ValueOf(Param)
	funcs["InParam"] = reflect.ValueOf(InParam)
	funcs["OutParam"] = reflect.ValueOf(OutParam)
	funcs["IsReceiver"] = reflect.ValueOf(IsReceiver)
	funcs["InIsReceiver"] = reflect.ValueOf(InIsReceiver)
	funcs["OutIsReceiver"] = reflect.ValueOf(OutIsReceiver)
}

type CommentGroupMetaData struct {
	ModelKind    string
	CommentFuncs []CommentFunc
}

type CommentFunc struct {
	Name   string
	Params []interface{}
}

type Funcs map[string]reflect.Value

func (f Funcs) Call(name string, modelKind string, sel *codemill.Selector, params ...interface{}) (result []reflect.Value, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	if _, ok := f[name]; !ok {
		err = fmt.Errorf("%s does not exist", name)
		return
	}
	if len(params) != f[name].Type().NumIn()-2 {
		err = fmt.Errorf("The number of params is not adapted")
		return
	}
	in := make([]reflect.Value, 0)
	//添加2个默认参数（modelKind string, sel *codemill.Selector）
	in = append(in, reflect.ValueOf(modelKind))
	in = append(in, reflect.ValueOf(sel))
	for _, param := range params {
		in = append(in, reflect.ValueOf(param))
	}
	result = f[name].Call(in)
	return
}

func Result(modelKind string, sel *codemill.Selector, args []int) {
	if sel.Kind == codemill.SelectorKindFunc {
		if fn, ok := sel.Qualifier.(*codemill.FuncQualifier); ok {
			if len(args) == 0 {
				fn.Results = []int{0}
			} else {
				fn.Results = args
			}
		}
	}
}

func InResult(modelKind string, sel *codemill.Selector, args []int) {
	if sel.Kind == codemill.SelectorKindFunc {
		if fn, ok := sel.Qualifier.(*codemill.FuncQualifier); ok {
			if len(args) == 0 {
				fn.Inp.Results = []int{0}
			} else {
				fn.Inp.Results = args
			}
		}
	}
}

func OutResult(modelKind string, sel *codemill.Selector, args []int) {
	if sel.Kind == codemill.SelectorKindFunc {
		if fn, ok := sel.Qualifier.(*codemill.FuncQualifier); ok {
			if len(args) == 0 {
				fn.Out.Results = []int{0}
			} else {
				fn.Out.Results = args
			}
		}
	}
}

func Param(modelKind string, sel *codemill.Selector, args []int) {
	if sel.Kind == codemill.SelectorKindFunc {
		if fn, ok := sel.Qualifier.(*codemill.FuncQualifier); ok {
			if len(args) == 0 {
				fn.Parameters = []int{0}
			} else {
				fn.Parameters = args
			}
		}
	}
}

func InParam(modelKind string, sel *codemill.Selector, args []int) {
	if sel.Kind == codemill.SelectorKindFunc {
		if fn, ok := sel.Qualifier.(*codemill.FuncQualifier); ok {
			if len(args) == 0 {
				fn.Inp.Parameters = []int{0}
			} else {
				fn.Inp.Parameters = args
			}
		}
	}
}

func OutParam(modelKind string, sel *codemill.Selector, args []int) {
	if sel.Kind == codemill.SelectorKindFunc {
		if fn, ok := sel.Qualifier.(*codemill.FuncQualifier); ok {
			if len(args) == 0 {
				fn.Out.Parameters = []int{0}
			} else {
				fn.Out.Parameters = args
			}
		}
	}
}

func IsReceiver(modelKind string, sel *codemill.Selector) {
	if sel.Kind == codemill.SelectorKindFunc {
		if fn, ok := sel.Qualifier.(*codemill.FuncQualifier); ok {
			fn.IsReceiver = true
		}
	}
}

func InIsReceiver(modelKind string, sel *codemill.Selector) {
	if sel.Kind == codemill.SelectorKindFunc {
		if fn, ok := sel.Qualifier.(*codemill.FuncQualifier); ok {
			fn.Inp.IsReceiver = true
		}
	}
}

func OutIsReceiver(modelKind string, sel *codemill.Selector) {
	if sel.Kind == codemill.SelectorKindFunc {
		if fn, ok := sel.Qualifier.(*codemill.FuncQualifier); ok {
			fn.Out.IsReceiver = true
		}
	}
}
