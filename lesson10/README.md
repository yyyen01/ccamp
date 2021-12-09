# 模块十作业

#### 1. 为 HTTPServer 添加 0-2 秒的随机延时
Add one new function in main.go

    func randInt(min int, max int) int {
        rand.Seed(time.Now().UTC().UnixNano())
        return min + rand.Intn(max-min)
    }

在rootHandler 添加以下(请看main.go : func rootHandler )

    delay := randInt(0, 2000)
    time.Sleep(time.Millisecond * time.Duration(delay))

#### 2. 为 HTTPServer 项目添加延时 Metric
请看metrics.go,并在rootHandler 添加以下(请看main.go : func rootHandler )

	timer := metrics.NewTimer()
	defer timer.ObserveTotal()

在main.go 的 func main() 添加

    metrics.Register()

#### 3. 将 HTTPServer 部署至测试集群，并完成 Prometheus 配置
在main.go 的 func main() 添加metrics handler

    mux.Handle("/metrics", promhttp.Handler())

Add Promtetheus annotation to Deployment yaml file (Refer to myhttpserver-all.yaml)

    template:
        metadata:
        annotations:
            prometheus.io/scrape: "true"
            prometheus.io/port: "9090"
        labels:
            app: myhttpserver
  
#### build docker image

    docker build -t yyyen01/myhttpserver .
    docker push yyyen01/myhttpserver:latest

#### Create tls.crt and tls.key for https ingress traffic

> openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout tls.key -out tls.crt -subj "/CN=myhttpserserver.com/O=myhttpserver"

#### Deploy myhttpserver deployment,service,tlskey secret and ingress

    k create -f myhttpserver-all.yaml

#### Access the httpserver metrics

> curl -k https://myhttpserver.com:<"ingress controller https port">/metrics

#### 从 Promethus 界面中查询延时指标数据
Obtain the IP and Port from "loki-prometheus-server" service information

    k get svc loki-prometheus-server
    NAME                     TYPE        CLUSTER-IP    EXTERNAL-IP   PORT(S)   AGE
    loki-prometheus-server   ClusterIP   10.97.73.33   <none>        80/TCP    19h


Enter IP and Port into browser, search for metrics keyword starts with "httpserver_execution_latency*"

![Alt text](img/prometheus_console.png?raw=true "P1")
![Alt text](img/prome-graph.png?raw=true "P2")

#### 创建一个 Grafana Dashboard 展现延时分配情况
Import a dashboard via panel Json

    {
    "annotations": {
        "list": [
        {
            "builtIn": 1,
            "datasource": "-- Grafana --",
            "enable": true,
            "hide": true,
            "iconColor": "rgba(0, 211, 255, 1)",
            "name": "Annotations & Alerts",
            "target": {
            "limit": 100,
            "matchAny": false,
            "tags": [],
            "type": "dashboard"
            },
            "type": "dashboard"
        }
        ]
    },
    "editable": true,
    "gnetId": null,
    "graphTooltip": 0,
    "id": 6,
    "links": [],
    "panels": [
        {
        "datasource": "Prometheus",
        "fieldConfig": {
            "defaults": {
            "color": {
                "mode": "palette-classic"
            },
            "custom": {
                "axisLabel": "",
                "axisPlacement": "auto",
                "barAlignment": 0,
                "drawStyle": "line",
                "fillOpacity": 0,
                "gradientMode": "none",
                "hideFrom": {
                "legend": false,
                "tooltip": false,
                "viz": false
                },
                "lineInterpolation": "linear",
                "lineWidth": 1,
                "pointSize": 5,
                "scaleDistribution": {
                "type": "linear"
                },
                "showPoints": "auto",
                "spanNulls": false,
                "stacking": {
                "group": "A",
                "mode": "none"
                },
                "thresholdsStyle": {
                "mode": "off"
                }
            },
            "mappings": [],
            "thresholds": {
                "mode": "absolute",
                "steps": [
                {
                    "color": "green",
                    "value": null
                },
                {
                    "color": "red",
                    "value": 80
                }
                ]
            }
            },
            "overrides": []
        },
        "gridPos": {
            "h": 9,
            "w": 12,
            "x": 0,
            "y": 0
        },
        "id": 2,
        "options": {
            "legend": {
            "calcs": [],
            "displayMode": "list",
            "placement": "bottom"
            },
            "tooltip": {
            "mode": "single"
            }
        },
        "targets": [
            {
            "exemplar": true,
            "expr": "histogram_quantile(0.95, sum(rate(httpserver_execution_latency_seconds_bucket[5m])) by (le))",
            "interval": "",
            "legendFormat": "",
            "refId": "A"
            },
            {
            "exemplar": true,
            "expr": "histogram_quantile(0.90, sum(rate(httpserver_execution_latency_seconds_bucket[5m])) by (le))",
            "hide": false,
            "interval": "",
            "legendFormat": "",
            "refId": "B"
            },
            {
            "exemplar": true,
            "expr": "histogram_quantile(0.50, sum(rate(httpserver_execution_latency_seconds_bucket[5m])) by (le))",
            "hide": false,
            "interval": "",
            "legendFormat": "",
            "refId": "C"
            }
        ],
        "title": "Panel Title",
        "type": "timeseries"
        }
    ],
    "refresh": "",
    "schemaVersion": 30,
    "style": "dark",
    "tags": [],
    "templating": {
        "list": []
    },
    "time": {
        "from": "now-1m",
        "to": "now"
    },
    "timepicker": {},
    "timezone": "",
    "title": "Http Server Latency",
    "uid": "mWgwgx5nz",
    "version": 1
    }

Write a while loop to continuosly hitting the httpserver url for one hour and display the statistic in Graphana:

![Alt text](img/graphana.jpg?raw=true "P3")

