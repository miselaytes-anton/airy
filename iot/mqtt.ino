void connectToMqttBroker(char host[], int port) {
  Serial.print("Attempting to connect to the MQTT host: ");
  Serial.println(host);

  while (!mqttClient.connect(host, port)) {
    Serial.print(".");
    delay(5000);
  }
 
  Serial.println("You're connected to the MQTT broker!");
  Serial.println();
}

String encodeMqttMessage (char sensorId[], float iaq, float co2Equivalent, float breathVocEquivalent, float pressure, float temperature, float humidity){
    String message = "";
    message += String(sensorId);
    message +=" ";
    message += String(iaq);
    message +=" ";
    message += String(co2Equivalent);
    message +=" ";
    message += String(breathVocEquivalent);
    message +=" ";
    message += String(pressure);
    message +=" ";
    message += String(temperature);
    message +=" ";
    message += String(humidity);

    return message;
}

void sendMqttMessage(char topic[], String message) { 
  Serial.print("Sending message to topic: ");
  Serial.println(topic);
  Serial.println(message);

  // send message, the Print interface can be used to set the message contents
  mqttClient.beginMessage(topic);
  mqttClient.print(message);
  mqttClient.endMessage();

  Serial.println();
}