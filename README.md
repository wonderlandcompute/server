
Development setup
---

Basically you need to do `dep ensure` (see [dep](https://github.com/golang/dep))

After that set `OPTIMUS_CONFIG` to some YAML file path, with contents similar to:

```
server_cert: /path/to/server/cert.crt
server_key: /path/to/server/key.key
ca_cert: /path/to/ca/cert.crt
listen_on: :50051
db_uri: postgres://localhost/optimus?sslmode=disable
```

After that you can launch server with `go run optimus_server.go` command

In order to run tests, you'll need to point `OPTIMUS_TESTS_CONFIG` env variable to file with contents like:
```
client_cert: /path/to/client/cert.crt
client_key: /path/to/client/key.key
ca_cert: /path/to/ca/cert.crt
connect_to: 127.0.0.1:50051
db_uri: postgres://localhost/optimus?sslmode=disable
```


Certificates
---

TLS auth supoort is inspired by this article: https://bbengfort.github.io/programmer/2017/03/03/secure-grpc.html

Full certificate creation workflow looks like this

```
certstrap

# For server
certstrap init --common-name "optimus"
certstrap request-cert -ip 127.0.0.1
certstrap sign 127.0.0.1 --CA optimus

# For client
certstrap request-cert -o ship-shield --cn test-user
certstrap sign test-user --CA optimus
```