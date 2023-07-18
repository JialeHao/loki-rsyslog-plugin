package utils

import (
    "fmt"
    "sync"

    "github.com/BurntSushi/toml"
)

var (
    GlobalConfig *Config
    cfgOnce      sync.Once
)

type Config struct {
    Server   Server   `toml:"server"`
    Loki     Loki     `toml:"loki"`
    Dingtalk Dingtalk `toml:"dingtalk"`
    Alarm    Alarm    `toml:"alarm"`
}

type Server struct {
    Udp     bool   `toml:"udp"`
    UdpBind string `toml:"udp_bind"`
    Tcp     bool   `toml:"tcp"`
    TcpBind string `toml:"tcp_bind"`
}

type Loki struct {
    Url string `toml:"url"`
}

type Dingtalk struct {
    Token  string `toml:"token"`
    Secret string `toml:"secret"`
}

type Alarm struct {
    AlarmLevel int `toml:"level"`
}

func InitConfig(cfgpath string) *Config {
    cfgOnce.Do(func() {
        if _, err := toml.DecodeFile(cfgpath, &GlobalConfig); err != nil {
            fmt.Println(err)
        }
    })
    
    logger.Infof("%v load success", cfgpath)

    return GlobalConfig
}
