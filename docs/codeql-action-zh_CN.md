强化GO开源项目安全检测&开源项目漏洞挖掘方法
===

## 现状

提到Go语言代码安全检测的现状，我们先分析一下现有开源方案的不足和局限，一些代码规范检查linter工具不在此次讨论范围。`gosec`和`gokart`这两款代码静态分析工具，都属于过程内分析，对函数调用的参数做最保守假设，过度假设导致丢失精度，所以误报很多。基于污点分析技术的CodeQL就强大多了，但要掌握CodeQL达到编写ql查询模块和qll库模块学习成本颇高。目前绝大部分Go的开源项目会在GitHub Actions中配置CodeQL检测，[配置](https://docs.microsoft.com/zh-cn/dotnet/architecture/devops-for-aspnet-developers/actions-codeql)过程超级简单，这里最大的问题就是CodeQL内置库模块支持的[框架](https://codeql.github.com/docs/codeql-overview/supported-languages-and-frameworks/)有局限，针对千变万化的代码，仅仅使用默认的CodeQL规则，检测能力大打折扣，甚至是形如虚设。

## codemillx设计

* 问题1: 如何为项目编写特定的CodeQL查询模块(ql)和库模块(qll)，提升检测能力
* 问题2: 在GitHub Actions中配置的CodeQL检测，如何能使用我们编写的qll模块

为了解决这2个问题，我们制造了`codemillx`这款辅助工具。

* 在项目代码中添加一些特定注释[(注释格式说明)](/docs/comment-zh_CN.md)，并可生成适配项目的qll库模块，比起掌握QL语法简单很多
* 运行参数-customizeCodeQLAction=true，就会把自定义规则写进CodeQL内置Customizations.qll文件

## 强化GO开源项目安全检测

只需两步：

1. 在项目代码中添加一些特定注释[(注释格式说明)](/docs/comment-zh_CN.md)
2. 修改Github Actions的CodeQL检测配置文件(.github/workflows/codeql.yml)

把以下配置添加在Initialize CodeQL步骤与Autobuild步骤之间。

```yaml
# ...
- name: Generate And Replace CodeQL Customizations
  run: go run github.com/hudangwei/codemillx/cmd/codemillx -customizeCodeQLAction=true ./...
# ...
```

#### 完整配置文件(.github/workflows/codeql.yml)

```yaml
# For most projects, this workflow file will not need changing; you simply need
# to commit it to your repository.
#
# You may wish to alter this file to override the set of languages analyzed,
# or to provide custom queries or build logic.
#
# ******** NOTE ********
# We have attempted to detect the languages in your repository. Please check
# the `language` matrix defined below to confirm you have the correct set of
# supported CodeQL languages.
#
name: "CodeQL"

on:
  push:
    branches: [ master ]
  pull_request:
    # The branches below must be a subset of the branches above
    branches: [ master ]

jobs:
  analyze:
    name: Analyze
    runs-on: ubuntu-latest
    permissions:
      actions: read
      contents: read
      security-events: write

    strategy:
      fail-fast: false
      matrix:
        language: [ 'go' ]
        # CodeQL supports [ 'cpp', 'csharp', 'go', 'java', 'javascript', 'python', 'ruby' ]
        # Learn more about CodeQL language support at https://git.io/codeql-language-support

    steps:
    - name: Checkout repository
      uses: actions/checkout@v2

    # Initializes the CodeQL tools for scanning.
    - name: Initialize CodeQL
      uses: github/codeql-action/init@v1
      with:
        languages: ${{ matrix.language }}
        # If you wish to specify custom queries, you can do so here or in a config file.
        # By default, queries listed here will override any specified in a config file.
        # Prefix the list here with "+" to use these queries and those in the config file.
        # queries: ./path/to/local/query, your-org/your-repo/queries@main
    - name: Generate And Replace CodeQL Customizations
      run: go run github.com/hudangwei/codemillx/cmd/codemillx -customizeCodeQLAction=true ./...
        
    # Autobuild attempts to build any compiled languages  (C/C++, C#, or Java).
    # If this step fails, then you should remove it and run the build manually (see below)
    - name: Autobuild
      uses: github/codeql-action/autobuild@v1

    # ℹ️ Command-line programs to run using the OS shell.
    # 📚 https://git.io/JvXDl

    # ✏️ If the Autobuild fails above, remove it and uncomment the following three lines
    #    and modify them (or add more) to build your code if your project
    #    uses a compiled language

    #- run: |
    #   make bootstrap
    #   make release

    - name: Perform CodeQL Analysis
      uses: github/codeql-action/analyze@v1

```

## 开源项目漏洞挖掘方法

如果你是一名白帽子，介绍使用`codemillx`去挖掘Go开源项目漏洞，应该是你感兴趣的。

### 步骤

* fork开源项目到自己仓库下
* 给开源项目添加注释(自定义标记污点源)
* 添加CodeQL配置文件.github/workflows/codeql.yml
* 提交代码，触发检测

如果你看不懂本页内容，先学习下[如何在GitHub工作流中使用CodeQL检测代码](https://docs.microsoft.com/zh-cn/dotnet/architecture/devops-for-aspnet-developers/actions-codeql)
