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

{"seq":5}

```

The value is persisted by a file named as {uuid} in folder STORAGE_DIR.  
The default value of STORAGE_DIR is "data". It can be override as environment variable.  
The first POST for a new uuid will get 1.  
Concurrent requests to same UUID get protected by lock.



