FROM debian:buster

WORKDIR /system-agent

COPY . /system-agent

RUN mkdir -p /etc/system-agent
RUN cp /system-agent/support/config/config.yaml /etc/system-agent/config.yaml

CMD ["/system-agent/bin/system-agentd"]