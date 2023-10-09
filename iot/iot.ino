#include "bsec.h"
#include <ArduinoMqttClient.h>
#include <WiFiNINA.h>
#include "arduino_secrets.h"

// wifi name
char ssid[] = SECRET_SSID;  
// wifi password     
char pass[] = SECRET_PASS;
char sensorId[] = SENSOR_ID;
char mqttHost[] = MQTT_HOST;
int mqttPort  = 1883;
char mqttTopic[] = "measurement";
long mqttMessageInterval = 60000;
long lastMqttMessageSentMillis = 0;

WiFiClient wifiClient;
MqttClient mqttClient(wifiClient);
Bsec iaqSensor;

void setup(void)
{
  /* Initializes the Serial communication */
  Serial.begin(115200);
  // while (!Serial) delay(10); // wait for console
  delay(5000);

  conectToWiFi(ssid, pass);
  connectToMqttBroker(mqttHost, mqttPort);
  setupSensors();
}

// Function that is looped forever
void loop(void)
{
  if (WiFi.status() != WL_CONNECTED) {
    conectToWiFi(ssid, pass);
  }
  if (!mqttClient.connected()) {
    connectToMqttBroker(mqttHost, mqttPort);
  }
  // call poll() regularly to allow the library to send MQTT keep alive which
  // avoids being disconnected by the broker
  mqttClient.poll();

  unsigned long currentMillis = millis();
  if (iaqSensor.run()) { // If new data is available
    digitalWrite(LED_BUILTIN, LOW);
    digitalWrite(LED_BUILTIN, HIGH);

    if (currentMillis - lastMqttMessageSentMillis >= mqttMessageInterval) {
      // save the last time a message was sent
      lastMqttMessageSentMillis = currentMillis;
      String message = encodeMqttMessage(sensorId, iaqSensor.iaq, iaqSensor.co2Equivalent, iaqSensor.breathVocEquivalent, iaqSensor.pressure, iaqSensor.temperature, iaqSensor.humidity);
      sendMqttMessage(mqttTopic, message); 
    }
  } else {
    checkIaqSensorStatus();
  }
}

void errLeds(void)
{
  pinMode(LED_BUILTIN, OUTPUT);
  digitalWrite(LED_BUILTIN, HIGH);
  delay(100);
  digitalWrite(LED_BUILTIN, LOW);
  delay(100);
}
