package main

import (
	"log"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Processor subscribes to messages in the mqtt broker and stores them in the database.
type Processor struct {
	Client   mqtt.Client
	LogError *log.Logger
	LogInfo  *log.Logger
}

// StartProcessing connects to mqtt broker and starts processing messages.
func (p Processor) StartProcessing() {

	if token := p.Client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	p.LogInfo.Println("Processing started.")
}

// StopProcessing disconnects from mqtt broker and stop processing messages.
func (p Processor) StopProcessing() {
	p.Client.Disconnect(1000)
	p.LogInfo.Println("Processing stopped.")
}

// EnableMqttLogging enables mqtt logging.
func (p Processor) EnableMqttLogging() {
	mqtt.ERROR = log.New(os.Stdout, "[ERROR] ", 0)
	mqtt.CRITICAL = log.New(os.Stdout, "[CRITICAL] ", 0)
	mqtt.WARN = log.New(os.Stdout, "[WARN]  ", 0)
	mqtt.DEBUG = log.New(os.Stdout, "[DEBUG] ", 0)
}
