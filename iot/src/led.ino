void blinkLed(int blinkMs)
{
  pinMode(LED_BUILTIN, OUTPUT);
  digitalWrite(LED_BUILTIN, HIGH);
  delay(blinkMs);
  digitalWrite(LED_BUILTIN, LOW);
  delay(blinkMs);
}

void ledOn(void)
{
  digitalWrite(LED_BUILTIN, HIGH);
}

void ledOff(void)
{
  digitalWrite(LED_BUILTIN, LOW);
}
