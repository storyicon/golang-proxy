Golang-Proxy
=========

Golang-Proxy -- 简单高效的免费代理抓取工具
通过抓取网络上公开的免费代理，来维护一个属于自己的高匿代理池，用于网络爬虫、资源下载等用途。
![GopherGoRequest](https://raw.githubusercontent.com/parnurzeal/gorequest/gh-pages/images/Gopher_GoRequest_400x300.jpg)

## 安装
#### 1. 通过编译源码  
```bash
$ git clone https://github.com/storyicon/golang-proxy.git
$ cd golang-proxy/publisher
$ go build publisher\publisher.go
$ go build consumer\consumer.go
$ go build assessor\assessor.go
```
分别启动编译好的publisher、consumer、assessor即可

#### 2. 开箱即用版本  
项目的 __./source/__目录下已经默认配置好了一些源，你只需要在__./config/local.yml__中正确配置你的数据库信息，就可以进行代理的抓取。

1. **对于windows用户**    
将./bin/windows目录下的publisher.exe, consumer.exe, assessor.exe 移动到golang-proxy的根目录，分别运行publisher.exe, consumer.exe, assessor.exe即可
2. **对于lunix用户**    
...
3. **对于Mac用户**    
...

__注意：__
想要让golang-proxy在你的机器上成功运行，你还需要一个MySQL和一个NSQ。     
NSQ下载地址：https://nsq.io/deployment/installing.html    
MySQL下载地址：https://www.mysql.com/downloads/    
安装完成后启动NSQ和MySQL服务    
在配置文件 **./config/local.yml** 中修改你的配置信息，在MySQL中相应的数据库（local.yml中的mysql.db）下执行如下SQL语句创建数据表：    
```sql
CREATE TABLE `valid_proxy` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `content` varchar(30) CHARACTER SET utf8mb4 NOT NULL DEFAULT '',
  `assess_times` int(11) NOT NULL DEFAULT '0',
  `last_assess_time` int(11) NOT NULL DEFAULT '0',
  `success_times` int(11) NOT NULL DEFAULT '0',
  `avg_response_time` float(5,3) NOT NULL DEFAULT '0.000',
  `continuous_failed_times` int(11) NOT NULL DEFAULT '0',
  `score` float(8,3) NOT NULL DEFAULT '0.000',
  PRIMARY KEY (`id`),
  UNIQUE KEY `content` (`content`)
) ENGINE=InnoDB AUTO_INCREMENT=10428 DEFAULT CHARSET=utf8;

```
在这些都全部完成后，你才可以正常的运行Golang-Proxy    
## 为什么要用 Golang-Proxy    
1. 稳定、快速。    
抓取模块，__单核并发可以到达1000个页面/秒__。    
2. 高可配置性、高拓展性
你不需要写任何代码，花__一两分钟__填写一个配置文件就可以添加一个新的网站源。
通过评估模块，周期性测试代理质量，根据代理的__抓取次数、成功率、响应速度、突变性__来进行综合评估，参数可以根据你的需要进行调节。
## 特色    
1. 高拓展性    
　　引入了 **源** 的概念，可以把一个源理解为一个待抓取的网站，有很多这样的网站，比如“kuaidaili.com”，“xiaohexia.cn”，他们都是源。优质的源越多，你可以抓取到的免费代理就越多，抓取效率也就越高。    
　　使用Golang-Proxy，添加一个新的网站源，不需要写任何代码，也不需要重新编译程序，一两分钟你就可以完成一个源的添加。
2. 模块化    
Golang-Proxy将代理的收集流程分为了三个模块:    
　　* __Publisher__ 用于抓取代理；    
　　* __Consumer__ 负责从队列中读取抓到的代理，并进行有效性检测，有效则插入数据库；    
　　* __Assessor__ 用于对数据库中的代理进行价值评估；     
　　为了提高某个模块的效率，你可以启动多个进程，比如你认为抓取速度太慢了，你甚至可以再启动一个Publisher。     
3. 简单配置     
　　./config目录下，你可以快捷地对MySQL和NSQ的参数进行配置。     
项目分为测试环境和线上环境，你可以通过设置环境变量**GOLANG_PROXY_ENV**为**local**或**prod**，来控制程序是使用local.yml还是prod.yml作为配置文件，**默认是读取local.yml**。当然，如果你只是在本地运行，那只需要修改local.yml的配置即可。    
　　你可以简单的通过在 __./source__ 目录下添加yml格式的源，Golang-Proxy的publisher在启动时，将自动读取并载入。    

## 配置项
### 1. Config配置项
　　**./config/** 下的**local.yml**和**prod.yml**分别用于配置**本机测试**和**线上运行**的数据库参数，运行前需要先行配置。
### 2. Source配置项：
　　**./source/**下的所有yml格式的文件都是**源**，你可以通过新建来增加一个源，也可以通过在文件名前加上一个 **.** 来告诉 **publisher** 不要读取这个源，当然你也可以直接删除，来让一个源永远的消失，下面进行Source参数介绍：
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
    protocal: "td:nth-child(3)"
    proxytype: ""
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
# 选择器为通用的JQuery选择器，iterator为循环对象，比如表格里的行，每行一条代理，那这个行的选择器就是iterator，而ip、port、protocal则是在iterator选择器的基础上进行子元素的查找。至于proxytype和filter忽略即可，这两个参数是为将来的功能拓展预留的项目。
#protocal为空，或protocal对应的元素无法找到，则默认是HTTP类型
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
### 3. 策略配置项：
你可以在 **./library/const.go** 中进行配置
```go
const (
	NSQTopic                       = "proxy"
	# 在Publisher每抓到一个代理，将会发送给Consumer进行PreAssess预检测，预检测程序将使用这个代理会访问httpbin.org，通过的才会插入到数据库，ProxyPreAssessTimes定义了访问时允许的trytimes次数
	ProxyPreAssessTimes            = 1
	# 定义了Consumer和Assessor在进行代理的测试时，能忍受的延迟
	ProxyAssessTimeOut             = 3
	# 定义了Assessor每隔多少秒对数据里的代理进行一次评估
	ProxyAssessInterval            = 60
	# 定义了当一条代理通过了Consumer的预检测，在插入数据库时给予的初始分数
	ProxyInitScore                 = 1
	# 定义了代理评估程序Assessor内存中允许暂存的代理数量
	ProxyAssessQueueMin            = 1000
	# 定义了能够忍受的代理测试的最小成功率，如果这个值设置为0.5，那么对于数据库里的每条代理，从概率上可以保证每使用两次，有一次可以成功。如果是0.8，则保证数据库中的每条代理都有80%的使用成功率
	AllowProxyAssessSuccessRateMin = 0.5
)
```

