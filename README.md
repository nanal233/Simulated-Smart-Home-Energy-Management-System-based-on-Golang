# 模拟智能家居能耗管理系统

## 编译

JetBrains 中，需要添加“构建”，并加入“服务端”和“客户端”。

其中，“服务端”构建参数如下：

- 运行种类：软件包
- 软件包路径：github.com/vistart/project20240227/server

其它均为默认。

“客户端”构建参数如下：

- 运行种类：软件包
- 软件包路径：github.com/vistart/project20240227/client

其它均为默认。

## 数据库

新建架构：
```sql
create database project20240227;
```

再执行 [server/models/database.sql](server/models/database.sql) 描述的数据表schema。

如有必要，需要修改 [server/conf/server1.toml](server/conf/server1.toml) 中数据库的连接信息。

## 调试

调试前需要准备配置文件，否则无法执行。

Client 和 Server 在各自文件夹下都有预知的配置文件，例如：

- Client: [client/conf](client/conf)
- Server: [Server/conf](server/conf)

启动时，需要指定 `--config` 参数为具体的配置文件路径，例如 Server 端：

```bash
go run server/main.go --config server/conf/server1.toml
```

Client 端：

```bash
go run client/main.go --config client/conf/type1device1.toml
```

如果是在 JetBrains 中调试，则需要在“程序实参”中填入该参数。