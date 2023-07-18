package main

import "flag"

func main() {
    var configFilePath string
    flag.StringVar(&configFilePath, "f", "config.toml", "configfile")
    flag.Parse()
}