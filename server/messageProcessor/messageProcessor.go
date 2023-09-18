package messageprocessor

import (
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	topic    = "measurement"
	qos      = 1
	clientID = "tatadata"
)

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

// StartProcessing processes messages from the mqtt broker and stores them in the database.
func StartProcessing(client mqtt.Client) {
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	fmt.Println("Connection is up")
}
