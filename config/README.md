# Config           


- [English Document](#english-document)
- [中文文档](#中文文档)


## English Document            

You can specify `golang-proxy` to use the `MySQL` database by adding a file called `config.yml` in the compiled binary directory. The directory structure is as follows:            

![directory structure]()

The contents of the `config.yml` file can be seen in [config.yml](). If the MySQL connection fails or fails to read the `config.yml` file, the `sqlite` database will be used by default.            

If you specify to use the mysql database, then the database needs you to build, but you don't need to create the data tables needed by golang-proxy at any time, because golang-proxy will automatically create when they don't exist.               

## 中文文档         

你可以通过在`golang-proxy`编译好的二进制文件的目录下建立一个名为 `config.yml` 的文件, 来指定 `golang-proxy` 使用 `MySQL` 数据库, 目录组织看起来像这样:       

![directory structure]()

在 `config.yml` 中你可以指定 `MySQL` 的 `HOST`, `PORT` 等等, `config.yml` 的示例见: [config.yml]()

如果你决定使用 `MySQL` 作为 `golang-proxy` 的储存引擎, 你需要自己建立数据库, 并在 `config.yml` 中指定它, 但是你不需要建立数据表, 因为 `golang-proxy` 会在它们不存在的时候自动创建

当 `golang-proxy` 连接 `MySQL` 出现错误, 或者没有找到同目录下的 `config.yml` 的时候, 会默认使用 `sqlite` 数据库, 你可以下载 `sqlite studio` 等软件来读取这个便携式数据库 (事实上, 使用 `golang-proxy` 提供的 [http接口]() 已经足够了)       