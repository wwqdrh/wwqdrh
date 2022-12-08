# 概览

- 服务调用: /v1.0/invoke
- 状态管理: /v1.0/state
- 发布订阅: /v1.0/publish、/v1.0/subscribe
- 资源绑定: /v1.0/bindings
- Actors: /v1.0/actors
- 可观测性
- secret: /v1.0/secrets

# 服务注册与发现

dapr在注册应用的时候，都会为应用指定一个标示，就是通过这个标示来作为服务的区分的

```bash
$ cd use_dapr/02-basic

$ dapr run --app-id cart --app-port 5000 --dapr-http-port 3500 python app.py
```

注册完一个服务就可以通过endpoint来进行与该服务相关的调用

> 上面的app.py暴露了一个`/add接口`

```bash
$ curl http://localhost:3500/v1.0/invoke/cart/method/add -X POST

$ curl http://localhost:3500/v1.0/invoke/cart/method/index -X GET

# 假设dapr运行在具有命名空间的环境中，需要把命名空间也带上
$ curl http://localhost:3500/v1.0/invoke/cart.production/method/add -X POST
```

# 状态管理

也就是数据存储，可插拔的后端存储组件，通过http接口去存储或者读取数据

下面通过http接口来演示如果做状态管理的

```bash
# 启动一个Dapr sidecar，这里主体程序是空的没关系，因为这里只需要sidecar就行
$ dapr run --app-id myapp --dapr-http-port 3500

# 添加状态
$ curl -X POST -H "Content-Type: application/json" -d '[{ "key": "key1", "value": "value1"}, { "key": "key2", "value": "value2"}]' http://localhost:3500/v1.0/state/statestore

# 查询状态
$ curl http://localhost:3500/v1.0/state/statestore/key1

# 设置查询条件
$ curl -X POST -H "Content-Type: application/json" -d '{"keys":["key1", "key2"]}' http://localhost:3500/v1.0/state/statestore/bulk

# 删除状态
$ curl -X DELETE 'http://localhost:3500/v1.0/state/statestore/key1'
```

# 测试服务之间进行调用

```bash
PORT=5000 dapr run --log-level debug --app-id myapp1 --app-port 5000 --dapr-http-port 3500 -- python3 app.py

PORT=5001 dapr run --log-level debug --app-id myapp2 --app-port 5001 --dapr-http-port 3501 -- python3 app.py

curl localhost:3500/index -H 'dapr-app-id: myapp1'
# "index app: 5000"

curl localhost:3501/index -H 'dapr-app-id: myapp2'
# "index app: 5001"

curl localhost:3500/index -H 'dapr-app-id: myapp2'
# "index app: 5001"
```

但是有一个坑，他们之间是通过placement组件通信的(通过设置运行日志为debug之后发现的)，也就是我这个容器内运行的环境需要修改下placement地址，如下

```bash
PORT=5000 dapr run --log-level debug --placement-host-address 192.168.0.111:50005 --app-id myapp1 --app-port 5000 --dapr-http-port 3500 -- python3 app.py

PORT=5001 dapr run --log-level debug --placement-host-address 192.168.0.111:50005 --app-id myapp2 --app-port 5001 --dapr-http-port 3501 -- python3 app.py

curl localhost:3500/index -H 'dapr-app-id: myapp1'
# "index app: 5000"

curl localhost:3501/index -H 'dapr-app-id: myapp2'
# "index app: 5001"

curl localhost:3500/index -H 'dapr-app-id: myapp2'
# "index app: 5001"
```

试验成功

# pub & sub

> pub-sub目录下，一个前端publish服务，两个后端subscriber服务

启动node以及python的订阅服务

```bash
# node服务
cd use_dapr/02-basic/pub-sub/node-subscriber

pnpm i

dapr run --placement-host-address 192.168.0.111:50005 --app-id node-subscriber --app-port 3000 node app.js

# python服务
cd use_dapr/02-basic/pub-sub/python-subscriber

pip install -r requirements.txt
cat requirements.txt | xargs poetry add

dapr run --placement-host-address 192.168.0.111:50005 --app-id python-subscriber --app-port 5001 python3 app.py
```

启动react，发布服务

```bash
cd use_dapr/02-basic/pub-sub/react-form

pnpm build

dapr run --placement-host-address 192.168.0.111:50005 --app-id react-form --app-port 8081 node server.js
```

打开浏览器并进行消息发布的测试，这里node订阅了A、B类型，python订阅了A、C类型，当出现这类消息发布时，对应的路由接口会得到通知

也可以用curl进行测试

```bash
curl -s http://localhost:8081/publish -H Content-Type:application/json --data @message_a.json

curl -s http://localhost:8081/publish -H Content-Type:application/json --data @message_b.json

curl -s http://localhost:8081/publish -H Content-Type:application/json --data @message_c.json
```

# binding

> 相当于一个管道，绑定输入与输出端，用于触发事件调用，这里的实验使用kafka组件，其实还是可以用redis，具体列表查看https://docs.dapr.io/zh-hans/reference/components-reference/supported-bindings/

启动kafaka服务

```bash
cd use_dapr/02-basic/bindings

docker stack deploy -c single-kafka.yml temp
```


# secret管理

> 查看secretstore文件夹

```bash
cd use_dapr/02-basic/secretstore/node

pnpm i

# 注意node-fetchv3只支持esm模式，需要安装2.6.6 node-fetch@^2.6.6
```

在components中定义了localsecret的描述文件，描述中声明的是secrets.json

```json
{
    "mysecret": "abcd"
}
```

除了使用文件之外还可以使用环境变量

```bash
export SECRET_STORE="localsecretstore"
```

启动服务

```bash
SECRET_STORE=localsecretstore dapr run --app-id nodeapp --components-path ./components --app-port 3000 --dapr-http-port 3500 node app.js
```

查看secret

```bash
$ curl -k http://localhost:3000/getsecret

$ curl -k http://localhost:3500/v1.0/secrets/localsecretstore/mysecret

{"mysecret":"abcd"}
```