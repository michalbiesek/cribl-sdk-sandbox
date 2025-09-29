# Use the official Python 3.11 devcontainer image as base
FROM mcr.microsoft.com/vscode/devcontainers/python:3.11

# Set environment variables for Python
ENV PYTHONUNBUFFERED=1
ENV PYTHONDONTWRITEBYTECODE=1
ENV PIP_NO_CACHE_DIR=1
ENV PIP_DISABLE_PIP_VERSION_CHECK=1

# Set working directory
WORKDIR /workspace

# Install additional system dependencies as root
RUN apt-get update && apt-get install -y \
    curl \
    git \
    build-essential \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/* \
    && apt-get clean

# The vscode user already exists in the base image, so we'll use that
# Switch to vscode user for any user-specific setup
USER vscode

# Set up Python user path for vscode user
ENV PATH="/home/vscode/.local/bin:$PATH"

# Create directories for Python packages
RUN mkdir -p /home/vscode/.local/lib/python3.11/site-packages

# Set default command
CMD ["bash"]
