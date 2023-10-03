# Use Ubuntu 22.04 as the base image
FROM ubuntu:22.04

# Copy the binary into the image
COPY bin/dtac-agentd-amd64 /tmp/dtac-agentd-amd64

# Ensure the binary has execute permissions
RUN chmod +x /tmp/dtac-agentd-amd64

# Expose port 8180
EXPOSE 8180

# Set the command the container will execute
CMD ["/tmp/dtac-agentd-amd64"]
