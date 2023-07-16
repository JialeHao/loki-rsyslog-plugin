package core

import (
	"sync"
	"time"

	// "github.com/BurntSushi/toml"
)

var logMsgPool = make(chan *LogMsg, 10240)

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

type LogMsg struct {
	ts    time.Time
	proto string
	ip    string
	msg   []byte
}

type lokiMsg struct {
	Streams []lokiStreams `json:"streams"`
}

type lokiStreams struct {
	Stream map[string]string `json:"stream"`
	Values [][2]string       `json:"values"`
}
