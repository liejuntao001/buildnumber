# A service to get an incremental number

Give your own [UUID](https://en.wikipedia.org/wiki/Universally_unique_identifier), do a simple POST to get an incremental number.

It's easy to [generate](https://www.uuidgenerator.net/) a random one and use the UUID as a key in below demos.

[GitHub](https://github.com/liejuntao001/buildnumber) [DockerHub](https://hub.docker.com/r/baibai/buildnumber) 

## Quick Demo
```
curl  -i  -H "Content-Type: application/json" -X POST https://buildnumber1.herokuapp.com/e9461f1c-ef78-4162-bcb7-e83da7287614

for i in {1..10}; do curl  -i  -H "Content-Type: application/json" -X POST https://buildnumber1.herokuapp.com/e9461f1c-ef78-4162-bcb7-e83da7287614; done
```

*This demo runs on a heroku free dyno. The first POST may take 10 seconds for the dyno to start up.  
As the free dyno will be released after 30 minutes idle, the demo data won't persist.*

## Initial purpose
The initial purpose of this project is to provide a Jenkins "build number" service.  
There are some sets of similar Jenkins tasks, but they are setup in different jobs, or even different Jenkins masters.  
For example, an open job for everyone and another restrited job for leaders, we could use a consistent incremental build number as part of the release binary version, if they are building same project.
It could be used in various scenarios where you need an incremental number.

## API
Request
```
POST /{uuid} HTTP/1.0

Parameters:
uuid, example e9461f1c-ef78-4162-bcb7-e83da7287614
```
Response
```
Content-Type: application/json; charset=UTF-8

{"bn":5}

```

uuid must be in [UUID](https://en.wikipedia.org/wiki/Universally_unique_identifier) format.  
The value is persisted by a file named as {uuid} in folder STORAGE_DIR.  
The default STORAGE_DIR is "data". It could be overridden by environment variable.  
The first POST for a new uuid will get 1, futher POSTs to same uuid get 2,3,4,5... 
Concurrent requests to same UUID get protected by lock.

## Run from command line
```
go get github.com/google/uuid
go get github.com/gorilla/mux
go install buildnumber
bin/buildnumber
```

## Build docker image
```
docker build -t buildnumber .
```

## Pull from docker hub
```
docker pull baibai/buildnumber
```

## Quick local deployment using docker-compose
```
docker-compose up -d
```
*In docker-compose.yml, it maps a volume to /data and use as STORAGE_DIR*

## Test from local
```
curl -i -H "Content-Type: application/json" -X POST http://127.0.0.1:8080/e9461f1c-ef78-4162-bcb7-e83da7287614
```

## Deploy with traefik as reverse proxy
Copy env.sample as .env, then set proper domain name in .env

```
docker-compose -f docker-compose-traefik.yml up -d
```

## Test the deployed server
```
curl -i -H "Content-Type: application/json" -X POST https://example.com/buildnumber/e9461f1c-ef78-4162-bcb7-e83da7287614
```

## Stress testing
```
go get github.com/tsenart/vegeta
echo "POST https://example.com/buildnumber/e9461f1c-ef78-4162-bcb7-e83da7287614" | vegeta attack -duration=5s -rate=200 | tee results.bin | vegeta report
```

## Deploy to Kubernetes
Modify the volumeClaimTemplates to match what the cluster has. This example uses rook-ceph.
```
kubectl apply -f buildnumber.yml
```
Add Ingress if necessary. The example uses traefik.
Be sure to modify the DOMAIN part.
```
kubectl apply -f buildnumber-ingress.yml
```

## Test deployment
```
curl -i -H "Content-Type: application/json" -X POST https://buildnumber.example.com/e9461f1c-ef78-4162-bcb7-e83da7287614
```