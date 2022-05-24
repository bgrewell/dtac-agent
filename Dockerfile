FROM debian:buster

WORKDIR /dtac-agent

COPY . /dtac-agent

RUN mkdir -p /etc/dtac-agent
RUN cp /dtac-agent/support/config/config.yaml /etc/dtac-agent/config.yaml

CMD ["/dtac-agent/bin/dtac-agentd"]