
# 自定义代码污点，生成codeql ql模块
## 使用说明
| | codeql标记 | 注释格式 | 内置方法 | 
| :----:| :----: | :----: | :----: |
| 标记污点源| UntrustedFlowSource | //@codeql untrust | Result() Param() IsReceiver |
| 标记污点传播函数| TaintTracking | //@codeql tainttrack | InParam() OutParam() OutResult() InIsReceiver OutIsReceiver |
| 标记SQL语句| SQL::QueryString | //@codeql sql | Param() |

## 目前支持的codeql model
* UntrustedFlowSource
* TaintTracking
* SQL::QueryString

## 注释格式
- //@codeql untrust Result(1)
- //@codeql tainttrack InParam(0) OutResult(0)
- //@codeql sql Param(0)

## 内置方法

- Result 标记返回值是污点源 在//@codeql untrust使用
- Param 
* 在//@codeql untrust使用,标记参数是污点源 
* 在//@codeql sql使用,标记参数是sql语句
- IsReceiver 标记receiver是污点源 在//@codeql untrust使用
- InParam 只能在//@codeql tainttrack使用
- OutParam 只能在//@codeql tainttrack使用
- OutResult 只能在//@codeql tainttrack使用
- InIsReceiver 只能在//@codeql tainttrack使用
- OutIsReceiver 只能在//@codeql tainttrack使用

## 参考项目 [codemill](https://github.com/gagliardetto/codemill)
