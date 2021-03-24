# FROM alpine:3.13
# MAINTAINER gcoder <gcoder@live.com>
# RUN apk add --no-cache --update ffmpeg

FROM gsf/ffmpeg:latest
MAINTAINER gcoder <gcoder@live.com>

COPY ./gsc /opt/bin/

CMD ["/opt/bin/gsc"]