package core

import "github.com/JialeHao/loki-rsyslog-plugin/utils"

var logger = utils.InitLogger()

func Run(configfile string) {
    cfg := loadcfg(configfile)

    if cfg.Server.Udp {
        coreWaitGroup.Add(1)
        go udpserver(cfg.Server.UdpBind)
    }

    if cfg.Server.Tcp {
        coreWaitGroup.Add(1)
        go tcpserver(cfg.Server.TcpBind)
    }

    coreWaitGroup.Add(1)
    go workshop(&cfg.Loki)

    coreWaitGroup.Wait()
}
