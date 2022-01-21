package main

import (
	server "github.com/anildogan/audit-log-to-kafka/kafka-log-processor/cmd"
	"github.com/anildogan/audit-log-to-kafka/kafka-log-processor/config"
)

func main() {
	config.Init()
	server.
		NewServer().
		InitLogger().
		InitRouter().
		Run()
}
