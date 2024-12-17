# Use a lightweight Debian image as the base
FROM debian:bullseye-slim

# Install Bash and GCC
RUN apt-get update && apt-get install -y \
    bash \
    gcc \
    make \
    && apt-get clean && rm -rf /var/lib/apt/lists/*

# Set Bash as the default shell
SHELL ["/bin/bash", "-c"]

# Set the working directory inside the container
WORKDIR /workspace
