# collector
A data collection program developed using Kafka, MongoDB, and Golang.

#### 1. Techs
- programming language: golang, go1.20.4
  - https://go.dev
- logging: uber-go/zap
  - https://github.com/uber-go/zap
- metrics: prometheus
  - https://prometheus.io/
  - https://github.com/d2jvkpn/deploy/tree/dev/productions/cloud-native
- tacing: opentelemetry
  - https://opentelemetry.io/
  - https://github.com/d2jvkpn/deploy/tree/dev/productions/cloud-native
- database: mongodb
  - https://www.mongodb.com/
  - https://github.com/d2jvkpn/deploy/tree/dev/swarm/mongo_sharded-cluster
  - https://github.com/d2jvkpn/deploy/tree/dev/productions/mongodb
- message queue: kafka
  - https://kafka.apache.org/
  - https://github.com/d2jvkpn/deploy/tree/dev/swarm/kafka-kraft
- RPC: grpc
  - https://grpc.io/
- docker images
  - registry.cn-shanghai.aliyuncs.com/d2jvkpn/collector:dev
- devops:
  - docker, docker-compose
  - TODO: kubernetes
- service discovery
  - TODO: consul

#### 2. Configuration
```yaml
service_name: collector_local

log:
  path: logs/collector.log
  size_mb: 256

mongodb:
  uri: mongodb://root:root@localhost:27017
  min_pool_size: 20
  max_pool_size: 500
  timeout_secs: 5
  db: collector

kafka:
  addrs: [localhost:9092]
  version: 3.4.0
  topic: collector
  # consumer
  group_id: default
  # producer
  key: key0001

bp:
  count: 1000
  interval: 1m

metrics:
  addr: :5011
  prometheus: true
  debug: true

otel:
  enable: false
  addr: otel-collector:4317

grpc:
  tls: false
  cert: configs/server.pem
  key: configs/server.key
```

#### 3. Run
```bash
go main.go --config=configs/local.yaml --addr=0.0.0.0:5021
```
