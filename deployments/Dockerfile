FROM golang:1-alpine AS builder
LABEL stage=collector_builder

ARG BuildLocal=${BuildLocal}
ARG GO_ldflags=${GO_ldflags}
ENV TZ="Asia/Shanghai"

RUN if [ "$BuildLocal" = "true" ]; then \
    echo "==> BuildLocal: true"; \
    sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories; \
  fi

RUN { apk --no-cache update && apk --no-cache upgrade; } &> /dev/null

WORKDIR /opt/collector

ADD ./main.go ./go.mod ./go.sum ./project.yaml  ./
COPY ./internal ./internal
# COPY ./models ./models
COPY ./pkg     ./pkg
# COPY ./vendor  ./vendor
# in alpine, date doesn't parse %:z
RUN go env -w GOPROXY="https://goproxy.cn,direct"
RUN go build -o main -ldflags="-w -s ${GO_ldflags}" main.go

####
FROM alpine

ENV TZ="Asia/Shanghai"

RUN if [ "$BuildLocal" = "true" ]; then \
    sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories; \
  fi

RUN { apk --no-cache update && apk --no-cache upgrade && apk --no-cache add tzdata; } &> /dev/null

RUN adduser -D -u 1000 hello
USER hello
WORKDIR /home/hello/collector

COPY --from=builder /opt/collector/main ./main
COPY ./migrations ./migrations

EXPOSE 5010
CMD ["./main", "-addr=:5010", "-config=configs/prod.yaml", "-release"]
