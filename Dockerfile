FROM alpine:3.4

MAINTAINER Pavel Paulau <pavel.paulau@gmail.com>

EXPOSE 8080

COPY static static
COPY rollercoaster /usr/local/bin/rollercoaster

CMD ["rollercoaster"]
