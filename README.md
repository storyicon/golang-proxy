# golang-proxy `v2.0`

![golang-proxy](https://img.shields.io/teamcity/codebetter/bt428.svg)
![download](https://img.shields.io/eclipse-marketplace/dt/notepad4e.svg)

Golang-Proxy -- 简单高效的免费代理抓取工具通过抓取网络上公开的免费代理，来维护一个属于自己的高匿代理池，用于网络爬虫、资源下载等用途。
![golang-proxy](https://raw.githubusercontent.com/parnurzeal/gorequest/gh-pages/images/Gopher_GoRequest_400x300.jpg)

## What's new in V2.0?

1.  **不再依赖 MySQL 和 NSQ**！
2.  之前需要分别启动`publisher`、`consumer`和`assessor`，现在 **只需要启动主程序** 即可！
3.  提供了高度灵活的 **API 接口**，在启动主程序后，即可通过在浏览器访问`localhost:9999/all` 与 `localhost:9999/random` 直接获取抓到的代理！甚至可以使用 `localhost:9999/sql?query=`来执行 SQL 语句来自定义代理筛选规则！
4.  提供 `Windows`、`Linux`、`Mac` **开箱即用版**！
    [Download Release v2.0](https://github.com/storyicon/golang-proxy/releases/)

## 安装

#### 1. 通过编译源码

```bash
go get github.com/storyicon/golang-proxy
```

进入到 `golang-proxy` 目录，执行 `go build main.go`，执行生成的二进制的执行程序即可。

**注意：**
在 `go build` 的过程中可能出现`cannot find package "github.com/gocolly/col1ly" in any of` 等找不到包的情况，根据提示的地址 `go get` 即可

```
# 比如如果在 go build main.go 的时候提示
business\publisher.go:8:2: cannot find package "github.com/gocolly/col1ly" in any of:
        F:\Go\src\github.com\gocolly\col1ly (from $GOROOT)
        D:\golang\src\github.com\gocolly\col1ly (from $GOPATH)
        C:\Users\Administrator\go\src\github.com\gocolly\col1ly
        D:\ivank\src\github.com\gocolly\col1ly
执行 go get github.com/gocolly/col1ly 即可
```

如果觉得麻烦，可以使用 `/bin` 目录中提供的 **`开箱即用`** 版本。

#### 2. 开箱即用版本

[Release 页面](https://github.com/storyicon/golang-proxy/releases/)根据系统环境提供了一些压缩包，将他们解压后执行即可。

开箱即用版下载地址: [Download Release v2.0](https://github.com/storyicon/golang-proxy/releases/)

#### 3. Tips

项目根目录下的 `./source` 是项目执行必须的文件夹，里面存储了各类网站源，其他的文件夹储存的均为项目源码。所以在编译后得到二进制程序 `main` 文件后，即可将 `main` 文件和 `source` 文件夹一同移动到任意地方，`main` 文件可以任意命名。  
如果提示找不到 `source`文件夹， 你可以在执行程序时加上`-source=`参数来指定`source`文件夹路径，例如：

```bash
# xxx填source文件夹的相对或者绝对路径
main -source=xxx
```

## API 接口

在程序运行后，可以通过在浏览器访问以下接口获取数据库中抓取到的代理。

#### 1. 随机获取一条代理

```json
地址: http://localhost:9999/random
返回示例：
{
    //状态码0表示成功，1表示错误
	"code": 0,
	"message": [{
		"id": 124,
		"content": "http://190.2.144.133:1080",
		//评估次数，次数越多，代表代理存活时间越长
		"assess_times": 13,
		//评估成功次数，success_times/assess_times可以得到评估成功率
		"success_times": 11,
		//平均响应时间，单位为秒
		"avg_response_time": 2.0831538461538464,
		//连续评估失败次数，是分数计算的重要指标
		"continuous_failed_times": 0,
		//分数，分数越高，代理质量越高
		"score": 3.2747991296955083,
		//插入时间戳（秒）
		"insert_time": 1532324791,
		//更新时间戳（秒）
		"update_time": 1532414960
	}]
}
```

#### 2. 获取所有可用代理

```json
地址: http://localhost:9999/all
```

#### 3. 执行 SQL

```json
地址: http://localhost:9999/sql/query=xxxx
将xxxx替换为要执行的sql语句即可，程序配置了两张表：
valid_proxy 存放高可用代理
crawl_proxy 抓取到的代理的缓存表（代理质量不能保证）

例如: http://localhost:9999/sql/query=SELECT * FROM VALID_PROXY WHERE 1 ORDER BY SCORE DESC
将会将所有的可用代理按照分数倒序并返回。
```

## 为什么要用 Golang-Proxy

1.  稳定、快速。  
    抓取模块，**单核并发可以到达 1000 个页面/秒**。
2.  高可配置性、高拓展性。  
    你不需要写任何代码，花**一两分钟**填写一个配置文件就可以添加一个新的网站源。
3.  评估功能。  
    通过 Assessor 评估模块，周期性测试代理质量，根据代理的**测试成功率、高匿性、测试次数、突变性、响应速度**等独立影响因子进行综合评分，算法具有高度可配置性，可以根据项目的需要可以对因子的权重进行独立调整。

## 如何配置一个新的源

`./source/`下的所有 yml 格式的文件都是**源**，你可以增加源，也可以通过在文件名前加上一个 **`.`** 来使程序忽略这个源，当然你也可以直接删除，来让一个源永远的消失，下面进行 Source 参数介绍：

```yml
#Page配置项
page:
    entry: "https://xxx/1.html"
    template: "https://xxx/{page}.html"
    from: 2
    to: 10
#publisher将会首先抓取entry，即 https://xxx/1.html
#然后根据 template、from 和 to 依次抓取
#　　https://xxx/2.html
#　　https://xxx/3.html
#　　https://xxx/4.html
#　　...
#　　https://xxx/10.html
```

```yml
#Selector配置项
selector:
    iterator: ".table tbody tr"
    ip: "td:nth-child(1)"
    port: "td:nth-child(2)"
    scheme: "td:nth-child(3)"
    filter: ""
# 以上配置用于抓取下面这种 HTML 结构
# <table class="table">
#     <tbody>
#         <tr>
#             <td>187.3.0.1</td>
#             <td>8080</td>
#             <td>HTTP</td>
#         <tr>
#         <tr>
#             <td>164.23.1.2</td>
#             <td>80</td>
#             <td>HTTPS</td>
#         <tr>
#         <tr>
#             <td>131.9.2.3</td>
#             <td>8080</td>
#             <td>HTTP</td>
#         <tr>
#     <tbody>
# <table>
# 选择器为通用的JQuery选择器，iterator为循环对象，比如表格里的行，每行一条代理，那这个行的选择器就是iterator，而ip、port、protocal则是在iterator选择器的基础上进行子元素的查找。
# protocal为空，或protocal对应的元素无法找到，则默认是HTTP类型
```

```yml
category:
    # 并行数
    parallelnumber: 1
    # 对于这个源，每抓取一个页面
    # 将会随机等待5~20s再抓下一个页面
    delayRange: [5, 20]
    # 间隔多长时间启用一次这个源
    # @every 10s ， @every 10h...
    interval: "@every 10m"
debug: true
```

## Request For Comments

1.  使用中任何问题提 `issues` 即可
2.  如果发现了新的好用的源，欢迎提交上来分享
3.  来都来了点个 Star 再走呗 : )
