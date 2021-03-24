# GSC
go stream core.

```
       RPC/Websocket           IPC
[API] <-------------> [gsc] <-------> [ffmpeg]
```

## Compile 
```bash
$ make
```

## Build
```bash
docker build -t gsf/gsc:latest .
```
## Run
```bash
docker run -d -p 8080:8080 --name gsc gsf/gsc:latest
```
