package messageprocessor

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	topic    = "measurement"
	qos      = 1
	clientID = "tatadata"
)

// MessageProcessor subscribes to messages in the mqtt broker and stores them in the database.
type MessageProcessor struct {
	Client mqtt.Client
}

// MakeMqttClient creates mqtt client.
func MakeMqttClient(mqttBrokerAddress string, handlers Handlers) mqtt.Client {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(mqttBrokerAddress)
	opts.SetClientID(clientID)

	opts.SetOrderMatters(false)       // Allow out of order messages (use this option unless in order delivery is essential)
	opts.ConnectTimeout = time.Second // Minimal delays on connect
	opts.WriteTimeout = time.Second   // Minimal delays on writes
	opts.KeepAlive = 10               // Keepalive every 10 seconds so we quickly detect network outages
	opts.PingTimeout = time.Second    // local broker so response should be quick

	// Automate connection management (will keep trying to connect and will reconnect if network drops)
	opts.ConnectRetry = true
	opts.AutoReconnect = true
	opts.DefaultPublishHandler = handlers.defaultPublishHandler
	opts.OnConnectionLost = handlers.onConnectionLost
	opts.OnConnect = handlers.onConnect
	opts.OnReconnecting = handlers.onReconnecting

	// Connect to the broker
	client := mqtt.NewClient(opts)

	// If using QOS2 and CleanSession = FALSE then messages may be transmitted to us before the subscribe completes.
	// Adding routes prior to connecting is a way of ensuring that these messages are processed
	client.AddRoute(topic, handlers.onMeasurementMessageHandler)

	return client
}

// NewMessageProcessor creates a new MessageProcessor.
func NewMessageProcessor(mqttBrokerAddress string, db *sql.DB) MessageProcessor {
	client := MakeMqttClient(mqttBrokerAddress, Handlers{Db: db})

	return MessageProcessor{
		Client: client,
	}
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
