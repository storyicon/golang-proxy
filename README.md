# golang-proxy `v3.1`

![golang-proxy](https://img.shields.io/teamcity/codebetter/bt428.svg)
[![download](https://img.shields.io/eclipse-marketplace/dt/notepad4e.svg)](https://github.com/storyicon/golang-proxy)

> Made some changes to the judgment logic of whether the proxy is highly anonymous. This problem has led to the absence of proxies that can be used. Now you can enjoy it. 😏

-   [English Document](#english-document)
    -   [1. Feature](#1-feature)
    -   [2. How to use](#2-how-to-use)
        -   [API interface](#api-interface)
    -   [3. Advanced](#3-advanced)
        -   [two `data tables`](#two-data-tables)
            -   [1. Table Crude Proxy](#1-table-crude-proxy)
            -   [2. Table Proxy](#2-table-proxy)
        -   [one `configuration file`](#one-configuration-file)
        -   [one `source folder`](#one-source-folder)
        -   [four `modules`](#four-modules)
    -   [Request for comments](#request-for-comments)
-   [中文文档](#中文文档)
    -   [在 `v3.0` 有哪些新特性](#在-v30-有哪些新特性)
    -   [如何使用 `golang-proxy`](#如何使用-golang-proxy)
        -   [1. 使用开箱即用版本](#1-使用开箱即用版本)
            -   [接口示例: `localhost:9999/sql`](#接口示例-localhost9999sql)
        -   [2. 使用源码编译](#2-使用源码编译)
    -   [为什么要用 Golang-Proxy](#为什么要用-golang-proxy)
    -   [如何配置一个新的源](#如何配置一个新的源)
    -   [征求意见](#征求意见)

![golang-proxy](https://raw.githubusercontent.com/parnurzeal/gorequest/gh-pages/images/Gopher_GoRequest_400x300.jpg)

# English Document

Golang-proxy is an efficient free proxy crawler that ensures that the captured proxies are highly anonymous and at the same time guarantee their quality. You can use these captured proxies to download network resources and ensure the privacy of your own identity.

## 1. Feature

-   Very high speed of proxy crawler, which can download 1000 pages per second.
-   You can customize the source of proxy crawler. The configuration file is extremely simple.
-   Provide a compiled version, comes with a SQLite database, and supports mysql
-   Comes with an API interface, all functions can be used with one click
-   Proxy evaluation system to ensure the quality of the proxy pool

## 2. How to use

`golang-proxy` provides compiled binary files so that you do not need `golang` on the machine. Download binary compression pack to [Release Page](https://github.com/storyicon/golang-proxy/releases/)  
According to your system type, download the corresponding compression package, unzip it and run it. After a few minutes, you can access `localhost:9999/all` in the browser to see the proxy's crawl results.

Before I go into the detailed introduction of golang-proxy, I think it's best to tell you the most useful information first.

### API interface

After you start the binary, you can access the following interface in the browser to get the proxy

| url                                    | description                                                                               |
| -------------------------------------- | ----------------------------------------------------------------------------------------- |
| `localhost:9999/all`                   | Get all highly available proxies                                                          |
| `localhost:9999/all?table=proxy`       | Get all highly available proxies                                                          |
| `localhost:9999/random`                | Randomly acquire a highly available proxy                                                 |
| `localhost:9999/all?table=crude_proxy` | Obtain the proxies in the temporary table (the quality of them cannot be guaranteed)      |
| `localhost:9999/random?table=proxy`    | Randomly get an proxy from the temporary table (the quality of them cannot be guaranteed) |
| `localhost:9999/sql?query=`            | Write the SQL statement you want to execute after `query=`, customize your filter rules.  |

Having mastered the above content, you have been able to use the 50% function of `golang-proxy`. But the last interface allows you to execute custom SQL statements, and you'll find that you need to know at least the structure of the tables. The following will tell you.

## 3. Advanced

golang-proxy consists of the following parts:

-   two `data tables`
-   one `configuration file`
-   one `source folder`
-   four `modules`

### two `data tables`

#### 1. Table Crude Proxy

In order to store temporary proxies, we designed the data table `crude_proxy`, the table is defined as follows.

| field       | type   | example         | description |
| ----------- | ------ | --------------- | ----------- |
| id          | int    | -               | -           |
| ip          | string | 192.168.0.1     | -           |
| port        | string | 255             | -           |
| content     | string | 192.168.0.1:255 | -           |
| insert_time | int    | 1540798717      | -           |
| update_time | int    | 1540798717      | -           |

table `crude_proxy` stores the proxies that are crawled out, and cannot guarantee their quality.

#### 2. Table Proxy

When the agent in the `crude_proxy` table passes through `pre assess` ( `pre assess` roughly verifies the availability of the proxy and tests the proxy's support for `https` and `http` ), it will enter the `proxy` table.

| field                   | type   | example         | description                                                                                                      |
| ----------------------- | ------ | --------------- | ---------------------------------------------------------------------------------------------------------------- |
| id                      | int    | -               | -                                                                                                                |
| ip                      | string | 192.168.0.1     | -                                                                                                                |
| port                    | string | 255             | -                                                                                                                |
| scheme_type             | int    | 2               | Identify the extent to which the proxy supports http and https, `0`: http only, `1` https only, `2` https & http |
| content                 | string | 192.168.0.1:255 |                                                                                                                  |
| assess_times            | int    | 5               | proxy evaluation times                                                                                           |
| success_times           | int    | 5               | The number of times the proxy successfully passed the evaluation                                                 |
| avg_response_time       | float  | 0.001           | -                                                                                                                |
| continuous_failed_times | int    | 0               | The number of consecutive failures during the proxy evaluation process                                           |
| score                   | float  | 25              | The higher the better                                                                                            |
| insert_time             | int    | 1540798717      | -                                                                                                                |
| update_time             | int    | 1540798717      | -                                                                                                                |

The proxy in the `proxy` table will be evaluated periodically and their scores will be modified. Low scores will be deleted.

### one `configuration file`

For convenience, the proxy in golang-proxy is stored in the portable database sqlite by default. You can make `golang-proxy` use the mysql database by adding the `config.yml` file in the executable directory.

For details, see [Config](https://github.com/storyicon/golang-proxy/tree/master/config) page.

### one `source folder`

golang-proxy needs `source` to define its crawling contents and rules. Therefore, the run directory of golang-proxy needs at least one `source` folder, and the source folder should have at least one source in `yml` format.
The source is defined as follows:

```yml
page:
    entry: "http://www.xxx.com/http/?page=1"
    template: "http://www.xxx.com/http/?page={page}"
    from: 1
    to: 2000
selector:
    iterator: ".list item"
    ip: ".ip"
    port: ".port"
category:
    parallelnumber: 3
    delayRange: [10, 30]
    interval: "@every 10m"
debug: true
```

In the definition above, `producer` will first crawl the entry page, then crawl:

```
http://www.xxx.com/http/?page=1
http://www.xxx.com/http/?page=2
http://www.xxx.com/http/?page=3
...
http://www.xxx.com/http/?page=2000
```

This source definition page expects this format:

```html
<html>
    ...
    <div class="list">
        <div class="item">
            <div class="ip">127.0.0.1</div>
            <div class="port">80</div>
            ...
        </div>
        <div class="item">
            <div class="ip">125.4.0.1</div>
            <div class="port">8080</div>
            ...
        </div>
        ...
    </div>
    ...
</html>
```

When `producer` parses a single page, it always traverses the nodes defined by iterator first, and then gets the elements defined by `ip` and `port` selectors from these nodes. The source definition above is still valid for the following HTML structure.

```html
<html>
    ...
    <div class="list">
        <div class="item">
            <div class="ip">127.0.0.1:80</div>
        </div>
        <div class="item">
            <div class="ip">125.4.0.1:8080</div>
        </div>
        ...
    </div>
    ...
</html>
```

Because when the `port` selector cannot get the content, it will try to parse the port from the text selected by the `ip` selector.

The source is stored in the source folder in yml format, and a source definition is completed. Golang-proxy will read it and crawl it the next time it starts. So you successfully define a source, store it in the source folder in YML format, and the next time you start golang-proxy, the source will enter the crawl list.

> If a source file name starts with a `.` , the source will not be read.

### four `modules`

golang-proxy consists of four modules, which cooperate to complete the task that golang-proxy wants to accomplish.

| module name | description                                                                                                                                                                                                                              |
| ----------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| producer    | Periodically fetch the source defined in the `source` directory, and write the fetched proxy to the `crude_proxy` table.                                                                                                                 |
| consumer    | Periodically read a certain number of proxies from `crude_proxy`, determine their proxy scheme type and availability, and write them to the `proxy` table.                                                                               |
| assessor    | Periodically read a number of proxies from the `proxy` table to evaluate their quality.                                                                                                                                                  |
| service     | Be responsible for the HTTP API interface provided by `golang-proxy`, allows you to filter and obtain the proxies in the `crude_proxy` and `proxy` tables by `localhost: 9999/all`, `localhost: 9999/random`, and `localhost: 9999/sql`. |

When you start the executable file of golang-proxy, you will start these module in turn. But you can add the `-mode` startup parameter after the golang-proxy executable to command golang-proxy to start only one module. Like below:

```bash
golang-proxy -mode=service
```

This will only start the HTTP API interface service.

At this point, you have mastered the 95% function of golang-proxy. If you want to find more, you can read the source code provided above, and improve them.

## Request for comments

Welcome to submit issue.
If you feel that golang-proxy is helping you, you can order a star or watch, thanks !

# 中文文档

Golang-Proxy -- 简单高效的免费代理抓取工具通过抓取网络上公开的免费代理，来维护一个属于自己的高匿代理池，用于网络爬虫、资源下载等用途。

## 在 `v3.0` 有哪些新特性

1. 依旧提供了高度灵活的 **API 接口**，在启动主程序后，即可通过在浏览器访问`localhost:9999/all` 与 `localhost:9999/random` 直接获取抓到的代理！甚至可以使用 `localhost:9999/sql?query=`来执行一些简单的 SQL 语句来自定义代理筛选规则！
2. 依旧提供 `Windows`、`Linux`、`Mac` **开箱即用版**！
   [Download Release v3.0](https://github.com/storyicon/golang-proxy/releases/)
3. 支持自动对代理类型进行判断, 可以通过 `schemeType` 判定代理对`http`和`https`的支持程度
4. 支持了 MySQL 数据库, 详情请见 [Config](https://github.com/storyicon/golang-proxy/tree/master/config)
5. 支持单独启动服务, 在启动编译好的二进制文件时, 通过 `-mode=` 来指定是否单独启动 `producer`/`consumer`/`assessor`/`service`
6. 重新设计了数据表, 请注意, 这意味着 `API` 接口发生了变动
7. 重新设计了 `源` 的数据结构, 去除了 `filter` 等字段, 请注意, 这意味着 `v2.0` 的源在直接提供给`v3.0` 使用时可能会出现一些问题
8. 更新了一些 `源`
9. 不再支持 `-source` 启动参数

## 如何使用 `golang-proxy`

### 1. 使用开箱即用版本

[Release 页面](https://github.com/storyicon/golang-proxy/releases/) 根据系统环境提供了一些压缩包，将他们解压后执行即可。

开箱即用版下载地址: [Download Release v3.0](https://github.com/storyicon/golang-proxy/releases/)

下载完成后, 将压缩包中的二进制文件和 `source` 目录解压到同一个位置, 启动二进制文件即可, 程序将会启动下面这些服务:

1. `producer` : 周期性的抓取`source`目录中定义的源, 将抓取到的代理写入到 `crude_proxy` 表中
2. `consumer` : 周期性的从 `crude_proxy` 中读取一定数量的代理, 判断它们的代理类型以及可用性, 将它们写入到 `proxy`表中
3. `assessor` : 周期性的从 `proxy` 表中读取一定数量的代理, 评估它们的质量
4. `service` : `golang-proxy` 提供的 http api 接口, 使你可以通过 `localhost:9999/all`, `localhost:9999/random`, `localhost:9999/sql?query=` 这三个接口来筛选和获取 `crude_proxy`和 `proxy` 表中的代理

当你启动编译好的二进制文件时, 默认这些服务会依次启动, 但是在 `v3.0` 版本, 你可以通过添加 `-mode` 启动参数来指定单独启动某个服务, 比如:

```
golang-proxy -mode=service
```

这样运行, 将只会启动 `service` 服务, 在启动了 `service` 之后, 你可以在浏览器中访问以下接口, 获得相应的代理:

| url                                       | description                                                            |
| ----------------------------------------- | ---------------------------------------------------------------------- |
| `localhost:9999/all`                      | 获取 `proxy` 表中所有已经抓取到的代理                                  |
| `localhost:9999/all?table=proxy`          | 获取 `proxy` 表中所有已经抓取到的代理                                  |
| `localhost:9999/all?table=crude_proxy`    | 获取 `crude_proxy` 表中所有已经抓取到的代理                            |
| `localhost:9999/random`                   | 从 `proxy` 表中随机获取一条代理                                        |
| `localhost:9999/random?table=proxy`       | 从 `proxy` 表中随机获取一条代理                                        |
| `localhost:9999/random?table=crude_proxy` | 从 `crude_proxy` 表中随机获取一条代理                                  |
| `localhost:9999/sql?query=`               | 在`query=`后加上`SQL`语句, 返回 SQL 执行结果, 只支持较为简单的查询语句 |

请注意, `crude_proxy` 只是抓取到的代理的临时储存表, 不能保证它们的质量, 而`proxy` 表中的代理将会不断得到 `assessor` 的评估, `proxy` 表中的 `score` 字段可以较为全面的反映一个代理的质量, 质量较低时会被删除

#### 接口示例: `localhost:9999/sql`

例如访问 `localhost:9999/sql?query=SELECT * FROM PROXY WHERE SCORE > 5 ORDER BY SCORE DESC`, 将会返回 `proxy` 表中所有分数大于 5 的代理, 并按照分数从高到低返回

```json
{
    "error": "",
    "message": [
        {
            "id": 2,
            "ip": "45.113.69.177",
            "port": "1080",
            // scheme_type 可以取以下值:
            // 0: 代理只支持 http
            // 1: 代理只支持 https
            // 2: 代理同时支持 http 和 https
            "scheme_type": 0,
            "content": "45.113.69.177:1080",
            // 评估次数
            "assess_times": 9,
            // 评估成功次数, 可以通过 success_times/assess_times获得代理连接成功率
            "success_times": 9,
            // 平均响应时间
            "avg_response_time": 0.098,
            // 连续失败次数
            "continuous_failed_times": 0,
            // 分数, 推荐使用 5 分以上的代理
            "score": 68.45106053570785,
            "insert_time": 1540793312,
            "update_time": 1540797880
        }
    ]
}
```

### 2. 使用源码编译

```bash
go get -u github.com/storyicon/golang-proxy
```

进入到 `golang-proxy` 目录，执行 `go build main.go`，执行生成的二进制的执行程序即可。

**注意：**

项目根目录下的 `./source` 是项目执行必须的文件夹，里面存储了各类网站源，其他的文件夹储存的均为项目源码。所以在编译后得到二进制程序 `main` 文件后，即可将 `main` 文件和 `source` 文件夹一同移动到任意地方，`main` 文件可以任意命名。

## 为什么要用 Golang-Proxy

1.  稳定、快速。  
    抓取模块，**单核并发可以到达 1000 个页面/秒**。
2.  高可配置性、高拓展性。  
    你不需要写任何代码，花**一两分钟**填写一个配置文件就可以添加一个新的网站源。
3.  评估功能。  
    通过 Assessor 评估模块，周期性测试代理质量，根据代理的**测试成功率、高匿性、测试次数、突变性、响应速度**等独立影响因子进行综合评分，算法具有高度可配置性，可以根据项目的需要可以对因子的权重进行独立调整。
4.  提供了高度灵活的 **API 接口**，在启动主程序后，即可通过在浏览器访问`localhost:9999/all` 与 `localhost:9999/random` 直接获取抓到的代理！甚至可以使用 `localhost:9999/sql?query=`来执行 SQL 语句来自定义代理筛选规则！
5.  不依赖任何服务型数据库，一键下载，开箱即用！

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

## 征求意见

1.  使用中任何问题提 `issues` 即可
2.  如果发现了新的好用的源，欢迎提交上来分享
3.  来都来了点个 Star 再走呗 : )
