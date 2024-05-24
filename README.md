# Masterservice Microservice

Master Service is a part of GUFO API Gateway. Store data of all microservices. it's host, port, version, etc. Also This microservice is coordinate cron jobs of other microservices, to avoid from double cron job in case if running few replics of the same microservice

Full API Documentation in docs/ folder

## Build Microservice

```
docker build --no-cache -t masterservice:latest -f Dockerfile .
```
or
```
docker build -t masterservice:latest -f Dockerfile .
```


## Run Microservice in Docker (in case if it in the same area with API Gateway)

```
docker run --name masterservice \
--restart=always \
-v $(pwd)/config:/var/gufo/config \
-v $(pwd)/logs:/var/gufo/log \
--network="gufo" \
-d masterservice:latest
```

Before run microservice need to add in API Gateway config next lines

```
[microservices]
[microservices.masterservice]
type = 'server'
host = 'masterservice'
port = '5300'
entrypointversion = '1.0.0'
internal = 'true'
```
