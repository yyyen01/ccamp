### deploy httpserver
```
kubectl create ns securesvc
kubectl label ns securesvc istio-injection=enabled
kubectl create -f myhttpserver-all.yaml -n securesvc
```
```
### Generate ssl certificates
openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 -subj '/O=cncamp Inc./CN=*.cncamp.io' -keyout cncamp.io.key -out cncamp.io.crt
kubectl create -n istio-system secret tls cncamp-credential --key=cncamp.io.key --cert=cncamp.io.crt
kubectl apply -f istio-spec.yaml -n securesvc
```

### check ingress ip
```
k get svc -nistio-system
NAME                   TYPE           CLUSTER-IP      EXTERNAL-IP   PORT(S)                                                                      AGE
istio-egressgateway    ClusterIP      10.109.69.5     <none>        80/TCP,443/TCP                                                               129m
istio-ingressgateway   LoadBalancer   10.97.75.71     <pending>     15021:30372/TCP,80:30230/TCP,443:31137/TCP,31400:32254/TCP,15443:30075/TCP   129m
istiod                 ClusterIP      10.99.204.228   <none>        15010/TCP,15012/TCP,443/TCP,15014/TCP                                        129m

```
### access the httpserver via ingress (with L7 url rewrite)
```

root@yen-VirtualBox:/usr/local/ccamp/homework/lesson12# curl --resolve httpsserver.cncamp.io:443:10.97.75.71 https://httpsserver.cncamp.io/ok -k
200


root@yen-VirtualBox:/usr/local/ccamp/homework/lesson12# curl --resolve httpsserver.cncamp.io:443:10.97.75.71 https://httpsserver.cncamp.io/hello -k
Congratulation! You have hit the page successfully

```
