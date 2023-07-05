FROM golang:1.20.5 as build

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /loki-rsyslog-plugin

COPY . .

RUN go build -o loki-rsyslog-plugin cmd/loki-rsyslog-plugin/loki-rsyslog-plugin.go


FROM debian:12.0-slim

WORKDIR /loki-rsyslog-plugin

COPY --from=build /loki-rsyslog-plugin/loki-rsyslog-plugin .

COPY conf/config.toml .

EXPOSE 5014/tcp 5014/udp

ENTRYPOINT [ "./loki-rsyslog-plugin" ]