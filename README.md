# README #

这个项目是基于 Echo Web 框架，这个框架成熟且文档完善。在这个框架的基础上，对应用配置、日志、中间件、自定义错误处理、自定义验证器、路由注册、数据库初始化、优雅关机等做了封装和加强。

### Lian Shuo ###

调用中台的示例如下:
```go
client := lian_shuo.NewClient("eyJhbGciOiJIUzUxMiJ9.eyJzdW")
request := requests.NewGetUserInfoRequest()
err := client.Request(request)
```