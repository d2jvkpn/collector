FROM golang:1-alpine AS builder
LABEL stage=collector_builder

ARG BuildLocal=${BuildLocal}
ARG GO_ldflags=${GO_ldflags}
ENV APP_NAME=collector
ENV TZ="Asia/Shanghai"

RUN if [ "$BuildLocal" = "true" ]; then \
    echo "==> BuildLocal: true"; \
    sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories; \
    go env -w GOPROXY="https://goproxy.cn,direct"; \
  fi

RUN { apk --no-cache update && apk --no-cache upgrade; } &> /dev/null

WORKDIR /opt/$APP_NAME

ADD ./main.go ./go.mod ./go.sum ./project.yaml  ./
COPY ./internal ./internal
COPY ./pkg      ./pkg
COPY ./proto    ./proto
# COPY ./vendor     ./vendor
# COPY ./migrations ./migrations

# in alpine, date doesn't parse %:z
RUN go build -o main -ldflags="-w -s ${GO_ldflags}" main.go

####
FROM alpine

ARG BuildLocal=${BuildLocal}
ENV APP_NAME=collector
ENV TZ="Asia/Shanghai"

RUN if [ "$BuildLocal" = "true" ]; then \
    echo "==> BuildLocal: true"; \
    sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories; \
  fi

RUN { apk --no-cache update && apk --no-cache upgrade && apk --no-cache add tzdata; } &> /dev/null

RUN adduser -D -u 1000 hello
USER hello
WORKDIR /home/hello/$APP_NAME

COPY --from=builder /opt/$APP_NAME/main ./main
# COPY ./migrations ./migrations

EXPOSE 5020
CMD ["./main", "-config=configs/prod.yaml", "-release", "-addr=:5020"]
