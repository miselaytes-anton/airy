#include "bsec.h"
#include <ArduinoMqttClient.h>
#include <WiFiNINA.h>
#include <ArduinoBearSSL.h>
#include <ArduinoECCX08.h>
#include "arduino_secrets.h"
#include "trust.h"

// wifi name
char ssid[] = SECRET_SSID;  
// wifi password     
char pass[] = SECRET_PASS;
char sensorId[] = SENSOR_ID;
char mqttHost[] = MQTT_HOST;
char mqttUsername[] = MQTT_USERNAME;
char mqttPassword[] = MQTT_PASSWORD;
int mqttPort  = MQTT_PORT;
char mqttTopic[] = "measurement";
long mqttMessageInterval = 60000;
long lastMqttMessageSentMillis = 0;

// Used for the TCP socket connection
WiFiClient    wifiClient;            
BearSSLClient sslClient(wifiClient, TAs, TAs_NUM);
// Used for SSL/TLS connection, integrates with ECC508
MqttClient    mqttClient(sslClient);

Bsec iaqSensor;

unsigned long getTime() {
  // Get the current time from the WiFi module
  return WiFi.getTime();
}

void setup(void)
{
  /* Initializes the Serial communication */
  Serial.begin(115200);
  // while (!Serial) delay(10); // wait for console
  delay(5000);

  // Check the crypto chip module presence (needed for BearSSL)
  if (!ECCX08.begin()) {
    Serial.println("No ECCX08 present!");
    while (1);
  }

  // Set a callback to get the current time
  // (used to validate the servers ssl certificate)
  ArduinoBearSSL.onGetTime(getTime);

  conectToWiFi(ssid, pass);
  connectToMqttBroker(mqttHost, mqttPort, mqttUsername, mqttPassword);
  setupSensors();
}

// Function that is looped forever
void loop(void)
{
  ledOn();

  if (WiFi.status() != WL_CONNECTED) {
    conectToWiFi(ssid, pass);
  }
  if (!mqttClient.connected()) {
    connectToMqttBroker(mqttHost, mqttPort, mqttUsername, mqttPassword);
  }
  // call poll() regularly to allow the library to send MQTT keep alive which
  // avoids being disconnected by the broker
  mqttClient.poll();

  unsigned long currentMillis = millis();
  // If new data is available
  if (iaqSensor.run()) { 
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
