<p align="center">
    <img src="./asset/plumber.png" alt="Plumber Logo" height="80px" width="auto" />
</p>

`plumber` is a tunnel for expose localhost http and ssh server.
<p align="center">
    <img src="./asset/plumber-proxy.png" alt="Plumber" height="450px" width="auto" />
</p>

## install
```
git clone git@github.com:lpxxn/plumber.git
cd plumber
make install
```
## quick start
`Plumber` can forward HTTP request to a specified local service, and can also forward different requests to different services through configuration. For example, forward `/api/v1/user` to `Srv1` and forward all requests of `/query/*` to `Srv2`. It also supports parameter forwarding, such as forwarding `/order/:orderID` to `Srv3`
<p align="center">
    <img src="./asset/plumber-http.png" alt="Plumber hppt" height="300px" width="auto" />
</p>

eg:    
server config:
```yaml
tcpAddr: :9870

httpProxy:
  - domain: lpxxn.com
    port: 9190
    defaultForwardTo: lpxxn # forward to client which uuid is lpxxn
    forwards: # if forwards is empty, then all requests will be forwarded to defaultForwardTo
      - path: /api/*
        forwardTo: abc # forward to abc server
      - path: /order/:orderNO
        forwardTo: http://127.0.0.1:7632  # if forwardTo is not empty, then forward to the server which name is forwardTo

```

client config:
```yaml
srvTcpAddr: 127.0.0.1:9870
http:
  remotePort: 9190 # remote port, same as server config port
  uid: lpxxn
  localSrvAddr: 127.0.0.1:7654
```