package core

import (
    "sync"
    "time"

    "github.com/BurntSushi/toml"
)

var logMsgPool = make(chan *logMsg, 10240)

var coreWaitGroup sync.WaitGroup


// RFC3164规定syslog长度不超过1024 bytes
// RFC5424规定syslog长度不超过2048 octets
// 为了兼容RFC3164和RFC5424设置syslogLength为2048
const syslogLength int = 2048

var severity = [8]string{
    "Emergency",
    "Alert",
    "Critical",
    "Error",
    "Warning",
    "Notice",
    "Informational",
    "Debug",
}

type logMsg struct {
    ts    time.Time
    proto string
    ip    string
    msg   []byte
}

type lokiFmt struct {
    Streams []lokiStreams `json:"streams"`
}

type lokiStreams struct {
    Stream map[string]string `json:"stream"`
    Values [][2]string       `json:"values"`
}

type config struct {
    Server server `toml:"server"`
    Loki   loki   `toml:"loki"`
}

type server struct {
    Udp     bool   `toml:"udp"`
    UdpBind string `toml:"udp_bind"`
    Tcp     bool   `toml:"tcp"`
    TcpBind string `toml:"tcp_bind"`
}

type loki struct {
    Url string `toml:"url"`
}

func loadcfg(path string) *config {
    var cfg config
    _, err := toml.DecodeFile(path, &cfg)

    if err != nil {
        logger.Fatal(err)
    }

    logger.Infof("%v load success", path)

    return &cfg
}
