NAME = gsc

bin = $(NAME)
objs = *.go

$(bin): $(objs)
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(NAME)

clean:
	@rm -rf $(NAME)