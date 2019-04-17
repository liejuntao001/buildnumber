# A service to get an incremental number

Do a simple POST to get an incremental number everytime.

The initial purpose of this project is to provide a "build number" service. Jenkins build tasks from different jobs, or even different Jenkins instances could call to get a consistently incremental number as its build number.

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

## Buiild docker image
```
docker build -t buildnumber .
```

## Run with docker-compose
In docker-compose.yml, it mount host folder /data/buildnumber to /data and use as STORAGE_DIR

```
docker-compose up -d
```

## Test from local
```
curl -i -H "Content-Type: application/json" -X POST http://127.0.0.1:8080/e9461f1c-ef78-4162-bcb7-e83da7287614
```

## Run with traefik as reverse proxy
Copy env.sample as .env, then set proper domain name in .env

```
docker-compose -f docker-compose-traefik.yml up -d
```

## Formal use
```
curl -i -H "Content-Type: application/json" -X POST https://example.com/buildnumber/e9461f1c-ef78-4162-bcb7-e83da7287614
```

## Stress test
```
go get github.com/tsenart/vegeta
echo "POST https://example.com/buildnumber/e9461f1c-ef78-4162-bcb7-e83da7287614" | vegeta attack -duration=5s -rate=200 | tee results.bin | vegeta report
```

## Deploy to kubernetes
Modify the volumeClaimTemplates to match the cluster. The example uses rook-ceph.
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