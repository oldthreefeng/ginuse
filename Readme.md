## ginuse

autoDeploy your project. 

### Webhooks干嘛用的？

> Github Webhooks提供了一堆事件，这些事件在用户特定的操作下会被触发，比如创建分支(Branch)、库被fork、项目被star、用户push了代码等等。
我们可以自己写一个服务,将服务的URL交给Webhooks，当上述事件被触发时，Webhook会向这个服务发送一个POST请求，请求中附带着该事件相关的详细描述信息(即Payload)。
这样，我们就可以在自己服务中知道Github的什么事件被触发了，事件的内容是什么？据此我们就可以干一些自己想干的事了。能干什么呢？官方说You're only limited by your imagination，就是说想干什么都行，就看你的想像力够不够 :)

> 当指定的事件发生时，我们将向您提供的每个URL发送POST请求。通过这个post请求，我们就能实现自动拉取仓库中的代码，更新到本地，最终实现自动化更新



### webHook Post

![](images/XHSign.jpg)

- Request URL

即前面配置中填写的"Payload URL"

- content-type

即前面配置中选择的"Content type"

- X-Hub-Signature

是对Payload计算得出的签名。当我们在前面的配置中输入了"Secret"后，Header中才会出现此项。官方文档对Secret作了详细说明，后面我们也会在代码中实现对它的校验

```cgo
1. 监听端口port:8000,监听的uri路径path,运行部署脚本sh,webhook的secret,
2. 使用-p 指定端口,使用-path 指定uri路径,使用-sh 指定运行脚本, 使用-s 指定密码,
3. 原理是通过webhook的Post,来校验sha1,通过校验则执行部署脚本
 
```

- 验证

![verify](./images/respone.jpg)

练手的go项目,学习golang15天~