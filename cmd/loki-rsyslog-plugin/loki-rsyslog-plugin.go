package main

import (
	"flag"

	"github.com/JialeHao/loki-rsyslog-plugin/core"
)

func main() {
	var configFilePath string
	flag.StringVar(&configFilePath, "f", "config.toml", "configfile")
	flag.Parse()
	core.Run(configFilePath)
}
