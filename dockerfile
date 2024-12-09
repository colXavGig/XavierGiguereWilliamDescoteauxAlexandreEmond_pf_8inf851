# Use a base image
FROM golang:1.23.3

# Set working directory
WORKDIR /app

# Install prerequisites
RUN apt-get update && apt-get install -y \
    libaio1 \
    unzip \
    curl \
    && rm -rf /var/lib/apt/lists/*

# Define Oracle Instant Client version
ENV ORACLE_VERSION=21_16

# Create necessary directories
RUN mkdir -p /usr/lib/oracle

# Copy the manually downloaded Oracle Instant Client ZIP file
COPY instantclient-basic-linux.x64-21.16.0.0.0dbru.zip /app/instantclient-basic.zip

# Unzip and set up Oracle Instant Client
RUN unzip /app/instantclient-basic.zip -d /usr/lib/oracle && \
    rm -f /app/instantclient-basic.zip && \
    ln -s /usr/lib/oracle/instantclient_${ORACLE_VERSION} /usr/lib/oracle/instantclient

# Configure library paths
RUN echo "/usr/lib/oracle/instantclient_${ORACLE_VERSION}" > /etc/ld.so.conf.d/oracle-instantclient.conf && \
    ldconfig

# Copy application code
COPY backend/go.mod backend/go.sum ./backend/
COPY backend/ ./backend/
COPY frontend/ ./frontend/

# Set working directory for backend
WORKDIR /app/backend

# Download dependencies
RUN go mod download

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o ./app.bin

# Command to run the application
CMD ["/app/backend/app.bin"]
