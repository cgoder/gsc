# GSC
go stream core


## Compile 
```bash
$ make
# or
$ CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build .
```

## Docker
### docker build
```bash
docker build -t gcoder/gsc:latest .
```
### docker run
```bash
docker run -d -p 8080:8080 --name gsc gsf/gsc:latest
```
