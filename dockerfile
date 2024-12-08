FROM golang:1.23.3

# Install Oracle Instant Client
# Update and install prerequisites
RUN apt-get update && apt-get install -y libaio1

# Define Oracle Instant Client version
ENV ORACLE_VERSION=21_8

# Download Oracle Instant Client
RUN curl -o instantclient-basic.zip https://download.oracle.com/otn_software/linux/instantclient/instantclient-basic-linux.x64-${ORACLE_VERSION}.0.0dbru.zip

# Unzip and set up Oracle Instant Client
RUN unzip instantclient-basic.zip -d /usr/lib/oracle && \
    rm -f instantclient-basic.zip && \
    ln -s /usr/lib/oracle/instantclient_${ORACLE_VERSION} /usr/lib/oracle/instantclient

# Configure library paths
RUN echo "/usr/lib/oracle/instantclient_${ORACLE_VERSION}" > /etc/ld.so.conf.d/oracle-instantclient.conf && \
    ldconfig


WORKDIR /app

# Copy backend code
COPY backend/go.mod backend/go.sum ./backend/
COPY backend/ ./backend/

# Copy frontend code if needed
COPY frontend/ ./frontend/

WORKDIR /app/backend

# Install dependencies
RUN go mod download

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o ./app.bin

CMD ["/app/backend/app.bin"]
