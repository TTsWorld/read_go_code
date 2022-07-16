# read_go_code


这个仓库用于翻译和注释 go 语言的一些源码，主要包含如下源码
- go1.18
  - io
  - context
  - http
- gin


# 关于阅读顺序和阅读包的选择
因为项目中使用到了 gin 框架，最开始只是想阅读 gin 框架的。但是发现 gin 框架 使用到了 http 包，发现如果不理解 http 包的话，是没有办法完整的理解 gin框架的，所以又开始阅读 http 包。
阅读http 包的过程中，再又发现底层调用了很多基础库，比如 io，context 等，如果对于基础库不够了解，阅读过程中，还是有太多黑盒。所以索性从基础库的最小粒度开始逐层阅读。当然最终的目标还是通过阅读基础库，解决阅读 gin 源码的障碍。
可以说通过阅读gin 源码，扣响阅读 go 源码的大门吧。

后续的阅读计划是，将阅读源码本身当成一个项目，在阅读过程中翻译+注释。比如阅读 gin 项目作为一个项目，但是 gin 项目依赖其它的 go 基础库，那就顺着 gin 的引用线将使用到的一些基础库也阅读完。
等 gin 项目阅读完后，大部分基础库应该也阅读的差不多了。阅读完 gin 后，可以选择其他的 go 项目继续进行阅读，从而不断的填补在 go 语言源码方面的知识空白。


# 关于阅读版本的选择
参考了[Go 各版本特性](https://github.com/guyan0319/golang_development_notes/blob/master/zh/1.6.md)  
可以选择阅读最新版源码

# 阅读源码的工具
- Goland
-  [一个可以生成Go源码UML结构的工具](https://www.dumels.com/)


# 关于源码阅读环境搭建
[搭建 Go 源码阅读环境](https://juejin.cn/post/6859225567028477966)

# 阅读过程中参考的文章
## 专栏
[知乎-Golang源码解析](https://www.zhihu.com/column/c_1305921732017573888)  
[GO夜读-如何高效阅读 Go 代码](https://www.bilibili.com/video/BV1XD4y1U7Pf?spm_id_from=333.999.0.0)  
[知乎-如何阅读 Golang 的源码](https://www.zhihu.com/question/327615791/answer/756625130)
[Github - Go源码分析](https://github.com/jianfengye/inside-go)
## HTTP
[Go语言 HTTP 标准库的视线原理](https://draveness.me/golang/docs/part4-advanced/ch09-stdlib/golang-net-http/)
## Gin
[详细讲解Go web框架之gin框架源码解析记录及思路流程和理解](https://blog.csdn.net/pythonstrat/article/details/121423122)
## Go版本选择
[https://juejin.cn/post/7018778089190588423](https://juejin.cn/post/7018778089190588423)