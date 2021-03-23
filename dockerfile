# FROM golang:1.16-alpine as build

# WORKDIR .
# # COPY . .

# RUN go get -d -v ./...
# RUN go install -v ./...

FROM alfg/ffmpeg:latest

# WORKDIR /home
ENV PATH=/opt/bin:$PATH

COPY ./gsc /opt/bin/gsc

EXPOSE 8080

CMD ["gsc"]