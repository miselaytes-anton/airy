# Local development

## Setup .env file

```
cp .env.sample .env 
```
Now modify env to provide correct values.

## Start docker with services and the server

```
make docker-dev
make server
```

# VM setup

## Access
Applications running in docker can be accessed from exteranl host as well.
This is because `sudo ufw status` gives us:

```
To                         Action      From
--                         ------      ----
22/tcp                     LIMIT       Anywhere
2375/tcp                   ALLOW       Anywhere
2376/tcp                   ALLOW       Anywhere
22/tcp (v6)                LIMIT       Anywhere (v6)
2375/tcp (v6)              ALLOW       Anywhere (v6)
2376/tcp (v6)              ALLOW       Anywhere (v6)
```

Where 2375/tcp and 2376/tcp rules are giving docker permission to route requests to any open port on the host machine. Docker usese it to expose container ports. See also [uwf](https://wiki.ubuntu.com/UncomplicatedFirewall) docs.

## SSL

### Reminder of how SSL works

- Browser connects to a web server (website) secured with SSL (https). Browser requests that the server identify itself.
- Server sends a copy of its SSL Certificate, including the server’s public key.
- Browser checks the certificate root against a list of trusted CAs and that the certificate is unexpired, unrevoked, and that its common name is valid for the website that it is connecting to. If the browser trusts the certificate, it creates, encrypts, and sends back a symmetric session key using the server’s public key.
- Server decrypts the symmetric session key using its private key and sends back an acknowledgement encrypted with the session key to start the encrypted session.
- Server and Browser now encrypt all transmitted data with the session key.

### SSL setup on the server

Certificates are configured using certbot

```
sudo certbot --nginx -d amiselaytes.com -d tatadata.amiselaytes.com
```
This command generates autrenewabale certificates stored in `/etc/letsencrypt/live/amiselaytes.com` folder.

See also [those docs](https://www.digitalocean.com/community/tutorials/how-to-secure-nginx-with-let-s-encrypt-on-ubuntu-20-04)

SSL connection is terminated in NGINX, then traffic from NGINX to MQTT in docker container is not encrypted.

![nginx ssl](./assets/nginx-mqtt-ssl.png "Nginx SSL")

The following nginx config is used:

```
stream {
  upstream mosquitto {
    server 127.0.0.1:1883;
  }

  server {
      listen 8883 ssl;
      ssl_certificate     /etc/letsencrypt/live/amiselaytes.com/cert.pem;
      ssl_certificate_key  /etc/letsencrypt/live/amiselaytes.com/privkey.pem;
      proxy_pass mosquitto;
  }
}
```

Command for testing SSL connection:
```
make test-publisher
```

Ensure correct values in .env file.

### SSL setup on IoT

- download [ROOT CA certificates chain](./assets/ca-chain.pem) from https://amiselaytes.com using browser 
- convert this file to a C header file using [brssl tool](./scripts/brssl). This tool can be downloaded using instructions [here](https://bearssl.org/#download-and-installation)
- then run 

```
brssl ta ./assets/ca-chain.pem > ./iot/trust.h
```
- include this file in the arudion code


# Air

- Static IAQ:
        The main difference between IAQ and static IAQ (sIAQ) relies in the scaling factor calculated based on the recent sensor history. The sIAQ output has been optimized for stationary applications (e.g. fixed indoor devices) whereas the IAQ output is ideal for mobile application (e.g. carry-on devices).
- bVOCeq estimate:
        The breath VOC equivalent output (bVOCeq) estimates the total VOC concentration [ppm] in the environment. It is calculated based on the sIAQ output and derived from lab tests.
- CO2eq estimate:
        Estimates a CO2-equivalent (CO2eq) concentration [ppm] in the environment. It is also calculated based on the sIAQ output and derived from VOC measurements and correlation from field studies.

Since bVOCeq and CO2eq are based on the sIAQ output, they are expected to perform optimally in stationnary applications where the main source of VOCs in the environment comes from human activity (e.g. in a bedroom).

# API

## Events

### Create event
POST /api/events

```json
{
  "startTimestamp": 1698090929,
  "eventType": "window:open",
  "locationId": "bedroom"
}
```

- `startTimestamp` required, must be unix timestamps in ms.
- `eventType` required, must be a string, can be anything
- `eventType` required, must be a string, one of `bedroom`, `livingroom`

```bash
curl -X POST -H "Content-Type: application/json" -d '{"startTimestamp": 1698090929, "eventType": "window:open", "locationId": "bedroom"}' http://localhost:8081/api/events
```

### Query events

GET /api/events?from=1698090929&to=1698090930

- `from` must be unix timestamps in ms.
- `to` must be unix timestamps in ms.
- `to` must be greater than `from`

```json
[{
 "id": "uuid",
 "startTimestamp": 1698090929,
 "endTimestamp": 1698090929,
  "eventType": "window:open",
  "locationId": "bedroom"
}]
```

### Add end timestamp to event

PATCH  /api/events/:eventId

```json
{"endTimestamp": 1698090929}
```

```bash
curl -X PATCH -H "Content-Type: application/json" -d '{"endTimestamp": 1698090929}' http://localhost:8081/api/events
```

## Measurements

### Query measurements

GET /api/measurements?resolution=86400&to=1702156335&from=1701810734

- `from` must be a unix timestamp in ms
- `to` must be a unix timestamp in ms
- `resolution` must be in ms, for example 86400 for a day, 3600 for an hour

```json
[
  {
    "timestamp": 1701734400,
    "sensorId": "bedroom",
    "iaq": 108.49368098159503,
    "co2": 949.001042944785,
    "voc": 1.4850920245398769,
    "pressure": 101128.12865030681,
    "temperature": 18.051226993865033,
    "humidity": 44.831656441717776
  },
  {
    "timestamp": 1701734400,
    "sensorId": "livingroom",
    "iaq": 98.50804878048783,
    "co2": 940.1318902439023,
    "voc": 1.377317073170732,
    "pressure": 101159.84170731703,
    "temperature": 20.659268292682928,
    "humidity": 42.94823170731707
  }
]

```
