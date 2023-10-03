# Use Ubuntu 22.04 as the base image
FROM ubuntu:22.04

# Create directories
RUN mkdir -p /opt/dtac/{bin,plugins}
RUN mkdir -p /etc/dtac/{certs,db}

# Copy the binary into the image
COPY bin/dtac-agentd-amd64 /opt/dtac/bin/dtac-agentd
COPY bin/plugins/ /opt/dtac/plugins/

# Ensure the binary has execute permissions
RUN chmod +x /opt/dtac/bin/dtac-agentd

# Expose port 8180
EXPOSE 8180

# Set the command the container will execute
CMD ["/opt/dtac/bin/dtac-agentd"]
