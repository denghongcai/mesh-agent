# Builder container
FROM golang:1-stretch AS builder

COPY . /go/src/code.aliyun.com/denghongcai/mesh-agent
WORKDIR /go
RUN cd /go/src/code.aliyun.com/denghongcai/mesh-agent && go get -v -t ./...
RUN go build -i -o /go/bin/agent /go/src/code.aliyun.com/denghongcai/mesh-agent/bin/bin.go

FROM registry.cn-hangzhou.aliyuncs.com/tianchi4-docker/tianchi4-services AS builder-1

# Runner container
FROM registry.cn-hangzhou.aliyuncs.com/tianchi4-docker/debian-jdk8

COPY --from=builder-1 /root/workspace/services/mesh-provider/target/mesh-provider-1.0-SNAPSHOT.jar /root/dists/mesh-provider.jar
COPY --from=builder-1 /root/workspace/services/mesh-consumer/target/mesh-consumer-1.0-SNAPSHOT.jar /root/dists/mesh-consumer.jar
COPY --from=builder /go/bin/agent /root/dists/agent

COPY --from=builder-1 /usr/local/bin/docker-entrypoint.sh /usr/local/bin
COPY start-agent.sh /usr/local/bin

RUN set -ex && mkdir -p /root/logs

ENTRYPOINT ["docker-entrypoint.sh"]
