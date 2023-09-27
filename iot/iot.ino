#include "bsec.h"
#include <ArduinoMqttClient.h>
#include <WiFiNINA.h>
#include "arduino_secrets.h"

char ssid[] = SECRET_SSID;        // your network SSID (name)
char pass[] = SECRET_PASS;    // your network password (use for WPA, or use as key for WEP)
char sensor_id[] = SENSOR_ID;
char broker[] = MQTT_HOST;

WiFiClient wifiClient;
MqttClient mqttClient(wifiClient);

int        port     = 1883;
const char topic[]  = "measurement";

//set interval for sending messages (milliseconds)
const long mqttInterval = 60000;
unsigned long previousMillis = 0;
String mqttMessage;

// Helper functions declarations
void checkIaqSensorStatus(void);
void errLeds(void);

// Create an object of the class Bsec
Bsec iaqSensor;

String output;

void conectToWiFi() {
  // attempt to connect to Wifi network:
  Serial.print("Attempting to connect to WPA SSID: ");
  Serial.println(ssid);
  while (WiFi.begin(ssid, pass) != WL_CONNECTED) {
    // failed, retry
    Serial.print(".");
    delay(5000);
  }

  Serial.println("You're connected to the network");
  Serial.println();
}

void connectToMqttBroker() {
  Serial.print("Attempting to connect to the MQTT broker: ");
  Serial.println(broker);

  while (!mqttClient.connect(broker, port)) {
    Serial.print(".");
    delay(5000);
  }
 
  Serial.println("You're connected to the MQTT broker!");
  Serial.println();
}

void setupSensors(void)
{
  pinMode(LED_BUILTIN, OUTPUT);
  iaqSensor.begin(0x77, Wire);
  output = "\nBSEC library version " + String(iaqSensor.version.major) + "." + String(iaqSensor.version.minor) + "." + String(iaqSensor.version.major_bugfix) + "." + String(iaqSensor.version.minor_bugfix);
  Serial.println(output);
  checkIaqSensorStatus();

  bsec_virtual_sensor_t sensorList[13] = {
    BSEC_OUTPUT_IAQ,
    BSEC_OUTPUT_STATIC_IAQ,
    BSEC_OUTPUT_CO2_EQUIVALENT,
    BSEC_OUTPUT_BREATH_VOC_EQUIVALENT,
    BSEC_OUTPUT_RAW_TEMPERATURE,
    BSEC_OUTPUT_RAW_PRESSURE,
    BSEC_OUTPUT_RAW_HUMIDITY,
    BSEC_OUTPUT_RAW_GAS,
    BSEC_OUTPUT_STABILIZATION_STATUS,
    BSEC_OUTPUT_RUN_IN_STATUS,
    BSEC_OUTPUT_SENSOR_HEAT_COMPENSATED_TEMPERATURE,
    BSEC_OUTPUT_SENSOR_HEAT_COMPENSATED_HUMIDITY,
    BSEC_OUTPUT_GAS_PERCENTAGE
  };

  iaqSensor.updateSubscription(sensorList, 13, BSEC_SAMPLE_RATE_LP);
  checkIaqSensorStatus();
}

void setup(void)
{
  /* Initializes the Serial communication */
  Serial.begin(115200);
  // while (!Serial) delay(10); // wait for console
  delay(5000);

  conectToWiFi();
  connectToMqttBroker();
  setupSensors();
}

// Function that is looped forever
void loop(void)
{
  if (WiFi.status() != WL_CONNECTED) {
    conectToWiFi();
  }
  if (!mqttClient.connected()) {
    connectToMqttBroker();
  }
  // call poll() regularly to allow the library to send MQTT keep alive which
  // avoids being disconnected by the broker
  mqttClient.poll();

  unsigned long currentMillis = millis();
  if (iaqSensor.run()) { // If new data is available
    digitalWrite(LED_BUILTIN, LOW);
    // printIaqSensorOutput(currentMillis);
    digitalWrite(LED_BUILTIN, HIGH);

    if (currentMillis - previousMillis >= mqttInterval) {
      // save the last time a message was sent
      previousMillis = currentMillis;
      sendMqttMessage(); 
    }
  } else {
    checkIaqSensorStatus();
  }
}

// Helper function definitions
void checkIaqSensorStatus(void)
{
  if (iaqSensor.bsecStatus != BSEC_OK) {
    if (iaqSensor.bsecStatus < BSEC_OK) {
      output = "BSEC error code : " + String(iaqSensor.bsecStatus);
      Serial.println(output);
      for (;;)
        errLeds(); /* Halt in case of failure */
    } else {
      output = "BSEC warning code : " + String(iaqSensor.bsecStatus);
      Serial.println(output);
    }
  }

  if (iaqSensor.bme68xStatus != BME68X_OK) {
    if (iaqSensor.bme68xStatus < BME68X_OK) {
      output = "BME68X error code : " + String(iaqSensor.bme68xStatus);
      Serial.println(output);
      for (;;)
        errLeds(); /* Halt in case of failure */
    } else {
      output = "BME68X warning code : " + String(iaqSensor.bme68xStatus);
      Serial.println(output);
    }
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

void sendMqttMessage() {
  mqttMessage = "";
  mqttMessage += String(sensor_id);
  mqttMessage +=" ";
  mqttMessage += String(iaqSensor.iaq);
  mqttMessage +=" ";
  mqttMessage += String(iaqSensor.co2Equivalent);
  mqttMessage +=" ";
  mqttMessage += String(iaqSensor.breathVocEquivalent);
  mqttMessage +=" ";
  mqttMessage += String(iaqSensor.pressure);
  mqttMessage +=" ";
  mqttMessage += String(iaqSensor.temperature);
  mqttMessage +=" ";
  mqttMessage += String(iaqSensor.humidity);
  
  Serial.print("Sending message to topic: ");
  Serial.println(topic);
  Serial.println(mqttMessage);

  // send message, the Print interface can be used to set the message contents
  mqttClient.beginMessage(topic);
  mqttClient.print(mqttMessage);
  mqttClient.endMessage();

  Serial.println();
}