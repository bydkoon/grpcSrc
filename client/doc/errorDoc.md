**1. tls 관련**
-

```shell

grpc options WithReturnConnectionError 사용
client pem키 없이 호출한경우  -> server은 tls사용한 경우

did not connect: grpc: no transport security set (use grpc.WithInsecure() explicitly or set credentials)
exit status 1

```

```shell

grpc options WithReturnConnectionError 사용
client에서 host, port,tls, connection 문제 -> server은 tls사용한 경우
 
did not connect: did not connect: context deadline exceeded
exit status 1

```
```shell

client tls파일이 정확하지 않을때 

 did not connect: grpc: no transport security set (use grpc.WithInsecure() explicitly or set credentials)
exit status 1

```


```shell

ssl subjectAltName ip 세팅문제
"transport: authentication handshake failed: x509: certificate signed by unknown authority"

```

**2. connection 관련**
-

```shell

grpc options WithReturnConnectionError 미사용
client에서 port, tls, connection 문제 -> server은 tls사용한 경우

rpc errors: code = Unavailable desc = connection errors: desc = "transport: Error while dialing dial tcp :0: connectex: The requested address is not valid in its context."


grpc options WithReturnConnectionError 미사용
client host 잘못기재한상태 -> server은 tls사용한 경우 
rpc errors: code = Unavailable desc = connection errors: desc = "transport: Error while dialing dial tcp 12.4.34.4:50051: i/o timeout"
exit status 1

```


