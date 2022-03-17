codemillx
===

[English](README.md) | 简体中文

codemillx是一款CodeQL辅助工具，通过提取代码中的注释，并可生成codeql qll库模块。

## 安装

通过 `go get` 命令安装:

```bash
go get github.com/hudangwei/codemillx/cmd/codemillx
```

## 运行

在你的项目根目录下执行命令:

```bash
cd mywebapp && codemillx ./...
```

运行命令后将生成`Customizations.qll`文件

## 依赖项

`codemillx`生成codeql模块时会调用`codeql`命令进行格式化，所以依赖本地PATH中有`codeql`。

```sh
codeql query format -qq -i Customizations.qll
```

## 如何添加代码注释

* [注释格式说明](docs/comment-zh_CN.md)

## 实践案例

* [强化GO开源项目安全检测(内含开源项目漏洞挖掘方法)](docs/codeql-action-zh_CN.md)

## 参考项目

* [codemill](https://github.com/gagliardetto/codemill)
