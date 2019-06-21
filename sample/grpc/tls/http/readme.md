
make certificate

```
openssl ecparam -genkey -name secp384r1 -out openserver.key
openssl req -new -x509 -sha256 -key openserver.key -out openserver.crt -days 3650
```