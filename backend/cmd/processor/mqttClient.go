package main

import (
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// MessageHandlers defines topics to which MQTT shoudl subscribe and corresponding handlers.
type MessageHandlers map[string]struct {
	Handler func(mqtt.Client, mqtt.Message)
	QOS     byte
}

// MqttClientOpts defines the options for creating a mqtt client.
type MqttClientOpts struct {
	BrokerAddress   string
	ClientID        string
	MessageHandlers MessageHandlers
}

// MakeMqttClient creates mqtt client.
func NewMqttClient(o MqttClientOpts) mqtt.Client {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(o.BrokerAddress)
	opts.SetClientID(o.ClientID)

	opts.SetOrderMatters(false)       // Allow out of order messages (use this option unless in order delivery is essential)
	opts.ConnectTimeout = time.Second // Minimal delays on connect
	opts.WriteTimeout = time.Second   // Minimal delays on writes
	opts.KeepAlive = 10               // Keepalive every 10 seconds so we quickly detect network outages
	opts.PingTimeout = time.Second    // local broker so response should be quick

	// Automate connection management (will keep trying to connect and will reconnect if network drops)
	opts.ConnectRetry = true
	opts.AutoReconnect = true

	// If using QOS2 and CleanSession = FALSE then it is possible that we will receive messages on topics that we
	// have not subscribed to here (if they were previously subscribed to they are part of the session and survive
	// disconnect/reconnect). Adding a DefaultPublishHandler lets us detect this.
	opts.DefaultPublishHandler = func(_ mqtt.Client, msg mqtt.Message) {
		fmt.Printf("Unexpected message: %s\n", msg)
	}
	opts.OnConnectionLost = func(_ mqtt.Client, err error) {
		fmt.Printf("Connection lost: %s\n", err)
	}
	opts.OnConnect = func(c mqtt.Client) {
		fmt.Println("Connection established")

		// Subscribe to the topic(s)
		for topic, handlerOpts := range o.MessageHandlers {
			t := c.Subscribe(topic, handlerOpts.QOS, handlerOpts.Handler)
			// the connection handler is called in a goroutine so blocking here would hot cause an issue. However as blocking
			// in other handlers does cause problems its best to just assume we should not block
			go func() {
				_ = t.Wait() // Can also use '<-t.Done()' in releases > 1.2.0
				if t.Error() != nil {
					fmt.Printf("Error subscribing: %s\n", t.Error())
				} else {
					fmt.Println("Subscribed to: ", topic)
				}
			}()
		}
	}
	opts.OnReconnecting = func(_ mqtt.Client, _ *mqtt.ClientOptions) {
		fmt.Println("Attempting to reconnect")
	}

	// Connect to the broker
	client := mqtt.NewClient(opts)

	// If using QOS2 and CleanSession = FALSE then messages may be transmitted to us before the subscribe completes.
	// Adding routes prior to connecting is a way of ensuring that these messages are processed
	for topic, handlerOpts := range o.MessageHandlers {
		client.AddRoute(topic, handlerOpts.Handler)
	}

	return client
}
