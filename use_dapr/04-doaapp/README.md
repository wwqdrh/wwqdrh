该服务包括4个组件，前端以及后端相对应的操作的服务

# 本地运行

go服务

> 处理加操作

```bash
cd use_dapr/04-doaapp/go

go get

go build app.go

dapr run --log-level debug --placement-host-address 192.168.0.111:50005 --app-id addapp --app-port 6000 --dapr-http-port 3503 ./app
```

node服务

> 处理除操作

```bash
cd  use_dapr/04-doaapp/node

npm install

dapr run --log-level debug --placement-host-address 192.168.0.111:50005 --app-id divideapp --app-port 4000 --dapr-http-port 3502 node app.js
```

python服务

> 处理乘操作

```bash
cd  use_dapr/04-doaapp/python

dapr run --log-level debug --placement-host-address 192.168.0.111:50005 --app-id multiplyapp --app-port 5000 --dapr-http-port 3501 python app.py
```

启动前端页面

```bash
cd  use_dapr/04-doaapp/react-calculator

npm i -g pnpm

pnpm install

pnpm run buildclient

dapr run --log-level debug --placement-host-address 192.168.0.111:50005 --app-id frontendapp --app-port 8081 --dapr-http-port 3500 node server.js
```

使用curl验证接口是否正常

```bash
curl -w "\n" -s 'http://localhost:8081/calculate/add' -H 'Content-Type: application/json' --data '{"operandOne":"56","operandTwo":"3"}'

curl -w "\n" -s 'http://localhost:3501/v1.0/invoke/multiplyapp/method/multiply' -H 'Content-Type: application/json' --data '{"operandOne":"56","operandTwo":"2"}'


curl -w "\n" -s 'http://localhost:3500/multiply' -H 'Content-Type: application/json' -H 'dapr-app-id: multiplyapp' --data '{"operandOne":"56","operandTwo":"2"}'
```

打开浏览器进行测试
