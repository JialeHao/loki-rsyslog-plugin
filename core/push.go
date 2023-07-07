package core

import (
    "bytes"
    "compress/gzip"
    "encoding/json"
    "net/http"
    "strconv"
    "time"
)

func workshop(lokicfg *loki) {
    defer coreWaitGroup.Done()
    for {
        LogData := <-logMsgPool
        go push(LogData, lokicfg)
    }
}

func push(log *logMsg, lc *loki) {
    var data []byte
    var httpReq *http.Request
    var httpResp *http.Response
    var err error

    if data, err = packline(log); err != nil {
        logger.Error(err)
    }

    client := &http.Client{}

    if httpReq, err = http.NewRequest("POST", lc.Url, bytes.NewReader(data)); err != nil {
        logger.Error(err)
    }

    httpReq.Header.Set("Content-Type", "application/json")
    httpReq.Header.Set("Content-Encoding", "gzip")
    if httpResp, err = client.Do(httpReq); err != nil {
        logger.Error(err)
    }

    respCode := httpResp.StatusCode

    switch respCode / 100 {
    case 2:
        logger.Infof("log push success, status_code: %d, proto: %v, ipv4: %v log: %v", respCode, log.proto, log.ip, string(log.msg))
    case 4, 5:
        logger.Errorf("an error was encountered when log was pushed to loki, status_code=%d", respCode)
    default:
        logger.Warnf("push state unknown, status_code: %d", respCode)
    }

}

func packline(log *logMsg) ([]byte, error) {
    var err error

    labels := make(map[string]string)
    labels["source"] = "rsyslog"
    labels["ip"] = log.ip
    labels["proto"] = log.proto
    labels["ts"] = log.ts.Format(time.RFC3339)

    if err = addLevelTag(labels, &log.msg); err != nil {
        labels["level"] = "unset"
        logger.Warn(err)
    }

    tsUnixNano := strconv.Itoa(int(log.ts.UnixNano()))

    lokiLog := [2]string{tsUnixNano, string(log.msg)}

    lokiLogs := [][2]string{lokiLog}

    lf := lokiFmt{
        Streams: []lokiStreams{
            {
                Stream: labels,
                Values: lokiLogs,
            },
        },
    }

    var jsonByte []byte
    if jsonByte, err = json.Marshal(lf); err != nil {
        logger.Error(err)
        return nil, err
    }

    var buf bytes.Buffer

    gzipWriter := gzip.NewWriter(&buf)

    _, err = gzipWriter.Write(jsonByte)

    if err != nil {
        logger.Error(err)
    }

    gzipWriter.Close()

    return buf.Bytes(), err
}
