FROM golang:1-alpine AS builder

ARG APP=${APP}
ARG TZ=${TZ:-Asia/Shanghai}
ARG BUILD_Region=${BUILD_Region}
ARG GO_ldflags=${GO_ldflags}

LABEL stage=${APP}_builder

# sed -i 's/dl-cdn.alpinelinux.org/mirror.tuna.tsinghua.edu.cn/g' /etc/apk/repositories;
RUN if [ "$BUILD_Region" = "cn" ]; then \
    echo "==> BUILD_Region: $BUILD_Region"; \
    sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories; \
    go env -w GOPROXY="https://goproxy.cn,direct"; \
  fi

RUN { apk --no-cache update && apk --no-cache upgrade; } &> /dev/null

WORKDIR /opt/$APP

ADD ./main.go ./go.mod ./go.sum ./project.yaml  ./
COPY ./internal ./internal
COPY ./pkg      ./pkg
COPY ./proto    ./proto
# COPY ./vendor     ./vendor

# in alpine, date doesn't parse %:z
RUN go build -o main -ldflags="-w -s ${GO_ldflags}" main.go

####
FROM alpine

ARG APP=${APP}
ARG TZ=${TZ:-Asia/Shanghai}
ARG BUILD_Region=${BUILD_Region}

RUN if [ "$BUILD_Region" = "cn" ]; then \
    echo "==> BUILD_Region: $BUILD_Region"; \
    sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories; \
  fi

RUN { apk --no-cache update && apk --no-cache upgrade && apk --no-cache add tzdata; } &> /dev/null

RUN adduser -D -u 1000 d2jvkpn
USER d2jvkpn
WORKDIR /home/d2jvkpn/$APP

COPY --from=builder /opt/$APP/main ./main
# COPY ./migrations ./migrations

EXPOSE 5020
CMD ["./main", "-config=configs/prod.yaml", "-release", "-addr=:5020"]
