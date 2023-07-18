package core

import (
	"net"
	"time"
)

func tcpserver(address string) {
    defer coreWaitGroup.Done()

    var listener net.Listener
    var err error

    if listener, err = net.Listen("tcp", address); err != nil {
        logger.Fatal(err)
    }

    defer listener.Close()

    logger.Infof("tcp server started, listen %v", address)

    for {
        var conn net.Conn
        if conn, err = listener.Accept(); err != nil {
            logger.Fatal(err)
        }
        go tcprecv(conn)
    }
}

func tcprecv(conn net.Conn) {
    defer conn.Close()

    var n int
    var err error

    buf := make([]byte, syslogLength)

    if n, err = conn.Read(buf); err != nil {
        logger.Error(err)
        return
    }

    ts := time.Now().Local()

    remoteIp, _ := net.ResolveTCPAddr("tcp", conn.RemoteAddr().String())

    lm := &LogMsg{
        ts:    ts,
        proto: "tcp",
        ip:    string(remoteIp.IP.String()),
        msg:   buf[:n],
    }

    logMsgPool <- lm

    logger.Infof("recv success from %v, proto=tcp, recv_ts: %v msg: %v", lm.ip, lm.ts, string(lm.msg))
}
