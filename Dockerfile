# Use Ubuntu as the base image
FROM ubuntu:22.04

# Set environment variables to avoid prompts during package installation
ENV DEBIAN_FRONTEND=noninteractive

# Install basic tools and utilities
RUN apt-get update && apt-get install -y \
    build-essential \
    curl \
    git \
    wget \
    vim \
    nano \
    python3 \
    python3-pip \
    unzip \
    zip \
    sudo \
    software-properties-common && \
    apt-get clean

# Set default user to root
USER root

# Set working directory
WORKDIR /workspace
