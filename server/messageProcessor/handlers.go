package messageprocessor

import (
	"database/sql"
	"encoding/json"
	"fmt"
	database "server/db"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Handlers struct defines the handlers for mqtt messages.
type Handlers struct {
	Db *sql.DB
}

// MeasurementMessage represents a single measurement sent by a sensor.
type MeasurementMessage struct {
	Temperature uint64
	Humidity    uint64
	CO2         uint64
	VOC         uint64
	SensorID    string
}

// handle is called when a message is received
func (h Handlers) onMeasurementMessageHandler(_ mqtt.Client, msg mqtt.Message) {
	// We extract the count and write that out first to simplify checking for missing values
	var m MeasurementMessage
	if err := json.Unmarshal(msg.Payload(), &m); err != nil {
		fmt.Printf("Message could not be parsed (%s): %s", msg.Payload(), err)
	}

	measurement := database.Measurement{
		Timestamp:   time.Now(),
		SensorID:    m.SensorID,
		Temperature: m.Temperature,
		Humidity:    m.Humidity,
		CO2:         m.CO2,
		VOC:         m.VOC,
	}

	database.InsertMeasurement(h.Db, measurement)

	fmt.Printf("Received message: %s\n", msg.Payload())
}

// If using QOS2 and CleanSession = FALSE then it is possible that we will receive messages on topics that we
// have not subscribed to here (if they were previously subscribed to they are part of the session and survive
// disconnect/reconnect). Adding a DefaultPublishHandler lets us detect this.
func (h Handlers) defaultPublishHandler(_ mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Unexpected message: %s\n", msg)
}

func (h Handlers) onConnectionLost(_ mqtt.Client, err error) {
	fmt.Printf("Connection lost: %s\n", err)
}

func (h Handlers) onReconnecting(_ mqtt.Client, _ *mqtt.ClientOptions) {
	fmt.Println("Attempting to reconnect")
}

func (h Handlers) onConnect(c mqtt.Client) {
	fmt.Println("Connection established")

	// Establish the subscription - doing this here means that it will happen every time a connection is established
	// (useful if opts.CleanSession is TRUE or the broker does not reliably store session data)
	t := c.Subscribe(topic, qos, h.onMeasurementMessageHandler)
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
