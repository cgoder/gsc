FROM alpine:3.13
MAINTAINER gcoder <gcoder@live.com>

RUN apk add --no-cache --update ffmpeg

ENV PATH=/opt/bin:$PATH

COPY ./gsc /opt/bin/gsc

EXPOSE 8080

CMD ["gsc"]