# 模块八作业

1. build docker image

> docker build -t yyyen01/myhttpserver .

> docker push yyyen01/myhttpserver:latest

2. 将配置写入app.properties.再创建一个configMap配置和代码分离.

> k create configmap propmap --from-env-file=./app.properties

3. Deploy nginx-ingress-controller

> k create -f nginx-ingress-deployment.yaml

4. Create tls.crt and tls.key for https ingress traffic

> openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout tls.key -out tls.crt -subj "/CN=myhttpserserver.com/O=myhttpserver"

5. Deploy myhttpserver deployment,service,tlskey secret and ingress

> k create -f myhttpserver-all.yaml

6. Access the httpserver using https ingress port

> curl -kv https://myhttpserver.com:<"ingress controller https port">
