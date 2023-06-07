# 介绍
Demo目录中是建议的项目目录结构，tools是代码工具，sample是一个实际的项目例子

# 项目框架
Demo目录中是建议的项目目录结构 

+-项目目录   
   +**lib**  第三方包
   |-+ gorm  mysql
   |-+ router 路由库
   |-+ avl 平衡二叉树
   |-+ OA OA服务：TOF，GAS，数据中转站
   |-+ workflow 精简工作流
   |-+ util 工具类
   |
   +**model** 数据模型
   |-+db 与数据库交互的模型
   |-+view 与前端交互的模型
   |
   +**router** 路由设置
   |-+middleware 路由中间件     
   |-- register.go 自动生成路由表
   |-- facade.go 自动生成接口实现
   |
   +**commcon** 项目中的通用处理代码
   |--const.go 常量，readonly变量
   |
   +**service** 服务类代码
   |
   +**start** 启动代码
   |--main.go 
   |--config.json 配置


目录之间的引用关系
lib目录下的包不可引用项目中其他目录，lib之间的包可以引用
router下的可以引用其他包
model包只引用lib包，只做类型定义
service包可以引用除 router和start 外的包

# 代码生成规则
代码生成的目的是将实际的Service类中的方法暴露为Http的接口，对需要暴露为Http方法的函数增加对应的注释即可自动生成接口代码

通过代码中注释进行代码生成，目前识别的命名如下：
1.@MethodGroup:在整个文件的package之前可以标注 @MethodGroup，此文件下所有方法都将归于此Group下。
2.@Method: Http方法，在函数前的注释中标识。
3.@MethodName: 路由名称，可以省略，默认为方法名，也可以配置成为动态路由/MyMethodName/:param1/:param2/MethodName,其中 param1,param2必须与方法签名中的参数名对应。在函数前的注释中标识。

在方法中如果要对Response、Request或者当前调用用户信息进行调用，可以在函数参数中增加 Context 参数，Context定义为
```
type User struct {
	UserID        string
	UserName      string
}

type Context struct{
	Res http.ResponseWriter
	Req *http.Request
	User
}
```

# 代码生成工具的使用
将tools中的CodeGen.go，编译成CodeGen.exe，
Usage of CodeGen:
  -d string
        要扫描的路径,
  -o string
         文件输出路径
  -r string
        当前路径在项目路径或go path中的位置
  -t string
        文件模版路径，此代码生成是基于模版的，目前生成register.go 和 facade.go 所以需要两个文档模版
		
在 sample目录下运行如下命令（此项目被go mod init为rice，所以下面的项目路径名称是"rice"）
```
 ..\tools\CodeGen -r "rice/" -d . -t  ..\tools -o router
```
 会自动在router目录下生成register.go  和 facade.go ,然后
```
 cd start 
 go run main.go
```
 就可以运行起项目来了。
 
 
# Router的功能简介
 支持如下的动态的路由定义
 r.AddRoute("/api/User/:id",HttpHandlerFunc1)
 r.AddRoute("/api/User/:id/Create",HttpHandlerFunc2)
 r.AddRoute("/api/User/:id/:name",HttpHandlerFunc3)
 r.AddRoute("/api/User/List/:name",HttpHandlerFunc4)
 不支持
 r.AddRoute("/api/User/:id/Create/:name",HttpHandlerFunc5)
 
 对动态路由的匹配，最长匹配优先，具名匹配优先，比如对 请求 /api/User/List/T3 会匹配到HttpHandlerFunc4 而不是HttpHandlerFunc3，/api/User/23/Create 会匹配到HttpHandlerFunc2 而不是HttpHandlerFunc3
 
# 后续计划
 1.自动生成单元测试
 
