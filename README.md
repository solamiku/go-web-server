# go-web-server
## 基础网络脚手架

该脚手架仅包含golang web网站构建的基础部分
由此可以扩展任何web服务
本脚手架使用fasthttp框架


### 以下是需要包含的第三方库

* https://github.com/go-sql-driver/mysql mysql
* https://github.com/cihub/seelog 日志库
* https://github.com/valyala/fasthttp fasthttp
	- https://github.com/klauspost/compress fasthttp需要的compress库
		- https://github.com/klauspost/cpuid compress需要的cpuid库
	- https://github.com/valyala/bytebufferpool fasthttp需要的bytebufferpool库

### 目录结构
* webserver
    - config        go-配置相关
    - db            go-数据库操作相关
    - log           日志存放-运行后自动生成
    - public        前端可访问的文件
        + css       css文件
        + js        js文件
        + img       图片
        + ...       其余第三方库等
    - router        go-路由功能相关
    - view          html模板