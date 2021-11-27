# ccamp

#docker build
docker build -t yyyen01/myhttpserver .
docker push yyyen01/myhttpserver:latest

#containerd pull image
sudo crictl pull yyyen01/myhttpserver:latest

配置和代码分离

k create configmap propmap --from-env-file=./app.properties 

