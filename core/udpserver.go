package core

import (
    "net"
    "time"
)

var udpPool = make(chan bool, 2048)

func udpserver(address string) {
    defer coreWaitGroup.Done()

    var addr *net.UDPAddr
    var conn *net.UDPConn
    var err error

    if addr, err = net.ResolveUDPAddr("udp", address); err != nil {
        logger.Fatal(err)
    }

    if conn, err = net.ListenUDP("udp", addr); err != nil {
        logger.Fatal(err)
    }

    defer conn.Close()

    logger.Infof("udp server started, listen %v", addr.String())

    for {
        udpPool <- true
        go udprecv(conn)
    }
}

func udprecv(c *net.UDPConn) {
    var n int
    var err error
    var udpAddr *net.UDPAddr

    data := make([]byte, syslogLength)

    if n, udpAddr, err = c.ReadFromUDP(data); err != nil {
        logger.Fatal(err)
    }

    ts := time.Now().Local()

    lm := &logMsg{
        ts:    ts,
        proto: "udp",
        ip:    udpAddr.IP.String(),
        msg:   data[:n],
    }

    logger.Infof("recv success from %v, proto=udp, recv_ts: %v msg: %v", lm.ip, lm.ts, string(lm.msg))

    logMsgPool <- lm

    <-udpPool
}
