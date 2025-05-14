# README #

这个项目是基于 Echo Web 框架，这个框架成熟且文档完善。在这个框架的基础上，对应用配置、日志、中间件、自定义错误处理、自定义验证器、路由注册、数据库初始化、优雅关机等做了封装和加强。

### 依赖注入 ###

依赖注入框架使用的是 Google 开元的 Wire，相比其他的框架它采用的是生成代码的方式，有别于反射的方式，所以并无运行时开销。

使用方法：在 `cmd/api/wire.go` 文件中，声明依赖关系后，执行如下命令：

```shell
$ cd cmd/api && $(go env GOPATH)/bin/wire
```

这条命令会生成 `wire_gen.go` 文件，然后在 main 方法中使用其中的 `InitializeApp` 方法来初始化应用的依赖关系。

### Lian Shuo ###

调用中台的示例如下:
```go
client := lian_shuo.NewClient("eyJhbGciOiJIUzUxMiJ9.eyJzdW")
request := requests.NewGetUserInfoRequest()
err := client.Request(request)
```