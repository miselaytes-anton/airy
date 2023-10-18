void conectToWiFi(char ssid[], char pass[]) {
  // attempt to connect to Wifi network:
  Serial.print("Attempting to connect to WPA SSID: ");
  Serial.println(ssid);
  int status = WiFi.begin(ssid, pass);
  while (status != WL_CONNECTED) {
    // failed, retry
    ledOff();
    Serial.print(status);
    delay(5000);
    status = WiFi.begin(ssid, pass);
  }

  Serial.println("You're connected to the network");
  Serial.println();
}