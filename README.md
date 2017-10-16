
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
certstrap request-cert -o ship-shield --cn test-user
certstrap sign test-user --CA disneyland
```

In order to use access to `ANY` project/kind jobs setup `openssl.cnf`, you could locate it with running folowing comand 
```
openssl version -a
```
You get output like this, and file `openssl.cnf` located in directory `OPENSSLDIR`
```
# macbook:disneyland macbook$ openssl version -a
OpenSSL 1.0.2l  25 May 2017
built on: reproducible build, date unspecified
platform: darwin64-x86_64-cc
options:  bn(64,64) rc4(ptr,int) des(idx,cisc,16,int) idea(int) blowfish(idx)
compiler: cc -I. -I.. -I../include  -fPIC -fno-common -DOPENSSL_PIC -DOPENSSL_THREADS -D_REENTRANT -DDSO_DLFCN -DHAVE_DLFCN_H -arch x86_64 -O3 -DL_ENDIAN -Wall -DOPENSSL_IA32_SSE2 -DOPENSSL_BN_ASM_MONT -DOPENSSL_BN_ASM_MONT5 -DOPENSSL_BN_ASM_GF2m -DSHA1_ASM -DSHA256_ASM -DSHA512_ASM -DMD5_ASM -DAES_ASM -DVPAES_ASM -DBSAES_ASM -DWHIRLPOOL_ASM -DGHASH_ASM -DECP_NISTZ256_ASM
OPENSSLDIR: "/Users/macbook/anaconda/envs/py3/ssl"
```
find and uncomment `[ new_oids ]` section and modify the `[ req_distinguished_name ]` section
```
[ new_oids ]
projectoid=1.2.3.4
kindoid=${projectoid}.5.6
...
[ req_distinguished_name ]
projectoid = access_to_project
kindoid = access_to_jobkind
```
now you could run the following script 
```
openssl

# for autorize docker-worker and project/kind access
openssl genrsa -out out/test-user.key 4096
openssl req -new -key out/test-user.key -out out/test-user.csr
openssl x509 -req -in out/test-user.csr -CA out/disneyland.crt -CAkey out/disneyland.key -CAcreateserial -out out/test-user.crt -days 5000
```