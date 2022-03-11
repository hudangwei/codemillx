package codemill

type SelectorKind string

const (
	SelectorKindStruct SelectorKind = "Struct" // Qualifier for structs.
	SelectorKindFunc   SelectorKind = "Func"   // Qualifier for funcs, type methods, interface methods.
	SelectorKindType   SelectorKind = "Type"   // Qualifier for types.
)

const (
	UntrustedFlowSourceKind = "untrust"
	TaintTrackingKind       = "tainttrack"
	SQLQueryStringSinkKind  = "sql"
	LoggerCallKind          = "logger"
	HTTPRedirectKind        = "redirect"
)

type CodeqlModuleSpec struct {
	ModuleName              string
	UntrustedFlowSourceSpec *UntrustedFlowSourceSpec
	TaintTrackingSpec       *TaintTrackingSpec
	SQLQueryStringSinkSpec  *SQLQueryStringSinkSpec
	LoggerCallSpec          *LoggerCallSpec
	HTTPRedirectSpec        *HTTPRedirectSpec
}

type UntrustedFlowSourceSpec struct {
	Funcs            map[string][]*FuncQualifier
	Methods          map[string]map[string][]*FuncQualifier
	InterfaceMethods map[string]map[string][]*FuncQualifier
	StructFieldsmap  map[string][]*StructQualifier
	Types            map[string][]*TypeQualifier
}

type TaintTrackingSpec struct {
	Funcs   map[string][]*FuncQualifier
	Methods map[string]map[string][]*FuncQualifier
}

type SQLQueryStringSinkSpec struct {
	Funcs            map[string][]*FuncQualifier
	Methods          map[string]map[string][]*FuncQualifier
	InterfaceMethods map[string]map[string][]*FuncQualifier
}

type LoggerCallSpec struct {
	Funcs            map[string][]*FuncQualifier
	Methods          map[string]map[string][]*FuncQualifier
	InterfaceMethods map[string]map[string][]*FuncQualifier
}

type HTTPRedirectSpec struct {
	Funcs   map[string][]*FuncQualifier
	Methods map[string]map[string][]*FuncQualifier
}

type Selector struct {
	Kind      SelectorKind
	Qualifier interface{}
}

type FuncQualifier struct {
	PkgPath      string
	FunctionName string
	Interface    string
	Receiver     string
	FuncDeclMetaData
	Inp FuncDeclMetaData
	Out FuncDeclMetaData
}

type FuncDeclMetaData struct {
	IsReceiver bool
	Parameters []int
	Results    []int
}

type StructQualifier struct {
	PkgPath    string
	StructName string
	Fields     []string
}

type TypeQualifier struct {
	PkgPath  string
	TypeName string
}
