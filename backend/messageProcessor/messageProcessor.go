package messageprocessor

import (
	"fmt"
	"log"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// MessageProcessor subscribes to messages in the mqtt broker and stores them in the database.
type MessageProcessor struct {
	Client mqtt.Client
}

// StartProcessing connects to mqtt broker and starts processing messages.
func (p MessageProcessor) StartProcessing() {

	if token := p.Client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	fmt.Println("Processing started.")
}

// StopProcessing disconnects from mqtt broker and stop processing messages.
func (p MessageProcessor) StopProcessing() {
	p.Client.Disconnect(1000)
	fmt.Println("Processing stopped.")
}

// EnableMqttLogging enables mqtt logging.
func (p MessageProcessor) EnableMqttLogging() {
	mqtt.ERROR = log.New(os.Stdout, "[ERROR] ", 0)
	mqtt.CRITICAL = log.New(os.Stdout, "[CRITICAL] ", 0)
	mqtt.WARN = log.New(os.Stdout, "[WARN]  ", 0)
	mqtt.DEBUG = log.New(os.Stdout, "[DEBUG] ", 0)
}
