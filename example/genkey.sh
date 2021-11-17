#!/bin/bash
mkdir certs
rm -rf certs/*
cd certs

HOST="localhost"

### 创建CA
openssl genrsa -aes256 -out ca-key.pem 4096
openssl req -nodes -new -x509 -days 365 -key ca-key.pem -sha256 -out ca.pem

### 创建 server key
openssl genrsa -out server-key.pem 4096
openssl req -subj "/CN=$HOST" -sha256 -new -key server-key.pem -out server.csr

echo subjectAltName = DNS:$HOST,IP:127.0.0.1 > extfile-server.cnf
echo extendedKeyUsage = serverAuth >> extfile-server.cnf

### 签名 server key
openssl x509 -req -days 365 -sha256 -in server.csr -CA ca.pem -CAkey ca-key.pem \
  -CAcreateserial -out server-cert.pem -extfile extfile-server.cnf

### 创建客户端CA
openssl genrsa -out key.pem 4096
openssl req -subj '/CN=client' -new -key key.pem -out client.csr
echo extendedKeyUsage = clientAuth > extfile-client.cnf

### 签名 client key
openssl x509 -req -days 365 -sha256 -in client.csr -CA ca.pem -CAkey ca-key.pem \
  -CAcreateserial -out cert.pem -extfile extfile-client.cnf


### 修改权限
chmod -v 0400 ca-key.pem key.pem server-key.pem
chmod -v 0444 ca.pem server-cert.pem cert.pem

### 各归其位
rm -rf server
mkdir -p server
cp ca.pem server-cert.pem server-key.pem ./server

rm -rf client
mkdir -p client
cp ca.pem cert.pem key.pem ./client

### 清理工作
rm -v client.csr server.csr extfile-server.cnf extfile-client.cnf
rm -v ca.srl
rm -v *.pem