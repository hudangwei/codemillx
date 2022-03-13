如何添加代码注释？
===

[English](comment.md) | 简体中文

## 目前支持的codeql-models

* [标记污点源-UntrustedFlowSource](#untrustedflowsource)
* [标记污点传播函数-TaintTracking::FunctionModel](#tainttrackingfunctionmodel)
* [标记SQL语句-SQL::QueryString](#sqlquerystring)
* [标记日志打印函数-LoggerCall](#loggercall)
* [标记HTTP重定向URL-HTTP::Redirect](#httpredirect)

## UntrustedFlowSource

标记污点源，可提升CodeQL检测能力

### 注释格式

//@codeql untrust (Param(0)| Result(0) | IsReceiver)

* Param(0) 表示第一个参数是不可信的污点源(用户可控)
* Result(0) 表示第一个返回值是不可信的污点源(用户可控)
* IsReceiver 表示方法接收者是不可信的污点源(用户可控)
* Param(0)，Result(0)，IsReceiver 可同时多个一起注释，表示第一个参数，第一个返回值，方法接收者都是不可信的污点源(用户可控)

### 例子

* 在type定义上标记，代表整个结构体的所有字段都是不可信的用户可控输入。

```go
//@codeql untrust
type Param struct {
 Key   string
 Value string
}
```

* 在结构体(struct)的字段(field)上标记，代表该字段是不可信的用户可控输入。

```go
type Context struct {
 app      *Application
 Params   Params //@codeql untrust
 Route    *Route
 Request  *http.Request
 Response http.ResponseWriter
 query    url.Values
}
```

* 在方法上添加注释

```go
//@codeql untrust Param(0)
func (c *Context) Decode(v interface{}) (err error) {
 if c.app.Decoder == nil {
  return ErrDecoderNotRegister
 }
 return c.app.Decoder.Decode(c.Request, v)
}

//@codeql untrust Result(0)
func (c *Context) QueryParam(key string) string {
 return c.QueryParams().Get(key)
}
```

* 在函数上添加注释

```go
// Vars returns the route variables for the current request, if any.
//@codeql untrust Result(0)
func Vars(r *http.Request) map[string]string {
 if rv := r.Context().Value(varsKey); rv != nil {
  return rv.(map[string]string)
 }
 return nil
}
```

* 在接口方法上添加注释

```go
// Decoder is an interface that decodes request's input.
type Decoder interface {
 // Decode decodes request's input and stores it in the value pointed to by v.
 Decode(req *http.Request, v interface{}) error //@codeql untrust Param(1)
}
```

## TaintTracking::FunctionModel

标记污点传播函数，可提升CodeQL检测能力

### 注释格式

//@codeql tainttrack InParam(0) OutResult(0)

* InParam(0) 表示第一个参数是污点源
* InIsReceiver 表示方法接收者是污点源
* OutParam(1) 表示第二个参数会被污点源传染，也将成为污点源
* OutResult(1) 表示第二个返回值会被污点源传染，也将成为污点源
* OutIsReceiver 表示方法接收者会被污点源传染，也将成为污点源

### 例子

* 在方法/函数/接口方法上标记污点传播函数

```go
//@codeql tainttrack InParam(1) OutResult(0)
func NewError(code int, message ...string) *Error {
 err := &Error{
  Code:    code,
  Message: utils.StatusMessage(code),
 }
 if len(message) > 0 {
  err.Message = message[0]
 }
 return err
}

//@codeql tainttrack InParam(0) OutResult(0)
func TrimBytes(b []byte, cutset byte) []byte {
 i, j := 0, len(b)-1
 for ; i <= j; i++ {
  if b[i] != cutset {
   break
  }
 }
 for ; i < j; j-- {
  if b[j] != cutset {
   break
  }
 }

 return b[i : j+1]
}
```

## SQL::QueryString

标记SQL语句，可用于发现SQL注入

### 注释格式

//@codeql sql Param(0)

### 例子

* 在方法/函数/接口方法上标记SQL语句参数位置
* 例子中Select方法的第二个参数就是SQL语句，所以是Param(1)

```go
//@codeql sql Param(1)
func (db *DB) Select(dest interface{}, query string, args ...interface{}) error {
 return Select(db, dest, query, args...)
}
```

## LoggerCall

标记日志打印函数，可用于发现敏感信息泄露

### 注释格式

//@codeql logger

### 例子

* 在日志库打印方法上添加注释

```go
//@codeql logger
func (logger *Logger) Log(level Level, args ...interface{}) {
 if logger.IsLevelEnabled(level) {
  entry := logger.newEntry()
  entry.Log(level, args...)
  logger.releaseEntry(entry)
 }
}
```

* 在日志库打印函数上添加注释

```go
//@codeql logger
func Debug(args ...interface{}) {
 std.Debug(args...)
}

//@codeql logger
func Print(args ...interface{}) {
 std.Print(args...)
}

//@codeql logger
func Info(args ...interface{}) {
 std.Info(args...)
}
```

* 为日志库接口方法添加注释

```go
type FieldLogger interface {
 WithField(key string, value interface{}) *Entry
 WithFields(fields Fields) *Entry
 WithError(err error) *Entry

 Debugf(format string, args ...interface{})//@codeql logger
 Infof(format string, args ...interface{})//@codeql logger
 Printf(format string, args ...interface{})//@codeql logger
 Warnf(format string, args ...interface{})//@codeql logger
 Warningf(format string, args ...interface{})//@codeql logger
 Errorf(format string, args ...interface{})//@codeql logger
 Fatalf(format string, args ...interface{})//@codeql logger
 Panicf(format string, args ...interface{})//@codeql logger

 Debug(args ...interface{})//@codeql logger
 Info(args ...interface{})//@codeql logger
 Print(args ...interface{})//@codeql logger
 Warn(args ...interface{})//@codeql logger
 Warning(args ...interface{})//@codeql logger
 Error(args ...interface{})//@codeql logger
 Fatal(args ...interface{})//@codeql logger
 Panic(args ...interface{})//@codeql logger

 Debugln(args ...interface{})//@codeql logger
 Infoln(args ...interface{})//@codeql logger
 Println(args ...interface{})//@codeql logger
 Warnln(args ...interface{})//@codeql logger
 Warningln(args ...interface{})//@codeql logger
 Errorln(args ...interface{})//@codeql logger
 Fatalln(args ...interface{})//@codeql logger
 Panicln(args ...interface{})//@codeql logger
}
```

## HTTP::Redirect

标记HTTP重定向URL，用于发现重定向注入攻击

### 注释格式

//@codeql redirect Param(0)

### 例子

* 在方法/函数/接口方法上标记HTTP redirect方法URL参数位置
* 例子中Redirect方法的第一个参数就是URL，所以是Param(0)

```go
//@codeql redirect Param(0)
func (c *Ctx) Redirect(location string, status ...int) error {
 c.setCanonical(HeaderLocation, location)
 if len(status) > 0 {
  c.Status(status[0])
 } else {
  c.Status(StatusFound)
 }
 return nil
}
```
