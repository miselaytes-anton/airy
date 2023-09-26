package messageprocessor

import (
	"database/sql"
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
	SensorID    string
	IAQ         float64
	CO2         float64
	VOC         float64
	Pressure    float64
	Temperature float64
	Humidity    float64
}

// parseMeasurementMessage parses a measurement message which comes in the form of "bedroom 51.86 607.44 0.52 100853 27.25 60.22"
func parseMeasurementMessage(msg mqtt.Message) (MeasurementMessage, error) {
	var m MeasurementMessage
	if _, err := fmt.Sscanf(string(msg.Payload()), "%s %g %g %g %g %g %g", &m.SensorID, &m.IAQ, &m.CO2, &m.VOC, &m.Pressure, &m.Temperature, &m.Humidity); err != nil {
		return m, err
	}

	return m, nil
}

// handle is called when a message is received
func (h Handlers) onMeasurementMessageHandler(_ mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s\n", msg.Payload())
	m, err := parseMeasurementMessage(msg)
	if err != nil {
		fmt.Printf("Message could not be parsed (%s): %s", msg.Payload(), err)
	}
	fmt.Printf("Parsed message: %+v\n", m)

	measurement := database.Measurement{
		Timestamp:   time.Now().Unix(),
		SensorID:    m.SensorID,
		IAQ:         m.IAQ,
		CO2:         m.CO2,
		VOC:         m.VOC,
		Pressure:    m.Pressure,
		Temperature: m.Temperature,
		Humidity:    m.Humidity,
	}

	_, err = database.InsertMeasurement(h.Db, measurement)
	if err != nil {
		fmt.Printf("Measurement could not be inserted into database: %s", err)
	}

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
