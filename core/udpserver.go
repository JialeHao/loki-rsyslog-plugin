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
    data := make([]byte, logDataSize)
    n, udpAddr, err := c.ReadFromUDP(data)
    ts := time.Now().Local()

    if err != nil {
        logger.Error(err)
    }

    lm := &logMsg{
        ts:    ts,
        proto: "udp",
        ip:    udpAddr.IP.String(),
        msg:   data[:n],
    }

    logger.Infof("log receive success from %v, proto=udp, recv_ts: %v msg: %v", lm.ip, lm.ts, string(lm.msg))

    logMsgPool <- lm

    <-udpPool
}
