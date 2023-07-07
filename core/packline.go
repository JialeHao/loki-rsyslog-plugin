package core

import (
    "bytes"
    "compress/gzip"
    "encoding/json"
    "errors"
    "regexp"
    "strconv"
    "time"
)

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

    // 压缩json
    var buf bytes.Buffer

    gzipWriter := gzip.NewWriter(&buf)

    if _, err = gzipWriter.Write(jsonByte); err != nil {
        logger.Error(err)
    }

    gzipWriter.Close()


    return buf.Bytes(), err
}

func addLevelTag(m map[string]string, log *[]byte) error {
    var pri int
    var err error

    expr := regexp.MustCompile(`^<(\d+)>.*`)
    matchResult := expr.FindSubmatch(*log)

    if len(matchResult) < 2 {
        err = errors.New("pri match error: can not add level tag")
        return err
    }

    if pri, err = strconv.Atoi(string(matchResult[1])); err != nil {
        return err
    }

    severityCode := pri % 8

    m["level"] = severity[severityCode]
    
    return nil
}
