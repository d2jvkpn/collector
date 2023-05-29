# collector
A data collection program developed using Kafka, MongoDB, and Golang.

#### 1. Dependent services
- mongodb: https://github.com/d2jvkpn/deploy/tree/dev/productions/mongodb
- kafka: https://github.com/d2jvkpn/deploy/tree/dev/productions/kafka-kraft

#### 2. Docker images
- registry.cn-shanghai.aliyuncs.com/d2jvkpn/collector:dev

#### 3. Configuration
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
```

#### 4. Run
```bash
go main.go --config=configs/local.yaml --addr=0.0.0.0:5011
```
