# GSC
go stream core.

```bash
# flow
       RPC/Websocket           IPC
[API] <-------------> [gsc] <-------> [ffmpeg]
```

## Compile 
```bash
# build gsc
$ make
```

## Build
```bash
# build ffmpeg image
$ docker build -f dockerfile.ffmpeg -t gsf/ffmpeg:latest .
# build gsc image base from ffmpeg image
$ docker build -t gsf/gsc:latest .
```
## Run
```bash
# run zk
$ docker run -d -e TZ="Asia/Shanghai" -p 2181:2181 -v $PWD/data:/data --name zookeeper --restart always zookeeper
# run gsc
$ docker run -d -p 8080:8080 --name gsc gsf/gsc:latest

# OR
$ docker-compose up -d
```

## Demo Test
### web
> http://localhost:8080/demo/web
