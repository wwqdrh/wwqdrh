## wasmedge

安装

```bash
curl -sSf https://raw.githubusercontent.com/WasmEdge/WasmEdge/master/utils/install.sh | bash
```

安装后包含的文件

- /root/.wasmedge


## wasm示例程序

代码查看当前文件夹，一个简单的rust httpserver

```bash
cd use_dapr/05-trywasm/wasi

rustup target add wasm32-wasi

cargo build --target wasm32-wasi --release

~/.wasmedge/bin/wasmedge target/wasm32-wasi/release/wasmedge_hyper_server.wasm

curl http://localhost:8081/echo -X POST -d "WasmEdge"

# 压力测试 80个连接 持续5秒
$ go-wrk -c 80 -d 5  http://localhost:8081/
  80 goroutine(s) running concurrently
8878 requests in 5.021039808s, 1.22MB read
Requests/sec:           1768.16
Transfer/sec:           248.65KB
Avg Req Time:           45.244783ms
Fastest Request:        8.057589ms
Slowest Request:        96.78264ms
Number of Errors:       0

$ go-wrk -c 1 -d 5  http://localhost:8081/
  1 goroutine(s) running concurrently
115 requests in 5.016781719s, 16.17KB read
Requests/sec:           22.92
Transfer/sec:           3.22KB
Avg Req Time:           43.624188ms
Fastest Request:        3.84232ms
Slowest Request:        44.528103ms
Number of Errors:       0
```

### 什么是wasi

> WebAssembly 是概念机的汇编语言，而不是物理机的汇编语言。 这就是它可以在各种不同计算机体系结构上运行的原因。

WebAssembly System Interface，正是因为wasm不是物理机汇编语言，因此要让wasm运行在web之外，需要系统接口

wasi专注于安全与兼容性，例如能够控制文件的打开操作，沙箱中的两个恶意应用。左边的使用 POSIX 并且成功打开了一个它本不应访问的文件。另一个使用 WASI，而它无法打开该文件。

## basic程序

```bash
cd use_dapr/05-trywasm/basic

cargo build --release

./target/release/basic_hyper_server

# 压力测试 80个连接 持续5秒
$ go-wrk -c 80 -d 5  http://localhost:8001
  80 goroutine(s) running concurrently
178795 requests in 4.964017791s, 11.94MB read
Requests/sec:           36018.20
Transfer/sec:           2.40MB
Avg Req Time:           2.221099ms
Fastest Request:        33.185µs
Slowest Request:        39.583135ms
Number of Errors:       0

$ go-wrk -c 1 -d 5  http://localhost:8001
  1 goroutine(s) running concurrently
31930 requests in 4.905510861s, 2.13MB read
Requests/sec:           6509.01
Transfer/sec:           444.95KB
Avg Req Time:           153.633µs
Fastest Request:        61.887µs
Slowest Request:        9.759471ms
Number of Errors:       0
```

> 反而wasi模式每秒处理的请求更慢了，毕竟wasi只是单线程，但是这个性能差异是不是太大了，十多倍了
