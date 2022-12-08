# openresty试用

一个nginx+lua开发工具，用于扩展nginx的能力

## 安装

基于docker

创建测试网络

```bash
docker network create --driver bridge temp_dev
```

运行openresty

```bash
cd use_dapr/03-gray

docker pull openresty/openresty

docker run -it -p 9000:9000 -d --name openresty-dev --network temp_dev -v `pwd`/gateway/simple.conf:/etc/nginx/conf.d/simple.conf openresty/openresty

docker exec -it openresty-dev bash

ls /usr/local/openresty

cat /usr/local/openresty/nginx/conf/nginx.conf
```

运行redis

```bash
docker pull redis:6

docker run -it -d --name redis-dev --network temp_dev redis:6

# 设置测试的key
docker exec -it redis-dev bash

$ redis-cli
   redis> set foo apache.org
   OK
   redis> set bar nginx.org
   OK
```

清理

```bash
docker stop openresty-dev && docker rm openresty-dev

docker stop redis-dev && docker rm redis-dev

docker network rm temp_dev
```

## 测试

根据redis的key来动态设置各个路由请求

依赖于下面的模块

- Redis2 Nginx Module
- Lua Nginx Module
- Lua Redis Parser Library
- Set Misc Nginx Module


```bash
docker run -it --rm --network temp_dev centos

curl --user-agent foo openresty-dev:9000

curl --user-agent bar openresty-dev:9000

curl --user-agent bar-empty openresty-dev:9000
```

查看日志信息

```bash
docker exec -it openresty-dev cat /usr/local/openresty/nginx/logs/access.log

docker exec -it openresty-dev cat /usr/local/openresty/nginx/logs/error.log

# 重定向到控制台了
docker logs openresty-dev
```

# 实现docker环境下的灰度部署

openresty作为ingress服务，通过redis获取旧版服务名字以及新版服务名字、新版的路由流量负载


打包镜像

```bash
cd use_dapr/03-gray

make build
```

运行容器

```bash
cd use_dapr/03-gray

make deploy
```

测试100与0的配比

```bash
curl localhost:9000/curapp\?name=app1\&port=8080\&weight=100

for i in {1..10..1}; do curl -XGET localhost:9000/index && sleep 1; done
```

测试50-50的配比
```bash
curl localhost:9000/curapp\?name=app1\&port=8080\&weight=90

curl localhost:9000/newapp\?name=app2\&port=8080\&weight=10

for i in {1..10..1}; do curl -XGET localhost:9000/index && sleep 1; done
# "app1""app1""app1""app1""app1""app1""app1""app2""app1""app1"

curl localhost:9000/curapp\?name=app1\&port=8080\&weight=50

curl localhost:9000/newapp\?name=app2\&port=8080\&weight=50

for i in {1..10..1}; do curl -XGET localhost:9000/index && sleep 1; done
# "app2""app2""app1""app2""app2""app2""app2""app1""app1""app2"
```

需要注意，lua的`os.time`默认只能精确到秒的单位，这样会导致相同时间下生成的随机数是一样的, 所以每次请求隔一段时间

