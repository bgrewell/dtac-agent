FROM debian:buster

WORKDIR /system-api

COPY . /system-api

RUN mkdir -p /etc/system-api
RUN cp /system-api/support/config/config.yaml /etc/system-api/config.yaml

CMD ["/system-api/bin/system-apid"]