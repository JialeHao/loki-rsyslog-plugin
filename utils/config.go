package utils

import (
	"fmt"
	"sync"

	"github.com/BurntSushi/toml"
)

var (
	cfg *Config
    cfgOnce sync.Once
)

type Config struct {
	Server Server `toml:"server"`
	Loki   Loki   `toml:"loki"`
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

func InitConfig(cfgpath string) *Config {
    cfgOnce.Do(func() {
        if _, err := toml.DecodeFile(cfgpath, &cfg); err!= nil {
            fmt.Println(err)
        }
    })
    logger.Infof("%v load success", cfgpath)
	return cfg
}
