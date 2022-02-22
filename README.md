# simproxy

Simple proxy to proxy an endpoint to localhost

## installation

```bash
git clone git@github.com:Cyronlee/simproxy.git
```

## usage example

```bash
# default port is 8080, default timeout is 15s
./simproxy -endpoint=http://httpbin.org

# define port and timeout
./simproxy -endpoint=http://httpbin.org -port=8888 -timeout=60
```

output:

```bash
./simproxy -endpoint=http://httpbin.org
2022/02/22 10:15:30 server start at 8080
2022/02/22 10:15:34 proxied: GET http://httpbin.org/get 200
```

## help doc

```bash
$ ./simproxy -h

Usage of ./simproxy:
  -endpoint string
        Target Endpoint (e.g: http://httpbin.org/get)
  -port int
        Local TCP port to listen on (default 8080)
  -timeout int
        Set a request timeout. Specify in seconds, defaults to 15 (default 15)
  -version
        Print version
```
