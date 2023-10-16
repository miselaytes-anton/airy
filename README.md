## Local development

### Setup .env file

```
cp .env.sample .env 
```
Now modify env to provide correct values.

### Start docker with services and the server

```
make docker-dev
make server
```

## VM setup

### Access
Applications running in docker can be accessed from exteranl host as well.
This is because ` sudo ufw status` gives us:

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

### SSL
Certificates are configured using certbot

```
sudo certbot --nginx -d amiselaytes.com -d tatadata.amiselaytes.com
```

See also [those docs](https://www.digitalocean.com/community/tutorials/how-to-secure-nginx-with-let-s-encrypt-on-ubuntu-20-04)
