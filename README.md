# go-http-erver
##Simple HTTP server written in Go Lang.  Used for testing.

##Building Docker Container:
```bash
docker build --no-cache=true --progress=plain --tag go-http-server --tag go-http-server:1.0.0 --tag go-http-server:latest .
```  

##Running Docker Container in local Docker:
```bash
docker container run --detach --name go-http-server --publish 8081:8081 --publish 8082:8082 go-http-server
```


curl  http://localhost:8082/healthz/liveness

curl  http://localhost:8082/healthz/readiness


curl  http://localhost:8082/debugz/pprof/

curl  http://localhost:8082/debugz/pprof/cmdline

curl  http://localhost:8082/debugz/pprof/profile

curl  http://localhost:8082/debugz/pprof/symbol

curl  http://localhost:8082/debugz/pprof/trace

curl  http://localhost:8082/debugz/pprof/goroutine

curl  http://localhost:8082/debugz/pprof/heap

curl  http://localhost:8082/debugz/pprof/threadcreate

curl  http://localhost:8082/debugz/pprof/block

curl  http://localhost:8082/debugz/vars
