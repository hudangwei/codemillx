codemillx
===

English | [简体中文](README-zh_CN.md)

`codemillx` is a tool for `CodeQL`, extract the comments in the code and generate codeql module.

## Installation

To install `codemillx` use the `go get` command:

```bash
go get github.com/hudangwei/codemillx/cmd/codemillx
```

## Run

Navigate to your web application folder and execute:

```bash
cd mywebapp && codemillx -module="mywebapp" ./...
```

## Reference

* [codemill](https://github.com/gagliardetto/codemill)
