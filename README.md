# 简介
loki-rsyslog-plugin是连接Grafana Loki与Remote syslog的中间插件

## 工作原理
loki-rsyslog-plugin默认监听系统端口`514/tcp` `514/udp`接收远程系统日志，日志接收后，对日志进行打标签，然后推送至Loki日志服务器，用于实现系统/设备远程日志集中管理，统一查询。

## 优点
- 使用Golang开发，构建部署简单，具有高并发优点；
- 相较于logstash等工具占用系统计算资源极低。

## 应用场景
### 场景1：硬件日志集中管理
数据中心机房服务器BMC、交换机、存储阵列等设备开启远程系统日志功能，由`loki-rsyslog-plugin`接收日志消息并处理，统一推送至Grafana Loki,实现日志集中管理，统一查询，并推送钉钉告警。

# 开始运行
## 开始前配置
修改配置文件conf/config.toml(二进制包的配置文件位于解压目录下)

配置文件很简单，说明如下：
```toml
[server]
udp = true  # 是否开启udp监听
udp_bind = "0.0.0.0:514"  # udp监听地址
tcp = true  # 是否开启tcp监听
tcp_bind = "0.0.0.0:514"  # tcp监听地址

[loki]
url = "http://<loki_ip>:3100/loki/api/v1/push" # loki推送接口

[dingtalk]
token = "<your_accress_token>"  # 钉钉群自定义webhook机器人access_token
secret = "<your_secret_code>"  # 钉钉群自定义webhook机器人secret,一般是SEC开头
```

## 方式1：快速开始（Docker）
### 构建镜像
Dockerfile base image可根据实际情况进行修改
`docker build -t loki-rsyslog-plugin:1.0 .`
### 运行容器
`docker run -dit --name loki-rsyslog-plugin -p 5014:5014/tcp -p 5014:5014/udp loki-rsyslog-plugin:1.0`

## 方式2：二进制部署(Linux 64位)
仅展示`linux amd64`架构部署，其它系统自行编译
```bash
# （1）创建运行目录
mkdir /usr/local/loki-rsyslog-plugin

# （2）Github Release下载打包文件
wget https://github.com/JialeHao/loki-rsyslog-plugin/releases/download/v1.0.1/loki-rsyslog-plugin-1.0.1-linux.amd64.tar.gz

# （3）解压至运行目录
tar zxvf loki-rsyslog-plugin-1.0.1-linux.amd64.tar.gz -C /usr/local/loki-rsyslog-plugin/

# （4）Systemd管理
mv /usr/local/loki-rsyslog-plugin/loki-rsyslog-plugin.service /usr/lib/systemd/system/
systemctl daemon-reload

# （5）启动并设置开机自启
systemctl enable loki-rsyslog-plugin.service --now
```

## 方式3：源码编译
`go build -o loki-rsyslog-plugin cmd/loki-rsyslog-plugin/loki-rsyslog-plugin` 

# 日志标签

## 默认标签
默认的日志标签有5个：
- `source`: 内置固定值`rsyslog`，定义日志来源
- `ip`: 远程日志发送方的IP地址
- `ts`: loki-rsyslog-plugin接收到远程日志并开始处理的时间戳
- `proto`: 远程日志发送方使用的协议，tcp或udp
- `level`: 根据pri解析出来的日志等级

备注：若level=unset说明远程日志未遵守`syslog` `RFC3164`或者`RFC5424`标准

## 新增标签
暂不支持简单的方式，需修改代码并重新编译（后续重点改善）。

需要修改core/packline.go，在packline函数中新增自定义逻辑，根据实际情况新增自定义标签。

## 设计理念
做好中间人角色，不过分臃肿，以ip标签为根据，联动企业内部其它系统的开放能力，查询标签并添加，超时或返回错误则不添加此标签，或者结合loki对外提供查询能力。
