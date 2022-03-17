codemillx
===

English | [简体中文](README.md)

`codemillx` is a tool for `CodeQL`, extract the comments in the code and generate codeql module.

## Installation

To install `codemillx` use the `go get` command:

```bash
go get github.com/hudangwei/codemillx/cmd/codemillx
```

## Run

Navigate to your web application folder and execute:

```bash
cd mywebapp && codemillx ./...
```

## Requirements

To allow cqlgen to format the generated codeql, you need a recent version of the codeql cli (otherwise it will not be formatted), and have it available as codeql in your PATH.

```sh
codeql query format -qq -i Customizations.qll
```

## How to add comments in your code?

* [Declarative Comments Format](docs/comment.md)

## Usage

* [How to use `Customizations.qll` file into Github CodeQL Action](docs/codeql-action.md)

## Reference

* [codemill](https://github.com/gagliardetto/codemill)
