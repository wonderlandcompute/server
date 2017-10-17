
Development setup
---

Basically you need to do `dep ensure` (see [dep](https://github.com/golang/dep))

After that set `DISNEYLAND_CONFIG` to some YAML file path, with contents similar to:

```
server_cert: /path/to/server/cert.crt
server_key: /path/to/server/key.key
ca_cert: /path/to/ca/cert.crt
listen_on: :50051
db_uri: postgres://localhost/disneyland?sslmode=disable
```

After that you can launch server with `go run disneyland_server.go` command

In order to run tests, you'll need to point `DISNEYLAND_TESTS_CONFIG` env variable to some YAML file with contents like:
```
client_cert: /path/to/client/cert.crt
client_key: /path/to/client/key.key
ca_cert: /path/to/ca/cert.crt
connect_to: 127.0.0.1:50051
db_uri: postgres://localhost/disneyland?sslmode=disable
```


Certificates
---

TLS auth supoort is inspired by this article: https://bbengfort.github.io/programmer/2017/03/03/secure-grpc.html

Full certificate creation workflow looks like this

```
certstrap

# For server
certstrap init --common-name "disneyland"
certstrap request-cert -ip 127.0.0.1
certstrap sign 127.0.0.1 --CA disneyland

# For client
certstrap request-cert -o ship-shield.ship-shield.docker --cn test-user
certstrap sign test-user --CA disneyland
```
