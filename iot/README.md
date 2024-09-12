# IoT

## SSL setup

- download [ROOT CA certificates chain](./certs/ca-chain.pem) from https://amiselaytes.com using browser 
- convert this file to a C header file using [brssl tool](./scripts/brssl). This tool can be downloaded using instructions [here](https://bearssl.org/#download-and-installation)
- then run to generate headers file and add it to Arduino code

```
./scripts/brssl ta ./certs/ca-chain.pem > ./src/trust.h
```

## Air measurements

- Static IAQ:
        The main difference between IAQ and static IAQ (sIAQ) relies in the scaling factor calculated based on the recent sensor history. The sIAQ output has been optimized for stationary applications (e.g. fixed indoor devices) whereas the IAQ output is ideal for mobile application (e.g. carry-on devices).
- bVOCeq estimate:
        The breath VOC equivalent output (bVOCeq) estimates the total VOC concentration [ppm] in the environment. It is calculated based on the sIAQ output and derived from lab tests.
- CO2eq estimate:
        Estimates a CO2-equivalent (CO2eq) concentration [ppm] in the environment. It is also calculated based on the sIAQ output and derived from VOC measurements and correlation from field studies.

Since bVOCeq and CO2eq are based on the sIAQ output, they are expected to perform optimally in stationary applications where the main source of VOCs in the environment comes from human activity (e.g. in a bedroom).

## Hardware used

- [Arduino nano IoT 33](https://docs.arduino.cc/hardware/nano-33-iot/)
- [Adafruit BME680 - Temperature, Humidity, Pressure and Gas Sensor](https://www.adafruit.com/product/3660)
