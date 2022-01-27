
```
openssl req -x509 -newkey rsa:4096 -days 365 -nodes -keyout ca-key.pem -out ca-cert.pem
```

```
openssl x509 -in ca-cert.pem -noout -text
```
```
openssl req -newkey rsa:4096 -nodes -keyout server-key.pem -out server-req.pem
```

```
Country Name (2 letter code) [AU]:KR
State or Province Name (full name) [Some-State]:Seoul
Locality Name (eg, city) []:Seoul
Organization Name (eg, company) [Internet Widgits Pty Ltd]:local
Organizational Unit Name (eg, section) []:local
Common Name (e.g. server FQDN or YOUR name) []:bemily
Email Address []:bykoon@bemily.com
```

```
Please enter the following 'extra' attributes
to be sent with your certificate request
A challenge password []:bemily1004!
An optional company name []:bemily
```

```
echo subjectAltName = IP:192.168.4.69 > extfile.cnf

openssl x509 -req -in server-req.pem -days 365 -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out server-cert.pem -extfile extfile.cnf
```