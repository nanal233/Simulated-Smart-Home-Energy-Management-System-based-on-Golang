# 数据模型和数据库

## MySQL

采用 Docker 分发的 MySQL 镜像。启动命令为：

```bash
docker run -itd --rm --name mysql -e TZ=Asia/Shanghai -e MYSQL_ROOT_PASSWORD=123456 -v mysql_data:/var/lib/mysql -p 3306:3306 mysql
```

> 根据实际需要调整时区、密码和端口。

